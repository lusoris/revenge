#!/usr/bin/env python3
"""Fix design_refs paths in YAML data files.

This script fixes incorrect relative paths in design_refs sections of YAML data files.

Common issues:
- architecture/INDEX.md → INDEX.md (same directory)
- services/INDEX.md → INDEX.md (same directory)
- operations/INDEX.md → INDEX.md (same directory)

Author: Automation System
Created: 2026-02-01
"""

import sys
from pathlib import Path

import yaml


def fix_design_refs(yaml_file: Path) -> bool:
    """Fix design_refs paths in a YAML file.

    Args:
        yaml_file: Path to YAML data file

    Returns:
        True if file was modified, False otherwise
    """
    with open(yaml_file) as f:
        data = yaml.safe_load(f)

    if not data or "design_refs" not in data:
        return False

    modified = False

    # Get the category from the file path
    # e.g., data/architecture/FILE.yaml → architecture
    # e.g., data/services/FILE.yaml → services
    parts = yaml_file.parts
    if "data" in parts:
        data_index = parts.index("data")
        category = parts[data_index + 1] if data_index + 1 < len(parts) else None
    else:
        category = None

    for ref in data["design_refs"]:
        if "path" not in ref:
            continue

        path = ref["path"]

        # Fix same-directory INDEX.md references
        # e.g., architecture/INDEX.md → INDEX.md
        if category and path == f"{category}/INDEX.md":
            ref["path"] = "INDEX.md"
            modified = True
            print(f"  Fixed: {path} → INDEX.md")

        # Fix same-directory file references
        # e.g., services/INDEX.md → INDEX.md
        elif "/" in path and path.endswith("/INDEX.md"):
            prefix = path.rsplit("/", 1)[0]
            if category and prefix == category:
                ref["path"] = "INDEX.md"
                modified = True
                print(f"  Fixed: {path} → INDEX.md")

    if modified:
        # Write back with preserved formatting
        with open(yaml_file, "w") as f:
            yaml.dump(
                data, f, default_flow_style=False, sort_keys=False, allow_unicode=True
            )

    return modified


def main():
    """Main entry point."""
    repo_root = Path(__file__).parent.parent
    data_dir = repo_root / "data"

    if not data_dir.exists():
        print(f"❌ Error: {data_dir} not found")
        sys.exit(1)

    print(f"🔍 Scanning {data_dir} for YAML files...")

    yaml_files = list(data_dir.glob("**/*.yaml"))
    yaml_files = [f for f in yaml_files if f.name != "shared-sot.yaml"]

    print(f"📁 Found {len(yaml_files)} YAML files")

    modified_count = 0

    for yaml_file in sorted(yaml_files):
        print(f"\n📄 {yaml_file.relative_to(repo_root)}")
        if fix_design_refs(yaml_file):
            modified_count += 1

    print(f"\n✅ Fixed {modified_count} files")


if __name__ == "__main__":
    main()
