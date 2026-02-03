// Package crypto provides cryptographic utilities for the application
package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/alexedwards/argon2id"
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

// VerifyPassword verifies a password against an Argon2id hash
func (h *PasswordHasher) VerifyPassword(password, hash string) (bool, error) {
	if password == "" {
		return false, fmt.Errorf("password cannot be empty")
	}
	if hash == "" {
		return false, fmt.Errorf("hash cannot be empty")
	}

	// Verify Argon2id hash
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("argon2id verification failed: %w", err)
	}

	return match, nil
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
