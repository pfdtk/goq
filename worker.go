package goq

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq/common/cst"
	"github.com/pfdtk/goq/iface"
	"github.com/pfdtk/goq/internal/common"
	rdq "github.com/pfdtk/goq/internal/queue/redis"
	sqsq "github.com/pfdtk/goq/internal/queue/sqs"
	"github.com/pfdtk/goq/internal/utils"
	"github.com/redis/go-redis/v9"
	"runtime/debug"
	"sync"
	"time"
)

var (
	ErrEmptyJob = errors.New("no jobs are ready for processing")
)

type Worker struct {
	conn       *sync.Map
	tasks      *sync.Map
	wg         *sync.WaitGroup
	stopRun    chan struct{}
	maxWorker  chan struct{}
	jobChannel chan *common.Job
	ctx        context.Context
	logger     iface.Logger
}

func (w *Worker) StartConsuming() error {
	w.consume()
	w.work()
	return nil
}

func (w *Worker) StopConsuming() {
	close(w.stopRun)
}

// consume only pop common from queue and send to work channel
func (w *Worker) consume() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.stopRun:
				w.logger.Info("receive stop consume sign, stopping...")
				close(w.jobChannel)
				return
			case w.maxWorker <- struct{}{}:
				// select is range, so check before run
				select {
				case <-w.stopRun:
					close(w.jobChannel)
					return
				default:
				}
				job, err := w.getNextJob()
				switch {
				case errors.Is(err, ErrEmptyJob):
					w.logger.Infof("all queue are empty")
					// sleep 1 second when all queue are empty
					time.Sleep(time.Second)
					// release token
					<-w.maxWorker
					continue
				case err != nil:
					w.logger.Error(err)
					// release token
					<-w.maxWorker
				}
				w.logger.Infof("get nex job, id=%s", job.Id)
				w.jobChannel <- job
			}
		}
	}()
}

// work process tasks
func (w *Worker) work() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			// TODO when to exit, maybe still some msg on job channel, we need ack!
			case <-w.stopRun:
				w.logger.Info("receive stop work sign, stopping...")
				return
			case job, ok := <-w.jobChannel:
				// stop working when channel was closed
				if !ok {
					w.logger.Info("job channel has been close")
					return
				}
				w.logger.Infof("start to run job, id=%s", job.Id)
				go w.runTask(job)
			}
		}
	}()
}

func (w *Worker) runTask(job *common.Job) {
	go func() {
		// goroutine timeout control
		ctx, cancel := context.WithDeadline(w.ctx, job.TimeoutAt())
		defer func() {
			w.logger.Info("task defer, release token")
			<-w.maxWorker
			cancel()
		}()
		prh := make(chan any, 1)
		go func() {
			// todo handle err
			res, _ := w.perform(job)
			w.logger.Infof("job has been processed, id=%s", job.Id)
			prh <- res
		}()
		// wait for response
		select {
		case <-prh:
			w.logger.Info("task result has been received")
			return
		case <-ctx.Done():
			// TODO retry if necessary
			// please note that, the task goroutine will not stop, even if deadline is reach
			w.logger.Warnf("tasks: %s has been reach it`s deadline", job.Id)
			return
		}
	}()
}

func (w *Worker) perform(job *common.Job) (res any, err error) {
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

func (w *Worker) getNextJob() (*common.Job, error) {
	tasks := utils.SortTask(w.tasks)
	for _, t := range tasks {
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

func (w *Worker) getJob(q iface.Queue, queueName string) (*common.Job, error) {
	message, err := q.Pop(w.ctx, queueName)
	if err == nil {
		return &common.Job{
			Id:      message.ID,
			Name:    message.Type,
			Queue:   message.Queue,
			Payload: message.Payload,
			Timeout: message.Timeout,
		}, nil
	}
	return nil, err
}
