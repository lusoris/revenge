#!/usr/bin/env python3
"""Batch Migration Tool - Migrate design docs to YAML format.

This tool:
1. Scans design documentation directory
2. Parses markdown files
3. Generates YAML data files
4. Preserves directory structure
5. Creates migration report

Author: Automation System
Created: 2026-01-31
"""

import sys
from pathlib import Path


# Import from same directory
sys.path.insert(0, str(Path(__file__).parent))
from md_parser import MarkdownParser


class BatchMigrator:
    """Batch migrate design docs to YAML format."""

    def __init__(self, repo_root: Path, dry_run: bool = True, force: bool = False):
        """Initialize migrator.

        Args:
            repo_root: Repository root path
            dry_run: If True, don't write files (just report)
            force: If True, overwrite existing YAML files
        """
        self.repo_root = repo_root
        self.dry_run = dry_run
        self.force = force
        self.docs_dir = repo_root / "docs" / "dev" / "design"
        self.data_dir = repo_root / "data"
        self.parser = MarkdownParser(repo_root)

    def find_design_docs(self) -> list[Path]:
        """Find all design documentation markdown files.

        Returns:
            List of markdown file paths
        """
        # Exclude certain files
        exclude_patterns = [
            "INDEX.md",
            "DESIGN_INDEX.md",
            "00_SOURCE_OF_TRUTH.md",
            "README.md",
        ]

        md_files = []
        for md_file in self.docs_dir.glob("**/*.md"):
            if md_file.name not in exclude_patterns:
                md_files.append(md_file)

        return sorted(md_files)

    def migrate_file(self, md_file: Path) -> tuple[bool, str, Path | None]:
        """Migrate a single markdown file to YAML.

        Args:
            md_file: Path to markdown file

        Returns:
            Tuple of (success, message, output_path)
        """
        try:
            # Parse markdown
            data = self.parser.parse_file(md_file)

            # Generate YAML
            yaml_content = self.parser.to_yaml(data)

            # Determine output path
            # Preserve directory structure: docs/dev/design/X/Y.md → data/X/Y.yaml
            rel_path = md_file.relative_to(self.docs_dir)
            output_path = self.data_dir / rel_path.with_suffix(".yaml")

            # Check if file exists
            if output_path.exists() and not self.force and not self.dry_run:
                return True, "⊘ Skipped (exists)", output_path

            if not self.dry_run:
                # Create parent directories
                output_path.parent.mkdir(parents=True, exist_ok=True)

                # Write YAML file
                with open(output_path, "w") as f:
                    f.write(yaml_content)

                return True, "✓ Migrated", output_path
            else:
                return True, "✓ Would migrate", output_path

        except Exception as e:
            return False, f"✗ Error: {e}", None

    def migrate_all(self) -> dict:
        """Migrate all design docs.

        Returns:
            Migration statistics
        """
        md_files = self.find_design_docs()

        print(f"\n{'='*70}")
        print(f"BATCH MIGRATION - {'DRY RUN' if self.dry_run else 'LIVE RUN'}")
        print(f"{'='*70}\n")
        print(f"Found {len(md_files)} design docs to migrate\n")

        stats = {
            "total": len(md_files),
            "success": 0,
            "skipped": 0,
            "failed": 0,
            "by_category": {},
        }

        results = []

        for md_file in md_files:
            success, message, output_path = self.migrate_file(md_file)

            if success:
                if "Skipped" in message:
                    stats["skipped"] += 1
                else:
                    stats["success"] += 1
            else:
                stats["failed"] += 1

            # Track by category
            category = md_file.parent.name
            if category not in stats["by_category"]:
                stats["by_category"][category] = 0
            stats["by_category"][category] += 1

            results.append(
                {
                    "file": md_file.relative_to(self.docs_dir),
                    "status": message,
                    "output": output_path.relative_to(self.repo_root)
                    if output_path
                    else None,
                }
            )

        # Print results
        print(f"\n{'='*70}")
        print("MIGRATION RESULTS")
        print(f"{'='*70}\n")

        for result in results:
            status_icon = "✓" if "✓" in result["status"] else "✗"
            print(f"{status_icon} {result['file']}")
            if result["output"]:
                print(f"   → {result['output']}")
            print()

        # Print summary
        print(f"{'='*70}")
        print("SUMMARY")
        print(f"{'='*70}")
        print(f"Total files: {stats['total']}")
        print(f"Migrated: {stats['success']}")
        print(f"Skipped: {stats['skipped']}")
        print(f"Failed: {stats['failed']}")
        print("\nBy category:")
        for category, count in sorted(stats["by_category"].items()):
            print(f"  {category}: {count} files")
        print(f"{'='*70}\n")

        if self.dry_run:
            print("⚠️  DRY RUN - No files were written")
            print("   Run with --live to perform actual migration\n")

        return stats


def main():
    """Main entry point."""
    repo_root = Path(__file__).parent.parent.parent

    # Check for flags
    dry_run = "--live" not in sys.argv
    force = "--force" in sys.argv

    # Initialize migrator
    migrator = BatchMigrator(repo_root, dry_run=dry_run, force=force)

    # Run migration
    stats = migrator.migrate_all()

    # Exit with appropriate code
    sys.exit(0 if stats["failed"] == 0 else 1)


if __name__ == "__main__":
    main()
