package settings

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Repository defines the data access interface for settings.
type Repository interface {
	// Server Settings
	GetServerSetting(ctx context.Context, key string) (*db.SharedServerSetting, error)
	ListServerSettings(ctx context.Context) ([]db.SharedServerSetting, error)
	ListServerSettingsByCategory(ctx context.Context, category string) ([]db.SharedServerSetting, error)
	ListPublicServerSettings(ctx context.Context) ([]db.SharedServerSetting, error)
	UpsertServerSetting(ctx context.Context, params db.UpsertServerSettingParams) (*db.SharedServerSetting, error)
	UpdateServerSetting(ctx context.Context, params db.UpdateServerSettingParams) (*db.SharedServerSetting, error)
	DeleteServerSetting(ctx context.Context, key string) error

	// User Settings
	GetUserSetting(ctx context.Context, userID uuid.UUID, key string) (*db.SharedUserSetting, error)
	ListUserSettings(ctx context.Context, userID uuid.UUID) ([]db.SharedUserSetting, error)
	ListUserSettingsByCategory(ctx context.Context, userID uuid.UUID, category string) ([]db.SharedUserSetting, error)
	UpsertUserSetting(ctx context.Context, params db.UpsertUserSettingParams) (*db.SharedUserSetting, error)
	UpdateUserSetting(ctx context.Context, params db.UpdateUserSettingParams) (*db.SharedUserSetting, error)
	DeleteUserSetting(ctx context.Context, userID uuid.UUID, key string) error
	DeleteAllUserSettings(ctx context.Context, userID uuid.UUID) error
}

// MarshalValue converts a Go value to JSONB for database storage.
func MarshalValue(v any) (json.RawMessage, error) {
	return json.Marshal(v)
}

// UnmarshalValue converts JSONB from database to a Go value.
func UnmarshalValue(data json.RawMessage, v any) error {
	return json.Unmarshal(data, v)
}
