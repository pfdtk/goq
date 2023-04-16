package goq

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/internal/event"
	"github.com/pfdtk/goq/logger"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	logger       logger.Logger
	eventManager *event.Manager
}

func NewClient(config *ClientConfig) *Client {
	return &Client{
		logger:       config.logger,
		eventManager: event.NewManager(),
	}
}

func (c *Client) AddConnect(name string, conn any) {
	connect.AddConnect(name, conn)
}

func (c *Client) AddRedisConnect(name string, conn *redis.Client) {
	connect.AddRedisConnect(name, conn)
}

func (c *Client) AddSqsConnect(name string, conn *sqs.Client) {
	connect.AddSqsConnect(name, conn)
}
