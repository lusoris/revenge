package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lusoris/revenge/internal/service/notification"
)

// WebhookConfig holds the configuration for a webhook notification agent
type WebhookConfig struct {
	notification.AgentConfig

	// URL is the webhook endpoint URL
	URL string `json:"url"`

	// Method is the HTTP method (default: POST)
	Method string `json:"method,omitempty"`

	// Headers are additional HTTP headers to send
	Headers map[string]string `json:"headers,omitempty"`

	// ContentType is the Content-Type header (default: application/json)
	ContentType string `json:"content_type,omitempty"`

	// AuthType is the authentication type (none, basic, bearer, header)
	AuthType string `json:"auth_type,omitempty"`

	// AuthValue is the authentication value (password for basic, token for bearer)
	AuthValue string `json:"auth_value,omitempty"`

	// AuthUsername is the username for basic auth
	AuthUsername string `json:"auth_username,omitempty"`

	// Timeout is the request timeout (default: 30s)
	Timeout time.Duration `json:"timeout,omitempty"`

	// RetryCount is the number of retries on failure (default: 3)
	RetryCount int `json:"retry_count,omitempty"`

	// PayloadTemplate is an optional Go template for custom payload formatting
	// If empty, the default JSON event payload is used
	PayloadTemplate string `json:"payload_template,omitempty"`
}

// WebhookPayload is the default payload sent to webhooks
type WebhookPayload struct {
	EventID   string            `json:"event_id"`
	EventType string            `json:"event_type"`
	Timestamp string            `json:"timestamp"`
	UserID    string            `json:"user_id,omitempty"`
	TargetID  string            `json:"target_id,omitempty"`
	Data      map[string]any    `json:"data,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Source    string            `json:"source"`
}

// WebhookAgent sends notifications to a generic webhook endpoint
type WebhookAgent struct {
	config WebhookConfig
	client *http.Client
}

// NewWebhookAgent creates a new webhook notification agent
func NewWebhookAgent(config WebhookConfig) (*WebhookAgent, error) {
	agent := &WebhookAgent{
		config: config,
	}

	// Set defaults
	if agent.config.Method == "" {
		agent.config.Method = http.MethodPost
	}
	if agent.config.ContentType == "" {
		agent.config.ContentType = "application/json"
	}
	if agent.config.Timeout == 0 {
		agent.config.Timeout = 30 * time.Second
	}
	if agent.config.RetryCount == 0 {
		agent.config.RetryCount = 3
	}

	// Create HTTP client with timeout
	agent.client = &http.Client{
		Timeout: agent.config.Timeout,
	}

	if err := agent.Validate(); err != nil {
		return nil, err
	}

	return agent, nil
}

// Type returns the agent type
func (a *WebhookAgent) Type() notification.AgentType {
	return notification.AgentWebhook
}

// Name returns the agent name
func (a *WebhookAgent) Name() string {
	if a.config.Name == "" {
		return "webhook"
	}
	return a.config.Name
}

// IsEnabled returns whether the agent is enabled
func (a *WebhookAgent) IsEnabled() bool {
	return a.config.Enabled
}

// Validate checks if the configuration is valid
func (a *WebhookAgent) Validate() error {
	if a.config.URL == "" {
		return fmt.Errorf("webhook URL is required")
	}

	// Validate URL format
	if len(a.config.URL) < 10 {
		return fmt.Errorf("invalid webhook URL: too short")
	}

	// Validate HTTP method
	switch a.config.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		// Valid methods
	default:
		return fmt.Errorf("invalid HTTP method: %s (must be POST, PUT, or PATCH)", a.config.Method)
	}

	// Validate auth type
	switch a.config.AuthType {
	case "", "none", "basic", "bearer", "header":
		// Valid auth types
	default:
		return fmt.Errorf("invalid auth type: %s", a.config.AuthType)
	}

	return nil
}

// Send sends a notification to the webhook
func (a *WebhookAgent) Send(ctx context.Context, event *notification.Event) error {
	if !a.config.ShouldSend(event.Type) {
		return nil // Skip this event type
	}

	// Build payload
	payload, err := a.buildPayload(event)
	if err != nil {
		return fmt.Errorf("failed to build payload: %w", err)
	}

	// Send with retries
	var lastErr error
	for attempt := 0; attempt <= a.config.RetryCount; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s...
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		if err := a.sendRequest(ctx, payload); err != nil {
			lastErr = err
			continue
		}

		return nil // Success
	}

	return fmt.Errorf("webhook delivery failed after %d attempts: %w", a.config.RetryCount+1, lastErr)
}

// buildPayload constructs the webhook payload
func (a *WebhookAgent) buildPayload(event *notification.Event) ([]byte, error) {
	// TODO: Support custom PayloadTemplate with Go templates

	payload := WebhookPayload{
		EventID:   event.ID.String(),
		EventType: event.Type.String(),
		Timestamp: event.Timestamp.Format(time.RFC3339),
		Data:      event.Data,
		Metadata:  event.Metadata,
		Source:    "revenge",
	}

	if event.UserID != nil {
		payload.UserID = event.UserID.String()
	}
	if event.TargetID != nil {
		payload.TargetID = event.TargetID.String()
	}

	return json.Marshal(payload)
}

// sendRequest sends the HTTP request
func (a *WebhookAgent) sendRequest(ctx context.Context, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, a.config.Method, a.config.URL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", a.config.ContentType)
	req.Header.Set("User-Agent", "Revenge/1.0")

	for key, value := range a.config.Headers {
		req.Header.Set(key, value)
	}

	// Set authentication
	switch a.config.AuthType {
	case "basic":
		req.SetBasicAuth(a.config.AuthUsername, a.config.AuthValue)
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+a.config.AuthValue)
	case "header":
		// Custom header auth - AuthUsername is header name, AuthValue is header value
		if a.config.AuthUsername != "" {
			req.Header.Set(a.config.AuthUsername, a.config.AuthValue)
		}
	}

	// Send request
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Ensure WebhookAgent implements Agent interface
var _ notification.Agent = (*WebhookAgent)(nil)
