// Package mfa provides multi-factor authentication services
package mfa

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"image/png"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// TOTPService handles TOTP (Time-based One-Time Password) operations
type TOTPService struct {
	queries   *db.Queries
	encryptor *crypto.Encryptor
	logger    *zap.Logger
	issuer    string // Application name shown in authenticator apps
}

// NewTOTPService creates a new TOTP service
func NewTOTPService(
	queries *db.Queries,
	encryptor *crypto.Encryptor,
	logger *zap.Logger,
	issuer string,
) *TOTPService {
	return &TOTPService{
		queries:   queries,
		encryptor: encryptor,
		logger:    logger,
		issuer:    issuer,
	}
}

// TOTPSetup contains the information needed for TOTP enrollment
type TOTPSetup struct {
	Secret string // Base32-encoded secret (user needs to enter manually if QR fails)
	QRCode []byte // PNG image of QR code
	URL    string // otpauth:// URL
}

// GenerateSecret creates a new TOTP secret and QR code for enrollment
func (s *TOTPService) GenerateSecret(ctx context.Context, userID uuid.UUID, accountName string) (*TOTPSetup, error) {
	// Generate random secret (160 bits / 20 bytes as recommended by RFC 6238)
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("failed to generate random secret: %w", err)
	}

	// Base32 encode (standard for TOTP)
	secretBase32 := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)

	// Generate OTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: accountName,
		Secret:      secret,
		Algorithm:   otp.AlgorithmSHA1, // Most compatible with authenticator apps
		Digits:      otp.DigitsSix,     // Standard 6-digit codes
		Period:      30,                // 30-second time window
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Generate QR code
	var qrBuf bytes.Buffer
	img, err := key.Image(256, 256) // 256x256 pixels
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code image: %w", err)
	}

	if err := png.Encode(&qrBuf, img); err != nil {
		return nil, fmt.Errorf("failed to encode QR code as PNG: %w", err)
	}

	// Encrypt secret before storing
	encryptedSecret, err := s.encryptor.EncryptString(secretBase32)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Store encrypted secret in database (upsert: update if exists, create if not)
	// Note: nonce is prepended to encrypted_secret by AES-256-GCM, no separate storage needed
	_, existsErr := s.queries.GetUserTOTPSecret(ctx, userID)
	if existsErr == nil {
		// User has existing secret, update it (re-enrollment)
		err = s.queries.UpdateTOTPSecret(ctx, db.UpdateTOTPSecretParams{
			UserID:          userID,
			EncryptedSecret: encryptedSecret,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update TOTP secret: %w", err)
		}
	} else if errors.Is(existsErr, pgx.ErrNoRows) {
		// No existing secret, create new one
		_, err = s.queries.CreateTOTPSecret(ctx, db.CreateTOTPSecretParams{
			UserID:          userID,
			EncryptedSecret: encryptedSecret,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to store TOTP secret: %w", err)
		}
	} else {
		return nil, fmt.Errorf("failed to check existing TOTP secret: %w", existsErr)
	}

	s.logger.Info("generated TOTP secret",
		zap.String("user_id", userID.String()),
		zap.String("issuer", s.issuer),
	)

	return &TOTPSetup{
		Secret: secretBase32,
		QRCode: qrBuf.Bytes(),
		URL:    key.URL(),
	}, nil
}

// VerifyCode verifies a TOTP code and enables TOTP if this is the first successful verification
func (s *TOTPService) VerifyCode(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	// Get encrypted secret from database
	totpSecret, err := s.queries.GetUserTOTPSecret(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get TOTP secret: %w", err)
	}

	// Decrypt secret
	secretBase32, err := s.encryptor.DecryptString(totpSecret.EncryptedSecret)
	if err != nil {
		return false, fmt.Errorf("failed to decrypt TOTP secret: %w", err)
	}

	// Verify code with time skew tolerance (±1 time step = ±30 seconds)
	valid := totp.Validate(code, secretBase32)
	if !valid {
		s.logger.Debug("invalid TOTP code",
			zap.String("user_id", userID.String()),
		)
		return false, nil
	}

	// Update last used timestamp
	if err := s.queries.UpdateTOTPLastUsed(ctx, userID); err != nil {
		s.logger.Error("failed to update TOTP last used",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		// Don't fail verification if we can't update timestamp
	}

	// If this is the first successful verification, mark as verified and enable
	if !totpSecret.VerifiedAt.Valid {
		if err := s.queries.VerifyTOTPSecret(ctx, userID); err != nil {
			return false, fmt.Errorf("failed to verify TOTP secret: %w", err)
		}

		s.logger.Info("TOTP verified and enabled",
			zap.String("user_id", userID.String()),
		)
	}

	return true, nil
}

// EnableTOTP enables TOTP for a user (must be already verified)
func (s *TOTPService) EnableTOTP(ctx context.Context, userID uuid.UUID) error {
	if err := s.queries.EnableTOTP(ctx, userID); err != nil {
		return fmt.Errorf("failed to enable TOTP: %w", err)
	}

	s.logger.Info("TOTP enabled",
		zap.String("user_id", userID.String()),
	)

	return nil
}

// DisableTOTP disables TOTP for a user
func (s *TOTPService) DisableTOTP(ctx context.Context, userID uuid.UUID) error {
	if err := s.queries.DisableTOTP(ctx, userID); err != nil {
		return fmt.Errorf("failed to disable TOTP: %w", err)
	}

	s.logger.Info("TOTP disabled",
		zap.String("user_id", userID.String()),
	)

	return nil
}

// DeleteTOTP completely removes TOTP for a user
func (s *TOTPService) DeleteTOTP(ctx context.Context, userID uuid.UUID) error {
	if err := s.queries.DeleteTOTPSecret(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete TOTP secret: %w", err)
	}

	s.logger.Info("TOTP deleted",
		zap.String("user_id", userID.String()),
	)

	return nil
}

// HasTOTP checks if a user has TOTP configured
func (s *TOTPService) HasTOTP(ctx context.Context, userID uuid.UUID) (bool, error) {
	_, err := s.queries.GetUserTOTPSecret(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check TOTP: %w", err)
	}
	return true, nil
}
