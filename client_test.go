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
	c.AddRedisConnect(test.Conn, conn)

	err := test.NewDemoTask().Dispatch(
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
	c.AddRedisConnect(test.Conn, conn)
	err := c.Dispatch(test.NewDemoTask(), []byte("test"))
	if err != nil {
		t.Error(err)
	}
}

func TestChain_Dispatch(t *testing.T) {
	// connect
	conn := test.GetRedis()
	connect.AddRedisConnect(test.Conn, conn)
	t1 := test.NewDemoTask().Message([]byte("test chain 1"))
	t2 := test.NewDemoTask().Message([]byte("test chain 2"))
	c := task.NewChain(t1, t2)
	err := c.Dispatch()
	if err != nil {
		t.Error(err)
	}
}
