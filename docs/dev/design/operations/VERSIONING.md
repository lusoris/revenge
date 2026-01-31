# Versioning Strategy



<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Semantic Versioning (SemVer)](#semantic-versioning-semver)
  - [Version Components](#version-components)
- [Release Phases](#release-phases)
  - [Pre-1.0 Development (Current)](#pre-10-development-current)
  - [Post-1.0 (Future)](#post-10-future)
- [Release Process](#release-process)
  - [Automated Releases (Release Please)](#automated-releases-release-please)
  - [Manual Releases](#manual-releases)
- [Version in Code](#version-in-code)
- [Pre-release Versions](#pre-release-versions)
- [API Versioning](#api-versioning)
- [Deprecation Policy](#deprecation-policy)
- [Version Compatibility Matrix](#version-compatibility-matrix)
- [Checking Version](#checking-version)
- [Related Design Docs](#related-design-docs)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🔴 |  |
| Sources | 🔴 |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |

---

This document describes the versioning strategy for revenge.

## Semantic Versioning (SemVer)

revenge follows [Semantic Versioning 2.0.0](https://semver.org/):

```
MAJOR.MINOR.PATCH

v0.1.0  →  v0.2.0  →  v0.10.0  →  v1.0.0  →  v1.1.0
```

### Version Components

| Component | When to Increment | Example |
|-----------|-------------------|---------|
| **MAJOR** | Breaking API changes | v1.0.0 → v2.0.0 |
| **MINOR** | New features (backward compatible) | v1.0.0 → v1.1.0 |
| **PATCH** | Bug fixes (backward compatible) | v1.0.0 → v1.0.1 |

## Release Phases

### Pre-1.0 Development (Current)

During the `v0.x` phase:
- **Breaking changes** increment MINOR (0.1.0 → 0.2.0)
- **New features** increment MINOR (0.1.0 → 0.2.0)
- **Bug fixes** increment PATCH (0.1.0 → 0.1.1)
- API stability is **not guaranteed**

```
v0.1.x  Phase 1: Core Foundation (Auth, Library, Direct Play)
v0.2.x  Phase 1: Media Management
v0.3.x  Phase 1: MVP Complete
v0.4.x  Phase 2: Transcoding
v0.5.x  Phase 2: Metadata & Search
v0.6.x  Phase 2: Advanced Features
v0.7.x  Phase 3: Plugins & Extensions
v0.8.x  Phase 3: Performance & Polish
v0.9.x  Phase 3: Release Candidates
v1.0.0  Stable: Feature Parity with Revenge
```

### Post-1.0 (Future)

After `v1.0.0`:
- Full backward compatibility within major versions
- Deprecation notices before breaking changes
- LTS (Long Term Support) for major versions

## Release Process

### Automated Releases (Release Please)

Releases are automated using [Release Please](https://github.com/googleapis/release-please):

1. **Conventional Commits** trigger release PRs:
   ```bash
   feat: add user authentication    # → Minor bump
   fix: resolve memory leak         # → Patch bump
   feat!: redesign API              # → Major bump (post-1.0)
   ```

2. **Release PR** is created automatically with:
   - Updated version in code
   - Generated CHANGELOG
   - Release notes

3. **Merging** the Release PR:
   - Creates GitHub Release
   - Triggers GoReleaser
   - Publishes Docker images

### Manual Releases

For special releases:

```bash
# Create release tag
git tag v0.1.0
git push origin v0.1.0

# GoReleaser will handle the rest
```

## Version in Code

Version information is embedded at build time:

```go
// internal/version/version.go
var (
    Version   = "dev"      // Set by ldflags
    GitCommit = "unknown"  // Set by ldflags
    BuildTime = "unknown"  // Set by ldflags
)
```

Build with version:
```bash
go build -ldflags "-X main.Version=v0.1.0 -X main.GitCommit=$(git rev-parse --short HEAD)" ./cmd/revenge
```

## Pre-release Versions

For testing before stable releases:

| Type | Format | Example | Use Case |
|------|--------|---------|----------|
| Alpha | `v0.1.0-alpha.1` | Early testing | Internal testing |
| Beta | `v0.1.0-beta.1` | Feature complete | Community testing |
| RC | `v0.1.0-rc.1` | Release candidate | Final validation |

## API Versioning

The Revenge API maintains compatibility with the original Revenge:

- **Current**: API v1 (Revenge compatible)
- **Future**: API v2 (Go-optimized, optional)

API version is indicated in:
- URL: `/api/v1/...`
- Header: `X-Revenge-API-Version: 1`

## Deprecation Policy

1. **Deprecation Notice**: Feature marked deprecated in documentation
2. **Warning Period**: Minimum 2 minor versions
3. **Removal**: Only in major version bumps (post-1.0)

## Version Compatibility Matrix

| revenge | Revenge API | Go Version | Database |
|-------------|--------------|------------|----------|
| v0.1.x | v1 | 1.23+ | PostgreSQL 18+ |
| v0.2.x | v1 | 1.24+ | PostgreSQL 18+ |
| v1.0.x | v1 | 1.24+ | PostgreSQL 18+ |

## Checking Version

```bash
# Binary
./revenge --version

# API
curl http://localhost:8096/System/Info
```

Response:
```json
{
  "version": "0.1.0",
  "git_commit": "abc1234",
  "go_version": "go1.24",
  "os": "linux",
  "arch": "amd64"
}
```


<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Planning](INDEX.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

---

<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../sources/database/postgresql-json.md) |
| [Semantic Versioning](https://semver.org/) | [Local](../../sources/standards/semver.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../sources/database/pgx.md) |

<!-- SOURCE-BREADCRUMBS-END -->