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
	Size(ctx context.Context, queue string) (int64, error)
	Push(ctx context.Context, message *Message) error
	Later(ctx context.Context, message *Message, at time.Time) error
	Pop(ctx context.Context, queue string) (*Message, error)
	Release(ctx context.Context, queue string, message *Message, at time.Time) error
	Delete(ctx context.Context, queue string, message *Message) error
}
