# Developer Documentation

> Design documentation, research, and external sources

**Last Updated**: 2026-01-29

---

## Design Document Template

All new design documents should follow the standard template for consistency:

- **[DESIGN_DOC_TEMPLATE.md](design/DESIGN_DOC_TEMPLATE.md)** - Template and guidelines for writing design docs

---

## Design Documentation 🔒

Project design documents - manually maintained, protected from auto-updates.

### [Architecture](design/architecture/)
Core system design and architectural decisions.

- [Overall Architecture V2](design/architecture/ARCHITECTURE_V2.md) - Complete system architecture
- [Design Principles](design/architecture/DESIGN_PRINCIPLES.md) - Guiding principles
- [Metadata System](design/architecture/METADATA_SYSTEM.md) - Metadata management
- [Player Architecture](design/architecture/PLAYER_ARCHITECTURE.md) - Media player design
- [Plugin Architecture Decision](design/architecture/PLUGIN_ARCHITECTURE_DECISION.md) - Native vs plugins

### [Features](design/features/)
Content modules and feature specifications. See [Features INDEX](design/features/INDEX.md) for full listing.

**Shared Features** (`design/features/shared/`):
- Analytics, RBAC, Access Controls, Request System
- News System, Wiki System, Ticketing
- Client Support, Voice Control, Scrobbling, i18n

**Playback Features** (`design/features/playback/`):
- SyncPlay, Trickplay, Skip Intro/Credits
- Watch Next, Release Calendar, Media Enhancements

**Module-Specific:**
- `photos/` - Photos Library
- `podcasts/` - RSS Podcasts
- `livetv/` - Live TV / DVR
- `comics/` - Comics/Manga
- `adult/` - Adult content (isolated in `c` schema)

### [Integrations](design/integrations/)
External service integration designs.

- [Index](design/integrations/INDEX.md) - Integration overview
- **Anime**: AniList, MyAnimeList, Kitsu
- **Auth**: Authelia, Authentik, Keycloak, OIDC
- **Audiobook**: Audiobookshelf
- **Casting**: Chromecast, DLNA
- **External**: Social media, adult platforms
- **Infrastructure**: PostgreSQL, Dragonfly, Typesense, River
- **Live TV**: TVHeadend, NextPVR
- **Metadata**: TMDb, MusicBrainz, StashDB, ComicVine, etc.
- **Scrobbling**: Trakt, Last.fm, ListenBrainz, Letterboxd, Simkl
- **Servarr**: Radarr, Sonarr, Lidarr, Whisparr, Chaptarr
- **Transcoding**: Blackbeard
- **Wiki**: Wikipedia, Fandom, TVTropes

### [Operations](design/operations/)
Setup, deployment, and operational guides.

- [Setup Guide](design/operations/SETUP.md) - Initial setup instructions
- [Development Guide](design/operations/DEVELOPMENT.md) - Development workflow
- [Database Auto-Healing](design/operations/DATABASE_AUTO_HEALING.md) - PostgreSQL corruption detection/repair
- [Reverse Proxy](design/operations/REVERSE_PROXY.md) - Nginx/Caddy configuration
- [Git Workflow](design/operations/GITFLOW.md) - Branch strategy
- [Branch Protection](design/operations/BRANCH_PROTECTION.md) - Repository protection rules
- [Best Practices](design/operations/BEST_PRACTICES.md) - Coding standards
- [Upstream Sync](design/operations/UPSTREAM_SYNC.md) - Syncing with upstream

### [Planning](design/planning/)
Roadmaps, versioning, and implementation phases.

- [Module Implementation Roadmap](design/planning/MODULE_IMPLEMENTATION_TODO.md) - Module rollout plan
- [Version Policy](design/planning/VERSION_POLICY.md) - Bleeding edge/latest stable policy
- [Versioning Strategy](design/planning/VERSIONING.md) - Semantic versioning approach

### [Research](design/research/)
Analysis, user research, and technology evaluations.

