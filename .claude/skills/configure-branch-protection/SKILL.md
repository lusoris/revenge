---
name: configure-branch-protection
description: Configure GitHub branch protection rules for main and develop branches
argument-hint: "[--init|--view|--update] [--branch main|develop|all] [--strict]"
disable-model-invocation: false
allowed-tools: Bash(python scripts/automation/github_security.py *)
---

# Configure Branch Protection

Set up and manage branch protection rules to enforce quality gates and security policies on protected branches.

## Usage

```
/configure-branch-protection --init                 # Set up default protections
/configure-branch-protection --view                 # View current protections
/configure-branch-protection --view --branch main   # View main branch rules
/configure-branch-protection --update --strict      # Enable strict mode
/configure-branch-protection --view --branch all    # View all protected branches
```

## Arguments

- `$0`: Action (--init to create, --view to display, --update to modify)
- `$1+`: Options (--branch for specific branch, --strict for strict enforcement)

## Protected Branches

| Branch | Purpose | Protection Level |
|--------|---------|------------------|
| main | Production releases | Strict |
| develop | Development/staging | Standard |

## Protection Rules

### Standard Rules (develop)

- Require pull request reviews (1 approval)
- Dismiss stale PR reviews
- Require status checks to pass
- Include administrators in restrictions
- Allow force pushes: No
- Allow deletions: No

### Strict Rules (main)

- Require pull request reviews (2 approvals)
- Dismiss stale PR reviews
- Require status checks to pass
- Require up-to-date branches before merge
- Include administrators in restrictions
- Allow force pushes: No
- Allow deletions: No
- Require signed commits

## Prerequisites

- GitHub repository
- Admin access to repository settings
- GitHub CLI (`gh`) installed and authenticated

## Task

Configure branch protection rules for code quality and security.

### Step 1: View Current Rules

**View all protections**:
```bash
python scripts/automation/github_security.py --view-branch-protection
```

**View specific branch**:
```bash
python scripts/automation/github_security.py --view-branch-protection --branch main
python scripts/automation/github_security.py --view-branch-protection --branch develop
```

### Step 2: Initialize Default Protections

**Set up standard protections**:
```bash
python scripts/automation/github_security.py --configure-branch-protection --init
```

**Creates**:
- Main branch: Strict protection (2 reviews, signed commits)
- Develop branch: Standard protection (1 review)

### Step 3: Update to Strict Mode

**Enable strict protection**:
```bash
python scripts/automation/github_security.py --configure-branch-protection --strict
```

**Applies to both branches**:
- Requires 2 approvals
- Requires signed commits
- Requires up-to-date branches
- Blocks force pushes and deletions

### Step 4: Configure Specific Branch

**Update develop branch rules**:
```bash
python scripts/automation/github_security.py --configure-branch-protection --branch develop
```

**Update main branch rules**:
```bash
python scripts/automation/github_security.py --configure-branch-protection --branch main
```

## Branch Protection Details

### Main Branch Protection

**Settings**:
```
Require a pull request before merging:
  ✓ Require approvals: 2
  ✓ Dismiss stale pull request approvals when new commits are pushed
  ✓ Require review from Code Owners

Require status checks to pass before merging:
  ✓ Require branches to be up to date before merging
  Checks required:
    - CI / tests
    - CI / lint
    - CI / coverage

Require signed commits:
  ✓ Enabled

Who has access:
  ✓ Include administrators
  ✓ Restrict who can push
```

### Develop Branch Protection

**Settings**:
```
Require a pull request before merging:
  ✓ Require approvals: 1
  ✓ Dismiss stale pull request approvals when new commits are pushed
  ✓ Require review from Code Owners

Require status checks to pass before merging:
  ✓ Require branches to be up to date before merging
  Checks required:
    - CI / tests
    - CI / lint

Who has access:
  ✓ Include administrators
```

## Status Checks Required

### Go Tests

**Check**: `CI / tests`

**What it does**:
- Runs `go test ./...`
- Requires 80% coverage
- No race conditions allowed
- Must pass to merge

### Linting

**Check**: `CI / lint`

**What it does**:
- Runs golangci-lint
- Runs ruff (Python)
- Runs prettier (frontend)
- Must pass to merge

### Coverage

**Check**: `CI / coverage`

**What it does**:
- Generates coverage report
- Requires minimum 80%
- Must pass to merge

## Code Owners

**File**: `.github/CODEOWNERS`

**Reviewers by area**:
```
# Default
* @lusoris

# Documentation
/docs/ @lusoris
*.md @lusoris

# Design docs
/docs/dev/design/ @lusoris

# CI/CD
/.github/ @lusoris

# Database
/migrations/ @lusoris
/internal/infra/database/ @lusoris

# API
/internal/api/ @lusoris

# Security
/SECURITY.md @lusoris
/internal/api/middleware/auth*.go @lusoris
```

**How it works**:
- PR touches `/internal/api/` → Review from @lusoris
- PR touches `/docs/dev/design/` → Review from @lusoris
- PR touches auth middleware → Review from @lusoris

