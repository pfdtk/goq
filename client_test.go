package goq

import (
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/task"
	"go.uber.org/zap"
	"testing"
)

func TestClient(t *testing.T) {
	z, err := zap.NewDevelopment()
	if err != nil {
		t.Error(err)
		return
	}
	log := z.Sugar()
	c := NewClient(&ClientConfig{Logger: log})
	// connect
	conn, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 1,
	})
	c.AddRedisConnect("test", conn)

	err = NewTestTask().Dispatch(
		[]byte("test payload"),
		task.WithDelay(10),
		task.WithPayloadUnique([]byte("test payload"), 10),
	)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDispatch(t *testing.T) {
	z, err := zap.NewDevelopment()
	if err != nil {
		t.Error(err)
		return
	}
	log := z.Sugar()
	c := NewClient(&ClientConfig{Logger: log})
	// connect
	conn, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 2,
	})
	c.AddRedisConnect("test", conn)
	err = c.Dispatch(NewTestTask(), []byte("test"), task.WithUnique("testUuid", 30))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestChain_Dispatch(t *testing.T) {
	// connect
	conn, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 1,
	})
	connect.AddRedisConnect("test", conn)

	t1 := NewTestTask().SetPayload([]byte("test chain 1"))
	t2 := NewTestTask().SetPayload([]byte("test chain 2"))
	t3 := NewTestTask2().SetPayload([]byte("test chain 3"))
	c := task.NewChain(t1, t2, t3)
	_ = c.Dispatch()
}
