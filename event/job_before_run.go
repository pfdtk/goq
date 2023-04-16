package event

import "github.com/pfdtk/goq/task"

type JobBeforeRunEvent struct {
	task task.Task
	job  *task.Job
}

func NewJobBeforeRunEvent(t task.Task, j *task.Job) *JobBeforeRunEvent {
	return &JobBeforeRunEvent{task: t, job: j}
}

func (j *JobBeforeRunEvent) Name() string {
	return "JobBeforeRunEvent"
}
