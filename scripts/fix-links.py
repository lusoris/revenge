#!/usr/bin/env python3
"""
Hybrid link fixer for design documentation.

This script:
1. Scans all markdown files for internal links
2. Auto-fixes common link patterns with high confidence
3. Flags ambiguous cases for manual review
4. Generates a detailed report

Usage:
    python3 scripts/fix-links.py [--fix] [--report]

Options:
    --fix       Apply auto-fixes (default: dry-run)
    --report    Generate detailed report file
    --verbose   Show all links, not just broken ones
"""

import argparse
import re
from collections import defaultdict
from datetime import datetime
from pathlib import Path

# Project paths
SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"
DOCS_DIR = PROJECT_ROOT / "docs" / "dev"

# Link pattern: [text](path)
LINK_PATTERN = re.compile(r'\[([^\]]+)\]\(([^)]+)\)')

# Known file locations for auto-fixing
KNOWN_LOCATIONS = {
    # technical/
    "TECH_STACK.md": "technical/TECH_STACK.md",
    "API.md": "technical/API.md",
    "FRONTEND.md": "technical/FRONTEND.md",
    "CONFIGURATION.md": "technical/CONFIGURATION.md",
    "AUDIO_STREAMING.md": "technical/AUDIO_STREAMING.md",
    "OFFLOADING.md": "technical/OFFLOADING.md",

    # architecture/
    "01_ARCHITECTURE.md": "architecture/01_ARCHITECTURE.md",
    "02_DESIGN_PRINCIPLES.md": "architecture/02_DESIGN_PRINCIPLES.md",
    "03_METADATA_SYSTEM.md": "architecture/03_METADATA_SYSTEM.md",
    "04_PLAYER_ARCHITECTURE.md": "architecture/04_PLAYER_ARCHITECTURE.md",
    "05_PLUGIN_ARCHITECTURE_DECISION.md": "architecture/05_PLUGIN_ARCHITECTURE_DECISION.md",

    # services/
    "AUTH.md": "services/AUTH.md",
    "USER.md": "services/USER.md",
    "SESSION.md": "services/SESSION.md",
    "RBAC.md": "services/RBAC.md",
    "ACTIVITY.md": "services/ACTIVITY.md",
    "SETTINGS.md": "services/SETTINGS.md",
    "APIKEYS.md": "services/APIKEYS.md",
    "OIDC.md": "services/OIDC.md",
    "GRANTS.md": "services/GRANTS.md",
    "FINGERPRINT.md": "services/FINGERPRINT.md",
    "LIBRARY.md": "services/LIBRARY.md",
    "METADATA.md": "services/METADATA.md",
    "SEARCH.md": "services/SEARCH.md",
    "ANALYTICS.md": "services/ANALYTICS.md",
    "NOTIFICATION.md": "services/NOTIFICATION.md",

    # operations/
    "SETUP.md": "operations/SETUP.md",
    "DEVELOPMENT.md": "operations/DEVELOPMENT.md",
    "BEST_PRACTICES.md": "operations/BEST_PRACTICES.md",
    "GITFLOW.md": "operations/GITFLOW.md",
    "REVERSE_PROXY.md": "operations/REVERSE_PROXY.md",
    "BRANCH_PROTECTION.md": "operations/BRANCH_PROTECTION.md",
    "DATABASE_AUTO_HEALING.md": "operations/DATABASE_AUTO_HEALING.md",

    # Root design docs
    "00_SOURCE_OF_TRUTH.md": "00_SOURCE_OF_TRUTH.md",
}

# Patterns that indicate relative path issues
FIX_PATTERNS = [
    # Missing ../ for going up from features/* to integrations/*
    (r'^\(integrations/', lambda m, depth: f"({'../' * depth}integrations/"),
    # Missing ../ for going up from features/* to services/*
    (r'^\(services/', lambda m, depth: f"({'../' * depth}services/"),
    # Missing ../ for going up from features/* to technical/*
    (r'^\(technical/', lambda m, depth: f"({'../' * depth}technical/"),
    # Missing ../ for going up from features/* to architecture/*
    (r'^\(architecture/', lambda m, depth: f"({'../' * depth}architecture/"),
    # Missing ../ for going up from integrations/* to features/*
    (r'^\(features/', lambda m, depth: f"({'../' * depth}features/"),
]


def get_depth_from_design(filepath: Path) -> int:
    """Calculate how many directories deep a file is from design/."""
    try:
        rel = filepath.relative_to(DESIGN_DIR)
        return len(rel.parts) - 1  # Subtract 1 for the filename itself
    except ValueError:
        return 0


