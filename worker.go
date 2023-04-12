package goq

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq/common/cst"
	"github.com/pfdtk/goq/common/job"
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
	EmptyJobError       = errors.New("no jobs are ready for processing")
	JobExecTimeoutError = errors.New("job exec timeout")
)

type worker struct {
	conn            *sync.Map
	tasks           *sync.Map
	sortTasks       []iface.Task
	wg              *sync.WaitGroup
	stopRun         chan struct{}
	maxWorker       chan struct{}
	jobChannel      chan *job.Job
	ctx             context.Context
	logger          iface.Logger
	taskErrorHandle []iface.ErrorJobHandler
	errorHandle     []iface.ErrorHandler
}

func newWorker(ctx context.Context, s *Server) *worker {
	return &worker{
		wg:         &s.wg,
		tasks:      &s.tasks,
		maxWorker:  make(chan struct{}, s.maxWorker),
		stopRun:    make(chan struct{}),
		jobChannel: make(chan *job.Job, s.maxWorker),
		ctx:        ctx,
		logger:     s.logger,
		conn:       &s.conn,
	}
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
	w.logger.Info("worker stopped")
}

func (w *worker) pop() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.stopRun:
				w.logger.Debug("received stop sign")
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
		return
	default:
	}
	j, err := w.getNextJob()
	switch {
	case errors.Is(err, EmptyJobError):
		w.handleError(err)
		time.Sleep(time.Second)
		<-w.maxWorker
		return
	case err != nil:
		w.handleError(err)
		<-w.maxWorker
		return
	}
	w.logger.Debugf("job received, id=%s, name=%s", j.Id(), j.Name())
	w.jobChannel <- j
}

func (w *worker) work() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.stopRun:
				return
			case j, ok := <-w.jobChannel:
				// stop working when channel was closed
				if !ok {
					return
				}
				w.logger.Debugf("job processing, id=%s, name=%s", j.Id(), j.Name())
				go w.runTask(j)
			}
		}
	}()
}

func (w *worker) runTask(job *job.Job) {
	go func() {
		ctx, cancel := context.WithDeadline(w.ctx, job.TimeoutAt())
		defer func() {
			<-w.maxWorker
			cancel()
		}()
		prh := make(chan any, 1)
		go func() {
			res, err := w.perform(job)
			if err == nil {
				prh <- res
			}
		}()
		select {
		case <-prh:
			// todo write response to backend
			return
		case <-ctx.Done():
			w.logger.Debugf("job timeout, id=%s, name=%s", job.Id(), job.Name())
			w.handleError(JobExecTimeoutError)
			return
		}
	}()
}

func (w *worker) perform(job *job.Job) (res any, err error) {
	// recover err from tasks, so that program will not exit
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(string(debug.Stack()))
		}
	}()
	name := job.Name()
	v, ok := w.tasks.Load(name)
	if ok {
		task := v.(iface.Task)
		res, err = task.Run(w.ctx, job)
		if err != nil {
			w.handleJobPerformError(task, job, err)
		} else {
			w.logger.Debugf("job processed, id=%s, name=%s", job.Id(), job.Name())
			w.handleError(err)
		}
	}
	return
}

func (w *worker) getQueue(t iface.Task) iface.Queue {
	c, ok := w.conn.Load(t.OnConnect())
	if !ok {
		return nil
	}
	switch t.QueueType() {
	case cst.Redis:
		return rdq.NewRedisQueue(c.(*redis.Client))
	case cst.Sqs:
		return sqsq.NewSqsQueue(c.(*sqs.Client))
	}
	return nil
}

func (w *worker) getNextJob() (*job.Job, error) {
	for _, t := range w.sortTasks {
		if t.GetStatus() == cst.Disable || !t.CanRun() {
			continue
		}
		q := w.getQueue(t)
		if q == nil {
			continue
		}
		j, err := w.getJob(q, t.OnQueue())
		if err == nil {
			return j, nil
		}
	}
	return nil, EmptyJobError
}

func (w *worker) getJob(q iface.Queue, qn string) (*job.Job, error) {
	msg, err := q.Pop(w.ctx, qn)
	if err == nil {
		return job.NewJob(q, msg), nil
	}
	return nil, err
}

func (w *worker) handleJobPerformError(task iface.Task, job *job.Job, _ error) {
	if !job.IsReachMacAttempts() {
		_ = w.retry(task, job)
	}
	if len(w.taskErrorHandle) != 0 {
		for _, h := range w.taskErrorHandle {
			h.Handle(w.ctx, task)
		}
	}
}

func (w *worker) handleError(err error) {
	if len(w.errorHandle) == 0 {
		return
	}
	for _, h := range w.errorHandle {
		h.Handle(w.ctx, err)
	}
}

func (w *worker) retry(task iface.Task, job *job.Job) (err error) {
	backoff := task.Backoff()
	return job.Release(w.ctx, backoff)
}
