---
name: update-dependencies
description: Update Go, npm, and Python dependencies to latest versions
argument-hint: "[--all|--go|--npm|--python] [--major] [--check] [--dry-run]"
disable-model-invocation: false
allowed-tools: Bash(python scripts/automation/update_dependencies.py *)
---

# Update Dependencies

Update dependencies across all package managers (Go, npm, Python) with options for major version updates and safety checks.

## Usage

```
/update-dependencies --all              # Update all (minor/patch only)
/update-dependencies --go               # Update Go dependencies only
/update-dependencies --npm              # Update npm (frontend) dependencies
/update-dependencies --python           # Update Python dependencies
/update-dependencies --all --major      # Update including major versions
/update-dependencies --all --check      # Check for available updates
/update-dependencies --all --dry-run    # Simulate updates without changes
```

## Arguments

- `$0`: Package manager (--all, --go, --npm, --python, or combination)
- `$1+`: Options (--major for major versions, --check for info only, --dry-run for simulation)

## Update Types

| Type | Go | npm | Python | Safety |
|------|----|----|--------|--------|
| Patch (1.0.0 -> 1.0.1) | Yes | Yes | Yes | High |
| Minor (1.0.0 -> 1.1.0) | Yes | Yes | Yes | High |
| Major (1.0.0 -> 2.0.0) | Optional | Optional | Optional | Low |

## Prerequisites

- Python 3.10+ installed
- `go` 1.25+ installed (for Go updates)
- `npm` or `pnpm` installed (for frontend updates)
- `pip` installed (for Python updates)
- Git configured for commits

## Task

Update dependencies safely with version constraints and testing.

### Step 1: Verify Prerequisites

```bash
if [ ! -f "scripts/automation/update_dependencies.py" ]; then
    echo "❌ Dependency update script not found"
    exit 1
fi
```

### Step 2: Check Available Updates

**Check all available updates**:
```bash
python scripts/automation/update_dependencies.py --all --check
```

**Check specific manager**:
```bash
python scripts/automation/update_dependencies.py --go --check
python scripts/automation/update_dependencies.py --npm --check
python scripts/automation/update_dependencies.py --python --check
```

### Step 3: Dry-Run Updates

**Simulate all updates**:
```bash
python scripts/automation/update_dependencies.py --all --dry-run
```

**Simulate specific manager**:
```bash
python scripts/automation/update_dependencies.py --go --dry-run
```

### Step 4: Update Minor/Patch Only

**Update all (safe)**:
```bash
python scripts/automation/update_dependencies.py --all
```

**Update specific manager**:
```bash
python scripts/automation/update_dependencies.py --go
python scripts/automation/update_dependencies.py --npm
python scripts/automation/update_dependencies.py --python
```

### Step 5: Update Including Major Versions

**Warning**: Major version updates may break compatibility

```bash
python scripts/automation/update_dependencies.py --all --major
```

**Specific manager with major versions**:
```bash
python scripts/automation/update_dependencies.py --go --major
```

### Step 6: Run Tests After Updates

```bash
# Run test suite to verify updates don't break anything
python scripts/automation/run_tests.py --all
```

## Examples

**Check for updates**:
```bash
/update-dependencies --all --check
```

**Safe update (minor/patch only)**:
```bash
/update-dependencies --all
```

**Dry-run to see what would change**:
```bash
/update-dependencies --all --dry-run
```

**Update Go dependencies**:
```bash
/update-dependencies --go
```

**Update npm with major versions**:
```bash
/update-dependencies --npm --major
```

**Full update with testing**:
```bash
/update-dependencies --all && /run-all-tests --all
```

## Go Dependencies

**File**: `go.mod`

**Update strategies**:
```bash
# Minor/patch only (safe)
python scripts/automation/update_dependencies.py --go

# Check available updates
python scripts/automation/update_dependencies.py --go --check

# Include major versions (caution)
python scripts/automation/update_dependencies.py --go --major

# Dry-run
python scripts/automation/update_dependencies.py --go --dry-run
```

**Common updates**:
- Framework updates (ogen, ogen-api)
- Database drivers (pgx, pgxpool)
- Cache libraries (rueidis, otter, sturdyc)
- Testing tools (testify, mockery)

**After update**:
```bash
go mod tidy
go test ./...
```

## npm Dependencies

**File**: `frontend/package.json`

**Update strategies**:
```bash
# Minor/patch only (safe)
python scripts/automation/update_dependencies.py --npm

# Check updates
python scripts/automation/update_dependencies.py --npm --check

# Include major versions
python scripts/automation/update_dependencies.py --npm --major

# Dry-run
python scripts/automation/update_dependencies.py --npm --dry-run
```

