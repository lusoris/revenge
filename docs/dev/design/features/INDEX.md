# Features Documentation

> Content modules and feature specifications

---

## Overview

This directory contains feature specifications for all Revenge modules and capabilities.

---

## Core Features

### Playback & Streaming

| Feature | Document | Status |
|---------|----------|--------|
| Client Support | [CLIENT_SUPPORT.md](CLIENT_SUPPORT.md) | 🟡 Planned |
| Media Enhancements | [MEDIA_ENHANCEMENTS.md](MEDIA_ENHANCEMENTS.md) | 🟡 Planned |
| Watch Next | [WATCH_NEXT_CONTINUE_WATCHING.md](WATCH_NEXT_CONTINUE_WATCHING.md) | 🟡 Planned |

### Content Management

| Feature | Document | Status |
|---------|----------|--------|
| Library Types | [LIBRARY_TYPES.md](LIBRARY_TYPES.md) | 🟡 Planned |
| Comics Module | [COMICS_MODULE.md](COMICS_MODULE.md) | 🔴 0% |
| Content Rating | [CONTENT_RATING.md](CONTENT_RATING.md) | 🟡 Planned |
| Release Calendar | [RELEASE_CALENDAR.md](RELEASE_CALENDAR.md) | ✅ Servarr |

### User Features

| Feature | Document | Status |
|---------|----------|--------|
| Request System | [REQUEST_SYSTEM.md](REQUEST_SYSTEM.md) | ✅ Designed |
| RBAC (Casbin) | [RBAC_CASBIN.md](RBAC_CASBIN.md) | ✅ Code exists |
| User Experience | [USER_EXPERIENCE_FEATURES.md](USER_EXPERIENCE_FEATURES.md) | 🟡 Partial |
| Ticketing System | [TICKETING_SYSTEM.md](TICKETING_SYSTEM.md) | 🟡 Planned |

### Analytics & Monitoring

| Feature | Document | Status |
|---------|----------|--------|
| Analytics Service | [ANALYTICS_SERVICE.md](ANALYTICS_SERVICE.md) | 🟡 Designed |
| Feature Comparison | [FEATURE_COMPARISON.md](FEATURE_COMPARISON.md) | ✅ Reference |

### Communication

| Feature | Document | Status |
|---------|----------|--------|
| News System | [NEWS_SYSTEM.md](NEWS_SYSTEM.md) | 🟡 Designed |
| Wiki System | [WIKI_SYSTEM.md](WIKI_SYSTEM.md) | 🟡 Designed |

### Integration Features

| Feature | Document | Status |
|---------|----------|--------|
| Scrobbling | [SCROBBLING.md](SCROBBLING.md) | 🟡 Planned |
| Internationalization | [I18N.md](I18N.md) | 🟡 Planned |

---

## Adult Content (Isolated)

All adult features are isolated in `c` schema with separate API namespace `/api/v1/c/`.

| Feature | Document | Status |
|---------|----------|--------|
| Adult Content System | [ADULT_CONTENT_SYSTEM.md](ADULT_CONTENT_SYSTEM.md) | 🟡 Designed |
| Adult Metadata | [ADULT_METADATA.md](ADULT_METADATA.md) | 🟡 Designed |
| Data Reconciliation | [adult/DATA_RECONCILIATION.md](adult/DATA_RECONCILIATION.md) | 🟡 Designed |
| NSFW Toggle | [NSFW_TOGGLE.md](NSFW_TOGGLE.md) | 🟡 Designed |
| Whisparr/StashDB Schema | [WHISPARR_STASHDB_SCHEMA.md](WHISPARR_STASHDB_SCHEMA.md) | 🟡 Designed |

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
