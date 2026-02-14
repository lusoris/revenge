package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TokenManager handles JWT access tokens and refresh tokens
type TokenManager interface {
	// GenerateAccessToken creates a signed JWT access token.
	// sessionID may be uuid.Nil when no session is associated (e.g. API-key flows).
	GenerateAccessToken(userID uuid.UUID, username string, sessionID uuid.UUID) (string, error)

	// GenerateRefreshToken creates a cryptographically secure refresh token
	GenerateRefreshToken() (string, error)

	// ValidateAccessToken verifies JWT signature and expiry
	ValidateAccessToken(token string) (*Claims, error)

	// HashRefreshToken creates SHA-256 hash for database storage
	HashRefreshToken(token string) string

	// ExtractClaims parses JWT payload without validation (for debugging)
	ExtractClaims(token string) (*Claims, error)
}

// Claims represents JWT token claims
type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	SessionID uuid.UUID `json:"session_id,omitempty"`
	IssuedAt  int64     `json:"iat"`
	ExpiresAt int64     `json:"exp"`
}

// jwtManager implements TokenManager using stdlib crypto
type jwtManager struct {
	secret    []byte
	jwtExpiry time.Duration
}

// NewTokenManager creates a new JWT token manager
func NewTokenManager(secret string, jwtExpiry time.Duration) TokenManager {
	return &jwtManager{
		secret:    []byte(secret),
		jwtExpiry: jwtExpiry,
	}
}

// GenerateAccessToken creates a JWT access token
// Format: header.payload.signature (all base64url encoded)
func (m *jwtManager) GenerateAccessToken(userID uuid.UUID, username string, sessionID uuid.UUID) (string, error) {
	now := time.Now()
	issuedAt := now.UnixNano() / int64(time.Millisecond) // Millisecond precision
	expiresAt := now.Add(m.jwtExpiry).UnixNano() / int64(time.Millisecond)

	// JWT Header (HS256)
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// JWT Payload
	payload := Claims{
		UserID:    userID,
		Username:  username,
		SessionID: sessionID,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Create signature: HMAC-SHA256(header.payload, secret)
	message := headerEncoded + "." + payloadEncoded
	signature := m.sign(message)

	// Final JWT: header.payload.signature
	return message + "." + signature, nil
}

// GenerateRefreshToken creates a cryptographically secure random token
// Uses crypto/rand as specified in AUTH.md
func (m *jwtManager) GenerateRefreshToken() (string, error) {
	// Generate 32 bytes (256 bits) of random data
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	// Encode as hex string for readability
	return hex.EncodeToString(b), nil
}

// ValidateAccessToken verifies JWT signature and expiry
func (m *jwtManager) ValidateAccessToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	signatureProvided := parts[2]

	// Verify signature
	message := headerEncoded + "." + payloadEncoded
	signatureExpected := m.sign(message)
	if !hmac.Equal([]byte(signatureProvided), []byte(signatureExpected)) {
		return nil, errors.New("invalid token signature")
	}

	// Decode and parse payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Check expiry (claims are in milliseconds, convert to compare)
	nowMs := time.Now().UnixNano() / int64(time.Millisecond)
	if nowMs > claims.ExpiresAt {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// HashRefreshToken creates SHA-256 hash for database storage
// Never store plaintext refresh tokens (security best practice)
func (m *jwtManager) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ExtractClaims parses JWT payload without validation
// Useful for debugging or extracting user info before validation
func (m *jwtManager) ExtractClaims(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return &claims, nil
}

// sign creates HMAC-SHA256 signature and returns base64url encoded string
func (m *jwtManager) sign(message string) string {
	h := hmac.New(sha256.New, m.secret)
	h.Write([]byte(message))
	signature := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(signature)
}
