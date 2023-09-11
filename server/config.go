package server

import "github.com/kovercjm/tool-go/logger"

type Config struct {
	RPCConfig
	APIConfig

	logger logger.Logger
}

type RPCConfig struct {
	Port        int `default:"4200" envconfig:"RPC_SERVER_PORT"`
	MessageSize int `default:"20971520" envconfig:"RPC_MSG_SIZE"`
}

type APIConfig struct {
	Port int `default:"4201" envconfig:"API_SERVER_PORT"`
}
