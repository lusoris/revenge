package activity

import (
	"context"
	"encoding/json"
	"net"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// RepositoryPg implements the Repository interface using PostgreSQL.
type RepositoryPg struct {
	queries *db.Queries
}

// NewRepositoryPg creates a new PostgreSQL-backed activity repository.
func NewRepositoryPg(queries *db.Queries) *RepositoryPg {
	return &RepositoryPg{queries: queries}
}

// Create logs a new activity entry.
func (r *RepositoryPg) Create(ctx context.Context, entry *Entry) error {
	var changesJSON, metadataJSON []byte
	var err error

	if entry.Changes != nil {
		changesJSON, err = json.Marshal(entry.Changes)
		if err != nil {
			return err
		}
	}

	if entry.Metadata != nil {
		metadataJSON, err = json.Marshal(entry.Metadata)
		if err != nil {
			return err
		}
	}

	params := db.CreateActivityLogParams{
		Username:     entry.Username,
		Action:       entry.Action,
		ResourceType: entry.ResourceType,
		Changes:      changesJSON,
		Metadata:     metadataJSON,
		UserAgent:    entry.UserAgent,
		Success:      &entry.Success,
		ErrorMessage: entry.ErrorMessage,
	}

	if entry.UserID != nil {
		params.UserID = uuidToPgtype(*entry.UserID)
	}

	if entry.ResourceID != nil {
		params.ResourceID = uuidToPgtype(*entry.ResourceID)
	}

	if entry.IPAddress != nil {
		addr, ok := netip.AddrFromSlice(*entry.IPAddress)
		if ok {
			params.IpAddress = addr
		}
	}

	result, err := r.queries.CreateActivityLog(ctx, params)
	if err != nil {
		return err
	}

	entry.ID = result.ID
	entry.CreatedAt = result.CreatedAt
	return nil
}

// Get retrieves a single activity log by ID.
func (r *RepositoryPg) Get(ctx context.Context, id uuid.UUID) (*Entry, error) {
	result, err := r.queries.GetActivityLog(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbActivityToEntry(result), nil
}

// List returns paginated activity logs.
func (r *RepositoryPg) List(ctx context.Context, limit, offset int32) ([]Entry, error) {
	results, err := r.queries.ListActivityLogs(ctx, db.ListActivityLogsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, len(results))
	for i, result := range results {
		entries[i] = *dbActivityToEntry(result)
	}
	return entries, nil
}

// Count returns total activity log count.
func (r *RepositoryPg) Count(ctx context.Context) (int64, error) {
	return r.queries.CountActivityLogs(ctx)
}

// Search returns activity logs matching filters.
func (r *RepositoryPg) Search(ctx context.Context, filters SearchFilters) ([]Entry, int64, error) {
	params := db.SearchActivityLogsParams{
		Limit:  filters.Limit,
		Offset: filters.Offset,
	}

	countParams := db.CountSearchActivityLogsParams{}

	if filters.UserID != nil {
		params.Column1 = *filters.UserID
		countParams.Column1 = *filters.UserID
	}

	if filters.Action != nil {
		params.Column2 = *filters.Action
		countParams.Column2 = *filters.Action
	}

	if filters.ResourceType != nil {
		params.Column3 = *filters.ResourceType
		countParams.Column3 = *filters.ResourceType
	}

	if filters.ResourceID != nil {
		params.Column4 = *filters.ResourceID
		countParams.Column4 = *filters.ResourceID
	}

	if filters.Success != nil {
		params.Column5 = *filters.Success
		countParams.Column5 = *filters.Success
	}

	if filters.StartTime != nil {
		params.Column6 = *filters.StartTime
		countParams.Column6 = *filters.StartTime
	}

	if filters.EndTime != nil {
		params.Column7 = *filters.EndTime
		countParams.Column7 = *filters.EndTime
	}

	count, err := r.queries.CountSearchActivityLogs(ctx, countParams)
	if err != nil {
		return nil, 0, err
	}

	results, err := r.queries.SearchActivityLogs(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	entries := make([]Entry, len(results))
	for i, result := range results {
		entries[i] = *dbActivityToEntry(result)
	}

	return entries, count, nil
}

// GetByUser returns activity logs for a specific user.
func (r *RepositoryPg) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]Entry, int64, error) {
	count, err := r.queries.CountUserActivityLogs(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, 0, err
	}

	results, err := r.queries.GetUserActivityLogs(ctx, db.GetUserActivityLogsParams{
		UserID: uuidToPgtype(userID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	entries := make([]Entry, len(results))
	for i, result := range results {
		entries[i] = *dbActivityToEntry(result)
	}

	return entries, count, nil
}

// GetByResource returns activity logs for a specific resource.
func (r *RepositoryPg) GetByResource(ctx context.Context, resourceType string, resourceID uuid.UUID, limit, offset int32) ([]Entry, int64, error) {
	count, err := r.queries.CountResourceActivityLogs(ctx, db.CountResourceActivityLogsParams{
		ResourceType: &resourceType,
		ResourceID:   uuidToPgtype(resourceID),
	})
	if err != nil {
		return nil, 0, err
	}

	results, err := r.queries.GetResourceActivityLogs(ctx, db.GetResourceActivityLogsParams{
		ResourceType: &resourceType,
		ResourceID:   uuidToPgtype(resourceID),
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		return nil, 0, err
	}

	entries := make([]Entry, len(results))
	for i, result := range results {
		entries[i] = *dbActivityToEntry(result)
	}

	return entries, count, nil
}

// GetByAction returns activity logs by action type.
func (r *RepositoryPg) GetByAction(ctx context.Context, action string, limit, offset int32) ([]Entry, error) {
	results, err := r.queries.GetActivityLogsByAction(ctx, db.GetActivityLogsByActionParams{
		Action: action,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, len(results))
	for i, result := range results {
		entries[i] = *dbActivityToEntry(result)
	}

	return entries, nil
}

// GetByIP returns activity logs from a specific IP.
func (r *RepositoryPg) GetByIP(ctx context.Context, ip net.IP, limit, offset int32) ([]Entry, error) {
	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return nil, nil
	}

	results, err := r.queries.GetActivityLogsByIP(ctx, db.GetActivityLogsByIPParams{
		IpAddress: addr,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, len(results))
	for i, result := range results {
		entries[i] = *dbActivityToEntry(result)
	}

	return entries, nil
}

// GetFailed returns failed activity logs.
func (r *RepositoryPg) GetFailed(ctx context.Context, limit, offset int32) ([]Entry, error) {
	results, err := r.queries.GetFailedActivityLogs(ctx, db.GetFailedActivityLogsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, len(results))
	for i, result := range results {
		entries[i] = *dbActivityToEntry(result)
	}

	return entries, nil
}

// DeleteOld deletes activity logs older than the given time.
func (r *RepositoryPg) DeleteOld(ctx context.Context, olderThan time.Time) (int64, error) {
	return r.queries.DeleteOldActivityLogs(ctx, olderThan)
}

// CountOld counts activity logs older than the given time.
func (r *RepositoryPg) CountOld(ctx context.Context, olderThan time.Time) (int64, error) {
	return r.queries.GetOldActivityLogsCount(ctx, olderThan)
}

// GetStats returns activity log statistics.
func (r *RepositoryPg) GetStats(ctx context.Context) (*Stats, error) {
	result, err := r.queries.GetActivityLogStats(ctx)
	if err != nil {
		return nil, err
	}

	stats := &Stats{
		TotalCount:   result.TotalCount,
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
	}

	// Handle interface{} types for nullable timestamps
	if oldest, ok := result.OldestEntry.(time.Time); ok {
		stats.OldestEntry = &oldest
	}
	if newest, ok := result.NewestEntry.(time.Time); ok {
		stats.NewestEntry = &newest
	}

	return stats, nil
}

// GetRecentActions returns recent distinct actions.
func (r *RepositoryPg) GetRecentActions(ctx context.Context, limit int32) ([]ActionCount, error) {
	results, err := r.queries.GetRecentActions(ctx, limit)
	if err != nil {
		return nil, err
	}

	actions := make([]ActionCount, len(results))
	for i, result := range results {
		actions[i] = ActionCount{
			Action: result.Action,
			Count:  result.Count,
		}
	}

	return actions, nil
}

// Helper functions

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func dbActivityToEntry(a db.ActivityLog) *Entry {
	entry := &Entry{
		ID:           a.ID,
		Username:     a.Username,
		Action:       a.Action,
		ResourceType: a.ResourceType,
		UserAgent:    a.UserAgent,
		ErrorMessage: a.ErrorMessage,
		CreatedAt:    a.CreatedAt,
	}

	if a.UserID.Valid {
		userID := uuid.UUID(a.UserID.Bytes)
		entry.UserID = &userID
	}

	if a.ResourceID.Valid {
		resourceID := uuid.UUID(a.ResourceID.Bytes)
		entry.ResourceID = &resourceID
	}

	if len(a.Changes) > 0 {
		var changes map[string]interface{}
		if err := json.Unmarshal(a.Changes, &changes); err == nil {
			entry.Changes = changes
		}
	}

	if len(a.Metadata) > 0 {
		var metadata map[string]interface{}
		if err := json.Unmarshal(a.Metadata, &metadata); err == nil {
			entry.Metadata = metadata
		}
	}

	if a.IpAddress.IsValid() {
		ip := net.IP(a.IpAddress.AsSlice())
		entry.IPAddress = &ip
	}

	if a.Success != nil {
		entry.Success = *a.Success
	}

	return entry
}
