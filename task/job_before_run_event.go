package task

type JobBeforeRunEventValue struct {
	Task Task
	Job  *Job
}

type JobBeforeRunEvent struct {
	value *JobBeforeRunEventValue
}

func NewJobBeforeRunEvent(t Task, j *Job) *JobBeforeRunEvent {
	return &JobBeforeRunEvent{&JobBeforeRunEventValue{Task: t, Job: j}}
}

func (j *JobBeforeRunEvent) Name() string {
	return "JobBeforeRunEvent"
}

func (j *JobBeforeRunEvent) Value() any {
	return j.value
}
