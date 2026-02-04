package agents

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
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

// ============================================================================
// Send Method Tests
// ============================================================================

func TestWebhookAgent_Send_Success(t *testing.T) {
	var receivedPayload WebhookPayload
	var receivedHeaders http.Header

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		err := json.NewDecoder(r.Body).Decode(&receivedPayload)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := WebhookConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		URL:         server.URL,
		Headers:     map[string]string{"X-Custom": "test-value"},
	}

	agent, err := NewWebhookAgent(config)
	require.NoError(t, err)

	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{"movie_title": "Test Movie"}
	userID := uuid.New()
	event.UserID = &userID

	err = agent.Send(context.Background(), event)
	require.NoError(t, err)

	// Verify payload
	assert.Equal(t, event.ID.String(), receivedPayload.EventID)
	assert.Equal(t, "movie.added", receivedPayload.EventType)
	assert.Equal(t, userID.String(), receivedPayload.UserID)
	assert.Equal(t, "revenge", receivedPayload.Source)

	// Verify headers
	assert.Equal(t, "application/json", receivedHeaders.Get("Content-Type"))
	assert.Equal(t, "test-value", receivedHeaders.Get("X-Custom"))
	assert.Equal(t, "Revenge/1.0", receivedHeaders.Get("User-Agent"))
}

func TestWebhookAgent_Send_BasicAuth(t *testing.T) {
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := WebhookConfig{
		AgentConfig:  notification.AgentConfig{Enabled: true},
		URL:          server.URL,
		AuthType:     "basic",
		AuthUsername: "testuser",
		AuthValue:    "testpass",
	}

	agent, err := NewWebhookAgent(config)
	require.NoError(t, err)

	err = agent.Send(context.Background(), notification.NewEvent(notification.EventMovieAdded))
	require.NoError(t, err)

	// Basic auth should be set
	assert.Contains(t, receivedAuth, "Basic")
}

func TestWebhookAgent_Send_BearerAuth(t *testing.T) {
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := WebhookConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		URL:         server.URL,
		AuthType:    "bearer",
		AuthValue:   "my-token-123",
	}

	agent, err := NewWebhookAgent(config)
	require.NoError(t, err)

	err = agent.Send(context.Background(), notification.NewEvent(notification.EventMovieAdded))
	require.NoError(t, err)

	assert.Equal(t, "Bearer my-token-123", receivedAuth)
}

func TestWebhookAgent_Send_HeaderAuth(t *testing.T) {
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("X-API-Key")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := WebhookConfig{
		AgentConfig:  notification.AgentConfig{Enabled: true},
		URL:          server.URL,
		AuthType:     "header",
		AuthUsername: "X-API-Key",
		AuthValue:    "secret-api-key",
	}

	agent, err := NewWebhookAgent(config)
	require.NoError(t, err)

	err = agent.Send(context.Background(), notification.NewEvent(notification.EventMovieAdded))
	require.NoError(t, err)

	assert.Equal(t, "secret-api-key", receivedAuth)
}

func TestWebhookAgent_Send_ServerError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))
	defer server.Close()

	config := WebhookConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		URL:         server.URL,
		RetryCount:  0, // No retries for faster test
	}

	agent, err := NewWebhookAgent(config)
	require.NoError(t, err)

	err = agent.Send(context.Background(), notification.NewEvent(notification.EventMovieAdded))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "webhook delivery failed")
	assert.Contains(t, err.Error(), "status 500")
}

func TestWebhookAgent_Send_EventFiltered(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := WebhookConfig{
		AgentConfig: notification.AgentConfig{
			Enabled:    true,
			EventTypes: []notification.EventType{notification.EventLoginFailed}, // Only login events
		},
		URL: server.URL,
	}

	agent, err := NewWebhookAgent(config)
	require.NoError(t, err)

	// Send movie event - should be filtered
	err = agent.Send(context.Background(), notification.NewEvent(notification.EventMovieAdded))
	require.NoError(t, err)
	assert.Equal(t, 0, requestCount)

	// Send login event - should be sent
	err = agent.Send(context.Background(), notification.NewEvent(notification.EventLoginFailed))
	require.NoError(t, err)
	assert.Equal(t, 1, requestCount)
}

