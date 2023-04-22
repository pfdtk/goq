package task

type JobAfterRunEventValue struct {
	task Task
	job  *Job
}

type JobAfterRunEvent struct {
	value *JobAfterRunEventValue
}

func NewJobAfterRunEvent(t Task, j *Job) *JobAfterRunEvent {
	return &JobAfterRunEvent{&JobAfterRunEventValue{task: t, job: j}}
}

func (j *JobAfterRunEvent) Name() string {
	return "JobAfterRunEvent"
}

func (j *JobAfterRunEvent) Value() any {
	return j.value
}
