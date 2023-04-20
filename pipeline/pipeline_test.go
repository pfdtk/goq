package pipeline

import "testing"

func TestNewPipeline(t *testing.T) {
	var mds = []Handler{HandlerFunc(func(p any, next func(passable any)) {
		println(1)
		next(p)
	}), Handler(HandlerFunc(func(p any, next func(passable any)) {
		println(2)
		next(p)
	}))}
	p := NewPipeline()
	p.Through(mds).Send("test").Then(func() {
		println("end")
	})
}
