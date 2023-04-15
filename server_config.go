package goq

import (
	"github.com/pfdtk/goq/logger"
)

type ServerConfig struct {
	MaxWorker int
	logger    logger.Logger
}
