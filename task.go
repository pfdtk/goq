package goq

type Task struct {
	typename string
	connect  string
	queue    string
	payload  []byte
}

func (t *Task) Connect() string {
	return t.connect
}

func (t *Task) Queue() string {
	return t.queue
}

func (t *Task) GetName() string {
	return t.typename
}

func (t *Task) GetPayload() []byte {
	return t.payload
}
