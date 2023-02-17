package server

import (
	"context"

	"github.com/kovercjm/tool-go/logger"
)

type Server interface {
	Init(*Config, logger.Logger) Server

	Start(context.Context) error
	Stop(context.Context) error

	RPC() any
	API() any
}
