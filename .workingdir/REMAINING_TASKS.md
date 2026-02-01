# Remaining Tasks (Post Content Gaps)

**Status**: On hold while filling YAML content gaps
**Last Updated**: 2026-02-01 17:00

---

## High Priority

### 1. Fix Status Overview Syncing Bugs
**Status**: Not Started
**Complexity**: Medium
**Description**: Investigate and fix bugs in status table synchronization across design docs

**Tasks**:
- [ ] Identify what's broken in status sync
- [ ] Check `scripts/doc-pipeline/03-status.py`
- [ ] Create tests for status sync
- [ ] Fix synchronization logic
- [ ] Verify all status tables are in sync

**Files Involved**:
- `scripts/doc-pipeline/03-status.py`
- Design docs with status tables

---

### 2. Verify All Remote Workflows Pass
**Status**: In Progress (monitoring)
**Current State**: 
- Code Coverage: in_progress
- Security Scanning: in_progress
- CodeQL: in_progress
- Development Build: queued
- Deploy Testing (k8s): failure (expected - no Helm charts)

**Tasks**:
- [ ] Monitor workflow completion
- [ ] Fix any failures in core workflows
- [ ] Deploy workflows will fail until Helm charts are created (v0.2.0+)
- [ ] Document expected failures vs actual issues

---

## Medium Priority

### 3. Create Additional Tests
**Status**: Partially Complete
**Progress**: 15/? test packages

**Remaining Go Packages**:
- [ ] `cmd/revenge` - Main command tests
- [ ] `internal/testutil` - Test utility tests
- [ ] More comprehensive integration tests

**Remaining Python**:
- [ ] `scripts/doc-pipeline/` - All 6 pipeline scripts
- [ ] `scripts/automation/` - Remaining 15+ scripts
- [ ] `scripts/deploy-pipeline/` - Deployment pipeline
- [ ] `scripts/generate-*.py` - Index, crossref, navigation generators

**Test Coverage Goals**:
- Go: 80% minimum (currently ~60%)
- Python: 70% minimum (currently ~40%)

---

### 4. Feature Testing (v0.1.0 Scope)
**Status**: Not Started
**Depends On**: All workflows passing

**Tasks**:
- [ ] Build binary: `GOEXPERIMENT=greenteagc,jsonv2 go build ./cmd/revenge`
- [ ] Test server startup
- [ ] Test health endpoints manually
  - GET /health/live
  - GET /health/ready
  - GET /health/startup
- [ ] Verify configuration loading
- [ ] Test database connection
- [ ] Test graceful shutdown
- [ ] Document any issues found

---

## Low Priority / Future

### 5. Create Helm Charts (v0.2.0)
**Status**: Not Started
**Blocked By**: v0.1.0 completion

**Tasks**:
- [ ] Create `charts/revenge/` structure
- [ ] Define values.yaml schema
- [ ] Create deployment templates
- [ ] Create service templates
- [ ] Create ingress templates
- [ ] Create ConfigMap templates
- [ ] Add horizontal pod autoscaling
- [ ] Test with kind/k3d
- [ ] Document deployment

---

### 6. Documentation Improvements
**Status**: Ongoing

**Tasks**:
- [ ] Add more code examples to design docs
- [ ] Create deployment guides
- [ ] Add troubleshooting sections
- [ ] Create API documentation (OpenAPI/Swagger)
- [ ] Add architecture diagrams
- [ ] Record demo videos/GIFs

---

## Completed ✅

- [x] Fix Dragonfly Helm repo URLs
- [x] Add GitHub Actions to SOURCE_OF_TRUTH
- [x] Add Helm repos to SOURCES.yaml
- [x] Fix doc generation pipeline (TOC + frontmatter)
- [x] Regenerate all 158 docs
- [x] Create tests for TOC generator
- [x] Create tests for app, cache, jobs, search modules
- [x] Add logging.NewTestLogger() helper
- [x] Commit and push all changes (6 commits)

---

## Notes

- **Token Budget**: ~86k remaining (as of last check)
- **Current Focus**: Filling ALL YAML content gaps (158 files)
- **Blocker**: None - all dependencies resolved
- **Risk Areas**: 
  - YAML content quality/consistency
  - Status sync bugs (unknown complexity)
  - Deployment workflows (need Helm charts)

---

## Commands Reference

```bash
# Run all tests
GOEXPERIMENT=greenteagc,jsonv2 go test ./...

# Run linters
golangci-lint run ./...
ruff check scripts/

# Regenerate docs
python -m scripts.automation.batch_regenerate

# Run doc pipeline
./scripts/doc-pipeline.sh --apply

# Check workflows
gh run list --branch develop --limit 10

# Feature testing
GOEXPERIMENT=greenteagc,jsonv2 go build ./cmd/revenge
./revenge --help
```

---

**Remember**: Focus on YAML content gaps FIRST, then return to this file for next priorities.
