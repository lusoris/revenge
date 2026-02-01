#!/usr/bin/env python3
"""Fix design_refs paths in YAML data files comprehensively.

This script fixes all incorrect relative paths in design_refs sections by:
1. Detecting same-directory INDEX.md references
2. Fixing cross-category references (e.g., from architecture to features)
3. Ensuring proper relative path depth

Author: Automation System
Created: 2026-02-01
"""

import sys
from pathlib import Path

import yaml


def fix_design_refs(yaml_file: Path, repo_root: Path) -> int:
    """Fix design_refs paths in a YAML file.

    Args:
        yaml_file: Path to YAML data file
        repo_root: Repository root path

    Returns:
        Number of fixes applied
    """
    with open(yaml_file) as f:
        data = yaml.safe_load(f)

    if not data or "design_refs" not in data:
        return 0

    # Get the output directory this YAML will generate to
    # e.g., data/architecture/FILE.yaml → docs/dev/design/architecture/FILE.md
    parts = yaml_file.relative_to(repo_root).parts
    if parts[0] != "data":
        return 0

    # Output path relative to docs/dev/design/
    # e.g., data/architecture/FILE.yaml → architecture/
    # e.g., data/features/adult/FILE.yaml → features/adult/
    output_subdir = "/".join(parts[1:-1])  # Skip 'data' and filename

    fixes_applied = 0
    refs_to_remove = []

    for ref in data["design_refs"]:
        if "path" not in ref:
            continue

        path = ref["path"]
        original_path = path

        # Remove self-referential category summary links
        # e.g., features/adult.md from features/adult/ (depth=2)
        # e.g., integrations/metadata/adult.md from integrations/metadata/adult/ (depth=3)
        # These category summary files don't exist and shouldn't be self-referenced
        if output_subdir:
            depth = output_subdir.count("/") + 1
            expected_prefix = "../" * depth
            self_ref_path = f"{expected_prefix}{output_subdir}.md"
            if path == self_ref_path:
                refs_to_remove.append(ref)
                fixes_applied += 1
                print(
                    f"    Removed: {original_path} (self-referential category summary)"
                )
                continue

        # Fix same-directory INDEX.md references
        # e.g., architecture/INDEX.md → INDEX.md (when file is in architecture/)
        if "/" in path and path.endswith("/INDEX.md"):
            prefix = path.rsplit("/", 1)[0]
            if output_subdir and prefix == output_subdir:
                ref["path"] = "INDEX.md"
                fixes_applied += 1
                print(f"    Fixed: {original_path} → INDEX.md (same directory)")
                continue

        # Fix bare filenames that should have path prefix
        # e.g., from architecture/, path "ADULT_CONTENT_SYSTEM.md" should be "../features/adult/ADULT_CONTENT_SYSTEM.md"
        # e.g., from features/playback/, path "ADULT_CONTENT_SYSTEM.md" should be "../../features/adult/ADULT_CONTENT_SYSTEM.md"
        if "/" not in path and path.endswith(".md") and path != "INDEX.md":
            # Try to find where this file actually exists
            target_found = False

            # Calculate depth for proper ../ prefix
            depth = output_subdir.count("/") + 1 if output_subdir else 0
            prefix = "../" * depth

            # Check if this looks like adult content or data reconciliation
            if (
                "ADULT" in path.upper()
                or "QAR" in path.upper()
                or "WHISPARR" in path.upper()
                or "DATA_RECONCILIATION" in path.upper()
            ):
                ref["path"] = f"{prefix}features/adult/{path}"
                fixes_applied += 1
                print(
                    f"    Fixed: {original_path} → {prefix}features/adult/{path} (depth={depth})"
                )
                target_found = True

            # Check for root-level files referenced from subdirectories
            # e.g., "00_SOURCE_OF_TRUTH.md" from technical/ → "../00_SOURCE_OF_TRUTH.md"
            if (
                not target_found
                and "/" not in path
                and path != "INDEX.md"
                and output_subdir
            ):
                # Special files that exist at root level
                root_files = ["00_SOURCE_OF_TRUTH.md", "NAVIGATION.md"]
                if path in root_files:
                    ref["path"] = f"../{path}"
                    fixes_applied += 1
                    print(f"    Fixed: {original_path} → ../{path} (root-level file)")
                    target_found = True

            # Check if this is a same-directory category reference
            # e.g., "technical.md" from technical/ → should be "INDEX.md"
            if (
                not target_found
                and "/" not in path
                and path != "INDEX.md"
                and output_subdir
            ):
                # Get the last part of output_subdir
                category_name = output_subdir.split("/")[-1]
                if path == f"{category_name}.md":
                    ref["path"] = "INDEX.md"
                    fixes_applied += 1
                    print(f"    Fixed: {original_path} → INDEX.md (category self-ref)")
                    target_found = True

            if not target_found and "/" not in path and path != "INDEX.md":
                # Keep as-is for now, might be valid same-directory reference
                pass

        # Fix cross-directory references with wrong depth
        # Calculate correct depth: architecture/ = 1, features/adult/ = 2, etc.
        # e.g., from features/adult/, "architecture/X.md" should be "../../architecture/X.md"
        # e.g., from features/adult/, "../architecture/X.md" should be "../../architecture/X.md"
        if "/" in path and output_subdir:
            depth = output_subdir.count("/") + 1  # features/adult = 2 levels
            correct_prefix = "../" * depth

            # Remove any existing ../ prefixes to get the target path
            target_path = path
            while target_path.startswith("../"):
                target_path = target_path[3:]  # Remove '../'

            # Build correct path with proper depth
            correct_path = f"{correct_prefix}{target_path}"

            if path != correct_path:
                ref["path"] = correct_path
                fixes_applied += 1
                print(f"    Fixed: {original_path} → {correct_path} (depth={depth})")
                continue

    # Remove refs that were marked for removal
    for ref in refs_to_remove:
        data["design_refs"].remove(ref)

    if fixes_applied > 0:
        # Write back
        with open(yaml_file, "w") as f:
            yaml.dump(
                data, f, default_flow_style=False, sort_keys=False, allow_unicode=True
            )

    return fixes_applied


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

    print(f"📁 Found {len(yaml_files)} YAML files\n")

    total_fixes = 0
    files_modified = 0

    for yaml_file in sorted(yaml_files):
        fixes = fix_design_refs(yaml_file, repo_root)
        if fixes > 0:
            print(f"📄 {yaml_file.relative_to(repo_root)}: {fixes} fixes")
            total_fixes += fixes
            files_modified += 1

    print(f"\n✅ Applied {total_fixes} fixes to {files_modified} files")


if __name__ == "__main__":
    main()
