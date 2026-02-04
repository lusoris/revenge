// Package email provides transactional email services for user verification,
// password reset, and other auth-related emails.
package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/zap"
)

// Service handles sending transactional emails.
type Service struct {
	cfg    config.EmailConfig
	logger *zap.Logger
}

// NewService creates a new email service.
func NewService(cfg config.EmailConfig, logger *zap.Logger) *Service {
	return &Service{
		cfg:    cfg,
		logger: logger.Named("email"),
	}
}

// SendVerificationEmail sends an email verification link to the user.
func (s *Service) SendVerificationEmail(ctx context.Context, toAddress, username, token string) error {
	if !s.cfg.Enabled {
		s.logger.Warn("Email disabled, skipping verification email",
			zap.String("to", toAddress),
			zap.String("username", username))
		return nil
	}

	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.cfg.BaseURL, token)

	subject := "Verify your email address - Revenge"
	body := buildVerificationEmail(username, verifyURL)

	return s.send(ctx, toAddress, subject, body)
}

// SendPasswordResetEmail sends a password reset link to the user.
func (s *Service) SendPasswordResetEmail(ctx context.Context, toAddress, username, token string) error {
	if !s.cfg.Enabled {
		s.logger.Warn("Email disabled, skipping password reset email",
			zap.String("to", toAddress),
			zap.String("username", username))
		return nil
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.BaseURL, token)

	subject := "Reset your password - Revenge"
	body := buildPasswordResetEmail(username, resetURL)

	return s.send(ctx, toAddress, subject, body)
}

// SendWelcomeEmail sends a welcome email to a newly verified user.
func (s *Service) SendWelcomeEmail(ctx context.Context, toAddress, username string) error {
	if !s.cfg.Enabled {
		s.logger.Warn("Email disabled, skipping welcome email",
			zap.String("to", toAddress),
			zap.String("username", username))
		return nil
	}

	subject := "Welcome to Revenge!"
	body := buildWelcomeEmail(username, s.cfg.BaseURL)

	return s.send(ctx, toAddress, subject, body)
}

// send dispatches an email using the configured provider.
func (s *Service) send(ctx context.Context, toAddress, subject, htmlBody string) error {
	switch s.cfg.Provider {
	case "smtp", "":
		return s.sendSMTP(ctx, toAddress, subject, htmlBody)
	case "sendgrid":
		return s.sendSendGrid(ctx, toAddress, subject, htmlBody)
	default:
		return fmt.Errorf("unknown email provider: %s", s.cfg.Provider)
	}
}

// sendSMTP sends an email via SMTP.
func (s *Service) sendSMTP(ctx context.Context, toAddress, subject, htmlBody string) error {
	cfg := s.cfg.SMTP

	if cfg.Host == "" {
		return fmt.Errorf("SMTP host not configured")
	}

	// Build email message
	msg := s.buildMessage(toAddress, subject, htmlBody)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// Create dialer with timeout
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	dialer := &net.Dialer{Timeout: timeout}

	var conn net.Conn
	var err error

	// Connect with or without TLS
	if cfg.UseTLS {
		// Direct TLS connection (port 465)
		tlsConfig := &tls.Config{
			ServerName: cfg.Host,
			// #nosec G402 -- InsecureSkipVerify is user-configurable for self-signed certs
			InsecureSkipVerify: cfg.SkipVerify,
		}
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer func() { _ = conn.Close() }()

	// Create SMTP client
	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer func() { _ = client.Close() }()

	// STARTTLS if configured (and not already using TLS)
	if cfg.UseStartTLS && !cfg.UseTLS {
		tlsConfig := &tls.Config{
			ServerName: cfg.Host,
			// #nosec G402 -- InsecureSkipVerify is user-configurable for self-signed certs
			InsecureSkipVerify: cfg.SkipVerify,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}

	// Authenticate if credentials provided
	if cfg.Username != "" && cfg.Password != "" {
		auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Send email
	if err := client.Mail(s.cfg.FromAddress); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	if err := client.Rcpt(toAddress); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close message: %w", err)
	}

	s.logger.Info("Email sent via SMTP",
		zap.String("to", toAddress),
		zap.String("subject", subject))

	return client.Quit()
}

// sendSendGrid sends an email via SendGrid API.
func (s *Service) sendSendGrid(ctx context.Context, toAddress, subject, htmlBody string) error {
	if s.cfg.SendGrid.APIKey == "" {
		return fmt.Errorf("SendGrid API key not configured")
	}

	// SendGrid API implementation
	// Using their v3 mail send endpoint
	// For now, we'll use a simple HTTP POST approach
	// In production, consider using the official sendgrid-go SDK

	s.logger.Info("Email sent via SendGrid",
		zap.String("to", toAddress),
		zap.String("subject", subject))

	// TODO: Implement SendGrid API call
	// For MVP, SMTP is sufficient. SendGrid can be added later.
	return fmt.Errorf("SendGrid provider not yet implemented, use SMTP")
}

// buildMessage creates the full email message with headers.
func (s *Service) buildMessage(toAddress, subject, htmlBody string) []byte {
	var buf bytes.Buffer

	// Headers
	from := s.cfg.FromAddress
	if s.cfg.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.FromAddress)
	}

	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", toAddress))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	buf.WriteString("X-Mailer: Revenge/1.0\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(htmlBody)

	return buf.Bytes()
}

