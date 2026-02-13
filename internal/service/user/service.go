package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/netip"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/storage"
)

// Service implements user business logic
type Service struct {
	pool           *pgxpool.Pool
	repo           Repository
	hasher         *crypto.PasswordHasher
	activityLogger activity.Logger
	storage        storage.Storage
	avatarConfig   config.AvatarConfig
}

// NewService creates a new user service
func NewService(pool *pgxpool.Pool, repo Repository, activityLogger activity.Logger, store storage.Storage, avatarCfg config.AvatarConfig) *Service {
	return &Service{
		pool:           pool,
		repo:           repo,
		hasher:         crypto.NewPasswordHasher(),
		activityLogger: activityLogger,
		storage:        store,
		avatarConfig:   avatarCfg,
	}
}

// ============================================================================
// User Management
// ============================================================================

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, userID uuid.UUID) (*db.SharedUser, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// GetUserByUsername retrieves a user by username
func (s *Service) GetUserByUsername(ctx context.Context, username string) (*db.SharedUser, error) {
	return s.repo.GetUserByUsername(ctx, username)
}

// GetUserByEmail retrieves a user by email
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*db.SharedUser, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

// ListUsers retrieves users with optional filters
func (s *Service) ListUsers(ctx context.Context, filters UserFilters) ([]db.SharedUser, int64, error) {
	return s.repo.ListUsers(ctx, filters)
}

// CreateUser creates a new user with hashed password
func (s *Service) CreateUser(ctx context.Context, params CreateUserParams) (*db.SharedUser, error) {
	// Validate input
	if params.Username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if params.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if params.PasswordHash == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Check if username already exists
	existing, err := s.repo.GetUserByUsername(ctx, params.Username)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Check if email already exists
	existing, err = s.repo.GetUserByEmail(ctx, params.Email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password if not already hashed
	hashedPassword, err := s.HashPassword(params.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	params.PasswordHash = hashedPassword

	// Create user
	user, err := s.repo.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create default preferences
	defaultPrefs := s.getDefaultPreferences(user.ID)
	_, err = s.repo.UpsertUserPreferences(ctx, defaultPrefs)
	if err != nil {
		// Log error but don't fail user creation
		_ = err
	}

	// Log user creation
	_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
		UserID:       user.ID,
		Username:     user.Username,
		Action:       activity.ActionUserCreate,
		ResourceType: activity.ResourceTypeUser,
		ResourceID:   user.ID,
	})

	return user, nil
}

// UpdateUser updates user information
func (s *Service) UpdateUser(ctx context.Context, userID uuid.UUID, params UpdateUserParams) (*db.SharedUser, error) {
	// Verify user exists
	oldUser, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update user
	updatedUser, err := s.repo.UpdateUser(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	// Build changes map
	changes := make(map[string]any)
	if params.DisplayName != nil && (oldUser.DisplayName == nil || *params.DisplayName != *oldUser.DisplayName) {
		changes["display_name"] = map[string]any{"old": oldUser.DisplayName, "new": *params.DisplayName}
	}
	if params.Email != nil && *params.Email != oldUser.Email {
		changes["email"] = map[string]any{"old": oldUser.Email, "new": *params.Email}
	}

	if len(changes) > 0 {
		_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
			UserID:       userID,
			Username:     updatedUser.Username,
			Action:       activity.ActionUserUpdate,
			ResourceType: activity.ResourceTypeUser,
			ResourceID:   userID,
			Changes:      changes,
		})
	}

	return updatedUser, nil
}

// UpdatePassword updates a user's password
func (s *Service) UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify old password
	if err := s.VerifyPassword(user.PasswordHash, oldPassword); err != nil {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	return s.repo.UpdatePassword(ctx, userID, hashedPassword)
}

// DeleteUser soft deletes a user
func (s *Service) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Get user for logging
	user, _ := s.repo.GetUserByID(ctx, userID)

	err := s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	// Log user deletion
	if user != nil {
		_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
			UserID:       userID,
			Username:     user.Username,
			Action:       activity.ActionUserDelete,
			ResourceType: activity.ResourceTypeUser,
			ResourceID:   userID,
		})
	}

	return nil
}

// HardDeleteUser permanently deletes a user (GDPR compliance)
func (s *Service) HardDeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Delete user preferences
	_ = s.repo.DeleteUserPreferences(ctx, userID)

	// Hard delete user (cascades to avatars)
	return s.repo.HardDeleteUser(ctx, userID)
}

// VerifyEmail marks a user's email as verified
func (s *Service) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	return s.repo.VerifyEmail(ctx, userID)
}

// RecordLogin updates the user's last login timestamp
func (s *Service) RecordLogin(ctx context.Context, userID uuid.UUID) error {
	return s.repo.UpdateLastLogin(ctx, userID)
}

// ============================================================================
// Password Management
// ============================================================================

// HashPassword hashes a password using Argon2id
func (s *Service) HashPassword(password string) (string, error) {
	return s.hasher.HashPassword(password)
}

