package task

import (
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/test"
	"testing"
)

func TestChain_Dispatch(t *testing.T) {
	connect.AddRedisConnect("test", test.GetRedis())

	t1 := NewDemoTask().Message([]byte("test chain 1"))
	t2 := NewDemoTask().Message([]byte("test chain 2"))
	t3 := NewDemoTask().Message([]byte("test chain 3"))
	err := NewChain(t1, t2, t3).Dispatch()
	if err != nil {
		t.Error(err)
	}
}
