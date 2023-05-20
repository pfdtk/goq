package task

type JobErrorEventValue struct {
	Task Task
	Job  *Job
	Err  error
}

type JobErrorEvent struct {
	value *JobErrorEventValue
}

func NewJobErrorEvent(t Task, j *Job, err error) *JobErrorEvent {
	return &JobErrorEvent{&JobErrorEventValue{Task: t, Job: j, Err: err}}
}

func (j *JobErrorEvent) Name() string {
	return "JobErrorEvent"
}

func (j *JobErrorEvent) Value() any {
	return j.value
}
