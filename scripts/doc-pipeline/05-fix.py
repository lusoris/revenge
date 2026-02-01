#!/usr/bin/env python3
"""Fix broken links and placeholders in design documentation.

Combines link fixing and placeholder fixing:
1. Auto-fixes common link patterns with high confidence
2. Flags ambiguous cases for manual review
3. Replaces common placeholder patterns

Usage:
    python scripts/doc-pipeline/05-fix.py               # Dry run
    python scripts/doc-pipeline/05-fix.py --apply       # Apply fixes
    python scripts/doc-pipeline/05-fix.py --report      # Generate report
"""

from __future__ import annotations

import argparse
import re
import sys
from datetime import datetime
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOURCES_DIR = PROJECT_ROOT / "docs" / "dev" / "sources"
DOCS_DIR = PROJECT_ROOT / "docs" / "dev"

LINK_PATTERN = re.compile(r"\[([^\]]+)\]\(([^)]+)\)")
YAML_SOURCE_PATTERN = re.compile(r"^\s+url:\s+(.+)$", re.MULTILINE)

KNOWN_LOCATIONS = {
    "TECH_STACK.md": "technical/TECH_STACK.md",
    "API.md": "technical/API.md",
    "FRONTEND.md": "technical/FRONTEND.md",
    "CONFIGURATION.md": "technical/CONFIGURATION.md",
    "AUDIO_STREAMING.md": "technical/AUDIO_STREAMING.md",
    "OFFLOADING.md": "technical/OFFLOADING.md",
    "01_ARCHITECTURE.md": "architecture/01_ARCHITECTURE.md",
    "02_DESIGN_PRINCIPLES.md": "architecture/02_DESIGN_PRINCIPLES.md",
    "03_METADATA_SYSTEM.md": "architecture/03_METADATA_SYSTEM.md",
    "04_PLAYER_ARCHITECTURE.md": "architecture/04_PLAYER_ARCHITECTURE.md",
    "05_PLUGIN_ARCHITECTURE_DECISION.md": "architecture/05_PLUGIN_ARCHITECTURE_DECISION.md",
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
    "SETUP.md": "operations/SETUP.md",
    "DEVELOPMENT.md": "operations/DEVELOPMENT.md",
    "BEST_PRACTICES.md": "operations/BEST_PRACTICES.md",
    "GITFLOW.md": "operations/GITFLOW.md",
    "REVERSE_PROXY.md": "operations/REVERSE_PROXY.md",
    "BRANCH_PROTECTION.md": "operations/BRANCH_PROTECTION.md",
    "DATABASE_AUTO_HEALING.md": "operations/DATABASE_AUTO_HEALING.md",
    "00_SOURCE_OF_TRUTH.md": "00_SOURCE_OF_TRUTH.md",
}

PLACEHOLDER_PATTERNS = [
    (r"\[TBD\]", "[To be determined]"),
    (r"\[TODO\]", "[Not yet implemented]"),
    (r"\[placeholder\]", "[Content pending]"),
]


def get_depth_from_design(filepath: Path) -> int:
    """Calculate directory depth from design/."""
    try:
        rel = filepath.relative_to(DESIGN_DIR)
        return len(rel.parts) - 1
    except ValueError:
        return 0


def is_internal_link(link: str) -> bool:
    """Check if a link is internal."""
    return not link.startswith(("http://", "https://", "mailto:", "#"))


def resolve_link(source_file: Path, link: str) -> Path:
    """Resolve a relative link to an absolute path."""
    link_path = link.split("#")[0]
    if not link_path:
        return source_file
    source_dir = source_file.parent
    return (source_dir / link_path).resolve()


