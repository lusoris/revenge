# MASTER TASK LIST - Complete Implementation Plan

**Date**: 2026-01-31
**Purpose**: Comprehensive collection of ALL tasks with dependencies and optimal sequence
**Target**: 99% perfection before coding begins

---

## OVERVIEW

**Total Phases**: 9 (Phase 0 complete + 8 remaining)
**Current Focus**: Phase 1 (Design Completion)

**Key Dependencies**:
1. Design docs MUST be complete before coding roadmap
2. Coding roadmap MUST be complete before templates/automation
3. Templates/automation before migration
4. Migration before validation

---

## PHASE 0: Planning & Approval ✅ COMPLETE

- [x] Analyze current documentation structure
- [x] Create comprehensive questions
- [x] Get all P0/P1 answers from user
- [x] Conduct gap analysis
- [x] Document decisions
- [x] Add external sources to SOURCES.yaml (Jinja2, GitHub Wiki, Go AST, markdownlint)
- [x] Start source fetching (fetch-sources.py)

---

## PHASE 1: Design Completion (CURRENT - MOST IMPORTANT)

**Goal**: Achieve 99% design perfection for ALL v1.0 features

### 1.1: Finalize Gap Analysis ⏸️
- [x] Create `.analysis/09_DESIGN_GAPS.md`
- [ ] Verify source fetch completed successfully
- [ ] Review gap analysis with user
- [ ] Confirm priority order

### 1.2: Create Scaffold Infrastructure
- [ ] Create `docs/dev/design/.templates/SCAFFOLD_TEMPLATE.md`
- [ ] Test scaffold template with one example
- [ ] Validate template structure

### 1.3: Scaffold Missing Documents
**For each missing/incomplete doc identified in gap analysis**:

#### Critical Technical Docs (HIGHEST PRIORITY)
- [ ] Scaffold API.md (if not already detailed enough)
- [ ] Scaffold FRONTEND.md (if not already detailed enough)

#### Content Modules (HIGH PRIORITY)
- [ ] Scaffold/complete MUSIC_MODULE.md (currently 🟡)
- [ ] Scaffold/complete AUDIOBOOK_MODULE.md (currently 🟡)
- [ ] Scaffold/complete BOOK_MODULE.md (currently 🟡)

#### Wiki Integrations (MEDIUM PRIORITY)
- [ ] Scaffold WIKIPEDIA.md complete design
- [ ] Scaffold FANDOM.md complete design
- [ ] Scaffold TVTROPES.md complete design
- [ ] Scaffold BABEPEDIA.md complete design (adult)
- [ ] Scaffold BOOBPEDIA.md complete design (adult)
- [ ] Scaffold IAFD.md complete design (adult)

#### Scrobbling Services (MEDIUM PRIORITY)
- [ ] Scaffold/complete TRAKT.md (20% → 100%)
- [ ] Scaffold/complete LETTERBOXD.md (20% → 100%)
- [ ] Scaffold/complete LISTENBRAINZ.md (20% → 100%)
- [ ] Scaffold/complete SIMKL.md (20% → 100%)

#### Auth Integrations (HIGH PRIORITY)
- [ ] Complete AUTHENTIK.md (50% → 100%)
- [ ] Complete KEYCLOAK.md (50% → 100%)

#### Adult Metadata (MEDIUM PRIORITY)
- [ ] Complete STASH.md (50% → 100%)
- [ ] Complete THEPORNDB.md (50% → 100%)
- [ ] Complete PORNHUB.md (50% → 100%)
- [ ] Complete THENUDE.md (50% → 100%)

#### Technical Documentation (HIGH PRIORITY)
- [ ] Complete NOTIFICATIONS.md
- [ ] Complete EMAIL.md
- [ ] Complete WEBHOOKS.md

#### Operations Documentation (MEDIUM PRIORITY)
- [ ] Complete SETUP.md
- [ ] Complete DEVELOPMENT.md
- [ ] Complete GITFLOW.md
- [ ] Complete BRANCH_PROTECTION.md

