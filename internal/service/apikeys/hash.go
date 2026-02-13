package apikeys

import (
	"crypto/sha256"
	"encoding/hex"
)

// hashAPIKey computes a SHA-256 fingerprint of a raw API key for storage and
// cache lookup. SHA-256 is appropriate here because API keys are high-entropy
// random tokens (not user-chosen passwords), making brute-force infeasible.
// This is NOT password hashing — see the auth service for password storage
// which uses Argon2id.
func hashAPIKey(token string) string { // codeql[go/weak-sensitive-data-hashing]: false positive — API key fingerprint, not password hash
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
