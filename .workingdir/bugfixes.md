# Bugfixes - v0.1.0 Implementation

> Issues encountered and their fixes (for future test creation)

---

## Issues

### [ISSUE-001] TestDefault fails - Database.URL has no default value
**Problem**: TestDefault test failure in CI: "Database.URL should have a default value"
**Cause**: Both `Defaults()` function (config.go:170) and `Default()` function (module.go:37) returned empty string for Database.URL
**Fix**: Changed both locations to use placeholder: `"postgres://revenge:changeme@localhost:5432/revenge?sslmode=disable"`
**Test Hint**: Test that default config values match expected placeholder values
**Files Changed**:
- [internal/config/config.go:170](internal/config/config.go#L170)
- [internal/config/module.go:37](internal/config/module.go#L37)

---

### [ISSUE-002] fx dependency injection - NewPool requires context.Context
**Problem**: Binary fails at startup with "missing type: context.Context" error
**Cause**: `NewPool(ctx context.Context, cfg *config.Config, logger *slog.Logger)` requires context as parameter, but fx cannot provide it automatically
**Fix**: Changed NewPool to create context internally with 30s timeout instead of accepting it as parameter
**Test Hint**: Test that NewPool creates pool successfully without external context
**Files Changed**:
- [internal/infra/database/pool.go:51](internal/infra/database/pool.go#L51)
- [cmd/revenge/migrate.go:35](cmd/revenge/migrate.go#L35)
- [internal/testutil/database.go:102](internal/testutil/database.go#L102)
- [internal/testutil/containers.go:119](internal/testutil/containers.go#L119)

---

### [ISSUE-003] Duplicate function declarations in wrap.go
**Problem**: Build fails with redeclared functions (Wrap, Unwrap, Is, As, New, Errorf)
**Cause**: wrap.go duplicated functions already present in errors.go
**Fix**: Removed duplicate basic functions, kept only additional utilities (Wrapf, WithStack, WrapSentinel, FormatError)
**Test Hint**: Test error wrapping utilities preserve stack traces correctly
**Files Changed**:
- [internal/errors/wrap.go](internal/errors/wrap.go)

---

### [ISSUE-004] testify assertion signature mismatch
**Problem**: Build fails with "too many arguments in call to assert.GreaterOrEqual/LessOrEqual"
**Cause**: assert/require methods use variadic ...interface{} for all args after first two, not separate message + args
**Fix**: Changed to use append([]interface{}{"message"}, msgAndArgs...)... pattern
**Test Hint**: Test custom time assertions work correctly with messages
**Files Changed**:
- [internal/testutil/assertions.go:71-72](internal/testutil/assertions.go#L71)
- [internal/testutil/assertions.go:86-87](internal/testutil/assertions.go#L86)

---

### [ISSUE-005] logging.NewLogger type mismatch in containers.go
**Problem**: Build fails with "cannot use cfg.Logging as logging.Config"
**Cause**: logging.NewLogger expects logging.Config, not config.LoggingConfig
**Fix**: Manually construct logging.Config from cfg.Logging fields
**Test Hint**: Test that testcontainer logger initialization works
**Files Changed**:
- [internal/testutil/containers.go:112](internal/testutil/containers.go#L112)

---

### [ISSUE-006] golangci-lint-action@v4 incompatible with golangci-lint v2.x
**Problem**: CI/CD "Lint & Code Quality" job fails with "Error: unknown flag: --out-format"
**Root Cause**: golangci-lint v2.x completely removed the `--out-format` flag. The golangci-lint-action@v4 is designed for v1.x and automatically adds `--out-format=github-actions` which doesn't exist in v2.x
**v1 vs v2 Syntax**:
  - v1: `golangci-lint run --out-format=json`
  - v2: `golangci-lint run --output.json.path=stdout`
**Fix**: Replaced golangci-lint-action@v4 with manual installation and execution using v2 syntax
**Test Hint**: Verify CI workflow completes successfully with golangci-lint v2.8.0
**Files Changed**:
- [.github/workflows/dev.yml:43-47](../.github/workflows/dev.yml#L43-L47)
**Workflow Runs**:
  - [21563257750](https://github.com/lusoris/revenge/actions/runs/21563257750) - Initial failure
  - [21563321778](https://github.com/lusoris/revenge/actions/runs/21563321778) - Still failed (action auto-adds flag)
**Sources**:
  - [v2 upgrade discussion](https://github.com/golangci/golangci-lint/discussions/5703)
  - [Output formats error](https://github.com/golangci/golangci-lint/issues/5605)
**Status**: FIXED - Using direct installation instead of action

---

### [ISSUE-007] 1,122 broken internal links in documentation
**Problem**: Documentation validation found 1,122 broken internal links out of 2,359 total links (47.5% broken)
**Cause**: Massive markdown restructuring and fixes may have broken relative links between documents
**Details**:
  - Files scanned: 214
  - Total internal links: 2,359
  - Broken links: 1,122
  - Auto-fixable (high confidence): 0
  - Placeholders found: 0
**Fix**: Need to review and fix broken links manually or improve auto-fix script
**Report Location**: docs/dev/design/.analysis/FIXES_REPORT.md (generated in CI, not available locally)
**Test Hint**: Run `python scripts/doc-pipeline/05-fix.py --report` to generate report locally
**Status**: WARNING - Not blocking but needs attention
**Workflow Run**: [21563257717](https://github.com/lusoris/revenge/actions/runs/21563257717)

---
