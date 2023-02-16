package zap

import (
	"syscall"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/kovercjm/tool-go/logger"
)

var _ logger.Logger = (*Logger)(nil)

type Logger struct {
	logger *zap.Logger
}

func (l Logger) Init(config *logger.Config) (logger.Logger, error) {
	var zapConfig zap.Config
	options := []zap.Option{zap.AddCallerSkip(1)}

	if config.Development {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("01-02T15:04:05.999")
	} else {
		zapConfig = zap.NewProductionConfig()
		options = append(options, zap.Fields(zap.String("_DEPLOY_", config.Deployment)))
		zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
	}

	if config.Debug {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	if config.StackTraceLevel != "" {
		level, err := zapcore.ParseLevel(config.StackTraceLevel)
		if err != nil {
			return nil, err
		}
		options = append(options, zap.AddStacktrace(level))
	} else {
		options = append(options, zap.AddStacktrace(zapcore.WarnLevel))
	}

	zapConfig.InitialFields = config.ExtraFields
	zapConfig.Sampling = nil

	zapConfig.EncoderConfig.MessageKey = "_MSG_"
	zapConfig.EncoderConfig.LevelKey = "_LEVEL_"
	zapConfig.EncoderConfig.TimeKey = "_TS_"
	zapConfig.EncoderConfig.NameKey = "_NAME_"
	zapConfig.EncoderConfig.CallerKey = "_CALLER_"
	zapConfig.EncoderConfig.StacktraceKey = "_STACKTRACE_"

	zapLogger, err := zapConfig.Build(options...)
	if err != nil {
		return nil, err
	}
	return Logger{logger: zapLogger}, nil
}

func (l Logger) Debug(msg string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Debug(msg, l.sweetenFields(args...)...)
}

func (l Logger) Info(msg string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Info(msg, l.sweetenFields(args...)...)
}

func (l Logger) Warn(msg string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Warn(msg, l.sweetenFields(args...)...)
}

func (l Logger) Error(msg string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Error(msg, l.sweetenFields(args...)...)
}

func (l Logger) sweetenFields(args ...interface{}) []zap.Field {
	if len(args) == 0 {
		return nil
	}

	fields := make([]zap.Field, 0, len(args))
	ignoredFields := make([]interface{}, 0)
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			if i >= len(args) {
				// even number of fields, reached end of slice
				break
			}
			// odd number of fields, reached last dangling element
			ignoredFields = append(ignoredFields, args[i])
			break
		}
		if key, ok := args[i].(string); !ok {
			ignoredFields = append(ignoredFields, args[i], args[i+1])
			break
		} else {
			switch args[i+1].(type) {
			case proto.Message:
				fields = append(fields, zap.Reflect(key, ProtoMessage{args[i+1]}))
			default:
				fields = append(fields, zap.Any(key, args[i+1]))
			}
		}
	}
	if len(ignoredFields) > 0 {
		l.logger.DPanic("invalid key-value pairs", zap.Any("ignored", ignoredFields))
	}
	return fields
}

type ProtoMessage struct {
	Value interface{}
}

func (pm ProtoMessage) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(pm.Value.(proto.Message))
}

func (l Logger) NoCaller() logger.Logger {
	if l.logger == nil {
		return nil
	}
	return Logger{logger: l.logger.WithOptions(zap.WithCaller(false))}
}

func (l Logger) Sync() error {
	if l.logger == nil {
		return nil
	}

	// ignore error like 'sync /dev/stderr: inappropriate ioctl for device' when using MacOS
	// refer to zap issue (https://github.com/uber-go/zap/issues/991#issuecomment-962098428)
	if err := l.logger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) { // ignore ENOTTY error when
		return err
	}
	return nil
}

func Zap(l logger.Logger) *zap.Logger {
	return l.(*Logger).logger
}
