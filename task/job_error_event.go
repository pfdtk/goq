package task

type JobErrorEventValue struct {
	task Task
	job  *Job
	err  error
}

type JobErrorEvent struct {
	value *JobErrorEventValue
}

func NewJobErrorEvent(t Task, j *Job, err error) *JobErrorEvent {
	return &JobErrorEvent{&JobErrorEventValue{task: t, job: j, err: err}}
}

func (j *JobErrorEvent) Name() string {
	return "JobErrorEvent"
}

func (j *JobErrorEvent) Value() any {
	return j.value
}
