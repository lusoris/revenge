package mfa

import (
	"fmt"
	"log/slog"
	"time"

	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/infra/cache"
	db "github.com/lusoris/revenge/internal/infra/database/db"
)

// Module provides MFA services.
var Module = fx.Module(
	"mfa",
	fx.Provide(
		NewTOTPServiceFromConfig,
		NewBackupCodesService,
		NewWebAuthnServiceFromConfig,
		NewMFAManager,
	),
)

// NewTOTPServiceFromConfig creates a TOTP service with issuer from config.
func NewTOTPServiceFromConfig(
	queries *db.Queries,
	encryptor *crypto.Encryptor,
	logger *slog.Logger,
	cfg *config.Config,
) *TOTPService {
	// Use application name as issuer
	issuer := "Revenge"
	return NewTOTPService(queries, encryptor, logger, issuer)
}

// NewWebAuthnServiceFromConfig creates a WebAuthn service with config.
func NewWebAuthnServiceFromConfig(
	queries *db.Queries,
	logger *slog.Logger,
	cfg *config.Config,
	cacheClient *cache.Client,
) (*WebAuthnService, error) {
	// Use server host as RP ID, or localhost
	rpID := cfg.Server.Host
	if rpID == "" || rpID == "0.0.0.0" {
		rpID = "localhost"
	}
	rpName := "Revenge"
	// Build origin from host:port
	origin := fmt.Sprintf("http://%s:%d", rpID, cfg.Server.Port)

	// Create dedicated cache for WebAuthn sessions (5 minute TTL)
	var sessionCache *cache.Cache
	if cacheClient != nil {
		var err error
		sessionCache, err = cache.NewNamedCache(cacheClient, 1000, 5*time.Minute, "webauthn")
		if err != nil {
			logger.Warn("failed to create webauthn session cache, sessions will not be cached", slog.Any("error", err))
		}
	}

	return NewWebAuthnService(queries, logger, sessionCache, rpName, rpID, []string{origin})
}
