package server

import (
	"context"
	"github.com/kovercjm/tool-go/logger"
)

type Server interface {
	Init(*RPCConfig, logger.Logger) Server

	Start(context.Context) error
	Stop(context.Context) error
}
