package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/service/settings"
	"github.com/lusoris/revenge/internal/service/user"
	"go.uber.org/zap"
)

// Handler implements the ogen.Handler interface for health check endpoints.
type Handler struct {
	logger          *zap.Logger
	healthService   *health.Service
	settingsService settings.Service
	userService     *user.Service
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

// ============================================================================
// User Endpoints
// ============================================================================

// GetCurrentUser returns the authenticated user's profile
func (h *Handler) GetCurrentUser(ctx context.Context) (ogen.GetCurrentUserRes, error) {
	// TODO: Get user ID from JWT token in context
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000") // Placeholder

	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		return &ogen.Error{}, fmt.Errorf("user not found: %w", err)
	}

	return &ogen.User{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		DisplayName:   ogen.NewOptString(stringPtrToString(user.DisplayName)),
		AvatarURL:     ogen.NewOptString(stringPtrToString(user.AvatarUrl)),
		Locale:        ogen.NewOptString(stringPtrToString(user.Locale)),
		Timezone:      ogen.NewOptString(stringPtrToString(user.Timezone)),
		QarEnabled:    ogen.NewOptBool(boolPtrToBool(user.QarEnabled)),
		IsActive:      boolPtrToBool(user.IsActive),
		IsAdmin:       ogen.NewOptBool(boolPtrToBool(user.IsAdmin)),
		EmailVerified: ogen.NewOptBool(boolPtrToBool(user.EmailVerified)),
		CreatedAt:     user.CreatedAt,
		LastLoginAt:   ogen.NewOptDateTime(user.LastLoginAt.Time),
	}, nil
}

// UpdateCurrentUser updates the authenticated user's profile
func (h *Handler) UpdateCurrentUser(ctx context.Context, req *ogen.UserUpdate) (ogen.UpdateCurrentUserRes, error) {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000") // Placeholder

	params := user.UpdateUserParams{}
	if email, ok := req.Email.Get(); ok {
		params.Email = &email
	}
	if displayName, ok := req.DisplayName.Get(); ok {
		params.DisplayName = &displayName
	}
	if timezone, ok := req.Timezone.Get(); ok {
		params.Timezone = &timezone
	}

	updatedUser, err := h.userService.UpdateUser(ctx, userID, params)
	if err != nil {
		return &ogen.UpdateCurrentUserBadRequest{}, fmt.Errorf("failed to update user: %w", err)
	}

	return &ogen.User{
		ID:            updatedUser.ID,
		Username:      updatedUser.Username,
		Email:         updatedUser.Email,
		DisplayName:   ogen.NewOptString(stringPtrToString(updatedUser.DisplayName)),
		AvatarURL:     ogen.NewOptString(stringPtrToString(updatedUser.AvatarUrl)),
		Locale:        ogen.NewOptString(stringPtrToString(updatedUser.Locale)),
		Timezone:      ogen.NewOptString(stringPtrToString(updatedUser.Timezone)),
		QarEnabled:    ogen.NewOptBool(boolPtrToBool(updatedUser.QarEnabled)),
		IsActive:      boolPtrToBool(updatedUser.IsActive),
		IsAdmin:       ogen.NewOptBool(boolPtrToBool(updatedUser.IsAdmin)),
		EmailVerified: ogen.NewOptBool(boolPtrToBool(updatedUser.EmailVerified)),
		CreatedAt:     updatedUser.CreatedAt,
		LastLoginAt:   ogen.NewOptDateTime(updatedUser.LastLoginAt.Time),
	}, nil
}

// GetUserById returns a user's public profile
func (h *Handler) GetUserById(ctx context.Context, params ogen.GetUserByIdParams) (ogen.GetUserByIdRes, error) {
	user, err := h.userService.GetUser(ctx, params.UserId)
	if err != nil {
		return &ogen.GetUserByIdNotFound{}, fmt.Errorf("user not found: %w", err)
	}

	return &ogen.User{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		DisplayName:   ogen.NewOptString(stringPtrToString(user.DisplayName)),
		AvatarURL:     ogen.NewOptString(stringPtrToString(user.AvatarUrl)),
		Locale:        ogen.NewOptString(stringPtrToString(user.Locale)),
		Timezone:      ogen.NewOptString(stringPtrToString(user.Timezone)),
		QarEnabled:    ogen.NewOptBool(boolPtrToBool(user.QarEnabled)),
		IsActive:      boolPtrToBool(user.IsActive),
		IsAdmin:       ogen.NewOptBool(boolPtrToBool(user.IsAdmin)),
		EmailVerified: ogen.NewOptBool(boolPtrToBool(user.EmailVerified)),
		CreatedAt:     user.CreatedAt,
		LastLoginAt:   ogen.NewOptDateTime(user.LastLoginAt.Time),
	}, nil
}

