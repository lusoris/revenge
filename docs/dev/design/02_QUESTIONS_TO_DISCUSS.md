# Questions & Gaps to Discuss

> Temporary file for tracking discrepancies, gaps, and questions that need resolution

**Created**: 2026-01-30
**Status**: 🔄 Active collection

---

## Critical Discrepancies Found

### Package Version Mismatches

| Package | go.mod (actual) | SOURCES.yaml | 00_SOURCE_OF_TRUTH.md | Resolution |
|---------|-----------------|--------------|---------------------|------------|
| otter | v1.2.4 | v2.3.0 | v1.2.4 ✅ | Update SOURCES.yaml |
| typesense-go | v4.0.0-alpha2 ❌ | v3.2.0 | v3.2.0 ✅ | **CODE FIX NEEDED**: go.mod must use v3.2.0 stable |
| resty | go-resty/resty/v2 v2.17.1 ✅ | resty v3.0.0-b6 | v2.17.1 ✅ | Update SOURCES.yaml |

### Resolved (2026-01-30)

1. **Typesense Go Client**: ✅ DECISION MADE
   - **Answer**: Use stable v3.2.0, NOT alpha v4
   - **Principle**: "Bleeding edge stable" = newest STABLE version, never alphas/RCs
   - **Action needed**: Fix go.mod when doing code work (change v4 alpha → v3 stable)
   - **Verified**: [typesense-go releases](https://github.com/typesense/typesense-go/releases) - v3.2.0 is latest stable (March 2025)

2. **Resty HTTP Client**: ✅ CORRECT
   - go.mod uses v2 stable which is correct
   - SOURCES.yaml mentioned v3 beta incorrectly - needs update

---

## Missing from SOURCE_OF_TRUTH

### Not Yet Documented
- [ ] Full API endpoint inventory
- [ ] Complete database table list
- [ ] All environment variables with defaults
- [ ] OpenAPI spec file locations
- [ ] Migration file naming conventions
- [ ] Test coverage requirements
- [ ] CI/CD pipeline stages

### Need Design Docs
- [ ] Fingerprint service design doc
- [ ] Grants service design doc
- [ ] Search service design doc
- [ ] Analytics service design doc
- [ ] Notification service design doc

---

## Outdated References Found

### Old Namespace 'c' → 'qar' ✅ COMPLETE

**All Fixed (2026-01-30):**
- [x] INDEX.md, docs/dev/INDEX.md, MODULE_IMPLEMENTATION_TODO.md
- [x] ARCHITECTURE.md - updated to `qar` schema + SOT reference
- [x] WHISPARR.md - all `/api/v1/c/` → `/api/v1/legacy/`
- [x] POSTGRESQL.md - updated to `qar` schema
- [x] features/adult/INDEX.md - SOT reference added
- [x] REQUEST_SYSTEM.md - QAR terminology (expedition, voyage)
- [x] RBAC_CASBIN.md - `qar` schema reference
- [x] All external/adult integrations (6 files) - `/api/v1/legacy/`
- [x] All wiki/adult integrations (3 files) - `/api/v1/legacy/`
- [x] NEWS_SYSTEM.md, WIKI_SYSTEM.md
- [x] WHISPARR_STASHDB_SCHEMA.md, WHISPARR_V3_ANALYSIS.md
- [x] PLUGIN_ARCHITECTURE_DECISION.md

### Broken Internal Links (RESOLVED)

All broken links have been updated to point to 00_SOURCE_OF_TRUTH.md sections:

| File | Old Link | New Link |
|------|----------|----------|
| ARCHITECTURE.md | `PROJECT_STRUCTURE.md` | `00_SOURCE_OF_TRUTH.md#project-structure` |
| DATA_RECONCILIATION.md | `RIVER_JOBS.md` | `00_SOURCE_OF_TRUTH.md#river-job-queue-patterns` |
| SKIP_INTRO.md | `RIVER_JOBS.md` | `00_SOURCE_OF_TRUTH.md#river-job-queue-patterns` |
| TRICKPLAY.md | `RIVER_JOBS.md` | `00_SOURCE_OF_TRUTH.md#river-job-queue-patterns` |
| NEWS_SYSTEM.md | `RIVER_JOBS.md` | `00_SOURCE_OF_TRUTH.md#river-job-queue-patterns` |

**Status:** RESOLVED - Links now point to SOT master sections

### Unreferenced Documentation (RESOLVED)

- [x] GALLERY_MODULE.md - added to features/adult/INDEX.md
- [x] STASH.md - added to metadata/adult/INDEX.md

---

## Questions for Project Owner ✅ ALL RESOLVED

### Q1: Database Strategy ✅
- **Question**: Is SQLite truly needed for single-user deployments?
- **Answer**: **PostgreSQL ONLY** - No SQLite support needed
- **Rationale**: Simplifies codebase, pgx is the best driver, no dual-DB complexity

### Q2: Typesense Version ✅
- **Question**: v4 alpha in go.mod - intentional bleeding edge?
- **Answer**: **Pin to v3.2.0 stable** - "Bleeding edge stable" = newest STABLE version, never alphas/RCs

### Q3: Package Update Policy ✅
- **Question**: How aggressive on updates?
- **Answer**: **1 Minor Behind** - Stay one minor version behind for stability
- **Changelog Monitoring**: **Dependabot** (GitHub built-in, free)

### Q4: Documentation Priority ✅
- **Question**: Docs first or start scaffolding?
- **Answer**: **Docs first, no scaffolding until consistent**
- **Code Strategy**: Delete all existing code, fresh start after docs complete

---

## Live Docs Verification Needed

These need to be verified against actual package/API documentation:

| Source | Last Verified | Action |
|--------|---------------|--------|
| Go 1.25.6 release notes | Never | Check new features |
| pgx v5.8.0 changelog | Never | Breaking changes? |
| River v0.30.2 changelog | Never | Breaking changes? |
| ogen v1.18.0 changelog | Never | Breaking changes? |
| rueidis v1.0.71 | Never | Breaking changes? |
| Typesense API v30 | Never | Verify latest |
| StashDB GraphQL schema | Never | Schema changes? |

---

## Resolution Log

| Date | Item | Resolution |
|------|------|------------|
| 2026-01-30 | otter version | SOURCES.yaml was wrong, go.mod is correct |
| 2026-01-30 | resty version | SOURCES.yaml references v3 beta, we use stable v2 |
| | | |

---

## Document Naming Convention Proposal

**Problem:** Current docs lack numbering and categorization in filenames, making order and relationships unclear.

**Proposed Convention:**
```
[NN]-[CATEGORY]-[name].md

Examples:
01-arch-overview.md          # Architecture overview (read first)
02-arch-principles.md        # Design principles
03-arch-data-flow.md         # Data flow patterns
10-svc-auth.md               # Auth service
11-svc-user.md               # User service
20-mod-movie.md              # Movie module
21-mod-tvshow.md             # TV Show module
30-int-tmdb.md               # TMDb integration
31-int-radarr.md             # Radarr integration
90-ref-api-endpoints.md      # Reference: API endpoints
91-ref-db-tables.md          # Reference: Database tables
```

**Categories:**
- `arch` - Architecture & design
- `svc` - Services
- `mod` - Content modules
- `int` - Integrations
- `ops` - Operations
- `ref` - Reference tables

**Benefits:**
- Clear reading order (numbers)
- Category visible in filename
- Alphabetical sorting = logical order
- Easy to identify doc purpose

**Decision needed:** Adopt this convention? Rename existing files?

### Detailed Renaming Plan

**Phase 1: Core Architecture (00-09)**
```
architecture/
├── 00-arch-source-of-truth.md      # 00_SOURCE_OF_TRUTH.md (master)
├── 01-arch-overview.md             # ARCHITECTURE.md
├── 02-arch-principles.md           # DESIGN_PRINCIPLES.md
├── 03-arch-player.md               # PLAYER_ARCHITECTURE.md
```

**Phase 2: Technical Infrastructure (10-19)**
```
technical/
├── 10-tech-stack.md                # TECH_STACK.md
├── 11-tech-configuration.md        # CONFIGURATION.md
├── 12-tech-audio-streaming.md      # AUDIO_STREAMING.md
├── 13-tech-offloading.md           # OFFLOADING.md
```

**Phase 3: Content Modules (20-39)**
```
features/video/
├── 20-mod-movie.md                 # MOVIE_MODULE.md
├── 21-mod-tvshow.md                # TVSHOW_MODULE.md

features/music/
├── 22-mod-music.md                 # MUSIC_MODULE.md

features/audiobooks/
├── 23-mod-audiobook.md             # (new)

features/podcasts/
├── 24-mod-podcast.md               # PODCASTS.md

features/photos/
├── 25-mod-photos.md                # PHOTOS_LIBRARY.md

features/comics/
├── 26-mod-comics.md                # COMICS_MODULE.md

features/livetv/
├── 27-mod-livetv.md                # LIVE_TV_DVR.md

features/adult/ (qar namespace)
├── 28-qar-performers.md            # PERFORMERS.md
├── 29-qar-scenes.md                # SCENES.md
├── 2A-qar-studios.md               # STUDIOS.md
├── 2B-qar-galleries.md             # GALLERY_MODULE.md
├── 2C-qar-reconciliation.md        # DATA_RECONCILIATION.md
```

**Phase 4: Shared Features (40-49)**
```
features/shared/
├── 40-feat-access-controls.md      # ACCESS_CONTROLS.md
├── 41-feat-analytics.md            # ANALYTICS_SERVICE.md
├── 42-feat-client-support.md       # CLIENT_SUPPORT.md
├── 43-feat-content-rating.md       # CONTENT_RATING.md
├── 44-feat-i18n.md                 # I18N.md
├── 45-feat-news.md                 # NEWS_SYSTEM.md
├── 46-feat-notifications.md        # NOTIFICATIONS.md
├── 47-feat-scrobbling.md           # SCROBBLING.md
├── 48-feat-ticketing.md            # TICKETING_SYSTEM.md
├── 49-feat-voice-control.md        # VOICE_CONTROL.md
```

**Phase 5: Playback Features (50-59)**
```
features/playback/
├── 50-play-media-enhancements.md   # MEDIA_ENHANCEMENTS.md
├── 51-play-skip-intro.md           # SKIP_INTRO.md
├── 52-play-syncplay.md             # SYNCPLAY.md
├── 53-play-trickplay.md            # TRICKPLAY.md
├── 54-play-watch-next.md           # WATCH_NEXT_CONTINUE_WATCHING.md
├── 55-play-release-calendar.md     # RELEASE_CALENDAR.md
```

**Phase 6: Integrations (60-79)**
```
integrations/metadata/video/
├── 60-int-tmdb.md                  # TMDB.md
├── 61-int-tvdb.md                  # THETVDB.md
├── 62-int-omdb.md                  # OMDB.md

integrations/servarr/
├── 63-int-radarr.md                # (exists in metadata)
├── 64-int-sonarr.md                # SONARR.md
├── 65-int-lidarr.md                # LIDARR.md

integrations/metadata/adult/
├── 70-qar-stashdb.md               # STASHDB.md
├── 71-qar-tpdb.md                  # THEPORNDB.md
├── 72-qar-whisparr.md              # WHISPARR.md
```

**Phase 7: Operations (80-89)**
```
operations/
├── 80-ops-setup.md                 # SETUP.md
├── 81-ops-gitflow.md               # GITFLOW.md
├── 82-ops-db-healing.md            # DATABASE_AUTO_HEALING.md
├── 83-ops-best-practices.md        # BEST_PRACTICES.md
├── 84-ops-reverse-proxy.md         # REVERSE_PROXY.md
```

**Phase 8: References (90-99)**
```
├── 90-ref-api-endpoints.md         # API reference (generate from OpenAPI)
├── 91-ref-db-schemas.md            # Database schemas
├── 92-ref-env-vars.md              # Environment variables
├── 93-ref-error-codes.md           # Error codes
├── 99-ref-glossary.md              # Terms and definitions
```

**INDEX.md files:** Keep as-is (no numbering), they serve as folder navigation.

**Execution:** This is a significant refactor. Recommend doing after initial MVP to avoid constant link updates during active development.

---

## Design Strategy Requirements (MANDATORY)

**Principle:** All advanced design patterns and coding strategies MUST be documented in 00_SOURCE_OF_TRUTH.md BEFORE implementation.

**Required in SOT:**
- [x] All performance patterns (caching, pooling, batching)
- [x] All resilience patterns (circuit breaker, retry, fallback)
- [x] All security patterns (auth, RBAC, isolation)
- [x] All async patterns (jobs, queues, workers)
- [x] All data patterns (transactions, consistency, partitioning)
- [x] All API patterns (versioning, errors, pagination)

**Why:**
- No "write first, fix later" approach
- Prevents wasted time on reiteration
- Every implementation inherits patterns from SOT
- Consistency across all modules

**Design Patterns (RESOLVED 2026-01-30):**

| Pattern | Decision | Notes |
|---------|----------|-------|
| Error Handling | **Sentinels (internal) + Custom APIError (external)** | Combo for type-safe errors + API responses |
| Testing | **Table-driven + testify + mockery** | mockery for auto-generated mocks |
| Logging | **Text (Dev, tint) + JSON (Prod)** | slog with tint handler |
| Metrics | **Prometheus + OpenTelemetry** | Both - Prometheus for K8s, OTel for traces |
| Validation | **ogen (API) + go-playground/validator (Business)** | Separate layers |
| Pagination | **Cursor (default) + Offset (option)** | Cursor for performance, Offset for compatibility |
| Integration Tests | **testcontainers-go** | Coder-compatible, real containers |
| Test Coverage | **80% minimum** | Required for all packages |

---

## Go Package Review (COMPREHENSIVE) ✅ RESOLVED

> After adding ErsatzTV, Raft, Network QoS, Container Orchestration - review ALL package choices

**Principle**: Best/fastest solution for each use case. No compromises. Bleeding edge STABLE only.

### CORE PRINCIPLE: Metadata Priority Chain

```
Priority Order (ALWAYS):
1. LOCAL CACHE     → First, instant UI display
2. ARR SERVICES    → Radarr, Sonarr, Whisparr (cached metadata)
3. INTERNAL        → Stash-App (if connected)
4. EXTERNAL        → TMDb, StashDB.org, MusicBrainz, etc.
5. ENRICHMENT      → Background jobs, lower priority, seamless
```

> **If data exists locally, ALWAYS use local first, then fallback to external. This applies to ALL data!**

### Core Infrastructure Packages

#### 1. Database Driver ✅
**Current**: `github.com/jackc/pgx/v5`
**Decision**: **pgx v5** - No replacement, best performance

| Question | Answer |
|----------|--------|
| Is pgx still the fastest PostgreSQL driver for Go? | **YES** - Undisputed |
| Any breaking changes in v5.8.0 we need to handle? | No, stable |
| Do we need pgx for SQLite or pure database/sql? | **N/A** - PostgreSQL only |

#### 2. Redis/Cache Client ✅
**Current**: `github.com/redis/rueidis` (14x faster than go-redis)
**Decision**: **rueidis** - 14x faster, no replacement needed

| Question | Answer |
|----------|--------|
| Is rueidis still the fastest option? | **YES** - 14x faster than go-redis |
| Does it support all Dragonfly features we need? | **YES** - Full compatibility |
| Any issues with client-side caching in our use case? | No - perfect for L1/L2 caching |

#### 3. In-Memory Cache ✅
**Current**: `github.com/maypok86/otter` (W-TinyLFU)
**Decision**: **otter** - W-TinyLFU, 50% less memory than ristretto

| Question | Answer |
|----------|--------|
| Otter vs Ristretto - which is actually faster now? | **otter** - Faster + 50% less RAM |
| Do we need the S3-FIFO algorithm or is LRU enough? | W-TinyLFU better than LRU |
| Memory overhead considerations for embedded/single-user? | otter is memory-efficient |

#### 4. Job Queue ✅
**Current**: `github.com/riverqueue/river` (PostgreSQL-backed)
**Decision**: **River** - PostgreSQL-native, transactional enqueueing

| Question | Answer |
|----------|--------|
| River vs Asynq - PostgreSQL vs Redis trade-offs? | **River** - Transactional with DB |
| Is transactional enqueueing critical for us? | **YES** - Atomic with DB changes |
| River's unique job constraints sufficient? | **YES** |
| Need for job scheduling (cron-like)? River has it? | **YES** - Periodic Jobs |

#### 5. HTTP Client ✅
**Current**: `github.com/go-resty/resty/v2` (stable)
**Decision**: **resty v2** - Middleware, stable, kein v3 Beta

| Question | Answer |
|----------|--------|
| Do we need resty's middleware or is stdlib enough? | **YES** - Middleware for retry, auth, logging |
| fasthttp for high-throughput external API calls? | No - resty is sufficient |
| Retry/backoff - resty built-in vs custom? | **resty built-in** |

#### 6. Configuration ✅
**Current**: `github.com/knadh/koanf/v2`
**Decision**: **koanf v2** - Hot reload via Watch(), better than viper

| Question | Answer |
|----------|--------|
| koanf vs viper performance/features? | **koanf** - Faster, more flexible |
| Hot reload necessary or just restart? | **YES** - Hot reload via Watch() |
| Environment variable precedence working correctly? | **YES** - koanf handles this well |

### New Feature Packages

#### 7. Raft Consensus (NEW) ✅
**Current**: `github.com/hashicorp/raft` (proposed)
**Decision**: **hashicorp/raft** - Echte Consensus, no half measures

| Question | Answer |
|----------|--------|
| hashicorp/raft vs etcd/raft - batteries-included vs minimal? | **hashicorp/raft** - Batteries-included |
| Do we actually need consensus or is Redis-based leader election enough? | **Real Raft** - Preparation for HA |
| Is Raft overkill for home server clustering? | No - "no half measures in core" |
| Dragonboat claims 10x faster - worth investigating? | No - hashicorp is more proven |

#### 8. IPTV/Streaming (NEW - ErsatzTV) ✅
**Current**: Custom REST client for ErsatzTV
**Decision**: **Nur REST Client** - Revenge steuert ErsatzTV via API

| Question | Answer |
|----------|--------|
| Should we rely on ErsatzTV or build native IPTV? | **ErsatzTV** - No native IPTV generation |
| HLS generation - `bluenviron/gohlslib` or FFmpeg? | N/A - ErsatzTV handles this |
| XMLTV generation - existing Go packages? | N/A - ErsatzTV handles this |

#### 9. WebSocket ✅
**Current**: `github.com/coder/websocket`
**Decision**: **gobwas/ws** - Maximum performance

| Question | Answer |
|----------|--------|
| coder/websocket = fork of nhooyr, is it maintained? | Yes, but gobwas faster |
| Do we need gobwas/ws for ultra-high performance? | **YES** - Zero-alloc, best performance |
| WebSocket compression needed? | Optional - gobwas supports it |

### API & Serialization Packages

#### 10. OpenAPI Codegen ✅
**Current**: `github.com/ogen-go/ogen`
**Decision**: **ogen** - Zero-alloc JSON, best code generation

| Question | Answer |
|----------|--------|
| ogen vs oapi-codegen - which generates cleaner code? | **ogen** - Cleaner, more type-safe |
| ogen's zero-allocation JSON - necessary for us? | **YES** - Performance for high-throughput |
| Client generation quality for Arr APIs? | **Good** - Type-safe client stubs |

#### 11. JSON Parsing ✅
**Current**: `github.com/go-faster/jx` (zero-allocation)
**Decision**: **go-faster/jx** - ogen-compatible, zero-allocation

| Question | Answer |
|----------|--------|
| jx vs go-json - which is faster for our use case? | **jx** - Zero-alloc, ogen integration |
| Do we need streaming JSON for large responses? | **YES** - jx supports streaming |
| Compatibility with ogen's generated code? | **Perfect** - Same author, native integration |

#### 12. GraphQL Client (StashDB, AniList) ✅
**Current**: Not specified
**Decision**: **Khan/genqlient** - Codegen, type-safe

| Question | Answer |
|----------|--------|
| Which GraphQL client is best for type-safe queries? | **genqlient** - Codegen |
| genqlient (code-gen) vs runtime clients? | **Codegen** - Type-safe, better |
| Need subscriptions for real-time updates? | No for metadata |

### Media Processing Packages

#### 13. Image Processing ✅
**Current**: `github.com/davidbyttow/govips/v2`
**Decision**: **govips** - Faster than bimg, CGo acceptable

| Question | Answer |
|----------|--------|
| govips vs bimg - both use libvips, which wrapper is better? | **govips** - Better maintained, faster |
| CGo dependency acceptable for image processing? | **YES** - Performance critical |
| Pure Go alternative for CGo-free builds? | Not needed - CGo acceptable for media server |

#### 14. FFmpeg Integration ✅
**Current**: `github.com/asticode/go-astiav` (proposed)
**Decision**: **go-astiav** - Native CGo bindings

| Question | Answer |
|----------|--------|
| go-astiav vs shell exec - performance difference? | **go-astiav** - Less overhead |
| CGo dependency for FFmpeg acceptable? | **YES** - For performance |
| Need Go bindings or is CLI sufficient? | **Bindings** - Better control |

#### 15. Audio Metadata ✅
**Current**: `github.com/dhowden/tag`
**Decision**: **go-taglib** - CGo, Read+Write ALL formats

| Question | Answer |
|----------|--------|
| dhowden/tag covers all formats (ID3, FLAC, MP4)? | No - go-taglib better |
| Performance for large music libraries? | **go-taglib** - Native performance |
| Write support needed or read-only? | **YES** - Read+Write for alignment |

### Observability Packages

#### 16. Logging ✅
**Current**: `log/slog` (stdlib) + `github.com/lmittmann/tint` (dev) + `uber-go/zap` (prod)
**Decision**: **slog/tint (dev) + zap (prod)** - Dev colors, prod performance

| Question | Answer |
|----------|--------|
| slog sufficient or need zap's performance? | **Both** - slog/tint dev, zap prod |
| tint for colored output - production-ready? | **Dev only** - zap for prod JSON |
| Structured logging format (JSON vs text)? | **Text dev, JSON prod** |

#### 17. Metrics ✅
**Current**: `github.com/prometheus/client_golang` + OpenTelemetry
**Decision**: **Both** - Prometheus for K8s, OTel for traces

| Question | Answer |
|----------|--------|
| Prometheus vs OpenTelemetry metrics? | **Both** - Different purposes |
| Both needed for different use cases? | **YES** - Prometheus K8s, OTel traces |
| Cardinality concerns for media server? | Managed - use labels wisely |

#### 18. Tracing ✅
**Current**: `go.opentelemetry.io/otel`
**Decision**: **OpenTelemetry** - Standard, vendor-agnostic

| Question | Answer |
|----------|--------|
| OpenTelemetry vs vendor-specific? | **OTel** - Vendor-agnostic standard |
| Trace sampling strategy for media server? | Head sampling, 10% default |
| Performance overhead acceptable? | **YES** - Minimal with sampling |

### Security Packages

#### 19. RBAC ✅
**Current**: `github.com/casbin/casbin/v2` + `open-policy-agent/opa`
**Decision**: **Casbin + OPA** - Casbin for roles, OPA for complex policies

| Question | Answer |
|----------|--------|
| Casbin RBAC/ABAC sufficient for our needs? | **Casbin for roles** - Simple RBAC |
| OPA for more complex policy decisions? | **YES** - Data-driven ABAC policies |
| Performance for frequent permission checks? | Good - Casbin fast, OPA cached |

#### 20. Cryptography ✅
**Current**: `golang.org/x/crypto` (stdlib)
**Decision**: **golang.org/x/crypto** - No alternative needed

| Question | Answer |
|----------|--------|
| Argon2id for password hashing confirmed? | **YES** - Argon2id |
| AES-256-GCM for field encryption? | **YES** - AES-256-GCM |
| Need hardware acceleration (AES-NI)? | Auto-detected by Go runtime |

### Testing Packages

#### 21. Testing Framework ✅
**Current**: `github.com/stretchr/testify` + `github.com/vektra/mockery`
**Decision**: **testify + mockery** - Table-driven, auto mocks

| Question | Answer |
|----------|--------|
| testify sufficient or need BDD framework? | **testify** - No BDD needed |
| Mock generation - testify/mock vs mockery vs gomock? | **mockery** - Auto-generation |
| Table-driven tests pattern adopted? | **YES** - Standard pattern |

#### 22. Integration Testing ✅
**Current**: `github.com/testcontainers/testcontainers-go` + `github.com/fergusstrange/embedded-postgres`
**Decision**: **Both** - Containers for integration, embedded for unit

| Question | Answer |
|----------|--------|
| Testcontainers vs Docker Compose test setup? | **testcontainers** - Coder-compatible |
| Embedded PostgreSQL for faster tests? | **YES** - embedded-postgres for unit tests |
| Test parallelization strategy? | Per-package with unique DB schemas |

### Container/Orchestration Packages

#### 23. Kubernetes Client ✅
**Current**: `sigs.k8s.io/controller-runtime`
**Decision**: **controller-runtime** - Operator pattern, real foundations

| Question | Answer |
|----------|--------|
| Do we need K8s client in the app itself? | **YES** - For Operator pattern |
| Or just Helm/kubectl for deployment? | **Both** - Helm + Operator |
| Operator pattern for self-healing? | **YES** - controller-runtime |

### Summary ✅ ALL RESOLVED (2026-01-30)

| Area | Decision | Confidence |
|------|----------|------------|
| Database | pgx v5 | HIGH |
| Cache (distributed) | rueidis | HIGH |
| Cache (local) | otter | HIGH |
| Job Queue | River | HIGH |
| HTTP Client | resty v2 | HIGH |
| Config | koanf v2 (hot reload) | HIGH |
| Raft | hashicorp/raft | HIGH |
| WebSocket | gobwas/ws | HIGH |
| OpenAPI | ogen + go-faster/jx | HIGH |
| JSON | go-faster/jx | HIGH |
| GraphQL | Khan/genqlient | HIGH |
| Image | govips (CGo OK) | HIGH |
| FFmpeg | go-astiav (CGo OK) | HIGH |
| Audio | go-taglib (CGo OK) | HIGH |
| Logging | slog/tint (dev) + zap (prod) | HIGH |
| Metrics | Prometheus + OpenTelemetry | HIGH |
| RBAC | Casbin + OPA | HIGH |
| Testing | testify + mockery | HIGH |
| Integration Tests | testcontainers + embedded-postgres | HIGH |
| K8s | controller-runtime | HIGH |

**All Action Items Complete** - Package decisions finalized and documented in 00_SOURCE_OF_TRUTH.md

---

---

## New Questions (From Integration Audit 2026-01-30)

### Q5: Multi-Provider OIDC Architecture
- **Context**: GENERIC_OIDC.md proposes `oidc_providers` table supporting multiple OIDC providers
- **Question**: Support multiple simultaneous OIDC providers (multiple login buttons)?
- **Options**: Single provider only | Multiple providers with selection

### Q6: Webhook Handler Architecture
- **Context**: RADARR.md references webhook handling for import/upgrade/rename/delete events
- **Question**: Centralized webhook dispatcher or per-service handlers?
- **Missing**: Validation strategy, idempotency handling, payload schema

### Q7: Rate Limiting Strategy
- **Context**: Different services have different limits (TMDb 40/10s, Trakt 1000/5min, StashDB undefined)
- **Question**: Global rate limiter or per-provider limiters?

### Q8: Watch History Sync Conflict Resolution
- **Context**: TRAKT.md specifies "Keep Trakt timestamp if earlier, Revenge if later"
- **Question**: Should conflict resolution be configurable per service or global policy?

---

## Notes

- This file is temporary - items should be resolved and moved to appropriate docs
- When resolved, update 00_SOURCE_OF_TRUTH.md and remove from here
- Questions for owner should be asked and answers documented in design docs
- **Live docs verification DELAYED** until package choices finalized
