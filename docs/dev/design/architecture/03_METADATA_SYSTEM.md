## Table of Contents

- [Revenge - Metadata System](#revenge-metadata-system)
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
- name: Dragonfly Documentation
    url: ../../sources/infrastructure/dragonfly.md
    note: Auto-resolved from dragonfly
- name: go-blurhash
    url: ../../sources/media/go-blurhash.md
    note: Auto-resolved from go-blurhash
- name: Last.fm API
    url: ../../sources/apis/lastfm.md
    note: Auto-resolved from lastfm-api
- name: pgx PostgreSQL Driver
    url: ../../sources/database/pgx.md
    note: Auto-resolved from pgx
- name: PostgreSQL Arrays
    url: ../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
- name: PostgreSQL JSON Functions
    url: ../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
- name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Auto-resolved from river
- name: rueidis
    url: ../../sources/tooling/rueidis.md
    note: Auto-resolved from rueidis
- name: rueidis GitHub README
    url: ../../sources/tooling/rueidis-guide.md
    note: Auto-resolved from rueidis-docs
design_refs:
- title: architecture
    path: architecture/INDEX.md
- title: ADULT_CONTENT_SYSTEM
    path: ADULT_CONTENT_SYSTEM.md
- title: ADULT_METADATA
    path: ADULT_METADATA.md
- title: DATA_RECONCILIATION
    path: DATA_RECONCILIATION.md
---

# Revenge - Metadata System

<!-- DESIGN: architecture, README, SCAFFOLD_TEMPLATE, test_output_claude -->

**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: architecture

> > Multi-source metadata system with caching and priority chain

Metadata handling:
- **Priority Chain**: Local cache → Arr services → Internal (Stash) → External APIs
- **Providers**: TMDb, TheTVDB, MusicBrainz, StashDB, and many more
- **Caching**: Two-tier with otter (L1 memory) and rueidis (L2 distributed)
- **Enrichment**: Background jobs for additional metadata, thumbnails, blurhash
- **Matching**: Fingerprinting for audio, hash matching for media

---

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | 🟡 | - |
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
- [architecture](architecture/INDEX.md)
- [ADULT_CONTENT_SYSTEM](ADULT_CONTENT_SYSTEM.md)
- [ADULT_METADATA](ADULT_METADATA.md)
- [DATA_RECONCILIATION](DATA_RECONCILIATION.md)

### External Sources
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [go-blurhash](../../sources/media/go-blurhash.md) - Auto-resolved from go-blurhash
- [Last.fm API](../../sources/apis/lastfm.md) - Auto-resolved from lastfm-api
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [rueidis](../../sources/tooling/rueidis.md) - Auto-resolved from rueidis
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md) - Auto-resolved from rueidis-docs
