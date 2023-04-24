package task

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/logger"
)

func NewLeakyBucketLimiter(
	conn string,
	name string,
	limit redis_rate.Limit) Middleware {
	return func(p any, next func(p any) any) any {
		logger.GetLogger().Debugf("try to get leaky bucket limiter lock, task=%s", name)
		redis := connect.GetRedis(conn)
		limiter := redis_rate.NewLimiter(redis)
		lockName := "goq-task-leaky-bucket-limiter:" + name
		res, err := limiter.Allow(context.Background(), lockName, limit)

		if err == nil && res.Allowed != 0 {
			logger.GetLogger().Debugf("got leaky bucket limiter lock, task=%s", name)
			return next(p)
		}

		logger.GetLogger().Debugf("fail to get leaky bucket limiter lock, task=%s", name)

		switch p.(type) {
		case *RunPassable:
			rp := p.(*RunPassable)
			if err = rp.job.Release(context.Background(), rp.task.Backoff()); err != nil {
				logger.GetLogger().Errorf("fail to release job when not get the leaky bucket limiter lock, task=%s", name)
			}
		}

		return false
	}
}
