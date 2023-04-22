package task

type JobBeforeRunEventValue struct {
	task Task
	job  *Job
}

type JobBeforeRunEvent struct {
	value *JobBeforeRunEventValue
}

func NewJobBeforeRunEvent(t Task, j *Job) *JobBeforeRunEvent {
	return &JobBeforeRunEvent{&JobBeforeRunEventValue{task: t, job: j}}
}

func (j *JobBeforeRunEvent) Name() string {
	return "JobBeforeRunEvent"
}

func (j *JobBeforeRunEvent) Value() any {
	return j.value
}
