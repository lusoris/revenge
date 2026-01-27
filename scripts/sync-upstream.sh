#!/usr/bin/env bash
# Sync script for Jellyfin upstream tracking
#
# This script maintains the following sync chain:
#   jellyfin/jellyfin (official) → lusoris/jellyfin (your fork) → lusoris/jellyfin-go (Go rewrite)
#
# Usage: ./scripts/sync-upstream.sh

set -euo pipefail

echo "🔄 Jellyfin Upstream Sync"
echo "========================="
echo ""

# 1. Fetch from official Jellyfin
echo "📥 Fetching from jellyfin/jellyfin (official)..."
git fetch upstream-official master

# 2. Fetch from your fork
echo "📥 Fetching from lusoris/jellyfin (your fork)..."
git fetch upstream master

# 3. Check if your fork is behind official
echo ""
echo "📊 Checking sync status..."
behind=$(git rev-list --count upstream/master..upstream-official/master)
ahead=$(git rev-list --count upstream-official/master..upstream/master)

echo ""
echo "Status:"
echo "  • Your fork is $behind commits behind official"
echo "  • Your fork is $ahead commits ahead of official"
echo ""

if [ "$behind" -gt 0 ]; then
    echo "⚠️  Your C# fork needs updating!"
    echo ""
    echo "To update your fork, run these commands:"
    echo "  cd /path/to/lusoris/jellyfin"
    echo "  git checkout master"
    echo "  git merge upstream-official/master"
    echo "  git push origin master"
    echo ""
else
    echo "✅ Your C# fork is up to date with official!"
    echo ""
fi

# 4. Update main branch with your fork's changes
current_branch=$(git branch --show-current)
echo "🔄 Syncing main branch with your C# fork..."

git checkout main
if ! git merge upstream/master --no-edit; then
    echo ""
    echo "❌ Merge conflicts detected!"
    echo "Keeping Go implementation..."
    git checkout --ours .
    git add -A
    git commit -m "chore: sync with upstream C# fork (kept Go implementation)"
fi

git push origin main
git checkout "$current_branch"

echo ""
echo "✅ Sync complete!"
echo ""
echo "Current remotes:"
git remote -v | grep fetch
