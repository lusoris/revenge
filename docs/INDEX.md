# Revenge Documentation Index

> Central navigation for all project documentation

**Last Updated**: 2026-01-29

## Quick Start

- **New to Revenge?** Start with [Setup Guide](dev/design/operations/SETUP.md)
- **Development?** See [Development Guide](dev/design/operations/DEVELOPMENT.md)
- **Architecture?** Read [Architecture](dev/design/architecture/ARCHITECTURE.md)
- **Current Status?** See [TODO.md](../TODO.md)

---

## Documentation Structure

```
docs/
├── INDEX.md              # This file
└── dev/                  # Developer documentation
    ├── INDEX.md          # Main dev docs index
    ├── design/           # Design documents
    │   ├── architecture/ # System architecture
    │   ├── features/     # Feature specifications
    │   ├── integrations/ # External integrations
    │   ├── operations/   # Setup & operations
    │   ├── planning/     # Roadmaps & versioning
    │   ├── research/     # User research & analysis
    │   └── technical/    # API & tech specs
    └── sources/          # Auto-fetched external docs
```

**Main Entry Point**: [Developer Documentation](dev/INDEX.md)

---

## Architecture

Core system design and architectural decisions.

- [Architecture](dev/design/architecture/ARCHITECTURE.md) - Complete system architecture
- [Design Principles](dev/design/architecture/DESIGN_PRINCIPLES.md) - Guiding principles
- [Player Architecture](dev/design/architecture/PLAYER_ARCHITECTURE.md) - Media player design
- [Metadata System](dev/design/architecture/METADATA_SYSTEM.md) - Metadata management
- [Plugin Architecture Decision](dev/design/architecture/PLUGIN_DECISION.md) - Native vs plugins

## Operations

Setup, deployment, and operational guides.

- [Setup Guide](dev/design/operations/SETUP.md) - Initial setup instructions
- [Development Guide](dev/design/operations/DEVELOPMENT.md) - Development workflow
- [Coding Standards](dev/design/operations/CODING_STANDARDS.md) - Coding standards
- [Database Auto-Healing](dev/design/operations/AUTO_HEALING.md) - PostgreSQL corruption detection/repair
- [Reverse Proxy](dev/design/operations/PROXY.md) - Nginx/Caddy configuration
- [Git Workflow](dev/design/operations/GITFLOW.md) - Branch strategy
- [Branch Protection](dev/design/operations/BRANCHES.md) - Repository protection rules
- [Upstream Sync](dev/design/operations/UPSTREAM_SYNC.md) - Syncing with upstream

## Features

Content modules and feature designs. See [Features Index](dev/design/features/INDEX.md) for full listing.

### Shared Features
- [Analytics Service](dev/design/features/shared/ANALYTICS_SERVICE.md) - Usage analytics
- [RBAC](dev/design/features/shared/RBAC.md) - Role-based access control
- [Access Controls](dev/design/features/shared/ACCESS_CONTROLS.md) - Permission system
- [Request System](dev/design/features/shared/REQUEST_SYSTEM.md) - Content requests
- [Content Rating](dev/design/features/shared/CONTENT_RATING.md) - Rating systems (MPAA/PEGI/etc.)
- [Scrobbling](dev/design/features/shared/SCROBBLING.md) - Trakt/Last.fm/ListenBrainz
- [Client Support](dev/design/features/shared/CLIENT_SUPPORT.md) - Client compatibility
- [Internationalization](dev/design/features/shared/I18N.md) - i18n/l10n support
- [Library Types](dev/design/features/shared/LIBRARIES.md) - Library management
- [News System](dev/design/features/shared/NEWS_SYSTEM.md) - News/announcements
- [Wiki System](dev/design/features/shared/WIKI_SYSTEM.md) - Wiki integration
- [Ticketing System](dev/design/features/shared/TICKETING_SYSTEM.md) - Support tickets
- [Voice Control](dev/design/features/shared/VOICE_CONTROL.md) - Voice commands
- [User Experience](dev/design/features/shared/USER_EXPERIENCE_FEATURES.md) - UX enhancements

### Playback Features
- [Watch Next / Continue Watching](dev/design/features/playback/WATCH_NEXT_CONTINUE_WATCHING.md) - Playback continuation
- [Release Calendar](dev/design/features/playback/RELEASE_CALENDAR.md) - Upcoming releases
- [Skip Intro/Credits](dev/design/features/playback/SKIP_INTRO.md) - Auto-skip
- [SyncPlay](dev/design/features/playback/SYNCPLAY.md) - Synchronized playback
- [Trickplay](dev/design/features/playback/TRICKPLAY.md) - Video thumbnails
- [Media Enhancements](dev/design/features/playback/MEDIA_ENHANCEMENTS.md) - Media-specific features

