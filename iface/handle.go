package iface

import "context"

type ErrorJobHandler interface {
	Handle(ctx context.Context, task Task)
}
