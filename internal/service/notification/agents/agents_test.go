package agents

import (
	"net/http"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/service/notification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWebhookAgent(t *testing.T) {
	config := WebhookConfig{
		AgentConfig: notification.AgentConfig{
			Enabled: true,
			Name:    "test-webhook",
		},
		URL: "https://example.com/webhook",
	}

	agent, err := NewWebhookAgent(config)
	require.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, notification.AgentWebhook, agent.Type())
	assert.Equal(t, "test-webhook", agent.Name())
	assert.True(t, agent.IsEnabled())

	// Check defaults
	assert.Equal(t, http.MethodPost, agent.config.Method)
	assert.Equal(t, "application/json", agent.config.ContentType)
	assert.Equal(t, 30*time.Second, agent.config.Timeout)
	assert.Equal(t, 3, agent.config.RetryCount)
}

func TestNewWebhookAgent_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config WebhookConfig
		errMsg string
	}{
		{
			name:   "missing URL",
			config: WebhookConfig{},
			errMsg: "URL is required",
		},
		{
			name: "URL too short",
			config: WebhookConfig{
				URL: "http://x",
			},
			errMsg: "too short",
		},
		{
			name: "invalid method",
			config: WebhookConfig{
				URL:    "https://example.com/webhook",
				Method: http.MethodGet,
			},
			errMsg: "invalid HTTP method",
		},
		{
			name: "invalid auth type",
			config: WebhookConfig{
				URL:      "https://example.com/webhook",
				AuthType: "unknown",
			},
			errMsg: "invalid auth type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWebhookAgent(tt.config)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestWebhookAgent_DefaultName(t *testing.T) {
	config := WebhookConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		URL:         "https://example.com/webhook",
	}

	agent, err := NewWebhookAgent(config)
	require.NoError(t, err)
	assert.Equal(t, "webhook", agent.Name())
}

func TestNewDiscordAgent(t *testing.T) {
	config := DiscordConfig{
		AgentConfig: notification.AgentConfig{
			Enabled: true,
			Name:    "test-discord",
		},
		WebhookURL: "https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz1234567890ABCDEF",
	}

	agent, err := NewDiscordAgent(config)
	require.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, notification.AgentDiscord, agent.Type())
	assert.Equal(t, "test-discord", agent.Name())
	assert.True(t, agent.IsEnabled())

	// Check defaults
	assert.Equal(t, "Revenge", agent.config.Username)
	assert.Equal(t, 30*time.Second, agent.config.Timeout)
}

func TestNewDiscordAgent_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config DiscordConfig
		errMsg string
	}{
		{
			name:   "missing webhook URL",
			config: DiscordConfig{},
			errMsg: "URL is required",
		},
		{
			name: "URL too short",
			config: DiscordConfig{
				WebhookURL: "https://discord.com/short",
			},
			errMsg: "too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDiscordAgent(tt.config)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestDiscordAgent_EventColors(t *testing.T) {
	config := DiscordConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		WebhookURL:  "https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz1234567890ABCDEF",
	}

	agent, _ := NewDiscordAgent(config)

	// Test color mapping
	assert.Equal(t, ColorError, agent.getEventColor(notification.EventLoginFailed))
	assert.Equal(t, ColorSuccess, agent.getEventColor(notification.EventMovieAdded))
	assert.Equal(t, ColorInfo, agent.getEventColor(notification.EventLibraryScanDone))
	assert.Equal(t, ColorError, agent.getEventColor(notification.EventRequestDenied))
}

func TestNewEmailAgent(t *testing.T) {
	config := EmailConfig{
		AgentConfig: notification.AgentConfig{
			Enabled: true,
			Name:    "test-email",
		},
		Host:        "smtp.example.com",
		FromAddress: "noreply@example.com",
		ToAddresses: []string{"admin@example.com"},
	}

	agent, err := NewEmailAgent(config)
	require.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, notification.AgentEmail, agent.Type())
	assert.Equal(t, "test-email", agent.Name())

	// Check defaults
	assert.Equal(t, 587, agent.config.Port)
	assert.Equal(t, "Revenge Media Server", agent.config.FromName)
	assert.Equal(t, 30*time.Second, agent.config.Timeout)
}

