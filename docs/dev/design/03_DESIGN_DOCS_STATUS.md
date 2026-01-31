# Design Documentation Status

<!-- SOURCES: authelia, authentik, casbin, dragonfly, keycloak, lastfm-api, pgx, postgresql-arrays, postgresql-json, river, typesense, typesense-go -->

> Auto-generated overview of design document completeness


<!-- TOC-START -->

## Table of Contents

- [Status Legend](#status-legend)
- [Overall Summary](#overall-summary)
- [Architecture](#architecture)
- [Features - Adult](#features---adult)
- [Features - Comics](#features---comics)
- [Features - Livetv](#features---livetv)
- [Features - Photos](#features---photos)
- [Features - Playback](#features---playback)
- [Features - Podcasts](#features---podcasts)
- [Features - Shared](#features---shared)
- [Features - Video](#features---video)
- [Integrations - Anime](#integrations---anime)
- [Integrations - Auth](#integrations---auth)
- [Integrations - Casting](#integrations---casting)
- [Integrations - External](#integrations---external)
- [Integrations - Infrastructure](#integrations---infrastructure)
- [Integrations - Livetv](#integrations---livetv)
- [Integrations - Metadata](#integrations---metadata)
- [Integrations - Scrobbling](#integrations---scrobbling)
- [Integrations - Servarr](#integrations---servarr)
- [Integrations - Transcoding](#integrations---transcoding)
- [Integrations - Wiki](#integrations---wiki)
- [Operations](#operations)
- [Planning](#planning)
- [Research](#research)
- [Services](#services)
- [Technical](#technical)
- [Notes](#notes)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Regenerate](#regenerate)

<!-- TOC-END -->

**Last Updated**: Auto-generated

---

## Status Legend

| Emoji | Meaning |
|-------|---------|
| ✅ | Complete |
| 🟡 | Partial |
| 🔴 | Not Started |

---

## Overall Summary

**Total Documents**: 126

| Category | Total | Design ✅ | Sources ✅ | Instructions ✅ |
|----------|-------|-----------|------------|-----------------|
| Architecture | 5 | 5 (100%) | 1 (20%) | 3 (60%) |
| Features - Adult | 5 | 5 (100%) | 2 (40%) | 4 (80%) |
| Features - Comics | 1 | 1 (100%) | 1 (100%) | 1 (100%) |
| Features - Livetv | 1 | 1 (100%) | 1 (100%) | 1 (100%) |
| Features - Photos | 1 | 1 (100%) | 0 (0%) | 1 (100%) |
| Features - Playback | 6 | 6 (100%) | 6 (100%) | 2 (33%) |
| Features - Podcasts | 1 | 1 (100%) | 1 (100%) | 1 (100%) |
| Features - Shared | 15 | 15 (100%) | 1 (6%) | 15 (100%) |
| Features - Video | 2 | 2 (100%) | 2 (100%) | 2 (100%) |
| Integrations - Anime | 3 | 3 (100%) | 0 (0%) | 3 (100%) |
| Integrations - Auth | 4 | 2 (50%) | 0 (0%) | 4 (100%) |
| Integrations - Casting | 2 | 2 (100%) | 0 (0%) | 2 (100%) |
| Integrations - External | 6 | 4 (66%) | 0 (0%) | 6 (100%) |
| Integrations - Infrastructure | 4 | 4 (100%) | 3 (75%) | 3 (75%) |
| Integrations - Livetv | 3 | 3 (100%) | 0 (0%) | 3 (100%) |
| Integrations - Metadata | 19 | 16 (84%) | 1 (5%) | 18 (94%) |
| Integrations - Scrobbling | 5 | 1 (20%) | 0 (0%) | 5 (100%) |
| Integrations - Servarr | 5 | 5 (100%) | 1 (20%) | 5 (100%) |
| Integrations - Transcoding | 1 | 1 (100%) | 0 (0%) | 1 (100%) |
| Integrations - Wiki | 6 | 0 (0%) | 0 (0%) | 6 (100%) |
| Operations | 7 | 3 (42%) | 4 (57%) | 0 (0%) |
| Planning | 1 | 1 (100%) | 1 (100%) | 0 (0%) |
| Research | 2 | 0 (0%) | 0 (0%) | 0 (0%) |
| Services | 15 | 15 (100%) | 0 (0%) | 14 (93%) |
| Technical | 6 | 4 (66%) | 1 (16%) | 0 (0%) |
| **TOTAL** | **126** | **101 (80%)** | **26 (20%)** | **100 (79%)** |

---

## Architecture

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md) | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [04_PLAYER_ARCHITECTURE](architecture/04_PLAYER_ARCHITECTURE.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [05_PLUGIN_ARCHITECTURE_DECISION](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) | ✅ | ✅ | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 5/5 Design ✅ | 1/5 Sources ✅ | 3/5 Instructions ✅

---

## Features - Adult

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [ADULT_CONTENT_SYSTEM](features/adult/ADULT_CONTENT_SYSTEM.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [ADULT_METADATA](features/adult/ADULT_METADATA.md) | ✅ | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [DATA_RECONCILIATION](features/adult/DATA_RECONCILIATION.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [GALLERY_MODULE](features/adult/GALLERY_MODULE.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [WHISPARR_STASHDB_SCHEMA](features/adult/WHISPARR_STASHDB_SCHEMA.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 5/5 Design ✅ | 2/5 Sources ✅ | 4/5 Instructions ✅

---

## Features - Comics

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [COMICS_MODULE](features/comics/COMICS_MODULE.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 1/1 Design ✅ | 1/1 Sources ✅ | 1/1 Instructions ✅

---

## Features - Livetv

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [LIVE_TV_DVR](features/livetv/LIVE_TV_DVR.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 1/1 Design ✅ | 1/1 Sources ✅ | 1/1 Instructions ✅

---

## Features - Photos

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [PHOTOS_LIBRARY](features/photos/PHOTOS_LIBRARY.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 1/1 Design ✅ | 0/1 Sources ✅ | 1/1 Instructions ✅

---

## Features - Playback

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [MEDIA_ENHANCEMENTS](features/playback/MEDIA_ENHANCEMENTS.md) | ✅ | ✅ | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 |
| [RELEASE_CALENDAR](features/playback/RELEASE_CALENDAR.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [SKIP_INTRO](features/playback/SKIP_INTRO.md) | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [SYNCPLAY](features/playback/SYNCPLAY.md) | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [TRICKPLAY](features/playback/TRICKPLAY.md) | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [WATCH_NEXT_CONTINUE_WATCHING](features/playback/WATCH_NEXT_CONTINUE_WATCHING.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 6/6 Design ✅ | 6/6 Sources ✅ | 2/6 Instructions ✅

---

## Features - Podcasts

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [PODCASTS](features/podcasts/PODCASTS.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 1/1 Design ✅ | 1/1 Sources ✅ | 1/1 Instructions ✅

---

## Features - Shared

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [ACCESS_CONTROLS](features/shared/ACCESS_CONTROLS.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [ANALYTICS_SERVICE](features/shared/ANALYTICS_SERVICE.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [CLIENT_SUPPORT](features/shared/CLIENT_SUPPORT.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [CONTENT_RATING](features/shared/CONTENT_RATING.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [I18N](features/shared/I18N.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [LIBRARY_TYPES](features/shared/LIBRARY_TYPES.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [NEWS_SYSTEM](features/shared/NEWS_SYSTEM.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [NSFW_TOGGLE](features/shared/NSFW_TOGGLE.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [RBAC_CASBIN](features/shared/RBAC_CASBIN.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [REQUEST_SYSTEM](features/shared/REQUEST_SYSTEM.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [SCROBBLING](features/shared/SCROBBLING.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [TICKETING_SYSTEM](features/shared/TICKETING_SYSTEM.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [USER_EXPERIENCE_FEATURES](features/shared/USER_EXPERIENCE_FEATURES.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [VOICE_CONTROL](features/shared/VOICE_CONTROL.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [WIKI_SYSTEM](features/shared/WIKI_SYSTEM.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 15/15 Design ✅ | 1/15 Sources ✅ | 15/15 Instructions ✅

---

## Features - Video

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [MOVIE_MODULE](features/video/MOVIE_MODULE.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [TVSHOW_MODULE](features/video/TVSHOW_MODULE.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 2/2 Design ✅ | 2/2 Sources ✅ | 2/2 Instructions ✅

---

## Integrations - Anime

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [ANILIST](integrations/anime/ANILIST.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [KITSU](integrations/anime/KITSU.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [MYANIMELIST](integrations/anime/MYANIMELIST.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 3/3 Design ✅ | 0/3 Sources ✅ | 3/3 Instructions ✅

---

## Integrations - Auth

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [AUTHELIA](integrations/auth/AUTHELIA.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [AUTHENTIK](integrations/auth/AUTHENTIK.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [GENERIC_OIDC](integrations/auth/GENERIC_OIDC.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [KEYCLOAK](integrations/auth/KEYCLOAK.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 2/4 Design ✅ | 0/4 Sources ✅ | 4/4 Instructions ✅

---

## Integrations - Casting

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [CHROMECAST](integrations/casting/CHROMECAST.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [DLNA](integrations/casting/DLNA.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 2/2 Design ✅ | 0/2 Sources ✅ | 2/2 Instructions ✅

---

## Integrations - External

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [FREEONES](integrations/metadata/adult/FREEONES.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [INSTAGRAM](integrations/metadata/adult/INSTAGRAM.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [ONLYFANS](integrations/metadata/adult/ONLYFANS.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [PORNHUB](integrations/metadata/adult/PORNHUB.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [THENUDE](integrations/metadata/adult/THENUDE.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [TWITTER_X](integrations/metadata/adult/TWITTER_X.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 4/6 Design ✅ | 0/6 Sources ✅ | 6/6 Instructions ✅

---

## Integrations - Infrastructure

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [DRAGONFLY](integrations/infrastructure/DRAGONFLY.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [POSTGRESQL](integrations/infrastructure/POSTGRESQL.md) | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [RIVER](integrations/infrastructure/RIVER.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [TYPESENSE](integrations/infrastructure/TYPESENSE.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 4/4 Design ✅ | 3/4 Sources ✅ | 3/4 Instructions ✅

---

## Integrations - Livetv

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [ERSATZTV](integrations/livetv/ERSATZTV.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [NEXTPVR](integrations/livetv/NEXTPVR.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [TVHEADEND](integrations/livetv/TVHEADEND.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 3/3 Design ✅ | 0/3 Sources ✅ | 3/3 Instructions ✅

---

## Integrations - Metadata

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [AUDIBLE](integrations/metadata/books/AUDIBLE.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [COMICVINE](integrations/metadata/comics/COMICVINE.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [DISCOGS](integrations/metadata/music/DISCOGS.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [GOODREADS](integrations/metadata/books/GOODREADS.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [GRAND_COMICS_DATABASE](integrations/metadata/comics/GRAND_COMICS_DATABASE.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [HARDCOVER](integrations/metadata/books/HARDCOVER.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [LASTFM](integrations/metadata/music/LASTFM.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [MARVEL_API](integrations/metadata/comics/MARVEL_API.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [MUSICBRAINZ](integrations/metadata/music/MUSICBRAINZ.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [OMDB](integrations/metadata/video/OMDB.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [OPENLIBRARY](integrations/metadata/books/OPENLIBRARY.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [SPOTIFY](integrations/metadata/music/SPOTIFY.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [STASH](integrations/metadata/adult/STASH.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [STASHDB](integrations/metadata/adult/STASHDB.md) | 🟡 | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [THEPORNDB](integrations/metadata/adult/THEPORNDB.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [THEPOSTERDB](integrations/metadata/video/THEPOSTERDB.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [THETVDB](integrations/metadata/video/THETVDB.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [TMDB](integrations/metadata/video/TMDB.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [WHISPARR_V3_ANALYSIS](integrations/metadata/adult/WHISPARR_V3_ANALYSIS.md) | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 16/19 Design ✅ | 1/19 Sources ✅ | 18/19 Instructions ✅

---

## Integrations - Scrobbling

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [LASTFM_SCROBBLE](integrations/scrobbling/LASTFM_SCROBBLE.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [LETTERBOXD](integrations/scrobbling/LETTERBOXD.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [LISTENBRAINZ](integrations/scrobbling/LISTENBRAINZ.md) | 🟡 | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [SIMKL](integrations/scrobbling/SIMKL.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [TRAKT](integrations/scrobbling/TRAKT.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 1/5 Design ✅ | 0/5 Sources ✅ | 5/5 Instructions ✅

---

## Integrations - Servarr

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [CHAPTARR](integrations/servarr/CHAPTARR.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [LIDARR](integrations/servarr/LIDARR.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [RADARR](integrations/servarr/RADARR.md) | ✅ | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [SONARR](integrations/servarr/SONARR.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [WHISPARR](integrations/servarr/WHISPARR.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 5/5 Design ✅ | 1/5 Sources ✅ | 5/5 Instructions ✅

---

## Integrations - Transcoding

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [BLACKBEARD](integrations/transcoding/BLACKBEARD.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 1/1 Design ✅ | 0/1 Sources ✅ | 1/1 Instructions ✅

---

## Integrations - Wiki

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [BABEPEDIA](integrations/wiki/adult/BABEPEDIA.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [BOOBPEDIA](integrations/wiki/adult/BOOBPEDIA.md) | 🟡 | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [FANDOM](integrations/wiki/FANDOM.md) | 🟡 | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [IAFD](integrations/wiki/adult/IAFD.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [TVTROPES](integrations/wiki/TVTROPES.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [WIKIPEDIA](integrations/wiki/WIKIPEDIA.md) | 🟡 | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 0/6 Design ✅ | 0/6 Sources ✅ | 6/6 Instructions ✅

---

## Operations

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [BEST_PRACTICES](operations/BEST_PRACTICES.md) | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [BRANCH_PROTECTION](operations/BRANCH_PROTECTION.md) | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [DATABASE_AUTO_HEALING](operations/DATABASE_AUTO_HEALING.md) | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [DEVELOPMENT](operations/DEVELOPMENT.md) | 🟡 | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [GITFLOW](operations/GITFLOW.md) | 🟡 | ✅ | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 |
| [REVERSE_PROXY](operations/REVERSE_PROXY.md) | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [SETUP](operations/SETUP.md) | 🟡 | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 3/7 Design ✅ | 4/7 Sources ✅ | 0/7 Instructions ✅

---

## Planning

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [VERSIONING](operations/VERSIONING.md) | ✅ | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 1/1 Design ✅ | 1/1 Sources ✅ | 0/1 Instructions ✅

---

## Research

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [USER_PAIN_POINTS_RESEARCH](research/USER_PAIN_POINTS_RESEARCH.md) | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [UX_UI_RESOURCES](research/UX_UI_RESOURCES.md) | 🟡 | 🔴 | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 0/2 Design ✅ | 0/2 Sources ✅ | 0/2 Instructions ✅

---

## Services

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [ACTIVITY](services/ACTIVITY.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [ANALYTICS](services/ANALYTICS.md) | ✅ | 🔴 | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 |
| [APIKEYS](services/APIKEYS.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [AUTH](services/AUTH.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [FINGERPRINT](services/FINGERPRINT.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [GRANTS](services/GRANTS.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [LIBRARY](services/LIBRARY.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [METADATA](services/METADATA.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [NOTIFICATION](services/NOTIFICATION.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [OIDC](services/OIDC.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [RBAC](services/RBAC.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [SEARCH](services/SEARCH.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [SESSION](services/SESSION.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [SETTINGS](services/SETTINGS.md) | ✅ | 🔴 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |
| [USER](services/USER.md) | ✅ | 🟡 | ✅ | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 15/15 Design ✅ | 0/15 Sources ✅ | 14/15 Instructions ✅

---

## Technical

| Document | Design | Sources | Instructions | Code | Linting | Unit | Integration |
|----------|--------|---------|--------------|------|---------|------|-------------|
| [API](technical/API.md) | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [AUDIO_STREAMING](technical/AUDIO_STREAMING.md) | ✅ | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [CONFIGURATION](technical/CONFIGURATION.md) | ✅ | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [FRONTEND](technical/FRONTEND.md) | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 | 🔴 |
| [OFFLOADING](technical/OFFLOADING.md) | ✅ | 🔴 | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 |
| [TECH_STACK](technical/TECH_STACK.md) | ✅ | ✅ | 🟡 | 🔴 | 🔴 | 🔴 | 🔴 |

**Summary**: 4/6 Design ✅ | 1/6 Sources ✅ | 0/6 Instructions ✅

---

## Notes

- **Code/Linting/Unit/Integration**: All 🔴 as codebase is at template stage
- **Design**: Schemas, tables, architecture diagrams, Go code examples
- **Sources**: Developer Resources section with external documentation links
- **Instructions**: Implementation Checklist with actionable items


---

## Regenerate

```bash
python scripts/audit-design-status.py --update
```