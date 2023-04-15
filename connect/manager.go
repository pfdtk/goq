package connect

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/redis/go-redis/v9"
	"sync"
)

type ConnType string

const (
	Redis ConnType = "redis"
	Sqs   ConnType = "sqs"
)

var conn sync.Map

func AddConnect(name string, c any) {
	conn.Store(name, c)
}

func AddRedisConnect(name string, r *redis.Client) {
	conn.Store(name, r)
}

func AddSqsConnect(name string, s *sqs.Client) {
	conn.Store(name, s)
}

func Get(name string) any {
	c, ok := conn.Load(name)
	if !ok {
		return nil
	}
	return c
}

func GetRedis(name string) *redis.Client {
	c, ok := conn.Load(name)
	if !ok {
		return nil
	}
	return c.(*redis.Client)
}

func GetSqs(name string) *sqs.Client {
	c, ok := conn.Load(name)
	if !ok {
		return nil
	}
	return c.(*sqs.Client)
}
