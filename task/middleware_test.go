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
	var mds = []Middleware{func(p any, next func(p any) any) any {
		println(1)
		return next(p)
	}, func(p any, next func(p any) any) any {
		println(2)
		return next(p)
	}}
	p := pipeline.NewPipeline()
	hds := CastMiddleware(mds)
	r := p.Through(hds).Send(NewRunPassable(NewTask(), nil)).Then(func(_ any) any {
		println("end")
		return "return end"
	})
	println(r)
}
