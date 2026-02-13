// Package activity provides audit logging and event tracking functionality.
package activity

import (
	"context"
	"net"
	"time"

	"github.com/google/uuid"
)

// Repository defines the interface for activity log persistence.
type Repository interface {
	// Create logs a new activity entry
	Create(ctx context.Context, entry *Entry) error

	// Get retrieves a single activity log by ID
	Get(ctx context.Context, id uuid.UUID) (*Entry, error)

	// List returns paginated activity logs
	List(ctx context.Context, limit, offset int32) ([]Entry, error)

	// Count returns total activity log count
	Count(ctx context.Context) (int64, error)

	// Search returns activity logs matching filters
	Search(ctx context.Context, filters SearchFilters) ([]Entry, int64, error)

	// GetByUser returns activity logs for a specific user
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]Entry, int64, error)

	// GetByResource returns activity logs for a specific resource
	GetByResource(ctx context.Context, resourceType string, resourceID uuid.UUID, limit, offset int32) ([]Entry, int64, error)

	// GetByAction returns activity logs by action type
	GetByAction(ctx context.Context, action string, limit, offset int32) ([]Entry, error)

	// GetByIP returns activity logs from a specific IP
	GetByIP(ctx context.Context, ip net.IP, limit, offset int32) ([]Entry, error)

	// GetFailed returns failed activity logs
	GetFailed(ctx context.Context, limit, offset int32) ([]Entry, error)

	// DeleteOld deletes activity logs older than the given time
	DeleteOld(ctx context.Context, olderThan time.Time) (int64, error)

	// CountOld counts activity logs older than the given time
	CountOld(ctx context.Context, olderThan time.Time) (int64, error)

	// GetStats returns activity log statistics
	GetStats(ctx context.Context) (*Stats, error)

	// GetRecentActions returns recent distinct actions
	GetRecentActions(ctx context.Context, limit int32) ([]ActionCount, error)
}

// Entry represents a single activity log entry.
type Entry struct {
	ID           uuid.UUID      `json:"id"`
	UserID       *uuid.UUID     `json:"user_id,omitempty"`
	Username     *string        `json:"username,omitempty"`
	Action       string         `json:"action"`
	ResourceType *string        `json:"resource_type,omitempty"`
	ResourceID   *uuid.UUID     `json:"resource_id,omitempty"`
	Changes      map[string]any `json:"changes,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	IPAddress    *net.IP        `json:"ip_address,omitempty"`
	UserAgent    *string        `json:"user_agent,omitempty"`
	Success      bool           `json:"success"`
	ErrorMessage *string        `json:"error_message,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// SearchFilters defines filters for searching activity logs.
type SearchFilters struct {
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	Action       *string    `json:"action,omitempty"`
	ResourceType *string    `json:"resource_type,omitempty"`
	ResourceID   *uuid.UUID `json:"resource_id,omitempty"`
	Success      *bool      `json:"success,omitempty"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	Limit        int32      `json:"limit"`
	Offset       int32      `json:"offset"`
}

// Stats represents activity log statistics.
type Stats struct {
	TotalCount   int64      `json:"total_count"`
	SuccessCount int64      `json:"success_count"`
	FailedCount  int64      `json:"failed_count"`
	OldestEntry  *time.Time `json:"oldest_entry,omitempty"`
	NewestEntry  *time.Time `json:"newest_entry,omitempty"`
}

// ActionCount represents an action with its count.
type ActionCount struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

// LogRequest represents a request to log an activity.
type LogRequest struct {
	UserID       *uuid.UUID     `json:"user_id,omitempty"`
	Username     *string        `json:"username,omitempty"`
	Action       string         `json:"action"`
	ResourceType *string        `json:"resource_type,omitempty"`
	ResourceID   *uuid.UUID     `json:"resource_id,omitempty"`
	Changes      map[string]any `json:"changes,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	IPAddress    *net.IP        `json:"ip_address,omitempty"`
	UserAgent    *string        `json:"user_agent,omitempty"`
	Success      bool           `json:"success"`
	ErrorMessage *string        `json:"error_message,omitempty"`
}

// Common action constants
const (
	// User actions
	ActionUserLogin         = "user.login"
	ActionUserLogout        = "user.logout"
	ActionUserCreate        = "user.create"
	ActionUserUpdate        = "user.update"
	ActionUserDelete        = "user.delete"
	ActionUserPasswordReset = "user.password_reset"

	// Session actions
	ActionSessionCreate  = "session.create"
	ActionSessionRevoke  = "session.revoke"
	ActionSessionExpired = "session.expired"

	// API Key actions
	ActionAPIKeyCreate = "apikey.create"
	ActionAPIKeyRevoke = "apikey.revoke"

	// OIDC actions
	ActionOIDCLogin          = "oidc.login"
	ActionOIDCLink           = "oidc.link"
	ActionOIDCUnlink         = "oidc.unlink"
	ActionOIDCProviderAdd    = "oidc.provider.add"
	ActionOIDCProviderUpdate = "oidc.provider.update"
	ActionOIDCProviderDelete = "oidc.provider.delete"

	// Settings actions
	ActionSettingsUpdate = "settings.update"

	// Library actions
	ActionLibraryCreate = "library.create"
	ActionLibraryUpdate = "library.update"
	ActionLibraryDelete = "library.delete"
	ActionLibraryScan   = "library.scan"

	// Admin actions
	ActionAdminRoleAssign = "admin.role.assign"
	ActionAdminRoleRevoke = "admin.role.revoke"
	ActionAdminUserBan    = "admin.user.ban"
	ActionAdminUserUnban  = "admin.user.unban"
)

// Resource type constants
const (
	ResourceTypeUser    = "user"
	ResourceTypeSession = "session"
	ResourceTypeAPIKey  = "apikey"
	ResourceTypeOIDC    = "oidc"
	ResourceTypeSetting = "setting"
	ResourceTypeLibrary = "library"
	ResourceTypeMovie   = "movie"
	ResourceTypeTVShow  = "tvshow"
	ResourceTypeEpisode = "episode"
)
