package none

import (
	"github.com/pfdtk/goq/backend/state"
	"github.com/pfdtk/goq/queue"
)

type Backend struct {
}

func (b Backend) Pending(_ *queue.Message) error {
	return nil
}

func (b Backend) Started(_ *queue.Message) error {
	return nil
}

func (b Backend) Success(_ *queue.Message) error {
	return nil
}

func (b Backend) Failure(_ *queue.Message, _ error) error {
	return nil
}

func (b Backend) State(_ string) (state.State, error) {
	return "", nil
}
