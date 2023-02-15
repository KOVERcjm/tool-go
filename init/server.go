package init

import (
	"fmt"

	"github.com/kovercjm/tool-go/logger"
	"github.com/kovercjm/tool-go/server"
	"github.com/kovercjm/tool-go/server/grpc"
)

type newServer struct {
	server.Server
	logger.Logger

	rpcChoice rpcServerImpl
	apiChoice apiServerImpl
}

type rpcServerImpl int
type apiServerImpl int

const (
	NoRPC rpcServerImpl = iota
	GRPCImpl
)

const (
	NoAPI apiServerImpl = iota
	GRPCGatewayImpl
)

func NewServer(config *server.Config, options ...Option) (server.Server, error) {
	s := &newServer{}
	for _, option := range options {
		option(s)
	}

	switch {
	case s.apiChoice == NoAPI && s.rpcChoice == NoRPC:
		return nil, fmt.Errorf("no server implementation selected")
	case s.apiChoice == NoAPI && s.rpcChoice == GRPCImpl:
		s.Server = grpc.GRPCServer{}
	case s.apiChoice == GRPCGatewayImpl && s.rpcChoice == NoRPC:
		return nil, fmt.Errorf("grpc gateway requires grpc server")
	case s.apiChoice == GRPCGatewayImpl && s.rpcChoice == GRPCImpl:
		// TODO implement grpc gateway
	default:
		return nil, fmt.Errorf("unknown combination of server choices")
	}

	return s.Server, nil
}

type Option func(*newServer)

func GRPC() Option {
	return func(s *newServer) {
		s.rpcChoice = GRPCImpl
	}
}

func GRPCGateway() Option {
	return func(s *newServer) {
		s.apiChoice = GRPCGatewayImpl
	}
}
