package goq

import (
	"github.com/pfdtk/goq/iface"
)

type ServerConfig struct {
	MaxWorker int
	logger    iface.Logger
}
