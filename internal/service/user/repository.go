package user

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Repository defines the data access interface for user operations
type Repository interface {
	// User management
	GetUserByID(ctx context.Context, id uuid.UUID) (*db.SharedUser, error)
	GetUserByUsername(ctx context.Context, username string) (*db.SharedUser, error)
	GetUserByEmail(ctx context.Context, email string) (*db.SharedUser, error)
	ListUsers(ctx context.Context, filters UserFilters) ([]db.SharedUser, int64, error)
	CreateUser(ctx context.Context, params CreateUserParams) (*db.SharedUser, error)
	UpdateUser(ctx context.Context, id uuid.UUID, params UpdateUserParams) (*db.SharedUser, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	VerifyEmail(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	HardDeleteUser(ctx context.Context, id uuid.UUID) error

	// User preferences
	GetUserPreferences(ctx context.Context, userID uuid.UUID) (*db.SharedUserPreference, error)
	UpsertUserPreferences(ctx context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error)
	DeleteUserPreferences(ctx context.Context, userID uuid.UUID) error

	// User avatars
	GetCurrentAvatar(ctx context.Context, userID uuid.UUID) (*db.SharedUserAvatar, error)
	GetAvatarByID(ctx context.Context, id uuid.UUID) (*db.SharedUserAvatar, error)
	ListUserAvatars(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.SharedUserAvatar, error)
	CreateAvatar(ctx context.Context, params CreateAvatarParams) (*db.SharedUserAvatar, error)
	UnsetCurrentAvatars(ctx context.Context, userID uuid.UUID) error
	SetCurrentAvatar(ctx context.Context, id uuid.UUID) error
	DeleteAvatar(ctx context.Context, id uuid.UUID) error
	HardDeleteAvatar(ctx context.Context, id uuid.UUID) error
	GetLatestAvatarVersion(ctx context.Context, userID uuid.UUID) (int32, error)
}

// UserFilters contains query filters for listing users
type UserFilters struct {
	Query    *string
	IsActive *bool
	IsAdmin  *bool
	Limit    int32
	Offset   int32
}

// CreateUserParams contains parameters for creating a new user
type CreateUserParams struct {
	Username     string
	Email        string
	PasswordHash string
	DisplayName  *string
	Timezone     *string
	QarEnabled   *bool
	IsActive     *bool
	IsAdmin      *bool
}

// UpdateUserParams contains optional parameters for updating a user
type UpdateUserParams struct {
	Email       *string
	DisplayName *string
	AvatarURL   *string
	Timezone    *string
	QarEnabled  *bool
	IsActive    *bool
	IsAdmin     *bool
}

// NotificationSettings represents email notification preferences
type NotificationSettings struct {
	Enabled   bool            `json:"enabled"`
	Frequency string          `json:"frequency"` // instant, daily, weekly
	Types     map[string]bool `json:"types,omitempty"`
}

// PushSettings represents push notification preferences
type PushSettings struct {
	Enabled      bool     `json:"enabled"`
	DeviceTokens []string `json:"device_tokens,omitempty"`
}

// DigestSettings represents digest notification preferences
type DigestSettings struct {
	Enabled   bool   `json:"enabled"`
	Frequency string `json:"frequency"` // daily, weekly, monthly
}

// UpsertPreferencesParams contains parameters for creating/updating user preferences
type UpsertPreferencesParams struct {
	UserID              uuid.UUID
	EmailNotifications  *json.RawMessage
	PushNotifications   *json.RawMessage
	DigestNotifications *json.RawMessage
	ProfileVisibility   *string
	ShowEmail           *bool
	ShowActivity        *bool
	Theme               *string
	DisplayLanguage     *string
	ContentLanguage     *string
	MetadataLanguage    *string
	ShowAdultContent    *bool
	ShowSpoilers        *bool
	AutoPlayVideos      *bool
}

// CreateAvatarParams contains parameters for creating a new avatar
type CreateAvatarParams struct {
	UserID                uuid.UUID
	FilePath              string
	FileSizeBytes         int64
	MimeType              string
	Width                 int32
	Height                int32
	IsAnimated            *bool
	Version               int32
	UploadedFromIP        *string
	UploadedFromUserAgent *string
}
