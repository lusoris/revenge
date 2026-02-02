## Table of Contents

- [Notification Service](#notification-service)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Service Structure](#service-structure)
    - [Dependencies](#dependencies)
    - [Provides](#provides)
    - [Component Diagram](#component-diagram)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
- [In-app notifications](#in-app-notifications)
- [Preferences](#preferences)
- [Push devices](#push-devices)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Notification Service


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service


> > Multi-channel notifications for users and admins

**Package**: `internal/service/notification`
**fx Module**: `notification.Module`

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

```mermaid
flowchart TD
    node1["Event<br/>Publisher"]
    node2[[Notification<br/>Service]]
    node3([Channels<br/>(Email,<br/>Push,])
    node4[(PostgreSQL<br/>(pgx))]
    node5([External<br/>Services<br/>(SMTP, FCM)])
    node1 --> node2
    node2 --> node3
    node4 --> node5
    node3 --> node4
```

### Service Structure

```
internal/service/notification/
├── module.go              # fx module definition
├── service.go             # Service implementation
├── repository.go          # Data access (if needed)
├── handler.go             # HTTP handlers (if exposed)
├── middleware.go          # Middleware (if needed)
├── types.go               # Domain types
└── service_test.go        # Tests
```

### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/riverqueue/river` - Background notification delivery
- `net/smtp` - Email delivery
- `firebase.google.com/go/v4/messaging` - Firebase Cloud Messaging (FCM)
- `text/template` - Template rendering
- `net/http` - Webhook delivery
- `go.uber.org/fx`

**External Services**:
- SMTP server for email
- Firebase Cloud Messaging (FCM) for push notifications


### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->
## Implementation

### Key Interfaces

```go
type NotificationService interface {
  // Send notifications
  Send(ctx context.Context, notification Notification) error
  SendToUser(ctx context.Context, userID uuid.UUID, notification Notification) error
  SendBulk(ctx context.Context, userIDs []uuid.UUID, notification Notification) error

  // In-app notifications
  GetNotifications(ctx context.Context, userID uuid.UUID, filters NotificationFilters) ([]InAppNotification, error)
  MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
  MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
  DeleteNotification(ctx context.Context, notificationID uuid.UUID) error

  // Preferences
  GetPreferences(ctx context.Context, userID uuid.UUID) (map[string]NotificationPreference, error)
  UpdatePreference(ctx context.Context, userID uuid.UUID, eventType string, pref NotificationPreference) error

  // Push devices
  RegisterPushDevice(ctx context.Context, userID uuid.UUID, fcmToken, deviceType string) error
  UnregisterPushDevice(ctx context.Context, fcmToken string) error
}

type Notification struct {
  EventType string                 `json:"event_type"`
  Title     string                 `json:"title"`
  Body      string                 `json:"body"`
  ActionURL *string                `json:"action_url,omitempty"`
  Data      map[string]interface{} `json:"data,omitempty"`
  Priority  string                 `json:"priority"`
}

type NotificationChannel interface {
  Name() string
  Send(ctx context.Context, recipient User, notification Notification) error
}
```


### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/riverqueue/river` - Background notification delivery
- `net/smtp` - Email delivery
- `firebase.google.com/go/v4/messaging` - Firebase Cloud Messaging (FCM)
- `text/template` - Template rendering
- `net/http` - Webhook delivery
- `go.uber.org/fx`

**External Services**:
- SMTP server for email
- Firebase Cloud Messaging (FCM) for push notifications

## Configuration

### Environment Variables

```bash
NOTIFICATION_SMTP_HOST=smtp.gmail.com
NOTIFICATION_SMTP_PORT=587
NOTIFICATION_SMTP_USER=noreply@example.com
NOTIFICATION_SMTP_PASSWORD=your_password
NOTIFICATION_FCM_CREDENTIALS_FILE=/config/fcm-service-account.json
NOTIFICATION_DELIVERY_WORKERS=5
```


### Config Keys
```yaml
notification:
  email:
    smtp_host: smtp.gmail.com
    smtp_port: 587
    smtp_user: noreply@example.com
    smtp_password: your_password
    from_address: Revenge <noreply@example.com>
  push:
    fcm_credentials_file: /config/fcm-service-account.json
  webhook:
    timeout: 10s
  delivery:
    workers: 5
    retry_attempts: 3
    retry_delay: 5m
```

## API Endpoints
```
# In-app notifications
GET    /api/v1/notifications              # List notifications
GET    /api/v1/notifications/unread/count # Get unread count
PATCH  /api/v1/notifications/:id/read     # Mark as read
POST   /api/v1/notifications/read-all     # Mark all as read
DELETE /api/v1/notifications/:id          # Delete notification

# Preferences
GET    /api/v1/notifications/preferences  # Get preferences
PUT    /api/v1/notifications/preferences/:event # Update preference

# Push devices
POST   /api/v1/notifications/devices      # Register push device
DELETE /api/v1/notifications/devices/:token # Unregister device
```

**Example Notification Response**:
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "New Movie Added",
  "body": "Inception (2010) has been added to your library",
  "icon_url": "https://image.tmdb.org/t/p/w200/...",
  "action_url": "/movies/27205",
  "is_read": false,
  "created_at": "2026-02-01T10:30:00Z"
}
```

**Example Email Template**:
```
Subject: {{ .Title }}

Hi {{ .UserDisplayName }},

{{ .Body }}

{{ if .ActionURL }}
View: {{ .ActionURL }}
{{ end }}

--
Revenge Media Server
```

## Related Documentation
### Design Documents
- [services](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx
- [ogen OpenAPI Generator](../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river

