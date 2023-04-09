package goq

import (
	"context"
	"github.com/pfdtk/goq/common"
	"github.com/pfdtk/goq/iface"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Server struct {
	tasks sync.Map
	// conn client
	conn      sync.Map
	maxWorker int
	worker    *worker
	wg        sync.WaitGroup
	logger    iface.Logger

	// todo add event control
	// todo add error handle
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		maxWorker: config.MaxWorker,
		logger:    config.logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("starting server...")
	err := s.startWorker(ctx)
	if err != nil {
		return err
	}
	err = s.startMigrate(ctx)
	if err != nil {
		return err
	}
	// wait for sign to exit
	s.waitSignals()
	// wait for all goroutine to finished
	s.wg.Wait()

	return nil
}

func (s *Server) startWorker(ctx context.Context) error {
	s.logger.Info("starting worker...")
	worker := &worker{
		wg:         &s.wg,
		tasks:      &s.tasks,
		maxWorker:  make(chan struct{}, s.maxWorker),
		stopRun:    make(chan struct{}),
		jobChannel: make(chan *common.Job, s.maxWorker),
		ctx:        ctx,
		logger:     s.logger,
		conn:       &s.conn,
	}
	s.worker = worker
	err := worker.StartConsuming()
	return err
}

func (s *Server) startMigrate(ctx context.Context) error {
	s.logger.Info("starting migrate...")
	migrate := &migrate{
		wg:       &s.wg,
		tasks:    &s.tasks,
		ctx:      ctx,
		logger:   s.logger,
		conn:     &s.conn,
		interval: 5 * time.Second,
	}
	return migrate.StartMigrate()
}

func (s *Server) waitSignals() {
	s.logger.Info("Send signal TERM or INT to terminate the process")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT)
	for {
		sig := <-sigs
		switch sig {
		case unix.SIGTERM:
		case unix.SIGINT:
			s.stopServer()
			break
		default:
		}
	}
}

func (s *Server) stopServer() {
	s.logger.Warn("Gracefully down server...")
	s.worker.StopConsuming()
}

func (s *Server) RegisterTask(task iface.Task) {
	s.tasks.Store(task.GetName(), task)
}

func (s *Server) AddConnect(name string, conn any) {
	s.conn.Store(name, conn)
}
