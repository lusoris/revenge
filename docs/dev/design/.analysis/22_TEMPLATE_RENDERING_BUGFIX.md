# Template Rendering Bugfix Report

**Date**: 2026-01-31
**Severity**: Critical (prevented 67% of docs from generating)
**Status**: ✅ Fixed and Tested

---

## Executive Summary

Fixed critical bug in Jinja2 templates that caused 95 out of 142 documentation files (67%) to fail generation with `UndefinedError`. All templates now handle missing optional YAML fields gracefully using `| default()` filters.

**Result**: 142/142 files (100%) now generate successfully ✅

---

## Bug Discovery

### Initial Symptoms
- Wiki generation completely broken (only 1-2 files vs expected 142)
- `batch_regenerate.py` failing for 95 files when run in --apply mode
- Error: `'dependencies' is undefined` for all service files
- Error: `'api_base_url' is undefined` for integration files
- Error: `'content_types' is undefined` for feature files

### Root Cause Analysis

**Problem**: Jinja2 templates used strict variable checks (`{%- if variable %}`) without providing default values for optional YAML fields.

**Impact**: When YAML data files didn't include optional fields (which is valid), Jinja2 raised `UndefinedError` and template rendering failed.

**Example Bad Code**:
```jinja2
{%- if dependencies %}  <!-- Fails if 'dependencies' key doesn't exist in YAML -->
  ...
{%- endif %}
```

**Example Good Code**:
```jinja2
{%- if dependencies | default([]) %}  <!-- Safe: uses empty list if undefined -->
  ...
{%- endif %}
```

---

## Affected Templates and Fixes

### 1. service.md.jinja2 (15 files affected)

**Undefined Variables Fixed**:
- `dependencies` (line 54)
- `provides` (line 65)
- `wiki_how_it_works` (line 80)

**Fix Applied**:
```diff
- {%- if dependencies %}
+ {%- if dependencies | default([]) %}

- {%- if provides %}
+ {%- if provides | default([]) %}

- {{ wiki_how_it_works }}
+ {{ wiki_how_it_works | default('<!-- How it works -->') }}
```

**Files Fixed**: All 15 service YAML files
- ACTIVITY, ANALYTICS, APIKEYS, AUTH, FINGERPRINT, GRANTS, LIBRARY, METADATA, NOTIFICATION, OIDC, RBAC, SEARCH, SESSION, SETTINGS, USER

---

### 2. integration.md.jinja2 (58 files affected)

**Undefined Variables Fixed**:
- `api_base_url` (line 27)
- `auth_method` (line 28)
- `provides_data` (line 60)
- `auth_config` (line 84)
- `rate_limits` (line 99)
- `api_endpoints` (line 124)
- `cache_ttl` (line 162)
- `prerequisites` (line 196)
- `wiki_how_it_works` (line 70)

**Fix Applied**:
```diff
+ {%- if api_base_url | default('') %}
  **API Base URL**: `{{ api_base_url }}`
+ {%- endif %}
+ {%- if auth_method | default('') %}
  **Authentication**: {{ auth_method }}
+ {%- endif %}

- {%- if provides_data %}
+ {%- if provides_data | default([]) %}

- {%- if auth_config %}
+ {%- if auth_config | default('') %}

- {%- if claude and rate_limits %}
+ {%- if claude and (rate_limits | default({})) %}

- {%- if api_endpoints %}
+ {%- if api_endpoints | default([]) %}

- {%- if cache_ttl %}
+ {%- if cache_ttl | default({}) %}

- {%- if prerequisites %}
+ {%- if prerequisites | default([]) %}

- {{ wiki_how_it_works }}
+ {{ wiki_how_it_works | default('<!-- How it works -->') }}
```

**Files Fixed**: All 58 integration YAML files across:
- anime/ (3 files)
- auth/ (4 files)
- casting/ (2 files)
- infrastructure/ (4 files)
- livetv/ (3 files)
- metadata/adult/ (10 files)
- metadata/books/ (4 files)
- metadata/comics/ (3 files)
- metadata/music/ (4 files)
- metadata/video/ (4 files)
- scrobbling/ (5 files)
- servarr/ (5 files)
- transcoding/ (1 file)
- wiki/ (3 files)
- wiki/adult/ (3 files)

---

### 3. feature.md.jinja2 (22 files affected)

**Undefined Variables Fixed**:
- `content_types` (line 22)

**Fix Applied**:
```diff
- > Content module for {{ content_types | join(', ') }}
+ > Content module for {{ (content_types | default([])) | join(', ') }}
```

**Files Fixed**: 22 feature YAML files across:
- adult/ (5 files)
- audiobook/ (1 file)
- book/ (1 file)
- comics/ (1 file)
- livetv/ (1 file)
- music/ (1 file)
- photos/ (1 file)
- playback/ (6 files)
- podcasts/ (1 file)
- shared/ (14 files)
- video/ (2 files)

Note: Some features worked because they had `content_types` defined in YAML, others failed.

---

## Testing Strategy

### 1. Comprehensive Template Tests Created

**File**: `tests/automation/test_template_rendering.py`
**Tests**: 16 comprehensive tests

