package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/kovercjm/tool-go/logger"
	"github.com/kovercjm/tool-go/server"
)

var _ server.Server = (*GRPCServer)(nil)

type GRPCServer struct {
	RPCServer *grpc.Server

	config *server.RPCConfig
	logger logger.Logger
}

func (s GRPCServer) Init(config *server.RPCConfig, logger logger.Logger) server.Server {
	return GRPCServer{
		RPCServer: grpc.NewServer(
			grpc.MaxSendMsgSize(config.MessageSize),
			grpc.MaxRecvMsgSize(config.MessageSize),
			grpc.ChainUnaryInterceptor(
				LoggerInterceptor(logger),
			),
		),
		config: &server.RPCConfig{Port: config.Port, MessageSize: config.MessageSize},
		logger: logger.NoCaller(),
	}
}

func (s GRPCServer) Start(ctx context.Context) error {
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
func (s GRPCServer) Stop(_ context.Context) error {
	s.RPCServer.GracefulStop()
	return nil
}
