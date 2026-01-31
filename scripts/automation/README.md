# Documentation Automation System

**Status**: ✅ Phase 1-3 Complete
**Created**: 2026-01-31
**Author**: Automation System

---

## Overview

Comprehensive documentation automation system for Revenge media server project. Implements a **Template → YAML → Dual Output** architecture for maintaining consistent, high-quality documentation.

## Architecture

```
SOURCE_OF_TRUTH.md
       ↓
   [SOT Parser] → shared-sot.yaml
       ↓
Design Docs (*.md)
       ↓
   [MD Parser] → YAML Data Files (data/**/*.yaml)
       ↓                                    ↓
   [Doc Generator] ←────────────────────────┘
       ↓
   ┌───┴────┐
   ↓        ↓
Claude    Wiki
(Tech)    (User)
```

## Components

### 1. SOT Parser ([sot_parser.py](sot_parser.py))
Extracts structured data from SOURCE_OF_TRUTH.md including:
- Content modules (12 modules)
- Backend services (18 services)
- Infrastructure components (4 components)
- Go dependencies (25 packages across 5 categories)

**Output**: `data/shared-sot.yaml`

### 2. Markdown Parser ([md_parser.py](md_parser.py))
Extracts metadata from existing design documentation:
- Document title and category
- 7-dimension status table
- Source references (auto-resolves URLs from SOURCES.yaml)
- Design doc cross-references
- Content sections

**Features**:
- Intelligent source ID resolution (exact + partial match)
- Design doc path resolution
- Category detection (feature/service/integration/operations/etc.)
- Clean YAML output with TODOs for manual completion

### 3. Batch Migrator ([batch_migrate.py](batch_migrate.py))
Bulk migration of design docs to YAML format:
- Processes 142 design docs across 28 categories
- Preserves directory structure
- Safety features:
  - Dry-run by default (`--live` to execute)
  - Skips existing files (`--force` to overwrite)
  - Detailed migration report

**Usage**:
```bash
# Dry run (default)
python batch_migrate.py

# Live migration
python batch_migrate.py --live

# Force overwrite existing files
python batch_migrate.py --live --force
```

### 4. Doc Generator ([doc_generator.py](doc_generator.py))
Generates dual documentation output from YAML:
- **Claude version**: Technical docs with implementation details
- **Wiki version**: User-friendly docs with screenshots/FAQs

**Features**:
- Jinja2 template engine with block inheritance
- Data merging (shared-sot.yaml + doc-specific YAML)
- Atomic file writes
- Valid YAML frontmatter

**Templates**:
- `base.md.jinja2` - Foundation template
- `feature.md.jinja2` - Content modules (movies, music, etc.)
- `service.md.jinja2` - Backend services (auth, user, etc.)
- `integration.md.jinja2` - External integrations (TMDb, Radarr, etc.)
- `generic.md.jinja2` - Operations/architecture/technical docs
- `wiki/base.md.jinja2` - User-friendly wiki docs

### 5. Validator ([validator.py](validator.py))
Validates YAML files against JSON schemas:
- Schema-aware validation (feature/service/integration)
- Detailed error reporting
- File path tracking

**Schemas**:
- `feature.schema.json` - Content modules
- `service.schema.json` - Backend services
- `integration.schema.json` - External integrations

**Usage**:
```bash
python validator.py
# Exit code 0 if all valid, 1 if any failures
```

### 6. YAML Analyzer ([yaml_analyzer.py](yaml_analyzer.py))
Analyzes YAML completion status:
- Identifies placeholder fields
- Checks missing required fields
- Calculates completion scores
- Prioritizes files needing completion

**Usage**:
```bash
python yaml_analyzer.py
```

### 7. YAML Completion Assistant ([yaml_completion_assistant.py](yaml_completion_assistant.py))
Auto-completes basic fields from original markdown:
- Extracts technical_summary from first paragraphs
- Generates wiki_tagline (shortened summary)
- Auto-detects schema_name (qar vs public)
- Generates module_name from doc_title

