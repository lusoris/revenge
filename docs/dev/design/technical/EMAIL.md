# Email System

<!-- SOURCES: go-mail-docs -->

<!-- DESIGN: technical, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Transactional email for notifications and account management

**Source of Truth**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md)

---

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🟡 | Scaffold |
| Sources | 🔴 |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
---

## Overview

Revenge uses email for:
- Password reset
- Account verification
- Admin notifications
- Scheduled reports
- Invitation links

---

## Email Types

| Type | Trigger | Template |
|------|---------|----------|
| `password-reset` | User requests reset | Reset link, expiry |
| `verify-email` | New registration | Verification link |
| `invite` | Admin creates invite | Invite link, expiry |
| `weekly-report` | Scheduled | Activity summary |
| `admin-alert` | System event | Alert details |

---

## Configuration

```yaml
email:
  enabled: true
  from: "Revenge <noreply@example.com>"
  smtp:
    host: smtp.example.com
    port: 587
    username: ${SMTP_USER}
    password: ${SMTP_PASS}
    tls: starttls
```

---

## Implementation

Using `github.com/wneessen/go-mail`:

```go
func (s *EmailService) SendPasswordReset(to, token string) error {
    m := mail.NewMsg()
    m.From(s.config.From)
    m.To(to)
    m.Subject("Reset Your Password")
    m.SetBodyString(mail.TypeTextHTML, s.renderTemplate("password-reset", map[string]any{
        "Token":  token,
        "Expiry": "1 hour",
    }))

    return s.client.DialAndSend(m)
}
```

---

## Templates

HTML email templates with inline CSS for compatibility:
- `templates/email/password-reset.html`
- `templates/email/verify-email.html`
- `templates/email/invite.html`
- `templates/email/weekly-report.html`

---

## Related

- [Notification Service](../services/NOTIFICATION.md)
- [User Service](../services/USER.md)
- [Auth Service](../services/AUTH.md)
- [go-mail Source](../../sources/tooling/go-mail.md)
