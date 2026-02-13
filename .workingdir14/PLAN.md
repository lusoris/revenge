# Workingdir 14 — Go 1.26 Version Bump Plan

**Date:** 2026-02-13
**Scope:** Upgrade from Go 1.25.7 to Go 1.26.0
**Status:** PLANNING (no implementation yet)

---

## 1. User Decisions (Collected)

### Q1: GOEXPERIMENT value
**Decision:** Remove `greenteagc` (now default in 1.26), keep `jsonv2` (still experimental, not shipped in 1.26).
User also wants all new relevant experiments enabled.

### Q2: New GOEXPERIMENT: goroutineleakprofile
**Decision:** YES — enable everywhere.
- Zero overhead unless actively queried
- Exposes `/debug/pprof/goroutineleak`
- Production-ready (experimental only for API feedback purposes)
- Aiming to be default in Go 1.27

### Q3: New GOEXPERIMENT: simd
**Decision:** YES — enable it. User wants to support/use it if anything can be made faster.
- Provides `simd/archsimd` package (amd64 only, 128/256/512-bit vectors)
- API is unstable but user wants to be forward-looking
- Could benefit image/video processing code paths

### Q4: New GOEXPERIMENT: runtimesecret
**Decision:** YES — enable it for security posture.
- Securely erases crypto temporaries (registers, stack, heap)
- amd64+arm64 on Linux only
- App handles API keys, MFA secrets, OIDC tokens
- User explicitly wants security-first approach

### Q5: go fix modernizers
**Decision:** YES — run `go fix ./...` as part of upgrade.

### Q6: errors.AsType migration
**Decision:** YES — migrate existing `errors.As` calls to `errors.AsType[T]()`.
- Generic, type-safe, faster replacement

---

## 2. Final GOEXPERIMENT Value

### Production (Dockerfile, CI, Makefile, .zed, .devcontainer, .coder):
```
GOEXPERIMENT=jsonv2,goroutineleakprofile,simd,runtimesecret
```

