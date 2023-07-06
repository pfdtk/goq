package redis

import (
	"github.com/pfdtk/goq/backend/state"
	"github.com/pfdtk/goq/queue"
	"github.com/redis/go-redis/v9"
)

type Backend struct {
	client *redis.Client
}

func NewRedisBackend(conn *redis.Client) *Backend {
	return &Backend{client: conn}
}

func (b Backend) Pending(message *queue.Message) error {
	return nil
}

func (b Backend) Started(message *queue.Message) error {
	return nil
}

func (b Backend) Success(message *queue.Message) error {
	return nil
}

func (b Backend) Failure(message *queue.Message) error {
	return nil
}

func (b Backend) State(messageId string) (state.State, error) {
	return state.Pending, nil
}
