package api

import (
	"context"
	"encoding/json"
	"errors"

	"log/slog"

	"github.com/go-faster/jx"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/validate"
)

// ============================================================================
// Activity Log Admin Endpoints (Admin only)
// ============================================================================

// SearchActivityLogs searches and filters activity logs with pagination.
// GET /api/v1/admin/activity
func (h *Handler) SearchActivityLogs(ctx context.Context, params ogen.SearchActivityLogsParams) (ogen.SearchActivityLogsRes, error) {
	// Admin check
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.SearchActivityLogsUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.SearchActivityLogsForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	// Build search filters
	filters := activity.SearchFilters{
		Limit:  50,
		Offset: 0,
	}

	if params.UserID.IsSet() {
		userID := params.UserID.Value
		filters.UserID = &userID
	}
	if params.Action.IsSet() {
		action := params.Action.Value
		filters.Action = &action
	}
	if params.ResourceType.IsSet() {
		rt := params.ResourceType.Value
		filters.ResourceType = &rt
	}
	if params.ResourceID.IsSet() {
		rid := params.ResourceID.Value
		filters.ResourceID = &rid
	}
	if params.Success.IsSet() {
		success := params.Success.Value
		filters.Success = &success
	}
	if params.StartTime.IsSet() {
		st := params.StartTime.Value
		filters.StartTime = &st
	}
	if params.EndTime.IsSet() {
		et := params.EndTime.Value
		filters.EndTime = &et
	}
	if params.Limit.IsSet() {
		limit, err := validate.SafeInt32(params.Limit.Value)
		if err != nil {
			h.logger.Error("invalid limit value", slog.Any("error",err))
			return &ogen.SearchActivityLogsForbidden{
				Code:    400,
				Message: "Invalid limit parameter",
			}, nil
		}
		filters.Limit = limit
	}
	if params.Offset.IsSet() {
		offset, err := validate.SafeInt32(params.Offset.Value)
		if err != nil {
			h.logger.Error("invalid offset value", slog.Any("error",err))
			return &ogen.SearchActivityLogsForbidden{
				Code:    400,
				Message: "Invalid offset parameter",
			}, nil
		}
		filters.Offset = offset
	}

	entries, total, err := h.activityService.Search(ctx, filters)
	if err != nil {
		h.logger.Error("failed to search activity logs", slog.Any("error",err))
		return &ogen.SearchActivityLogsForbidden{
			Code:    500,
			Message: "Failed to search activity logs",
		}, nil
	}

	return &ogen.ActivityLogListResponse{
		Entries:  convertActivityEntries(entries),
		Total:    total,
		Page:     ogen.NewOptInt(int(filters.Offset/filters.Limit) + 1),
		PageSize: ogen.NewOptInt(int(filters.Limit)),
	}, nil
}

// GetUserActivityLogs returns activity logs for a specific user.
// GET /api/v1/admin/activity/users/{userId}
func (h *Handler) GetUserActivityLogs(ctx context.Context, params ogen.GetUserActivityLogsParams) (ogen.GetUserActivityLogsRes, error) {
	// Admin check
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.GetUserActivityLogsUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.GetUserActivityLogsForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	limit := int32(50)
	offset := int32(0)
	if params.Limit.IsSet() {
		l, err := validate.SafeInt32(params.Limit.Value)
		if err != nil {
			h.logger.Error("invalid limit value", slog.Any("error",err))
			return &ogen.GetUserActivityLogsForbidden{
				Code:    400,
				Message: "Invalid limit parameter",
			}, nil
		}
		limit = l
	}
	if params.Offset.IsSet() {
		o, err := validate.SafeInt32(params.Offset.Value)
		if err != nil {
			h.logger.Error("invalid offset value", slog.Any("error",err))
			return &ogen.GetUserActivityLogsForbidden{
				Code:    400,
				Message: "Invalid offset parameter",
			}, nil
		}
		offset = o
	}

	entries, total, err := h.activityService.GetUserActivity(ctx, params.UserID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get user activity logs",
			slog.String("user_id", params.UserID.String()),
			slog.Any("error",err),
		)
		return &ogen.GetUserActivityLogsForbidden{
			Code:    500,
			Message: "Failed to get user activity logs",
		}, nil
	}

	return &ogen.ActivityLogListResponse{
		Entries:  convertActivityEntries(entries),
		Total:    total,
		Page:     ogen.NewOptInt(int(offset/limit) + 1),
		PageSize: ogen.NewOptInt(int(limit)),
	}, nil
}

