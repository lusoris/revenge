# User Pain Points Research - Existing Media Servers

> Analysis of what users most complain about (and love) in Plex, Jellyfin, Emby



<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Research Sources](#research-sources)
- [Major Complaints (What Users Hate)](#major-complaints-what-users-hate)
  - [1. **Library Scanning Performance**](#1-library-scanning-performance)
  - [2. **Metadata Accuracy & Matching**](#2-metadata-accuracy-matching)
  - [3. **Transcoding Quality & Performance**](#3-transcoding-quality-performance)
  - [4. **Remote Access Complexity**](#4-remote-access-complexity)
  - [5. **Mobile App Quality**](#5-mobile-app-quality)
  - [6. **Subtitle Support**](#6-subtitle-support)
  - [7. **Collections & Organization**](#7-collections-organization)
  - [8. **Live TV & DVR**](#8-live-tv-dvr)
  - [9. **User Management & Permissions**](#9-user-management-permissions)
  - [10. **Resource Usage**](#10-resource-usage)
- [What Users LOVE (Keep These)](#what-users-love-keep-these)
  - [1. **Ease of Initial Setup** (Plex wins)](#1-ease-of-initial-setup-plex-wins)
  - [2. **Automatic Metadata** (Plex/Jellyfin)](#2-automatic-metadata-plexjellyfin)
  - [3. **Free & Open Source** (Jellyfin)](#3-free-open-source-jellyfin)
  - [4. **Hardware Transcoding** (Plex Pass, Emby Plus)](#4-hardware-transcoding-plex-pass-emby-plus)
  - [5. **Client Availability** (Plex)](#5-client-availability-plex)
  - [6. **Watch Together** (Plex)](#6-watch-together-plex)
  - [7. **Intro Skip** (Plex, Jellyfin with plugins)](#7-intro-skip-plex-jellyfin-with-plugins)
- [Competitive Positioning](#competitive-positioning)
- [Lessons for Revenge](#lessons-for-revenge)
  - [DO:](#do)
  - [DON'T:](#dont)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Next Steps](#next-steps)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🔴 |  |
| Sources | 🔴 |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |

---

## Research Sources

- **Reddit**: r/Plex, r/jellyfin, r/selfhosted
- **Forums**: Plex Forums, Jellyfin Forums, Emby Community
- **GitHub Issues**: Jellyfin/jellyfin, MediaBrowser/Emby (public issues)
- **Review Sites**: Reddit posts, TrustPilot, alternativeto.net

---

## Major Complaints (What Users Hate)

### 1. **Library Scanning Performance**

**Jellyfin 10.11.x Specific**:
- "Moved from 10.10.7 to 10.11.x, library scanning now taking forever" (Reddit)
- Scan times: 45 minutes → **8-14 hours** even with little changed
- Trickplay generation: Estimated **270 days** for large libraries (disabled by default)

**General Issues**:
- Slow initial scans on large libraries (100k+ items)
- No incremental scan (re-scans entire library vs changed files)
- High CPU/disk I/O during scans blocks other operations
- Metadata fetching bottlenecks (rate limits from TMDb, TheTVDB)

**Revenge Solution**:
- Incremental scans only (detect changed files via filesystem events)
- Parallel metadata fetching with rate limiting per provider
- Background priority for scans (don't block playback)
- Dedicated scanning service (separate from main server)

### 2. **Metadata Accuracy & Matching**

**Common Complaints**:
- "Wrong posters/fanart for movies" (especially for foreign films)
- "TV episodes mismatched to wrong season"
- "Anime matching is terrible" (English titles vs Japanese)
- "Can't manually override incorrect matches"
- "Metadata refreshes overwrite manual edits"

**Revenge Solution**:
- Multi-provider fallback (TMDb → TheTVDB → OMDb → manual)
- Anime-specific matching (AniList, MyAnimeList, Kitsu)
- Manual overrides persist across refreshes
- Metadata moderation queue (admin review conflicts)
- User-reported corrections (ticketing system)

### 3. **Transcoding Quality & Performance**

**Plex Issues**:
- "Plex transcoding uses 100% CPU for 1080p → 720p"
- "Hardware transcoding locked behind Plex Pass ($)"
- "Transcoder crashes on 4K HDR content"
- "Subtitles burn-in is slow/buggy"

**Jellyfin Issues**:
- "Transcoding randomly stops mid-playback"
- "HDR tone mapping is broken (washed out colors)"
- "Audio transcoding artifacts (crackling, sync issues)"

**Revenge Solution**:
- **External transcoding** (Blackbeard service, no CPU load on Revenge)
- Hardware transcoding free (Blackbeard handles GPU)
- Client capability detection (avoid unnecessary transcodes)
- Bandwidth-aware quality selection (external clients)
- HLS/DASH segment buffering (smooth playback during issues)

### 4. **Remote Access Complexity**

**Plex Issues**:
- "Plex relay is slow (limited to 1 Mbps free tier)"
- "Remote access randomly stops working"
- "Port forwarding confusing for non-technical users"

**Jellyfin/Emby Issues**:
- "No built-in remote access (must self-configure)"
- "Reverse proxy setup is complex (Nginx, Caddy)"
- "HTTPS certificates are manual (Let's Encrypt renewal)"

**Revenge Solution**:
- Built-in reverse proxy (Traefik integration)
- Automatic HTTPS (Let's Encrypt auto-renewal)
- Cloudflare Tunnel option (no port forwarding)
- OIDC/SSO for secure external access
- Optional relay service (self-hosted, no fees)

### 5. **Mobile App Quality**

**Plex**:
- "iOS app crashes on 4K content"
- "Offline sync is unreliable"
- "Downloads fail silently"

**Jellyfin**:
- "Mobile apps are buggy (crashes, playback stops)"
- "No Picture-in-Picture (PiP) on iOS"
- "UI is clunky (not native feel)"

**Revenge Solution**:
- Native mobile apps (Swift iOS, Kotlin Android) vs webviews
- Offline sync with conflict resolution
- Background download with retry logic
- PiP support (iOS, Android)

### 6. **Subtitle Support**

**Common Issues**:
- "External SRT files not detected"
- "Subtitle sync is off (timing issues)"
- "No support for PGS/VobSub in browser"
- "Embedded subtitle extraction fails"

**Revenge Solution**:
- Auto-detect external subtitles (multiple naming conventions)
- Subtitle offset adjustment UI
- WebVTT conversion for browser playback
- Embedded subtitle extraction (FFmpeg)

### 7. **Collections & Organization**

**Plex**:
- "Collections are manual (no auto-grouping)"
- "Can't mix movies + shows in collections"
- "No nested collections"

**Jellyfin**:
- "Collections don't sync metadata (posters, etc.)"
- "Collection visibility settings are confusing"

**Revenge Solution**:
- Auto-collections (franchises, directors, actors)
- Cross-module collections (video pool: movies + episodes)
- Nested collections (Marvel → MCU → Phase 1)
- Per-user collection visibility

### 8. **Live TV & DVR**

**Plex**:
- "Live TV requires Plex Pass ($)"
- "DVR recordings are buggy (cut off endings)"
- "EPG data is missing/incorrect"

**Jellyfin**:
- "Live TV setup is complex (HDHomeRun, TVHeadend)"
- "EPG mapping is manual"
- "Recording conflicts not handled"

**Revenge Solution**:
- Free Live TV (no paywall)
- Auto-EPG mapping (fuzzy match)
- Smart recording conflict resolution
- Post-recording commercial skip (Blackbeard)

### 9. **User Management & Permissions**

**Plex**:
- "Home users share watch history (no privacy)"
- "Can't restrict specific libraries per user"
- "Parental controls are weak"

**Jellyfin/Emby**:
- "No SSO/OIDC support (local accounts only)"
- "Permission granularity is limited"

**Revenge Solution**:
- OIDC/SSO integration (Authelia, Authentik, Keycloak)
- Per-library permissions (read, write, admin)
- Profile-based restrictions (Netflix-style)
- NSFW mode toggle (hide adult content)

### 10. **Resource Usage**

**Plex**:
- "Plex uses 2GB RAM idle (excessive)"
- "Database grows to 10GB+ (slow queries)"

**Jellyfin**:
- "Memory leaks after 7 days uptime"
- "CPU spikes during library scans"

**Revenge Solution**:
- Efficient PostgreSQL (indexed queries, partitioning)
- Memory-aware transcode cache (eviction strategies)
- Background job prioritization (River queue)
- Idle resource release (graceful shutdown)

---

## What Users LOVE (Keep These)

### 1. **Ease of Initial Setup** (Plex wins)
- "Plex setup is 5 minutes (wizard, auto-detect)"
- Revenge: Simple installer, auto-configuration, sane defaults

### 2. **Automatic Metadata** (Plex/Jellyfin)
- "Just drop files, server handles rest"
- Revenge: Multi-provider fetching, background jobs

### 3. **Free & Open Source** (Jellyfin)
- "No paywalls, community-driven"
- Revenge: 100% open source, no premium tiers

### 4. **Hardware Transcoding** (Plex Pass, Emby Plus)
- "GPU transcoding is fast"
- Revenge: Free via Blackbeard service

### 5. **Client Availability** (Plex)
- "Works on everything (TV, phone, browser)"
- Revenge: Web (primary), mobile apps (native), Kodi plugin

### 6. **Watch Together** (Plex)
- "Sync playback with friends"
- Revenge: WebSocket sync for Watch Party

### 7. **Intro Skip** (Plex, Jellyfin with plugins)
- "Auto-skip intros/credits"
- Revenge: ML-based intro detection (Blackbeard)

---

## Competitive Positioning

| Feature | Plex | Jellyfin | Emby | **Revenge** |
|---------|------|----------|------|-------------|
| **Performance** | ⚠️ Slow scans | ⚠️ Very slow (10.11.x) | ✅ Fast | ✅ **Incremental, parallel** |
| **Metadata** | ⚠️ Inaccurate | ⚠️ Manual fixes lost | ⚠️ Limited providers | ✅ **Multi-provider + moderation** |
| **Transcoding** | ⚠️ CPU-heavy | ⚠️ Buggy HDR | ⚠️ Paid feature | ✅ **External (Blackbeard)** |
| **Remote Access** | ⚠️ Relay ($) | ❌ Manual | ⚠️ Complex | ✅ **Built-in reverse proxy** |
| **Mobile Apps** | ⚠️ Crashes | ❌ Buggy | ✅ Good | ✅ **Native (planned)** |
| **Subtitles** | ⚠️ Limited | ⚠️ Buggy | ✅ Good | ✅ **Auto-detect + sync** |
| **Collections** | ⚠️ Manual | ⚠️ Limited | ⚠️ Manual | ✅ **Auto + cross-module** |
| **Live TV** | 💰 Plex Pass | ⚠️ Complex | 💰 Emby Plus | ✅ **Free + smart DVR** |
| **OIDC/SSO** | ❌ No | ❌ No | ❌ No | ✅ **Yes** |
| **Open Source** | ❌ Proprietary | ✅ Yes | ⚠️ Core only | ✅ **Yes** |

---

## Lessons for Revenge

### DO:
1. ✅ **Incremental library scanning** (don't re-scan entire library)
2. ✅ **Multi-provider metadata** (TMDb, TheTVDB, OMDb, AniList, etc.)
3. ✅ **External transcoding** (avoid Plex/Jellyfin CPU bottleneck)
4. ✅ **Built-in reverse proxy** (Traefik, automatic HTTPS)
5. ✅ **OIDC/SSO** (enterprise-friendly)
6. ✅ **Free hardware transcoding** (Blackbeard GPU support)
7. ✅ **NSFW mode toggle** (hide adult content)
8. ✅ **Metadata moderation queue** (admin review conflicts)
9. ✅ **Auto-collections** (franchises, directors, actors)
10. ✅ **Watch Party** (WebSocket sync)

### DON'T:
1. ❌ **Slow library scans** (Jellyfin 10.11.x mistake)
2. ❌ **CPU-heavy transcoding** (Plex bottleneck)
3. ❌ **Paid features** (keep 100% free)
4. ❌ **Manual remote access setup** (Jellyfin complexity)
5. ❌ **Buggy mobile apps** (Jellyfin WebView issue)
6. ❌ **Overwrite manual metadata** (Plex/Jellyfin problem)
7. ❌ **Complex permissions** (Emby granularity overload)


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Authelia Documentation](https://www.authelia.com/overview/) | [Local](../../sources/security/authelia.md) |
| [Authentik Documentation](https://goauthentik.io/docs/) | [Local](../../sources/security/authentik.md) |
| [FFmpeg Codecs](https://ffmpeg.org/ffmpeg-codecs.html) | [Local](../../sources/media/ffmpeg-codecs.md) |
| [FFmpeg Documentation](https://ffmpeg.org/ffmpeg.html) | [Local](../../sources/media/ffmpeg.md) |
| [FFmpeg Formats](https://ffmpeg.org/ffmpeg-formats.html) | [Local](../../sources/media/ffmpeg-formats.md) |
| [Keycloak Documentation](https://www.keycloak.org/documentation) | [Local](../../sources/security/keycloak.md) |
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../sources/tooling/river.md) |
| [go-astiav (FFmpeg bindings)](https://pkg.go.dev/github.com/asticode/go-astiav) | [Local](../../sources/media/go-astiav.md) |
| [gohlslib (HLS)](https://pkg.go.dev/github.com/bluenviron/gohlslib/v2) | [Local](../../sources/media/gohlslib.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../sources/database/pgx.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Research](INDEX.md)

### In This Section

- [UX/UI Design & Frontend Resources](UX_UI_RESOURCES.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

---

## Next Steps

1. **Incremental Scanning**: Filesystem watcher + changed files only
2. **Metadata Moderation**: Queue for conflicts (admin review)
3. **Blackbeard Reliability**: Ensure external transcoding is rock-solid
4. **Reverse Proxy**: Traefik integration with auto-HTTPS
5. **Mobile Apps**: Native Swift/Kotlin (not WebView wrappers)
6. **User Research**: Continuous Reddit/forum monitoring

