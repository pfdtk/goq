package goq

import "context"

type Task interface {
	Run(context.Context, *Job) error
	OnConnect() string
	OnQueue() string
	GetName() string
}