func TestWebhookAgent_Send_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // Slow server
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := WebhookConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		URL:         server.URL,
		Timeout:     100 * time.Millisecond,
		RetryCount:  0,
	}

	agent, err := NewWebhookAgent(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = agent.Send(ctx, notification.NewEvent(notification.EventMovieAdded))
	require.Error(t, err)
}

func TestDiscordAgent_Send_Success(t *testing.T) {
	var receivedPayload DiscordWebhookPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		err := json.NewDecoder(r.Body).Decode(&receivedPayload)
		require.NoError(t, err)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	config := DiscordConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		WebhookURL:  server.URL + "/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz",
		Username:    "TestBot",
		AvatarURL:   "https://example.com/avatar.png",
	}

	agent, err := NewDiscordAgent(config)
	require.NoError(t, err)

	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{
		"movie_title": "The Matrix",
		"year":        1999,
		"quality":     "4K",
	}

	err = agent.Send(context.Background(), event)
	require.NoError(t, err)

	// Verify payload
	assert.Equal(t, "TestBot", receivedPayload.Username)
	assert.Equal(t, "https://example.com/avatar.png", receivedPayload.AvatarURL)
	require.Len(t, receivedPayload.Embeds, 1)

	embed := receivedPayload.Embeds[0]
	assert.Contains(t, embed.Title, "The Matrix")
	assert.Equal(t, ColorSuccess, embed.Color)
	assert.NotNil(t, embed.Footer)

	// Verify fields
	hasYear := false
	hasQuality := false
	for _, field := range embed.Fields {
		if field.Name == "Year" {
			hasYear = true
			assert.Equal(t, "1999", field.Value)
		}
		if field.Name == "Quality" {
			hasQuality = true
			assert.Equal(t, "4K", field.Value)
		}
	}
	assert.True(t, hasYear, "Year field should be present")
	assert.True(t, hasQuality, "Quality field should be present")
}

func TestDiscordAgent_Send_WithThumbnail(t *testing.T) {
	var receivedPayload DiscordWebhookPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&receivedPayload)
		require.NoError(t, err)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	config := DiscordConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		WebhookURL:  server.URL + "/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz",
	}

	agent, err := NewDiscordAgent(config)
	require.NoError(t, err)

	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{
		"poster_url": "https://example.com/poster.jpg",
	}

	err = agent.Send(context.Background(), event)
	require.NoError(t, err)

	require.Len(t, receivedPayload.Embeds, 1)
	require.NotNil(t, receivedPayload.Embeds[0].Thumbnail)
	assert.Equal(t, "https://example.com/poster.jpg", receivedPayload.Embeds[0].Thumbnail.URL)
}

func TestDiscordAgent_Send_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"message": "rate limited"}`))
	}))
	defer server.Close()

	config := DiscordConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		WebhookURL:  server.URL + "/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz",
	}

	agent, err := NewDiscordAgent(config)
	require.NoError(t, err)

	err = agent.Send(context.Background(), notification.NewEvent(notification.EventMovieAdded))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "discord returned status 429")
}

func TestGotifyAgent_Send_Success(t *testing.T) {
	var receivedMsg GotifyMessage
	var receivedToken string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/message", r.URL.Path)
		receivedToken = r.Header.Get("X-Gotify-Key")
		err := json.NewDecoder(r.Body).Decode(&receivedMsg)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": 1}`))
	}))
	defer server.Close()

	config := GotifyConfig{
		AgentConfig:     notification.AgentConfig{Enabled: true},
		ServerURL:       server.URL,
		AppToken:        "test-app-token",
		DefaultPriority: 5,
	}

	agent, err := NewGotifyAgent(config)
	require.NoError(t, err)

	event := notification.NewEvent(notification.EventLoginFailed)
	event.Data = map[string]any{
		"username":   "admin",
		"ip_address": "192.168.1.100",
	}

	err = agent.Send(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, "test-app-token", receivedToken)
	assert.Equal(t, "⚠️ Failed Login Attempt", receivedMsg.Title)
	assert.Contains(t, receivedMsg.Message, "admin")
	assert.Contains(t, receivedMsg.Message, "192.168.1.100")
	assert.Equal(t, 8, receivedMsg.Priority) // High priority for login failure
}

