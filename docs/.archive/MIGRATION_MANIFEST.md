# Documentation Migration Manifest

> Created: 2026-01-28
> Purpose: Backup reference for /docs restructuring to /docs/dev

## Migration Summary

| Metric | Count |
|--------|-------|
| Total Files | 121 |
| Architecture | 5 |
| Features | 12 |
| Integrations | 66 |
| Operations | 8 |
| Planning | 3 |
| Research | 3 |
| Technical | 5 |
| Root-level | 5 |
| Index files | 14 |

## Original File Locations

All files below will be moved from `/docs/` to `/docs/dev/design/`.

### Root-Level Files → design/
- DOCUMENTATION_ANALYSIS.md
- EXTERNAL_INTEGRATIONS_TODO.md
- INDEX.md
- PREPARATION_MASTER_PLAN.md
- RESTRUCTURING_PLAN.md

### architecture/ → design/architecture/
- ARCHITECTURE_V2.md
- DESIGN_PRINCIPLES.md
- METADATA_SYSTEM.md
- PLAYER_ARCHITECTURE.md
- PLUGIN_ARCHITECTURE_DECISION.md

### features/ → design/features/
- ADULT_CONTENT_SYSTEM.md
- ADULT_METADATA.md
- ANALYTICS_SERVICE.md
- CLIENT_SUPPORT.md
- COMICS_MODULE.md
- CONTENT_RATING.md
- I18N.md
- LIBRARY_TYPES.md
- MEDIA_ENHANCEMENTS.md
- REQUEST_SYSTEM.md
- SCROBBLING.md
- TICKETING_SYSTEM.md
- USER_EXPERIENCE_FEATURES.md
- WHISPARR_STASHDB_SCHEMA.md

### integrations/ → design/integrations/
All 66 integration docs remain in their current subcategory structure:
- anime/ (4 files)
- audiobook/ (2 files)
- auth/ (5 files)
- casting/ (3 files)
- external/ (9 files including adult/)
- infrastructure/ (5 files)
- livetv/ (3 files)
- metadata/ (24 files across subcategories)
- scrobbling/ (6 files)
- servarr/ (6 files)
- transcoding/ (2 files)
- wiki/ (8 files)

### operations/ → design/operations/
- BEST_PRACTICES.md
- BRANCH_PROTECTION.md
- DATABASE_AUTO_HEALING.md
- DEVELOPMENT.md
- GITFLOW.md
- REVERSE_PROXY.md
- SETUP.md
- UPSTREAM_SYNC.md

### planning/ → design/planning/
- MODULE_IMPLEMENTATION_TODO.md
- VERSION_POLICY.md
- VERSIONING.md

### research/ → design/research/
- DOCUMENTATION_GAP_ANALYSIS.md
- GO_PACKAGES_RESEARCH.md
- USER_PAIN_POINTS_RESEARCH.md

### technical/ → design/technical/
- API.md
- AUDIO_STREAMING.md
- FRONTEND.md
- OFFLOADING.md
- TECH_STACK.md

## New Structure

```
/docs/
├── INDEX.md                    # Router (new)
└── dev/
    ├── INDEX.md                # Dev docs index (new)
    ├── MIGRATION_MANIFEST.md   # This file
    ├── design/                 # 🔒 PROTECTED
    │   ├── INDEX.md
    │   ├── architecture/
    │   ├── features/
    │   ├── integrations/
    │   ├── operations/
    │   ├── planning/
    │   ├── research/
    │   └── technical/
    └── sources/                # 🔄 AUTO-FETCH
        ├── SOURCES.yaml
        ├── INDEX.yaml
        └── {category folders}
```

## Verification Checklist

- [ ] All 121 files accounted for
- [ ] No broken internal links
- [ ] INDEX.md files updated with new paths
- [ ] SOURCES.yaml created
- [ ] Fetcher script operational
