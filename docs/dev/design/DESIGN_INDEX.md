# Design Documentation Index

> Index of all active design documentation for Revenge

---

## Architecture

Core system design and architectural decisions.

| Document | Description |
|----------|-------------|
| [ARCHITECTURE.md](architecture/ARCHITECTURE.md) | System structure, fx modules, dependency layers |
| [DESIGN_PRINCIPLES.md](architecture/DESIGN_PRINCIPLES.md) | Patterns, error handling, testing conventions |
| [METADATA_SYSTEM.md](architecture/METADATA_SYSTEM.md) | Provider chain, adapters, caching strategy |
| [PLAYER_ARCHITECTURE.md](architecture/PLAYER_ARCHITECTURE.md) | Media playback with HLS streaming |
| [PLUGIN_DECISION.md](architecture/PLUGIN_DECISION.md) | ADR: integrations over plugins |

---

## Services

Backend services in `internal/service/`.

| Document | Package | Description |
|----------|---------|-------------|
| [AUTH.md](services/AUTH.md) | `service/auth` | Authentication, registration, password management |
| [SESSION.md](services/SESSION.md) | `service/session` | Session tokens and device tracking |
| [MFA.md](services/MFA.md) | `service/mfa` | TOTP, WebAuthn, backup codes |
| [OIDC.md](services/OIDC.md) | `service/oidc` | OpenID Connect / SSO providers |
| [USER.md](services/USER.md) | `service/user` | User account management |
| [RBAC.md](services/RBAC.md) | `service/rbac` | Role-based access control (Casbin) |
| [APIKEYS.md](services/APIKEYS.md) | `service/apikeys` | API key generation and validation |
| [METADATA.md](services/METADATA.md) | `service/metadata` | External metadata provider orchestration |
| [LIBRARY.md](services/LIBRARY.md) | `service/library` | Library management and access control |
| [SEARCH.md](services/SEARCH.md) | `service/search` | Full-text search via Typesense |
| [SETTINGS.md](services/SETTINGS.md) | `service/settings` | Server settings persistence |
| [ACTIVITY.md](services/ACTIVITY.md) | `service/activity` | Audit logging and event tracking |
| [EMAIL.md](services/EMAIL.md) | `service/email` | Email sending (SMTP/SendGrid) |
| [NOTIFICATION.md](services/NOTIFICATION.md) | `service/notification` | Multi-channel notifications |
| [STORAGE.md](services/STORAGE.md) | `service/storage` | File storage (local/S3) |

---

## Content Modules

Content domain modules in `internal/content/`.

| Document | Description |
|----------|-------------|
| [MOVIE_MODULE.md](features/video/MOVIE_MODULE.md) | Movie entities, service, repository, jobs |
| [TVSHOW_MODULE.md](features/video/TVSHOW_MODULE.md) | TV show/season/episode entities, service, jobs |
| [SHARED_CONTENT.md](features/shared/SHARED_CONTENT.md) | Shared content packages (progress, watched, files) |
| [ADULT_CONTENT_SYSTEM.md](features/adult/ADULT_CONTENT_SYSTEM.md) | QAR module (pirate-themed adult content) |
| [LIBRARIES.md](features/shared/LIBRARIES.md) | Library management overview |
| [RBAC.md](features/shared/RBAC.md) | Dynamic RBAC with Casbin |

---

## Infrastructure

Internal infrastructure packages in `internal/infra/`.

| Document | Package | Description |
|----------|---------|-------------|
| [DATABASE.md](infrastructure/DATABASE.md) | `infra/database` | PostgreSQL pooling, migrations, sqlc |
| [CACHE.md](infrastructure/CACHE.md) | `infra/cache` | L1 (otter) + L2 (rueidis/Dragonfly) cache |
| [JOBS.md](infrastructure/JOBS.md) | `infra/jobs` | River job queue, priority queues, workers |
| [HEALTH.md](infrastructure/HEALTH.md) | `infra/health` | K8s probes, dependency checks |
| [IMAGE.md](infrastructure/IMAGE.md) | `infra/image` | Image download, proxy, caching |
| [LOGGING.md](infrastructure/LOGGING.md) | `infra/logging` | slog + zap, tint dev mode |
| [OBSERVABILITY.md](infrastructure/OBSERVABILITY.md) | `infra/observability` | Prometheus metrics, pprof |
| [SEARCH_INFRA.md](infrastructure/SEARCH_INFRA.md) | `infra/search` | Typesense client wrapper |

---

## Integrations

External service integrations.

### Metadata Providers

| Document | Provider | Description |
|----------|----------|-------------|
| [TMDB.md](integrations/metadata/video/TMDB.md) | TMDb | Movie/TV metadata (API v3, priority 100) |
| [THETVDB.md](integrations/metadata/video/THETVDB.md) | TheTVDB | TV metadata (API v4, JWT auth, priority 80) |

### Servarr

| Document | Service | Description |
|----------|---------|-------------|
| [RADARR.md](integrations/servarr/RADARR.md) | Radarr | Movie library sync + webhooks |
| [SONARR.md](integrations/servarr/SONARR.md) | Sonarr | TV library sync + webhooks |

### Authentication

| Document | Provider | Description |
|----------|----------|-------------|
| [GENERIC_OIDC.md](integrations/auth/GENERIC_OIDC.md) | OIDC | Generic OpenID Connect provider |

