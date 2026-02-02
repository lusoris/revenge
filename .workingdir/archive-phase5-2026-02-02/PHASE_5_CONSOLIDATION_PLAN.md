# Phase 5 Consolidation Plan: Detailed Execution Strategy

## Overview

This document provides line-by-line instructions for consolidating duplicate content across ~900 lines in 15+ files. Each step includes exact file locations, text to remove, and replacement cross-references.

---

## Step 5.3: Cache Architecture Consolidation (~60 lines, 4 files)

### Canonical Location
**File**: `data/integrations/caching/DRAGONFLY.yaml`
**Keep**: Full two-tier caching architecture description (L1 otter + L2 Dragonfly)

### Files to Update

#### 1. data/architecture/01_ARCHITECTURE.yaml
**Search for**: Section describing "Two-tier caching: L1 (otter) + L2 (Dragonfly)"
**Action**: Replace detailed description with:
```yaml
  Cache layer uses two-tier architecture (L1 otter + L2 Dragonfly).
  See data/integrations/caching/DRAGONFLY.yaml for complete architecture details.
```
**Estimated removal**: ~15 lines

#### 2. data/services/METADATA.yaml
**Search for**: `cache_architecture:` or `dependencies:` section mentioning otter + Dragonfly
**Action**: Replace with:
```yaml
cache_architecture: |
  Uses two-tier caching for metadata responses.
  See data/integrations/caching/DRAGONFLY.yaml for implementation details.
```
**Estimated removal**: ~20 lines

#### 3. data/services/SESSION.yaml
**Search for**: Architecture diagram or text mentioning both otter and Dragonfly
**Action**: Update diagram/text to reference DRAGONFLY.yaml instead of inline description
**Estimated removal**: ~15-20 lines

#### 4. data/architecture/03_METADATA_SYSTEM.yaml
**Search for**: Cache layer description in system components
**Action**: Replace with cross-reference to DRAGONFLY.yaml
**Estimated removal**: ~10 lines

### Verification
- [ ] Check DRAGONFLY.yaml has complete cache architecture description
- [ ] Ensure all 4 files reference DRAGONFLY.yaml in design_refs
- [ ] No duplicate text remains
- [ ] Run yamllint on all modified files
- [ ] Commit with message: "refactor: consolidate cache architecture to DRAGONFLY.yaml"

---

## Step 5.4: Arr Dual-Role Consolidation (~300 lines, 7 files)

### Canonical Location
**File**: `data/architecture/03_METADATA_SYSTEM.yaml`
**Section**: `arr_dual_role_architecture:` (or similar)
**Keep**: Full description of "Arr services serve DUAL purposes: 1) PRIMARY Metadata Aggregator, 2) Download Automation Manager"

### Files to Update

#### 1. data/architecture/01_ARCHITECTURE.yaml
**Search for**: Lines 275-288 (per analysis), "Arr services serve DUAL purposes"
**Current content**: Full explanation of dual role
**Replace with**:
```yaml
  Arr services (Radarr, Sonarr, etc.) serve dual purposes as both metadata aggregators
  and download automation managers.
  See data/architecture/03_METADATA_SYSTEM.yaml#arr_dual_role_architecture for details.
```
**Estimated removal**: ~50 lines

#### 2. data/services/METADATA.yaml
**Search for**: `dual_role_architecture:` section
**Action**: Replace with cross-reference
**Estimated removal**: ~60 lines

#### 3. data/integrations/servarr/RADARR.yaml
**Search for**: `dual_role_metadata_and_downloads:` section
**Current**: Full explanation of dual role
**Replace with**:
```yaml
dual_role_metadata_and_downloads: |
  Radarr serves dual purposes in the Revenge architecture.
  See data/architecture/03_METADATA_SYSTEM.yaml#arr_dual_role_architecture for complete
  explanation of how Arr services function as both metadata aggregators and download managers.
```
**Estimated removal**: ~50 lines

#### 4. data/integrations/metadata/video/TMDB.yaml
**Search for**: `supplementary_role:` section mentioning Arr dual role
**Action**: Keep supplementary role description, remove Arr dual role explanation, add reference
**Estimated removal**: ~30 lines

#### 5. data/features/video/MOVIE_MODULE.yaml
**Search for**: wiki_overview mentioning Arr dual role
**Action**: Simplify to brief mention with reference
**Estimated removal**: ~40 lines

