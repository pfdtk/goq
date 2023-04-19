package task

import (
	"github.com/pfdtk/goq/queue"
)

type JobAddEvent struct {
	task    Task
	message *queue.Message
}

func NewJobAddEvent(t Task, message *queue.Message) *JobAddEvent {
	return &JobAddEvent{task: t, message: message}
}

func (j *JobAddEvent) Name() string {
	return "JobAddEvent"
}
