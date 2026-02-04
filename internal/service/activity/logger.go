// Package activity provides audit logging and event tracking functionality.
package activity

import (
	"context"
	"net"

	"github.com/google/uuid"
)

// Logger provides a simplified interface for activity logging.
// This interface is designed to be injected into other services
// for audit logging without tight coupling.
type Logger interface {
	// LogAction logs a successful action
	LogAction(ctx context.Context, req LogActionRequest) error

	// LogFailure logs a failed action
	LogFailure(ctx context.Context, req LogFailureRequest) error
}

// LogActionRequest contains data for logging a successful action.
type LogActionRequest struct {
	UserID       uuid.UUID
	Username     string
	Action       string
	ResourceType string
	ResourceID   uuid.UUID
	Changes      map[string]interface{}
	Metadata     map[string]interface{}
	IPAddress    net.IP
	UserAgent    string
}

// LogFailureRequest contains data for logging a failed action.
type LogFailureRequest struct {
	UserID       *uuid.UUID
	Username     *string
	Action       string
	ErrorMessage string
	IPAddress    *net.IP
	UserAgent    *string
}

// ServiceLogger wraps the activity Service to implement Logger interface.
type ServiceLogger struct {
	svc *Service
}

// NewLogger creates a new activity logger wrapping the service.
func NewLogger(svc *Service) Logger {
	return &ServiceLogger{svc: svc}
}

// LogAction logs a successful action.
func (l *ServiceLogger) LogAction(ctx context.Context, req LogActionRequest) error {
	// Handle nil UUID (uuid.Nil) - don't store it as a pointer
	var userIDPtr *uuid.UUID
	if req.UserID != uuid.Nil {
		userIDPtr = &req.UserID
	}

	var resourceIDPtr *uuid.UUID
	if req.ResourceID != uuid.Nil {
		resourceIDPtr = &req.ResourceID
	}

	// Handle empty strings
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

	return l.svc.Log(ctx, LogRequest{
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
	})
}

// LogFailure logs a failed action.
func (l *ServiceLogger) LogFailure(ctx context.Context, req LogFailureRequest) error {
	return l.svc.Log(ctx, LogRequest{
		UserID:       req.UserID,
		Username:     req.Username,
		Action:       req.Action,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		Success:      false,
		ErrorMessage: &req.ErrorMessage,
	})
}

// NoopLogger is a no-op logger for testing or when logging is disabled.
type NoopLogger struct{}

// NewNoopLogger creates a no-op activity logger.
func NewNoopLogger() Logger {
	return &NoopLogger{}
}

// LogAction does nothing.
func (l *NoopLogger) LogAction(_ context.Context, _ LogActionRequest) error {
	return nil
}

// LogFailure does nothing.
func (l *NoopLogger) LogFailure(_ context.Context, _ LogFailureRequest) error {
	return nil
}
