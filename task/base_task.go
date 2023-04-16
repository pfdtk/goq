package task

import (
	"context"
	"errors"
	"github.com/google/uuid"
	qm "github.com/pfdtk/goq/internal/queue"
	"github.com/pfdtk/goq/queue"
	"time"
)

var (
	defaultConnect   = "default"
	defaultQueueType = queue.Redis
	defaultQueue     = "default"
)

type BaseTask struct {
	Option *Option
}

func (b *BaseTask) Run(_ context.Context, _ *Job) (any, error) {
	return nil, errors.New("func Run is not defined")
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

func (b *BaseTask) Timeout() int64 {
	return b.Option.Timeout
}

func (b *BaseTask) Dispatch(ctx context.Context, payload []byte) (err error) {
	return b.EnQueue(ctx, payload, 0)
}

func (b *BaseTask) EnQueue(ctx context.Context, payload []byte, delay time.Duration) (err error) {
	q := qm.GetQueue(b.OnConnect(), b.QueueType())
	if q == nil {
		return errors.New("fail to get queue")
	}
	message := &queue.Message{
		ID:      uuid.NewString(),
		Type:    b.GetName(),
		Payload: payload,
		Queue:   b.OnQueue(),
		Timeout: b.Timeout(),
		Retries: b.Retries(),
	}
	if delay == 0 {
		err = q.Push(ctx, message)
	} else {
		at := time.Now().Add(delay)
		err = q.Later(ctx, message, at)
	}
	return
}
