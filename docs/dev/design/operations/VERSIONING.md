## Table of Contents

- [Semantic Versioning & Releases](#semantic-versioning-releases)
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
- name: Semantic Versioning
    url: ../../sources/standards/semver.md
    note: Auto-resolved from semver
design_refs:
- title: operations
    path: operations/INDEX.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
- title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Semantic Versioning & Releases

<!-- DESIGN: operations, README, SCAFFOLD_TEMPLATE, test_output_claude -->

**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: operations

> > Semantic versioning (semver) and automated release process

  Versioning strategy:
- **Format**: vMAJOR.MINOR.PATCH (e.g., v1.2.3)
- **MAJOR**: Breaking changes to API or behavior
- **MINOR**: New features, backwards compatible
- **PATCH**: Bug fixes, security patches
- **Automation**: Release Please for changelog and version bumps

---

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete versioning guide |
| Sources | 🔴 | - |
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
- [operations](operations/INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [pgx PostgreSQL Driver](../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river
- [Semantic Versioning](../sources/standards/semver.md) - Auto-resolved from semver
