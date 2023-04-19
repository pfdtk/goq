package goq

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/internal/event"
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
	cronTasks         []*CronTask
	maxWorker         int
	worker            *worker
	wg                sync.WaitGroup
	logger            logger.Logger
	migrate           *migrate
	scheduler         *scheduler
	schedulerLocation *time.Location
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
	err = s.startScheduler(ctx)
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

func (s *Server) startScheduler(ctx context.Context) error {
	s.logger.Info("starting scheduler...")
	scheduler := newScheduler(ctx, s)
	s.scheduler = scheduler
	return scheduler.startScheduler()
}

func (s *Server) stopServer() {
	s.logger.Info("Gracefully down server...")
	s.migrate.stopMigrating()
	s.worker.stopConsuming()
	s.scheduler.stopScheduler()
}

func (s *Server) RegisterTask(task task.Task) {
	s.tasks.Store(task.GetName(), task)
}

func (s *Server) RegisterCronTask(spec string, task task.Task) {
	s.cronTasks = append(s.cronTasks, &CronTask{spec: spec, task: task})
}

func (s *Server) AddConnect(name string, conn any) {
	connect.AddConnect(name, conn)
}

func (s *Server) AddRedisConnect(name string, conn *redis.Client) {
	connect.AddRedisConnect(name, conn)
}

func (s *Server) AddSqsConnect(name string, conn *sqs.Client) {
	connect.AddSqsConnect(name, conn)
}

func (s *Server) Listen(e event.Event, h event.Handler) {
	event.Listen(e, h)
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
