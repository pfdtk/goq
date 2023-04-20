package task

type Middleware interface {
	Handle(t Task, j Job) (bool, error)
}
