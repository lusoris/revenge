# Planning vs. Reality Analysis

**Date**: 2026-02-02
**Purpose**: Compare current TODO planning with actual design work completed

---

## Summary of Design Work Completed

We have **159 YAML design documents** covering:

### ✅ What We HAVE Designed (Beyond v0.3.0 MVP scope):

**Services** (19 designed):
- ✅ Activity, Analytics, API Keys, Auth, EPG, Fingerprint, Grants
- ✅ HTTP Client & Proxy, Library, Metadata, Notification, OIDC
- ✅ RBAC, Search, Session, Settings, Transcoding, User

**Content Modules** (11 designed):
- ✅ Movies (v0.3.0 scope)
- ✅ TV Shows (v0.4.0 scope)
- ✅ Music (v0.5.0+ scope)
- ✅ Audiobooks (v0.5.0+ scope)
- ✅ Books (v0.5.0+ scope)
- ✅ Comics (v0.5.0+ scope)
- ✅ Podcasts (v0.5.0+ scope)
- ✅ Photos (v0.5.0+ scope)
- ✅ Live TV & DVR (v0.5.0+ scope)
- ✅ Adult Content System (v0.7.0+ scope)
- ✅ Adult Gallery Module (v0.7.0+ scope)

**Playback Features** (6 designed):
- ✅ Trickplay (Timeline Thumbnails)
- ✅ Skip Intro/Credits Detection
- ✅ SyncPlay (Watch Together)
- ✅ Media Enhancements
- ✅ Watch Next & Continue Watching
- ✅ Release Calendar

**Shared Features** (17 designed):
- ✅ Collections & Playlists
- ✅ Content Rating System
- ✅ RBAC with Casbin
- ✅ Client Support & Device Capabilities
- ✅ User Experience Features
- ✅ Library Types
- ✅ NSFW Toggle
- ✅ Time-Based Access Controls
- ✅ Request System
- ✅ News System
- ✅ Ticketing System
- ✅ Wiki System
- ✅ Voice Control
- ✅ External Scrobbling & Sync
- ✅ Analytics Service (Tracearr)
- ✅ Internationalization (i18n)

**Integrations** (58 designed):
- Infrastructure: PostgreSQL, Dragonfly, River, Typesense
- Auth: Authelia, Authentik, Keycloak, Generic OIDC
- Servarr: Radarr, Sonarr, Lidarr, Chaptarr, Whisparr
- Metadata Providers: TMDb, TheTVDB, OMDb, ThePosterDB
- Anime: AniList, Kitsu, MyAnimeList
- Music: Spotify, Last.fm, MusicBrainz, Discogs, ListenBrainz
- Books: Audible, Goodreads, Hardcover, OpenLibrary
- Comics: ComicVine, Grand Comics Database, Marvel API
- Adult: StashDB, ThePornDB, FreeOnes, Stash, Pornhub, OnlyFans, Instagram, Twitter/X, TheNude
- Adult Wikis: IAFD, Babepedia, Boobpedia
- Scrobbling: Trakt, Simkl, Letterboxd, Last.fm, ListenBrainz
- Wiki: Wikipedia, FANDOM, TVTropes
- Live TV: TVHeadend, NextPVR, ErsatzTV
- Casting: Chromecast, DLNA/UPnP
- Transcoding: Blackbeard

**Architecture & Technical** (27 designed):
- Architecture: System Architecture v2, Design Principles, Metadata System, Player Architecture, Plugin Architecture
- Patterns: Arr Integration, Metadata Enrichment, Observability, Testing, Webhooks, HTTP Client/Proxy
- Technical: API, Configuration, Email, Webhooks, WebSockets, Frontend, Observability, Offloading, Tech Stack, Testing, Audio Streaming, Notifications
- Design System: Brand Identity, Color System, Typography, Components, Layout, Motion, Navigation, Accessibility, Pirate Mode
- Operations: Git Workflow, Branch Protection, Versioning, Development Setup, Production Setup, Reverse Proxy, Best Practices, Database Auto-Healing
- Research: User Pain Points