func TestNewEmailAgent_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config EmailConfig
		errMsg string
	}{
		{
			name:   "missing host",
			config: EmailConfig{},
			errMsg: "host is required",
		},
		{
			name: "missing from address",
			config: EmailConfig{
				Host: "smtp.example.com",
			},
			errMsg: "from address is required",
		},
		{
			name: "missing recipients",
			config: EmailConfig{
				Host:        "smtp.example.com",
				FromAddress: "noreply@example.com",
			},
			errMsg: "at least one recipient",
		},
		{
			name: "invalid from address",
			config: EmailConfig{
				Host:        "smtp.example.com",
				FromAddress: "invalid",
				ToAddresses: []string{"admin@example.com"},
			},
			errMsg: "invalid from address",
		},
		{
			name: "invalid recipient address",
			config: EmailConfig{
				Host:        "smtp.example.com",
				FromAddress: "noreply@example.com",
				ToAddresses: []string{"invalid"},
			},
			errMsg: "invalid recipient address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmailAgent(tt.config)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestNewGotifyAgent(t *testing.T) {
	config := GotifyConfig{
		AgentConfig: notification.AgentConfig{
			Enabled: true,
			Name:    "test-gotify",
		},
		ServerURL: "https://gotify.example.com",
		AppToken:  "secret-token",
	}

	agent, err := NewGotifyAgent(config)
	require.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, notification.AgentGotify, agent.Type())
	assert.Equal(t, "test-gotify", agent.Name())

	// Check defaults
	assert.Equal(t, 5, agent.config.DefaultPriority)
	assert.Equal(t, 30*time.Second, agent.config.Timeout)
}

func TestNewGotifyAgent_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config GotifyConfig
		errMsg string
	}{
		{
			name:   "missing server URL",
			config: GotifyConfig{},
			errMsg: "server URL is required",
		},
		{
			name: "missing app token",
			config: GotifyConfig{
				ServerURL: "https://gotify.example.com",
			},
			errMsg: "app token is required",
		},
		{
			name: "invalid priority",
			config: GotifyConfig{
				ServerURL:       "https://gotify.example.com",
				AppToken:        "token",
				DefaultPriority: 15,
			},
			errMsg: "priority must be between",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGotifyAgent(tt.config)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestNewNtfyAgent(t *testing.T) {
	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{
			Enabled: true,
			Name:    "test-ntfy",
		},
		Topic: "revenge-notifications",
	}

	agent, err := NewNtfyAgent(config)
	require.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, notification.AgentNtfy, agent.Type())
	assert.Equal(t, "test-ntfy", agent.Name())

	// Check defaults
	assert.Equal(t, "https://ntfy.sh", agent.config.ServerURL)
	assert.Equal(t, 3, agent.config.DefaultPriority)
	assert.Equal(t, 30*time.Second, agent.config.Timeout)
}

func TestNewNtfyAgent_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config NtfyConfig
		errMsg string
	}{
		{
			name:   "missing topic",
			config: NtfyConfig{},
			errMsg: "topic is required",
		},
		{
			name: "invalid priority",
			config: NtfyConfig{
				Topic:           "test",
				DefaultPriority: 10,
			},
			errMsg: "priority must be between",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewNtfyAgent(tt.config)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestNtfyAgent_Tags(t *testing.T) {
	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Topic:       "test",
	}

	agent, _ := NewNtfyAgent(config)

	// Test tag mapping
	assert.Equal(t, "movie_camera", agent.getTags(notification.NewEvent(notification.EventMovieAdded)))
	assert.Equal(t, "warning", agent.getTags(notification.NewEvent(notification.EventLoginFailed)))
	assert.Equal(t, "lock", agent.getTags(notification.NewEvent(notification.EventMFAEnabled)))
	assert.Equal(t, "rocket", agent.getTags(notification.NewEvent(notification.EventSystemStartup)))
}

func TestGotifyAgent_Priority(t *testing.T) {
	config := GotifyConfig{
		AgentConfig:     notification.AgentConfig{Enabled: true},
		ServerURL:       "https://gotify.example.com",
		AppToken:        "token",
		DefaultPriority: 5,
	}

	agent, _ := NewGotifyAgent(config)

	// High priority for security events
	assert.Equal(t, 8, agent.getPriority(notification.NewEvent(notification.EventLoginFailed)))
	assert.Equal(t, 7, agent.getPriority(notification.NewEvent(notification.EventPasswordChanged)))
	assert.Equal(t, 6, agent.getPriority(notification.NewEvent(notification.EventUserCreated)))

	// Default priority for other events
	assert.Equal(t, 5, agent.getPriority(notification.NewEvent(notification.EventMovieAdded)))
}
