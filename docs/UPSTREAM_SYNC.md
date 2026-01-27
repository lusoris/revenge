# Upstream Sync Guide

This document explains how to keep the Go rewrite in sync with the original Revenge C# project.

## Repository Structure

```
revenge/revenge (official C# - upstream-official)
    ↓ (you sync manually)
lusoris/revenge (your C# fork - upstream)
    ↓ (tracked for reference)
lusoris/revenge (Go rewrite - origin)
```

## Remotes Configuration

```bash
origin              https://github.com/lusoris/revenge.git    # Go rewrite
upstream            https://github.com/lusoris/revenge.git       # Your C# fork
upstream-official   https://github.com/revenge/revenge.git      # Official C#
```

## Automated Sync (Recommended)

Run the sync script to check for updates:

**Windows (PowerShell):**
```powershell
.\scripts\sync-upstream.ps1
```

**Linux/Mac:**
```bash
./scripts/sync-upstream.sh
```

This script will:
1. ✅ Fetch latest from official Revenge
2. ✅ Fetch latest from your C# fork
3. ✅ Report sync status
4. ✅ Merge your fork's changes into `main` branch
5. ✅ Handle conflicts automatically (keeps Go implementation)

## Manual Sync

### Step 1: Update Your C# Fork

```bash
# In your C# fork repository (lusoris/revenge)
cd /path/to/lusoris/revenge
git checkout master
git fetch upstream master  # If you have official as 'upstream'
git merge upstream/master
git push origin master
```

### Step 2: Update Go Rewrite Main Branch

```bash
# In this repository (lusoris/revenge)
cd /path/to/revenge

# Fetch latest from your fork
git fetch upstream master

# Merge into main branch
git checkout main
git merge upstream/master --allow-unrelated-histories

# Handle conflicts (keep Go implementation)
git checkout --ours .
git add -A
git commit -m "chore: sync with upstream C# fork"

# Push changes
git push origin main
```

### Step 3: Check Status in Develop

```bash
git checkout develop
# Continue development
```

## Tracking Important Changes

When syncing, pay attention to:

### API Changes
- New endpoints in Revenge.Api controllers
- Changed request/response models
- Authentication/authorization updates

### Database Schema Changes
- New tables or columns
- Migration scripts in Revenge.Server.Implementations

### Business Logic Changes
- Service implementations
- Media processing updates
- Library scanning improvements

### Configuration Changes
- New settings
- Changed defaults
- Environment variables

## Sync Frequency

**Recommended sync schedule:**
- **Weekly**: Run automated sync script to check status
- **Monthly**: Review changes for new features to implement
- **Before releases**: Ensure compatibility with latest Revenge version

## GitHub Actions (Future)

Consider setting up automated sync via GitHub Actions:

```yaml
name: Upstream Sync Check
on:
  schedule:
    - cron: '0 0 * * 0'  # Weekly on Sunday
  workflow_dispatch:

jobs:
  sync-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Check Upstream Status
        run: ./scripts/sync-upstream.sh
```

## Useful Commands

```bash
# Check how many commits behind/ahead
git fetch upstream-official master
git rev-list --count upstream/master..upstream-official/master  # behind
git rev-list --count upstream-official/master..upstream/master  # ahead

# View new commits in official
git log upstream/master..upstream-official/master --oneline

# View changes to specific files
git diff upstream/master upstream-official/master -- path/to/file

# Find commits by keyword
git log upstream-official/master --grep="keyword" --oneline
```

## Notes

- The `main` branch contains the C# reference code for historical context
- Active Go development happens on `develop` branch
- Merge conflicts are expected and resolved by keeping Go implementation
- Focus on tracking API contracts and business logic changes, not implementation details
