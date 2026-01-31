## Table of Contents

- [Whisparr v3 & StashDB Schema Integration](#whisparr-v3-stashdb-schema-integration)
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



---
sources:
  - name: FFmpeg Documentation
    url: https://ffmpeg.org/ffmpeg.html
    note: Auto-resolved from ffmpeg
  - name: FFmpeg Codecs
    url: https://ffmpeg.org/ffmpeg-codecs.html
    note: Auto-resolved from ffmpeg-codecs
  - name: FFmpeg Formats
    url: https://ffmpeg.org/ffmpeg-formats.html
    note: Auto-resolved from ffmpeg-formats
  - name: go-astiav (FFmpeg bindings)
    url: https://pkg.go.dev/github.com/asticode/go-astiav
    note: Auto-resolved from go-astiav
  - name: go-astiav GitHub README
    url: https://github.com/asticode/go-astiav
    note: Auto-resolved from go-astiav-docs
  - name: pgx PostgreSQL Driver
    url: https://pkg.go.dev/github.com/jackc/pgx/v5
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: https://www.postgresql.org/docs/current/arrays.html
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: https://www.postgresql.org/docs/current/functions-json.html
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: https://pkg.go.dev/github.com/riverqueue/river
    note: Auto-resolved from river
  - name: shadcn-svelte
    url: https://www.shadcn-svelte.com/docs
    note: Auto-resolved from shadcn-svelte
  - name: StashDB GraphQL API
    url: https://stashdb.org/graphql
    note: Auto-resolved from stashdb
  - name: Svelte 5 Runes
    url: https://svelte.dev/docs/svelte/$state
    note: Auto-resolved from svelte-runes
  - name: Svelte 5 Documentation
    url: https://svelte.dev/docs/svelte/overview
    note: Auto-resolved from svelte5
  - name: SvelteKit Documentation
    url: https://svelte.dev/docs/kit/introduction
    note: Auto-resolved from sveltekit
  - name: Whisparr OpenAPI Spec
    url: https://raw.githubusercontent.com/Whisparr/Whisparr/develop/src/Whisparr.Api.V3/openapi.json
    note: Auto-resolved from whisparr
design_refs:
  - title: features/adult
    path: features/adult.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Whisparr v3 & StashDB Schema Integration


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for 

> Custom UI/UX approach for adult content scenes using Whisparr cache

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

### Database Schema

**Schema**: `qar`

<!-- Schema diagram -->

### Module Structure

```
internal/content/whisparr_v3_&_stashdb_schema_integration/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── whisparr_v3_&_stashdb_schema_integration_test.go
```

### Component Interaction

<!-- Component interaction diagram -->


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

### Content Management
<!-- API endpoints placeholder -->


## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [features/adult](features/adult.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [FFmpeg Documentation](https://ffmpeg.org/ffmpeg.html) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](https://ffmpeg.org/ffmpeg-codecs.html) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](https://ffmpeg.org/ffmpeg-formats.html) - Auto-resolved from ffmpeg-formats
- [go-astiav (FFmpeg bindings)](https://pkg.go.dev/github.com/asticode/go-astiav) - Auto-resolved from go-astiav
- [go-astiav GitHub README](https://github.com/asticode/go-astiav) - Auto-resolved from go-astiav-docs
- [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) - Auto-resolved from pgx
- [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) - Auto-resolved from postgresql-json
- [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) - Auto-resolved from river
- [shadcn-svelte](https://www.shadcn-svelte.com/docs) - Auto-resolved from shadcn-svelte
- [StashDB GraphQL API](https://stashdb.org/graphql) - Auto-resolved from stashdb
- [Svelte 5 Runes](https://svelte.dev/docs/svelte/$state) - Auto-resolved from svelte-runes
- [Svelte 5 Documentation](https://svelte.dev/docs/svelte/overview) - Auto-resolved from svelte5
- [SvelteKit Documentation](https://svelte.dev/docs/kit/introduction) - Auto-resolved from sveltekit
- [Whisparr OpenAPI Spec](https://raw.githubusercontent.com/Whisparr/Whisparr/develop/src/Whisparr.Api.V3/openapi.json) - Auto-resolved from whisparr

