package goq

import (
	"context"
	"github.com/pfdtk/goq/internal/queue"
	"time"
)

type Task interface {
	Run(context.Context, *Job) error
	// QueueType which queue type, e.g.: RedisQueue
	QueueType() queue.Type
	// OnConnect queue on which connect
	OnConnect() string
	// OnQueue these are some queue on one connect, specify which one to consuming
	OnQueue() string
	// GetName task name
	GetName() string
	// CanRun check if task can run, e.g.: rate limit
	CanRun() bool
	// Backoff retry delay when exception throw
	Backoff() time.Time
	// Priority task priority
	Priority() int
}
