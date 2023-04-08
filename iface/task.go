package iface

import (
	"context"
	"github.com/pfdtk/goq/internal/common"
	"github.com/pfdtk/goq/internal/queue"
	"time"
)

type Task interface {
	Run(context.Context, *common.Job) error
	// QueueType which queue type, e.g.: RedisQueue
	QueueType() queue.Type
	// OnConnect queue on which connect
	OnConnect() string
	// OnQueue these are some queue on one connect, specify which one to consuming
	OnQueue() string
	// GetName task name
	GetName() string
	// GetStatus 0 is disabled, 1 is active
	GetStatus() uint32
	// CanRun check if task can run, e.g.: rate limit
	CanRun() bool
	// Backoff retry delay when exception throw
	Backoff() time.Time
	// Priority task priority
	Priority() int
	// Retries how many times to retry process
	Retries() int
}
