# Questions & Gaps to Discuss

> Temporary file for tracking discrepancies, gaps, and questions that need resolution

**Created**: 2026-01-30
**Status**: 🔄 Active collection

---

## Critical Discrepancies Found

### Package Version Mismatches

| Package | go.mod (actual) | SOURCES.yaml | SOURCE_OF_TRUTH.md | Resolution Needed |
|---------|-----------------|--------------|---------------------|-------------------|
| otter | v1.2.4 | v2.3.0 | v1.2.4 ✅ | Update SOURCES.yaml |
| typesense-go | v4.0.0-alpha2 | v3.2.0 ("NO v4 exists!") | v3.2.0 | **VERIFY**: Is v4 alpha stable enough? |
| resty | go-resty/resty/v2 v2.17.1 | resty v3.0.0-b6 | v2.17.1 ✅ | SOURCES.yaml references wrong package |

### Questions

1. **Typesense Go Client**:
   - go.mod uses `typesense-go/v4 v4.0.0-alpha2`
   - SOURCES.yaml explicitly says "NO v4 exists!"
   - Is v4 alpha suitable for production? Or should we downgrade to v3.2.0?
   - Need to verify latest stable from live docs

2. **Resty HTTP Client**:
   - go.mod uses v2 (`go-resty/resty/v2`)
   - SOURCES.yaml mentions v3 beta (`resty.dev/v3`)
   - Should we upgrade to v3 when stable? What's the migration effort?

---

## Missing from SOURCE_OF_TRUTH

### Not Yet Documented
- [ ] Full API endpoint inventory
- [ ] Complete database table list
- [ ] All environment variables with defaults
- [ ] OpenAPI spec file locations
- [ ] Migration file naming conventions
- [ ] Test coverage requirements
- [ ] CI/CD pipeline stages

### Need Design Docs
- [ ] Fingerprint service design doc
- [ ] Grants service design doc
- [ ] Search service design doc
- [ ] Analytics service design doc
- [ ] Notification service design doc

---

## Outdated References Found

### Old Namespace 'c' → 'qar'
- [x] Already fixed in previous session

### ARCHITECTURE_V2.md → ARCHITECTURE.md
- [x] Already fixed in this session

### Other Potential Issues
- [ ] Some adult module docs still reference `c` schema? (need verification)
- [ ] INDEX.md still has `c` schema reference? (line 90: `adult/ADULT_CONTENT_SYSTEM.md` mentions `c` schema)

---

## Questions for Project Owner

1. **Database Strategy**:
   - Is SQLite truly needed for single-user deployments?
   - What's the priority for dual DB support vs other features?

2. **Typesense Version**:
   - v4 alpha in go.mod - intentional bleeding edge?
   - Should we pin to v3.2.0 stable?

3. **Package Update Policy**:
   - How aggressive on updates? (bleeding edge vs 1 version behind)
   - Who monitors changelogs/breaking changes?

4. **Documentation Priority**:
   - Should we complete SOURCE_OF_TRUTH first or start scaffolding?
   - User clarified: docs first, no scaffolding until consistent ✅

---

## Live Docs Verification Needed

These need to be verified against actual package/API documentation:

| Source | Last Verified | Action |
|--------|---------------|--------|
| Go 1.25.6 release notes | Never | Check new features |
| pgx v5.8.0 changelog | Never | Breaking changes? |
| River v0.30.2 changelog | Never | Breaking changes? |
| ogen v1.18.0 changelog | Never | Breaking changes? |
| rueidis v1.0.71 | Never | Breaking changes? |
| Typesense API v30 | Never | Verify latest |
| StashDB GraphQL schema | Never | Schema changes? |

---

## Resolution Log

| Date | Item | Resolution |
|------|------|------------|
| 2026-01-30 | otter version | SOURCES.yaml was wrong, go.mod is correct |
| 2026-01-30 | resty version | SOURCES.yaml references v3 beta, we use stable v2 |
| | | |

---

## Notes

- This file is temporary - items should be resolved and moved to appropriate docs
- When resolved, update SOURCE_OF_TRUTH.md and remove from here
- Questions for owner should be asked and answers documented in design docs
