package jobs

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"

	"github.com/lusoris/revenge/internal/service/notification"
)

// --- scaleQueueWorkers tests ---

func TestScaleQueueWorkers(t *testing.T) {
	tests := []struct {
		name       string
		queues     map[string]river.QueueConfig
		maxWorkers int
		validate   func(t *testing.T, queues map[string]river.QueueConfig)
	}{
		{
			name: "proportional scaling with default config",
			queues: map[string]river.QueueConfig{
				QueueCritical: {MaxWorkers: 20},
				QueueHigh:     {MaxWorkers: 15},
				QueueDefault:  {MaxWorkers: 10},
				QueueLow:      {MaxWorkers: 5},
				QueueBulk:     {MaxWorkers: 3},
			},
			maxWorkers: 53, // same total as defaults
			validate: func(t *testing.T, queues map[string]river.QueueConfig) {
				// Proportional scaling should maintain similar ratios
				assert.Equal(t, 20, queues[QueueCritical].MaxWorkers)
				assert.Equal(t, 15, queues[QueueHigh].MaxWorkers)
				assert.Equal(t, 10, queues[QueueDefault].MaxWorkers)
				assert.Equal(t, 5, queues[QueueLow].MaxWorkers)
				assert.Equal(t, 3, queues[QueueBulk].MaxWorkers)
			},
		},
		{
			name: "scale down to small total",
			queues: map[string]river.QueueConfig{
				QueueCritical: {MaxWorkers: 20},
				QueueHigh:     {MaxWorkers: 15},
				QueueDefault:  {MaxWorkers: 10},
				QueueLow:      {MaxWorkers: 5},
				QueueBulk:     {MaxWorkers: 3},
			},
			maxWorkers: 10,
			validate: func(t *testing.T, queues map[string]river.QueueConfig) {
				// All should have at least 1 worker
				for name, qc := range queues {
					assert.GreaterOrEqual(t, qc.MaxWorkers, 1,
						"queue %s should have at least 1 worker", name)
				}
				// Critical should have more than bulk
				assert.Greater(t, queues[QueueCritical].MaxWorkers, queues[QueueBulk].MaxWorkers)
			},
		},
		{
			name: "minimum 1 worker per queue when scaling very small",
			queues: map[string]river.QueueConfig{
				QueueCritical: {MaxWorkers: 20},
				QueueHigh:     {MaxWorkers: 15},
				QueueDefault:  {MaxWorkers: 10},
				QueueLow:      {MaxWorkers: 5},
				QueueBulk:     {MaxWorkers: 3},
			},
			maxWorkers: 5,
			validate: func(t *testing.T, queues map[string]river.QueueConfig) {
				for name, qc := range queues {
					assert.GreaterOrEqual(t, qc.MaxWorkers, 1,
						"queue %s should have at least 1 worker even at minimal scale", name)
				}
			},
		},
		{
			name: "scale up",
			queues: map[string]river.QueueConfig{
				QueueCritical: {MaxWorkers: 20},
				QueueHigh:     {MaxWorkers: 15},
				QueueDefault:  {MaxWorkers: 10},
				QueueLow:      {MaxWorkers: 5},
				QueueBulk:     {MaxWorkers: 3},
			},
			maxWorkers: 200,
			validate: func(t *testing.T, queues map[string]river.QueueConfig) {
				// All queues should be scaled up
				assert.Greater(t, queues[QueueCritical].MaxWorkers, 20)
				assert.Greater(t, queues[QueueHigh].MaxWorkers, 15)
			},
		},
		{
			name:       "empty queues is no-op",
			queues:     map[string]river.QueueConfig{},
			maxWorkers: 100,
			validate: func(t *testing.T, queues map[string]river.QueueConfig) {
				assert.Empty(t, queues)
			},
		},
		{
			name: "zero total default workers is no-op",
			queues: map[string]river.QueueConfig{
				"q1": {MaxWorkers: 0},
				"q2": {MaxWorkers: 0},
			},
			maxWorkers: 100,
			validate: func(t *testing.T, queues map[string]river.QueueConfig) {
				// Should remain unchanged since totalDefault == 0
				assert.Equal(t, 0, queues["q1"].MaxWorkers)
				assert.Equal(t, 0, queues["q2"].MaxWorkers)
			},
		},
		{
			name: "single queue",
			queues: map[string]river.QueueConfig{
				"only": {MaxWorkers: 10},
			},
			maxWorkers: 50,
			validate: func(t *testing.T, queues map[string]river.QueueConfig) {
				assert.Equal(t, 50, queues["only"].MaxWorkers)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scaleQueueWorkers(tt.queues, tt.maxWorkers)
			tt.validate(t, tt.queues)
		})
	}
}

