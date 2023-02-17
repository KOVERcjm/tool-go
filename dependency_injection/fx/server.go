package fx

import (
	"context"

	"go.uber.org/fx"

	"github.com/kovercjm/tool-go/server"
)

func ServerLifecycle(lifecycle fx.Lifecycle, server server.Server) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return server.Stop(ctx)
		},
	})
}