## Examples

**View current branch protections**:
```bash
/configure-branch-protection --view
```

**Set up default protections**:
```bash
/configure-branch-protection --init
```

**View main branch rules**:
```bash
/configure-branch-protection --view --branch main
```

**Enable strict mode**:
```bash
/configure-branch-protection --update --strict
```

## Managing Protected Branches

### Bypassing Protection (Emergency Only)

**Temporarily disable for hotfix**:
1. Go to Settings → Branches
2. Find branch rule
3. Click "Edit"
4. Temporarily uncheck rule
5. Merge hotfix
6. Re-enable immediately

**Should be rare** - use admin review instead

### Updating Protection Rules

**Add new required check**:
1. Settings → Branches → Branch protection rules
2. Edit rule for branch
3. Add status check under "Status checks"
4. Select check to require

**Change approval requirements**:
1. Edit branch rule
2. Update "Require approvals" number
3. Save changes

**Add Code Owner requirement**:
1. Edit branch rule
2. Check "Require review from Code Owners"
3. Ensure CODEOWNERS file updated
4. Save

## Workflow for PRs

### Before PR

1. **Create feature branch** from develop
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/my-feature
   ```

2. **Make changes and commit**
   ```bash
   git add .
   git commit -m "feat: description"
   ```

3. **Push to remote**
   ```bash
   git push origin feature/my-feature
   ```

### PR Created

1. **Branch protection checks trigger**:
   - CI tests run
   - Linting runs
   - Coverage checked

2. **Required reviews requested**:
   - Code Owners notified (if applicable)
   - Approvers assigned

3. **Status checks appear**:
   - Green check = passed
   - Red X = failed
   - Yellow dot = running

### Before Merge

1. **All checks must pass**:
   - ✓ CI tests
   - ✓ Linting
   - ✓ Coverage
   - ✓ Required approvals

2. **Branch must be up-to-date**:
   - Update branch button available
   - Click to merge main into feature
   - Re-run checks

3. **Required reviewers approve**:
   - At least 1 approval for develop
   - At least 2 approvals for main

4. **Merge PR**:
   - Squash and merge (recommended)
   - Create merge commit (alternative)
   - Rebase and merge (alternative)

## Troubleshooting

**"Branch protection rule not applying"**:
1. Verify rule is enabled
2. Check branch name pattern
3. Ensure rule saved successfully
4. Try creating new rule

**"Can't merge PR despite all checks passing"**:
1. Check for required approvals
2. Verify no "dismiss reviews" setting
3. Ensure stale review dismissal configured
4. Check if up-to-date branch required

**"Status check not appearing"**:
1. Check workflow runs
2. Verify workflow has correct triggers
3. Check workflow permissions
4. Add check to protection rule

**"Code Owner review not requested"**:
1. Verify CODEOWNERS file syntax
2. Check if "Require Code Owner review" enabled
3. Ensure users exist and have access
4. Check if PR touches CODEOWNERS file

**"Can't bypass protection (emergency)"**:
1. Admin only feature
2. Check your permissions
3. Go to Settings → Branches
4. Edit rule → Temporarily disable

## Tips

1. **Keep main highly protected**:
   - 2 approvals minimum
   - Signed commits required
   - Strict enforcement

2. **Develop can be less strict**:
   - 1 approval acceptable
   - Faster merges
   - Still maintains quality

3. **Use Code Owners**:
   - Automatic review requests
   - Domain-specific expertise
   - Reduces review time

4. **Run tests locally**:
   - Before pushing
   - Catch issues early
   - Reduce PR cycles
   ```bash
   go test ./...
   go fmt ./...
   golangci-lint run
   ```

5. **Keep branches updated**:
   - Merge main → develop regularly
   - Use "Update branch" in PR
   - Avoid merge conflicts

## Best Practices

1. **Require status checks**:
   - Tests must pass
   - Linting must pass
   - Coverage must meet threshold

2. **Require approvals**:
   - At least 1 for develop
   - At least 2 for main
   - From designated reviewers

3. **Enable stale review dismissal**:
   - New commits dismiss old reviews
   - Ensures fresh approval
   - Catches hidden changes

4. **Require up-to-date branches**:
   - No sneaking in code
   - All tests re-run
   - Ensures compatibility

5. **Sign commits** (main only):
   - Proves author identity
   - Increases security
   - Supports compliance

## Integration with Workflows

**Status checks come from**:
- `.github/workflows/ci.yml` - Tests and lint
- `.github/workflows/coverage.yml` - Coverage report
- Third-party services (Codecov, etc.)

**All must pass before merge**:
```
Status checks required:
✓ CI / tests
✓ CI / lint
✓ CI / coverage
```

## Related Skills

- `/manage-ci-workflows` - Manage CI/CD workflows
- `/run-all-tests` - Run test suite
- `/run-linters` - Run linters
- `/setup-codeql` - Code security scanning
