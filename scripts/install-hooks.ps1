# Script to install Git hooks on Windows

$ScriptDir = Split-Path -Parent $PSCommandPath
$ProjectRoot = Split-Path -Parent $ScriptDir
$HooksDir = Join-Path $ProjectRoot ".githooks"
$GitHooksDir = Join-Path $ProjectRoot ".git\hooks"

Write-Host "🔧 Installing Git hooks..." -ForegroundColor Green

# Create hooks directory if it doesn't exist
if (-not (Test-Path $GitHooksDir)) {
    New-Item -ItemType Directory -Path $GitHooksDir | Out-Null
}

# Install hooks
Get-ChildItem $HooksDir | ForEach-Object {
    $hookName = $_.Name
    Write-Host "  📌 Installing $hookName hook..." -ForegroundColor Cyan
    
    $source = $_.FullName
    $target = Join-Path $GitHooksDir $hookName
    
    # Remove existing hook if it exists
    if (Test-Path $target) {
        Remove-Item $target -Force
    }
    
    # Create symlink (requires admin) or copy file
    try {
        New-Item -ItemType SymbolicLink -Path $target -Target $source -Force | Out-Null
    } catch {
        # If symlink fails (no admin), copy instead
        Copy-Item $source $target -Force
        Write-Host "    ⚠️  Created copy instead of symlink (requires admin for symlinks)" -ForegroundColor Yellow
    }
}

Write-Host "✅ Git hooks installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Installed hooks:" -ForegroundColor White
Write-Host "  - pre-commit: Runs formatting, vetting, and tests"
Write-Host "  - commit-msg: Validates commit message format"
Write-Host "  - pre-push: Runs full checks before pushing"
Write-Host ""
Write-Host "To bypass hooks temporarily, use: git commit --no-verify" -ForegroundColor Yellow
