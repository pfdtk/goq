package task

type JobBeforeRunEvent struct {
	task Task
	job  *Job
}

func NewJobBeforeRunEvent(t Task, j *Job) *JobBeforeRunEvent {
	return &JobBeforeRunEvent{task: t, job: j}
}

func (j *JobBeforeRunEvent) Name() string {
	return "JobBeforeRunEvent"
}
