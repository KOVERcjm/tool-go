package gin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/fx"
)

func Lifecycle(lifecycle fx.Lifecycle, server *Server) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := fmt.Sprintf(":%d", server.config.Port)
			server.HTTPServer = &http.Server{
				Addr:    address,
				Handler: server.GinEngine,
			}
			go func() {
				server.logger.Info("gin gen server starting", "listening", address)
				if err := server.HTTPServer.ListenAndServe(); err != nil {
					server.logger.Error("gin gen server failed to serve", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.logger.Info("gin gen server is shutting down")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO set graceful stop time period
			defer cancel()
			if err := server.HTTPServer.Shutdown(ctx); err != nil {
				server.logger.Error("gin gen server shutdown failed", "error", err)
			}

			server.logger.Info("gin gen server stopped gracefully")
			return nil
		},
	})
}
