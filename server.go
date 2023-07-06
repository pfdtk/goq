package goq

import (
	"context"
	"github.com/pfdtk/goq/backend"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/event"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/task"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Server struct {
	tasks             sync.Map
	cronTasks         []*cronTask
	maxWorker         int
	worker            *worker
	wg                sync.WaitGroup
	logger            logger.Logger
	migrate           *migrate
	scheduler         *scheduler
	schedulerLocation *time.Location
}

func NewServer(config *ServerConfig) *Server {
	logger.SetLogger(config.Logger)
	return &Server{
		maxWorker: config.MaxWorker,
		logger:    logger.GetLogger(),
	}
}

func (s *Server) MustStart(ctx context.Context) {
	s.logger.Info("starting server...")
	s.mustStartWorker(ctx)
	s.mustStartMigrate(ctx)
	s.mustStartScheduler(ctx)
	// wait for sign to exit
	s.waitSignals()
	// wait for all goroutine to finished
	s.wg.Wait()
}

func (s *Server) mustStartWorker(ctx context.Context) {
	s.logger.Info("starting worker...")
	s.worker = newWorker(ctx, s)
	s.worker.mustStartConsuming()
}

func (s *Server) mustStartMigrate(ctx context.Context) {
	s.logger.Info("starting migrate...")
	s.migrate = newMigrate(ctx, s)
	s.migrate.mustStartMigrate()
}

func (s *Server) mustStartScheduler(ctx context.Context) {
	s.logger.Info("starting scheduler...")
	s.scheduler = newScheduler(ctx, s)
	s.scheduler.mustStartScheduler()
}

func (s *Server) stopServer() {
	s.logger.Info("Gracefully down server...")
	s.migrate.stopMigrating()
	s.scheduler.stopScheduler()
	s.worker.stopConsuming()
}

func (s *Server) RegisterTask(tasks ...task.Task) {
	for _, t := range tasks {
		s.tasks.Store(t.GetName(), t)
	}
}

func (s *Server) RegisterCronTask(spec string, task task.Task) {
	s.cronTasks = append(s.cronTasks, &cronTask{spec: spec, task: task})
}

func (s *Server) AddConnect(name string, conn any) {
	connect.AddConnect(name, conn)
}

func (s *Server) AddRedisConnect(name string, conn *redis.Client) {
	connect.AddRedisConnect(name, conn)
}

func (s *Server) AddSqsConnect(name string, conn *connect.SqsClient) {
	connect.AddSqsConnect(name, conn)
}

func (s *Server) Listen(e event.Event, h event.Handler) {
	event.Listen(e, h)
}

func (s *Server) RegisterBackend(b backend.Backend) {
	backend.RegisterBackend(b)
}

func (s *Server) waitSignals() {
	s.logger.Info("Send signal TERM or INT to terminate the process")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT)
	for {
		sig := <-sigs
		switch sig {
		case unix.SIGTERM, unix.SIGINT:
			s.stopServer()
			return
		default:
		}
	}
}
