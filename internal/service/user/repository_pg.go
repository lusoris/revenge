package user

import (
	"context"
	"database/sql"
	"fmt"
	"net/netip"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// postgresRepository is a PostgreSQL implementation of the Repository interface
type postgresRepository struct {
	queries *db.Queries
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(queries *db.Queries) Repository {
	return &postgresRepository{queries: queries}
}

// ============================================================================
// User Management
// ============================================================================

func (r *postgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*db.SharedUser, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *postgresRepository) GetUserByUsername(ctx context.Context, username string) (*db.SharedUser, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

func (r *postgresRepository) GetUserByEmail(ctx context.Context, email string) (*db.SharedUser, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *postgresRepository) ListUsers(ctx context.Context, filters UserFilters) ([]db.SharedUser, int64, error) {
	// Count total users matching filters
	// Pass nil pointers to queries when filters are not set - SQL handles NULL checking
	count, err := r.queries.CountUsers(ctx, db.CountUsersParams{
		IsActive: filters.IsActive,
		IsAdmin:  filters.IsAdmin,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// List users with pagination
	users, err := r.queries.ListUsers(ctx, db.ListUsersParams{
		IsActive: filters.IsActive,
		IsAdmin:  filters.IsAdmin,
		Limit:    filters.Limit,
		Offset:   filters.Offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, count, nil
}

func (r *postgresRepository) CreateUser(ctx context.Context, params CreateUserParams) (*db.SharedUser, error) {
	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		DisplayName:  params.DisplayName,
		Timezone:     params.Timezone,
		QarEnabled:   params.QarEnabled,
		IsActive:     params.IsActive,
		IsAdmin:      params.IsAdmin,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *postgresRepository) UpdateUser(ctx context.Context, id uuid.UUID, params UpdateUserParams) (*db.SharedUser, error) {
	user, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		UserID:      id,
		Email:       params.Email,
		DisplayName: params.DisplayName,
		AvatarUrl:   params.AvatarURL,
		Timezone:    params.Timezone,
		QarEnabled:  params.QarEnabled,
		IsActive:    params.IsActive,
		IsAdmin:     params.IsAdmin,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func (r *postgresRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	err := r.queries.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:           id,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

func (r *postgresRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	err := r.queries.UpdateLastLogin(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

func (r *postgresRepository) VerifyEmail(ctx context.Context, id uuid.UUID) error {
	err := r.queries.VerifyEmail(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

func (r *postgresRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *postgresRepository) HardDeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.queries.HardDeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete user: %w", err)
	}
	return nil
}

// ============================================================================
// User Preferences
// ============================================================================

func (r *postgresRepository) GetUserPreferences(ctx context.Context, userID uuid.UUID) (*db.SharedUserPreference, error) {
	prefs, err := r.queries.GetUserPreferences(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user preferences not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}
	return &prefs, nil
}

func (r *postgresRepository) UpsertUserPreferences(ctx context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error) {
	upsertParams := db.UpsertUserPreferencesParams{
		UserID:            params.UserID,
		ProfileVisibility: params.ProfileVisibility,
		ShowEmail:         params.ShowEmail,
		ShowActivity:      params.ShowActivity,
		Theme:             params.Theme,
		DisplayLanguage:   params.DisplayLanguage,
		ContentLanguage:   params.ContentLanguage,
		MetadataLanguage:  params.MetadataLanguage,
		ShowAdultContent:  params.ShowAdultContent,
		ShowSpoilers:      params.ShowSpoilers,
		AutoPlayVideos:    params.AutoPlayVideos,
	}

	// Convert JSONB fields (json.RawMessage is []byte)
	if params.EmailNotifications != nil {
		upsertParams.EmailNotifications = []byte(*params.EmailNotifications)
	}
	if params.PushNotifications != nil {
		upsertParams.PushNotifications = []byte(*params.PushNotifications)
	}
	if params.DigestNotifications != nil {
		upsertParams.DigestNotifications = []byte(*params.DigestNotifications)
	}

	prefs, err := r.queries.UpsertUserPreferences(ctx, upsertParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user preferences: %w", err)
	}

	return &prefs, nil
}

func (r *postgresRepository) DeleteUserPreferences(ctx context.Context, userID uuid.UUID) error {
	err := r.queries.DeleteUserPreferences(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user preferences: %w", err)
	}
	return nil
}

// ============================================================================
// User Avatars
// ============================================================================

func (r *postgresRepository) GetCurrentAvatar(ctx context.Context, userID uuid.UUID) (*db.SharedUserAvatar, error) {
	avatar, err := r.queries.GetCurrentAvatar(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("current avatar not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get current avatar: %w", err)
	}
	return &avatar, nil
}

func (r *postgresRepository) GetAvatarByID(ctx context.Context, id uuid.UUID) (*db.SharedUserAvatar, error) {
	avatar, err := r.queries.GetAvatarByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("avatar not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get avatar by ID: %w", err)
	}
	return &avatar, nil
}

func (r *postgresRepository) ListUserAvatars(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.SharedUserAvatar, error) {
	avatars, err := r.queries.ListUserAvatars(ctx, db.ListUserAvatarsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list user avatars: %w", err)
	}
	return avatars, nil
}

func (r *postgresRepository) CreateAvatar(ctx context.Context, params CreateAvatarParams) (*db.SharedUserAvatar, error) {
	createParams := db.CreateAvatarParams{
		UserID:                params.UserID,
		FilePath:              params.FilePath,
		FileSizeBytes:         params.FileSizeBytes,
		MimeType:              params.MimeType,
		Width:                 params.Width,
		Height:                params.Height,
		Version:               params.Version,
		IsAnimated:            params.IsAnimated,
		UploadedFromUserAgent: params.UploadedFromUserAgent,
	}

	// Parse IP if provided
	if params.UploadedFromIP != nil {
		addr, err := netip.ParseAddr(*params.UploadedFromIP)
		if err != nil {
			return nil, fmt.Errorf("failed to parse IP address: %w", err)
		}
		createParams.UploadedFromIp = addr
	}

	avatar, err := r.queries.CreateAvatar(ctx, createParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create avatar: %w", err)
	}

	return &avatar, nil
}

func (r *postgresRepository) UnsetCurrentAvatars(ctx context.Context, userID uuid.UUID) error {
	err := r.queries.UnsetCurrentAvatars(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to unset current avatars: %w", err)
	}
	return nil
}

func (r *postgresRepository) SetCurrentAvatar(ctx context.Context, id uuid.UUID) error {
	err := r.queries.SetCurrentAvatar(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to set current avatar: %w", err)
	}
	return nil
}

func (r *postgresRepository) DeleteAvatar(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteAvatar(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete avatar: %w", err)
	}
	return nil
}

func (r *postgresRepository) HardDeleteAvatar(ctx context.Context, id uuid.UUID) error {
	err := r.queries.HardDeleteAvatar(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete avatar: %w", err)
	}
	return nil
}

func (r *postgresRepository) GetLatestAvatarVersion(ctx context.Context, userID uuid.UUID) (int32, error) {
	version, err := r.queries.GetLatestAvatarVersion(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest avatar version: %w", err)
	}
	return int32(version), nil
}
