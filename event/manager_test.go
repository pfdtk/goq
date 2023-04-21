package event

import "testing"

type TestEvent struct {
}

func (t *TestEvent) Name() string {
	return "test"
}

func NewTestHandle() Handler {
	return func(e Event) bool {
		return true
	}
}

func TestManager_Listen(t *testing.T) {
	m := newManager()
	m.Listen(&TestEvent{}, NewTestHandle())
}

func TestManager_Dispatch(t *testing.T) {
	m := newManager()
	e := &TestEvent{}
	m.Listen(e, NewTestHandle())
	m.Dispatch(e)
}
