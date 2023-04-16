package event

import "sync"

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

func NewManager() *Manager {
	return &Manager{
		event: make(map[string][]Handler),
	}
}

func (m *Manager) Listen(e Event, h Handler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	handles := m.event[e.Name()]
	handles = append(handles, h)
	m.event[e.Name()] = handles
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
