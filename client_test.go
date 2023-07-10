package goq

import (
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/task"
	"github.com/pfdtk/goq/test"
	"testing"
)

func TestClient(t *testing.T) {
	c := NewClient(&ClientConfig{})
	// connect
	conn := test.GetRedis()
	c.AddRedisConnect(task.Conn, conn)

	err := task.NewDemoTask().Dispatch(
		[]byte("test payload"),
		task.WithDelay(10),
		task.WithPayloadUnique([]byte("test payload"), 10),
	)
	if err != nil {
		t.Error(err)
	}
}

func TestDispatch(t *testing.T) {
	c := NewClient(&ClientConfig{})
	// connect
	conn := test.GetRedis()
	c.AddRedisConnect(task.Conn, conn)
	err := c.Dispatch(task.NewDemoTask(), []byte("test"))
	if err != nil {
		t.Error(err)
	}
}

func TestChain_Dispatch(t *testing.T) {
	// connect
	conn := test.GetRedis()
	connect.AddRedisConnect(task.Conn, conn)
	t1 := task.NewDemoTask().Message([]byte("test chain 1"))
	t2 := task.NewDemoTask().Message([]byte("test chain 2"))
	c := task.NewChain(t1, t2)
	err := c.Dispatch()
	if err != nil {
		t.Error(err)
	}
}
