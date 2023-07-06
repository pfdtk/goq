package redis

import (
	"context"
	"encoding/json"
	"github.com/pfdtk/goq/backend/state"
	"github.com/pfdtk/goq/queue"
	"github.com/redis/go-redis/v9"
	"time"
)

type Backend struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration
}

func NewRedisBackend(conn *redis.Client) *Backend {
	return &Backend{
		client: conn,
		ctx:    context.Background(),
		ttl:    24 * time.Hour,
	}
}

func (b Backend) Pending(message *queue.Message) error {
	s := state.NewTaskState(message, state.Pending)
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return b.client.Set(b.ctx, s.TaskId, bytes, b.ttl).Err()
}

func (b Backend) Started(message *queue.Message) error {
	s := state.NewTaskState(message, state.Started)
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return b.client.Set(b.ctx, s.TaskId, bytes, b.ttl).Err()
}

func (b Backend) Success(message *queue.Message) error {
	s := state.NewTaskState(message, state.Success)
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return b.client.Set(b.ctx, s.TaskId, bytes, b.ttl).Err()
}

func (b Backend) Failure(message *queue.Message, err error) error {
	s := state.NewFailureTaskState(message, err)
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return b.client.Set(b.ctx, s.TaskId, bytes, b.ttl).Err()
}

func (b Backend) State(messageId string) (state.State, error) {
	result, err := b.client.Get(b.ctx, messageId).Result()
	if err != nil {
		return "", err
	}
	st := state.TaskState{}
	err = json.Unmarshal([]byte(result), &st)
	if err != nil {
		return "", err
	}
	return st.State, nil
}
