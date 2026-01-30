// Package request provides QAR content request domain models.
package request

import (
	"go.uber.org/fx"
)

// Module provides request (content request) dependencies for fx.
var Module = fx.Module("qar.request",
	fx.Provide(
		NewSQLCRepository,
		NewService,
	),
)
