# Content Inconsistency & Duplication Analysis

**Generated**: 2026-02-02
**Source**: Analysis of SOURCE_OF_TRUTH.md vs all YAML data files

---

## Executive Summary

After analyzing the SOURCE_OF_TRUTH (SOT) against 100+ YAML data files, the following categories of issues were identified:

| Category | Count | Severity |
|----------|-------|----------|
| Status Mismatches | 8 | Medium |
| Version Inconsistencies | 3 | Low |
| Duplicate Content | 15+ | High |
| Missing Cross-References | 5 | Low |
| Terminology Inconsistencies | 4 | Medium |
| Structural Inconsistencies | 6 | Low |

---

## 1. Status Mismatches (SOT vs YAML Files)

### 1.1 Backend Services Status

**SOT declares these services as "Complete":**

| Service | SOT Status | YAML overall_status | Mismatch? |
|---------|------------|---------------------|-----------|
| Auth | ✅ Complete | ✅ Complete | No |
| User | ✅ Complete | ✅ Complete | No |
| Session | ✅ Complete | ✅ Complete | No |
| RBAC | ✅ Complete | ✅ Complete | No |
| Activity | ✅ Complete | ✅ Complete | No |
| Settings | ✅ Complete | ✅ Complete | No |
| API Keys | ✅ Complete | ✅ Complete | No |
| OIDC | ✅ Complete | ✅ Complete | No |
| Library | ✅ Complete | ✅ Complete | No |

**SOT declares these as "Partial" but YAML says "Complete":**

| Service | SOT Status | YAML overall_status | Issue |
|---------|------------|---------------------|-------|
| Metadata | 🟡 Partial | ✅ Complete | **MISMATCH** |
| Search | 🟡 Partial | ✅ Complete | **MISMATCH** |

**Recommendation**: Update SOT to reflect actual YAML status, or update YAML if SOT is authoritative.

### 1.2 Content Modules Status

| Module | SOT Status | YAML overall_status | Mismatch? |
|--------|------------|---------------------|-----------|
| Movie | ✅ Complete | ✅ Complete | No |
| TV Show | ✅ Complete | ✅ Complete | No |
| Music | 🟡 Scaffold | Not checked | - |
| QAR Voyages | 🟡 Scaffold | ✅ Complete | **MISMATCH** |

### 1.3 Integration Status

**SOT Metadata Providers vs YAML:**

| Provider | SOT Status | YAML Status | Issue |
|----------|------------|-------------|-------|
| TMDb | ✅ | ✅ Complete | OK |
| TheTVDB | 🔴 | ✅ Complete | **MISMATCH** |
| MusicBrainz | 🔴 | Not checked | - |
| StashDB | 🟡 | ✅ Complete | **MISMATCH** |

**SOT Arr Ecosystem vs YAML:**

| Service | SOT Status | YAML Status | Issue |
|---------|------------|-------------|-------|
| Radarr | ✅ | ✅ Complete | OK |
| Sonarr | 🔴 | ✅ Complete | **MISMATCH** |
| Whisparr | 🟡 | ✅ Complete | **MISMATCH** |

---

## 2. Version Inconsistencies

### 2.1 Package Version References

The SOT defines specific versions that should be referenced (not duplicated) in YAML files:

| Package | SOT Version | Issue |
|---------|-------------|-------|
| pgx | v5.8.0 | Some YAML files reference "pgx/v5" without version |
| River | v0.26.0 | Consistent |
| rueidis | v1.0.49 | Consistent |
| otter | v1.2.4 | Consistent |

**Note**: YAML files correctly don't duplicate versions, they just reference package names.

### 2.2 PostgreSQL Version

- **SOT**: PostgreSQL 18.1
- **POSTGRESQL.yaml**: "PostgreSQL 18.0+" (correct - minimum version)
- **No inconsistency** - YAML correctly specifies minimum, SOT specifies exact

### 2.3 Dragonfly Version

- **SOT**: v1.36.0
- **DRAGONFLY.yaml**: No version specified
- **Recommendation**: Add version reference to DRAGONFLY.yaml

---

## 3. Duplicate Content (High Priority)

### 3.1 Metadata Priority Chain (Duplicated 5+ times)

The metadata priority chain is described in detail in:
1. `docs/dev/design/00_SOURCE_OF_TRUTH.md` (lines 68-80)
2. `data/architecture/01_ARCHITECTURE.yaml` (system_components section)
3. `data/architecture/03_METADATA_SYSTEM.yaml`
4. `data/services/METADATA.yaml` (provider_priority_chain section)
5. `data/features/video/MOVIE_MODULE.yaml` (config_keys section)
6. `data/features/video/TVSHOW_MODULE.yaml` (config_keys section)
7. `data/integrations/metadata/video/TMDB.yaml` (supplementary_role section)
8. `data/integrations/servarr/RADARR.yaml` (dual_role_metadata_and_downloads section)

**Each file contains nearly identical descriptions of:**
- L1 Cache (otter) → L2 Cache (Dragonfly) → Arr Services → External APIs
- PRIMARY vs SUPPLEMENTARY terminology
- Proxy/VPN optional routing

