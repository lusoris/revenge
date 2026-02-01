## Table of Contents

- [Configuration Reference](#configuration-reference)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Components](#components)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)



---
sources:
  - name: koanf
    url: ../sources/tooling/koanf.md
    note: Configuration management
  - name: go-playground/validator
    url: ../sources/tooling/validator.md
    note: Config validation
design_refs:
  - title: technical
    path: technical/INDEX.md
  - title: TECH_STACK
    path: technical/TECH_STACK.md
  - title: 00_SOURCE_OF_TRUTH
    path: 00_SOURCE_OF_TRUTH.md
---

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



---


## Architecture

<!-- Architecture diagram placeholder -->

### Components

<!-- Component description -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

<!-- Interface definitions -->

### Dependencies

<!-- Dependency list -->





## Configuration
### Environment Variables

<!-- Environment variables -->

### Config Keys

<!-- Configuration keys -->




## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [technical](technical/INDEX.md)
- [TECH_STACK](TECH_STACK.md)
- [00_SOURCE_OF_TRUTH](../00_SOURCE_OF_TRUTH.md)

### External Sources
- [koanf](../sources/tooling/koanf.md) - Configuration management
- [go-playground/validator](../sources/tooling/validator.md) - Config validation

