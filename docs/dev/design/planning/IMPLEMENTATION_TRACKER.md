# Revenge - Preparation Master Plan

> **CRITICAL**: Complete this BEFORE any implementation begins!
> Prevents breaking changes, rework, and architectural mistakes.

**Generated**: 2026-01-28
**Status**: 🔴 IN PROGRESS
**Completion**: 27% (11/41 external APIs, 14/14 UX/UI sources, 0/133 docs/instructions created)

**Latest Update** (2026-01-28):
- ✅ Fetched 14 UX/UI design & frontend resources
- ✅ Added Phase 5: UX/UI Design & Frontend Architecture (Week 9-11)
- ✅ Expanded critical docs from 14 to 23 files
- ✅ Created Appendix A: UX/UI Design & Frontend Resources
- ✅ Defined 9 new frontend instruction files
- ✅ Established frontend quality gates

---

## Executive Summary

**Problem**: Starting implementation without complete documentation leads to:
- Breaking API changes mid-development
- Inconsistent patterns across modules
- Security vulnerabilities discovered late
- Performance issues requiring refactoring
- Multiple rewrites due to missing requirements

**Solution**: Complete ALL documentation, instructions, and API research FIRST.

**Timeline**: ~10-12 weeks (full-time) before coding begins (expanded from 8-10 weeks with UX/UI research)
**Resources**: Documentation Gap Analysis + External Integrations TODO + Live API docs + UX/UI Best Practices

**Total Deliverables**: 133 files
- 41 Integration docs (`docs/integrations/<SERVICE>.md`)
- 41 Client instructions (`.github/instructions/<service>-client.instructions.md`)
- 23 Critical user docs (up from 14, includes 5 new UX/UI docs)
- 16 Critical instruction files (up from 7, includes 9 new frontend files)
- 12 High-priority instruction files (up from 9)

---

## Phase 1: External API Documentation (Week 1-2)

### 1.1 Fetch ALL Live API Documentation

**Status**: 🟡 27% Complete (11/41 services)

#### ✅ Already Fetched (11 services)

| Service | API Version | Auth Method | Rate Limit | Status |
|---------|-------------|-------------|------------|--------|
| Radarr | v3 | X-Api-Key | Self-hosted (none) | ✅ DONE |
| Sonarr | v3 | X-Api-Key | Self-hosted (none) | ✅ DONE |
| Lidarr | v1 | X-Api-Key | Self-hosted (none) | ✅ DONE |
| TMDb | v3 | api_key query param | 40 req/10 sec | ✅ DONE |
| TheTVDB | v4 | JWT (1 month) | Not specified | ✅ DONE |
| MusicBrainz | v2 | None (User-Agent required) | **1 req/sec CRITICAL** | ✅ DONE |
| Audiobookshelf | REST | Bearer token | Self-hosted (none) | ✅ DONE |
| Trakt | v2 | OAuth2 | Implemented | ✅ DONE |
| Fanart.tv | v3 | api-key header OR query | Not specified | ✅ DONE |
| Last.fm | v2 | API key | Reasonable usage | ✅ DONE |

#### ❌ To Fetch (28 services)

**Servarr Ecosystem (2)**:
- [ ] Whisparr (adult content management)
- [ ] Readarr (book management)

**Metadata Providers (11)**:
- [ ] OMDb (movie metadata, IMDb alternative)
- [ ] ThePosterDB (posters, fanart)
- [ ] Spotify (music metadata)
- [ ] Discogs (music metadata, vinyl focus)
- [ ] Goodreads (book ratings, reviews)
- [ ] OpenLibrary (book metadata)
- [ ] Audible (audiobook metadata)
- [ ] Hardcover (book tracking) - **ISSUE: Docs not accessible**
- [ ] Stash/StashDB (adult content metadata)
- [ ] ThePornDB (adult content metadata)

