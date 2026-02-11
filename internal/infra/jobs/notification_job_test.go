package jobs

import (
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/service/notification"
)

func TestNotificationArgs_Kind(t *testing.T) {
	t.Parallel()

	args := NotificationArgs{}
	assert.Equal(t, NotificationJobKind, args.Kind())
	assert.Equal(t, "notification", args.Kind())
}

func TestNotificationArgs_InsertOpts(t *testing.T) {
	t.Parallel()

	args := NotificationArgs{}
	opts := args.InsertOpts()

	assert.Equal(t, QueueHigh, opts.Queue)
	assert.Equal(t, 5, opts.MaxAttempts)
	assert.True(t, opts.UniqueOpts.ByArgs, "should deduplicate by args")
	assert.Equal(t, 1*time.Hour, opts.UniqueOpts.ByPeriod, "should deduplicate within 1 hour")
}

func TestNotificationArgs_ToEvent(t *testing.T) {
	t.Parallel()

	eventID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	targetID := uuid.Must(uuid.NewV7())
	timestamp := time.Now()

	args := &NotificationArgs{
		EventID:   eventID,
		EventType: "test.event",
		Timestamp: timestamp,
		UserID:    &userID,
		TargetID:  &targetID,
		Data: map[string]any{
			"key": "value",
		},
		Metadata: map[string]string{
			"source": "test",
		},
	}

	event := args.ToEvent()

	assert.Equal(t, eventID, event.ID)
	assert.Equal(t, notification.EventType("test.event"), event.Type)
	assert.Equal(t, timestamp, event.Timestamp)
	assert.Equal(t, &userID, event.UserID)
	assert.Equal(t, &targetID, event.TargetID)
	assert.Equal(t, "value", event.Data["key"])
	assert.Equal(t, "test", event.Metadata["source"])
}

func TestNotificationArgs_ToEvent_MinimalArgs(t *testing.T) {
	t.Parallel()

	eventID := uuid.Must(uuid.NewV7())
	timestamp := time.Now()

	args := &NotificationArgs{
		EventID:   eventID,
		EventType: "minimal.event",
		Timestamp: timestamp,
	}

	event := args.ToEvent()

	assert.Equal(t, eventID, event.ID)
	assert.Equal(t, notification.EventType("minimal.event"), event.Type)
	assert.Nil(t, event.UserID)
	assert.Nil(t, event.TargetID)
	assert.Nil(t, event.Data)
	assert.Nil(t, event.Metadata)
}

func TestNewNotificationArgs(t *testing.T) {
	t.Parallel()

	t.Run("with event and no agents", func(t *testing.T) {
		eventID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		event := &notification.Event{
			ID:        eventID,
			Type:      notification.EventType("user.created"),
			Timestamp: time.Now(),
			UserID:    &userID,
			Data: map[string]any{
				"username": "testuser",
			},
			Metadata: map[string]string{
				"ip": "127.0.0.1",
			},
		}

		args := NewNotificationArgs(event)

		assert.Equal(t, eventID, args.EventID)
		assert.Equal(t, "user.created", args.EventType)
		assert.Equal(t, event.Timestamp, args.Timestamp)
		assert.Equal(t, &userID, args.UserID)
		assert.Equal(t, "testuser", args.Data["username"])
		assert.Equal(t, "127.0.0.1", args.Metadata["ip"])
		assert.Empty(t, args.AgentNames)
	})

	t.Run("with specific agents", func(t *testing.T) {
		event := &notification.Event{
			ID:        uuid.Must(uuid.NewV7()),
			Type:      notification.EventType("movie.added"),
			Timestamp: time.Now(),
		}

		args := NewNotificationArgs(event, "discord", "webhook")

		assert.Len(t, args.AgentNames, 2)
		assert.Contains(t, args.AgentNames, "discord")
		assert.Contains(t, args.AgentNames, "webhook")
	})
}

func TestNewNotificationWorker(t *testing.T) {
	t.Parallel()

	t.Run("with nil logger uses default", func(t *testing.T) {
		worker := NewNotificationWorker(nil, nil)

		assert.NotNil(t, worker)
		assert.NotNil(t, worker.logger)
		assert.Nil(t, worker.dispatcher)
	})

	t.Run("with custom logger", func(t *testing.T) {
		logger := slog.Default()
		worker := NewNotificationWorker(nil, logger)

		assert.NotNil(t, worker)
		assert.NotNil(t, worker.logger)
	})
}

func TestNotificationWorker_Timeout(t *testing.T) {
	t.Parallel()

	worker := NewNotificationWorker(nil, nil)

	// Create a mock job (we just need to call Timeout)
	job := &river.Job[NotificationArgs]{
		Args: NotificationArgs{
			EventID:   uuid.Must(uuid.NewV7()),
			EventType: "test.event",
		},
	}

	timeout := worker.Timeout(job)

	assert.Equal(t, 2*time.Minute, timeout)
}

func TestNotificationJobResult_MarshalJSON(t *testing.T) {
	t.Parallel()

	result := NotificationJobResult{
		AgentResults: []AgentResult{
			{
				AgentName: "discord",
				AgentType: "discord",
				Success:   true,
				SentAt:    time.Now(),
			},
			{
				AgentName: "webhook",
				AgentType: "webhook",
				Success:   false,
				Error:     "connection refused",
				SentAt:    time.Now(),
			},
		},
		CompletedAt: time.Now(),
	}

	data, err := result.MarshalJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Verify it's valid JSON
	var parsed map[string]any
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	// Check structure
	agentResults, ok := parsed["agent_results"].([]any)
	require.True(t, ok)
	assert.Len(t, agentResults, 2)
}

func TestAgentResult_Structure(t *testing.T) {
	t.Parallel()

	result := AgentResult{
		AgentName: "test_agent",
		AgentType: "webhook",
		Success:   true,
		Error:     "",
		SentAt:    time.Now(),
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	var parsed AgentResult
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "test_agent", parsed.AgentName)
	assert.Equal(t, "webhook", parsed.AgentType)
	assert.True(t, parsed.Success)
	assert.Empty(t, parsed.Error)
}
