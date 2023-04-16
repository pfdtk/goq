package goq

import (
	"context"
	"github.com/pfdtk/goq/connect"
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

	err = NewTask().Dispatch(context.Background(), []byte("test payload"))
	if err != nil {
		t.Error(err)
		return
	}
}
