package task

type JobErrorEvent struct {
	task Task
	job  *Job
	err  error
}

func NewJobErrorEvent(t Task, j *Job, err error) *JobErrorEvent {
	return &JobErrorEvent{task: t, job: j, err: err}
}

func (j *JobErrorEvent) Name() string {
	return "JobErrorEvent"
}
