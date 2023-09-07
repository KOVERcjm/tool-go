package init

import (
	"github.com/kovercjm/tool-go/server"
	"github.com/kovercjm/tool-go/server/gin"
	"github.com/kovercjm/tool-go/server/grpc"
)

type newServer struct {
	server.Server
}

func NewServer(options ...serverOption) (server.Server, error) {
	s := &newServer{}
	for _, option := range options {
		option(s)
	}

	return s.Server, nil
}

type serverOption func(*newServer)

var (
	GRPC serverOption = func(s *newServer) {
		s.Server = &grpc.Server{}
	}
	GIN serverOption = func(s *newServer) {
		s.Server = &gin.Server{}
	}
)
