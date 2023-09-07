package gin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kovercjm/tool-go/logger"
	"github.com/kovercjm/tool-go/server"
	"net/http"
	"time"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	GinEngine  *gin.Engine
	HTTPServer *http.Server

	config *server.APIConfig
	logger logger.Logger
}

func (s Server) Init(config *server.Config, logger logger.Logger) server.Server {
	if config == nil || logger == nil {
		return nil
	}
	s.config = &server.APIConfig{Port: config.APIConfig.Port}
	s.logger = logger.NoCaller()

	s.GinEngine = gin.New()
	s.GinEngine.Use(
		APILogging(s.logger),
		ErrorFormatter(),
		PanicRecovery(s.logger),
	)
	return &s
}

func (s Server) Start(ctx context.Context) error {
	address := fmt.Sprintf(":%d", s.config.Port)
	s.HTTPServer = &http.Server{
		Addr:    address,
		Handler: s.GinEngine,
	}
	go func() {
		s.logger.Info("gin api server starting", "listening", address)
		if err := s.HTTPServer.ListenAndServe(); err != nil {
			s.logger.Error("gin api server failed to serve", "error", err)
		}
	}()
	return nil
}

func (s Server) Stop(_ context.Context) error {
	s.logger.Info("gin api server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		s.logger.Error("gin api server shutdown failed", "error", err)
	}

	s.logger.Info("gin api server stopped gracefully")
	return nil
}
