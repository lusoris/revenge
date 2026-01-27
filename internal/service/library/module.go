// Package library provides the library management service.
package library

import (
	"go.uber.org/fx"
)

// Module provides library service dependencies for fx.
var Module = fx.Options(
	fx.Provide(NewService),
)
