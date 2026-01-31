#!/usr/bin/env python3
"""
Validate all internal markdown links in documentation.

Checks:
- Relative links to .md files exist
- Anchor links (#section) are valid
- No broken references

Usage:
    python scripts/validate-links.py                    # Check all docs
    python scripts/validate-links.py --fix             # Auto-fix simple issues
    python scripts/validate-links.py docs/dev/design   # Check specific directory
"""

import argparse
import re
import sys
from pathlib import Path


SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DOCS_DIR = PROJECT_ROOT / "docs"

# Link pattern: [text](url)
LINK_PATTERN = re.compile(r"\[([^\]]*)\]\(([^)]+)\)")

# Heading pattern for anchor validation
HEADING_PATTERN = re.compile(r"^#{1,6}\s+(.+)$", re.MULTILINE)


def slugify(text: str) -> str:
    """Convert heading text to anchor slug."""
    # GitHub-style slugification
    slug = text.lower()
    slug = re.sub(r"[^\w\s-]", "", slug)
    slug = re.sub(r"[\s_]+", "-", slug)
    slug = slug.strip("-")
    return slug


def extract_headings(content: str) -> set[str]:
    """Extract all heading anchors from content."""
    headings = set()
    for match in HEADING_PATTERN.finditer(content):
        heading_text = match.group(1).strip()
        # Remove markdown formatting
        heading_text = re.sub(r"\*\*([^*]+)\*\*", r"\1", heading_text)
        heading_text = re.sub(r"\*([^*]+)\*", r"\1", heading_text)
        heading_text = re.sub(r"`([^`]+)`", r"\1", heading_text)
        headings.add(slugify(heading_text))
    return headings


def find_markdown_files(directory: Path) -> list[Path]:
    """Find all markdown files in directory."""
    files = []
    for md_file in sorted(directory.rglob("*.md")):
        if ".archive" in str(md_file):
            continue
        files.append(md_file)
    return files


def validate_link(
    source_file: Path,
    link_text: str,
    link_url: str,
    file_headings_cache: dict,
) -> tuple[bool, str]:
    """
    Validate a single link.
    Returns (is_valid, error_message).
    """
    # Skip external links
    if link_url.startswith(("http://", "https://", "mailto:")):
        return True, ""

    # Skip special links
    if link_url.startswith("#"):
        # Anchor-only link - check in same file
        anchor = link_url[1:]
        if source_file not in file_headings_cache:
            content = source_file.read_text(encoding="utf-8")
            file_headings_cache[source_file] = extract_headings(content)

        headings = file_headings_cache[source_file]
        if anchor not in headings:
            return False, f"Anchor #{anchor} not found in file"
        return True, ""

    # Parse link
    if "#" in link_url:
        file_path, anchor = link_url.split("#", 1)
    else:
        file_path = link_url
        anchor = None

    # Skip non-markdown links (images, etc)
    if not file_path.endswith((".md", ".json", ".yaml", ".graphql")):
        # Could be directory link
        if file_path.endswith("/"):
            target = (source_file.parent / file_path).resolve()
            if target.is_dir():
                return True, ""
            return False, f"Directory not found: {file_path}"
        return True, ""  # Assume valid (could be image, etc)

    # Resolve relative path
    try:
        target = (source_file.parent / file_path).resolve()
    except (ValueError, OSError) as e:
        return False, f"Invalid path: {file_path} ({e})"

    # Check file exists
    if not target.exists():
        return False, f"File not found: {file_path}"

    # Check anchor if present
    if anchor:
        if target not in file_headings_cache:
            try:
                content = target.read_text(encoding="utf-8")
                file_headings_cache[target] = extract_headings(content)
            except Exception as e:
                return False, f"Cannot read target file: {e}"

        headings = file_headings_cache[target]
        if anchor not in headings:
            # Try common variations
            anchor_lower = anchor.lower()
            if anchor_lower not in headings:
                return False, f"Anchor #{anchor} not found in {file_path}"

    return True, ""


def validate_file(
    file_path: Path,
    file_headings_cache: dict,
) -> list[dict]:
    """Validate all links in a file. Returns list of errors."""
    errors = []

    try:
        content = file_path.read_text(encoding="utf-8")
    except Exception as e:
        return [{"line": 0, "link": "", "error": f"Cannot read file: {e}"}]

    lines = content.split("\n")
    for line_num, line in enumerate(lines, 1):
        for match in LINK_PATTERN.finditer(line):
            link_text = match.group(1)
            link_url = match.group(2)

            is_valid, error = validate_link(
                file_path, link_text, link_url, file_headings_cache,
            )

            if not is_valid:
                errors.append(
                    {
                        "line": line_num,
                        "text": link_text,
                        "url": link_url,
                        "error": error,
                    },
                )

    return errors


def main():
    parser = argparse.ArgumentParser(description="Validate markdown links")
    parser.add_argument(
        "path",
        nargs="?",
        default=str(DOCS_DIR),
        help="Directory or file to check (default: docs/)",
    )
    parser.add_argument(
        "--fix", "-f", action="store_true", help="Attempt to fix simple issues",
    )
    parser.add_argument(
        "--verbose", "-v", action="store_true", help="Show all checked files",
    )
    args = parser.parse_args()

    target = Path(args.path)
    if not target.is_absolute():
        target = PROJECT_ROOT / target
    target = target.resolve()

    if not target.exists():
        print(f"Error: {target} does not exist")
        return 1

    files = [target] if target.is_file() else find_markdown_files(target)

    print(f"Checking {len(files)} markdown files...")

    file_headings_cache = {}
    total_errors = 0
    files_with_errors = 0

    for file_path in files:
        errors = validate_file(file_path, file_headings_cache)

        if errors:
            files_with_errors += 1
            try:
                rel_path = file_path.relative_to(PROJECT_ROOT)
            except ValueError:
                rel_path = file_path
            print(f"\n{rel_path}:")
            for err in errors:
                print(f"  Line {err['line']}: [{err.get('text', '')}]({err['url']})")
                print(f"    → {err['error']}")
                total_errors += 1
        elif args.verbose:
            try:
                rel_path = file_path.relative_to(PROJECT_ROOT)
            except ValueError:
                rel_path = file_path
            print(f"  ✓ {rel_path}")

    print("\n=== SUMMARY ===")
    print(f"Files checked: {len(files)}")
    print(f"Files with errors: {files_with_errors}")
    print(f"Total errors: {total_errors}")

    if total_errors > 0:
        print("\nRun with --verbose to see all checked files")
        return 1

    print("\n✓ All links valid!")
    return 0


if __name__ == "__main__":
    sys.exit(main())
