# TODO: Documentation & Sources Pipeline Fixes

**Created**: 2026-02-01
**Priority**: HIGH
**Status**: In Progress

---

## Overview

Multiple critical issues discovered during CI/CD monitoring that need systematic fixes:

1. Missing sources for tools we use (had to web search for golangci-lint)
2. Missing changelog sources for ALL dependency versions
3. Wiki template URL generation broken (internal URLs don't link to actual sources)
4. 1,122 broken internal links in documentation (47.5% of all links!)

---

## Task 1: Add Missing Tool Sources to SOURCES.yaml

**Priority**: HIGH - We shouldn't need to web search for tools we use regularly

### Tools Missing from Sources

Based on recent work, add these to `docs/dev/sources/SOURCES.yaml`:

- [ ] **golangci-lint**
  - Releases: `https://github.com/golangci/golangci-lint/releases`
  - Changelog: `https://golangci-lint.run/docs/product/changelog/`
  - Docs: `https://golangci-lint.run/docs/`
  - v2 Migration: `https://github.com/golangci/golangci-lint/discussions/5703`

- [ ] **markdownlint-cli2**
  - Releases: `https://github.com/DavidAnson/markdownlint-cli2/releases`
  - Changelog: `https://github.com/DavidAnson/markdownlint-cli2/blob/main/CHANGELOG.md`
  - Rules: `https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md`

- [ ] **testcontainers-go**
  - Releases: `https://github.com/testcontainers/testcontainers-go/releases`
  - Changelog: `https://github.com/testcontainers/testcontainers-go/blob/main/CHANGELOG.md`
  - Docs: `https://golang.testcontainers.org/`

### Verification

After adding:
- [ ] Run `python scripts/fetch-sources.py` to verify fetching works
- [ ] Check that sources appear in `docs/dev/sources/`
- [ ] Update SOURCES_INDEX.md
- [ ] Verify cross-references in DESIGN_CROSSREF.md

---

## Task 2: Add Changelog Sources for ALL Dependency Versions

**Priority**: HIGH - Need automated tracking of breaking changes

### Philosophy

For EVERY package/tool with a version in SOURCE_OF_TRUTH, we should have:
1. Release page URL
2. Changelog URL
3. Documentation URL
4. Migration guides (if applicable)

### Categories to Cover

#### Go Dependencies (Core)

From SOURCE_OF_TRUTH section "Go Dependencies (Core)":

- [ ] **koanf** (v2.3.0)
  - Releases, changelog, migration from v1

- [ ] **validator** (v10.28.0)
  - Releases, changelog, docs

- [ ] **pgx/v5** (v5.8.0)
  - Releases, changelog, v4→v5 migration guide

- [ ] **ogen** (v1.18.0)
  - Releases, changelog, docs

- [ ] **fx** (v1.24.0)
  - Releases, changelog, docs

- [ ] **zap** (v1.27.1)
  - Releases, changelog, docs

#### Go Dependencies (Caching)

- [ ] **otter** (v2.3.1)
- [ ] **rueidis** (v1.0.54)
- [ ] **sturdyc** (v1.4.0)

#### Go Dependencies (Search)

- [ ] **typesense** (client + server versions)

#### Go Dependencies (Jobs)

- [ ] **river** (v0.20.1)
- [ ] **riverdriver/riverpgxv5** (v0.20.1)

#### Infrastructure

- [ ] **PostgreSQL** (18.x)
  - Release notes for 18.0, 18.1, etc.
  - Breaking changes from 17→18

- [ ] **Dragonfly** (v1.26.1)
  - Releases, changelog

- [ ] **Go** (1.25.6)
  - Release notes
  - Breaking changes from 1.24

#### Frontend

- [ ] **SvelteKit** (2.x)
- [ ] **Svelte** (5.x)
- [ ] **Tailwind CSS** (4.x)
- [ ] **shadcn-svelte**
- [ ] **TanStack Query**

#### Development Tools

- [ ] **sqlc** (v1.30.0)
- [ ] **mockery** (v3.3.0)
- [ ] **migrate** (v4.19.1)

### Implementation Plan

1. **Audit SOURCE_OF_TRUTH**: Extract ALL packages with versions
2. **Create SOURCES entries**: For each package, find official sources
3. **Test fetching**: Verify all URLs are fetchable
4. **Document format**: Standardize changelog extraction
5. **Automate checks**: Add CI job to verify sources exist for all SOT versions

---

## Task 3: Fix Wiki Template URL Generation

**Priority**: CRITICAL - "Major flaw" per user

### Problem

Wiki template is not creating links to actual sources. Internal URLs remain as internal references instead of linking to external documentation.

### Investigation Needed

- [ ] **Identify the template**: Find wiki template file
  - Likely in `docs/dev/design/.templates/`
  - Check `test_output_wiki.md`

- [ ] **Analyze URL generation**: How are internal URLs supposed to become external links?
  - Review template variables
  - Check pipeline scripts (doc-pipeline/)

- [ ] **Find broken examples**: Get specific examples of broken links
  - Compare output in wiki vs expected
  - Document the URL transformation that should happen

### Fix Implementation

- [ ] Update template to properly generate external links
- [ ] Update pipeline script if needed
- [ ] Test on sample document
- [ ] Regenerate all wiki outputs
- [ ] Verify external links work

---

## Task 4: Fix 1,122 Broken Internal Links

**Priority**: HIGH - 47.5% of all documentation links are broken!

### Current State

From Documentation Validation workflow (run 21563257717):

```
Files scanned: 214
Total internal links: 2,359
Broken links: 1,122 (47.5%)
Auto-fixable (high confidence): 0
Placeholders found: 0
```

Report: `docs/dev/design/.analysis/FIXES_REPORT.md` (generated in CI)

### Root Cause Analysis

Likely causes:
1. Massive markdown auto-fix (134k→0 errors) may have broken relative paths
2. File moves/renames not reflected in links
3. Anchor links to non-existent sections
4. Case sensitivity issues (Linux vs macOS)

### Investigation Steps

- [ ] **Get the report locally**:
  ```bash
  python scripts/doc-pipeline/05-fix.py --report
  ```

- [ ] **Analyze broken link patterns**:
  - Group by error type
  - Identify most common broken targets
  - Check if recent file moves caused bulk breakage

- [ ] **Categorize fixes needed**:
  - Simple path fixes (file moved)
  - Anchor fixes (section renamed/removed)
  - File deletions (need to remove links)
  - Case sensitivity fixes

### Fix Strategy

**Option A: Automated Fix Script Enhancement**

Currently 0 are auto-fixable. Why?

- [ ] Review `scripts/doc-pipeline/05-fix.py` logic
- [ ] Add fuzzy matching for file moves
- [ ] Add anchor validation
- [ ] Improve confidence scoring
- [ ] Re-run to see how many become auto-fixable

**Option B: Manual Fix + Documentation**

For links that can't be auto-fixed:

- [ ] Create categorized fix list
- [ ] Fix systematically by category
- [ ] Document link conventions to prevent future breaks
- [ ] Add pre-commit hook to validate links

**Option C: Hybrid Approach** (RECOMMENDED)

1. Improve auto-fix script first
2. Auto-fix what we can (target: >80%)
3. Manually fix remainder
4. Add validation to prevent regression

### Implementation Plan

1. [ ] Generate report locally
2. [ ] Analyze top 10 broken link patterns
3. [ ] Enhance auto-fix script for those patterns
4. [ ] Run auto-fix
5. [ ] Manually fix remainder
6. [ ] Add link validation to CI
7. [ ] Document link guidelines

---

## Task 5: Add Doc Pipeline PR Creation for Changes

**Priority**: MEDIUM - Automation improvement

### Goal

When doc pipeline detects changes from fetched sources:
- Auto-create PR with changes
- Include diff summary
- Link to upstream changelogs
- Label for review

### Implementation

- [ ] Add GitHub Actions workflow for doc updates
- [ ] Use `gh pr create` in pipeline
- [ ] Generate meaningful commit messages
- [ ] Add labels: `documentation`, `automated`, `sources`
- [ ] Configure auto-merge rules (or require review?)

### Integration Points

- Documentation Validation workflow
- Source fetcher script
- Version comparison logic

---

## Success Criteria

### Task 1: Sources
- [ ] No need to web search for any tool in SOURCE_OF_TRUTH
- [ ] All tools have: releases, changelog, docs URLs
- [ ] Fetcher successfully retrieves all sources

### Task 2: Changelogs
- [ ] Every version in SOT has changelog source
- [ ] Breaking changes automatically detected
- [ ] Migration guides linked where applicable

### Task 3: Wiki URLs
- [ ] All internal URLs in wiki link to actual sources
- [ ] Template generates correct external links
- [ ] Verified across multiple doc types

### Task 4: Broken Links
- [ ] < 5% broken links (down from 47.5%)
- [ ] Auto-fix handles >80% of common patterns
- [ ] CI prevents new broken links
- [ ] Documentation on link conventions

### Task 5: Automation
- [ ] Doc changes auto-create PRs
- [ ] PRs include meaningful context
- [ ] Process documented for reviewers

---

## Next Steps

1. **Immediate**: Generate broken links report locally
2. **Today**: Add missing tool sources to SOURCES.yaml
3. **This week**: Fix wiki template URL generation
4. **This week**: Fix broken links (auto + manual)
5. **Next week**: Add PR automation for doc updates

---

## Notes

- User emphasized this is high priority ("finally")
- Wiki template issue is "major flaw"
- Need systematic approach, not one-off fixes
- Goal: Never manually search for docs we already use

---

## Related Files

- `docs/dev/sources/SOURCES.yaml` - Source registry
- `docs/dev/design/00_SOURCE_OF_TRUTH.md` - All versions
- `scripts/fetch-sources.py` - Fetcher script
- `scripts/doc-pipeline/05-fix.py` - Link fixer
- `docs/dev/design/.analysis/FIXES_REPORT.md` - Broken links report
- `.workingdir/bugfixes.md` - ISSUE-007 details