**Recommendation**:
- Keep full description ONLY in `03_METADATA_SYSTEM.yaml`
- Other files should reference: "See 03_METADATA_SYSTEM for priority chain details"

### 3.2 Arr Dual-Role Description (Duplicated 6+ times)

The concept "Arr services serve DUAL purposes: metadata aggregation + download automation" appears in:
1. `01_ARCHITECTURE.yaml` (system_components, lines ~275-288)
2. `03_METADATA_SYSTEM.yaml`
3. `METADATA.yaml` (dual_role_architecture section)
4. `RADARR.yaml` (dual_role_metadata_and_downloads section)
5. `TMDB.yaml` (supplementary_role section)
6. `MOVIE_MODULE.yaml` (wiki_overview)
7. `TVSHOW_MODULE.yaml` (wiki_overview)

**Nearly identical text in each:**
```
Arr services serve DUAL purposes:
1. PRIMARY Metadata Aggregator
2. Download Automation Manager
```

**Recommendation**: Define once in `03_METADATA_SYSTEM.yaml`, reference elsewhere.

### 3.3 Database Schema Definitions (Duplicated)

User-related tables defined in multiple places:
1. `AUTH.yaml` - defines `users` table
2. `USER.yaml` - references `users` table, defines `user_profiles`
3. `SESSION.yaml` - defines `sessions` table with FK to `users`

**Issue**: The `users` table schema is partially defined in AUTH.yaml but USER.yaml also references it.

**Recommendation**:
- Define all user-related tables in a single `USERS_SCHEMA.yaml` or keep in AUTH.yaml
- Other files should reference, not redefine

### 3.4 Cache Architecture (Duplicated 4+ times)

Two-tier caching (L1 otter + L2 Dragonfly) described in:
1. `01_ARCHITECTURE.yaml` (system_components section)
2. `DRAGONFLY.yaml` (cache_architecture section)
3. `METADATA.yaml` (dependencies section)
4. `SESSION.yaml` (architecture diagram shows both)

**Recommendation**: Full description in DRAGONFLY.yaml, brief reference elsewhere.

### 3.5 Proxy/VPN Support (Duplicated 5+ times)

Optional proxy/VPN routing for external APIs described in:
1. `01_ARCHITECTURE.yaml` (system_components, lines ~291-299)
2. `METADATA.yaml` (proxy_vpn_support section)
3. `TMDB.yaml` (proxy_vpn_support section)
4. `MOVIE_MODULE.yaml` (dependencies section)
5. `TVSHOW_MODULE.yaml` (dependencies section)

**Each contains similar text about:**
- HTTP/SOCKS5/Tor/VPN options
- "OPTIONAL - must be explicitly enabled"
- Configuration examples

**Recommendation**: Define fully in `HTTP_CLIENT.yaml`, brief reference elsewhere.

---

## 4. Missing Cross-References

### 4.1 SOT References in YAML Files

Many YAML files should reference SOT for versions but don't:

| File | Missing Reference |
|------|-------------------|
| `DRAGONFLY.yaml` | Should reference SOT for rueidis version |
| `POSTGRESQL.yaml` | Should reference SOT for pgx version |
| `MOVIE_MODULE.yaml` | Has inline version refs, should use SOT |

### 4.2 Missing design_refs

Some YAML files have incomplete `design_refs`:

| File | Missing Reference |
|------|-------------------|
| `TMDB.yaml` | Missing reference to HTTP_CLIENT.yaml |
| `RADARR.yaml` | Missing reference to MOVIE_MODULE.yaml |
| `SESSION.yaml` | Missing reference to AUTH.yaml |

---

## 5. Terminology Inconsistencies

### 5.1 "Supplementary" vs "Fallback"

| File | Term Used |
|------|-----------|
| SOT | "External APIs (SUPPLEMENTARY via optional proxy/VPN)" |
| TMDB.yaml | "SUPPLEMENTARY metadata provider (fallback + enrichment)" |
| MOVIE_MODULE.yaml | "TMDb (fallback)" in diagram, "Supplementary" in text |
| TVSHOW_MODULE.yaml | "TheTVDB (fallback)" in diagram |

**Issue**: Diagrams use "fallback" while text uses "SUPPLEMENTARY"

**Recommendation**: Standardize on "SUPPLEMENTARY" everywhere, update diagrams.

### 5.2 Environment Variable Prefixes

| File | Prefix Used | Example |
|------|-------------|---------|
| SOT | `REVENGE_` | `REVENGE_SERVER_PORT` |
| AUTH.yaml | `AUTH_` | `AUTH_PASSWORD_MIN_LENGTH` |
| SESSION.yaml | `SESSION_` | `SESSION_TOKEN_LENGTH` |
| MOVIE_MODULE.yaml | `REVENGE_` | `REVENGE_MOVIE_CACHE_TTL` |

**Issue**: Inconsistent env var prefixes

**Recommendation**: All should use `REVENGE_` prefix as per SOT.

