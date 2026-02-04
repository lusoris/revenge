# TODO v0.3.0 - Index

**Last Updated**: 2026-02-04
**Status**: Pre-Frontend Phase
**Next**: Critical Fixes → Implementation Tasks → Frontend → MVP Release

---

## Quick Reference - Decisions Made (2026-02-04)

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
| A0 | [TODO_A0_CRITICAL_FIXES.md](TODO_A0_CRITICAL_FIXES.md) | P0-P2 | 15-20h | ✅ Complete |
| A1 | [TODO_A1_MOVIE_REPOSITORY.md](TODO_A1_MOVIE_REPOSITORY.md) | P2 | 6-8h | ✅ Complete |
| A2 | [TODO_A2_MOVIE_JOBS.md](TODO_A2_MOVIE_JOBS.md) | P2 | 4-6h | ✅ Complete |
| A3 | [TODO_A3_TEST_INFRASTRUCTURE.md](TODO_A3_TEST_INFRASTRUCTURE.md) | P2 | 3-4h | ✅ Complete |
| A4 | [TODO_A4_WEBHOOKS.md](TODO_A4_WEBHOOKS.md) | P2 | 2-3h | ✅ Complete |
| A5 | [TODO_A5_LIBRARY_MATCHER.md](TODO_A5_LIBRARY_MATCHER.md) | P2 | 4-6h | ✅ Complete |
| A6 | [TODO_A6_PERMISSIONS_SESSION.md](TODO_A6_PERMISSIONS_SESSION.md) | HIGH | 20-30h | Pending |
| B | [TODO_B_FRONTEND.md](TODO_B_FRONTEND.md) | - | 40-60h | Pending |
| C | [TODO_C_INFRASTRUCTURE.md](TODO_C_INFRASTRUCTURE.md) | - | 8-16h | Pending |

**Completed**: [TODO_COMPLETED.md](TODO_COMPLETED.md)

---

## Timeline Estimate

| Phase | Effort | Dependencies |
|-------|--------|--------------|
| A0: Critical Fixes | 15-20h | None |
| A1-A5: Stubs Completion | 20-30h | A0 |
| A6: Existing Tasks | 20-30h | A0 |
| B: Frontend | 40-60h | A0-A6 |
| C: Infrastructure | 8-16h | Phase B |
| **Total** | **103-156h** | |

**Estimated completion**: ~3-4 weeks full-time

---

## Priority Order

1. **A0.1-A0.3** (Critical blockers) - Must fix before any testing
2. **A0.4-A0.7** (P1 issues) - Required for functional MVP
3. **A6.1-A6.6** (Existing tasks) - Can parallel with A1-A5
4. **A1-A5** (P2 stubs) - Can defer some to post-MVP

---

## Notes

- Run `golangci-lint run ./...` before commits
- Run `go test ./... -short` frequently
- Update phase files as tasks complete
- All design docs in `docs/dev/design/`
