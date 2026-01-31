# Design Documentation Index

> Auto-generated cross-reference index for all design documents

**Source of Truth**: [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md)

---

## Start Here

> Recommended reading order for new developers

### 1. Understand the System

| Step | Document | Time | Description |
|------|----------|------|-------------|
| 1 | [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md) | 10 min | Package versions, module status, core decisions |
| 2 | [Architecture](architecture/01_ARCHITECTURE.md) | 20 min | System overview, components, data flow |
| 3 | [Design Principles](architecture/02_DESIGN_PRINCIPLES.md) | 10 min | Coding conventions, patterns, best practices |

### 2. Pick Your Area

| If you're working on... | Start with |
|-------------------------|------------|
| Content modules (Movies, TV, Music) | [Features Index](features/INDEX.md) → Module doc |
| Backend services | [Services Index](services/INDEX.md) → Service doc |
| External integrations | [Integrations Index](integrations/INDEX.md) → Provider doc |
| Deployment/Operations | [Operations Index](operations/INDEX.md) |
| Frontend | [FRONTEND.md](technical/FRONTEND.md) |
| API design | [API.md](technical/API.md) |

### 3. Navigation Aids

- [NAVIGATION.md](NAVIGATION.md) - Full navigation map with all categories
- [SOURCES_INDEX.md](../sources/SOURCES_INDEX.md) - External documentation sources
- [DESIGN_CROSSREF.md](../sources/DESIGN_CROSSREF.md) - Which docs link to which sources

---

## Quick Stats

- **Total Documents**: 126
- **Categories**: 30
- **Topics**: 16

---

## By Category

### Architecture

| Document | Topics | Links |
|----------|--------|-------|
| [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md) | authentication, metadata, playback | 16 |
| [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md) | authentication, metadata, playback | 10 |
| [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md) | authentication, metadata, playback | 8 |
| [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md) | authentication, metadata, playback | 11 |
| [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) | authentication, metadata, playback | 5 |

### Features → Adult

| Document | Topics | Links |
|----------|--------|-------|
| [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md) | authentication, metadata, search | 5 |
| [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md) | metadata, search, database | 3 |
| [Adult Data Reconciliation](features/adult/DATA_RECONCILIATION.md) | metadata, search, database | 2 |
| [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md) | authentication, metadata, search | 3 |
| [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md) | authentication, metadata, playback | 3 |

### Features → Comics

| Document | Topics | Links |
|----------|--------|-------|
| [Comics Module](features/comics/COMICS_MODULE.md) | authentication, metadata, playback | 3 |

### Features → Livetv

| Document | Topics | Links |
|----------|--------|-------|
| [Live TV & DVR](features/livetv/LIVE_TV_DVR.md) | authentication, metadata, playback | 3 |

### Features → Photos

| Document | Topics | Links |
|----------|--------|-------|
| [Photos Library](features/photos/PHOTOS_LIBRARY.md) | authentication, metadata, playback | 2 |

### Features → Playback

| Document | Topics | Links |
|----------|--------|-------|
| [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md) | authentication, metadata, playback | 2 |
| [Release Calendar System](features/playback/RELEASE_CALENDAR.md) | authentication, metadata, search | 0 |
| [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md) | metadata, playback, search | 2 |
| [SyncPlay (Watch Together)](features/playback/SYNCPLAY.md) | authentication, playback, search | 0 |
| [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md) | metadata, playback, search | 2 |
| [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md) | authentication, metadata, playback | 1 |

### Features → Podcasts

| Document | Topics | Links |
|----------|--------|-------|
| [Podcasts](features/podcasts/PODCASTS.md) | authentication, metadata, playback | 2 |

### Features → Shared

| Document | Topics | Links |
|----------|--------|-------|
| [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md) | authentication, playback, search | 3 |
| [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md) | authentication, playback, search | 3 |
| [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md) | authentication, metadata, playback | 4 |
| [Content Rating System](features/shared/CONTENT_RATING.md) | authentication, metadata, playback | 3 |
| [Revenge - Internationalization (i18n)](features/shared/I18N.md) | authentication, metadata, playback | 0 |
| [Library Types](features/shared/LIBRARY_TYPES.md) | authentication, metadata, search | 8 |
| [News System](features/shared/NEWS_SYSTEM.md) | authentication, metadata, playback | 3 |
| [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md) | authentication, metadata, search | 2 |
| [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md) | authentication, metadata, playback | 9 |
| [Native Request System](features/shared/REQUEST_SYSTEM.md) | metadata, playback, search | 5 |
| [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md) | authentication, metadata, playback | 2 |
| [Ticketing System](features/shared/TICKETING_SYSTEM.md) | authentication, metadata, playback | 0 |
| [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md) | authentication, metadata, playback | 7 |
| [Voice Control](features/shared/VOICE_CONTROL.md) | authentication, playback, search | 4 |
| [Internal Wiki System](features/shared/WIKI_SYSTEM.md) | authentication, metadata, playback | 2 |

### Features → Video

| Document | Topics | Links |
|----------|--------|-------|
| [Movie Module](features/video/MOVIE_MODULE.md) | authentication, metadata, search | 7 |
| [TV Show Module](features/video/TVSHOW_MODULE.md) | authentication, metadata, playback | 7 |

