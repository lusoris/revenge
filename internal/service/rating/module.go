// Package rating provides the content rating service.
package rating

import (
	"go.uber.org/fx"
)

// Module provides rating service dependencies for fx.
var Module = fx.Options(
	fx.Provide(NewService),
)
