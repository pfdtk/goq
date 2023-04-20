package event

type Event interface {
	Name() string
}

type Handler interface {
	Handle(e Event) (bool, error)
}

type FuncHandler func(e Event) (bool, error)

func (f FuncHandler) Handle(e Event) (bool, error) {
	return f(e)
}
