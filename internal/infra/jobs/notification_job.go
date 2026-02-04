package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"github.com/lusoris/revenge/internal/service/notification"
)

// NotificationJobKind is the unique identifier for notification dispatch jobs.
const NotificationJobKind = "notification"

// NotificationArgs defines the arguments for notification jobs.
type NotificationArgs struct {
	// EventID is the unique ID of the event
	EventID uuid.UUID `json:"event_id"`

	// EventType is the type of notification event
	EventType string `json:"event_type"`

	// Timestamp is when the event occurred
	Timestamp time.Time `json:"timestamp"`

	// UserID is the user who triggered the event (optional)
	UserID *uuid.UUID `json:"user_id,omitempty"`

	// TargetID is the target resource ID (optional)
	TargetID *uuid.UUID `json:"target_id,omitempty"`

	// Data contains event-specific data
	Data map[string]any `json:"data,omitempty"`

	// Metadata contains additional metadata
	Metadata map[string]string `json:"metadata,omitempty"`

	// AgentNames is an optional list of specific agents to send to
	// If empty, all enabled agents will be used
	AgentNames []string `json:"agent_names,omitempty"`
}

// Kind returns the job kind identifier.
func (NotificationArgs) Kind() string {
	return NotificationJobKind
}

// InsertOpts returns the default insert options for notification jobs.
// Notifications use QueueHigh for responsive user experience.
func (NotificationArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       QueueHigh,
		MaxAttempts: 5,
	}
}

// ToEvent converts NotificationArgs to a notification.Event
func (a *NotificationArgs) ToEvent() *notification.Event {
	return &notification.Event{
		ID:        a.EventID,
		Type:      notification.EventType(a.EventType),
		Timestamp: a.Timestamp,
		UserID:    a.UserID,
		TargetID:  a.TargetID,
		Data:      a.Data,
		Metadata:  a.Metadata,
	}
}

// NewNotificationArgs creates NotificationArgs from a notification.Event
func NewNotificationArgs(event *notification.Event, agentNames ...string) NotificationArgs {
	return NotificationArgs{
		EventID:    event.ID,
		EventType:  event.Type.String(),
		Timestamp:  event.Timestamp,
		UserID:     event.UserID,
		TargetID:   event.TargetID,
		Data:       event.Data,
		Metadata:   event.Metadata,
		AgentNames: agentNames,
	}
}

// NotificationWorker dispatches notification events to configured agents.
type NotificationWorker struct {
	river.WorkerDefaults[NotificationArgs]
	dispatcher *notification.Dispatcher
	logger     *slog.Logger
}

// NewNotificationWorker creates a new notification worker.
func NewNotificationWorker(dispatcher *notification.Dispatcher, logger *slog.Logger) *NotificationWorker {
	if logger == nil {
		logger = slog.Default()
	}
	return &NotificationWorker{
		dispatcher: dispatcher,
		logger:     logger.With("component", "notification_worker"),
	}
}

// Work executes the notification dispatch job.
func (w *NotificationWorker) Work(ctx context.Context, job *river.Job[NotificationArgs]) error {
	args := job.Args

	w.logger.Info("processing notification job",
		"job_id", job.ID,
		"event_id", args.EventID,
		"event_type", args.EventType,
		"attempt", job.Attempt,
	)

	// Convert args to event
	event := args.ToEvent()

	// Get agents to notify
	var agents []notification.Agent
	if len(args.AgentNames) > 0 {
		// Specific agents requested
		for _, name := range args.AgentNames {
			if agent, ok := w.dispatcher.GetAgent(name); ok && agent.IsEnabled() {
				agents = append(agents, agent)
			}
		}
	} else {
		// All enabled agents
		for _, agent := range w.dispatcher.ListAgents() {
			if agent.IsEnabled() {
				agents = append(agents, agent)
			}
		}
	}

	if len(agents) == 0 {
		w.logger.Warn("no enabled agents for notification",
			"event_id", args.EventID,
			"event_type", args.EventType,
		)
		return nil // Not an error, just nothing to do
	}

	// Send to each agent
	var lastErr error
	successCount := 0

	for _, agent := range agents {
		// Create timeout context for each agent
		agentCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		if err := agent.Send(agentCtx, event); err != nil {
			w.logger.Error("failed to send notification to agent",
				"agent", agent.Name(),
				"agent_type", agent.Type(),
				"event_type", args.EventType,
				"error", err,
			)
			lastErr = err
		} else {
			successCount++
			w.logger.Debug("notification sent successfully",
				"agent", agent.Name(),
				"event_type", args.EventType,
			)
		}

		cancel()
	}

	w.logger.Info("notification job completed",
		"job_id", job.ID,
		"event_type", args.EventType,
		"agents_total", len(agents),
		"agents_success", successCount,
	)

	// If all agents failed, return the last error for retry
	if successCount == 0 && lastErr != nil {
		return fmt.Errorf("all notification agents failed: %w", lastErr)
	}

	return nil
}

// Timeout returns the maximum execution time for notification jobs.
func (w *NotificationWorker) Timeout(job *river.Job[NotificationArgs]) time.Duration {
	// Allow 2 minutes total (for multiple agents with 30s timeout each)
	return 2 * time.Minute
}

// NotificationJobResult is stored with the completed job for auditing
type NotificationJobResult struct {
	AgentResults []AgentResult `json:"agent_results"`
	CompletedAt  time.Time     `json:"completed_at"`
}

// AgentResult records the outcome for a single agent
type AgentResult struct {
	AgentName string    `json:"agent_name"`
	AgentType string    `json:"agent_type"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	SentAt    time.Time `json:"sent_at"`
}

// MarshalJSON implements json.Marshaler for NotificationJobResult
func (r NotificationJobResult) MarshalJSON() ([]byte, error) {
	type Alias NotificationJobResult
	return json.Marshal((Alias)(r))
}
