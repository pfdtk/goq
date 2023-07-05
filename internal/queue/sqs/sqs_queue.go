package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/queue"
	"strconv"
	"time"
)

var EmptyMessageError = errors.New("empty message")

type Queue struct {
	client *sqs.Client
	conf   *connect.SqsConf
}

func NewSqsQueue(c *connect.SqsClient) *Queue {
	return &Queue{client: c.SqsClient(), conf: c.SqsConf()}
}

func (s *Queue) getQueueUrl(queue string) *string {
	url := s.conf.Prefix + "/" + queue
	return &url
}

func (s *Queue) Size(ctx context.Context, queue string) (int64, error) {
	p := &sqs.GetQueueAttributesInput{
		QueueUrl:       s.getQueueUrl(queue),
		AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameApproximateNumberOfMessages},
	}
	res, err := s.client.GetQueueAttributes(ctx, p)
	if err != nil {
		return 0, err
	}
	count, err := strconv.ParseInt(res.Attributes["ApproximateNumberOfMessages"], 10, 64)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Queue) Push(ctx context.Context, message *queue.Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	body := string(bytes)
	p := &sqs.SendMessageInput{
		QueueUrl:    s.getQueueUrl(message.Queue),
		MessageBody: &body,
	}
	_, err = s.client.SendMessage(ctx, p)
	return err
}

func (s *Queue) Later(ctx context.Context, message *queue.Message, at time.Time) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	body := string(bytes)
	delay := at.Unix() - time.Now().Unix()
	p := &sqs.SendMessageInput{
		QueueUrl:     s.getQueueUrl(message.Queue),
		MessageBody:  &body,
		DelaySeconds: int32(delay),
	}
	_, err = s.client.SendMessage(ctx, p)
	return err
}

func (s *Queue) Pop(ctx context.Context, qname string) (*queue.Message, error) {
	attrNames := []types.QueueAttributeName{types.QueueAttributeName("ApproximateReceiveCount")}
	p := &sqs.ReceiveMessageInput{
		QueueUrl:       s.getQueueUrl(qname),
		AttributeNames: attrNames,
	}
	res, err := s.client.ReceiveMessage(ctx, p)
	if err != nil {
		return nil, err
	}
	if len(res.Messages) == 0 {
		return nil, EmptyMessageError
	}
	message := res.Messages[0]
	attempts, err := strconv.ParseUint(message.Attributes["ApproximateReceiveCount"], 10, 64)
	if err != nil {
		return nil, err
	}
	msg := queue.Message{}
	err = json.Unmarshal([]byte(*message.Body), &msg)
	if err != nil {
		return nil, err
	}
	msg.Attempts = uint(attempts)
	msg.Reserved = *message.Body
	msg.ReceiptHandle = *message.ReceiptHandle
	return &msg, nil
}

func (s *Queue) Release(
	ctx context.Context,
	queue string,
	message *queue.Message,
	at time.Time) error {

	delay := at.Unix() - time.Now().Unix()
	p := &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          s.getQueueUrl(queue),
		ReceiptHandle:     &message.ReceiptHandle,
		VisibilityTimeout: int32(delay),
	}
	_, err := s.client.ChangeMessageVisibility(ctx, p)
	return err
}

func (s *Queue) Delete(ctx context.Context, queue string, message *queue.Message) error {
	p := &sqs.DeleteMessageInput{
		QueueUrl:      s.getQueueUrl(queue),
		ReceiptHandle: &message.ReceiptHandle,
	}
	_, err := s.client.DeleteMessage(ctx, p)
	return err
}
