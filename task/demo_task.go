package task

import (
	"context"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/queue"
	"time"
)

var Conn = "test"

type DemoTask struct {
	BaseTask
	logger logger.Logger
}

func NewDemoTask() *DemoTask {
	option := &Option{
		Name:      "test",
		OnConnect: Conn,
		QueueType: queue.Redis,
		OnQueue:   "default",
		Status:    Active,
		Priority:  0,
		Retries:   0,
		Timeout:   500,
	}
	return &DemoTask{
		BaseTask: BaseTask{Option: option},
		logger:   logger.GetLogger(),
	}
}

func (t *DemoTask) Run(_ context.Context, j *Job) (any, error) {
	time.Sleep(8 * time.Second)
	t.logger.Info(string(j.RawMessage().Payload))
	return "test", nil
}

func (t *DemoTask) Beforeware() []Middleware {
	return []Middleware{
		NewMaxWorkerLimiter("test", t.GetName(), 1, 10),
		//task.NewLeakyBucketLimiter("test", t.GetName(), redis_rate.PerMinute(1)),
	}
}

func (t *DemoTask) Processware() []Middleware {
	return nil
}
