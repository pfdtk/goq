package goq

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq/internal/queue"
	"github.com/redis/go-redis/v9"
	"runtime/debug"
	"time"
)

type Worker struct {
	server     *Server
	stopRun    chan struct{}
	maxWorker  chan struct{}
	jobChannel chan *Job
	ctx        context.Context
	// todo logger
}

func (w *Worker) StartConsuming() error {
	// TODO single coroutine, maybe a litter slow
	w.consume()
	w.work()
	return nil
}

func (w *Worker) StopConsuming() {
	close(w.stopRun)
}

// consume only pop message from queue and send to work channel
func (w *Worker) consume() {
	w.server.wg.Add(1)
	go func() {
		defer w.server.wg.Done()
		for {
			select {
			case <-w.stopRun:
				close(w.jobChannel)
				return
			case w.maxWorker <- struct{}{}:
				// select is range, so check before run
				select {
				case <-w.stopRun:
					close(w.jobChannel)
					return
				default:
				}
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
	}()
}

// work process task
func (w *Worker) work() {
	w.server.wg.Add(1)
	go func() {
		defer w.server.wg.Done()
		for {
			select {
			// TODO when to exit, maybe still some msg on job channel, we need ack!
			case <-w.stopRun:
				return
			case job, ok := <-w.jobChannel:
				// stop working when channel was closed
				if !ok {
					return
				}
				go func() {
					defer func() {
						// release token
						<-w.maxWorker
					}()
					// todo handle err
					_ = w.runTask(job)
				}()
			}
		}
	}()
}

func (w *Worker) runTask(job *Job) (err error) {
	// recover err from task, so that program will not exit
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(string(debug.Stack()))
		}
	}()
	name := job.GetName()
	v, ok := w.server.task.Load(name)
	if ok {
		task := v.(Task)
		err = task.Run(w.ctx, job)
	}
	return err
}

func (w *Worker) getNextJob() *Job {
	tasks := sortTask(&w.server.task)
	for _, t := range tasks {
		if !t.CanRun() {
			continue
		}
		c, ok := w.server.conn.Load(t.OnConnect())
		if !ok {
			continue
		}
		switch t.QueueType() {
		case queue.Redis:
			q := queue.NewRedisQueue(c.(*redis.Client))
			job := w.getJob(q, t.OnQueue())
			if job != nil {
				return job
			}
		case queue.Sqs:
			q := queue.NewSqsQueue(c.(*sqs.Client))
			job := w.getJob(q, t.OnQueue())
			if job != nil {
				return job
			}
		}
	}
	return nil
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
