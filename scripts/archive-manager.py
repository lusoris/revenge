#!/usr/bin/env python3
"""
Manage archived documentation files.

Operations:
- Archive deprecated docs to .archive/ with redirect stubs
- List currently archived files
- Restore archived files
- Clean up old archives

Usage:
    python scripts/archive-manager.py list                     # List archived
    python scripts/archive-manager.py archive path/to/doc.md   # Archive a doc
    python scripts/archive-manager.py restore path/to/doc.md   # Restore from archive
    python scripts/archive-manager.py cleanup --days 90        # Remove old archives
"""

import argparse
import shutil
from datetime import datetime
from pathlib import Path

SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"
ARCHIVE_DIR = DESIGN_DIR / ".archive"

# Redirect stub template
REDIRECT_STUB = """# {title}

> **This document has been archived.**

This document was archived on {date} and moved to:
[.archive/{filename}](.archive/{filename})

## Reason

{reason}

---

*If you need this document restored, run:*
```bash
python scripts/archive-manager.py restore {original_path}
```
"""


def ensure_archive_dir():
    """Ensure archive directory exists."""
    ARCHIVE_DIR.mkdir(parents=True, exist_ok=True)

    # Create .gitkeep if empty
    gitkeep = ARCHIVE_DIR / ".gitkeep"
    if not gitkeep.exists():
        gitkeep.touch()


def list_archives():
    """List all archived files."""
    if not ARCHIVE_DIR.exists():
        print("No archive directory found.")
        return

    archives = list(ARCHIVE_DIR.glob("*.md"))
    if not archives:
        print("No archived files found.")
        return

    print(f"Archived files ({len(archives)}):\n")

    for archive in sorted(archives):
        stat = archive.stat()
        date = datetime.fromtimestamp(stat.st_mtime).strftime("%Y-%m-%d")
        size = stat.st_size

        # Try to extract original location from content
        content = archive.read_text(encoding="utf-8")
        original = "unknown"
        if "Original location:" in content:
            for line in content.split("\n"):
                if "Original location:" in line:
                    original = line.split(":", 1)[1].strip()
                    break

        print(f"  {archive.name}")
        print(f"    Archived: {date} | Size: {size:,} bytes")
        print(f"    Original: {original}")
        print()


def archive_document(doc_path: Path, reason: str = "Document deprecated", dry_run: bool = True):
    """Archive a document, leaving a redirect stub."""
    if not doc_path.exists():
        print(f"Error: {doc_path} does not exist")
        return False

    if ".archive" in str(doc_path):
        print(f"Error: {doc_path} is already in archive")
        return False

    ensure_archive_dir()

    # Read original content
    content = doc_path.read_text(encoding="utf-8")

    # Extract title
    title = doc_path.stem.replace("_", " ").replace("-", " ")
    for line in content.split("\n"):
        if line.startswith("# "):
            title = line[2:].strip()
            break

    # Generate archive filename (with timestamp to avoid conflicts)
    timestamp = datetime.now().strftime("%Y%m%d")
    archive_name = f"{doc_path.stem}_{timestamp}.md"
    archive_path = ARCHIVE_DIR / archive_name

    # Add metadata header to archived content
    archived_content = f"""<!--
Archived: {datetime.now().isoformat()}
Original location: {doc_path.relative_to(PROJECT_ROOT)}
Reason: {reason}
-->

{content}
"""

    # Generate redirect stub
    rel_path = doc_path.relative_to(DESIGN_DIR)
    stub = REDIRECT_STUB.format(
        title=title,
        date=datetime.now().strftime("%Y-%m-%d"),
        filename=archive_name,
        reason=reason,
        original_path=rel_path,
    )

    if dry_run:
        print(f"Would archive: {doc_path.relative_to(PROJECT_ROOT)}")
        print(f"  → {archive_path.relative_to(PROJECT_ROOT)}")
        print(f"  Reason: {reason}")
        return True

    # Write archived content
    archive_path.write_text(archived_content, encoding="utf-8")

    # Write redirect stub
    doc_path.write_text(stub, encoding="utf-8")

    print(f"Archived: {doc_path.relative_to(PROJECT_ROOT)}")
    print(f"  → {archive_path.relative_to(PROJECT_ROOT)}")
    return True


