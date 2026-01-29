// Package auth provides authentication services.
package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/user"
)

var (
	// ErrSetupRequired indicates initial setup has not been completed.
	ErrSetupRequired = errors.New("initial setup required")
)

// Service provides authentication operations.
type Service struct {
	userService    *user.Service
	sessionService *session.Service
	logger         *slog.Logger
}

// NewService creates a new auth service.
func NewService(
	userService *user.Service,
	sessionService *session.Service,
	logger *slog.Logger,
) *Service {
	return &Service{
		userService:    userService,
		sessionService: sessionService,
		logger:         logger.With(slog.String("service", "auth")),
	}
}

// LoginParams contains parameters for login.
type LoginParams struct {
	Username      string
	Password      string
	DeviceName    *string
	DeviceType    *string
	ClientName    *string
	ClientVersion *string
	IPAddress     netip.Addr
	UserAgent     *string
}

// LoginResult contains the result of a successful login.
type LoginResult struct {
	User    *db.User
	Session *db.Session
	Token   string // Raw token to return to client
}

// Login authenticates a user and creates a session.
func (s *Service) Login(ctx context.Context, params LoginParams) (*LoginResult, error) {
	// Authenticate user
	usr, err := s.userService.Authenticate(ctx, params.Username, params.Password)
	if err != nil {
		return nil, err
	}

	// Create session
	result, err := s.sessionService.Create(ctx, session.CreateParams{
		UserID:        usr.ID,
		DeviceName:    params.DeviceName,
		DeviceType:    params.DeviceType,
		ClientName:    params.ClientName,
		ClientVersion: params.ClientVersion,
		IPAddress:     params.IPAddress,
		UserAgent:     params.UserAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	s.logger.Info("User logged in",
		slog.String("user_id", usr.ID.String()),
		slog.String("session_id", result.Session.ID.String()),
	)

	return &LoginResult{
		User:    usr,
		Session: result.Session,
		Token:   result.Token,
	}, nil
}

// Logout deactivates the session for the given token.
func (s *Service) Logout(ctx context.Context, token string) error {
	sess, err := s.sessionService.ValidateToken(ctx, token)
	if err != nil {
		// Already logged out or invalid token - not an error
		return nil
	}

	if err := s.sessionService.Deactivate(ctx, sess.ID); err != nil {
		return fmt.Errorf("deactivate session: %w", err)
	}

	s.logger.Info("User logged out",
		slog.String("session_id", sess.ID.String()),
		slog.String("user_id", sess.UserID.String()),
	)

	return nil
}

// LogoutAll deactivates all sessions for a user.
func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	if err := s.sessionService.DeactivateAllForUser(ctx, userID); err != nil {
		return fmt.Errorf("deactivate all sessions: %w", err)
	}

	s.logger.Info("All sessions logged out for user",
		slog.String("user_id", userID.String()),
	)

	return nil
}

// ValidateToken validates a token and returns the associated user and session.
func (s *Service) ValidateToken(ctx context.Context, token string) (*db.User, *db.Session, error) {
	sess, err := s.sessionService.ValidateToken(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	usr, err := s.userService.GetByID(ctx, sess.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("get user: %w", err)
	}

	// Check if user is disabled
	if usr.IsDisabled {
		// Deactivate the session
		_ = s.sessionService.Deactivate(ctx, sess.ID)
		return nil, nil, user.ErrUserDisabled
	}

	// Update session activity
	_ = s.sessionService.UpdateActivity(ctx, sess.ID, nil)

	return usr, sess, nil
}

// RegisterParams contains parameters for user registration.
type RegisterParams struct {
	Username          string
	Email             *string
	Password          string
	PreferredLanguage *string
}

// Register creates a new user account.
// The first user registered becomes an admin.
func (s *Service) Register(ctx context.Context, params RegisterParams) (*db.User, error) {
	// Check if this is the first user
	hasUsers, err := s.userService.HasAnyUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("check existing users: %w", err)
	}

	// First user is always admin
	isAdmin := !hasUsers

	usr, err := s.userService.Create(ctx, user.CreateParams{
		Username:          params.Username,
		Email:             params.Email,
		Password:          params.Password,
		IsAdmin:           isAdmin,
		MaxRatingLevel:    100, // Full access by default
		AdultEnabled:      false,
		PreferredLanguage: params.PreferredLanguage,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("User registered",
		slog.String("user_id", usr.ID.String()),
		slog.String("username", usr.Username),
		slog.Bool("is_admin", isAdmin),
	)

	return usr, nil
}

// ChangePassword changes a user's password.
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	// Get user
	usr, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	// Validate current password
	if err := s.userService.ValidatePassword(ctx, usr, currentPassword); err != nil {
		return err
	}

	// Update password
	if err := s.userService.UpdatePassword(ctx, userID, newPassword); err != nil {
		return err
	}

	// Optionally: deactivate all sessions to force re-login
	// This is a security best practice after password change
	// Uncomment if desired:
	// _ = s.sessionService.DeactivateAllForUser(ctx, userID)

	s.logger.Info("Password changed",
		slog.String("user_id", userID.String()),
	)

	return nil
}

// IsSetupRequired returns true if initial setup has not been completed.
// Initial setup is required if there are no users in the system.
func (s *Service) IsSetupRequired(ctx context.Context) (bool, error) {
	hasUsers, err := s.userService.HasAnyUsers(ctx)
	if err != nil {
		return false, err
	}
	return !hasUsers, nil
}
