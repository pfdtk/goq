package queue

import (
	"context"
	"encoding/json"
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
	size, err := r.client.LLen(ctx, queue).Result()
	return size, err
}

func (r RedisQueue) Push(ctx context.Context, message *goq.Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = r.client.RPush(ctx, message.Queue, bytes).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisQueue) Later(ctx context.Context, message *goq.Message, at time.Time) error {
	return nil
}

// Pop TODO get and backup msg to a sort set, then ack will delete the msg from sort set
func (r RedisQueue) Pop(ctx context.Context, queue string) (*goq.Message, error) {
	val, err := r.client.LPop(ctx, queue).Result()
	if err != nil {
		return nil, err
	}
	msg := goq.Message{}
	err = json.Unmarshal([]byte(val), &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
