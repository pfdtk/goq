package task

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/event"
	"github.com/pfdtk/goq/limiter"
)

var namePrefix = "goq-task-max-worker-control"
var MaxWorkerError = errors.New("max worker error")

// NewMaxWorkerControl you must add a redis connect to server first
func NewMaxWorkerControl(conn string, name string, maxWorker int, timeout int) Middleware {
	return func(p any, next func(p any) any) any {
		redis := connect.GetRedis(conn)
		name = namePrefix + ":" + name
		l := limiter.NewConcurrency(redis, name, maxWorker, timeout)
		id := uuid.NewString()
		_, err := l.Acquire(id)
		if err != nil {
			event.Dispatch(NewMaxWorkerErrorEvent(MaxWorkerError))
			return false
		}
		return next(p)
	}
}
