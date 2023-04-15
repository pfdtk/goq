package task

import (
	"github.com/pfdtk/goq/queue"
)

var (
	defaultConnect   = "default"
	defaultQueueType = queue.Redis
	defaultQueue     = "default"
)

type BaseTask struct {
	Option *Option
}

func (b *BaseTask) GetName() string {
	return b.Option.Name
}

func (b *BaseTask) OnConnect() string {
	if b.Option.OnConnect == "" {
		return defaultConnect
	}
	return b.Option.OnConnect
}

func (b *BaseTask) QueueType() queue.Type {
	if b.Option.QueueType == "" {
		return defaultQueueType
	}
	return b.Option.QueueType
}

func (b *BaseTask) OnQueue() string {
	if b.Option.OnQueue == "" {
		return defaultQueue
	}
	return b.Option.OnQueue
}

func (b *BaseTask) Status() Status {
	if b.Option.Status == 0 {
		return Active
	}
	return b.Option.Status
}

func (b *BaseTask) CanRun() bool {
	return true
}

func (b *BaseTask) Backoff() uint {
	return b.Option.Backoff
}

func (b *BaseTask) Priority() int {
	return b.Option.Priority
}

func (b *BaseTask) Retries() uint {
	return b.Option.Retries
}