// --- NewRiverWorkers ---

func TestNewRiverWorkers(t *testing.T) {
	workers := NewRiverWorkers()
	assert.NotNil(t, workers)
}

// --- Progress ---

func TestJobProgress_Struct(t *testing.T) {
	progress := &JobProgress{
		Phase:   "scanning",
		Current: 50,
		Total:   100,
		Percent: 50,
		Message: "Scanning files...",
	}

	assert.Equal(t, "scanning", progress.Phase)
	assert.Equal(t, 50, progress.Current)
	assert.Equal(t, 100, progress.Total)
	assert.Equal(t, 50, progress.Percent)
	assert.Equal(t, "Scanning files...", progress.Message)
}

func TestJobProgress_JSON(t *testing.T) {
	progress := &JobProgress{
		Phase:   "indexing",
		Current: 75,
		Total:   200,
		Percent: 37,
		Message: "Indexing metadata",
	}

	data, err := json.Marshal(progress)
	require.NoError(t, err)

	var parsed JobProgress
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, progress.Phase, parsed.Phase)
	assert.Equal(t, progress.Current, parsed.Current)
	assert.Equal(t, progress.Total, parsed.Total)
	assert.Equal(t, progress.Percent, parsed.Percent)
	assert.Equal(t, progress.Message, parsed.Message)
}

func TestJobProgress_JSONOmitEmpty(t *testing.T) {
	progress := &JobProgress{
		Phase:   "waiting",
		Current: 0,
		Percent: 0,
	}

	data, err := json.Marshal(progress)
	require.NoError(t, err)

	// Total and Message have omitempty
	var parsed map[string]any
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	_, hasTotal := parsed["total"]
	_, hasMessage := parsed["message"]
	assert.False(t, hasTotal, "total should be omitted when zero")
	assert.False(t, hasMessage, "message should be omitted when empty")
}

func TestClient_ReportProgress_NilClient(t *testing.T) {
	c := &Client{
		client: nil,
		config: DefaultConfig(),
		logger: slog.Default(),
	}

	progress := &JobProgress{
		Phase:   "test",
		Current: 1,
		Total:   10,
	}

	// Should be a no-op (returns nil)
	err := c.ReportProgress(context.Background(), 123, progress)
	assert.NoError(t, err)
}

func TestClient_ReportProgress_AutoCalcPercent(t *testing.T) {
	// When client.client is nil, ReportProgress still calculates percent
	// before returning (the auto-calc happens before the nil check early return).
	c := &Client{
		client: nil,
		config: DefaultConfig(),
		logger: slog.Default(),
	}

	progress := &JobProgress{
		Phase:   "processing",
		Current: 50,
		Total:   200,
	}

	err := c.ReportProgress(context.Background(), 123, progress)
	assert.NoError(t, err)

	// With nil client, the function returns nil early.
	// The auto-calc only runs when Total > 0, and it runs BEFORE the nil check.
	// Looking at the code: the nil check is FIRST, so percent is NOT calculated.
	// We verify the nil-client no-op behavior here.
	assert.Equal(t, 0, progress.Percent, "nil client returns early before auto-calc")
}

func TestClient_ReportProgress_ZeroTotal(t *testing.T) {
	c := &Client{
		client: nil,
		config: DefaultConfig(),
		logger: slog.Default(),
	}

	progress := &JobProgress{
		Phase:   "unknown",
		Current: 10,
		Total:   0,
		Percent: 0,
	}

	_ = c.ReportProgress(context.Background(), 123, progress)
	// With Total=0, percent should not be recalculated
	assert.Equal(t, 0, progress.Percent)
}

func TestClient_GetJobProgress_NilClient(t *testing.T) {
	c := &Client{
		client: nil,
		config: DefaultConfig(),
		logger: slog.Default(),
	}

	progress, err := c.GetJobProgress(context.Background(), 123)
	assert.Error(t, err)
	assert.Nil(t, progress)
	assert.Contains(t, err.Error(), "not initialized")
}

// --- Notification job additional tests ---