### 1.4: Design Writing (Iterative - by Priority)

#### Tier 1: FOUNDATION (DO FIRST - BLOCKS EVERYTHING)
**Estimated**: 40-60 hours

1. **API.md** - Complete endpoint specifications
   - [ ] Research: Review ogen docs, existing service designs
   - [ ] Architecture: Define API structure, versioning, namespaces
   - [ ] Database: N/A (endpoints don't have their own DB)
   - [ ] API Endpoints: ALL endpoints for ALL services (this IS the doc)
     - [ ] Auth endpoints (login, logout, refresh, etc.)
     - [ ] User management endpoints
     - [ ] Library endpoints
     - [ ] Metadata endpoints
     - [ ] Search endpoints
     - [ ] Playback endpoints
     - [ ] Settings endpoints
     - [ ] All content module endpoints (movies, TV, music, etc.)
   - [ ] Request/Response schemas: Define JSON structures
   - [ ] Error handling: Define error codes, formats
   - [ ] Rate limiting: Define policies
   - [ ] Authentication: Define auth flow, token handling
   - [ ] Pagination: Define patterns
   - [ ] Filtering/Sorting: Define query parameters
   - [ ] Testing Strategy: API testing approach
   - [ ] Review and validate
   - [ ] Mark as ✅

2. **FRONTEND.md** - Complete component architecture
   - [ ] Research: Review SvelteKit 2, Svelte 5 runes, shadcn-svelte
   - [ ] Architecture: Define component structure, routing, layouts
   - [ ] Component Patterns:
     - [ ] shadcn-svelte usage patterns
     - [ ] Custom component guidelines
     - [ ] Composition patterns
   - [ ] State Management:
     - [ ] Svelte stores patterns
     - [ ] TanStack Query integration
     - [ ] Local vs server state
   - [ ] API Integration:
     - [ ] fetch vs query patterns
     - [ ] Error handling
     - [ ] Loading states
   - [ ] Route Structure:
     - [ ] Page organization
     - [ ] Navigation patterns
     - [ ] Protected routes
   - [ ] Testing Strategies:
     - [ ] Vitest unit tests
     - [ ] Playwright E2E tests
     - [ ] Component testing
   - [ ] Performance:
     - [ ] Code splitting
     - [ ] Lazy loading
     - [ ] SSR vs CSR decisions
   - [ ] Accessibility:
     - [ ] ARIA patterns
     - [ ] Keyboard navigation
     - [ ] Screen reader support
   - [ ] Theming:
     - [ ] Light/dark mode
     - [ ] CSS custom properties
     - [ ] Tailwind config
   - [ ] Forms:
     - [ ] Validation patterns
     - [ ] Error display
     - [ ] Superforms integration
   - [ ] Review and validate
   - [ ] Mark as ✅

#### Tier 2: CONTENT MODULES (MVP-CRITICAL)
**Estimated**: 30-40 hours

3. **MUSIC_MODULE.md** (🟡 → ✅)
   - [ ] Research: Check Lidarr, MusicBrainz, Spotify APIs
   - [ ] Architecture: Artists, albums, tracks, playlists
   - [ ] Database Schema:
     - [ ] artists table
     - [ ] albums table
     - [ ] tracks table
     - [ ] genres table
     - [ ] playlists table
     - [ ] artist_album relationships
     - [ ] album_track relationships
     - [ ] ER diagram
   - [ ] API Endpoints: Music CRUD, playlist management
   - [ ] Integrations: Lidarr, MusicBrainz, Spotify, Discogs, Last.fm
   - [ ] Metadata Priority Chain
   - [ ] Audio Fingerprinting (chromaprint)
   - [ ] Scrobbling Integration
   - [ ] Smart Playlists/Collections
   - [ ] Testing Strategy
   - [ ] Review and validate
   - [ ] Mark as ✅

4. **AUDIOBOOK_MODULE.md** (🟡 → ✅)
   - [ ] Research: Check Chaptarr, Audnexus, OpenLibrary
   - [ ] Architecture: Books, authors, narrators, chapters
   - [ ] Database Schema:
     - [ ] audiobooks table
     - [ ] authors table
     - [ ] narrators table
     - [ ] chapters table
     - [ ] bookmarks table
     - [ ] ER diagram
   - [ ] API Endpoints: Audiobook CRUD, chapter/bookmark management
   - [ ] Integrations: Chaptarr, Audnexus, OpenLibrary
   - [ ] Chapter Management
   - [ ] Bookmark System
   - [ ] Playback Position Tracking
   - [ ] Testing Strategy
   - [ ] Review and validate
   - [ ] Mark as ✅

5. **BOOK_MODULE.md** (🟡 → ✅)
   - [ ] Research: Check Chaptarr, OpenLibrary, Hardcover, Goodreads
   - [ ] Architecture: Books, authors, publishers, series
   - [ ] Database Schema:
     - [ ] books table
     - [ ] authors table
     - [ ] publishers table
     - [ ] series table
     - [ ] ER diagram
   - [ ] API Endpoints: Book CRUD, reading progress
   - [ ] Integrations: Chaptarr, OpenLibrary, Hardcover, Goodreads
   - [ ] Format Support (EPUB, PDF, MOBI, AZW3)
   - [ ] Reading Progress Tracking
   - [ ] Annotation/Highlighting System
   - [ ] Testing Strategy
   - [ ] Review and validate
   - [ ] Mark as ✅

#### Tier 3: HIGH-PRIORITY INTEGRATIONS
**Estimated**: 40-50 hours

6. **WIKIPEDIA.md** (🟡 → ✅)
   - [ ] Research: MediaWiki API
   - [ ] Complete design following integration template
   - [ ] API integration patterns
   - [ ] Content mapping (articles to media items)
   - [ ] Caching strategy
   - [ ] Review and validate
   - [ ] Mark as ✅

7. **AUTHENTIK.md** (50% → 100%)
   - [ ] Complete Group → Role mapping specifications
   - [ ] Property mapping examples
   - [ ] Outpost proxy setup details
   - [ ] Complete implementation checklist
   - [ ] Review and validate
   - [ ] Mark as ✅

8. **KEYCLOAK.md** (50% → 100%)
   - [ ] Complete OIDC configuration details
   - [ ] Realm and client setup specifications
   - [ ] Role/group mapping strategies
   - [ ] Implementation examples
   - [ ] Review and validate
   - [ ] Mark as ✅

9. **TRAKT.md** (20% → 100%)
   - [ ] Research: Trakt API v2
   - [ ] Complete API design
   - [ ] OAuth flow
   - [ ] Scrobbling patterns
   - [ ] Rating/watchlist sync
   - [ ] Review and validate
   - [ ] Mark as ✅

10. **NOTIFICATIONS.md** (🟡 → ✅)
    - [ ] Complete notification type specifications
    - [ ] Channel routing rules (email, push, websocket)
    - [ ] Delivery guarantees
    - [ ] User preference system
    - [ ] WebSocket integration details
    - [ ] Review and validate
    - [ ] Mark as ✅

11. **SETUP.md** (🟡 → ✅)
    - [ ] Complete installation guide
    - [ ] Step-by-step deployment
    - [ ] Configuration examples
    - [ ] Troubleshooting section
    - [ ] Review and validate
    - [ ] Mark as ✅

#### Tier 4: MEDIUM-PRIORITY FEATURES
**Estimated**: 40-50 hours

12. **LETTERBOXD.md** (20% → 100%)
13. **LISTENBRAINZ.md** (20% → 100%)
14. **SIMKL.md** (20% → 100%)
15. **FANDOM.md** (🟡 → ✅)
16. **TVTROPES.md** (🟡 → ✅)
17. **BABEPEDIA.md** (🟡 → ✅)
18. **BOOBPEDIA.md** (🟡 → ✅)
19. **IAFD.md** (🟡 → ✅)
20. **STASH.md** (50% → 100%)
21. **THEPORNDB.md** (50% → 100%)
22. **PORNHUB.md** (50% → 100%)
23. **THENUDE.md** (50% → 100%)
24. **EMAIL.md** (🟡 → ✅)
25. **WEBHOOKS.md** (🟡 → ✅)

#### Tier 5: POLISH (v1.0)
**Estimated**: 20-30 hours

26. **DEVELOPMENT.md** (🟡 → ✅)
27. **GITFLOW.md** (🟡 → ✅)
28. **BRANCH_PROTECTION.md** (🟡 → ✅)
29. Add source references to ALL docs
30. Add implementation checklists to ALL docs
31. Verify all version numbers against SOURCE_OF_TRUTH
32. Final review of all designs

### 1.5: Create Progress Tracking
- [ ] Create `.analysis/10_DESIGN_PROGRESS.md`
- [ ] Set up weekly goals tracking
- [ ] Track blockers and questions

### 1.6: Design Validation
- [ ] All designs follow template structure
- [ ] All designs have architecture diagrams
- [ ] All designs link to external sources
- [ ] All cross-references are valid
- [ ] Status tables updated
- [ ] Run markdown linter on all design docs
- [ ] User review and approval

**Exit Criteria**:
- ✅ Zero "🔴 PLANNED" status
- ✅ All v1.0 features have ✅ or 🟡 design docs
- ✅ API.md and FRONTEND.md are ✅ (foundational)
- ✅ All 3 content modules ✅ (music, audiobook, book)
- ✅ User explicitly approves design completeness

**Estimated Total**: 170-230 hours collaborative work

---

## PHASE 2: MVP Definition & Coding Roadmap

**Dependencies**: Phase 1 complete (all designs at 99%)

### 2.1: Create MVP Definition
- [ ] Create `planning/MVP_DEFINITION.md`
- [ ] Define v0.3.x scope clearly
- [ ] List included features
- [ ] List explicitly excluded features
- [ ] Define success criteria

### 2.2: Create Implementation Roadmap
- [ ] Create `planning/IMPLEMENTATION_ROADMAP.md`
- [ ] Define milestone structure (v0.1.x → v1.0.0)
- [ ] Create dependency graph between milestones
- [ ] Define exit criteria for each milestone

### 2.3: Create Detailed Phase TODOs
**For each milestone** (v0.1.x, v0.2.x, v0.3.x, v0.4.x, v0.5.x, v0.6-0.9.x, v1.0.0):

- [ ] Create `planning/v0.1.x_TODO.md` (100+ tasks)
  - Infrastructure (PostgreSQL, Dragonfly, River)
  - Core Services (Auth, User, Session, RBAC, Library)
  - API Layer (ogen setup, endpoints)
  - Observability (slog, Prometheus, OpenTelemetry)
  - Testing (unit, integration)
  - DevOps (Docker, CI/CD)

- [ ] Create `planning/v0.2.x_TODO.md` (100+ tasks)
  - First content type (Movie OR TV)
  - Metadata service
  - First integration (TMDB OR TheTVDB)
  - First *arr integration (Radarr OR Sonarr)
  - Basic SvelteKit UI

- [ ] Create `planning/v0.3.x_TODO.md` (100+ tasks - MVP)
  - Second content type (TV OR Movie)
  - Complete metadata service
  - Both integrations (TMDB + TheTVDB)
  - Both *arr (Radarr + Sonarr)
  - Full SvelteKit UI
  - Direct play + HLS/DASH
  - Search (Typesense)
  - OIDC auth

- [ ] Create `planning/v0.4.x_TODO.md` - Music module
- [ ] Create `planning/v0.5.x_TODO.md` - Transcoding (Blackbeard)
- [ ] Create `planning/v0.6.x_TODO.md` - Audiobooks + Books
- [ ] Create `planning/v0.7.x_TODO.md` - Photos + Comics
- [ ] Create `planning/v0.8.x_TODO.md` - LiveTV/DVR
- [ ] Create `planning/v0.9.x_TODO.md` - QAR/Adult + Advanced Playback
- [ ] Create `planning/v1.0.0_TODO.md` - Polish and release prep

### 2.4: Update SOURCE_OF_TRUTH
- [ ] Add "Planning & Roadmap" section to SOT
- [ ] Link to MVP_DEFINITION.md
- [ ] Link to IMPLEMENTATION_ROADMAP.md
- [ ] Link to all phase TODO files
- [ ] Add design completion status
- [ ] Link to gap analysis and progress tracker

**Exit Criteria**:
- ✅ MVP clearly defined
- ✅ Roadmap has 8+ milestones
- ✅ Each milestone has 100+ detailed tasks
- ✅ All linked from SOURCE_OF_TRUTH
- ✅ User approves roadmap

**Estimated Duration**: 3-5 days

---

## PHASE 3: Template System

**Dependencies**: Phase 2 complete, Jinja2 docs fetched

### 3.1: Design Template Structure
- [ ] Research Jinja2 best practices
- [ ] Define variable system (from SOURCE_OF_TRUTH)
- [ ] Define conditional blocks ({{ if claude }}, {{ if wiki }})
- [ ] Define common macros/includes

### 3.2: Create Claude Document Template
- [ ] Create `docs/dev/design/.templates/CLAUDE_DOC_TEMPLATE.template.md`
- [ ] Include: file paths, DB schema, code patterns, internal APIs
- [ ] Exclude: screenshots, user tutorials, step-by-step guides
- [ ] Add conditionals for Claude-specific content

### 3.3: Create Wiki Document Template
- [ ] Create `docs/dev/design/.templates/WIKI_DOC_TEMPLATE.template.md`
- [ ] Include: screenshots, tutorials, user perspective, external APIs
- [ ] Exclude: internal file paths, implementation details
- [ ] Add conditionals for Wiki-specific content

### 3.4: Test Templates
- [ ] Convert 2-3 sample docs to templates
- [ ] Generate both Claude and Wiki versions
- [ ] Validate output quality
- [ ] Lint generated docs
- [ ] User review and approval

**Exit Criteria**:
- ✅ Templates render without errors
- ✅ Output matches quality standards
- ✅ Variables substitute correctly
- ✅ Conditionals work as expected
- ✅ User approves template approach

**Estimated Duration**: 2-3 days

---

## PHASE 4: Automation Scripts

**Dependencies**: Phase 3 complete, Go AST docs fetched, GitHub Wiki docs fetched

### 4.1: Wiki Generation Script
- [ ] Create `scripts/generate-wiki.py`
- [ ] Integrate with Jinja2 templates
- [ ] Parse template variables from SOURCE_OF_TRUTH
- [ ] Generate wiki markdown for all docs
- [ ] Optimize for speed (< 30s for 179 docs)
- [ ] Test with all design docs

### 4.2: GitHub Wiki Sync
- [ ] Research GitHub Wiki API/Git interface
- [ ] Implement sync mechanism (git push to wiki repo)
- [ ] Add to generate-wiki.py or separate script
- [ ] Test sync to actual GitHub Wiki
- [ ] Add error handling and retry logic

### 4.3: Settings Sync Script
- [ ] Create `scripts/sync-tool-settings.py`
- [ ] Read SOURCE_OF_TRUTH.md for:
  - Go version
  - Python version
  - Node version
  - Formatter configs (gofmt, ruff, prettier)
  - Linter configs (golangci-lint, eslint)
  - LSP settings (gopls, typescript-language-server)
- [ ] Update all tool configs:
  - `.vscode/settings.json`
  - `.zed/settings.json`
  - `.jetbrains/` configs
  - `.coder/` configs
- [ ] Add validation mode (check without updating)
- [ ] Test on all config files

### 4.4: Code Status Verification Script
- [ ] Create `scripts/verify-code-status.py`
- [ ] Use Go AST (go/ast, go/parser, golang.org/x/tools/go/packages)
- [ ] Parse Go codebase to detect:
  - Implemented services (check for service.go files)
  - Implemented handlers (check for handler.go files)
  - Database tables (check migrations)
  - API endpoints (check OpenAPI spec + ogen generated code)
- [ ] Compare with design doc status
- [ ] Generate status report
- [ ] Update status tables in design docs automatically
- [ ] Add "last verified" timestamps

### 4.5: Integration Testing
- [ ] Test wiki generation end-to-end
- [ ] Test settings sync end-to-end
- [ ] Test code verification end-to-end
- [ ] Test all scripts together
- [ ] Fix any integration issues

**Exit Criteria**:
- ✅ All 3 scripts working correctly
- ✅ Wiki generation < 30s
- ✅ Settings sync 100% accurate
- ✅ Code status verification accurate
- ✅ Full test coverage on scripts

**Estimated Duration**: 3-4 days

---

## PHASE 5: Pipeline Integration

**Dependencies**: Phase 4 complete

### 5.1: Extend Doc Pipeline
- [ ] Review existing `scripts/doc-pipeline.sh`
- [ ] Add Stage 7: Wiki Generation
  - Call generate-wiki.py
  - Sync to GitHub Wiki
  - Validate wiki output
- [ ] Modify Stage 3 (status generation):
  - Add MVP filtering
  - Add code status verification
  - Update status tables with latest info
- [ ] Update pipeline runner script
- [ ] Test full pipeline

### 5.2: Extend Source Pipeline
- [ ] Review existing `scripts/source-pipeline.sh`
- [ ] Verify new sources included (Jinja2, GitHub Wiki, Go AST, markdownlint)
- [ ] Test source fetching
- [ ] Verify output locations

### 5.3: Create/Update CI/CD Workflows
- [ ] Create `.github/workflows/wiki-sync.yml`
  - Trigger on design doc changes
  - Run wiki generation
  - Sync to GitHub Wiki
  - Notify on failure

- [ ] Create `.github/workflows/settings-sync-validation.yml`
  - Trigger on SOURCE_OF_TRUTH changes
  - Run settings sync in validation mode
  - Fail if configs don't match SOT
  - Auto-fix option (create PR)

- [ ] Update `.github/workflows/doc-validation.yml`
  - Add wiki validation
  - Add code status check

- [ ] Update `.github/workflows/fetch-sources.yml`
  - Include new sources (Jinja2, etc.)

### 5.4: End-to-End Testing
- [ ] Test: Design doc change → wiki generation → GitHub Wiki sync
- [ ] Test: SOT change → settings validation → failure on mismatch
- [ ] Test: Code change → status update → design doc update
- [ ] Test: Source update → fetch → docs regenerate
- [ ] Fix any issues

**Exit Criteria**:
- ✅ Doc pipeline extended (7 stages)
- ✅ Source pipeline includes new sources
- ✅ 2 new CI/CD workflows working
- ✅ Full end-to-end tests passing
- ✅ Wiki syncs automatically on commit

**Estimated Duration**: 2-3 days

---

## PHASE 6: Documentation Migration

**Dependencies**: Phase 5 complete

### 6.1: Convert Design Docs to Templates
**Priority order** (do in batches):

1. **Feature docs** (41 docs)
   - Start with video/MOVIE_MODULE.md, video/TVSHOW_MODULE.md
   - Then music, audiobooks, books, etc.

2. **Service docs** (15 docs)
   - AUTH.md, USER.md, SESSION.md, etc.

3. **Integration docs** (65 docs)
   - Metadata providers
   - Auth providers
   - Arr stack
   - Scrobbling
   - Wiki integrations

4. **Operations docs** (7 docs)
   - SETUP.md, DEVELOPMENT.md, etc.

5. **Technical docs** (6 docs)
   - API.md, FRONTEND.md, etc.

**For each doc**:
- [ ] Identify Claude-specific content (tag with {{ if claude }})
- [ ] Identify Wiki-specific content (tag with {{ if wiki }})
- [ ] Extract variables (versions, URLs, etc.)
- [ ] Convert to .template.md
- [ ] Test generation
- [ ] Validate output

### 6.2: Generate Both Versions
- [ ] Run generate-wiki.py on all converted templates
- [ ] Verify Claude versions in docs/dev/design/
- [ ] Verify Wiki versions in docs/wiki/
- [ ] Compare Claude versions with originals (should match closely)
- [ ] Review Wiki versions for user-friendliness

### 6.3: Link Validation
- [ ] Run link checker on all generated docs
- [ ] Fix broken internal links
- [ ] Fix broken external links
- [ ] Validate cross-references
- [ ] Update link mappings if needed

### 6.4: Quality Pass
- [ ] Lint all generated docs (markdownlint)
- [ ] Check formatting consistency
- [ ] Verify all status tables
- [ ] Verify all breadcrumbs
- [ ] User spot-check and approval

**Exit Criteria**:
- ✅ All 136+ design docs converted to templates
- ✅ Both Claude and Wiki versions generated
- ✅ Zero broken links
- ✅ All linting clean
- ✅ User approves migration quality

**Estimated Duration**: 3-5 days

---

## PHASE 7: Skills & Tooling

**Dependencies**: Phase 6 complete

### 7.1: Create Claude Code Skills

1. **mvp-status skill**
   - [ ] Create `.claude/skills/mvp-status/SKILL.md`
   - [ ] Implement: Show MVP completion percentage
   - [ ] List what's done vs pending (from roadmap TODOs)
   - [ ] Suggest next task based on priority
   - [ ] Test skill

2. **generate-wiki skill**
   - [ ] Create `.claude/skills/generate-wiki/SKILL.md`
   - [ ] Wrapper for generate-wiki.py
   - [ ] Allow regenerating specific docs or all
   - [ ] Deploy to GitHub Wiki
   - [ ] Test skill

3. **verify-settings skill**
   - [ ] Create `.claude/skills/verify-settings/SKILL.md`
   - [ ] Run sync-tool-settings.py in check mode
   - [ ] Report mismatches
   - [ ] Optionally auto-fix
   - [ ] Test skill

4. **code-status skill**
   - [ ] Create `.claude/skills/code-status/SKILL.md`
   - [ ] Run verify-code-status.py
   - [ ] Show code implementation status
   - [ ] Compare with design status
   - [ ] Highlight gaps
   - [ ] Test skill

### 7.2: Update .claude/CLAUDE.md
- [ ] Document new skills
- [ ] Add usage examples
- [ ] Update workflow section
- [ ] Add MVP/roadmap references
- [ ] Add links to planning docs

### 7.3: Test All Skills
- [ ] Test each skill individually
- [ ] Test skill integration with existing workflow
- [ ] Verify output quality
- [ ] User acceptance testing

**Exit Criteria**:
- ✅ All 4 skills created and working
- ✅ .claude/CLAUDE.md updated
- ✅ Skills documentation complete
- ✅ User can successfully use all skills

**Estimated Duration**: 2-3 days

---

## PHASE 8: Validation & Rollout

**Dependencies**: Phase 7 complete

### 8.1: Full Validation
- [ ] Run full doc pipeline on all docs
- [ ] Run full source pipeline
- [ ] Lint all markdown files
- [ ] Lint all Python files (ruff)
- [ ] Lint all YAML files
- [ ] Validate all links
- [ ] Check all CI/CD workflows
- [ ] Test all 4 new skills
- [ ] Verify settings sync
- [ ] Run code status verification

### 8.2: Quality Checks
- [ ] Zero markdown linting errors
- [ ] Zero Python linting errors
- [ ] Zero broken links
- [ ] All CI/CD workflows pass
- [ ] All design docs have ✅ or 🟡
- [ ] Wiki renders correctly on GitHub
- [ ] Settings match SOURCE_OF_TRUTH 100%

### 8.3: User Review
- [ ] User reviews all changes
- [ ] User tests wiki generation
- [ ] User tests new skills
- [ ] User tests settings sync
- [ ] User approves rollout

### 8.4: Git Operations
- [ ] Create feature branch: `git checkout -b phase-1-design-completion`
- [ ] Commit Phase 1 work (design docs)
- [ ] Create tag: `git tag pre-automation-2026-01-31`
- [ ] Push tag: `git push origin pre-automation-2026-01-31`
- [ ] Commit Phases 2-7 work (automation)
- [ ] Push to feature branch
- [ ] Create PR to develop
- [ ] Merge after CI passes

### 8.5: Cleanup
- [ ] Archive .analysis/ directory to `.backup/analysis-2026-01-31/`
- [ ] Update TODO.md with post-implementation tasks
- [ ] Create GitHub issues for any follow-up work
- [ ] Update documentation indexes

**Exit Criteria**:
- ✅ All validation passed
- ✅ Git tags created
- ✅ All changes committed and pushed
- ✅ PR merged to develop
- ✅ .analysis/ archived
- ✅ No regressions

**Estimated Duration**: 1-2 days

---

## DEPENDENCY GRAPH

```
Phase 0 (Planning) ✅
    ↓
Phase 1 (Design Completion) ← CURRENT, BLOCKS EVERYTHING
    ↓
Phase 2 (MVP Definition & Roadmap) ← needs complete designs
    ↓
Phase 3 (Template System) ← needs roadmap for MVP filtering
    ↓
Phase 4 (Automation Scripts) ← needs templates
    ↓
Phase 5 (Pipeline Integration) ← needs scripts
    ↓
Phase 6 (Documentation Migration) ← needs pipelines
    ↓
Phase 7 (Skills & Tooling) ← needs migration complete
    ↓
Phase 8 (Validation & Rollout) ← needs everything
```

**Critical Path**: Phase 1 → Phase 2 → ... → Phase 8 (linear)

**NO parallel work possible** - each phase depends on previous

---

## ESTIMATED TIMELINE

| Phase | Estimated Duration | Cumulative |
|-------|-------------------|------------|
| Phase 0 | ✅ Complete | 0 days |
| Phase 1 | Variable (design completion) | ? days |
| Phase 2 | 3-5 days | +3-5 |
| Phase 3 | 2-3 days | +2-3 |
| Phase 4 | 3-4 days | +3-4 |
| Phase 5 | 2-3 days | +2-3 |
| Phase 6 | 3-5 days | +3-5 |
| Phase 7 | 2-3 days | +2-3 |
| Phase 8 | 1-2 days | +1-2 |
| **Total (after Phase 1)** | **16-25 days** | - |

**Phase 1 is variable** - depends on design writing speed and user involvement

---

## CRITICAL SUCCESS FACTORS

1. **Complete API.md and FRONTEND.md FIRST** - everything depends on these
2. **Don't skip design detail** - 99% perfection target means thorough docs
3. **Test each phase before moving to next** - no moving ahead with broken foundation
4. **User review at each phase boundary** - ensure alignment
5. **Maintain git safety** - feature branch, tags, backups

---

## NEXT IMMEDIATE STEPS

1. ✅ Collect all tasks (THIS DOCUMENT)
2. → Verify source fetch completed
3. → Review this master plan with user
4. → Start Phase 1.2: Create scaffold template
5. → Begin design writing (API.md first)

---

**STATUS**: Master plan complete - awaiting user review and approval to proceed
