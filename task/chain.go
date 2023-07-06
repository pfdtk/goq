package task

import (
	"context"
	"errors"
	"github.com/google/uuid"
	qm "github.com/pfdtk/goq/internal/queue"
	"github.com/pfdtk/goq/queue"
	"time"
)

type Chain struct {
	task []Task
}

func NewChain(t ...Task) *Chain {
	return &Chain{task: t}
}

func (c *Chain) Dispatch(opt ...DispatchOptFunc) (err error) {
	return c.DispatchContext(context.Background(), opt...)
}

func (c *Chain) DispatchContext(ctx context.Context, optsFun ...DispatchOptFunc) (err error) {
	if len(c.task) < 1 {
		return errors.New("task < 1")
	}
	first := c.convertMessage(c.task[0])
	next := c.convertMessages(c.task[1:]...)
	first.Chain = next

	q := qm.GetQueue(first.OnConnect, first.QueueType)
	if q == nil {
		return errors.New("fail to get queue")
	}

	opt := &DispatchOpt{}
	for i := range optsFun {
		optsFun[i](opt)
	}
	if opt.Delay == 0 {
		err = q.Push(ctx, first)
	} else {
		sec := time.Duration(opt.Delay) * time.Second
		at := time.Now().Add(sec)
		err = q.Later(ctx, first, at)
	}
	return err
}

func (c *Chain) convertMessage(t Task) *queue.Message {
	message := &queue.Message{
		ID:         uuid.NewString(),
		Type:       t.GetName(),
		Payload:    t.Payload(),
		Queue:      t.OnQueue(),
		Timeout:    t.Timeout(),
		Retries:    t.Retries(),
		OnConnect:  t.OnConnect(),
		QueueType:  t.QueueType(),
		DispatchAt: time.Now().Unix(),
	}
	return message
}

func (c *Chain) convertMessages(tasks ...Task) []*queue.Message {
	var ms []*queue.Message
	for i := range tasks {
		m := c.convertMessage(tasks[i])
		ms = append(ms, m)
	}
	return ms
}
