package goq

type Job struct {
	// id job uuid
	id string
	// name use to find task to process
	name string
	// on which queue
	queue string
	// payload job payload
	payload []byte
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
