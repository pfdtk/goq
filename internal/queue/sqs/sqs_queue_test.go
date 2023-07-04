package sqs

import (
	"context"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/queue"
	"testing"
	"time"
)

func getQueue() *Queue {
	client, err := connect.NewSqsConn(&connect.SqsConf{
		Region: "ap-southeast-1",
		Prefix: "https://sqs.ap-southeast-1.amazonaws.com/123",
	})
	if err != nil {
		panic(err)
	}
	q := NewSqsQueue(client)
	return q
}

func TestQueue_Size(t *testing.T) {
	q := getQueue()
	s, err := q.Size(context.Background(), "name")
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}

func TestQueue_Push(t *testing.T) {
	q := getQueue()
	err := q.Push(context.Background(), &queue.Message{
		Type:    "testsqs",
		Payload: []byte("payload"),
		ID:      "uuid-13",
		Queue:   "name",
		Timeout: 10,
		Retries: 2,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestQueue_Later(t *testing.T) {
	q := getQueue()
	err := q.Later(context.Background(), &queue.Message{
		Type:    "testsqs",
		Payload: []byte("payload"),
		ID:      "uuid-13",
		Queue:   "name",
		Timeout: 10,
		Retries: 2,
	}, time.Now().Add(30*time.Second))
	if err != nil {
		t.Error(err)
	}
}

func TestQueue_Pop(t *testing.T) {
	q := getQueue()
	msg, err := q.Pop(context.Background(), "name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(msg)
}
