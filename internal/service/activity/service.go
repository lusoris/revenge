// Package activity provides activity logging services.
package activity

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Activity type constants matching the database enum.
const (
	TypeUserLogin       = "user_login"
	TypeUserLogout      = "user_logout"
	TypeUserCreated     = "user_created"
	TypeUserUpdated     = "user_updated"
	TypeUserDeleted     = "user_deleted"
	TypePasswordChanged = "password_changed"
	TypeSessionCreated  = "session_created"
	TypeSessionExpired  = "session_expired"
	TypeLibraryCreated  = "library_created"
	TypeLibraryUpdated  = "library_updated"
	TypeLibraryDeleted  = "library_deleted"
	TypeLibraryScanned  = "library_scanned"
	TypeContentPlayed   = "content_played"
	TypeContentRated    = "content_rated"
	TypeSettingsChanged = "settings_changed"
	TypeAPIError        = "api_error"
	TypeSecurityEvent   = "security_event"
)

// Severity constants matching the database enum.
const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityError    = "error"
	SeverityCritical = "critical"
)

// Service provides activity logging operations.
type Service struct {
	queries *db.Queries
	logger  *slog.Logger
}

// NewService creates a new activity service.
func NewService(queries *db.Queries, logger *slog.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger.With(slog.String("service", "activity")),
	}
}

// LogParams contains parameters for logging an activity.
type LogParams struct {
	UserID    *uuid.UUID
	Type      string
	Severity  string
	Message   string
	Metadata  map[string]any
	IPAddress netip.Addr
	UserAgent *string
}

// Log creates a new activity log entry.
func (s *Service) Log(ctx context.Context, params LogParams) (*db.ActivityLog, error) {
	// Default severity to info
	severity := params.Severity
	if severity == "" {
		severity = SeverityInfo
	}

	// Marshal metadata to JSON
	var metadata []byte
	if params.Metadata != nil {
		data, err := json.Marshal(params.Metadata)
		if err != nil {
			s.logger.Warn("Failed to marshal activity metadata", "error", err)
			metadata = []byte("{}")
		} else {
			metadata = data
		}
	} else {
		metadata = []byte("{}")
	}

	// Convert userID to pgtype.UUID
	var userID pgtype.UUID
	if params.UserID != nil {
		userID = pgtype.UUID{Bytes: *params.UserID, Valid: true}
	}

	log, err := s.queries.CreateActivityLog(ctx, db.CreateActivityLogParams{
		UserID:    userID,
		Type:      params.Type,
		Severity:  severity,
		Message:   params.Message,
		Metadata:  metadata,
		IpAddress: params.IPAddress,
		UserAgent: params.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	// Also log to slog for observability
	s.logger.Log(ctx, severityToLevel(severity), params.Message,
		slog.String("type", params.Type),
		slog.Any("user_id", params.UserID),
	)

	return &log, nil
}

// LogUserLogin logs a user login event.
func (s *Service) LogUserLogin(ctx context.Context, userID uuid.UUID, ip netip.Addr, userAgent *string) error {
	_, err := s.Log(ctx, LogParams{
		UserID:    &userID,
		Type:      TypeUserLogin,
		Severity:  SeverityInfo,
		Message:   "User logged in",
		IPAddress: ip,
		UserAgent: userAgent,
	})
	return err
}

// LogUserLogout logs a user logout event.
func (s *Service) LogUserLogout(ctx context.Context, userID uuid.UUID) error {
	_, err := s.Log(ctx, LogParams{
		UserID:   &userID,
		Type:     TypeUserLogout,
		Severity: SeverityInfo,
		Message:  "User logged out",
	})
	return err
}

// LogSecurityEvent logs a security-related event.
func (s *Service) LogSecurityEvent(ctx context.Context, userID *uuid.UUID, message string, metadata map[string]any, ip netip.Addr) error {
	_, err := s.Log(ctx, LogParams{
		UserID:    userID,
		Type:      TypeSecurityEvent,
		Severity:  SeverityWarning,
		Message:   message,
		Metadata:  metadata,
		IPAddress: ip,
	})
	return err
}

// LogAPIError logs an API error event.
func (s *Service) LogAPIError(ctx context.Context, userID *uuid.UUID, message string, metadata map[string]any) error {
	_, err := s.Log(ctx, LogParams{
		UserID:   userID,
		Type:     TypeAPIError,
		Severity: SeverityError,
		Message:  message,
		Metadata: metadata,
	})
	return err
}

// ListByUser returns activity logs for a specific user.
func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.ActivityLog, error) {
	return s.queries.ListActivityLogByUser(ctx, db.ListActivityLogByUserParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Limit:  limit,
		Offset: offset,
	})
}

// ListByType returns activity logs filtered by type.
func (s *Service) ListByType(ctx context.Context, activityType string, limit, offset int32) ([]db.ActivityLog, error) {
	return s.queries.ListActivityLogByType(ctx, db.ListActivityLogByTypeParams{
		Type:   activityType,
		Limit:  limit,
		Offset: offset,
	})
}

// ListBySeverity returns activity logs filtered by severity.
func (s *Service) ListBySeverity(ctx context.Context, severity string, limit, offset int32) ([]db.ActivityLog, error) {
	return s.queries.ListActivityLogBySeverity(ctx, db.ListActivityLogBySeverityParams{
		Severity: severity,
		Limit:    limit,
		Offset:   offset,
	})
}

// ListRecent returns the most recent activity logs.
func (s *Service) ListRecent(ctx context.Context, limit, offset int32) ([]db.ActivityLog, error) {
	return s.queries.ListRecentActivity(ctx, db.ListRecentActivityParams{
		Limit:  limit,
		Offset: offset,
	})
}

// DeleteOlderThan deletes activity logs older than the specified time.
func (s *Service) DeleteOlderThan(ctx context.Context, before time.Time) error {
	return s.queries.DeleteOldActivityLogs(ctx, before)
}

// severityToLevel converts activity severity to slog level.
func severityToLevel(severity string) slog.Level {
	switch severity {
	case SeverityCritical:
		return slog.LevelError + 4 // Custom level above Error
	case SeverityError:
		return slog.LevelError
	case SeverityWarning:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
