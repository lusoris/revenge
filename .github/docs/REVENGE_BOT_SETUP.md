# Revenge Bot Setup Guide

> Documentation for setting up the `revenge-bot` GitHub account for automated repository operations

**Last Updated**: 2026-01-31
**Bot Account**: `revenge-bot` (to be created)
**Purpose**: Automated commits, PR management, issue triage, and GitHub Actions operations

---

## Overview

The revenge-bot is a dedicated GitHub account used for automated operations in the Revenge repository. Using a bot account instead of personal accounts provides better audit trails, clearer attribution, and separates automated actions from human contributions.

---

## Prerequisites

- Admin access to the Revenge GitHub repository
- Email address for bot account (e.g., `revenge-bot@example.com`)
- GitHub organization (optional but recommended)

---

## Step 1: Create GitHub Account

### 1.1 Account Creation

1. **Log out** of all GitHub sessions
2. Navigate to [github.com/signup](https://github.com/signup)
3. Use bot email: `revenge-bot@example.com`
4. Username: `revenge-bot`
5. Complete signup process
6. Verify email address

### 1.2 Profile Configuration

Set up the bot profile:

```
Username: revenge-bot
Name: Revenge Bot
Bio: 🤖 Automated assistant for the Revenge media server project
Company: @revenge (or your org)
Location: (leave empty)
Website: https://github.com/yourusername/revenge
```

### 1.3 Profile Picture

Create a simple robot avatar or use a service like:
- [RoboHash](https://robohash.org/) - Generate robot avatars
- Use project logo with "BOT" overlay
- Simple emoji-based avatar (🤖)

---

## Step 2: Repository Access

### 2.1 Add as Collaborator

**Option A: Direct Collaborator**
```bash
# Using gh CLI
gh api repos/OWNER/REPO/collaborators/revenge-bot -X PUT \
  -f permission=write
```

**Option B: Organization Member**
1. Invite bot to organization
2. Create "Bots" team
3. Add `revenge-bot` to "Bots" team
4. Grant team write access to repository

### 2.2 Accept Invitation

1. Log in as `revenge-bot`
2. Navigate to repository
3. Accept collaboration invitation

---

## Step 3: Generate Tokens

### 3.1 Personal Access Token (Classic)

Create a classic PAT with these scopes:

```
✓ repo (all)
  ✓ repo:status
  ✓ repo_deployment
  ✓ public_repo
  ✓ repo:invite
  ✓ security_events

✓ workflow
✓ write:packages
  ✓ read:packages

✓ admin:org
  ✓ write:org
  ✓ read:org

✓ admin:repo_hook
  ✓ write:repo_hook
  ✓ read:repo_hook
```

**Steps**:
1. Go to Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Note: `revenge-bot repository automation`
4. Select scopes above
5. Set expiration: 1 year (or no expiration if trusted)
6. Generate token
7. **SAVE TOKEN SECURELY** - you won't see it again

### 3.2 Fine-Grained Personal Access Token

Alternative to classic token (recommended for better security):

```
Repository access: Only select repositories (revenge)

Repository permissions:
- Actions: Read and write
- Checks: Read and write
- Contents: Read and write
- Deployments: Read and write
- Issues: Read and write
- Metadata: Read (automatically granted)
- Pull requests: Read and write
- Secrets: Read and write
- Workflows: Read and write
```

---

## Step 4: Configure Repository Secrets

Add bot token to repository secrets:

```bash
# Using gh CLI
gh secret set BOT_TOKEN --body "ghp_YOUR_TOKEN_HERE"
gh secret set REVENGE_BOT_TOKEN --body "ghp_YOUR_TOKEN_HERE"
```

Or via GitHub UI:
1. Go to Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Name: `BOT_TOKEN`
4. Secret: Paste the token
5. Add secret

---

## Step 5: Configure Git for Bot

### 5.1 Git Config for Workflows

In GitHub Actions workflows, configure git to use bot:

```yaml
- name: Configure git
  run: |
    git config --global user.name "revenge-bot"
    git config --global user.email "revenge-bot@users.noreply.github.com"
```

### 5.2 GitHub Actions Bot User

For commits made by GitHub Actions (without token):

```yaml
- name: Configure git (GitHub Actions bot)
  run: |
    git config --global user.name "github-actions[bot]"
    git config --global user.email "github-actions[bot]@users.noreply.github.com"
```

---

## Step 6: Update Workflows

### 6.1 Use Bot Token in Workflows

Update `.github/workflows/*.yml` to use bot token:

```yaml
name: Automated Updates

on:
  schedule:
    - cron: '0 2 * * 1'  # Weekly Monday 2am
  workflow_dispatch:

jobs:
  update-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          token: ${{ secrets.BOT_TOKEN }}

      - name: Configure git
        run: |
          git config user.name "revenge-bot"
          git config user.email "revenge-bot@users.noreply.github.com"

      - name: Update documentation
        run: python scripts/fetch-sources.py

      - name: Commit and push
        run: |
          git add .
          git commit -m "docs: auto-update from sources" || exit 0
          git push
```

### 6.2 Create PR with Bot

```yaml
      - name: Create Pull Request
        uses: peter-evans/create-pull-request@v6
        with:
          token: ${{ secrets.BOT_TOKEN }}
          commit-message: "docs: automated source updates"
          title: "docs: weekly source documentation update"
          body: |
            Automated update of external documentation sources.

            Changes:
            - Updated sources from external APIs
            - Regenerated documentation indexes

            🤖 Generated by revenge-bot
          branch: auto-update-sources
          author: "revenge-bot <revenge-bot@users.noreply.github.com>"
          committer: "revenge-bot <revenge-bot@users.noreply.github.com>"
```

---

## Step 7: Configure CODEOWNERS

Exclude bot from review requirements:

```
# CODEOWNERS

# Default owner
* @kilian

# Backend code
internal/** @kilian

# Frontend code
frontend/** @kilian

# Documentation - bot can auto-update
docs/dev/sources/** @revenge-bot
```

---

## Step 8: Branch Protection Rules

Configure branch protection to allow bot pushes:

**Settings → Branches → Branch protection rules → develop**

```
✓ Require a pull request before merging
  □ Require approvals: 1
  ✓ Dismiss stale reviews
  ✓ Require review from Code Owners
  ✓ Allow specified actors to bypass pull request requirements:
    + Add: revenge-bot

✓ Require status checks to pass
  ✓ Require branches to be up to date
  Status checks:
    - test
    - lint
    - build

✓ Require signed commits (optional)

✓ Include administrators

Allowed to push:
  + Add: revenge-bot
```

---

## Step 9: Test Bot Operations

### 9.1 Test Commit

```bash
# Clone as bot (locally)
git clone https://ghp_TOKEN@github.com/OWNER/revenge.git
cd revenge
git config user.name "revenge-bot"
git config user.email "revenge-bot@users.noreply.github.com"

# Make test commit
echo "# Bot Test" >> test-bot.md
git add test-bot.md
git commit -m "test: bot commit verification"
git push

# Clean up
git rm test-bot.md
git commit -m "test: remove bot test file"
git push
```

### 9.2 Test Workflow

Trigger a workflow manually:

```bash
gh workflow run "Source Refresh" --ref develop
gh workflow view "Source Refresh"
```

### 9.3 Verify Commits

Check that commits are attributed to bot:

```bash
git log --author="revenge-bot" --oneline | head -5
```

Should show:
```
abc1234 docs: auto-update from sources
def5678 docs: regenerate indexes
```

---

## Step 10: Security Best Practices

### 10.1 Token Rotation

**Schedule**: Rotate bot token every 6-12 months

**Process**:
1. Generate new token as `revenge-bot`
2. Update `BOT_TOKEN` secret in repository
3. Revoke old token after confirming new one works
4. Document rotation in team notes

### 10.2 Token Storage

- **Never** commit tokens to repository
- Store tokens in GitHub Secrets or secure vault
- Use environment variables for local testing
- Revoke immediately if leaked

### 10.3 Audit Log

Regularly review bot actions:

```bash
# Check recent bot commits
git log --author="revenge-bot" --since="1 month ago"

# Check workflow runs
gh run list --user revenge-bot
```

### 10.4 Least Privilege

Only grant permissions the bot actually needs:
- Don't use admin tokens unless required
- Scope tokens to specific repositories
- Use fine-grained tokens when possible

---

## Step 11: Automation Use Cases

### 11.1 Source Documentation Updates

**Workflow**: `.github/workflows/source-refresh.yml`

```yaml
name: Source Refresh
on:
  schedule:
    - cron: '0 2 * * 1'  # Weekly Monday 2am
  workflow_dispatch:

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          token: ${{ secrets.BOT_TOKEN }}

      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.12'

      - name: Install dependencies
        run: pip install -r scripts/requirements.txt

      - name: Fetch sources
        run: python scripts/fetch-sources.py

      - name: Commit changes
        run: |
          git config user.name "revenge-bot"
          git config user.email "revenge-bot@users.noreply.github.com"
          git add docs/dev/sources/
          git commit -m "docs: auto-refresh external sources" || exit 0
          git push
```

### 11.2 Dependency Updates

**Workflow**: `.github/workflows/dependency-update.yml`

Auto-create PRs for dependency updates.

### 11.3 Issue Triage

**Workflow**: `.github/workflows/issue-triage.yml`

Auto-label and assign issues based on content.

---

## Troubleshooting

### Bot Token Not Working

**Symptoms**: 403 errors, "Resource not accessible by integration"

**Solutions**:
1. Verify token has correct scopes
2. Check token hasn't expired
3. Confirm bot has repository access
4. Try regenerating token

```bash
# Test token
curl -H "Authorization: token ghp_YOUR_TOKEN" \
  https://api.github.com/user
```

### Commits Not Attributed to Bot

**Symptoms**: Commits show as different user

**Solution**:
```bash
# Verify git config in workflow
- run: |
    git config --list | grep user
    # Should show: user.name=revenge-bot
```

### Branch Protection Blocks Bot

**Symptoms**: Bot can't push to protected branches

**Solution**:
1. Settings → Branches → Branch protection
2. Add `revenge-bot` to "Allow specified actors to bypass"
3. Or create PR instead of direct push

### Workflow Permissions Error

**Symptoms**: "Resource not accessible", "Forbidden"

**Solution**:

Update workflow permissions:

```yaml
permissions:
  contents: write
  pull-requests: write
  issues: write
```

---

## Monitoring

### Daily Checks

```bash
# Check recent bot activity
gh api /repos/OWNER/REPO/commits --jq '.[] | select(.commit.author.name=="revenge-bot") | .commit.message' | head -5

# Check workflow runs
gh run list --workflow=source-refresh.yml --limit 5
```

### Weekly Review

1. Review all bot commits: `git log --author="revenge-bot" --since="1 week ago"`
2. Check failed workflow runs: `gh run list --status failure`
3. Verify PRs created by bot are reasonable
4. Check for any security alerts

---

## Decommissioning

If bot account needs to be removed:

1. **Revoke all tokens**: Settings → Developer settings → Delete tokens
2. **Remove from repository**: Settings → Collaborators → Remove
3. **Update workflows**: Remove bot token references
4. **Delete secrets**: Settings → Secrets → Delete `BOT_TOKEN`
5. **Archive account**: Consider keeping account but revoking access

---

## Related Documentation

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Managing Bots in Organizations](https://docs.github.com/en/organizations/managing-organization-settings/managing-github-actions-settings-for-your-organization)
- [Personal Access Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)
- [Branch Protection Rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)

---

**Note**: The revenge-bot account has not been created yet. This documentation provides the setup process for when it's needed. For now, automated operations can use the `github-actions[bot]` account or personal tokens.
