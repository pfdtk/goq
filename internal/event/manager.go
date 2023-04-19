package event

import "sync"

var manager = newManager()

type Event interface {
	Name() string
}

type Handler interface {
	Handle(e Event) (bool, error)
}

type Manager struct {
	event map[string][]Handler
	lock  sync.Mutex
}

func newManager() *Manager {
	return &Manager{
		event: make(map[string][]Handler),
	}
}

func Listen(e Event, h Handler) {
	manager.Listen(e, h)
}

func (m *Manager) Listen(e Event, h Handler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	en := e.Name()
	handles := m.event[en]
	handles = append(handles, h)
	m.event[en] = handles
}

func Dispatch(e Event) {
	manager.Dispatch(e)
}

func (m *Manager) Dispatch(e Event) {
	handlers := m.event[e.Name()]
	for _, h := range handlers {
		ct, _ := h.Handle(e)
		if ct == false {
			break
		}
	}
}
