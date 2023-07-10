package task

import (
	"github.com/pfdtk/goq/pipeline"
	"testing"
)

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
	r := p.Through(hds...).Send(NewRunPassable(NewDemoTask(), nil)).Then(func(_ any) any {
		return true
	})
	rs := r.(bool)
	if rs != true {
		t.Error(rs)
	}
}
