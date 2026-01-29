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
- [Player Architecture](design/architecture/PLAYER_ARCHITECTURE.md) - Media player design
- [Metadata System](design/architecture/METADATA_SYSTEM.md) - Metadata management
- [Plugin Architecture Decision](design/architecture/PLUGIN_ARCHITECTURE_DECISION.md) - Native vs plugins

### [Features](design/features/)
Content modules and feature specifications.

- [Adult Content System](design/features/ADULT_CONTENT_SYSTEM.md) - Adult module isolation (`c` schema)
- [Adult Metadata](design/features/ADULT_METADATA.md) - Whisparr/StashDB integration
- [Analytics Service](design/features/ANALYTICS_SERVICE.md) - Analytics design
- [Client Support](design/features/CLIENT_SUPPORT.md) - Client compatibility
- [Comics Module](design/features/COMICS_MODULE.md) - Digital comics/manga support
- [Content Rating](design/features/CONTENT_RATING.md) - Rating systems (MPAA/PEGI/etc.)
- [Internationalization](design/features/I18N.md) - i18n/l10n support
- [Library Types](design/features/LIBRARY_TYPES.md) - Library management
- [Media Enhancements](design/features/MEDIA_ENHANCEMENTS.md) - Media-specific features
- [NSFW Toggle](design/features/NSFW_TOGGLE.md) - Adult content visibility toggle
- [RBAC with Casbin](design/features/RBAC_CASBIN.md) - Dynamic role-based access control
- [Release Calendar](design/features/RELEASE_CALENDAR.md) - Upcoming releases via Servarr
- [Request System](design/features/REQUEST_SYSTEM.md) - Content request system
- [Scrobbling](design/features/SCROBBLING.md) - Trakt/Last.fm/ListenBrainz
- [Watch Next / Continue Watching](design/features/WATCH_NEXT_CONTINUE_WATCHING.md) - Playback continuation system
- [Ticketing System](design/features/TICKETING_SYSTEM.md) - Support tickets
- [User Experience Features](design/features/USER_EXPERIENCE_FEATURES.md) - UX enhancements
- [Whisparr/StashDB Schema](design/features/WHISPARR_STASHDB_SCHEMA.md) - Adult metadata schema

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
- **Servarr**: Radarr, Sonarr, Lidarr, Whisparr, Readarr
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
- [Go Packages Research](design/research/GO_PACKAGES_RESEARCH.md) - Awesome-Go analysis (500+ packages)
- [Documentation Gap Analysis](../archive/reports/DOCUMENTATION_GAP_ANALYSIS.md) - Documentation coverage audit (archived)

### [Technical](design/technical/)
API design, frontend architecture, and technical documentation.

- [API Documentation](design/technical/API.md) - REST API design
- [Frontend Architecture](design/technical/FRONTEND.md) - Svelte 5 + shadcn-svelte
- [Tech Stack](design/technical/TECH_STACK.md) - Technologies used
- [Audio Streaming](design/technical/AUDIO_STREAMING.md) - Audio streaming architecture
- [Offloading](design/technical/OFFLOADING.md) - Task offloading patterns

### Archived Planning Docs
Historical planning and analysis documents.

- [Documentation Analysis](../archive/reports/DOCUMENTATION_ANALYSIS.md) - Historical analysis
- [Preparation Master Plan](../archive/planning/PREPARATION_MASTER_PLAN.md) - Pre-implementation checklist
- [Module Implementation TODO](../archive/planning/MODULE_IMPLEMENTATION_TODO.md) - Module rollout (archived)
- [Restructuring Plan](../archive/planning/RESTRUCTURING_PLAN.md) - Codebase restructuring plan

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
