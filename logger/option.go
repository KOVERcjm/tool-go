package logger

import (
	zapLogger "github.com/kovercjm/tool-go/logger/zap"
)

type Option func(*newLogger)

func Zap() Option {
	return func(o *newLogger) {
		o.Logger = zapLogger.Logger{}
	}
}
