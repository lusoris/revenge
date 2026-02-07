package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides the crypto Encryptor for fx dependency injection.
var Module = fx.Options(
	fx.Provide(provideEncryptor),
)

// provideEncryptor creates an Encryptor from configuration.
//
// Priority:
//  1. auth.encryption_key (hex-encoded 32-byte key) — recommended for production
//  2. Derived from auth.jwt_secret via SHA-256 — acceptable for development
//  3. Deterministic dev key — only when no secrets are configured at all
func provideEncryptor(cfg *config.Config, logger *zap.Logger) (*Encryptor, error) {
	// Option 1: Dedicated encryption key (recommended)
	if cfg.Auth.EncryptionKey != "" {
		key, err := hex.DecodeString(cfg.Auth.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("auth.encryption_key must be hex-encoded: %w", err)
		}
		if len(key) != 32 {
			return nil, fmt.Errorf("auth.encryption_key must be exactly 32 bytes (64 hex chars), got %d bytes", len(key))
		}
		return NewEncryptor(key)
	}

	// Option 2: Derive from JWT secret
	if cfg.Auth.JWTSecret != "" {
		logger.Warn("auth.encryption_key not set, deriving from jwt_secret (set a dedicated key for production)")
		hash := sha256.Sum256([]byte(cfg.Auth.JWTSecret))
		return NewEncryptor(hash[:])
	}

	// Option 3: Development fallback
	logger.Warn("no encryption key configured, using insecure development key")
	hash := sha256.Sum256([]byte("revenge-dev-encryption-key-do-not-use-in-production"))
	return NewEncryptor(hash[:])
}
