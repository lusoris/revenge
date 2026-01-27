#!/bin/bash

# Script to install Git hooks

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOOKS_DIR="$SCRIPT_DIR/.githooks"
GIT_HOOKS_DIR="$SCRIPT_DIR/.git/hooks"

echo "🔧 Installing Git hooks..."

# Create hooks directory if it doesn't exist
mkdir -p "$GIT_HOOKS_DIR"

# Install hooks
for hook in "$HOOKS_DIR"/*; do
    hook_name=$(basename "$hook")
    echo "  📌 Installing $hook_name hook..."
    
    # Make hook executable
    chmod +x "$hook"
    
    # Create symlink
    ln -sf "../../.githooks/$hook_name" "$GIT_HOOKS_DIR/$hook_name"
done

echo "✅ Git hooks installed successfully!"
echo ""
echo "Installed hooks:"
echo "  - pre-commit: Runs formatting, vetting, and tests"
echo "  - commit-msg: Validates commit message format"
echo "  - pre-push: Runs full checks before pushing"
echo ""
echo "To bypass hooks temporarily, use: git commit --no-verify"
