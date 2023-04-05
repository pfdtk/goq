package queue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq"
	"time"
)

type SqsQueue struct {
	client *sqs.Client
}

func (s SqsQueue) Size(ctx context.Context, queue string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s SqsQueue) Push(ctx context.Context, message *goq.Message) error {
	//TODO implement me
	panic("implement me")
}

func (s SqsQueue) Later(ctx context.Context, message *goq.Message, at time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (s SqsQueue) Pop(ctx context.Context, queue string) (*goq.Message, error) {
	//TODO implement me
	panic("implement me")
}

func NewSqsQueue(client *sqs.Client) *SqsQueue {
	return &SqsQueue{client: client}
}
