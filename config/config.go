package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func InitFromEnv(target interface{}, envFiles ...string) {
	deployment := os.Getenv("DEPLOYMENT")
	if deployment == "" {
		deployment = "default"
	}

	if len(envFiles) > 0 {
		if err := godotenv.Load(envFiles...); err != nil {
			fmt.Printf("cannot load env files (%s): %v", envFiles, err)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			fmt.Println("no .env file found")
		}
	}
	if err := envconfig.Process(deployment, target); err != nil {
		fmt.Printf("cannot load config from env: %v", err)
	}
	fmt.Println("Load config from env success", "config", target)
}