### Integrations → Anime

| Document | Topics | Links |
|----------|--------|-------|
| [AniList Integration](integrations/anime/ANILIST.md) | authentication, metadata, playback | 6 |
| [Kitsu Integration](integrations/anime/KITSU.md) | authentication, metadata, search | 5 |
| [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md) | authentication, metadata, search | 5 |

### Integrations → Auth

| Document | Topics | Links |
|----------|--------|-------|
| [Authelia Integration](integrations/auth/AUTHELIA.md) | authentication, metadata, playback | 6 |
| [Authentik Integration](integrations/auth/AUTHENTIK.md) | authentication, metadata, playback | 6 |
| [Generic OIDC Integration](integrations/auth/GENERIC_OIDC.md) | authentication, metadata, search | 6 |
| [Keycloak Integration](integrations/auth/KEYCLOAK.md) | authentication, metadata, playback | 6 |

### Integrations → Casting

| Document | Topics | Links |
|----------|--------|-------|
| [Chromecast Integration](integrations/casting/CHROMECAST.md) | authentication, metadata, playback | 3 |
| [DLNA/UPnP Integration](integrations/casting/DLNA.md) | authentication, metadata, playback | 3 |

### Integrations → External → Adult

| Document | Topics | Links |
|----------|--------|-------|
| [FreeOnes Integration](integrations/metadata/adult/FREEONES.md) | authentication, metadata, playback | 5 |
| [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md) | authentication, metadata, playback | 2 |
| [OnlyFans Integration](integrations/metadata/adult/ONLYFANS.md) | metadata, search, database | 3 |
| [Pornhub Integration](integrations/metadata/adult/PORNHUB.md) | authentication, metadata, playback | 4 |
| [TheNude Integration](integrations/metadata/adult/THENUDE.md) | authentication, metadata, search | 2 |
| [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md) | authentication, metadata, playback | 2 |

### Integrations → Infrastructure

| Document | Topics | Links |
|----------|--------|-------|
| [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md) | authentication, metadata, playback | 6 |
| [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md) | authentication, metadata, playback | 7 |
| [River Integration](integrations/infrastructure/RIVER.md) | authentication, metadata, playback | 7 |
| [Typesense Integration](integrations/infrastructure/TYPESENSE.md) | authentication, playback, search | 6 |

### Integrations → Livetv

| Document | Topics | Links |
|----------|--------|-------|
| [ErsatzTV Integration](integrations/livetv/ERSATZTV.md) | authentication, metadata, playback | 4 |
| [NextPVR Integration](integrations/livetv/NEXTPVR.md) | authentication, metadata, playback | 3 |
| [TVHeadend Integration](integrations/livetv/TVHEADEND.md) | authentication, metadata, playback | 4 |

### Integrations → Metadata → Adult

| Document | Topics | Links |
|----------|--------|-------|
| [Stash Integration](integrations/metadata/adult/STASH.md) | authentication, metadata, playback | 7 |
| [StashDB Integration](integrations/metadata/adult/STASHDB.md) | authentication, metadata, playback | 9 |
| [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md) | authentication, metadata, search | 7 |
| [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md) | metadata, playback, search | 5 |

### Integrations → Metadata → Books

| Document | Topics | Links |
|----------|--------|-------|
| [Audible Integration](integrations/metadata/books/AUDIBLE.md) | authentication, metadata, playback | 2 |
| [Goodreads Integration](integrations/metadata/books/GOODREADS.md) | authentication, metadata, search | 6 |
| [Hardcover Integration](integrations/metadata/books/HARDCOVER.md) | authentication, metadata, search | 5 |
| [OpenLibrary Integration](integrations/metadata/books/OPENLIBRARY.md) | authentication, metadata, search | 5 |

### Integrations → Metadata → Comics

| Document | Topics | Links |
|----------|--------|-------|
| [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md) | authentication, metadata, search | 0 |
| [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md) | metadata, playback, search | 0 |
| [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md) | authentication, metadata, search | 0 |

### Integrations → Metadata → Music

| Document | Topics | Links |
|----------|--------|-------|
| [Discogs Integration](integrations/metadata/music/DISCOGS.md) | authentication, metadata, playback | 1 |
| [Last.fm Integration](integrations/metadata/music/LASTFM.md) | authentication, metadata, playback | 6 |
| [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md) | authentication, metadata, search | 8 |
| [Spotify Integration](integrations/metadata/music/SPOTIFY.md) | authentication, metadata, playback | 4 |

### Integrations → Metadata → Video

| Document | Topics | Links |
|----------|--------|-------|
| [OMDb (Open Movie Database) Integration](integrations/metadata/video/OMDB.md) | authentication, metadata, playback | 2 |
| [ThePosterDB Integration](integrations/metadata/video/THEPOSTERDB.md) | authentication, metadata, playback | 1 |
| [TheTVDB Integration](integrations/metadata/video/THETVDB.md) | authentication, metadata, search | 4 |
| [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md) | authentication, metadata, search | 8 |

### Integrations → Scrobbling

