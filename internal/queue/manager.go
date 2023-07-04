package queue

import (
	"github.com/pfdtk/goq/connect"
	rdq "github.com/pfdtk/goq/internal/queue/redis"
	sqsq "github.com/pfdtk/goq/internal/queue/sqs"
	"github.com/pfdtk/goq/queue"
	"github.com/redis/go-redis/v9"
)

func GetQueue(conn string, qt queue.Type) queue.Queue {
	c := connect.Get(conn)
	if c == nil {
		return nil
	}
	switch qt {
	case queue.Redis:
		return rdq.NewRedisQueue(c.(*redis.Client))
	case queue.Sqs:
		return sqsq.NewSqsQueue(c.(*connect.SqsClient))
	}
	return nil
}
