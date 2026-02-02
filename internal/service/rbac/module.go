package rbac

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/config"
)

// Module provides the RBAC service.
var Module = fx.Module("rbac",
	fx.Provide(
		NewAdapter,
		NewEnforcer,
		NewService,
	),
	fx.Invoke(func(*Service) {}), // Ensure service is created
)

// NewEnforcer creates a new Casbin enforcer with PostgreSQL adapter.
func NewEnforcer(pool *pgxpool.Pool, cfg *config.Config, logger *zap.Logger) (*casbin.Enforcer, error) {
	adapter := NewAdapter(pool)

	// Load model from config path
	modelPath := cfg.RBAC.ModelPath
	if modelPath == "" {
		modelPath = "config/casbin_model.conf"
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		logger.Error("failed to create Casbin enforcer",
			zap.String("model_path", modelPath),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create Casbin enforcer: %w", err)
	}

	// Load policy from database
	if err := enforcer.LoadPolicy(); err != nil {
		logger.Error("failed to load policy", zap.Error(err))
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	logger.Info("Casbin enforcer initialized",
		zap.String("model_path", modelPath),
	)

	return enforcer, nil
}
