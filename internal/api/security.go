package api

import (
	"context"
	"log/slog"

	"github.com/ogen-go/ogen/ogenerrors"

	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/service/auth"
)

// Operations that don't require authentication.
var publicOperations = map[gen.OperationName]bool{
	gen.LoginOperation:         true,
	gen.RegisterOperation:      true,
	gen.GetServerInfoOperation: true,
	gen.GetHealthOperation:     true,
	gen.GetLivenessOperation:   true,
	gen.GetReadinessOperation:  true,
}

// SecurityHandler implements bearer token authentication.
type SecurityHandler struct {
	authService *auth.Service
	logger      *slog.Logger
}

// NewSecurityHandler creates a new security handler.
func NewSecurityHandler(authService *auth.Service, logger *slog.Logger) *SecurityHandler {
	return &SecurityHandler{
		authService: authService,
		logger:      logger.With(slog.String("component", "security")),
	}
}

// HandleBearerAuth validates the bearer token and adds user/session to context.
func (s *SecurityHandler) HandleBearerAuth(ctx context.Context, operationName gen.OperationName, t gen.BearerAuth) (context.Context, error) {
	// Skip auth for public operations
	if publicOperations[operationName] {
		return ctx, ogenerrors.ErrSkipServerSecurity
	}

	if t.Token == "" {
		s.logger.Debug("No token provided",
			slog.String("operation", operationName),
		)
		return ctx, ErrUnauthorized
	}

	// Validate token and get user/session
	user, session, err := s.authService.ValidateToken(ctx, t.Token)
	if err != nil {
		s.logger.Debug("Token validation failed",
			slog.String("operation", operationName),
			slog.String("error", err.Error()),
		)
		return ctx, ErrUnauthorized
	}

	// Add user and session to context
	ctx = contextWithUser(ctx, user)
	ctx = contextWithSession(ctx, session)

	s.logger.Debug("Request authenticated",
		slog.String("operation", operationName),
		slog.String("user_id", user.ID.String()),
		slog.String("session_id", session.ID.String()),
	)

	return ctx, nil
}
