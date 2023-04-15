package task

import "github.com/pfdtk/goq/base"

type Option struct {
	Name      string
	OnConnect string
	QueueType base.QueueType
	OnQueue   string
	Status    base.TaskStatus
	Backoff   uint
	Priority  int
	Retries   uint
}
