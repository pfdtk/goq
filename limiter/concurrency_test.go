package limiter

import (
	"github.com/pfdtk/goq/connect"
	"testing"
)

func TestNewConcurrency(t *testing.T) {
	c, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 1,
	})
	id := "id-of-lock"
	limit := NewConcurrency(c, "test-lock", 10, 20)
	key, err := limit.Acquire(id)
	if err != nil {
		t.Error(err)
		return
	}
	println(key)
	limit.Release(key, id)
}
