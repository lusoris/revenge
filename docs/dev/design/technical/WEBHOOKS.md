# Webhooks

<!-- DESIGN: technical, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Outgoing webhook system for event notifications

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

Revenge supports outgoing webhooks to notify external services of events. This enables integration with:
- Discord/Slack for notifications
- Home automation systems
- Custom monitoring dashboards
- Third-party services

---

## Webhook Events

| Event | Trigger | Payload |
|-------|---------|---------|
| `playback.start` | Media playback begins | Item, user, device |
| `playback.stop` | Media playback ends | Item, user, progress |
| `library.new` | New item added | Item metadata |
| `library.updated` | Item metadata updated | Item metadata |
| `user.created` | New user registered | User info |
| `session.created` | New session started | Session info |

---

## Configuration

```yaml
webhooks:
  endpoints:
    - name: discord
      url: https://discord.com/api/webhooks/...
      events: [playback.start, library.new]
      secret: ${WEBHOOK_SECRET}
    - name: home-assistant
      url: http://homeassistant.local:8123/api/webhook/...
      events: [playback.start, playback.stop]
```

---

## Payload Format

```json
{
  "event": "playback.start",
  "timestamp": "2026-01-31T12:00:00Z",
  "data": {
    "item": { ... },
    "user": { ... },
    "device": { ... }
  },
  "signature": "sha256=..."
}
```

---

## Security

- HMAC-SHA256 signatures for payload verification
- TLS required for all webhook URLs
- Retry with exponential backoff on failure

---

## Related

- [Notification Service](../services/NOTIFICATION.md)
- [Webhook Patterns](../patterns/WEBHOOK_PATTERNS.md)
- [External Services](../integrations/external/INDEX.md)