#### 6. data/features/video/TVSHOW_MODULE.yaml
**Search for**: wiki_overview mentioning Arr dual role
**Action**: Simplify to brief mention with reference
**Estimated removal**: ~40 lines

#### 7. data/integrations/servarr/SONARR.yaml
**Search for**: Similar dual_role section as RADARR.yaml
**Action**: Replace with cross-reference
**Estimated removal**: ~30 lines

### Verification
- [ ] 03_METADATA_SYSTEM.yaml has complete arr_dual_role_architecture section
- [ ] All 7 files have design_refs pointing to 03_METADATA_SYSTEM.yaml
- [ ] Brief mentions remain where contextually relevant
- [ ] No full duplicate descriptions remain
- [ ] Run yamllint on all modified files
- [ ] Commit with message: "refactor: consolidate Arr dual-role architecture to METADATA_SYSTEM.yaml"

---

## Step 5.5: Metadata Priority Chain Consolidation (~125 lines, 8 files)

### Canonical Location
**File**: `data/architecture/03_METADATA_SYSTEM.yaml`
**Section**: Metadata priority chain description (L1 Cache → L2 Cache → Arr → External APIs)
**Keep**: Complete flow diagram and explanation of PRIMARY vs SUPPLEMENTARY

### Files to Update

#### 1. docs/dev/design/00_SOURCE_OF_TRUTH.md
**Search for**: Lines 68-80 (per analysis), metadata priority chain
**Current**: Text description of priority chain
**Action**: Keep brief summary, add reference to 03_METADATA_SYSTEM.yaml
**Replace with**:
```markdown
### Metadata Priority Chain

L1 Cache (otter) → L2 Cache (Dragonfly) → Arr Services → External APIs

See `data/architecture/03_METADATA_SYSTEM.yaml` for complete architecture details,
including PRIMARY vs SUPPLEMENTARY provider terminology.
```
**Estimated removal**: ~10 lines (keep concise summary)

#### 2. data/architecture/01_ARCHITECTURE.yaml
**Search for**: system_components section describing priority chain
**Action**: Replace detailed flow with cross-reference
**Estimated removal**: ~20 lines

#### 3. data/services/METADATA.yaml
**Search for**: `provider_priority_chain:` section
**Action**: Replace with:
```yaml
provider_priority_chain: |
  Metadata requests follow strict priority chain.
  See data/architecture/03_METADATA_SYSTEM.yaml for complete chain definition
  and PRIMARY vs SUPPLEMENTARY provider roles.
```
**Estimated removal**: ~25 lines

#### 4. data/features/video/MOVIE_MODULE.yaml
**Search for**: `config_keys:` or similar section describing priority chain
**Action**: Replace chain description with reference
**Estimated removal**: ~15 lines

#### 5. data/features/video/TVSHOW_MODULE.yaml
**Search for**: Similar priority chain description as MOVIE_MODULE
**Action**: Replace with reference
**Estimated removal**: ~15 lines

#### 6. data/integrations/metadata/video/TMDB.yaml
**Search for**: `supplementary_role:` section explaining priority chain context
**Action**: Keep "SUPPLEMENTARY" definition, remove priority chain flow, add reference
**Estimated removal**: ~15 lines

#### 7. data/integrations/servarr/RADARR.yaml
**Search for**: Priority chain context in metadata aggregation section
**Action**: Simplify to reference
**Estimated removal**: ~15 lines

#### 8. data/integrations/servarr/SONARR.yaml
**Search for**: Similar to RADARR.yaml
**Action**: Simplify to reference
**Estimated removal**: ~10 lines

### Verification
- [ ] 03_METADATA_SYSTEM.yaml has complete priority chain with diagrams
- [ ] SOT maintains high-level summary
- [ ] All feature modules reference 03_METADATA_SYSTEM.yaml
- [ ] All integration files reference canonical source
- [ ] design_refs updated in all 8 files
- [ ] Run yamllint on all modified files
- [ ] Commit with message: "refactor: consolidate metadata priority chain to METADATA_SYSTEM.yaml"

---

## Step 5.6: Update Proxy/VPN References (Bonus - leverages Step 5.2)

Now that HTTP_CLIENT.yaml exists, update files to reference it:

### Files to Update

