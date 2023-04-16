package event

import "github.com/pfdtk/goq/task"

type JobErrorEvent struct {
	task task.Task
	job  *task.Job
}

func NewJobErrorEvent(t task.Task, j *task.Job) *JobErrorEvent {
	return &JobErrorEvent{task: t, job: j}
}

func (j *JobErrorEvent) Name() string {
	return "JobErrorEvent"
}
