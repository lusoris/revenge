package email

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/revenge/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewService(t *testing.T) {
	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "smtp",
		FromAddress: "test@example.com",
		FromName:    "Test",
		BaseURL:     "http://localhost:8080",
	}

	logger := zap.NewNop()
	svc := NewService(cfg, logger)

	require.NotNil(t, svc)
	assert.True(t, svc.IsEnabled())
}

func TestService_IsEnabled(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
		want    bool
	}{
		{"enabled", true, true},
		{"disabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.EmailConfig{Enabled: tt.enabled}
			svc := NewService(cfg, zap.NewNop())
			assert.Equal(t, tt.want, svc.IsEnabled())
		})
	}
}

func TestService_SendVerificationEmail_Disabled(t *testing.T) {
	cfg := config.EmailConfig{Enabled: false}
	svc := NewService(cfg, zap.NewNop())

	// Should return nil when disabled (no-op)
	err := svc.SendVerificationEmail(context.Background(), "user@example.com", "testuser", "token123")
	assert.NoError(t, err)
}

func TestService_SendPasswordResetEmail_Disabled(t *testing.T) {
	cfg := config.EmailConfig{Enabled: false}
	svc := NewService(cfg, zap.NewNop())

	// Should return nil when disabled (no-op)
	err := svc.SendPasswordResetEmail(context.Background(), "user@example.com", "testuser", "token123")
	assert.NoError(t, err)
}

func TestService_SendWelcomeEmail_Disabled(t *testing.T) {
	cfg := config.EmailConfig{Enabled: false}
	svc := NewService(cfg, zap.NewNop())

	// Should return nil when disabled (no-op)
	err := svc.SendWelcomeEmail(context.Background(), "user@example.com", "testuser")
	assert.NoError(t, err)
}

func TestService_SendSMTP_NoHost(t *testing.T) {
	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "smtp",
		FromAddress: "test@example.com",
		SMTP: config.SMTPConfig{
			Host: "", // No host configured
		},
	}
	svc := NewService(cfg, zap.NewNop())

	err := svc.SendVerificationEmail(context.Background(), "user@example.com", "testuser", "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SMTP host not configured")
}

func TestService_SendSendGrid_Success(t *testing.T) {
	var receivedReq sendGridRequest
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		err := json.NewDecoder(r.Body).Decode(&receivedReq)
		require.NoError(t, err)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	// Override the endpoint for testing
	origEndpoint := sendGridEndpoint
	sendGridEndpoint = server.URL
	defer func() { sendGridEndpoint = origEndpoint }()

	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "sendgrid",
		FromAddress: "sender@example.com",
		FromName:    "Test Sender",
		BaseURL:     "http://localhost:8080",
		SendGrid: config.SendGridConfig{
			APIKey: "SG.test-api-key",
		},
	}
	svc := NewService(cfg, zap.NewNop())

	err := svc.SendVerificationEmail(context.Background(), "user@example.com", "testuser", "token123")
	require.NoError(t, err)

	assert.Equal(t, "Bearer SG.test-api-key", receivedAuth)
	require.Len(t, receivedReq.Personalizations, 1)
	require.Len(t, receivedReq.Personalizations[0].To, 1)
	assert.Equal(t, "user@example.com", receivedReq.Personalizations[0].To[0].Email)
	assert.Equal(t, "sender@example.com", receivedReq.From.Email)
	assert.Equal(t, "Test Sender", receivedReq.From.Name)
	assert.Equal(t, "Verify your email address - Revenge", receivedReq.Subject)
	require.Len(t, receivedReq.Content, 1)
	assert.Equal(t, "text/html", receivedReq.Content[0].Type)
	assert.Contains(t, receivedReq.Content[0].Value, "Verify Email")
}

