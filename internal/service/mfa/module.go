package mfa

import (
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/crypto"
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
	logger *zap.Logger,
	cfg *config.Config,
) *TOTPService {
	// Use application name as issuer
	issuer := "Revenge"
	return NewTOTPService(queries, encryptor, logger, issuer)
}

// NewWebAuthnServiceFromConfig creates a WebAuthn service with config.
func NewWebAuthnServiceFromConfig(
	queries *db.Queries,
	logger *zap.Logger,
	cfg *config.Config,
) (*WebAuthnService, error) {
	// Use server host as RP ID, or localhost
	rpID := cfg.Server.Host
	if rpID == "" || rpID == "0.0.0.0" {
		rpID = "localhost"
	}
	rpName := "Revenge"
	// Build origin from host:port
	origin := fmt.Sprintf("http://%s:%d", rpID, cfg.Server.Port)
	return NewWebAuthnService(queries, logger, rpName, rpID, []string{origin})
}
