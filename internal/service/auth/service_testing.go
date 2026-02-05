package auth

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/email"
)

// NewServiceForTesting creates a Service instance for testing purposes.
// Email service is optional (nil disables email sending in tests).
func NewServiceForTesting(
	pool *pgxpool.Pool,
	repo Repository,
	tokenManager TokenManager,
	activityLogger activity.Logger,
	jwtExpiry time.Duration,
	refreshExpiry time.Duration,
) *Service {
	return &Service{
		pool:           pool,
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
	pool *pgxpool.Pool,
	repo Repository,
	tokenManager TokenManager,
	activityLogger activity.Logger,
	emailService *email.Service,
	jwtExpiry time.Duration,
	refreshExpiry time.Duration,
) *Service {
	return &Service{
		pool:           pool,
		repo:           repo,
		tokenManager:   tokenManager,
		hasher:         crypto.NewPasswordHasher(),
		activityLogger: activityLogger,
		emailService:   emailService,
		jwtExpiry:      jwtExpiry,
		refreshExpiry:  refreshExpiry,
	}
}
