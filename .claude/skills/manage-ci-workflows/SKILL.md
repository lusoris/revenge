---
name: manage-ci-workflows
description: Manage GitHub Actions workflows (list, validate, trigger, monitor, download logs)
argument-hint: "[--list|--validate|--trigger WORKFLOW|--status|--watch RUN_ID|--logs RUN_ID] [--branch BRANCH]"
disable-model-invocation: false
allowed-tools: Bash(*), Read(*), Write(*)
---

# Manage CI/CD Workflows

Manages GitHub Actions workflows including listing, validation, triggering runs, monitoring status, and downloading logs.

## Usage

```
/manage-ci-workflows --list              # List all workflows
/manage-ci-workflows --validate          # Validate workflow syntax
/manage-ci-workflows --trigger ci.yml    # Trigger workflow run
/manage-ci-workflows --status            # Show recent workflow runs
/manage-ci-workflows --watch 123456      # Watch workflow run in real-time
/manage-ci-workflows --logs 123456       # Download workflow logs
/manage-ci-workflows --view 123456       # View workflow log
/manage-ci-workflows --cancel 123456     # Cancel workflow run
/manage-ci-workflows --rerun 123456      # Rerun workflow
```

## Arguments

- `$0`: Action (--list, --validate, --trigger, --status, --watch, --logs, --view, --cancel, --rerun)
- `$1`: Target (WORKFLOW name or RUN_ID depending on action)
- `$2+`: Options (--branch, --job, --output, --limit, --failed-only)

## Prerequisites

- Python 3.10+ installed
- `gh` CLI installed and authenticated
- Repository access permissions
- actionlint (optional, for advanced validation)

## Task

Manage CI/CD workflows with comprehensive control over workflow execution and monitoring.

### Step 1: Verify Prerequisites

```bash
if [ ! -f "scripts/automation/manage_ci.py" ]; then
    echo "❌ CI management script not found"
    exit 1
fi

if ! command -v gh &> /dev/null; then
    echo "❌ gh CLI not found. Install: https://cli.github.com/"
    exit 1
fi

# Check authentication
if ! gh auth status &> /dev/null; then
    echo "❌ Not authenticated with GitHub"
    echo "Run: gh auth login"
    exit 1
fi
```

### Step 2: List and Validate Workflows

**List all workflows**:
```bash
python scripts/automation/manage_ci.py --list
```

**Validate workflow syntax**:
```bash
# Basic validation (YAML syntax)
python scripts/automation/manage_ci.py --validate

# Advanced validation with actionlint
# Install: brew install actionlint (macOS) or download from GitHub
actionlint
```

**List workflow files**:
```bash
ls -la .github/workflows/
```

### Step 3: Trigger Workflow Runs

**Trigger workflow**:
```bash
# Trigger on default branch (main/develop)
python scripts/automation/manage_ci.py --trigger ci.yml

# Trigger on specific branch
python scripts/automation/manage_ci.py --trigger ci.yml --branch feature/new-feature

# Trigger with inputs (for workflow_dispatch)
gh workflow run ci.yml --ref main -f environment=staging
```

**Available workflows** (check with --list):
- `ci.yml` - Continuous Integration
- `release.yml` - Release automation
- `coverage.yml` - Code coverage reporting
- `security.yml` - Security scanning
- `dev.yml` - Development builds
- `doc-validation.yml` - Documentation validation

### Step 4: Monitor Workflow Runs

**Show recent runs**:
```bash
# Show last 10 runs
python scripts/automation/manage_ci.py --status

# Show last 20 runs
python scripts/automation/manage_ci.py --status --limit 20
```

**Watch run in real-time**:
```bash
# Get run ID from --status, then watch
python scripts/automation/manage_ci.py --watch 123456

# This will stream logs and show progress
```

**Check specific job status**:
```bash
# View run details
gh run view 123456

# View specific job
gh run view 123456 --job test
```

### Step 5: View and Download Logs

**View logs**:
```bash
# View all logs for run
python scripts/automation/manage_ci.py --view 123456

# View specific job logs
python scripts/automation/manage_ci.py --view 123456 --job "test"
```

