package redis

import (
	"context"
	"encoding/json"
	"github.com/pfdtk/goq/common"
	"github.com/redis/go-redis/v9"
	"time"
)

type Queue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *Queue {
	return &Queue{client: client}
}

func (r Queue) Size(ctx context.Context, queue string) (int64, error) {
	size, err := r.client.LLen(ctx, queue).Result()
	return size, err
}

func (r Queue) Push(ctx context.Context, message *common.Message) error {
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

func (r Queue) Later(ctx context.Context, message *common.Message, at time.Time) error {
	return nil
}

// Pop TODO get and backup msg to a sort set, then ack will delete the msg from sort set
func (r Queue) Pop(ctx context.Context, queue string) (*common.Message, error) {
	val, err := r.client.LPop(ctx, queue).Result()
	if err != nil {
		return nil, err
	}
	msg := common.Message{}
	err = json.Unmarshal([]byte(val), &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
