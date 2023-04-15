package task

import (
	"context"
	"github.com/pfdtk/goq/base"
)

type TaskStatus uint32

const (
	Active  TaskStatus = 1
	Disable TaskStatus = 2
)

type Task interface {
	Run(context.Context, *Job) (any, error)
	// QueueType which queue type, e.g.: RedisQueue
	QueueType() base.QueueType
	// OnConnect queue on which connect
	OnConnect() string
	// OnQueue these are some queue on one connect, specify which one to consuming
	OnQueue() string
	// GetName task name
	GetName() string
	// Status 0 is disabled, 1 is active
	Status() TaskStatus
	// CanRun check if task can run, e.g.: rate limit
	CanRun() bool
	// Backoff retry delay when exception throw
	Backoff() uint
	// Priority task priority
	Priority() int
	// Retries how many times to retry process
	Retries() uint
}
