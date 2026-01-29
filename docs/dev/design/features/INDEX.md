# Features Documentation

> Content modules and feature specifications

---

## Overview

This directory contains feature specifications organized by:
- **shared/** - Features that apply across all modules
- **playback/** - Playback and streaming features
- **Module folders** - Module-specific features (photos, podcasts, livetv, comics, adult)

---

## Shared Features

Features that apply to all or most modules.

| Feature | Document | Status |
|---------|----------|--------|
| Analytics | [shared/ANALYTICS_SERVICE.md](shared/ANALYTICS_SERVICE.md) | 🟡 Designed |
| RBAC | [shared/RBAC_CASBIN.md](shared/RBAC_CASBIN.md) | ✅ Code exists |
| Access Controls | [shared/ACCESS_CONTROLS.md](shared/ACCESS_CONTROLS.md) | 🟡 Designed |
| Request System | [shared/REQUEST_SYSTEM.md](shared/REQUEST_SYSTEM.md) | ✅ Designed |
| News System | [shared/NEWS_SYSTEM.md](shared/NEWS_SYSTEM.md) | 🟡 Designed |
| Wiki System | [shared/WIKI_SYSTEM.md](shared/WIKI_SYSTEM.md) | 🟡 Designed |
| Ticketing | [shared/TICKETING_SYSTEM.md](shared/TICKETING_SYSTEM.md) | 🟡 Planned |
| Client Support | [shared/CLIENT_SUPPORT.md](shared/CLIENT_SUPPORT.md) | 🟡 Planned |
| User Experience | [shared/USER_EXPERIENCE_FEATURES.md](shared/USER_EXPERIENCE_FEATURES.md) | 🟡 Partial |
| Voice Control | [shared/VOICE_CONTROL.md](shared/VOICE_CONTROL.md) | 🟡 Designed |
| Scrobbling | [shared/SCROBBLING.md](shared/SCROBBLING.md) | 🟡 Planned |
| i18n | [shared/I18N.md](shared/I18N.md) | 🟡 Planned |
| Library Types | [shared/LIBRARY_TYPES.md](shared/LIBRARY_TYPES.md) | 🟡 Planned |
| Content Rating | [shared/CONTENT_RATING.md](shared/CONTENT_RATING.md) | 🟡 Planned |
| NSFW Toggle | [shared/NSFW_TOGGLE.md](shared/NSFW_TOGGLE.md) | 🟡 Designed |

---

## Playback Features

Features related to media playback and streaming.

| Feature | Document | Status |
|---------|----------|--------|
| SyncPlay | [playback/SYNCPLAY.md](playback/SYNCPLAY.md) | 🟡 Designed |
| Trickplay | [playback/TRICKPLAY.md](playback/TRICKPLAY.md) | 🟡 Designed |
| Skip Intro/Credits | [playback/SKIP_INTRO.md](playback/SKIP_INTRO.md) | 🟡 Designed |
| Watch Next | [playback/WATCH_NEXT_CONTINUE_WATCHING.md](playback/WATCH_NEXT_CONTINUE_WATCHING.md) | 🟡 Planned |
| Release Calendar | [playback/RELEASE_CALENDAR.md](playback/RELEASE_CALENDAR.md) | ✅ Servarr |
| Media Enhancements | [playback/MEDIA_ENHANCEMENTS.md](playback/MEDIA_ENHANCEMENTS.md) | 🟡 Planned |

---

## Module-Specific Features

### Video (Movies & TV Shows)

| Feature | Document | Status |
|---------|----------|--------|
| Video Overview | [video/INDEX.md](video/INDEX.md) | ✅ Code exists |
| Movie Module | [video/MOVIE_MODULE.md](video/MOVIE_MODULE.md) | ✅ Code exists |
| TV Show Module | [video/TVSHOW_MODULE.md](video/TVSHOW_MODULE.md) | ✅ Code exists |

### Photos

| Feature | Document | Status |
|---------|----------|--------|
| Photos Library | [photos/PHOTOS_LIBRARY.md](photos/PHOTOS_LIBRARY.md) | 🟡 Designed |

### Podcasts

| Feature | Document | Status |
|---------|----------|--------|
| Podcasts | [podcasts/PODCASTS.md](podcasts/PODCASTS.md) | 🟡 Designed |

### Live TV

| Feature | Document | Status |
|---------|----------|--------|
| Live TV / DVR | [livetv/LIVE_TV_DVR.md](livetv/LIVE_TV_DVR.md) | 🟡 Designed |

### Comics

| Feature | Document | Status |
|---------|----------|--------|
| Comics Module | [comics/COMICS_MODULE.md](comics/COMICS_MODULE.md) | 🔴 0% |

### Adult (Isolated in `c` schema)

| Feature | Document | Status |
|---------|----------|--------|
| Adult Content System | [adult/ADULT_CONTENT_SYSTEM.md](adult/ADULT_CONTENT_SYSTEM.md) | 🟡 Designed |
| Adult Metadata | [adult/ADULT_METADATA.md](adult/ADULT_METADATA.md) | 🟡 Designed |
| Data Reconciliation | [adult/DATA_RECONCILIATION.md](adult/DATA_RECONCILIATION.md) | 🟡 Designed |
| Whisparr/StashDB Schema | [adult/WHISPARR_STASHDB_SCHEMA.md](adult/WHISPARR_STASHDB_SCHEMA.md) | 🟡 Designed |

---

## Feature Status Legend

| Status | Meaning |
|--------|---------|
| ✅ | Complete or code exists |
| 🟡 | Designed / Partial implementation |
| 🔴 | Not started |

---

## Related Documentation

- [Architecture](../architecture/) - System design
- [Integrations](../integrations/) - External services
- [Operations](../operations/) - Deployment guides
- [Technical](../technical/) - API and frontend docs
