package task

import (
	"github.com/pfdtk/goq/queue"
)

type JobAddEventValue struct {
	Task    Task
	Message *queue.Message
}

type JobAddEvent struct {
	value *JobAddEventValue
}

func NewJobAddEvent(t Task, message *queue.Message) *JobAddEvent {
	return &JobAddEvent{&JobAddEventValue{Task: t, Message: message}}
}

func (j *JobAddEvent) Name() string {
	return "JobAddEvent"
}

func (j *JobAddEvent) Value() any {
	return j.value
}
