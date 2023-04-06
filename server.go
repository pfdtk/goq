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
	wg        sync.WaitGroup
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
	if err != nil {
		return err
	}
	// wait for all goroutine to finished
	s.wg.Wait()
	return nil
}

func (s *Server) StopServer() {
	s.worker.StopConsuming()
}

func (s *Server) RegisterTask(task Task) {
	s.task.Store(task.GetName(), task)
}

func (s *Server) AddConnect(name string, conn any) {
	s.conn.Store(name, conn)
}
