package task

import (
	"context"
	"github.com/pfdtk/goq/queue"
	"time"
)

var defaultTimeout = 24 * time.Hour

type Job struct {
	// Id job uuid
	id string
	// Name use to find task to process
	name string
	// on which queue
	queueName string
	// timeout task`s max running time in second
	timeout int64
	// payload job payload
	payload []byte
	// retries max retry times
	retries uint
	// how many times job has been tried
	attempts uint
	// rawMessage raw message
	rawMessage *queue.Message
	queue      queue.Queue
}

func NewJob(q queue.Queue, msg *queue.Message) *Job {
	return &Job{
		id:         msg.ID,
		name:       msg.Type,
		queueName:  msg.Queue,
		payload:    msg.Payload,
		timeout:    msg.Timeout,
		attempts:   msg.Attempts,
		retries:    msg.Retries,
		rawMessage: msg,
		queue:      q,
	}
}

func (j *Job) Release(ctx context.Context, backoff uint) (err error) {
	at := time.Now().Add(time.Duration(backoff) * time.Second)
	return j.queue.Release(ctx, j.QueueName(), j.RawMessage().Reserved, at)
}

func (j *Job) Id() string {
	return j.id
}

func (j *Job) Name() string {
	return j.name
}

func (j *Job) QueueName() string {
	return j.queueName
}

func (j *Job) Payload() []byte {
	return j.payload
}

func (j *Job) Retries() uint {
	return j.retries
}

func (j *Job) Attempts() uint {
	return j.attempts
}

func (j *Job) RawMessage() *queue.Message {
	return j.rawMessage
}

func (j *Job) Queue() queue.Queue {
	return j.queue
}

func (j *Job) TimeoutAt() time.Time {
	if j.timeout == 0 {
		return time.Now().Add(defaultTimeout)
	}
	return time.Now().Add(time.Duration(j.timeout) * time.Second)
}

func (j *Job) IsReachMacAttempts() bool {
	return j.attempts >= j.retries
}
