package pipeline

import "testing"

func TestNewPipeline(t *testing.T) {
	var mds = []Handler{func(p any, next func(p any) any) any {
		println(1)
		return next(p)
	}, func(p any, next func(p any) any) any {
		println(2)
		return next(p)
	}}
	p := NewPipeline()
	p.Through(mds...).Send("test").Then(func(_ any) any {
		println("end")
		return nil
	})
}
