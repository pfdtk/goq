package event

type SchedulerErrorEvent struct {
	err error
}

func NewSchedulerErrorEvent(e error) *SchedulerErrorEvent {
	return &SchedulerErrorEvent{err: e}
}

func (j *SchedulerErrorEvent) Name() string {
	return "SchedulerErrorEvent"
}