func TestGotifyAgent_Send_MovieEvent(t *testing.T) {
	var receivedMsg GotifyMessage

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&receivedMsg)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := GotifyConfig{
		AgentConfig:     notification.AgentConfig{Enabled: true},
		ServerURL:       server.URL,
		AppToken:        "token",
		DefaultPriority: 5,
	}

	agent, err := NewGotifyAgent(config)
	require.NoError(t, err)

	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{
		"movie_title": "Inception",
		"year":        2010,
	}

	err = agent.Send(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, "New Movie Added", receivedMsg.Title)
	assert.Equal(t, "Inception (2010)", receivedMsg.Message)
	assert.Equal(t, 5, receivedMsg.Priority) // Default priority
}

func TestGotifyAgent_Send_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("invalid token"))
	}))
	defer server.Close()

	config := GotifyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   server.URL,
		AppToken:    "wrong-token",
	}

	agent, err := NewGotifyAgent(config)
	require.NoError(t, err)

	err = agent.Send(context.Background(), notification.NewEvent(notification.EventMovieAdded))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gotify returned status 401")
}

func TestNtfyAgent_Send_Success(t *testing.T) {
	var receivedBody string
	var receivedHeaders http.Header

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   server.URL,
		Topic:       "test-topic",
	}

	agent, err := NewNtfyAgent(config)
	require.NoError(t, err)

	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{
		"movie_title": "Blade Runner",
	}

	err = agent.Send(context.Background(), event)
	require.NoError(t, err)

	// Verify headers
	assert.Contains(t, receivedHeaders.Get("Title"), "Blade Runner")
	assert.Equal(t, "3", receivedHeaders.Get("Priority"))
	assert.Equal(t, "movie_camera", receivedHeaders.Get("Tags"))

	// Verify body
	assert.Contains(t, receivedBody, "Blade Runner")
}

func TestNtfyAgent_Send_WithAuth(t *testing.T) {
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   server.URL,
		Topic:       "test-topic",
		AccessToken: "ntfy-secret-token",
	}

	agent, err := NewNtfyAgent(config)
	require.NoError(t, err)

	err = agent.Send(context.Background(), notification.NewEvent(notification.EventMovieAdded))
	require.NoError(t, err)

	assert.Equal(t, "Bearer ntfy-secret-token", receivedAuth)
}

func TestNtfyAgent_Send_BasicAuth(t *testing.T) {
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   server.URL,
		Topic:       "test-topic",
		Username:    "user",
		Password:    "pass",
	}

	agent, err := NewNtfyAgent(config)
	require.NoError(t, err)

	err = agent.Send(context.Background(), notification.NewEvent(notification.EventMovieAdded))
	require.NoError(t, err)

	assert.Contains(t, receivedAuth, "Basic")
}

func TestNtfyAgent_Send_SecurityEvent(t *testing.T) {
	var receivedHeaders http.Header

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := NtfyConfig{
		AgentConfig:     notification.AgentConfig{Enabled: true},
		ServerURL:       server.URL,
		Topic:           "alerts",
		DefaultPriority: 3,
	}

	agent, err := NewNtfyAgent(config)
	require.NoError(t, err)

	event := notification.NewEvent(notification.EventLoginFailed)
	event.Data = map[string]any{
		"username": "attacker",
	}

	err = agent.Send(context.Background(), event)
	require.NoError(t, err)

	// Should be urgent priority for security events
	assert.Equal(t, "5", receivedHeaders.Get("Priority"))
	assert.Equal(t, "warning", receivedHeaders.Get("Tags"))
	assert.Contains(t, receivedHeaders.Get("Title"), "Security Alert")
}

