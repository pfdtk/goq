package iface

import (
	"context"
	"github.com/pfdtk/goq/common"
	"time"
)

type Queue interface {
	Size(ctx context.Context, queue string) (int64, error)
	Push(ctx context.Context, message *common.Message) error
	Later(ctx context.Context, message *common.Message, at time.Time) error
	Pop(ctx context.Context, queue string) (*common.Message, error)
	Release(ctx context.Context, message *common.Message, at time.Time) error
}
