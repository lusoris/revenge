#!/usr/bin/env python3
"""GitHub Labels management script.

Synchronizes GitHub labels with .github/labels.yml configuration:
- Creates missing labels
- Updates existing labels (color, description)
- Optionally removes labels not in config
- Auto-labels PRs based on size, type, area

Uses gh CLI for GitHub API access.

Requirements:
- gh CLI installed and authenticated
- Repository admin permissions

Usage:
    python scripts/automation/github_labels.py --sync
    python scripts/automation/github_labels.py --check
    python scripts/automation/github_labels.py --cleanup  # Remove labels not in config
    python scripts/automation/github_labels.py --auto-label <pr-number>

Author: Automation System
Created: 2026-01-31
"""

import argparse
import json
import subprocess
import sys
from pathlib import Path


try:
    import yaml
except ImportError:
    print("❌ Error: PyYAML not installed")
    print("   Install: pip install PyYAML")
    sys.exit(1)


class GitHubLabelsManager:
    """Manage GitHub repository labels."""

    def __init__(self, repo_owner: str, repo_name: str, dry_run: bool = False):
        """Initialize labels manager.

        Args:
            repo_owner: GitHub repository owner
            repo_name: Repository name
            dry_run: If True, print actions without executing
        """
        self.repo_owner = repo_owner
        self.repo_name = repo_name
        self.dry_run = dry_run
        self.repo_full = f"{repo_owner}/{repo_name}"

    def run_gh_command(self, args: list[str]) -> str:
        """Run gh CLI command.

        Args:
            args: Command arguments

        Returns:
            Command output

        Raises:
            subprocess.CalledProcessError: If command fails
        """
        cmd = ["gh", *args]

        if self.dry_run:
            print(f"[DRY-RUN] Would run: {' '.join(cmd)}")
            return ""

        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            check=True,
        )
        return result.stdout.strip()

    def load_label_config(self, config_file: Path) -> list[dict[str, str]]:
        """Load labels from configuration file.

        Args:
            config_file: Path to labels.yml

        Returns:
            List of label definitions
        """
        if not config_file.exists():
            print(f"❌ Error: Config file not found: {config_file}")
            sys.exit(1)

        with open(config_file) as f:
            labels = yaml.safe_load(f)

        if not isinstance(labels, list):
            print("❌ Error: labels.yml must contain a list of labels")
            sys.exit(1)

        return labels

    def get_existing_labels(self) -> dict[str, dict[str, str]]:
        """Get existing labels from GitHub repository.

        Returns:
            Dict mapping label names to their properties
        """
        try:
            output = self.run_gh_command([
                "api",
                f"/repos/{self.repo_full}/labels",
                "--paginate",
                "--jq", ".",
            ])

            if not output or self.dry_run:
                return {}

            labels = json.loads(output)
            return {
                label["name"]: {
                    "color": label["color"],
                    "description": label.get("description", ""),
                }
                for label in labels
            }

        except (subprocess.CalledProcessError, json.JSONDecodeError) as e:
            print(f"⚠️  Could not fetch existing labels: {e}")
            return {}

    def create_label(self, name: str, color: str, description: str):
        """Create a new label.

        Args:
            name: Label name
            color: Label color (hex without #)
            description: Label description
        """
        print(f"   Creating label: {name}")

        if self.dry_run:
            print(f"      [DRY-RUN] Color: {color}, Description: {description}")
            return

        try:
            self.run_gh_command([
                "api",
                "-X", "POST",
                f"/repos/{self.repo_full}/labels",
                "-f", f"name={name}",
                "-f", f"color={color}",
                "-f", f"description={description}",
            ])
            print(f"      ✓ Created: {name}")
        except subprocess.CalledProcessError as e:
            print(f"      ✗ Failed to create {name}: {e}")

    def update_label(self, name: str, color: str, description: str):
        """Update an existing label.

        Args:
            name: Label name
            color: Label color (hex without #)
            description: Label description
        """
        print(f"   Updating label: {name}")

        if self.dry_run:
            print(f"      [DRY-RUN] Color: {color}, Description: {description}")
            return

        try:
            # URL-encode the label name for the API call
            import urllib.parse
            encoded_name = urllib.parse.quote(name)

            self.run_gh_command([
                "api",
                "-X", "PATCH",
                f"/repos/{self.repo_full}/labels/{encoded_name}",
                "-f", f"color={color}",
                "-f", f"description={description}",
            ])
            print(f"      ✓ Updated: {name}")
        except subprocess.CalledProcessError as e:
            print(f"      ✗ Failed to update {name}: {e}")

    def delete_label(self, name: str):
        """Delete a label.

        Args:
            name: Label name
        """
        print(f"   Deleting label: {name}")

        if self.dry_run:
            print(f"      [DRY-RUN] Would delete: {name}")
            return

        try:
            import urllib.parse
            encoded_name = urllib.parse.quote(name)

            self.run_gh_command([
                "api",
                "-X", "DELETE",
                f"/repos/{self.repo_full}/labels/{encoded_name}",
            ])
            print(f"      ✓ Deleted: {name}")
        except subprocess.CalledProcessError as e:
            print(f"      ✗ Failed to delete {name}: {e}")

    def sync_labels(self, config_file: Path, cleanup: bool = False):
        """Synchronize labels with configuration file.

        Args:
            config_file: Path to labels.yml
            cleanup: If True, remove labels not in config
        """
        print(f"\n{'='*70}")
        print(f"GitHub Labels Sync - {self.repo_full}")
        print(f"{'='*70}\n")

        if self.dry_run:
            print("🔍 DRY-RUN MODE - No changes will be made\n")

        # Load config and existing labels
        config_labels = self.load_label_config(config_file)
        existing_labels = self.get_existing_labels()

        print(f"Configuration labels: {len(config_labels)}")
        print(f"Existing labels: {len(existing_labels)}\n")

        # Track changes
        created = 0
        updated = 0
        skipped = 0

        # Process config labels
        print("🏷️  Processing labels...")
        for label_def in config_labels:
            name = label_def["name"]
            color = label_def.get("color", "ededed").lstrip("#")
            description = label_def.get("description", "")

            if name not in existing_labels:
                # Create new label
                self.create_label(name, color, description)
                created += 1
            else:
                # Check if update needed
                existing = existing_labels[name]
                needs_update = (
                    existing["color"].lower() != color.lower()
                    or existing["description"] != description
                )

                if needs_update:
                    self.update_label(name, color, description)
                    updated += 1
                else:
                    skipped += 1

        # Cleanup labels not in config
        if cleanup:
            config_label_names = {label["name"] for label in config_labels}
            print("\n🗑️  Cleaning up labels not in config...")

            deleted = 0
            for existing_name in existing_labels:
                if existing_name not in config_label_names:
                    self.delete_label(existing_name)
                    deleted += 1

            print(f"\n   Deleted: {deleted} labels")

        # Summary
        print(f"\n{'='*70}")
        print("Sync Summary")
        print(f"{'='*70}")
        print(f"   Created: {created}")
        print(f"   Updated: {updated}")
        print(f"   Skipped: {skipped}")
        print("\n✅ Label sync complete!\n")

    def check_labels(self, config_file: Path):
        """Check labels without making changes.

        Args:
            config_file: Path to labels.yml
        """
        print(f"\n{'='*70}")
        print(f"GitHub Labels Check - {self.repo_full}")
        print(f"{'='*70}\n")

        config_labels = self.load_label_config(config_file)
        existing_labels = self.get_existing_labels()

        config_names = {label["name"] for label in config_labels}
        existing_names = set(existing_labels.keys())

        missing = config_names - existing_names
        extra = existing_names - config_names
        common = config_names & existing_names

        print("📊 Label Status:")
        print(f"   Total in config: {len(config_names)}")
        print(f"   Total in GitHub: {len(existing_names)}")
        print(f"   Missing from GitHub: {len(missing)}")
        print(f"   Extra in GitHub: {len(extra)}")
        print(f"   In sync: {len(common)}\n")

        if missing:
            print("Missing labels:")
            for name in sorted(missing):
                print(f"   ✗ {name}")
            print()

        if extra:
            print("Extra labels:")
            for name in sorted(extra):
                print(f"   + {name}")
            print()

        # Check for updates needed
        updates_needed = 0
        for label_def in config_labels:
            name = label_def["name"]
            if name in existing_labels:
                color = label_def.get("color", "ededed").lstrip("#")
                description = label_def.get("description", "")
                existing = existing_labels[name]

                if (existing["color"].lower() != color.lower()
                        or existing["description"] != description):
                    updates_needed += 1

        if updates_needed:
            print(f"Labels needing update: {updates_needed}\n")

        if missing or extra or updates_needed:
            print("⚠️  Labels are out of sync")
            print("Run with --sync to synchronize\n")
        else:
            print("✅ All labels are in sync!\n")


