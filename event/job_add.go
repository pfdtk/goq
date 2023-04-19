package event

import (
	"github.com/pfdtk/goq/queue"
	"github.com/pfdtk/goq/task"
)

type JobAddEvent struct {
	task    task.Task
	message *queue.Message
}

func NewJobAddEvent(t task.Task, message *queue.Message) *JobAddEvent {
	return &JobAddEvent{task: t, message: message}
}

func (j *JobAddEvent) Name() string {
	return "JobAddEvent"
}
