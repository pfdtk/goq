package goq

import (
	"context"
	"errors"
	"github.com/pfdtk/goq/event"
	qm "github.com/pfdtk/goq/internal/queue"
	"github.com/pfdtk/goq/internal/utils"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/pipeline"
	"github.com/pfdtk/goq/queue"
	"github.com/pfdtk/goq/task"
	"runtime/debug"
	"sync"
	"time"
)

var (
	JobEmptyError       = errors.New("no jobs are ready for processing")
	JobExecTimeoutError = errors.New("job exec timeout")
)

type worker struct {
	tasks       *sync.Map
	sortedTasks []task.Task
	wg          *sync.WaitGroup
	stopRun     chan struct{}
	maxWorker   chan struct{}
	jobChannel  chan *task.Job
	ctx         context.Context
	logger      logger.Logger
	pl          *pipeline.Pipeline
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
		pl:         pipeline.NewPipeline(),
	}
	return w
}

func (w *worker) startConsuming() error {
	w.sortedTasks = utils.SortTask(w.tasks)
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
				w.logger.Info("worker received stop sign")
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
	case errors.Is(err, JobEmptyError):
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
			} else {
				w.handleError(err)
			}
		}
	}()
	name := job.Name()
	v, ok := w.tasks.Load(name)
	if ok {
		t = v.(task.Task)
		res, err = w.performThroughMiddleware(t, job)
	}
	return
}

func (w *worker) performThroughMiddleware(t task.Task, job *task.Job) (res any, err error) {
	fn := func(passable any) any {
		event.Dispatch(task.NewJobBeforeRunEvent(t, job))
		res, err = t.Run(w.ctx, job)
		event.Dispatch(task.NewJobAfterRunEvent(t, job))
		if err != nil {
			w.handleJobError(t, job, err)
		} else {
			w.handleJobDone(t, job)
		}
		return nil
	}
	mds := task.CastMiddleware(t.Processware())
	// run task through middleware
	passable := task.NewRunPassable(t, job)
	w.pl.Send(passable).Through(mds).Then(fn)
	return
}

func (w *worker) getNextJob() (*task.Job, error) {
	for _, t := range w.sortedTasks {
		// check if task can pop message through middleware,
		// and middleware handle should return a bool value
		mds := task.CastMiddleware(t.Beforeware())
		passable := task.NewPopPassable()
		skip := w.pl.Send(passable).Through(mds).Then(func(_ any) any {
			// pop by default
			return true
		})
		canPop, ok := skip.(bool)
		if t.Status() == task.Disable || !ok || !canPop {
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
	return nil, JobEmptyError
}

func (w *worker) getJob(q queue.Queue, qn string) (*task.Job, error) {
	msg, err := q.Pop(w.ctx, qn)
	if err == nil {
		return task.NewJob(q, msg), nil
	}
	return nil, err
}

func (w *worker) handleJobDone(_ task.Task, job *task.Job) {
	w.logger.Infof("job processed, id=%s, name=%s", job.Id(), job.Name())
	err := job.Delete(w.ctx)
	if err != nil {
		event.Dispatch(NewWorkErrorEvent(err))
	}
}

func (w *worker) handleJobError(t task.Task, job *task.Job, err error) {
	w.logger.Infof("job fail, id=%s, name=%s", job.Id(), job.Name())
	var e error
	if !job.IsReachMacAttempts() {
		w.logger.Infof("job retry, id=%s, name=%s", job.Id(), job.Name())
		e = w.retry(t, job)
	} else {
		e = job.Delete(w.ctx)
	}
	if e != nil {
		err = errors.Join(err, e)
	}
	event.Dispatch(task.NewJobErrorEvent(t, job, err))
}

func (w *worker) handleError(err error) {
	event.Dispatch(NewWorkErrorEvent(err))
}

func (w *worker) retry(task task.Task, job *task.Job) (err error) {
	backoff := task.Backoff()
	return job.Release(w.ctx, backoff)
}