func TestNtfyAgent_Send_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("access denied"))
	}))
	defer server.Close()

	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   server.URL,
		Topic:       "test-topic",
	}

	agent, err := NewNtfyAgent(config)
	require.NoError(t, err)

	err = agent.Send(context.Background(), notification.NewEvent(notification.EventMovieAdded))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ntfy returned status 403")
}

func TestDiscordAgent_GetEventTitle(t *testing.T) {
	config := DiscordConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		WebhookURL:  "https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz1234567890ABCDEF",
	}
	agent, _ := NewDiscordAgent(config)

	tests := []struct {
		eventType notification.EventType
		data      map[string]any
		contains  string
	}{
		{notification.EventMovieAdded, map[string]any{"movie_title": "Test"}, "Test"},
		{notification.EventMovieAdded, nil, "New Movie Added"},
		{notification.EventMovieAvailable, map[string]any{"movie_title": "Available"}, "Available"},
		{notification.EventRequestCreated, nil, "New Request"},
		{notification.EventRequestApproved, nil, "Request Approved"},
		{notification.EventRequestDenied, nil, "Request Denied"},
		{notification.EventLibraryScanStarted, nil, "Scan Started"},
		{notification.EventLibraryScanDone, nil, "Scan Complete"},
		{notification.EventUserCreated, nil, "New User"},
		{notification.EventLoginFailed, nil, "Failed Login"},
		{notification.EventMFAEnabled, nil, "MFA Enabled"},
		{notification.EventPasswordChanged, nil, "Password Changed"},
		{notification.EventPlaybackStarted, nil, "Playback Started"},
		{notification.EventSystemStartup, nil, "System Started"},
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			event := notification.NewEvent(tt.eventType)
			if tt.data != nil {
				event.Data = tt.data
			}
			title := agent.getEventTitle(event)
			assert.Contains(t, title, tt.contains)
		})
	}
}

func TestGotifyAgent_GetMessage(t *testing.T) {
	config := GotifyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   "https://gotify.example.com",
		AppToken:    "token",
	}
	agent, _ := NewGotifyAgent(config)

	// Test custom message
	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{"message": "Custom message"}
	msg := agent.getMessage(event)
	assert.Equal(t, "Custom message", msg)

	// Test description fallback
	event.Data = map[string]any{"description": "A description"}
	msg = agent.getMessage(event)
	assert.Equal(t, "A description", msg)

	// Test library scan message
	event = notification.NewEvent(notification.EventLibraryScanDone)
	event.Data = map[string]any{"movies_added": 42}
	msg = agent.getMessage(event)
	assert.Equal(t, "42 new movies added", msg)
}

func TestNtfyAgent_GetMessage(t *testing.T) {
	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Topic:       "test",
	}
	agent, _ := NewNtfyAgent(config)

	// Test custom message
	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{"message": "Custom ntfy message"}
	msg := agent.getMessage(event)
	assert.Equal(t, "Custom ntfy message", msg)

	// Test login failed message
	event = notification.NewEvent(notification.EventLoginFailed)
	event.Data = map[string]any{"username": "hacker"}
	msg = agent.getMessage(event)
	assert.Contains(t, msg, "hacker")
}

func TestNtfyAgent_GetPriority(t *testing.T) {
	config := NtfyConfig{
		AgentConfig:     notification.AgentConfig{Enabled: true},
		Topic:           "test",
		DefaultPriority: 2,
	}
	agent, _ := NewNtfyAgent(config)

	tests := []struct {
		eventType notification.EventType
		priority  int
	}{
		{notification.EventLoginFailed, 5},         // Urgent
		{notification.EventPasswordChanged, 4},    // High
		{notification.EventMFAEnabled, 4},         // High
		{notification.EventMovieAvailable, 3},     // Default
		{notification.EventRequestApproved, 3},    // Default
		{notification.EventLibraryScanDone, 2},    // Uses config default
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			event := notification.NewEvent(tt.eventType)
			priority := agent.getPriority(event)
			assert.Equal(t, tt.priority, priority)
		})
	}
}

