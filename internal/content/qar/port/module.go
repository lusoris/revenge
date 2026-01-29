// Package port provides adult studio domain models (QAR obfuscation: studios → ports).
package port

import (
	"go.uber.org/fx"
)

// Module provides port (studio) dependencies for fx.
var Module = fx.Module("qar.port",
	fx.Provide(
		NewSQLCRepository,
		NewService,
	),
)
