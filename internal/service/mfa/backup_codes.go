package mfa

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/crypto"
	db "github.com/lusoris/revenge/internal/infra/database/db"
)

const (
	// BackupCodeCount is the number of backup codes to generate
	BackupCodeCount = 10
	// BackupCodeLength is the length of each backup code in bytes (hex encoded to 16 chars)
	BackupCodeLength = 8
)

var (
	ErrNoBackupCodes     = errors.New("no backup codes available")
	ErrInvalidBackupCode = errors.New("invalid backup code")
	ErrBackupCodeUsed    = errors.New("backup code already used")
)

// BackupCodesService handles backup code generation and verification.
type BackupCodesService struct {
	queries *db.Queries
	hasher  *crypto.PasswordHasher
	logger  *slog.Logger
}

// NewBackupCodesService creates a new backup codes service.
func NewBackupCodesService(queries *db.Queries, logger *slog.Logger) *BackupCodesService {
	return &BackupCodesService{
		queries: queries,
		hasher:  crypto.NewPasswordHasher(),
		logger:  logger,
	}
}

// BackupCode represents a generated backup code before hashing.
type BackupCode struct {
	Code      string // Plain text code (shown to user only once)
	Hash      string // Argon2id hash (stored in database)
	CreatedAt time.Time
}

// GenerateCodes creates a new set of backup codes for a user.
// Returns the plain text codes (should be shown to user only once).
// Existing codes are NOT deleted - use RegenerateCodes for that.
func (s *BackupCodesService) GenerateCodes(ctx context.Context, userID uuid.UUID) ([]string, error) {
	codes := make([]string, BackupCodeCount)
	params := make([]db.CreateBackupCodesParams, BackupCodeCount)

	// Generate random codes
	for i := range BackupCodeCount {
		code, err := generateRandomCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate code: %w", err)
		}

		// Hash the code with Argon2id
		hash, err := s.hasher.HashPassword(code)
		if err != nil {
			return nil, fmt.Errorf("failed to hash code: %w", err)
		}

		codes[i] = formatCode(code)
		params[i] = db.CreateBackupCodesParams{
			UserID:   userID,
			CodeHash: hash,
		}
	}

	// Bulk insert codes
	count, err := s.queries.CreateBackupCodes(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to store backup codes: %w", err)
	}

	s.logger.Info("backup codes generated",
		slog.String("user_id", userID.String()),
		slog.Int64("count", count))

	return codes, nil
}

// RegenerateCodes deletes all existing backup codes and generates a new set.
func (s *BackupCodesService) RegenerateCodes(ctx context.Context, userID uuid.UUID) ([]string, error) {
	// Delete all existing codes
	err := s.queries.DeleteAllBackupCodes(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing codes: %w", err)
	}

	s.logger.Info("deleted existing backup codes",
		slog.String("user_id", userID.String()))

	// Generate new codes
	return s.GenerateCodes(ctx, userID)
}

// VerifyCode verifies a backup code and marks it as used.
// Uses constant-time comparison to prevent timing attacks.
func (s *BackupCodesService) VerifyCode(ctx context.Context, userID uuid.UUID, code string, clientIP string) (bool, error) {
	// Normalize code (remove spaces, uppercase)
	normalizedCode := normalizeCode(code)

	// Get all unused backup codes for the user
	backupCodes, err := s.queries.GetUnusedBackupCodes(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get backup codes: %w", err)
	}

	if len(backupCodes) == 0 {
		return false, ErrNoBackupCodes
	}

	// Try to match the code against all hashes
	// Use constant-time comparison to prevent timing attacks
	var matchedCodeID *uuid.UUID
	for _, bc := range backupCodes {
		match, err := s.hasher.VerifyPassword(normalizedCode, bc.CodeHash)
		if err == nil && match {
			// Found a match
			id := bc.ID
			matchedCodeID = &id
			break
		}
	}

	if matchedCodeID == nil {
		s.logger.Warn("invalid backup code attempt",
			slog.String("user_id", userID.String()),
			slog.String("client_ip", clientIP))
		return false, nil
	}

	// Mark the code as used
	ipAddr, err := netip.ParseAddr(clientIP)
	if err != nil {
		// If IP parsing fails, use unspecified address
		ipAddr = netip.Addr{}
		s.logger.Warn("failed to parse client IP",
			slog.String("client_ip", clientIP),
			slog.Any("error", err))
	}

	err = s.queries.UseBackupCode(ctx, db.UseBackupCodeParams{
		ID:         *matchedCodeID,
		UserID:     userID,
		UsedFromIp: ipAddr,
	})
	if err != nil {
		return false, fmt.Errorf("failed to mark code as used: %w", err)
	}

	s.logger.Info("backup code verified and used",
		slog.String("user_id", userID.String()),
		slog.String("code_id", matchedCodeID.String()),
		slog.String("client_ip", clientIP))

	return true, nil
}

// GetRemainingCount returns the number of unused backup codes for a user.
func (s *BackupCodesService) GetRemainingCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.queries.CountUnusedBackupCodes(ctx, userID)
}

// HasBackupCodes checks if a user has any backup codes (used or unused).
func (s *BackupCodesService) HasBackupCodes(ctx context.Context, userID uuid.UUID) (bool, error) {
	count, err := s.GetRemainingCount(ctx, userID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteAllCodes removes all backup codes for a user.
func (s *BackupCodesService) DeleteAllCodes(ctx context.Context, userID uuid.UUID) error {
	err := s.queries.DeleteAllBackupCodes(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete backup codes: %w", err)
	}

	s.logger.Info("deleted all backup codes",
		slog.String("user_id", userID.String()))

	return nil
}

// Helper functions

// generateRandomCode generates a cryptographically secure random code.
func generateRandomCode() (string, error) {
	bytes := make([]byte, BackupCodeLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// formatCode formats a code for display (e.g., "1234-5678-90ab-cdef").
func formatCode(code string) string {
	// Insert dashes every 4 characters for readability
	if len(code) != 16 {
		return code
	}
	return fmt.Sprintf("%s-%s-%s-%s",
		code[0:4],
		code[4:8],
		code[8:12],
		code[12:16])
}

// normalizeCode removes dashes and converts to lowercase for comparison.
func normalizeCode(code string) string {
	// Remove dashes and spaces
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, " ", "")
	// Convert to lowercase for consistent comparison
	return strings.ToLower(code)
}

// ConstantTimeCompare compares two strings in constant time.
// This prevents timing attacks when comparing sensitive values.
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