- [User Pain Points Research](design/research/USER_PAIN_POINTS_RESEARCH.md) - Jellyfin/Plex/Emby issues
- [UX/UI Resources](design/research/UX_UI_RESOURCES.md) - Design patterns and resources

### [Technical](design/technical/)
API design, frontend architecture, and technical documentation.

- [API Documentation](design/technical/API.md) - REST API design
- [Frontend Architecture](design/technical/FRONTEND.md) - Svelte 5 + shadcn-svelte
- [Tech Stack](design/technical/TECH_STACK.md) - Technologies used
- [Audio Streaming](design/technical/AUDIO_STREAMING.md) - Audio streaming architecture
- [Offloading](design/technical/OFFLOADING.md) - Task offloading patterns

### [Services](design/services/)
Core application services implementing business logic.

- [Services Overview](design/services/INDEX.md) - Service layer architecture
- [Auth Service](design/services/AUTH.md) - Authentication, registration, password management
- [User Service](design/services/USER.md) - User CRUD, roles, profile management
- [Session Service](design/services/SESSION.md) - Session tokens, device tracking
- [Library Service](design/services/LIBRARY.md) - Library CRUD, access control
- [Metadata Service](design/services/METADATA.md) - TMDb, Radarr providers
- [RBAC Service](design/services/RBAC.md) - Casbin permission management
- [OIDC Service](design/services/OIDC.md) - SSO provider management
- [API Keys Service](design/services/APIKEYS.md) - API key management
- [Activity Service](design/services/ACTIVITY.md) - Audit logging
- [Settings Service](design/services/SETTINGS.md) - Server settings

---

## External Sources 🔄

Auto-fetched external documentation. Updated weekly via CI.

**Registry**: [sources/SOURCES.yaml](sources/SOURCES.yaml)
**Status**: [sources/INDEX.yaml](sources/INDEX.yaml)

### Source Categories

| Category | Description | Content |
|----------|-------------|---------|
| [go](sources/go/) | Go language & stdlib | context, slog, net/http, testing |
| [apis](sources/apis/) | External APIs | TMDb, MusicBrainz, Trakt, AniList |
| [protocols](sources/protocols/) | Streaming protocols | HLS, DASH, HTTP Range |
| [database](sources/database/) | Database documentation | PostgreSQL, sqlc patterns |
| [frontend](sources/frontend/) | Frontend technologies | Svelte 5, TanStack, shadcn |
| [tooling](sources/tooling/) | Go tooling | ogen, river, koanf, fx |
| [media](sources/media/) | Media handling | FFmpeg, codecs, containers |
| [security](sources/security/) | Security standards | OIDC, OAuth 2.0, PKCE |
| [testing](sources/testing/) | Testing patterns | Go testing, testify |
| [observability](sources/observability/) | Monitoring | Prometheus, OpenTelemetry |
| [infrastructure](sources/infrastructure/) | Infrastructure | Dragonfly, Typesense |

### Fetcher Usage

```bash
# Fetch all sources (runs weekly via CI)
python scripts/fetch-sources.py

# Fetch specific category
python scripts/fetch-sources.py --category go

# Fetch single source
python scripts/fetch-sources.py --id tmdb

# Dry run (show what would be fetched)
python scripts/fetch-sources.py --dry-run
```

**Requirements**: `pip install -r scripts/requirements-fetch.txt`

### How It Works

1. **Weekly CI** runs `fetch-sources.py` on Sunday 03:00 UTC
2. **Safety check** ensures only `sources/` is modified
3. **PR created** with changes for review
4. **CODEOWNERS** requires @lusoris approval
5. **Protected files** in `design/` are never touched

### Adding New Sources

Edit `sources/SOURCES.yaml`:

```yaml
- id: new-source
  name: "New Source Name"
  url: "https://example.com/docs"
  type: html  # or graphql_schema
  selectors: [".content"]  # CSS selectors (optional)
  output: "category/filename.md"
```

---

## Migration Info

This documentation was restructured on 2026-01-28.
See [MIGRATION_MANIFEST.md](MIGRATION_MANIFEST.md) for details.
