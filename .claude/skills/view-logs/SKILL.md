---
name: view-logs
description: View and search logs from GitHub Actions, Docker containers, and local files
argument-hint: "[--workflow|--docker CONTAINER|--local FILE|--search PATTERN] [--follow] [--lines N]"
disable-model-invocation: false
allowed-tools: Bash(*), Read(*), Write(*)
---

# View Logs

View and search logs from multiple sources: GitHub Actions workflows, Docker containers, and local log files.

## Usage

```
/view-logs --workflow                        # List recent workflow runs
/view-logs --view RUN_ID                     # View specific workflow run log
/view-logs --search PATTERN                  # Search local logs for pattern
/view-logs --search PATTERN --run-id RUN_ID  # Search workflow logs
/view-logs --docker postgres                 # View Docker container logs
/view-logs --docker postgres --follow        # Follow Docker logs (tail -f)
/view-logs --local automation.log            # View local log file
/view-logs --local automation.log --follow   # Follow local log file
/view-logs --list                            # List all local log files
```

## Arguments

- `$0`: Action (--workflow, --view, --search, --docker, --local, --list, --download)
- `$1`: Target (RUN_ID, PATTERN, CONTAINER, FILE depending on action)
- `$2+`: Options (--follow, --lines N, --run-id, --job, --status, --case-sensitive)

## Prerequisites

- Python 3.10+ installed
- `gh` CLI installed and authenticated (for GitHub Actions logs)
- Docker and docker-compose (for container logs)
- Access to logs directory (for local logs)

## Task

View and search logs from various sources with filtering and real-time following capabilities.

### Step 1: Verify Script Exists

```bash
if [ ! -f "scripts/automation/view_logs.py" ]; then
    echo "❌ Log viewer script not found"
    exit 1
fi
```

### Step 2: GitHub Actions Workflow Logs

**List recent workflow runs**:
```bash
python scripts/automation/view_logs.py --workflow

# With status filter
python scripts/automation/view_logs.py --workflow --status failure

# Limit results
python scripts/automation/view_logs.py --workflow --limit 10
```

**View specific workflow run**:
```bash
# Get run ID from list command
python scripts/automation/view_logs.py --view 123456

# View specific job
python scripts/automation/view_logs.py --view 123456 --job "test"
```

**Search workflow logs**:
```bash
# Search for pattern in workflow run
python scripts/automation/view_logs.py --search "error" --run-id 123456

# Case-sensitive search
python scripts/automation/view_logs.py --search "Error" --run-id 123456 --case-sensitive
```

**Download workflow logs**:
```bash
# Download to default location (logs/run-{id}/)
python scripts/automation/view_logs.py --download 123456

# Download to custom location
python scripts/automation/view_logs.py --download 123456 --output /tmp/logs
```

### Step 3: Docker Container Logs

**View container logs**:
```bash
# View last 100 lines (default)
python scripts/automation/view_logs.py --docker postgres

# View specific number of lines
python scripts/automation/view_logs.py --docker postgres --lines 50

# Follow logs in real-time (tail -f)
python scripts/automation/view_logs.py --docker postgres --follow
```

**Available containers** (from docker-compose.yml):
- `postgres` - PostgreSQL database
- `dragonfly` - Dragonfly cache
- `typesense` - Typesense search
- `backend` - Go backend (if running in Docker)
- `frontend` - Frontend dev server (if running in Docker)

### Step 4: Local Log Files

**List all local logs**:
```bash
python scripts/automation/view_logs.py --list
```

**View local log file**:
```bash
# View entire file
python scripts/automation/view_logs.py --local automation.log

# View last N lines
python scripts/automation/view_logs.py --local automation.log --lines 50

# Follow log file (tail -f)
python scripts/automation/view_logs.py --local automation.log --follow
```

**Search local logs**:
```bash
# Search all log files in logs/ directory
python scripts/automation/view_logs.py --search "error"

# Case-sensitive search
python scripts/automation/view_logs.py --search "ERROR" --case-sensitive
```

### Step 5: Analyze Results

**Common patterns to search for**:
- `error` or `ERROR` - Error messages
- `failed` or `FAILED` - Failed operations
- `warning` or `WARN` - Warnings
- `exception` or `Exception` - Exceptions
- `timeout` - Timeout issues
- `connection refused` - Network issues
- `permission denied` - Permission issues

**Workflow status filters**:
- `success` - Successful runs only
- `failure` - Failed runs only
- `in_progress` - Currently running

## Examples

**Check failed workflow runs**:
```bash
# List failed runs
/view-logs --workflow --status failure

# View specific failed run
/view-logs --view 123456

# Search for errors
/view-logs --search "error" --run-id 123456
```

**Monitor Docker container**:
```bash
# Follow PostgreSQL logs
/view-logs --docker postgres --follow

# View recent errors
/view-logs --docker postgres --lines 100 | grep -i error
```

**Search automation logs**:
```bash
# List all local logs
/view-logs --list

# Search for specific error
/view-logs --search "validation failed"

# Follow automation log
/view-logs --local automation.log --follow
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

**"Log file not found"**:
```bash
# List available logs
python scripts/automation/view_logs.py --list

# Check logs directory
ls -la logs/
```

**"Container not found"**:
```bash
# Check running containers
docker-compose ps

# Start container
docker-compose up -d postgres
```

**"Permission denied" accessing logs**:
```bash
# Check file permissions
ls -la logs/

# Fix permissions if needed
chmod 644 logs/*.log
```

## Tips

1. **Use grep for additional filtering**:
   ```bash
   /view-logs --docker postgres | grep -i error
   ```

2. **Combine with other tools**:
   ```bash
   # Download logs and analyze
   /view-logs --download 123456
   cd logs/run-123456/
   grep -r "error" .
   ```

3. **Real-time monitoring**:
   ```bash
   # Follow multiple logs in separate terminals
   /view-logs --docker postgres --follow  # Terminal 1
   /view-logs --docker backend --follow   # Terminal 2
   /view-logs --local automation.log --follow  # Terminal 3
   ```

4. **Create log archives**:
   ```bash
   # Download workflow logs for archival
   /view-logs --download 123456 --output /backup/logs/
   ```

## Exit Codes

- `0`: Success (logs retrieved)
- `1`: Failure (error accessing logs)

## Related Skills

- `/check-health` - Check system health status
- `/manage-ci-workflows` - Manage CI/CD workflows
- `/manage-docker-config` - Manage Docker configurations
