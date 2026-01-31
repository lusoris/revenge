---
name: check-licenses
description: Check software licenses in dependencies (Go, npm, Python)
argument-hint: "[--all|--go|--npm|--python] [--report] [--excluded]"
disable-model-invocation: false
allowed-tools: Bash(python scripts/automation/check_licenses.py *)
---

# Check Licenses

Validate software licenses in all dependencies to ensure compliance with project policy and legal requirements.

## Usage

```
/check-licenses --all                   # Check all dependencies
/check-licenses --go                    # Check Go dependencies only
/check-licenses --npm                   # Check npm (frontend) dependencies
/check-licenses --python                # Check Python dependencies
/check-licenses --all --report          # Generate detailed license report
/check-licenses --all --excluded        # Show excluded/problematic licenses
```

## Arguments

- `$0`: Package manager (--all, --go, --npm, --python, or combination)
- `$1+`: Options (--report for detailed report, --excluded for problematic licenses)

## License Categories

### Permissive Licenses (Approved)
- MIT
- Apache 2.0
- BSD (2-Clause, 3-Clause)
- ISC
- MPL 2.0

### Copyleft Licenses (Review Required)
- GPL v2/v3
- AGPL v3
- LGPL v2/v3

### Proprietary/Commercial
- Commercial licenses
- Proprietary software

## Prerequisites

- Python 3.10+ installed
- `go mod` command available
- `npm` or `pnpm` installed (for frontend)
- `pip` installed (for Python dependencies)

## Task

Check software licenses across package managers and generate compliance reports.

### Step 1: Verify Prerequisites

```bash
if [ ! -f "scripts/automation/check_licenses.py" ]; then
    echo "❌ License checker script not found"
    exit 1
fi
```

### Step 2: Check Go Dependencies

**Check all Go dependencies**:
```bash
python scripts/automation/check_licenses.py --go
```

**Show only problematic licenses**:
```bash
python scripts/automation/check_licenses.py --go --excluded
```

**Generate detailed Go license report**:
```bash
python scripts/automation/check_licenses.py --go --report
```

### Step 3: Check npm Dependencies

**Check all npm dependencies (frontend)**:
```bash
python scripts/automation/check_licenses.py --npm
```

**Check npm with report**:
```bash
python scripts/automation/check_licenses.py --npm --report
```

### Step 4: Check Python Dependencies

**Check all Python dependencies**:
```bash
python scripts/automation/check_licenses.py --python
```

**Check Python with report**:
```bash
python scripts/automation/check_licenses.py --python --report
```

### Step 5: Check All Dependencies

**Complete license audit**:
```bash
python scripts/automation/check_licenses.py --all
```

**Complete license audit with detailed report**:
```bash
python scripts/automation/check_licenses.py --all --report
```

**Find all problematic licenses**:
```bash
python scripts/automation/check_licenses.py --all --excluded
```

### Step 6: Review Results

The script will output:
- Summary of checked dependencies
- License classification (permissive, copyleft, proprietary)
- Problematic licenses (if any)
- Recommendations
- Detailed report (if --report used)

## Examples

**Quick license check before release**:
```bash
/check-licenses --all
```

**Check Go dependencies only**:
```bash
/check-licenses --go
```

**Generate full compliance report**:
```bash
/check-licenses --all --report
```

**Find problematic licenses**:
```bash
/check-licenses --all --excluded
```

**Check frontend (npm) dependencies**:
```bash
/check-licenses --npm --report
```

## License Compliance Policy

### Approved Licenses
- **MIT** - Fully permissive
- **Apache 2.0** - Patent protection included
- **BSD** - Permissive, requires notice
- **ISC** - Minimal, permissive
- **MPL 2.0** - File-level copyleft

### Requires Review
- **GPL v2/v3** - Strong copyleft (full codebase must be GPL)
- **AGPL v3** - Network copyleft (service usage triggers GPL)
- **LGPL v2/v3** - Weaker copyleft (can link, but library updates required)

### Typically Rejected
- Unknown/Custom licenses
- Commercial licenses (without purchase)
- Proprietary software

## Go Dependencies

**Check location**: `go.mod`

**Tools used**:
- `go mod graph` - Dependency graph
- `go list -m all` - All modules
- License detection (scanning source)

**Example**:
```bash
python scripts/automation/check_licenses.py --go
```

**Output shows**:
- Module name
- Version
- Detected license
- License type (permissive/copyleft)
- Status (OK/review/error)

## npm Dependencies

**Check location**: `frontend/package.json`

**Tools used**:
- `npm ls --json` - Dependency tree
- License field parsing
- Dual license detection

**Example**:
```bash
python scripts/automation/check_licenses.py --npm
```

**Output shows**:
- Package name
- Version
- License(s)
- Status for each license

## Python Dependencies

**Check location**: `requirements.txt` and `setup.py`

**Tools used**:
- `pip` metadata
- License field from package metadata
- SPDX license detection

**Example**:
```bash
python scripts/automation/check_licenses.py --python
```

## Report Generation

**Detailed compliance report**:
```bash
python scripts/automation/check_licenses.py --all --report
```

**Report includes**:
- Summary statistics
- License distribution
- Risk assessment
- Problematic licenses
- Remediation steps
- Export formats (JSON, CSV)

## Workflow Integration

### Before Release

```bash
# Check all licenses
/check-licenses --all

# Review if any issues
if [ $? -ne 0 ]; then
    # Generate full report
    /check-licenses --all --report
    exit 1
fi
```

### CI/CD Integration

```bash
# In GitHub Actions
python scripts/automation/check_licenses.py --all
```

### License Whitelist

**Update for approved licenses**:
1. Review flagged license in report
2. Verify licensing is acceptable
3. Add to whitelist in configuration
4. Re-run check to confirm

## Troubleshooting

**"License detection failed"**:
1. Manual license inspection
2. Check if source code includes LICENSE file
3. Check package metadata manually
4. Mark as manual review

**Dual licensed packages**:
- Checks all listed licenses
- Needs only one to be approved
- Report shows all options

**Unknown license detected**:
```bash
# Generate report for manual review
python scripts/automation/check_licenses.py --all --report

# Then manually review and whitelist if appropriate
```

**Missing license information**:
1. Check package documentation
2. Review GitHub repository
3. Contact package maintainer
4. Mark as "to be determined"

**Permission denied checking packages**:
```bash
# Verify access to node_modules, vendor, site-packages
ls -la frontend/node_modules/
ls -la vendor/
python -m site
```

## Tips

1. **Run before every release**:
   ```bash
   /check-licenses --all
   ```

2. **Generate report for compliance documentation**:
   ```bash
   /check-licenses --all --report > LICENSE_REPORT.md
   ```

3. **Check specific package ecosystem**:
   ```bash
   /check-licenses --go
   /check-licenses --npm
   /check-licenses --python
   ```

4. **Monitor problematic licenses**:
   ```bash
   /check-licenses --all --excluded
   ```

5. **Regular audits (weekly/monthly)**:
   - Schedule in CI/CD
   - Alert on new problematic licenses
   - Track license changes

## Exit Codes

- `0`: Success (all licenses approved)
- `1`: Problematic licenses found
- `2`: Configuration error
- `3`: Unable to detect licenses

## Performance Notes

- Go check: ~1-5 seconds
- npm check: ~2-10 seconds (with node_modules)
- Python check: ~1-3 seconds
- Full audit: ~10-20 seconds

## Related Skills

- `/validate-tools` - Validate development tools
- `/run-all-tests` - Run test suite
- `/manage-ci-workflows` - Manage CI/CD workflows
