#!/usr/bin/env python3
"""Batch Doc Regenerator - Regenerate all docs from YAML data.

This tool:
1. Finds all YAML data files
2. Determines appropriate template for each
3. Regenerates Claude and Wiki docs
4. Outputs to designated directories
5. Provides safety options (preview mode, backup, etc.)

Author: Automation System
Created: 2026-01-31
"""

import shutil
import sys
from datetime import datetime
from pathlib import Path


# Import from same directory
sys.path.insert(0, str(Path(__file__).parent))
from doc_generator import DocGenerator


class BatchRegenerator:
    """Batch regenerate documentation from YAML files."""

    def __init__(
        self,
        repo_root: Path,
        backup_originals: bool = False,
    ):
        """Initialize regenerator.

        Args:
            repo_root: Repository root path
            backup_originals: If True, backup existing files before overwriting
        """
        self.repo_root = repo_root
        self.backup_originals = backup_originals
        self.data_dir = repo_root / "data"
        self.generator = DocGenerator(repo_root)

        # Output directories
        self.output_claude = repo_root / "docs" / "dev" / "design"
        self.output_wiki = repo_root / "docs" / "wiki"
        if backup_originals:
            print("💾 BACKUP MODE - Will backup existing files")

    def find_yaml_files(self) -> list[Path]:
        """Find all YAML data files.

        Returns:
            List of YAML file paths
        """
        yaml_files = list(self.data_dir.glob("**/*.yaml"))

        # Exclude shared-sot.yaml
        yaml_files = [f for f in yaml_files if f.name != "shared-sot.yaml"]

        return sorted(yaml_files)

    def determine_template(self, category: str) -> str:
        """Determine which template to use based on category.

        Args:
            category: Document category

        Returns:
            Template filename
        """
        template_map = {
            "feature": "feature.md.jinja2",
            "service": "service.md.jinja2",
            "integration": "integration.md.jinja2",
        }

        return template_map.get(category, "generic.md.jinja2")

    def backup_file(self, file_path: Path):
        """Backup a file if it exists.

        Args:
            file_path: File to backup
        """
        if not file_path.exists():
            return

        # Create backup with timestamp
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        backup_path = file_path.with_suffix(f".backup_{timestamp}.md")

        shutil.copy2(file_path, backup_path)
        print(f"    💾 Backed up to: {backup_path.name}")

    def regenerate_file(self, yaml_file: Path) -> tuple[bool, str, dict]:
        """Regenerate docs for a single YAML file.

        Args:
            yaml_file: Path to YAML data file

        Returns:
            Tuple of (success, message, generated_paths)
        """
        try:
            # Load YAML to get category
            import yaml

            with open(yaml_file) as f:
                data = yaml.safe_load(f)

            category = data.get("doc_category", "other")
            template = self.determine_template(category)

            # Determine output subpath (preserve directory structure)
            rel_path = yaml_file.relative_to(self.data_dir)
            output_subpath = str(rel_path.parent)

            # Backup if needed
            if self.backup_originals:
                claude_path = self.output_claude / rel_path.with_suffix(".md")
                wiki_path = self.output_wiki / rel_path.with_suffix(".md")
                self.backup_file(claude_path)
                self.backup_file(wiki_path)

            # Generate docs
            generated = self.generator.generate_doc(
                data_file=yaml_file,
                template_name=template,
                output_subpath=output_subpath,
                render_both=True,
            )

            return True, "✓ Generated", generated

        except Exception as e:
            return False, f"✗ Error: {e}", {}

    def regenerate_all(self) -> dict:
        """Regenerate all documentation.

        Returns:
            Statistics dict
        """
        yaml_files = self.find_yaml_files()

        print(f"\n{'=' * 70}")
        print("BATCH REGENERATION")
        print(f"{'=' * 70}\n")
        print(f"Found {len(yaml_files)} YAML files to process\n")

        stats = {
            "total": len(yaml_files),
            "success": 0,
            "failed": 0,
            "by_category": {},
        }

        for yaml_file in yaml_files:
            rel_path = yaml_file.relative_to(self.data_dir)
            print(f"📄 {rel_path}")

            success, message, generated = self.regenerate_file(yaml_file)

            if success:
                stats["success"] += 1
                for target, path in generated.items():
                    print(f"  ✓ {target}: {path.relative_to(self.repo_root)}")
            else:
                stats["failed"] += 1
                print(f"  {message}")

            # Track by category (if available)
            try:
                import yaml

                with open(yaml_file) as f:
                    data = yaml.safe_load(f)
                category = data.get("doc_category", "unknown")
                stats["by_category"][category] = (
                    stats["by_category"].get(category, 0) + 1
                )
            except Exception:
                pass

            print()

        # Summary
        print(f"{'=' * 70}")
        print("SUMMARY")
        print(f"{'=' * 70}")
        print(f"Total files: {stats['total']}")
        print(f"Success: {stats['success']}")
        print(f"Failed: {stats['failed']}")
        print("\nBy category:")
        for category, count in sorted(stats["by_category"].items()):
            print(f"  {category}: {count} files")
        print(f"{'=' * 70}\n")

        return stats


def main():
    """Main entry point."""
    import argparse

    parser = argparse.ArgumentParser(
        description="Batch regenerate documentation from YAML data files",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Regenerate all docs
  python batch_regenerate.py

  # Regenerate with backups
  python batch_regenerate.py --backup
        """,
    )

    parser.add_argument(
        "--backup",
        action="store_true",
        help="Backup existing files before overwriting",
    )

    args = parser.parse_args()

    repo_root = Path(__file__).parent.parent.parent

    # Initialize regenerator
    regenerator = BatchRegenerator(
        repo_root,
        backup_originals=args.backup,
    )

    # Run regeneration (includes TOC generation via doc_generator.py)
    stats = regenerator.regenerate_all()

    # Exit code
    sys.exit(0 if stats["failed"] == 0 else 1)


if __name__ == "__main__":
    main()
