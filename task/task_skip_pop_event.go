package task

type SkipPopEvent struct {
}

func NewSkipPopEvent() *SkipPopEvent {
	return &SkipPopEvent{}
}

func (j *SkipPopEvent) Name() string {
	return "TaskSkipPopEvent"
}

func (j *SkipPopEvent) Value() any {
	return nil
}
