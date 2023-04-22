package task

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/pfdtk/goq/event"
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

func (b *BaseTask) UniqueId() string {
	return b.Option.UniqueId
}

func (b *BaseTask) Beforeware() []Middleware {
	return nil
}

func (b *BaseTask) Processware() []Middleware {
	return nil
}

func (b *BaseTask) Dispatch(payload []byte, opt ...DispatchOptFunc) (err error) {
	return b.DispatchContext(context.Background(), payload, opt...)
}

func (b *BaseTask) DispatchContext(
	ctx context.Context,
	payload []byte,
	optsFun ...DispatchOptFunc) (err error) {

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
	opt := &DispatchOpt{}
	for i := range optsFun {
		optsFun[i](opt)
	}
	if opt.Delay == 0 {
		err = q.Push(ctx, message)
	} else {
		sec := time.Duration(opt.Delay) * time.Second
		at := time.Now().Add(sec)
		err = q.Later(ctx, message, at)
	}
	if err == nil {
		event.Dispatch(NewJobAddEvent(b, message))
	}
	return
}
