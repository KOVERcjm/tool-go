package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/kovercjm/tool-go/logger"
	"github.com/kovercjm/tool-go/server"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	RPCServer *grpc.Server

	config *server.RPCConfig
	logger logger.Logger
}

func (s Server) Init(config *server.Config, logger logger.Logger) {
	s.RPCServer = grpc.NewServer(
		grpc.MaxSendMsgSize(config.MessageSize),
		grpc.MaxRecvMsgSize(config.MessageSize),
		grpc.ChainUnaryInterceptor(
			LoggerInterceptor(logger),
		),
	)
	s.config = &server.RPCConfig{Port: config.RPCConfig.Port, MessageSize: config.MessageSize}
	s.logger = logger.NoCaller()
}

func (s Server) Start(ctx context.Context) error {
	address := fmt.Sprintf("0.0.0.0:%d", s.config.Port)
	netListener, err := net.Listen("tcp4", address)
	if err != nil {
		return fmt.Errorf("listen grpc endpoint failed: %w", err)
	}
	go func(ctx context.Context) {
		s.logger.Info("starting", "address to listen", address)
		if err = s.RPCServer.Serve(netListener); err != nil {
			s.logger.Error("grpc server failed to serve", "error", err)
		}
	}(ctx)
	return nil
}

func (s Server) Stop(_ context.Context) error {
	s.RPCServer.GracefulStop()
	return nil
}
