package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq/common/message"
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

func (s *Queue) Push(ctx context.Context, message *message.Message) error {
	return nil
}

func (s *Queue) Later(ctx context.Context, message *message.Message, at time.Time) error {
	return nil
}

func (s *Queue) Pop(ctx context.Context, queue string) (*message.Message, error) {
	return nil, nil
}

func (s *Queue) Release(ctx context.Context, queue string, message string, at time.Time) error {
	return nil
}