**Usage**:
```bash
# Dry-run (preview)
python yaml_completion_assistant.py --feature

# Auto-apply
python yaml_completion_assistant.py --feature --auto
python yaml_completion_assistant.py --service --auto
python yaml_completion_assistant.py --integration --auto
```

### 8. Enhanced Completion Assistant ([enhanced_completion_assistant.py](enhanced_completion_assistant.py))
Extracts category-specific fields from markdown content:

**For Integrations**:
- integration_name (from doc title)
- integration_id (slug version)
- external_service (usually same as integration_name)
- api_base_url (extracted from content)
- auth_method (detected from keywords: oauth, api_key, bearer, basic, none)

**For Services**:
- service_name (from doc title)
- package_path (extracted or generated: internal/service/{name})
- fx_module (generated from package_path)

**For Features**:
- content_types (inferred from title or extracted from content)

**Usage**:
```bash
# Dry-run (preview)
python enhanced_completion_assistant.py --integration --dry-run

# Auto-apply
python enhanced_completion_assistant.py --integration --auto
python enhanced_completion_assistant.py --service --auto
python enhanced_completion_assistant.py --feature --auto
```

**Success Rate**: Achieved 90% completion across all 142 files!

---

## Current Status

### Migration Complete ✅
- **142 design docs** migrated to YAML
- **28 categories** processed
- **100% migration success rate**

### Auto-Completion Complete ✅
- **Basic fields** auto-completed for all 108 core files (features + services + integrations)
- **Category-specific fields** extracted for all 108 core files via enhanced completion assistant
- **90% overall completion** achieved through automation

### Completion Status (LATEST) 🎉
```
76-100% complete:  140 files (99%)  ⬆️ +338% from initial 32 files!
51-75% complete:    1 file   (1%)
26-50% complete:    1 file   (1%)
0-25% complete:     0 files  (0%)
```

### Issues to Address
- **210 placeholder fields** remaining (down 51% from 426)
- **107 missing required fields** remaining (down 79% from 510!) 🎉
- **2 files** missing cross-references
- Some validation format issues (overall_status format, fx_module pattern)

### Priority Files
**Top 20 features** (65% complete each):
1. MOVIE_MODULE - Missing: feature_name, module_name, schema_name, content_types, metadata_providers
2. ADULT_CONTENT_SYSTEM - Same missing fields
3. AUDIOBOOK_MODULE - Same missing fields
4. BOOK_MODULE - Same missing fields
5. COMICS_MODULE - Same missing fields
6. MUSIC_MODULE - Same missing fields
... (15 more)

All feature files need manual completion of:
- `feature_name` - Display name (e.g., "Movie Module")
- `module_name` - Code module name (e.g., "movie")
- `schema_name` - Database schema ("public" or "qar")
- `content_types` - Content types list (e.g., ["Movies", "Collections"])
- `metadata_providers` - Provider configurations

---

## Workflow

### Full Documentation Workflow

1. **Update SOURCE_OF_TRUTH.md** with latest tech stack/design decisions
2. **Run SOT parser** to generate shared-sot.yaml:
   ```bash
   python sot_parser.py
   ```

3. **Create/update YAML data files** in `data/`:
   ```yaml
   # Example: data/features/video/MOVIE_MODULE.yaml
   doc_title: Movie Module
   doc_category: feature
   feature_name: Movie Module
   module_name: movie
   schema_name: public
   content_types:
     - Movies
     - Collections
   metadata_providers:
     - name: TMDb
       purpose: Primary metadata source
       priority: 1
   # ... more fields
   ```

4. **Validate YAML**:
   ```bash
   python validator.py
   ```

5. **Generate documentation**:
   ```bash
   python doc_generator.py
   # Or programmatically for specific files
   ```

6. **Review and commit** generated docs

### New Document Creation

1. Create YAML file in appropriate `data/` subdirectory
2. Follow schema for category (feature/service/integration)
3. Run validator to check compliance
4. Generate docs
5. Review output

---

## File Structure

