package task

import "github.com/pfdtk/goq/pipeline"

type Middleware interface {
	pipeline.Handler
}

type Passable struct {
	task Task
	job  *Job
}

func NewPassable(t Task, j *Job) *Passable {
	return &Passable{task: t, job: j}
}

func CastMiddleware(m []Middleware) []pipeline.Handler {
	h := make([]pipeline.Handler, len(m), len(m))
	for i, v := range m {
		h[i] = pipeline.Handler(v)
	}
	return h
}
