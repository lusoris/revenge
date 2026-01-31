---
name: setup-codeql
description: Set up GitHub CodeQL security code scanning and analysis
argument-hint: "[--init|--view|--enable LANGUAGE] [--schedule cron|push] [--strict]"
disable-model-invocation: false
allowed-tools: Bash(python scripts/automation/github_security.py *)
---

# Setup CodeQL

Configure GitHub Advanced Security with CodeQL for automated security vulnerability scanning and code quality analysis.

## Usage

```
/setup-codeql --init                           # Initialize CodeQL scanning
/setup-codeql --view                           # View CodeQL configuration
/setup-codeql --enable go                      # Enable Go analysis
/setup-codeql --enable javascript --schedule push  # Enable JS on push
/setup-codeql --view --strict                  # View strict mode settings
```

## Arguments

- `$0`: Action (--init to create, --view to display, --enable for languages)
- `$1+`: Options (--schedule for scan timing, --strict for strict mode)

## Supported Languages

| Language | Database | Risk Level |
|----------|----------|-----------|
| Go | go | High (backend) |
| JavaScript | javascript | Medium (frontend) |
| Python | python | Medium (scripts) |
| YAML | yaml | Low (config) |

## Scan Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| Schedule | Weekly | When to scan (push or schedule) |
| Trigger | Push + schedule | Scan on code push + weekly |
| Upload results | Yes | Send to GitHub Security tab |
| Severity filter | all | Report all severities |

## Prerequisites

- GitHub Advanced Security enabled
- Repository admin access
- CodeQL CLI installed (optional for local scanning)
- Python 3.10+ for automation script

## Task

Set up CodeQL for continuous security scanning.

### Step 1: Check CodeQL Status

**View current CodeQL setup**:
```bash
python scripts/automation/github_security.py --view-codeql
```

**Shows**:
- Enabled languages
- Scan schedule
- Alert settings
- Historical scans

### Step 2: Initialize CodeQL

**Set up default configuration**:
```bash
python scripts/automation/github_security.py --configure-codeql --init
```

**Creates**:
- `codeql-analysis.yml` workflow
- Scans Go, JavaScript, Python
- Runs on push and weekly schedule
- Uploads results to GitHub

### Step 3: Enable Additional Languages

**Enable Go analysis**:
```bash
python scripts/automation/github_security.py --configure-codeql --enable go
```

**Enable JavaScript analysis**:
```bash
python scripts/automation/github_security.py --configure-codeql --enable javascript
```

**Enable Python analysis**:
```bash
python scripts/automation/github_security.py --configure-codeql --enable python
```

### Step 4: Configure Scan Schedule

**Run on every push**:
```bash
python scripts/automation/github_security.py --configure-codeql --schedule push
```

**Run on schedule (weekly)**:
```bash
python scripts/automation/github_security.py --configure-codeql --schedule weekly
```

### Step 5: Enable Strict Mode

**Enhanced analysis**:
```bash
python scripts/automation/github_security.py --configure-codeql --strict
```

**Strict mode includes**:
- Security query suites
- Quality query suites
- Custom queries
- Stricter alerts

## CodeQL Workflow

**File**: `.github/workflows/codeql-analysis.yml`

**Configuration**:
```yaml
name: CodeQL Analysis

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]
  schedule:
    - cron: '0 0 * * 0'  # Weekly

jobs:
  analyze:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        language: [go, javascript, python]
    steps:
      - uses: actions/checkout@v4

      - name: Initialize CodeQL
        uses: github/codeql-action/init@v2
        with:
          languages: ${{ matrix.language }}

      - name: Autobuild
        uses: github/codeql-action/autobuild@v2

      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@v2
```

## Analyzing Results

### Security Alerts

**Location**: Settings → Security → Code scanning

**Shows**:
- Vulnerability type
- Severity (critical, high, medium, low)
- File and line number
- CWE ID (Common Weakness Enumeration)
- Remediation advice

### Alert Severities

| Severity | Impact | Action |
|----------|--------|--------|
| Critical | Immediate security risk | Fix immediately |
| High | Significant risk | Fix ASAP |
| Medium | Moderate risk | Fix in next sprint |
| Low | Minor issue | Fix eventually |

### Common Vulnerabilities Found

**Go**:
- SQL injection
- Command injection
- Insecure cryptography
- Race conditions
- Null pointer dereferences

**JavaScript**:
- XSS (Cross-site scripting)
- CSRF (Cross-site request forgery)
- Unsafe DOM methods
- Missing input validation

**Python**:
- SQL injection
- Command injection
- Unsafe deserialization
- Missing input validation

## Examples

**View CodeQL configuration**:
```bash
/setup-codeql --view
```

**Initialize CodeQL scanning**:
```bash
/setup-codeql --init
```

**Enable Go security scanning**:
```bash
/setup-codeql --enable go
```

**Set up with weekly schedule**:
```bash
/setup-codeql --init --schedule weekly
```

**Enable strict mode analysis**:
```bash
/setup-codeql --init --strict
```

