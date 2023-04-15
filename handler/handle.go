package handler

import (
	"context"
	"github.com/pfdtk/goq/task"
)

type ErrorHandler interface {
	Handle(ctx context.Context, err error)
}

type ErrorJobHandler interface {
	Handle(ctx context.Context, task task.Task)
}
