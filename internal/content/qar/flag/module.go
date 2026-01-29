// Package flag provides adult tag domain models (QAR obfuscation: tags → flags).
package flag

import (
	"go.uber.org/fx"
)

// Module provides flag (tag) dependencies for fx.
var Module = fx.Module("qar.flag",
	fx.Provide(
		NewSQLCRepository,
		NewService,
	),
)