def find_correct_path(filename: str, source_file: Path) -> str | None:
    """Try to find the correct relative path for a file."""
    try:
        source_file.relative_to(DESIGN_DIR)
    except ValueError:
        return None

    if filename in KNOWN_LOCATIONS:
        target = DESIGN_DIR / KNOWN_LOCATIONS[filename]
        if target.exists():
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
        content = filepath.read_text(encoding="utf-8")
    except Exception as e:
        return {"error": str(e), "links": [], "placeholders": []}

    links = []
    placeholders = []

    # Check markdown links
    for match in LINK_PATTERN.finditer(content):
        text = match.group(1)
        link = match.group(2)
        start = match.start()
        line_num = content[:start].count("\n") + 1

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
            "type": "markdown",
        }

        if not exists:
            filename = Path(link.split("#")[0]).name
            if filename:
                suggestion = find_correct_path(filename, filepath)
                if suggestion:
                    link_info["suggestion"] = suggestion
                    link_info["confidence"] = (
                        "high" if filename in KNOWN_LOCATIONS else "medium"
                    )

        links.append(link_info)

    # Check YAML source url: fields
    for match in YAML_SOURCE_PATTERN.finditer(content):
        link = match.group(1).strip()
        start = match.start()
        line_num = content[:start].count("\n") + 1

        if not is_internal_link(link):
            continue

        resolved = resolve_link(filepath, link)
        exists = resolved.exists()

        link_info = {
            "text": f"YAML source: {link}",
            "link": link,
            "line": line_num,
            "exists": exists,
            "resolved": str(resolved),
            "suggestion": None,
            "confidence": None,
            "type": "yaml",
        }

        # For YAML sources, we don't auto-suggest fixes since the structure is more rigid
        if not exists:
            link_info["suggestion"] = "-"
            link_info["confidence"] = "manual"

        links.append(link_info)

    # Check for placeholders
    for pattern, replacement in PLACEHOLDER_PATTERNS:
        for match in re.finditer(pattern, content, re.IGNORECASE):
            placeholders.append(
                {
                    "pattern": pattern,
                    "replacement": replacement,
                    "line": content[: match.start()].count("\n") + 1,
                },
            )

    return {
        "filepath": filepath,
        "links": links,
        "placeholders": placeholders,
        "total": len(links),
        "broken": len([link for link in links if not link["exists"]]),
        "fixable": len(
            [
                link
                for link in links
                if link["suggestion"] and link["confidence"] == "high"
            ],
        ),
        "placeholder_count": len(placeholders),
    }


def fix_file(filepath: Path, analysis: dict, *, dry_run: bool = True) -> list:
    """Apply fixes to a file."""
    fixes = []

    try:
        content = filepath.read_text(encoding="utf-8")
    except Exception:
        return fixes

    # Fix broken links
    for link_info in analysis["links"]:
        if link_info["exists"]:
            continue

        if link_info["suggestion"] and link_info["confidence"] == "high":
            old_link = link_info["link"]
            new_link = link_info["suggestion"]

            if "#" in old_link:
                anchor = "#" + old_link.split("#")[1]
                new_link = new_link + anchor

            old_pattern = f"]({old_link})"
            new_pattern = f"]({new_link})"
            content = content.replace(old_pattern, new_pattern)

            fixes.append(
                {
                    "type": "link",
                    "line": link_info["line"],
                    "old": old_link,
                    "new": new_link,
                    "text": link_info["text"],
                },
            )

    # Fix placeholders
    for placeholder in analysis["placeholders"]:
        content = re.sub(
            placeholder["pattern"],
            placeholder["replacement"],
            content,
            flags=re.IGNORECASE,
        )
        fixes.append(
            {
                "type": "placeholder",
                "line": placeholder["line"],
                "old": placeholder["pattern"],
                "new": placeholder["replacement"],
            },
        )

    if fixes and not dry_run:
        filepath.write_text(content, encoding="utf-8")

    return fixes


