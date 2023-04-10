package goq

import (
	"context"
	"github.com/pfdtk/goq/common"
	"github.com/pfdtk/goq/internal/connect"
	"go.uber.org/zap"
	"testing"
	"time"
)

type TestTask struct {
	common.BaseTask
}

func (t TestTask) Run(_ context.Context, job *common.Job) (any, error) {
	log := zap.S()
	time.Sleep(10 * time.Second)
	log.Info("touch test task")
	return nil, nil
}

func (t TestTask) OnConnect() string {
	return "test"
}

func (t TestTask) OnQueue() string {
	return "default"
}

func (t TestTask) GetName() string {
	return "test"
}

func TestServer_Start(t *testing.T) {
	z, err := zap.NewDevelopment()
	if err != nil {
		t.Error(err)
		return
	}
	log := z.Sugar()
	log.Info("xx")
	server := NewServer(&ServerConfig{
		MaxWorker: 2,
		logger:    log,
	})
	// connect
	conn, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 1,
	})
	server.AddConnect("test", conn)
	server.RegisterTask(&TestTask{})
	err = server.Start(context.Background())
	if err != nil {
		t.Error(err)
	}
}
