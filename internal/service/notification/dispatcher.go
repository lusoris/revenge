package notification

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Dispatcher routes events to registered notification agents
type Dispatcher struct {
	mu     sync.RWMutex
	agents map[string]Agent
	logger *slog.Logger
	wg     sync.WaitGroup   // Track goroutines for graceful shutdown
	stopCh chan struct{}    // Signal shutdown to goroutines
}

// NewDispatcher creates a new notification dispatcher
func NewDispatcher(logger *slog.Logger) *Dispatcher {
	if logger == nil {
		logger = slog.Default()
	}
	return &Dispatcher{
		agents: make(map[string]Agent),
		logger: logger.With("component", "notification_dispatcher"),
		stopCh: make(chan struct{}),
	}
}

// RegisterAgent registers a notification agent
func (d *Dispatcher) RegisterAgent(agent Agent) error {
	if agent == nil {
		return fmt.Errorf("agent cannot be nil")
	}

	if err := agent.Validate(); err != nil {
		return fmt.Errorf("invalid agent configuration: %w", err)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	name := agent.Name()
	if _, exists := d.agents[name]; exists {
		return fmt.Errorf("agent with name %q already registered", name)
	}

	d.agents[name] = agent
	d.logger.Info("registered notification agent",
		"name", name,
		"type", agent.Type(),
		"enabled", agent.IsEnabled(),
	)

	return nil
}

// UnregisterAgent removes an agent by name
func (d *Dispatcher) UnregisterAgent(name string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.agents[name]; !exists {
		return fmt.Errorf("agent with name %q not found", name)
	}

	delete(d.agents, name)
	d.logger.Info("unregistered notification agent", "name", name)

	return nil
}

// ListAgents returns all registered agents
func (d *Dispatcher) ListAgents() []Agent {
	d.mu.RLock()
	defer d.mu.RUnlock()

	agents := make([]Agent, 0, len(d.agents))
	for _, agent := range d.agents {
		agents = append(agents, agent)
	}
	return agents
}

// GetAgent returns an agent by name
func (d *Dispatcher) GetAgent(name string) (Agent, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	agent, exists := d.agents[name]
	return agent, exists
}

// Dispatch sends an event to all enabled agents (async, fire-and-forget)
// For production use, this should enqueue to River job queue
func (d *Dispatcher) Dispatch(ctx context.Context, event *Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	d.logger.Debug("dispatching event",
		"event_id", event.ID,
		"event_type", event.Type,
	)

	// Get snapshot of agents
	d.mu.RLock()
	agents := make([]Agent, 0, len(d.agents))
	for _, agent := range d.agents {
		if agent.IsEnabled() {
			agents = append(agents, agent)
		}
	}
	d.mu.RUnlock()

	// Dispatch asynchronously with goroutine tracking for graceful shutdown
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()

		for _, agent := range agents {
			// Check for shutdown signal before processing each agent
			select {
			case <-d.stopCh:
				d.logger.Info("dispatcher shutting down, skipping remaining notifications")
				return
			default:
			}

			// Create new context with timeout for each agent
			agentCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			if err := agent.Send(agentCtx, event); err != nil {
				d.logger.Error("failed to send notification",
					"agent", agent.Name(),
					"agent_type", agent.Type(),
					"event_type", event.Type,
					"error", err,
				)
			} else {
				d.logger.Debug("notification sent",
					"agent", agent.Name(),
					"event_type", event.Type,
				)
			}

			cancel()
		}
	}()

	return nil
}

// DispatchSync sends an event synchronously and returns results for all agents
func (d *Dispatcher) DispatchSync(ctx context.Context, event *Event) ([]NotificationResult, error) {
	if event == nil {
		return nil, fmt.Errorf("event cannot be nil")
	}

	d.logger.Debug("dispatching event synchronously",
		"event_id", event.ID,
		"event_type", event.Type,
	)

	// Get snapshot of enabled agents
	d.mu.RLock()
	agents := make([]Agent, 0, len(d.agents))
	for _, agent := range d.agents {
		if agent.IsEnabled() {
			agents = append(agents, agent)
		}
	}
	d.mu.RUnlock()

	results := make([]NotificationResult, 0, len(agents))
	var wg sync.WaitGroup
	resultChan := make(chan NotificationResult, len(agents))

	for _, agent := range agents {
		wg.Add(1)
		go func(a Agent) {
			defer wg.Done()

			result := NotificationResult{
				AgentType: a.Type(),
				AgentName: a.Name(),
				SentAt:    time.Now().UTC(),
			}

			// Create context with timeout
			agentCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			if err := a.Send(agentCtx, event); err != nil {
				result.Success = false
				result.Error = err.Error()
				d.logger.Error("failed to send notification",
					"agent", a.Name(),
					"agent_type", a.Type(),
					"event_type", event.Type,
					"error", err,
				)
			} else {
				result.Success = true
				d.logger.Debug("notification sent",
					"agent", a.Name(),
					"event_type", event.Type,
				)
			}

			resultChan <- result
		}(agent)
	}

	// Wait for all agents to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		results = append(results, result)
	}

	return results, nil
}

// TestAgent tests a specific agent with a test event
func (d *Dispatcher) TestAgent(ctx context.Context, agentName string) (*NotificationResult, error) {
	agent, exists := d.GetAgent(agentName)
	if !exists {
		return nil, fmt.Errorf("agent with name %q not found", agentName)
	}

	// Create test event
	testEvent := NewEvent(EventType("test.notification")).
		WithData("message", "This is a test notification from Revenge").
		WithData("agent_name", agentName).
		WithMetadata("test", "true")

	result := &NotificationResult{
		AgentType: agent.Type(),
		AgentName: agent.Name(),
		SentAt:    time.Now().UTC(),
	}

	// Send test notification
	agentCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := agent.Send(agentCtx, testEvent); err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	return result, nil
}

// Close gracefully shuts down the dispatcher by signaling all goroutines
// to stop and waiting for them to complete. This prevents goroutine leaks
// and ensures all in-flight notifications are processed or cancelled cleanly.
func (d *Dispatcher) Close() error {
	d.logger.Info("shutting down notification dispatcher")

	// Signal all goroutines to stop
	close(d.stopCh)

	// Wait for all goroutines to complete
	d.wg.Wait()

	d.logger.Info("notification dispatcher shutdown complete")
	return nil
}

// Ensure Dispatcher implements Service interface
var _ Service = (*Dispatcher)(nil)
