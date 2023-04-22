package task

type MaxWorkerErrorEvent struct {
	err error
}

func NewMaxWorkerErrorEvent(err error) *MaxWorkerErrorEvent {
	return &MaxWorkerErrorEvent{err: err}
}

func (j *MaxWorkerErrorEvent) Name() string {
	return "TaskMaxWorkerErrorEvent"
}

func (j *MaxWorkerErrorEvent) Value() any {
	return j.err
}
