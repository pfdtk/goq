package goq

import (
	"context"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/queue"
	"github.com/pfdtk/goq/task"
	"go.uber.org/zap"
	"testing"
	"time"
)

type TestTask struct {
	task.BaseTask
	logger logger.Logger
}

func NewTask(logger logger.Logger) *TestTask {
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
		BaseTask: task.BaseTask{Option: option},
		logger:   logger,
	}
}

func (t *TestTask) Run(_ context.Context, j *task.Job) (any, error) {
	time.Sleep(2 * time.Second)
	t.logger.Info(string(j.RawMessage().Payload))
	return "test", nil
}

func (t *TestTask) Beforeware() []task.Middleware {
	var mds = []task.Middleware{func(p any, next func(p any) any) any {
		t.logger.Info("before middleware 1, can run: true")
		return next(p)
	}, func(p any, next func(p any) any) any {
		t.logger.Info("before middleware 2, can run: true")
		return next(p)
	}}
	return mds
}

func (t *TestTask) Processware() []task.Middleware {
	var mds = []task.Middleware{func(p any, next func(p any) any) any {
		t.logger.Info("process middleware 1")
		return next(p)
	}, func(p any, next func(p any) any) any {
		t.logger.Info("process middleware 2")
		r := next(p)
		t.logger.Info("process middleware 2: after")
		return r
	}}
	return mds
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
	server.RegisterTask(NewTask(log))
	server.RegisterCronTask("* * * * *", NewTask(log))
	err = server.Start(context.Background())
	if err != nil {
		t.Error(err)
	}
}
