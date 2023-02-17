package server

type Config struct {
	RPCConfig
	APIConfig
}

type RPCConfig struct {
	Port        int `default:"4200" envconfig:"GRPC_BIND_PORT"`
	MessageSize int `default:"20971520" envconfig:"GRPC_MSG_SIZE"`
}

type APIConfig struct {
	Port int `default:"4201" envconfig:"HTTP_BIND_PORT"`
}

// TODO rename config key name
