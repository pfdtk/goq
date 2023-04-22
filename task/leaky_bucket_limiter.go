package task

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/logger"
)

func NewLeakyBucketLimiter(conn string, name string, limit redis_rate.Limit) Middleware {
	return func(p any, next func(p any) any) any {
		_, ok := p.(*PopPassable)
		// ensure from Beforeware
		if !ok {
			return next(p)
		}
		logger.GetLogger().Infof("try to get leaky bucket limiter lock, task=%s", name)
		redis := connect.GetRedis(conn)
		limiter := redis_rate.NewLimiter(redis)
		lockName := namePrefix + ":" + name
		res, err := limiter.Allow(context.Background(), lockName, limit)
		if err != nil || res.Allowed == 0 {
			logger.GetLogger().Infof("fail to get leaky bucket limiter lock, task=%s", name)
			return false
		}
		logger.GetLogger().Infof("got leaky bucket limiter lock, task=%s", name)
		return next(p)
	}
}
