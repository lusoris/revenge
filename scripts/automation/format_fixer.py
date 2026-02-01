#!/usr/bin/env python3
"""Validation Format Fixer - Fix common validation format issues.

This tool fixes known validation format issues:
1. fx_module: "AuthModule" → "auth.Module" (schema expects package.Module format)
2. overall_status: "✅" → "✅ Complete" (schema expects emoji + description)

Author: Automation System
Created: 2026-01-31
"""

import sys
from pathlib import Path

import yaml


class FormatFixer:
    """Fix validation format issues in YAML files."""

    def __init__(self, repo_root: Path):
        """Initialize fixer.

        Args:
            repo_root: Repository root path
        """
        self.repo_root = repo_root
        self.data_dir = repo_root / "data"

    def fix_fx_module_format(self, data: dict) -> bool:
        """Fix fx_module format from 'AuthModule' to 'auth.Module'.

        Args:
            data: YAML data dict

        Returns:
            True if fixed, False if no change needed
        """
        if "fx_module" not in data:
            return False

        fx_module = data["fx_module"]

        # Check if already in correct format (package.Module)
        if "." in fx_module:
            return False

        # Convert: "AuthModule" → "auth.Module"
        # Extract package name from module name
        # "AuthModule" → "auth", "UserModule" → "user"
        if fx_module.endswith("Module"):
            package_name = fx_module[:-6].lower()  # Remove "Module" suffix
            new_format = f"{package_name}.Module"
            data["fx_module"] = new_format
            return True

        return False

    def fix_overall_status_format(self, data: dict) -> bool:
        """Fix overall_status format from '✅' to '✅ Complete'.

        Args:
            data: YAML data dict

        Returns:
            True if fixed, False if no change needed
        """
        if "overall_status" not in data:
            return False

        status = str(data["overall_status"])

        # If already has description, no fix needed
        if len(status) > 2:  # More than just emoji
            return False

        # Map emoji to description
        status_map = {
            "✅": "✅ Complete",
            "🟡": "🟡 In Progress",
            "🔴": "🔴 Not Started",
            "🔵": "🔵 Planning",
        }

        # Fix if it's a lone emoji
        if status in status_map:
            data["overall_status"] = status_map[status]
            return True

        return False

    def fix_file(self, yaml_file: Path, dry_run: bool = False) -> dict:
        """Fix validation format issues in a YAML file.

        Args:
            yaml_file: Path to YAML file
            dry_run: If True, don't write changes

        Returns:
            Dict with fix statistics
        """
        stats = {"fx_module_fixed": False, "overall_status_fixed": False}

        # Load YAML
        with open(yaml_file) as f:
            data = yaml.safe_load(f)

        # Apply fixes
        stats["fx_module_fixed"] = self.fix_fx_module_format(data)
        stats["overall_status_fixed"] = self.fix_overall_status_format(data)

        # Write if any fixes applied and not dry-run
        if (stats["fx_module_fixed"] or stats["overall_status_fixed"]) and not dry_run:
            with open(yaml_file, "w") as f:
                yaml.safe_dump(
                    data,
                    f,
                    default_flow_style=False,
                    allow_unicode=True,
                    sort_keys=False,
                )

        return stats

    def fix_all(self, category: str | None = None, dry_run: bool = False) -> dict:
        """Fix all YAML files or specific category.

        Args:
            category: Category to fix (feature/service/integration), or None for all
            dry_run: If True, don't write changes

        Returns:
            Statistics dict
        """
        # Find YAML files
        if category:
            category_dirs = {
                "feature": ["features"],
                "service": ["services"],
                "integration": ["integrations"],
            }
            yaml_files = []
            for dir_name in category_dirs.get(category, []):
                category_path = self.data_dir / dir_name
                if category_path.exists():
                    yaml_files.extend(category_path.glob("**/*.yaml"))
        else:
            # All YAML files except shared-sot.yaml
            yaml_files = [
                f
                for f in self.data_dir.glob("**/*.yaml")
                if f.name != "shared-sot.yaml"
            ]

        if not yaml_files:
            print("\n⚠️  No YAML files found")
            return {"total": 0, "fixed": 0, "fx_module_fixes": 0, "status_fixes": 0}

        print(f"\n{'=' * 70}")
        print("VALIDATION FORMAT FIXER")
        if category:
            print(f"Category: {category}")
        else:
            print("Processing: All files")
        print(f"Mode: {'DRY-RUN' if dry_run else 'LIVE'}")
        print(f"{'=' * 70}\n")
        print(f"Found {len(yaml_files)} files\n")

        stats = {
            "total": len(yaml_files),
            "fixed": 0,
            "fx_module_fixes": 0,
            "status_fixes": 0,
        }

        for yaml_file in sorted(yaml_files):
            rel_path = yaml_file.relative_to(self.data_dir)

            file_stats = self.fix_file(yaml_file, dry_run)

            if file_stats["fx_module_fixed"] or file_stats["overall_status_fixed"]:
                stats["fixed"] += 1
                if file_stats["fx_module_fixed"]:
                    stats["fx_module_fixes"] += 1
                if file_stats["overall_status_fixed"]:
                    stats["status_fixes"] += 1

                fixes = []
                if file_stats["fx_module_fixed"]:
                    fixes.append("fx_module")
                if file_stats["overall_status_fixed"]:
                    fixes.append("overall_status")

                print(f"📝 {rel_path}")
                print(f"   Fixed: {', '.join(fixes)}")
                if not dry_run:
                    print("   ✅ Updated")
                else:
                    print("   ℹ️  Would update")
                print()

        # Summary
        print(f"{'=' * 70}")
        print("SUMMARY")
        print(f"{'=' * 70}")
        print(f"Total files: {stats['total']}")
        print(f"Files fixed: {stats['fixed']}")
        print(f"  • fx_module fixes: {stats['fx_module_fixes']}")
        print(f"  • overall_status fixes: {stats['status_fixes']}")
        if dry_run:
            print("\n⚠️  DRY-RUN MODE - No changes written")
        print(f"{'=' * 70}\n")

        return stats


def main():
    """Main entry point."""
    repo_root = Path(__file__).parent.parent.parent

    # Parse arguments
    args = sys.argv[1:]

    category = None
    dry_run = True

    for arg in args:
        if arg == "--feature":
            category = "feature"
        elif arg == "--service":
            category = "service"
        elif arg == "--integration":
            category = "integration"
        elif arg == "--all":
            category = None
        elif arg == "--live":
            dry_run = False
        elif arg == "--dry-run":
            dry_run = True

    # Default to dry-run unless --live specified
    if "--live" not in args:
        dry_run = True

    # Show usage if no category specified
    if category is None and "--all" not in args:
        print("Usage: python format_fixer.py <category|--all> [--live]")
        print()
        print("Categories:")
        print("  --feature      Fix feature modules only")
        print("  --service      Fix backend services only")
        print("  --integration  Fix integrations only")
        print("  --all          Fix all files")
        print()
        print("Options:")
        print("  --live         Apply fixes (default: dry-run)")
        print("  --dry-run      Show what would be fixed (default)")
        print()
        print("Examples:")
        print("  python format_fixer.py --service --dry-run")
        print("  python format_fixer.py --all --live")
        sys.exit(1)

    # Initialize fixer
    fixer = FormatFixer(repo_root)

    # Fix files
    fixer.fix_all(category, dry_run)

    # Exit code based on results
    sys.exit(0)


if __name__ == "__main__":
    main()
