# Documentation Rewrite - Incremental Action Plan

## Problem Statement

**218 design docs** (~44K lines) massively out of sync with codebase. **6 root-level MDs** also outdated. Every layer of documentation is wrong, from README.md down to individual service docs.

### Root-Level MDs (Broken)
| File | Issue |
|------|-------|
| `README.md` | Claims 12 content types (only 2+QAR exist), lists Watch Next/Scrobbling/Audio Streaming (none implemented), claims Lidarr/Whisparr/Chaptarr integration (only Radarr+Sonarr exist), references `docs/wiki/` (doesn't exist) |
| `TODO.md` | Says River workers are "stub" (fully implemented), stale security items, no mention of recent metadata/CI fixes |
| `CONTRIBUTING.md` | References "Revenge (C#)" (wrong project), wrong project structure (`domain/`, `pkg/`, `migrations/`, `configs/` - none exist), mentions "Discord" (no server), "Phase 1-4" don't match roadmap |
| `SUPPORT.md` | References doc paths that will change during restructure |
| `SECURITY.md` | Version table says 0.1.x supported (project is pre-0.1) |
| `CHANGELOG.md` | Unknown state, likely auto-generated |

### Design Docs (Broken)
| Issue | Count |
|-------|-------|
| Status tables wrong (Code says 🔴 but code exists) | ~20 docs |
| Describe unimplemented features/integrations | ~150 docs |
| Code exists but no design doc | ~10 packages |
| Auto-generated shells with no real content | ~80% of docs |

### Key Mismatches Found
| Design Doc Says | Reality |
|-----------------|---------|
| Metadata Service `Code: 🔴 Not Started` | Full implementation with TMDb+TVDb providers, adapters, caching, ClearCache |
| MFA Service `Planned` | 9 Go files implementing TOTP+WebAuthn+backup codes |
| Notification Service `Planned` | 10 Go files with dispatcher+agent pattern |
| River workers `stub` | 9 workers fully implemented (metadata refresh, library scan, file match, search index, cleanup) |
| Grants/Fingerprint `✅ Complete` | Zero code exists |
| 12 content types | Only movie, tvshow, qar scaffold |

---

## Strategy: Incremental Steps

**Principle**: Each step is one commit, leaves the project in a consistent state, and makes things strictly better. Never worse.

**Order**: Fix what's most visible and impactful first (README, TODO), then work inward.

---

## Step 1: Fix Root-Level MDs

**Scope**: README.md, TODO.md, CONTRIBUTING.md
**Goal**: These files are the first thing anyone sees. Make them reflect reality.

### README.md Changes
- Remove features section claims for unimplemented things
- Split "Features" into "Implemented" vs "Planned"
- Fix architecture table (accurate)
- Remove `docs/wiki/` reference
- Update "Current Phase" description
- Keep Quick Start, Development, Testing sections (those are accurate)

### TODO.md Changes
- Mark River workers as ✅ DONE (fully implemented: 9 workers)
- Mark metadata system fixes as ✅ DONE (Force/Languages plumbing)
- Mark CI fixes as ✅ DONE (govulncheck CGO deps, migration paths)
- Update current priorities
- Remove stale security items if already fixed

### CONTRIBUTING.md Changes
- Remove "Revenge (C#)" reference
- Fix project structure to match actual layout (`content/`, `service/`, `infra/`, not `domain/`, `pkg/`)
- Remove Discord mention
- Update phase references
- Fix development setup commands (mention `make` commands, GOEXPERIMENT)

**Checkpoint**: `git diff` shows only root MDs changed. All doc links still work.

---

## Step 2: Fix Design Doc Status Tables

**Scope**: Only the `## Status` tables in existing design docs
**Goal**: Make "Code" status accurate without rewriting content.

### Services with wrong Code status
| Doc | Current | Correct |
|-----|---------|---------|
| services/METADATA.md | 🔴 Not Started | 🟡 Partial |
| services/MFA.md (if exists) | 🔴 Planned | 🟡 Partial |
| services/NOTIFICATION.md | 🔴 Planned | 🟡 Partial |
| services/GRANTS.md | ✅ Complete | 🔴 Not Started (no code) |
| services/FINGERPRINT.md | ✅ Complete | 🔴 Not Started (no code) |
| services/TRANSCODING.md | ✅ Complete | 🔴 Not Started (no code) |

### Content modules with wrong status
| Doc | Current | Correct |
|-----|---------|---------|
| features/video/MOVIE_MODULE.md | Code: 🟡 | Code: 🟡 (correct but Unit Testing should be 🟡) |
| architecture/METADATA_SYSTEM.md | Code: 🔴 | Code: 🟡 Partial |

**Checkpoint**: Only status tables changed. No structural changes. Content preserved.

---

## Step 3: Create `planned/` Directory + Move Unimplemented Docs

**Scope**: File moves only, no content changes
**Goal**: Clearly separate "what exists" from "what's designed but not built"

### Create structure
```
docs/dev/design/planned/
  features/       # Unimplemented content modules
  integrations/   # Unimplemented integrations
  services/       # Unimplemented services
  technical/      # Unimplemented technical features
```

### Move to `planned/features/`
- features/music/MUSIC_MODULE.md
- features/audiobook/AUDIOBOOK_MODULE.md
- features/book/BOOK_MODULE.md
- features/comics/COMICS_MODULE.md
- features/livetv/LIVE_TV_DVR.md
- features/photos/PHOTOS_LIBRARY.md
- features/podcasts/PODCASTS.md
- features/playback/ (ALL files - Watch Next, Skip Intro, SyncPlay, Trickplay, etc.)
- features/shared/ (most files - Collections, Scrobbling, Voice Control, Wiki, News, Ticketing, etc.)
- features/adult/ (most files - keep ADULT_CONTENT_SYSTEM.md in place since QAR has code)

### Move to `planned/integrations/`
- ALL integration docs EXCEPT:
  - integrations/infrastructure/ (Dragonfly, PostgreSQL, River, Typesense - all implemented)
  - integrations/metadata/video/TMDB.md and THETVDB.md (implemented)
  - integrations/servarr/RADARR.md and SONARR.md (implemented)
  - integrations/auth/GENERIC_OIDC.md (implemented)

### Move to `planned/services/`
- services/ANALYTICS.md
- services/EPG.md
- services/FINGERPRINT.md
- services/GRANTS.md
- services/TRANSCODING.md
- services/USER_SETTINGS.md (if exists)

### Move to `planned/technical/`
- technical/AUDIO_STREAMING.md
- technical/WEBHOOKS.md
- technical/WEBSOCKETS.md
- technical/design/ (ALL frontend design docs - BRAND_IDENTITY, COLOR_SYSTEM, COMPONENTS, etc.)
- technical/DESIGN_SYSTEM.md
- technical/PIRATE_MODE.md

**Checkpoint**: `git status` shows only file moves. All content preserved. Planned features accessible. Links in moved docs may break (acceptable - they'll be fixed in later steps).

---

## Step 4: Clean Up Auto-Generated Shells

**Scope**: Delete files that have no useful content
**Goal**: Remove noise, keep signal

### Delete
- docs/dev/design/01_DESIGN_DOC_TEMPLATE.md (replace with simple TEMPLATE.md)
- docs/dev/design/03_DESIGN_DOCS_STATUS.md (stale tracking file)
- docs/dev/design/.templates/ (auto-generation templates no longer needed)
- All empty INDEX.md files in moved directories (will regenerate)

### Create
- docs/dev/design/TEMPLATE.md (simple 30-line template for new docs)

**Checkpoint**: Less noise. Core docs untouched.

---

## Step 5: Rewrite Architecture Docs (from code)

**Scope**: 3 architecture docs
**Goal**: Architecture docs match actual implementation

### 5a: architecture/ARCHITECTURE.md
- Rewrite from actual fx module wiring (`internal/app/module.go`)
- Document real layers: API → Service → Content → Infra → Integration
- Accurate dependency graph
- Max 400 lines

### 5b: architecture/DESIGN_PRINCIPLES.md
- Document actual patterns used in codebase
- fx DI, repository pattern, handler pattern, worker pattern
- Remove aspirational patterns not yet implemented

### 5c: architecture/METADATA_SYSTEM.md
- Rewrite from actual provider chain (metadata.Service → TMDb/TVDb → adapters → content modules)
- Document ClearCache flow, language handling, Force refresh
- Max 300 lines

**Checkpoint**: Architecture docs are accurate. Each is a standalone commit.

---

## Step 6: Rewrite Service Docs (from code, one at a time)

**Scope**: 15 service docs, one per commit
**Goal**: Each service doc accurately describes its interface, config, and usage

### Priority order (dependencies first)
1. **METADATA.md** - Most complex, recently fixed
2. **AUTH.md** - Core service, many dependents
3. **SESSION.md** - Auth companion
4. **MFA.md** - Auth companion
5. **OIDC.md** - Auth companion
6. **USER.md** - Core service
7. **RBAC.md** - Access control
8. **APIKEYS.md** - API access
9. **LIBRARY.md** - Content management
10. **SEARCH.md** - Content discovery
11. **SETTINGS.md** - Configuration
12. **ACTIVITY.md** - Audit logging
13. **EMAIL.md** - NEW doc (code exists, no doc)
14. **NOTIFICATION.md** - Partial implementation
15. **STORAGE.md** - NEW doc (code exists, no doc)

### Per-doc template (max 250 lines)
```markdown
# {Service Name}

Package: `internal/service/{name}/`
Module: `{name}.Module` (fx)

## Purpose
1-2 sentences.

## Interface
Key methods with signatures.

## Configuration
Config keys from koanf.

## Dependencies
What this service depends on.

## Current Status
What's implemented vs what's planned.
```

**Checkpoint**: Each doc is a standalone commit. Old doc replaced, not appended to.

---

## Step 7: Write Infrastructure Docs (new)

**Scope**: 8 new docs for `infra/` section
**Goal**: Document infrastructure that currently has zero docs

1. **DATABASE.md** - pgxpool config, 30+ migrations, sqlc codegen, schemas
2. **CACHE.md** - otter L1 (W-TinyLFU), rueidis L2 (Dragonfly), graceful degradation
3. **JOBS.md** - River setup, 9 workers, retry config, queue management
4. **HEALTH.md** - K8s probes (liveness/readiness/startup), dependency checks
5. **IMAGE.md** - govips processing, avatar/poster optimization
6. **LOGGING.md** - slog (dev) + zap (prod), structured logging
7. **OBSERVABILITY.md** - pprof, OpenTelemetry, metrics
8. **SEARCH_INFRA.md** - Typesense client, collection management, health checks

**Checkpoint**: Each doc standalone commit. New section, no conflicts with existing.

---

## Step 8: Rewrite Content Module Docs (from code)

**Scope**: Movie, TV Show, QAR, Shared
**Goal**: Accurate module documentation

### 8a: Movie module
- `modules/movie/OVERVIEW.md` - entities, service, repo, handler
- `modules/movie/API.md` - endpoints from ogen handlers
- `modules/movie/JOBS.md` - 4 workers (metadata refresh, library scan, file match, search index)

### 8b: TV Show module
- `modules/tvshow/OVERVIEW.md` - entities (show/season/episode), service, repo
- `modules/tvshow/API.md` - endpoints
- `modules/tvshow/JOBS.md` - 5 workers

### 8c: Shared module
- `modules/shared/LIBRARY.md` - library scanning, cleanup
- `modules/shared/METADATA.md` - shared metadata types

**Checkpoint**: Content docs accurate. Old feature docs already moved to planned/.

---

## Step 9: Rewrite Integration Docs (implemented only)

**Scope**: 4-5 docs for implemented integrations
**Goal**: Document only what's actually built

1. **integrations/metadata/TMDB.md** - TMDb provider, client, caching (sync.Map, 24h TTL)
2. **integrations/metadata/TVDB.md** - TVDb provider, client, caching
3. **integrations/radarr/OVERVIEW.md** - Radarr sync, client, mapper
4. **integrations/sonarr/OVERVIEW.md** - Sonarr sync, client, mapper

**Checkpoint**: Integration docs match code. Unimplemented integrations safely in planned/.

---

## Step 10: Rewrite Operations + Technical Docs

**Scope**: CI/CD, deployment, dev setup, API, config, testing
**Goal**: Operational docs match actual setup

1. **operations/CI_CD.md** - from 8 GitHub Actions workflows
2. **operations/DEPLOYMENT.md** - from Dockerfile, Helm chart, 5 compose files
3. **operations/DEVELOPMENT.md** - from Makefile, Coder setup
4. **technical/API.md** - from OpenAPI spec, ogen config
5. **technical/CONFIGURATION.md** - from koanf config struct (~19KB)
6. **technical/TESTING.md** - from testcontainers, 144 test files

---

## Step 11: Finalize

1. Update `00_SOURCE_OF_TRUTH.md` - verify all dependency versions against go.mod
2. Regenerate `DESIGN_INDEX.md` from new structure (only list implemented docs + planned section)
3. Update `CLAUDE.md` doc references to new paths
4. Update `SUPPORT.md` + `SECURITY.md` with correct paths/versions
5. Verify all cross-references (no broken links)
6. Clean up empty directories from moved files

---

## Doc Size Strategy

**Target**: Each doc readable in a single AI context pass alongside source code.

| Doc Type | Max Lines | Rationale |
|----------|-----------|-----------|
| Module OVERVIEW | 300 | Core interfaces + architecture |
| Module API | 200 | Endpoint list + examples |
| Module JOBS | 150 | Worker list + flow |
| Service doc | 250 | Interface + config + usage |
| Infra doc | 200 | Setup + patterns |
| Architecture doc | 400 | Diagrams + decisions |
| SOURCE_OF_TRUTH | 1500 | Reference (exception, kept as-is) |

**Rules**:
1. No doc exceeds 400 lines (except SOURCE_OF_TRUTH)
2. Split large topics into sub-docs (e.g., movie/ has OVERVIEW, API, JOBS)
3. Cross-reference instead of duplicate (link to SOURCE_OF_TRUTH for versions)
4. Keep mermaid diagrams simple (max 15 nodes)
5. Code examples: max 20 lines inline, link to source for longer
6. Status is a single "Current Status" section, not a multi-row table

---

## Scope Summary

| Step | What | Commits | Effort | Status |
|------|------|---------|--------|--------|
| 1 | Fix root MDs (README, TODO, CONTRIBUTING) | 1 | Low | ✅ Done (0a6f838) |
| 2 | Fix design doc status tables | 1 | Low | ✅ Done (0a6f838) |
| 3 | Create planned/ + move ~150 docs | 1 | Low (file moves) | ✅ Done (0a6f838) |
| 4 | Delete auto-generated shells + stale artifacts | 1 | Low | ✅ Done (0a6f838) - deleted docs/wiki/ (161), .templates/ (9), tests/automation/ (23), .shared/ (9) |
| 5 | Rewrite 3 architecture docs | 1 | Medium | ✅ Done (14935fe) - ARCHITECTURE, DESIGN_PRINCIPLES, METADATA_SYSTEM |
| 6 | Align 15 service docs (+ 2 new) | 17 | High (biggest step) | ✅ Done |
| 7 | Write 8 new infra docs (+ INDEX) | 9 | Medium | ✅ Done |
| 8 | Rewrite content module docs | ~8 | Medium | ✅ Done |
| 9 | Rewrite 4 integration docs | 4 | Low-Medium | Pending |
| 10 | Rewrite 6 operations/technical docs | 6 | Medium | Pending |
| 11 | Finalize (SOURCE_OF_TRUTH, INDEX, links) | 1 | Low | Pending |

**Also completed** (not in original plan):
- Wiki sync workflow (.github/workflows/wiki-sync.yml) - auto-syncs docs to GitHub wiki
- Fixed all cross-references broken by Step 3 renames (78 active docs + 7 root/config files)
- Fixed tool folder issues: git hooks (Jellyfin→Revenge), .vscode (SQLite→PostgreSQL, YAML→JSON), .github (phantom workflows), .zed (wrong org name)
- Cleaned up CLAUDE.md: removed stale SOT dependency, added accurate architecture doc links
- Deleted 4 duplicate CLAUDE.md files, entire .shared/ folder
- Fixed link aliases in README.md (raw filenames → human-readable)
- Notification system split: designed 4 independent systems (dispatcher, announcements, helpdesk/wiki, push). Created `planned/services/ANNOUNCEMENTS.md`. Design decisions in `.workingdir3/QUESTIONS_INBOX_HELPDESK.md`

**Total**: ~45 docs to write/rewrite, ~150 to move, ~15 to delete, ~50 commits

---

## Rules

- All documentation in English (CLAUDE.md language policy)
- Each step is 1 commit that leaves everything consistent
- Write docs from code (read Go source, document what IS, not what SHOULD BE)
- Planned docs preserved for future reference
- Never delete user-written content without moving it first
- Steps 1-4 are prep work (can be done in one session)
- Steps 5-10 are the main work (can be spread across sessions)
- Any step can be paused and resumed without leaving a mess
