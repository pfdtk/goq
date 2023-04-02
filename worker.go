package goq

import "context"

type Task struct {
	typename string
	payload  []byte
}

type Handler interface {
	Run(context.Context, *Task) error
}

type HandlerFunc func(context.Context, *Task) error

func (fn HandlerFunc) Run(ctx context.Context, task *Task) error {
	return fn(ctx, task)
}

type Worker struct {
}

func (srv *Worker) Start(ctx context.Context, handler Handler) error {
	return nil
}
