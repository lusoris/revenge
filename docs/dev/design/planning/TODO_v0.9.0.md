# TODO v0.9.0 - RC1 (Release Candidate)

<!-- DESIGN: planning, README, test_output_claude, test_output_wiki -->


<!-- TOC-START -->

## Table of Contents

- [Overview](#overview)
- [Deliverables](#deliverables)
  - [QAR Module (Adult Content)](#qar-module-adult-content)
  - [StashDB Integration](#stashdb-integration)
  - [Whisparr Integration](#whisparr-integration)
  - [Live TV/DVR](#live-tvdvr)
  - [Photos Module](#photos-module)
  - [Comics Module](#comics-module)
  - [Performance Optimization](#performance-optimization)
  - [Bug Fixes](#bug-fixes)
  - [Documentation](#documentation)
  - [Security Audit](#security-audit)
  - [Helm Chart Polish](#helm-chart-polish)
  - [Docker Configuration Polish](#docker-configuration-polish)
- [Verification Checklist](#verification-checklist)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->


> QAR Module, Live TV, Polish

**Status**: 🔴 Not Started
**Tag**: `v0.9.0-rc1`
**Focus**: QAR Module, Live TV/DVR, Performance, Bug Fixes

**Depends On**: [v0.8.0](TODO_v0.8.0.md) (All features must be implemented)

---

## Overview

This is the **Release Candidate** milestone. It adds the QAR (adult content) module with complete isolation, Live TV/DVR support, photos module, comics module, and focuses heavily on performance optimization, bug fixes, and polish.

---

## Deliverables

### QAR Module (Adult Content)

> **Pirate-themed obfuscation** - See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology)

- [ ] **Database Schema** (`migrations/`)
  - [ ] `qar.crew` table (performers)
    - [ ] id, name, aliases
    - [ ] birth_date, scuttled (death_date)
    - [ ] origin (ethnicity), nationality
    - [ ] mast_color (hair), eye_color
    - [ ] height, ballast (weight)
    - [ ] cargo (measurements)
    - [ ] rigged (fake_tits), trimmed (circumcised)
    - [ ] cutlass (penis_length)
    - [ ] ink (tattoos), hooks (piercings)
    - [ ] bio, image_paths
    - [ ] stashdb_id, tpdb_id
  - [ ] `qar.ports` table (studios)
  - [ ] `qar.flags` table (tags)
  - [ ] `qar.voyages` table (scenes)
  - [ ] `qar.expeditions` table (movies)
  - [ ] `qar.treasures` table (galleries)
  - [ ] `qar.fleet` table (libraries)
  - [ ] `qar.voyage_crew` table (scene performers)
  - [ ] `qar.expedition_crew` table (movie performers)
  - [ ] `qar.voyage_flags` table (scene tags)
  - [ ] `qar.voyage_files` table
  - [ ] `qar.treasure_images` table
  - [ ] `qar.watch_progress` table
  - [ ] All indexes on external IDs

- [ ] **Entity** (`internal/content/qar/entity.go`)
  - [ ] Crew (performer)
  - [ ] Port (studio)
  - [ ] Flag (tag)
  - [ ] Voyage (scene)
  - [ ] Expedition (movie)
  - [ ] Treasure (gallery)
  - [ ] Fleet (library)

- [ ] **Repository** (`internal/content/qar/`)
  - [ ] All CRUD operations
  - [ ] Performer search/filter
  - [ ] Scene search/filter
  - [ ] Gallery operations

- [ ] **Service** (`internal/content/qar/service.go`)
  - [ ] Full CRUD for all entities
  - [ ] Search functionality
  - [ ] Progress tracking
  - [ ] Library scanning

- [ ] **Handler** (`internal/api/qar_handler.go`)
  - [ ] All endpoints under `/api/v1/legacy/*`
  - [ ] `GET /api/v1/legacy/crew`
  - [ ] `GET /api/v1/legacy/crew/:id`
  - [ ] `GET /api/v1/legacy/ports`
  - [ ] `GET /api/v1/legacy/voyages`
  - [ ] `GET /api/v1/legacy/voyages/:id`
  - [ ] `GET /api/v1/legacy/expeditions`
  - [ ] `GET /api/v1/legacy/expeditions/:id`
  - [ ] `GET /api/v1/legacy/treasures`
  - [ ] `GET /api/v1/legacy/treasures/:id`
  - [ ] `GET /api/v1/legacy/fleets`
  - [ ] Progress endpoints
  - [ ] Refresh endpoints

- [ ] **RBAC Integration**
  - [ ] `legacy:read` scope required
  - [ ] PIN verification middleware
  - [ ] Audit logging for all QAR access

- [ ] **Privacy Features**
  - [ ] Separate encryption key (`legacy.encryption_key`)
  - [ ] Encrypted field support
  - [ ] Session isolation
  - [ ] Quick hide shortcut

### StashDB Integration

- [ ] **StashDB Client** (`internal/service/metadata/stashdb/client.go`)
  - [ ] GraphQL API (genqlient)
  - [ ] API key authentication
  - [ ] Rate limiting

- [ ] **StashDB Service** (`internal/service/metadata/stashdb/service.go`)
  - [ ] Search performer
  - [ ] Get performer details
  - [ ] Search scene
  - [ ] Get scene details
  - [ ] Get scene fingerprints
  - [ ] Get studio details
  - [ ] Match by fingerprint

### Whisparr Integration

- [ ] **Whisparr Client** (`internal/service/metadata/whisparr/client.go`)
  - [ ] API v3 implementation
  - [ ] Authentication

- [ ] **Whisparr Service** (`internal/service/metadata/whisparr/service.go`)
  - [ ] Sync library
  - [ ] Get scenes/movies
  - [ ] Webhook handling
  - [ ] Full sync logic

### Live TV/DVR

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.livetv_channels` table
  - [ ] `public.livetv_programs` table (EPG)
  - [ ] `public.livetv_recordings` table
  - [ ] `public.livetv_timers` table

- [ ] **Service** (`internal/content/livetv/service.go`)
  - [ ] Channel management
  - [ ] EPG data handling
  - [ ] Recording scheduling
  - [ ] Live stream proxying

- [ ] **TVHeadend Integration** (`internal/service/metadata/tvheadend/`)
  - [ ] Channel sync
  - [ ] EPG sync
  - [ ] Recording management
  - [ ] Stream proxy

- [ ] **XMLTV Parser** (`internal/content/livetv/xmltv.go`)
  - [ ] Parse XMLTV files
  - [ ] Schedule EPG refresh
  - [ ] Map to programs

- [ ] **Handler** (`internal/api/livetv_handler.go`)
  - [ ] `GET /api/v1/livetv/channels`
  - [ ] `GET /api/v1/livetv/channels/:id/stream`
  - [ ] `GET /api/v1/livetv/guide`
  - [ ] `GET /api/v1/livetv/guide/now`
  - [ ] `GET /api/v1/livetv/recordings`
  - [ ] `POST /api/v1/livetv/recordings`
  - [ ] `DELETE /api/v1/livetv/recordings/:id`
  - [ ] `GET /api/v1/livetv/export/m3u`
  - [ ] `GET /api/v1/livetv/export/epg`

### Photos Module

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.photo_albums` table
  - [ ] `public.photos` table
    - [ ] EXIF data fields
    - [ ] Location (GPS)
    - [ ] Taken date
    - [ ] Camera/lens info
  - [ ] `public.photo_tags` table

- [ ] **Service** (`internal/content/photo/service.go`)
  - [ ] EXIF extraction
  - [ ] Thumbnail generation
  - [ ] Album organization
  - [ ] Timeline view
  - [ ] Map view support

- [ ] **Immich Integration** (optional)
  - [ ] Sync photos
  - [ ] Import albums

- [ ] **Handler** (`internal/api/photo_handler.go`)
  - [ ] `GET /api/v1/photos`
  - [ ] `GET /api/v1/photos/:id`
  - [ ] `GET /api/v1/photos/:id/thumbnail`
  - [ ] `GET /api/v1/photo-albums`
  - [ ] `GET /api/v1/photos/timeline`

### Comics Module

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.comics` table
  - [ ] `public.comic_series` table
  - [ ] `public.comic_publishers` table
  - [ ] `public.comic_reading_progress` table

- [ ] **Service** (`internal/content/comic/service.go`)
  - [ ] CBZ/CBR parsing
  - [ ] Page extraction
  - [ ] Reading progress

- [ ] **ComicVine Integration** (`internal/service/metadata/comicvine/`)
  - [ ] Search comics
  - [ ] Get comic details
  - [ ] Get series details

- [ ] **Handler** (`internal/api/comic_handler.go`)
  - [ ] `GET /api/v1/comics`
  - [ ] `GET /api/v1/comics/:id`
  - [ ] `GET /api/v1/comics/:id/pages`
  - [ ] `GET /api/v1/comics/:id/pages/:pageNum`
  - [ ] `GET /api/v1/comic-series`

### Performance Optimization

- [ ] **Profiling**
  - [ ] CPU profiling (pprof)
  - [ ] Memory profiling
  - [ ] Identify bottlenecks
  - [ ] Database query optimization

- [ ] **Caching Optimization**
  - [ ] Review cache TTLs
  - [ ] Add missing caches
  - [ ] Optimize cache invalidation
  - [ ] Monitor hit rates

- [ ] **Database Optimization**
  - [ ] Index analysis
  - [ ] Query plan review
  - [ ] Connection pool tuning
  - [ ] Vacuum/analyze scheduling

- [ ] **API Response Optimization**
  - [ ] Response compression (gzip, brotli)
  - [ ] Pagination review
  - [ ] Field selection (sparse fieldsets)
  - [ ] Batch endpoints

### Bug Fixes

- [ ] **Issue Triage**
  - [ ] Review all open issues
  - [ ] Prioritize by severity
  - [ ] Fix critical bugs
  - [ ] Fix major bugs

- [ ] **Stability**
  - [ ] Memory leak detection
  - [ ] Connection leak detection
  - [ ] Error handling review
  - [ ] Graceful degradation

### Documentation

- [ ] **User Guides**
  - [ ] Installation guide (Docker, K8s)
  - [ ] Configuration reference
  - [ ] Integration guides (Radarr, Sonarr, etc.)
  - [ ] Troubleshooting guide

- [ ] **API Documentation**
  - [ ] Complete OpenAPI spec
  - [ ] Swagger UI
  - [ ] API changelog

- [ ] **Admin Documentation**
  - [ ] Deployment guide
  - [ ] Backup/restore guide
  - [ ] Monitoring guide
  - [ ] Scaling guide

### Security Audit

- [ ] **Code Review**
  - [ ] Authentication flow review
  - [ ] Authorization review
  - [ ] Input validation review
  - [ ] SQL injection check
  - [ ] XSS prevention check

- [ ] **Penetration Testing**
  - [ ] API endpoint testing
  - [ ] Authentication testing
  - [ ] Session management testing
  - [ ] File upload testing

- [ ] **Dependency Audit**
  - [ ] Review all dependencies
  - [ ] Check for known vulnerabilities
  - [ ] Update vulnerable packages

### Helm Chart Polish

- [ ] **Values Review**
  - [ ] All configurable values documented
  - [ ] Sensible defaults
  - [ ] Resource limits/requests

- [ ] **Template Review**
  - [ ] Ingress annotations
  - [ ] Service annotations
  - [ ] Pod security context
  - [ ] Network policies

- [ ] **Documentation**
  - [ ] README with examples
  - [ ] Values reference
  - [ ] Upgrade notes

### Docker Configuration Polish

- [ ] **Docker Compose Review**
  - [ ] All services documented
  - [ ] Health checks configured
  - [ ] Resource limits set
  - [ ] Volume management

- [ ] **Docker Swarm Review**
  - [ ] Stack templates complete
  - [ ] Secrets management
  - [ ] Rolling update config
  - [ ] Placement constraints

- [ ] **Dockerfile Optimization**
  - [ ] Multi-stage builds optimized
  - [ ] Layer caching optimal
  - [ ] Image size minimized
  - [ ] Security hardening

---

## Verification Checklist

- [ ] QAR module fully functional and isolated
- [ ] Live TV channels stream
- [ ] EPG guide displays correctly
- [ ] DVR recording works
- [ ] Photos display with EXIF
- [ ] Comics reader works
- [ ] Performance targets met
- [ ] All critical/major bugs fixed
- [ ] Security audit passed
- [ ] Documentation complete
- [ ] All tests pass (80%+ coverage)
- [ ] CI pipeline passes

---

## Related Documentation

- [ROADMAP.md](ROADMAP.md) - Full roadmap overview
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Authoritative versions
- [ADULT_CONTENT_SYSTEM.md](../features/adult/ADULT_CONTENT_SYSTEM.md) - QAR design
- [STASHDB.md](../integrations/metadata/adult/STASHDB.md) - StashDB integration
- [WHISPARR.md](../integrations/servarr/WHISPARR.md) - Whisparr integration
- [LIVE_TV_DVR.md](../features/livetv/LIVE_TV_DVR.md) - Live TV design
- [PHOTOS_LIBRARY.md](../features/photos/PHOTOS_LIBRARY.md) - Photos design
- [COMICS_MODULE.md](../features/comics/COMICS_MODULE.md) - Comics design
