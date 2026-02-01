## Table of Contents

- [Email System](#email-system)
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
    - [POST /api/v1/admin/email/send](#post-apiv1adminemailsend)
    - [GET /api/v1/user/email/preferences](#get-apiv1useremailpreferences)
    - [PUT /api/v1/user/email/preferences](#put-apiv1useremailpreferences)
    - [GET /unsubscribe](#get-unsubscribe)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Email System


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > SMTP email system with go-mail, TLS/STARTTLS, templates, async delivery

Complete email infrastructure:
- **Library**: go-mail (wneessen/go-mail) with connection pooling
- **Security**: TLS/STARTTLS support, SMTP auth (PLAIN, LOGIN, CRAM-MD5)
- **Templates**: HTML with text fallback using Go templates
- **Delivery**: Async via River job queue with retry logic
- **Features**: Bounce handling, unsubscribe links, rate limiting

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete email system design |
| Sources | ✅ | go-mail documentation included |
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
### POST /api/v1/admin/email/send

Send test email (admin only)

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/user/email/preferences

Get email notification preferences

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### PUT /api/v1/user/email/preferences

Update email notification preferences

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /unsubscribe

Unsubscribe from emails (public, no auth)

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
- [NOTIFICATIONS](../technical/NOTIFICATIONS.md)

### External Sources
- [go-mail GitHub README](../../sources/tooling/go-mail-guide.md) - Auto-resolved from go-mail-docs
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx

