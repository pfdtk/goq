package connect

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/redis/go-redis/v9"
)

type ConnType string

const (
	ConnRedis ConnType = "redis"
	ConnSqs   ConnType = "sqs"
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

// SqsConf TODO
type SqsConf struct {
}

// NewSqsConn new an aws sqs client
func NewSqsConn(_ *SqsConf) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := sqs.NewFromConfig(cfg)
	return client, nil
}
