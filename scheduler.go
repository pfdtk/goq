package goq

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/pfdtk/goq/event"
	e "github.com/pfdtk/goq/internal/errors"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/task"
	"github.com/robfig/cron/v3"
	"time"
)

type cronTask struct {
	spec string
	task task.Task
}

type scheduler struct {
	cron   *cron.Cron
	ctx    context.Context
	tasks  []*cronTask
	logger logger.Logger
}

type payload struct {
	CreatedAt int64 `json:"created_at"`
}

func newScheduler(ctx context.Context, s *Server) *scheduler {
	loc := s.schedulerLocation
	if loc == nil {
		loc = time.UTC
	}
	return &scheduler{
		cron:   cron.New(cron.WithLocation(loc)),
		ctx:    ctx,
		logger: s.logger,
		tasks:  s.cronTasks,
	}
}

func (s *scheduler) mustStartScheduler() {
	for i := range s.tasks {
		t := s.tasks[i]
		err := s.register(t.spec, t.task)
		if err != nil {
			panic(err)
		}
	}
	s.cron.Start()
}

func (s *scheduler) stopScheduler() {
	s.logger.Info("stopping scheduler...")
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.logger.Info("scheduler stopped")
}

func (s *scheduler) register(spec string, t task.Task) error {
	_, err := s.cron.AddFunc(spec, func() {
		var err error
		defer func() {
			if x := recover(); x != nil {
				err = errors.Join(err, e.NewPanicError(x))
			}
			if err != nil {
				event.Dispatch(NewSchedulerErrorEvent(err))
			}
		}()
		s.logger.Infof("scheduler trigger, name=%s", t.GetName())
		// when scheduler, we will dispatch task to queue, then process by worker
		payload, err := json.Marshal(&payload{CreatedAt: time.Now().Unix()})
		if err == nil {
			err = client.DispatchContext(s.ctx, t, payload)
		}
	})
	return err
}
