# Revenge Roadmap

<!-- DESIGN: planning, README, test_output_claude, test_output_wiki -->


<!-- TOC-START -->

## Table of Contents

- [Milestone Overview](#milestone-overview)
- [Dependency Graph](#dependency-graph)
- [v0.0.0 - Foundation (COMPLETE)](#v000---foundation-complete)
  - [Deliverables](#deliverables)
  - [Detailed TODO](#detailed-todo)
- [v0.1.0 - Skeleton](#v010---skeleton)
  - [Deliverables](#deliverables)
  - [Dependencies](#dependencies)
  - [Detailed TODO](#detailed-todo)
- [v0.2.0 - Core](#v020---core)
  - [Deliverables](#deliverables)
  - [Dependencies](#dependencies)
  - [Detailed TODO](#detailed-todo)
- [v0.3.0 - MVP (Movies)](#v030---mvp-movies)
  - [Deliverables](#deliverables)
    - [Backend](#backend)
    - [Frontend (Basic)](#frontend-basic)
    - [Infrastructure](#infrastructure)
  - [Dependencies](#dependencies)
  - [Detailed TODO](#detailed-todo)
- [v0.4.0 - Shows](#v040---shows)
  - [Deliverables](#deliverables)
  - [Dependencies](#dependencies)
  - [Detailed TODO](#detailed-todo)
- [v0.5.0 - Audio](#v050---audio)
  - [Deliverables](#deliverables)
  - [Dependencies](#dependencies)
  - [Detailed TODO](#detailed-todo)
- [v0.6.0 - Playback](#v060---playback)
  - [Deliverables](#deliverables)
  - [Dependencies](#dependencies)
  - [Detailed TODO](#detailed-todo)
- [v0.7.0 - Media](#v070---media)
  - [Deliverables](#deliverables)
  - [Dependencies](#dependencies)
  - [Detailed TODO](#detailed-todo)
- [v0.8.0 - Intelligence](#v080---intelligence)
  - [Deliverables](#deliverables)
  - [Dependencies](#dependencies)
  - [Detailed TODO](#detailed-todo)
- [v0.9.0 - RC1 (Release Candidate)](#v090---rc1-release-candidate)
  - [Deliverables](#deliverables)
  - [Dependencies](#dependencies)
  - [Detailed TODO](#detailed-todo)
- [v1.0.0 - Revenge (First Stable)](#v100---revenge-first-stable)
  - [Deliverables](#deliverables)
  - [Dependencies](#dependencies)
  - [Detailed TODO](#detailed-todo)
- [Post-1.0 Roadmap (Future)](#post-10-roadmap-future)
- [Related Documentation](#related-documentation)
- [Changelog](#changelog)

<!-- TOC-END -->


> Version progression from v0.0.0 to v1.0.0

**Last Updated**: 2026-02-02
**Current Phase**: Implementation Phase (Design Phase ✅ Complete)
**Current Version**: v0.1.3 (Skeleton Complete + CI Fixes)
**MVP Milestone**: v0.3.0
**Release Candidate**: v0.9.0
**First Stable**: v1.0.0

---

## Project Phases

### Phase 1: Design Phase ✅ COMPLETE (2026-02-02)

**Deliverables**: 159 design documents covering all features, integrations, services, and architecture

**Scope**:
- ✅ 19 Services (Auth, User, Session, RBAC, Library, Metadata, Search, etc.)
- ✅ 11 Content Modules (Movies, TV, Music, Books, Comics, Audiobooks, Podcasts, Photos, Live TV, Adult)
- ✅ 58 Integrations (Metadata providers, Servarr, Scrobbling, Auth, Infrastructure)
- ✅ 23 Features (Playback, Collections, RBAC, i18n, Request System, etc.)
- ✅ 27 Architecture & Technical docs (API, Frontend, Patterns, Design System, Operations)

**Status**: All major features and integrations are **fully designed**. Implementation can reference YAML docs in `data/` directory.

### Phase 2: Implementation Phase 🔵 IN PROGRESS

**Current Focus**: Building the foundation and MVP (v0.0.0 → v0.3.0)

**Note**: Implementation timeline is **intentionally smaller in scope** than design work. We implement features incrementally (MVP-first), even though designs exist for all features.

---

## Milestone Overview

| Version | Codename | Focus | Key Deliverables | Design Status |
|---------|----------|-------|------------------|---------------|
| v0.0.0 | **Foundation** | CI/CD + Documentation | Pipelines, Deploy configs, YAML data structure, Doc generation | ✅ Complete |
| v0.1.0 | **Skeleton** | Project Structure | Go modules, fx setup, Database schema | ✅ Complete |
| v0.1.1 | Skeleton | Test Coverage Sprint | Database 78% coverage, testcontainers | ✅ Complete |
| v0.1.2 | Skeleton | Errors Coverage | Errors package 100% coverage | ✅ Complete |
| v0.1.3 | Skeleton | CI Fixes | Port conflicts, lint, macOS/Windows | ✅ Complete |
| v0.2.0 | **Core** | Backend Services | Auth, User, Session, RBAC, Library | ✅ Designed |
| **v0.3.0** | **MVP** | Movie Module | Full backend + Movies + Basic UI | ✅ Designed |
| v0.4.0 | **Shows** | TV Shows Module | Series, Seasons, Episodes, Sonarr | ✅ Designed |
| v0.5.0 | **Audio** | Music Module | Artists, Albums, Tracks, Lidarr | ✅ Designed |
| v0.6.0 | **Playback** | Playback Features | Trickplay, Skip Intro, Watch Next, SyncPlay | ✅ Designed |
| v0.7.0 | **Media** | Additional Modules | Audiobook, Book, Podcast | ✅ Designed |
| v0.8.0 | **Intelligence** | Advanced Features | Scrobbling, Analytics, Notifications | ✅ Designed |
| v0.9.0 | **RC1** | Release Candidate | QAR module, Live TV, Polish, Bug fixes | ✅ Designed |
| v1.0.0 | **Revenge** | First Stable | All features complete, Production ready | ✅ Designed |

---

## Dependency Graph

```
v0.0.0 Foundation
    │
    v
v0.1.0 Skeleton ─────────────────────────────────────────┐
    │                                                    │
    v                                                    │
v0.2.0 Core ─────────────────────────────────────────────┤
    │                                                    │
    ├───────────────────────┐                            │
    v                       v                            │
v0.3.0 MVP (Movies)    v0.4.0 Shows (TV)                │
    │                       │                            │
    └───────┬───────────────┘                            │
            v                                            │
       v0.5.0 Audio (Music)                              │
            │                                            │
            v                                            │
       v0.6.0 Playback ◄─────────────────────────────────┤
            │                                            │
            v                                            │
       v0.7.0 Media (Books/Podcasts)                     │
            │                                            │
            v                                            │
       v0.8.0 Intelligence                               │
            │                                            │
            v                                            │
       v0.9.0 RC1 (QAR, LiveTV) ◄────────────────────────┘
            │
            v
       v1.0.0 Revenge (Stable)
```

---

## v0.0.0 - Foundation (COMPLETE)

**Status**: ✅ Complete
**Tag**: `v0.0.0`
**Focus**: CI/CD Infrastructure + Documentation System
**Completed**: 2026-02-02

### Deliverables

- [x] GitHub Actions CI/CD pipelines
- [x] Multi-platform builds (linux, darwin, windows)
- [x] Docker image builds
- [x] Helm chart generation
- [x] Docker Compose / Swarm configs
- [x] Security scanning (CodeQL, Trivy, gosec)
- [x] Documentation validation
- [x] Dependabot configuration
- [x] Release Please setup
- [x] **Phase 5: Documentation Infrastructure**
  - [x] YAML data structure (159 files)
  - [x] Doc generation pipeline (batch_regenerate.py)
  - [x] UTF-8 encoding fixes
  - [x] Windows compatibility fixes
  - [x] CI pipeline validation
  - [x] Yamllint configuration

### Detailed TODO

→ See [TODO_v0.0.0.md](TODO_v0.0.0.md)

---

## v0.1.0 - Skeleton (COMPLETE)

**Status**: ✅ Complete
**Tag**: `v0.1.3`
**Completed**: 2026-02-02
**Focus**: Project Structure

### Patch Releases

| Version | Focus | Key Changes |
|---------|-------|-------------|
| v0.1.0 | Skeleton | Go modules, fx, koanf, pgx, ogen, health, logging, errors |
| v0.1.1 | Test Coverage | Database 22%→78%, testcontainers integration |
| v0.1.2 | Errors Coverage | Errors package 44%→100% coverage |
| v0.1.3 | CI Fixes | Port conflicts resolved, lint fixes, macOS/Windows `-short` |

### Deliverables

- [x] Go module structure (`internal/`, `cmd/`, `pkg/`)
- [x] fx dependency injection setup
- [x] Configuration system (koanf)
- [x] Database migrations framework
- [x] OpenAPI spec skeleton (ogen)
- [x] Basic health endpoints
- [x] Logging infrastructure (tint/JSON)
- [x] Error handling patterns

### Dependencies

| Depends On | Reason |
|------------|--------|
| v0.0.0 | CI/CD must work for builds |

### Detailed TODO

→ See [TODO_v0.1.0.md](TODO_v0.1.0.md)

---

## v0.2.0 - Core

**Status**: 🔴 Not Started
**Focus**: Backend Services

### Deliverables

- [ ] **Auth Service**: Login, logout, registration, password reset
- [ ] **User Service**: Profile management, preferences
- [ ] **Session Service**: Token management, device tracking
- [ ] **RBAC Service**: Casbin integration, role management
- [ ] **API Keys Service**: Key generation, validation
- [ ] **OIDC Service**: SSO provider support
- [ ] **Settings Service**: Server configuration
- [ ] **Activity Service**: Audit logging
- [ ] **Library Service**: Library CRUD, access control
- [ ] **Health Service**: Liveness, readiness probes
- [ ] PostgreSQL integration (pgx)
- [ ] Dragonfly/Redis integration (rueidis)
- [ ] River job queue setup

### Dependencies

| Depends On | Reason |
|------------|--------|
| v0.1.0 | Project structure must exist |

### Detailed TODO

→ See [TODO_v0.2.0.md](TODO_v0.2.0.md)

---

## v0.3.0 - MVP (Movies)

**Status**: 🔴 Not Started
**Focus**: Movie Module + Basic Frontend

### Deliverables

#### Backend
- [ ] **Movie Module**: Repository, Service, Handler
- [ ] **Collection Support**: Movie collections
- [ ] **Metadata Service**: TMDb integration
- [ ] **Search Service**: Typesense integration
- [ ] **Radarr Integration**: Library sync

#### Frontend (Basic)
- [ ] SvelteKit project setup
- [ ] Authentication flow (login/logout)
- [ ] Library browser
- [ ] Movie detail page
- [ ] Basic player integration

#### Infrastructure
- [ ] Typesense deployment configs
- [ ] Full Docker Compose stack
- [ ] Basic documentation

### Dependencies

| Depends On | Reason |
|------------|--------|
| v0.2.0 | Core services required |

### Detailed TODO

→ See [TODO_v0.3.0.md](TODO_v0.3.0.md)

---

## v0.4.0 - Shows

**Status**: 🔴 Not Started
**Focus**: TV Shows Module

### Deliverables

- [ ] **TV Show Module**: Series, Seasons, Episodes
- [ ] **TheTVDB Integration**: Metadata provider
- [ ] **Sonarr Integration**: Library sync
- [ ] **Episode Progress**: Watch tracking
- [ ] **Series Continue Watching**: Resume logic
- [ ] Frontend: Series browser, Episode list

### Dependencies

| Depends On | Reason |
|------------|--------|
| v0.3.0 | Movie module patterns to follow |

### Detailed TODO

→ See [TODO_v0.4.0.md](TODO_v0.4.0.md)

---

## v0.5.0 - Audio

**Status**: 🔴 Not Started
**Focus**: Music Module

### Deliverables

- [ ] **Music Module**: Artists, Albums, Tracks
- [ ] **MusicBrainz Integration**: Metadata
- [ ] **Last.fm Integration**: Metadata enrichment
- [ ] **Lidarr Integration**: Library sync
- [ ] **Audio Player**: Web audio playback
- [ ] **Lyrics Support**: Synced lyrics display
- [ ] Frontend: Music browser, Album view, Player

### Dependencies

| Depends On | Reason |
|------------|--------|
| v0.4.0 | Pattern established for content modules |

### Detailed TODO

→ See [TODO_v0.5.0.md](TODO_v0.5.0.md)

---

## v0.6.0 - Playback

**Status**: 🔴 Not Started
**Focus**: Playback Features

### Deliverables

- [ ] **Trickplay**: Preview thumbnails
- [ ] **Skip Intro/Credits**: Chapter detection
- [ ] **Watch Next**: Continue watching logic
- [ ] **SyncPlay**: Synchronized group playback
- [ ] **Media Enhancements**: Audio boost, subtitles
- [ ] **Chromecast Support**: Cast integration
- [ ] **DLNA Support**: Local network streaming

### Dependencies

| Depends On | Reason |
|------------|--------|
| v0.5.0 | Content modules needed |
| v0.1.0 | Shared playback infrastructure |

### Detailed TODO

→ See [TODO_v0.6.0.md](TODO_v0.6.0.md)

---

## v0.7.0 - Media

**Status**: 🔴 Not Started
**Focus**: Additional Content Modules

### Deliverables

- [ ] **Audiobook Module**: Books, Chapters
- [ ] **Book Module**: eBooks, Reading progress
- [ ] **Podcast Module**: RSS feeds, Episodes
- [ ] **Audnexus Integration**: Audiobook metadata
- [ ] **OpenLibrary Integration**: Book metadata
- [ ] **iTunes/RSS Integration**: Podcast feeds
- [ ] Frontend: Readers, Podcast player

### Dependencies

| Depends On | Reason |
|------------|--------|
| v0.6.0 | Playback infrastructure |

### Detailed TODO

→ See [TODO_v0.7.0.md](TODO_v0.7.0.md)

---

## v0.8.0 - Intelligence

**Status**: 🔴 Not Started
**Focus**: Advanced Features

### Deliverables

- [ ] **Scrobbling Service**: Trakt, Last.fm, ListenBrainz
- [ ] **Analytics Service**: Usage statistics
- [ ] **Notification Service**: Push, Email, Webhooks
- [ ] **Request System**: User content requests
- [ ] **Fingerprint Service**: Media identification
- [ ] **Grants Service**: Sharing with external users
- [ ] **i18n Support**: Multi-language

### Dependencies

| Depends On | Reason |
|------------|--------|
| v0.7.0 | All content modules needed |

### Detailed TODO

→ See [TODO_v0.8.0.md](TODO_v0.8.0.md)

---

## v0.9.0 - RC1 (Release Candidate)

**Status**: 🔴 Not Started
**Focus**: QAR Module, Live TV, Polish

### Deliverables

- [ ] **QAR Module (Adult)**: Isolated schema, full feature set
- [ ] **StashDB Integration**: Adult metadata
- [ ] **Whisparr Integration**: Adult library sync
- [ ] **Live TV/DVR**: TVHeadend, XMLTV
- [ ] **Photos Module**: EXIF, Immich integration
- [ ] **Comics Module**: ComicVine integration
- [ ] **Performance Optimization**: Profiling, caching
- [ ] **Bug Fixes**: All known issues resolved
- [ ] **Documentation**: Complete user guides
- [ ] **Security Audit**: Penetration testing

### Dependencies

| Depends On | Reason |
|------------|--------|
| v0.8.0 | All features must be implemented |

### Detailed TODO

→ See [TODO_v0.9.0.md](TODO_v0.9.0.md)

---

## v1.0.0 - Revenge (First Stable)

**Status**: 🔴 Not Started
**Focus**: Production Ready

### Deliverables

- [ ] All v0.9.0 issues resolved
- [ ] Performance benchmarks met
- [ ] Security audit passed
- [ ] Documentation complete
- [ ] Migration guides from Jellyfin/Plex
- [ ] Official Docker images
- [ ] Helm chart published to GHCR
- [ ] Community contribution guidelines

### Dependencies

| Depends On | Reason |
|------------|--------|
| v0.9.0 | RC testing complete |

### Detailed TODO

→ See [TODO_v1.0.0.md](TODO_v1.0.0.md)

---

## Post-1.0 Roadmap (Future)

Features for consideration after v1.0.0:

- **v1.1.0**: Mobile apps (iOS, Android)
- **v1.2.0**: Voice control (Alexa, Google Home)
- **v1.3.0**: AI-powered recommendations
- **v1.4.0**: Multi-server federation
- **v1.5.0**: Hardware transcoding (NVENC, QSV, VAAPI)

---

## Design Documentation Reference

All features and integrations mentioned in this roadmap have **complete design documentation** in the `data/` directory:

- **159 YAML design documents** covering all planned features
- **Designs grouped by category**: architecture, features, integrations, operations, patterns, research, services, technical
- **Each TODO file** now includes a "Design Documentation" section linking to relevant YAML docs

### Quick Navigation

- [Architecture](../architecture/) - System architecture, design principles, metadata system
- [Services](../services/) - 19 backend services (Auth, User, RBAC, Library, etc.)
- [Features](../features/) - Content modules, playback features, shared features
- [Integrations](../integrations/) - 58 external service integrations
- [Technical](../technical/) - API, Frontend, Design System, Observability
- [Patterns](../patterns/) - Reusable design patterns
- [Operations](../operations/) - Development, deployment, best practices

**Note**: The design phase is **complete**. Implementation follows the milestone schedule above (MVP-first approach).

---

## Related Documentation

- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Versions and dependencies
- [VERSIONING.md](../operations/VERSIONING.md) - Semantic versioning strategy
- [DESIGN_INDEX.md](../DESIGN_INDEX.md) - Full design documentation
- [.workingdir/PLANNING_ANALYSIS.md](../../../.workingdir/PLANNING_ANALYSIS.md) - Design vs. Planning analysis

---

## Changelog

| Date | Version | Change |
|------|---------|--------|
| 2026-02-02 | Update | Added Phase 1 (Design) completion, Design Status column, Design Documentation section |
| 2026-02-01 | Initial | Created roadmap from v0.0.0 to v1.0.0 |
