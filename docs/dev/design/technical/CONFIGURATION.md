## Table of Contents

- [Configuration Reference](#configuration-reference)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Configuration Reference


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > Configuration system using koanf (YAML + environment variables + hot reload)

Revenge configuration management:
- **koanf**: Unified config from YAML files, env vars, flags
- **Hot Reload**: Runtime config updates (subset of keys)
- **Validation**: Schema validation with go-playground/validator
- **Secrets**: Environment variable expansion and file-based secrets
- **Env Prefix**: All env vars use `REVENGE_` prefix


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete configuration reference |
| Sources | ✅ | All config tools documented |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete




## Related Documentation
### Design Documents
- [technical](INDEX.md)
- [TECH_STACK](../technical/TECH_STACK.md)
- [00_SOURCE_OF_TRUTH](../00_SOURCE_OF_TRUTH.md)

### External Sources
- [koanf](../../sources/tooling/koanf.md) - Configuration management
- [go-playground/validator](../../sources/tooling/validator.md) - Config validation

