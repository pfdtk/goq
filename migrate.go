package goq

import (
	"context"
	"github.com/pfdtk/goq/iface"
	rdq "github.com/pfdtk/goq/internal/queue/redis"
	"github.com/pfdtk/goq/internal/utils"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type migrate struct {
	conn     *sync.Map
	tasks    *sync.Map
	wg       *sync.WaitGroup
	stopRun  chan struct{}
	ctx      context.Context
	logger   iface.Logger
	interval time.Duration
}

func (m migrate) StartMigrate() error {
	redisTask := utils.GetRedisTask(m.tasks)
	if len(redisTask) != 0 {
		for _, t := range redisTask {
			m.migrateRedisTasks(t, "ack")
			m.migrateRedisTasks(t, "delay")
		}
	}
	return nil
}

func (m migrate) migrateRedisTasks(t iface.Task, cat string) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		// tick
		timer := time.NewTimer(m.interval)
		for {
			select {
			case <-m.stopRun:
				m.logger.Info("stopping migrate redis ack timeout tasks, name=%s", t.GetName())
			case <-timer.C:
				m.performMigrateTasks(t, cat)
			}
		}
	}()
}

func (m migrate) performMigrateTasks(t iface.Task, cat string) {
	c, ok := m.conn.Load(t.OnConnect())
	if !ok {
		m.logger.Errorf("unable to find connect, name=%s", t.OnConnect())
		return
	}
	q := rdq.NewRedisQueue(c.(*redis.Client))
	err := q.Migrate(m.ctx, t.OnQueue(), q.GetDelayedKey(t.OnQueue()))
	if err != nil {
		m.logger.Errorf("execute migrate %s task, queue=%s", cat, t.OnQueue())
	}
}