| Document | Topics | Links |
|----------|--------|-------|
| [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md) | authentication, metadata, playback | 6 |
| [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md) | authentication, metadata, search | 5 |
| [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md) | authentication, metadata, playback | 6 |
| [Simkl Integration](integrations/scrobbling/SIMKL.md) | authentication, metadata, playback | 5 |
| [Trakt Integration](integrations/scrobbling/TRAKT.md) | authentication, metadata, playback | 8 |

### Integrations → Servarr

| Document | Topics | Links |
|----------|--------|-------|
| [Chaptarr Integration](integrations/servarr/CHAPTARR.md) | authentication, metadata, playback | 8 |
| [Lidarr Integration](integrations/servarr/LIDARR.md) | authentication, metadata, playback | 7 |
| [Radarr Integration](integrations/servarr/RADARR.md) | authentication, metadata, playback | 7 |
| [Sonarr Integration](integrations/servarr/SONARR.md) | authentication, metadata, playback | 7 |
| [Whisparr v3 Integration](integrations/servarr/WHISPARR.md) | authentication, metadata, playback | 7 |

### Integrations → Transcoding

| Document | Topics | Links |
|----------|--------|-------|
| [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md) | authentication, metadata, playback | 3 |

### Integrations → Wiki

| Document | Topics | Links |
|----------|--------|-------|
| [FANDOM Integration](integrations/wiki/FANDOM.md) | authentication, metadata, playback | 4 |
| [TVTropes Integration](integrations/wiki/TVTROPES.md) | authentication, metadata, playback | 4 |
| [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md) | authentication, metadata, playback | 5 |

### Integrations → Wiki → Adult

| Document | Topics | Links |
|----------|--------|-------|
| [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md) | authentication, metadata, search | 4 |
| [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md) | authentication, metadata, playback | 5 |
| [IAFD Integration](integrations/wiki/adult/IAFD.md) | authentication, metadata, search | 4 |

### Operations

| Document | Topics | Links |
|----------|--------|-------|
| [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md) | authentication, metadata, playback | 0 |
| [Branch Protection Rules](operations/BRANCH_PROTECTION.md) | authentication, search, api | 0 |
| [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md) | metadata, search, database | 1 |
| [Clone repository](operations/DEVELOPMENT.md) | authentication, playback, search | 3 |
| [GitFlow Workflow Guide](operations/GITFLOW.md) | authentication, playback, search | 0 |
| [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md) | authentication, metadata, playback | 0 |
| [revenge - Setup Guide](operations/SETUP.md) | metadata, playback, search | 2 |

### Planning

| Document | Topics | Links |
|----------|--------|-------|
| [Versioning Strategy](operations/VERSIONING.md) | authentication, metadata, playback | 0 |

### Research

| Document | Topics | Links |
|----------|--------|-------|
| [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md) | authentication, metadata, playback | 0 |
| [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md) | authentication, playback, search | 0 |

### Services

| Document | Topics | Links |
|----------|--------|-------|
| [Activity Service](services/ACTIVITY.md) | authentication, metadata, search | 10 |
| [Analytics Service](services/ANALYTICS.md) | authentication, metadata, playback | 4 |
| [API Keys Service](services/APIKEYS.md) | authentication, playback, search | 6 |
| [Auth Service](services/AUTH.md) | authentication, search, api | 12 |
| [Fingerprint Service](services/FINGERPRINT.md) | metadata, playback, search | 4 |
| [Grants Service](services/GRANTS.md) | playback, search, database | 7 |
| [Library Service](services/LIBRARY.md) | authentication, metadata, playback | 16 |
| [Metadata Service](services/METADATA.md) | metadata, playback, search | 6 |
| [Notification Service](services/NOTIFICATION.md) | playback, search, database | 4 |
| [OIDC Service](services/OIDC.md) | authentication, metadata, search | 6 |
| [RBAC Service](services/RBAC.md) | search, database, api | 8 |
| [Search Service](services/SEARCH.md) | authentication, metadata, playback | 6 |
| [Session Service](services/SESSION.md) | authentication, search, database | 11 |
| [Settings Service](services/SETTINGS.md) | authentication, playback, search | 5 |
| [User Service](services/USER.md) | authentication, search, database | 14 |

### Technical

| Document | Topics | Links |
|----------|--------|-------|
| [API Reference](technical/API.md) | authentication, search, database | 1 |
| [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md) | authentication, metadata, playback | 1 |
| [Configuration Reference](technical/CONFIGURATION.md) | authentication, metadata, search | 3 |
| [Revenge - Frontend Architecture](technical/FRONTEND.md) | authentication, metadata, playback | 2 |
| [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md) | authentication, metadata, playback | 1 |
| [Revenge - Technology Stack](technical/TECH_STACK.md) | authentication, playback, search | 0 |

---

## By Topic

