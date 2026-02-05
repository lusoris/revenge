# TODO v0.3.0 - Index

**Last Updated**: 2026-02-05
**Status**: Extended Phase A (Security, Cluster, Multi-Language, TV Module)
**Next**: Security Fixes → Cluster → Multi-Language → Shared Code → TV Module → Frontend

---

## Quick Reference - Decisions Made (2026-02-05)

| Topic | Decision | Date |
|-------|----------|------|
| Storage (Media) | Hybrid: NFS for media files | 2026-02-05 |
| Storage (Avatars) | S3/MinIO | 2026-02-05 |
| Leader Election | River + Raft | 2026-02-05 |
| Multi-Language Schema | Hybrid JSONB | 2026-02-05 |
| Languages Priority | en, de, fr, es, ja | 2026-02-05 |
| Age Ratings | All TMDb systems | 2026-02-05 |
| TV Metadata | TMDb primary | 2026-02-05 |
| Refactoring Timing | Jetzt komplett (before TV) | 2026-02-05 |
| Multi-Source Matching | Multiple strategies needed | 2026-02-05 |

### Previous Decisions (2026-02-04)

| Topic | Decision |
|-------|----------|
| Permissions | Fine-grained (`movie:list`, `movie:get`, etc.) |
| Cache Tests | testcontainers for Redis L2 paths |
| Session Storage | Hybrid (Dragonfly L1 + PostgreSQL L2) |
| Queue Priorities | 5 levels: critical, high, default, low, bulk |
| Log Retention | 90 days default, configurable |
| Pagination | Both cursor + offset, cursor default |
| Device Fingerprinting | Deferred to v0.6.0 (Transcoding) |

---

## Phase Overview

| Phase | File | Priority | Effort | Status |
|-------|------|----------|--------|--------|
| **Phase A: Core Backend (Complete)** |
| A0 | [TODO_A0_CRITICAL_FIXES.md](TODO_A0_CRITICAL_FIXES.md) | P0-P2 | 15-20h | ✅ Complete |
| A1 | [TODO_A1_MOVIE_REPOSITORY.md](TODO_A1_MOVIE_REPOSITORY.md) | P2 | 6-8h | ✅ Complete |
| A2 | [TODO_A2_MOVIE_JOBS.md](TODO_A2_MOVIE_JOBS.md) | P2 | 4-6h | ✅ Complete |
| A3 | [TODO_A3_TEST_INFRASTRUCTURE.md](TODO_A3_TEST_INFRASTRUCTURE.md) | P2 | 3-4h | ✅ Complete |
| A4 | [TODO_A4_WEBHOOKS.md](TODO_A4_WEBHOOKS.md) | P2 | 2-3h | ✅ Complete |
| A5 | [TODO_A5_LIBRARY_MATCHER.md](TODO_A5_LIBRARY_MATCHER.md) | P2 | 4-6h | ✅ Complete |
| A6 | [TODO_A6_PERMISSIONS_SESSION.md](TODO_A6_PERMISSIONS_SESSION.md) | HIGH | 20-30h | 🔶 A6.1-A6.5 ✅, A6.6 Pending |
| **Phase A: Extended (New - 2026-02-05)** |
| A7 | [TODO_A7_SECURITY_FIXES.md](TODO_A7_SECURITY_FIXES.md) | P0 | 16-24h | 🔴 Pending |
| A8 | [TODO_A8_CLUSTER_READINESS.md](TODO_A8_CLUSTER_READINESS.md) | P0 | 24-40h | 🔴 Pending |
| A9 | [TODO_A9_MULTILANGUAGE.md](TODO_A9_MULTILANGUAGE.md) | P0 | 32-48h | ✅ Complete |
| A10 | [TODO_A10_SHARED_ABSTRACTIONS.md](TODO_A10_SHARED_ABSTRACTIONS.md) | P1 | 40-60h | 🔴 Pending |
| A11 | [TODO_A11_TV_MODULE.md](TODO_A11_TV_MODULE.md) | P1 | 32-48h | 🔴 Pending |
| **Phase B & C (Deferred)** |
| B | [TODO_B_FRONTEND.md](TODO_B_FRONTEND.md) | - | 40-60h | ⏸️ Deferred |
| C | [TODO_C_INFRASTRUCTURE.md](TODO_C_INFRASTRUCTURE.md) | - | 8-16h | ⏸️ Deferred |

