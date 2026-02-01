## Table of Contents

- [Webhooks](#webhooks)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Webhooks


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > Bidirectional webhook system for integrations and event notifications

Complete webhook infrastructure:
- **Incoming**: Receive events from Radarr, Sonarr, Lidarr, Whisparr, Stash
- **Outgoing**: Send events to Discord, Slack, custom endpoints
- **Security**: HMAC SHA-256 signature verification
- **Reliability**: Exponential backoff retries via River queue
- **Events**: 30+ event types for all content changes

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete webhook system design |
| Sources | ✅ | All webhook tools documented |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---



















## Related Documentation
### Design Documents
- [technical](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [WEBHOOK_PATTERNS](../patterns/WEBHOOK_PATTERNS.md)
- [ARR_INTEGRATION](../patterns/ARR_INTEGRATION.md)
- [NOTIFICATIONS](../technical/NOTIFICATIONS.md)

### External Sources
- [Radarr API Docs](../../sources/apis/radarr-docs.md) - Auto-resolved from radarr-docs
- [Sonarr API Docs](../../sources/apis/sonarr-docs.md) - Auto-resolved from sonarr-docs
- [Lidarr API Docs](../../sources/apis/lidarr-docs.md) - Auto-resolved from lidarr-docs
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx

