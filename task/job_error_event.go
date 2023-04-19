package task

type JobErrorEvent struct {
	task Task
	job  *Job
}

func NewJobErrorEvent(t Task, j *Job) *JobErrorEvent {
	return &JobErrorEvent{task: t, job: j}
}

func (j *JobErrorEvent) Name() string {
	return "JobErrorEvent"
}