```
revenge/
├── data/                          # YAML data files
│   ├── shared-sot.yaml           # From SOURCE_OF_TRUTH.md
│   ├── features/
│   │   ├── video/
│   │   │   ├── MOVIE_MODULE.yaml
│   │   │   └── TVSHOW_MODULE.yaml
│   │   ├── music/
│   │   └── adult/
│   ├── services/
│   ├── integrations/
│   ├── operations/
│   └── architecture/
│
├── templates/                     # Jinja2 templates
│   ├── base.md.jinja2
│   ├── feature.md.jinja2
│   ├── service.md.jinja2
│   ├── integration.md.jinja2
│   ├── generic.md.jinja2
│   └── wiki/
│       └── base.md.jinja2
│
├── schemas/                       # JSON validation schemas
│   ├── feature.schema.json
│   ├── service.schema.json
│   └── integration.schema.json
│
├── scripts/automation/           # Automation scripts
│   ├── sot_parser.py
│   ├── md_parser.py
│   ├── batch_migrate.py
│   ├── doc_generator.py
│   ├── validator.py
│   ├── yaml_analyzer.py
│   └── README.md (this file)
│
└── docs/
    ├── dev/design/               # Generated Claude docs
    │   ├── features/
    │   ├── services/
    │   ├── integrations/
    │   └── operations/
    └── wiki/                     # Generated Wiki docs
        └── features/
```

---

## Next Steps

### Immediate (Manual)
1. **Complete high-priority YAML files**:
   - Start with features (35 files)
   - Add missing required fields
   - Replace placeholders with actual content

2. **Validate and regenerate**:
   ```bash
   # After completing YAML files
   python validator.py
   python doc_generator.py  # Or batch script
   ```

### Future Automation (Phase 4+)
- [ ] Bulk doc generation script
- [ ] CI/CD integration (validate on PR, regenerate on merge)
- [ ] Change detection (only regenerate modified files)
- [ ] Cross-reference validation
- [ ] Markdown linting integration
- [ ] Link checking integration
- [ ] Screenshot placeholder generation
- [ ] API endpoint documentation extraction from code
- [ ] Auto-update from SOURCE_OF_TRUTH changes

---

## Development

### Requirements
```bash
pip install -r ../../requirements.txt
```

Dependencies:
- PyYAML >= 6.0
- Jinja2 >= 3.1
- jsonschema >= 4.0

### Testing
```bash
# Test SOT parser
python sot_parser.py

# Test markdown parser
python md_parser.py

# Test doc generator (with MOVIE_MODULE)
python doc_generator.py

# Test validator
python validator.py

# Test analyzer
python yaml_analyzer.py

# Test batch migration (dry run)
python batch_migrate.py
```

---

## Design Decisions

### Why YAML as intermediate format?
- **Human-readable** and easy to edit
- **Version control friendly** (line-by-line diffs)
- **Schema-validatable** (via JSON Schema)
- **Separates data from presentation** (templates)
- **Enables dual output** (technical + user-friendly)
- **Supports structured data** (lists, nested objects)

### Why Jinja2 templates?
- **Industry standard** template engine
- **Powerful inheritance** (base + child templates)
- **Conditional rendering** (claude vs wiki)
- **Filter support** (map, join, default, etc.)
- **Python ecosystem** integration

### Why dual output?
- **Claude docs**: Implementation-focused, status tracking, technical details
- **Wiki docs**: User-focused, screenshots, step-by-step guides, FAQs
- Different audiences have different needs
- Single source of truth (YAML) generates both

---

## Troubleshooting

### Template errors
- Check for undefined variables (use `| default('value')`)
- Verify whitespace control (`{%-` vs `{%`)
- Test with minimal data first

### YAML validation failures
- Run validator to see specific errors
- Check required fields for category
- Verify YAML syntax (indentation, lists, etc.)

### Migration issues
- Use dry-run first (`python batch_migrate.py`)
- Check logs for parse errors
- Verify SOURCE_OF_TRUTH structure if SOT parser fails

---

## Contributing

When adding new automation:
1. Follow existing code style
2. Add docstrings and type hints
3. Test with sample data
4. Update this README
5. Add to CI/CD pipeline (future)

---

**Questions?** Check [../docs/dev/design/DESIGN_INDEX.md](../../../docs/dev/design/DESIGN_INDEX.md) or open an issue.