**Completed**: [TODO_COMPLETED.md](TODO_COMPLETED.md)

---

## Timeline Estimate

### Phase A (Original + Extended)

| Phase | Effort | Dependencies | Status |
|-------|--------|--------------|--------|
| A0: Critical Fixes | 15-20h | None | ✅ Complete |
| A1-A5: Stubs Completion | 20-30h | A0 | ✅ Complete |
| A6: Existing Tasks | 20-30h | A0 | 🔶 Partial (A6.6 pending) |
| **A7: Security Fixes** | **16-24h** | **None (parallel)** | 🔴 **Pending** |
| **A8: Cluster Readiness** | **24-40h** | **A7** | 🔴 **Pending** |
| **A9: Multi-Language** | **32-48h** | **A7, A8** | ✅ **Complete** |
| **A10: Shared Abstractions** | **40-60h** | **A9** | 🔴 **Pending** |
| **A11: TV Module** | **32-48h** | **A9, A10** | 🔴 **Pending** |
| **Phase A Total** | **199-300h** | | **60% Complete** |

### Phase B & C (Deferred)

| Phase | Effort | Dependencies | Status |
|-------|--------|--------------|--------|
| B: Frontend | 40-60h | A0-A11 | ⏸️ Deferred |
| C: Infrastructure | 8-16h | Phase B | ⏸️ Deferred |

### Total Estimate

| Category | Original | New (Extended) | Change |
|----------|----------|----------------|--------|
| Phase A | 55-80h | 199-300h | +144-220h |
| Phase B | 40-60h | 40-60h | No change |
| Phase C | 8-16h | 8-16h | No change |
| **Total** | **103-156h** | **247-376h** | **+144-220h** |

**Estimated completion (Phase A)**: ~5-8 weeks full-time (instead of original 2-3 weeks)
**Reason for extension**: Security hardening, cluster readiness, multi-language support, code deduplication, and TV module

---

## Priority Order (Updated 2026-02-05)

### Completed ✅
- ~~A0: Critical Fixes~~ ✅
- ~~A1-A5: Stubs Completion~~ ✅
- ~~A6.1-A6.5: Permissions & Session~~ ✅

### Current Focus (Sequential)
1. **A7: Security Fixes** (P0) - Critical vulnerabilities
   - Transactions, timing attacks, goroutine leaks
   - 16-24h effort

2. **A8: Cluster Readiness** (P0) - Production deployment blocker
   - NFS + S3 storage, River + Raft, request correlation
   - 24-40h effort

3. ~~**A9: Multi-Language** (P0)~~ ✅ **Complete**
   - Hybrid JSONB schema, TMDb multi-language, age ratings
   - 32-48h effort
   - **Multi-language support now in place**

4. **A10: Shared Abstractions** (P1) - Code deduplication
   - Scanner, matcher, provider, library frameworks
   - 40-60h effort
   - **Reduces TV module effort by 60-70%**

5. **A11: TV Module** (P1) - First expansion module
   - Reuses A10 abstractions, includes A9 multi-language
   - 32-48h effort (vs 100+ without A10!)

6. **A6.6: Test Coverage** (P2) - Ongoing
   - Target: 80% across all services
   - Can be done in parallel

### Deferred
- **Phase B: Frontend** - After A7-A11 complete
- **Phase C: Infrastructure** - After Phase B

---

## Notes

- Run `golangci-lint run ./...` before commits
- Run `go test ./... -short` frequently
- Update phase files as tasks complete
- All design docs in `docs/dev/design/`
