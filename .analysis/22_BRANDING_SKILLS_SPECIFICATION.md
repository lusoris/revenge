# Branding & Assets Skills Specification

**Created**: 2026-01-31
**Purpose**: Skills for managing branding, assets, placeholders, and badges
**Extends**: 17_CLAUDE_SKILLS_SPECIFICATION.md (now 30 total skills)

---

## New Skills for Branding & Assets (5 skills)

### Skill 26: generate-placeholder-assets

**Purpose**: Generate all placeholder assets (logos, screenshots, diagrams, social media)

**Usage**:
```bash
/generate-placeholder-assets
/generate-placeholder-assets --type logos
/generate-placeholder-assets --type screenshots
/generate-placeholder-assets --force  # Overwrite existing
```

**Features**:
- Generate placeholder logos (SVG + PNG + ICO)
- Generate placeholder screenshots (1200x800, 1600x900, 375x667)
- Generate placeholder diagrams (Mermaid templates)
- Generate social media assets (OpenGraph, Twitter cards)
- Generate favicons (ICO, SVG, Apple touch icon)
- Optimize images (compress, resize)
- Validate generated assets

**Interactive Flow**:
1. Ask: Generate all assets or specific type?
2. Ask: Color scheme (use default or custom?)
3. Ask: Project name for text overlay
4. Generate assets using Pillow (Python)
5. Show summary of generated files
6. Validate all assets created

**Parameters**:
- `--type` (`-t`): Asset type (logos, screenshots, diagrams, social, all)
- `--color-scheme`: Color scheme (default, blue, purple, custom)
- `--project-name`: Project name for text overlay
- `--force`: Overwrite existing assets
- `--optimize`: Optimize/compress images

**Output**:
```
🎨 Generating placeholder assets...

[1/5] Placeholder Logos
✅ logo.svg (512x512)
✅ logo.png (1024x1024)
✅ logo-light.svg
✅ logo-dark.svg
✅ logo-icon.svg (256x256)
✅ logo-wordmark.svg
✅ favicon.ico (16x16, 32x32, 48x48)
✅ favicon.svg
✅ apple-touch-icon.png (180x180)

[2/5] Placeholder Screenshots
✅ placeholder-movie-library.png (1200x800)
✅ placeholder-settings-general.png (1200x800)
✅ placeholder-player-controls.png (1200x800)
... (15 screenshots generated)

[3/5] Placeholder Diagrams
✅ system-architecture.svg (Mermaid template)
✅ data-flow.svg (Mermaid template)
... (5 diagrams generated)

[4/5] Social Media Assets
✅ og-image.png (1200x630)
✅ og-image-home.png (1200x630)
✅ og-image-docs.png (1200x630)
✅ og-image-wiki.png (1200x630)
✅ twitter-card.png (1200x600)

[5/5] Optimizing
✅ Compressed 32 images (saved 1.2 MB)

Summary:
- Logos: 9 files
- Screenshots: 15 files
- Diagrams: 5 files
- Social: 5 files
- Total: 34 assets generated
- Total size: 3.4 MB (optimized)
- Location: docs/assets/

Next steps:
1. Review assets: ls -lh docs/assets/
2. Test in docs: /generate-docs --doc README
3. Replace with real assets when ready (post-v1)
```

