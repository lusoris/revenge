// Package crypto provides cryptographic utilities for the application
package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher provides password hashing and verification
type PasswordHasher struct {
	params *argon2id.Params
}

// NewPasswordHasher creates a new password hasher with default Argon2id parameters
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		params: argon2id.DefaultParams,
	}
}

// NewPasswordHasherWithParams creates a password hasher with custom parameters
func NewPasswordHasherWithParams(params *argon2id.Params) *PasswordHasher {
	return &PasswordHasher{
		params: params,
	}
}

// HashPassword hashes a password using Argon2id
// Returns the hash in PHC string format: $argon2id$v=19$m=65536,t=3,p=2$...
func (h *PasswordHasher) HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hash, err := argon2id.CreateHash(password, h.params)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return hash, nil
}

// VerifyPassword verifies a password against a hash
// Supports both bcrypt (legacy) and argon2id (current) formats
// Bcrypt hashes start with $2a$, $2b$, or $2y$
// Argon2id hashes start with $argon2id$
func (h *PasswordHasher) VerifyPassword(password, hash string) (bool, error) {
	if password == "" {
		return false, fmt.Errorf("password cannot be empty")
	}
	if hash == "" {
		return false, fmt.Errorf("hash cannot be empty")
	}

	// Check if this is a bcrypt hash (legacy format)
	if strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$") {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf("bcrypt verification failed: %w", err)
		}
		return true, nil
	}

	// Argon2id format (current)
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("argon2id verification failed: %w", err)
	}

	return match, nil
}

// NeedsMigration checks if a password hash needs to be migrated from bcrypt to argon2id
func (h *PasswordHasher) NeedsMigration(hash string) bool {
	return strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$")
}

// GenerateSecureToken generates a cryptographically secure random token
// Returns a hex-encoded string of the specified byte length
func GenerateSecureToken(byteLength int) (string, error) {
	if byteLength <= 0 {
		return "", fmt.Errorf("byte length must be positive")
	}

	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
