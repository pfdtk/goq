package goq

import (
	"context"
	"sync"
)

type App struct {
	task sync.Map
	// multi connect, info from Task
	connect sync.Map
}

func (app *App) Start(ctx context.Context) error {
	worker := &Worker{app: app}
	err := worker.StartConsuming(ctx)
	return err
}

func (app *App) RegisterTask(task *Task) {
	app.task.Store(task.GetName(), task)
}
