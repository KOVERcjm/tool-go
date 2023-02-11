package logger

type Config struct {
	Debug       bool   `default:"false" envconfig:"DEBUG"`
	Development bool   `default:"false" envconfig:"DEV_MODE"`
	Deployment  string `default:"unknown" envconfig:"DEPLOYMENT"`

	ExtraFields     map[string]interface{}
	StackTraceLevel string
}
