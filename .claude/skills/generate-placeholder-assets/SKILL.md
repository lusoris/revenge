---
name: generate-placeholder-assets
description: Generate placeholder branding assets (logos, screenshots, social media images)
argument-hint: "[--all|--branding|--screenshots|--social]"
disable-model-invocation: false
allowed-tools: Bash(*), Read(*)
---

# Generate Placeholder Assets

Generates placeholder assets for documentation using the asset generation script. Creates temporary branding assets, screenshot placeholders, and social media images for use until final designs are ready.

## Usage

```
/generate-placeholder-assets              # Generate all assets
/generate-placeholder-assets --branding   # Generate branding only (logos, icons, favicons)
/generate-placeholder-assets --screenshots  # Generate screenshot placeholders only
/generate-placeholder-assets --social     # Generate social media images only
```

## Arguments

- `$0`: Asset type (optional: --all, --branding, --screenshots, --social)

## Prerequisites

- Python 3.10+ with Pillow installed
- `docs/assets/` directory structure exists

## Task

Generate placeholder assets and report what was created.

### Step 1: Verify Prerequisites

Check that the asset generation script and dependencies are available:

```bash
# Check Python is available
if ! command -v python &> /dev/null; then
    echo "❌ Python not found"
    exit 1
fi

# Check Pillow is installed
if ! python -c "import PIL" 2>/dev/null; then
    echo "❌ Pillow not installed"
    echo "Install: pip install Pillow"
    exit 1
fi

# Verify asset generation script exists
if [ ! -f "scripts/automation/generate_placeholder_assets.py" ]; then
    echo "❌ Asset generation script not found"
    exit 1
fi

# Verify assets directory exists
if [ ! -d "docs/assets" ]; then
    echo "⚠️ Assets directory not found, creating..."
    mkdir -p docs/assets/{branding,placeholders/{screenshots,diagrams},social,badges}
fi
```

### Step 2: Generate Assets

Run the asset generation script based on the requested type:

**For all assets** (default or --all):
```bash
echo "🎨 Generating all placeholder assets..."
python scripts/automation/generate_placeholder_assets.py --all
```

**For branding only** (--branding):
```bash
echo "🎨 Generating branding assets (logos, icons, favicons)..."
python scripts/automation/generate_placeholder_assets.py --branding
```

**For screenshots only** (--screenshots):
```bash
echo "📸 Generating screenshot placeholders..."
python scripts/automation/generate_placeholder_assets.py --screenshots
```

**For social media only** (--social):
```bash
echo "🌐 Generating social media images..."
python scripts/automation/generate_placeholder_assets.py --social
```

### Step 3: Verify Generation

Check that assets were created successfully:

```bash
echo ""
echo "📊 Generated assets:"
echo ""
echo "Branding:"
ls -lh docs/assets/branding/ 2>/dev/null | tail -n +2 | awk '{print "  ", $9, "(" $5 ")"}'

echo ""
echo "Screenshots:"
ls -lh docs/assets/placeholders/screenshots/ 2>/dev/null | tail -n +2 | awk '{print "  ", $9, "(" $5 ")"}'

echo ""
echo "Social Media:"
ls -lh docs/assets/social/ 2>/dev/null | tail -n +2 | awk '{print "  ", $9, "(" $5 ")"}'

echo ""
echo "✅ Asset generation complete!"
```

### Step 4: Report Summary

Provide a summary of what was generated:

```bash
echo ""
echo "Summary:"
echo "  Location: docs/assets/"
echo "  Total files: $(find docs/assets -type f | wc -l)"
echo "  Total size: $(du -sh docs/assets | cut -f1)"
echo ""
echo "Next steps:"
echo "  1. Review assets: ls -R docs/assets/"
echo "  2. Test in documentation (README, wiki, etc.)"
echo "  3. Replace with real assets when UI is implemented (post-v1)"
```

## Output Format

The script should output:
- List of generated files with sizes
- Total count and size
- Next steps for using the assets

## Error Handling

- If Pillow is not installed, instruct user to install it
- If asset directory doesn't exist, create it automatically
- If generation fails, show the error and suggest manual installation

## Notes

- These are **placeholder assets** - temporary until real branding is designed
- Uses default color scheme (blue/purple theme)
- All assets are optimized and compressed
- Safe to run multiple times (will overwrite existing placeholders)