### Security workflow (govulncheck):
The govulncheck + jsonv2 panic (golang/go#74846, dup of #73871) is **FIXED in Go 1.26**.
The `go/types.NewSignatureType` panic was resolved in Go 1.26rc1 (confirmed by users).

**Therefore:** The security.yml workaround that excluded jsonv2 from GOEXPERIMENT can be **REMOVED**.
govulncheck can now use the same GOEXPERIMENT as everything else.

---

## 3. Go 1.26 Release Notes Summary (from https://go.dev/doc/go1.26)

### Language Changes
- `new(expr)` — `new` now accepts an expression operand (initial value), e.g. `new(yearsSince(born))`
- Self-referential generic type constraints now allowed (`type Adder[A Adder[A]]`)

### Tools
- **`go fix`** completely revamped — dozens of modernizers, push-button code updates
- `go mod init` now defaults to `go 1.(N-1).0` (so Go 1.26 creates `go 1.25.0`)
- `cmd/doc` deleted, use `go doc` instead
- pprof web UI defaults to flame graph view

### Runtime
- **Green Tea GC now default** — opt-out with `GOEXPERIMENT=nogreenteagc` (removed in 1.27)
- **~30% faster cgo calls** — directly benefits go-astiav (FFmpeg) and govips
- **Heap base address randomization** on 64-bit (security enhancement for cgo apps)
- **Goroutine leak profile** (experimental) — see Q2 above

### Compiler
- More stack allocations for slices (better perf)

### Linker
- Various section changes (transparent to programs)
- windows/arm64 internal linking for cgo

### Standard Library — New Packages
- `crypto/hpke` — Hybrid Public Key Encryption (RFC 9180)
- `simd/archsimd` (experimental) — SIMD vector operations
- `runtime/secret` (experimental) — secure memory erasure

### Standard Library — Key Changes
- `bytes.Buffer.Peek(n)` — peek without advancing
- **`errors.AsType[T]()`** — generic, type-safe, faster `errors.As`
- `fmt.Errorf("x")` — allocates less, matches `errors.New("x")`
- `io.ReadAll` — ~2x faster, ~half memory
- **`log/slog.NewMultiHandler`** — fan-out to multiple handlers
- `net.Dialer` — new `DialIP`, `DialTCP`, `DialUDP`, `DialUnix` methods with context
- `net/http.HTTP2Config.StrictMaxConcurrentRequests`
- `net/http.Transport.NewClientConn` — custom connection management
- `net/url.Parse` — rejects malformed hostnames with colons (GODEBUG `urlstrictcolons`)
- `os.Process.WithHandle` — access internal process handle
- `os/signal.NotifyContext` — sets cancel cause with signal info
- `reflect` — iterator methods: `Type.Fields`, `Type.Methods`, `Value.Fields`, etc.
- `runtime/metrics` — new goroutine/scheduler metrics
- `testing.T.ArtifactDir` — write test output artifacts
- `testing.B.Loop` — no longer prevents inlining
- `image/jpeg` — new faster encoder/decoder
- `crypto/tls` — post-quantum key exchanges enabled by default

### Breaking Changes / GODEBUG
- `urlstrictcolons=1` (default) — rejects `http://::1/`, use `http://[::1]/`
- `httpcookiemaxnum=3000` — max cookies in HTTP headers
- `urlmaxqueryparams=10000` — max query params
- `cryptocustomrand=0` (default) — crypto funcs ignore random parameter
- `tracebacklabels` — goroutine labels in tracebacks
- `asynctimerchan` — will be removed in Go 1.27

---

## 4. Files Requiring Version Bumps

### 4.1 Core Build Files

| # | File | Current | New | Change |
|---|------|---------|-----|--------|
| 1 | `go.mod` | `go 1.25.7` | `go 1.26.0` | go directive |
| 2 | `Makefile:13` | `export GOEXPERIMENT=greenteagc,jsonv2` | `export GOEXPERIMENT=jsonv2,goroutineleakprofile,simd,runtimesecret` | remove greenteagc, add new experiments |
| 3 | `Dockerfile:3` | `FROM golang:1.25-alpine AS builder` | `FROM golang:1.26-alpine AS builder` | base image |
| 4 | `Dockerfile:31` | `ENV GOEXPERIMENT=greenteagc,jsonv2` | `ENV GOEXPERIMENT=jsonv2,goroutineleakprofile,simd,runtimesecret` | same as Makefile |

### 4.2 Dev Container

| # | File | Current | New | Change |
|---|------|---------|-----|--------|
| 5 | `.devcontainer/Dockerfile:1` | `FROM golang:1.25-alpine` | `FROM golang:1.26-alpine` | base image |
| 6 | `.devcontainer/Dockerfile` | `ENV GOEXPERIMENT=greenteagc,jsonv2` | `ENV GOEXPERIMENT=jsonv2,goroutineleakprofile,simd,runtimesecret` | GOEXPERIMENT |
| 7 | `.devcontainer/devcontainer.json` | `"version": "1.25"` | `"version": "1.26"` | Go feature version |

### 4.3 GitHub Actions (5 workflows)

| # | File | Current | New | Change |
|---|------|---------|-----|--------|
| 8 | `.github/workflows/ci.yml` | `GO_VERSION: "1.25"` | `GO_VERSION: "1.26"` | env var |
| 9 | `.github/workflows/ci.yml` | 3x `image: golang:1.25-alpine` | `image: golang:1.26-alpine` | container images |
| 10 | `.github/workflows/ci.yml` | `GOEXPERIMENT: greenteagc,jsonv2` | `GOEXPERIMENT: jsonv2,goroutineleakprofile,simd,runtimesecret` | GOEXPERIMENT |
| 11 | `.github/workflows/develop.yml` | `GO_VERSION: "1.25"` | `GO_VERSION: "1.26"` | env var |
| 12 | `.github/workflows/develop.yml` | `GOEXPERIMENT: greenteagc,jsonv2` | `GOEXPERIMENT: jsonv2,goroutineleakprofile,simd,runtimesecret` | GOEXPERIMENT |
| 13 | `.github/workflows/release-please.yml` | `GO_VERSION: "1.25"` | `GO_VERSION: "1.26"` | env var |
| 14 | `.github/workflows/release-please.yml` | `GOEXPERIMENT: greenteagc,jsonv2` | `GOEXPERIMENT: jsonv2,goroutineleakprofile,simd,runtimesecret` | GOEXPERIMENT |
| 15 | `.github/workflows/security.yml` | `GO_VERSION: "1.25"` | `GO_VERSION: "1.26"` | env var |
| 16 | `.github/workflows/security.yml` | `image: golang:1.25-alpine` | `image: golang:1.26-alpine` | container image |
| 17 | `.github/workflows/security.yml` | `GOEXPERIMENT: greenteagc,jsonv2` | `GOEXPERIMENT: jsonv2,goroutineleakprofile,simd,runtimesecret` | GOEXPERIMENT |
| 18 | `.github/workflows/security.yml:123` | `GOEXPERIMENT: greenteagc  # exclude jsonv2...` | `GOEXPERIMENT: jsonv2,goroutineleakprofile,simd,runtimesecret` | **REMOVE WORKAROUND** — govulncheck panic fixed in 1.26 |
| 19 | `.github/workflows/coverage.yml` | `GO_VERSION: "1.25"` | `GO_VERSION: "1.26"` | env var |
| 20 | `.github/workflows/coverage.yml` | `GOEXPERIMENT: greenteagc,jsonv2` | `GOEXPERIMENT: jsonv2,goroutineleakprofile,simd,runtimesecret` | GOEXPERIMENT |

### 4.4 Coder Template

| # | File | Current | New | Change |
|---|------|---------|-----|--------|
| 21 | `.coder/template.tf:33` | `go_version = "1.25.6"` | `go_version = "1.26.0"` | Coder var |
| 22 | `.coder/template.tf:37` | `workspace_image = "golang:1.25-alpine"` | `workspace_image = "golang:1.26-alpine"` | image |
| 23 | `.coder/template.tf` | `GOEXPERIMENT = "greenteagc,jsonv2"` | `GOEXPERIMENT = "jsonv2,goroutineleakprofile,simd,runtimesecret"` | GOEXPERIMENT |

### 4.5 Editor Configuration

| # | File | Current | New | Change |
|---|------|---------|-----|--------|
| 24 | `.zed/tasks.json` | 6x `"GOEXPERIMENT": "greenteagc,jsonv2"` | `"GOEXPERIMENT": "jsonv2,goroutineleakprofile,simd,runtimesecret"` | all task envs |
| 25 | `.zed/settings.json:240` | `"GOEXPERIMENT": "greenteagc,jsonv2"` | `"GOEXPERIMENT": "jsonv2,goroutineleakprofile,simd,runtimesecret"` | terminal env |

### 4.6 Labels

| # | File | Current | New | Change |
|---|------|---------|-----|--------|
| 26 | `.github/labels.yml` | `go/1.25` label | Add `go/1.26` label | new label |

### 4.7 Documentation (content updates)

| # | File | What to Update |
|---|------|----------------|
| 27 | `.claude/CLAUDE.md:4` | "Go 1.25" → "Go 1.26" |
| 28 | `.claude/CLAUDE.md:180` | "Go 1.25" → "Go 1.26" + update feature list |
| 29 | `README.md:75` | Table row: Go 1.25+ → Go 1.26+, update GOEXPERIMENT |
| 30 | `docs/dev/design/technical/TECH_STACK.md:25` | Go 1.25.7 → Go 1.26.0, update GOEXPERIMENT |
| 31 | `docs/dev/design/operations/DEVELOPMENT.md:5,107` | Go 1.25+ → Go 1.26+, update GOEXPERIMENT |
| 32 | `docs/dev/design/operations/CI_CD.md:98-99` | GO_VERSION and GOEXPERIMENT |
| 33 | `docs/dev/design/operations/SETUP.md:18` | Update GOEXPERIMENT |
| 34 | `CONTRIBUTING.md:43` | Go 1.25+ → Go 1.26+ |
| 35 | `.githooks/pre-commit:12` | "Go 1.25+" → "Go 1.26+" |
| 36 | `.zed/docs/INDEX.md:27,71` | Go 1.25.6 → Go 1.26.0 |
| 37 | `.zed/docs/SETUP.md:108` | Go 1.25.6+ → Go 1.26.0+ |

---

## 5. Free Performance Wins (No Code Changes Needed)

These improvements apply automatically just by upgrading to Go 1.26:

| Improvement | Impact on This App |
|-------------|-------------------|
| Green Tea GC default | 10–40% less GC overhead (already using via experiment, now better integrated) |
| ~30% faster cgo calls | Directly benefits `go-astiav` (FFmpeg transcoding) and `govips` (image processing) |
| `io.ReadAll` ~2x faster | Any place reading full HTTP bodies, file contents |
| `fmt.Errorf` less alloc | Reduces allocations in error paths |
| Stack-allocated slices | Compiler optimizes more slice backing stores onto stack |
| Heap base randomization | ASLR-like security for cgo memory (go-astiav, govips) |
| JPEG encoder/decoder | Faster, more accurate image processing |

---

## 6. Code Changes to Make During Upgrade

### 6.1 `go fix ./...` (automatic)
Run the new `go fix` modernizers to auto-update code to Go 1.26 idioms.

### 6.2 `errors.As` → `errors.AsType` migration
Migrate existing `errors.As(err, &target)` calls to `errors.AsType[T](err)`.
Need to audit all `errors.As` usage first.

### 6.3 Potential `slog.NewMultiHandler` adoption
If the app has multi-handler slog setup, can simplify with stdlib.

### 6.4 Potential `new(expr)` adoption
For pointer-to-value fields in structs (JSON optional fields, etc.).

---

## 7. Dependency Compatibility Notes

### Direct Dependencies to Watch
All deps should be compatible with Go 1.26. Key ones:
- `go-astiav v0.40.0` — CGO, benefits from 30% faster cgo calls
- `govips/v2` — CGO, same benefit
- `ogen v1.18.0` — code generator, may need regen
- `pgx/v5 v5.8.0` — should be fine
- `rueidis v1.0.71` — should be fine
- `river v0.30.2` — should be fine
- `testcontainers-go v0.40.0` — should be fine

### golang.org/x packages (indirect)
- `golang.org/x/crypto v0.47.0` — may bump
- `golang.org/x/sys v0.40.0` — may bump
- `golang.org/x/text v0.33.0` — may bump
- `golang.org/x/net v0.48.0` — may bump
- `golang.org/x/sync v0.19.0` — may bump
- `golang.org/x/image v0.35.0` — may bump

These will be resolved by `go mod tidy` after updating `go.mod`.

---

## 8. GODEBUG Changes for Go 1.26

New GODEBUG defaults that apply automatically:
- `urlstrictcolons=1` — rejects malformed URLs with colons in host. **CHECK**: scan codebase for URL parsing that might break
- `httpcookiemaxnum=3000` — limits cookies. Should be fine.
- `urlmaxqueryparams=10000` — limits query params. Should be fine.
- `cryptocustomrand=0` — crypto funcs ignore random io.Reader param. Should be fine (we use standard crypto).
- Post-quantum TLS key exchanges enabled by default.

### Upcoming Removals (Go 1.27)
These GODEBUG settings will be removed in 1.27:
- `tlsunsafeekm`, `tlsrsakex`, `tls10server`, `tls3des`, `x509keypairleaf`
- `asynctimerchan` — timer channels become synchronous permanently

---

## 9. Implementation Order (When Ready)

1. **Update `go.mod`** — `go 1.26.0`
2. **Update all GOEXPERIMENT** values (Makefile, Dockerfile, devcontainer, CI, Coder, Zed)
3. **Update all version references** (Docker images, GO_VERSION vars, docs)
4. **Remove govulncheck workaround** in security.yml
5. **Add `go/1.26` label** to labels.yml
6. **Run `go mod tidy`** — resolve dependency changes
7. **Run `go fix ./...`** — apply modernizers
8. **Migrate `errors.As` → `errors.AsType`** — manual or via go fix if available
9. **Run tests** — `make test`
10. **Run linter** — `make lint`
11. **Test Docker build** — `make docker-build`
12. **Update documentation** — all .md files

---

## 10. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `urlstrictcolons` breaks URL parsing | Low | Medium | Grep for URL parsing, test |
| Dependency incompatibility with Go 1.26 | Low | Medium | `go mod tidy` + test |
| `go fix` changes semantics | Very Low | Low | Review diff before commit |
| JPEG encoder output differs | Low | Low | Only affects bit-exact comparisons |
| `simd` experiment instability | Low | Low | API unstable but build-gated |
| `runtimesecret` on non-Linux | N/A | None | Linux-only, no-op elsewhere |
