package activity

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrNotFound is returned when an activity log entry is not found.
	ErrNotFound = errors.New("activity log not found")
)

// Service provides activity logging functionality.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new activity service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With("component", "activity"),
	}
}

// Log records an activity entry.
func (s *Service) Log(ctx context.Context, req LogRequest) error {
	entry := &Entry{
		UserID:       req.UserID,
		Username:     req.Username,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Changes:      req.Changes,
		Metadata:     req.Metadata,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		Success:      req.Success,
		ErrorMessage: req.ErrorMessage,
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		s.logger.Error("failed to log activity",
			slog.String("action", req.Action),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Debug("activity logged",
		slog.String("id", entry.ID.String()),
		slog.String("action", req.Action),
	)

	return nil
}

// LogWithContext logs an activity with context from the request.
func (s *Service) LogWithContext(
	ctx context.Context,
	userID uuid.UUID,
	username string,
	action string,
	resourceType string,
	resourceID uuid.UUID,
	changes map[string]any,
	ipAddress net.IP,
	userAgent string,
) error {
	req := LogRequest{
		UserID:       &userID,
		Username:     &username,
		Action:       action,
		ResourceType: &resourceType,
		ResourceID:   &resourceID,
		Changes:      changes,
		IPAddress:    &ipAddress,
		UserAgent:    &userAgent,
		Success:      true,
	}
	return s.Log(ctx, req)
}

// LogFailure logs a failed action.
func (s *Service) LogFailure(
	ctx context.Context,
	userID *uuid.UUID,
	username *string,
	action string,
	errorMessage string,
	ipAddress *net.IP,
	userAgent *string,
) error {
	req := LogRequest{
		UserID:       userID,
		Username:     username,
		Action:       action,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Success:      false,
		ErrorMessage: &errorMessage,
	}
	return s.Log(ctx, req)
}

// Get retrieves a single activity log by ID.
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Entry, error) {
	entry, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// List returns paginated activity logs.
func (s *Service) List(ctx context.Context, limit, offset int32) ([]Entry, int64, error) {
	entries, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return entries, count, nil
}

// Search returns activity logs matching filters.
func (s *Service) Search(ctx context.Context, filters SearchFilters) ([]Entry, int64, error) {
	if filters.Limit <= 0 {
		filters.Limit = 50
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	return s.repo.Search(ctx, filters)
}

// GetUserActivity returns activity logs for a specific user.
func (s *Service) GetUserActivity(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]Entry, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetByUser(ctx, userID, limit, offset)
}

// GetResourceActivity returns activity logs for a specific resource.
func (s *Service) GetResourceActivity(ctx context.Context, resourceType string, resourceID uuid.UUID, limit, offset int32) ([]Entry, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetByResource(ctx, resourceType, resourceID, limit, offset)
}

// GetFailedActivity returns failed activity logs.
func (s *Service) GetFailedActivity(ctx context.Context, limit, offset int32) ([]Entry, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetFailed(ctx, limit, offset)
}

// GetStats returns activity log statistics.
func (s *Service) GetStats(ctx context.Context) (*Stats, error) {
	return s.repo.GetStats(ctx)
}

// GetRecentActions returns recent distinct actions.
func (s *Service) GetRecentActions(ctx context.Context, limit int32) ([]ActionCount, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	return s.repo.GetRecentActions(ctx, limit)
}

// CleanupOldLogs deletes activity logs older than the retention period.
func (s *Service) CleanupOldLogs(ctx context.Context, olderThan time.Time) (int64, error) {
	count, err := s.repo.DeleteOld(ctx, olderThan)
	if err != nil {
		s.logger.Error("failed to cleanup old activity logs",
			slog.Time("older_than", olderThan),
			slog.Any("error", err),
		)
		return 0, err
	}

	s.logger.Info("cleaned up old activity logs",
		slog.Int64("deleted_count", count),
		slog.Time("older_than", olderThan),
	)

	return count, nil
}

// CountOldLogs counts activity logs older than the given time (for dry-run).
func (s *Service) CountOldLogs(ctx context.Context, olderThan time.Time) (int64, error) {
	return s.repo.CountOld(ctx, olderThan)
}