def generate_report(results: list) -> str:
    """Generate a detailed report."""
    lines = [
        "# Link and Placeholder Analysis Report",
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
    placeholder_count = sum(r.get("placeholder_count", 0) for r in results)
    files_with_broken = len([r for r in results if r["broken"] > 0])

    lines.extend(
        [
            f"- **Total internal links**: {total_links}",
            f"- **Broken links**: {broken_links}",
            f"- **Auto-fixable (high confidence)**: {fixable_links}",
            f"- **Files with broken links**: {files_with_broken}",
            f"- **Placeholder patterns found**: {placeholder_count}",
            "",
            "---",
            "",
            "## Broken Links by File",
            "",
        ],
    )

    for result in sorted(results, key=lambda x: -x["broken"]):
        if result["broken"] == 0:
            continue

        rel_path = result["filepath"].relative_to(PROJECT_ROOT)
        lines.extend(
            [
                f"### {rel_path}",
                "",
                "| Line | Link | Status | Suggestion |",
                "|------|------|--------|------------|",
            ],
        )

        for link in result["links"]:
            if link["exists"]:
                continue

            status = (
                "Auto-fix"
                if link["confidence"] == "high"
                else ("Review" if link["suggestion"] else "Manual")
            )
            suggestion = link["suggestion"] or "-"

            lines.append(
                f"| {link['line']} | `{link['link']}` | {status} | `{suggestion}` |",
            )

        lines.append("")

    lines.extend(
        [
            "---",
            "",
            "*Report generated by `scripts/doc-pipeline/05-fix.py`*",
        ],
    )

    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Fix broken links and placeholders in documentation",
    )
    parser.add_argument(
        "--apply",
        action="store_true",
        help="Apply fixes (default: dry-run)",
    )
    parser.add_argument(
        "--report",
        action="store_true",
        help="Generate detailed report file",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Show all links, not just broken ones",
    )
    parser.add_argument(
        "--path",
        type=Path,
        default=DESIGN_DIR,
        help="Path to scan (default: docs/dev/design)",
    )
    args = parser.parse_args()

    dry_run = not args.apply

    if dry_run:
        print("=== DRY RUN (use --apply to write) ===\n")

    target = args.path
    if not target.is_absolute():
        target = PROJECT_ROOT / target
    target = target.resolve()

    print(f"Scanning: {target}")
    print()

    results = []
    md_files = list(target.rglob("*.md"))

    for filepath in md_files:
        if ".archive" in str(filepath) or ".analysis" in str(filepath):
            continue

        result = analyze_file(filepath)
        results.append(result)

    total_links = sum(r["total"] for r in results)
    broken_links = sum(r["broken"] for r in results)
    fixable_links = sum(r["fixable"] for r in results)
    placeholder_count = sum(r.get("placeholder_count", 0) for r in results)

    print(f"Files scanned: {len(results)}")
    print(f"Total internal links: {total_links}")
    print(f"Broken links: {broken_links}")
    print(f"Auto-fixable (high confidence): {fixable_links}")
    print(f"Placeholders found: {placeholder_count}")
    print()

    if broken_links == 0 and placeholder_count == 0:
        print("✓ No issues found!")
        return 0

    total_fixes = 0
    for result in results:
        if result["fixable"] > 0 or result.get("placeholder_count", 0) > 0:
            fixes = fix_file(result["filepath"], result, dry_run=dry_run)
            if fixes:
                rel_path = result["filepath"].relative_to(PROJECT_ROOT)
                action = "Fixed" if not dry_run else "Would fix"
                print(f"  {action} {len(fixes)} issues in {rel_path}")
                total_fixes += len(fixes)

    if total_fixes > 0:
        action = "Applied" if not dry_run else "Would apply"
        print(f"\n{action} {total_fixes} fixes")

    if args.report:
        report = generate_report(results)
        report_path = (
            PROJECT_ROOT / "docs" / "dev" / "design" / ".analysis" / "FIXES_REPORT.md"
        )
        report_path.parent.mkdir(parents=True, exist_ok=True)
        report_path.write_text(report, encoding="utf-8")
        print(f"\nReport saved to: {report_path}")

    if dry_run and total_fixes > 0:
        print("\n=== DRY RUN complete. Use --apply to write. ===")

    return 0


if __name__ == "__main__":
    sys.exit(main())
