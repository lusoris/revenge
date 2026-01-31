---
name: configure-release-please
description: Configure Release Please for automated semantic versioning and releases
argument-hint: "[--init|--view|--update] [--release-type major|minor|patch] [--monorepo]"
disable-model-invocation: false
allowed-tools: Bash(cat .github/release-please-config.json *, cat .github/workflows/release-please.yml *)
---

# Configure Release Please

Set up and manage Release Please configuration for automated semantic versioning, changelog generation, and GitHub releases.

## Usage

```
/configure-release-please --view                    # View current configuration
/configure-release-please --init                    # Initialize Release Please config
/configure-release-please --view --release-type     # Show release strategy
/configure-release-please --update --release-type patch  # Update release strategy
```

## Arguments

- `$0`: Action (--init to create, --view to display, --update to modify)
- `$1+`: Options (--release-type, --monorepo flag)

## Release Please Features

| Feature | Status |
|---------|--------|
| Semantic Versioning | Enabled |
| Conventional Commits | Required |
| Changelog Generation | Automatic |
| Release PR Creation | Automatic |
| GitHub Releases | Automatic |
| Git Tags | Automatic |

## Prerequisites

- GitHub repository with Release Please workflow enabled
- `.github/release-please-config.json` configuration file
- `.github/workflows/release-please.yml` workflow
- Commits following conventional commit format

## Task

View and configure Release Please for automated versioning and releases.

### Step 1: View Current Configuration

**View Release Please config**:
```bash
cat .github/release-please-config.json
```

**View Release Please workflow**:
```bash
cat .github/workflows/release-please.yml
```

### Step 2: Understand Current Settings

**Configuration file structure**:
```json
{
  "release-type": "go",
  "bump-minor-pre-major": true,
  "bump-patch-for-minor-pre-major": true,
  "changelog-path": "CHANGELOG.md",
  "include-vfile-updates": true,
  "packages": {
    ".": {
      "changelog-path": "CHANGELOG.md"
    }
  }
}
```

**Workflow triggers**:
- Runs on push to `develop` branch
- Analyzes conventional commits
- Creates release PR if version bump needed

### Step 3: View Configuration

**View entire configuration**:
```bash
/configure-release-please --view
```

**Show release strategy**:
```bash
/configure-release-please --view --release-type
```

### Step 4: Understanding Release Detection

**Release Please analyzes commits**:

| Commit Type | Impact | Example |
|-------------|--------|---------|
| `fix:` | Patch version | `fix: resolve login bug` |
| `feat:` | Minor version | `feat: add dark mode` |
| `BREAKING CHANGE:` | Major version | `feat: redesign API\n\nBREAKING CHANGE: ...` |

**Example commit impacts**:
- `fix: bug fix` → 1.0.0 → 1.0.1
- `feat: new feature` → 1.0.0 → 1.1.0
- `feat: refactor API\n\nBREAKING CHANGE: ...` → 1.0.0 → 2.0.0

## Current Configuration

**File**: `.github/release-please-config.json`

```json
{
  "release-type": "go",
  "bump-minor-pre-major": true,
  "bump-patch-for-minor-pre-major": true,
  "changelog-path": "CHANGELOG.md",
  "include-vfile-updates": true,
  "packages": {
    ".": {
      "changelog-path": "CHANGELOG.md"
    }
  }
}
```

**Configuration options**:

- `release-type`: "go" (for Go projects)
- `bump-minor-pre-major`: Allow minor bumps before major
- `bump-patch-for-minor-pre-major`: Allow patch bumps before major
- `changelog-path`: Where to write changelog
- `include-vfile-updates`: Update version files automatically

## Workflow Integration

**Workflow**: `.github/workflows/release-please.yml`

**Triggers**:
- `push` to `develop` branch
- Manual trigger via `workflow_dispatch`

**Actions**:
1. Run Release Please action
2. Create or update release PR
3. On merge, create GitHub release
4. Tag commit with version
5. Push to remote

## Release Workflow

### Automatic Flow

1. **Developer pushes commit** to `develop`:
   ```bash
   git commit -m "feat: add new feature"
   git push origin develop
   ```

2. **Release Please checks commit**:
   - Analyzes conventional commit message
   - Calculates version bump (minor in this case)
   - Current version: 1.0.0 → 1.1.0

3. **Release PR created** (if needed):
   - Title: "chore: release 1.1.0"
   - Body: Generated changelog
   - Commits: Updated package.json, CHANGELOG.md, etc.

4. **Developers review PR**:
   - Check changelog looks good
   - Verify version correct
   - Approve and merge

5. **Automated Release**:
   - Merge triggers workflow
   - GitHub release created
   - Git tag created (v1.1.0)
   - Tagged commit pushed

### Manual Workflow

**If Release Please skips a commit**:
```bash
# Force release PR creation
gh workflow run release-please.yml
```

## Examples

**View current Release Please configuration**:
```bash
/configure-release-please --view
```

**Check release strategy**:
```bash
/configure-release-please --view --release-type
```