### 5.3 Cache Client Names

| File | Client Name |
|------|-------------|
| SOT | "rueidis" for Dragonfly |
| 01_ARCHITECTURE.yaml | "valkey-go" mentioned in dependencies |
| DRAGONFLY.yaml | "rueidis" |

**Issue**: SOT and DRAGONFLY.yaml say "rueidis", but 01_ARCHITECTURE.yaml mentions "valkey-go"

**Recommendation**: Standardize on "rueidis" as per SOT.

---

## 6. Structural Inconsistencies

### 6.1 Status Field Formats

**Inconsistent status field values:**

| File | overall_status Format |
|------|----------------------|
| AUTH.yaml | `✅ Complete` |
| MOVIE_MODULE.yaml | `✅ Complete` |
| RADARR.yaml | `✅ Complete` |

| File | status_code Format |
|------|-------------------|
| AUTH.yaml | `🔴` (emoji only) |
| MOVIE_MODULE.yaml | `🔴 Not Started` (emoji + text) |

**Recommendation**: Standardize - either emoji only or emoji + text for all files.

### 6.2 Source URL Formats

Some sources use different URL patterns:

| Type | Example |
|------|---------|
| pkg.go.dev | `https://pkg.go.dev/github.com/jackc/pgx/v5` |
| GitHub | `https://github.com/typesense/typesense-go` |
| Official docs | `https://docs.sqlc.dev/en/stable/` |

**This is acceptable** - different source types have different URLs.

### 6.3 Design Refs Path Formats

| File | Path Format |
|------|-------------|
| AUTH.yaml | `../architecture/01_ARCHITECTURE.md` |
| TMDB.yaml | `../../../architecture/03_METADATA_SYSTEM.md` |

**Issue**: Relative paths vary based on file location - this is correct but makes maintenance harder.

**Recommendation**: Consider using absolute paths from docs root.

---

## 7. Content Gaps

### 7.1 YAML Files Missing Content

| File | Missing Section |
|------|-----------------|
| `03_DESIGN_DOCS_STATUS.yaml` | No status tracking tables (just sources) |
| Several integration files | Missing `rate_limits` section |

### 7.2 SOT Sections Not Reflected in YAML

| SOT Section | Missing YAML Coverage |
|-------------|----------------------|
| GitHub Actions | No corresponding YAML file |
| Helm Repositories | No corresponding YAML file |
| Container Orchestration | Partial coverage in operations/ |

---

## 8. Recommendations Summary

### High Priority (Duplications)

1. **Consolidate Metadata Priority Chain** - Keep in 03_METADATA_SYSTEM.yaml only
2. **Consolidate Arr Dual-Role Description** - Keep in 03_METADATA_SYSTEM.yaml only
3. **Consolidate Proxy/VPN Documentation** - Keep in HTTP_CLIENT.yaml only
4. **Consolidate Cache Architecture** - Keep in DRAGONFLY.yaml only

### Medium Priority (Status/Terminology)

5. **Sync Status Values** - Align SOT with YAML files (or vice versa)
6. **Standardize "SUPPLEMENTARY" terminology** - Update diagrams
7. **Standardize env var prefixes** - All should use `REVENGE_`
8. **Fix rueidis vs valkey-go inconsistency** - Use rueidis

### Low Priority (Structural)

9. **Standardize status field formats** - Emoji + text everywhere
10. **Add missing design_refs** - Cross-reference related files
11. **Add version references to SOT** - Don't duplicate versions

---

## 9. Files Requiring Updates

### Immediate Updates Needed

| File | Action Required |
|------|-----------------|
| `01_ARCHITECTURE.yaml` | Remove duplicate priority chain, fix valkey-go→rueidis |
| `METADATA.yaml` | Reduce duplication, reference 03_METADATA_SYSTEM |
| `TMDB.yaml` | Reference HTTP_CLIENT for proxy details |
| `RADARR.yaml` | Reference 03_METADATA_SYSTEM for dual-role |
| `MOVIE_MODULE.yaml` | Reference 03_METADATA_SYSTEM, fix status format |
| `TVSHOW_MODULE.yaml` | Reference 03_METADATA_SYSTEM, fix diagram label |
| `AUTH.yaml` | Change env vars to REVENGE_ prefix |
| `SESSION.yaml` | Change env vars to REVENGE_ prefix |

### SOT Updates Needed

| Section | Action Required |
|---------|-----------------|
| Backend Services | Update Metadata status to ✅ or note discrepancy |
| Backend Services | Update Search status to ✅ or note discrepancy |
| Metadata Providers | Update TheTVDB status to ✅ |
| Arr Ecosystem | Update Sonarr status to ✅ |

---

## 10. Metrics

| Metric | Value |
|--------|-------|
| Total YAML files analyzed | ~100 |
| Files with duplications | ~15 |
| Status mismatches found | 8 |
| Total duplicated sections | ~25 |
| Estimated lines of duplicate content | ~500+ |

---

*Analysis complete. Recommend addressing High Priority items first to reduce maintenance burden and improve consistency.*