---

## Issues Found in Current Planning

### 🔴 CRITICAL: Scope Creep in TODO Planning

The current TODO files (v0.1.0 - v1.0.0) are **significantly smaller in scope** than what we've actually designed:

#### v0.3.0 MVP (Movies)
**TODO says**: Basic movie module + TMDb + Radarr
**Reality**: We've designed ALL content modules, not just movies

#### v0.4.0 (TV Shows)
**TODO says**: Just TV shows
**Reality**: TV shows design exists, but so do 9 OTHER content modules

#### v0.5.0+
**TODO says**: Music, Photos, Live TV, Audiobooks (spread across multiple versions)
**Reality**: ALL of these are already fully designed

#### v0.7.0+ (Adult Content)
**TODO says**: Adult content as late-stage feature
**Reality**: Fully designed adult system with 11 integrations

### 🟡 MODERATE: Missing from Planning

**What's Designed but NOT in TODOs**:

1. **Playback Features** (6 features):
   - Trickplay
   - Skip Intro/Credits
   - SyncPlay
   - Media Enhancements
   - Watch Next/Continue Watching
   - Release Calendar

2. **Shared Features** (many missing):
   - Collections (mentioned, but underspecified)
   - Content Rating System
   - Client Support & Device Capabilities
   - NSFW Toggle
   - Time-Based Access Controls
   - Request System (native, not just Overseerr)
   - News System
   - Ticketing System
   - Wiki System
   - Voice Control
   - Analytics Service (Tracearr)

3. **Integrations** (many underspecified):
   - TODOs mention some integrations but don't reflect the **58 integration designs**
   - Missing: Most anime providers, adult providers, scrobbling services, wikis

4. **Services**:
   - EPG Service (for Live TV)
   - Fingerprint Service
   - Grants Service
   - HTTP Client & Proxy Service (critical pattern)
   - Transcoding Service (mentioned but not scoped)

5. **Architecture**:
   - Plugin Architecture Decision (critical for extensibility)
   - Offloading Architecture
   - HTTP Client Pattern (proxy/VPN routing)

---

## Recommendations

### Option 1: Keep Current Phasing, Add References
- Keep v0.3.0 as MVP (Movies only)
- Add notes in each TODO referencing existing design docs
- Make it clear: "Design complete, implementation pending"

### Option 2: Restructure Milestones to Reflect Reality
- v0.3.0: MVP (Movies + Core Playback Features)
- v0.4.0: TV Shows + Enhanced Playback
- v0.5.0: Multi-Content Support (Music, Books, Comics, Audiobooks)
- v0.6.0: Live TV, Photos, Podcasts
- v0.7.0: Adult Content
- v0.8.0: Advanced Features (Collections, Request System, etc.)

### Option 3: Two-Phase Approach (RECOMMENDED)
- **Phase 1: Design** ✅ COMPLETE (159 docs)
- **Phase 2: Implementation** (follow current TODO phasing)
- Add "Design References" section to each TODO linking to relevant YAML docs

---

## What Needs to Happen

1. ✅ Acknowledge that design work is WAY ahead of planning
2. ✅ Add "Design Documentation" sections to each TODO file
3. ✅ Link TODOs to their corresponding YAML design docs
4. ✅ Make it clear: design ≠ implementation timeline
5. ✅ Keep implementation phasing realistic (MVP-first approach)

---

## Next Steps

**Immediate**:
- [ ] Update TODO_v0.3.0.md with "Design References" section
- [ ] Update TODO_v0.4.0.md with "Design References" section
- [ ] Update TODO_v0.5.0.md+ with content module references
- [ ] Update ROADMAP.md to clarify: Design phase complete, now in Implementation phase

**Optional** (for clarity):
- [ ] Create DESIGN_INVENTORY.md listing all 159 design docs by category
- [ ] Create IMPLEMENTATION_BACKLOG.md separating designed features from planned features
