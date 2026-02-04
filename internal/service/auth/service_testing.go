package auth

import (
	"time"

	"github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/email"
)

// NewServiceForTesting creates a Service instance for testing purposes.
// Email service is optional (nil disables email sending in tests).
func NewServiceForTesting(
	repo Repository,
	tokenManager TokenManager,
	activityLogger activity.Logger,
	jwtExpiry time.Duration,
	refreshExpiry time.Duration,
) *Service {
	return &Service{
		repo:           repo,
		tokenManager:   tokenManager,
		hasher:         crypto.NewPasswordHasher(),
		activityLogger: activityLogger,
		emailService:   nil, // Email disabled in tests by default
		jwtExpiry:      jwtExpiry,
		refreshExpiry:  refreshExpiry,
	}
}

// NewServiceForTestingWithEmail creates a Service instance for testing with email support.
func NewServiceForTestingWithEmail(
	repo Repository,
	tokenManager TokenManager,
	activityLogger activity.Logger,
	emailService *email.Service,
	jwtExpiry time.Duration,
	refreshExpiry time.Duration,
) *Service {
	return &Service{
		repo:           repo,
		tokenManager:   tokenManager,
		hasher:         crypto.NewPasswordHasher(),
		activityLogger: activityLogger,
		emailService:   emailService,
		jwtExpiry:      jwtExpiry,
		refreshExpiry:  refreshExpiry,
	}
}
