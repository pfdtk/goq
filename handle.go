package goq

import "context"

type Handler interface {
	Run(context.Context, *Task) error
}
