# Quick Start - Immediate Fixes for Automation System

**Date**: 2026-01-31
**Goal**: Get automation system from 65% → 90% completion
**Time**: 4-12 hours

---

## Pre-Flight Check

```bash
# Check current status
cd /home/kilian/dev/revenge

# Verify Python environment
python3 --version  # Should be 3.12+
which python3      # Should be in .venv

# Verify dependencies
pip3 list | grep -E "PyYAML|Jinja2|jsonschema|ruff"

# Check git status
git status
```

---

## Fix 1: Update JSON Schemas (30-60 min)

### Problem
Schemas reject valid YAML files due to missing categories and strict regexes.

### Solution

Edit `schemas/feature.schema.json`, `schemas/service.schema.json`, `schemas/integration.schema.json`:

#### 1. Add missing doc_category values

Find the `doc_category` enum and add:

```json
{
  "doc_category": {
    "enum": [
      "feature",
      "service",
      "integration",
      "architecture",
      "operations",
      "technical",
      "pattern",
      "research",
      "other"
    ]
  }
}
```

#### 2. Relax module_name regex

Change from:
```json
"module_name": {
  "pattern": "^[a-z][a-z0-9_]*$"
}
```

To:
```json
"module_name": {
  "pattern": "^[a-z][a-z0-9_\\s\\(\\)&/:,-]*$"
}
```

Or make it optional and sanitize during generation.

#### 3. Make integration fields optional

In `schemas/integration.schema.json`, move `api_base_url` and `auth_method` from `required` array:

```json
{
  "required": [
    "doc_type",
    "doc_category",
    "title",
    "description",
    "provider_name"
  ]
}
```

Remove: `"api_base_url"`, `"auth_method"`

#### 4. Fix service schema patterns

In `schemas/service.schema.json`:

```json
{
  "package_path": {
    "pattern": "^internal/service/[a-z][a-z0-9_/]*$"
  },
  "fx_module": {
    "pattern": "^[a-z][a-z0-9_]*\\.Module$|^\\.Module$"
  }
}
```

### Verify

```bash
python3 scripts/automation/validator.py
```

Expected: Significantly fewer errors (target: <10 invalid files)

---

## Fix 2: Fix Doc Generator Imports (15-30 min)

### Problem
```python
ImportError: attempted relative import with no known parent package
```

### Solution Option A: Add __init__.py

```bash
touch scripts/__init__.py
touch scripts/automation/__init__.py
```

### Solution Option B: Fix imports in doc_generator.py

Edit `scripts/automation/doc_generator.py`:

Change:
```python
from .toc_generator import TOCGenerator
```

To:
```python
import sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).parent.parent.parent))
from scripts.automation.toc_generator import TOCGenerator
```

### Solution Option C: Run as module

Document that scripts must be run as:
```bash
python3 -m scripts.automation.doc_generator --help
```

### Verify

```bash
python3 scripts/automation/doc_generator.py --help
# OR
python3 -m scripts.automation.doc_generator --help
```

Expected: Help text displayed, no import errors

---

## Fix 3: Install Missing npm Tools (5-10 min)

### Problem
Validation pipeline incomplete without npm tools.

### Solution

```bash
# Check if npm is installed
npm --version

# Install tools globally
npm install -g markdownlint-cli markdown-link-check

# Verify installation
which markdownlint
which markdown-link-check

# Test
markdownlint --version
markdown-link-check --version
```

### Verify

```bash
# Test markdownlint
markdownlint docs/dev/design/00_SOURCE_OF_TRUTH.md

# Test link checker
markdown-link-check docs/dev/design/00_SOURCE_OF_TRUTH.md
```

---

## Fix 4: Create CODEOWNERS (10 min)

### Problem
Missing CODEOWNERS file prevents auto-reviewer assignment.

### Solution

Create `CODEOWNERS` file:

```bash
cat > CODEOWNERS <<'EOF'
# CODEOWNERS - Auto-assign reviewers
# See: https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners

# Default owner for everything
* @kilian

# Backend (Go code)
/cmd/ @kilian
/internal/ @kilian
/pkg/ @kilian
*.go @kilian
go.mod @kilian
go.sum @kilian

# Frontend (SvelteKit/Svelte)
/web/ @kilian
*.svelte @kilian
*.ts @kilian

# Documentation
/docs/ @kilian
*.md @kilian

# Infrastructure
/.github/ @kilian
/.coder/ @kilian
/docker-compose.yml @kilian
/Dockerfile @kilian

# Automation scripts
/scripts/ @kilian
*.py @kilian

# Configuration
/.vscode/ @kilian
/.zed/ @kilian
/.editorconfig @kilian
/ruff.toml @kilian
/.markdownlint* @kilian

# Database
/migrations/ @kilian
*.sql @kilian

# Templates & Schemas
/templates/ @kilian
/schemas/ @kilian
/data/ @kilian
EOF
```

