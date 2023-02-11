package logger

import (
	zapLogger "github.com/kovercjm/tool-go/logger/zap"
	"os"
)

type newLogger struct {
	Logger
}

func New(config *Config, options ...Option) (Logger, error) {
	l := &newLogger{}
	for _, option := range options {
		option(l)
	}
	if l.Logger == nil {
		l.Logger = zapLogger.Logger{}
	}

	return l.Logger.Init(config)
}

func Default() (Logger, error) {
	deployment := os.Getenv("DEPLOYMENT") // try to get deployment name from env
	if deployment == "" {
		deployment = "default"
	}
	return New(&Config{Deployment: deployment})
}
