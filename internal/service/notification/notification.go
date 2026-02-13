// Package notification provides a notification service for dispatching events
// to various notification agents (webhook, discord, email, gotify/ntfy).
package notification

import (
	"context"
	"slices"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of event that triggered a notification
type EventType string

// Event types for notifications
const (
	// Content events
	EventMovieAdded     EventType = "movie.added"
	EventMovieAvailable EventType = "movie.available"
	EventMovieUpdated   EventType = "movie.updated"
	EventMovieDeleted   EventType = "movie.deleted"

	// Request events
	EventRequestCreated  EventType = "request.created"
	EventRequestApproved EventType = "request.approved"
	EventRequestDenied   EventType = "request.denied"

	// Library events
	EventLibraryScanStarted EventType = "library.scan_started"
	EventLibraryScanDone    EventType = "library.scan_done"
	EventLibraryCreated     EventType = "library.created"
	EventLibraryDeleted     EventType = "library.deleted"

	// User events
	EventUserCreated EventType = "user.created"
	EventUserDeleted EventType = "user.deleted"
	EventUserUpdated EventType = "user.updated"

	// Auth/Security events
	EventLoginSuccess    EventType = "auth.login_success"
	EventLoginFailed     EventType = "auth.login_failed"
	EventMFAEnabled      EventType = "auth.mfa_enabled"
	EventMFADisabled     EventType = "auth.mfa_disabled"
	EventPasswordChanged EventType = "auth.password_changed"
	EventPasswordReset   EventType = "auth.password_reset"

	// Playback events
	EventPlaybackStarted EventType = "playback.started"
	EventPlaybackStopped EventType = "playback.stopped"

	// System events
	EventSystemStartup  EventType = "system.startup"
	EventSystemShutdown EventType = "system.shutdown"
	EventBackupCreated  EventType = "system.backup_created"

	// Integration events
	EventRadarrSync EventType = "integration.radarr_sync"
)

// EventCategory groups events for notification preferences
type EventCategory string

const (
	CategoryContent  EventCategory = "content"
	CategoryRequests EventCategory = "requests"
	CategoryLibrary  EventCategory = "library"
	CategoryUser     EventCategory = "user"
	CategoryAuth     EventCategory = "auth"
	CategoryPlayback EventCategory = "playback"
	CategorySystem   EventCategory = "system"
)

// GetCategory returns the category for an event type
func (e EventType) GetCategory() EventCategory {
	switch e {
	case EventMovieAdded, EventMovieAvailable, EventMovieUpdated, EventMovieDeleted:
		return CategoryContent
	case EventRequestCreated, EventRequestApproved, EventRequestDenied:
		return CategoryRequests
	case EventLibraryScanStarted, EventLibraryScanDone, EventLibraryCreated, EventLibraryDeleted:
		return CategoryLibrary
	case EventUserCreated, EventUserDeleted, EventUserUpdated:
		return CategoryUser
	case EventLoginSuccess, EventLoginFailed, EventMFAEnabled, EventMFADisabled, EventPasswordChanged, EventPasswordReset:
		return CategoryAuth
	case EventPlaybackStarted, EventPlaybackStopped:
		return CategoryPlayback
	default:
		return CategorySystem
	}
}

// String returns the string representation of the event type
func (e EventType) String() string {
	return string(e)
}

// Event represents a notification event
type Event struct {
	ID        uuid.UUID         `json:"id"`
	Type      EventType         `json:"type"`
	Timestamp time.Time         `json:"timestamp"`
	UserID    *uuid.UUID        `json:"user_id,omitempty"`   // User who triggered the event (if applicable)
	TargetID  *uuid.UUID        `json:"target_id,omitempty"` // Target resource ID (movie, user, etc.)
	Data      map[string]any    `json:"data,omitempty"`      // Event-specific data
	Metadata  map[string]string `json:"metadata,omitempty"`  // Additional metadata
}

// NewEvent creates a new event with a generated ID and current timestamp
func NewEvent(eventType EventType) *Event {
	return &Event{
		ID:        uuid.Must(uuid.NewV7()),
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Data:      make(map[string]any),
		Metadata:  make(map[string]string),
	}
}

// WithUser sets the user ID for the event
func (e *Event) WithUser(userID uuid.UUID) *Event {
	e.UserID = &userID
	return e
}

// WithTarget sets the target resource ID
func (e *Event) WithTarget(targetID uuid.UUID) *Event {
	e.TargetID = &targetID
	return e
}

// WithData sets arbitrary data for the event
func (e *Event) WithData(key string, value any) *Event {
	e.Data[key] = value
	return e
}

// WithMetadata sets metadata for the event
func (e *Event) WithMetadata(key, value string) *Event {
	e.Metadata[key] = value
	return e
}

// AgentType represents the type of notification agent
type AgentType string

const (
	AgentWebhook AgentType = "webhook"
	AgentDiscord AgentType = "discord"
	AgentEmail   AgentType = "email"
	AgentGotify  AgentType = "gotify"
	AgentNtfy    AgentType = "ntfy"
)

// String returns the string representation of the agent type
func (a AgentType) String() string {
	return string(a)
}

// Agent is the interface that all notification agents must implement
type Agent interface {
	// Type returns the agent type
	Type() AgentType

	// Name returns a human-readable name for this agent instance
	Name() string

	// Send sends a notification for the given event
	// Returns an error if the notification could not be sent
	Send(ctx context.Context, event *Event) error

	// Validate checks if the agent configuration is valid
	Validate() error

	// IsEnabled returns whether this agent is enabled
	IsEnabled() bool
}

// AgentConfig is the base configuration for all agents
type AgentConfig struct {
	Enabled         bool            `json:"enabled"`
	Name            string          `json:"name"`
	EventTypes      []EventType     `json:"event_types,omitempty"`      // Empty = all events
	EventCategories []EventCategory `json:"event_categories,omitempty"` // Empty = all categories
}

// ShouldSend checks if this agent should send for the given event type
func (c *AgentConfig) ShouldSend(eventType EventType) bool {
	if !c.Enabled {
		return false
	}

	// If no specific types or categories configured, send all
	if len(c.EventTypes) == 0 && len(c.EventCategories) == 0 {
		return true
	}

	// Check if event type is explicitly listed
	if slices.Contains(c.EventTypes, eventType) {
		return true
	}

	// Check if event category is listed
	eventCategory := eventType.GetCategory()
	return slices.Contains(c.EventCategories, eventCategory)
}

// NotificationResult represents the result of a notification send attempt
type NotificationResult struct {
	AgentType AgentType `json:"agent_type"`
	AgentName string    `json:"agent_name"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	SentAt    time.Time `json:"sent_at"`
}

// Service is the notification service interface
type Service interface {
	// Dispatch sends an event to all configured notification agents
	// This is typically async - events are queued for delivery
	Dispatch(ctx context.Context, event *Event) error

	// DispatchSync sends an event synchronously and returns results
	DispatchSync(ctx context.Context, event *Event) ([]NotificationResult, error)

	// RegisterAgent registers a notification agent
	RegisterAgent(agent Agent) error

	// UnregisterAgent removes an agent by name
	UnregisterAgent(name string) error

	// ListAgents returns all registered agents
	ListAgents() []Agent

	// TestAgent tests a specific agent with a test event
	TestAgent(ctx context.Context, agentName string) (*NotificationResult, error)
}
