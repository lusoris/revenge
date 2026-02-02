# Questions for User

**Date**: 2026-02-02

---

## Immediate Questions (Blocking Progress)

### Q1: HTTP_CLIENT.yaml exact structure
**Context**: Need to create `data/patterns/HTTP_CLIENT.yaml` for proxy/VPN consolidation

**Question**: Should HTTP_CLIENT.yaml follow the same structure as other pattern files (like ARR_INTEGRATION.yaml)?
- Include all standard fields (doc_title, doc_category, status fields, etc.)?
- Or simpler structure since it's just documentation consolidation?

**Impact**: Blocks Phase 5, Step 5.3

---

### Q2: Dragonfly version field name
**Context**: Need to add version reference to DRAGONFLY.yaml

**Question**: What field name should I use?
- A. `sot_version_ref: "v1.36.0 (see SOURCE_OF_TRUTH.md)"`
- B. Add to existing `sources:` section as a note
- C. Add to `technical_summary`
- D. Other?

**Impact**: Affects Phase 1, Step 1.1

---

## Nice-to-Have Clarifications

### Q3: Test running command
**Question**: What's the exact command to run tests?
- `pytest tests/`
- `pytest tests/ -v`
- `python -m pytest tests/`
- Other flags needed?

**Current Assumption**: `pytest tests/ -v --tb=short` (from doc-validation.yml line 133)

---

### Q4: Virtual environment activation
**Question**: After creating `.venv`, what's the activation command preference?
- Windows PowerShell: `.venv\Scripts\Activate.ps1`
- Windows CMD: `.venv\Scripts\activate.bat`
- Does it matter?

**Current Assumption**: PowerShell since that's the active terminal

---

## Answered Questions

### A1: HTTP_CLIENT.yaml structure
**Answer**: Must follow standard template structure like ARR_INTEGRATION.yaml
- Include all standard fields: doc_title, doc_category, created_date, overall_status, all status fields
- Include sources, design_refs, technical_summary, wiki_tagline, wiki_overview
- Pattern-specific fields as needed

### A2: Dragonfly version field
**Answer**: Must follow standard template - no custom fields
- Version info should NOT be hardcoded in YAML files
- All versions only in SOURCE_OF_TRUTH.md
- YAML files reference packages in `sources:` section only
- No version field needed in DRAGONFLY.yaml

---

## Parking Lot (For Later)

- Should we add more comprehensive CI checks beyond what we're fixing now?
- Should duplicate consolidation be done in smaller PRs or one large PR?
