package redis

import (
	"context"
	"errors"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/queue"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
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
	err := q.Push(context.Background(), &queue.Message{
		Type:    "test",
		Payload: []byte("payload"),
		ID:      "uuid-13",
		Queue:   queueName,
		Timeout: 15,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestQueue_Later(t *testing.T) {
	q := getQueue()
	err := q.Later(context.Background(), &queue.Message{
		Type:    "test",
		Payload: []byte("payload"),
		ID:      "uuid-13",
		Queue:   queueName,
		Timeout: 15,
	}, time.Now().Add(1*time.Hour))
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

func TestQueue_Release(t *testing.T) {
	q := getQueue()

	err := q.Push(context.Background(), &queue.Message{
		Type:    "test",
		Payload: []byte("payload"),
		ID:      "uuid-13",
		Queue:   queueName,
		Timeout: 15,
	})
	if err != nil {
		t.Error(err)
		return
	}

	s, err := q.Pop(context.Background(), queueName)
	if err != nil && !errors.Is(err, redis.Nil) {
		t.Error(err)
		return
	}

	at := time.Now().Add(10 * time.Minute)
	err = q.Release(context.Background(), queueName, s.Reserved, at)
	if err != nil && !errors.Is(err, redis.Nil) {
		t.Error(err)
	} else {
		t.Log(s)
	}
}
