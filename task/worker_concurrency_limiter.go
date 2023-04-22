package task

import (
	"github.com/google/uuid"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/limiter"
	"github.com/pfdtk/goq/logger"
)

var namePrefix = "goq-task-max-worker-control"

//var MaxWorkerError = errors.New("max worker error")

// NewMaxWorkerControl you must add a redis connect to server first
func NewMaxWorkerControl(conn string, name string, maxWorker int, timeout int) Middleware {
	return func(p any, next func(p any) any) any {
		pp, ok := p.(*PopPassable)
		// ensure from Beforeware
		if !ok {
			return next(p)
		}
		logger.GetLogger().Infof("try to get max worker lock, task=%s", name)
		redis := connect.GetRedis(conn)
		lockname := namePrefix + ":" + name
		l := limiter.NewConcurrency(redis, lockname, maxWorker, timeout)
		id := uuid.NewString()
		token, err := l.Acquire(id)
		if err != nil {
			logger.GetLogger().Infof("reach max worker error, task=%s", name)
			return false
		}
		logger.GetLogger().Infof("got max worker lock, task=%s", name)
		prev := pp.Callback
		pp.Callback = func() {
			logger.GetLogger().Infof("release max worker lock, task=%s", name)
			l.Release(token, id)
			if prev != nil {
				prev()
			}
		}
		return next(p)
	}
}