// ============================================================================
// Email Agent Helper Tests
// ============================================================================

func TestEmailAgent_DefaultName(t *testing.T) {
	config := EmailConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Host:        "smtp.example.com",
		FromAddress: "test@example.com",
		ToAddresses: []string{"admin@example.com"},
	}

	agent, err := NewEmailAgent(config)
	require.NoError(t, err)
	assert.Equal(t, "email", agent.Name())
	assert.True(t, agent.IsEnabled())
	assert.Equal(t, notification.AgentEmail, agent.Type())
}

func TestEmailAgent_NamedAgent(t *testing.T) {
	config := EmailConfig{
		AgentConfig: notification.AgentConfig{Enabled: false, Name: "admin-alerts"},
		Host:        "smtp.example.com",
		FromAddress: "test@example.com",
		ToAddresses: []string{"admin@example.com"},
	}

	agent, err := NewEmailAgent(config)
	require.NoError(t, err)
	assert.Equal(t, "admin-alerts", agent.Name())
	assert.False(t, agent.IsEnabled())
}

func TestEmailAgent_PortDefaults(t *testing.T) {
	// Test default port 587 (STARTTLS)
	config1 := EmailConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Host:        "smtp.example.com",
		FromAddress: "test@example.com",
		ToAddresses: []string{"admin@example.com"},
	}

	agent1, err := NewEmailAgent(config1)
	require.NoError(t, err)
	assert.Equal(t, 587, agent1.config.Port)

	// Test default port 465 (TLS)
	config2 := EmailConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Host:        "smtp.example.com",
		FromAddress: "test@example.com",
		ToAddresses: []string{"admin@example.com"},
		UseTLS:      true,
	}

	agent2, err := NewEmailAgent(config2)
	require.NoError(t, err)
	assert.Equal(t, 465, agent2.config.Port)
}

func TestEmailAgent_GetSubject(t *testing.T) {
	config := EmailConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Host:        "smtp.example.com",
		FromAddress: "test@example.com",
		ToAddresses: []string{"admin@example.com"},
	}

	agent, _ := NewEmailAgent(config)

	tests := []struct {
		eventType notification.EventType
		data      map[string]any
		contains  string
	}{
		{notification.EventMovieAdded, map[string]any{"movie_title": "Interstellar"}, "Interstellar"},
		{notification.EventMovieAdded, nil, "New Movie Added"},
		{notification.EventMovieAvailable, map[string]any{"movie_title": "Dune"}, "Dune"},
		{notification.EventMovieAvailable, nil, "Now Available"},
		{notification.EventRequestCreated, nil, "New Content Request"},
		{notification.EventRequestApproved, nil, "Request Approved"},
		{notification.EventRequestDenied, nil, "Request Denied"},
		{notification.EventLibraryScanDone, nil, "Library Scan Complete"},
		{notification.EventUserCreated, nil, "New User Account"},
		{notification.EventLoginFailed, nil, "Security Alert"},
		{notification.EventMFAEnabled, nil, "MFA Enabled"},
		{notification.EventPasswordChanged, nil, "Password Changed"},
		{notification.EventSystemStartup, nil, "Server Started"},
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			event := notification.NewEvent(tt.eventType)
			if tt.data != nil {
				event.Data = tt.data
			}
			subject := agent.getSubject(event)
			assert.Contains(t, subject, "[Revenge]")
			assert.Contains(t, subject, tt.contains)
		})
	}
}

func TestEmailAgent_GetSubject_CustomSubject(t *testing.T) {
	config := EmailConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Host:        "smtp.example.com",
		FromAddress: "test@example.com",
		ToAddresses: []string{"admin@example.com"},
	}

	agent, _ := NewEmailAgent(config)

	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{"email_subject": "Custom Subject Line"}

	subject := agent.getSubject(event)
	assert.Equal(t, "Custom Subject Line", subject)
}

