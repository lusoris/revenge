---
name: format-code
description: Format code across all languages (Go, Python, TypeScript, JSON, YAML, Markdown)
argument-hint: "[--all|--go|--python|--frontend] [--check] [--parallel]"
disable-model-invocation: false
allowed-tools: Bash(python scripts/automation/format_code.py *)
---

# Format Code

Format code across all languages in the codebase to maintain consistent style and conventions.

## Usage

```
/format-code --all                      # Format all code (apply changes)
/format-code --all --check              # Check formatting without changes
/format-code --go                       # Format Go code only
/format-code --python                   # Format Python code only
/format-code --frontend                 # Format TypeScript/JavaScript/JSON/YAML
/format-code --go --python --check      # Check Go and Python formatting
/format-code --all --parallel           # Format all code in parallel
```

## Arguments

- `$0`: Language selection (--all, --go, --python, --frontend, or combination)
- `$1+`: Options (--check to verify without changes, --parallel for parallel execution)

## Supported Formatters

| Formatter | Language | Files | Tools |
|-----------|----------|-------|-------|
| gofmt + goimports | Go | `*.go` | gofmt, goimports |
| ruff | Python | `*.py` | ruff format |
| prettier | TypeScript/JavaScript/JSON/YAML | `*.ts`, `*.js`, `*.json`, `*.yaml`, `*.yml` | prettier |
| prettier | Markdown | `*.md` | prettier |

## Prerequisites

- Python 3.10+ installed
- Go 1.25+ installed (includes gofmt)
- `goimports` installed (`go install golang.org/x/tools/cmd/goimports@latest`)
- `ruff` installed (`pip install ruff`)
- `prettier` installed (`npm install -g prettier`)

## Task

Format code across multiple languages with check mode and parallel execution support.

### Step 1: Verify Prerequisites

```bash
if [ ! -f "scripts/automation/format_code.py" ]; then
    echo "❌ Format script not found"
    exit 1
fi
```

### Step 2: Check Formatting (Dry-Run)

**Check all code without changes**:
```bash
python scripts/automation/format_code.py --all --check
```

**Check specific languages**:
```bash
# Go only
python scripts/automation/format_code.py --go --check

# Python only
python scripts/automation/format_code.py --python --check

# Frontend only
python scripts/automation/format_code.py --frontend --check
```

### Step 3: Apply Formatting

**Format all code**:
```bash
python scripts/automation/format_code.py --all
```

**Format specific languages**:
```bash
# Go only
python scripts/automation/format_code.py --go

# Python only
python scripts/automation/format_code.py --python

# Frontend only
python scripts/automation/format_code.py --frontend
```

**Format multiple languages**:
```bash
python scripts/automation/format_code.py --go --python --frontend
```

### Step 4: Run in Parallel

**Format all code in parallel for speed**:
```bash
python scripts/automation/format_code.py --all --parallel
```

### Step 5: Review Results

The script will output:
- Number of files formatted
- Files that were changed
- Summary report
- Exit code indicating success/failure

## Examples

**Check formatting before committing**:
```bash
/format-code --all --check
```

**Auto-format Go code**:
```bash
/format-code --go
```

**Format Python and frontend together**:
```bash
/format-code --python --frontend
```

**Quick parallel format**:
```bash
/format-code --all --parallel
```

**Check only without changes**:
```bash
/format-code --all --check
```

## Formatter Details

### Go (gofmt + goimports)

**Configuration**: None (Go standard formatting)

**Features**:
- Consistent indentation (tabs)
- Brace placement
- Import organization and cleanup
- Line length preservation

**Example**:
```bash
python scripts/automation/format_code.py --go
```

**Common issues fixed**:
- Unused imports removed
- Imports sorted alphabetically
- Inconsistent spacing

### Python (ruff)

**Configuration**: `ruff.toml`

**Features**:
- PEP 8 compliance
- Consistent quote styles
- Import sorting
- Indentation consistency

**Example**:
```bash
python scripts/automation/format_code.py --python
```

**Common issues fixed**:
- Quote normalization
- Whitespace cleanup
- Import organization

### Frontend (prettier)

**Configuration**: `.prettierrc.json`

**Features**:
- TypeScript/JavaScript formatting
- JSON formatting
- YAML formatting
- Markdown formatting
- Consistent line length (80 chars)

**Example**:
```bash
python scripts/automation/format_code.py --frontend
```

**Common issues fixed**:
- Indentation (2 spaces)
- Quote styles (double quotes)
- Trailing commas
- Line wrapping

## Workflow Integration

### Before Committing

```bash
# Check formatting
/format-code --all --check

# If issues, apply fixes
/format-code --all

# Then commit
git add .
git commit -m "fix: format code"
```

### In CI/CD

**Check mode** (fails if formatting needed):
```bash
python scripts/automation/format_code.py --all --check
```

**Fix mode** (applies fixes):
```bash
python scripts/automation/format_code.py --all
```

## Troubleshooting

**"goimports not found"**:
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

**"ruff not found"**:
```bash
pip install ruff
```

**"prettier not found"**:
```bash
npm install -g prettier
```

**Formatter conflicts with linter**:
1. Ensure formatter configs match linter configs
2. Check `.golangci.yml`, `ruff.toml`, `.prettierrc.json`
3. Run linters after formatting to verify no conflicts

**Files not being formatted**:
1. Check file extensions match supported types
2. Verify formatter is in PATH
3. Review script output for errors
4. Check configuration for exclude patterns

**"Permission denied" errors**:
```bash
# Check file permissions
ls -la path/to/file

# Fix if needed
chmod 644 path/to/file
```

## Tips

1. **Format before linting**:
   ```bash
   /format-code --all
   /run-linters --all
   ```

2. **Use check mode in pre-commit hook**:
   ```bash
   python scripts/automation/format_code.py --all --check
   ```

3. **Format specific language**:
   ```bash
   /format-code --go
   /format-code --python
   /format-code --frontend
   ```

4. **Run parallel for faster processing**:
   ```bash
   /format-code --all --parallel
   ```

5. **Combine with git commands**:
   ```bash
   /format-code --all && git add -u && git commit -m "chore: format code"
   ```

## Exit Codes

- `0`: Success (formatted or already formatted)
- `1`: Formatting errors
- `2`: Configuration error

## Performance Notes

- Without `--parallel`: Sequential formatting (slower but safer)
- With `--parallel`: Concurrent formatting (faster, safe for most cases)
- Check mode is faster than fix mode

## Related Skills

- `/run-linters` - Run code linters
- `/run-all-tests` - Run test suite
- `/manage-ci-workflows` - Manage CI/CD workflows
