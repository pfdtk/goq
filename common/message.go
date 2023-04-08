package common

type Message struct {
	// witch task
	Type string
	// task common
	Payload []byte
	// common id
	ID string
	// witch queue to dispatch
	Queue string
	// Timeout task`s max running time in second
	Timeout int64
}
