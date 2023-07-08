package goq

import (
	"context"
	"github.com/pfdtk/goq/test"
	"testing"
)

func TestServer_Start(t *testing.T) {
	server := NewServer(&ServerConfig{MaxWorker: 5})
	// connect
	server.AddRedisConnect(test.Conn, test.GetRedis())
	server.RegisterTask(test.NewDemoTask())
	server.RegisterCronTask("* * * * *", test.NewDemoTask())
	server.MustStart(context.Background())
}
