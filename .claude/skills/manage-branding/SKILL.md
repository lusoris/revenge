---
name: manage-branding
description: Manage project branding assets (generate, validate, optimize, report)
argument-hint: "[--generate|--validate|--optimize|--report]"
disable-model-invocation: false
allowed-tools: Bash(*), Read(*), Glob(*)
---

# Manage Branding

Comprehensive branding asset management - generate placeholders, validate existing assets, optimize file sizes, and report on branding status.

## Usage

```
/manage-branding              # Interactive menu
/manage-branding --generate   # Generate all missing assets
/manage-branding --validate   # Validate all assets exist and are correct
/manage-branding --optimize   # Optimize/compress all assets
/manage-branding --report     # Generate branding status report
```

## Arguments

- `$0`: Action (optional: --generate, --validate, --optimize, --report)

## Prerequisites

- Python 3.10+ with Pillow and requests libraries
- `docs/assets/` directory structure

## Task

Manage all project branding assets with comprehensive validation and reporting.

### Step 1: Assess Current State

Check what assets exist and their status:

```bash
echo "🔍 Assessing branding assets..."
echo ""

# Check branding directory
if [ -d "docs/assets/branding" ]; then
    BRANDING_COUNT=$(find docs/assets/branding -type f | wc -l)
    echo "Branding assets: $BRANDING_COUNT files"
else
    echo "⚠️  Branding directory missing"
    BRANDING_COUNT=0
fi

# Check screenshots
if [ -d "docs/assets/placeholders/screenshots" ]; then
    SCREENSHOT_COUNT=$(find docs/assets/placeholders/screenshots -type f | wc -l)
    echo "Screenshots: $SCREENSHOT_COUNT files"
else
    echo "⚠️  Screenshots directory missing"
    SCREENSHOT_COUNT=0
fi

# Check social media
if [ -d "docs/assets/social" ]; then
    SOCIAL_COUNT=$(find docs/assets/social -type f | wc -l)
    echo "Social media: $SOCIAL_COUNT files"
else
    echo "⚠️  Social media directory missing"
    SOCIAL_COUNT=0
fi

# Check badges
if [ -d "docs/assets/badges" ]; then
    BADGE_COUNT=$(find docs/assets/badges -type f | wc -l)
    echo "Badges: $BADGE_COUNT files"
else
    echo "⚠️  Badges directory missing"
    BADGE_COUNT=0
fi

echo ""
TOTAL_ASSETS=$((BRANDING_COUNT + SCREENSHOT_COUNT + SOCIAL_COUNT + BADGE_COUNT))
echo "Total assets: $TOTAL_ASSETS files"
```

### Step 2: Execute Requested Action

**Interactive menu** (no argument):
```bash
echo ""
echo "What would you like to do?"
echo "  1) Generate all assets (placeholders + badges)"
echo "  2) Validate existing assets"
echo "  3) Optimize/compress assets"
echo "  4) Generate status report"
echo "  5) Regenerate missing assets only"
echo ""
echo "Recommendation:"
if [ $TOTAL_ASSETS -eq 0 ]; then
    echo "  → Run option 1 to generate all assets"
elif [ $TOTAL_ASSETS -lt 30 ]; then
    echo "  → Some assets missing, run option 5 to fill gaps"
else
    echo "  → Assets look complete, run option 4 for status report"
fi
```

**Generate all** (--generate):
```bash
echo "🎨 Generating all branding assets..."

# Generate placeholders
python scripts/automation/generate_placeholder_assets.py --all

# Generate badges
python scripts/automation/generate_badges.py --all

echo "✅ All assets generated!"
```