### Module-Specific
- [Adult Content System](dev/design/features/adult/ADULT_CONTENT_SYSTEM.md) - QAR module isolation (`qar` schema)
- [Adult Metadata](dev/design/features/adult/ADULT_METADATA.md) - Whisparr/StashDB integration
- [Comics Module](dev/design/features/comics/COMICS_MODULE.md) - Digital comics/manga
- [Live TV / DVR](dev/design/features/livetv/LIVE_TV_DVR.md) - Live television
- [Photos Library](dev/design/features/photos/PHOTOS_LIBRARY.md) - Photo management
- [Podcasts](dev/design/features/podcasts/PODCASTS.md) - Podcast support

## Integrations

External service integrations. See [Integrations Index](dev/design/integrations/INDEX.md) for full listing.

- **Servarr**: [Radarr](dev/design/integrations/servarr/RADARR.md), [Sonarr](dev/design/integrations/servarr/SONARR.md), [Lidarr](dev/design/integrations/servarr/LIDARR.md), [Whisparr](dev/design/integrations/servarr/WHISPARR.md), [Chaptarr](dev/design/integrations/servarr/CHAPTARR.md)
- **Metadata**: [TMDb](dev/design/integrations/metadata/video/TMDB.md), [TheTVDB](dev/design/integrations/metadata/video/THETVDB.md), [MusicBrainz](dev/design/integrations/metadata/music/MUSICBRAINZ.md), [StashDB](dev/design/integrations/metadata/adult/STASHDB.md), [ComicVine](dev/design/integrations/metadata/comics/COMICVINE.md)
- **Scrobbling**: [Trakt](dev/design/integrations/scrobbling/TRAKT.md), [Last.fm](dev/design/integrations/scrobbling/LASTFM_SCROBBLE.md), [ListenBrainz](dev/design/integrations/scrobbling/LISTENBRAINZ.md), [Letterboxd](dev/design/integrations/scrobbling/LETTERBOXD.md)
- **Auth**: [Authelia](dev/design/integrations/auth/AUTHELIA.md), [Authentik](dev/design/integrations/auth/AUTHENTIK.md), [Keycloak](dev/design/integrations/auth/KEYCLOAK.md), [OIDC](dev/design/integrations/auth/GENERIC_OIDC.md)
- **Anime**: [AniList](dev/design/integrations/anime/ANILIST.md), [MyAnimeList](dev/design/integrations/anime/MYANIMELIST.md), [Kitsu](dev/design/integrations/anime/KITSU.md)
- **Infrastructure**: [PostgreSQL](dev/design/integrations/infrastructure/POSTGRESQL.md), [Dragonfly](dev/design/integrations/infrastructure/DRAGONFLY.md), [Typesense](dev/design/integrations/infrastructure/TYPESENSE.md), [River](dev/design/integrations/infrastructure/RIVER.md)

## Technical

API design, frontend architecture, and technical documentation.

- [API Documentation](dev/design/technical/API.md) - REST API design
- [Frontend Architecture](dev/design/technical/FRONTEND.md) - Svelte 5 + shadcn-svelte
- [Tech Stack](dev/design/technical/TECH_STACK.md) - Technologies used
- [Audio Streaming](dev/design/technical/AUDIO_STREAMING.md) - Audio streaming architecture
- [Offloading](dev/design/technical/OFFLOADING.md) - Task offloading patterns

## Research

Analysis, user research, and technology evaluations.

- [User Pain Points Research](dev/design/research/USER_PAIN_POINTS_RESEARCH.md) - Jellyfin/Plex/Emby issues
- [UX/UI Resources](dev/design/research/UX_UI_RESOURCES.md) - Design patterns and resources

## Planning

Roadmaps, versioning, and implementation tracking.

- [Module Implementation Roadmap](dev/design/planning/MODULE_IMPLEMENTATION_TODO.md) - Module rollout plan
- [Version Policy](dev/design/planning/VERSION_POLICY.md) - Bleeding edge/latest stable policy
- [Versioning Strategy](dev/design/planning/VERSIONING.md) - Semantic versioning approach

**Current Progress**: See [TODO.md](../TODO.md) for active tasks

## External Sources

Auto-fetched external documentation. Updated weekly via CI.

- **Registry**: [sources/SOURCES.yaml](dev/sources/SOURCES.yaml)
- **Fetcher**: `python scripts/fetch-sources.py`
- See [Developer Documentation](dev/INDEX.md#external-sources-) for details

---

## Related Resources

- [README](../README.md) - Project overview
- [TODO List](../TODO.md) - Project backlog
- [Contributing Guide](../CONTRIBUTING.md) - Contribution guidelines
- [AGENTS.md](../AGENTS.md) - Automated coding agent rules
