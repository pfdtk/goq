package task

import (
	"context"
	"github.com/google/uuid"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/limiter"
	"github.com/pfdtk/goq/logger"
)

// NewMaxWorkerLimiter you must add a redis connect to server first
func NewMaxWorkerLimiter(
	conn string,
	name string,
	maxWorker int,
	timeout int) Middleware {
	return func(p any, next func(p any) any) any {
		logger.GetLogger().Debugf("try to get max worker limiter lock, task=%s", name)
		redis := connect.GetRedis(conn)
		lockName := "goq-task-max-worker-limiter:" + name
		l := limiter.NewConcurrency(redis, lockName, maxWorker, timeout)
		id := uuid.NewString()
		token, err := l.Acquire(id)

		if err != nil {
			logger.GetLogger().Debugf("fail to get max worker limiter lock, task=%s", name)
			switch p.(type) {
			case *RunPassable:
				rp := p.(*RunPassable)
				if err = rp.job.Release(context.Background(), rp.task.Backoff()); err != nil {
					logger.GetLogger().Errorf("fail to release job when not get the max worker limiter lock, task=%s", name)
				}
			}
			return false
		}

		logger.GetLogger().Debugf("got max worker limiter lock, task=%s", name)

		switch p.(type) {
		case *PopPassable:
			pp := p.(*PopPassable)
			prev := pp.Callback
			pp.Callback = func() {
				logger.GetLogger().Debugf("release the max worker limiter lock, task=%s", name)
				l.Release(token, id)
				if prev != nil {
					prev()
				}
			}
		}

		return next(p)
	}
}
