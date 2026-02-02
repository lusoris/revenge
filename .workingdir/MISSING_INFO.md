# Missing Information for Sources System

**Date**: 2026-02-02
**Purpose**: Track information that should be fetched/documented in `docs/dev/design/sources/`

---

## Currently Missing

### MISS-001: yamllint documentation
**Topic**: yamllint configuration and usage
**Why Needed**: We have `.yamllint.yml` but no source docs explaining the tool
**Priority**: Low
**Where**: Should be in `docs/dev/design/sources/` or tooling docs
**Action**: Add to fetch queue for doc-pipeline

---

### MISS-002: JSON Schema validation tools
**Topic**: jsonschema Python library and validation patterns
**Why Needed**: We have `schemas/` directory but no docs on validation approach
**Priority**: Medium
**Where**: Should document validation strategy
**Action**: Research and document

---

### MISS-003: pytest best practices for script testing
**Topic**: Patterns for testing Python scripts that modify files
**Why Needed**: Writing tests for sync-sot-status.py
**Priority**: Medium
**Where**: Could reference existing test patterns in `tests/unit/test_breadcrumbs.py`
**Action**: Review existing tests and document patterns

---

## Future Sources to Add

### FUT-001: GitHub Actions best practices
**Topic**: Workflow composition, job dependencies, caching strategies
**Why Needed**: Creating new workflows and enhancing existing ones
**Priority**: Low
**Source**: https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions

---

### FUT-002: Markdown table parsing libraries
**Topic**: Python libraries for parsing/manipulating Markdown tables
**Why Needed**: sync-sot-status.py needs to safely parse and update SOT tables
**Priority**: High (for current work)
**Research Needed**: Yes

---

## Resolved

_(Will be populated as information is found/documented)_

---

## Notes

- This file feeds into the doc-fetcher pipeline
- High priority items should be researched immediately
- Low priority items can be added to backlog
