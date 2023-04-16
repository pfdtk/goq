package event

import "github.com/pfdtk/goq/internal/event"

type JobErrorHandler struct {
}

func NewJobErrorHandler() *JobErrorHandler {
	return &JobErrorHandler{}
}

func (j *JobErrorHandler) Handle(e event.Event) (ct bool, err error) {
	return true, nil
}
