package event

import (
	"testing"
)

type TestEvent struct {
}

func (t *TestEvent) Name() string {
	return "test"
}

func NewTestHandle() Handler {
	return func(e Event) bool {
		println(e.Name())
		return true
	}
}

func TestManager_Listen(t *testing.T) {
	Listen(&TestEvent{}, NewTestHandle())
}

func TestManager_Dispatch(t *testing.T) {
	e := &TestEvent{}
	Listen(e, NewTestHandle())
	IListens([]Event{e}, func(e Event) bool {
		println(e.Name())
		return true
	})
	Dispatch(e)
}
