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
| Overseerr | v1 | X-Api-Key OR Cookie | Self-hosted (none) | ✅ DONE |
| Fanart.tv | v3 | api-key header OR query | Not specified | ✅ DONE |
| Last.fm | v2 | API key | Reasonable usage | ✅ DONE |

#### ❌ To Fetch (30 services)

**Servarr Ecosystem (2)**:
- [ ] Whisparr (adult content management)
- [ ] Readarr (book management)

**Request Management (1)**:
- [ ] Jellyseerr (Jellyfin-focused)

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

**Other (3)**:
- [ ] OIDC Generic (OpenID Connect specification)
- [ ] Plex API (for Overseerr/Jellyseerr Plex auth)
- [ ] Jellyfin API (compatibility reference)

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
- gorilla/websocket documentation
- nhooyr.io/websocket (modern alternative)
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
- go-redis/v9 documentation

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
- Jellyfin data export
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
- [x] Overseerr
- [x] Fanart.tv
- [x] Last.fm

**Remaining (30/41)**:
- [ ] Whisparr, Readarr, Jellyseerr
- [ ] OMDb, ThePosterDB, Spotify, Discogs, Goodreads, OpenLibrary, Audible, Hardcover, Stash/StashDB, ThePornDB
- [ ] ListenBrainz, Letterboxd, Simkl
- [ ] Authelia, Authentik, Keycloak
- [ ] Blackbeard, TVHeadend, NextPVR, Chromecast, DLNA
- [ ] OIDC Generic, Plex API, Jellyfin API

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

1. **Continue fetching external APIs** (30 remaining)
   - Batch 1: Whisparr, Readarr, Jellyseerr
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

## Appendix A: UX/UI Design & Frontend Resources

> **Status**: ✅ ALL 14 SOURCES FETCHED (2026-01-28)
>
> These authoritative sources inform frontend instruction files, component design, and user experience patterns for Revenge.

### A.1 Accessibility Standards

#### W3C WCAG 2.2 (Web Content Accessibility Guidelines)
- **URL**: https://www.w3.org/WAI/standards-guidelines/wcag/
- **Version**: WCAG 2.2 (published October 2023, updated December 2024)
- **Status**: International Standard, ISO/IEC 40500:2025
- **Scope**: 13 guidelines under 4 principles (Perceivable, Operable, Understandable, Robust)
- **Conformance Levels**: A, AA, AAA (9 new success criteria in 2.2)
- **Key Changes**:
  - Added mobile accessibility
  - Cognitive accessibility improvements
  - 4.1.1 Parsing obsolete
  - Focus indicators, dragging movements, target size
- **Resources**:
  - Quick Reference: https://www.w3.org/WAI/WCAG22/quickref/
  - WCAG 2.2 Standard: https://www.w3.org/TR/WCAG22/
  - Understanding WCAG 2: Detailed guidance
  - Techniques for WCAG 2: Implementation patterns
  - Supplemental Guidance: Beyond baseline
- **Revenge Application**:
  - Level AA compliance target
  - Keyboard navigation for all controls
  - Screen reader optimization
  - Color contrast 4.5:1 (text), 3:1 (UI)
  - Focus visible indicators
  - Alt text for all images
  - Skip links for navigation
  - ARIA landmarks and roles

#### ISO 9241-11:2018 (Ergonomics of Human-System Interaction)
- **URL**: https://www.iso.org/standard/63500.html
- **Edition**: 2 (March 2018, confirmed 2023)
- **Scope**: Usability definitions and concepts
- **Key Concepts**:
  - Usability as outcome of use (not inherent property)
  - Effectiveness, efficiency, satisfaction in context
  - User needs drive design
  - System, product, service applicability
- **Revenge Application**:
  - User research validation
  - Usability testing metrics (task success, time, satisfaction)
  - Context of use analysis (device types, environments)

---

### A.2 Usability Heuristics & Laws

