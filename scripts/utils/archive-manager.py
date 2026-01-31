#!/usr/bin/env python3
"""Manage archived documentation files.

Provides utilities for:
- Moving deprecated docs to .archive/
- Listing archived files
- Restoring archived files

Usage:
    python scripts/utils/archive-manager.py list
    python scripts/utils/archive-manager.py archive FILE [--apply]
    python scripts/utils/archive-manager.py restore FILE [--apply]
"""

from __future__ import annotations

import argparse
import sys
from datetime import datetime
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
ARCHIVE_DIR = DESIGN_DIR / ".archive"


def list_archived() -> list[Path]:
    """List all archived files."""
    if not ARCHIVE_DIR.exists():
        return []
    return sorted(ARCHIVE_DIR.rglob("*.md"))


def archive_file(filepath: Path, *, dry_run: bool = True) -> bool:
    """Move a file to the archive."""
    if not filepath.exists():
        print(f"Error: File not found: {filepath}")
        return False

    try:
        rel_path = filepath.relative_to(DESIGN_DIR)
    except ValueError:
        print(f"Error: File must be in design directory: {filepath}")
        return False

    archive_path = ARCHIVE_DIR / rel_path
    archive_path.parent.mkdir(parents=True, exist_ok=True)

    if dry_run:
        print(f"Would archive: {rel_path} -> .archive/{rel_path}")
        return True

    # Add archive header
    content = filepath.read_text(encoding="utf-8")
    header = f"""<!-- ARCHIVED: {datetime.now().strftime("%Y-%m-%d")} -->
<!-- Original location: {rel_path} -->

"""
    archive_path.write_text(header + content, encoding="utf-8")
    filepath.unlink()
    print(f"Archived: {rel_path} -> .archive/{rel_path}")
    return True


def restore_file(filepath: Path, *, dry_run: bool = True) -> bool:
    """Restore a file from the archive."""
    if not filepath.exists():
        # Try to find in archive
        if not str(filepath).startswith(str(ARCHIVE_DIR)):
            filepath = ARCHIVE_DIR / filepath
        if not filepath.exists():
            print(f"Error: Archived file not found: {filepath}")
            return False

    try:
        rel_path = filepath.relative_to(ARCHIVE_DIR)
    except ValueError:
        print(f"Error: File must be in archive directory: {filepath}")
        return False

    restore_path = DESIGN_DIR / rel_path
    if restore_path.exists():
        print(f"Error: Target already exists: {restore_path}")
        return False

    if dry_run:
        print(f"Would restore: .archive/{rel_path} -> {rel_path}")
        return True

    restore_path.parent.mkdir(parents=True, exist_ok=True)

    # Remove archive header
    content = filepath.read_text(encoding="utf-8")
    lines = content.split("\n")
    clean_lines = []
    skip_header = True
    for line in lines:
        if skip_header and line.startswith("<!-- ARCHIVED:"):
            continue
        if skip_header and line.startswith("<!-- Original location:"):
            continue
        if skip_header and line.strip() == "":
            skip_header = False
            continue
        skip_header = False
        clean_lines.append(line)

    restore_path.write_text("\n".join(clean_lines), encoding="utf-8")
    filepath.unlink()
    print(f"Restored: .archive/{rel_path} -> {rel_path}")
    return True


def main() -> int:
    parser = argparse.ArgumentParser(description="Manage archived documentation")
    subparsers = parser.add_subparsers(dest="command", required=True)

    # List command
    list_parser = subparsers.add_parser("list", help="List archived files")
    list_parser.add_argument("--verbose", "-v", action="store_true")

    # Archive command
    archive_parser = subparsers.add_parser("archive", help="Archive a file")
    archive_parser.add_argument("file", type=Path, help="File to archive")
    archive_parser.add_argument(
        "--apply", action="store_true", help="Actually archive (default: dry-run)"
    )

    # Restore command
    restore_parser = subparsers.add_parser("restore", help="Restore an archived file")
    restore_parser.add_argument("file", type=Path, help="File to restore")
    restore_parser.add_argument(
        "--apply", action="store_true", help="Actually restore (default: dry-run)"
    )

    args = parser.parse_args()

    if args.command == "list":
        archived = list_archived()
        if not archived:
            print("No archived files found.")
            return 0

        print(f"Archived files ({len(archived)}):\n")
        for path in archived:
            rel = path.relative_to(ARCHIVE_DIR)
            if args.verbose:
                content = path.read_text(encoding="utf-8")
                lines = len(content.split("\n"))
                print(f"  {rel} ({lines} lines)")
            else:
                print(f"  {rel}")
        return 0

    elif args.command == "archive":
        filepath = args.file
        if not filepath.is_absolute():
            filepath = DESIGN_DIR / filepath
        dry_run = not args.apply
        if dry_run:
            print("=== DRY RUN (use --apply to archive) ===\n")
        success = archive_file(filepath, dry_run=dry_run)
        return 0 if success else 1

    elif args.command == "restore":
        filepath = args.file
        if not filepath.is_absolute():
            filepath = ARCHIVE_DIR / filepath
        dry_run = not args.apply
        if dry_run:
            print("=== DRY RUN (use --apply to restore) ===\n")
        success = restore_file(filepath, dry_run=dry_run)
        return 0 if success else 1

    return 0


if __name__ == "__main__":
    sys.exit(main())