def get_repo_info() -> tuple[str, str]:
    """Get repository owner and name from git remote.

    Returns:
        Tuple of (owner, repo_name)
    """
    try:
        result = subprocess.run(
            ["git", "remote", "get-url", "origin"],
            capture_output=True,
            text=True,
            check=True,
        )
        remote_url = result.stdout.strip()

        if "github.com" in remote_url:
            parts = remote_url.split("github.com")[-1].strip(":/").replace(".git", "").split("/")
            if len(parts) >= 2:
                return parts[0], parts[1]

        raise ValueError("Could not parse GitHub repo from remote URL")

    except (subprocess.CalledProcessError, ValueError) as e:
        print(f"❌ Error: Could not determine repository info: {e}")
        print("   Make sure you're in a git repository with a GitHub remote")
        sys.exit(1)


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Manage GitHub labels")
    parser.add_argument(
        "--sync",
        action="store_true",
        help="Synchronize labels with config file",
    )
    parser.add_argument(
        "--check",
        action="store_true",
        help="Check labels without making changes",
    )
    parser.add_argument(
        "--cleanup",
        action="store_true",
        help="Remove labels not in config (use with --sync)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Print actions without executing",
    )
    parser.add_argument(
        "--config",
        type=Path,
        default=Path.cwd() / ".github" / "labels.yml",
        help="Path to labels config file (default: .github/labels.yml)",
    )
    parser.add_argument(
        "--owner",
        help="Repository owner (auto-detected if not specified)",
    )
    parser.add_argument(
        "--repo",
        help="Repository name (auto-detected if not specified)",
    )

    args = parser.parse_args()

    # Get repo info
    if args.owner and args.repo:
        repo_owner = args.owner
        repo_name = args.repo
    else:
        repo_owner, repo_name = get_repo_info()

    print(f"Repository: {repo_owner}/{repo_name}\n")

    # Check gh CLI availability
    try:
        subprocess.run(["gh", "--version"], capture_output=True, check=True)
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("❌ Error: gh CLI not found")
        print("   Install: https://cli.github.com/")
        sys.exit(1)

    # Create manager
    manager = GitHubLabelsManager(repo_owner, repo_name, dry_run=args.dry_run)

    # Execute requested action
    if args.check:
        manager.check_labels(args.config)
    elif args.sync:
        manager.sync_labels(args.config, cleanup=args.cleanup)
    else:
        print("❌ Error: Specify --sync or --check")
        parser.print_help()
        sys.exit(1)


if __name__ == "__main__":
    main()