func TestEmailAgent_BuildBody(t *testing.T) {
	config := EmailConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Host:        "smtp.example.com",
		FromAddress: "test@example.com",
		ToAddresses: []string{"admin@example.com"},
	}

	agent, _ := NewEmailAgent(config)

	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{
		"description":  "A sci-fi epic",
		"movie_title":  "Arrival",
		"year":         2016,
		"quality":      "1080p",
	}

	body := agent.buildBody(event)

	// Verify HTML structure
	assert.Contains(t, body, "<!DOCTYPE html>")
	assert.Contains(t, body, "<html>")
	assert.Contains(t, body, "A sci-fi epic")
	assert.Contains(t, body, "Revenge Media Server")
	assert.Contains(t, body, "movie.added")
}

func TestEmailAgent_BuildBody_MessageFallback(t *testing.T) {
	config := EmailConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Host:        "smtp.example.com",
		FromAddress: "test@example.com",
		ToAddresses: []string{"admin@example.com"},
	}

	agent, _ := NewEmailAgent(config)

	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{
		"message": "Fallback message content",
	}

	body := agent.buildBody(event)
	assert.Contains(t, body, "Fallback message content")
}

func TestEmailAgent_BuildMessage(t *testing.T) {
	config := EmailConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Host:        "smtp.example.com",
		FromAddress: "test@example.com",
		FromName:    "Test Sender",
		ToAddresses: []string{"admin@example.com"},
	}

	agent, _ := NewEmailAgent(config)

	msg := agent.buildMessage(
		"Test Subject",
		"<html><body>Test body</body></html>",
		[]string{"admin@example.com", "user@example.com"},
	)

	msgStr := string(msg)
	assert.Contains(t, msgStr, "From: Test Sender <test@example.com>")
	assert.Contains(t, msgStr, "To: admin@example.com, user@example.com")
	assert.Contains(t, msgStr, "Subject: Test Subject")
	assert.Contains(t, msgStr, "MIME-Version: 1.0")
	assert.Contains(t, msgStr, "Content-Type: text/html; charset=UTF-8")
	assert.Contains(t, msgStr, "X-Mailer: Revenge/1.0")
	assert.Contains(t, msgStr, "Test body")
}

func TestEmailAgent_BuildMessage_NoFromName(t *testing.T) {
	config := EmailConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Host:        "smtp.example.com",
		FromAddress: "test@example.com",
		FromName:    "",
		ToAddresses: []string{"admin@example.com"},
	}

	agent, _ := NewEmailAgent(config)
	// FromName defaults to "Revenge Media Server"
	assert.Equal(t, "Revenge Media Server", agent.config.FromName)

	msg := agent.buildMessage("Subject", "Body", []string{"admin@example.com"})
	msgStr := string(msg)
	assert.Contains(t, msgStr, "From: Revenge Media Server <test@example.com>")
}

func TestFormatFieldName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"movie_title", "Movie Title"},
		{"ip_address", "Ip Address"},
		{"username", "Username"},
		{"scan_duration_seconds", "Scan Duration Seconds"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatFieldName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================================
// Additional Agent Coverage Tests
// ============================================================================

func TestNtfyAgent_IsEnabled(t *testing.T) {
	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Topic:       "test",
	}
	agent, _ := NewNtfyAgent(config)
	assert.True(t, agent.IsEnabled())

	config2 := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: false},
		Topic:       "test",
	}
	agent2, _ := NewNtfyAgent(config2)
	assert.False(t, agent2.IsEnabled())
}

func TestNtfyAgent_DefaultName(t *testing.T) {
	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Topic:       "test",
	}
	agent, _ := NewNtfyAgent(config)
	assert.Equal(t, "ntfy", agent.Name())

	config2 := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true, Name: "my-ntfy"},
		Topic:       "test",
	}
	agent2, _ := NewNtfyAgent(config2)
	assert.Equal(t, "my-ntfy", agent2.Name())
}

