package goq

import (
	"context"
	"time"
)

type Queue interface {
	Size(ctx context.Context, queue string) (int64, error)
	Push(ctx context.Context, message *Message) error
	Later(ctx context.Context, message *Message, at time.Time) error
	Pop(ctx context.Context, queue string) (*Message, error)
}
