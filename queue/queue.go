package queue

import (
	"context"
	"time"
)

type Type string

const (
	Redis Type = "redis"
	Sqs   Type = "sqs"
)

type Queue interface {
	// Size get size of queue
	Size(ctx context.Context, queue string) (int64, error)
	// Push a message to queue
	Push(ctx context.Context, message *Message) error
	// Later push a delay message to queue
	Later(ctx context.Context, message *Message, at time.Time) error
	// Pop a message from queue
	Pop(ctx context.Context, queue string) (*Message, error)
	// Release a message back to queue for retry reason
	Release(ctx context.Context, queue string, message *Message, at time.Time) error
	// Delete ack message from queue
	Delete(ctx context.Context, queue string, message *Message) error
}
