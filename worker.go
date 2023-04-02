package goq

import "context"

type Worker struct {
	app *App
}

func (srv *Worker) StartConsuming(ctx context.Context) error {
	return nil
}
