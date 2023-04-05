package queue

import (
	"context"
	"github.com/pfdtk/goq"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{client: client}
}

func (r RedisQueue) Size(ctx context.Context, queue string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisQueue) Push(ctx context.Context, message *goq.Message) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisQueue) Later(ctx context.Context, message *goq.Message, at time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisQueue) Pop(ctx context.Context, queue string) (*goq.Message, error) {
	//TODO implement me
	panic("implement me")
}