**Implementation**:
```python
from PIL import Image, ImageDraw, ImageFont
import cairosvg

def create_placeholder_logo(size: int, text: str, color_scheme: dict):
    """Create placeholder logo with text and geometric shapes."""
    img = Image.new('RGBA', (size, size), color=(0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # Draw circle background
    draw.ellipse(
        [(size*0.1, size*0.1), (size*0.9, size*0.9)],
        fill=color_scheme['primary']
    )

    # Draw text (centered)
    font = ImageFont.truetype('/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf', size//4)
    draw.text(
        (size/2, size/2),
        text,
        fill='white',
        font=font,
        anchor='mm'
    )

    return img

def create_placeholder_screenshot(size: tuple, feature: str):
    """Create placeholder screenshot with feature name."""
    img = Image.new('RGB', size, color='#f9fafb')
    draw = ImageDraw.Draw(img)

    # Draw border
    draw.rectangle(
        [(20, 20), (size[0]-20, size[1]-20)],
        outline='#d1d5db',
        width=4
    )

    # Draw title
    font_title = ImageFont.truetype('...', 48)
    draw.text(
        (size[0]/2, size[1]/2 - 40),
        feature,
        fill='#111827',
        font=font_title,
        anchor='mm'
    )

    # Draw subtitle
    font_sub = ImageFont.truetype('...', 24)
    draw.text(
        (size[0]/2, size[1]/2 + 20),
        'Screenshot Placeholder',
        fill='#6b7280',
        font=font_sub,
        anchor='mm'
    )

    # Draw note
    font_note = ImageFont.truetype('...', 20)
    draw.text(
        (size[0]/2, size[1]/2 + 60),
        '(Coming post-v1.0)',
        fill='#9ca3af',
        font=font_note,
        anchor='mm'
    )

    return img
```

---

### Skill 27: generate-badges

**Purpose**: Generate shields.io badges (coverage, build, version, license)

**Usage**:
```bash
/generate-badges
/generate-badges --type coverage
/generate-badges --update  # Update all badges
```

**Features**:
- Generate coverage badge from coverage report
- Generate build badge from GitHub Actions status
- Generate version badge from git tags / go.mod
- Generate license badge from LICENSE file
- Custom badges (downloads, contributors, etc.)
- Download from shields.io API
- Save to `docs/assets/badges/`

**Interactive Flow**:
1. Ask: Generate all badges or specific type?
2. Read data sources (coverage.out, git tags, LICENSE)
3. Generate badge URLs
4. Download from shields.io
5. Save to docs/assets/badges/
6. Show summary

**Parameters**:
- `--type` (`-t`): Badge type (coverage, build, version, license, all)
- `--update`: Update existing badges
- `--custom`: Create custom badge (label + value + color)

**Output**:
```
📛 Generating badges...

[1/4] Coverage Badge
📊 Reading coverage report: coverage.out
   Coverage: 84.2%
   Color: green (>80%)
✅ docs/assets/badges/coverage.svg

[2/4] Build Badge
🔨 Checking GitHub Actions: lusoris/revenge
   Latest run: #1234 (success)
   Color: brightgreen
✅ docs/assets/badges/build.svg

[3/4] Version Badge
🏷️  Reading version from git tags
   Latest tag: v0.3.0
   Color: blue
✅ docs/assets/badges/version.svg

[4/4] License Badge
⚖️  Reading LICENSE file
   License: AGPL-3.0
   Color: blue
✅ docs/assets/badges/license.svg

Summary:
- Total badges: 4
- Updated: 4
- Failed: 0
- Location: docs/assets/badges/

Badge URLs (for README):
[![Coverage](docs/assets/badges/coverage.svg)](https://codecov.io/gh/lusoris/revenge)
[![Build](docs/assets/badges/build.svg)](https://github.com/lusoris/revenge/actions)
[![Version](docs/assets/badges/version.svg)](https://github.com/lusoris/revenge/releases)
[![License](docs/assets/badges/license.svg)](LICENSE)

Next steps:
1. Badges auto-update on CI runs
2. Regenerate docs: /generate-docs --doc README
```

