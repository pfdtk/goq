package redis

import (
	"errors"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/queue"
	"testing"
)

func TestNewRedisBackend(t *testing.T) {
	conn, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 2,
	})
	connect.AddRedisConnect("test", conn)
	b := NewRedisBackend("test")
	msg := &queue.Message{ID: "x1", Type: "type"}

	_ = b.Pending(msg)
	ss, _ := b.State("x1")
	t.Log(ss)
	_ = b.Started(msg)
	ss, _ = b.State("x1")
	t.Log(ss)
	_ = b.Success(msg)
	ss, _ = b.State("x1")
	t.Log(ss)
	_ = b.Failure(msg, errors.New("test"))
	ss, _ = b.State("x1")
	t.Log(ss)
}
