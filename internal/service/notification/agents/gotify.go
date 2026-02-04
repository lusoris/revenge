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

// GotifyConfig holds the configuration for a Gotify notification agent
type GotifyConfig struct {
	notification.AgentConfig

	// ServerURL is the Gotify server URL (e.g., https://gotify.example.com)
	ServerURL string `json:"server_url"`

	// AppToken is the Gotify application token
	AppToken string `json:"app_token"`

	// DefaultPriority is the default message priority (0-10, default: 5)
	DefaultPriority int `json:"default_priority,omitempty"`

	// Timeout is the request timeout (default: 30s)
	Timeout time.Duration `json:"timeout,omitempty"`
}

// GotifyMessage is the payload sent to Gotify
type GotifyMessage struct {
	Title    string         `json:"title"`
	Message  string         `json:"message"`
	Priority int            `json:"priority"`
	Extras   map[string]any `json:"extras,omitempty"`
}

// GotifyAgent sends notifications to Gotify server
type GotifyAgent struct {
	config GotifyConfig
	client *http.Client
}

// NewGotifyAgent creates a new Gotify notification agent
func NewGotifyAgent(config GotifyConfig) (*GotifyAgent, error) {
	agent := &GotifyAgent{
		config: config,
	}

	// Set defaults
	if agent.config.DefaultPriority == 0 {
		agent.config.DefaultPriority = 5
	}
	if agent.config.Timeout == 0 {
		agent.config.Timeout = 30 * time.Second
	}

	// Create HTTP client
	agent.client = &http.Client{
		Timeout: agent.config.Timeout,
	}

	if err := agent.Validate(); err != nil {
		return nil, err
	}

	return agent, nil
}

// Type returns the agent type
func (a *GotifyAgent) Type() notification.AgentType {
	return notification.AgentGotify
}

// Name returns the agent name
func (a *GotifyAgent) Name() string {
	if a.config.Name == "" {
		return "gotify"
	}
	return a.config.Name
}

// IsEnabled returns whether the agent is enabled
func (a *GotifyAgent) IsEnabled() bool {
	return a.config.Enabled
}

// Validate checks if the configuration is valid
func (a *GotifyAgent) Validate() error {
	if a.config.ServerURL == "" {
		return fmt.Errorf("gotify server URL is required")
	}
	if a.config.AppToken == "" {
		return fmt.Errorf("gotify app token is required")
	}
	if a.config.DefaultPriority < 0 || a.config.DefaultPriority > 10 {
		return fmt.Errorf("priority must be between 0 and 10")
	}

	return nil
}

