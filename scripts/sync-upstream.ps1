#!/usr/bin/env pwsh
# Sync script for Jellyfin upstream tracking
#
# This script maintains the following sync chain:
#   jellyfin/jellyfin (official) → lusoris/jellyfin (your fork) → lusoris/revenge (Go rewrite)
#
# Usage: .\scripts\sync-upstream.ps1

$ErrorActionPreference = "Stop"

Write-Host "🔄 Revenge Upstream Sync" -ForegroundColor Cyan
Write-Host "=========================" -ForegroundColor Cyan
Write-Host ""

# 1. Fetch from official Jellyfin
Write-Host "📥 Fetching from jellyfin/jellyfin (official)..." -ForegroundColor Yellow
git fetch upstream-official master
if ($LASTEXITCODE -ne 0) { exit 1 }

# 2. Fetch from your fork
Write-Host "📥 Fetching from lusoris/jellyfin (your fork)..." -ForegroundColor Yellow
git fetch upstream master
if ($LASTEXITCODE -ne 0) { exit 1 }

# 3. Check if your fork is behind official
Write-Host ""
Write-Host "📊 Checking sync status..." -ForegroundColor Yellow
$behind = git rev-list --count upstream/master..upstream-official/master
$ahead = git rev-list --count upstream-official/master..upstream/master

Write-Host ""
Write-Host "Status:" -ForegroundColor Cyan
Write-Host "  • Your fork is $behind commits behind official" -ForegroundColor $(if ($behind -gt 0) { "Red" } else { "Green" })
Write-Host "  • Your fork is $ahead commits ahead of official" -ForegroundColor $(if ($ahead -gt 0) { "Yellow" } else { "Green" })
Write-Host ""

if ($behind -gt 0) {
    Write-Host "⚠️  Your C# fork needs updating!" -ForegroundColor Red
    Write-Host ""
    Write-Host "To update your fork, run these commands:" -ForegroundColor Yellow
    Write-Host "  cd /path/to/lusoris/jellyfin" -ForegroundColor Gray
    Write-Host "  git checkout master" -ForegroundColor Gray
    Write-Host "  git merge upstream-official/master" -ForegroundColor Gray
    Write-Host "  git push origin master" -ForegroundColor Gray
    Write-Host ""
} else {
    Write-Host "✅ Your C# fork is up to date with official!" -ForegroundColor Green
    Write-Host ""
}

# 4. Update main branch with your fork's changes
$currentBranch = git branch --show-current
Write-Host "🔄 Syncing main branch with your C# fork..." -ForegroundColor Yellow

git checkout main
if ($LASTEXITCODE -ne 0) { exit 1 }

git merge upstream/master --no-edit
if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "❌ Merge conflicts detected!" -ForegroundColor Red
    Write-Host "Keeping Go implementation..." -ForegroundColor Yellow
    git checkout --ours .
    git add -A
    git commit -m "chore: sync with upstream C# fork (kept Go implementation)"
}

git push origin main
if ($LASTEXITCODE -ne 0) { exit 1 }

git checkout $currentBranch
if ($LASTEXITCODE -ne 0) { exit 1 }

Write-Host ""
Write-Host "✅ Sync complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Current remotes:" -ForegroundColor Cyan
git remote -v | Select-String "fetch"
