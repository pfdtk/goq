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
	UpdateAt  time.Time `json:"update_at"`
	TTL       int64     `json:"-"`
}

func NewTaskState(msg *queue.Message) *TaskState {
	return nil
}
