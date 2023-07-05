package event

import (
	"github.com/pfdtk/goq/internal/errors"
	"sync"
)

var manager = newManager()

type Event interface {
	Name() string
	Value() any
}

type Handler func(e Event) bool

type Manager struct {
	event map[string][]Handler
	// internalEvent internal event
	internalEvent map[string][]Handler
	lock          sync.Mutex
}

func newManager() *Manager {
	return &Manager{
		event:         make(map[string][]Handler),
		internalEvent: make(map[string][]Handler),
	}
}

func Listen(e Event, h Handler) {
	manager.Listen(e, h)
}

// Listens multi events
func Listens(events []Event, h Handler) {
	for i := range events {
		Listen(events[i], h)
	}
}

// IListen Do not use this function!, this function is use for internal
func IListen(e Event, h Handler) {
	manager.IListen(e, h)
}

func IListens(events []Event, h Handler) {
	for i := range events {
		IListen(events[i], h)
	}
}

func (m *Manager) Listen(e Event, h Handler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	en := e.Name()
	handles := m.event[en]
	handles = append(handles, h)
	m.event[en] = handles
}

func (m *Manager) IListen(e Event, h Handler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	en := e.Name()
	handles := m.internalEvent[en]
	handles = append(handles, h)
	m.internalEvent[en] = handles
}

func Dispatch(e Event) {
	manager.Dispatch(e)
}

func RecoverDispatch(e Event) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.NewPanicError(x)
		}
	}()
	manager.Dispatch(e)
	return
}

func (m *Manager) Dispatch(e Event) {
	handlers := m.event[e.Name()]
	internalHandlers := m.internalEvent[e.Name()]
	size := len(handlers) + len(internalHandlers)
	mergedHandlers := make([]Handler, size)
	// internal event first
	copy(mergedHandlers, internalHandlers)
	copy(mergedHandlers[len(internalHandlers):], handlers)
	for i := range mergedHandlers {
		ct := mergedHandlers[i](e)
		if ct == false {
			break
		}
	}
}
