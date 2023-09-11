package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/kovercjm/tool-go/logger"
	"github.com/kovercjm/tool-go/server"
)

type Server struct {
	RPCServer *grpc.Server

	config *server.RPCConfig
	logger logger.Logger
}

func NewServer(config *server.Config, logger logger.Logger) *Server {
	if config == nil || logger == nil {
		panic("Missing critical arguments to init a server")
	}
	s := Server{
		config: &server.RPCConfig{Port: config.RPCConfig.Port, MessageSize: config.RPCConfig.MessageSize},
		logger: logger.NoCaller(),
	}
	return &s
}

func (s *Server) WithDefaultInterceptors() *Server {
	s.RPCServer = grpc.NewServer(
		grpc.MaxSendMsgSize(s.config.MessageSize),
		grpc.MaxRecvMsgSize(s.config.MessageSize),
		grpc.ChainUnaryInterceptor(
			ErrorInterceptor(),
			LoggerInterceptor(s.logger),
			ValidateInterceptor(),
		),
	)
	return s
}

func (s *Server) Start(_ context.Context) error {
	address := fmt.Sprintf("0.0.0.0:%d", s.config.Port)
	netListener, err := net.Listen("tcp4", address)
	if err != nil {
		return fmt.Errorf("listen grpc endpoint failed: %w", err)
	}
	go func() {
		s.logger.Info("starting", "address to listen", address)
		if err = s.RPCServer.Serve(netListener); err != nil {
			s.logger.Error("grpc server failed to serve", "error", err)
		}
	}()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.RPCServer.GracefulStop()
	return nil
}
