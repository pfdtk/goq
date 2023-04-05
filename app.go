package goq

import (
	"context"
	"sync"
)

type App struct {
	task sync.Map
	// conn client
	conn      sync.Map
	maxWorker int
}

func (app *App) Start(ctx context.Context) error {
	worker := &Worker{
		app:        app,
		maxWorker:  make(chan struct{}, app.maxWorker),
		stopRun:    make(chan struct{}),
		jobChannel: make(chan *Job, app.maxWorker),
		ctx:        ctx,
	}
	err := worker.StartConsuming()
	return err
}

func (app *App) RegisterTask(task Task) {
	app.task.Store(task.GetName(), task)
}

func (app *App) AddConnect(name string, conn any) {
	app.conn.Store(name, conn)
}
