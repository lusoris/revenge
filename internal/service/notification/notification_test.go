package notification

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEventType_GetCategory(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  EventCategory
	}{
		{EventMovieAdded, CategoryContent},
		{EventMovieAvailable, CategoryContent},
		{EventMovieUpdated, CategoryContent},
		{EventMovieDeleted, CategoryContent},
		{EventRequestCreated, CategoryRequests},
		{EventRequestApproved, CategoryRequests},
		{EventRequestDenied, CategoryRequests},
		{EventLibraryScanStarted, CategoryLibrary},
		{EventLibraryScanDone, CategoryLibrary},
		{EventUserCreated, CategoryUser},
		{EventLoginSuccess, CategoryAuth},
		{EventLoginFailed, CategoryAuth},
		{EventMFAEnabled, CategoryAuth},
		{EventPasswordChanged, CategoryAuth},
		{EventPlaybackStarted, CategoryPlayback},
		{EventPlaybackStopped, CategoryPlayback},
		{EventSystemStartup, CategorySystem},
		{EventRadarrSync, CategorySystem},
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.eventType.GetCategory())
		})
	}
}

func TestNewEvent(t *testing.T) {
	event := NewEvent(EventMovieAdded)

	assert.NotEqual(t, uuid.Nil, event.ID)
	assert.Equal(t, EventMovieAdded, event.Type)
	assert.NotZero(t, event.Timestamp)
	assert.NotNil(t, event.Data)
	assert.NotNil(t, event.Metadata)
	assert.Nil(t, event.UserID)
	assert.Nil(t, event.TargetID)
}

func TestEvent_WithUser(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	event := NewEvent(EventMovieAdded).WithUser(userID)

	assert.NotNil(t, event.UserID)
	assert.Equal(t, userID, *event.UserID)
}

func TestEvent_WithTarget(t *testing.T) {
	targetID := uuid.Must(uuid.NewV7())
	event := NewEvent(EventMovieAdded).WithTarget(targetID)

	assert.NotNil(t, event.TargetID)
	assert.Equal(t, targetID, *event.TargetID)
}

func TestEvent_WithData(t *testing.T) {
	event := NewEvent(EventMovieAdded).
		WithData("movie_title", "Test Movie").
		WithData("year", 2024)

	assert.Equal(t, "Test Movie", event.Data["movie_title"])
	assert.Equal(t, 2024, event.Data["year"])
}

func TestEvent_WithMetadata(t *testing.T) {
	event := NewEvent(EventMovieAdded).
		WithMetadata("source", "radarr").
		WithMetadata("library_id", "123")

	assert.Equal(t, "radarr", event.Metadata["source"])
	assert.Equal(t, "123", event.Metadata["library_id"])
}

func TestAgentConfig_ShouldSend(t *testing.T) {
	tests := []struct {
		name       string
		config     AgentConfig
		eventType  EventType
		shouldSend bool
	}{
		{
			name:       "disabled agent never sends",
			config:     AgentConfig{Enabled: false},
			eventType:  EventMovieAdded,
			shouldSend: false,
		},
		{
			name:       "enabled with no filters sends all",
			config:     AgentConfig{Enabled: true},
			eventType:  EventMovieAdded,
			shouldSend: true,
		},
		{
			name: "enabled with matching event type sends",
			config: AgentConfig{
				Enabled:    true,
				EventTypes: []EventType{EventMovieAdded, EventMovieAvailable},
			},
			eventType:  EventMovieAdded,
			shouldSend: true,
		},
		{
			name: "enabled with non-matching event type does not send",
			config: AgentConfig{
				Enabled:    true,
				EventTypes: []EventType{EventMovieAdded},
			},
			eventType:  EventLoginFailed,
			shouldSend: false,
		},
		{
			name: "enabled with matching category sends",
			config: AgentConfig{
				Enabled:         true,
				EventCategories: []EventCategory{CategoryContent},
			},
			eventType:  EventMovieAdded,
			shouldSend: true,
		},
		{
			name: "enabled with non-matching category does not send",
			config: AgentConfig{
				Enabled:         true,
				EventCategories: []EventCategory{CategoryAuth},
			},
			eventType:  EventMovieAdded,
			shouldSend: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.shouldSend, tt.config.ShouldSend(tt.eventType))
		})
	}
}

func TestAgentType_String(t *testing.T) {
	assert.Equal(t, "webhook", AgentWebhook.String())
	assert.Equal(t, "discord", AgentDiscord.String())
	assert.Equal(t, "email", AgentEmail.String())
	assert.Equal(t, "gotify", AgentGotify.String())
	assert.Equal(t, "ntfy", AgentNtfy.String())
}

func TestEventType_String(t *testing.T) {
	assert.Equal(t, "movie.added", EventMovieAdded.String())
	assert.Equal(t, "auth.login_failed", EventLoginFailed.String())
}
