package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/jellyfin/jellyfin-go/internal/domain"
)

func TestTokenService_GenerateAccessToken(t *testing.T) {
	svc := newTokenService("super-secret-key-for-testing-32chars", 15*time.Minute)

	now := time.Now()
	claims := domain.TokenClaims{
		UserID:    uuid.New(),
		SessionID: uuid.New(),
		Username:  "testuser",
		IsAdmin:   true,
		IssuedAt:  now,
		ExpiresAt: now.Add(15 * time.Minute),
	}

	token, err := svc.GenerateAccessToken(claims)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateAccessToken() returned empty token")
	}

	// Token should be a valid JWT (three parts separated by dots)
	parts := 0
	for _, c := range token {
		if c == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("GenerateAccessToken() token does not look like a JWT, got %d dots instead of 2", parts)
	}
}

func TestTokenService_ValidateAccessToken(t *testing.T) {
	svc := newTokenService("super-secret-key-for-testing-32chars", 15*time.Minute)

	userID := uuid.New()
	sessionID := uuid.New()
	now := time.Now()

	originalClaims := domain.TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		Username:  "testuser",
		IsAdmin:   true,
		IssuedAt:  now,
		ExpiresAt: now.Add(15 * time.Minute),
	}

	token, err := svc.GenerateAccessToken(originalClaims)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	parsedClaims, err := svc.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if parsedClaims.UserID != userID {
		t.Errorf("ValidateAccessToken() UserID = %v, want %v", parsedClaims.UserID, userID)
	}
	if parsedClaims.SessionID != sessionID {
		t.Errorf("ValidateAccessToken() SessionID = %v, want %v", parsedClaims.SessionID, sessionID)
	}
	if parsedClaims.Username != "testuser" {
		t.Errorf("ValidateAccessToken() Username = %v, want %v", parsedClaims.Username, "testuser")
	}
	if parsedClaims.IsAdmin != true {
		t.Errorf("ValidateAccessToken() IsAdmin = %v, want %v", parsedClaims.IsAdmin, true)
	}
}

func TestTokenService_ValidateAccessToken_Expired(t *testing.T) {
	svc := newTokenService("super-secret-key-for-testing-32chars", 15*time.Minute)

	// Create an expired token
	now := time.Now()
	claims := domain.TokenClaims{
		UserID:    uuid.New(),
		SessionID: uuid.New(),
		Username:  "testuser",
		IsAdmin:   false,
		IssuedAt:  now.Add(-1 * time.Hour),
		ExpiresAt: now.Add(-30 * time.Minute), // Expired 30 minutes ago
	}

	token, err := svc.GenerateAccessToken(claims)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	_, err = svc.ValidateAccessToken(token)
	if err == nil {
		t.Error("ValidateAccessToken() should return error for expired token")
	}
}

func TestTokenService_ValidateAccessToken_InvalidSignature(t *testing.T) {
	svc1 := newTokenService("secret-key-one-for-testing-32chars", 15*time.Minute)
	svc2 := newTokenService("secret-key-two-for-testing-32chars", 15*time.Minute)

	now := time.Now()
	claims := domain.TokenClaims{
		UserID:    uuid.New(),
		SessionID: uuid.New(),
		Username:  "testuser",
		IsAdmin:   false,
		IssuedAt:  now,
		ExpiresAt: now.Add(15 * time.Minute),
	}

	// Generate with svc1
	token, err := svc1.GenerateAccessToken(claims)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	// Validate with svc2 (different secret)
	_, err = svc2.ValidateAccessToken(token)
	if err == nil {
		t.Error("ValidateAccessToken() should return error for token signed with different secret")
	}
}

func TestTokenService_ValidateAccessToken_InvalidToken(t *testing.T) {
	svc := newTokenService("super-secret-key-for-testing-32chars", 15*time.Minute)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "random string",
			token: "not-a-valid-jwt-token",
		},
		{
			name:  "malformed JWT",
			token: "header.payload",
		},
		{
			name:  "tampered token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.tampered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ValidateAccessToken(tt.token)
			if err == nil {
				t.Errorf("ValidateAccessToken(%q) should return error", tt.token)
			}
		})
	}
}

func TestTokenService_GenerateRefreshToken(t *testing.T) {
	svc := newTokenService("super-secret-key-for-testing-32chars", 15*time.Minute)

	token1, err := svc.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	token2, err := svc.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	// Tokens should not be empty
	if token1 == "" {
		t.Error("GenerateRefreshToken() returned empty token")
	}

	// Tokens should be unique
	if token1 == token2 {
		t.Error("GenerateRefreshToken() should return unique tokens")
	}

	// Token should be base64 encoded (32 bytes -> ~43 chars)
	if len(token1) < 40 {
		t.Errorf("GenerateRefreshToken() token too short: %d chars", len(token1))
	}
}

func TestTokenService_HashToken(t *testing.T) {
	svc := newTokenService("super-secret-key-for-testing-32chars", 15*time.Minute)

	token := "test-token-to-hash"

	hash1 := svc.HashToken(token)
	hash2 := svc.HashToken(token)

	// Same input should produce same hash (deterministic)
	if hash1 != hash2 {
		t.Error("HashToken() should be deterministic")
	}

	// Different input should produce different hash
	hash3 := svc.HashToken("different-token")
	if hash1 == hash3 {
		t.Error("HashToken() should produce different hashes for different inputs")
	}

	// Hash should be hex-encoded SHA-256 (64 chars)
	if len(hash1) != 64 {
		t.Errorf("HashToken() hash length = %d, want 64", len(hash1))
	}
}

func TestTokenService_DefaultDuration(t *testing.T) {
	// Test that zero duration defaults to 15 minutes
	svc := newTokenService("super-secret-key-for-testing-32chars", 0)
	if svc.accessTokenDuration != 15*time.Minute {
		t.Errorf("newTokenService() accessTokenDuration = %v, want %v", svc.accessTokenDuration, 15*time.Minute)
	}
}

func BenchmarkTokenService_GenerateAccessToken(b *testing.B) {
	svc := newTokenService("super-secret-key-for-testing-32chars", 15*time.Minute)
	claims := domain.TokenClaims{
		UserID:    uuid.New(),
		SessionID: uuid.New(),
		Username:  "testuser",
		IsAdmin:   false,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	for b.Loop() {
		_, _ = svc.GenerateAccessToken(claims)
	}
}

func BenchmarkTokenService_ValidateAccessToken(b *testing.B) {
	svc := newTokenService("super-secret-key-for-testing-32chars", 15*time.Minute)
	claims := domain.TokenClaims{
		UserID:    uuid.New(),
		SessionID: uuid.New(),
		Username:  "testuser",
		IsAdmin:   false,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	token, _ := svc.GenerateAccessToken(claims)

	b.ResetTimer()
	for b.Loop() {
		_, _ = svc.ValidateAccessToken(token)
	}
}

func BenchmarkTokenService_HashToken(b *testing.B) {
	svc := newTokenService("super-secret-key-for-testing-32chars", 15*time.Minute)
	token := "test-token-to-hash-in-benchmark"

	for b.Loop() {
		_ = svc.HashToken(token)
	}
}
