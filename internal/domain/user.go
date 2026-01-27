// Package domain contains core business entities and repository interfaces.
// These interfaces define the contract for data access without implementation details.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents a Jellyfin user entity.
// This is the domain model, separate from database or API representations.
type User struct {
	ID             uuid.UUID
	Username       string
	Email          *string // Nullable
	PasswordHash   *string // Nullable (OIDC users may not have password)
	DisplayName    *string
	IsAdmin        bool
	IsDisabled     bool
	LastLoginAt    *time.Time
	LastActivityAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// CreateUserParams contains parameters for creating a new user.
type CreateUserParams struct {
	Username     string
	Email        *string
	PasswordHash *string
	DisplayName  *string
	IsAdmin      bool
}

// UpdateUserParams contains parameters for updating an existing user.
type UpdateUserParams struct {
	ID          uuid.UUID
	Username    *string
	Email       *string
	DisplayName *string
	IsAdmin     *bool
	IsDisabled  *bool
}

// UserRepository defines the interface for user data access.
// Implementations may use PostgreSQL, caching layers, etc.
type UserRepository interface {
	// GetByID retrieves a user by their unique ID.
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetByUsername retrieves a user by their username.
	GetByUsername(ctx context.Context, username string) (*User, error)

	// GetByEmail retrieves a user by their email address.
	GetByEmail(ctx context.Context, email string) (*User, error)

	// List retrieves users with pagination.
	List(ctx context.Context, limit, offset int32) ([]*User, error)

	// Create creates a new user and returns the created entity.
	Create(ctx context.Context, params CreateUserParams) (*User, error)

	// Update updates an existing user.
	Update(ctx context.Context, params UpdateUserParams) error

	// Delete removes a user by their ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateLastLogin updates the user's last login timestamp.
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error

	// UpdateLastActivity updates the user's last activity timestamp.
	UpdateLastActivity(ctx context.Context, id uuid.UUID) error

	// SetPassword updates the user's password hash.
	SetPassword(ctx context.Context, id uuid.UUID, passwordHash string) error

	// Count returns the total number of users.
	Count(ctx context.Context) (int64, error)

	// CountAdmins returns the number of admin users.
	CountAdmins(ctx context.Context) (int64, error)

	// UsernameExists checks if a username is already taken.
	UsernameExists(ctx context.Context, username string) (bool, error)

	// EmailExists checks if an email is already registered.
	EmailExists(ctx context.Context, email string) (bool, error)
}
