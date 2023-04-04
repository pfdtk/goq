package goq

import (
	"context"
	"github.com/pfdtk/goq/internal/connect"
)

type Task interface {
	Run(context.Context, *Job) error
	OnConnect() connect.ConnType
	OnQueue() string
	GetName() string
}
