// Package repository provides PostgreSQL implementations of domain repositories.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jellyfin/jellyfin-go/internal/domain"
	"github.com/jellyfin/jellyfin-go/internal/infra/database/db"
)

// UserRepository implements domain.UserRepository using PostgreSQL.
type UserRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewUserRepository creates a new PostgreSQL user repository.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// GetByID retrieves a user by their unique ID.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return mapDBUserToDomain(&user), nil
}

// GetByUsername retrieves a user by their username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return mapDBUserToDomain(&user), nil
}

// GetByEmail retrieves a user by their email address.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return mapDBUserToDomain(&user), nil
}

// List retrieves users with pagination.
func (r *UserRepository) List(ctx context.Context, limit, offset int32) ([]*domain.User, error) {
	users, err := r.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	result := make([]*domain.User, len(users))
	for i := range users {
		result[i] = mapDBUserToDomain(&users[i])
	}
	return result, nil
}

// Create creates a new user and returns the created entity.
func (r *UserRepository) Create(ctx context.Context, params domain.CreateUserParams) (*domain.User, error) {
	dbParams := db.CreateUserParams{
		Username: params.Username,
		IsAdmin:  params.IsAdmin,
	}

	if params.Email != nil {
		dbParams.Email = pgtype.Text{String: *params.Email, Valid: true}
	}
	if params.PasswordHash != nil {
		dbParams.PasswordHash = pgtype.Text{String: *params.PasswordHash, Valid: true}
	}
	if params.DisplayName != nil {
		dbParams.DisplayName = pgtype.Text{String: *params.DisplayName, Valid: true}
	}

	user, err := r.queries.CreateUser(ctx, dbParams)
	if err != nil {
		// Check for unique constraint violations
		if isUniqueViolation(err, "users_username_key") {
			return nil, domain.ErrDuplicateUsername
		}
		if isUniqueViolation(err, "users_email_key") {
			return nil, domain.ErrDuplicateEmail
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return mapDBUserToDomain(&user), nil
}

// Update updates an existing user.
func (r *UserRepository) Update(ctx context.Context, params domain.UpdateUserParams) error {
	// Check if user exists
	_, err := r.queries.GetUserByID(ctx, params.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to get user for update: %w", err)
	}

	// Update basic fields (email, display_name)
	if params.Email != nil || params.DisplayName != nil {
		dbParams := db.UpdateUserParams{
			ID: params.ID,
		}
		if params.Email != nil {
			dbParams.Email = pgtype.Text{String: *params.Email, Valid: true}
		}
		if params.DisplayName != nil {
			dbParams.DisplayName = pgtype.Text{String: *params.DisplayName, Valid: true}
		}
		_, err = r.queries.UpdateUser(ctx, dbParams)
		if err != nil {
			if isUniqueViolation(err, "users_email_key") {
				return domain.ErrDuplicateEmail
			}
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Update admin status separately
	if params.IsAdmin != nil {
		err = r.queries.UpdateUserAdmin(ctx, db.UpdateUserAdminParams{
			ID:      params.ID,
			IsAdmin: *params.IsAdmin,
		})
		if err != nil {
			return fmt.Errorf("failed to update user admin status: %w", err)
		}
	}

	// Update disabled status separately
	if params.IsDisabled != nil {
		err = r.queries.UpdateUserDisabled(ctx, db.UpdateUserDisabledParams{
			ID:         params.ID,
			IsDisabled: *params.IsDisabled,
		})
		if err != nil {
			return fmt.Errorf("failed to update user disabled status: %w", err)
		}
	}

	return nil
}

// Delete removes a user by their ID.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// UpdateLastLogin updates the user's last login timestamp.
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	err := r.queries.UpdateUserLastLogin(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// UpdateLastActivity updates the user's last activity timestamp.
func (r *UserRepository) UpdateLastActivity(ctx context.Context, id uuid.UUID) error {
	err := r.queries.UpdateUserLastActivity(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to update last activity: %w", err)
	}
	return nil
}

// SetPassword updates the user's password hash.
func (r *UserRepository) SetPassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	err := r.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: pgtype.Text{String: passwordHash, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}
	return nil
}

// Count returns the total number of users.
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// CountAdmins returns the number of admin users.
func (r *UserRepository) CountAdmins(ctx context.Context) (int64, error) {
	count, err := r.queries.CountAdminUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count admin users: %w", err)
	}
	return count, nil
}

// UsernameExists checks if a username is already taken.
func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	exists, err := r.queries.UsernameExists(ctx, username)
	if err != nil {
		return false, fmt.Errorf("failed to check username: %w", err)
	}
	return exists, nil
}

// EmailExists checks if an email is already registered.
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := r.queries.EmailExists(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		return false, fmt.Errorf("failed to check email: %w", err)
	}
	return exists, nil
}

// mapDBUserToDomain converts a database user to a domain user.
func mapDBUserToDomain(u *db.User) *domain.User {
	user := &domain.User{
		ID:         u.ID,
		Username:   u.Username,
		IsAdmin:    u.IsAdmin,
		IsDisabled: u.IsDisabled,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}

	if u.Email.Valid {
		user.Email = &u.Email.String
	}
	if u.PasswordHash.Valid {
		user.PasswordHash = &u.PasswordHash.String
	}
	if u.DisplayName.Valid {
		user.DisplayName = &u.DisplayName.String
	}
	if u.LastLoginAt.Valid {
		user.LastLoginAt = &u.LastLoginAt.Time
	}
	if u.LastActivityAt.Valid {
		user.LastActivityAt = &u.LastActivityAt.Time
	}

	return user
}

// Ensure UserRepository implements domain.UserRepository.
var _ domain.UserRepository = (*UserRepository)(nil)
