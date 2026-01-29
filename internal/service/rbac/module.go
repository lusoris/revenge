package rbac

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Module provides RBAC service dependencies.
var Module = fx.Module("rbac",
	fx.Provide(func(pool *pgxpool.Pool, queries *db.Queries, logger *slog.Logger) (*CasbinService, error) {
		return NewCasbinService(pool, queries, logger)
	}),
)
