## Table of Contents

- [EPG Service](#epg-service)
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
  - [API Endpoints](#api-endpoints)
    - [GET /api/v1/epg/channels](#get-apiv1epgchannels)
    - [GET /api/v1/epg/channels/{id}/schedule](#get-apiv1epgchannelsidschedule)
    - [GET /api/v1/epg/programs/{id}](#get-apiv1epgprogramsid)
    - [GET /api/v1/epg/search](#get-apiv1epgsearch)
    - [POST /api/v1/epg/refresh](#post-apiv1epgrefresh)
    - [GET /api/v1/epg/stats](#get-apiv1epgstats)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# EPG Service


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service


> > TV program schedule and guide data management service

EPG capabilities:
- **Format**: XMLTV standard for program data
- **Sources**: TVHeadend, NextPVR, ErsatzTV integration
- **Refresh**: Automatic scheduled updates every 6 hours
- **Search**: Full-text search via Typesense for programs
- **Cache**: Dragonfly cache for frequently accessed schedules
- **API**: RESTful endpoints for channel listings and program queries

**Package**: `internal/service/epg`
**fx Module**: `epg.Module`

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete EPG service design |
| Sources | ✅ | XMLTV and Live TV sources documented |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Service Structure

```
internal/service/epg/
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


## API Endpoints
### GET /api/v1/epg/channels

List all EPG channels

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/epg/channels/{id}/schedule

Get channel schedule

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/epg/programs/{id}

Get program details

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/epg/search

Search EPG programs

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### POST /api/v1/epg/refresh

Trigger EPG refresh (admin only)

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/epg/stats

Get EPG statistics

**Request**:
```json
{}
```

**Response**:
```json
{}
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
- [services](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [LIVE_TV_DVR](../features/livetv/LIVE_TV_DVR.md)
- [TVHEADEND](../integrations/livetv/TVHEADEND.md)
- [NEXTPVR](../integrations/livetv/NEXTPVR.md)
- [ERSATZTV](../integrations/livetv/ERSATZTV.md)

### External Sources
- [XMLTV Format](http://wiki.xmltv.org/index.php/XMLTVFormat) - EPG data format standard
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md) - Full-text search for programs
- [River Job Queue](../../sources/tooling/river.md) - Scheduled EPG refresh jobs
- [Uber fx](../../sources/tooling/fx.md) - Dependency injection

