package state

import (
	"github.com/pfdtk/goq/queue"
	"time"
)

type State string

var (
	Pending State = "pending"
	Started State = "started"
	Success State = "success"
	Failure State = "failure"
)

type TaskState struct {
	TaskId    string    `json:"task_id"`
	TaskName  string    `json:"task_name"`
	State     State     `json:"state"`
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
}

func NewTaskState(msg *queue.Message, s State) *TaskState {
	st := &TaskState{
		TaskId:    msg.ID,
		TaskName:  msg.Type,
		State:     s,
		CreatedAt: time.Now(),
	}
	return st
}

func NewFailureTaskState(msg *queue.Message, err error) *TaskState {
	st := &TaskState{
		TaskId:    msg.ID,
		TaskName:  msg.Type,
		State:     Failure,
		Error:     err.Error(),
		CreatedAt: time.Now(),
	}
	return st
}
