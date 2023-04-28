package task

import (
	"context"
	"github.com/pfdtk/goq/queue"
)

type Status uint32

const (
	Active  Status = 1
	Disable Status = 2
)

type Task interface {
	// Run a job
	Run(context.Context, *Job) (any, error)
	// QueueType which queue type, e.g.: RedisQueue
	QueueType() queue.Type
	// OnConnect queue on which connect
	OnConnect() string
	// OnQueue these are some queue on one connect, specify which one to consuming
	OnQueue() string
	// GetName task name
	GetName() string
	// Status 0 is disabled, 1 is active
	Status() Status
	// Backoff retry delay when exception throw
	Backoff() uint
	// Priority task priority
	Priority() int
	// Retries how many times to retry process
	Retries() uint
	// Timeout of job
	Timeout() int64
	// Beforeware a list of middleware decide whether to start to process this task
	Beforeware() []Middleware
	// Processware a list of middleware a task run through
	Processware() []Middleware
}
