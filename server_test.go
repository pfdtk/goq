package goq

import (
	"context"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/queue"
	"github.com/pfdtk/goq/task"
	"go.uber.org/zap"
	"testing"
	"time"
)

type TestTask struct {
	task.BaseTask
}

func NewTask() *TestTask {
	option := &task.Option{
		Name:      "test",
		OnConnect: "test",
		QueueType: queue.Redis,
		OnQueue:   "default",
		Status:    task.Active,
		Priority:  0,
		Retries:   0,
		Timeout:   11,
	}
	return &TestTask{
		task.BaseTask{Option: option},
	}
}

func (t *TestTask) Run(_ context.Context, _ *task.Job) (any, error) {
	time.Sleep(2 * time.Second)
	return "test", nil
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
	server.AddRedisConnect("test", conn)
	server.RegisterTask(NewTask())
	server.RegisterCronTask("* * * * *", NewTask())
	err = server.Start(context.Background())
	if err != nil {
		t.Error(err)
	}
}
