package goq

import (
	"context"
	"github.com/pfdtk/goq/common/job"
	"github.com/pfdtk/goq/common/task"
	"github.com/pfdtk/goq/internal/connect"
	"go.uber.org/zap"
	"testing"
	"time"
)

type TestTask struct {
	task.BaseTask
}

func (t TestTask) Run(_ context.Context, j *job.Job) (any, error) {
	log := zap.S()
	time.Sleep(10 * time.Second)
	log.Info("touch test task")
	log.Info(j)
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
		MaxWorker: 1,
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
