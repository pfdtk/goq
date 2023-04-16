package event

type WorkErrorEvent struct {
	err error
}

func NewWorkErrorEvent(e error) *WorkErrorEvent {
	return &WorkErrorEvent{err: e}
}

func (j *WorkErrorEvent) Name() string {
	return "JobErrorEvent"
}
