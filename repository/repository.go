package repository

import (
	"context"

	"github.com/kovercjm/tool-go/logger"
)

type Repository interface {
	Init(*Config, logger.Logger) (Repository, error)

	ToCtx(context.Context, interface{}) context.Context
	Ctx(context.Context) Repository
}