// GetUserPreferences returns user preferences
func (h *Handler) GetUserPreferences(ctx context.Context) (ogen.GetUserPreferencesRes, error) {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000") // Placeholder

	prefs, err := h.userService.GetUserPreferences(ctx, userID)
	if err != nil {
		return &ogen.Error{}, fmt.Errorf("failed to get preferences: %w", err)
	}

	return &ogen.UserPreferences{
		UserID: prefs.UserID,
		ProfileVisibility: ogen.NewOptUserPreferencesProfileVisibility(
			ogen.UserPreferencesProfileVisibility(stringPtrToString(prefs.ProfileVisibility)),
		),
		ShowEmail:        ogen.NewOptBool(boolPtrToBool(prefs.ShowEmail)),
		ShowActivity:     ogen.NewOptBool(boolPtrToBool(prefs.ShowActivity)),
		Theme:            ogen.NewOptUserPreferencesTheme(ogen.UserPreferencesTheme(stringPtrToString(prefs.Theme))),
		DisplayLanguage:  ogen.NewOptString(stringPtrToString(prefs.DisplayLanguage)),
		ContentLanguage:  ogen.NewOptString(stringPtrToString(prefs.ContentLanguage)),
		ShowAdultContent: ogen.NewOptBool(boolPtrToBool(prefs.ShowAdultContent)),
		ShowSpoilers:     ogen.NewOptBool(boolPtrToBool(prefs.ShowSpoilers)),
		AutoPlayVideos:   ogen.NewOptBool(boolPtrToBool(prefs.AutoPlayVideos)),
	}, nil
}

// UpdateUserPreferences updates user preferences
func (h *Handler) UpdateUserPreferences(ctx context.Context, req *ogen.UserPreferencesUpdate) (ogen.UpdateUserPreferencesRes, error) {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000") // Placeholder

	params := user.UpsertPreferencesParams{
		UserID: userID,
	}

	// TODO: Handle notification settings (JSONB fields)

	if vis, ok := req.ProfileVisibility.Get(); ok {
		v := string(vis)
		params.ProfileVisibility = &v
	}
	if showEmail, ok := req.ShowEmail.Get(); ok {
		params.ShowEmail = &showEmail
	}
	if showActivity, ok := req.ShowActivity.Get(); ok {
		params.ShowActivity = &showActivity
	}
	if theme, ok := req.Theme.Get(); ok {
		t := string(theme)
		params.Theme = &t
	}
	if lang, ok := req.DisplayLanguage.Get(); ok {
		params.DisplayLanguage = &lang
	}
	if contentLang, ok := req.ContentLanguage.Get(); ok {
		params.ContentLanguage = &contentLang
	}
	if showAdult, ok := req.ShowAdultContent.Get(); ok {
		params.ShowAdultContent = &showAdult
	}
	if showSpoilers, ok := req.ShowSpoilers.Get(); ok {
		params.ShowSpoilers = &showSpoilers
	}
	if autoPlay, ok := req.AutoPlayVideos.Get(); ok {
		params.AutoPlayVideos = &autoPlay
	}

	prefs, err := h.userService.UpdateUserPreferences(ctx, params)
	if err != nil {
		return &ogen.UpdateUserPreferencesBadRequest{}, fmt.Errorf("failed to update preferences: %w", err)
	}

	return &ogen.UserPreferences{
		UserID: prefs.UserID,
		ProfileVisibility: ogen.NewOptUserPreferencesProfileVisibility(
			ogen.UserPreferencesProfileVisibility(stringPtrToString(prefs.ProfileVisibility)),
		),
		ShowEmail:        ogen.NewOptBool(boolPtrToBool(prefs.ShowEmail)),
		ShowActivity:     ogen.NewOptBool(boolPtrToBool(prefs.ShowActivity)),
		Theme:            ogen.NewOptUserPreferencesTheme(ogen.UserPreferencesTheme(stringPtrToString(prefs.Theme))),
		DisplayLanguage:  ogen.NewOptString(stringPtrToString(prefs.DisplayLanguage)),
		ContentLanguage:  ogen.NewOptString(stringPtrToString(prefs.ContentLanguage)),
		ShowAdultContent: ogen.NewOptBool(boolPtrToBool(prefs.ShowAdultContent)),
		ShowSpoilers:     ogen.NewOptBool(boolPtrToBool(prefs.ShowSpoilers)),
		AutoPlayVideos:   ogen.NewOptBool(boolPtrToBool(prefs.AutoPlayVideos)),
	}, nil
}

// UploadAvatar handles avatar upload
func (h *Handler) UploadAvatar(ctx context.Context, req *ogen.UploadAvatarReq) (ogen.UploadAvatarRes, error) {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000") // Placeholder

	// TODO: Parse multipart form and get file metadata
	// For now, return a placeholder response
	h.logger.Info("Avatar upload requested", zap.String("user_id", userID.String()))

	return &ogen.UploadAvatarBadRequest{}, fmt.Errorf("avatar upload not yet implemented")
}
