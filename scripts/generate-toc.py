#!/usr/bin/env python3
"""
Generate table of contents for long markdown documents.

Adds/updates a TOC section based on document headings.
Only processes documents with 5+ headings or 100+ lines.

Usage:
    python scripts/generate-toc.py                         # Dry run
    python scripts/generate-toc.py --update                # Write changes
    python scripts/generate-toc.py --file path/to/doc.md   # Single file
    python scripts/generate-toc.py --min-headings 3        # Lower threshold
"""

import argparse
import re
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DESIGN_DIR = PROJECT_ROOT / "docs" / "dev" / "design"

# TOC markers
TOC_START = "<!-- TOC-START -->"
TOC_END = "<!-- TOC-END -->"

# Files to skip
SKIP_FILES = {
    "00_SOURCE_OF_TRUTH.md",
    "01_DESIGN_DOC_TEMPLATE.md",
    "INDEX.md",
    "00_INDEX.md",
    "DESIGN_INDEX.md",
}

# Heading pattern
HEADING_PATTERN = re.compile(r"^(#{2,4})\s+(.+)$", re.MULTILINE)


def slugify(text: str) -> str:
    """Convert heading text to anchor slug."""
    slug = text.lower()
    # Remove markdown formatting
    slug = re.sub(r"\*\*([^*]+)\*\*", r"\1", slug)
    slug = re.sub(r"\*([^*]+)\*", r"\1", slug)
    slug = re.sub(r"`([^`]+)`", r"\1", slug)
    # Remove special characters
    slug = re.sub(r"[^\w\s-]", "", slug)
    slug = re.sub(r"[\s_]+", "-", slug)
    slug = slug.strip("-")
    return slug


def extract_headings(content: str) -> list[dict]:
    """Extract headings from markdown content."""
    headings = []

    # Skip content inside code blocks
    code_block_pattern = re.compile(r"```[\s\S]*?```", re.MULTILINE)
    clean_content = code_block_pattern.sub("", content)

    for match in HEADING_PATTERN.finditer(clean_content):
        level = len(match.group(1))
        text = match.group(2).strip()

        # Skip certain headings
        if text.lower() in ["table of contents", "contents", "toc"]:
            continue

        headings.append(
            {
                "level": level,
                "text": text,
                "slug": slugify(text),
            },
        )

    return headings


def generate_toc(headings: list[dict], min_level: int = 2) -> str:
    """Generate table of contents markdown."""
    lines = [
        TOC_START,
        "",
        "## Table of Contents",
        "",
    ]

    for heading in headings:
        level = heading["level"]
        text = heading["text"]
        slug = heading["slug"]

        # Calculate indentation (level 2 = no indent, level 3 = 2 spaces, etc.)
        indent = "  " * (level - min_level)
        lines.append(f"{indent}- [{text}](#{slug})")

    lines.extend(
        [
            "",
            TOC_END,
        ],
    )

    return "\n".join(lines)


def needs_toc(content: str, min_headings: int, min_lines: int) -> bool:
    """Determine if document needs a TOC."""
    headings = extract_headings(content)
    lines = content.count("\n")

    # Already has TOC
    if TOC_START in content:
        return True  # Needs update check

    # Check thresholds
    return len(headings) >= min_headings or lines >= min_lines


