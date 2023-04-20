package pipeline

import "testing"

func TestNewPipeline(t *testing.T) {
	var mds = []Handler{HandlerFunc(func(p any, next Next) any {
		println(1)
		return next(p)
	}), Handler(HandlerFunc(func(p any, next Next) any {
		println(2)
		return next(p)
	}))}
	p := NewPipeline()
	p.Through(mds).Send("test").Then(func() any {
		println("end")
		return nil
	})
}
