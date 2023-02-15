package init

import (
	"os"

	"github.com/kovercjm/tool-go/logger"
	"github.com/kovercjm/tool-go/logger/zap"
)

type newLogger struct {
	logger.Logger
}

func New(config *logger.Config, options ...LoggerOption) (logger.Logger, error) {
	l := &newLogger{}
	for _, option := range options {
		option(l)
	}
	if l.Logger == nil {
		l.Logger = zap.Logger{}
	}

	return l.Logger.Init(config)
}

func Default() (logger.Logger, error) {
	deployment := os.Getenv("DEPLOYMENT") // try to get deployment name from env
	if deployment == "" {
		deployment = "default"
	}
	return New(&logger.Config{Deployment: deployment})
}

type LoggerOption func(logger.Logger)

func Zap() LoggerOption {
	return func(l logger.Logger) {
		l = zap.Logger{}
	}
}
