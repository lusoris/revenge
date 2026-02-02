package api

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/settings"
	"github.com/lusoris/revenge/internal/service/user"
	"go.uber.org/zap"
)

// Handler implements the ogen.Handler interface for health check endpoints.
type Handler struct {
	logger          *zap.Logger
	cfg             *config.Config
	healthService   *health.Service
	settingsService settings.Service
	userService     *user.Service
	authService     *auth.Service
	sessionService  *session.Service
	rbacService     *rbac.Service
	apikeyService   *apikeys.Service
	tokenManager    auth.TokenManager
}

// HandleBearerAuth implements the SecurityHandler interface.
// Validates JWT access tokens and injects user context.
func (h *Handler) HandleBearerAuth(ctx context.Context, operationName ogen.OperationName, t ogen.BearerAuth) (context.Context, error) {
	h.logger.Debug("Bearer auth requested", zap.String("operation", string(operationName)))

	// Validate JWT token
	claims, err := h.tokenManager.ValidateAccessToken(t.Token)
	if err != nil {
		h.logger.Warn("Invalid JWT token", zap.Error(err))
		return nil, errors.Wrap(err, "invalid token")
	}

	// Inject user data into context
	ctx = WithUserID(ctx, claims.UserID)
	ctx = WithUsername(ctx, claims.Username)

	h.logger.Debug("JWT validated successfully",
		zap.String("user_id", claims.UserID.String()),
		zap.String("username", claims.Username))

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

// getUserID retrieves the user ID from the context (convenience wrapper).
func (h *Handler) getUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, err := GetUserID(ctx)
	return userID, err == nil
}

// getSessionID retrieves the session ID from the context.
func (h *Handler) getSessionID(ctx context.Context) (uuid.UUID, bool) {
	sessionID, ok := ctx.Value(sessionIDKey).(uuid.UUID)
	return sessionID, ok
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
	// Get user ID from JWT token in context
	userID, err := GetUserID(ctx)
	if err != nil {
		h.logger.Warn("No user ID in context", zap.Error(err))
		return &ogen.Error{
			Code:    401,
			Message: "Unauthorized",
		}, nil
	}

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

// ============================================================================
// Auth Handlers
// ============================================================================

// Register handles user registration
func (h *Handler) Register(ctx context.Context, req *ogen.RegisterRequest) (ogen.RegisterRes, error) {
	h.logger.Info("User registration requested", zap.String("username", req.Username))

	// Extract display name if set
	var displayName *string
	if req.DisplayName.Set {
		displayName = &req.DisplayName.Value
	}

	// Create user via auth service
	user, err := h.authService.Register(ctx, auth.RegisterRequest{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: displayName,
	})
	if err != nil {
		h.logger.Warn("Registration failed", zap.Error(err))
		return &ogen.Error{
			Code:    400,
			Message: fmt.Sprintf("Registration failed: %v", err),
		}, nil
	}

	h.logger.Info("User registered successfully", zap.String("user_id", user.ID.String()))

	return &ogen.User{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		DisplayName:   ogen.NewOptString(stringPtrToString(user.DisplayName)),
		EmailVerified: ogen.NewOptBool(boolPtrToBool(user.EmailVerified)),
		IsActive:      boolPtrToBool(user.IsActive),
		CreatedAt:     user.CreatedAt,
	}, nil
}

// Login handles user authentication
func (h *Handler) Login(ctx context.Context, req *ogen.LoginRequest) (ogen.LoginRes, error) {
	h.logger.Info("Login requested", zap.String("username", req.Username))

	// Extract device name if set
	var deviceName *string
	if req.DeviceName.Set {
		deviceName = &req.DeviceName.Value
	}

	// Authenticate user (TODO: extract IP, user agent, fingerprint from request)
	loginResp, err := h.authService.Login(ctx, req.Username, req.Password, nil, nil, deviceName, nil)
	if err != nil {
		h.logger.Warn("Login failed", zap.Error(err), zap.String("username", req.Username))
		return &ogen.LoginUnauthorized{
			Code:    401,
			Message: "Invalid username or password",
		}, nil
	}

	h.logger.Info("Login successful", zap.String("user_id", loginResp.User.ID.String()))

	return &ogen.LoginResponse{
		User: ogen.User{
			ID:            loginResp.User.ID,
			Username:      loginResp.User.Username,
			Email:         loginResp.User.Email,
			DisplayName:   ogen.NewOptString(stringPtrToString(loginResp.User.DisplayName)),
			EmailVerified: ogen.NewOptBool(boolPtrToBool(loginResp.User.EmailVerified)),
			IsActive:      boolPtrToBool(loginResp.User.IsActive),
			CreatedAt:     loginResp.User.CreatedAt,
		},
		AccessToken:  loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
		ExpiresIn:    int(loginResp.ExpiresIn),
	}, nil
}

// Logout handles user logout
func (h *Handler) Logout(ctx context.Context, req *ogen.LogoutRequest) (ogen.LogoutRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		h.logger.Warn("Logout: no user in context", zap.Error(err))
		return &ogen.Error{Code: 401, Message: "Unauthorized"}, nil
	}

	h.logger.Info("Logout requested", zap.String("user_id", userID.String()))

	// Logout logic
	if req.LogoutAll.Value {
		// Logout from all devices
		if err := h.authService.LogoutAll(ctx, userID); err != nil {
			h.logger.Error("Logout all failed", zap.Error(err))
			return &ogen.Error{}, fmt.Errorf("logout failed: %w", err)
		}
	} else {
		// Logout from current device only
		if err := h.authService.Logout(ctx, req.RefreshToken); err != nil {
			h.logger.Error("Logout failed", zap.Error(err))
			return &ogen.Error{}, fmt.Errorf("logout failed: %w", err)
		}
	}

	return &ogen.LogoutNoContent{}, nil
}

