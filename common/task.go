package common

import (
	"github.com/pfdtk/goq/common/cst"
	"time"
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

func (b BaseTask) Backoff() time.Time {
	return time.Now().Add(0 * time.Second)
}

func (b BaseTask) Priority() int {
	return cst.P0
}

func (b BaseTask) Retries() int {
	return 0
}