func TestGotifyAgent_IsEnabled(t *testing.T) {
	config := GotifyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   "https://gotify.example.com",
		AppToken:    "token",
	}
	agent, _ := NewGotifyAgent(config)
	assert.True(t, agent.IsEnabled())
}

func TestGotifyAgent_DefaultName(t *testing.T) {
	config := GotifyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   "https://gotify.example.com",
		AppToken:    "token",
	}
	agent, _ := NewGotifyAgent(config)
	assert.Equal(t, "gotify", agent.Name())

	config2 := GotifyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true, Name: "my-gotify"},
		ServerURL:   "https://gotify.example.com",
		AppToken:    "token",
	}
	agent2, _ := NewGotifyAgent(config2)
	assert.Equal(t, "my-gotify", agent2.Name())
}

func TestDiscordAgent_DefaultName(t *testing.T) {
	config := DiscordConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		WebhookURL:  "https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz1234567890ABCDEF",
	}
	agent, _ := NewDiscordAgent(config)
	assert.Equal(t, "discord", agent.Name())
}

func TestDiscordAgent_IsEnabled(t *testing.T) {
	config := DiscordConfig{
		AgentConfig: notification.AgentConfig{Enabled: false},
		WebhookURL:  "https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz1234567890ABCDEF",
	}
	agent, _ := NewDiscordAgent(config)
	assert.False(t, agent.IsEnabled())
}

func TestNtfyAgent_GetTitle(t *testing.T) {
	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Topic:       "test",
	}
	agent, _ := NewNtfyAgent(config)

	tests := []struct {
		eventType notification.EventType
		data      map[string]any
		contains  string
	}{
		{notification.EventMovieAdded, map[string]any{"movie_title": "Matrix"}, "Matrix"},
		{notification.EventMovieAdded, nil, "New Movie Added"},
		{notification.EventMovieAvailable, nil, "Movie Now Available"},
		{notification.EventRequestCreated, nil, "New Content Request"},
		{notification.EventLibraryScanDone, nil, "Library Scan Complete"},
		{notification.EventLoginFailed, nil, "Security Alert"},
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			event := notification.NewEvent(tt.eventType)
			if tt.data != nil {
				event.Data = tt.data
			}
			title := agent.getTitle(event)
			assert.Contains(t, title, tt.contains)
		})
	}
}

func TestNtfyAgent_GetTitle_CustomTitle(t *testing.T) {
	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Topic:       "test",
	}
	agent, _ := NewNtfyAgent(config)

	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{"title": "Custom Title"}

	title := agent.getTitle(event)
	assert.Equal(t, "Custom Title", title)
}

func TestNtfyAgent_Tags_AllTypes(t *testing.T) {
	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Topic:       "test",
	}
	agent, _ := NewNtfyAgent(config)

	tests := []struct {
		eventType notification.EventType
		tag       string
	}{
		{notification.EventMovieAdded, "movie_camera"},
		{notification.EventMovieAvailable, "movie_camera"},
		{notification.EventRequestCreated, "memo"},
		{notification.EventRequestApproved, "white_check_mark"},
		{notification.EventRequestDenied, "x"},
		{notification.EventLibraryScanDone, "file_folder"},
		{notification.EventUserCreated, "bust_in_silhouette"},
		{notification.EventLoginFailed, "warning"},
		{notification.EventMFAEnabled, "lock"},
		{notification.EventPlaybackStarted, "arrow_forward"},
		{notification.EventSystemStartup, "rocket"},
		{notification.EventPasswordChanged, "bell"}, // Falls through to default
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			event := notification.NewEvent(tt.eventType)
			tag := agent.getTags(event)
			assert.Equal(t, tt.tag, tag)
		})
	}
}

