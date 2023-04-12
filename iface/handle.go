package iface

import "context"

type ErrorHandler interface {
	Handle(ctx context.Context, err error)
}

type ErrorJobHandler interface {
	Handle(ctx context.Context, task Task)
}
