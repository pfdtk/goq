package goq

import (
	"context"
	"sync"
)

type Server struct {
	task sync.Map
	// conn client
	conn      sync.Map
	maxWorker int
	worker    *Worker
}

func NewServer(config *ServerConfig) *Server {
	return &Server{maxWorker: config.MaxWorker}
}

func (s *Server) Start(ctx context.Context) error {
	worker := &Worker{
		server:     s,
		maxWorker:  make(chan struct{}, s.maxWorker),
		stopRun:    make(chan struct{}),
		jobChannel: make(chan *Job, s.maxWorker),
		ctx:        ctx,
	}
	s.worker = worker
	err := worker.StartConsuming()
	return err
}

func (s *Server) RegisterTask(task Task) {
	s.task.Store(task.GetName(), task)
}

func (s *Server) AddConnect(name string, conn any) {
	s.conn.Store(name, conn)
}