## Managing Alerts

### View Alerts

**In GitHub UI**:
1. Repository → Settings → Security
2. Code scanning alerts
3. Filter by severity, language, status

**Via CLI**:
```bash
gh secret-scanning alerts
gh code-scanning alerts
```

### Close Alerts

**Close with reason**:
1. Click alert
2. Select reason:
   - Fixed (code change)
   - False positive
   - Won't fix (accepted risk)
3. Save

**Programmatically**:
```bash
gh api repos/:owner/:repo/code-scanning/alerts/:alert-number \
  -X PATCH \
  -F state=dismissed \
  -F dismissed_reason=false_positive
```

### Dismissing False Positives

**Mark as false positive**:
1. Review alert details
2. Click "Dismiss alert"
3. Select "False positive"
4. Add comment explaining why
5. Save

**Reopening dismissed alert**:
1. Click "Show dismissed alerts"
2. Click alert
3. Click "Reopen alert"

## Integration with Branch Protection

**Configure to block merges on critical issues**:

1. Settings → Branches
2. Edit branch protection rule
3. Add status check: "CodeQL / Go analysis"
4. Make required

**Status checks**:
- PR must pass CodeQL scans
- No critical issues allowed
- Fail if new vulnerability introduced

## Local Scanning

**Install CodeQL CLI**:
```bash
gh codeql
# Or download from GitHub releases
```

**Run local scan**:
```bash
codeql database create /tmp/codeql-db --language=go
codeql database analyze /tmp/codeql-db go-security-and-quality.qls
```

**Compare local vs server**:
- Local: Full control, offline
- Server: Integrated with GitHub, automatic

## Customizing Queries

**Create custom query**:
```yaml
# .github/codeql/custom-queries.yml

name: My Custom Queries
version: 1.0

queries:
  - id: my-custom-security-check
    uses: .github/codeql-queries/my-check.ql
    severity: error
```

**Query file** (`.github/codeql-queries/my-check.ql`):
```ql
import go

from FunctionCall fc
where fc.getCallee().getName() = "eval"
select fc, "Use of eval is dangerous"
```

## Performance Tuning

### Large Repositories

**Parallel analysis**:
```yaml
strategy:
  matrix:
    language: [go, javascript, python]
  max-parallel: 3  # 3 languages in parallel
```

**Limit paths scanned**:
```yaml
with:
  paths: src/
  paths-ignore: vendor/, vendor-libs/
```

### Reducing Scan Time

1. **Exclude vendor directories**:
   ```yaml
   paths-ignore: |
     vendor/
     node_modules/
   ```

2. **Cache dependencies**:
   ```yaml
   - uses: actions/cache@v3
     with:
       path: ~/.cache/codeql
   ```

3. **Limit to push only** (not schedule):
   ```yaml
   on:
     push:
       branches: [main, develop]
   ```

## Troubleshooting

**"CodeQL analysis failing"**:
1. Check CodeQL logs in Actions
2. Verify language support
3. Check for build errors
4. Review error messages

**"No results being uploaded"**:
1. Verify GitHub Advanced Security enabled
2. Check workflow permissions
3. Ensure upload step present
4. Check for API rate limits

**"Too many false positives"**:
1. Review and dismiss false positives
2. Consider custom query tuning
3. Adjust severity thresholds
4. Exclude known safe patterns

**"High memory usage during scan"**:
1. Use path filtering to scan subsets
2. Run scans sequentially, not parallel
3. Increase runner memory if available
4. Split large repos into parts

**"Scan timeout"**:
1. Use path limiting
2. Run on larger runner
3. Run scans sequentially
4. Increase timeout in workflow

## Best Practices

1. **Scan all code**:
   - Go backend
   - JavaScript frontend
   - Python automation
   - All security-critical code

2. **Review alerts promptly**:
   - Within 24 hours for critical
   - Within 1 week for high
   - Track in project management

3. **Fix or dismiss**:
   - Fix security issues
   - Dismiss clearly documented false positives
   - Document "won't fix" decisions

4. **Use branch protection**:
   - Require CodeQL checks pass
   - Prevent merging with critical issues
   - Maintain security baseline

5. **Regular reviews**:
   - Weekly alert reviews
   - Monthly trend analysis
   - Quarterly security audits

## Monitoring and Reporting

### Alert Trends

**View over time**:
1. Settings → Security → Code scanning
2. See graph of alerts by date
3. Track improvement

### Export for Reporting

**Via API**:
```bash
gh api repos/:owner/:repo/code-scanning/alerts \
  --paginate \
  > alerts.json
```

**Create report**:
```bash
# Convert to CSV, filter by severity, etc.
jq '.[] | {rule: .rule.id, severity: .rule.severity, path: .most_recent_instance.location.path}' alerts.json > alerts.csv
```

## Related Skills

- `/configure-branch-protection` - Branch protection rules
- `/manage-ci-workflows` - CI/CD workflows
- `/run-linters` - Code quality linters
- `/check-health` - System health checks
