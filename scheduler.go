package goq

import (
	"context"
	"encoding/json"
	"github.com/pfdtk/goq/internal/event"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/task"
	"github.com/robfig/cron/v3"
	"time"
)

type CronTask struct {
	spec string
	task task.Task
}

type scheduler struct {
	cron   *cron.Cron
	ctx    context.Context
	tasks  []*CronTask
	logger logger.Logger
}

type Payload struct {
	CreatedAt int64 `json:"created_at"`
}

func newScheduler(ctx context.Context, s *Server) *scheduler {
	return &scheduler{
		cron:   cron.New(cron.WithLocation(time.UTC)),
		ctx:    ctx,
		logger: s.logger,
		tasks:  s.cronTasks,
	}
}

func (s *scheduler) startScheduler() error {
	for _, t := range s.tasks {
		err := s.register(t.spec, t.task)
		if err != nil {
			return err
		}
	}
	s.cron.Start()
	return nil
}

func (s *scheduler) stopScheduler() {
	ctx := s.cron.Stop()
	<-ctx.Done()
}

func (s *scheduler) register(spec string, t task.Task) error {
	_, err := s.cron.AddFunc(spec, func() {
		s.logger.Infof("scheduler trigger, name=%s", t.GetName())
		// when scheduler, we will dispatch task to queue, then process by worker
		payload, err := json.Marshal(&Payload{CreatedAt: time.Now().Unix()})
		if err != nil {
			event.Dispatch(NewSchedulerErrorEvent(err))
		}
		err = client.DispatchContext(s.ctx, t, payload)
		if err != nil {
			event.Dispatch(NewSchedulerErrorEvent(err))
		}
	})
	return err
}
