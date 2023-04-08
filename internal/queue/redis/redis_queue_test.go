package redis

import (
	"context"
	"errors"
	"github.com/pfdtk/goq/common"
	"github.com/pfdtk/goq/internal/connect"
	"github.com/redis/go-redis/v9"
	"testing"
)

var queueName = "default"

func getQueue() *Queue {
	c, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 1,
	})
	q := NewRedisQueue(c)
	return q
}

func TestRedisQueue_Push(t *testing.T) {
	q := getQueue()
	err := q.Push(context.Background(), &common.Message{
		Type:    "test",
		Payload: []byte("payload"),
		ID:      "uuid-13",
		Queue:   queueName,
		Timeout: 10,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestQueue_Size(t *testing.T) {
	q := getQueue()
	s, err := q.Size(context.Background(), queueName)
	if err != nil {
		t.Error(err)
	}
	t.Logf("size: %d", s)
}

func TestQueue_Pop(t *testing.T) {
	q := getQueue()
	s, err := q.Pop(context.Background(), queueName)
	if err != nil && !errors.Is(err, redis.Nil) {
		t.Error(err)
	} else {
		t.Log(s)
	}
}