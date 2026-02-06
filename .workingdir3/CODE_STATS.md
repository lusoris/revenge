# Revenge - Code Statistics

> Generated 2026-02-06 via `scc` from working tree + `git show` for generated files

---

## TL;DR

| Metric | Count |
|--------|-------|
| **Hand-written Go (production)** | **42,878 lines of code** |
| **Hand-written Go (tests)** | **41,219 lines of code** |
| **Generated Go (ogen)** | ~169,000 lines (not on disk, in git) |
| **SQL migrations** | 2,969 lines across 64 files |
| **OpenAPI spec** | 8,546 lines (YAML) |
| **Total hand-written Go** | **84,097 lines of code** |

Production-to-test ratio: **1 : 0.96** (nearly 1:1 test coverage by volume)

---

## Hand-Written Go Breakdown

### By Layer

| Layer | Files | Lines | Code | Complexity |
|-------|-------|-------|------|------------|
| `internal/service/` | 132 | 50,129 | 37,893 | 4,130 |
| `internal/content/` | 73 | 19,950 | 15,740 | 1,962 |
| `internal/api/` (handlers, not ogen) | 41 | 14,383 | 11,252 | 1,767 |
| `internal/infra/` (Go only) | 60 | 11,194 | 8,167 | 626 |
| `internal/errors/` | 4 | 665 | 472 | 13 |
| `internal/app/` | 1 | 81 | 60 | 0 |
| `cmd/` + `tests/` | 15 | 3,956 | 2,812 | 178 |
| **Total** | **359+16** | **~110,000** | **~84,000** | **~9,700** |

### Services (14 services, 132 files)

| Service | Files | Lines | Code |
|---------|-------|-------|------|
| metadata | 18 | 8,087 | 6,110 |
| auth | 11 | 4,824 | 3,649 |
| oidc | 7 | 4,268 | 3,253 |
| user | 9 | 4,397 | 3,433 |
| notification | 9 | 4,128 | 3,171 |
| mfa | 9 | 3,875 | 2,833 |
| rbac | 12 | 3,722 | 2,747 |
| activity | 11 | 3,011 | 2,285 |
| session | 10 | 2,824 | 2,169 |
| settings | 8 | 2,094 | 1,521 |
| apikeys | 7 | 1,807 | 1,338 |
| search | 5 | 1,677 | 1,324 |
| email | 3 | 707 | 564 |
| storage | 5 | 689 | 499 |

### Content Modules (73 files)

| Module | Files | Lines | Code |
|--------|-------|-------|------|
| movie | 34 | 7,820 | 6,155 |
| tvshow | 15 | 7,271 | 5,866 |
| shared | 24 | 4,859 | 3,719 |
| qar | 0 | 0 | 0 (scaffold only, no Go files yet) |

### Production vs Test

| Category | Files | Code Lines |
|----------|-------|------------|
| Production Go (no ogen) | 231 | 42,878 |
| Test Go | 144 | 41,219 |
| **Total hand-written** | **375** | **84,097** |

---

## Non-Go Code

| Asset | Lines (code) | Notes |
|-------|-------------|-------|
| SQL migrations | 2,969 | 64 migration files in `internal/infra/database/migrations/` |
| OpenAPI spec | 8,546 | Single YAML file, drives ogen codegen |
| Helm chart | 594 | 7 YAML + 1 Smarty template |
| GitHub Actions | ~1,640 | 9 workflow files |
| Makefile | 111 | 25 targets |
| Dockerfile | 48 | Multi-stage build |
| Terraform (Coder) | 783 | Workspace template |
| Python scripts | 551 | 7 utility scripts |

---

## Generated Code (not on disk)

| Generator | Files | ~Lines | Notes |
|-----------|-------|--------|-------|
| ogen (OpenAPI) | 20 | ~169,000 | Tracked in git, generated from OpenAPI spec |

The ogen-generated code is **2x the entire hand-written codebase**. It's committed to git but not present on disk (regenerated via `make ogen`).

---

## Full Project (all languages)

| Language | Files | Code Lines |
|----------|-------|------------|
| Go | 375 | 84,097 |
| Markdown | 391 | 72,162 |
| JSON | 26 | 22,092 |
| YAML | 39 | 11,476 |
| Plain Text | 2 | 6,530 |
| SQL | 85 | 2,969 |
| Terraform | 1 | 783 |
| Jinja | 6 | 766 |
| Python | 7 | 551 |
| Shell | 4 | 167 |
| TOML | 3 | 157 |
| Makefile | 1 | 111 |
| Dockerfile | 2 | 66 |
| **Total** | **944** | **~202,000** |

---

## Key Ratios

| Metric | Value |
|--------|-------|
| Test : Production code | 0.96 : 1 |
| Services : Content modules | 37,893 : 15,740 (2.4 : 1) |
| Largest service | metadata (6,110 loc) |
| Largest content module | movie (6,155 loc) |
| Avg lines per Go file | ~224 |
| Avg lines per production Go file | ~186 |
| SQL per migration | ~46 lines avg |
