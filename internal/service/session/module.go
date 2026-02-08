package session

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"go.uber.org/fx"
)

// Module provides session service dependencies
var Module = fx.Module("session",
	fx.Provide(
		NewService,
		NewRepositoryPG,
	),
	fx.Invoke(func(*Service) {}), // Ensure service is instantiated
)

// NewService creates a new session service with configuration
func NewService(
	repo Repository,
	logger *slog.Logger,
	cfg *config.Config,
) *Service {
	// Use session config with fallbacks
	tokenLength := cfg.Session.TokenLength
	if tokenLength == 0 {
		tokenLength = 32 // 32 bytes = 64 hex chars
	}

	maxPerUser := cfg.Session.MaxPerUser
	if maxPerUser == 0 {
		maxPerUser = 10 // Max 10 sessions per user
	}

	expiry := cfg.Auth.RefreshExpiry            // Reuse auth refresh expiry
	refreshExpiry := cfg.Auth.RefreshExpiry * 3 // 3x session expiry

	return &Service{
		repo:          repo,
		logger:        logger.With("component", "session"),
		tokenLength:   tokenLength,
		expiry:        expiry,
		refreshExpiry: refreshExpiry,
		maxPerUser:    maxPerUser,
	}
}

// NewRepositoryPG creates a new PostgreSQL session repository
func NewRepositoryPG(queries *db.Queries) Repository {
	return &RepositoryPG{
		queries: queries,
	}
}
