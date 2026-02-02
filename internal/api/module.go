package api

import "go.uber.org/fx"

// Module provides the HTTP API server.
var Module = fx.Module(
	"api",
	fx.Provide(NewServer),
	// Invoke ensures the server is created and lifecycle hooks are registered
	fx.Invoke(func(*Server) {}),
)
