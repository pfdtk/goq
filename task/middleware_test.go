package task

import (
	"github.com/pfdtk/goq/pipeline"
	"github.com/pfdtk/goq/queue"
	"testing"
)

type TestTask struct {
	BaseTask
}

func NewTask() *TestTask {
	option := &Option{
		Name:      "test",
		OnConnect: "test",
		QueueType: queue.Redis,
		OnQueue:   "default",
		Status:    Active,
		Priority:  0,
		Retries:   0,
		Timeout:   11,
	}
	return &TestTask{
		BaseTask{Option: option},
	}
}

func TestNewMiddlewarePipeline(t *testing.T) {
	var mds = []Middleware{MiddlewareFunc(func(p any, next pipeline.Next) any {
		println(1)
		return next(p)
	}), MiddlewareFunc(func(p any, next pipeline.Next) any {
		println(2)
		return next(p)
	})}
	p := pipeline.NewPipeline()
	hds := CastMiddlewareAsPipelineHandler(mds)
	p.Through(hds).Send(NewPassable(NewTask(), nil)).Then(func() any {
		println("end")
		return nil
	})
}
