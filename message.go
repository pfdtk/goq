package goq

type Message struct {
	// witch task
	Type string
	// task message
	Payload []byte
	// message id
	ID string
	// witch queue to dispatch
	Queue string
	// Timeout task`s max running time in second
	Timeout int64
}