def is_internal_link(link: str) -> bool:
    """Check if a link is internal (not URL, not anchor-only)."""
    if link.startswith(('http://', 'https://', 'mailto:', '#')):
        return False
    return True


def resolve_link(source_file: Path, link: str) -> Path:
    """Resolve a relative link to an absolute path."""
    # Remove anchor
    link_path = link.split('#')[0]
    if not link_path:
        return source_file  # Anchor-only link

    source_dir = source_file.parent
    return (source_dir / link_path).resolve()


def find_correct_path(filename: str, source_file: Path) -> str | None:
    """Try to find the correct relative path for a file."""
    # Only process files within design directory
    try:
        source_file.relative_to(DESIGN_DIR)
    except ValueError:
        # Source file is not in design directory, skip
        return None

    # Check known locations
    if filename in KNOWN_LOCATIONS:
        target = DESIGN_DIR / KNOWN_LOCATIONS[filename]
        if target.exists():
            try:
                return str(target.relative_to(source_file.parent))
            except ValueError:
                # Need to go up directories
                try:
                    source_rel = source_file.relative_to(DESIGN_DIR)
                    target_rel = target.relative_to(DESIGN_DIR)
                    up_count = len(source_rel.parts) - 1
                    return "../" * up_count + str(target_rel)
                except ValueError:
                    return None

    # Search in design directory (exclude .archive)
    matches = [m for m in DESIGN_DIR.rglob(filename) if ".archive" not in str(m)]
    if len(matches) == 1:
        target = matches[0]
        try:
            return str(target.relative_to(source_file.parent))
        except ValueError:
            try:
                source_rel = source_file.relative_to(DESIGN_DIR)
                target_rel = target.relative_to(DESIGN_DIR)
                up_count = len(source_rel.parts) - 1
                return "../" * up_count + str(target_rel)
            except ValueError:
                return None

    return None


def analyze_file(filepath: Path) -> dict:
    """Analyze a markdown file for broken links."""
    try:
        content = filepath.read_text(encoding='utf-8')
    except Exception as e:
        return {"error": str(e), "links": []}

    links = []
    depth = get_depth_from_design(filepath)

    for match in LINK_PATTERN.finditer(content):
        text = match.group(1)
        link = match.group(2)
        start = match.start()
        line_num = content[:start].count('\n') + 1

        if not is_internal_link(link):
            continue

        resolved = resolve_link(filepath, link)
        exists = resolved.exists()

        link_info = {
            "text": text,
            "link": link,
            "line": line_num,
            "exists": exists,
            "resolved": str(resolved),
            "suggestion": None,
            "confidence": None,
        }

        if not exists:
            # Try to find correct path
            filename = Path(link.split('#')[0]).name
            if filename:
                suggestion = find_correct_path(filename, filepath)
                if suggestion:
                    link_info["suggestion"] = suggestion
                    link_info["confidence"] = "high" if filename in KNOWN_LOCATIONS else "medium"

        links.append(link_info)

    return {
        "filepath": filepath,
        "depth": depth,
        "links": links,
        "total": len(links),
        "broken": len([l for l in links if not l["exists"]]),
        "fixable": len([l for l in links if l["suggestion"] and l["confidence"] == "high"]),
    }


def fix_file(filepath: Path, analysis: dict, dry_run: bool = True) -> list:
    """Apply fixes to a file."""
    fixes = []

    try:
        content = filepath.read_text(encoding='utf-8')
        original = content
    except Exception:
        return fixes

    for link_info in analysis["links"]:
        if link_info["exists"]:
            continue

        if link_info["suggestion"] and link_info["confidence"] == "high":
            old_link = link_info["link"]
            new_link = link_info["suggestion"]

            # Preserve anchor
            if '#' in old_link:
                anchor = '#' + old_link.split('#')[1]
                new_link = new_link + anchor

            # Replace in content
            old_pattern = f"]({old_link})"
            new_pattern = f"]({new_link})"
            content = content.replace(old_pattern, new_pattern)

            fixes.append({
                "line": link_info["line"],
                "old": old_link,
                "new": new_link,
                "text": link_info["text"],
            })

    if fixes and not dry_run:
        filepath.write_text(content, encoding='utf-8')

    return fixes


