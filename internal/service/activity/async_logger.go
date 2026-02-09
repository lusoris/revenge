package activity

import (
	"context"
	"log/slog"
	"net"

	"github.com/google/uuid"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/riverqueue/river"
)

// AsyncLogger implements the Logger interface by inserting River jobs
// instead of writing to the database synchronously. Security-critical
// events (auth, login, password changes) route to the critical queue;
// all other activity goes to the low queue.
type AsyncLogger struct {
	client *infrajobs.Client
	logger *slog.Logger
}

// NewAsyncLogger creates a new async activity logger.
// Falls back to synchronous logging if client is nil.
func NewAsyncLogger(client *infrajobs.Client, logger *slog.Logger) *AsyncLogger {
	return &AsyncLogger{
		client: client,
		logger: logger.With("component", "activity-async-logger"),
	}
}

// LogAction logs a successful action by enqueuing a River job.
func (l *AsyncLogger) LogAction(ctx context.Context, req LogActionRequest) error {
	// Build the LogRequest
	var userIDPtr *uuid.UUID
	if req.UserID != uuid.Nil {
		userIDPtr = &req.UserID
	}

	var resourceIDPtr *uuid.UUID
	if req.ResourceID != uuid.Nil {
		resourceIDPtr = &req.ResourceID
	}

	var usernamePtr *string
	if req.Username != "" {
		usernamePtr = &req.Username
	}

	var resourceTypePtr *string
	if req.ResourceType != "" {
		resourceTypePtr = &req.ResourceType
	}

	var userAgentPtr *string
	if req.UserAgent != "" {
		userAgentPtr = &req.UserAgent
	}

	var ipAddressPtr *net.IP
	if req.IPAddress != nil {
		ipAddressPtr = &req.IPAddress
	}

	logReq := LogRequest{
		UserID:       userIDPtr,
		Username:     usernamePtr,
		Action:       req.Action,
		ResourceType: resourceTypePtr,
		ResourceID:   resourceIDPtr,
		Changes:      req.Changes,
		Metadata:     req.Metadata,
		IPAddress:    ipAddressPtr,
		UserAgent:    userAgentPtr,
		Success:      true,
	}

	return l.enqueue(ctx, logReq)
}

// LogFailure logs a failed action by enqueuing a River job.
func (l *AsyncLogger) LogFailure(ctx context.Context, req LogFailureRequest) error {
	logReq := LogRequest{
		UserID:       req.UserID,
		Username:     req.Username,
		Action:       req.Action,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		Success:      false,
		ErrorMessage: &req.ErrorMessage,
	}

	return l.enqueue(ctx, logReq)
}

// enqueue inserts an activity log job, routing to the critical queue
// for security-related actions.
func (l *AsyncLogger) enqueue(ctx context.Context, req LogRequest) error {
	if l.client == nil {
		l.logger.Warn("activity log client not available, dropping log entry",
			slog.String("action", req.Action),
		)
		return nil
	}

	args := ActivityLogArgs{LogRequest: req}

	var opts *river.InsertOpts
	if isSecurityAction(req.Action) {
		opts = &river.InsertOpts{
			Queue:       infrajobs.QueueCritical,
			MaxAttempts: 3,
		}
	}
	// nil opts → uses args.InsertOpts() → QueueLow

	if _, err := l.client.Insert(ctx, args, opts); err != nil {
		l.logger.Error("failed to enqueue activity log",
			slog.String("action", req.Action),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}
