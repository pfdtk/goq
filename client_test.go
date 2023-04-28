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
	c := NewClient(&ClientConfig{logger: log})
	// connect
	conn, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 1,
	})
	c.AddRedisConnect("test", conn)

	err = NewTask(log).Dispatch([]byte("test payload"), task.WithDelay(10), task.WithUnique("testUuid", 10))
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
	c := NewClient(&ClientConfig{logger: log})
	// connect
	conn, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 2,
	})
	c.AddRedisConnect("test", conn)
	err = c.Dispatch(NewTask(log), []byte("test"), task.WithUnique("testUuid", 30))
	if err != nil {
		t.Error(err)
		return
	}
}