def generate_report(results: list) -> str:
    """Generate a detailed report."""
    lines = [
        "# Link Analysis Report",
        "",
        f"> Generated: {datetime.now().strftime('%Y-%m-%d %H:%M')}",
        "",
        "---",
        "",
        "## Summary",
        "",
    ]

    total_links = sum(r["total"] for r in results)
    broken_links = sum(r["broken"] for r in results)
    fixable_links = sum(r["fixable"] for r in results)
    files_with_broken = len([r for r in results if r["broken"] > 0])

    lines.extend([
        f"- **Total internal links**: {total_links}",
        f"- **Broken links**: {broken_links}",
        f"- **Auto-fixable (high confidence)**: {fixable_links}",
        f"- **Files with broken links**: {files_with_broken}",
        "",
        "---",
        "",
        "## Broken Links by File",
        "",
    ])

    for result in sorted(results, key=lambda x: -x["broken"]):
        if result["broken"] == 0:
            continue

        rel_path = result["filepath"].relative_to(PROJECT_ROOT)
        lines.extend([
            f"### {rel_path}",
            "",
            "| Line | Link | Status | Suggestion |",
            "|------|------|--------|------------|",
        ])

        for link in result["links"]:
            if link["exists"]:
                continue

            status = "Auto-fix" if link["confidence"] == "high" else (
                "Review" if link["suggestion"] else "Manual"
            )
            suggestion = link["suggestion"] or "-"

            lines.append(
                f"| {link['line']} | `{link['link']}` | {status} | `{suggestion}` |"
            )

        lines.append("")

    # Manual review section
    lines.extend([
        "---",
        "",
        "## Manual Review Required",
        "",
        "These links need manual investigation:",
        "",
    ])

    manual_count = 0
    for result in results:
        for link in result["links"]:
            if not link["exists"] and not link["suggestion"]:
                rel_path = result["filepath"].relative_to(PROJECT_ROOT)
                lines.append(f"- `{rel_path}:{link['line']}` → `{link['link']}`")
                manual_count += 1

    if manual_count == 0:
        lines.append("*None - all broken links have suggestions*")

    lines.extend([
        "",
        "---",
        "",
        "*Report generated by `scripts/fix-links.py`*",
    ])

    return "\n".join(lines)


def main():
    parser = argparse.ArgumentParser(
        description="Hybrid link fixer for design documentation"
    )
    parser.add_argument(
        "--fix",
        action="store_true",
        help="Apply auto-fixes (default: dry-run)"
    )
    parser.add_argument(
        "--report",
        action="store_true",
        help="Generate detailed report file"
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Show all links, not just broken ones"
    )
    parser.add_argument(
        "--path",
        type=Path,
        default=DESIGN_DIR,
        help="Path to scan (default: docs/dev/design)"
    )

    args = parser.parse_args()

    # Resolve path
    target = args.path
    if not target.is_absolute():
        target = PROJECT_ROOT / target
    target = target.resolve()

    print(f"Scanning: {target}")
    print()

    # Analyze all markdown files
    results = []
    md_files = list(target.rglob("*.md"))

    for filepath in md_files:
        # Skip archive directories
        if ".archive" in str(filepath) or ".analysis" in str(filepath):
            continue

        result = analyze_file(filepath)
        results.append(result)

    # Summary
    total_links = sum(r["total"] for r in results)
    broken_links = sum(r["broken"] for r in results)
    fixable_links = sum(r["fixable"] for r in results)

    print(f"Files scanned: {len(results)}")
    print(f"Total internal links: {total_links}")
    print(f"Broken links: {broken_links}")
    print(f"Auto-fixable (high confidence): {fixable_links}")
    print()

    if broken_links == 0:
        print("No broken links found!")
        return

    # Apply fixes if requested
    if args.fix:
        print("Applying auto-fixes...")
        total_fixes = 0
        for result in results:
            if result["fixable"] > 0:
                fixes = fix_file(result["filepath"], result, dry_run=False)
                if fixes:
                    rel_path = result["filepath"].relative_to(PROJECT_ROOT)
                    print(f"  Fixed {len(fixes)} links in {rel_path}")
                    total_fixes += len(fixes)

        print(f"\nApplied {total_fixes} fixes")
    else:
        print("Dry-run mode. Use --fix to apply changes.")

        # Show what would be fixed
        print("\nWould fix:")
        for result in results:
            for link in result["links"]:
                if link["suggestion"] and link["confidence"] == "high":
                    rel_path = result["filepath"].relative_to(PROJECT_ROOT)
                    print(f"  {rel_path}:{link['line']}")
                    print(f"    {link['link']} → {link['suggestion']}")

    # Generate report if requested
    if args.report:
        report = generate_report(results)
        report_path = PROJECT_ROOT / "docs" / "dev" / "design" / ".analysis" / "LINK_FIXES_REPORT.md"
        report_path.parent.mkdir(parents=True, exist_ok=True)
        report_path.write_text(report, encoding='utf-8')
        print(f"\nReport saved to: {report_path}")


if __name__ == "__main__":
    main()
