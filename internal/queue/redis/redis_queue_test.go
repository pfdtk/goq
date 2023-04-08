package redis

import (
	"context"
	"github.com/pfdtk/goq/common"
	"github.com/pfdtk/goq/internal/connect"
	"testing"
)

func TestRedisQueue_Push(t *testing.T) {
	c, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 1,
	})
	defer func() {
		_ = c.Close()
	}()
	q := NewRedisQueue(c)
	err := q.Push(context.Background(), &common.Message{
		Type:    "test",
		Payload: []byte("payload"),
		ID:      "uuid-13",
		Queue:   "default",
		Timeout: 10,
	})
	if err != nil {
		t.Error(err)
	}
}
