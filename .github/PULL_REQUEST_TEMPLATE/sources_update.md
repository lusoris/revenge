---
name: External Sources Update
about: Automated PR for external documentation sources
title: 'docs: Update external documentation sources'
labels: documentation, automated
---

## Automated Documentation Update

This PR updates external documentation sources fetched from upstream.

### Fetch Summary

<!-- Auto-populated by workflow -->

### Changes

- Updated documentation from external sources
- See `docs/dev/sources/INDEX.yaml` for fetch status

### Review Checklist

- [ ] No unexpected content changes in fetched files
- [ ] No modifications to `docs/dev/design/` (protected)
- [ ] Check for any failed fetches in INDEX.yaml
- [ ] Ensure no sensitive data was accidentally fetched
- [ ] Verify links in fetched content are still valid

### ⚠️ Important

**This PR should ONLY modify files in `docs/dev/sources/`.**

If any files in `docs/dev/design/` are modified, **reject this PR** - the fetcher has a bug or was misconfigured.

---
*This PR was automatically created by the fetch-sources workflow.*
