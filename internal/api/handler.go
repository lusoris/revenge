package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/service/settings"
	"go.uber.org/zap"
)

// Handler implements the ogen.Handler interface for health check endpoints.
type Handler struct {
	logger          *zap.Logger
	healthService   *health.Service
	settingsService settings.Service
}

// HandleBearerAuth implements the SecurityHandler interface.
// TODO: Implement JWT verification
func (h *Handler) HandleBearerAuth(ctx context.Context, operationName ogen.OperationName, t ogen.BearerAuth) (context.Context, error) {
	h.logger.Debug("Bearer auth requested", zap.String("operation", string(operationName)))
	// TODO: Verify JWT token and extract user ID
	// For now, just accept any token
	return ctx, nil
}

// GetLiveness implements the liveness probe endpoint.
// This always returns healthy unless the process is deadlocked.
func (h *Handler) GetLiveness(ctx context.Context) (*ogen.HealthCheck, error) {
	h.logger.Debug("Liveness check requested")

	return &ogen.HealthCheck{
		Name:    "liveness",
		Status:  ogen.HealthCheckStatusHealthy,
		Message: ogen.NewOptString("Service is alive"),
	}, nil
}

// GetReadiness implements the readiness probe endpoint.
// Returns healthy only if all dependencies are available.
func (h *Handler) GetReadiness(ctx context.Context) (ogen.GetReadinessRes, error) {
	h.logger.Debug("Readiness check requested")

	// Check if service is ready
	result := h.healthService.Readiness(ctx)

	healthCheck := &ogen.HealthCheck{
		Name: result.Name,
	}

	if result.Status == health.StatusHealthy {
		healthCheck.Status = ogen.HealthCheckStatusHealthy
		healthCheck.Message = ogen.NewOptString(result.Message)
		return (*ogen.GetReadinessOK)(healthCheck), nil
	}

	healthCheck.Status = ogen.HealthCheckStatusUnhealthy
	healthCheck.Message = ogen.NewOptString(result.Message)
	return (*ogen.GetReadinessServiceUnavailable)(healthCheck), nil
}

// GetStartup implements the startup probe endpoint.
// Returns healthy only after initialization is complete.
func (h *Handler) GetStartup(ctx context.Context) (ogen.GetStartupRes, error) {
	h.logger.Debug("Startup check requested")

	// Check if service has started
	result := h.healthService.Startup(ctx)

	healthCheck := &ogen.HealthCheck{
		Name: result.Name,
	}

	if result.Status == health.StatusHealthy {
		healthCheck.Status = ogen.HealthCheckStatusHealthy
		healthCheck.Message = ogen.NewOptString(result.Message)
		return (*ogen.GetStartupOK)(healthCheck), nil
	}

	healthCheck.Status = ogen.HealthCheckStatusUnhealthy
	healthCheck.Message = ogen.NewOptString(result.Message)
	return (*ogen.GetStartupServiceUnavailable)(healthCheck), nil
}

// NewError creates an error response for failed requests.
func (h *Handler) NewError(ctx context.Context, err error) *ogen.ErrorStatusCode {
	h.logger.Error("Request error", zap.Error(err))

	return &ogen.ErrorStatusCode{
		StatusCode: 500,
		Response: ogen.Error{
			Code:    500,
			Message: fmt.Sprintf("Internal server error: %v", err),
		},
	}
}

// ============================================================================
// Server Settings Handlers
// ============================================================================

// ListServerSettings retrieves all server settings.
func (h *Handler) ListServerSettings(ctx context.Context) (ogen.ListServerSettingsRes, error) {
	h.logger.Debug("List server settings requested")

	settingsList, err := h.settingsService.ListServerSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list server settings: %w", err)
	}

	result := make([]ogen.ServerSetting, len(settingsList))
	for i, s := range settingsList {
		result[i] = toOgenServerSetting(&s)
	}

	return (*ogen.ListServerSettingsOKApplicationJSON)(&result), nil
}

// GetServerSetting retrieves a specific server setting.
func (h *Handler) GetServerSetting(ctx context.Context, params ogen.GetServerSettingParams) (ogen.GetServerSettingRes, error) {
	h.logger.Debug("Get server setting requested", zap.String("key", params.Key))

	setting, err := h.settingsService.GetServerSetting(ctx, params.Key)
	if err != nil {
		return &ogen.GetServerSettingNotFound{}, nil
	}

	result := toOgenServerSetting(setting)
	return &result, nil
}

// UpdateServerSetting updates a server setting value.
func (h *Handler) UpdateServerSetting(ctx context.Context, req *ogen.SettingValue, params ogen.UpdateServerSettingParams) (ogen.UpdateServerSettingRes, error) {
	h.logger.Debug("Update server setting requested", zap.String("key", params.Key))

	// TODO: Get user ID from auth context
	userID := uuid.New() // Placeholder

	setting, err := h.settingsService.SetServerSetting(ctx, params.Key, req.Value, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update server setting: %w", err)
	}

	result := toOgenServerSetting(setting)
	return &result, nil
}

