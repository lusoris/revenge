package settings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Service defines the business logic interface for settings management.
type Service interface {
	// Server Settings
	GetServerSetting(ctx context.Context, key string) (*ServerSetting, error)
	ListServerSettings(ctx context.Context) ([]ServerSetting, error)
	ListServerSettingsByCategory(ctx context.Context, category string) ([]ServerSetting, error)
	ListPublicServerSettings(ctx context.Context) ([]ServerSetting, error)
	SetServerSetting(ctx context.Context, key string, value interface{}, updatedBy uuid.UUID) (*ServerSetting, error)
	DeleteServerSetting(ctx context.Context, key string) error

	// User Settings
	GetUserSetting(ctx context.Context, userID uuid.UUID, key string) (*UserSetting, error)
	ListUserSettings(ctx context.Context, userID uuid.UUID) ([]UserSetting, error)
	ListUserSettingsByCategory(ctx context.Context, userID uuid.UUID, category string) ([]UserSetting, error)
	SetUserSetting(ctx context.Context, userID uuid.UUID, key string, value interface{}) (*UserSetting, error)
	SetUserSettingsBulk(ctx context.Context, userID uuid.UUID, settings map[string]interface{}) error
	DeleteUserSetting(ctx context.Context, userID uuid.UUID, key string) error
}

// ServerSetting represents a server-wide configuration setting.
type ServerSetting struct {
	Key           string      `json:"key"`
	Value         interface{} `json:"value"`
	Description   *string     `json:"description,omitempty"`
	Category      *string     `json:"category,omitempty"`
	DataType      string      `json:"data_type"`
	IsSecret      *bool       `json:"is_secret,omitempty"`
	IsPublic      *bool       `json:"is_public,omitempty"`
	AllowedValues []string    `json:"allowed_values,omitempty"`
}

// UserSetting represents a user-specific configuration setting.
type UserSetting struct {
	UserID      uuid.UUID   `json:"user_id"`
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Description *string     `json:"description,omitempty"`
	Category    *string     `json:"category,omitempty"`
	DataType    string      `json:"data_type"`
}

// service implements the Service interface.
type service struct {
	repo Repository
}

// NewService creates a new settings service instance.
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// ============================================================================
// Server Settings
// ============================================================================

func (s *service) GetServerSetting(ctx context.Context, key string) (*ServerSetting, error) {
	dbSetting, err := s.repo.GetServerSetting(ctx, key)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("setting %q not found", key)
		}
		return nil, fmt.Errorf("failed to get server setting: %w", err)
	}

	return s.toServerSetting(dbSetting)
}

func (s *service) ListServerSettings(ctx context.Context) ([]ServerSetting, error) {
	dbSettings, err := s.repo.ListServerSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list server settings: %w", err)
	}

	settings := make([]ServerSetting, 0, len(dbSettings))
	for _, dbSetting := range dbSettings {
		setting, err := s.toServerSetting(&dbSetting)
		if err != nil {
			return nil, err
		}
		settings = append(settings, *setting)
	}

	return settings, nil
}

func (s *service) ListServerSettingsByCategory(ctx context.Context, category string) ([]ServerSetting, error) {
	dbSettings, err := s.repo.ListServerSettingsByCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to list server settings by category: %w", err)
	}

	settings := make([]ServerSetting, 0, len(dbSettings))
	for _, dbSetting := range dbSettings {
		setting, err := s.toServerSetting(&dbSetting)
		if err != nil {
			return nil, err
		}
		settings = append(settings, *setting)
	}

	return settings, nil
}

func (s *service) ListPublicServerSettings(ctx context.Context) ([]ServerSetting, error) {
	dbSettings, err := s.repo.ListPublicServerSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list public server settings: %w", err)
	}

	settings := make([]ServerSetting, 0, len(dbSettings))
	for _, dbSetting := range dbSettings {
		setting, err := s.toServerSetting(&dbSetting)
		if err != nil {
			return nil, err
		}
		settings = append(settings, *setting)
	}

	return settings, nil
}