**Download logs**:
```bash
# Download to default location (logs/run-{id}/)
python scripts/automation/manage_ci.py --logs 123456

# Download to custom location
python scripts/automation/manage_ci.py --logs 123456 --output /tmp/workflow-logs
```

**Search logs**:
```bash
# Download first, then search
/manage-ci-workflows --logs 123456
cd logs/run-123456/
grep -r "error" .
```

### Step 6: Control Workflow Runs

**Cancel running workflow**:
```bash
python scripts/automation/manage_ci.py --cancel 123456
```

**Rerun workflow**:
```bash
# Rerun all jobs
python scripts/automation/manage_ci.py --rerun 123456

# Rerun only failed jobs
python scripts/automation/manage_ci.py --rerun 123456 --failed-only
```

## Common Workflows

### CI/CD Pipeline
1. **Trigger**: Runs on push/PR to main branches
2. **Jobs**:
   - Lint (golangci-lint, ruff, markdownlint)
   - Test (Go, Python, frontend)
   - Build (compile binaries, build Docker images)
   - Security (CodeQL, gitleaks, dependency scan)
3. **Artifacts**: Coverage reports, binaries

### Release Workflow
1. **Trigger**: Automated via Release Please or manual
2. **Jobs**:
   - Build release artifacts
   - Generate changelog
   - Create GitHub release
   - Push Docker images to registry
   - Deploy to staging (optional)

### Documentation Validation
1. **Trigger**: Runs on changes to docs/
2. **Jobs**:
   - Lint markdown
   - Check links
   - Validate YAML
   - Check SOT references

## Examples

**Check failed CI runs**:
```bash
# List recent runs
/manage-ci-workflows --status

# View failed run
/manage-ci-workflows --view 123456

# Download logs for analysis
/manage-ci-workflows --logs 123456

# Rerun failed jobs
/manage-ci-workflows --rerun 123456 --failed-only
```

**Manual deployment**:
```bash
# Trigger release workflow
/manage-ci-workflows --trigger release.yml --branch main

# Watch progress
/manage-ci-workflows --watch 123456

# If it fails, view logs
/manage-ci-workflows --view 123456
```

**Development testing**:
```bash
# Trigger CI on feature branch
/manage-ci-workflows --trigger ci.yml --branch feature/my-feature

# Monitor in real-time
/manage-ci-workflows --watch 123456

# If tests fail, download logs
/manage-ci-workflows --logs 123456
```

## Troubleshooting

**"gh CLI not found"**:
```bash
# Install gh CLI
# macOS: brew install gh
# Linux: https://cli.github.com/

# Authenticate
gh auth login
```

**"Workflow not found"**:
```bash
# List available workflows
/manage-ci-workflows --list

# Check workflow files
ls -la .github/workflows/
```

**"Permission denied"**:
```bash
# Check repository permissions
gh auth status

# May need workflow:write permission
# Re-authenticate with correct scopes
gh auth refresh -s workflow
```

**"Workflow validation failed"**:
```bash
# Check YAML syntax
yamllint .github/workflows/ci.yml

# Install actionlint for better validation
brew install actionlint  # macOS
actionlint
```

**"Run failed immediately"**:
```bash
# Common causes:
# 1. Workflow syntax error
/manage-ci-workflows --validate

# 2. Missing secrets
gh secret list

# 3. Branch protection preventing run
# Check repository settings
```

## Tips

1. **Use dry-run mode**:
   ```bash
   python scripts/automation/manage_ci.py --trigger ci.yml --dry-run
   ```

2. **Monitor long-running workflows**:
   ```bash
   /manage-ci-workflows --watch 123456
   # Ctrl+C to stop watching (doesn't cancel run)
   ```

3. **Analyze failures**:
   ```bash
   # Download logs
   /manage-ci-workflows --logs 123456

   # Use view-logs for searching
   /view-logs --search "error" --run-id 123456
   ```

4. **Automate with scripts**:
   ```bash
   # Trigger + wait for completion
   RUN_ID=$(gh workflow run ci.yml --json databaseId --jq '.databaseId')
   gh run watch $RUN_ID
   ```

## Exit Codes

- `0`: Success
- `1`: Failure (error in operation)

## Related Skills

- `/view-logs` - View and search workflow logs
- `/check-health` - Check system health including CI/CD status
- `/run-all-tests` - Run tests locally before pushing
