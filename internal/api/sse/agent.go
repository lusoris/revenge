package sse

import (
	"context"

	"github.com/lusoris/revenge/internal/service/notification"
)

const AgentTypeSSE notification.AgentType = "sse"

// Agent is a notification agent that broadcasts events to SSE clients.
type Agent struct {
	broker *Broker
}

// NewAgent creates a new SSE notification agent.
func NewAgent(broker *Broker) *Agent {
	return &Agent{broker: broker}
}

func (a *Agent) Type() notification.AgentType { return AgentTypeSSE }
func (a *Agent) Name() string                   { return "sse-broadcast" }
func (a *Agent) Validate() error                { return nil }
func (a *Agent) IsEnabled() bool                { return true }

// Send broadcasts an event to all connected SSE clients.
func (a *Agent) Send(_ context.Context, event *notification.Event) error {
	a.broker.Broadcast(event)
	return nil
}