func (s *service) SetServerSetting(ctx context.Context, key string, value interface{}, updatedBy uuid.UUID) (*ServerSetting, error) {
	// Marshal value to JSONB
	jsonValue, err := MarshalValue(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value: %w", err)
	}

	// Try to get existing setting to preserve metadata
	existing, err := s.repo.GetServerSetting(ctx, key)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing setting: %w", err)
	}

	var dbSetting *db.SharedServerSetting
	if existing != nil {
		// Update existing setting
		dbSetting, err = s.repo.UpdateServerSetting(ctx, db.UpdateServerSettingParams{
			Key:       key,
			Value:     jsonValue,
			UpdatedBy: pgtype.UUID{Bytes: updatedBy, Valid: true},
		})
	} else {
		// Create new setting with defaults
		dbSetting, err = s.repo.UpsertServerSetting(ctx, db.UpsertServerSettingParams{
			Key:       key,
			Value:     jsonValue,
			DataType:  "string", // Default type
			UpdatedBy: pgtype.UUID{Bytes: updatedBy, Valid: true},
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to set server setting: %w", err)
	}

	return s.toServerSetting(dbSetting)
}

func (s *service) DeleteServerSetting(ctx context.Context, key string) error {
	if err := s.repo.DeleteServerSetting(ctx, key); err != nil {
		return fmt.Errorf("failed to delete server setting: %w", err)
	}
	return nil
}

// ============================================================================
// User Settings
// ============================================================================

func (s *service) GetUserSetting(ctx context.Context, userID uuid.UUID, key string) (*UserSetting, error) {
	dbSetting, err := s.repo.GetUserSetting(ctx, userID, key)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user setting %q not found for user %s", key, userID)
		}
		return nil, fmt.Errorf("failed to get user setting: %w", err)
	}

	return s.toUserSetting(dbSetting)
}

func (s *service) ListUserSettings(ctx context.Context, userID uuid.UUID) ([]UserSetting, error) {
	dbSettings, err := s.repo.ListUserSettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user settings: %w", err)
	}

	settings := make([]UserSetting, 0, len(dbSettings))
	for _, dbSetting := range dbSettings {
		setting, err := s.toUserSetting(&dbSetting)
		if err != nil {
			return nil, err
		}
		settings = append(settings, *setting)
	}

	return settings, nil
}

func (s *service) ListUserSettingsByCategory(ctx context.Context, userID uuid.UUID, category string) ([]UserSetting, error) {
	dbSettings, err := s.repo.ListUserSettingsByCategory(ctx, userID, category)
	if err != nil {
		return nil, fmt.Errorf("failed to list user settings by category: %w", err)
	}

	settings := make([]UserSetting, 0, len(dbSettings))
	for _, dbSetting := range dbSettings {
		setting, err := s.toUserSetting(&dbSetting)
		if err != nil {
			return nil, err
		}
		settings = append(settings, *setting)
	}

	return settings, nil
}

func (s *service) SetUserSetting(ctx context.Context, userID uuid.UUID, key string, value interface{}) (*UserSetting, error) {
	// Marshal value to JSONB
	jsonValue, err := MarshalValue(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value: %w", err)
	}

	// Upsert setting (create or update)
	dbSetting, err := s.repo.UpsertUserSetting(ctx, db.UpsertUserSettingParams{
		UserID:   userID,
		Key:      key,
		Value:    jsonValue,
		DataType: "string", // Default type
	})

	if err != nil {
		return nil, fmt.Errorf("failed to set user setting: %w", err)
	}

	return s.toUserSetting(dbSetting)
}

func (s *service) SetUserSettingsBulk(ctx context.Context, userID uuid.UUID, settings map[string]interface{}) error {
	for key, value := range settings {
		if _, err := s.SetUserSetting(ctx, userID, key, value); err != nil {
			return fmt.Errorf("failed to set user setting %q: %w", key, err)
		}
	}
	return nil
}

func (s *service) DeleteUserSetting(ctx context.Context, userID uuid.UUID, key string) error {
	if err := s.repo.DeleteUserSetting(ctx, userID, key); err != nil {
		return fmt.Errorf("failed to delete user setting: %w", err)
	}
	return nil
}

// ============================================================================
// Helper Methods
// ============================================================================

func (s *service) toServerSetting(dbSetting *db.SharedServerSetting) (*ServerSetting, error) {
	var value interface{}
	if err := json.Unmarshal(dbSetting.Value, &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal setting value: %w", err)
	}

	return &ServerSetting{
		Key:         dbSetting.Key,
		Value:       value,
		Description: dbSetting.Description,
		Category:    dbSetting.Category,
		DataType:    dbSetting.DataType,
		IsSecret:    dbSetting.IsSecret,
		IsPublic:    dbSetting.IsPublic,
	}, nil
}

func (s *service) toUserSetting(dbSetting *db.SharedUserSetting) (*UserSetting, error) {
	var value interface{}
	if err := json.Unmarshal(dbSetting.Value, &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal setting value: %w", err)
	}

	return &UserSetting{
		UserID:      dbSetting.UserID,
		Key:         dbSetting.Key,
		Value:       value,
		Description: dbSetting.Description,
		Category:    dbSetting.Category,
		DataType:    dbSetting.DataType,
	}, nil
}
