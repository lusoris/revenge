## Table of Contents

- [Webhook Patterns](#webhook-patterns)
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
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


---
sources:
  - name: River Job Queue
    url: ../sources/tooling/river.md
    note: Background job processing
  - name: crypto/hmac
    url: https://pkg.go.dev/crypto/hmac
    note: HMAC signature validation
  - name: resty
    url: ../sources/tooling/resty.md
    note: HTTP client for webhook delivery
  - name: gobreaker
    url: ../sources/tooling/gobreaker.md
    note: Circuit breaker pattern
design_refs:
  - title: patterns
    path: patterns.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Webhook Patterns

<!-- DESIGN: patterns, README, SCAFFOLD_TEMPLATE, test_output_claude -->


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


## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**


## Related Documentation
### Design Documents
- [patterns](patterns.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [River Job Queue](../sources/tooling/river.md) - Background job processing
- [crypto/hmac](https://pkg.go.dev/crypto/hmac) - HMAC signature validation
- [resty](../sources/tooling/resty.md) - HTTP client for webhook delivery
- [gobreaker](../sources/tooling/gobreaker.md) - Circuit breaker pattern

