package goq

import (
	"context"
	"errors"
	"fmt"
	"github.com/pfdtk/goq/event"
	qm "github.com/pfdtk/goq/internal/queue"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/pipeline"
	"github.com/pfdtk/goq/task"
	"runtime/debug"
	"sync"
	"time"
)

var (
	JobEmptyError      = errors.New("no jobs are ready for processing")
	QueueNotFoundError = errors.New("queue not found")
)

type worker struct {
	tasks        *sync.Map
	sortedTasks  []task.Task
	delayConsume map[string]time.Time
	wg           *sync.WaitGroup
	lock         sync.Mutex
	stopRun      chan struct{}
	maxWorker    chan struct{}
	jobChannel   chan *task.Job
	ctx          context.Context
	logger       logger.Logger
	pl           *pipeline.Pipeline
}

func newWorker(ctx context.Context, s *Server) *worker {
	w := &worker{
		wg:           &s.wg,
		tasks:        &s.tasks,
		sortedTasks:  task.SortTask(&s.tasks),
		maxWorker:    make(chan struct{}, s.maxWorker),
		stopRun:      make(chan struct{}),
		jobChannel:   make(chan *task.Job, s.maxWorker),
		ctx:          ctx,
		logger:       s.logger,
		pl:           pipeline.NewPipeline(),
		delayConsume: make(map[string]time.Time),
	}
	return w
}

func (w *worker) mustStartConsuming() {
	w.startPop()
	w.startWork()
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
	var err error
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(fmt.Sprintf("Panic Error: %+v;\nStack: %s", x, string(debug.Stack())))
		}
		if err != nil {
			w.handleError(err)
		}
	}()
	j, err := w.getNextJob()
	switch {
	case errors.Is(err, JobEmptyError):
		time.Sleep(time.Second)
		<-w.maxWorker
		return
	case err != nil:
		<-w.maxWorker
		return
	}
	w.logger.Infof("job received, name=%s, id=%s", j.Name(), j.Id())
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
				w.logger.Infof("job processing, name=%s, id=%s", j.Name(), j.Id())
				go w.runJob(j)
			}
		}
	}()
}

func (w *worker) runJob(job *task.Job) {
	defer func() {
		if x := recover(); x != nil {
			err := errors.New(fmt.Sprintf("Panic Error: %+v;\nStack: %s", x, string(debug.Stack())))
			w.handleError(err)
		}
	}()
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
		w.handleJobTimeoutError(job)
		return
	}
}

func (w *worker) perform(job *task.Job) (res any, err error) {
	var t task.Task = nil
	// recover err from tasks, so that program will not exit
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(fmt.Sprintf("Panic Error: %+v;\nStack: %s", x, string(debug.Stack())))
			if t != nil {
				w.handleJobError(t, job, err)
			} else {
				w.handleError(err)
			}
		}
	}()
	v, ok := w.tasks.Load(job.Name())
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
			return err
		} else {
			w.handleJobDone(t, job)
			return nil
		}
	}
	mds := task.CastMiddleware(t.Processware())
	// run task through middleware
	passable := task.NewRunPassable(t, job)
	w.pl.Send(passable).Through(mds).Then(fn)
	return
}

func (w *worker) getNextJob() (*task.Job, error) {
	for i := range w.sortedTasks {
		t := w.sortedTasks[i]
		pp := task.NewPopPassable()
		if !w.shouldGetNextJob(t, pp) {
			continue
		}
		w.logger.Infof("start to get job, name=%s", t.GetName())
		j, err := w.getJob(t)
		if err == nil {
			j.Then(pp.Callback)
			return j, nil
		}
		w.logger.Infof("no job for process, name=%s", t.GetName())
		// no message on queue, delay some seconds before next time
		w.delayGetNextJob(t)
		// exec callback func
		pp.ExecCallback()
	}
	return nil, JobEmptyError
}

func (w *worker) shouldGetNextJob(t task.Task, pp *task.PopPassable) bool {
	if t.Status() == task.Disable {
		return false
	}
	delayAt, ok := w.delayConsume[t.OnQueue()]
	if ok && time.Now().Before(delayAt) {
		return false
	}
	// check if task can pop message through middleware,
	// and middleware handle should return a bool value
	mds := task.CastMiddleware(t.Beforeware())
	res := w.pl.Send(pp).Through(mds).Then(func(_ any) any {
		return true
	})
	can, ok := res.(bool)
	if !ok || !can {
		w.delayGetNextJob(t)
		return false
	}
	return true
}

func (w *worker) delayGetNextJob(task task.Task) {
	w.logger.Infof("wait for 3 second before next time, name=%s", task.GetName())
	w.lock.Lock()
	defer w.lock.Unlock()
	w.delayConsume[task.OnQueue()] = time.Now().Add(3 * time.Second)
}

func (w *worker) getJob(t task.Task) (*task.Job, error) {
	q := qm.GetQueue(t.OnConnect(), t.QueueType())
	if q == nil {
		w.logger.Errorf("queue not found, name=%s, conn=%s", t.OnQueue(), t.OnConnect())
		w.handleError(QueueNotFoundError)
		return nil, QueueNotFoundError
	}
	msg, err := q.Pop(w.ctx, t.OnQueue())
	if err == nil {
		return task.NewJob(q, msg), nil
	}
	return nil, err
}

func (w *worker) handleJobDone(_ task.Task, job *task.Job) {
	w.logger.Infof("job processed, name=%s, id=%s", job.Name(), job.Id())
	job.Success()
	err := job.Delete(w.ctx)
	if err != nil {
		event.Dispatch(NewWorkErrorEvent(err))
	}
}

func (w *worker) handleJobError(t task.Task, job *task.Job, err error) {
	w.logger.Infof("job fail, name=%s, id=%s", job.Name(), job.Id())
	job.Fail()
	var e error
	if !job.IsReachMacAttempts() {
		w.logger.Infof("job retry, name=%s, id=%s", job.Name(), job.Id())
		e = job.Release(w.ctx, t.Backoff())
	} else {
		e = job.Delete(w.ctx)
	}
	if e != nil {
		err = errors.Join(err, e)
	}
	event.Dispatch(task.NewJobErrorEvent(t, job, err))
}

func (w *worker) handleJobTimeoutError(job *task.Job) {
	w.logger.Infof("job exec timeout, id=%s, name=%s", job.Id(), job.Name())
	job.Fail()
	event.Dispatch(task.NewJobExecTimeoutEvent(job))
}

func (w *worker) handleError(err error) {
	event.Dispatch(NewWorkErrorEvent(err))
}
