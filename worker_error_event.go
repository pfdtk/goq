package goq

type WorkErrorEvent struct {
	err error
}

func NewWorkErrorEvent(e error) *WorkErrorEvent {
	return &WorkErrorEvent{err: e}
}

func (j *WorkErrorEvent) Name() string {
	return "JobErrorEvent"
}

func (j *WorkErrorEvent) Value() any {
	return j.err
}
