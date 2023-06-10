package task

import (
	"github.com/pfdtk/goq/connect"
	"testing"
)

func TestChain_Dispatch(t *testing.T) {
	// connect
	conn, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 1,
	})
	connect.AddRedisConnect("test", conn)

	t1 := NewTask().SetPayload([]byte("test chain 1"))
	t2 := NewTask().SetPayload([]byte("test chain 2"))
	t3 := NewTask().SetPayload([]byte("test chain 3"))
	c := NewChain(t1, t2, t3)
	_ = c.Dispatch()
}