// ============================================================================
// User Settings Handlers
// ============================================================================

// ListUserSettings retrieves all settings for the current user.
func (h *Handler) ListUserSettings(ctx context.Context) (ogen.ListUserSettingsRes, error) {
	h.logger.Debug("List user settings requested")

	// TODO: Get user ID from auth context
	userID := uuid.New() // Placeholder

	settingsList, err := h.settingsService.ListUserSettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user settings: %w", err)
	}

	result := make([]ogen.UserSetting, len(settingsList))
	for i, s := range settingsList {
		result[i] = toOgenUserSetting(&s)
	}

	return (*ogen.ListUserSettingsOKApplicationJSON)(&result), nil
}

// GetUserSetting retrieves a specific user setting.
func (h *Handler) GetUserSetting(ctx context.Context, params ogen.GetUserSettingParams) (ogen.GetUserSettingRes, error) {
	h.logger.Debug("Get user setting requested", zap.String("key", params.Key))

	// TODO: Get user ID from auth context
	userID := uuid.New() // Placeholder

	setting, err := h.settingsService.GetUserSetting(ctx, userID, params.Key)
	if err != nil {
		return &ogen.GetUserSettingNotFound{}, nil
	}

	result := toOgenUserSetting(setting)
	return &result, nil
}

// UpdateUserSetting updates a user setting value.
func (h *Handler) UpdateUserSetting(ctx context.Context, req *ogen.SettingValue, params ogen.UpdateUserSettingParams) (ogen.UpdateUserSettingRes, error) {
	h.logger.Debug("Update user setting requested", zap.String("key", params.Key))

	// TODO: Get user ID from auth context
	userID := uuid.New() // Placeholder

	setting, err := h.settingsService.SetUserSetting(ctx, userID, params.Key, req.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to update user setting: %w", err)
	}

	result := toOgenUserSetting(setting)
	return &result, nil
}

// DeleteUserSetting deletes a user setting.
func (h *Handler) DeleteUserSetting(ctx context.Context, params ogen.DeleteUserSettingParams) (ogen.DeleteUserSettingRes, error) {
	h.logger.Debug("Delete user setting requested", zap.String("key", params.Key))

	// TODO: Get user ID from auth context
	userID := uuid.New() // Placeholder

	if err := h.settingsService.DeleteUserSetting(ctx, userID, params.Key); err != nil {
		return &ogen.DeleteUserSettingNotFound{}, nil
	}

	return &ogen.DeleteUserSettingNoContent{}, nil
}

// ============================================================================
// Helper Functions
// ============================================================================

func toOgenServerSetting(s *settings.ServerSetting) ogen.ServerSetting {
	// Convert value to ogen union type
	// For complex types (objects), serialize to string
	var value ogen.ServerSettingValue
	switch v := s.Value.(type) {
	case string:
		value = ogen.NewStringServerSettingValue(v)
	case float64:
		value = ogen.NewFloat64ServerSettingValue(v)
	case bool:
		value = ogen.NewBoolServerSettingValue(v)
	default:
		// For objects and other types, convert to string
		value = ogen.NewStringServerSettingValue(fmt.Sprintf("%v", v))
	}

	return ogen.ServerSetting{
		Key:         s.Key,
		Value:       value,
		Description: ogen.NewOptString(stringPtrToString(s.Description)),
		Category:    ogen.NewOptString(stringPtrToString(s.Category)),
		DataType:    ogen.ServerSettingDataType(s.DataType),
		IsSecret:    ogen.NewOptBool(boolPtrToBool(s.IsSecret)),
		IsPublic:    ogen.NewOptBool(boolPtrToBool(s.IsPublic)),
	}
}

func toOgenUserSetting(s *settings.UserSetting) ogen.UserSetting {
	// Convert value to ogen union type
	// For complex types (objects), serialize to string
	var value ogen.UserSettingValue
	switch v := s.Value.(type) {
	case string:
		value = ogen.NewStringUserSettingValue(v)
	case float64:
		value = ogen.NewFloat64UserSettingValue(v)
	case bool:
		value = ogen.NewBoolUserSettingValue(v)
	default:
		// For objects and other types, convert to string
		value = ogen.NewStringUserSettingValue(fmt.Sprintf("%v", v))
	}

	return ogen.UserSetting{
		UserID:      s.UserID,
		Key:         s.Key,
		Value:       value,
		Description: ogen.NewOptString(stringPtrToString(s.Description)),
		Category:    ogen.NewOptString(stringPtrToString(s.Category)),
		DataType:    ogen.UserSettingDataType(s.DataType),
	}
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func boolPtrToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
