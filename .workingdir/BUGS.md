# Bugs Found During Content Inconsistency Fixes

**Date**: 2026-02-02

---

## Critical Bugs

### BUG-001: valkey-go instead of rueidis
- **File**: `data/architecture/01_ARCHITECTURE.yaml:415`
- **Issue**: Wrong package `github.com/valkey-io/valkey-go` used instead of `github.com/redis/rueidis`
- **Impact**: Incorrect dependency information
- **Status**: Not Fixed
- **Fix**: Change to `github.com/redis/rueidis` per SOT line 159

---

## CI/Pipeline Flaws

### FLAW-001: yamllint not running in CI
- **File**: `.github/workflows/doc-validation.yml`
- **Issue**: `.yamllint.yml` config exists but yamllint is NOT run in any workflow
- **Impact**: YAML syntax errors not caught in CI
- **Status**: Not Fixed
- **Fix**: Add yamllint job to doc-validation.yml

### FLAW-002: Incomplete version validation
- **File**: `.github/workflows/validate-sot.yml`
- **Issue**: Only validates Go versions, not all packages (pgx, rueidis, otter, etc.)
- **Impact**: Hardcoded package versions in other files not detected
- **Status**: Not Fixed
- **Fix**: Extend to validate ALL packages from SOT

### FLAW-003: No JSON schema validation in CI
- **Issue**: `schemas/` directory exists but no schema validation in CI
- **Impact**: Invalid YAML files may pass CI
- **Status**: Not Fixed
- **Fix**: Add jsonschema validation to doc-validation.yml

---

## Discovered During Work

### BUG-004: Version drift in 7 markdown files
- **Files**: LIVE_TV_DVR.md, PHOTOS_LIBRARY.md, SKIP_INTRO.md, SYNCPLAY.md, TRICKPLAY.md, WATCH_NEXT_CONTINUE_WATCHING.md, TECH_STACK.md
- **Issue**: 60 version references out of sync with SOT (outdated package versions)
- **Impact**: Documentation shows wrong versions
- **Status**: ✅ FIXED with sync-versions.py --fix
- **Fix**: Ran `python scripts/sync-versions.py --fix`

### BUG-005: Systemic YAML indentation errors across ALL data files
- **Scope**: ALL YAML files in `data/` directory 
- **Issue**: yamllint reports wrong indentation for `sources:`, `design_refs:`, and nested structures
  - Example: Expected 2 spaces but found 0 for top-level keys
  - Example: Expected 4/6 spaces but found 2/4 for nested structures
- **Root Cause**: Either .yamllint.yml config is wrong OR all YAML files have incorrect indentation
- **Impact**: HIGH - This affects hundreds of files, not just 01_ARCHITECTURE.yaml
- **Status**: Found - This is systemic, not an isolated issue
- **Decision Needed**: 
  1. Fix .yamllint.yml config (if config is wrong)
  2. OR bulk-fix all YAML files (if files are wrong)
  3. OR disable yamllint indentation rule until proper fix
- **Recommendation**: This should be a SEPARATE task/PR, not part of current content fix

### BUG-006: Duplicate pgx versions in SOT
- **File**: `docs/dev/design/00_SOURCE_OF_TRUTH.md`
- **Issue**: 
  - Line 145: `pgx/v5` listed as v5.7.5 (Infrastructure table)
  - Line 157: `github.com/jackc/pgx/v5` listed as v5.8.0 (Dependencies table)
- **Impact**: sync-versions.py picks up BOTH entries causing permanent drift
- **Status**: Found during --strict mode testing
- **Fix**: Remove duplicate entry, keep only `github.com/jackc/pgx/v5` v5.8.0

---

_(Will be populated as bugs are fixed)_

---

## Notes

- All bugs should have corresponding tests added after fix
- Regression tests are mandatory for all fixes