// VerifyPassword verifies a password against a hash
func (s *Service) VerifyPassword(hashedPassword, password string) error {
	match, err := s.hasher.VerifyPassword(password, hashedPassword)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("password verification failed")
	}
	return nil
}

// ============================================================================
// User Preferences
// ============================================================================

// GetUserPreferences retrieves user preferences
func (s *Service) GetUserPreferences(ctx context.Context, userID uuid.UUID) (*db.SharedUserPreference, error) {
	prefs, err := s.repo.GetUserPreferences(ctx, userID)
	if err != nil {
		// Return default preferences if not found
		defaultPrefs := s.getDefaultPreferences(userID)
		return s.repo.UpsertUserPreferences(ctx, defaultPrefs)
	}
	return prefs, nil
}

// UpdateUserPreferences updates user preferences
func (s *Service) UpdateUserPreferences(ctx context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error) {
	// Validate preferences
	if err := s.validatePreferences(params); err != nil {
		return nil, err
	}

	return s.repo.UpsertUserPreferences(ctx, params)
}

// UpdateNotificationPreferences updates notification settings
func (s *Service) UpdateNotificationPreferences(ctx context.Context, userID uuid.UUID, email, push, digest *NotificationSettings) error {
	params := UpsertPreferencesParams{
		UserID: userID,
	}

	if email != nil {
		emailJSON, err := json.Marshal(email)
		if err != nil {
			return fmt.Errorf("failed to marshal email notifications: %w", err)
		}
		raw := json.RawMessage(emailJSON)
		params.EmailNotifications = &raw
	}

	if push != nil {
		pushJSON, err := json.Marshal(push)
		if err != nil {
			return fmt.Errorf("failed to marshal push notifications: %w", err)
		}
		raw := json.RawMessage(pushJSON)
		params.PushNotifications = &raw
	}

	if digest != nil {
		digestJSON, err := json.Marshal(digest)
		if err != nil {
			return fmt.Errorf("failed to marshal digest notifications: %w", err)
		}
		raw := json.RawMessage(digestJSON)
		params.DigestNotifications = &raw
	}

	_, err := s.repo.UpsertUserPreferences(ctx, params)
	return err
}

// ============================================================================
// Avatar Management
// ============================================================================

// GetCurrentAvatar retrieves the user's current avatar
func (s *Service) GetCurrentAvatar(ctx context.Context, userID uuid.UUID) (*db.SharedUserAvatar, error) {
	return s.repo.GetCurrentAvatar(ctx, userID)
}

// ListUserAvatars retrieves avatar history for a user
func (s *Service) ListUserAvatars(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.SharedUserAvatar, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.ListUserAvatars(ctx, userID, limit, offset)
}

