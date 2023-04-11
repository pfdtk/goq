package iface

import (
	"context"
	"github.com/pfdtk/goq/common/cst"
	"github.com/pfdtk/goq/common/job"
)

type Task interface {
	Run(context.Context, *job.Job) (any, error)
	// QueueType which queue type, e.g.: RedisQueue
	QueueType() cst.Type
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
	Backoff() uint
	// Priority task priority
	Priority() int
	// Retries how many times to retry process
	Retries() int
}
