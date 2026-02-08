package rbac

import (
	"context"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/service/activity"
)

// Module provides the RBAC service.
var Module = fx.Module("rbac",
	fx.Provide(
		NewAdapter,
		NewEnforcer,
		newService,
	),
	fx.Invoke(startAutoReload),
)

// newService creates a new RBAC service with activity logger.
func newService(enforcer *casbin.SyncedEnforcer, logger *zap.Logger, activityLogger activity.Logger) *Service {
	return NewService(enforcer, logger, activityLogger)
}

// startAutoReload enables periodic policy reload from the database.
// This ensures all server instances eventually see policy changes made by
// other instances or direct DB modifications (e.g., admin menu, migrations).
// Changes made through the RBAC Service API (AssignRole, AddPolicy, etc.)
// are visible immediately on the current instance via Casbin's in-memory update.
func startAutoReload(lc fx.Lifecycle, enforcer *casbin.SyncedEnforcer, logger *zap.Logger) {
	const reloadInterval = 10 * time.Second

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			enforcer.StartAutoLoadPolicy(reloadInterval)
			logger.Info("casbin auto-reload started", zap.Duration("interval", reloadInterval))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			enforcer.StopAutoLoadPolicy()
			logger.Info("casbin auto-reload stopped")
			return nil
		},
	})
}

// NewEnforcer creates a new thread-safe Casbin enforcer with PostgreSQL adapter.
// SyncedEnforcer is used instead of Enforcer for:
//   - Thread-safe concurrent access (RWMutex around all operations)
//   - Built-in StartAutoLoadPolicy for periodic DB sync across instances
func NewEnforcer(pool *pgxpool.Pool, cfg *config.Config, logger *zap.Logger) (*casbin.SyncedEnforcer, error) {
	adapter := NewAdapter(pool)

	// Load model from config path
	modelPath := cfg.RBAC.ModelPath
	if modelPath == "" {
		modelPath = "config/casbin_model.conf"
	}

	enforcer, err := casbin.NewSyncedEnforcer(modelPath, adapter)
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
