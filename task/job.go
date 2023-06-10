package task

import (
	"context"
	qm "github.com/pfdtk/goq/internal/queue"
	"github.com/pfdtk/goq/queue"
	"sync"
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
	// the time when message was dispatched on queue
	dispatchAt int64
	// rawMessage raw message
	rawMessage *queue.Message
	// queue client
	queue queue.Queue
	// lock
	lock sync.Mutex
	// callback
	errFunc     []func()
	successFunc []func()
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
		dispatchAt: msg.DispatchAt,
		rawMessage: msg,
		queue:      q,
	}
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

func (j *Job) DispatchAt() int64 {
	return j.dispatchAt
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

func (j *Job) Release(ctx context.Context, backoff uint) (err error) {
	at := time.Now().Add(time.Duration(backoff) * time.Second)
	return j.queue.Release(ctx, j.QueueName(), j.RawMessage(), at)
}

// Delete the job from the queue
func (j *Job) Delete(ctx context.Context) error {
	return j.queue.Delete(ctx, j.QueueName(), j.RawMessage())
}

func (j *Job) WhenSuccess(fn func()) {
	if fn == nil {
		return
	}
	j.lock.Lock()
	defer j.lock.Unlock()
	j.successFunc = append(j.successFunc, fn)
}

func (j *Job) WhenFail(fn func()) {
	if fn == nil {
		return
	}
	j.lock.Lock()
	defer j.lock.Unlock()
	j.errFunc = append(j.errFunc, fn)
}

func (j *Job) Then(fn func()) {
	if fn == nil {
		return
	}
	j.lock.Lock()
	defer j.lock.Unlock()
	j.errFunc = append(j.errFunc, fn)
	j.successFunc = append(j.successFunc, fn)
}

func (j *Job) Success() {
	for i := range j.successFunc {
		j.successFunc[i]()
	}
}

func (j *Job) Fail() {
	for i := range j.errFunc {
		j.errFunc[i]()
	}
}

func (j *Job) DispatchNextJobInChain(ctx context.Context) error {
	if len(j.rawMessage.Chain) == 0 {
		return nil
	}
	next := j.rawMessage.Chain[0]
	if len(j.rawMessage.Chain) > 1 {
		chain := j.rawMessage.Chain[1:]
		next.Chain = chain
	}
	q := qm.GetQueue(next.OnConnect, next.QueueType)
	return q.Push(ctx, next)
}
