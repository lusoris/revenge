// Package user provides user management services for Revenge Go.
package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/domain"
)

// Service implements user management operations.
type Service struct {
	users     domain.UserRepository
	sessions  domain.SessionRepository
	passwords domain.PasswordService
}

// newService creates a new user service.
// Use NewService from module.go for fx integration.
func newService(
	users domain.UserRepository,
	sessions domain.SessionRepository,
	passwords domain.PasswordService,
) *Service {
	return &Service{
		users:     users,
		sessions:  sessions,
		passwords: passwords,
	}
}

// GetByID retrieves a user by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByUsername retrieves a user by username.
func (s *Service) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// List retrieves users with pagination.
func (s *Service) List(ctx context.Context, limit, offset int32) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.users.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// CreateParams contains parameters for creating a new user.
type CreateParams struct {
	Username    string
	Password    string // Will be hashed
	Email       *string
	DisplayName *string
	IsAdmin     bool
}

// Create creates a new user.
func (s *Service) Create(ctx context.Context, params CreateParams) (*domain.User, error) {
	// Validate username
	username := strings.TrimSpace(params.Username)
	if username == "" {
		return nil, errors.New("username is required")
	}
	if len(username) < 3 {
		return nil, errors.New("username must be at least 3 characters")
	}
	if len(username) > 64 {
		return nil, errors.New("username must be at most 64 characters")
	}

	// Check username uniqueness
	exists, err := s.users.UsernameExists(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, domain.ErrDuplicateUsername
	}

	// Check email uniqueness if provided
	var email *string
	if params.Email != nil && *params.Email != "" {
		e := strings.TrimSpace(*params.Email)
		email = &e

		emailExists, err := s.users.EmailExists(ctx, e)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if emailExists {
			return nil, domain.ErrDuplicateEmail
		}
	}

	// Hash password if provided
	var passwordHash *string
	if params.Password != "" {
		if len(params.Password) < 8 {
			return nil, errors.New("password must be at least 8 characters")
		}
		hash, err := s.passwords.Hash(params.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		passwordHash = &hash
	}

	// Create user
	user, err := s.users.Create(ctx, domain.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		DisplayName:  params.DisplayName,
		IsAdmin:      params.IsAdmin,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	slog.Info("user created",
		slog.String("user_id", user.ID.String()),
		slog.String("username", user.Username))

	return user, nil
}

// UpdateParams contains parameters for updating a user.
type UpdateParams struct {
	ID          uuid.UUID
	Username    *string
	Email       *string
	DisplayName *string
	IsAdmin     *bool
	IsDisabled  *bool
}

// Update updates an existing user.
func (s *Service) Update(ctx context.Context, params UpdateParams) error {
	// Get existing user
	user, err := s.users.GetByID(ctx, params.ID)
	if err != nil {
		return err
	}

	// Validate and check username uniqueness if changing
	if params.Username != nil {
		username := strings.TrimSpace(*params.Username)
		if username == "" {
			return errors.New("username cannot be empty")
		}
		if len(username) < 3 {
			return errors.New("username must be at least 3 characters")
		}
		if username != user.Username {
			exists, err := s.users.UsernameExists(ctx, username)
			if err != nil {
				return fmt.Errorf("failed to check username: %w", err)
			}
			if exists {
				return domain.ErrDuplicateUsername
			}
		}
		params.Username = &username
	}

	// Check email uniqueness if changing
	if params.Email != nil && *params.Email != "" {
		email := strings.TrimSpace(*params.Email)
		if user.Email == nil || email != *user.Email {
			exists, err := s.users.EmailExists(ctx, email)
			if err != nil {
				return fmt.Errorf("failed to check email: %w", err)
			}
			if exists {
				return domain.ErrDuplicateEmail
			}
		}
		params.Email = &email
	}

	// Prevent removing last admin
	if params.IsAdmin != nil && !*params.IsAdmin && user.IsAdmin {
		count, err := s.users.CountAdmins(ctx)
		if err != nil {
			return fmt.Errorf("failed to count admins: %w", err)
		}
		if count <= 1 {
			return errors.New("cannot remove the last administrator")
		}
	}

	// Update user
	if err := s.users.Update(ctx, domain.UpdateUserParams{
		ID:          params.ID,
		Username:    params.Username,
		Email:       params.Email,
		DisplayName: params.DisplayName,
		IsAdmin:     params.IsAdmin,
		IsDisabled:  params.IsDisabled,
	}); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	slog.Info("user updated",
		slog.String("user_id", params.ID.String()))

	// If user is being disabled, invalidate all their sessions
	if params.IsDisabled != nil && *params.IsDisabled && !user.IsDisabled {
		if err := s.sessions.DeleteByUser(ctx, params.ID); err != nil {
			slog.Warn("failed to delete sessions for disabled user",
				slog.String("user_id", params.ID.String()),
				slog.Any("error", err))
		}
	}

	return nil
}

// Delete removes a user and all their sessions.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// Get user to check if exists and if last admin
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Prevent deleting last admin
	if user.IsAdmin {
		count, err := s.users.CountAdmins(ctx)
		if err != nil {
			return fmt.Errorf("failed to count admins: %w", err)
		}
		if count <= 1 {
			return errors.New("cannot delete the last administrator")
		}
	}

	// Delete all sessions first
	if err := s.sessions.DeleteByUser(ctx, id); err != nil {
		slog.Warn("failed to delete user sessions",
			slog.String("user_id", id.String()),
			slog.Any("error", err))
	}

	// Delete user
	if err := s.users.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	slog.Info("user deleted",
		slog.String("user_id", id.String()),
		slog.String("username", user.Username))

	return nil
}

// Count returns the total number of users.
func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.users.Count(ctx)
}

// SetPassword sets a user's password.
func (s *Service) SetPassword(ctx context.Context, userID uuid.UUID, password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hash, err := s.passwords.Hash(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.users.SetPassword(ctx, userID, hash); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	// Invalidate all sessions
	if err := s.sessions.DeleteByUser(ctx, userID); err != nil {
		slog.Warn("failed to delete sessions after password set",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
	}

	return nil
}

// Ensure Service implements UserService if we had an interface.
// For now, handlers will use the concrete type.
