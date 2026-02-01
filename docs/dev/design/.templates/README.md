# Design Documentation Templates

<!-- DESIGN: .templates, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


**Purpose**: Jinja2 templates for generating both Claude-optimized and Wiki documentation from a single source.

**Master Reference**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - All versions, dependencies, and architecture decisions

---

## Files

| File | Purpose |
|------|---------|
| `DESIGN_TEMPLATE.md.jinja2` | Main template with Claude/Wiki conditionals |
| `VARIABLES.yaml` | Documentation of all template variables |
| `README.md` | This file - usage guide |

---

## Template Philosophy

**Single Source, Dual Output**:
- **Claude Docs** (`.claude/docs/`): Implementation details, file paths, technical architecture
- **Wiki Docs** (`docs/wiki/`): User guides, screenshots, tutorials, getting started

**Powered by Conditionals**:
```jinja2
{% if claude %}
  Technical implementation details...
{% endif %}

{% if wiki %}
  User-friendly explanation...
{% endif %}
```

---

## Variable Sources

**CRITICAL**: Many variables must come from SOURCE_OF_TRUTH.md:

- **Go package versions**: `go_packages[].version` → SOURCE_OF_TRUTH.md
- **Dependency versions**: PostgreSQL, Dragonfly, River versions
- **API namespaces**: `/api/v1/*` structure
- **Database schemas**: Table structures, migrations
- **Caching strategy**: otter (L1), rueidis (L2) patterns
- **Testing coverage**: 80%+ requirement

**Generated Variables**:
- `last_updated`: Auto-generated timestamp
- `generation_date`: When doc was generated
- `sot_ref_checked`: Validation flag

---

## Usage

### 1. Create Variable Data File

For each design doc, create a YAML data file with all variables:

```yaml
# data/MUSIC_MODULE.yaml
feature_name: "Music Module"
category: "Features - Music"
category_path: "features/music"

status_design: "🟡"
status_design_notes: "Scaffold - needs detailed spec"
# ... all other variables
```

### 2. Generate Documentation

**Claude Version**:
```bash
python scripts/generate-docs.py \
  --template DESIGN_TEMPLATE.md.jinja2 \
  --data data/MUSIC_MODULE.yaml \
  --output .claude/docs/features/music/MUSIC_MODULE.md \
  --target claude
```

**Wiki Version**:
```bash
python scripts/generate-docs.py \
  --template DESIGN_TEMPLATE.md.jinja2 \
  --data data/MUSIC_MODULE.yaml \
  --output docs/wiki/features/music/Music-Module.md \
  --target wiki
```

### 3. Validate Against SOURCE_OF_TRUTH

```bash
python scripts/validate-sot-refs.py \
  --doc .claude/docs/features/music/MUSIC_MODULE.md \
  --sot ../../00_SOURCE_OF_TRUTH.md
```

Checks:
- Go package versions match SOT
- Database schemas match SOT
- API namespaces match SOT
- Caching patterns match SOT

---

## Template Sections

### Claude-Only Sections

These appear ONLY in Claude docs (`.claude/docs/`):

1. **Architecture**
   - Component structure with file paths
   - Design patterns used
   - Data flow diagrams

2. **Database Schema**
   - CREATE TABLE statements
   - Indexes and constraints
   - Migration files
   - ER diagrams

3. **API Endpoints**
   - Request/response examples
   - Authentication requirements
   - RBAC scopes
   - Rate limiting rules

4. **External Integrations**
   - Integration points (technical)
   - Data sync strategies
   - Error handling code patterns

5. **Testing Strategy**
   - Unit test file locations
   - Integration test scenarios
   - Mocking strategies
   - Coverage targets (80%+)

6. **Security Considerations**
   - RBAC scope definitions
   - Sensitive data handling
   - Input validation rules

7. **Performance Considerations**
   - Query optimizations
   - Monitoring metrics
   - Bottleneck analysis

8. **Implementation Checklist**
   - Phase-by-phase tasks
   - Go packages with versions
   - File locations

9. **Dependencies**
   - Design doc dependencies
   - Blocks/blocked by
   - Related docs

10. **Source Documentation**
    - Links to `docs/dev/sources/`
    - Cross-references to design docs

### Wiki-Only Sections

These appear ONLY in Wiki docs (`docs/wiki/`):

1. **Getting Started**
   - Prerequisites
   - Installation steps
   - Configuration guide

2. **Usage Examples**
   - Step-by-step tutorials
   - Example code (user-facing)
   - Common workflows

3. **Screenshots**
   - UI screenshots with descriptions
   - Visual guides
   - Before/after comparisons

4. **Troubleshooting**
   - Common issues
   - Symptoms and solutions
   - FAQ

5. **Community**
   - GitHub issues link
   - Discussions link
   - Wiki link

### Shared Sections

These appear in BOTH versions (with different content):

1. **Status Table**
   - Claude: Includes "Notes" column with technical details
   - Wiki: Simple status only

2. **Overview**
   - Claude: Technical purpose, architecture pattern
   - Wiki: User benefits, key features, use cases

---

## Example: Music Module

### Claude Doc Preview

```markdown
# Music Module

<!-- BREADCRUMB: [Design Index](../DESIGN_INDEX.md) > [Features - Music](../features/music/INDEX.md) > Music Module -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🟡 | Scaffold - needs detailed database schema |
| Code | 🔴 | Not started |## Architecture

```
┌─────────────────────────────────────────┐
│           API Layer (ogen)              │
└─────────────┬───────────────────────────┘
              │
┌─────────────▼───────────────────────────┐
│     Music Service (otter cache)         │
└─────────────┬───────────────────────────┘
              │
