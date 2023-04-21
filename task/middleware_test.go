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
	var mds = []Middleware{pipeline.HandlerFunc(func(p any, next pipeline.Next) any {
		println(1)
		return next(p)
	}), pipeline.HandlerFunc(func(p any, next pipeline.Next) any {
		println(2)
		nr := next(p)
		return nr
	})}
	p := pipeline.NewPipeline()
	hds := CastMiddleware(mds)
	r := p.Through(hds).Send(NewPassable(NewTask(), nil)).Then(func(_ any) any {
		println("end")
		return "return end"
	})
	println(r)
}
