package email

import (
	"context"
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

func TestService_SendSendGrid_NotImplemented(t *testing.T) {
	cfg := config.EmailConfig{
		Enabled:     true,
		Provider:    "sendgrid",
		FromAddress: "test@example.com",
		SendGrid: config.SendGridConfig{
			APIKey: "test-key",
		},
	}
	svc := NewService(cfg, zap.NewNop())

	err := svc.SendVerificationEmail(context.Background(), "user@example.com", "testuser", "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SendGrid provider not yet implemented")
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
