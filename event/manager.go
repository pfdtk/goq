package event

import (
	"sync"
)

var manager = newManager()

type Event interface {
	Name() string
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
	manager.Listen(e, h, false)
}

func Listens(events []Event, h Handler) {
	for i := range events {
		e := events[i]
		manager.Listen(e, h, false)
	}
}

// IListen Do not use this function!, this function is use for internal
func IListen(e Event, h Handler) {
	manager.Listen(e, h, true)
}

func IListens(events []Event, h Handler) {
	for i := range events {
		e := events[i]
		manager.Listen(e, h, true)
	}
}

func (m *Manager) Listen(e Event, h Handler, internal bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	en := e.Name()
	if !internal {
		handles := m.event[en]
		handles = append(handles, h)
		m.event[en] = handles
	} else {
		handles := m.internalEvent[en]
		handles = append(handles, h)
		m.internalEvent[en] = handles
	}
}

func Dispatch(e Event) {
	manager.Dispatch(e)
}

func (m *Manager) Dispatch(e Event) {
	handlers := m.event[e.Name()]
	internalHandlers := m.internalEvent[e.Name()]
	size := len(handlers) + len(internalHandlers)
	mergedHandlers := make([]Handler, size)
	// internal event first
	copy(mergedHandlers, internalHandlers)
	copy(mergedHandlers[len(internalHandlers):], handlers)
	for _, h := range mergedHandlers {
		ct := h(e)
		if ct == false {
			break
		}
	}
}