#### 1. data/architecture/03_METADATA_SYSTEM.yaml
**Search for**: Lines 152-250, entire `proxy_vpn_support:` block
**Replace with**:
```yaml
proxy_vpn_support: |
  External metadata API calls support optional proxy/VPN routing for privacy and geo-access.
  See data/patterns/HTTP_CLIENT.yaml for complete implementation guide including:
  - HTTP/HTTPS proxy configuration
  - SOCKS5 and Tor support
  - VPN interface binding
  - Security considerations
```
**Estimated removal**: ~90 lines

#### 2. data/services/METADATA.yaml
**Search for**: `proxy_vpn_support:` section (line ~581)
**Replace with**: Similar concise reference to HTTP_CLIENT.yaml
**Estimated removal**: ~50 lines

#### 3. data/integrations/metadata/video/TMDB.yaml
**Search for**: Lines 437-480, `proxy_vpn_support:` block
**Replace with**:
```yaml
proxy_vpn_support: |
  TMDb API calls can be routed through proxy/VPN when needed.
  See data/patterns/HTTP_CLIENT.yaml for complete configuration guide.
```
**Estimated removal**: ~40 lines

#### 4. data/integrations/metadata/video/THETVDB.yaml
**Search for**: Similar `proxy_vpn_support:` block
**Replace with**: Reference to HTTP_CLIENT.yaml
**Estimated removal**: ~40 lines

#### 5-9. Other integration files (FREEONES, PORNHUB, THEPORNDB, THENUDE, etc.)
**Action**: Replace `proxy_vpn_support:` blocks with references
**Estimated removal**: ~30-40 lines each

### Verification
- [ ] All files reference data/patterns/HTTP_CLIENT.yaml in proxy sections
- [ ] All files have HTTP_CLIENT.yaml in design_refs
- [ ] No duplicate proxy implementation code remains
- [ ] Run yamllint on all modified files
- [ ] Commit with message: "refactor: replace proxy/VPN duplicates with HTTP_CLIENT.yaml references"

---

## Execution Order

1. **Step 5.3**: Cache Architecture (smallest, 4 files, ~60 lines)
2. **Step 5.6**: Proxy/VPN References (leverages existing HTTP_CLIENT.yaml)
3. **Step 5.4**: Arr Dual-Role (7 files, ~300 lines)
4. **Step 5.5**: Metadata Priority Chain (8 files, ~125 lines)

**Rationale**: Start with smallest consolidation to establish pattern, then leverage newly created HTTP_CLIENT.yaml, then tackle larger consolidations.

---

## Safety Checklist (Per Step)

Before each commit:
- [ ] Read original content to understand context
- [ ] Verify canonical file has complete information
- [ ] Ensure replacement text maintains necessary context
- [ ] Check no unique information is lost
- [ ] Update design_refs in all affected files
- [ ] Run `yamllint` on all modified files
- [ ] Run `python -c "import yaml; yaml.safe_load(open('file.yaml'))"`  on each
- [ ] Review git diff carefully
- [ ] Commit with descriptive message

After all steps:
- [ ] Run sync-versions.py --strict (should still pass)
- [ ] Run 04-sync-sot-status.py --strict (should still pass)
- [ ] Run all unit tests
- [ ] Create summary document of Phase 5 changes

---

## Risk Assessment

**Low Risk**:
- Cache Architecture (Step 5.3) - Technical implementation details, low coupling
- Proxy/VPN (Step 5.6) - Already created canonical source, straightforward replacement

**Medium Risk**:
- Arr Dual-Role (Step 5.4) - Conceptual explanation, spread across many files
- Metadata Priority Chain (Step 5.5) - Core architecture, referenced in many contexts

**Mitigation**:
- Commit after each step for easy rollback
- Keep brief contextual mentions where necessary
- Don't over-consolidate - references should make sense in context
- Test YAML parsing after each file modification

---

## Success Criteria

- **Lines Removed**: ~900+ lines of duplicate content
- **Commits**: 4 clean commits (Steps 5.3, 5.4, 5.5, 5.6)
- **Files Modified**: 15-20 YAML files
- **Validation**: All tests passing, no YAML errors
- **Documentation**: Clear cross-references established
- **Maintainability**: Single source of truth for each architectural concept

---

**Status**: Plan ready for execution
**Next Action**: Begin Step 5.3 (Cache Architecture)
