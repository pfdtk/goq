package event

import "github.com/pfdtk/goq/task"

type JobAfterRunEvent struct {
	task task.Task
	job  *task.Job
}

func NewJobAfterRunEvent(t task.Task, j *task.Job) *JobAfterRunEvent {
	return &JobAfterRunEvent{task: t, job: j}
}

func (j *JobAfterRunEvent) Name() string {
	return "JobAfterRunEvent"
}
