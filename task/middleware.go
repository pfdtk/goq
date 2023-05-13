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
	callback func()
}

func NewPopPassable() *PopPassable {
	return &PopPassable{}
}

func (p *PopPassable) SetCallback(fn func()) {
	p.callback = fn
}

func (p *PopPassable) GetCallback() func() {
	return p.callback
}

func (p *PopPassable) Callback() {
	if p.callback != nil {
		p.callback()
	}
}

func CastMiddleware(m []Middleware) []pipeline.Handler {
	h := make([]pipeline.Handler, len(m), len(m))
	for i, v := range m {
		h[i] = pipeline.Handler(v)
	}
	return h
}