def restore_document(doc_path: Path, dry_run: bool = True):
    """Restore a document from archive."""
    # Check if it's a stub
    if doc_path.exists():
        content = doc_path.read_text(encoding="utf-8")
        if "This document has been archived" not in content:
            print(f"Error: {doc_path} is not an archive stub")
            return False

        # Find archive reference
        archive_name = None
        for line in content.split("\n"):
            if ".archive/" in line and ".md" in line:
                match = line.split(".archive/")[1].split(")")[0].split("]")[0]
                archive_name = match.strip()
                break

        if not archive_name:
            print(f"Error: Could not find archive reference in {doc_path}")
            return False

        archive_path = ARCHIVE_DIR / archive_name
    else:
        # Assume doc_path is the archive path
        archive_path = doc_path
        if not archive_path.exists():
            archive_path = ARCHIVE_DIR / doc_path.name

    if not archive_path.exists():
        print(f"Error: Archive file not found: {archive_path}")
        return False

    # Read archived content
    content = archive_path.read_text(encoding="utf-8")

    # Extract original location and content
    original_path = None
    if "Original location:" in content:
        for line in content.split("\n"):
            if "Original location:" in line:
                original_path = line.split(":", 1)[1].strip()
                break

    # Remove metadata header
    if content.startswith("<!--"):
        end_comment = content.find("-->")
        if end_comment != -1:
            content = content[end_comment + 3:].strip()

    # Determine restore path
    if original_path:
        restore_path = PROJECT_ROOT / original_path
    else:
        # Default to design dir with original name (minus timestamp)
        name = archive_path.stem
        if "_20" in name:  # Remove timestamp suffix
            name = name.rsplit("_", 1)[0]
        restore_path = DESIGN_DIR / f"{name}.md"

    if dry_run:
        print(f"Would restore: {archive_path.relative_to(PROJECT_ROOT)}")
        print(f"  → {restore_path.relative_to(PROJECT_ROOT)}")
        return True

    # Write restored content
    restore_path.parent.mkdir(parents=True, exist_ok=True)
    restore_path.write_text(content, encoding="utf-8")

    # Remove archive file
    archive_path.unlink()

    print(f"Restored: {restore_path.relative_to(PROJECT_ROOT)}")
    return True


def cleanup_archives(days: int, dry_run: bool = True):
    """Remove archives older than specified days."""
    if not ARCHIVE_DIR.exists():
        print("No archive directory found.")
        return

    cutoff = datetime.now().timestamp() - (days * 24 * 60 * 60)
    removed = 0

    for archive in ARCHIVE_DIR.glob("*.md"):
        stat = archive.stat()
        if stat.st_mtime < cutoff:
            age_days = (datetime.now().timestamp() - stat.st_mtime) / (24 * 60 * 60)
            if dry_run:
                print(f"Would remove: {archive.name} ({age_days:.0f} days old)")
            else:
                archive.unlink()
                print(f"Removed: {archive.name} ({age_days:.0f} days old)")
            removed += 1

    print(f"\n{'Would remove' if dry_run else 'Removed'}: {removed} archives")


def main():
    parser = argparse.ArgumentParser(description="Manage documentation archives")
    subparsers = parser.add_subparsers(dest="command", help="Command to run")

    # List command
    list_parser = subparsers.add_parser("list", help="List archived files")

    # Archive command
    archive_parser = subparsers.add_parser("archive", help="Archive a document")
    archive_parser.add_argument("path", type=Path, help="Document to archive")
    archive_parser.add_argument(
        "--reason", "-r", default="Document deprecated", help="Reason for archiving"
    )
    archive_parser.add_argument(
        "--execute", "-x", action="store_true", help="Actually archive (default: dry run)"
    )

    # Restore command
    restore_parser = subparsers.add_parser("restore", help="Restore from archive")
    restore_parser.add_argument("path", type=Path, help="Document or archive to restore")
    restore_parser.add_argument(
        "--execute", "-x", action="store_true", help="Actually restore (default: dry run)"
    )

    # Cleanup command
    cleanup_parser = subparsers.add_parser("cleanup", help="Remove old archives")
    cleanup_parser.add_argument(
        "--days", "-d", type=int, default=90, help="Remove archives older than N days"
    )
    cleanup_parser.add_argument(
        "--execute", "-x", action="store_true", help="Actually remove (default: dry run)"
    )

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        return

    if args.command == "list":
        list_archives()

    elif args.command == "archive":
        if not args.path.is_absolute():
            args.path = PROJECT_ROOT / args.path
        archive_document(args.path, args.reason, dry_run=not args.execute)

    elif args.command == "restore":
        if not args.path.is_absolute():
            args.path = PROJECT_ROOT / args.path
        restore_document(args.path, dry_run=not args.execute)

    elif args.command == "cleanup":
        cleanup_archives(args.days, dry_run=not args.execute)


if __name__ == "__main__":
    main()
