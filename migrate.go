package goq

import (
	"context"
	"errors"
	"fmt"
	"github.com/pfdtk/goq/connect"
	"github.com/pfdtk/goq/event"
	rdq "github.com/pfdtk/goq/internal/queue/redis"
	"github.com/pfdtk/goq/logger"
	"github.com/pfdtk/goq/task"
	"runtime/debug"
	"sync"
	"time"
)

type MigrateType string

const (
	MigrateAck   MigrateType = "ack"
	MigrateDelay MigrateType = "delay"
)

type migrate struct {
	tasks    *sync.Map
	wg       *sync.WaitGroup
	stopRun  chan struct{}
	ctx      context.Context
	logger   logger.Logger
	interval time.Duration
}

func newMigrate(ctx context.Context, s *Server) *migrate {
	return &migrate{
		wg:       &s.wg,
		tasks:    &s.tasks,
		ctx:      ctx,
		logger:   s.logger,
		interval: 5 * time.Second,
		stopRun:  make(chan struct{}),
	}
}

func (m *migrate) mustStartMigrate() {
	redisTask := task.GetRedisTask(m.tasks)
	if len(redisTask) != 0 {
		for i := range redisTask {
			m.migrateRedisTasks(redisTask[i], MigrateAck)
			m.migrateRedisTasks(redisTask[i], MigrateDelay)
		}
	}
}

func (m *migrate) stopMigrating() {
	m.logger.Info("stopping migration...")
	close(m.stopRun)
}

func (m *migrate) migrateRedisTasks(t task.Task, cat MigrateType) {
	m.wg.Add(1)
	go func() {
		defer func() {
			if x := recover(); x != nil {
				err := errors.New(fmt.Sprintf("Panic Error: %+v;\nStack: %s", x, string(debug.Stack())))
				m.logger.Error(err)
				m.handleError(err)
			}
		}()
		defer m.wg.Done()
		// tick
		timer := time.NewTimer(m.interval)
		for {
			select {
			case <-m.stopRun:
				m.logger.Infof("migrate stopped, task=%s, name=%s", cat, t.GetName())
				timer.Stop()
				return
			case <-timer.C:
				m.performMigrateTasks(t, cat)
				timer.Reset(m.interval)
			}
		}
	}()
}

func (m *migrate) performMigrateTasks(t task.Task, cat MigrateType) {
	c := connect.GetRedis(t.OnConnect())
	if c == nil {
		m.logger.Errorf("connect not found, name=%s", t.OnConnect())
		m.handleError(errors.New("connect not found"))
		return
	}
	q := rdq.NewRedisQueue(c)
	from := m.getMigrateQueueKey(q, t.OnQueue(), cat)
	moveTo := t.OnQueue()
	err := q.Migrate(m.ctx, from, moveTo)
	if err != nil {
		m.logger.Errorf("migrate error, task=%s, queue=%s", cat, t.OnQueue())
		m.handleError(err)
	} else {
		m.logger.Infof("migrate success, task=%s, queue=%s", cat, t.OnQueue())
	}
}

func (m *migrate) getMigrateQueueKey(
	q *rdq.Queue,
	qn string,
	cat MigrateType) string {
	switch cat {
	case MigrateAck:
		return q.GetReservedKey(qn)
	case MigrateDelay:
		return q.GetDelayedKey(qn)
	default:
		panic("invalid migrate type")
	}
}

func (m *migrate) handleError(err error) {
	event.Dispatch(NewMigrateErrorEvent(err))
}