**Test Coverage**:
- ✅ All template files exist and load
- ✅ Templates render with minimal YAML data
- ✅ Each optional field handles undefined gracefully
- ✅ All optional fields use `| default()` filters
- ✅ No UndefinedError exceptions raised

**Test Classes**:
1. `TestTemplateLoading` - Verify template files exist and load
2. `TestServiceTemplate` - Test service template with missing fields
3. `TestIntegrationTemplate` - Test integration template with missing fields
4. `TestFeatureTemplate` - Test feature template with missing fields
5. `TestTemplateDefaults` - Verify all templates use defaults

**Test Results**: 16/16 passing ✅

---

### 2. Full Pipeline Validation

**Test**: Run `batch_regenerate.py --apply` on all 142 YAML files

**Before Fix**:
```
Total files: 142
Success: 47
Failed: 95 (67% failure rate)
```

**After Fix**:
```
Total files: 142
Success: 142
Failed: 0 (100% success rate ✅)
```

**Files Generated**:
- 142 Claude design docs → `docs/dev/design/`
- 142 Wiki docs → `docs/wiki/`
- Total: 284 markdown files

---

### 3. Test Suite Results

**Total Tests**: 615 tests
**Passing**: 615 ✅
**Skipped**: 2
**Failed**: 0

**New Tests Added**:
- 16 template rendering tests
- Previous: ~599 tests
- Current: 615 tests

---

## Files Modified

### Templates (3 files)
- `templates/service.md.jinja2` - Fixed 3 undefined variables
- `templates/integration.md.jinja2` - Fixed 9 undefined variables
- `templates/feature.md.jinja2` - Fixed 1 undefined variable

### Tests (1 file)
- `tests/automation/test_template_rendering.py` - 16 new tests (379 lines)

### Generated Documentation (284 files)
- `docs/dev/design/**/*.md` - 142 Claude design docs regenerated
- `docs/wiki/**/*.md` - 142 Wiki docs created

---

## Commits

### Commit 1: Template Fixes and Tests
```
fix(templates): handle undefined variables with | default()

- Fixed service, integration, and feature templates
- Added 16 comprehensive template rendering tests
- All 615 tests passing
```

### Commit 2: Regenerated Documentation
```
docs: regenerate all docs from YAML with template fixes

- ✅ 142/142 files generated successfully
- 142 Claude design docs + 142 Wiki docs
- Created complete docs/wiki/ structure
```

---

## Wiki Structure Created

```
docs/wiki/
├── .templates/           # 3 files
├── architecture/         # 5 files
├── features/
│   ├── adult/           # 5 files
│   ├── audiobook/       # 1 file
│   ├── book/            # 1 file
│   ├── comics/          # 1 file
│   ├── livetv/          # 1 file
│   ├── music/           # 1 file
│   ├── photos/          # 1 file
│   ├── playback/        # 6 files
│   ├── podcasts/        # 1 file
│   ├── shared/          # 14 files
│   └── video/           # 2 files
├── integrations/
│   ├── anime/           # 3 files
│   ├── auth/            # 4 files
│   ├── casting/         # 2 files
│   ├── infrastructure/  # 4 files
│   ├── livetv/          # 3 files
│   ├── metadata/        # 31 files (across subdirs)
│   ├── scrobbling/      # 5 files
│   ├── servarr/         # 5 files
│   ├── transcoding/     # 1 file
│   └── wiki/            # 6 files (across subdirs)
├── operations/          # 8 files
├── patterns/            # 3 files
├── research/            # 2 files
├── services/            # 15 files
└── technical/           # 10 files

Total: 142 wiki markdown files
```

---

## Prevention Measures

### 1. Comprehensive Tests
- All templates now have dedicated test coverage
- Tests verify minimal YAML data renders successfully
- Tests check each optional field handles undefined gracefully

### 2. Template Best Practices
- **ALWAYS** use `| default()` for optional YAML fields
- **NEVER** use bare `{%- if variable %}` without default
- Use appropriate default types:
  - Lists: `| default([])`
  - Dicts: `| default({})`
  - Strings: `| default('')`
  - Comments: `| default('<!-- Comment -->')`

### 3. CI/CD Integration
- Tests run on every commit
- Template rendering failures caught before merge
- 615 tests ensure no regressions

---

## Lessons Learned

1. **Jinja2 Strict Mode**: Templates should always handle undefined variables gracefully
2. **Optional Fields**: Not all YAML files include all optional fields - this is valid
3. **Test Early**: Template bugs affect ALL generated files - test comprehensively
4. **User Impact**: 67% failure rate completely broke wiki generation
5. **Default Values**: Always provide sensible defaults for optional template variables

---

## Related Issues

- **Original Bug**: #wiki-generation-bug
- **Related**: Step 0 was missing from doc-pipeline.sh (fixed separately)

---

## Verification Checklist

- [x] All 142 YAML files generate successfully
- [x] Claude design docs regenerated (142 files)
- [x] Wiki docs created (142 files)
- [x] All 615 tests passing
- [x] No linting errors
- [x] Template best practices documented
- [x] Comprehensive tests added (16 tests)
- [x] Changes committed and ready to push

---

**Status**: ✅ **RESOLVED**
**Impact**: Critical bug affecting 67% of documentation generation
**Fix Verified**: Full pipeline test, 615 tests passing, all files generated successfully
