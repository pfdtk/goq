package task

type Middleware interface {
	Handle(t Task, j Job) (bool, error)
}

type FuncHandler func(t Task, j Job) (bool, error)

func (f FuncHandler) Handle(t Task, j Job) (bool, error) {
	return f(t, j)
}
