# Revenge Roadmap

<!-- DESIGN: planning, README, SCAFFOLD_TEMPLATE, test_output_claude -->

> Version progression from v0.0.0 to v1.0.0

**Last Updated**: 2026-02-01
**Current Version**: v0.0.0 (CI/CD Infrastructure)
**MVP Milestone**: v0.3.0
**Release Candidate**: v0.9.0
**First Stable**: v1.0.0

---

## Milestone Overview

| Version | Codename | Focus | Key Deliverables |
|---------|----------|-------|------------------|
| v0.0.0 | **Foundation** | CI/CD Infrastructure | Pipelines, Deploy configs, Testing framework |
| v0.1.0 | **Skeleton** | Project Structure | Go modules, fx setup, Database schema |
| v0.2.0 | **Core** | Backend Services | Auth, User, Session, RBAC, Library |
| **v0.3.0** | **MVP** | Movie Module | Full backend + Movies + Basic UI |
| v0.4.0 | **Shows** | TV Shows Module | Series, Seasons, Episodes, Sonarr |
| v0.5.0 | **Audio** | Music Module | Artists, Albums, Tracks, Lidarr |
| v0.6.0 | **Playback** | Playback Features | Trickplay, Skip Intro, Watch Next, SyncPlay |
| v0.7.0 | **Media** | Additional Modules | Audiobook, Book, Podcast |
| v0.8.0 | **Intelligence** | Advanced Features | Scrobbling, Analytics, Notifications |
| v0.9.0 | **RC1** | Release Candidate | QAR module, Live TV, Polish, Bug fixes |
| v1.0.0 | **Revenge** | First Stable | All features complete, Production ready |

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
**Focus**: CI/CD Infrastructure

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

### Detailed TODO

→ See [TODO_v0.0.0.md](TODO_v0.0.0.md)

---

## v0.1.0 - Skeleton

**Status**: 🔴 Not Started
**Focus**: Project Structure

### Deliverables

- [ ] Go module structure (`internal/`, `cmd/`, `pkg/`)
- [ ] fx dependency injection setup
- [ ] Configuration system (koanf)
- [ ] Database migrations framework
- [ ] OpenAPI spec skeleton (ogen)
- [ ] Basic health endpoints
- [ ] Logging infrastructure (tint/zap)
- [ ] Error handling patterns

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

## Related Documentation

- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Versions and dependencies
- [VERSIONING.md](../operations/VERSIONING.md) - Semantic versioning strategy
- [DESIGN_INDEX.md](../DESIGN_INDEX.md) - Full design documentation

---

## Changelog

| Date | Version | Change |
|------|---------|--------|
| 2026-02-01 | Initial | Created roadmap from v0.0.0 to v1.0.0 |
