package oidc

import (
	"fmt"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides OIDC service dependencies
var Module = fx.Module("oidc",
	fx.Provide(
		newRepositoryPg,
		newService,
		provideConfig,
	),
)

// Config holds OIDC service configuration
type Config struct {
	CallbackURL string
	EncryptKey  []byte
}

// provideConfig extracts OIDC config from app config
func provideConfig(cfg *config.Config) Config {
	// Build base URL from server config
	baseURL := fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)
	return Config{
		CallbackURL: baseURL + "/api/v1/oidc/callback",
		EncryptKey:  []byte(cfg.Auth.JWTSecret), // Reuse JWT secret for now
	}
}

// newRepositoryPg creates a new PostgreSQL repository
func newRepositoryPg(q *db.Queries) Repository {
	return &RepositoryPg{q: q}
}

// newService creates a new OIDC service with fx dependencies
func newService(repo Repository, logger *zap.Logger, cfg Config) *Service {
	return &Service{
		repo:        repo,
		logger:      logger,
		callbackURL: cfg.CallbackURL,
		encryptKey:  cfg.EncryptKey,
	}
}