### Adult

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Data Reconciliation](features/adult/DATA_RECONCILIATION.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md)
- [OnlyFans Integration](integrations/metadata/adult/ONLYFANS.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [Auth Service](services/AUTH.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [Search Service](services/SEARCH.md)
- [Settings Service](services/SETTINGS.md)
- [User Service](services/USER.md)
- [API Reference](technical/API.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Api

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Data Reconciliation](features/adult/DATA_RECONCILIATION.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md)
- [SyncPlay (Watch Together)](features/playback/SYNCPLAY.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Voice Control](features/shared/VOICE_CONTROL.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Authelia Integration](integrations/auth/AUTHELIA.md)
- [Authentik Integration](integrations/auth/AUTHENTIK.md)
- [Generic OIDC Integration](integrations/auth/GENERIC_OIDC.md)
- [Keycloak Integration](integrations/auth/KEYCLOAK.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md)
- [OnlyFans Integration](integrations/metadata/adult/ONLYFANS.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [NextPVR Integration](integrations/livetv/NEXTPVR.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Goodreads Integration](integrations/metadata/books/GOODREADS.md)
- [Hardcover Integration](integrations/metadata/books/HARDCOVER.md)
- [OpenLibrary Integration](integrations/metadata/books/OPENLIBRARY.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Discogs Integration](integrations/metadata/music/DISCOGS.md)
- [Last.fm Integration](integrations/metadata/music/LASTFM.md)
- [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [OMDb (Open Movie Database) Integration](integrations/metadata/video/OMDB.md)
- [ThePosterDB Integration](integrations/metadata/video/THEPOSTERDB.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Sonarr Integration](integrations/servarr/SONARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Branch Protection Rules](operations/BRANCH_PROTECTION.md)
- [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [GitFlow Workflow Guide](operations/GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [Versioning Strategy](operations/VERSIONING.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md)
- [Activity Service](services/ACTIVITY.md)
- [Analytics Service](services/ANALYTICS.md)
- [API Keys Service](services/APIKEYS.md)
- [Auth Service](services/AUTH.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Grants Service](services/GRANTS.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [Notification Service](services/NOTIFICATION.md)
- [OIDC Service](services/OIDC.md)
- [RBAC Service](services/RBAC.md)
- [Search Service](services/SEARCH.md)
- [Session Service](services/SESSION.md)
- [Settings Service](services/SETTINGS.md)
- [User Service](services/USER.md)
- [API Reference](technical/API.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Authentication

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [SyncPlay (Watch Together)](features/playback/SYNCPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Voice Control](features/shared/VOICE_CONTROL.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Authelia Integration](integrations/auth/AUTHELIA.md)
- [Authentik Integration](integrations/auth/AUTHENTIK.md)
- [Generic OIDC Integration](integrations/auth/GENERIC_OIDC.md)
- [Keycloak Integration](integrations/auth/KEYCLOAK.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [NextPVR Integration](integrations/livetv/NEXTPVR.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Goodreads Integration](integrations/metadata/books/GOODREADS.md)
- [Hardcover Integration](integrations/metadata/books/HARDCOVER.md)
- [OpenLibrary Integration](integrations/metadata/books/OPENLIBRARY.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Discogs Integration](integrations/metadata/music/DISCOGS.md)
- [Last.fm Integration](integrations/metadata/music/LASTFM.md)
- [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [OMDb (Open Movie Database) Integration](integrations/metadata/video/OMDB.md)
- [ThePosterDB Integration](integrations/metadata/video/THEPOSTERDB.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Sonarr Integration](integrations/servarr/SONARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Branch Protection Rules](operations/BRANCH_PROTECTION.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [GitFlow Workflow Guide](operations/GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [Versioning Strategy](operations/VERSIONING.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md)
- [Activity Service](services/ACTIVITY.md)
- [Analytics Service](services/ANALYTICS.md)
- [API Keys Service](services/APIKEYS.md)
- [Auth Service](services/AUTH.md)
- [Library Service](services/LIBRARY.md)
- [OIDC Service](services/OIDC.md)
- [Search Service](services/SEARCH.md)
- [Session Service](services/SESSION.md)
- [Settings Service](services/SETTINGS.md)
- [User Service](services/USER.md)
- [API Reference](technical/API.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Books

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Keycloak Integration](integrations/auth/KEYCLOAK.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Goodreads Integration](integrations/metadata/books/GOODREADS.md)
- [Hardcover Integration](integrations/metadata/books/HARDCOVER.md)
- [OpenLibrary Integration](integrations/metadata/books/OPENLIBRARY.md)
- [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md)
- [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [Search Service](services/SEARCH.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)

### Caching

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Last.fm Integration](integrations/metadata/music/LASTFM.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [Analytics Service](services/ANALYTICS.md)
- [Grants Service](services/GRANTS.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [Settings Service](services/SETTINGS.md)
- [User Service](services/USER.md)
- [API Reference](technical/API.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Database

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Data Reconciliation](features/adult/DATA_RECONCILIATION.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md)
- [SyncPlay (Watch Together)](features/playback/SYNCPLAY.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Voice Control](features/shared/VOICE_CONTROL.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Authelia Integration](integrations/auth/AUTHELIA.md)
- [Authentik Integration](integrations/auth/AUTHENTIK.md)
- [Generic OIDC Integration](integrations/auth/GENERIC_OIDC.md)
- [Keycloak Integration](integrations/auth/KEYCLOAK.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md)
- [OnlyFans Integration](integrations/metadata/adult/ONLYFANS.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [NextPVR Integration](integrations/livetv/NEXTPVR.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Goodreads Integration](integrations/metadata/books/GOODREADS.md)
- [Hardcover Integration](integrations/metadata/books/HARDCOVER.md)
- [OpenLibrary Integration](integrations/metadata/books/OPENLIBRARY.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Discogs Integration](integrations/metadata/music/DISCOGS.md)
- [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [OMDb (Open Movie Database) Integration](integrations/metadata/video/OMDB.md)
- [ThePosterDB Integration](integrations/metadata/video/THEPOSTERDB.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Sonarr Integration](integrations/servarr/SONARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [GitFlow Workflow Guide](operations/GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [Versioning Strategy](operations/VERSIONING.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [Activity Service](services/ACTIVITY.md)
- [Analytics Service](services/ANALYTICS.md)
- [API Keys Service](services/APIKEYS.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Grants Service](services/GRANTS.md)
- [Library Service](services/LIBRARY.md)
- [Notification Service](services/NOTIFICATION.md)
- [OIDC Service](services/OIDC.md)
- [RBAC Service](services/RBAC.md)
- [Search Service](services/SEARCH.md)
- [Session Service](services/SESSION.md)
- [Settings Service](services/SETTINGS.md)
- [User Service](services/USER.md)
- [API Reference](technical/API.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Frontend

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Data Reconciliation](features/adult/DATA_RECONCILIATION.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md)
- [SyncPlay (Watch Together)](features/playback/SYNCPLAY.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Voice Control](features/shared/VOICE_CONTROL.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Authelia Integration](integrations/auth/AUTHELIA.md)
- [Authentik Integration](integrations/auth/AUTHENTIK.md)
- [Generic OIDC Integration](integrations/auth/GENERIC_OIDC.md)
- [Keycloak Integration](integrations/auth/KEYCLOAK.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md)
- [OnlyFans Integration](integrations/metadata/adult/ONLYFANS.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [NextPVR Integration](integrations/livetv/NEXTPVR.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Goodreads Integration](integrations/metadata/books/GOODREADS.md)
- [Hardcover Integration](integrations/metadata/books/HARDCOVER.md)
- [OpenLibrary Integration](integrations/metadata/books/OPENLIBRARY.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Discogs Integration](integrations/metadata/music/DISCOGS.md)
- [Last.fm Integration](integrations/metadata/music/LASTFM.md)
- [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [OMDb (Open Movie Database) Integration](integrations/metadata/video/OMDB.md)
- [ThePosterDB Integration](integrations/metadata/video/THEPOSTERDB.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Sonarr Integration](integrations/servarr/SONARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Branch Protection Rules](operations/BRANCH_PROTECTION.md)
- [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [GitFlow Workflow Guide](operations/GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [Versioning Strategy](operations/VERSIONING.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md)
- [Activity Service](services/ACTIVITY.md)
- [Analytics Service](services/ANALYTICS.md)
- [API Keys Service](services/APIKEYS.md)
- [Auth Service](services/AUTH.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Grants Service](services/GRANTS.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [Notification Service](services/NOTIFICATION.md)
- [OIDC Service](services/OIDC.md)
- [Search Service](services/SEARCH.md)
- [Session Service](services/SESSION.md)
- [Settings Service](services/SETTINGS.md)
- [User Service](services/USER.md)
- [API Reference](technical/API.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Jobs

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Data Reconciliation](features/adult/DATA_RECONCILIATION.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Discogs Integration](integrations/metadata/music/DISCOGS.md)
- [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md)
- [Analytics Service](services/ANALYTICS.md)
- [API Keys Service](services/APIKEYS.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Grants Service](services/GRANTS.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [Notification Service](services/NOTIFICATION.md)
- [Search Service](services/SEARCH.md)
- [Session Service](services/SESSION.md)
- [User Service](services/USER.md)
- [API Reference](technical/API.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Livetv

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Keycloak Integration](integrations/auth/KEYCLOAK.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [NextPVR Integration](integrations/livetv/NEXTPVR.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [Library Service](services/LIBRARY.md)
- [Notification Service](services/NOTIFICATION.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Metadata

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Data Reconciliation](features/adult/DATA_RECONCILIATION.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Authelia Integration](integrations/auth/AUTHELIA.md)
- [Authentik Integration](integrations/auth/AUTHENTIK.md)
- [Generic OIDC Integration](integrations/auth/GENERIC_OIDC.md)
- [Keycloak Integration](integrations/auth/KEYCLOAK.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md)
- [OnlyFans Integration](integrations/metadata/adult/ONLYFANS.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [NextPVR Integration](integrations/livetv/NEXTPVR.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Goodreads Integration](integrations/metadata/books/GOODREADS.md)
- [Hardcover Integration](integrations/metadata/books/HARDCOVER.md)
- [OpenLibrary Integration](integrations/metadata/books/OPENLIBRARY.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Discogs Integration](integrations/metadata/music/DISCOGS.md)
- [Last.fm Integration](integrations/metadata/music/LASTFM.md)
- [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [OMDb (Open Movie Database) Integration](integrations/metadata/video/OMDB.md)
- [ThePosterDB Integration](integrations/metadata/video/THEPOSTERDB.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Sonarr Integration](integrations/servarr/SONARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [Versioning Strategy](operations/VERSIONING.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [Activity Service](services/ACTIVITY.md)
- [Analytics Service](services/ANALYTICS.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [OIDC Service](services/OIDC.md)
- [Search Service](services/SEARCH.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)

### Music

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md)
- [SyncPlay (Watch Together)](features/playback/SYNCPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md)
- [OnlyFans Integration](integrations/metadata/adult/ONLYFANS.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Discogs Integration](integrations/metadata/music/DISCOGS.md)
- [Last.fm Integration](integrations/metadata/music/LASTFM.md)
- [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Sonarr Integration](integrations/servarr/SONARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [Activity Service](services/ACTIVITY.md)
- [Analytics Service](services/ANALYTICS.md)
- [API Keys Service](services/APIKEYS.md)
- [Auth Service](services/AUTH.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [Search Service](services/SEARCH.md)
- [Session Service](services/SESSION.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Photos

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Authentik Integration](integrations/auth/AUTHENTIK.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [OnlyFans Integration](integrations/metadata/adult/ONLYFANS.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Hardcover Integration](integrations/metadata/books/HARDCOVER.md)
- [OpenLibrary Integration](integrations/metadata/books/OPENLIBRARY.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Discogs Integration](integrations/metadata/music/DISCOGS.md)
- [Last.fm Integration](integrations/metadata/music/LASTFM.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [OMDb (Open Movie Database) Integration](integrations/metadata/video/OMDB.md)
- [ThePosterDB Integration](integrations/metadata/video/THEPOSTERDB.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [Versioning Strategy](operations/VERSIONING.md)
- [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Library Service](services/LIBRARY.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Playback

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md)
- [SyncPlay (Watch Together)](features/playback/SYNCPLAY.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Voice Control](features/shared/VOICE_CONTROL.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Authelia Integration](integrations/auth/AUTHELIA.md)
- [Authentik Integration](integrations/auth/AUTHENTIK.md)
- [Keycloak Integration](integrations/auth/KEYCLOAK.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [NextPVR Integration](integrations/livetv/NEXTPVR.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Discogs Integration](integrations/metadata/music/DISCOGS.md)
- [Last.fm Integration](integrations/metadata/music/LASTFM.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [OMDb (Open Movie Database) Integration](integrations/metadata/video/OMDB.md)
- [ThePosterDB Integration](integrations/metadata/video/THEPOSTERDB.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Sonarr Integration](integrations/servarr/SONARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [GitFlow Workflow Guide](operations/GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [Versioning Strategy](operations/VERSIONING.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md)
- [Analytics Service](services/ANALYTICS.md)
- [API Keys Service](services/APIKEYS.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Grants Service](services/GRANTS.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [Notification Service](services/NOTIFICATION.md)
- [Search Service](services/SEARCH.md)
- [Settings Service](services/SETTINGS.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Scrobbling

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Data Reconciliation](features/adult/DATA_RECONCILIATION.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [SyncPlay (Watch Together)](features/playback/SYNCPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Authelia Integration](integrations/auth/AUTHELIA.md)
- [OnlyFans Integration](integrations/metadata/adult/ONLYFANS.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [NextPVR Integration](integrations/livetv/NEXTPVR.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Hardcover Integration](integrations/metadata/books/HARDCOVER.md)
- [OpenLibrary Integration](integrations/metadata/books/OPENLIBRARY.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Last.fm Integration](integrations/metadata/music/LASTFM.md)
- [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Sonarr Integration](integrations/servarr/SONARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md)
- [GitFlow Workflow Guide](operations/GITFLOW.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [Metadata Service](services/METADATA.md)
- [Search Service](services/SEARCH.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Search

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Data Reconciliation](features/adult/DATA_RECONCILIATION.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Photos Library](features/photos/PHOTOS_LIBRARY.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md)
- [SyncPlay (Watch Together)](features/playback/SYNCPLAY.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Time-Based Access Controls](features/shared/ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Voice Control](features/shared/VOICE_CONTROL.md)
- [Internal Wiki System](features/shared/WIKI_SYSTEM.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Authelia Integration](integrations/auth/AUTHELIA.md)
- [Authentik Integration](integrations/auth/AUTHENTIK.md)
- [Generic OIDC Integration](integrations/auth/GENERIC_OIDC.md)
- [Keycloak Integration](integrations/auth/KEYCLOAK.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [FreeOnes Integration](integrations/metadata/adult/FREEONES.md)
- [Instagram Integration](integrations/metadata/adult/INSTAGRAM.md)
- [OnlyFans Integration](integrations/metadata/adult/ONLYFANS.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [TheNude Integration](integrations/metadata/adult/THENUDE.md)
- [Twitter/X Integration](integrations/metadata/adult/TWITTER_X.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [NextPVR Integration](integrations/livetv/NEXTPVR.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Stash Integration](integrations/metadata/adult/STASH.md)
- [StashDB Integration](integrations/metadata/adult/STASHDB.md)
- [ThePornDB Integration](integrations/metadata/adult/THEPORNDB.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Goodreads Integration](integrations/metadata/books/GOODREADS.md)
- [Hardcover Integration](integrations/metadata/books/HARDCOVER.md)
- [OpenLibrary Integration](integrations/metadata/books/OPENLIBRARY.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Discogs Integration](integrations/metadata/music/DISCOGS.md)
- [Last.fm Integration](integrations/metadata/music/LASTFM.md)
- [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md)
- [Spotify Integration](integrations/metadata/music/SPOTIFY.md)
- [OMDb (Open Movie Database) Integration](integrations/metadata/video/OMDB.md)
- [ThePosterDB Integration](integrations/metadata/video/THEPOSTERDB.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Sonarr Integration](integrations/servarr/SONARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [Boobpedia Integration](integrations/wiki/adult/BOOBPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Branch Protection Rules](operations/BRANCH_PROTECTION.md)
- [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [GitFlow Workflow Guide](operations/GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [Versioning Strategy](operations/VERSIONING.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md)
- [Activity Service](services/ACTIVITY.md)
- [Analytics Service](services/ANALYTICS.md)
- [API Keys Service](services/APIKEYS.md)
- [Auth Service](services/AUTH.md)
- [Fingerprint Service](services/FINGERPRINT.md)
- [Grants Service](services/GRANTS.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [Notification Service](services/NOTIFICATION.md)
- [OIDC Service](services/OIDC.md)
- [RBAC Service](services/RBAC.md)
- [Search Service](services/SEARCH.md)
- [Session Service](services/SESSION.md)
- [Settings Service](services/SETTINGS.md)
- [User Service](services/USER.md)
- [API Reference](technical/API.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

### Video

- [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md)
- [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md)
- [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md)
- [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md)
- [Revenge - Adult Content System](features/adult/ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](features/adult/ADULT_METADATA.md)
- [Adult Gallery Module (QAR: Treasures)](features/adult/GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](features/adult/WHISPARR_STASHDB_SCHEMA.md)
- [Comics Module](features/comics/COMICS_MODULE.md)
- [Live TV & DVR](features/livetv/LIVE_TV_DVR.md)
- [Revenge - Media Enhancement Features](features/playback/MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](features/playback/SKIP_INTRO.md)
- [Trickplay (Timeline Thumbnails)](features/playback/TRICKPLAY.md)
- [Watch Next & Continue Watching System](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md)
- [Podcasts](features/podcasts/PODCASTS.md)
- [Tracearr Analytics Service](features/shared/ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](features/shared/CLIENT_SUPPORT.md)
- [Content Rating System](features/shared/CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Library Types](features/shared/LIBRARY_TYPES.md)
- [News System](features/shared/NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](features/shared/NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md)
- [Native Request System](features/shared/REQUEST_SYSTEM.md)
- [Revenge - External Scrobbling & Sync](features/shared/SCROBBLING.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md)
- [Voice Control](features/shared/VOICE_CONTROL.md)
- [Movie Module](features/video/MOVIE_MODULE.md)
- [TV Show Module](features/video/TVSHOW_MODULE.md)
- [AniList Integration](integrations/anime/ANILIST.md)
- [Kitsu Integration](integrations/anime/KITSU.md)
- [MyAnimeList (MAL) Integration](integrations/anime/MYANIMELIST.md)
- [Chromecast Integration](integrations/casting/CHROMECAST.md)
- [DLNA/UPnP Integration](integrations/casting/DLNA.md)
- [Pornhub Integration](integrations/metadata/adult/PORNHUB.md)
- [Dragonfly Integration](integrations/infrastructure/DRAGONFLY.md)
- [PostgreSQL Integration](integrations/infrastructure/POSTGRESQL.md)
- [River Integration](integrations/infrastructure/RIVER.md)
- [Typesense Integration](integrations/infrastructure/TYPESENSE.md)
- [ErsatzTV Integration](integrations/livetv/ERSATZTV.md)
- [NextPVR Integration](integrations/livetv/NEXTPVR.md)
- [TVHeadend Integration](integrations/livetv/TVHEADEND.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md)
- [Audible Integration](integrations/metadata/books/AUDIBLE.md)
- [Hardcover Integration](integrations/metadata/books/HARDCOVER.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [OMDb (Open Movie Database) Integration](integrations/metadata/video/OMDB.md)
- [ThePosterDB Integration](integrations/metadata/video/THEPOSTERDB.md)
- [TheTVDB Integration](integrations/metadata/video/THETVDB.md)
- [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md)
- [Last.fm Scrobbling Integration](integrations/scrobbling/LASTFM_SCROBBLE.md)
- [Letterboxd Integration](integrations/scrobbling/LETTERBOXD.md)
- [ListenBrainz Integration](integrations/scrobbling/LISTENBRAINZ.md)
- [Simkl Integration](integrations/scrobbling/SIMKL.md)
- [Trakt Integration](integrations/scrobbling/TRAKT.md)
- [Chaptarr Integration](integrations/servarr/CHAPTARR.md)
- [Lidarr Integration](integrations/servarr/LIDARR.md)
- [Radarr Integration](integrations/servarr/RADARR.md)
- [Sonarr Integration](integrations/servarr/SONARR.md)
- [Whisparr v3 Integration](integrations/servarr/WHISPARR.md)
- [Blackbeard Integration](integrations/transcoding/BLACKBEARD.md)
- [FANDOM Integration](integrations/wiki/FANDOM.md)
- [TVTropes Integration](integrations/wiki/TVTROPES.md)
- [Wikipedia Integration](integrations/wiki/WIKIPEDIA.md)
- [Babepedia Integration](integrations/wiki/adult/BABEPEDIA.md)
- [IAFD Integration](integrations/wiki/adult/IAFD.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md)
- [Clone repository](operations/DEVELOPMENT.md)
- [revenge - Setup Guide](operations/SETUP.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [Grants Service](services/GRANTS.md)
- [Library Service](services/LIBRARY.md)
- [Metadata Service](services/METADATA.md)
- [Notification Service](services/NOTIFICATION.md)
- [Search Service](services/SEARCH.md)
- [API Reference](technical/API.md)
- [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md)
- [Configuration Reference](technical/CONFIGURATION.md)
- [Revenge - Frontend Architecture](technical/FRONTEND.md)
- [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)

---

## Most Connected Documents

> Documents with the most internal cross-references

| Document | Links To | Linked From | Total |
|----------|----------|-------------|-------|
| [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md) | 7 | 9 | 16 |
| [Library Service](services/LIBRARY.md) | 6 | 10 | 16 |
| [User Service](services/USER.md) | 6 | 8 | 14 |
| [Auth Service](services/AUTH.md) | 6 | 6 | 12 |
| [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md) | 4 | 7 | 11 |
| [Session Service](services/SESSION.md) | 5 | 6 | 11 |
| [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md) | 5 | 5 | 10 |
| [Activity Service](services/ACTIVITY.md) | 5 | 5 | 10 |
| [Dynamic RBAC with Casbin](features/shared/RBAC_CASBIN.md) | 3 | 6 | 9 |
| [StashDB Integration](integrations/metadata/adult/STASHDB.md) | 5 | 4 | 9 |
| [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md) | 4 | 4 | 8 |
| [Library Types](features/shared/LIBRARY_TYPES.md) | 4 | 4 | 8 |
| [MusicBrainz Integration](integrations/metadata/music/MUSICBRAINZ.md) | 2 | 6 | 8 |
| [TMDb (The Movie Database) Integration](integrations/metadata/video/TMDB.md) | 2 | 6 | 8 |
| [Trakt Integration](integrations/scrobbling/TRAKT.md) | 4 | 4 | 8 |
| [Chaptarr Integration](integrations/servarr/CHAPTARR.md) | 3 | 5 | 8 |
| [RBAC Service](services/RBAC.md) | 5 | 3 | 8 |
| [Revenge - User Experience Features](features/shared/USER_EXPERIENCE_FEATURES.md) | 2 | 5 | 7 |
| [Movie Module](features/video/MOVIE_MODULE.md) | 5 | 2 | 7 |
| [TV Show Module](features/video/TVSHOW_MODULE.md) | 6 | 1 | 7 |

---

## Orphan Documents

> Documents with no internal cross-references (may need linking)

- [Release Calendar System](features/playback/RELEASE_CALENDAR.md)
- [SyncPlay (Watch Together)](features/playback/SYNCPLAY.md)
- [Revenge - Internationalization (i18n)](features/shared/I18N.md)
- [Ticketing System](features/shared/TICKETING_SYSTEM.md)
- [ComicVine API Integration](integrations/metadata/comics/COMICVINE.md)
- [Grand Comics Database (GCD) Integration](integrations/metadata/comics/GRAND_COMICS_DATABASE.md)
- [Marvel API Integration](integrations/metadata/comics/MARVEL_API.md)
- [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md)
- [Branch Protection Rules](operations/BRANCH_PROTECTION.md)
- [GitFlow Workflow Guide](operations/GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md)
- [Versioning Strategy](operations/VERSIONING.md)
- [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md)
- [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md)
- [Revenge - Technology Stack](technical/TECH_STACK.md)


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Authelia Documentation](https://www.authelia.com/overview/) | [Local](../sources/security/authelia.md) |
| [Authentik Documentation](https://goauthentik.io/docs/) | [Local](../sources/security/authentik.md) |
| [Casbin](https://pkg.go.dev/github.com/casbin/casbin/v2) | [Local](../sources/security/casbin.md) |
| [Dragonfly Documentation](https://www.dragonflydb.io/docs) | [Local](../sources/infrastructure/dragonfly.md) |
| [Keycloak Documentation](https://www.keycloak.org/documentation) | [Local](../sources/security/keycloak.md) |
| [Last.fm API](https://www.last.fm/api/intro) | [Local](../sources/apis/lastfm.md) |
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../sources/tooling/river.md) |
| [Typesense API](https://typesense.org/docs/latest/api/) | [Local](../sources/infrastructure/typesense.md) |
| [Typesense Go Client](https://github.com/typesense/typesense-go) | [Local](../sources/infrastructure/typesense-go.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../sources/database/pgx.md) |

<!-- SOURCE-BREADCRUMBS-END -->

---

*Generated by `scripts/generate-design-crossref.py`*