**Scrobbling Services (3)**:
- [ ] ListenBrainz (music scrobbling)
- [ ] Letterboxd (movie ratings/reviews)
- [ ] Simkl (anime/TV tracking)

**Authentication/SSO (3)**:
- [ ] Authelia (authentication proxy)
- [ ] Authentik (identity provider)
- [ ] Keycloak (identity & access management)

**Infrastructure (4)**:
- [ ] Blackbeard (Revenge's transcoding service)
- [ ] TVHeadend (Live TV backend)
- [ ] NextPVR (Live TV/DVR)
- [ ] Chromecast (casting protocol)

**Device Protocols (1)**:
- [ ] DLNA (media streaming protocol)

**Search/Discovery (2)**:
- [ ] Spotify API (detailed search/recommendations)
- [ ] Last.fm recommendations (already have scrobbling)

**Other (1)**:
- [ ] OIDC Generic (OpenID Connect specification)

### 1.2 Create Integration Documentation

**For EACH service**, create:

1. **`docs/integrations/<SERVICE>.md`** (user-facing)
   - Purpose & use cases
   - Authentication setup
   - Configuration examples
   - Rate limits & best practices
   - Troubleshooting

2. **`.github/instructions/<service>-client.instructions.md`** (developer)
   - Client implementation patterns
   - Error handling
   - Retry logic
   - Testing patterns
   - Code examples

**Estimated**: 41 services × 2 files = **82 files to create**
**Time**: ~3-4 days per service (research + write) = **~120-160 days**
**Parallel Work**: Can batch similar services (e.g., all Servarr, all scrobblers)

---

## Phase 2: Core Infrastructure Documentation (Week 3-4)

### 2.1 Critical Foundation Docs (BLOCKERS)

#### ✅ River Job Queue

**Files**:
- [ ] `docs/JOB_QUEUE.md` (user guide)
- [ ] `.github/instructions/river-job-queue.instructions.md` (patterns)

**Content**:
- River architecture (workers, queues, jobs)
- Job definition patterns
- Worker registration with fx
- Error handling & retries
- Priority queues
- Dead letter queue
- Monitoring & metrics
- Testing job workers
- Example jobs for each module

**Research**:
- River documentation (riverqueue.com)
- Sidekiq patterns (Ruby, applicable)
- Temporal workflows (inspiration)

**Time**: 4 days

---

#### ✅ Typesense Search

**Files**:
- [ ] `docs/SEARCH.md` (user guide)
- [ ] `.github/instructions/typesense-integration.instructions.md` (patterns)

**Content**:
- Typesense architecture
- Per-module collections design
- Schema definitions
- Indexing patterns (create/update/delete)
- Query syntax & filters
- Faceting & sorting
- Synonyms & typo tolerance
- Reindexing strategies
- Performance optimization
- Testing search integration

**Research**:
- Typesense 0.25+ documentation
- Algolia patterns (inspiration)
- Elasticsearch migration guides

**Time**: 4 days

---

#### ✅ Database Schema

**Files**:
- [ ] `docs/DATABASE_SCHEMA.md` (complete reference)

**Content**:
- Full ER diagrams (Mermaid)
- Per-module table definitions
- Shared infrastructure tables
- Index definitions & rationale
- Foreign key relationships
- Migration history
- Performance considerations
- Adult schema (`c`) isolation

**Research**:
- PostgreSQL 18 features
- ER diagram best practices

**Time**: 3 days

---

#### ✅ Security Architecture

**Files**:
- [ ] `docs/SECURITY.md` (user guide)
- [ ] `.github/instructions/security-best-practices.instructions.md` (patterns)

**Content**:
- JWT implementation
- OIDC authentication flow
- Adult content isolation (schema `c`)
- Encryption at rest
- Audit logging
- Security headers (CSP, HSTS, etc.)
- CORS configuration
- Input validation
- SQL injection prevention (sqlc)
- XSS prevention
- Rate limiting
- File upload security

**Research**:
- OWASP recommendations
- Go security best practices
- securityheaders.com guidelines

**Time**: 4 days

---

#### ✅ WebSocket Protocol

**Files**:
- [ ] `docs/WEBSOCKET.md` (protocol specification)
- [ ] `.github/instructions/websocket-handlers.instructions.md` (patterns)

**Content**:
- WebSocket architecture
- Authentication (JWT on upgrade)
- Message types (JSON schema)
- Pub/sub patterns
- Watch Party protocol
- Progress tracking
- Quality switching
- Reconnection strategy
- Scaling (multiple instances)
- Testing WebSocket handlers

**Research**:
- coder/websocket documentation (modern, zero deps, ISC license)
- Centrifugo patterns
- Phoenix Channels (Elixir) inspiration

**Time**: 3 days

---

### 2.2 API & Code Generation Patterns

#### ✅ ogen API Patterns

**Files**:
- [ ] `.github/instructions/ogen-api-patterns.instructions.md`

**Content**:
- OpenAPI spec structure
- Handler interface implementation
- Validation handling
- Error response formatting
- Testing ogen handlers
- Code generation workflow
- Best practices

**Time**: 1 day

---

#### ✅ Dragonfly Cache Patterns

**Files**:
- [ ] `.github/instructions/dragonfly-cache-patterns.instructions.md`

**Content**:
- Cache key patterns
- TTL strategies
- Invalidation patterns
- Cache-aside pattern
- Error handling (cache miss)
- Testing caching logic
- Performance monitoring

**Research**:
- Redis patterns (applicable to Dragonfly)
- rueidis documentation (14x faster than go-redis, auto-pipelining)

**Time**: 1 day

---

**Phase 2 Total**: ~20 days

---

## Phase 3: Development Standards (Week 5-6)

### 3.1 Testing Strategy

**Files**:
- [ ] `docs/TESTING_STRATEGY.md` (overall strategy)
- [ ] `.github/instructions/content-module-testing.instructions.md` (module patterns)
- [ ] Update `.github/instructions/testing-patterns.instructions.md`

**Content**:
- Unit testing standards
- Integration testing (testcontainers)
- E2E testing approach
- Per-module test structure
- Test fixtures & data management
- Coverage thresholds
- CI/CD integration
- Testing River jobs
- Testing search integration
- Testing WebSocket handlers

**Research**:
- testcontainers-go documentation
- Go testing best practices
- pgx test helpers

**Time**: 4 days

---

### 3.2 Error Handling & Logging

**Files**:
- [ ] `.github/instructions/error-handling-patterns.instructions.md`
- [ ] `.github/instructions/logging-standards.instructions.md`

**Content**:
- Error wrapping patterns
- Domain error types
- HTTP status code mapping
- User-friendly error messages
- Structured logging with slog
- Log levels (Debug/Info/Warn/Error)
- Sensitive data redaction
- Request ID propagation
- Error testing

**Research**:
- Go error handling best practices
- slog documentation
- Log aggregation patterns

**Time**: 2 days

---

**Phase 3 Total**: ~6 days

---

## Phase 4: Production Readiness (Week 7-8)

### 4.1 Deployment & Operations

**Files**:
- [ ] `docs/MONITORING.md`
- [ ] `docs/BACKUP_RESTORE.md`
- [ ] `docs/PERFORMANCE_TUNING.md`
- [ ] `docs/CICD.md`
- [ ] `docs/TROUBLESHOOTING.md`

**Content**:
- Prometheus metrics
- Grafana dashboards
- Structured logging (Loki)
- Alerting rules
- Backup strategies (PostgreSQL, media files)
- Disaster recovery procedures
- Database optimization
- Query analysis
- Connection pooling
- Index tuning
- GitHub Actions workflows
- Release process (Release Please)
- Docker builds (multi-stage)
- Common issues & solutions
- Debug mode
- Health check troubleshooting

**Research**:
- PostgreSQL performance tuning
- Docker multi-stage builds
- GitHub Actions best practices
- GoReleaser patterns

**Time**: 10 days

---

### 4.2 API Design & Best Practices

**Files**:
- [ ] Update `docs/API.md`
- [ ] `.github/instructions/api-pagination.instructions.md`
- [ ] `.github/instructions/rate-limiting.instructions.md`

**Content**:
- REST API design principles
- Pagination (cursor-based)
- Filtering & sorting
- Partial responses
- Error response format
- Versioning strategy
- Rate limiting patterns
- CORS configuration

**Research**:
- Microsoft REST API Guidelines
- Google API Design Guide
- Stripe API patterns

**Time**: 3 days

---

**Phase 4 Total**: ~13 days

---

## Phase 5: UX/UI Design & Frontend Architecture (Week 9-11)

### 5.1 UX/UI Design Principles & Standards

**Files**:
- [ ] `docs/UX_DESIGN_PRINCIPLES.md`
- [ ] `.github/instructions/frontend-ux-guidelines.instructions.md`
- [ ] `.github/instructions/accessibility-standards.instructions.md`

**Content**:
- W3C WCAG 2.2 compliance (A, AA, AAA levels)
- Nielsen's 10 Usability Heuristics
- Laws of UX (Hick's Law, Fitts's Law, Miller's Law, Jakob's Law)
- Government Design Principles (UK GDS patterns)
- Material Design 3 Expressive patterns
- Apple Human Interface Guidelines
- Fluent Design System (Microsoft)
- Atlassian Design System principles
- Carbon Design System (IBM)
- Gestalt principles
- Progressive disclosure
- Recognition over recall
- User control & freedom
- Error prevention & recovery
- Consistency across platforms

**Research Sources** (all fetched):
- ✅ W3C WCAG 2.2 (Perceivable, Operable, Understandable, Robust)
- ✅ ISO 9241-11:2018 (Usability definition)
- ✅ Nielsen Norman Group (10 Heuristics, UX best practices)
- ✅ Laws of UX (26 laws: Aesthetic-Usability, Choice Overload, Chunking, Cognitive Load, etc.)
- ✅ Google Material Design 3 (Expressive components, motion, shape library)
- ✅ Apple HIG (Hierarchy, Harmony, Consistency, Accessibility)
- ✅ Microsoft Fluent 2 (Web, iOS, Android, Windows components)
- ✅ Atlassian Design System (Rovo AI patterns, unified language)
- ✅ IBM Carbon Design System (Web Components, React, Angular, Vue, Svelte)
- ✅ UK Government Design Principles (11 principles: start with user needs, do less, iterate)
- ✅ Interaction Design Foundation (UX research methods, user-centered design)
- ✅ Baymard Institute (E-commerce UX research, 18,000+ design examples)
- ✅ Smashing Magazine (UX design patterns, accessibility, mobile-first)
- ✅ web.dev Patterns (Animation, Layout, Components, Media, Theming)

**Time**: 5 days

---

### 5.2 Frontend Component Library & Design System

**Files**:
- [ ] `docs/FRONTEND_COMPONENTS.md`
- [ ] `.github/instructions/svelte-component-patterns.instructions.md`
- [ ] `.github/instructions/shadcn-svelte-customization.instructions.md`

**Content**:
- SvelteKit 2 architecture
- shadcn-svelte component catalog
- Component composition patterns
- State management (Svelte stores + TanStack Query)
- Responsive design patterns
- Mobile-first approach
- Dark/light mode theming
- Design tokens (CSS variables)
- Component accessibility (ARIA, keyboard navigation)
- Animation & transitions (motion design)
- Icon system (SVG sprites)
- Typography scale
- Color system (semantic colors)
- Spacing & layout grid

**Research Sources**:
- SvelteKit 2 documentation
- shadcn-svelte component library
- Tailwind CSS 4 utilities
- Material Design 3 components
- Radix UI primitives (shadcn foundation)

**Time**: 4 days

---

### 5.3 Player Architecture & Media UX

**Files**:
- [ ] `docs/PLAYER_UX.md`
- [ ] `.github/instructions/video-player-patterns.instructions.md`
- [ ] `.github/instructions/audio-player-patterns.instructions.md`

**Content**:
- Shaka Player (DASH) integration
- hls.js (HLS) fallback patterns
- Web Audio API (gapless audio, crossfade)
- Custom controls design
- Subtitle/caption rendering (WebVTT)
- Quality switching UI
- Playback speed controls
- Picture-in-Picture (PiP)
- Chromecast integration UI
- Watch Party synchronization UI
- Progress tracking visualization
- Keyboard shortcuts (accessibility)
- Touch gestures (mobile)
- Scrubbing/seeking UX
- Buffer visualization
- Error state handling

**Research Sources**:
- Shaka Player documentation
- hls.js best practices
- Web Audio API guides
- YouTube player patterns (inspiration)
- Netflix player patterns (inspiration)

**Time**: 3 days

---

### 5.4 Accessibility & Internationalization

**Files**:
- [ ] `docs/ACCESSIBILITY_IMPLEMENTATION.md`
- [ ] `docs/I18N_STRATEGY.md`
- [ ] `.github/instructions/wcag-compliance.instructions.md`
- [ ] `.github/instructions/i18n-best-practices.instructions.md`

**Content**:
- WCAG 2.2 Level AA compliance checklist
- ARIA landmarks & roles
- Keyboard navigation patterns
- Screen reader optimization
- Focus management
- Color contrast requirements (4.5:1 text, 3:1 UI)
- Alt text guidelines
- Form accessibility
- Skip links
- Language detection
- Translation file structure (JSON)
- Date/time localization
- Number/currency formatting
- RTL (right-to-left) support
- Content translation workflow

**Research Sources**:
- W3C WAI guidelines
- ARIA Authoring Practices Guide (APG)
- axe-core accessibility testing
- i18next patterns

**Time**: 3 days

---

### 5.5 Media Processing & Advanced Features

**Files**:
- [ ] `.github/instructions/file-handling.instructions.md`
- [ ] `.github/instructions/image-processing.instructions.md`

**Content**:
- Multipart file uploads
- File validation
- Storage patterns
- Image resizing
- Blurhash generation
- Format conversion
- Cleanup strategies

**Time**: 2 days

---

### 5.6 Migration & Extension

**Files**:
- [ ] `docs/MIGRATION_GUIDE.md`
- [ ] `docs/MODULE_DEVELOPMENT.md`
- [ ] `docs/EXTERNAL_API_INTEGRATION.md`

**Content**:
- External source data migration
- Import to Revenge
- Watch history migration
- Metadata preservation
- Module structure guide
- Domain entity patterns
- Repository implementation
- Handler patterns
- Job definitions
- Testing modules
- API authentication
- Rate limiting
- Caching
- Error handling
- Fallback chains

**Time**: 6 days

---

**Phase 5 Total**: ~23 days (expanded from 8 days due to UX/UI research integration)

---

## Completion Checklist

### Documentation Files

**Critical (23 files)** (expanded from 14):
- [ ] `docs/JOB_QUEUE.md`
- [ ] `docs/SEARCH.md`
- [ ] `docs/SECURITY.md`
- [ ] `docs/DATABASE_SCHEMA.md`
- [ ] `docs/WEBSOCKET.md`
- [ ] `docs/MONITORING.md`
- [ ] `docs/BACKUP_RESTORE.md`
- [ ] `docs/PERFORMANCE_TUNING.md`
- [ ] `docs/UX_DESIGN_PRINCIPLES.md` ⭐ **NEW**
- [ ] `docs/FRONTEND_COMPONENTS.md` ⭐ **NEW**
- [ ] `docs/PLAYER_UX.md` ⭐ **NEW**
- [ ] `docs/ACCESSIBILITY_IMPLEMENTATION.md` ⭐ **NEW**
- [ ] `docs/I18N_STRATEGY.md` ⭐ **NEW**
- [ ] `docs/TESTING_STRATEGY.md`
- [ ] `docs/CICD.md`
- [ ] `docs/TROUBLESHOOTING.md`
- [ ] `docs/MIGRATION_GUIDE.md`
- [ ] `docs/MODULE_DEVELOPMENT.md`
- [ ] `docs/EXTERNAL_API_INTEGRATION.md`

**Integration Docs (41 files)**:
- [ ] 41× `docs/integrations/<SERVICE>.md`

### Instruction Files

**Critical (7 files)**:
- [ ] `.github/instructions/river-job-queue.instructions.md`
- [ ] `.github/instructions/typesense-integration.instructions.md`
- [ ] `.github/instructions/websocket-handlers.instructions.md`
- [ ] `.github/instructions/ogen-api-patterns.instructions.md`
- [ ] `.github/instructions/dragonfly-cache-patterns.instructions.md`
- [ ] `.github/instructions/security-best-practices.instructions.md`
- [ ] `.github/instructions/error-handling-patterns.instructions.md`

**High Priority (9 files)**:
- [ ] `.github/instructions/logging-standards.instructions.md`
- [ ] `.github/instructions/content-module-testing.instructions.md`
- [ ] `.github/instructions/api-pagination.instructions.md`
- [ ] `.github/instructions/rate-limiting.instructions.md`
- [ ] `.github/instructions/file-handling.instructions.md`
- [ ] `.github/instructions/image-processing.instructions.md`
- [ ] 41× `.github/instructions/<service>-client.instructions.md`

### External API Research

**Completed (11/41)**:
- [x] Radarr
- [x] Sonarr
- [x] Lidarr
- [x] TMDb
- [x] TheTVDB
- [x] MusicBrainz
- [x] Audiobookshelf
- [x] Trakt
- [x] Fanart.tv
- [x] Last.fm

**Remaining (28/39)**:
- [ ] Whisparr, Readarr
- [ ] OMDb, ThePosterDB, Spotify, Discogs, Goodreads, OpenLibrary, Audible, Hardcover, Stash/StashDB, ThePornDB
- [ ] ListenBrainz, Letterboxd, Simkl
- [ ] Authelia, Authentik, Keycloak
- [ ] Blackbeard, TVHeadend, NextPVR, Chromecast, DLNA
- [ ] OIDC Generic

---

## Success Criteria

### Before ANY Implementation Begins

- [ ] ✅ All 41 external APIs documented (live docs fetched)
- [ ] ✅ All 14 critical user docs written
- [ ] ✅ All 7 critical instruction files written
- [ ] ✅ All 41 integration guides created
- [ ] ✅ All 41 client instruction files created
- [ ] ✅ Database schema fully documented with ER diagrams
- [ ] ✅ Security architecture reviewed & approved
- [ ] ✅ API design patterns established
- [ ] ✅ Testing strategy defined

### Quality Gates

**Each documentation file must have**:
- [ ] Clear purpose statement
- [ ] Code examples (working, not pseudo-code)
- [ ] Diagrams where applicable
- [ ] Cross-references to related docs
- [ ] Troubleshooting section
- [ ] Last updated date

**Each instruction file must have**:
- [ ] `applyTo` paths defined
- [ ] DO/DON'T examples
- [ ] Copy-pasteable code snippets
- [ ] Anti-patterns listed
- [ ] Testing guidance

---

## Risk Mitigation

### Risks of Starting Implementation Too Early

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Breaking API changes mid-development | HIGH | 90% | Complete API research first |
| Inconsistent module patterns | HIGH | 80% | Establish patterns via instructions |
| Security vulnerabilities | CRITICAL | 70% | Document security architecture |
| Performance issues | MEDIUM | 60% | Database schema + tuning guide |
| Testing gaps | HIGH | 75% | Testing strategy document |

### Benefits of Complete Preparation

1. **No rework**: Patterns established once, used consistently
2. **Faster development**: Developers follow clear instructions
3. **Higher quality**: Standards enforced from day one
4. **Easier onboarding**: New contributors have complete docs
5. **Production-ready**: Deployment, monitoring, security covered upfront

---

## Timeline Summary

| Phase | Duration | Deliverables |
|-------|----------|--------------|
| **Phase 1: External APIs** | 2 weeks | 82 files (41 docs + 41 instructions) |
| **Phase 2: Core Infrastructure** | 2 weeks | 7 critical docs + 2 instruction files |
| **Phase 3: Development Standards** | 1 week | 3 docs + 2 instruction files |
| **Phase 4: Production Readiness** | 2 weeks | 5 docs + 2 instruction files |
| **Phase 5: UX/UI & Advanced Features** | 3 weeks | 8 docs + 11 instruction files (expanded from 1 week) |
| **TOTAL** | **10-11 weeks** | **133 files** (up from 104) |

**Note**: Phase 5 expanded significantly due to comprehensive UX/UI research integration (14 authoritative sources).

---

## Current Status (2026-01-28)

**Completed**:
- ✅ 11 external APIs fetched (27%)
- ✅ Architecture documentation (ARCHITECTURE_V2.md, etc.)
- ✅ Basic instruction files (fx, sqlc, koanf, etc.)

**In Progress**:
- 🟡 Fetching remaining 30 external APIs

**Not Started**:
- ❌ Core infrastructure docs (River, Typesense, Security, etc.)
- ❌ Integration guides (41 services)
- ❌ Client instruction files (41 services)
- ❌ Development standards
- ❌ Production readiness docs

---

## Next Actions

### Immediate (This Week)

1. **Continue fetching external APIs** (28 remaining)
   - Batch 1: Whisparr, Readarr
   - Batch 2: Metadata providers (OMDb, ThePosterDB, Spotify, etc.)
   - Batch 3: Scrobbling (ListenBrainz, Letterboxd, Simkl)
   - Batch 4: Auth/SSO (Authelia, Authentik, Keycloak)
   - Batch 5: Infrastructure (Blackbeard, TVHeadend, NextPVR, etc.)

2. **Update EXTERNAL_INTEGRATIONS_TODO.md** with fetched API details

### Week 2

1. **Start Phase 2** (Core Infrastructure Docs)
   - River Job Queue documentation
   - Typesense Search documentation
   - Database Schema with ER diagrams

### Week 3-8

1. Follow master plan phases
2. Track progress in this document
3. Update completion percentages

---

## Approval & Sign-Off

**Prepared By**: GitHub Copilot
**Reviewed By**: TBD
**Approved By**: TBD
**Date**: 2026-01-28

**Status**: ⚠️ AWAITING APPROVAL TO PROCEED

---

**THIS MASTER PLAN MUST BE COMPLETED BEFORE ANY CODING BEGINS!**

No shortcuts. No "we'll document it later". Complete preparation prevents catastrophic rework.

---

## Appendix

### A. UX/UI Design & Frontend Resources

**Extracted to**: [docs/research/UX_UI_RESOURCES.md](../research/UX_UI_RESOURCES.md)

Contains 14 authoritative sources:
- W3C WCAG 2.2, ISO 9241-11:2018
- Nielsen's 10 Heuristics, Laws of UX
- Material Design 3, Apple HIG, Fluent 2, Atlassian, Carbon
- UK GDS Principles
- IDF, Baymard, Smashing Magazine
- web.dev Patterns

---

**END OF IMPLEMENTATION TRACKER**
