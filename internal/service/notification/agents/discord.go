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

// DiscordConfig holds the configuration for a Discord webhook notification agent
type DiscordConfig struct {
	notification.AgentConfig

	// WebhookURL is the Discord webhook URL
	WebhookURL string `json:"webhook_url"`

	// Username is the override username for the webhook (optional)
	Username string `json:"username,omitempty"`

	// AvatarURL is the override avatar URL for the webhook (optional)
	AvatarURL string `json:"avatar_url,omitempty"`

	// Timeout is the request timeout (default: 30s)
	Timeout time.Duration `json:"timeout,omitempty"`
}

// DiscordEmbed represents a Discord embed object
type DiscordEmbed struct {
	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	URL         string               `json:"url,omitempty"`
	Color       int                  `json:"color,omitempty"`
	Timestamp   string               `json:"timestamp,omitempty"`
	Footer      *DiscordEmbedFooter  `json:"footer,omitempty"`
	Thumbnail   *DiscordEmbedMedia   `json:"thumbnail,omitempty"`
	Image       *DiscordEmbedMedia   `json:"image,omitempty"`
	Author      *DiscordEmbedAuthor  `json:"author,omitempty"`
	Fields      []DiscordEmbedField  `json:"fields,omitempty"`
}

// DiscordEmbedFooter represents a Discord embed footer
type DiscordEmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// DiscordEmbedMedia represents a Discord embed image/thumbnail
type DiscordEmbedMedia struct {
	URL string `json:"url"`
}