**Implementation**:
```python
import requests
import re

def generate_coverage_badge(coverage_file: Path) -> str:
    """Generate coverage badge URL and download SVG."""
    # Parse coverage from coverage.out
    with open(coverage_file) as f:
        content = f.read()
        match = re.search(r'total:\s+\(statements\)\s+([\d.]+)%', content)
        coverage = float(match.group(1)) if match else 0.0

    # Determine color
    if coverage >= 80:
        color = 'brightgreen'
    elif coverage >= 60:
        color = 'yellow'
    else:
        color = 'red'

    # Generate URL
    url = f"https://img.shields.io/badge/coverage-{coverage:.1f}%25-{color}"

    # Download SVG
    response = requests.get(url)
    output = Path('docs/assets/badges/coverage.svg')
    output.write_bytes(response.content)

    return str(output)

def generate_build_badge(repo: str) -> str:
    """Generate build badge from GitHub Actions."""
    # Use GitHub API to get latest workflow run
    url = f"https://api.github.com/repos/{repo}/actions/runs?per_page=1"
    response = requests.get(url)
    data = response.json()

    if data['workflow_runs']:
        status = data['workflow_runs'][0]['conclusion']
        color = 'brightgreen' if status == 'success' else 'red'
        label = 'passing' if status == 'success' else 'failing'
    else:
        color = 'lightgrey'
        label = 'unknown'

    # Generate URL
    badge_url = f"https://img.shields.io/badge/build-{label}-{color}"

    # Download SVG
    response = requests.get(badge_url)
    output = Path('docs/assets/badges/build.svg')
    output.write_bytes(response.content)

    return str(output)
```

---

### Skill 28: update-branding

**Purpose**: Update branding assets across all documentation

**Usage**:
```bash
/update-branding
/update-branding --regenerate-docs
/update-branding --asset logo.svg
```

**Features**:
- Update logo references in all docs
- Update color schemes in docs
- Regenerate docs with new branding
- Validate all asset references
- Update social media assets

**Interactive Flow**:
1. Ask: What changed? (logo, colors, both)
2. If logo: Ask for new logo path
3. If colors: Ask for color scheme YAML
4. Update all doc data files with new asset paths
5. Regenerate affected docs
6. Validate all references work
7. Show summary

**Parameters**:
- `--asset`: Asset file to update (logo.svg, colors.yml)
- `--regenerate-docs`: Regenerate all docs with new branding
- `--validate`: Only validate, don't update

**Output**:
```
🎨 Updating branding...

Changes detected:
- New logo: docs/assets/branding/logo-new.svg
- Updated colors: Updated 5 colors

[1/4] Updating asset references
✅ Updated 136 data files
  - Updated logo path: docs/assets/branding/logo-new.svg
  - Updated color references

[2/4] Updating SOT
✅ Updated branding section in SOURCE_OF_TRUTH.md

[3/4] Regenerating documentation
🔄 Regenerating all human-readable docs (42 files)
✅ README.md
✅ CONTRIBUTING.md
✅ docs/wiki/Home.md
... (39 more files)

[4/4] Validating
✅ All asset paths valid
✅ All colors valid hex codes
✅ All docs render correctly

Summary:
- Data files updated: 136
- Docs regenerated: 42
- Assets updated: 1 logo, 5 colors
- Validation: Passed

Next steps:
1. Review updated docs: git diff
2. Commit changes: git add . && git commit -m "chore: update branding"
3. Push: git push
```

---

### Skill 29: capture-screenshots

**Purpose**: Capture real screenshots from running application (POST-v1 feature)

**Usage**:
```bash
/capture-screenshots
/capture-screenshots --feature movies
/capture-screenshots --theme dark
/capture-screenshots --all-themes
```

**Features**:
- Use Playwright to automate browser
- Navigate to feature pages
- Capture screenshots at defined breakpoints
- Capture multiple themes (light/dark)
- Capture multiple resolutions (desktop, tablet, mobile)
- Save to `docs/assets/screenshots/`
- Update YAML data files with real paths
- Regenerate docs with real screenshots

**Interactive Flow**:
1. Ask: Capture all features or specific?
2. Ask: Which themes? (light, dark, both)
3. Ask: Which resolutions? (desktop, tablet, mobile, all)
4. Start application (check if running)
5. Launch Playwright browser
6. Navigate to each feature page
7. Capture screenshots
8. Update data files
9. Regenerate docs
10. Show summary

