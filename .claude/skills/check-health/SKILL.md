---
name: check-health
description: Check system health across all components (automation, services, frontend, resources)
argument-hint: "[--all|--automation|--services|--frontend|--resources] [--json] [--alert]"
disable-model-invocation: false
allowed-tools: Bash(*), Read(*), Write(*)
---

# Check System Health

Checks health of various system components including automation system, backend services, frontend build, and resource usage.

## Usage

```
/check-health                     # Check all components
/check-health --all               # Check all components (explicit)
/check-health --automation        # Check automation system only
/check-health --services          # Check backend services only
/check-health --frontend          # Check frontend build only
/check-health --resources         # Check resource usage only
/check-health --json              # Output in JSON format
/check-health --alert             # Create GitHub issue on failure
```

## Arguments

- `$0`: Check scope (optional: --all, --automation, --services, --frontend, --resources)
- `$1`: Output format (optional: --json)
- `$2`: Alert flag (optional: --alert)

## Prerequisites

- Python 3.10+ with PyYAML, pytest installed
- Docker and docker-compose (for service checks)
- Access to system resources (df command)

## Task

Check system health and report status with clear indicators (healthy/degraded/unhealthy).

### Step 1: Verify Script Exists

```bash
if [ ! -f "scripts/automation/check_health.py" ]; then
    echo "❌ Health checker script not found"
    exit 1
fi
```

### Step 2: Run Health Checks

**Check all components** (default):
```bash
python scripts/automation/check_health.py --all
```

**Check specific components**:
```bash
# Automation system (dependencies, templates, schemas)
python scripts/automation/check_health.py --automation

# Backend services (database, cache, search)
python scripts/automation/check_health.py --services

# Frontend build status
python scripts/automation/check_health.py --frontend

# Resource usage (disk space, etc.)
python scripts/automation/check_health.py --resources
```

**JSON output** (for monitoring tools):
```bash
python scripts/automation/check_health.py --all --json
```

**With GitHub issue creation** on failure:
```bash
python scripts/automation/check_health.py --all --alert
```

### Step 3: Interpret Results

The script will output health status for each component:

- ✅ **healthy**: Component is functioning normally
- ⚠️  **degraded**: Component has warnings but is functional
- ❌ **unhealthy**: Component is not functioning properly

**Health checks performed**:

**Automation System**:
- Python dependencies (PyYAML, pytest)
- Jinja2 templates existence
- JSON schemas existence

**Backend Services**:
- PostgreSQL running (docker-compose)
- Dragonfly cache running
- Typesense search running

**Frontend Build**:
- Frontend directory exists
- node_modules installed
- Build directory exists

**Resource Usage**:
- Disk space (warns at 75%, critical at 90%)

### Step 4: Take Action on Failures

If any components are unhealthy:

1. **Review the error message** for specific failure reason
2. **Check logs** for more details: `/view-logs`
3. **Fix the issue**:
   - Missing dependencies: `pip install -r scripts/requirements.txt`
   - Services not running: `docker-compose up -d`
   - Frontend not built: `cd frontend && npm install && npm run build`
4. **Re-run health check** to verify fix

If `--alert` flag was used, a GitHub issue will be created automatically for unhealthy components.

## Examples

**Quick health check**:
```bash
/check-health
```

**Detailed automation system check**:
```bash
/check-health --automation --verbose
```

**Monitoring integration**:
```bash
# Get JSON output for external monitoring
/check-health --json

# Create alert on failure
/check-health --all --alert
```

## Troubleshooting

**"Python dependencies missing"**:
```bash
pip install -r scripts/requirements.txt
```

**"Templates directory not found"**:
```bash
# Templates should be in templates/ directory
ls -la templates/
```

**"PostgreSQL is not running"**:
```bash
# Start services
docker-compose up -d postgres

# Or start all services
docker-compose up -d
```

**"Disk usage above 90%"**:
```bash
# Check disk usage
df -h

# Clean up Docker resources
docker system prune -a

# Clean up old logs
find logs/ -type f -mtime +30 -delete
```

## Exit Codes

- `0`: All checks passed (healthy or degraded)
- `1`: At least one check failed (unhealthy)

## Related Skills

- `/view-logs` - View and search system logs
- `/manage-docker-config` - Manage Docker configurations
- `/run-all-tests` - Run comprehensive test suites
