package logger

type Option func(Logger)

func Zap() Option {
	return func(l Logger) {
		l = ZapLogger{}
	}
}
