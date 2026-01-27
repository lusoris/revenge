// Package auth provides authentication services for Jellyfin Go.
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/jellyfin/jellyfin-go/internal/domain"
)

// TokenService implements domain.TokenService using JWT.
type TokenService struct {
	secret              []byte
	accessTokenDuration time.Duration
}

// jwtClaims extends jwt.RegisteredClaims with custom fields.
type jwtClaims struct {
	jwt.RegisteredClaims
	UserID    string `json:"uid"`
	SessionID string `json:"sid"`
	Username  string `json:"usr"`
	IsAdmin   bool   `json:"adm"`
}

// newTokenService creates a new JWT-based token service.
// Use NewTokenService from module.go for fx integration.
func newTokenService(secret string, accessTokenDuration time.Duration) *TokenService {
	if accessTokenDuration <= 0 {
		accessTokenDuration = 15 * time.Minute
	}

	return &TokenService{
		secret:              []byte(secret),
		accessTokenDuration: accessTokenDuration,
	}
}

// GenerateAccessToken creates a new JWT access token.
func (s *TokenService) GenerateAccessToken(claims domain.TokenClaims) (string, error) {
	now := time.Now()

	jwtc := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "jellyfin-go",
			Subject:   claims.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		UserID:    claims.UserID.String(),
		SessionID: claims.SessionID.String(),
		Username:  claims.Username,
		IsAdmin:   claims.IsAdmin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtc)

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a new cryptographically random refresh token.
func (s *TokenService) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// ValidateAccessToken validates a JWT access token and extracts claims.
func (s *TokenService) ValidateAccessToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrSessionExpired
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID in token: %w", err)
	}

	return &domain.TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		Username:  claims.Username,
		IsAdmin:   claims.IsAdmin,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// HashToken creates a SHA-256 hash of a token for storage.
func (s *TokenService) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Ensure TokenService implements domain.TokenService.
var _ domain.TokenService = (*TokenService)(nil)
