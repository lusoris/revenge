## Table of Contents

- [Generic OIDC](#generic-oidc)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
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
    url: ../sources/security/authelia.md
    note: Auto-resolved from authelia
  - name: Authentik Documentation
    url: ../sources/security/authentik.md
    note: Auto-resolved from authentik
  - name: Keycloak Documentation
    url: ../sources/security/keycloak.md
    note: Auto-resolved from keycloak
  - name: OpenID Connect Core
    url: ../sources/security/oidc-core.md
    note: Auto-resolved from oidc
design_refs:
  - title: integrations/auth
    path: integrations/auth.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Generic OIDC


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Generic OIDC

> Support for any OpenID Connect compatible provider
**API Base URL**: `https://revenge.example.com/api/v1/auth/oidc/callback`
**Authentication**: oauth

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

### Integration Structure

```
internal/integration/generic_oidc/
├── client.go              # API client
├── types.go               # Response types
├── mapper.go              # Map external → internal types
├── cache.go               # Response caching
└── client_test.go         # Tests
```

### Data Flow

<!-- Data flow diagram -->

### Provides

This integration provides:
<!-- Data provided by integration -->


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
- [integrations/auth](integrations/auth.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Authelia Documentation](../sources/security/authelia.md) - Auto-resolved from authelia
- [Authentik Documentation](../sources/security/authentik.md) - Auto-resolved from authentik
- [Keycloak Documentation](../sources/security/keycloak.md) - Auto-resolved from keycloak
- [OpenID Connect Core](../sources/security/oidc-core.md) - Auto-resolved from oidc

