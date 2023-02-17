package gorm

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/kovercjm/tool-go/logger"
)

var _ gormLogger.Interface = (*Logger)(nil)

type Logger struct {
	l        logger.Logger
	logLevel gormLogger.LogLevel

	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func (gl Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	gl.logLevel = level
	return gl
}

func (gl Logger) Info(_ context.Context, s string, i ...interface{}) {
	if gl.logLevel < gormLogger.Info {
		return
	}
	gl.l.Info(fmt.Sprintf(s, i...))
}

func (gl Logger) Warn(_ context.Context, s string, i ...interface{}) {
	if gl.logLevel < gormLogger.Warn {
		return
	}
	gl.l.Warn(fmt.Sprintf(s, i...))
}

func (gl Logger) Error(_ context.Context, s string, i ...interface{}) {
	if gl.logLevel < gormLogger.Error {
		return
	}
	gl.l.Error(fmt.Sprintf(s, i...))
}

func (gl Logger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if gl.logLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []interface{}{"sql", sql, "rows", rows, "ms", float64(elapsed.Nanoseconds()) / 1e6}
	if gl.SlowThreshold != 0 && elapsed > gl.SlowThreshold {
		fields = append(fields, "SLOW SQL over threshold", gl.SlowThreshold)
	}

	switch {
	case err != nil:
		fields = append(fields, "error", err)
		if err == gorm.ErrRecordNotFound && !gl.IgnoreRecordNotFoundError {
			if gl.logLevel < gormLogger.Warn {
				return
			}
			gl.l.Warn("gorm record not found", fields...)
		} else {
			if gl.logLevel < gormLogger.Error {
				return
			}
			gl.l.Error("gorm error", fields...)
		}
	case gl.SlowThreshold != 0 && elapsed > gl.SlowThreshold:
		if gl.logLevel < gormLogger.Warn {
			return
		}
		fields = append(fields, "SLOW SQL over threshold", gl.SlowThreshold)
		gl.l.Warn("gorm slow query", fields...)
	default:
		if gl.logLevel < gormLogger.Info {
			return
		}
		gl.l.Debug("gorm trace", fields...)
	}
}
