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

func NewTestTask2() *TestTask {
	option := &task.Option{
		Name:      "test2",
		OnConnect: "test",
		QueueType: queue.Redis,
		OnQueue:   "default2",
		Status:    task.Active,
		Priority:  0,
		Retries:   0,
		Timeout:   500,
	}
	return &TestTask{
		BaseTask: task.BaseTask{Option: option},
		logger:   logger.GetLogger(),
	}
}

func NewTestTask() *TestTask {
	option := &task.Option{
		Name:      "test",
		OnConnect: "test",
		QueueType: queue.Redis,
		OnQueue:   "default",
		Status:    task.Active,
		Priority:  0,
		Retries:   0,
		Timeout:   500,
	}
	return &TestTask{
		BaseTask: task.BaseTask{Option: option},
		logger:   logger.GetLogger(),
	}
}

func (t *TestTask) Run(_ context.Context, j *task.Job) (any, error) {
	time.Sleep(8 * time.Second)
	t.logger.Info(string(j.RawMessage().Payload))
	return "test", nil
}

func (t *TestTask) Beforeware() []task.Middleware {
	return []task.Middleware{
		//task.NewMaxWorkerLimiter("test", t.GetName(), 1, 10),
		//task.NewLeakyBucketLimiter("test", t.GetName(), redis_rate.PerMinute(1)),
	}
}

func (t *TestTask) Processware() []task.Middleware {
	return nil
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
		Logger:    log,
	})
	// connect
	conn, err := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}
	server.AddRedisConnect("test", conn)
	server.RegisterTask(NewTestTask())
	server.RegisterTask(NewTestTask2())
	//server.RegisterCronTask("* * * * *", NewTestTask(log))
	server.MustStart(context.Background())
}