// IsEnabled returns whether email sending is enabled.
func (s *Service) IsEnabled() bool {
	return s.cfg.Enabled
}

// buildVerificationEmail creates the HTML body for verification emails.
func buildVerificationEmail(username, verifyURL string) string {
	return buildEmailTemplate(
		"Verify Your Email Address",
		fmt.Sprintf("Hi %s,", username),
		"Thank you for registering with Revenge. Please click the button below to verify your email address.",
		verifyURL,
		"Verify Email",
		"This link will expire in 24 hours. If you didn't create an account, you can safely ignore this email.",
	)
}

// buildPasswordResetEmail creates the HTML body for password reset emails.
func buildPasswordResetEmail(username, resetURL string) string {
	return buildEmailTemplate(
		"Reset Your Password",
		fmt.Sprintf("Hi %s,", username),
		"We received a request to reset your password. Click the button below to create a new password.",
		resetURL,
		"Reset Password",
		"This link will expire in 1 hour. If you didn't request a password reset, you can safely ignore this email.",
	)
}

// buildWelcomeEmail creates the HTML body for welcome emails.
func buildWelcomeEmail(username, baseURL string) string {
	return buildEmailTemplate(
		"Welcome to Revenge!",
		fmt.Sprintf("Hi %s,", username),
		"Your email has been verified and your account is now active. You can now start using Revenge to manage your media library.",
		baseURL,
		"Go to Revenge",
		"Thank you for joining us!",
	)
}

// buildEmailTemplate creates a consistent HTML email template.
func buildEmailTemplate(title, greeting, message, actionURL, actionText, footer string) string {
	// Escape HTML in dynamic content
	title = escapeHTML(title)
	greeting = escapeHTML(greeting)
	message = escapeHTML(message)
	actionText = escapeHTML(actionText)
	footer = escapeHTML(footer)

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 600px;
            margin: 40px auto;
            background: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 30px 20px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
            font-weight: 600;
        }
        .content {
            padding: 30px 20px;
        }
        .greeting {
            font-size: 18px;
            font-weight: 500;
            margin-bottom: 15px;
        }
        .message {
            color: #555;
            margin-bottom: 25px;
        }
        .button {
            display: inline-block;
            padding: 14px 28px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white !important;
            text-decoration: none;
            border-radius: 6px;
            font-weight: 500;
            font-size: 16px;
        }
        .button:hover {
            opacity: 0.9;
        }
        .button-container {
            text-align: center;
            margin: 25px 0;
        }
        .footer {
            text-align: center;
            padding: 20px;
            color: #888;
            font-size: 13px;
            border-top: 1px solid #eee;
        }
        .link-fallback {
            margin-top: 20px;
            padding: 15px;
            background: #f9f9f9;
            border-radius: 4px;
            font-size: 12px;
            color: #666;
            word-break: break-all;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Revenge Media Server</h1>
        </div>
        <div class="content">
            <p class="greeting">%s</p>
            <p class="message">%s</p>
            <div class="button-container">
                <a href="%s" class="button">%s</a>
            </div>
            <div class="link-fallback">
                If the button doesn't work, copy and paste this link into your browser:<br>
                <a href="%s">%s</a>
            </div>
        </div>
        <div class="footer">
            <p>%s</p>
            <p>This email was sent by Revenge Media Server</p>
        </div>
    </div>
</body>
</html>`, title, greeting, message, actionURL, actionText, actionURL, actionURL, footer)
}

// escapeHTML escapes HTML special characters.
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
