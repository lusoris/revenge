// Package settings provides server settings management services.
package settings

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

var (
	// ErrSettingNotFound indicates the setting was not found.
	ErrSettingNotFound = errors.New("setting not found")
)

// Setting category constants.
const (
	CategoryGeneral  = "general"
	CategorySecurity = "security"
	CategoryMedia    = "media"
	CategoryCache    = "cache"
	CategorySearch   = "search"
	CategoryAdult    = "adult"
)

// Common setting keys.
const (
	KeyServerName               = "server.name"
	KeyServerVersion            = "server.version"
	KeyServerTimezone           = "server.timezone"
	KeySecurityRequireAuth      = "security.require_authentication"
	KeySecurityAllowRegistration = "security.allow_registration"
	KeyMediaDefaultProfile      = "media.default_transcoding_profile"
	KeyMediaHWAccel             = "media.enable_hardware_acceleration"
	KeyAdultGloballyEnabled     = "adult.globally_enabled"
)

// Service provides server settings management operations.
type Service struct {
	queries *db.Queries
	logger  *slog.Logger
}

// NewService creates a new settings service.
func NewService(queries *db.Queries, logger *slog.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger.With(slog.String("service", "settings")),
	}
}

// Get retrieves a setting by key.
func (s *Service) Get(ctx context.Context, key string) (*db.ServerSetting, error) {
	setting, err := s.queries.GetServerSetting(ctx, key)
	if err != nil {
		return nil, ErrSettingNotFound
	}
	return &setting, nil
}

// GetValue retrieves and unmarshals a setting value.
func (s *Service) GetValue(ctx context.Context, key string, dest any) error {
	setting, err := s.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(setting.Value, dest)
}

// GetString retrieves a string setting value.
func (s *Service) GetString(ctx context.Context, key string) (string, error) {
	var value string
	if err := s.GetValue(ctx, key, &value); err != nil {
		return "", err
	}
	return value, nil
}

// GetBool retrieves a boolean setting value.
func (s *Service) GetBool(ctx context.Context, key string) (bool, error) {
	var value bool
	if err := s.GetValue(ctx, key, &value); err != nil {
		return false, err
	}
	return value, nil
}

// GetInt retrieves an integer setting value.
func (s *Service) GetInt(ctx context.Context, key string) (int, error) {
	var value int
	if err := s.GetValue(ctx, key, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// SetParams contains parameters for setting a value.
type SetParams struct {
	Key         string
	Value       any
	Category    string
	Description *string
	IsPublic    *bool
}

// Set creates or updates a setting.
func (s *Service) Set(ctx context.Context, params SetParams) (*db.ServerSetting, error) {
	// Marshal value to JSON
	valueJSON, err := json.Marshal(params.Value)
	if err != nil {
		return nil, err
	}

	// Default category
	category := params.Category
	if category == "" {
		category = CategoryGeneral
	}

	// Default isPublic to false
	isPublic := false
	if params.IsPublic != nil {
		isPublic = *params.IsPublic
	}

	setting, err := s.queries.UpsertServerSetting(ctx, db.UpsertServerSettingParams{
		Key:         params.Key,
		Value:       valueJSON,
		Category:    category,
		Description: params.Description,
		IsPublic:    isPublic,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("Setting updated",
		slog.String("key", params.Key),
		slog.String("category", category),
	)

	return &setting, nil
}

// Delete removes a setting.
func (s *Service) Delete(ctx context.Context, key string) error {
	if err := s.queries.DeleteServerSetting(ctx, key); err != nil {
		return err
	}

	s.logger.Info("Setting deleted", slog.String("key", key))
	return nil
}

// ListAll returns all settings.
func (s *Service) ListAll(ctx context.Context) ([]db.ServerSetting, error) {
	return s.queries.ListServerSettings(ctx)
}

// ListByCategory returns settings filtered by category.
func (s *Service) ListByCategory(ctx context.Context, category string) ([]db.ServerSetting, error) {
	return s.queries.GetServerSettingsByCategory(ctx, category)
}

// ListPublic returns all public settings (visible to non-admin users).
func (s *Service) ListPublic(ctx context.Context) ([]db.ServerSetting, error) {
	return s.queries.GetPublicServerSettings(ctx)
}

// GetServerName returns the server display name.
func (s *Service) GetServerName(ctx context.Context) (string, error) {
	return s.GetString(ctx, KeyServerName)
}

// IsRegistrationAllowed checks if user registration is enabled.
func (s *Service) IsRegistrationAllowed(ctx context.Context) (bool, error) {
	return s.GetBool(ctx, KeySecurityAllowRegistration)
}

// IsAuthRequired checks if authentication is required.
func (s *Service) IsAuthRequired(ctx context.Context) (bool, error) {
	return s.GetBool(ctx, KeySecurityRequireAuth)
}

// IsAdultContentGloballyEnabled checks if adult content is enabled server-wide.
func (s *Service) IsAdultContentGloballyEnabled(ctx context.Context) (bool, error) {
	return s.GetBool(ctx, KeyAdultGloballyEnabled)
}
