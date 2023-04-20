package task

import (
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
	var mds = []Middleware{FuncHandler(func(p *Passable, next func(passable *Passable)) {
		println(p.t.GetName(), 1)
		next(p)
	}), FuncHandler(func(p *Passable, next func(passable *Passable)) {
		println(p.t.GetName(), 2)
		next(p)
	})}
	p := NewMiddlewarePipeline()
	p.Through(mds).Send(NewPassable(NewTask(), nil)).Then(func() {
		println("end")
	})
}
