## Table of Contents

- [OIDC Service](#oidc-service)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Service Structure](#service-structure)
    - [Dependencies](#dependencies)
    - [Provides](#provides)
    - [Component Diagram](#component-diagram)
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
- name: Authelia Documentation
    url: ../../sources/security/authelia.md
    note: Auto-resolved from authelia
- name: Authentik Documentation
    url: ../../sources/security/authentik.md
    note: Auto-resolved from authentik
- name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Auto-resolved from fx
- name: Keycloak Documentation
    url: ../../sources/security/keycloak.md
    note: Auto-resolved from keycloak
- name: ogen OpenAPI Generator
    url: ../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
- name: sqlc
    url: ../../sources/database/sqlc.md
    note: Auto-resolved from sqlc
- name: sqlc Configuration
    url: ../../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
design_refs:
- title: services
    path: services/INDEX.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
- title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# OIDC Service

<!-- DESIGN: services, README, SCAFFOLD_TEMPLATE, test_output_claude -->

**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service

> > OpenID Connect / SSO provider management

**Package**: `internal/service/oidc`
**fx Module**: `oidc.Module`

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

### Service Structure

```
internal/service/oidc/
├── module.go              # fx module definition
├── service.go             # Service implementation
├── repository.go          # Data access (if needed)
├── handler.go             # HTTP handlers (if exposed)
├── middleware.go          # Middleware (if needed)
├── types.go               # Domain types
└── service_test.go        # Tests
```

### Dependencies
No external service dependencies.

### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->

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
- [services](services/INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Authelia Documentation](../sources/security/authelia.md) - Auto-resolved from authelia
- [Authentik Documentation](../sources/security/authentik.md) - Auto-resolved from authentik
- [Uber fx](../sources/tooling/fx.md) - Auto-resolved from fx
- [Keycloak Documentation](../sources/security/keycloak.md) - Auto-resolved from keycloak
- [ogen OpenAPI Generator](../sources/tooling/ogen.md) - Auto-resolved from ogen
- [sqlc](../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config
