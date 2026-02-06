## Table of Contents

- [API Reference](#api-reference)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# API Reference

<!-- DESIGN: technical, README, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES -->


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


## Related Documentation
### Design Documents
- [technical](INDEX.md)
- [01_ARCHITECTURE](../architecture/ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/DESIGN_PRINCIPLES.md)

### External Sources
- [ogen OpenAPI Generator](../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [OpenAPI 3.1 Specification](../../sources/apis/openapi-spec.md) - API spec standard
- [RFC 7807 Problem Details](https://datatracker.ietf.org/doc/html/rfc7807) - Standardized error format

