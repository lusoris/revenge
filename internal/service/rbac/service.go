// Package rbac provides role-based access control services.
package rbac

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

var (
	// ErrPermissionDenied indicates the user lacks required permissions.
	ErrPermissionDenied = errors.New("permission denied")
	// ErrUserNotFound indicates the user was not found.
	ErrUserNotFound = errors.New("user not found")
)

// Permission constants for common permission names.
const (
	// System permissions
	PermSystemSettingsRead  = "system.settings.read"
	PermSystemSettingsWrite = "system.settings.write"
	PermSystemLogsRead      = "system.logs.read"
	PermSystemJobsRead      = "system.jobs.read"
	PermSystemJobsManage    = "system.jobs.manage"
	PermSystemAPIKeysManage = "system.apikeys.manage"

	// User permissions
	PermUsersRead           = "users.read"
	PermUsersCreate         = "users.create"
	PermUsersUpdate         = "users.update"
	PermUsersDelete         = "users.delete"
	PermUsersSessionsManage = "users.sessions.manage"

	// Library permissions
	PermLibrariesRead   = "libraries.read"
	PermLibrariesCreate = "libraries.create"
	PermLibrariesUpdate = "libraries.update"
	PermLibrariesDelete = "libraries.delete"
	PermLibrariesScan   = "libraries.scan"

	// Content permissions
	PermContentBrowse        = "content.browse"
	PermContentMetadataRead  = "content.metadata.read"
	PermContentMetadataWrite = "content.metadata.write"
	PermContentImagesManage  = "content.images.manage"
	PermContentDelete        = "content.delete"

	// Playback permissions
	PermPlaybackStream    = "playback.stream"
	PermPlaybackDownload  = "playback.download"
	PermPlaybackTranscode = "playback.transcode"

	// Social permissions
	PermSocialRate              = "social.rate"
	PermSocialPlaylistsCreate   = "social.playlists.create"
	PermSocialPlaylistsManage   = "social.playlists.manage"
	PermSocialCollectionsCreate = "social.collections.create"
	PermSocialCollectionsManage = "social.collections.manage"
	PermSocialHistoryRead       = "social.history.read"
	PermSocialFavoritesManage   = "social.favorites.manage"

	// Adult content permissions
	PermAdultBrowse        = "adult.browse"
	PermAdultStream        = "adult.stream"
	PermAdultMetadataWrite = "adult.metadata.write"
)

// Service provides RBAC operations.
type Service struct {
	queries *db.Queries
	logger  *slog.Logger
}

// NewService creates a new RBAC service.
func NewService(queries *db.Queries, logger *slog.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger.With(slog.String("service", "rbac")),
	}
}

// HasPermission checks if a user has a specific permission.
func (s *Service) HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	has, err := s.queries.UserHasPermission(ctx, db.UserHasPermissionParams{
		ID:   userID,
		Name: permission,
	})
	if err != nil {
		return false, err
	}
	return has, nil
}

// HasAnyPermission checks if a user has any of the specified permissions.
func (s *Service) HasAnyPermission(ctx context.Context, userID uuid.UUID, permissions []string) (bool, error) {
	has, err := s.queries.UserHasAnyPermission(ctx, db.UserHasAnyPermissionParams{
		ID:      userID,
		Column2: permissions,
	})
	if err != nil {
		return false, err
	}
	return has, nil
}

// RequirePermission checks if a user has the required permission and returns an error if not.
func (s *Service) RequirePermission(ctx context.Context, userID uuid.UUID, permission string) error {
	has, err := s.HasPermission(ctx, userID, permission)
	if err != nil {
		return err
	}
	if !has {
		s.logger.Debug("Permission denied",
			slog.String("user_id", userID.String()),
			slog.String("permission", permission),
		)
		return ErrPermissionDenied
	}
	return nil
}

// RequireAnyPermission checks if a user has any of the required permissions.
func (s *Service) RequireAnyPermission(ctx context.Context, userID uuid.UUID, permissions []string) error {
	has, err := s.HasAnyPermission(ctx, userID, permissions)
	if err != nil {
		return err
	}
	if !has {
		s.logger.Debug("Permission denied",
			slog.String("user_id", userID.String()),
			slog.Any("permissions", permissions),
		)
		return ErrPermissionDenied
	}
	return nil
}

// GetUserPermissions returns all permissions for a user.
func (s *Service) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]db.Permission, error) {
	return s.queries.GetUserPermissions(ctx, userID)
}

// GetUserPermissionNames returns all permission names for a user.
func (s *Service) GetUserPermissionNames(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.queries.GetUserPermissionNames(ctx, userID)
}

// GetPermissionsForRole returns all permissions assigned to a role.
func (s *Service) GetPermissionsForRole(ctx context.Context, role string) ([]db.Permission, error) {
	return s.queries.GetPermissionsForRole(ctx, role)
}

// CanAccessAdultContent checks if a user can access adult content.
// This requires both the adult.browse permission AND adult_enabled flag on the user.
func (s *Service) CanAccessAdultContent(ctx context.Context, user *db.User) (bool, error) {
	// First check if user has adult content enabled
	if !user.AdultEnabled {
		return false, nil
	}

	// Then check if user has the permission
	return s.HasPermission(ctx, user.ID, PermAdultBrowse)
}