// DiscordEmbedAuthor represents a Discord embed author
type DiscordEmbedAuthor struct {
	Name    string `json:"name"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// DiscordEmbedField represents a Discord embed field
type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// DiscordWebhookPayload is the payload sent to Discord webhooks
type DiscordWebhookPayload struct {
	Username  string         `json:"username,omitempty"`
	AvatarURL string         `json:"avatar_url,omitempty"`
	Content   string         `json:"content,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

// Discord colors for different event types
const (
	ColorSuccess = 0x2ECC71 // Green
	ColorInfo    = 0x3498DB // Blue
	ColorWarning = 0xF39C12 // Orange
	ColorError   = 0xE74C3C // Red
	ColorDefault = 0x9B59B6 // Purple (Revenge brand color)
)

// DiscordAgent sends notifications to Discord via webhooks
type DiscordAgent struct {
	config DiscordConfig
	client *http.Client
}

// NewDiscordAgent creates a new Discord notification agent
func NewDiscordAgent(config DiscordConfig) (*DiscordAgent, error) {
	agent := &DiscordAgent{
		config: config,
	}

	// Set defaults
	if agent.config.Timeout == 0 {
		agent.config.Timeout = 30 * time.Second
	}
	if agent.config.Username == "" {
		agent.config.Username = "Revenge"
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
func (a *DiscordAgent) Type() notification.AgentType {
	return notification.AgentDiscord
}

// Name returns the agent name
func (a *DiscordAgent) Name() string {
	if a.config.Name == "" {
		return "discord"
	}
	return a.config.Name
}

// IsEnabled returns whether the agent is enabled
func (a *DiscordAgent) IsEnabled() bool {
	return a.config.Enabled
}

// Validate checks if the configuration is valid
func (a *DiscordAgent) Validate() error {
	if a.config.WebhookURL == "" {
		return fmt.Errorf("discord webhook URL is required")
	}

	// Validate Discord webhook URL format
	if len(a.config.WebhookURL) < 50 {
		return fmt.Errorf("invalid discord webhook URL: too short")
	}

	return nil
}

// Send sends a notification to Discord
func (a *DiscordAgent) Send(ctx context.Context, event *notification.Event) error {
	if !a.config.ShouldSend(event.Type) {
		return nil // Skip this event type
	}

	// Build Discord payload
	payload := a.buildPayload(event)

	// Marshal to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.config.WebhookURL, bytes.NewReader(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Revenge/1.0")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Discord returns 204 No Content on success
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("discord returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// buildPayload constructs the Discord webhook payload
func (a *DiscordAgent) buildPayload(event *notification.Event) DiscordWebhookPayload {
	payload := DiscordWebhookPayload{
		Username:  a.config.Username,
		AvatarURL: a.config.AvatarURL,
	}

	embed := DiscordEmbed{
		Title:     a.getEventTitle(event),
		Color:     a.getEventColor(event.Type),
		Timestamp: event.Timestamp.Format(time.RFC3339),
		Footer: &DiscordEmbedFooter{
			Text: "Revenge Media Server",
		},
	}

	// Add description from event data
	if desc, ok := event.Data["description"].(string); ok {
		embed.Description = desc
	} else if msg, ok := event.Data["message"].(string); ok {
		embed.Description = msg
	}

	// Add thumbnail if available (e.g., movie poster)
	if posterURL, ok := event.Data["poster_url"].(string); ok && posterURL != "" {
		embed.Thumbnail = &DiscordEmbedMedia{URL: posterURL}
	}

	// Add common fields based on event data
	fields := a.buildFields(event)
	if len(fields) > 0 {
		embed.Fields = fields
	}

	payload.Embeds = []DiscordEmbed{embed}

	return payload
}

// getEventTitle returns a human-readable title for the event
func (a *DiscordAgent) getEventTitle(event *notification.Event) string {
	// Check for custom title in data
	if title, ok := event.Data["title"].(string); ok {
		return title
	}

	// Generate title based on event type
	switch event.Type {
	case notification.EventMovieAdded:
		if name, ok := event.Data["movie_title"].(string); ok {
			return fmt.Sprintf("🎬 New Movie Added: %s", name)
		}
		return "🎬 New Movie Added"
	case notification.EventMovieAvailable:
		if name, ok := event.Data["movie_title"].(string); ok {
			return fmt.Sprintf("✅ Movie Available: %s", name)
		}
		return "✅ Movie Now Available"
	case notification.EventRequestCreated:
		return "📝 New Request"
	case notification.EventRequestApproved:
		return "✅ Request Approved"
	case notification.EventRequestDenied:
		return "❌ Request Denied"
	case notification.EventLibraryScanStarted:
		return "🔄 Library Scan Started"
	case notification.EventLibraryScanDone:
		return "✅ Library Scan Complete"
	case notification.EventUserCreated:
		return "👤 New User Created"
	case notification.EventLoginFailed:
		return "⚠️ Failed Login Attempt"
	case notification.EventMFAEnabled:
		return "🔐 MFA Enabled"
	case notification.EventPasswordChanged:
		return "🔑 Password Changed"
	case notification.EventPlaybackStarted:
		return "▶️ Playback Started"
	case notification.EventSystemStartup:
		return "🚀 System Started"
	default:
		return fmt.Sprintf("📢 %s", event.Type)
	}
}

// getEventColor returns the Discord embed color for an event type
func (a *DiscordAgent) getEventColor(eventType notification.EventType) int {
	category := eventType.GetCategory()

	switch category {
	case notification.CategoryAuth:
		// Security events - orange/warning
		if eventType == notification.EventLoginFailed {
			return ColorError
		}
		return ColorWarning
	case notification.CategoryContent:
		return ColorSuccess
	case notification.CategoryRequests:
		if eventType == notification.EventRequestDenied {
			return ColorError
		}
		return ColorInfo
	case notification.CategoryLibrary:
		return ColorInfo
	case notification.CategoryPlayback:
		return ColorDefault
	default:
		return ColorDefault
	}
}

// buildFields creates embed fields from event data
func (a *DiscordAgent) buildFields(event *notification.Event) []DiscordEmbedField {
	var fields []DiscordEmbedField

	// Add movie-specific fields
	if year, ok := event.Data["year"].(int); ok {
		fields = append(fields, DiscordEmbedField{
			Name:   "Year",
			Value:  fmt.Sprintf("%d", year),
			Inline: true,
		})
	}

	if quality, ok := event.Data["quality"].(string); ok {
		fields = append(fields, DiscordEmbedField{
			Name:   "Quality",
			Value:  quality,
			Inline: true,
		})
	}

	// Add user info for relevant events
	if username, ok := event.Data["username"].(string); ok {
		fields = append(fields, DiscordEmbedField{
			Name:   "User",
			Value:  username,
			Inline: true,
		})
	}

	// Add IP for security events
	if ip, ok := event.Data["ip_address"].(string); ok {
		fields = append(fields, DiscordEmbedField{
			Name:   "IP Address",
			Value:  fmt.Sprintf("`%s`", ip),
			Inline: true,
		})
	}

	// Add library scan stats
	if added, ok := event.Data["movies_added"].(int); ok {
		fields = append(fields, DiscordEmbedField{
			Name:   "Movies Added",
			Value:  fmt.Sprintf("%d", added),
			Inline: true,
		})
	}

	if duration, ok := event.Data["scan_duration"].(string); ok {
		fields = append(fields, DiscordEmbedField{
			Name:   "Duration",
			Value:  duration,
			Inline: true,
		})
	}

	return fields
}

// Ensure DiscordAgent implements Agent interface
var _ notification.Agent = (*DiscordAgent)(nil)
