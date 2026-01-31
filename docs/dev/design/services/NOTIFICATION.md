# Notification Service

<!-- SOURCES: fx, ogen, pgx, postgresql-arrays, postgresql-json, river -->

<!-- DESIGN: services, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Multi-channel notifications for users and admins


<!-- TOC-START -->

## Table of Contents

- [Developer Resources](#developer-resources)
- [Status](#status)
- [Overview](#overview)
- [Goals](#goals)
- [Non-Goals](#non-goals)
- [Technical Design](#technical-design)
  - [Notification Channels](#notification-channels)
  - [Notification Types](#notification-types)
  - [Service Interface](#service-interface)
  - [WebSocket Integration](#websocket-integration)
- [Database Schema](#database-schema)
- [River Jobs](#river-jobs)
- [Configuration](#configuration)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: In-App Channel](#phase-3-in-app-channel)
  - [Phase 4: Email Channel](#phase-4-email-channel)
  - [Phase 5: Webhook Channel](#phase-5-webhook-channel)
  - [Phase 6: Background Jobs](#phase-6-background-jobs)
  - [Phase 7: API Integration](#phase-7-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documents](#related-documents)

<!-- TOC-END -->

**Module**: `internal/service/notification`

## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Package | Purpose |
|---------|---------|
| River | Background job delivery |
| gobwas/ws | WebSocket connections |
| go-mail | SMTP email sending |
| go-fcm | Firebase Cloud Messaging |
| text/template | Notification templates |

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |---

## Overview

The Notification service handles all outbound notifications across multiple channels:
- In-app notifications (real-time via WebSocket)
- Email notifications
- Webhook notifications (for integrations)
- Push notifications (mobile apps)

## Goals

- Unified notification API for all services
- User-configurable notification preferences
- Reliable delivery with retries
- Template-based notification content

## Non-Goals

- Marketing/promotional notifications
- Third-party notification aggregators (Pushover, etc.) - future consideration
- SMS notifications

---

## Technical Design

### Notification Channels

| Channel | Transport | Use Cases |
|---------|-----------|-----------|
| `in_app` | WebSocket | Real-time alerts, activity |
| `email` | SMTP | Account, weekly digests |
| `webhook` | HTTP POST | Integrations, automation |
| `push` | FCM/APNs | Mobile app alerts |

### Notification Types

| Type | Channels | Example |
|------|----------|---------|
| `library.new_content` | in_app, push | "New movie added: Dune 2" |
| `playback.continue` | in_app | "Continue watching: Breaking Bad S3E5" |
| `system.update` | in_app, email | "Server update available" |
| `admin.alert` | in_app, email, webhook | "Library scan failed" |
| `user.mention` | in_app, push | "John shared a playlist with you" |

### Service Interface

```go
type NotificationService interface {
    // Send notifications
    Send(ctx context.Context, notification Notification) error
    SendBatch(ctx context.Context, notifications []Notification) error
    SendToUser(ctx context.Context, userID uuid.UUID, notificationType string, data map[string]interface{}) error
    SendToAdmins(ctx context.Context, notificationType string, data map[string]interface{}) error

    // User preferences
    GetPreferences(ctx context.Context, userID uuid.UUID) (*NotificationPreferences, error)
    UpdatePreferences(ctx context.Context, userID uuid.UUID, prefs NotificationPreferences) error

    // In-app notifications
    MarkAsRead(ctx context.Context, userID uuid.UUID, notificationID uuid.UUID) error
    GetUnread(ctx context.Context, userID uuid.UUID) ([]Notification, error)
}

type Notification struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Type      string
    Title     string
    Body      string
    Data      map[string]interface{}
    Channels  []string
    CreatedAt time.Time
    ReadAt    *time.Time
}

type NotificationPreferences struct {
    UserID   uuid.UUID
    Channels map[string]ChannelPrefs
}

type ChannelPrefs struct {
    Enabled bool
    Types   map[string]bool // Which notification types are enabled
}
```

### WebSocket Integration

```go
// Real-time notifications via WebSocket
type NotificationHub struct {
    connections map[uuid.UUID][]*websocket.Conn
    broadcast   chan Notification
}

func (h *NotificationHub) NotifyUser(userID uuid.UUID, notification Notification) {
    conns := h.connections[userID]
    for _, conn := range conns {
        conn.WriteJSON(notification)
    }
}
```

---

## Database Schema

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    body TEXT,
    data JSONB,
    channels TEXT[] NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ
);

CREATE INDEX idx_notifications_user_unread ON notifications (user_id, created_at)
    WHERE read_at IS NULL;

CREATE TABLE notification_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    preferences JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE notification_templates (
    type VARCHAR(50) PRIMARY KEY,
    title_template TEXT NOT NULL,
    body_template TEXT NOT NULL,
    default_channels TEXT[] NOT NULL
);
```

---

## River Jobs

```go
type SendNotificationArgs struct {
    NotificationID uuid.UUID `json:"notification_id"`
    Channel        string    `json:"channel"`
}

func (SendNotificationArgs) Kind() string { return "notification.send" }

type SendDigestArgs struct {
    UserID uuid.UUID `json:"user_id"`
    Period string    `json:"period"` // "daily", "weekly"
}

func (SendDigestArgs) Kind() string { return "notification.digest" }
```

---

## Configuration

```yaml
notification:
  enabled: true

  channels:
    in_app:
      enabled: true
      retention: 30d

    email:
      enabled: true
      smtp:
        host: "smtp.example.com"
        port: 587
        username: ""
        password: ""
        from: "Revenge <noreply@example.com>"

    webhook:
      enabled: true
      timeout: 10s
      max_retries: 3

    push:
      enabled: false
      fcm_key: ""
      apns_key: ""
```

---

## Implementation Checklist

### Phase 1: Core Infrastructure

- [ ] Create `internal/service/notification/` package structure
- [ ] Define notification entity and types
- [ ] Create repository interface
- [ ] Implement PostgreSQL repository
- [ ] Add fx module wiring

### Phase 2: Database

- [ ] Create migration for notifications table
- [ ] Create notification_preferences table
- [ ] Create notification_templates table
- [ ] Add indexes for unread queries

### Phase 3: In-App Channel

- [ ] Implement WebSocket hub
- [ ] Add connection management
- [ ] Implement real-time delivery
- [ ] Add notification read status

### Phase 4: Email Channel

- [ ] Implement SMTP sender
- [ ] Add email templates
- [ ] Implement digest notifications

### Phase 5: Webhook Channel

- [ ] Implement HTTP POST sender
- [ ] Add retry logic
- [ ] Add webhook configuration per user

### Phase 6: Background Jobs

- [ ] Create River job definitions
- [ ] Implement SendNotificationJob
- [ ] Implement SendDigestJob
- [ ] Add cleanup job for old notifications

### Phase 7: API Integration

- [ ] Define OpenAPI endpoints
- [ ] Generate ogen handlers
- [ ] Wire handlers to service
- [ ] Add preference management endpoints

---


## Related Documents

- [Activity Service](ACTIVITY.md) - Event logging for notifications
- [User Service](USER.md) - User preferences
- [Settings Service](SETTINGS.md) - Server-wide notification settings
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