### Verify

```bash
cat CODEOWNERS
git add CODEOWNERS
```

---

## Fix 5: Update SOURCE_OF_TRUTH.md (20-30 min)

### Problem
Config sync can't extract versions from SOT.

### Solution

Add Development Tools section to `docs/dev/design/00_SOURCE_OF_TRUTH.md`:

```markdown
## Development Tools

Tools used for development, with exact versions synced to configs.

| Tool | Version | Config Files | Notes |
|------|---------|--------------|-------|
| **Go** | 1.25.6 | `go.mod`, `.tool-versions`, `.github/workflows/_versions.yml` | Latest stable (2026-01-15) |
| **Python** | 3.12+ | `.python-version`, `.tool-versions`, `scripts/requirements.txt` | 3.14.2 recommended |
| **Node** | 20.x | `.nvmrc`, `.tool-versions`, `package.json` | LTS version |
| **PostgreSQL** | 18.1 | `docker-compose.yml`, `.github/workflows/*.yml` | Latest stable |
| **Docker** | 27+ | `.github/workflows/*.yml` | Docker Engine |
| **Dragonfly** | v1.36.0 | `docker-compose.yml`, `.coder/template.tf` | Redis-compatible cache |
| **Typesense** | v30.1 | `docker-compose.yml`, `.coder/template.tf` | Search engine |
| **golangci-lint** | v1.61.0 | `.github/workflows/ci.yml`, `.golangci.yml` | Go linter |
| **ruff** | 0.4+ | `scripts/requirements.txt`, `.github/workflows/*.yml`, `ruff.toml` | Python linter/formatter |
| **Coder** | v2.17.2+ | `.coder/template.tf` | Dev environments |
| **markdownlint** | 0.39+ | `.markdownlint.json`, `.github/workflows/*.yml` | Markdown linter |
| **gitleaks** | 8.18+ | `.github/workflows/security.yml`, `.gitleaksignore` | Secret scanner |
```

### Verify

```bash
python3 scripts/automation/sot_parser.py
# Check that data/shared-sot.yaml includes development_tools section

python3 scripts/automation/config_sync.py --dry-run
# Should now extract versions correctly
```

---

## Fix 6: Test Key Scripts (60-120 min)

### Problem
Unknown if automation scripts work.

### Solution

Test each critical script:

```bash
# 1. SOT Parser (should already work)
python3 scripts/automation/sot_parser.py
ls -lh data/shared-sot.yaml

# 2. Validator (should work after schema fixes)
python3 scripts/automation/validator.py

# 3. Doc Generator (should work after import fixes)
python3 -m scripts.automation.doc_generator \
  --file data/features/video/MOVIE_MODULE.yaml \
  --output /tmp/test_movie.md

# 4. Config Sync (should work after SOT update)
python3 scripts/automation/config_sync.py --dry-run

# 5. Batch Migration (should already work)
python3 scripts/automation/batch_migrate.py --dry-run

# 6. Batch Regenerate (test after generator works)
python3 scripts/automation/batch_regenerate.py --dry-run

# 7. Health Check
python3 scripts/automation/check_health.py

# 8. Format Code
python3 scripts/automation/format_code.py --dry-run

# 9. Run Linters
python3 scripts/automation/run_linters.py --dry-run

# 10. Check Licenses
python3 scripts/automation/check_licenses.py --dry-run
```

Document any errors in `.analysis/TEST_RESULTS.md`

---

## Fix 7: Test GitHub Scripts (30-60 min)

### Problem
GitHub management scripts untested.

### Solution

```bash
# Requires gh CLI authenticated
gh auth status

# Test each GitHub script in dry-run mode
python3 scripts/automation/github_labels.py --dry-run
python3 scripts/automation/github_milestones.py --dry-run
python3 scripts/automation/github_projects.py --dry-run
python3 scripts/automation/github_discussions.py --dry-run
python3 scripts/automation/github_security.py --dry-run
```

Note: Some may require `--dry-run` flag to be added if missing.

---

## Verification Checklist

After completing fixes, verify:

### ✅ Validation Fixed
- [ ] Run `python3 scripts/automation/validator.py`
- [ ] Target: <10 invalid files (preferably 0)
- [ ] Document remaining issues

### ✅ Generation Working
- [ ] Run doc generator on 1 file
- [ ] Verify markdown output
- [ ] Verify TOC generated
- [ ] Verify frontmatter correct

### ✅ Tools Installed
- [ ] `markdownlint --version` works
- [ ] `markdown-link-check --version` works
- [ ] Test both tools on a doc

### ✅ Configs Updated
- [ ] CODEOWNERS file created
- [ ] SOURCE_OF_TRUTH.md has Development Tools table
- [ ] Config sync runs without errors

### ✅ Scripts Tested
- [ ] At least 10 critical scripts tested
- [ ] Document any failures
- [ ] Create issue for each failure

---

## Quick Commands Reference

```bash
# Validate all YAML
python3 scripts/automation/validator.py

# Generate single doc
python3 -m scripts.automation.doc_generator --file data/features/video/MOVIE_MODULE.yaml

# Batch regenerate all docs
python3 scripts/automation/batch_regenerate.py --dry-run

# Sync configs from SOT
python3 scripts/automation/config_sync.py --dry-run

# Parse SOT
python3 scripts/automation/sot_parser.py

# Check system health
python3 scripts/automation/check_health.py

# Format all code
python3 scripts/automation/format_code.py --dry-run

# Run all linters
python3 scripts/automation/run_linters.py --dry-run

# Validate links
markdown-link-check docs/dev/design/00_SOURCE_OF_TRUTH.md

# Lint markdown
markdownlint docs/dev/design/*.md
```

---

## After Fixes Complete

### Next Steps

1. **Create test suite** (Priority 6)
   - Start with `tests/test_validator.py`
   - Add `tests/test_doc_generator.py`
   - Target: 80% coverage

2. **Configure GitHub** (Priority 5)
   - Create revenge-bot account
   - Enable branch protection
   - Set up Projects
   - Enable Discussions

3. **Full integration test**
   - Change SOURCE_OF_TRUTH.md
   - Verify config sync triggers
   - Verify doc generation triggers
   - Verify PR creation

4. **Documentation**
   - Document each script's usage
   - Create troubleshooting guide
   - Update .claude/CLAUDE.md

---

## Troubleshooting

### Import errors persist

```bash
# Ensure scripts is a package
ls scripts/__init__.py scripts/automation/__init__.py

# Run as module
python3 -m scripts.automation.script_name

# Check PYTHONPATH
echo $PYTHONPATH
export PYTHONPATH=/home/kilian/dev/revenge:$PYTHONPATH
```

### Validator still fails

```bash
# Check schema syntax
python3 -m json.tool schemas/feature.schema.json
python3 -m json.tool schemas/service.schema.json
python3 -m json.tool schemas/integration.schema.json

# Validate individual file
python3 scripts/automation/validator.py --file data/features/video/MOVIE_MODULE.yaml
```

### npm tools not found

```bash
# Check npm prefix
npm prefix -g

# Add to PATH
export PATH="$(npm prefix -g)/bin:$PATH"

# Or install locally
npm install markdownlint-cli markdown-link-check
./node_modules/.bin/markdownlint --version
```

---

## Success Metrics

**Before Fixes**:
- ❌ 83/143 YAML files invalid (58%)
- ❌ Doc generator doesn't run
- ❌ 0 automation scripts tested
- ❌ Missing critical configs

**After Fixes** (Target):
- ✅ <10 YAML files invalid (<7%)
- ✅ Doc generator works end-to-end
- ✅ 10+ automation scripts tested (29%)
- ✅ Critical configs created

**Overall Completion**:
- Before: 65%
- After: 85-90%

---

## Time Estimates

| Fix | Time | Priority |
|-----|------|----------|
| 1. Update schemas | 30-60 min | 🔴 Critical |
| 2. Fix imports | 15-30 min | 🔴 Critical |
| 3. Install npm tools | 5-10 min | 🟡 High |
| 4. Create CODEOWNERS | 10 min | 🟡 High |
| 5. Update SOT | 20-30 min | 🟡 High |
| 6. Test scripts | 60-120 min | 🟡 High |
| 7. Test GitHub | 30-60 min | 🟢 Medium |
| **Total** | **2.5-5 hours** | |

**Recommended**: Do fixes 1-5 in order (90 min - 2 hours), then 6-7 time permitting.

---

**Status**: Ready to execute
**Prerequisites**: Python venv activated, git clean working tree
**Output**: 85-90% completion, working doc generation pipeline

**See also**:
- [EXECUTIVE_SUMMARY.md](EXECUTIVE_SUMMARY.md) - High-level overview
- [PHASE_AUDIT_REPORT.md](PHASE_AUDIT_REPORT.md) - Detailed audit
