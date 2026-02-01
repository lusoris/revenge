#!/usr/bin/env python3
"""
Sync package versions from SOT to all design documents.

Extracts version information from 00_SOURCE_OF_TRUTH.md and updates
mentions of those packages in other design documents.

Usage:
    python3 scripts/sync-versions.py              # Check for drift
    python3 scripts/sync-versions.py --fix       # Update versions
    python3 scripts/sync-versions.py --report    # Generate report
"""

import argparse
import re
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOT_FILE = DESIGN_DIR / "00_SOURCE_OF_TRUTH.md"

# Version pattern: matches "vX.Y.Z" or "X.Y.Z"
VERSION_PATTERN = re.compile(r"v?(\d+\.\d+(?:\.\d+)?(?:-\w+)?)")


def extract_sot_versions() -> dict:
    """Extract package versions from SOT."""
    content = SOT_FILE.read_text(encoding="utf-8")
    versions = {}

    # Pattern for table rows with package and version
    # | `package/path` | vX.Y.Z | description |
    table_pattern = re.compile(
        r"\|\s*`([^`]+)`\s*\|\s*(v?\d+\.\d+(?:\.\d+)?(?:-\w+)?)\s*\|",
        re.MULTILINE,
    )

    for match in table_pattern.finditer(content):
        package = match.group(1).strip()
        version = match.group(2).strip()
        versions[package] = version

    # Also extract from inline mentions like "Package: vX.Y.Z"
    re.compile(
        r"(?:go\.uber\.org|github\.com|golang\.org)/[\w\-/]+[`\s]+[`]?(v?\d+\.\d+(?:\.\d+)?)",
    )

    return versions


def find_version_mentions(filepath: Path, sot_versions: dict) -> list:
    """Find version mentions in a file that may need updating."""
    content = filepath.read_text(encoding="utf-8")
    mentions = []

    for package, sot_version in sot_versions.items():
        # Look for the package name followed by a version
        # Pattern: package name + optional backticks/spaces + version
        pkg_escaped = re.escape(package)

        # Try different patterns
        patterns = [
            # `package` | vX.Y.Z |
            rf"`{pkg_escaped}`\s*\|\s*(v?\d+\.\d+(?:\.\d+)?)",
            # package vX.Y.Z
            rf"{pkg_escaped}\s+(v?\d+\.\d+(?:\.\d+)?)",
            # package: vX.Y.Z
            rf"{pkg_escaped}:\s*(v?\d+\.\d+(?:\.\d+)?)",
        ]

        for pattern in patterns:
            for match in re.finditer(pattern, content):
                found_version = match.group(1)
                if found_version != sot_version:
                    line_num = content[: match.start()].count("\n") + 1
                    mentions.append(
                        {
                            "package": package,
                            "line": line_num,
                            "found": found_version,
                            "expected": sot_version,
                            "match": match.group(0),
                        },
                    )

    return mentions


def fix_versions(filepath: Path, mentions: list) -> bool:
    """Fix version mentions in a file."""
    if not mentions:
        return False

    content = filepath.read_text(encoding="utf-8")
    original = content

    for mention in mentions:
        old = mention["match"]
        new = old.replace(mention["found"], mention["expected"])
        content = content.replace(old, new, 1)

    if content != original:
        filepath.write_text(content, encoding="utf-8")
        return True
    return False


def main():
    parser = argparse.ArgumentParser(description="Sync versions from SOT")
    parser.add_argument("--fix", action="store_true", help="Apply fixes")
    parser.add_argument("--report", action="store_true", help="Generate report")
    parser.add_argument("--verbose", "-v", action="store_true", help="Verbose output")
    args = parser.parse_args()

    print("Extracting versions from SOT...")
    sot_versions = extract_sot_versions()
    print(f"Found {len(sot_versions)} package versions in SOT")

    if args.verbose:
        for pkg, ver in sorted(sot_versions.items()):
            print(f"  {pkg}: {ver}")

    print("\nScanning design documents...")
    total_drift = 0
    files_with_drift = 0
    all_mentions = []

    for filepath in sorted(DESIGN_DIR.rglob("*.md")):
        if ".archive" in str(filepath) or ".analysis" in str(filepath):
            continue
        if filepath.name == "00_SOURCE_OF_TRUTH.md":
            continue

        mentions = find_version_mentions(filepath, sot_versions)

        if mentions:
            files_with_drift += 1
            total_drift += len(mentions)
            rel_path = filepath.relative_to(PROJECT_ROOT)

            for m in mentions:
                m["file"] = rel_path
                all_mentions.append(m)

            if args.verbose or not args.fix:
                print(f"\n{rel_path}:")
                for m in mentions:
                    print(f"  Line {m['line']}: {m['package']}")
                    print(f"    Found: {m['found']} → Expected: {m['expected']}")

            if args.fix and fix_versions(filepath, mentions):
                print(f"  Fixed {len(mentions)} version(s) in {rel_path}")

    print(
        f"\n{'Fixed' if args.fix else 'Found'} {total_drift} version drift(s) in {files_with_drift} file(s)",
    )

    if args.report:
        report_path = DESIGN_DIR / ".analysis" / "VERSION_DRIFT_REPORT.md"
        report_path.parent.mkdir(parents=True, exist_ok=True)

        lines = [
            "# Version Drift Report",
            "",
            f"Found {total_drift} version drift(s) in {files_with_drift} file(s)",
            "",
            "| File | Line | Package | Found | Expected |",
            "|------|------|---------|-------|----------|",
        ]

        for m in all_mentions:
            lines.append(
                f"| {m['file']} | {m['line']} | `{m['package']}` | {m['found']} | {m['expected']} |",
            )

        report_path.write_text("\n".join(lines), encoding="utf-8")
        print(f"\nReport saved to: {report_path}")

    if not args.fix and total_drift > 0:
        print("\nRun with --fix to update versions")


if __name__ == "__main__":
    main()
