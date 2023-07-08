package test

import (
	"github.com/pfdtk/goq/connect"
	"github.com/redis/go-redis/v9"
)

func GetRedis() *redis.Client {
	conn, err := connect.NewRedisConn(&connect.RedisConf{
		Addr:     "127.0.0.1",
		Port:     "6379",
		DB:       1,
		PoolSize: 1,
	})
	if err != nil {
		panic("get redis client error")
	}
	return conn
}
