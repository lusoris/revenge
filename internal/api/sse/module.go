package sse

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/notification"
	"go.uber.org/fx"
)

// Module provides SSE components for dependency injection.
var Module = fx.Module("sse",
	fx.Provide(
		func(logger *slog.Logger) *Broker {
			return NewBroker(logger)
		},
		func(broker *Broker, tm auth.TokenManager, logger *slog.Logger) *Handler {
			return NewHandler(broker, tm, logger)
		},
	),
	fx.Invoke(func(broker *Broker, svc notification.Service) {
		agent := NewAgent(broker)
		_ = svc.RegisterAgent(agent)
	}),
)
