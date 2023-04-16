package task

import (
	"github.com/pfdtk/goq/queue"
)

type Option struct {
	// Name of task
	Name string
	// OnConnect which connect
	OnConnect string
	// QueueType queue type, redis or sqs...
	QueueType queue.Type
	// OnQueue maybe multi queue on a connection
	OnQueue string
	// Status of task, if is disabled, task will not run
	Status Status
	// Backoff seconds to process task again when fail
	Backoff uint
	// Priority of task, smaller is more priority
	Priority int
	// Retries times to retry when task fail
	Retries uint
	// Timeout job timeout
	Timeout int64
	// UniqueId use for unique message for a time period
	UniqueId string
}
