#!/usr/bin/env python3
"""Sync package versions across documentation.

Ensures version consistency between:
- 00_SOURCE_OF_TRUTH.md (master)
- go.mod
- package.json
- SOURCES.yaml

Usage:
    python scripts/utils/sync-versions.py           # Check versions
    python scripts/utils/sync-versions.py --apply   # Update SOT from go.mod
"""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
SOT_FILE = DESIGN_DIR / "00_SOURCE_OF_TRUTH.md"
GO_MOD = PROJECT_ROOT / "go.mod"
PACKAGE_JSON = PROJECT_ROOT / "web" / "package.json"


def parse_go_mod() -> dict[str, str]:
    """Parse go.mod for package versions."""
    if not GO_MOD.exists():
        return {}

    versions = {}
    content = GO_MOD.read_text(encoding="utf-8")

    # Match require block
    require_match = re.search(r"require\s*\((.*?)\)", content, re.DOTALL)
    if require_match:
        for line in require_match.group(1).split("\n"):
            line = line.strip()
            if not line or line.startswith("//"):
                continue
            parts = line.split()
            if len(parts) >= 2:
                pkg = parts[0]
                version = parts[1]
                # Extract package name
                pkg_name = pkg.split("/")[-1]
                versions[pkg_name] = version

    # Match single require statements
    for match in re.finditer(r"require\s+(\S+)\s+(\S+)", content):
        pkg = match.group(1)
        version = match.group(2)
        pkg_name = pkg.split("/")[-1]
        versions[pkg_name] = version

    return versions


def parse_package_json() -> dict[str, str]:
    """Parse package.json for frontend dependencies."""
    if not PACKAGE_JSON.exists():
        return {}

    import json

    with open(PACKAGE_JSON, encoding="utf-8") as f:
        data = json.load(f)

    versions = {}
    for section in ["dependencies", "devDependencies"]:
        for pkg, version in data.get(section, {}).items():
            # Remove version prefixes
            version = version.lstrip("^~")
            versions[pkg] = version

    return versions


def parse_sot_versions() -> dict[str, str]:
    """Parse versions from SOT file."""
    if not SOT_FILE.exists():
        return {}

    content = SOT_FILE.read_text(encoding="utf-8")
    versions = {}

    # Match version table rows: | package | version |
    pattern = re.compile(r"\|\s*`?([^|`]+)`?\s*\|\s*`?([v\d][^|`]*)`?\s*\|")
    for match in pattern.finditer(content):
        pkg = match.group(1).strip()
        version = match.group(2).strip()
        versions[pkg] = version

    return versions


def compare_versions(
    go_versions: dict[str, str],
    pkg_versions: dict[str, str],
    sot_versions: dict[str, str],
) -> list[dict]:
    """Compare versions and find discrepancies."""
    discrepancies = []

    # Check Go packages
    for pkg, go_version in go_versions.items():
        sot_version = sot_versions.get(pkg)
        if sot_version and sot_version != go_version:
            discrepancies.append(
                {
                    "package": pkg,
                    "source": "go.mod",
                    "actual": go_version,
                    "sot": sot_version,
                },
            )

    # Check frontend packages
    for pkg, pkg_version in pkg_versions.items():
        sot_version = sot_versions.get(pkg)
        if sot_version and sot_version != pkg_version:
            discrepancies.append(
                {
                    "package": pkg,
                    "source": "package.json",
                    "actual": pkg_version,
                    "sot": sot_version,
                },
            )

    return discrepancies


def main() -> int:
    parser = argparse.ArgumentParser(description="Sync package versions")
    parser.add_argument(
        "--apply",
        action="store_true",
        help="Update SOT with actual versions (not implemented yet)",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Show all versions",
    )
    args = parser.parse_args()

    print("Parsing version sources...")

    go_versions = parse_go_mod()
    print(f"  go.mod: {len(go_versions)} packages")

    pkg_versions = parse_package_json()
    print(f"  package.json: {len(pkg_versions)} packages")

    sot_versions = parse_sot_versions()
    print(f"  SOT: {len(sot_versions)} packages")

    if args.verbose:
        print("\n=== Go Packages ===")
        for pkg, version in sorted(go_versions.items()):
            print(f"  {pkg}: {version}")

        print("\n=== Frontend Packages ===")
        for pkg, version in sorted(pkg_versions.items()):
            print(f"  {pkg}: {version}")

    print("\nComparing versions...")
    discrepancies = compare_versions(go_versions, pkg_versions, sot_versions)

    if not discrepancies:
        print("✓ All versions in sync!")
        return 0

    print(f"\n⚠️  Found {len(discrepancies)} version discrepancies:\n")
    print("| Package | Source | Actual | SOT |")
    print("|---------|--------|--------|-----|")
    for d in discrepancies:
        print(f"| {d['package']} | {d['source']} | {d['actual']} | {d['sot']} |")

    if args.apply:
        print("\n--apply not yet implemented. Please update SOT manually.")

    return 1


if __name__ == "__main__":
    sys.exit(main())