// RefreshToken handles token refresh
func (h *Handler) RefreshToken(ctx context.Context, req *ogen.RefreshRequest) (ogen.RefreshTokenRes, error) {
	h.logger.Debug("Token refresh requested")

	// Refresh access token
	loginResp, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		h.logger.Warn("Token refresh failed", zap.Error(err))
		return &ogen.RefreshTokenUnauthorized{
			Code:    401,
			Message: "Invalid or expired refresh token",
		}, nil
	}

	return &ogen.RefreshResponse{
		AccessToken: loginResp.AccessToken,
		ExpiresIn:   int(loginResp.ExpiresIn),
	}, nil
}

// VerifyEmail handles email verification
func (h *Handler) VerifyEmail(ctx context.Context, req *ogen.VerifyEmailRequest) (ogen.VerifyEmailRes, error) {
	h.logger.Info("Email verification requested")

	// Verify email token
	if err := h.authService.VerifyEmail(ctx, req.Token); err != nil {
		h.logger.Warn("Email verification failed", zap.Error(err))
		return &ogen.Error{}, fmt.Errorf("verification failed: %w", err)
	}

	h.logger.Info("Email verified successfully")
	return &ogen.VerifyEmailNoContent{}, nil
}

// ResendVerification handles resending verification email
func (h *Handler) ResendVerification(ctx context.Context) (ogen.ResendVerificationRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		h.logger.Warn("Resend verification: no user in context", zap.Error(err))
		return &ogen.Error{Code: 401, Message: "Unauthorized"}, nil
	}

	h.logger.Info("Resend verification requested", zap.String("user_id", userID.String()))

	// Resend verification email
	if err := h.authService.ResendVerification(ctx, userID); err != nil {
		h.logger.Error("Resend verification failed", zap.Error(err))
		return &ogen.Error{}, fmt.Errorf("resend failed: %w", err)
	}

	return &ogen.ResendVerificationNoContent{}, nil
}

// ForgotPassword handles password reset request
func (h *Handler) ForgotPassword(ctx context.Context, req *ogen.ForgotPasswordRequest) (ogen.ForgotPasswordRes, error) {
	h.logger.Info("Password reset requested", zap.String("email", req.Email))

	// Request password reset (always returns success to prevent email enumeration)
	// TODO: Extract IP address and user agent from request
	_, err := h.authService.RequestPasswordReset(ctx, req.Email, nil, nil)
	if err != nil {
		h.logger.Error("Password reset request failed", zap.Error(err))
		// Still return success to avoid email enumeration
	}

	return &ogen.ForgotPasswordNoContent{}, nil
}

// ResetPassword handles password reset with token
func (h *Handler) ResetPassword(ctx context.Context, req *ogen.ResetPasswordRequest) (ogen.ResetPasswordRes, error) {
	h.logger.Info("Password reset with token")

	// Reset password using token
	if err := h.authService.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		h.logger.Warn("Password reset failed", zap.Error(err))
		return &ogen.Error{}, fmt.Errorf("reset failed: %w", err)
	}

	h.logger.Info("Password reset successfully")
	return &ogen.ResetPasswordNoContent{}, nil
}

// ChangePassword handles password change for authenticated user
func (h *Handler) ChangePassword(ctx context.Context, req *ogen.ChangePasswordRequest) (ogen.ChangePasswordRes, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		h.logger.Warn("Change password: no user in context", zap.Error(err))
		return &ogen.ChangePasswordUnauthorized{
			Code:    401,
			Message: "Unauthorized",
		}, nil
	}

	h.logger.Info("Password change requested", zap.String("user_id", userID.String()))

	// Change password
	if err := h.authService.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword); err != nil {
		h.logger.Warn("Password change failed", zap.Error(err))
		return &ogen.ChangePasswordBadRequest{
			Code:    400,
			Message: fmt.Sprintf("Password change failed: %v", err),
		}, nil
	}

	h.logger.Info("Password changed successfully")
	return &ogen.ChangePasswordNoContent{}, nil
}
