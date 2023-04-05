package goq

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq/internal/queue"
	"github.com/redis/go-redis/v9"
)

type Worker struct {
	app       *App
	stopRun   chan struct{}
	maxWorker chan struct{}
	ctx       context.Context
}

func (w *Worker) StartConsuming() error {
	go w.consume()
	return nil
}

func (w *Worker) consume() {
	for {
		select {
		case <-w.stopRun:
			return
		case w.maxWorker <- struct{}{}:

		}
	}
}

func (w *Worker) getNextJob() *Job {
	var job *Job = nil
	// TODO sort tasks
	w.app.task.Range(func(key, value any) bool {
		t := value.(Task)
		queueType := t.QueueType()
		switch queueType {
		case queue.Redis:
			job = w.getRedisJob(t)
			if job != nil {
				return false
			}
		case queue.Sqs:
			job = w.getSqsJob(t)
			if job != nil {
				return false
			}
		}
		return true
	})
	return job
}

func (w *Worker) getSqsJob(task Task) *Job {
	c, ok := w.app.conn.Load(task.OnConnect())
	if !ok {
		return nil
	}
	sqsClient := queue.NewSqsQueue(c.(*sqs.Client))
	message, err := sqsClient.Pop(w.ctx, task.OnQueue())
	if err != nil {
		return nil
	}
	return &Job{
		id:      message.ID,
		name:    message.Type,
		queue:   message.Queue,
		payload: message.Payload,
	}
}

func (w *Worker) getRedisJob(task Task) *Job {
	c, ok := w.app.conn.Load(task.OnConnect())
	if !ok {
		return nil
	}
	rdc := queue.NewRedisQueue(c.(*redis.Client))
	message, err := rdc.Pop(w.ctx, task.OnQueue())
	if err != nil {
		return nil
	}
	return &Job{
		id:      message.ID,
		name:    message.Type,
		queue:   message.Queue,
		payload: message.Payload,
	}
}
