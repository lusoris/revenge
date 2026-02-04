package agents

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/notification"
)

// EmailConfig holds the configuration for an email notification agent
type EmailConfig struct {
	notification.AgentConfig

	// SMTP server configuration
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`

	// TLS configuration
	UseTLS      bool `json:"use_tls"`       // Use TLS from the start (port 465)
	UseStartTLS bool `json:"use_starttls"`  // Use STARTTLS upgrade (port 587)
	SkipVerify  bool `json:"skip_verify"`   // Skip TLS certificate verification

	// Email configuration
	FromAddress string `json:"from_address"`
	FromName    string `json:"from_name,omitempty"`

	// Default recipients (can be overridden per event)
	ToAddresses []string `json:"to_addresses"`

	// Timeout configuration
	Timeout time.Duration `json:"timeout,omitempty"`
}

// EmailAgent sends notifications via email (SMTP)
type EmailAgent struct {
	config EmailConfig
}

// NewEmailAgent creates a new email notification agent
func NewEmailAgent(config EmailConfig) (*EmailAgent, error) {
	agent := &EmailAgent{
		config: config,
	}

	// Set defaults
	if agent.config.Port == 0 {
		if agent.config.UseTLS {
			agent.config.Port = 465
		} else {
			agent.config.Port = 587
		}
	}
	if agent.config.Timeout == 0 {
		agent.config.Timeout = 30 * time.Second
	}
	if agent.config.FromName == "" {
		agent.config.FromName = "Revenge Media Server"
	}

	if err := agent.Validate(); err != nil {
		return nil, err
	}

	return agent, nil
}

// Type returns the agent type
func (a *EmailAgent) Type() notification.AgentType {
	return notification.AgentEmail
}

// Name returns the agent name
func (a *EmailAgent) Name() string {
	if a.config.Name == "" {
		return "email"
	}
	return a.config.Name
}

// IsEnabled returns whether the agent is enabled
func (a *EmailAgent) IsEnabled() bool {
	return a.config.Enabled
}

// Validate checks if the configuration is valid
func (a *EmailAgent) Validate() error {
	if a.config.Host == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if a.config.FromAddress == "" {
		return fmt.Errorf("from address is required")
	}
	if len(a.config.ToAddresses) == 0 {
		return fmt.Errorf("at least one recipient address is required")
	}

	// Validate email addresses
	if !strings.Contains(a.config.FromAddress, "@") {
		return fmt.Errorf("invalid from address: %s", a.config.FromAddress)
	}
	for _, addr := range a.config.ToAddresses {
		if !strings.Contains(addr, "@") {
			return fmt.Errorf("invalid recipient address: %s", addr)
		}
	}

	return nil
}

// Send sends an email notification
func (a *EmailAgent) Send(ctx context.Context, event *notification.Event) error {
	if !a.config.ShouldSend(event.Type) {
		return nil // Skip this event type
	}

	// Build email content
	subject := a.getSubject(event)
	body := a.buildBody(event)

	// Determine recipients
	recipients := a.config.ToAddresses
	if customRecipients, ok := event.Data["email_recipients"].([]string); ok && len(customRecipients) > 0 {
		recipients = customRecipients
	}

	// Build email message
	msg := a.buildMessage(subject, body, recipients)

	// Send email
	return a.sendEmail(ctx, recipients, msg)
}

// getSubject returns the email subject for an event
func (a *EmailAgent) getSubject(event *notification.Event) string {
	// Check for custom subject
	if subject, ok := event.Data["email_subject"].(string); ok {
		return subject
	}

	// Generate subject based on event type
	switch event.Type {
	case notification.EventMovieAdded:
		if name, ok := event.Data["movie_title"].(string); ok {
			return fmt.Sprintf("[Revenge] New Movie: %s", name)
		}
		return "[Revenge] New Movie Added"
	case notification.EventMovieAvailable:
		if name, ok := event.Data["movie_title"].(string); ok {
			return fmt.Sprintf("[Revenge] Now Available: %s", name)
		}
		return "[Revenge] Movie Now Available"
	case notification.EventRequestCreated:
		return "[Revenge] New Content Request"
	case notification.EventRequestApproved:
		return "[Revenge] Request Approved"
	case notification.EventRequestDenied:
		return "[Revenge] Request Denied"
	case notification.EventLibraryScanDone:
		return "[Revenge] Library Scan Complete"
	case notification.EventUserCreated:
		return "[Revenge] New User Account Created"
	case notification.EventLoginFailed:
		return "[Revenge] Security Alert: Failed Login"
	case notification.EventMFAEnabled:
		return "[Revenge] MFA Enabled on Your Account"
	case notification.EventPasswordChanged:
		return "[Revenge] Password Changed"
	case notification.EventSystemStartup:
		return "[Revenge] Server Started"
	default:
		return fmt.Sprintf("[Revenge] %s", event.Type)
	}
}

// buildBody creates the email body
func (a *EmailAgent) buildBody(event *notification.Event) string {
	var buf bytes.Buffer

	// HTML email template
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 20px; border: 1px solid #ddd; border-top: none; }
        .footer { text-align: center; padding: 20px; color: #888; font-size: 12px; }
        .field { margin: 10px 0; }
        .field-name { font-weight: bold; color: #555; }
        .field-value { color: #333; }
        .button { display: inline-block; padding: 10px 20px; background: #667eea; color: white; text-decoration: none; border-radius: 4px; margin-top: 15px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2 style="margin:0;">{{.Title}}</h2>
        </div>
        <div class="content">
            {{if .Description}}<p>{{.Description}}</p>{{end}}
            {{range .Fields}}
            <div class="field">
                <span class="field-name">{{.Name}}:</span>
                <span class="field-value">{{.Value}}</span>
            </div>
            {{end}}
            {{if .ActionURL}}<a href="{{.ActionURL}}" class="button">{{.ActionText}}</a>{{end}}
        </div>
        <div class="footer">
            <p>This email was sent by Revenge Media Server</p>
            <p>Event: {{.EventType}} | Time: {{.Timestamp}}</p>
        </div>
    </div>
</body>
</html>`

	// Build template data
	data := struct {
		Title       string
		Description string
		Fields      []struct{ Name, Value string }
		ActionURL   string
		ActionText  string
		EventType   string
		Timestamp   string
	}{
		Title:     a.getSubject(event),
		EventType: event.Type.String(),
		Timestamp: event.Timestamp.Format("2006-01-02 15:04:05 UTC"),
	}

	// Add description
	if desc, ok := event.Data["description"].(string); ok {
		data.Description = desc
	} else if msg, ok := event.Data["message"].(string); ok {
		data.Description = msg
	}

	// Add fields from event data
	for key, value := range event.Data {
		switch key {
		case "description", "message", "email_recipients", "email_subject":
			continue // Skip these
		default:
			data.Fields = append(data.Fields, struct{ Name, Value string }{
				Name:  formatFieldName(key),
				Value: fmt.Sprintf("%v", value),
			})
		}
	}

	// Parse and execute template
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return fmt.Sprintf("Error generating email: %v", err)
	}

	if err := t.Execute(&buf, data); err != nil {
		return fmt.Sprintf("Error generating email: %v", err)
	}

	return buf.String()
}

