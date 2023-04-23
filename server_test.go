package goq

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
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
	time.Sleep(8 * time.Second)
	t.logger.Info(string(j.RawMessage().Payload))
	return "test", nil
}

func (t *TestTask) Beforeware() []task.Middleware {
	return nil
}

func (t *TestTask) Processware() []task.Middleware {
	return []task.Middleware{
		//task.NewMaxWorkerLimiter("test", t.GetName(), 1, 10),
		task.NewLeakyBucketLimiter("test", t.GetName(), redis_rate.PerMinute(10)),
	}
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
	server.AddRedisConnect("test", conn)
	server.RegisterTask(NewTask(log))
	//server.RegisterCronTask("* * * * *", NewTask(log))
	err = server.Start(context.Background())
	if err != nil {
		t.Error(err)
	}
}
