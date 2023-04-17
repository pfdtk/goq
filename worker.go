package goq

import (
	"context"
	"errors"
	evt "github.com/pfdtk/goq/event"
	"github.com/pfdtk/goq/internal/event"
	qm "github.com/pfdtk/goq/internal/queue"
	"github.com/pfdtk/goq/internal/utils"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/queue"
	"github.com/pfdtk/goq/task"
	"runtime/debug"
	"sync"
	"time"
)

var (
	EmptyJobError       = errors.New("no jobs are ready for processing")
	JobExecTimeoutError = errors.New("job exec timeout")
)

type worker struct {
	tasks      *sync.Map
	sortTasks  []task.Task
	wg         *sync.WaitGroup
	stopRun    chan struct{}
	maxWorker  chan struct{}
	jobChannel chan *task.Job
	ctx        context.Context
	logger     logger.Logger
	em         *event.Manager
}

func newWorker(ctx context.Context, s *Server) *worker {
	w := &worker{
		wg:         &s.wg,
		tasks:      &s.tasks,
		maxWorker:  make(chan struct{}, s.maxWorker),
		stopRun:    make(chan struct{}),
		jobChannel: make(chan *task.Job, s.maxWorker),
		ctx:        ctx,
		logger:     s.logger,
		em:         s.eventManager,
	}
	w.registerEvents()
	return w
}

func (w *worker) registerEvents() {
	w.em.Listen(&evt.WorkErrorEvent{}, evt.NewWorkerErrorHandler())
}

func (w *worker) startConsuming() error {
	w.sortTasks = utils.SortTask(w.tasks)
	w.startPop()
	w.startWork()
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

func (w *worker) startPop() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.stopRun:
				w.logger.Debug("received stop sign")
				return
			case w.maxWorker <- struct{}{}:
				go w.pop()
			}
		}
	}()
}

func (w *worker) pop() {
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
	w.logger.Infof("job received, id=%s, name=%s", j.Id(), j.Name())
	w.jobChannel <- j
}

func (w *worker) startWork() {
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
				w.logger.Infof("job processing, id=%s, name=%s", j.Id(), j.Name())
				go w.runJob(j)
			}
		}
	}()
}

func (w *worker) runJob(job *task.Job) {
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
		} else {
			prh <- err
		}
	}()
	select {
	case <-prh:
		// todo write response to backend, prh may be an error
		return
	case <-ctx.Done():
		w.logger.Infof("job timeout, id=%s, name=%s", job.Id(), job.Name())
		w.handleError(JobExecTimeoutError)
		return
	}
}

func (w *worker) perform(job *task.Job) (res any, err error) {
	var t task.Task = nil
	// recover err from tasks, so that program will not exit
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(string(debug.Stack()))
			if t != nil {
				w.handleJobError(t, job, err)
			}
		}
	}()
	name := job.Name()
	v, ok := w.tasks.Load(name)
	if ok {
		t = v.(task.Task)
		w.em.Dispatch(evt.NewJobBeforeRunEvent(t, job))
		res, err = t.Run(w.ctx, job)
		w.em.Dispatch(evt.NewJobAfterRunEvent(t, job))
		if err != nil {
			w.handleJobError(t, job, err)
		} else {
			w.logger.Infof("job processed, id=%s, name=%s", job.Id(), job.Name())
			_ = w.handleJobDone(t, job)
		}
	}
	return
}

func (w *worker) getNextJob() (*task.Job, error) {
	for _, t := range w.sortTasks {
		if t.Status() == task.Disable || !t.CanRun() {
			continue
		}
		q := qm.GetQueue(t.OnConnect(), t.QueueType())
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

func (w *worker) getJob(q queue.Queue, qn string) (*task.Job, error) {
	msg, err := q.Pop(w.ctx, qn)
	if err == nil {
		return task.NewJob(q, msg), nil
	}
	return nil, err
}

func (w *worker) handleJobDone(_ task.Task, job *task.Job) (err error) {
	return job.Delete(w.ctx)
}

func (w *worker) handleJobError(task task.Task, job *task.Job, _ error) {
	w.logger.Infof("job fail, id=%s, name=%s", job.Id(), job.Name())
	if !job.IsReachMacAttempts() {
		w.logger.Infof("job retry, id=%s, name=%s", job.Id(), job.Name())
		_ = w.retry(task, job)
	} else {
		_ = job.Delete(w.ctx)
	}
	w.em.Dispatch(evt.NewJobErrorEvent(task, job))
}

func (w *worker) handleError(err error) {
	w.em.Dispatch(evt.NewWorkErrorEvent(err))
}

func (w *worker) retry(task task.Task, job *task.Job) (err error) {
	backoff := task.Backoff()
	return job.Release(w.ctx, backoff)
}
