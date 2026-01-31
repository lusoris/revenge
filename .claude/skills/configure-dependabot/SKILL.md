---
name: configure-dependabot
description: Configure GitHub Dependabot for automated dependency updates
argument-hint: "[--init|--update|--view] [--ecosystem gomod|docker|npm|pip] [--schedule weekly|daily|monthly]"
disable-model-invocation: false
allowed-tools: Bash(cat .github/dependabot.yml *, python scripts/automation/update_dependencies.py *)
---

# Configure Dependabot

Set up and manage GitHub Dependabot configuration for automated dependency updates across all package managers.

## Usage

```
/configure-dependabot --view                    # View current Dependabot configuration
/configure-dependabot --view --ecosystem gomod  # View Go configuration only
/configure-dependabot --update --ecosystem npm --schedule daily   # Update npm schedule to daily
/configure-dependabot --init                    # Initialize default Dependabot config
```

## Arguments

- `$0`: Action (--init to create, --view to display, --update to modify)
- `$1+`: Options (--ecosystem for specific manager, --schedule for update frequency)

## Supported Package Managers

| Ecosystem | Files | Schedule | Reviewers |
|-----------|-------|----------|-----------|
| gomod | `go.mod` | Weekly (Monday 6am) | @lusoris |
| docker | `Dockerfile` | Weekly (Monday 6am) | @lusoris |
| npm | `frontend/package.json` | Weekly (Monday 6am) | @lusoris |
| pip | `requirements.txt` | Weekly (Monday 6am) | @lusoris |
| github-actions | `.github/workflows` | Weekly (Monday 6am) | @lusoris |

## Configuration File

**Location**: `.github/dependabot.yml`

**Current config includes**:
- Go modules (gomod)
- Docker images
- npm packages (frontend)
- Python (pip)
- GitHub Actions

## Prerequisites

- GitHub repository with Dependabot enabled
- Admin access to repository settings
- `.github/dependabot.yml` file

## Task

View and configure Dependabot settings for automated dependency management.

### Step 1: View Current Configuration

**View entire Dependabot config**:
```bash
cat .github/dependabot.yml
```

**View specific ecosystem**:
```bash
# Go dependencies
grep -A 10 'package-ecosystem: "gomod"' .github/dependabot.yml

# npm dependencies
grep -A 10 'package-ecosystem: "npm"' .github/dependabot.yml

# Docker images
grep -A 10 'package-ecosystem: "docker"' .github/dependabot.yml
```

### Step 2: Understand Current Settings

**Default schedule** (all ecosystems):
- Interval: Weekly
- Day: Monday
- Time: 06:00 UTC
- Timezone: UTC (or "Europe/Berlin" if configured)

**Default limits**:
- Go: 10 open PRs max
- npm: 10 open PRs max
- Docker: 5 open PRs max
- pip: 5 open PRs max
- github-actions: 5 open PRs max

**Default labels**:
- All PRs get `dependencies` label
- Ecosystem-specific labels (go, docker, npm, python, github-actions)

### Step 3: Modify Schedule (If Needed)

**To update frequency**, edit `.github/dependabot.yml`:

```yaml
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "daily"        # Change to: daily, weekly, monthly
      day: "monday"            # day: (for weekly): monday, tuesday, etc.
      time: "06:00"            # time: HH:MM in UTC
```

**Schedule options**:
- `interval: "daily"` - Check daily
- `interval: "weekly"` - Check weekly (specify day)
- `interval: "monthly"` - Check monthly

### Step 4: Modify PR Limits (If Needed)

**To change max open PRs**, edit section:

```yaml
updates:
  - package-ecosystem: "gomod"
    open-pull-requests-limit: 20    # Change from 10 to 20
```

### Step 5: Update Reviewers (If Needed)

**To add reviewers**, modify:

```yaml
updates:
  - package-ecosystem: "gomod"
    reviewers:
      - "@lusoris"           # Existing
      - "@newuser"           # Add new
```

### Step 6: Configure Auto-Merge (Optional)

**Enable auto-merge for safe updates**:

```yaml
updates:
  - package-ecosystem: "gomod"
    auto-merge-options:
      update-types:
        - "minor"
        - "patch"
      automerge-type: "squash"
      automerge-options:
        commit-message-prefix: "chore(deps):"
```

## View Configuration

**Current Dependabot config**:
```bash
/configure-dependabot --view
```

**Go ecosystem only**:
```bash
/configure-dependabot --view --ecosystem gomod
```

**npm ecosystem**:
```bash
/configure-dependabot --view --ecosystem npm
```

## Examples

**View current settings**:
```bash
/configure-dependabot --view
```

**View Go dependency configuration**:
```bash
/configure-dependabot --view --ecosystem gomod
```

**Understand npm settings**:
```bash
/configure-dependabot --view --ecosystem npm
```

**View all configured ecosystems**:
```bash
/configure-dependabot --view
```

## Dependabot Configuration Details

### Go Modules Configuration

```yaml
- package-ecosystem: "gomod"
  directory: "/"
  schedule:
    interval: "weekly"
    day: "monday"
    time: "06:00"
  open-pull-requests-limit: 10
  reviewers:
    - "lusoris"
  labels:
    - "dependencies"
    - "go"
  commit-message:
    prefix: "chore(deps)"
    include: "scope"
  groups:
    go-dependencies:
      patterns:
        - "*"
      update-types:
        - "minor"
        - "patch"
```

