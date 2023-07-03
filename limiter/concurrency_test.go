package limiter

import (
	"github.com/pfdtk/goq/connect"
	"sync"
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
	limit := NewConcurrency(c, "test-lock", 10, 200)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {

			key, _ := limit.Acquire(id)
			println(key + "x")
			wg.Done()
		}()
	}
	wg.Wait()
	//limit.Release(key, id)
}