#### Nielsen Norman Group - 10 Usability Heuristics
- **URL**: https://www.nngroup.com/articles/ten-usability-heuristics/
- **Author**: Jakob Nielsen (1994, updated 2020)
- **Status**: Industry standard for heuristic evaluation
- **The 10 Heuristics**:
  1. **Visibility of System Status** - Feedback within reasonable time
  2. **Match System & Real World** - User language, natural mapping
  3. **User Control & Freedom** - Undo/redo, emergency exits
  4. **Consistency & Standards** - Platform conventions, Jakob's Law
  5. **Error Prevention** - Eliminate error-prone conditions, confirm destructive actions
  6. **Recognition > Recall** - Visible options, minimize memory load
  7. **Flexibility & Efficiency** - Shortcuts for experts, personalization
  8. **Aesthetic & Minimalist Design** - No irrelevant information
  9. **Error Recognition & Recovery** - Plain language errors, solutions
  10. **Help & Documentation** - Easy to search, contextual, concrete steps
- **Resources**:
  - Free posters (summary + 10 detailed): https://media.nngroup.com/media/articles/attachments/Jakob's10UsabilityHeuristics_AllPosters_5.zip
  - Video explanations (2-3 min each)
  - Application examples (complex apps, VR, video games)
- **Revenge Application**:
  - Heuristic evaluation of all UIs
  - Progress indicators for transcoding
  - Undo/cancel for destructive actions
  - Consistent navigation across modules
  - Error messages with solutions
  - Contextual help tooltips