**What it does**:
- Checks Go modules weekly
- Groups minor/patch updates together
- Auto-creates PR with `dependencies` and `go` labels
- Assigns to @lusoris for review
- Commit message starts with `chore(deps):`

### npm Configuration

```yaml
- package-ecosystem: "npm"
  directory: "/frontend"
  schedule:
    interval: "weekly"
    day: "monday"
    time: "06:00"
  open-pull-requests-limit: 10
  reviewers:
    - "lusoris"
  labels:
    - "dependencies"
    - "frontend"
  commit-message:
    prefix: "chore(deps)"
    include: "scope"
```

### Docker Configuration

```yaml
- package-ecosystem: "docker"
  directory: "/"
  schedule:
    interval: "weekly"
    day: "monday"
    time: "06:00"
  open-pull-requests-limit: 5
  reviewers:
    - "lusoris"
  labels:
    - "dependencies"
    - "docker"
```

### Python (pip) Configuration

```yaml
- package-ecosystem: "pip"
  directory: "/"
  schedule:
    interval: "weekly"
    day: "monday"
    time: "06:00"
  open-pull-requests-limit: 5
  reviewers:
    - "lusoris"
  labels:
    - "dependencies"
    - "python"
```

## Managing Dependabot PRs

### Workflow for Dependabot PRs

1. **Dependabot creates PR** (Monday 6am)
2. **CI/CD runs tests** (automatic)
3. **Review PR** (human review)
4. **Merge or close** (manual or auto-merge if configured)

### Merging Dependabot PRs

```bash
# Option 1: Manual merge
gh pr merge <pr-number> --squash

# Option 2: Auto-merge (if configured)
gh pr merge <pr-number> --auto --squash

# Option 3: Command line git
git checkout dependabot/go_modules/...
git merge main
git push
```

### Dismissing Dependency Alerts

**If you want to ignore an update**:
1. Go to Security → Dependabot alerts
2. Click "Dismiss"
3. Select reason for dismissal
4. Optionally add comment

## Troubleshooting

**"Dependabot not creating PRs"**:
1. Verify `.github/dependabot.yml` is valid YAML
2. Check syntax: `yamllint .github/dependabot.yml`
3. Ensure Dependabot is enabled in repository settings
4. Check for branch protection rules blocking PRs
5. Review Dependabot logs: Settings → Security → Dependabot

**"Schedule not being followed"**:
1. Dependabot may batch updates
2. May wait for PR queue to clear
3. May be affected by rate limits
4. Check Dependabot logs for details

**"Too many open PRs"**:
1. Check `open-pull-requests-limit` in config
2. Merge or close some existing PRs
3. Adjust limit in `.github/dependabot.yml`
4. Dependabot will resume after limit drops

**"PR keeps getting updated"**:
- Normal behavior if new versions released
- Dependabot updates PR as new versions available
- Review and merge when ready

**"Auto-merge not working"**:
1. Verify `auto-merge-options` configured
2. Check branch protection rules
3. Ensure tests pass
4. Verify PR labels match conditions

## Best Practices

1. **Keep Dependabot enabled**:
   - Ensures dependencies stay current
   - Reduces security vulnerability window
   - Easier updates than major version jumps

2. **Review regularly**:
   ```bash
   # Check for open Dependabot PRs
   gh pr list --search "is:open author:dependabot"
   ```

3. **Test before merging**:
   - Always let CI/CD complete
   - Review changes before merge
   - Test breaking changes manually

4. **Batch updates**:
   - Use grouping for related packages
   - Reduces PR volume
   - Easier to track

5. **Set up auto-merge cautiously**:
   - Only for truly safe updates (patch/minor)
   - Require tests to pass
   - Add extra review layers

## Monitoring Dependabot Activity

**View recent Dependabot activity**:
```bash
gh pr list --search "is:open author:dependabot"
gh pr list --search "is:closed author:dependabot" --limit 10
```

**Check specific ecosystem PRs**:
```bash
gh pr list --search "is:open author:dependabot label:go"
gh pr list --search "is:open author:dependabot label:frontend"
```

## Configuration Validation

**Validate YAML syntax**:
```bash
yamllint .github/dependabot.yml
```

**Test configuration locally**:
```bash
# Simulate what Dependabot would do
python scripts/automation/update_dependencies.py --all --dry-run
```

## Tips

1. **Review changes monthly**:
   ```bash
   /configure-dependabot --view
   ```

2. **Keep grouped updates for efficiency**:
   - Less PR fatigue
   - Easier to test related changes
   - Cleaner commit history

3. **Adjust schedule for your timezone**:
   - Edit `time:` field (UTC)
   - Or add `timezone: "Europe/Berlin"` field

4. **Use labels for filtering**:
   ```bash
   gh pr list --search "label:dependencies label:go"
   ```

## Related Skills

- `/update-dependencies` - Manual dependency updates
- `/check-licenses` - Verify license compliance
- `/manage-ci-workflows` - CI/CD workflow management
- `/validate-tools` - Tool version validation
