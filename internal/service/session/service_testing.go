package session

import (
	"time"

	"go.uber.org/zap"
)

// NewServiceForTesting creates a Service instance for testing purposes
// This is exported to allow test packages to create Service instances with mock dependencies
func NewServiceForTesting(
	repo Repository,
	logger *zap.Logger,
	tokenLength int,
	expiry time.Duration,
	refreshExpiry time.Duration,
	maxPerUser int,
) *Service {
	return &Service{
		repo:          repo,
		logger:        logger,
		tokenLength:   tokenLength,
		expiry:        expiry,
		refreshExpiry: refreshExpiry,
		maxPerUser:    maxPerUser,
	}
}
