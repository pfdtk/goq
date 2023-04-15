package task

import (
	"github.com/pfdtk/goq/base"
)

type BaseTask struct {
}

func (b *BaseTask) GetJob() *Job {
	return nil
}

func (b *BaseTask) QueueType() base.QueueType {
	return base.Redis
}

func (b *BaseTask) OnQueue() string {
	return "default"
}

func (b *BaseTask) GetStatus() uint32 {
	return base.Active
}

func (b *BaseTask) CanRun() bool {
	return true
}

func (b *BaseTask) Backoff() uint {
	return 0
}

func (b *BaseTask) Priority() int {
	return base.P0
}

func (b *BaseTask) Retries() uint {
	return 0
}
