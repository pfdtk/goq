package goq

import "context"

type Worker struct {
	app *App
}

func (srv *Worker) StartConsuming(ctx context.Context) error {
	return nil
}

// Connect get connect by name, then get next job to run
func (srv *Worker) Connect(name string) Queue {
	return nil
}
