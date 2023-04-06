package goq

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq/internal/queue"
	"github.com/redis/go-redis/v9"
	"time"
)

type Worker struct {
	server     *Server
	stopRun    chan struct{}
	maxWorker  chan struct{}
	jobChannel chan *Job
	ctx        context.Context
}

func (w *Worker) StartConsuming() error {
	// TODO single coroutine, maybe a litter slow
	go w.consume()
	go w.work()
	return nil
}

func (w *Worker) StopConsuming() {

}

// consume only pop message from queue and send to work channel
func (w *Worker) consume() {
	for {
		select {
		case <-w.stopRun:
			return
		case w.maxWorker <- struct{}{}:
			job := w.getNextJob()
			if job == nil {
				// release token
				<-w.maxWorker
				// sleep 1 second when all queue are empty
				time.Sleep(time.Second)
				continue
			}
			w.jobChannel <- job
		}
	}
}

// work process task
func (w *Worker) work() {
	for {
		select {
		// TODO when to exit, maybe still some msg on job channel
		case <-w.stopRun:
			return
		case job := <-w.jobChannel:
			go func() {
				defer func() {
					// release token
					<-w.maxWorker
				}()
				name := job.GetName()
				v, ok := w.server.task.Load(name)
				if ok {
					task := v.(Task)
					// TODO err handle
					_ = task.Run(w.ctx, job)
				}
			}()
		}
	}
}

func (w *Worker) getNextJob() *Job {
	var job *Job = nil
	// TODO sort tasks
	w.server.task.Range(func(key, value any) bool {
		t := value.(Task)
		if !t.CanRun() {
			return true
		}
		c, ok := w.server.conn.Load(t.OnConnect())
		if !ok {
			return true
		}
		switch t.QueueType() {
		case queue.Redis:
			q := queue.NewRedisQueue(c.(*redis.Client))
			job = w.getJob(q, t.OnQueue())
			if job != nil {
				return false
			}
		case queue.Sqs:
			q := queue.NewSqsQueue(c.(*sqs.Client))
			job = w.getJob(q, t.OnQueue())
			if job != nil {
				return false
			}
		}
		return true
	})
	return job
}

func (w *Worker) getJob(q Queue, queueName string) *Job {
	message, err := q.Pop(w.ctx, queueName)
	if err == nil {
		return &Job{
			id:      message.ID,
			name:    message.Type,
			queue:   message.Queue,
			payload: message.Payload,
		}
	}
	return nil
}
