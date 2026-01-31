# Original TODO.md Backup

**Saved**: 2026-01-31
**Reason**: Backup before clearing for restructuring

---

# Revenge - Development Roadmap

> Modular media server with complete content isolation

**Last Updated**: 2026-01-30
**Status**: Fresh Start - Design Phase

---

## Current Focus

**Milestone 1: Complete Design Scaffold**

All design documentation must be complete, consistent, and English-only before any code implementation begins.

### Workflow (Sequential)

1. [x] **Review Core Docs** - 00_SOURCE_OF_TRUTH.md, 01_ARCHITECTURE.md, 02_QUESTIONS_TO_DISCUSS.md
2. [x] **Create Missing Design Docs** - Fingerprint, Grants, Search, Analytics, Notification
3. [x] **Second Review** - Check for additional gaps, problems, or new questions
4. [x] **Quality Pass** - English only, cross-refs, deduplication, fix broken links
5. [x] **Correct Sources** - Update SOURCES.yaml with real web data (current versions, URLs)
6. [ ] **Restructure Instructions** - Using Claude Agent SDK + Claude Code live docs as SOT:
   - Review `.github/instructions/*.instructions.md` files
   - Cross-reference with SOURCES.yaml and SOT
   - Deduplicate, index, and link instruction files properly
   - Create instruction templates (based on Claude Agent SDK patterns)
   - Optimize for AI-assisted development workflows
7. [ ] **Expand Documentation** - Dev docs, full specs, template code patterns
8. [ ] **Recreate Instructions** - Using Claude Code design specs (online source only, no memory)

---

## M1: Design Scaffold

> Complete all design documentation before implementation

### Core Documents

- [x] 00_SOURCE_OF_TRUTH.md - Final review + remove SQLite references
- [x] 01_ARCHITECTURE.md - Verify alignment with SOT, deduplicate versions
- [x] 02_QUESTIONS_TO_DISCUSS.md - Translate German, resolve decisions

### Design Patterns (Add to SOT)

- [x] Metadata Priority Chain (Local → Arr → Internal → External → Enrichment)
- [x] Error Handling (Sentinels + APIError)
- [x] Testing (Table-driven + testify + mockery)
- [x] Logging (Text dev + JSON prod)
- [x] Metrics (Prometheus + OpenTelemetry)
- [x] Validation (ogen API + go-playground business)
- [x] Pagination (Cursor default + Offset option)

### Package Decisions (Finalized)

- [x] Database: pgx v5 (PostgreSQL only)
- [x] Cache: rueidis (distributed) + otter (local)
- [x] Jobs: River (PostgreSQL-native)
- [x] HTTP: resty v2
- [x] WebSocket: gobwas/ws
- [x] GraphQL: Khan/genqlient
- [x] FFmpeg: go-astiav
- [x] Audio Tags: go-taglib
- [x] Raft: hashicorp/raft

### Missing Design Docs

- [x] Fingerprint Service design doc (`services/FINGERPRINT.md`)
- [x] Grants Service design doc (`services/GRANTS.md`)
- [x] Search Service design doc (`services/SEARCH.md`)
- [x] Analytics Service design doc (`services/ANALYTICS.md`)
- [x] Notification Service design doc (`services/NOTIFICATION.md`)

### Documentation Quality

- [x] All docs in English (no German)
- [x] Consistent cross-references (SOT ↔ detail docs)
- [x] No duplicate information (versions in SOT only)
- [x] Obsolete docs moved to `.archive/`
- [x] File renaming (see `.archive/FILE_RENAMING_STRATEGY.md`) - deferred to after MVP

### OpenAPI Specs

- [ ] Design new OpenAPI specs (aligned with SOT)
- [ ] Per-module API definitions
- [ ] Shared components/schemas

---

## Completion Criteria

Before starting M2 (Code Implementation):

1. **SOT Complete** - All packages, versions, patterns documented
2. **Consistency** - All docs reference SOT, no contradictions
3. **English Only** - No German text in any documentation
4. **No Gaps** - All services have design docs or SOT entries

---

## Future Milestones (After Design Complete)

- M2: Code Structure Scaffold (folders, interfaces, modules, fx wiring)
- M3: Infrastructure (pgx, rueidis, River, health)
- M4: Core Services (auth, user, session, RBAC)
- M5: Content Modules (movie, tv, music, etc.)
- M6: External Integrations (metadata providers)
- M7: Playback & Streaming
- M8: Frontend (SvelteKit)

---

## Key Decisions (Q&A Session 2026-01-30)

| Decision | Answer |
|----------|--------|
| Database | PostgreSQL ONLY (no SQLite) |
| Updates | 1 minor behind, Dependabot |
| Test Coverage | 80% minimum |
| ErsatzTV | REST client only (no native IPTV) |
| K8s | controller-runtime (Operator pattern) |
| Code Strategy | Docs first, then fresh implementation |
| Multi-OIDC | Multiple providers (allow multiple login buttons) |
| Webhooks | Per-service endpoints (`/webhooks/radarr`, etc.) |
| Rate Limits | Per-provider (each API has own limits) |
| Watch History Sync | Full history with watch counts (like Trakt) |
| Logging | slog/tint (dev) + zap (prod) |
| Policy Engine | OPA alongside Casbin |
| Testing | testcontainers + embedded-postgres |