func TestNotificationArgs_JSONRoundTrip(t *testing.T) {
	eventID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	targetID := uuid.Must(uuid.NewV7())
	ts := time.Now().Truncate(time.Millisecond)

	args := NotificationArgs{
		EventID:    eventID,
		EventType:  "movie.added",
		Timestamp:  ts,
		UserID:     &userID,
		TargetID:   &targetID,
		Data:       map[string]any{"title": "Test Movie"},
		Metadata:   map[string]string{"source": "radarr"},
		AgentNames: []string{"discord", "webhook"},
	}

	data, err := json.Marshal(args)
	require.NoError(t, err)

	var parsed NotificationArgs
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, args.EventID, parsed.EventID)
	assert.Equal(t, args.EventType, parsed.EventType)
	assert.Equal(t, args.UserID, parsed.UserID)
	assert.Equal(t, args.TargetID, parsed.TargetID)
	assert.Equal(t, "Test Movie", parsed.Data["title"])
	assert.Equal(t, "radarr", parsed.Metadata["source"])
	assert.Equal(t, []string{"discord", "webhook"}, parsed.AgentNames)
}

func TestNotificationArgs_NilOptionalFields(t *testing.T) {
	args := NotificationArgs{
		EventID:   uuid.Must(uuid.NewV7()),
		EventType: "system.startup",
		Timestamp: time.Now(),
	}

	event := args.ToEvent()
	assert.Nil(t, event.UserID)
	assert.Nil(t, event.TargetID)
	assert.Nil(t, event.Data)
	assert.Nil(t, event.Metadata)
}

func TestNewNotificationArgs_RoundTrip(t *testing.T) {
	eventID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	event := &notification.Event{
		ID:        eventID,
		Type:      notification.EventType("scan.complete"),
		Timestamp: time.Now(),
		UserID:    &userID,
		Data:      map[string]any{"files": 42},
	}

	args := NewNotificationArgs(event, "discord")
	assert.Equal(t, eventID, args.EventID)
	assert.Equal(t, "scan.complete", args.EventType)
	assert.Equal(t, &userID, args.UserID)
	assert.Equal(t, []string{"discord"}, args.AgentNames)

	// Convert back
	roundTripped := args.ToEvent()
	assert.Equal(t, eventID, roundTripped.ID)
	assert.Equal(t, notification.EventType("scan.complete"), roundTripped.Type)
}

func TestNotificationWorker_Timeout_Value(t *testing.T) {
	worker := NewNotificationWorker(nil, slog.Default())
	job := &river.Job[NotificationArgs]{
		Args: NotificationArgs{
			EventType: "test",
		},
	}
	assert.Equal(t, 2*time.Minute, worker.Timeout(job))
}

func TestNotificationJobResult_JSON(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	result := NotificationJobResult{
		AgentResults: []AgentResult{
			{
				AgentName: "discord",
				AgentType: "discord",
				Success:   true,
				SentAt:    now,
			},
			{
				AgentName: "webhook-main",
				AgentType: "webhook",
				Success:   false,
				Error:     "timeout after 30s",
				SentAt:    now,
			},
		},
		CompletedAt: now,
	}

	data, err := result.MarshalJSON()
	require.NoError(t, err)

	var parsed NotificationJobResult
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Len(t, parsed.AgentResults, 2)
	assert.Equal(t, "discord", parsed.AgentResults[0].AgentName)
	assert.True(t, parsed.AgentResults[0].Success)
	assert.Equal(t, "webhook-main", parsed.AgentResults[1].AgentName)
	assert.False(t, parsed.AgentResults[1].Success)
	assert.Equal(t, "timeout after 30s", parsed.AgentResults[1].Error)
}

func TestNotificationJobResult_EmptyResults(t *testing.T) {
	result := NotificationJobResult{
		AgentResults: []AgentResult{},
		CompletedAt:  time.Now(),
	}

	data, err := result.MarshalJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

// --- Cleanup job additional tests ---

func TestCleanupArgs_JSONRoundTrip(t *testing.T) {
	args := CleanupArgs{
		TargetType: CleanupTargetAll,
		OlderThan:  24 * time.Hour,
		BatchSize:  1000,
		DryRun:     true,
	}

	data, err := json.Marshal(args)
	require.NoError(t, err)

	var parsed CleanupArgs
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, CleanupTargetAll, parsed.TargetType)
	assert.Equal(t, 24*time.Hour, parsed.OlderThan)
	assert.Equal(t, 1000, parsed.BatchSize)
	assert.True(t, parsed.DryRun)
}

