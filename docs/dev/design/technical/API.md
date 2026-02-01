

---
sources:
  - name: ogen OpenAPI Generator
    url: ../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
  - name: OpenAPI 3.1 Specification
    url: ../../sources/apis/openapi-spec.md
    note: API spec standard
  - name: RFC 7807 Problem Details
    url: https://datatracker.ietf.org/doc/html/rfc7807
    note: Standardized error format
design_refs:
  - title: technical
    path: INDEX.md
  - title: 01_ARCHITECTURE
    path: ../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../architecture/02_DESIGN_PRINCIPLES.md
---

## Table of Contents

- [API Reference](#api-reference)
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


# API Reference


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > REST API reference with OpenAPI specification and ogen code generation

Complete API documentation for Revenge:
- **OpenAPI 3.1**: Full specification with ogen generator
- **Authentication**: Bearer token (JWT) and API keys
- **Versioning**: `/api/v1/` with backward compatibility
- **Rate Limiting**: User-based and IP-based limits
- **Error Handling**: Standardized RFC 7807 Problem Details


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete API reference |
| Sources | ✅ | All API tools documented |
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
- [technical](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)

### External Sources
- [ogen OpenAPI Generator](../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [OpenAPI 3.1 Specification](../../sources/apis/openapi-spec.md) - API spec standard
- [RFC 7807 Problem Details](https://datatracker.ietf.org/doc/html/rfc7807) - Standardized error format

