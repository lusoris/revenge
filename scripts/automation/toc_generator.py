#!/usr/bin/env python3
"""TOC Generator - Generate table of contents for markdown files.

This tool generates TOC (table of contents) for markdown files by:
- Parsing markdown headers (## and ###)
- Creating anchor links
- Inserting TOC after frontmatter

Author: Automation System
Created: 2026-01-31
"""

import re
import sys
from pathlib import Path


class TOCGenerator:
    """Generate table of contents for markdown files."""

    def __init__(self):
        """Initialize TOC generator."""
        # Match h1-h6 headers (all levels)
        self.header_pattern = re.compile(r"^(#{1,6})\s+(.+)$", re.MULTILINE)

    def extract_headers(self, content: str) -> list[tuple[int, str]]:
        """Extract headers from markdown content.

        Args:
            content: Markdown content

        Returns:
            List of (level, header_text) tuples
        """
        headers = []
        for match in self.header_pattern.finditer(content):
            level = len(match.group(1))  # Number of # symbols
            text = match.group(2).strip()
            headers.append((level, text))
        return headers

    def generate_anchor(self, text: str) -> str:
        """Generate GitHub-flavored markdown anchor from header text.

        Args:
            text: Header text

        Returns:
            Anchor slug
        """
        # Convert to lowercase
        anchor = text.lower()

        # Remove special characters except spaces and hyphens
        anchor = re.sub(r"[^\w\s-]", "", anchor)

        # Replace spaces with hyphens
        anchor = re.sub(r"\s+", "-", anchor)

        # Remove multiple consecutive hyphens
        anchor = re.sub(r"-+", "-", anchor)

        # Strip leading/trailing hyphens
        anchor = anchor.strip("-")

        return anchor

    def generate_toc(self, headers: list[tuple[int, str]]) -> str:
        """Generate TOC markdown from headers.

        Args:
            headers: List of (level, header_text) tuples

        Returns:
            TOC markdown string
        """
        if not headers:
            return ""

        toc_lines = ["## Table of Contents\n"]

        for level, text in headers:
            # Skip ## Table of Contents itself
            if text == "Table of Contents":
                continue

            # h1 = no indent, h2 = 2 spaces, h3 = 4 spaces, etc.
            indent = "  " * (level - 1)
            anchor = self.generate_anchor(text)
            toc_lines.append(f"{indent}- [{text}](#{anchor})")

        return "\n".join(toc_lines) + "\n"

    def split_frontmatter(self, content: str) -> tuple[str, str]:
        """Split content into frontmatter and body.

        Args:
            content: Markdown content

        Returns:
            Tuple of (frontmatter, body)
        """
        # Check for YAML frontmatter
        if content.startswith("---\n"):
            # Find closing ---
            parts = content.split("---\n", 2)
            if len(parts) >= 3:
                frontmatter = f"---\n{parts[1]}---\n"
                body = parts[2]
                return frontmatter, body

        # No frontmatter
        return "", content

    def has_toc(self, content: str) -> bool:
        """Check if content already has a TOC.

        Args:
            content: Markdown content

        Returns:
            True if TOC exists
        """
        return "## Table of Contents" in content

    def remove_existing_toc(self, content: str) -> str:
        """Remove existing TOC from content.

        Args:
            content: Markdown content

        Returns:
            Content with TOC removed
        """
        # Find TOC section
        toc_start = content.find("## Table of Contents")
        if toc_start == -1:
            return content

        # Find next header (any level) or horizontal rule after TOC
        # Look for patterns: \n# , \n## , \n### , or \n---
        search_from = toc_start + len("## Table of Contents")

        # Try to find next section marker
        patterns = [
            content.find("\n---", search_from),
            content.find("\n# ", search_from),
            content.find("\n## ", search_from),
        ]

        # Filter out -1 (not found) and get minimum position
        valid_positions = [p for p in patterns if p != -1]

        if valid_positions:
            next_section = min(valid_positions)
            # Remove TOC section up to (but not including) the newline before next section
            return content[:toc_start] + content[next_section + 1:]
        # TOC is at the end
        return content[:toc_start]

    def add_toc(self, content: str) -> str:
        """Add TOC to markdown content.

        Args:
            content: Markdown content

        Returns:
            Content with TOC added
        """
        # Split frontmatter and body
        frontmatter, body = self.split_frontmatter(content)

        # Remove existing TOC if present
        body = self.remove_existing_toc(body)

        # Extract headers
        headers = self.extract_headers(body)

        # Generate TOC
        toc = self.generate_toc(headers)

        if not toc:
            # No headers, return original
            return content

        # Insert TOC after frontmatter
        if frontmatter:
            return f"{frontmatter}\n{toc}\n{body}"
        return f"{toc}\n{body}"

    def process_file(self, file_path: Path, dry_run: bool = False) -> bool:
        """Process a markdown file to add/update TOC.

        Args:
            file_path: Path to markdown file
            dry_run: If True, don't write changes

        Returns:
            True if file was modified
        """
        # Read file
        with open(file_path) as f:
            original_content = f.read()

        # Add TOC
        new_content = self.add_toc(original_content)

        # Check if changed
        if new_content == original_content:
            return False

        # Write if not dry-run
        if not dry_run:
            with open(file_path, "w") as f:
                f.write(new_content)

        return True

    def process_directory(
        self, directory: Path, pattern: str = "**/*.md", dry_run: bool = False,
    ) -> dict:
        """Process all markdown files in directory.

        Args:
            directory: Directory to process
            pattern: Glob pattern for markdown files
            dry_run: If True, don't write changes

        Returns:
            Statistics dict
        """
        stats = {"processed": 0, "modified": 0, "unchanged": 0}

        md_files = list(directory.glob(pattern))

        for md_file in md_files:
            stats["processed"] += 1

            modified = self.process_file(md_file, dry_run)

            if modified:
                stats["modified"] += 1
                if dry_run:
                    print(f"  Would update: {md_file.relative_to(directory)}")
                else:
                    print(f"✓ Updated: {md_file.relative_to(directory)}")
            else:
                stats["unchanged"] += 1

        return stats