func TestService_SendSendGrid_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"errors":[{"message":"The provided authorization grant is invalid"}]}`))
	}))
	defer server.Close()

	origEndpoint := sendGridEndpoint
	sendGridEndpoint = server.URL
	defer func() { sendGridEndpoint = origEndpoint }()

	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "sendgrid",
		FromAddress: "test@example.com",
		BaseURL:     "http://localhost:8080",
		SendGrid: config.SendGridConfig{
			APIKey: "bad-key",
		},
	}
	svc := NewService(cfg, zap.NewNop())

	err := svc.SendVerificationEmail(context.Background(), "user@example.com", "testuser", "token123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SendGrid API error (status 401)")
	assert.Contains(t, err.Error(), "invalid")
}

func TestBuildVerificationEmail(t *testing.T) {
	body := buildVerificationEmail("testuser", "http://localhost:8080/verify?token=abc123")

	assert.Contains(t, body, "Verify Your Email Address")
	assert.Contains(t, body, "testuser")
	assert.Contains(t, body, "http://localhost:8080/verify?token=abc123")
	assert.Contains(t, body, "Verify Email")
}

func TestBuildPasswordResetEmail(t *testing.T) {
	body := buildPasswordResetEmail("testuser", "http://localhost:8080/reset?token=abc123")

	assert.Contains(t, body, "Reset Your Password")
	assert.Contains(t, body, "testuser")
	assert.Contains(t, body, "http://localhost:8080/reset?token=abc123")
	assert.Contains(t, body, "Reset Password")
}

func TestBuildWelcomeEmail(t *testing.T) {
	body := buildWelcomeEmail("testuser", "http://localhost:8080")

	assert.Contains(t, body, "Welcome to Revenge!")
	assert.Contains(t, body, "testuser")
	assert.Contains(t, body, "http://localhost:8080")
}

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"Hello & World", "Hello &amp; World"},
		{"\"quoted\"", "&quot;quoted&quot;"},
		{"normal text", "normal text"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := escapeHTML(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_SendUnknownProvider(t *testing.T) {
	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "unknown_provider",
		FromAddress: "test@example.com",
	}
	svc := NewService(cfg, zap.NewNop())

	err := svc.SendVerificationEmail(context.Background(), "user@example.com", "testuser", "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown email provider")
}

func TestService_SendSendGrid_NoAPIKey(t *testing.T) {
	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "sendgrid",
		FromAddress: "test@example.com",
		SendGrid: config.SendGridConfig{
			APIKey: "", // No API key
		},
	}
	svc := NewService(cfg, zap.NewNop())

	err := svc.SendVerificationEmail(context.Background(), "user@example.com", "testuser", "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SendGrid API key not configured")
}

func TestBuildMessage(t *testing.T) {
	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "smtp",
		FromAddress: "test@example.com",
		FromName:    "Test Sender",
	}
	svc := NewService(cfg, zap.NewNop())

	msg := svc.buildMessage("recipient@example.com", "Test Subject", "<p>Test Body</p>")

	// Verify headers are present
	msgStr := string(msg)
	assert.Contains(t, msgStr, "From: Test Sender <test@example.com>")
	assert.Contains(t, msgStr, "To: recipient@example.com")
	assert.Contains(t, msgStr, "Subject: Test Subject")
	assert.Contains(t, msgStr, "MIME-Version: 1.0")
	assert.Contains(t, msgStr, "Content-Type: text/html; charset=UTF-8")
	assert.Contains(t, msgStr, "X-Mailer: Revenge/1.0")
	assert.Contains(t, msgStr, "<p>Test Body</p>")
}

func TestBuildMessage_NoFromName(t *testing.T) {
	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "smtp",
		FromAddress: "test@example.com",
		FromName:    "", // No from name
	}
	svc := NewService(cfg, zap.NewNop())

	msg := svc.buildMessage("recipient@example.com", "Test Subject", "<p>Test Body</p>")

	// Verify From header uses just the address
	msgStr := string(msg)
	assert.Contains(t, msgStr, "From: test@example.com\r\n")
	assert.NotContains(t, msgStr, "From:  <") // No empty name with angle brackets
}

func TestBuildEmailTemplate_EscapesContent(t *testing.T) {
	// Test that HTML content is properly escaped
	body := buildEmailTemplate(
		"Title with <script>",
		"<Greeting>",
		"Message with & special chars",
		"http://example.com",
		"Click <here>",
		"Footer with \"quotes\"",
	)

	// The content should be escaped
	assert.Contains(t, body, "&lt;script&gt;")
	assert.Contains(t, body, "&lt;Greeting&gt;")
	assert.Contains(t, body, "&amp; special chars")
	assert.Contains(t, body, "&lt;here&gt;")
	assert.Contains(t, body, "&quot;quotes&quot;")

	// URL should not be escaped (it's used in href)
	assert.Contains(t, body, `href="http://example.com"`)
}

func TestService_SendPasswordResetEmail_Enabled(t *testing.T) {
	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "smtp",
		FromAddress: "test@example.com",
		BaseURL:     "http://localhost:8080",
		SMTP: config.SMTPConfig{
			Host: "", // Will fail due to no host
		},
	}
	svc := NewService(cfg, zap.NewNop())

	// Should fail when trying to send (no SMTP host)
	err := svc.SendPasswordResetEmail(context.Background(), "user@example.com", "testuser", "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SMTP host not configured")
}

func TestService_SendWelcomeEmail_Enabled(t *testing.T) {
	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "smtp",
		FromAddress: "test@example.com",
		BaseURL:     "http://localhost:8080",
		SMTP: config.SMTPConfig{
			Host: "", // Will fail due to no host
		},
	}
	svc := NewService(cfg, zap.NewNop())

	// Should fail when trying to send (no SMTP host)
	err := svc.SendWelcomeEmail(context.Background(), "user@example.com", "testuser")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SMTP host not configured")
}

func TestService_EmptyProvider(t *testing.T) {
	// Empty provider should default to SMTP
	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "", // Empty
		FromAddress: "test@example.com",
		SMTP: config.SMTPConfig{
			Host: "", // Will fail due to no host
		},
	}
	svc := NewService(cfg, zap.NewNop())

	err := svc.SendVerificationEmail(context.Background(), "user@example.com", "testuser", "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SMTP host not configured") // Uses SMTP as default
}
