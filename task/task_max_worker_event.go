package task

type MaxWorkerEvent struct {
	err error
}

func NewMaxWorkerErrorEvent(err error) *MaxWorkerEvent {
	return &MaxWorkerEvent{err: err}
}

func (j *MaxWorkerEvent) Name() string {
	return "TaskMaxWorkerEvent"
}

func (j *MaxWorkerEvent) Value() any {
	return j.err
}
