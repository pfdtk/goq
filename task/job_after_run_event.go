package task

type JobAfterRunEventValue struct {
	Task Task
	Job  *Job
}

type JobAfterRunEvent struct {
	value *JobAfterRunEventValue
}

func NewJobAfterRunEvent(t Task, j *Job) *JobAfterRunEvent {
	return &JobAfterRunEvent{&JobAfterRunEventValue{Task: t, Job: j}}
}

func (j *JobAfterRunEvent) Name() string {
	return "JobAfterRunEvent"
}

func (j *JobAfterRunEvent) Value() any {
	return j.value
}
