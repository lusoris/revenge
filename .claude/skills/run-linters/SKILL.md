---
name: run-linters
description: Run code linters (golangci-lint, ruff, markdownlint, prettier)
argument-hint: "[--all|--go|--python|--markdown|--frontend] [--fix] [--parallel]"
disable-model-invocation: false
allowed-tools: Bash(python scripts/automation/run_linters.py *)
---

# Run Linters

Run code linters across the codebase to check for code quality issues, style violations, and potential bugs.

## Usage

```
/run-linters --all                      # Run all linters (dry-run)
/run-linters --all --fix                # Run all linters with auto-fix
/run-linters --go                       # Run golangci-lint only
/run-linters --python                   # Run ruff only
/run-linters --go --python --fix        # Run Go and Python linters with fixes
/run-linters --markdown                 # Run markdownlint only
/run-linters --frontend                 # Run prettier (TypeScript/JSON/YAML)
/run-linters --all --parallel           # Run linters in parallel
```

## Arguments

- `$0`: Linter selection (--all, --go, --python, --markdown, --frontend, or combination)
- `$1+`: Options (--fix for auto-fix, --parallel for parallel execution)

## Supported Linters

| Linter | Language | Files | Fix Mode |
|--------|----------|-------|----------|
| golangci-lint | Go | `*.go` | Yes |
| ruff | Python | `*.py` | Yes |
| markdownlint | Markdown | `*.md` | Yes |
| prettier | TypeScript/JSON/YAML | `*.ts`, `*.js`, `*.json`, `*.yaml`, `*.yml` | Yes |

## Prerequisites

- Python 3.10+ installed
- `golangci-lint` installed (`brew install golangci-lint` or `https://golangci-lint.run/usage/install/`)
- `ruff` installed (`pip install ruff`)
- `markdownlint-cli` installed (`npm install -g markdownlint-cli`)
- `prettier` installed (`npm install -g prettier`)

## Task

Run code linters with optional auto-fixing and parallel execution.

### Step 1: Verify Prerequisites

```bash
if [ ! -f "scripts/automation/run_linters.py" ]; then
    echo "❌ Linter script not found"
    exit 1
fi
```

### Step 2: Run Selected Linters

**Run all linters (check mode)**:
```bash
python scripts/automation/run_linters.py --all
```

**Run all linters with auto-fix**:
```bash
python scripts/automation/run_linters.py --all --fix
```

**Run specific linters**:
```bash
# Go only
python scripts/automation/run_linters.py --go

# Python only
python scripts/automation/run_linters.py --python

# Markdown only
python scripts/automation/run_linters.py --markdown

# Frontend (TypeScript, JSON, YAML)
python scripts/automation/run_linters.py --frontend
```

**Run multiple linters with fix**:
```bash
python scripts/automation/run_linters.py --go --python --fix
```

**Run linters in parallel**:
```bash
python scripts/automation/run_linters.py --all --parallel
```

### Step 3: Review Results

The script will output:
- Number of issues found by each linter
- File paths with violations
- Severity levels (error, warning)
- Auto-fix results (if --fix used)
- Summary report

### Step 4: Fix Issues Automatically

For issues that support auto-fixing:
```bash
python scripts/automation/run_linters.py --all --fix
```

## Examples

**Check all linters before committing**:
```bash
/run-linters --all
```

**Auto-fix Go code issues**:
```bash
/run-linters --go --fix
```

**Check Python and Markdown**:
```bash
/run-linters --python --markdown
```

**Fix everything in parallel**:
```bash
/run-linters --all --fix --parallel
```

**Check frontend files only**:
```bash
/run-linters --frontend
```

## Linter Details

### golangci-lint

**Configuration**: `.golangci.yml`

**Checks**:
- Syntax errors
- Unused variables and imports
- Ineffective code patterns
- Race conditions
- Security issues
- Code complexity

**Example**:
```bash
python scripts/automation/run_linters.py --go
python scripts/automation/run_linters.py --go --fix
```

### ruff

**Configuration**: `ruff.toml`

**Checks**:
- PEP 8 compliance
- Unused imports
- Undefined names
- Logic errors
- Security issues

**Example**:
```bash
python scripts/automation/run_linters.py --python
python scripts/automation/run_linters.py --python --fix
```

### markdownlint

**Configuration**: `.markdownlint.json`

**Checks**:
- Markdown formatting
- Link validity (local)
- Heading structure
- List formatting

**Example**:
```bash
python scripts/automation/run_linters.py --markdown
python scripts/automation/run_linters.py --markdown --fix
```

### prettier

**Configuration**: `.prettierrc.json`

**Checks**:
- Code formatting (TypeScript, JSON)
- Indentation consistency
- Line length
- Quote styles

**Example**:
```bash
python scripts/automation/run_linters.py --frontend
python scripts/automation/run_linters.py --frontend --fix
```

## Troubleshooting

**"golangci-lint not found"**:
```bash
# Install golangci-lint
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

**"ruff not found"**:
```bash
pip install ruff
```

**"markdownlint not found"**:
```bash
npm install -g markdownlint-cli
```

**"prettier not found"**:
```bash
npm install -g prettier
```

**Linter reports false positives**:
1. Check configuration files (`.golangci.yml`, `ruff.toml`, etc.)
2. Review specific linter documentation
3. Consider adding exclusions for specific files
4. Check for version conflicts

**Fix mode doesn't work for some linters**:
- Some issues require manual fixes
- Check linter documentation for which issues support auto-fix
- Review error messages for guidance

## Tips

1. **Run before committing**:
   ```bash
   /run-linters --all
   ```

2. **Auto-fix during development**:
   ```bash
   /run-linters --all --fix
   ```

3. **Check specific language**:
   ```bash
   /run-linters --go
   /run-linters --python
   /run-linters --frontend
   ```

4. **Use parallel for speed**:
   ```bash
   /run-linters --all --parallel
   ```

5. **Combine with git pre-commit**:
   - See `.githooks/` for pre-commit hook setup

## Exit Codes

- `0`: Success (no issues or all fixed)
- `1`: Linting errors found
- `2`: Configuration error

## Related Skills

- `/format-code` - Format code across all languages
- `/run-all-tests` - Run test suite
- `/manage-ci-workflows` - Manage CI/CD workflows
