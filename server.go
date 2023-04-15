package goq

import (
	"context"
	"github.com/pfdtk/goq/handler"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/task"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"sync"
)

type Server struct {
	tasks sync.Map
	// conn client
	conn      sync.Map
	maxWorker int
	worker    *worker
	wg        sync.WaitGroup
	logger    logger.Logger
	migrate   *migrate
	// todo handle
	taskErrorHandle []handler.ErrorJobHandler
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
	worker := newWorker(ctx, s)
	s.worker = worker
	err := worker.startConsuming()
	return err
}

func (s *Server) startMigrate(ctx context.Context) error {
	s.logger.Info("starting migrate...")
	migrate := newMigrate(ctx, s)
	s.migrate = migrate
	return migrate.startMigrate()
}

func (s *Server) stopServer() {
	s.logger.Warn("Gracefully down server...")
	s.migrate.stopMigrating()
	s.worker.stopConsuming()
}

func (s *Server) RegisterTask(task task.Task) {
	s.tasks.Store(task.GetName(), task)
}

func (s *Server) AddConnect(name string, conn any) {
	s.conn.Store(name, conn)
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
