package notification

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockAgent is a mock notification agent for testing
type mockAgent struct {
	name       string
	agentType  AgentType
	enabled    bool
	sendFunc   func(ctx context.Context, event *Event) error
	sendCalls  int
	mu         sync.Mutex
}

func newMockAgent(name string, enabled bool) *mockAgent {
	return &mockAgent{
		name:      name,
		agentType: AgentWebhook,
		enabled:   enabled,
		sendFunc:  func(ctx context.Context, event *Event) error { return nil },
	}
}

func (m *mockAgent) Type() AgentType    { return m.agentType }
func (m *mockAgent) Name() string       { return m.name }
func (m *mockAgent) IsEnabled() bool    { return m.enabled }
func (m *mockAgent) Validate() error    { return nil }

func (m *mockAgent) Send(ctx context.Context, event *Event) error {
	m.mu.Lock()
	m.sendCalls++
	m.mu.Unlock()
	return m.sendFunc(ctx, event)
}

func (m *mockAgent) getSendCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sendCalls
}

func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher(nil)
	assert.NotNil(t, d)
	assert.NotNil(t, d.agents)
}

func TestDispatcher_RegisterAgent(t *testing.T) {
	d := NewDispatcher(slog.Default())

	agent := newMockAgent("test-agent", true)
	err := d.RegisterAgent(agent)
	require.NoError(t, err)

	// Verify agent is registered
	agents := d.ListAgents()
	assert.Len(t, agents, 1)
	assert.Equal(t, "test-agent", agents[0].Name())
}

func TestDispatcher_RegisterAgent_NilAgent(t *testing.T) {
	d := NewDispatcher(slog.Default())

	err := d.RegisterAgent(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestDispatcher_RegisterAgent_DuplicateName(t *testing.T) {
	d := NewDispatcher(slog.Default())

	agent1 := newMockAgent("test-agent", true)
	agent2 := newMockAgent("test-agent", true)

	err := d.RegisterAgent(agent1)
	require.NoError(t, err)

	err = d.RegisterAgent(agent2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestDispatcher_UnregisterAgent(t *testing.T) {
	d := NewDispatcher(slog.Default())

	agent := newMockAgent("test-agent", true)
	require.NoError(t, d.RegisterAgent(agent))

	err := d.UnregisterAgent("test-agent")
	require.NoError(t, err)

	agents := d.ListAgents()
	assert.Len(t, agents, 0)
}

func TestDispatcher_UnregisterAgent_NotFound(t *testing.T) {
	d := NewDispatcher(slog.Default())

	err := d.UnregisterAgent("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDispatcher_GetAgent(t *testing.T) {
	d := NewDispatcher(slog.Default())

	agent := newMockAgent("test-agent", true)
	require.NoError(t, d.RegisterAgent(agent))

	found, exists := d.GetAgent("test-agent")
	assert.True(t, exists)
	assert.Equal(t, "test-agent", found.Name())

	_, exists = d.GetAgent("nonexistent")
	assert.False(t, exists)
}

func TestDispatcher_DispatchSync(t *testing.T) {
	d := NewDispatcher(slog.Default())

	agent1 := newMockAgent("agent1", true)
	agent2 := newMockAgent("agent2", true)
	agent3 := newMockAgent("agent3", false) // Disabled

	require.NoError(t, d.RegisterAgent(agent1))
	require.NoError(t, d.RegisterAgent(agent2))
	require.NoError(t, d.RegisterAgent(agent3))

	event := NewEvent(EventMovieAdded).
		WithData("movie_title", "Test Movie")

	results, err := d.DispatchSync(context.Background(), event)
	require.NoError(t, err)

	// Only enabled agents should be called
	assert.Len(t, results, 2)
	assert.Equal(t, 1, agent1.getSendCalls())
	assert.Equal(t, 1, agent2.getSendCalls())
	assert.Equal(t, 0, agent3.getSendCalls())

	// All results should be successful
	for _, r := range results {
		assert.True(t, r.Success)
		assert.Empty(t, r.Error)
	}
}

func TestDispatcher_DispatchSync_WithErrors(t *testing.T) {
	d := NewDispatcher(slog.Default())

	agent1 := newMockAgent("agent1", true)
	agent2 := newMockAgent("agent2", true)

	agent2.sendFunc = func(ctx context.Context, event *Event) error {
		return errors.New("send failed")
	}

	require.NoError(t, d.RegisterAgent(agent1))
	require.NoError(t, d.RegisterAgent(agent2))

	event := NewEvent(EventMovieAdded)

	results, err := d.DispatchSync(context.Background(), event)
	require.NoError(t, err) // Dispatch itself succeeds

	// Check that one succeeded and one failed
	successCount := 0
	failCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
		} else {
			failCount++
			assert.Contains(t, r.Error, "send failed")
		}
	}
	assert.Equal(t, 1, successCount)
	assert.Equal(t, 1, failCount)
}

func TestDispatcher_DispatchSync_NilEvent(t *testing.T) {
	d := NewDispatcher(slog.Default())

	_, err := d.DispatchSync(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestDispatcher_Dispatch_Async(t *testing.T) {
	d := NewDispatcher(slog.Default())

	agent := newMockAgent("agent1", true)
	require.NoError(t, d.RegisterAgent(agent))

	event := NewEvent(EventMovieAdded)

	err := d.Dispatch(context.Background(), event)
	require.NoError(t, err)

	// Wait for async dispatch to complete
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, agent.getSendCalls())
}

func TestDispatcher_TestAgent(t *testing.T) {
	d := NewDispatcher(slog.Default())

	agent := newMockAgent("test-agent", true)
	require.NoError(t, d.RegisterAgent(agent))

	result, err := d.TestAgent(context.Background(), "test-agent")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Success)
	assert.Equal(t, "test-agent", result.AgentName)
	assert.Equal(t, 1, agent.getSendCalls())
}

func TestDispatcher_TestAgent_NotFound(t *testing.T) {
	d := NewDispatcher(slog.Default())

	_, err := d.TestAgent(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDispatcher_TestAgent_SendError(t *testing.T) {
	d := NewDispatcher(slog.Default())

	agent := newMockAgent("test-agent", true)
	agent.sendFunc = func(ctx context.Context, event *Event) error {
		return errors.New("test error")
	}
	require.NoError(t, d.RegisterAgent(agent))

	result, err := d.TestAgent(context.Background(), "test-agent")
	require.NoError(t, err) // TestAgent doesn't return error, just unsuccessful result
	require.NotNil(t, result)

	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "test error")
}
