package task

import "github.com/pfdtk/goq/pipeline"

type Middleware pipeline.Handler

type RunPassable struct {
	task Task
	job  *Job
}

func NewRunPassable(t Task, j *Job) *RunPassable {
	return &RunPassable{task: t, job: j}
}

type PopPassable struct {
	Callback func()
}

func NewPopPassable() *PopPassable {
	return &PopPassable{}
}

func (p *PopPassable) ExecCallback() {
	if p.Callback != nil {
		p.Callback()
	}
}

func CastMiddleware(m []Middleware) []pipeline.Handler {
	h := make([]pipeline.Handler, len(m), len(m))
	for i, v := range m {
		h[i] = pipeline.Handler(v)
	}
	return h
}
