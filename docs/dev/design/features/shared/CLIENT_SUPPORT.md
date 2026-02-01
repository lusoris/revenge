

---
sources:
  - name: Casbin
    url: ../../../sources/security/casbin.md
    note: Auto-resolved from casbin
  - name: Uber fx
    url: ../../../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: gohlslib (HLS)
    url: ../../../sources/media/gohlslib.md
    note: Auto-resolved from gohlslib
  - name: ogen OpenAPI Generator
    url: ../../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
  - name: pgx PostgreSQL Driver
    url: ../../../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: Svelte 5 Runes
    url: ../../../sources/frontend/svelte-runes.md
    note: Auto-resolved from svelte-runes
  - name: Svelte 5 Documentation
    url: ../../../sources/frontend/svelte5.md
    note: Auto-resolved from svelte5
  - name: SvelteKit Documentation
    url: ../../../sources/frontend/sveltekit.md
    note: Auto-resolved from sveltekit
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [Revenge - Client Support & Device Capabilities](#revenge-client-support-device-capabilities)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Database Schema](#database-schema)
    - [Module Structure](#module-structure)
    - [Component Interaction](#component-interaction)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [Content Management](#content-management)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Revenge - Client Support & Device Capabilities


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for 

> Multi-platform client support with intelligent capability detection.

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | 🟡 | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Database Schema

**Schema**: `public`

<!-- Schema diagram -->

### Module Structure

```
internal/content/revenge___client_support_&_device_capabilities/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── revenge___client_support_&_device_capabilities_test.go
```

### Component Interaction

<!-- Component interaction diagram -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

```go
type ClientService interface {
  RegisterDevice(ctx context.Context, userID uuid.UUID, device ClientDeviceRegister) (*ClientDevice, error)
  GetDevice(ctx context.Context, deviceID string) (*ClientDevice, error)
  UpdateCapabilities(ctx context.Context, deviceID string, capabilities DeviceCapabilities) error
  GetQualityProfile(ctx context.Context, deviceID string) (*QualityProfile, error)
  UpdateQualityProfile(ctx context.Context, deviceID string, profile QualityProfile) error
  DetectCapabilities(ctx context.Context, userAgent string) (*DeviceCapabilities, error)
}
```


### Dependencies

**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/mssola/user_agent` - User agent parsing
- `go.uber.org/fx`






## Configuration
### Environment Variables

```bash
CLIENTS_AUTO_DETECT_CAPABILITIES=true
```


### Config Keys

```yaml
clients:
  auto_detect_capabilities: true
  default_max_resolution: 1080p
  default_max_bitrate_mbps: 20
```



## API Endpoints

### Content Management
```
POST /api/v1/clients/register     # Register device
GET  /api/v1/clients/:id           # Get device info
PUT  /api/v1/clients/:id/profile   # Update quality profile
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
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Casbin](../../../sources/security/casbin.md) - Auto-resolved from casbin
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [gohlslib (HLS)](../../../sources/media/gohlslib.md) - Auto-resolved from gohlslib
- [ogen OpenAPI Generator](../../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [Svelte 5 Runes](../../../sources/frontend/svelte-runes.md) - Auto-resolved from svelte-runes
- [Svelte 5 Documentation](../../../sources/frontend/svelte5.md) - Auto-resolved from svelte5
- [SvelteKit Documentation](../../../sources/frontend/sveltekit.md) - Auto-resolved from sveltekit

