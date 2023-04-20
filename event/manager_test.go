package event

import "testing"

type TestEvent struct {
}

func (t *TestEvent) Name() string {
	return "test"
}

type TestHandle struct {
}

func (t *TestHandle) Handle(e Event) bool {
	println(e.Name())
	return true
}

func TestManager_Listen(t *testing.T) {
	m := newManager()
	m.Listen(&TestEvent{}, &TestHandle{})
}

func TestManager_Dispatch(t *testing.T) {
	m := newManager()
	e := &TestEvent{}
	m.Listen(e, &TestHandle{})
	m.Dispatch(e)
}
