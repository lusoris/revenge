# Go 1.26 Upgrade — TODO

**GOEXPERIMENT:** `jsonv2,goroutineleakprofile,simd,runtimesecret`

---

## Phase 1: Core Build Files
- [x] 1.1 `go.mod` — bump `go 1.25.7` → `go 1.26.0`
- [x] 1.2 `Makefile:13` — GOEXPERIMENT update
- [x] 1.3 `Dockerfile` — base image `golang:1.26-alpine` + GOEXPERIMENT
- [x] 1.4 `go mod tidy` — resolve dependency changes

## Phase 2: Dev Container
- [x] 2.1 `.devcontainer/Dockerfile` — base image + GOEXPERIMENT
- [x] 2.2 `.devcontainer/devcontainer.json` — Go feature version

## Phase 3: CI Workflows
- [x] 3.1 `ci.yml` — GO_VERSION, container images (3x), GOEXPERIMENT
- [x] 3.2 `develop.yml` — GO_VERSION, GOEXPERIMENT
- [x] 3.3 `release-please.yml` — GO_VERSION, GOEXPERIMENT
- [x] 3.4 `security.yml` — GO_VERSION, container image, GOEXPERIMENT, **remove govulncheck workaround**
- [x] 3.5 `coverage.yml` — GO_VERSION, GOEXPERIMENT

## Phase 4: Coder & Editor Config
- [x] 4.1 `.coder/template.tf` — go_version, workspace_image, GOEXPERIMENT
- [x] 4.2 `.zed/tasks.json` — GOEXPERIMENT (6 occurrences)
- [x] 4.3 `.zed/settings.json` — GOEXPERIMENT

## Phase 5: Labels
- [x] 5.1 `.github/labels.yml` — add `go/1.26` label

## Phase 6: Documentation
- [x] 6.1 `.claude/CLAUDE.md` — Go version refs (lines 4, 180)
- [x] 6.2 `README.md` — tech table row (line 75)
- [x] 6.3 `docs/dev/design/technical/TECH_STACK.md` — version + GOEXPERIMENT
- [x] 6.4 `docs/dev/design/operations/DEVELOPMENT.md` — version + GOEXPERIMENT (lines 5, 107)
- [x] 6.5 `docs/dev/design/operations/CI_CD.md` — GO_VERSION + GOEXPERIMENT
- [x] 6.6 `docs/dev/design/operations/SETUP.md` — GOEXPERIMENT
- [x] 6.7 `CONTRIBUTING.md` — Go version
- [x] 6.8 `.githooks/pre-commit` — Go version
- [x] 6.9 `.zed/docs/INDEX.md` — Go version (lines 27, 71)
- [x] 6.10 `.zed/docs/SETUP.md` — Go version

## Phase 7: Code Modernization
- [x] 7.1 Run `go fix ./...` — 153 files updated (interface{} → any, for range N, //go:fix inline, alignment)
- [x] 7.2 Audit & migrate `errors.As` → `errors.AsType[T]()` — 4 sites in handler.go, 1 in s3.go, 1 in metadata/errors.go
- [x] 7.3 Fix `errors.Errorf` %w vet errors in errors_test.go (Go 1.26 stricter vet)
- [x] 7.4 Add `errors.AsType[E]` to internal/errors package (delegates to stdlib)
- [x] 7.5 Switch handler.go from go-faster/errors to internal/errors package

## Phase 8: Validation
- [x] 8.1 `go vet` — all non-CGO packages pass
- [x] 8.2 `go test` — all non-CGO packages pass (40+ package suites)
- [ ] 8.3 Full CI validation — requires Linux (CGO + FFmpeg + libvips)
- [ ] 8.4 Docker build verification — requires Linux
