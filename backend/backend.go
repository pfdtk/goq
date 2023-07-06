package backend

import (
	"github.com/pfdtk/goq/backend/none"
	"github.com/pfdtk/goq/backend/state"
	"github.com/pfdtk/goq/queue"
)

type Backend interface {
	// Pending Mark message as pending
	Pending(message *queue.Message) error
	Started(message *queue.Message) error
	Success(message *queue.Message) error
	Failure(message *queue.Message, err error) error
	// State get state of message
	State(messageId string) (state.State, error)
}

var backend Backend

var defaultBackend Backend = &none.Backend{}

func RegisterBackend(b Backend) {
	backend = b
}

func Get() Backend {
	if backend == nil {
		return defaultBackend
	}
	return backend
}