┌─────────────▼───────────────────────────┐
│   Repository (PostgreSQL + sqlc)        │
└─────────────────────────────────────────┘
```

## Database Schema

**Reference**: See SOURCE_OF_TRUTH.md database section.

### `artists`

```sql
CREATE TABLE artists (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  sort_name TEXT,
  musicbrainz_id UUID,
  ...
);
```

## API Endpoints

### GET `/api/v1/music/artists`

List all artists.

**Authentication**: Bearer token
**RBAC Scope**: `music:read`

**Request**:
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8096/api/v1/music/artists
```

...
```

### Wiki Doc Preview

```markdown
# Music Module

> Organize and stream your music collection

## Getting Started

### Prerequisites

- Revenge server installed
- Music library added to Revenge
- (Optional) Lidarr for automatic metadata

### Installation

Music support is built-in! Just add your music library:

1. Go to Settings → Libraries
2. Click "Add Library"
3. Select type: "Music"
4. Choose your music folder
5. Click "Save"

Revenge will automatically:
- Scan your music files
- Fetch metadata from MusicBrainz
- Organize by Artist/Album/Track
- Enable streaming

### Usage Examples

#### Play Music

1. Browse to "Music" in sidebar
2. Select an artist
3. Choose an album
4. Click a track to play

![Music Player Screenshot](screenshots/music-player.png)

#### Create Playlist

1. Browse music library
2. Click "Add to Playlist" on tracks
3. Go to "Playlists"
4. Name your playlist
5. Enjoy!

...
```

---

## Best Practices

### 1. Always Reference SOURCE_OF_TRUTH

**DO**:
```yaml
go_packages:
  - import_path: "go.uber.org/fx"
    version: "v1.23.0"  # From SOURCE_OF_TRUTH.md
```

**DON'T**:
```yaml
go_packages:
  - import_path: "go.uber.org/fx"
    version: "latest"  # ❌ No hard-coded versions!
```

### 2. Keep Claude Docs Technical

**DO** (Claude):
```markdown
**Location**: `internal/content/music/service.go`

**Caching**: L1 (otter) 5m TTL, L2 (rueidis) 1h TTL
```

**DON'T** (Claude):
```markdown
Music is stored in a special place and cached for faster loading!
```

### 3. Keep Wiki Docs User-Friendly

**DO** (Wiki):
```markdown
Revenge automatically organizes your music by artist, album, and track.
```

**DON'T** (Wiki):
```markdown
The music service implements a repository pattern with otter L1 caching.
```

### 4. Use Diagrams

**Claude**: ASCII/Mermaid architecture diagrams
**Wiki**: Screenshots, visual guides

### 5. Validate Everything

After generation:
- Run `validate-sot-refs.py` (check SOURCE_OF_TRUTH alignment)
- Run `markdownlint` (check formatting)
- Run `validate-links.py` (check all links)
- Test template rendering with both `--target claude` and `--target wiki`

---

## Automation

### Doc Pipeline Integration

The template system integrates with the existing doc pipeline:

**Stage 7: Wiki Generation**
```bash
scripts/doc-pipeline.sh
  └─ Stage 1-6: Existing stages
  └─ Stage 7: NEW - Generate wiki docs from templates
      ├─ Generate all wiki versions
      ├─ Sync to docs/wiki/
      └─ Push to GitHub Wiki
```

### CI/CD Hooks

**On Design Doc Change**:
```yaml
# .github/workflows/doc-validation.yml
on:
  push:
    paths:
      - 'docs/dev/design/**/*.md'

jobs:
  regenerate:
    - Generate Claude docs → .claude/docs/
    - Generate Wiki docs → docs/wiki/
    - Validate SOT references
    - Lint all generated docs
    - Sync to GitHub Wiki
```

---

## Template Maintenance

### Updating the Template

1. Edit `DESIGN_TEMPLATE.md.jinja2`
2. Update `VARIABLES.yaml` documentation
3. Test with existing data files
4. Regenerate ALL docs (to apply template changes)
5. Validate with linting

### Adding New Variables

1. Add to `VARIABLES.yaml` with documentation
2. Add to template with conditional if needed
3. Update this README with usage examples
4. Update generation scripts if needed

### Version Control

**Template files are version controlled**:
- `DESIGN_TEMPLATE.md.jinja2` → Git
- `VARIABLES.yaml` → Git
- `README.md` → Git

**Data files are version controlled**:
- `data/*.yaml` → Git (design-specific data)

**Generated files**:
- `.claude/docs/*.md` → Git (tracked)
- `docs/wiki/*.md` → Git (tracked, synced to GitHub Wiki)

---

## Troubleshooting

### Template Rendering Error

**Problem**: `jinja2.exceptions.UndefinedError: 'feature_name' is undefined`

**Solution**: Ensure data YAML has all required variables from `VARIABLES.yaml`

### SOT Validation Fails

**Problem**: `Package version mismatch: fx v1.22.0 (doc) vs v1.23.0 (SOT)`

**Solution**: Update data YAML to match SOURCE_OF_TRUTH.md versions

### Links Break After Generation

**Problem**: Relative links don't resolve

**Solution**: Check `category_path` is correct in data YAML, run `validate-links.py`

---

## Future Enhancements

**Skill Ideas** (Claude Code):
1. **scaffold-doc** - Interactive doc creation from template
2. **validate-doc** - Comprehensive doc validation (SOT, lint, links)
3. **template-convert** - Convert existing doc to template format
4. **doc-coverage** - Show design coverage gaps

**Automation**:
- Auto-generate data YAML from existing docs
- Auto-sync variables from SOURCE_OF_TRUTH
- Auto-update all docs when SOT changes
- Auto-screenshot capture for wiki docs

---

**Last Updated**: 2026-01-31
**Template Version**: 1.0.0