def main():
    """Main entry point."""
    # Parse arguments
    args = sys.argv[1:]

    if not args or "--help" in args:
        print("Usage: python toc_generator.py <file_or_directory> [--live]")
        print()
        print("Arguments:")
        print("  file_or_directory  Path to markdown file or directory")
        print()
        print("Options:")
        print("  --live             Apply changes (default: dry-run)")
        print("  --dry-run          Show what would be changed (default)")
        print()
        print("Examples:")
        print("  python toc_generator.py docs/dev/design/ --dry-run")
        print("  python toc_generator.py docs/dev/design/features/video/MOVIE_MODULE.md --live")
        sys.exit(0)

    path_arg = args[0]
    dry_run = "--live" not in args

    path = Path(path_arg)

    if not path.exists():
        print(f"❌ Path not found: {path}")
        sys.exit(1)

    # Initialize generator
    generator = TOCGenerator()

    print(f"\n{'='*70}")
    print(f"TOC GENERATOR - {'DRY RUN' if dry_run else 'LIVE'}")
    print(f"{'='*70}\n")

    if path.is_file():
        # Process single file
        print(f"Processing: {path}")
        modified = generator.process_file(path, dry_run)

        if modified:
            if dry_run:
                print("  Would update file")
            else:
                print("✓ Updated file")
        else:
            print("  No changes needed")

    elif path.is_dir():
        # Process directory
        print(f"Processing: {path}/**/*.md\n")
        stats = generator.process_directory(path, dry_run=dry_run)

        print(f"\n{'='*70}")
        print("SUMMARY")
        print(f"{'='*70}")
        print(f"Processed: {stats['processed']}")
        print(f"Modified: {stats['modified']}")
        print(f"Unchanged: {stats['unchanged']}")

        if dry_run:
            print("\n⚠️  DRY RUN MODE - No changes written")

        print(f"{'='*70}\n")

    else:
        print(f"❌ Invalid path: {path}")
        sys.exit(1)

    sys.exit(0)


if __name__ == "__main__":
    main()
