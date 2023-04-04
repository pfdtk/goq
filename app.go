package goq

import (
	"context"
	"github.com/pfdtk/goq/internal/connect"
	"sync"
)

type App struct {
	task sync.Map
	// conn client
	conn      sync.Map
	redisConf *connect.RedisConf
	sqsConf   *connect.SqsConf
}

func (app *App) Start(ctx context.Context) error {
	worker := &Worker{app: app}
	err := worker.StartConsuming(ctx)
	return err
}

func (app *App) RegisterTask(task Task) {
	app.task.Store(task.GetName(), task)
}

func (app *App) AddConnect(name string, conn any) {
	app.conn.Store(name, conn)
}
