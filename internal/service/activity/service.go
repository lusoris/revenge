// Package activity provides activity/audit logging services.
package activity

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/netip"
	"time"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Action constants for audit logging.
const (
	ActionMetadataEdit    = "metadata.edit"
	ActionMetadataLock    = "metadata.lock"
	ActionMetadataUnlock  = "metadata.unlock"
	ActionMetadataRefresh = "metadata.refresh"
	ActionImageUpload     = "image.upload"
	ActionImageSelect     = "image.select"
	ActionImageDelete     = "image.delete"
	ActionContentDelete   = "content.delete"
	ActionUserLogin       = "user.login"
	ActionUserLogout      = "user.logout"
	ActionUserCreate      = "user.create"
	ActionUserUpdate      = "user.update"
	ActionUserDelete      = "user.delete"
	ActionLibraryCreate   = "library.create"
	ActionLibraryUpdate   = "library.update"
	ActionLibraryDelete   = "library.delete"
	ActionLibraryScan     = "library.scan"
	ActionSettingsChange  = "settings.change"
	ActionSecurityEvent   = "security.event"
)

// Module constants for categorizing audit entries.
const (
	ModuleMovie   = "movie"
	ModuleTVShow  = "tvshow"
	ModuleQAR     = "qar"
	ModuleUser    = "user"
	ModuleLibrary = "library"
	ModuleSystem  = "system"
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
	UserID     uuid.UUID
	Action     string
	Module     string
	EntityID   uuid.UUID
	EntityType string
	Changes    map[string]any
	IPAddress  netip.Addr
	UserAgent  *string
}

// Log creates a new activity log entry.
func (s *Service) Log(ctx context.Context, params LogParams) (*db.ActivityLog, error) {
	// Marshal changes to JSON
	var changesJSON json.RawMessage
	if params.Changes != nil {
		data, err := json.Marshal(params.Changes)
		if err != nil {
			s.logger.Warn("Failed to marshal activity changes", "error", err)
			changesJSON = json.RawMessage("{}")
		} else {
			changesJSON = data
		}
	} else {
		changesJSON = json.RawMessage("{}")
	}

	log, err := s.queries.CreateActivityLog(ctx, db.CreateActivityLogParams{
		UserID:     params.UserID,
		Action:     params.Action,
		Module:     params.Module,
		EntityID:   params.EntityID,
		EntityType: params.EntityType,
		Changes:    changesJSON,
		IpAddress:  params.IPAddress,
		UserAgent:  params.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	// Also log to slog for observability
	s.logger.Info("Activity logged",
		slog.String("action", params.Action),
		slog.String("module", params.Module),
		slog.String("user_id", params.UserID.String()),
		slog.String("entity_id", params.EntityID.String()),
	)

	return &log, nil
}

// LogUserLogin logs a user login event.
func (s *Service) LogUserLogin(ctx context.Context, userID uuid.UUID, ip netip.Addr, userAgent *string) error {
	_, err := s.Log(ctx, LogParams{
		UserID:     userID,
		Action:     ActionUserLogin,
		Module:     ModuleUser,
		EntityID:   userID,
		EntityType: "user",
		Changes:    map[string]any{},
		IPAddress:  ip,
		UserAgent:  userAgent,
	})
	return err
}

// LogUserLogout logs a user logout event.
func (s *Service) LogUserLogout(ctx context.Context, userID uuid.UUID) error {
	_, err := s.Log(ctx, LogParams{
		UserID:     userID,
		Action:     ActionUserLogout,
		Module:     ModuleUser,
		EntityID:   userID,
		EntityType: "user",
		Changes:    map[string]any{},
	})
	return err
}

// LogMetadataEdit logs a metadata edit event.
func (s *Service) LogMetadataEdit(ctx context.Context, userID uuid.UUID, module, entityType string, entityID uuid.UUID, changes map[string]any) error {
	_, err := s.Log(ctx, LogParams{
		UserID:     userID,
		Action:     ActionMetadataEdit,
		Module:     module,
		EntityID:   entityID,
		EntityType: entityType,
		Changes:    changes,
	})
	return err
}

// LogSecurityEvent logs a security-related event.
func (s *Service) LogSecurityEvent(ctx context.Context, userID uuid.UUID, message string, metadata map[string]any, ip netip.Addr) error {
	changes := metadata
	if changes == nil {
		changes = map[string]any{}
	}
	changes["message"] = message

	_, err := s.Log(ctx, LogParams{
		UserID:     userID,
		Action:     ActionSecurityEvent,
		Module:     ModuleSystem,
		EntityID:   userID,
		EntityType: "user",
		Changes:    changes,
		IPAddress:  ip,
	})
	return err
}

// ListByUser returns activity logs for a specific user.
func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.ActivityLog, error) {
	return s.queries.ListActivityLogByUser(ctx, db.ListActivityLogByUserParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
}

// ListByAction returns activity logs filtered by action.
func (s *Service) ListByAction(ctx context.Context, action string, limit, offset int32) ([]db.ActivityLog, error) {
	return s.queries.ListActivityLogByAction(ctx, db.ListActivityLogByActionParams{
		Action: action,
		Limit:  limit,
		Offset: offset,
	})
}

// ListByModule returns activity logs filtered by module.
func (s *Service) ListByModule(ctx context.Context, module string, limit, offset int32) ([]db.ActivityLog, error) {
	return s.queries.ListActivityLogByModule(ctx, db.ListActivityLogByModuleParams{
		Module: module,
		Limit:  limit,
		Offset: offset,
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