func TestGotifyAgent_GetTitle(t *testing.T) {
	config := GotifyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   "https://gotify.example.com",
		AppToken:    "token",
	}
	agent, _ := NewGotifyAgent(config)

	tests := []struct {
		eventType notification.EventType
		expected  string
	}{
		{notification.EventMovieAdded, "New Movie Added"},
		{notification.EventMovieAvailable, "Movie Now Available"},
		{notification.EventRequestCreated, "New Content Request"},
		{notification.EventRequestApproved, "Request Approved"},
		{notification.EventRequestDenied, "Request Denied"},
		{notification.EventLibraryScanDone, "Library Scan Complete"},
		{notification.EventUserCreated, "New User Created"},
		{notification.EventLoginFailed, "⚠️ Failed Login Attempt"},
		{notification.EventMFAEnabled, "MFA Enabled"},
		{notification.EventPasswordChanged, "Password Changed"},
		{notification.EventPlaybackStarted, "Playback Started"},
		{notification.EventSystemStartup, "Revenge Started"},
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			event := notification.NewEvent(tt.eventType)
			title := agent.getTitle(event)
			assert.Equal(t, tt.expected, title)
		})
	}
}

func TestGotifyAgent_GetTitle_CustomTitle(t *testing.T) {
	config := GotifyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   "https://gotify.example.com",
		AppToken:    "token",
	}
	agent, _ := NewGotifyAgent(config)

	event := notification.NewEvent(notification.EventMovieAdded)
	event.Data = map[string]any{"title": "Custom Gotify Title"}

	title := agent.getTitle(event)
	assert.Equal(t, "Custom Gotify Title", title)
}

func TestGotifyAgent_GetClickURL(t *testing.T) {
	config := GotifyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		ServerURL:   "https://gotify.example.com",
		AppToken:    "token",
	}
	agent, _ := NewGotifyAgent(config)

	// No URL
	event := notification.NewEvent(notification.EventMovieAdded)
	url := agent.getClickURL(event)
	assert.Equal(t, "", url)

	// With URL
	event.Data = map[string]any{"url": "https://example.com/movie/123"}
	url = agent.getClickURL(event)
	assert.Equal(t, "https://example.com/movie/123", url)
}

func TestNtfyAgent_GetClickURL(t *testing.T) {
	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Topic:       "test",
	}
	agent, _ := NewNtfyAgent(config)

	// No URL
	event := notification.NewEvent(notification.EventMovieAdded)
	url := agent.getClickURL(event)
	assert.Equal(t, "", url)

	// With URL
	event.Data = map[string]any{"url": "https://example.com/click"}
	url = agent.getClickURL(event)
	assert.Equal(t, "https://example.com/click", url)
}

func TestNtfyAgent_GetIcon(t *testing.T) {
	config := NtfyConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		Topic:       "test",
	}
	agent, _ := NewNtfyAgent(config)

	// No icon
	event := notification.NewEvent(notification.EventMovieAdded)
	icon := agent.getIcon(event)
	assert.Equal(t, "", icon)

	// With icon
	event.Data = map[string]any{"icon_url": "https://example.com/icon.png"}
	icon = agent.getIcon(event)
	assert.Equal(t, "https://example.com/icon.png", icon)
}

func TestDiscordAgent_BuildFields(t *testing.T) {
	config := DiscordConfig{
		AgentConfig: notification.AgentConfig{Enabled: true},
		WebhookURL:  "https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz1234567890ABCDEF",
	}
	agent, _ := NewDiscordAgent(config)

	event := notification.NewEvent(notification.EventLoginFailed)
	event.Data = map[string]any{
		"year":           2024,
		"quality":        "4K",
		"username":       "testuser",
		"ip_address":     "192.168.1.50",
		"movies_added":   5,
		"scan_duration":  "2m30s",
	}

	fields := agent.buildFields(event)

	// Find fields by name
	fieldMap := make(map[string]string)
	for _, f := range fields {
		fieldMap[f.Name] = f.Value
	}

	assert.Equal(t, "2024", fieldMap["Year"])
	assert.Equal(t, "4K", fieldMap["Quality"])
	assert.Equal(t, "testuser", fieldMap["User"])
	assert.Contains(t, fieldMap["IP Address"], "192.168.1.50")
	assert.Equal(t, "5", fieldMap["Movies Added"])
	assert.Equal(t, "2m30s", fieldMap["Duration"])
}
