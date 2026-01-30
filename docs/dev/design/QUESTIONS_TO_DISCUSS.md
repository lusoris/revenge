# Questions & Gaps to Discuss

> Temporary file for tracking discrepancies, gaps, and questions that need resolution

**Created**: 2026-01-30
**Status**: 🔄 Active collection

---

## Critical Discrepancies Found

### Package Version Mismatches

| Package | go.mod (actual) | SOURCES.yaml | SOURCE_OF_TRUTH.md | Resolution Needed |
|---------|-----------------|--------------|---------------------|-------------------|
| otter | v1.2.4 | v2.3.0 | v1.2.4 ✅ | Update SOURCES.yaml |
| typesense-go | v4.0.0-alpha2 | v3.2.0 ("NO v4 exists!") | v3.2.0 | **VERIFY**: Is v4 alpha stable enough? |
| resty | go-resty/resty/v2 v2.17.1 | resty v3.0.0-b6 | v2.17.1 ✅ | SOURCES.yaml references wrong package |

### Questions

1. **Typesense Go Client**:
   - go.mod uses `typesense-go/v4 v4.0.0-alpha2`
   - SOURCES.yaml explicitly says "NO v4 exists!"
   - Is v4 alpha suitable for production? Or should we downgrade to v3.2.0?
   - Need to verify latest stable from live docs

2. **Resty HTTP Client**:
   - go.mod uses v2 (`go-resty/resty/v2`)
   - SOURCES.yaml mentions v3 beta (`resty.dev/v3`)
   - Should we upgrade to v3 when stable? What's the migration effort?

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

### Broken Internal Links (CRITICAL)

| File | Line | Broken Link |
|------|------|-------------|
| ARCHITECTURE.md | 1020 | `PROJECT_STRUCTURE.md` - **does not exist** |
| DATA_RECONCILIATION.md | 426 | `RIVER_JOBS.md` - **does not exist** |
| SKIP_INTRO.md | 404 | `RIVER_JOBS.md` - **does not exist** |
| TRICKPLAY.md | 361 | `RIVER_JOBS.md` - **does not exist** |
| NEWS_SYSTEM.md | 320 | `RIVER_JOBS.md` - **does not exist** |

**Decision needed:** Create these files or update references?

### Unreferenced Documentation

- [ ] GALLERY_MODULE.md - not in features/adult/INDEX.md
- [ ] STASH.md - not in metadata/adult/INDEX.md

---

## Questions for Project Owner

1. **Database Strategy**:
   - Is SQLite truly needed for single-user deployments?
   - What's the priority for dual DB support vs other features?

2. **Typesense Version**:
   - v4 alpha in go.mod - intentional bleeding edge?
   - Should we pin to v3.2.0 stable?

3. **Package Update Policy**:
   - How aggressive on updates? (bleeding edge vs 1 version behind)
   - Who monitors changelogs/breaking changes?

4. **Documentation Priority**:
   - Should we complete SOURCE_OF_TRUTH first or start scaffolding?
   - User clarified: docs first, no scaffolding until consistent ✅

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

---

## Design Strategy Requirements (MANDATORY)

**Principle:** All advanced design patterns and coding strategies MUST be documented in SOURCE_OF_TRUTH.md BEFORE implementation.

**Required in SOT:**
- [ ] All performance patterns (caching, pooling, batching)
- [ ] All resilience patterns (circuit breaker, retry, fallback)
- [ ] All security patterns (auth, RBAC, isolation)
- [ ] All async patterns (jobs, queues, workers)
- [ ] All data patterns (transactions, consistency, partitioning)
- [ ] All API patterns (versioning, errors, pagination)

**Why:**
- No "write first, fix later" approach
- Prevents wasted time on reiteration
- Every implementation inherits patterns from SOT
- Consistency across all modules

**Current gaps to add to SOT:**
- [ ] Error handling patterns (Go errors, API errors)
- [ ] Testing patterns (unit, integration, mocks)
- [ ] Logging patterns (slog, structured, levels)
- [ ] Metrics patterns (Prometheus, OTel)
- [ ] Validation patterns (input, business rules)
- [ ] Pagination patterns (cursor, offset)

---

## Notes

- This file is temporary - items should be resolved and moved to appropriate docs
- When resolved, update SOURCE_OF_TRUTH.md and remove from here
- Questions for owner should be asked and answers documented in design docs
