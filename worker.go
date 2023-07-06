package goq

import (
	"context"
	"errors"
	"github.com/pfdtk/goq/backend"
	"github.com/pfdtk/goq/event"
	e "github.com/pfdtk/goq/internal/errors"
	qm "github.com/pfdtk/goq/internal/queue"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/pipeline"
	"github.com/pfdtk/goq/task"
	"sync"
	"time"
)

var (
	JobEmptyError      = errors.New("no jobs are ready for processing")
	JobTimeoutError    = errors.New("job process timeout")
	QueueNotFoundError = errors.New("queue not found")
	TaskNotFoundError  = errors.New("task not found")
)

type worker struct {
	tasks            *sync.Map
	sortedTasks      []task.Task
	consumeDelayFlag map[string]time.Time
	wg               *sync.WaitGroup
	lock             sync.Mutex
	stopRun          chan struct{}
	maxWorker        chan struct{}
	jobChannel       chan *task.Job
	ctx              context.Context
	logger           logger.Logger
	pl               *pipeline.Pipeline
}

func newWorker(ctx context.Context, s *Server) *worker {
	w := &worker{
		wg:               &s.wg,
		tasks:            &s.tasks,
		sortedTasks:      task.SortTask(&s.tasks),
		maxWorker:        make(chan struct{}, s.maxWorker),
		stopRun:          make(chan struct{}),
		jobChannel:       make(chan *task.Job, s.maxWorker),
		ctx:              ctx,
		logger:           s.logger,
		pl:               pipeline.NewPipeline(),
		consumeDelayFlag: make(map[string]time.Time),
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
				w.logger.Info("stop popping message")
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
			err = errors.Join(err, e.NewPanicError(x))
		}
		if err != nil {
			time.Sleep(time.Second)
			// we should release token when error
			<-w.maxWorker
			w.handleError(err)
		}
	}()
	// get next job to process
	job, err := w.getNextJob()
	if err != nil {
		return
	}
	// send to channel wait for process
	w.logger.Infof("job received, name=%s, id=%s", job.Name(), job.Id())
	w.jobChannel <- job
}

func (w *worker) startWork() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.stopRun:
				w.logger.Info("stop working")
				return
			case job, ok := <-w.jobChannel:
				// stop working when channel was closed
				if !ok {
					return
				}
				w.logger.Infof("job processing, name=%s, id=%s", job.Name(), job.Id())
				go w.runJob(job)
			}
		}
	}()
}

func (w *worker) runJob(job *task.Job) {
	var t task.Task
	defer func() {
		<-w.maxWorker
		if x := recover(); x != nil {
			if t != nil {
				w.handleJobError(t, job, e.NewPanicError(x))
			} else {
				w.handleError(e.NewPanicError(x))
			}
		}
	}()
	v, ok := w.tasks.Load(job.Name())
	if !ok {
		w.handleError(TaskNotFoundError)
		return
	}
	t = v.(task.Task)
	w.runJobWithTimeout(t, job)
}

func (w *worker) runJobWithTimeout(t task.Task, job *task.Job) {
	// defined timeout
	ctx, cancel := context.WithDeadline(w.ctx, job.TimeoutAt())
	defer func() {
		cancel()
	}()
	// perform in goroutine
	prh := make(chan any, 1)
	go func(prh chan any) {
		w.perform(t, job)
		prh <- struct{}{}
	}(prh)
	select {
	case <-prh:
		return
	case <-ctx.Done():
		w.handleJobError(t, job, JobTimeoutError)
		return
	}
}

func (w *worker) perform(t task.Task, job *task.Job) {
	// recover err from tasks, so that program will not exit
	defer func() {
		if x := recover(); x != nil {
			w.handleJobError(t, job, e.NewPanicError(x))
		}
	}()
	// task runner
	fn := func(_ any) any {
		w.handleJobStart(t, job)
		_, err := t.Run(w.ctx, job)
		if err == nil {
			err = job.DispatchChain(w.ctx)
			if err == nil {
				w.handleJobDone(t, job)
			}
		}
		if err != nil {
			w.handleJobError(t, job, err)
		}
		return err
	}
	// call func through middleware
	w.pl.Send(task.NewRunPassable(t, job)).
		Through(task.CastMiddleware(t.Processware())...).
		Then(fn)

	return
}

func (w *worker) getNextJob() (*task.Job, error) {
	for i := range w.sortedTasks {
		t := w.sortedTasks[i]
		passable := task.NewPopPassable()
		if !w.canGetNextJob(t, passable) {
			continue
		}
		w.logger.Debugf("start to get job, name=%s", t.GetName())
		job, err := w.getJob(t)
		if err == nil {
			job.Then(passable.GetCallback())
			return job, nil
		}
		w.logger.Debugf("no job for process, name=%s", t.GetName())
		// no message on queue, delay some seconds before next time
		w.waitSecondNextJob(t)
		// exec callback func
		passable.ExecCallback()
	}
	return nil, JobEmptyError
}

func (w *worker) canGetNextJob(t task.Task, pp *task.PopPassable) bool {
	if t.Status() == task.Disable {
		return false
	}
	delayAt, ok := w.consumeDelayFlag[t.OnQueue()]
	if ok && time.Now().Before(delayAt) {
		return false
	}
	// check if task can pop message through middleware,
	res := w.pl.Send(pp).
		Through(task.CastMiddleware(t.Beforeware())...).
		Then(func(_ any) any { return true })
	// middleware should return a bool value
	can, ok := res.(bool)
	if !ok || !can {
		w.waitSecondNextJob(t)
		return false
	}
	return true
}

func (w *worker) waitSecondNextJob(task task.Task) {
	w.logger.Debugf("wait for 3 second before next time, name=%s", task.GetName())
	w.lock.Lock()
	defer w.lock.Unlock()
	w.consumeDelayFlag[task.OnQueue()] = time.Now().Add(3 * time.Second)
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

func (w *worker) handleJobStart(t task.Task, job *task.Job) {
	_ = backend.Get().Started(job.RawMessage())
	event.Dispatch(task.NewJobBeforeRunEvent(t, job))
}

func (w *worker) handleJobDone(t task.Task, job *task.Job) {
	w.logger.Infof("job processed, name=%s, id=%s", job.Name(), job.Id())
	job.Success()
	// ack
	if err := job.Delete(w.ctx); err != nil {
		w.handleError(err)
	}
	_ = backend.Get().Success(job.RawMessage())
	event.Dispatch(task.NewJobAfterRunEvent(t, job))
}

func (w *worker) handleJobError(t task.Task, job *task.Job, err error) {
	w.logger.Warnf("job fail, name=%s, id=%s", job.Name(), job.Id())
	job.Failure()
	if !job.IsReachMaxAttempts() {
		w.logger.Infof("job retry, name=%s, id=%s", job.Name(), job.Id())
		// if not success, message will visibility again after ack timeout
		if err := job.Release(w.ctx, t.Backoff()); err != nil {
			w.handleError(err)
		}
	} else {
		// ack
		if err := job.Delete(w.ctx); err != nil {
			w.handleError(err)
		}
		_ = backend.Get().Failure(job.RawMessage(), err)
	}
	event.Dispatch(task.NewJobErrorEvent(t, job, err))
}

func (w *worker) handleError(err error) {
	event.Dispatch(NewWorkErrorEvent(err))
}
