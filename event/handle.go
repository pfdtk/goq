package event

type Event interface {
	Name() string
}

type Handler func(e Event) bool
