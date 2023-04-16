package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq/queue"
	"time"
)

type Queue struct {
	client *sqs.Client
}

func NewSqsQueue(client *sqs.Client) *Queue {
	return &Queue{client: client}
}

func (s *Queue) Size(ctx context.Context, queue string) (int64, error) {
	return 0, nil
}

func (s *Queue) Push(ctx context.Context, message *queue.Message) error {
	return nil
}

func (s *Queue) Later(ctx context.Context, message *queue.Message, at time.Time) error {
	return nil
}

func (s *Queue) Pop(ctx context.Context, queue string) (*queue.Message, error) {
	return nil, nil
}

func (s *Queue) Release(
	ctx context.Context,
	queue string,
	message *queue.Message,
	at time.Time) error {

	return nil
}

func (s *Queue) Delete(ctx context.Context, queue string, message *queue.Message) error {
	return nil
}
