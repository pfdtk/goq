package task

import (
	"context"
	"errors"
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

func (b *BaseTask) Dispatch(ctx context.Context, payload []byte) error {
	q := qm.GetQueue(b.OnConnect(), b.QueueType())
	if q == nil {
		return errors.New("fail to get queue")
	}
	err := q.Push(ctx, &queue.Message{
		ID:      "uuid-13", // todo create random uuid
		Type:    b.GetName(),
		Payload: payload,
		Queue:   b.OnQueue(),
		Timeout: b.Timeout(),
		Retries: b.Retries(),
	})
	return err
}

func (b *BaseTask) DispatchDelay(ctx context.Context, payload []byte, delay time.Duration) error {
	q := qm.GetQueue(b.OnConnect(), b.QueueType())
	if q == nil {
		return errors.New("fail to get queue")
	}
	err := q.Later(ctx, &queue.Message{
		ID:      "uuid-13", // todo create random uuid
		Type:    b.GetName(),
		Payload: payload,
		Queue:   b.OnQueue(),
		Timeout: b.Timeout(),
		Retries: b.Retries(),
	}, time.Now().Add(delay))
	return err
}