## Conventional Commit Guide

**Commit message format**:
```
type(scope): subject

body

footer
```

**Types that trigger releases**:
- `feat:` → Minor version bump
- `fix:` → Patch version bump
- `BREAKING CHANGE:` → Major version bump

**Examples**:

Patch release:
```
fix: resolve database connection issue
```

Minor release:
```
feat: add user authentication system
```

Major release:
```
feat: redesign API endpoints

BREAKING CHANGE: User endpoints moved from /api/v1/users to /api/v2/users
```

## Release Please PR Details

**What Release Please PR includes**:

1. **CHANGELOG.md update**:
   - Lists all commits since last release
   - Grouped by type (Features, Bug Fixes, etc.)
   - Formatted markdown

2. **Version bumps**:
   - `package.json` (for frontend)
   - `go.mod` version comments (if applicable)
   - `VERSION` file (if exists)

3. **Git information**:
   - Commit message with version
   - Reference to release commit

## Managing Releases

### Review Release PR

```bash
# List open Release Please PRs
gh pr list --search "is:open author:release-please-bot"

# View specific PR
gh pr view <pr-number>

# Review changes
git checkout -b release-<version>
git pull origin pull/<pr-number>/head
```

### Approve Release PR

```bash
# Approve PR
gh pr review <pr-number> --approve

# Merge PR (triggers release)
gh pr merge <pr-number> --squash
```

### Check Release History

```bash
# List releases
gh release list

# View specific release
gh release view v1.0.0

# Check git tags
git tag -l
git log v1.0.0^..v1.0.0
```

## Troubleshooting

**"Release Please not creating PR"**:
1. Check commits use conventional format
2. Verify workflow file exists: `.github/workflows/release-please.yml`
3. Review workflow logs: Actions tab
4. Ensure branch is `develop`
5. Check for existing open release PR

**"Wrong version calculated"**:
1. Review commit messages
2. Ensure BREAKING CHANGE is in body (not subject)
3. Check for typos in commit type
4. Verify config file is correct

**"Changelog missing some commits"**:
1. Check commit format: `type(scope): message`
2. Verify commits are on release branch
3. Look for commits without conventional format
4. Review Release Please logs for filtering

**"Release PR has merge conflicts"**:
1. Merge conflicts in CHANGELOG.md are normal
2. Release Please may need to rebase
3. Comment on PR requesting rebase
4. Or manually fix and push to PR branch

**"Can't merge release PR"**:
1. Check branch protection rules
2. Ensure all CI checks pass
3. Verify required approvals met
4. Check for status check failures

## Best Practices

1. **Use conventional commits consistently**:
   ```bash
   # Good
   feat: add user authentication
   fix: resolve login bug

   # Bad
   added auth
   fixed bug
   ```

2. **Major versions need clear communication**:
   ```bash
   feat: redesign API

   BREAKING CHANGE: All endpoints now require API key authentication.
   Migration: See MIGRATION.md for upgrade guide.
   ```

3. **Review release PRs carefully**:
   - Verify changelog is accurate
   - Check version number makes sense
   - Review generated changes

4. **Don't skip releases for minor commits**:
   - Patch releases are good
   - Keeps version history clean
   - Users can track changes

5. **Test before releasing**:
   - CI/CD runs on release PR
   - Wait for all checks before merge
   - Manual testing for major releases

## Customization

### Change Release Type

**If moving from pure Go to monorepo**:

```json
{
  "release-type": "node",
  "monorepo-type": "rust",
  "packages": {
    ".": {
      "changelog-path": "CHANGELOG.md",
      "release-type": "go"
    }
  }
}
```

### Customize Changelog Path

**Default**: `CHANGELOG.md` in root

**To change**:
```json
{
  "changelog-path": "docs/RELEASES.md"
}
```

### Add Version File Updates

**To auto-update VERSION file**:
```json
{
  "include-vfile-updates": true,
  "version-file": "VERSION"
}
```

## Integration with Version Commands

**Check current version**:
```bash
# From git tags
git describe --tags --abbrev=0

# From Release Please config
jq -r '.packages["."].version' .github/release-please-config.json
```

## Monitoring Releases

**Check for pending releases**:
```bash
gh pr list --search "is:open author:release-please-bot"
```

**View recent releases**:
```bash
gh release list --limit 10
```

**Check release notes**:
```bash
gh release view v1.0.0
```

## Tips

1. **Automatic releases free up time**:
   - No manual version bumping
   - Consistent versioning
   - Automated changelog

2. **Review before release**:
   ```bash
   /configure-release-please --view
   ```

3. **Monitor release PRs**:
   ```bash
   gh pr list --search "author:release-please-bot"
   ```

4. **Use git tags for reference**:
   ```bash
   git tag
   git checkout v1.0.0
   ```

## Related Skills

- `/update-dependencies` - Dependency management
- `/manage-ci-workflows` - CI/CD workflows
- `/validate-tools` - Tool validation
- `/run-all-tests` - Test suite before release
