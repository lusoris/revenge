package auth

import (
	"time"

	"github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/service/activity"
)

// NewServiceForTesting creates a Service instance for testing purposes
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
		jwtExpiry:      jwtExpiry,
		refreshExpiry:  refreshExpiry,
	}
}
