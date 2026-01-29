// Package user provides user management services.
package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Role constants matching the database enum.
const (
	RoleAdmin     = "admin"
	RoleModerator = "moderator"
	RoleUser      = "user"
	RoleGuest     = "guest"
)

var (
	// ErrUserNotFound indicates the user was not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserExists indicates a user with that username/email already exists.
	ErrUserExists = errors.New("user already exists")
	// ErrInvalidCredentials indicates invalid login credentials.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserDisabled indicates the user account is disabled.
	ErrUserDisabled = errors.New("user account is disabled")
	// ErrInvalidRole indicates an invalid role was provided.
	ErrInvalidRole = errors.New("invalid role")
)

// ValidRoles contains all valid user roles.
var ValidRoles = []string{RoleAdmin, RoleModerator, RoleUser, RoleGuest}

// IsValidRole checks if the given role is valid.
func IsValidRole(role string) bool {
	for _, r := range ValidRoles {
		if r == role {
			return true
		}
	}
	return false
}

// Service provides user management operations.
type Service struct {
	queries    *db.Queries
	logger     *slog.Logger
	bcryptCost int
}

// NewService creates a new user service.
func NewService(queries *db.Queries, logger *slog.Logger) *Service {
	return &Service{
		queries:    queries,
		logger:     logger.With(slog.String("service", "user")),
		bcryptCost: 12, // Default bcrypt cost
	}
}

// CreateParams contains parameters for creating a user.
type CreateParams struct {
	Username          string
	Email             *string
	Password          string // Plain text password
	Role              string // User role (admin, moderator, user, guest)
	IsAdmin           bool   // Deprecated: use Role instead
	MaxRatingLevel    int32
	AdultEnabled      bool
	PreferredLanguage *string
}

// Create creates a new user with a hashed password.
func (s *Service) Create(ctx context.Context, params CreateParams) (*db.User, error) {
	// Check if username exists
	exists, err := s.queries.UserExistsByUsername(ctx, params.Username)
	if err != nil {
		return nil, fmt.Errorf("check username: %w", err)
	}
	if exists {
		return nil, ErrUserExists
	}

	// Check if email exists (if provided)
	if params.Email != nil {
		exists, err = s.queries.UserExistsByEmail(ctx, params.Email)
		if err != nil {
			return nil, fmt.Errorf("check email: %w", err)
		}
		if exists {
			return nil, ErrUserExists
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), s.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	passwordHash := string(hashedPassword)

	// Determine role - use explicit role if set, fallback to IsAdmin for backwards compat
	role := params.Role
	if role == "" {
		if params.IsAdmin {
			role = RoleAdmin
		} else {
			role = RoleUser
		}
	} else if !IsValidRole(role) {
		return nil, ErrInvalidRole
	}

	// Sync IsAdmin with role for backwards compatibility
	isAdmin := role == RoleAdmin

	// Create user
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Username:          params.Username,
		Email:             params.Email,
		PasswordHash:      &passwordHash,
		IsAdmin:           isAdmin,
		Role:              role,
		MaxRatingLevel:    params.MaxRatingLevel,
		AdultEnabled:      params.AdultEnabled,
		PreferredLanguage: params.PreferredLanguage,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	s.logger.Info("User created",
		slog.String("user_id", user.ID.String()),
		slog.String("username", user.Username),
	)

	return &user, nil
}

// GetByID retrieves a user by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*db.User, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &user, nil
}

// GetByUsername retrieves a user by username.
func (s *Service) GetByUsername(ctx context.Context, username string) (*db.User, error) {
	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email.
func (s *Service) GetByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := s.queries.GetUserByEmail(ctx, &email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &user, nil
}

// ValidatePassword checks if the provided password matches the user's password.
func (s *Service) ValidatePassword(ctx context.Context, user *db.User, password string) error {
	if user.PasswordHash == nil {
		return ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}

	return nil
}

// Authenticate validates credentials and returns the user if valid.
func (s *Service) Authenticate(ctx context.Context, username, password string) (*db.User, error) {
	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.IsDisabled {
		return nil, ErrUserDisabled
	}

	if err := s.ValidatePassword(ctx, &user, password); err != nil {
		return nil, err
	}

	// Update last login
	_ = s.queries.UpdateUserLastLogin(ctx, user.ID)

	s.logger.Info("User authenticated",
		slog.String("user_id", user.ID.String()),
		slog.String("username", user.Username),
	)

	return &user, nil
}

// UpdateParams contains parameters for updating a user.
type UpdateParams struct {
	ID                    uuid.UUID
	Username              *string
	Email                 *string
	Role                  *string // User role (admin, moderator, user, guest)
	IsAdmin               *bool   // Deprecated: use Role instead
	IsDisabled            *bool
	MaxRatingLevel        *int32
	AdultEnabled          *bool
	PreferredLanguage     *string
	PreferredRatingSystem *string
}

// Update updates a user's information.
func (s *Service) Update(ctx context.Context, params UpdateParams) (*db.User, error) {
	// Validate role if provided
	if params.Role != nil && !IsValidRole(*params.Role) {
		return nil, ErrInvalidRole
	}

	// Sync IsAdmin with role for backwards compatibility
	isAdmin := params.IsAdmin
	if params.Role != nil {
		adminVal := *params.Role == RoleAdmin
		isAdmin = &adminVal
	}

	// Convert role to NullUserRole
	var role db.NullUserRole
	if params.Role != nil {
		role = db.NullUserRole{UserRole: db.UserRole(*params.Role), Valid: true}
	}

	user, err := s.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:                    params.ID,
		Username:              params.Username,
		Email:                 params.Email,
		Role:                  role,
		IsAdmin:               isAdmin,
		IsDisabled:            params.IsDisabled,
		MaxRatingLevel:        params.MaxRatingLevel,
		AdultEnabled:          params.AdultEnabled,
		PreferredLanguage:     params.PreferredLanguage,
		PreferredRatingSystem: params.PreferredRatingSystem,
	})
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	s.logger.Info("User updated",
		slog.String("user_id", user.ID.String()),
	)

	return &user, nil
}

// UpdatePassword updates a user's password.
func (s *Service) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.bcryptCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	passwordHash := string(hashedPassword)

	_, err = s.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:           userID,
		PasswordHash: &passwordHash,
	})
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	s.logger.Info("User password updated",
		slog.String("user_id", userID.String()),
	)

	return nil
}

// Delete deletes a user.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	s.logger.Info("User deleted",
		slog.String("user_id", id.String()),
	)

	return nil
}

// List returns a paginated list of users.
func (s *Service) List(ctx context.Context, limit, offset int32) ([]db.User, error) {
	users, err := s.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

// Count returns the total number of users.
func (s *Service) Count(ctx context.Context) (int64, error) {
	count, err := s.queries.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return count, nil
}

// HasAnyUsers checks if there are any users in the system.
func (s *Service) HasAnyUsers(ctx context.Context) (bool, error) {
	count, err := s.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
