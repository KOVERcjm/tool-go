package logger

type Config struct {
	Debug           bool   `default:"false" envconfig:"DEBUG"`
	Development     bool   `default:"false" envconfig:"DEVELOPMENT"`
	Deployment      string `default:"unknown" envconfig:"DEPLOYMENT"`
	StackTraceLevel string `default:"error" envconfig:"LOG_STACK_TRACE_LEVEL"`

	ExtraFields map[string]interface{}
}
