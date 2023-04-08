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

var (
	ErrEmptyJob = errors.New("no jobs are ready for processing")
)

type Worker struct {
	server     *Server
	stopRun    chan struct{}
	maxWorker  chan struct{}
	jobChannel chan *Job
	ctx        context.Context
	logger     Logger
}

func (w *Worker) StartConsuming() error {
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
				job, err := w.getNextJob()
				switch {
				case errors.Is(err, ErrEmptyJob):
					// sleep 1 second when all queue are empty
					time.Sleep(time.Second)
					// release token
					<-w.maxWorker
					continue
				case err != nil:
					// release token
					<-w.maxWorker
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
				go w.runTask(job)
			}
		}
	}()
}

func (w *Worker) runTask(job *Job) {
	defer func() {
		// release token
		<-w.maxWorker
	}()
	go func() {
		// goroutine timeout control
		ctx, cancel := context.WithDeadline(w.ctx, job.TimeoutAt())
		defer func() {
			cancel()
		}()
		select {
		// exit if timeout
		case <-ctx.Done():
			// TODO retry if necessary
			w.logger.Warnf("task: %s has been reach it`s deadline", job.Id())
			return
		default:
		}
		// if parent goroutine exit, sub goroutine will be destroy
		go func() {
			// todo handle err
			_ = w.perform(job)
		}()
	}()
}

func (w *Worker) perform(job *Job) (err error) {
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

func (w *Worker) getNextJob() (*Job, error) {
	tasks := sortTask(&w.server.task)
	for _, t := range tasks {
		if t.GetStatus() == 0 || !t.CanRun() {
			continue
		}
		c, ok := w.server.conn.Load(t.OnConnect())
		if !ok {
			continue
		}
		switch t.QueueType() {
		case queue.Redis:
			q := queue.NewRedisQueue(c.(*redis.Client))
			job, err := w.getJob(q, t.OnQueue())
			if err == nil {
				return job, nil
			}
		case queue.Sqs:
			q := queue.NewSqsQueue(c.(*sqs.Client))
			job, err := w.getJob(q, t.OnQueue())
			if err == nil {
				return job, nil
			}
		}
	}
	return nil, ErrEmptyJob
}

func (w *Worker) getJob(q Queue, queueName string) (*Job, error) {
	message, err := q.Pop(w.ctx, queueName)
	if err == nil {
		return &Job{
			id:      message.ID,
			name:    message.Type,
			queue:   message.Queue,
			payload: message.Payload,
			timeout: message.Timeout,
		}, nil
	}
	return nil, err
}
