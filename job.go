package goq

import "time"

var defaultTimeout = 24 * time.Hour

type Job struct {
	// id job uuid
	id string
	// name use to find task to process
	name string
	// on which queue
	queue string
	// timeout task`s max running time in second
	timeout int64
	// payload job payload
	payload []byte
}

func (j *Job) Timeout() int64 {
	return j.timeout
}

func (j *Job) TimeoutAt() time.Time {
	if j.timeout == 0 {
		return time.Now().Add(defaultTimeout)
	}
	return time.Now().Add(time.Duration(j.timeout) * time.Second)
}

func (j *Job) Id() string {
	return j.id
}

func (j *Job) GetName() string {
	return j.name
}

func (j *Job) GetPayload() []byte {
	return j.payload
}

func (j *Job) GetQueue() string {
	return j.name
}