def update_toc(content: str, headings: list[dict]) -> tuple[str, bool]:
    """Update or insert TOC in content. Returns (new_content, changed)."""
    new_toc = generate_toc(headings)

    # Check if TOC already exists
    toc_pattern = re.compile(
        rf"{re.escape(TOC_START)}.*?{re.escape(TOC_END)}",
        re.DOTALL,
    )

    if toc_pattern.search(content):
        # Replace existing TOC
        new_content = toc_pattern.sub(new_toc, content)
        return new_content, new_content != content

    # Insert TOC after first heading and description
    # Find position after title and optional blockquote
    lines = content.split("\n")
    insert_pos = 0

    for i, line in enumerate(lines):
        # Skip title
        if line.startswith("# "):
            insert_pos = i + 1
            continue
        # Skip empty lines after title
        if insert_pos > 0 and not line.strip():
            insert_pos = i + 1
            continue
        # Skip blockquote description
        if insert_pos > 0 and line.startswith(">"):
            insert_pos = i + 1
            continue
        # Skip empty line after blockquote
        if insert_pos > 1 and not line.strip():
            insert_pos = i + 1
            break
        # Found content, insert before it
        if insert_pos > 0:
            break

    # Insert TOC
    new_lines = [*lines[:insert_pos], "", new_toc, "", *lines[insert_pos:]]
    new_content = "\n".join(new_lines)

    return new_content, True


def find_documents(directory: Path) -> list[Path]:
    """Find all markdown documents."""
    docs = []
    for md_file in sorted(directory.rglob("*.md")):
        if ".archive" in str(md_file):
            continue
        if md_file.name in SKIP_FILES:
            continue
        docs.append(md_file)
    return docs


def main():
    parser = argparse.ArgumentParser(description="Generate table of contents")
    parser.add_argument(
        "--update",
        "-u",
        action="store_true",
        help="Write changes (default: dry run)",
    )
    parser.add_argument("--file", "-f", type=Path, help="Process single file")
    parser.add_argument(
        "--min-headings",
        type=int,
        default=5,
        help="Minimum headings to trigger TOC (default: 5)",
    )
    parser.add_argument(
        "--min-lines",
        type=int,
        default=150,
        help="Minimum lines to trigger TOC (default: 150)",
    )
    parser.add_argument(
        "--force",
        action="store_true",
        help="Generate TOC even for small docs",
    )
    parser.add_argument(
        "--remove",
        action="store_true",
        help="Remove TOC from documents",
    )
    args = parser.parse_args()

    docs = [args.file] if args.file else find_documents(DESIGN_DIR)

    print(f"Checking {len(docs)} documents...")

    updated = 0
    skipped = 0
    removed = 0

    for doc_path in docs:
        content = doc_path.read_text(encoding="utf-8")
        headings = extract_headings(content)

        # Handle removal
        if args.remove:
            if TOC_START in content:
                toc_pattern = re.compile(
                    rf"{re.escape(TOC_START)}.*?{re.escape(TOC_END)}\n*",
                    re.DOTALL,
                )
                new_content = toc_pattern.sub("", content)
                if new_content != content:
                    rel_path = doc_path.relative_to(PROJECT_ROOT)
                    if args.update:
                        doc_path.write_text(new_content, encoding="utf-8")
                        print(f"  Removed TOC: {rel_path}")
                    else:
                        print(f"  Would remove TOC: {rel_path}")
                    removed += 1
            continue

        # Check if needs TOC
        if not args.force and not needs_toc(content, args.min_headings, args.min_lines):
            skipped += 1
            continue

        # Skip if too few headings (unless forced or already has TOC)
        if len(headings) < 2 and TOC_START not in content:
            skipped += 1
            continue

        # Update TOC
        new_content, changed = update_toc(content, headings)

        if changed:
            rel_path = doc_path.relative_to(PROJECT_ROOT)
            if args.update:
                doc_path.write_text(new_content, encoding="utf-8")
                print(f"  Updated: {rel_path} ({len(headings)} headings)")
            else:
                print(f"  Would update: {rel_path} ({len(headings)} headings)")
            updated += 1
        else:
            skipped += 1

    print("\n=== SUMMARY ===")
    if args.remove:
        print(f"{'Removed' if args.update else 'Would remove'}: {removed}")
    else:
        print(f"{'Updated' if args.update else 'Would update'}: {updated}")
    print(f"Skipped: {skipped}")

    if not args.update and (updated > 0 or removed > 0):
        print("\nRun with --update to write changes")


if __name__ == "__main__":
    main()