// UploadAvatar uploads a new avatar and sets it as current
func (s *Service) UploadAvatar(ctx context.Context, userID uuid.UUID, file io.Reader, metadata AvatarMetadata) (*db.SharedUserAvatar, error) {
	// Validate metadata
	if err := s.validateAvatarMetadata(metadata); err != nil {
		return nil, err
	}

	// Generate storage key
	storageKey := storage.GenerateAvatarKey(userID, metadata.FileName)

	// Store the file first (outside transaction to avoid blocking DB)
	storedKey, err := s.storage.Store(ctx, storageKey, file, metadata.MimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to store avatar: %w", err)
	}

	// Get the URL for the stored file
	filePath := s.storage.GetURL(storedKey)

	// Begin transaction to ensure atomicity of all DB operations
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		_ = s.storage.Delete(ctx, storedKey)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Create transaction-scoped queries
	txQueries := db.New(tx)

	// Get next version number (within transaction)
	latestVersion, err := txQueries.GetLatestAvatarVersion(ctx, userID)
	if err != nil {
		_ = s.storage.Delete(ctx, storedKey)
		return nil, fmt.Errorf("failed to get latest avatar version: %w", err)
	}
	nextVersion := latestVersion + 1

	// Unset current avatars (within transaction)
	if err := txQueries.UnsetCurrentAvatars(ctx, userID); err != nil {
		_ = s.storage.Delete(ctx, storedKey)
		return nil, fmt.Errorf("failed to unset current avatars: %w", err)
	}

	// Prepare CreateAvatar parameters with IP address parsing
	createParams := db.CreateAvatarParams{
		UserID:                userID,
		FilePath:              filePath,
		FileSizeBytes:         metadata.FileSizeBytes,
		MimeType:              metadata.MimeType,
		Width:                 metadata.Width,
		Height:                metadata.Height,
		IsAnimated:            metadata.IsAnimated,
		Version:               nextVersion,
		UploadedFromUserAgent: metadata.UploadedFromUserAgent,
	}

	// Parse IP address if provided
	if metadata.UploadedFromIP != nil {
		addr, err := netip.ParseAddr(*metadata.UploadedFromIP)
		if err != nil {
			_ = s.storage.Delete(ctx, storedKey)
			return nil, fmt.Errorf("failed to parse IP address: %w", err)
		}
		createParams.UploadedFromIp = addr
	}

	// Create new avatar record (within transaction)
	avatar, err := txQueries.CreateAvatar(ctx, createParams)
	if err != nil {
		_ = s.storage.Delete(ctx, storedKey)
		return nil, fmt.Errorf("failed to create avatar: %w", err)
	}

	// Update user avatar_url (within transaction)
	// All operations must succeed together or all fail
	avatarURL := filePath
	_, err = txQueries.UpdateUser(ctx, db.UpdateUserParams{
		UserID:    userID,
		AvatarUrl: &avatarURL,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, storedKey)
		return nil, fmt.Errorf("failed to update user avatar: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		_ = s.storage.Delete(ctx, storedKey)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &avatar, nil
}

// SetCurrentAvatar sets an existing avatar as current
func (s *Service) SetCurrentAvatar(ctx context.Context, userID, avatarID uuid.UUID) error {
	// Verify avatar belongs to user
	avatar, err := s.repo.GetAvatarByID(ctx, avatarID)
	if err != nil {
		return fmt.Errorf("avatar not found: %w", err)
	}
	if avatar.UserID != userID {
		return fmt.Errorf("avatar does not belong to user")
	}

	// Unset current avatars
	if err := s.repo.UnsetCurrentAvatars(ctx, userID); err != nil {
		return fmt.Errorf("failed to unset current avatars: %w", err)
	}

	// Set as current
	if err := s.repo.SetCurrentAvatar(ctx, avatarID); err != nil {
		return fmt.Errorf("failed to set current avatar: %w", err)
	}

	// Update user avatar_url
	_, err = s.repo.UpdateUser(ctx, userID, UpdateUserParams{
		AvatarURL: &avatar.FilePath,
	})

	return err
}

// DeleteAvatar soft deletes an avatar
func (s *Service) DeleteAvatar(ctx context.Context, userID, avatarID uuid.UUID) error {
	// Verify avatar belongs to user
	avatar, err := s.repo.GetAvatarByID(ctx, avatarID)
	if err != nil {
		return fmt.Errorf("avatar not found: %w", err)
	}
	if avatar.UserID != userID {
		return fmt.Errorf("avatar does not belong to user")
	}

	return s.repo.DeleteAvatar(ctx, avatarID)
}

// ============================================================================
// Validation and Helpers
// ============================================================================

// AvatarMetadata contains metadata for avatar upload
type AvatarMetadata struct {
	FileName              string
	FileSizeBytes         int64
	MimeType              string
	Width                 int32
	Height                int32
	IsAnimated            *bool
	UploadedFromIP        *string
	UploadedFromUserAgent *string
}

func (s *Service) getDefaultPreferences(userID uuid.UUID) UpsertPreferencesParams {
	emailNotif := json.RawMessage(`{"enabled": true, "frequency": "instant"}`)
	pushNotif := json.RawMessage(`{"enabled": false}`)
	digestNotif := json.RawMessage(`{"enabled": true, "frequency": "weekly"}`)

	profileVis := "private"
	showEmail := false
	showActivity := true
	theme := "system"
	displayLang := "en-US"
	showAdult := false
	showSpoilers := false
	autoPlay := true

	return UpsertPreferencesParams{
		UserID:              userID,
		EmailNotifications:  &emailNotif,
		PushNotifications:   &pushNotif,
		DigestNotifications: &digestNotif,
		ProfileVisibility:   &profileVis,
		ShowEmail:           &showEmail,
		ShowActivity:        &showActivity,
		Theme:               &theme,
		DisplayLanguage:     &displayLang,
		ShowAdultContent:    &showAdult,
		ShowSpoilers:        &showSpoilers,
		AutoPlayVideos:      &autoPlay,
	}
}

func (s *Service) validatePreferences(params UpsertPreferencesParams) error {
	// Validate profile visibility
	if params.ProfileVisibility != nil {
		allowed := map[string]bool{"public": true, "friends": true, "private": true}
		if !allowed[*params.ProfileVisibility] {
			return fmt.Errorf("invalid profile visibility: must be public, friends, or private")
		}
	}

	// Validate theme
	if params.Theme != nil {
		allowed := map[string]bool{"light": true, "dark": true, "system": true}
		if !allowed[*params.Theme] {
			return fmt.Errorf("invalid theme: must be light, dark, or system")
		}
	}

	return nil
}

func (s *Service) validateAvatarMetadata(metadata AvatarMetadata) error {
	// Max 5MB
	if metadata.FileSizeBytes > 5*1024*1024 {
		return fmt.Errorf("avatar size exceeds maximum (5MB)")
	}

	// Valid MIME types
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedMimeTypes[metadata.MimeType] {
		return fmt.Errorf("invalid MIME type: must be JPEG, PNG, GIF, or WebP")
	}

	// Reasonable dimensions
	if metadata.Width < 16 || metadata.Width > 4096 {
		return fmt.Errorf("invalid width: must be between 16 and 4096 pixels")
	}
	if metadata.Height < 16 || metadata.Height > 4096 {
		return fmt.Errorf("invalid height: must be between 16 and 4096 pixels")
	}

	return nil
}