**Parameters**:
- `--feature`: Feature to capture (movies, music, settings, all)
- `--theme`: Theme (light, dark, both)
- `--resolution`: Resolution (desktop, tablet, mobile, all)
- `--app-url`: Application URL (default: http://localhost:8096)
- `--update-docs`: Regenerate docs after capture

**Output**:
```
📸 Capturing screenshots...

Starting Playwright...
✅ Browser launched (Chromium)
✅ Application running at http://localhost:8096

[1/3] Desktop (1920x1080)
  Light Theme:
    ✅ movies-library.png
    ✅ movie-details.png
    ✅ settings-general.png
  Dark Theme:
    ✅ movies-library-dark.png
    ✅ movie-details-dark.png
    ✅ settings-general-dark.png

[2/3] Tablet (768x1024)
  Light Theme:
    ✅ movies-library-tablet.png
    ✅ movie-details-tablet.png
  Dark Theme:
    ✅ movies-library-tablet-dark.png
    ✅ movie-details-tablet-dark.png

[3/3] Mobile (375x667)
  Light Theme:
    ✅ movies-library-mobile.png
  Dark Theme:
    ✅ movies-library-mobile-dark.png

Updating data files...
✅ Updated 5 data files with real screenshot paths

Regenerating docs...
✅ Regenerated 8 wiki docs
✅ Regenerated 3 user docs

Summary:
- Screenshots captured: 12
- Themes: light + dark
- Resolutions: desktop, tablet, mobile
- Data files updated: 5
- Docs regenerated: 11
- Location: docs/assets/screenshots/

Next steps:
1. Review screenshots: ls -lh docs/assets/screenshots/
2. Review updated docs: git diff docs/wiki/
3. Commit: git add . && git commit -m "docs: add real screenshots"
```

**Implementation**:
```python
from playwright.sync_api import sync_playwright

def capture_screenshot(
    url: str,
    output_path: Path,
    viewport_size: tuple,
    theme: str = 'light'
):
    """Capture screenshot using Playwright."""
    with sync_playwright() as p:
        browser = p.chromium.launch()
        context = browser.new_context(
            viewport={'width': viewport_size[0], 'height': viewport_size[1]},
            color_scheme=theme
        )
        page = context.new_page()
        page.goto(url)
        page.wait_for_load_state('networkidle')
        page.screenshot(path=str(output_path))
        browser.close()

# Usage
capture_screenshot(
    url='http://localhost:8096/movies',
    output_path=Path('docs/assets/screenshots/movies-library.png'),
    viewport_size=(1920, 1080),
    theme='light'
)
```

**Prerequisites**:
- Application must be running
- Playwright installed (`pip install playwright && playwright install`)
- Valid test data in application

---

### Skill 30: validate-assets

**Purpose**: Validate all assets exist and are properly referenced

**Usage**:
```bash
/validate-assets
/validate-assets --fix  # Fix broken references
/validate-assets --report  # Generate detailed report
```

**Features**:
- Check all asset paths in data files exist
- Check all asset references in generated docs exist
- Check image dimensions (correct sizes)
- Check file sizes (not too large)
- Check image formats (valid SVG, PNG, etc.)
- Find orphaned assets (not referenced anywhere)
- Find missing assets (referenced but don't exist)
- Generate validation report

**Interactive Flow**:
1. Read all data files
2. Extract all asset references
3. Check if files exist
4. Check file properties (size, dimensions, format)
5. Find orphaned assets
6. Show validation report
7. If `--fix`: Offer to fix broken references

**Parameters**:
- `--fix`: Attempt to fix broken references
- `--report`: Generate detailed validation report
- `--clean`: Remove orphaned assets

**Output**:
```
🔍 Validating assets...

[1/5] Loading data files
✅ Loaded 136 data files

[2/5] Extracting asset references
✅ Found 247 asset references
  - Logos: 42 references
  - Screenshots: 156 references
  - Diagrams: 38 references
  - Social: 11 references

[3/5] Checking files exist
✅ 245/247 assets exist
❌ 2 missing:
  - docs/assets/screenshots/music-player.png (referenced in MUSIC_MODULE.yaml)
  - docs/assets/diagrams/auth-flow.svg (referenced in AUTH.yaml)

[4/5] Validating file properties
✅ Image dimensions: 245/245 correct
✅ File sizes: 245/245 under 2MB
✅ File formats: 245/245 valid

[5/5] Finding orphaned assets
⚠️  3 orphaned assets (not referenced):
  - docs/assets/placeholders/old-logo.png (1.2 MB)
  - docs/assets/screenshots/test-screenshot.png (0.5 MB)
  - docs/assets/social/old-og-image.png (0.3 MB)

Summary:
- Total references: 247
- Valid: 245 (99.2%)
- Missing: 2 (0.8%)
- Orphaned: 3
- Total size: 12.4 MB

Issues to fix:
1. Add missing screenshots: music-player.png, auth-flow.svg
2. Remove orphaned assets (saves 2.0 MB)

Fix issues? [y/n]: y

Fixing...
✅ Removed 3 orphaned assets (saved 2.0 MB)
⚠️  Missing assets cannot be auto-fixed (need manual creation)

Validation complete!
```

---

## Updated Skills Summary

### Total Skills: 30 (was 25)

**Documentation Automation** (6):
1. scaffold-doc
2. generate-docs
3. validate-doc
4. migrate-doc
5. sync-configs
6. check-automation

**GitHub Project Management** (7):
7. setup-github-projects
8. setup-github-discussions
9. configure-branch-protection
10. setup-codeql
11. manage-labels
12. assign-reviewers
13. manage-milestones

**Dependency & Release** (3):
14. configure-dependabot
15. configure-release-please
16. update-dependencies

**Code Quality** (4):
17. run-linters
18. run-tests
19. format-code
20. check-licenses

**Infrastructure** (3):
21. manage-coder-workspace
22. manage-docker-config
23. manage-ci-workflows

**Monitoring** (2):
24. check-health
25. view-logs

**Branding & Assets** (5): 🆕
26. generate-placeholder-assets
27. generate-badges
28. update-branding
29. capture-screenshots (post-v1)
30. validate-assets

---

## Implementation Priority

### P0: Critical (before Phase 2)
- generate-placeholder-assets (needed for template system)
- generate-badges (needed for README)
- validate-assets (needed for validation pipeline)

### P1: High (Phase 2-6)
- update-branding (useful during development)

### P2: Post-v1
- capture-screenshots (only when UI exists)

---

## Integration with Phases

### Phase 2: Template System
- **Add**: generate-placeholder-assets skill
- **Add**: generate-badges skill
- **Add**: validate-assets skill
- Generate all placeholders before pilot migration

### Phase 5: Generation Pipeline
- Integrate validate-assets into validation pipeline

### Phase 14: Claude Code Skills
- Implement all 30 skills (including 5 branding skills)

### Post-v1: Screenshots
- Implement capture-screenshots skill
- Replace all placeholders with real screenshots

---

## Dependencies

**Python Packages** (add to requirements.txt):
```txt
# Existing
pyyaml>=6.0
jinja2>=3.1.5
yamale>=4.0

# NEW for branding/assets
Pillow>=10.0          # Image generation
cairosvg>=2.7         # SVG to PNG conversion
playwright>=1.40      # Screenshot capture (post-v1)
requests>=2.31        # Badge downloading
```

**System Dependencies**:
```bash
# Fonts for text rendering
sudo apt-get install fonts-dejavu fonts-liberation

# Playwright browsers (post-v1)
playwright install chromium
```

---

## Success Criteria

✅ All 30 skills specified
✅ All 5 branding skills integrated into plan
✅ Placeholder generation automated
✅ Badge generation automated
✅ Asset validation automated
✅ Screenshot capture planned (post-v1)
✅ Branding system complete

---

**Status**: Complete
**Total Skills**: 30 (5 added for branding/assets)
**Ready**: YES - All skills specified, ready to implement

