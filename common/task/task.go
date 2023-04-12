package task

import (
	"github.com/pfdtk/goq/common/cst"
)

type BaseTask struct {
}

func (b BaseTask) QueueType() cst.Type {
	return cst.Redis
}

func (b BaseTask) OnQueue() string {
	return "default"
}

func (b BaseTask) GetStatus() uint32 {
	return cst.Active
}

func (b BaseTask) CanRun() bool {
	return true
}

func (b BaseTask) Backoff() uint {
	return 0
}

func (b BaseTask) Priority() int {
	return cst.P0
}

func (b BaseTask) Retries() uint {
	return 0
}
