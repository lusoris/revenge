// Package auth provides authentication services for Revenge Go.
package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService implements domain.PasswordService using bcrypt.
type PasswordService struct {
	cost int
}

// newPasswordService creates a new bcrypt-based password service.
// Use NewPasswordService from module.go for fx integration.
func newPasswordService(cost int) *PasswordService {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	if cost > bcrypt.MaxCost {
		cost = bcrypt.MaxCost
	}
	return &PasswordService{cost: cost}
}

// Hash creates a bcrypt hash of the password.
func (s *PasswordService) Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

// Verify checks if a password matches a hash.
// Returns nil on success, error on mismatch or failure.
func (s *PasswordService) Verify(password, hash string) error {
	if password == "" || hash == "" {
		return errors.New("password and hash cannot be empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("password does not match")
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}

	return nil
}
