package redis

import (
	"github.com/pfdtk/goq/connect"
	"testing"
)

func TestNewRedisBackend(t *testing.T) {
	conn, _ := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 2,
	})
	NewRedisBackend(conn)
}
