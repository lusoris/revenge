---
name: generate-badges
description: Generate shields.io badges (coverage, build, version, license)
argument-hint: "[--all|--coverage|--build|--version|--license]"
disable-model-invocation: false
allowed-tools: Bash(*), Read(*)
---

# Generate Badges

Generates shields.io badges for documentation (README, wiki) including coverage, build status, version, and license badges.

## Usage

```
/generate-badges              # Generate all badges
/generate-badges --coverage   # Generate coverage badge only
/generate-badges --build      # Generate build status badge only
/generate-badges --version    # Generate version badge only
/generate-badges --license    # Generate license badge only
```

## Arguments

- `$0`: Badge type (optional: --all, --coverage, --build, --version, --license)

## Prerequisites

- Python 3.10+ with requests library installed
- Git repository with GitHub remote (for build badge)
- coverage.out file (for coverage badge)
- LICENSE file (for license badge)

## Task

Generate shields.io badges and save them to `docs/assets/badges/`.

### Step 1: Verify Prerequisites

Check that the badge generation script and dependencies are available:

```bash
# Check Python is available
if ! command -v python &> /dev/null; then
    echo "❌ Python not found"
    exit 1
fi

# Check requests library is installed
if ! python -c "import requests" 2>/dev/null; then
    echo "❌ requests library not installed"
    echo "Install: pip install requests"
    exit 1
fi

# Verify badge generation script exists
if [ ! -f "scripts/automation/generate_badges.py" ]; then
    echo "❌ Badge generation script not found"
    exit 1
fi

# Verify badges directory exists
if [ ! -d "docs/assets/badges" ]; then
    echo "⚠️ Badges directory not found, creating..."
    mkdir -p docs/assets/badges
fi

# Verify git repository
if [ ! -d ".git" ]; then
    echo "❌ Not a git repository"
    exit 1
fi
```

### Step 2: Generate Badges

Run the badge generation script based on the requested type:

**For all badges** (default or --all):
```bash
echo "🏷️  Generating all badges..."
python scripts/automation/generate_badges.py --all
```

**For coverage badge** (--coverage):
```bash
echo "📊 Generating coverage badge..."
if [ -f "coverage.out" ]; then
    python scripts/automation/generate_badges.py --coverage
else
    echo "⚠️  coverage.out not found, generating placeholder..."
    python scripts/automation/generate_badges.py --coverage
fi
```

**For build badge** (--build):
```bash
echo "🔨 Generating build status badge..."
python scripts/automation/generate_badges.py --build
```

**For version badge** (--version):
```bash
echo "🏷️  Generating version badge..."
python scripts/automation/generate_badges.py --version
```

**For license badge** (--license):
```bash
echo "📜 Generating license badge..."
python scripts/automation/generate_badges.py --license
```

### Step 3: Verify Generation

Check that badges were created successfully:

```bash
echo ""
echo "📊 Generated badges:"
ls -lh docs/assets/badges/ 2>/dev/null | tail -n +2 | awk '{print "  ", $9, "(" $5 ")"}'

echo ""
echo "✅ Badge generation complete!"
```

### Step 4: Show Usage Examples

Provide markdown examples for using the badges:

```bash
echo ""
echo "Usage in README.md:"
echo ""
echo "[![Coverage](docs/assets/badges/coverage.svg)](https://codecov.io/gh/OWNER/REPO)"
echo "[![Build](docs/assets/badges/build.svg)](https://github.com/OWNER/REPO/actions)"
echo "[![Version](docs/assets/badges/version.svg)](https://github.com/OWNER/REPO/releases)"
echo "[![License](docs/assets/badges/license.svg)](LICENSE)"
echo ""
echo "Next steps:"
echo "  1. Add badges to README.md, CONTRIBUTING.md, or wiki pages"
echo "  2. Update badges when coverage/version changes"
echo "  3. Badges auto-update from shields.io (for build status)"
```

## Output Format

The script should output:
- List of generated badge files
- Markdown examples for using badges
- Next steps for integrating badges

## Error Handling

- If requests library is not installed, instruct user to install it
- If coverage.out is missing, generate a placeholder badge (0% coverage)
- If not a GitHub repository, generate a generic build badge
- If LICENSE file is missing, use MIT as default

## Notes

- Badges are downloaded from shields.io API
- Coverage badge color depends on percentage (green ≥80%, yellow ≥60%, red <60%)
- Build badge pulls status from GitHub Actions
- Version badge reads from git tags or go.mod
- License badge auto-detects from LICENSE file
- All badges are SVG format