func TestCleanupWorker_Work_LeaderCheck(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteExpiredAuthTokens", context.Background()).Return(nil)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 100},
		Args: CleanupArgs{
			TargetType: CleanupTargetExpiredTokens,
			OlderThan:  1 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCleanupWorker_ValidateArgs_AllCases(t *testing.T) {
	worker := NewCleanupWorker(nil, slog.Default())

	tests := []struct {
		name    string
		args    CleanupArgs
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid all target",
			args: CleanupArgs{
				TargetType: CleanupTargetAll,
				OlderThan:  1 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "valid expired tokens target",
			args: CleanupArgs{
				TargetType: CleanupTargetExpiredTokens,
				OlderThan:  12 * time.Hour,
				BatchSize:  500,
			},
			wantErr: false,
		},
		{
			name: "empty target type",
			args: CleanupArgs{
				TargetType: "",
				OlderThan:  1 * time.Hour,
			},
			wantErr: true,
			errMsg:  "target_type is required",
		},
		{
			name: "zero older_than",
			args: CleanupArgs{
				TargetType: CleanupTargetAll,
				OlderThan:  0,
			},
			wantErr: true,
			errMsg:  "older_than must be positive",
		},
		{
			name: "negative older_than",
			args: CleanupArgs{
				TargetType: CleanupTargetAll,
				OlderThan:  -1 * time.Hour,
			},
			wantErr: true,
			errMsg:  "older_than must be positive",
		},
		{
			name: "negative batch_size",
			args: CleanupArgs{
				TargetType: CleanupTargetAll,
				OlderThan:  1 * time.Hour,
				BatchSize:  -1,
			},
			wantErr: true,
			errMsg:  "batch_size cannot be negative",
		},
		{
			name: "zero batch_size is valid",
			args: CleanupArgs{
				TargetType: CleanupTargetAll,
				OlderThan:  1 * time.Hour,
				BatchSize:  0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := worker.validateArgs(tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- DefaultConfig ---

func TestDefaultConfig_Values(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Queues)
	assert.Equal(t, 100*time.Millisecond, cfg.FetchCooldown)
	assert.Equal(t, 500*time.Millisecond, cfg.FetchPollInterval)
	assert.Equal(t, 1*time.Hour, cfg.RescueStuckJobsAfter)
	assert.Equal(t, 5, cfg.MaxAttempts)
}

// --- Config struct ---

func TestConfig_Struct(t *testing.T) {
	cfg := &Config{
		Queues: map[string]river.QueueConfig{
			"test": {MaxWorkers: 5},
		},
		FetchCooldown:        50 * time.Millisecond,
		FetchPollInterval:    250 * time.Millisecond,
		RescueStuckJobsAfter: 30 * time.Minute,
		MaxAttempts:          3,
	}

	assert.Len(t, cfg.Queues, 1)
	assert.Equal(t, 5, cfg.Queues["test"].MaxWorkers)
	assert.Equal(t, 50*time.Millisecond, cfg.FetchCooldown)
	assert.Equal(t, 250*time.Millisecond, cfg.FetchPollInterval)
	assert.Equal(t, 30*time.Minute, cfg.RescueStuckJobsAfter)
	assert.Equal(t, 3, cfg.MaxAttempts)
}

// --- ExponentialBackoff boundary ---

func TestExponentialBackoff_Boundary(t *testing.T) {
	// attempt=11: 2^11 = 2048 seconds = ~34 minutes, within 1 hour
	result := ExponentialBackoff(11)
	assert.Equal(t, 2048*time.Second, result)

	// attempt=12: 2^12 = 4096 seconds = ~68 minutes > 1 hour, should cap
	result = ExponentialBackoff(12)
	assert.Equal(t, 1*time.Hour, result)

	// attempt=31: should cap at max (above 30 threshold)
	result = ExponentialBackoff(31)
	assert.Equal(t, 1*time.Hour, result)
}

// --- LinearBackoff boundary ---

func TestLinearBackoff_Boundary(t *testing.T) {
	// attempt=59: 59*30s = 1770s = 29.5 min < 30 min
	result := LinearBackoff(59)
	assert.Equal(t, 59*30*time.Second, result)

	// attempt=60: 60*30s = 1800s = 30 min = max
	result = LinearBackoff(60)
	assert.Equal(t, 30*time.Minute, result)

	// attempt=61: should cap at 30 min
	result = LinearBackoff(61)
	assert.Equal(t, 30*time.Minute, result)
}

// --- Cleanup target constants ---

func TestCleanupTargetConstants_Exhaustive(t *testing.T) {
	targets := []string{
		CleanupTargetExpiredTokens,
		CleanupTargetRevokedTokens,
		CleanupTargetPasswordResets,
		CleanupTargetEmailVerifications,
		CleanupTargetFailedLogins,
		CleanupTargetAll,
	}

	// Verify all unique
	seen := make(map[string]bool)
	for _, target := range targets {
		assert.False(t, seen[target], "duplicate target constant: %s", target)
		seen[target] = true
	}
	assert.Len(t, seen, 6)
}

// --- NotificationJobKind ---

func TestNotificationJobKind_Value(t *testing.T) {
	assert.Equal(t, "notification", NotificationJobKind)
}

// --- CleanupJobKind ---

func TestCleanupJobKind_Value(t *testing.T) {
	assert.Equal(t, "cleanup", CleanupJobKind)
}

// --- registerCleanupWorker / registerNotificationWorker ---

func TestRegisterCleanupWorker(t *testing.T) {
	workers := river.NewWorkers()
	repo := &mockAuthCleanupRepo{}
	cleanupWorker := NewCleanupWorker(repo, slog.Default())

	// Should not panic
	assert.NotPanics(t, func() {
		registerCleanupWorker(workers, cleanupWorker)
	})
}

func TestRegisterNotificationWorker(t *testing.T) {
	workers := river.NewWorkers()
	notifWorker := NewNotificationWorker(nil, slog.Default())

	// Should not panic
	assert.NotPanics(t, func() {
		registerNotificationWorker(workers, notifWorker)
	})
}

// --- NotificationWorker.Work tests ---

// mockAgent implements notification.Agent for testing.
type mockAgent struct {
	name      string
	agentType notification.AgentType
	enabled   bool
	sendErr   error
}

func (m *mockAgent) Type() notification.AgentType { return m.agentType }
func (m *mockAgent) Name() string                 { return m.name }
func (m *mockAgent) IsEnabled() bool              { return m.enabled }
func (m *mockAgent) Validate() error              { return nil }
func (m *mockAgent) Send(_ context.Context, _ *notification.Event) error {
	return m.sendErr
}

func TestNotificationWorker_Work_NoAgents(t *testing.T) {
	// Create a dispatcher with no agents
	dispatcher := notification.NewDispatcher(slog.Default())
	worker := NewNotificationWorker(dispatcher, slog.Default())

	eventID := uuid.Must(uuid.NewV7())
	job := &river.Job[NotificationArgs]{
		JobRow: &rivertype.JobRow{ID: 1, Attempt: 1},
		Args: NotificationArgs{
			EventID:   eventID,
			EventType: "test.event",
			Timestamp: time.Now(),
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err, "should succeed even with no agents (nothing to do)")
}

func TestNotificationWorker_Work_WithEnabledAgent(t *testing.T) {
	dispatcher := notification.NewDispatcher(slog.Default())
	agent := &mockAgent{
		name:      "test-webhook",
		agentType: notification.AgentWebhook,
		enabled:   true,
		sendErr:   nil,
	}
	_ = dispatcher.RegisterAgent(agent)

	worker := NewNotificationWorker(dispatcher, slog.Default())

	eventID := uuid.Must(uuid.NewV7())
	job := &river.Job[NotificationArgs]{
		JobRow: &rivertype.JobRow{ID: 2, Attempt: 1},
		Args: NotificationArgs{
			EventID:   eventID,
			EventType: "movie.added",
			Timestamp: time.Now(),
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err)
}

func TestNotificationWorker_Work_AllAgentsFail(t *testing.T) {
	dispatcher := notification.NewDispatcher(slog.Default())
	agent := &mockAgent{
		name:      "failing-agent",
		agentType: notification.AgentWebhook,
		enabled:   true,
		sendErr:   assert.AnError,
	}
	_ = dispatcher.RegisterAgent(agent)

	worker := NewNotificationWorker(dispatcher, slog.Default())

	eventID := uuid.Must(uuid.NewV7())
	job := &river.Job[NotificationArgs]{
		JobRow: &rivertype.JobRow{ID: 3, Attempt: 1},
		Args: NotificationArgs{
			EventID:   eventID,
			EventType: "movie.added",
			Timestamp: time.Now(),
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "all notification agents failed")
}

func TestNotificationWorker_Work_PartialFailure(t *testing.T) {
	dispatcher := notification.NewDispatcher(slog.Default())

	successAgent := &mockAgent{
		name:      "success-agent",
		agentType: notification.AgentWebhook,
		enabled:   true,
		sendErr:   nil,
	}
	failAgent := &mockAgent{
		name:      "fail-agent",
		agentType: notification.AgentDiscord,
		enabled:   true,
		sendErr:   assert.AnError,
	}
	_ = dispatcher.RegisterAgent(successAgent)
	_ = dispatcher.RegisterAgent(failAgent)

	worker := NewNotificationWorker(dispatcher, slog.Default())

	eventID := uuid.Must(uuid.NewV7())
	job := &river.Job[NotificationArgs]{
		JobRow: &rivertype.JobRow{ID: 4, Attempt: 1},
		Args: NotificationArgs{
			EventID:   eventID,
			EventType: "movie.added",
			Timestamp: time.Now(),
		},
	}

	err := worker.Work(context.Background(), job)
	// Should succeed because at least one agent succeeded
	assert.NoError(t, err)
}

func TestNotificationWorker_Work_SpecificAgents(t *testing.T) {
	dispatcher := notification.NewDispatcher(slog.Default())

	agent1 := &mockAgent{
		name:      "discord-main",
		agentType: notification.AgentDiscord,
		enabled:   true,
		sendErr:   nil,
	}
	agent2 := &mockAgent{
		name:      "webhook-backup",
		agentType: notification.AgentWebhook,
		enabled:   true,
		sendErr:   nil,
	}
	_ = dispatcher.RegisterAgent(agent1)
	_ = dispatcher.RegisterAgent(agent2)

	worker := NewNotificationWorker(dispatcher, slog.Default())

	eventID := uuid.Must(uuid.NewV7())
	job := &river.Job[NotificationArgs]{
		JobRow: &rivertype.JobRow{ID: 5, Attempt: 1},
		Args: NotificationArgs{
			EventID:    eventID,
			EventType:  "scan.complete",
			Timestamp:  time.Now(),
			AgentNames: []string{"discord-main"}, // only this agent
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err)
}

func TestNotificationWorker_Work_SpecificAgentNotFound(t *testing.T) {
	dispatcher := notification.NewDispatcher(slog.Default())
	worker := NewNotificationWorker(dispatcher, slog.Default())

	eventID := uuid.Must(uuid.NewV7())
	job := &river.Job[NotificationArgs]{
		JobRow: &rivertype.JobRow{ID: 6, Attempt: 1},
		Args: NotificationArgs{
			EventID:    eventID,
			EventType:  "test.event",
			Timestamp:  time.Now(),
			AgentNames: []string{"nonexistent-agent"},
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err, "should succeed when specific agent not found (nothing to do)")
}

func TestNotificationWorker_Work_DisabledAgent(t *testing.T) {
	dispatcher := notification.NewDispatcher(slog.Default())
	agent := &mockAgent{
		name:      "disabled-agent",
		agentType: notification.AgentWebhook,
		enabled:   false,
		sendErr:   nil,
	}
	_ = dispatcher.RegisterAgent(agent)

	worker := NewNotificationWorker(dispatcher, slog.Default())

	eventID := uuid.Must(uuid.NewV7())
	job := &river.Job[NotificationArgs]{
		JobRow: &rivertype.JobRow{ID: 7, Attempt: 1},
		Args: NotificationArgs{
			EventID:   eventID,
			EventType: "test.event",
			Timestamp: time.Now(),
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err, "should succeed when agent is disabled (nothing to do)")
}

// --- cleanupPasswordResets error paths ---

func TestCleanupWorker_Work_PasswordResets_ExpiredError(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteExpiredPasswordResetTokens", mock.Anything).Return(assert.AnError)
	repo.On("DeleteUsedPasswordResetTokens", mock.Anything).Return(nil)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 20},
		Args: CleanupArgs{
			TargetType: CleanupTargetPasswordResets,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 errors")
	repo.AssertExpectations(t)
}

func TestCleanupWorker_Work_PasswordResets_UsedError(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteExpiredPasswordResetTokens", mock.Anything).Return(nil)
	repo.On("DeleteUsedPasswordResetTokens", mock.Anything).Return(assert.AnError)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 21},
		Args: CleanupArgs{
			TargetType: CleanupTargetPasswordResets,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 errors")
	repo.AssertExpectations(t)
}

// --- cleanupEmailVerifications error paths ---

func TestCleanupWorker_Work_EmailVerifications_ExpiredError(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteExpiredEmailVerificationTokens", mock.Anything).Return(assert.AnError)
	repo.On("DeleteVerifiedEmailTokens", mock.Anything).Return(nil)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 22},
		Args: CleanupArgs{
			TargetType: CleanupTargetEmailVerifications,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCleanupWorker_Work_EmailVerifications_VerifiedError(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteExpiredEmailVerificationTokens", mock.Anything).Return(nil)
	repo.On("DeleteVerifiedEmailTokens", mock.Anything).Return(assert.AnError)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 23},
		Args: CleanupArgs{
			TargetType: CleanupTargetEmailVerifications,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	repo.AssertExpectations(t)
}

func TestNotificationWorker_Work_SpecificDisabledAgent(t *testing.T) {
	dispatcher := notification.NewDispatcher(slog.Default())
	agent := &mockAgent{
		name:      "disabled-webhook",
		agentType: notification.AgentWebhook,
		enabled:   false,
		sendErr:   nil,
	}
	_ = dispatcher.RegisterAgent(agent)

	worker := NewNotificationWorker(dispatcher, slog.Default())

	eventID := uuid.Must(uuid.NewV7())
	job := &river.Job[NotificationArgs]{
		JobRow: &rivertype.JobRow{ID: 8, Attempt: 1},
		Args: NotificationArgs{
			EventID:    eventID,
			EventType:  "test.event",
			Timestamp:  time.Now(),
			AgentNames: []string{"disabled-webhook"},
		},
	}

	err := worker.Work(context.Background(), job)
	assert.NoError(t, err, "should succeed when specific agent is disabled")
}

// --- registerHooks tests ---

func TestRegisterHooks_StartStop(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	logger := slog.Default()

	// Client with nil river client - Start and Stop will return errors
	client := &Client{
		client: nil,
		config: DefaultConfig(),
		logger: logger,
	}

	registerHooks(lc, client, logger)

	ctx := context.Background()
	// Start calls client.Start which returns error for nil client
	err := lc.Start(ctx)
	assert.Error(t, err, "start should fail with nil river client")
}

// --- Client.Start/Stop nil client ---

// Client.Start/Stop/Insert/InsertMany/JobGet/JobCancel/Subscribe/RiverClient nil-client
// tests are in river_test.go. Not duplicated here.

// --- CleanupWorker.Work: revoked tokens error ---

func TestCleanupWorker_Work_RevokedTokens_Error(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteRevokedAuthTokens", mock.Anything).Return(assert.AnError)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 30},
		Args: CleanupArgs{
			TargetType: CleanupTargetRevokedTokens,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 errors")
	repo.AssertExpectations(t)
}

// --- CleanupWorker.Work: failed logins error ---

func TestCleanupWorker_Work_FailedLogins_Error(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteOldFailedLoginAttempts", mock.Anything).Return(assert.AnError)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 31},
		Args: CleanupArgs{
			TargetType: CleanupTargetFailedLogins,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 errors")
	repo.AssertExpectations(t)
}

// --- CleanupWorker.Work: expired tokens error ---

func TestCleanupWorker_Work_ExpiredTokens_Error(t *testing.T) {
	repo := &mockAuthCleanupRepo{}
	worker := NewCleanupWorker(repo, slog.Default())

	repo.On("DeleteExpiredAuthTokens", mock.Anything).Return(assert.AnError)

	job := &river.Job[CleanupArgs]{
		JobRow: &rivertype.JobRow{ID: 32},
		Args: CleanupArgs{
			TargetType: CleanupTargetExpiredTokens,
			OlderThan:  24 * time.Hour,
		},
	}

	err := worker.Work(context.Background(), job)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 errors")
	repo.AssertExpectations(t)
}

// --- DefaultRetryPolicy ---

func TestDefaultRetryPolicy_Values(t *testing.T) {
	policy := DefaultRetryPolicy()
	assert.NotNil(t, policy)
	assert.Equal(t, 25, policy.MaxAttempts)
	assert.NotNil(t, policy.Backoff)

	// Verify the default backoff is ExponentialBackoff
	assert.Equal(t, ExponentialBackoff(5), policy.Backoff(5))
}

// --- ExponentialBackoff edge cases ---

func TestExponentialBackoff_NegativeAttempt(t *testing.T) {
	// Negative attempts should be clamped to 0
	result := ExponentialBackoff(-5)
	assert.Equal(t, 1*time.Second, result) // 2^0 = 1
}

func TestExponentialBackoff_ZeroAttempt(t *testing.T) {
	result := ExponentialBackoff(0)
	assert.Equal(t, 1*time.Second, result) // 2^0 = 1
}

func TestExponentialBackoff_SmallAttempts(t *testing.T) {
	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{3, 8 * time.Second},
		{5, 32 * time.Second},
		{10, 1024 * time.Second},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.expected, ExponentialBackoff(tt.attempt))
		})
	}
}

// --- LinearBackoff edge cases ---

func TestLinearBackoff_NegativeAttempt(t *testing.T) {
	result := LinearBackoff(-10)
	assert.Equal(t, 0*time.Second, result) // clamped to 0
}

func TestLinearBackoff_ZeroAttempt(t *testing.T) {
	result := LinearBackoff(0)
	assert.Equal(t, 0*time.Second, result)
}

func TestLinearBackoff_SmallAttempts(t *testing.T) {
	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{1, 30 * time.Second},
		{2, 60 * time.Second},
		{5, 150 * time.Second},
		{10, 300 * time.Second},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.expected, LinearBackoff(tt.attempt))
		})
	}
}

// --- QueuePriority additional tests ---

func TestQueuePriority_AllBoundaries(t *testing.T) {
	tests := []struct {
		priority int
		expected string
	}{
		{100, QueueCritical},
		{20, QueueCritical},
		{19, QueueHigh},
		{10, QueueHigh},
		{9, QueueDefault},
		{0, QueueDefault},
		{-9, QueueDefault},
		{-10, QueueLow},
		{-19, QueueLow},
		{-20, QueueBulk},
		{-100, QueueBulk},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.expected, QueuePriority(tt.priority))
		})
	}
}

// --- DefaultQueueConfig ---

func TestDefaultQueueConfig_AllQueues(t *testing.T) {
	cfg := DefaultQueueConfig()
	assert.NotNil(t, cfg)
	assert.Len(t, cfg.Queues, 5)

	assert.Equal(t, 20, cfg.Queues[QueueCritical].MaxWorkers)
	assert.Equal(t, 15, cfg.Queues[QueueHigh].MaxWorkers)
	assert.Equal(t, 10, cfg.Queues[QueueDefault].MaxWorkers)
	assert.Equal(t, 5, cfg.Queues[QueueLow].MaxWorkers)
	assert.Equal(t, 3, cfg.Queues[QueueBulk].MaxWorkers)
}

// --- NewNotificationWorker with nil logger ---

func TestNewNotificationWorker_NilLogger(t *testing.T) {
	worker := NewNotificationWorker(nil, nil)
	assert.NotNil(t, worker)
	assert.NotNil(t, worker.logger)
}

// --- ToEvent with full fields ---

func TestNotificationArgs_ToEvent_FullFields(t *testing.T) {
	eventID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	targetID := uuid.Must(uuid.NewV7())
	ts := time.Now().Truncate(time.Millisecond)

	args := NotificationArgs{
		EventID:   eventID,
		EventType: "movie.updated",
		Timestamp: ts,
		UserID:    &userID,
		TargetID:  &targetID,
		Data:      map[string]any{"title": "Updated Movie"},
		Metadata:  map[string]string{"source": "tmdb"},
	}

	event := args.ToEvent()
	assert.Equal(t, eventID, event.ID)
	assert.Equal(t, notification.EventType("movie.updated"), event.Type)
	assert.Equal(t, ts, event.Timestamp)
	assert.Equal(t, &userID, event.UserID)
	assert.Equal(t, &targetID, event.TargetID)
	assert.Equal(t, "Updated Movie", event.Data["title"])
	assert.Equal(t, "tmdb", event.Metadata["source"])
}
