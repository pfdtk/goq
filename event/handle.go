package event

type Event interface {
	Name() string
}

type Handler interface {
	Handle(e Event) bool
}

type FuncHandler func(e Event) bool

func (f FuncHandler) Handle(e Event) bool {
	return f(e)
}
