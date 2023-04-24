package queue

type Message struct {
	// witch task
	Type string `json:"type"`
	// task common
	Payload []byte `json:"payload"`
	// common id
	ID string `json:"id"`
	// witch queue to dispatch
	Queue string `json:"queue"`
	// Timeout task`s max running time in second
	Timeout int64 `json:"timeout"`
	// Retries max retry times
	Retries uint `json:"retries"`
	// how many times job has been tried
	Attempts uint `json:"attempts"`
	// the time when message was dispatched in queue
	DispatchAt int64 `json:"dispatch_at"`
	// UniqueId of a message
	UniqueId string `json:"unique_id"`
	// keep message unique before UniqueTTL
	UniqueTTL int64 `json:"unique_ttl"`
	// use for visibility timeout and move raw message back to list
	Reserved string `json:"-"`
}