### Infrastructure

| Document | Service | Description |
|----------|---------|-------------|
| [POSTGRESQL.md](integrations/infrastructure/POSTGRESQL.md) | PostgreSQL 18 | Primary database |
| [DRAGONFLY.md](integrations/infrastructure/DRAGONFLY.md) | Dragonfly | Redis-compatible cache |
| [TYPESENSE.md](integrations/infrastructure/TYPESENSE.md) | Typesense | Full-text search engine |
| [RIVER.md](integrations/infrastructure/RIVER.md) | River | PostgreSQL-backed job queue |

---

## Technical

API specs, configuration, testing, and cross-cutting concerns.

| Document | Description |
|----------|-------------|
| [API.md](technical/API.md) | REST API (124 endpoints, ogen codegen) |
| [CONFIGURATION.md](technical/CONFIGURATION.md) | koanf config system with all defaults |
| [TESTING.md](technical/TESTING.md) | 157 test files, testcontainers integration |
| [TECH_STACK.md](technical/TECH_STACK.md) | Technology choices and rationale |
| [FRONTEND.md](technical/FRONTEND.md) | SvelteKit frontend architecture |
| [EMAIL.md](technical/EMAIL.md) | Email system design |
| [NOTIFICATION_CHANNELS.md](technical/NOTIFICATION_CHANNELS.md) | Notification channel design |
| [OBSERVABILITY.md](technical/OBSERVABILITY.md) | Observability stack design |
| [OFFLOADING.md](technical/OFFLOADING.md) | Background worker offloading |

---

## Operations

Development setup, deployment, and CI/CD.

| Document | Description |
|----------|-------------|
| [DEVELOPMENT.md](operations/DEVELOPMENT.md) | Development environment, Makefile, devcontainer |
| [SETUP.md](operations/SETUP.md) | Production deployment (Docker, Helm) |
| [CI_CD.md](operations/CI_CD.md) | GitHub Actions workflows |
| [CODING_STANDARDS.md](operations/CODING_STANDARDS.md) | Development best practices |
| [GITFLOW.md](operations/GITFLOW.md) | Git workflow and branching strategy |
| [BRANCHES.md](operations/BRANCHES.md) | Branch protection rules |
| [VERSIONING.md](operations/VERSIONING.md) | Semantic versioning and releases |
| [AUTO_HEALING.md](operations/AUTO_HEALING.md) | Database auto-healing |
| [PROXY.md](operations/PROXY.md) | Reverse proxy configuration |

---

## Patterns

Reusable implementation patterns.

| Document | Description |
|----------|-------------|
| [SERVARR.md](patterns/SERVARR.md) | Arr integration pattern |
| [METADATA.md](patterns/METADATA.md) | Metadata enrichment pattern |
| [HTTP_CLIENT.md](patterns/HTTP_CLIENT.md) | HTTP client with proxy/VPN |
| [OBSERVABILITY.md](patterns/OBSERVABILITY.md) | Observability pattern |
| [TESTING.md](patterns/TESTING.md) | Testing patterns |
| [WEBHOOKS.md](patterns/WEBHOOKS.md) | Webhook patterns |

---

## Planning

Roadmap and milestone tracking.

| Document | Description |
|----------|-------------|
| [ROADMAP.md](planning/ROADMAP.md) | Version progression v0.0.0 to v1.0.0 |
| [TODO_v0.0.0.md](planning/TODO_v0.0.0.md) | Foundation (CI/CD) |
| [TODO_v0.1.0.md](planning/TODO_v0.1.0.md) | Skeleton (project structure) |
| [TODO_v0.2.0.md](planning/TODO_v0.2.0.md) | Core (backend services) |
| [TODO_v0.3.0.md](planning/TODO_v0.3.0.md) | MVP (movies) |
| [TODO_v0.4.0.md](planning/TODO_v0.4.0.md) | Shows (TV) |
| [TODO_v0.5.0.md](planning/TODO_v0.5.0.md) | Audio (music) |
| [TODO_v0.6.0.md](planning/TODO_v0.6.0.md) | Playback features |
| [TODO_v0.7.0.md](planning/TODO_v0.7.0.md) | Additional content modules |
| [TODO_v0.8.0.md](planning/TODO_v0.8.0.md) | Advanced features |
| [TODO_v0.9.0.md](planning/TODO_v0.9.0.md) | RC1 (QAR, Live TV, polish) |
| [TODO_v1.0.0.md](planning/TODO_v1.0.0.md) | First stable release |

---

## Research

| Document | Description |
|----------|-------------|
| [USER_PAIN_POINTS_RESEARCH.md](research/USER_PAIN_POINTS_RESEARCH.md) | User problems with existing media servers |
| [UX_UI_RESOURCES.md](research/UX_UI_RESOURCES.md) | Frontend design resources |

---

## Planned Features

Designed but not yet implemented. Located in [planned/](planned/).

Includes: music, audiobook, book, comics, photos, podcasts, live TV modules; playback features (skip intro, SyncPlay, trickplay); additional integrations (MusicBrainz, Last.fm, Trakt, StashDB); design system components; and more.

---

## Status Legend

| Status | Meaning |
|--------|---------|
| Active | Document describes implemented code |
| Planned | Document describes future features (in `planned/`) |
