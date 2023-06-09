package task

import (
	"context"
	"github.com/google/uuid"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/limiter"
	"github.com/pfdtk/goq/logger"
)

// NewMaxWorkerLimiter you must add a redis connect to server first
// This middleware can be used in beforeware & processware
func NewMaxWorkerLimiter(
	conn string,
	name string,
	maxWorker int,
	timeout int) Middleware {
	return func(p any, next func(p any) any) any {
		redis := connect.GetRedis(conn)
		lockName := "goq-task-max-worker-limiter:" + name
		l := limiter.NewConcurrency(redis, lockName, maxWorker, timeout)
		id := uuid.NewString()
		token, err := l.Acquire(id)

		if err != nil {
			switch p.(type) {
			case *RunPassable:
				rp := p.(*RunPassable)
				if err = rp.job.Release(context.Background(), rp.task.Backoff()); err != nil {
					logger.GetLogger().Errorf("fail to release job when not get the max worker limiter lock, task=%s", name)
				}
			}
			return false
		}

		switch p.(type) {
		case *PopPassable:
			pp := p.(*PopPassable)
			prev := pp.GetCallback()
			pp.SetCallback(func() {
				l.Release(token, id)
				if prev != nil {
					prev()
				}
			})
		}

		return next(p)
	}
}
