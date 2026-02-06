## Table of Contents

- [Webhook Patterns](#webhook-patterns)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Webhook Patterns

<!-- DESIGN: patterns, README, test_output_claude, test_output_wiki -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: pattern


> > Secure, reliable webhook handling with HMAC validation, async processing, and retry logic

Standard webhook handling pattern:
- **Security**: HMAC signature validation, API key auth, IP whitelisting
- **Async Processing**: Queue events via River for non-blocking response
- **Deduplication**: Event ID tracking to prevent duplicate processing
- **Retry Logic**: Exponential backoff for failed webhook processing
- **Logging**: Comprehensive webhook event and error logging

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete webhook handling patterns |
| Sources | ✅ | All webhook tools documented |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete


## Related Documentation
### Design Documents
- [patterns](INDEX.md)
- [01_ARCHITECTURE](../architecture/ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/METADATA_SYSTEM.md)

### External Sources
- [River Job Queue](../../sources/tooling/river.md) - Background job processing
- [crypto/hmac](https://pkg.go.dev/crypto/hmac) - HMAC signature validation
- [resty](../../sources/tooling/resty.md) - HTTP client for webhook delivery
- [gobreaker](../../sources/tooling/gobreaker.md) - Circuit breaker pattern