// formatFieldName converts snake_case to Title Case
func formatFieldName(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// buildMessage creates the full email message with headers
func (a *EmailAgent) buildMessage(subject, body string, recipients []string) []byte {
	var buf bytes.Buffer

	// Headers
	from := a.config.FromAddress
	if a.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", a.config.FromName, a.config.FromAddress)
	}

	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(recipients, ", ")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	buf.WriteString("X-Mailer: Revenge/1.0\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(body)

	return buf.Bytes()
}

// sendEmail sends the email via SMTP
func (a *EmailAgent) sendEmail(ctx context.Context, recipients []string, msg []byte) error {
	addr := fmt.Sprintf("%s:%d", a.config.Host, a.config.Port)

	// Create dialer with timeout
	dialer := &net.Dialer{Timeout: a.config.Timeout}

	var conn net.Conn
	var err error

	// Connect with or without TLS
	if a.config.UseTLS {
		// Direct TLS connection (port 465)
		tlsConfig := &tls.Config{
			ServerName: a.config.Host,
			// #nosec G402 -- InsecureSkipVerify is user-configurable for self-signed certs
			InsecureSkipVerify: a.config.SkipVerify,
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
	client, err := smtp.NewClient(conn, a.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer func() { _ = client.Close() }()

	// STARTTLS if configured (and not already using TLS)
	if a.config.UseStartTLS && !a.config.UseTLS {
		tlsConfig := &tls.Config{
			ServerName: a.config.Host,
			// #nosec G402 -- InsecureSkipVerify is user-configurable for self-signed certs
			InsecureSkipVerify: a.config.SkipVerify,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}

	// Authenticate if credentials provided
	if a.config.Username != "" && a.config.Password != "" {
		auth := smtp.PlainAuth("", a.config.Username, a.config.Password, a.config.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Send email
	if err := client.Mail(a.config.FromAddress); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	for _, rcpt := range recipients {
		if err := client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("RCPT TO failed for %s: %w", rcpt, err)
		}
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

	return client.Quit()
}

// Ensure EmailAgent implements Agent interface
var _ notification.Agent = (*EmailAgent)(nil)
