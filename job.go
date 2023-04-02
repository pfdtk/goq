package goq

type Job struct {
	name    string
	queue   string
	payload []byte
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