// Send sends a notification to Gotify
func (a *GotifyAgent) Send(ctx context.Context, event *notification.Event) error {
	if !a.config.ShouldSend(event.Type) {
		return nil // Skip this event type
	}

	// Build Gotify message
	msg := a.buildMessage(event)

	// Marshal to JSON
	jsonPayload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send request
	url := fmt.Sprintf("%s/message", a.config.ServerURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gotify-Key", a.config.AppToken)
	req.Header.Set("User-Agent", "Revenge/1.0")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("gotify returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// buildMessage creates the Gotify message
func (a *GotifyAgent) buildMessage(event *notification.Event) GotifyMessage {
	msg := GotifyMessage{
		Title:    a.getTitle(event),
		Message:  a.getMessage(event),
		Priority: a.getPriority(event),
	}

	// Add extras for rich notifications
	msg.Extras = map[string]any{
		"client::notification": map[string]any{
			"click": map[string]any{
				"url": a.getClickURL(event),
			},
		},
	}

	return msg
}

// getTitle returns the notification title
func (a *GotifyAgent) getTitle(event *notification.Event) string {
	if title, ok := event.Data["title"].(string); ok {
		return title
	}

	switch event.Type {
	case notification.EventMovieAdded:
		return "New Movie Added"
	case notification.EventMovieAvailable:
		return "Movie Now Available"
	case notification.EventRequestCreated:
		return "New Content Request"
	case notification.EventRequestApproved:
		return "Request Approved"
	case notification.EventRequestDenied:
		return "Request Denied"
	case notification.EventLibraryScanDone:
		return "Library Scan Complete"
	case notification.EventUserCreated:
		return "New User Created"
	case notification.EventLoginFailed:
		return "⚠️ Failed Login Attempt"
	case notification.EventMFAEnabled:
		return "MFA Enabled"
	case notification.EventPasswordChanged:
		return "Password Changed"
	case notification.EventPlaybackStarted:
		return "Playback Started"
	case notification.EventSystemStartup:
		return "Revenge Started"
	default:
		return string(event.Type)
	}
}

// getMessage returns the notification message body
func (a *GotifyAgent) getMessage(event *notification.Event) string {
	// Check for custom message
	if msg, ok := event.Data["message"].(string); ok {
		return msg
	}
	if desc, ok := event.Data["description"].(string); ok {
		return desc
	}

	// Build message from event data
	switch event.Type {
	case notification.EventMovieAdded, notification.EventMovieAvailable:
		if name, ok := event.Data["movie_title"].(string); ok {
			year := ""
			if y, ok := event.Data["year"].(int); ok {
				year = fmt.Sprintf(" (%d)", y)
			}
			return fmt.Sprintf("%s%s", name, year)
		}
	case notification.EventLibraryScanDone:
		if added, ok := event.Data["movies_added"].(int); ok {
			return fmt.Sprintf("%d new movies added", added)
		}
	case notification.EventLoginFailed:
		if username, ok := event.Data["username"].(string); ok {
			ip := ""
			if ipAddr, ok := event.Data["ip_address"].(string); ok {
				ip = fmt.Sprintf(" from %s", ipAddr)
			}
			return fmt.Sprintf("Failed login attempt for user: %s%s", username, ip)
		}
	}

	return fmt.Sprintf("Event: %s at %s", event.Type, event.Timestamp.Format("15:04:05"))
}

// getPriority determines the notification priority
func (a *GotifyAgent) getPriority(event *notification.Event) int {
	// High priority for security events
	switch event.Type {
	case notification.EventLoginFailed:
		return 8
	case notification.EventPasswordChanged, notification.EventMFAEnabled, notification.EventMFADisabled:
		return 7
	case notification.EventUserCreated, notification.EventUserDeleted:
		return 6
	}

	return a.config.DefaultPriority
}

// getClickURL returns the URL to open when notification is clicked
func (a *GotifyAgent) getClickURL(event *notification.Event) string {
	if url, ok := event.Data["url"].(string); ok {
		return url
	}
	// Default to nothing (or could be configured server URL)
	return ""
}

// NtfyConfig holds the configuration for an ntfy notification agent
type NtfyConfig struct {
	notification.AgentConfig

	// ServerURL is the ntfy server URL (default: https://ntfy.sh)
	ServerURL string `json:"server_url,omitempty"`

	// Topic is the ntfy topic to publish to
	Topic string `json:"topic"`

	// AccessToken is optional access token for authentication
	AccessToken string `json:"access_token,omitempty"`

	// Username and Password for basic auth (alternative to token)
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`

	// DefaultPriority is the default message priority (1-5, default: 3)
	DefaultPriority int `json:"default_priority,omitempty"`

	// Timeout is the request timeout (default: 30s)
	Timeout time.Duration `json:"timeout,omitempty"`
}

// NtfyAgent sends notifications to ntfy server
type NtfyAgent struct {
	config NtfyConfig
	client *http.Client
}

// NewNtfyAgent creates a new ntfy notification agent
func NewNtfyAgent(config NtfyConfig) (*NtfyAgent, error) {
	agent := &NtfyAgent{
		config: config,
	}

	// Set defaults
	if agent.config.ServerURL == "" {
		agent.config.ServerURL = "https://ntfy.sh"
	}
	if agent.config.DefaultPriority == 0 {
		agent.config.DefaultPriority = 3
	}
	if agent.config.Timeout == 0 {
		agent.config.Timeout = 30 * time.Second
	}

	// Create HTTP client
	agent.client = &http.Client{
		Timeout: agent.config.Timeout,
	}

	if err := agent.Validate(); err != nil {
		return nil, err
	}

	return agent, nil
}

// Type returns the agent type
func (a *NtfyAgent) Type() notification.AgentType {
	return notification.AgentNtfy
}

// Name returns the agent name
func (a *NtfyAgent) Name() string {
	if a.config.Name == "" {
		return "ntfy"
	}
	return a.config.Name
}

// IsEnabled returns whether the agent is enabled
func (a *NtfyAgent) IsEnabled() bool {
	return a.config.Enabled
}

// Validate checks if the configuration is valid
func (a *NtfyAgent) Validate() error {
	if a.config.Topic == "" {
		return fmt.Errorf("ntfy topic is required")
	}
	if a.config.DefaultPriority < 1 || a.config.DefaultPriority > 5 {
		return fmt.Errorf("priority must be between 1 and 5")
	}

	return nil
}

// Send sends a notification to ntfy
func (a *NtfyAgent) Send(ctx context.Context, event *notification.Event) error {
	if !a.config.ShouldSend(event.Type) {
		return nil // Skip this event type
	}

	// Build URL
	url := fmt.Sprintf("%s/%s", a.config.ServerURL, a.config.Topic)

	// Create request with message as body
	body := a.getMessage(event)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader([]byte(body)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Title", a.getTitle(event))
	req.Header.Set("Priority", fmt.Sprintf("%d", a.getPriority(event)))
	req.Header.Set("Tags", a.getTags(event))
	req.Header.Set("User-Agent", "Revenge/1.0")

	// Set click URL if available
	if clickURL := a.getClickURL(event); clickURL != "" {
		req.Header.Set("Click", clickURL)
	}

	// Set icon if available
	if icon := a.getIcon(event); icon != "" {
		req.Header.Set("Icon", icon)
	}

	// Authentication
	if a.config.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+a.config.AccessToken)
	} else if a.config.Username != "" && a.config.Password != "" {
		req.SetBasicAuth(a.config.Username, a.config.Password)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("ntfy returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// getTitle returns the notification title
func (a *NtfyAgent) getTitle(event *notification.Event) string {
	if title, ok := event.Data["title"].(string); ok {
		return title
	}

	switch event.Type {
	case notification.EventMovieAdded:
		if name, ok := event.Data["movie_title"].(string); ok {
			return fmt.Sprintf("New Movie: %s", name)
		}
		return "New Movie Added"
	case notification.EventMovieAvailable:
		return "Movie Now Available"
	case notification.EventRequestCreated:
		return "New Content Request"
	case notification.EventLibraryScanDone:
		return "Library Scan Complete"
	case notification.EventLoginFailed:
		return "Security Alert: Failed Login"
	default:
		return "Revenge Notification"
	}
}

// getMessage returns the notification message
func (a *NtfyAgent) getMessage(event *notification.Event) string {
	if msg, ok := event.Data["message"].(string); ok {
		return msg
	}
	if desc, ok := event.Data["description"].(string); ok {
		return desc
	}

	switch event.Type {
	case notification.EventMovieAdded, notification.EventMovieAvailable:
		if name, ok := event.Data["movie_title"].(string); ok {
			return name
		}
	case notification.EventLoginFailed:
		if username, ok := event.Data["username"].(string); ok {
			return fmt.Sprintf("Failed login attempt for: %s", username)
		}
	}

	return event.Type.String()
}

// getPriority returns the ntfy priority (1-5)
func (a *NtfyAgent) getPriority(event *notification.Event) int {
	switch event.Type {
	case notification.EventLoginFailed:
		return 5 // Urgent
	case notification.EventPasswordChanged, notification.EventMFAEnabled:
		return 4 // High
	case notification.EventMovieAvailable, notification.EventRequestApproved:
		return 3 // Default
	default:
		return a.config.DefaultPriority
	}
}

// getTags returns emoji tags for the notification
func (a *NtfyAgent) getTags(event *notification.Event) string {
	switch event.Type {
	case notification.EventMovieAdded, notification.EventMovieAvailable:
		return "movie_camera"
	case notification.EventRequestCreated:
		return "memo"
	case notification.EventRequestApproved:
		return "white_check_mark"
	case notification.EventRequestDenied:
		return "x"
	case notification.EventLibraryScanDone:
		return "file_folder"
	case notification.EventUserCreated:
		return "bust_in_silhouette"
	case notification.EventLoginFailed:
		return "warning"
	case notification.EventMFAEnabled:
		return "lock"
	case notification.EventPlaybackStarted:
		return "arrow_forward"
	case notification.EventSystemStartup:
		return "rocket"
	default:
		return "bell"
	}
}

// getClickURL returns the URL for clicking the notification
func (a *NtfyAgent) getClickURL(event *notification.Event) string {
	if url, ok := event.Data["url"].(string); ok {
		return url
	}
	return ""
}

// getIcon returns the icon URL for the notification
func (a *NtfyAgent) getIcon(event *notification.Event) string {
	if icon, ok := event.Data["icon_url"].(string); ok {
		return icon
	}
	// Could return a default Revenge icon URL
	return ""
}

// Ensure NtfyAgent implements Agent interface
var _ notification.Agent = (*NtfyAgent)(nil)
var _ notification.Agent = (*GotifyAgent)(nil)
