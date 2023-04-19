package goq

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/internal/event"
	qm "github.com/pfdtk/goq/internal/queue"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/queue"
	"github.com/pfdtk/goq/task"
	"github.com/redis/go-redis/v9"
	"time"
)

var client *Client

type Client struct {
	logger       logger.Logger
	eventManager *event.Manager
}

func Dispatch(
	task task.Task,
	payload []byte,
	opt ...task.DispatchOptFunc) (err error) {

	return DispatchContext(context.Background(), task, payload, opt...)
}

func DispatchContext(
	ctx context.Context,
	t task.Task,
	payload []byte,
	opts ...task.DispatchOptFunc) (err error) {

	q := qm.GetQueue(t.OnConnect(), t.QueueType())
	if q == nil {
		return errors.New("fail to get queue")
	}
	message := &queue.Message{
		ID:      uuid.NewString(),
		Type:    t.GetName(),
		Payload: payload,
		Queue:   t.OnQueue(),
		Timeout: t.Timeout(),
		Retries: t.Retries(),
	}
	opt := &task.DispatchOpt{}
	for _, fn := range opts {
		fn(opt)
	}
	if opt.Delay == 0 {
		err = q.Push(ctx, message)
	} else {
		sec := time.Duration(opt.Delay) * time.Second
		at := time.Now().Add(sec)
		err = q.Later(ctx, message, at)
	}
	return
}

func NewClient(config *ClientConfig) *Client {
	client = &Client{
		logger: config.logger,
	}
	return client
}

func (c *Client) AddConnect(name string, conn any) {
	connect.AddConnect(name, conn)
}

func (c *Client) AddRedisConnect(name string, conn *redis.Client) {
	connect.AddRedisConnect(name, conn)
}

func (c *Client) AddSqsConnect(name string, conn *sqs.Client) {
	connect.AddSqsConnect(name, conn)
}

func (c *Client) Dispatch(
	task task.Task,
	payload []byte,
	opt ...task.DispatchOptFunc) (err error) {

	return Dispatch(task, payload, opt...)
}

func (c *Client) DispatchContext(
	ctx context.Context,
	task task.Task,
	payload []byte,
	opts ...task.DispatchOptFunc) (err error) {

	return DispatchContext(ctx, task, payload, opts...)
}
