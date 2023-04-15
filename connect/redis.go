package connect

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisConf struct {
	Addr     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

// NewRedisConn new a redis client
func NewRedisConn(conf *RedisConf) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Addr + ":" + conf.Port,
		Password: conf.Password,
		DB:       conf.DB,
		PoolSize: conf.PoolSize,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
