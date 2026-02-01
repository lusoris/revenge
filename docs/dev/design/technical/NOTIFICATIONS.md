## Table of Contents

- [Notifications System](#notifications-system)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Components](#components)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [GET /api/v1/user/notifications/preferences](#get-apiv1usernotificationspreferences)
    - [PUT /api/v1/user/notifications/preferences](#put-apiv1usernotificationspreferences)
    - [POST /api/v1/user/push/register](#post-apiv1userpushregister)
    - [DELETE /api/v1/user/push/unregister](#delete-apiv1userpushunregister)
    - [GET /api/v1/user/notifications/history](#get-apiv1usernotificationshistory)
    - [POST /api/v1/admin/webhooks](#post-apiv1adminwebhooks)
    - [GET /api/v1/admin/webhooks](#get-apiv1adminwebhooks)
    - [DELETE /api/v1/admin/webhooks/:id](#delete-apiv1adminwebhooksid)
    - [POST /api/v1/admin/notifications/test](#post-apiv1adminnotificationstest)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Notifications System


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > Multi-channel notification system: Email, Push (FCM), Webhooks

Unified notification delivery:
- **Email**: SMTP via go-mail, templates, bounce handling
- **Push**: Firebase Cloud Messaging (FCM) via go-fcm
- **Webhooks**: HTTP callbacks to external services
- **Queue**: River async job queue with retries
- **Preferences**: Per-user, per-channel filtering
- **Features**: Rate limiting, deduplication, batching

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete notification system design |
| Sources | ✅ | All notification tools documented |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

<!-- Architecture diagram placeholder -->

### Components

<!-- Component description -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

<!-- Interface definitions -->

### Dependencies
<!-- Dependency list -->





## Configuration
### Environment Variables

<!-- Environment variables -->

### Config Keys
<!-- Configuration keys -->


## API Endpoints
### GET /api/v1/user/notifications/preferences

Get user notification preferences

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### PUT /api/v1/user/notifications/preferences

Update notification preferences

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### POST /api/v1/user/push/register

Register FCM device token

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### DELETE /api/v1/user/push/unregister

Unregister device from push

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/user/notifications/history

Get notification history

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### POST /api/v1/admin/webhooks

Register webhook endpoint

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/admin/webhooks

List registered webhooks

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### DELETE /api/v1/admin/webhooks/:id

Delete webhook

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### POST /api/v1/admin/notifications/test

Send test notification (admin only)

**Request**:
```json
{}
```

**Response**:
```json
{}
```


## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [technical](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [EMAIL](../technical/EMAIL.md)
- [WEBHOOKS](../technical/WEBHOOKS.md)

### External Sources
- [go-mail GitHub README](../../sources/tooling/go-mail-guide.md) - Auto-resolved from go-mail
- [go-fcm](../../sources/tooling/go-fcm.md) - FCM push notifications
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx

