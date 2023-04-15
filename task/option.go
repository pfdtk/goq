package task

import (
	"github.com/pfdtk/goq/queue"
)

type Option struct {
	Name      string
	OnConnect string
	QueueType queue.Type
	OnQueue   string
	Status    Status
	Backoff   uint
	Priority  int
	Retries   uint
}
