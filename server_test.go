package goq

import (
	"context"
	"github.com/pfdtk/goq/common/cst"
	"github.com/pfdtk/goq/internal/common"
	"github.com/pfdtk/goq/internal/connect"
	"go.uber.org/zap"
	"testing"
	"time"
)

type TestTask struct {
}

func (t TestTask) Run(_ context.Context, job *common.Job) error {
	time.Sleep(10 * time.Second)
	println(job)
	return nil
}

func (t TestTask) QueueType() cst.Type {
	return cst.Redis
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

func (t TestTask) GetStatus() uint32 {
	return cst.Active
}

func (t TestTask) CanRun() bool {
	return true
}

func (t TestTask) Backoff() time.Time {
	return time.Now().Add(1 * time.Second)
}

func (t TestTask) Priority() int {
	return cst.P0
}

func (t TestTask) Retries() int {
	return 0
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
