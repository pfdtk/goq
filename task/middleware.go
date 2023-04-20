package task

import "github.com/pfdtk/goq/pipeline"

type Middleware interface {
	// Handle continue if return True, else break
	Handle(p any, next func(passable any))
}

type MiddlewareFunc func(p any, next func(passable any))

func (f MiddlewareFunc) Handle(p any, next func(passable any)) {
	f(p, next)
}

type Passable struct {
	task Task
	job  *Job
}

func NewPassable(t Task, j *Job) *Passable {
	return &Passable{task: t, job: j}
}

func CastMiddlewareAsPipelineHandler(m []Middleware) []pipeline.Handler {
	h := make([]pipeline.Handler, len(m), len(m))
	for i, v := range m {
		h[i] = pipeline.Handler(v)
	}
	return h
}