**Validate assets** (--validate):
```bash
echo "🔍 Validating branding assets..."
echo ""

# Required branding assets
REQUIRED_BRANDING=(
    "logo.png"
    "logo-light.png"
    "logo-dark.png"
    "logo-icon.png"
    "logo-wordmark.png"
    "favicon.ico"
    "apple-touch-icon.png"
)

echo "Checking branding assets:"
MISSING_BRANDING=0
for asset in "${REQUIRED_BRANDING[@]}"; do
    if [ -f "docs/assets/branding/$asset" ]; then
        SIZE=$(du -h "docs/assets/branding/$asset" | cut -f1)
        echo "  ✓ $asset ($SIZE)"
    else
        echo "  ✗ $asset MISSING"
        MISSING_BRANDING=$((MISSING_BRANDING + 1))
    fi
done

# Required social media assets
REQUIRED_SOCIAL=(
    "og-image.png"
    "og-image-home.png"
    "og-image-docs.png"
    "twitter-card.png"
)

echo ""
echo "Checking social media assets:"
MISSING_SOCIAL=0
for asset in "${REQUIRED_SOCIAL[@]}"; do
    if [ -f "docs/assets/social/$asset" ]; then
        SIZE=$(du -h "docs/assets/social/$asset" | cut -f1)
        echo "  ✓ $asset ($SIZE)"
    else
        echo "  ✗ $asset MISSING"
        MISSING_SOCIAL=$((MISSING_SOCIAL + 1))
    fi
done

# Required badges
REQUIRED_BADGES=(
    "coverage.svg"
    "build.svg"
    "version.svg"
    "license.svg"
)

echo ""
echo "Checking badges:"
MISSING_BADGES=0
for asset in "${REQUIRED_BADGES[@]}"; do
    if [ -f "docs/assets/badges/$asset" ]; then
        SIZE=$(du -h "docs/assets/badges/$asset" | cut -f1)
        echo "  ✓ $asset ($SIZE)"
    else
        echo "  ✗ $asset MISSING"
        MISSING_BADGES=$((MISSING_BADGES + 1))
    fi
done

echo ""
TOTAL_MISSING=$((MISSING_BRANDING + MISSING_SOCIAL + MISSING_BADGES))
if [ $TOTAL_MISSING -eq 0 ]; then
    echo "✅ All required assets present!"
else
    echo "⚠️  $TOTAL_MISSING assets missing"
    echo "Run: /manage-branding --generate"
fi
```

**Optimize assets** (--optimize):
```bash
echo "🗜️  Optimizing assets..."
echo ""

# Check if optipng is available
if command -v optipng &> /dev/null; then
    echo "Optimizing PNG files with optipng..."
    find docs/assets -name "*.png" -exec optipng -o2 {} \;
elif command -v pngcrush &> /dev/null; then
    echo "Optimizing PNG files with pngcrush..."
    find docs/assets -name "*.png" -exec pngcrush -ow {} \;
else
    echo "⚠️  No PNG optimizer found (install optipng or pngcrush)"
    echo "Using Python Pillow for optimization..."
    python -c "
from PIL import Image
from pathlib import Path

for png in Path('docs/assets').rglob('*.png'):
    img = Image.open(png)
    img.save(png, optimize=True, quality=95)
    print(f'  ✓ Optimized: {png.name}')
"
fi

echo ""
echo "✅ Optimization complete!"
echo "New total size: $(du -sh docs/assets | cut -f1)"
```

**Status report** (--report):
```bash
echo "📊 Branding Status Report"
echo "=========================="
echo ""

# Asset counts
echo "Asset Inventory:"
echo "  Branding:   $BRANDING_COUNT files"
echo "  Screenshots: $SCREENSHOT_COUNT files"
echo "  Social:     $SOCIAL_COUNT files"
echo "  Badges:     $BADGE_COUNT files"
echo "  Total:      $TOTAL_ASSETS files"
echo ""

# Total size
TOTAL_SIZE=$(du -sh docs/assets 2>/dev/null | cut -f1)
echo "Total Size: $TOTAL_SIZE"
echo ""

# Recent changes
echo "Recent Changes:"
find docs/assets -type f -mtime -7 | head -5 | while read file; do
    echo "  $(stat -c '%y %n' "$file" | cut -d' ' -f1,4)"
done
echo ""

# Status
echo "Status:"
if [ $TOTAL_ASSETS -ge 30 ]; then
    echo "  ✅ Complete - all asset categories populated"
elif [ $TOTAL_ASSETS -ge 15 ]; then
    echo "  🟡 Partial - some assets generated"
else
    echo "  🔴 Incomplete - run /manage-branding --generate"
fi
echo ""

# Next steps
echo "Next Steps:"
echo "  • Review assets: ls -R docs/assets/"
echo "  • Update README.md with badges"
echo "  • Add logo to project templates"
echo "  • Replace placeholders with real designs (post-v1)"
```

### Step 3: Summary

Provide a summary of actions taken:

```bash
echo ""
echo "=========================================="
echo "Branding Management Complete"
echo "=========================================="
echo ""
echo "Location: docs/assets/"
echo "Total Assets: $(find docs/assets -type f | wc -l) files"
echo "Total Size: $(du -sh docs/assets 2>/dev/null | cut -f1)"
```

## Output Format

The skill should output:
- Current asset inventory
- Action results (what was generated/validated/optimized)
- Status summary
- Next steps

## Error Handling

- If asset directories don't exist, create them automatically
- If generation fails, report specific errors
- If optimization tools aren't installed, use fallback methods

## Notes

- All assets are placeholders until final branding is designed
- Safe to run multiple times
- Validates required assets vs. actual assets
- Optimization reduces file sizes without quality loss
- Report shows asset health and completeness
