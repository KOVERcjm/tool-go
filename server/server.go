package server

import (
	"context"

	"github.com/kovercjm/tool-go/logger"
)

type Server interface {
	Init(*Config, logger.Logger) error
	Customize(func(Server) error)

	Start(context.Context) error
	Stop(context.Context) error
}
