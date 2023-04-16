package event

import "github.com/pfdtk/goq/internal/event"

type WorkerErrorHandler struct {
}

func NewWorkerErrorHandler() *WorkerErrorHandler {
	return &WorkerErrorHandler{}
}

func (j *WorkerErrorHandler) Handle(e event.Event) (ct bool, err error) {
	return true, nil
}
