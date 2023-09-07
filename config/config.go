package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func InitFromEnv(target interface{}, envFiles ...string) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(fmt.Sprintf("cannot load .env file: %v\n", err))
	}
	if len(envFiles) > 0 {
		if err := godotenv.Load(envFiles...); err != nil && !errors.Is(err, os.ErrNotExist) {
			panic(fmt.Sprintf("cannot load env files [%s]: %v\n", envFiles, err))
		}
	}

	prefixKey := os.Getenv("ENV_PREFIX")
	if prefixKey == "" {
		prefixKey = os.Getenv("DEPLOYMENT")
	}
	if err := envconfig.Process(prefixKey, target); err != nil {
		panic(fmt.Sprintf("cannot load config from os env: %v", err))
	}
}
