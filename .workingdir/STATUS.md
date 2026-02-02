# Content Inconsistency Fix - Status Tracker

**Started**: 2026-02-02
**Goal**: Fix critical content errors, build status sync, consolidate duplicates

---

## Progress Overview

- [ ] **Phase 1**: Fix Critical Issues (valkey-go, versions)
- [ ] **Phase 2**: Build Status Sync Script with Tests
- [ ] **Phase 3**: Integrate into Pipeline & CI
- [ ] **Phase 4**: Run Status Sync & Verify
- [ ] **Phase 5**: Consolidate Duplicate Content

---

## Phase 1: Fix Critical Issues

### Step 1.1: Fix valkey-go → rueidis
- [x] Update line 415 in `data/architecture/01_ARCHITECTURE.yaml`
- [x] Run `scripts/sync-versions.py` to verify no other issues (FOUND 60 drifts, FIXED)
- [x] Run yamllint on the file (FIXED line-length warning; indentation tracked as BUG-005)
- [x] Commit changes

**Completed**: Phase 1, Step 1.1 ✅

### Step 1.2: Enhance sync-versions.py
- [x] Add `--strict` mode that errors on violations
- [x] Test on current codebase (found BUG-006 duplicate pgx version, fixed)
- [x] Fix any violations found (fixed SOT duplicate)
- [x] Add tests for the new mode (6 tests, all passing)
- [x] Run tests until passing
- [x] Commit changes

**Completed**: Phase 1, Step 1.2 ✅

---

## Phase 2: Build Status Sync Script

### Step 2.1: Create sync-sot-status.py
- [x] Create `scripts/doc-pipeline/04-sync-sot-status.py`
- [x] Implement backup creation
- [x] Implement SOT table parsing
- [x] Implement status-only updates (safe mode)
- [x] Add dry-run mode
- [x] Test dry-run - found 13 status mismatches

**Status**: Complete ✅
**Findings**: Script detected 13 status inconsistencies including 4 undocumented in Arr Ecosystem table

### Step 2.2: Write tests for sync-sot-status.py
- [x] Create `tests/unit/test_sync_sot_status.py`
- [x] Test backup functionality
- [x] Test SOT table parsing
- [x] Test status-only updates
- [x] Test dry-run vs apply modes
- [x] Test strict mode
- [x] Run all tests until passing

**Status**: Complete ✅ (10/10 tests passing)

---

## Phase 3: Integrate into Pipeline & CI

### Step 3.1: Update doc-pipeline.sh
- [x] Add status sync step to scripts/doc-pipeline.sh
- [x] Place after status table validation, before cross-reference generation
- [x] Test pipeline execution (manual test successful)

**Status**: Complete ✅

### Step 3.2: Enhance CI/CD workflows
- [x] Add yamllint job to doc-validation.yml (fix FLAW-001)
- [x] Add status sync verification job
- [x] Add status sync to dry-run pipeline
- [ ] Extend validate-sot.yml to check ALL package versions (address FLAW-002)
- [ ] Test workflow execution

**Current**: Finishing Phase 3, Step 3.2

### Step 2.2: Write Tests
- [ ] Create test file in `tests/unit/`
- [ ] Test backup functionality
- [ ] Test parsing of SOT tables
- [ ] Test safe status updates
- [ ] Test error handling
- [ ] Run tests until all passing

---

## Phase 3: Pipeline Integration

### Step 3.1: Add to doc-pipeline.sh
- [ ] Update `scripts/doc-pipeline.sh`
- [ ] Test locally

### Step 3.2: Create GitHub Action
- [ ] Add yamllint job to `doc-validation.yml` (fix existing flaw)
- [ ] Add status sync job
- [ ] Test in CI

### Step 3.3: Enhance validate-sot.yml
- [ ] Extend version validation to all packages (not just Go)
- [ ] Add tests
- [ ] Verify

---

## Phase 4: Run Status Sync

### Step 4.1: Execute Sync
- [ ] Run `04-sync-sot-status.py`
- [ ] Review git diff
- [ ] Verify ONLY status cells changed
- [ ] Run all validation scripts

### Step 4.2: Fix Any Issues
- [ ] Document bugs found
- [ ] Fix script
- [ ] Add regression tests
- [ ] Re-run until clean

---

## Phase 5: Consolidate Duplicates

### Step 5.1: Metadata Priority Chain (~125 lines)
- [ ] Keep canonical in `03_METADATA_SYSTEM.yaml`
- [ ] Remove from 6 files, add cross-references
- [ ] Validate all files
- [ ] Check for broken links
- [ ] Verify content integrity

### Step 5.2: Arr Dual-Role Description (~300 lines)
- [ ] Keep canonical in `03_METADATA_SYSTEM.yaml`
- [ ] Remove from 6 files, add cross-references
- [ ] Validate, verify

### Step 5.3: Proxy/VPN Documentation (~425 lines)
- [ ] Create `data/patterns/HTTP_CLIENT.yaml`
- [ ] Move canonical content
- [ ] Update 5 files with cross-references
- [ ] Validate, verify

### Step 5.4: Cache Architecture (~60 lines)
- [ ] Keep canonical in `DRAGONFLY.yaml`
- [ ] Update 4 files with cross-references
- [ ] Validate, verify

---

## Current Status

**Current Phase**: Phase 1 - Setup
**Current Step**: Creating tracking files
**Blockers**: None
**Next Action**: Create other tracking files, then start Step 1.1
