package test

import (
	"context"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/queue"
	"github.com/pfdtk/goq/task"
	"time"
)

var Conn = "test"

type DemoTask struct {
	task.BaseTask
	logger logger.Logger
}

func NewDemoTask() *DemoTask {
	option := &task.Option{
		Name:      "test",
		OnConnect: Conn,
		QueueType: queue.Redis,
		OnQueue:   "default",
		Status:    task.Active,
		Priority:  0,
		Retries:   0,
		Timeout:   500,
	}
	return &DemoTask{
		BaseTask: task.BaseTask{Option: option},
		logger:   logger.GetLogger(),
	}
}

func (t *DemoTask) Run(_ context.Context, j *task.Job) (any, error) {
	time.Sleep(8 * time.Second)
	t.logger.Info(string(j.RawMessage().Payload))
	return "test", nil
}

func (t *DemoTask) Beforeware() []task.Middleware {
	return []task.Middleware{
		task.NewMaxWorkerLimiter("test", t.GetName(), 1, 10),
		//task.NewLeakyBucketLimiter("test", t.GetName(), redis_rate.PerMinute(1)),
	}
}

func (t *DemoTask) Processware() []task.Middleware {
	return nil
}
