# Notifications System

<!-- DESIGN: technical, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Multi-channel notification delivery

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

Unified notification system supporting multiple delivery channels:
- In-app notifications (WebSocket)
- Email
- Push notifications (FCM)
- Webhooks
- Discord/Telegram bots

---

## Notification Types

| Type | Channels | Priority |
|------|----------|----------|
| `library.new` | In-app, Push | Normal |
| `playback.recommendation` | In-app | Low |
| `account.security` | Email, In-app | High |
| `system.maintenance` | All | High |
| `admin.alert` | Email, Webhook | Critical |

---

## Configuration

```yaml
notifications:
  channels:
    inapp:
      enabled: true
    email:
      enabled: true
    push:
      enabled: true
      fcm:
        project_id: ${FCM_PROJECT_ID}
        credentials: ${FCM_CREDENTIALS}
    discord:
      enabled: false
      webhook_url: ${DISCORD_WEBHOOK}
```

---

## User Preferences

Users can configure per-channel, per-type preferences:

```json
{
  "userId": "uuid",
  "preferences": {
    "library.new": {
      "inapp": true,
      "email": false,
      "push": true
    },
    "account.security": {
      "inapp": true,
      "email": true,
      "push": true
    }
  }
}
```

---

## Architecture

```
┌─────────────────┐
│ Notification    │
│ Service         │
└────────┬────────┘
         │
    ┌────┴────┐
    │ Router  │
    └────┬────┘
         │
┌────────┼────────┬────────┬────────┐
│        │        │        │        │
▼        ▼        ▼        ▼        ▼
In-App  Email   Push   Discord  Webhook
```

---

## Implementation

```go
type NotificationService struct {
    channels []NotificationChannel
    router   *NotificationRouter
    prefs    NotificationPreferences
}

type NotificationChannel interface {
    Name() string
    Send(ctx context.Context, n *Notification) error
    Supports(notificationType string) bool
}
```

---

## Related

- [Notification Service](../services/NOTIFICATION.md)
- [Webhooks](WEBHOOKS.md)
- [Email System](EMAIL.md)
- [WebSockets](WEBSOCKETS.md)
