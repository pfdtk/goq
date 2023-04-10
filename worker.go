package goq

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq/common"
	"github.com/pfdtk/goq/common/cst"
	"github.com/pfdtk/goq/iface"
	rdq "github.com/pfdtk/goq/internal/queue/redis"
	sqsq "github.com/pfdtk/goq/internal/queue/sqs"
	"github.com/pfdtk/goq/internal/utils"
	"github.com/redis/go-redis/v9"
	"runtime/debug"
	"sync"
	"time"
)

var (
	ErrEmptyJob    = errors.New("no jobs are ready for processing")
	ErrWorkExpired = errors.New("task process expired")
)

type worker struct {
	conn       *sync.Map
	tasks      *sync.Map
	sortTasks  []iface.Task
	wg         *sync.WaitGroup
	stopRun    chan struct{}
	maxWorker  chan struct{}
	jobChannel chan *common.Job
	ctx        context.Context
	logger     iface.Logger
}

func (w *worker) startConsuming() error {
	w.sortTasks = utils.SortTask(w.tasks)
	w.pop()
	w.work()
	return nil
}

func (w *worker) stopConsuming() {
	w.logger.Info("stopping worker...")
	close(w.stopRun)
	// block until all workers have released the token
	for i := 0; i < cap(w.maxWorker); i++ {
		w.maxWorker <- struct{}{}
	}
}

func (w *worker) pop() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.stopRun:
				w.logger.Info("got stop sign from pop")
				return
			case w.maxWorker <- struct{}{}:
				w.readyToWork()
			}
		}
	}()
}

func (w *worker) readyToWork() {
	// select is range, so check before run
	select {
	case <-w.stopRun:
		w.logger.Info("got stop sign from work")
		return
	default:
	}
	job, err := w.getNextJob()
	switch {
	case errors.Is(err, ErrEmptyJob):
		w.logger.Info("no jobs are ready for processing on all queue")
		// sleep 1 second when all queue are empty
		time.Sleep(time.Second)
		// release token
		<-w.maxWorker
		return
	case err != nil:
		w.logger.Error(err)
		// release token
		<-w.maxWorker
		return
	}
	w.logger.Infof("got next job to process, id=%s, name=%s", job.Id, job.Name)
	w.jobChannel <- job
}

func (w *worker) work() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.stopRun:
				w.handleWorkerStopRun()
				w.logger.Info("worker has been stopped")
				return
			case job, ok := <-w.jobChannel:
				// stop working when channel was closed
				if !ok {
					w.logger.Info("worker has been stopped")
					return
				}
				w.logger.Infof("start to process job, id=%s, name=%s", job.Id, job.Name)
				go w.runTask(job)
			}
		}
	}()
}

func (w *worker) runTask(job *common.Job) {
	go func() {
		// goroutine timeout control
		ctx, cancel := context.WithDeadline(w.ctx, job.TimeoutAt())
		defer func() {
			w.logger.Infof("task defer, release token, id=%s, name=%s", job.Id, job.Name)
			<-w.maxWorker
			cancel()
		}()
		prh := make(chan any, 1)
		go func() {
			res, err := w.perform(job)
			if err != nil {
				w.handleJobError(job, err)
			}
			w.logger.Infof("job has been processed, id=%s, name=%s", job.Id, job.Name)
			prh <- res
		}()
		// wait for response
		select {
		case <-prh:
			// todo write to backend
			return
		case <-ctx.Done():
			w.logger.Warnf("tasks has been reach it`s deadline, id=%s, name=%s", job.Id, job.Name)
			w.handleJobError(job, ErrWorkExpired)
			return
		}
	}()
}

func (w *worker) perform(job *common.Job) (res any, err error) {
	// recover err from tasks, so that program will not exit
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(string(debug.Stack()))
		}
	}()
	name := job.Name
	v, ok := w.tasks.Load(name)
	if ok {
		task := v.(iface.Task)
		res, err = task.Run(w.ctx, job)
	}
	return
}

func (w *worker) getNextJob() (*common.Job, error) {
	for _, t := range w.sortTasks {
		if t.GetStatus() == cst.Disable || !t.CanRun() {
			continue
		}
		c, ok := w.conn.Load(t.OnConnect())
		if !ok {
			continue
		}
		switch t.QueueType() {
		case cst.Redis:
			q := rdq.NewRedisQueue(c.(*redis.Client))
			job, err := w.getJob(q, t.OnQueue())
			if err == nil {
				return job, nil
			}
		case cst.Sqs:
			q := sqsq.NewSqsQueue(c.(*sqs.Client))
			job, err := w.getJob(q, t.OnQueue())
			if err == nil {
				return job, nil
			}
		}
	}
	return nil, ErrEmptyJob
}

func (w *worker) getJob(q iface.Queue, qn string) (*common.Job, error) {
	msg, err := q.Pop(w.ctx, qn)
	if err == nil {
		return &common.Job{
			Id:       msg.ID,
			Name:     msg.Type,
			Queue:    msg.Queue,
			Payload:  msg.Payload,
			Timeout:  msg.Timeout,
			Attempts: msg.Attempts,
		}, nil
	}
	return nil, err
}

func (w *worker) handleWorkerStopRun() {

}

func (w *worker) handleJobError(_ *common.Job, _ error) {

}
