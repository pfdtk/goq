package task

import (
	"github.com/pfdtk/goq/base"
)

var (
	defaultConnect   = "default"
	defaultQueueType = base.Redis
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

func (b *BaseTask) QueueType() base.QueueType {
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

func (b *BaseTask) Status() uint32 {
	if b.Option.Status == 0 {
		return base.Active
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
