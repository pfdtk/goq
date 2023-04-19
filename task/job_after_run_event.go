package task

type JobAfterRunEvent struct {
	task Task
	job  *Job
}

func NewJobAfterRunEvent(t Task, j *Job) *JobAfterRunEvent {
	return &JobAfterRunEvent{task: t, job: j}
}

func (j *JobAfterRunEvent) Name() string {
	return "JobAfterRunEvent"
}
