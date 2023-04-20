package task

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
