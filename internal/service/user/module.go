package user

import (
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Module provides the user service and its dependencies
var Module = fx.Module("user",
	fx.Provide(
		// Repository
		func(queries *db.Queries) Repository {
			return NewPostgresRepository(queries)
		},
		// Service
		NewService,
	),
)
