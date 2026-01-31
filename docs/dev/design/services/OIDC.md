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
    url: https://www.authelia.com/overview/
    note: Auto-resolved from authelia
  - name: Authentik Documentation
    url: https://goauthentik.io/docs/
    note: Auto-resolved from authentik
  - name: Uber fx
    url: https://pkg.go.dev/go.uber.org/fx
    note: Auto-resolved from fx
  - name: Keycloak Documentation
    url: https://www.keycloak.org/documentation
    note: Auto-resolved from keycloak
  - name: ogen OpenAPI Generator
    url: https://pkg.go.dev/github.com/ogen-go/ogen
    note: Auto-resolved from ogen
  - name: sqlc
    url: https://docs.sqlc.dev/en/stable/
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: https://docs.sqlc.dev/en/stable/reference/config.html
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
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Authelia Documentation](https://www.authelia.com/overview/) - Auto-resolved from authelia
- [Authentik Documentation](https://goauthentik.io/docs/) - Auto-resolved from authentik
- [Uber fx](https://pkg.go.dev/go.uber.org/fx) - Auto-resolved from fx
- [Keycloak Documentation](https://www.keycloak.org/documentation) - Auto-resolved from keycloak
- [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) - Auto-resolved from ogen
- [sqlc](https://docs.sqlc.dev/en/stable/) - Auto-resolved from sqlc
- [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) - Auto-resolved from sqlc-config