#### Laws of UX
- **URL**: https://lawsofux.com/
- **Author**: Jon Yablonski (O'Reilly book)
- **Scope**: 26 laws/principles for UI design
- **Key Laws for Revenge**:
  - **Aesthetic-Usability Effect**: Beautiful design perceived as more usable
  - **Fitts's Law**: Target acquisition time = f(distance, size) → Larger touch targets
  - **Hick's Law**: Decision time increases with choices → Limit navigation options
  - **Jakob's Law**: Users expect site to work like others → Use conventions
  - **Miller's Law**: Working memory 7±2 items → Chunk information
  - **Serial Position Effect**: Remember first & last items → Key actions at ends
  - **Von Restorff Effect**: Distinctive items remembered → Highlight primary CTA
  - **Doherty Threshold**: <400ms response time boosts productivity
  - **Goal-Gradient Effect**: Motivation increases near goal → Show progress
  - **Peak-End Rule**: Experiences judged by peak & end → Optimize critical moments
  - **Tesler's Law**: Complexity conservation → Simplify user-facing, accept backend complexity
  - **Pareto Principle**: 80% effects from 20% causes → Focus on common tasks
  - **Postel's Law**: Be liberal in inputs, conservative in outputs
  - **Zeigarnik Effect**: Unfinished tasks remembered → Save drafts, show incomplete
  - **Choice Overload**: Too many options → paralysis
  - **Cognitive Load**: Mental resources to understand UI → Minimize
  - **Flow**: Immersed energized focus → Remove friction
- **Revenge Application**:
  - Large tap targets (48×48px minimum) for mobile
  - Limit main navigation to 5-7 items
  - Chunk settings into logical groups
  - Primary actions (Play, Add) visually prominent
  - Progress bars for uploads/transcoding
  - <400ms API response target
  - Autosave for playlist editing
  - Conservative validation (accept variations), strict output

---

### A.3 Design Systems & Component Libraries

#### Google Material Design 3 (M3 Expressive)
- **URL**: https://m3.material.io/
- **Status**: Latest evolution (I/O 2025 update)
- **Philosophy**: Emotion-driven UX with vibrant colors, intuitive motion, adaptive components
- **Key Features**:
  - **M3 Expressive update** (2025):
    - Vibrant color system (extended palettes)
    - Motion physics (easier-to-implement, token-powered transitions)
    - Shape library (35 shapes with built-in morph motion)
    - Flexible typography
  - **New Components**:
    - Toolbars (flexible action containers, pairs with FAB)
    - Split buttons (button + menu, expressive shape/motion)
    - Progress indicators (waveform, customizable thickness)
    - Button groups (shape-shifting, reactive buttons)
  - **Updated Components**: 14 total (existing components refreshed)
- **Libraries**:
  - Web (Material Web Components)
  - Android (Compose, MDC)
  - Flutter
  - Figma UI Kit (latest M3 Expressive)
- **Resources**:
  - Blog: Building with M3 Expressive
  - Motion physics guide
  - Figma plugin
- **Revenge Application**:
  - Inspiration for component animations (FAB transitions, button states)
  - Motion design for quality switching
  - Progress indicators for transcoding
  - Color palette generation
  - **NOT using Material directly** (using shadcn-svelte), but adopting motion principles

#### Apple Human Interface Guidelines (HIG)
- **URL**: https://developer.apple.com/design/human-interface-guidelines/
- **Platforms**: iOS, iPadOS, macOS, watchOS, tvOS, visionOS
- **Core Principles**:
  - **Hierarchy**: Controls elevate content, clear visual hierarchy
  - **Harmony**: Align with concentric hardware/software design
  - **Consistency**: Platform conventions, adapt across displays
- **Design Fundamentals**:
  - App icons, color, materials, layout
  - Typography, icons, accessibility
  - Generative AI patterns (NEW)
- **Topics**: Foundations, Patterns, Components, Inputs, Technologies
- **Revenge Application**:
  - iOS/macOS native app patterns (future)
  - Touch gesture conventions
  - System integration (PiP, AirPlay)
  - SF Symbols icon style inspiration

#### Microsoft Fluent Design System (Fluent 2)
- **URL**: https://fluent2.microsoft.design/
- **Platforms**: Web (React), iOS, Android, Windows
- **Philosophy**: Let creativity flow, accessible & inclusive
- **Components**: Web (React, extensive catalog), mobile (iOS/Android), Windows (WinUI 3)
- **Resources**:
  - Figma UI Kit (employee access)
  - Component documentation per platform
- **Revenge Application**:
  - Cross-platform component patterns
  - Windows native app (future)
  - **NOT using Fluent directly**, but reference for Windows UX conventions

#### Atlassian Design System
- **URL**: https://atlassian.design/
- **Philosophy**: Better teamwork by design, unified design language across all apps
- **Key Features**:
  - Rovo AI patterns (AI integration UX)
  - Unified design language across Jira, Confluence, Trello
  - Foundations: Color, Typography, Iconography, Grid, Accessibility, Tokens
  - Release phases (alpha, beta, stable)
  - Contact/feedback form
- **Libraries**: Atlaskit (React), Forge UI Kit
- **Revenge Application**:
  - AI chat patterns (Rovo) for future AI features
  - Token-based theming
  - Component composition patterns

#### IBM Carbon Design System
- **URL**: https://www.carbondesignsystem.com/
- **Philosophy**: Adaptable system, best practices of UI design, open-source
- **Commitment**: Web Components (September 2024)
- **Libraries**:
  - Web Components
  - React (primary)
  - Angular, Vue, Svelte
- **Key Features**:
  - AI Chat v1 (October 2025)
  - Comprehensive component catalog
  - Accessibility-first design
  - Design tokens
- **Resources**:
  - Figma kit
  - Medium blog
  - Community contribution guidelines
- **Revenge Application**:
  - Svelte component patterns
  - AI chat interface (future)
  - Design token structure

---

### A.4 Government & Institutional Standards

#### UK Government Design Principles (GDS)
- **URL**: https://www.gov.uk/guidance/government-design-principles
- **Published**: April 2012, updated April 2025
- **The 11 Principles**:
  1. **Start with user needs** - Research, don't assume
  2. **Do less** - Reusable platforms, link to others
  3. **Design with data** - Analytics built-in, always on
  4. **Do hard work to make it simple** - Simplicity > "always been that way"
  5. **Iterate. Then iterate again** - MVP, alpha→beta→live, learn from failures
  6. **This is for everyone** - Accessible design = good design
  7. **Understand context** - Not designing for screen, for people (library? phone? Facebook-only?)
  8. **Build digital services, not websites** - Connect real world (beyond UI)
  9. **Be consistent, not uniform** - Shared language/patterns, but improve when needed
  10. **Make things open** - Share code, designs, ideas, failures
  11. **Minimise environmental impact** - Reduce energy, water, materials (NEW 2025)
- **Resources**:
  - GOV.UK Design System
  - Poster download (GitHub)
- **Revenge Application**:
  - User research-driven features
  - Iterative development (alpha modules first)
  - Accessibility priority
  - Open-source ethos
  - Reusable module patterns
  - Energy-efficient transcoding (Blackbeard optimization)

---

### A.5 UX Research & Best Practices

#### Interaction Design Foundation (IDF)
- **URL**: https://www.interaction-design.org/literature/topics/ux-design
- **Scope**: World's largest UX education community (1.2M+ enrollments)
- **Key Content**:
  - User Experience definition (Don Norman)
  - UX vs UI distinction
  - ISO 9241-210 (human-centered design)
  - Who/Why/What/How framework
  - User-centered design process (iterative)
  - Multidisciplinary field (psychology, interaction design, IA, research)
  - Typical tasks: User research, personas, wireframes, prototypes, testing
- **Courses**:
  - User Experience: The Beginner's Guide
  - User Research - Methods and Best Practices
  - Get Your First Job as a UX Designer
- **Revenge Application**:
  - User research methodology
  - Persona creation for target users (media enthusiasts, families, power users)
  - Iterative design process
  - Usability testing protocols

#### Baymard Institute (E-commerce UX Research)
- **URL**: https://baymard.com/blog
- **Scope**: 200,000+ hours of UX research, 18,000+ design examples
- **Focus**: E-commerce usability (applicable to media libraries)
- **Research Topics**:
  - Homepage & category navigation (11 best practices, 67% sites mediocre)
  - Product list UX
  - Product page optimization
  - Checkout UX (10 pitfalls)
  - Mobile UX trends 2025 (9 common pitfalls)
  - Accessibility
  - Search & filtering
- **Resources**:
  - UX Benchmark (326 top sites, 52,000+ scores)
  - Design examples (annotated)
  - Cart abandonment stats
  - Credit card patterns
- **Revenge Application**:
  - Library navigation patterns
  - Filtering & sorting best practices
  - Mobile UX optimization
  - Search interface design
  - Checkout flow → "Add to Watchlist" flow

#### Smashing Magazine (UX Design)
- **URL**: https://www.smashingmagazine.com/category/ux-design/
- **Scope**: Professional web design & UX articles (57+ UX design articles)
- **Key Topics**:
  - Design patterns for design systems
  - AI in design (skills AI can't replicate)
  - Infinite scroll best practices
  - Rapid research programs
  - Accessibility (global developments during COVID)
  - B2B UX design
  - Timing in design
- **Revenge Application**:
  - Infinite scroll for media libraries
  - Design system collaboration patterns
  - Research-driven feature development
  - Accessibility best practices

---

### A.6 Web Standards & Patterns

#### web.dev Patterns (Google)
- **URL**: https://web.dev/patterns/
- **Scope**: Modern web API patterns with browser support (Baseline)
- **Pattern Categories**:
  - **Animation**: CSS/JS animations with accessibility, user preferences (prefers-reduced-motion)
  - **Clipboard**: Copy/paste patterns
  - **Components**: Cross-browser UI components, design system inspiration
  - **Files & Directories**: File upload, drag-drop, directory access
  - **Layout**: Modern CSS (Grid, Flexbox, Container Queries)
    - Cards, dynamic grids, full-page layouts
  - **Media**: Video, audio, images (lazy loading, responsive)
  - **Theming**: Color management, dark mode, CSS custom properties
  - **Web Apps**: PWA patterns (service workers, manifest, offline)
- **Revenge Application**:
  - Responsive layout patterns (media grids)
  - Dark/light mode theming
  - File upload for custom artwork
  - Video player controls
  - PWA offline support
  - Animation performance (GPU-accelerated)
  - Clipboard API for sharing

---

### A.7 Summary: UX/UI Resources Application to Revenge

| Resource | Primary Use in Revenge |
|----------|------------------------|
| **WCAG 2.2** | Accessibility compliance (AA level), keyboard navigation, screen readers |
| **ISO 9241-11** | Usability testing framework, effectiveness/efficiency metrics |
| **Nielsen's Heuristics** | Heuristic evaluation of all UIs, error handling, consistency |
| **Laws of UX** | Touch target sizing, navigation chunking, progress visualization, cognitive load reduction |
| **Material Design 3** | Motion design inspiration, progress indicators, animation principles |
| **Apple HIG** | iOS/macOS app patterns (future), touch gestures, system integration |
| **Fluent 2** | Cross-platform component patterns, Windows UX (future) |
| **Atlassian** | AI chat patterns, token-based theming, component composition |
| **Carbon** | Svelte component patterns, AI features, design tokens |
| **GDS Principles** | User research-driven development, iterative design, accessibility priority, open-source |
| **IDF** | User research methodology, persona creation, iterative design process |
| **Baymard** | Library navigation, filtering/sorting, mobile UX, search interface |
| **Smashing** | Infinite scroll, design patterns, accessibility, research-driven features |
| **web.dev** | Responsive layouts, theming, PWA, file handling, animation performance |

### A.8 Frontend Instruction Files Informed by These Sources

The following instruction files will be created using patterns from these UX/UI resources:

1. **`.github/instructions/frontend-ux-guidelines.instructions.md`**
   - Nielsen's 10 Heuristics implementation
   - Laws of UX checklist (Fitts, Hick, Miller, Jakob)
   - GDS principles application

2. **`.github/instructions/accessibility-standards.instructions.md`**
   - WCAG 2.2 Level AA compliance checklist
   - ARIA patterns
   - Keyboard navigation requirements
   - Screen reader optimization
   - Color contrast tools

3. **`.github/instructions/svelte-component-patterns.instructions.md`**
   - shadcn-svelte customization
   - Carbon/Atlassian component patterns
   - Composition over inheritance
   - Accessibility built-in

4. **`.github/instructions/shadcn-svelte-customization.instructions.md`**
   - Tailwind CSS 4 utilities
   - Design token structure
   - Theming system
   - Component variants

5. **`.github/instructions/video-player-patterns.instructions.md`**
   - Shaka Player integration
   - hls.js fallback
   - Custom controls (Baymard patterns)
   - Keyboard shortcuts (Apple HIG)
   - Touch gestures

6. **`.github/instructions/audio-player-patterns.instructions.md`**
   - Web Audio API (gapless, crossfade)
   - Material Design 3 progress indicators
   - Playback visualization

7. **`.github/instructions/wcag-compliance.instructions.md`**
   - WCAG 2.2 audit checklist
   - axe-core integration
   - Testing workflow
   - Remediation patterns

8. **`.github/instructions/i18n-best-practices.instructions.md`**
   - Translation workflow
   - Date/time localization
   - RTL support
   - Content translation strategy

---

### A.9 Quality Gates for Frontend Development

Before ANY frontend component is merged:

- [ ] ✅ WCAG 2.2 Level AA compliant (axe-core passes)
- [ ] ✅ Keyboard navigable (tab order, focus visible)
- [ ] ✅ Screen reader tested (NVDA/JAWS/VoiceOver)
- [ ] ✅ Touch targets ≥48×48px (Fitts's Law)
- [ ] ✅ Color contrast ≥4.5:1 text, ≥3:1 UI
- [ ] ✅ Dark/light mode support
- [ ] ✅ Mobile-first responsive
- [ ] ✅ `prefers-reduced-motion` respected
- [ ] ✅ Error messages with solutions (Nielsen #9)
- [ ] ✅ Undo/cancel for destructive actions (Nielsen #3)
- [ ] ✅ Consistent with design system
- [ ] ✅ i18n keys used (no hardcoded strings)
- [ ] ✅ Documented in Storybook (component catalog)

---

**END OF APPENDIX A**
