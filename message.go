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
}
