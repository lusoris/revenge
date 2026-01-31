# Documentation Automation - Quick Start Guide

Get started with the Revenge documentation automation system in 5 minutes.

---

## Prerequisites

```bash
# Install Python dependencies
pip install -r ../requirements.txt

# Verify installation
python validator.py --help  # Should not error
```

---

## Common Workflows

### 1. Create New Documentation

**For a new feature module (e.g., Podcast Module):**

```bash
# 1. Create YAML data file
cat > ../../data/features/podcasts/PODCASTS.yaml << 'EOF'
doc_title: Podcast Module
doc_category: feature
feature_name: Podcast Module
module_name: podcast
schema_name: public
content_types:
  - Podcasts
  - Episodes
metadata_providers:
  - name: Podcast Index
    purpose: Primary metadata source
    priority: 1

# ... add more fields (see feature.schema.json)
EOF

# 2. Validate YAML
python validator.py

# 3. Generate docs (preview first)
python batch_regenerate.py --preview

# 4. Review preview output
ls -la ../../docs/dev/design-preview/features/podcasts/
ls -la ../../docs/wiki-preview/features/podcasts/

# 5. If good, regenerate to actual locations
python batch_regenerate.py --live --backup

# 6. Commit
git add data/features/podcasts/PODCASTS.yaml
git add docs/dev/design/features/podcasts/PODCASTS.md
git add docs/wiki/features/podcasts/PODCASTS.md
git commit -m "docs: add podcast module documentation"
```

---

### 2. Update Existing Documentation

**Update YAML, regenerate docs:**

```bash
# 1. Edit YAML file
vim ../../data/features/video/MOVIE_MODULE.yaml

# 2. Validate
python validator.py

# 3. Regenerate that specific doc
python -c "
from doc_generator import DocGenerator
from pathlib import Path

gen = DocGenerator(Path.cwd().parent.parent)
gen.generate_doc(
    data_file=Path('../../data/features/video/MOVIE_MODULE.yaml'),
    template_name='feature.md.jinja2',
    output_subpath='features/video',
    render_both=True
)
"

# 4. Review and commit
git diff docs/dev/design/features/video/MOVIE_MODULE.md
git add data/features/video/MOVIE_MODULE.yaml docs/dev/design/features/video/MOVIE_MODULE.md
git commit -m "docs: update movie module metadata providers"
```

---

### 3. Migrate Existing Markdown to YAML

**Already done for 142 docs, but if you have new ones:**

```bash
# 1. Parse markdown to YAML
python md_parser.py

# 2. This creates data/**/*.yaml files with placeholders

# 3. Complete the placeholders manually
vim data/features/new-feature/NEW_FEATURE.yaml

# 4. Validate
python validator.py

# 5. Regenerate docs
python batch_regenerate.py --preview
```

---

### 4. Check Completion Status

**See what needs work:**

```bash
# Run analyzer
python yaml_analyzer.py

# Output shows:
# - Files by category
# - Completion scores
# - Top priority files
# - Missing cross-references
```

---

### 5. Validate Everything (CI Mode)

**Run all checks:**

```bash
# Lenient mode (warnings OK)
python ci_validate.py

# Strict mode (fail on warnings)
python ci_validate.py --strict
```

---

## File Structure Quick Reference

```
revenge/
├── data/                    # YAML source files (edit these!)
│   ├── features/
│   │   └── video/
│   │       └── MOVIE_MODULE.yaml
│   ├── services/
│   ├── integrations/
│   └── operations/
│
├── docs/
│   ├── dev/design/          # Generated Claude docs (don't edit directly)
│   │   └── features/
│   │       └── video/
│   │           └── MOVIE_MODULE.md
│   └── wiki/                # Generated Wiki docs (don't edit directly)
│       └── features/
│           └── video/
│               └── MOVIE_MODULE.md
│
└── scripts/automation/
    ├── doc_generator.py     # YAML → Markdown
    ├── validator.py         # YAML validation
    ├── yaml_analyzer.py     # Completion analysis
    ├── batch_regenerate.py  # Bulk regeneration
    └── ci_validate.py       # CI/CD checks
```

---

## Key Concepts

### YAML is Source of Truth

✅ **DO**: Edit YAML files in `data/`
❌ **DON'T**: Edit generated markdown in `docs/`

### Dual Output

Every YAML file generates:
- **Claude docs** (`docs/dev/design/`) - Technical, implementation-focused
- **Wiki docs** (`docs/wiki/`) - User-friendly, screenshots, FAQs

### Preview Before Live

Always use `--preview` first:
```bash
python batch_regenerate.py --preview  # Safe, outputs to preview dirs
python batch_regenerate.py --live     # Overwrites actual docs
```

### Validation is Required

Before committing:
```bash
python validator.py  # Must pass!
```

---

## Templates

Choose the right template for your doc type:

| Doc Type    | Template                 | Use For                                         |
| ----------- | ------------------------ | ----------------------------------------------- |
| Feature     | `feature.md.jinja2`      | Content modules (movies, music, etc.)           |
| Service     | `service.md.jinja2`      | Backend services (auth, user, etc.)             |
| Integration | `integration.md.jinja2`  | External APIs (TMDb, Radarr, etc.)              |
| Generic     | `generic.md.jinja2`      | Operations, architecture, technical docs        |

---

## Troubleshooting

### "UndefinedError: 'field_name' is undefined"

**Fix**: Add the field to your YAML file or use `| default('value')` in template

### "YAML validation failed"

**Fix**: Check required fields for your category in `schemas/`
```bash
# For features, see: schemas/feature.schema.json
# For services, see: schemas/service.schema.json
# For integrations, see: schemas/integration.schema.json
```

### "Placeholders detected"

**Fix**: Replace `PLACEHOLDER:` values in YAML files
```bash
# Find all placeholders
grep -r "PLACEHOLDER" data/
```

---

## Next Steps

1. **Complete high-priority YAML files**:
   ```bash
   python yaml_analyzer.py  # See what needs work
   ```

2. **Set up CI/CD** (optional):
   ```bash
   # Copy example workflow
   cp .github/workflows/doc-validation.yml.example .github/workflows/doc-validation.yml
   ```

3. **Read full documentation**:
   - [README.md](README.md) - Complete system overview
   - [../../../docs/dev/design/DESIGN_INDEX.md](../../../docs/dev/design/DESIGN_INDEX.md) - Design docs index

---

## Help & Support

- **Check analyzer**: `python yaml_analyzer.py`
- **Validate everything**: `python ci_validate.py`
- **Read full docs**: [README.md](README.md)
- **Open issue**: [GitHub Issues](https://github.com/revenge-project/revenge/issues)

---

**Last Updated**: 2026-01-31
