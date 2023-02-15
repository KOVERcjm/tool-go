package repository

import (
	"context"

	"github.com/kovercjm/tool-go/logger"
)

type Repository interface {
	ToCtx(context.Context, interface{}) context.Context
	Ctx(context.Context) Repository

	Init(*Config, logger.Logger) (Repository, error)
}
