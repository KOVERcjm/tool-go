package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kovercjm/tool-go/logger"
)

func InitFromEnv(target interface{}, envFiles ...string) {
	l, err := logger.Default()
	if err != nil {
		panic(err)
	}

	if len(envFiles) > 0 {
		if err = godotenv.Load(envFiles...); err != nil {
			l.Error("Cannot load env files", "files", envFiles, "error", err)
		}
	} else {
		if err = godotenv.Load(); err != nil {
			l.Info("no .env file found")
		}
	}
	if err = envconfig.Process(os.Getenv("DEPLOYMENT"), target); err != nil {
		l.Error("Cannot load config from env", "error", err)
	}
	l.Info("Load config from env success", "config", target)
}