// GetResourceActivityLogs returns activity logs for a specific resource.
// GET /api/v1/admin/activity/resources/{resourceType}/{resourceId}
func (h *Handler) GetResourceActivityLogs(ctx context.Context, params ogen.GetResourceActivityLogsParams) (ogen.GetResourceActivityLogsRes, error) {
	// Admin check
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.GetResourceActivityLogsUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.GetResourceActivityLogsForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	limit := int32(50)
	offset := int32(0)
	if params.Limit.IsSet() {
		l, err := validate.SafeInt32(params.Limit.Value)
		if err != nil {
			h.logger.Error("invalid limit value", slog.Any("error",err))
			return &ogen.GetResourceActivityLogsForbidden{
				Code:    400,
				Message: "Invalid limit parameter",
			}, nil
		}
		limit = l
	}
	if params.Offset.IsSet() {
		o, err := validate.SafeInt32(params.Offset.Value)
		if err != nil {
			h.logger.Error("invalid offset value", slog.Any("error",err))
			return &ogen.GetResourceActivityLogsForbidden{
				Code:    400,
				Message: "Invalid offset parameter",
			}, nil
		}
		offset = o
	}

	entries, total, err := h.activityService.GetResourceActivity(ctx, params.ResourceType, params.ResourceID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get resource activity logs",
			slog.String("resource_type", params.ResourceType),
			slog.String("resource_id", params.ResourceID.String()),
			slog.Any("error",err),
		)
		return &ogen.GetResourceActivityLogsForbidden{
			Code:    500,
			Message: "Failed to get resource activity logs",
		}, nil
	}

	return &ogen.ActivityLogListResponse{
		Entries:  convertActivityEntries(entries),
		Total:    total,
		Page:     ogen.NewOptInt(int(offset/limit) + 1),
		PageSize: ogen.NewOptInt(int(limit)),
	}, nil
}

// GetActivityStats returns activity log statistics.
// GET /api/v1/admin/activity/stats
func (h *Handler) GetActivityStats(ctx context.Context) (ogen.GetActivityStatsRes, error) {
	// Admin check
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.GetActivityStatsUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.GetActivityStatsForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	stats, err := h.activityService.GetStats(ctx)
	if err != nil {
		h.logger.Error("failed to get activity stats", slog.Any("error",err))
		return &ogen.GetActivityStatsForbidden{
			Code:    500,
			Message: "Failed to get activity statistics",
		}, nil
	}

	result := &ogen.ActivityStats{
		TotalCount:   stats.TotalCount,
		SuccessCount: stats.SuccessCount,
		FailedCount:  stats.FailedCount,
	}

	if stats.OldestEntry != nil {
		result.OldestEntry = ogen.NewOptDateTime(*stats.OldestEntry)
	}
	if stats.NewestEntry != nil {
		result.NewestEntry = ogen.NewOptDateTime(*stats.NewestEntry)
	}

	return result, nil
}

// GetRecentActions returns recent distinct action types for filtering.
// GET /api/v1/admin/activity/actions
func (h *Handler) GetRecentActions(ctx context.Context, params ogen.GetRecentActionsParams) (ogen.GetRecentActionsRes, error) {
	// Admin check
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.GetRecentActionsUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.GetRecentActionsForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	limit := int32(20)
	if params.Limit.IsSet() {
		l, err := validate.SafeInt32(params.Limit.Value)
		if err != nil {
			h.logger.Error("invalid limit value", slog.Any("error",err))
			return &ogen.GetRecentActionsForbidden{
				Code:    400,
				Message: "Invalid limit parameter",
			}, nil
		}
		limit = l
	}

	actions, err := h.activityService.GetRecentActions(ctx, limit)
	if err != nil {
		h.logger.Error("failed to get recent actions", slog.Any("error",err))
		return &ogen.GetRecentActionsForbidden{
			Code:    500,
			Message: "Failed to get recent actions",
		}, nil
	}

	ogenActions := make([]ogen.ActionCount, len(actions))
	for i, a := range actions {
		ogenActions[i] = ogen.ActionCount{
			Action: a.Action,
			Count:  a.Count,
		}
	}

	return &ogen.ActionCountListResponse{
		Actions: ogenActions,
	}, nil
}

// ============================================================================
// Helper Functions
// ============================================================================

// convertActivityEntries converts domain activity entries to ogen types.
func convertActivityEntries(entries []activity.Entry) []ogen.ActivityLogEntry {
	result := make([]ogen.ActivityLogEntry, len(entries))
	for i, e := range entries {
		entry := ogen.ActivityLogEntry{
			ID:        e.ID,
			Action:    e.Action,
			Success:   e.Success,
			CreatedAt: e.CreatedAt,
		}

		if e.UserID != nil {
			entry.UserID = ogen.NewOptUUID(*e.UserID)
		}
		if e.Username != nil {
			entry.Username = ogen.NewOptString(*e.Username)
		}
		if e.ResourceType != nil {
			entry.ResourceType = ogen.NewOptString(*e.ResourceType)
		}
		if e.ResourceID != nil {
			entry.ResourceID = ogen.NewOptUUID(*e.ResourceID)
		}
		if e.IPAddress != nil {
			entry.IPAddress = ogen.NewOptString(e.IPAddress.String())
		}
		if e.UserAgent != nil {
			entry.UserAgent = ogen.NewOptString(*e.UserAgent)
		}
		if e.ErrorMessage != nil {
			entry.ErrorMessage = ogen.NewOptString(*e.ErrorMessage)
		}

		// Convert changes map to ogen format
		if len(e.Changes) > 0 {
			changes := make(ogen.ActivityLogEntryChanges)
			for k, v := range e.Changes {
				if data, err := json.Marshal(v); err == nil {
					changes[k] = jx.Raw(data)
				}
			}
			entry.Changes = ogen.NewOptActivityLogEntryChanges(changes)
		}

		// Convert metadata map to ogen format
		if len(e.Metadata) > 0 {
			metadata := make(ogen.ActivityLogEntryMetadata)
			for k, v := range e.Metadata {
				if data, err := json.Marshal(v); err == nil {
					metadata[k] = jx.Raw(data)
				}
			}
			entry.Metadata = ogen.NewOptActivityLogEntryMetadata(metadata)
		}

		result[i] = entry
	}
	return result
}
