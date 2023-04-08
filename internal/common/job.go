package common

import "time"

var defaultTimeout = 24 * time.Hour

type Job struct {
	// Id job uuid
	Id string
	// Name use to find task to process
	Name string
	// on which Queue
	Queue string
	// Timeout task`s max running time in second
	Timeout int64
	// Payload job Payload
	Payload []byte
	// how many times job has been tried
	Attempts uint
}

func (j *Job) TimeoutAt() time.Time {
	if j.Timeout == 0 {
		return time.Now().Add(defaultTimeout)
	}
	return time.Now().Add(time.Duration(j.Timeout) * time.Second)
}