**Common updates**:
- SvelteKit and Svelte
- TanStack Query
- Tailwind CSS
- shadcn-svelte components
- UI libraries

**After update**:
```bash
npm install
npm run build
npm run test
```

## Python Dependencies

**Files**: `requirements.txt`, `setup.py`

**Update strategies**:
```bash
# Minor/patch only
python scripts/automation/update_dependencies.py --python

# Check updates
python scripts/automation/update_dependencies.py --python --check

# Include major versions
python scripts/automation/update_dependencies.py --python --major

# Dry-run
python scripts/automation/update_dependencies.py --python --dry-run
```

**Common updates**:
- Testing frameworks (pytest, testify)
- Code quality tools (ruff, black, flake8)
- Documentation generators
- Automation scripts

**After update**:
```bash
pip install -r requirements.txt
python -m pytest
```

## Safety Recommendations

### Before Updating

1. **Commit current code**:
   ```bash
   git add -A
   git commit -m "chore: pre-dependency-update backup"
   ```

2. **Check available updates**:
   ```bash
   /update-dependencies --all --check
   ```

3. **Review breaking changes**:
   - Check CHANGELOG for updated packages
   - Review migration guides for major versions
   - Check for deprecated APIs

### During Update

1. **Use --dry-run first**:
   ```bash
   /update-dependencies --all --dry-run
   ```

2. **Update minor/patch only**:
   ```bash
   /update-dependencies --all
   ```

3. **Review changes**:
   ```bash
   git diff go.mod
   git diff frontend/package.json
   git diff requirements.txt
   ```

### After Update

1. **Run full test suite**:
   ```bash
   /run-all-tests --all
   ```

2. **Build application**:
   ```bash
   go build ./...
   npm run build
   ```

3. **Manual testing**:
   - Test critical features
   - Check for warnings/deprecations
   - Verify no breaking changes

4. **Commit updates**:
   ```bash
   git add -A
   git commit -m "chore(deps): update dependencies"
   ```

## Workflow Integration

### Dependabot Alternative

Instead of relying on Dependabot, use this skill for controlled updates:
```bash
# Weekly updates (manual or scheduled)
/update-dependencies --all --check
/update-dependencies --all
/run-all-tests --all
```

### CI/CD Integration

```bash
# In GitHub Actions
python scripts/automation/update_dependencies.py --all --check
```

## Troubleshooting

**"go: version not found" errors**:
1. Check Go version is 1.25+: `go version`
2. Run `go mod tidy`: `go mod tidy`
3. Run `go mod vendor`: `go mod vendor` (if using vendor)

**npm conflicts after update**:
1. Delete `package-lock.json`: `rm frontend/package-lock.json`
2. Reinstall: `npm install`
3. Resolve remaining conflicts manually

**Major version breaks tests**:
1. Review breaking changes in package docs
2. Update code to match new API
3. Check for migration guides
4. Consider keeping older version if upgrade too complex

**Circular dependencies detected**:
1. Review `go mod graph` output
2. Check for import cycles
3. Refactor to break cycles
4. Consider splitting modules

**Memory issues during update**:
```bash
# For large dependency trees
GOGC=50 python scripts/automation/update_dependencies.py --go
```

## Tips

1. **Regular updates keep codebase healthy**:
   ```bash
   # Weekly or monthly
   /update-dependencies --all --check
   ```

2. **Separate commits for each manager**:
   ```bash
   /update-dependencies --go
   git commit -m "chore(deps): update go dependencies"
   /update-dependencies --npm
   git commit -m "chore(deps): update npm dependencies"
   ```

3. **Major version updates need careful review**:
   ```bash
   /update-dependencies --all --major --dry-run
   # Review changes carefully before applying
   ```

4. **Always test after updates**:
   ```bash
   /update-dependencies --all && /run-all-tests --all
   ```

5. **Create branch for updates**:
   ```bash
   git checkout -b chore/update-dependencies
   /update-dependencies --all
   /run-all-tests --all
   # Then create PR
   ```

## Exit Codes

- `0`: Success (dependencies updated or no updates available)
- `1`: Update failed
- `2`: Configuration error
- `3`: Test failures after update

## Related Skills

- `/check-licenses` - Verify license compliance after updates
- `/run-all-tests` - Test suite validation
- `/manage-ci-workflows` - CI/CD workflows
- `/validate-tools` - Verify tool versions
