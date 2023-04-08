package goq

import (
	"context"
	"github.com/pfdtk/goq/iface"
	"github.com/pfdtk/goq/internal/common"
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
	worker    *Worker
	wg        sync.WaitGroup
	logger    iface.Logger

	// todo add event control
}

func NewServer(config *ServerConfig) *Server {
	return &Server{maxWorker: config.MaxWorker}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.StartWorker(ctx)
	if err != nil {
		return err
	}
	// TODO move delay tasks to list
	// TODO move priority tasks to list
	// TODO move un ack tasks to list
	// wait for sign to exit
	s.waitSignals()
	// wait for all goroutine to finished
	s.wg.Wait()

	return nil
}

func (s *Server) StartWorker(ctx context.Context) error {
	worker := &Worker{
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
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) waitSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	for {
		sig := <-sigs
		switch sig {
		case unix.SIGTERM:
		case unix.SIGINT:
		case unix.SIGTSTP:
			s.logger.Info("Gracefully down server...")
			s.stopServer()
			break
		default:
			continue
		}
	}
}

func (s *Server) stopServer() {
	s.worker.StopConsuming()
}

func (s *Server) RegisterTask(task iface.Task) {
	s.tasks.Store(task.GetName(), task)
}

func (s *Server) AddConnect(name string, conn any) {
	s.conn.Store(name, conn)
}
