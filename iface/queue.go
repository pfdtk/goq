package iface

import (
	"context"
	"github.com/pfdtk/goq/common/message"
	"time"
)

type Queue interface {
	Size(ctx context.Context, queue string) (int64, error)
	Push(ctx context.Context, message *message.Message) error
	Later(ctx context.Context, message *message.Message, at time.Time) error
	Pop(ctx context.Context, queue string) (*message.Message, error)
	Release(ctx context.Context, queue string, message string, at time.Time) error
}
