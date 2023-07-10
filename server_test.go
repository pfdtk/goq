package goq

import (
	"context"
	"github.com/pfdtk/goq/task"
	"github.com/pfdtk/goq/test"
	"testing"
)

func TestServer_Start(t *testing.T) {
	server := NewServer(&ServerConfig{MaxWorker: 5})
	// connect
	server.AddRedisConnect(task.Conn, test.GetRedis())
	server.RegisterTask(task.NewDemoTask())
	server.RegisterCronTask("* * * * *", task.NewDemoTask())
	server.MustStart(context.Background())
}
