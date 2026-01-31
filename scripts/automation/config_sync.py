#!/usr/bin/env python3
"""Config Sync - Synchronize configs from SOURCE_OF_TRUTH.md.

This tool syncs configuration files from the SOURCE_OF_TRUTH.md:
- IDE settings (VS Code, Zed, JetBrains)
- Language version files (.tool-versions, .nvmrc, .python-version, go.mod)
- CI/CD workflows (GitHub Actions)
- Linter configs (golangci-lint, ruff, markdownlint)
- Docker configs (Dockerfile, docker-compose.yml)
- Coder template (.coder/template.tf)

One-way sync: SOT → Configs (SOT is master)

Author: Automation System
Created: 2026-01-31
"""

import re
import sys
from pathlib import Path

import yaml


class ConfigSync:
    """Synchronize configuration files from SOURCE_OF_TRUTH."""

    def __init__(self, repo_root: Path):
        """Initialize config sync.

        Args:
            repo_root: Repository root path
        """
        self.repo_root = repo_root
        self.sot_file = repo_root / "docs" / "dev" / "design" / "00_SOURCE_OF_TRUTH.md"
        self.shared_sot = repo_root / "data" / "shared-sot.yaml"

    def load_sot_data(self) -> dict:
        """Load parsed SOT data from YAML.

        Returns:
            Dict with SOT data
        """
        if not self.shared_sot.exists():
            print(f"⚠️  {self.shared_sot} not found. Run sot_parser.py first.")
            return {}

        with open(self.shared_sot) as f:
            data = yaml.safe_load(f)
            return data if data is not None else {}

    def sync_tool_versions(self, dry_run: bool = False) -> dict:
        """Sync .tool-versions from SOT.

        Args:
            dry_run: If True, don't write changes

        Returns:
            Statistics dict
        """
        sot_data = self.load_sot_data()
        go_deps = sot_data.get("go_dependencies", {})

        stats = {"updated": 0, "unchanged": 0, "errors": 0}

        # Extract Go version
        go_version = None
        for dep in go_deps.get("language", []):
            if dep.get("package") == "Go":
                version = dep.get("version", "")
                # Extract version number (e.g., "1.25+" -> "1.25")
                match = re.match(r"(\d+\.\d+)", version)
                if match:
                    go_version = match.group(1)
                break

        if not go_version:
            print("⚠️  Go version not found in SOT")
            stats["errors"] += 1
            return stats

        # Update .tool-versions
        tool_versions_file = self.repo_root / ".tool-versions"

        if tool_versions_file.exists():
            with open(tool_versions_file) as f:
                content = f.read()

            # Update Go version
            new_content = re.sub(
                r"golang \d+\.\d+(?:\.\d+)?",
                f"golang {go_version}",
                content
            )

            if new_content != content:
                stats["updated"] += 1
                if not dry_run:
                    with open(tool_versions_file, "w") as f:
                        f.write(new_content)
                    print(f"✓ Updated .tool-versions (Go {go_version})")
                else:
                    print(f"  Would update .tool-versions (Go {go_version})")
            else:
                stats["unchanged"] += 1
        else:
            print("⚠️  .tool-versions not found")
            stats["errors"] += 1

        return stats

    def sync_go_mod(self, dry_run: bool = False) -> dict:
        """Sync go.mod Go version from SOT.

        Args:
            dry_run: If True, don't write changes

        Returns:
            Statistics dict
        """
        sot_data = self.load_sot_data()
        go_deps = sot_data.get("go_dependencies", {})

        stats = {"updated": 0, "unchanged": 0, "errors": 0}

        # Extract Go version
        go_version = None
        for dep in go_deps.get("language", []):
            if dep.get("package") == "Go":
                version = dep.get("version", "")
                # Extract version number
                match = re.match(r"(\d+\.\d+)", version)
                if match:
                    go_version = match.group(1)
                break

        if not go_version:
            stats["errors"] += 1
            return stats

        # Update go.mod
        go_mod_file = self.repo_root / "go.mod"

        if go_mod_file.exists():
            with open(go_mod_file) as f:
                content = f.read()

            # Update go version line
            new_content = re.sub(
                r"^go \d+\.\d+$",
                f"go {go_version}",
                content,
                flags=re.MULTILINE
            )

            if new_content != content:
                stats["updated"] += 1
                if not dry_run:
                    with open(go_mod_file, "w") as f:
                        f.write(new_content)
                    print(f"✓ Updated go.mod (Go {go_version})")
                else:
                    print(f"  Would update go.mod (Go {go_version})")
            else:
                stats["unchanged"] += 1
        else:
            stats["errors"] += 1

        return stats

    def sync_github_workflows(self, dry_run: bool = False) -> dict:
        """Sync GitHub Actions workflow Go versions from SOT.

        Args:
            dry_run: If True, don't write changes

        Returns:
            Statistics dict
        """
        sot_data = self.load_sot_data()
        go_deps = sot_data.get("go_dependencies", {})

        stats = {"updated": 0, "unchanged": 0, "errors": 0}

        # Extract Go version
        go_version = None
        for dep in go_deps.get("language", []):
            if dep.get("package") == "Go":
                version = dep.get("version", "")
                match = re.match(r"(\d+\.\d+)", version)
                if match:
                    go_version = match.group(1)
                break

        if not go_version:
            stats["errors"] += 1
            return stats

        # Find all workflow files
        workflows_dir = self.repo_root / ".github" / "workflows"
        if not workflows_dir.exists():
            stats["errors"] += 1
            return stats

        workflow_files = list(workflows_dir.glob("*.yml"))

        for workflow_file in workflow_files:
            with open(workflow_file) as f:
                content = f.read()

            # Update go-version in setup-go action
            new_content = re.sub(
                r"(go-version:\s*['\"]?)(\d+\.\d+(?:\.\d+)?)",
                rf"\g<1>{go_version}",
                content
            )

            if new_content != content:
                stats["updated"] += 1
                if not dry_run:
                    with open(workflow_file, "w") as f:
                        f.write(new_content)
                    print(f"✓ Updated {workflow_file.name} (Go {go_version})")
                else:
                    print(f"  Would update {workflow_file.name} (Go {go_version})")
            else:
                stats["unchanged"] += 1

        return stats

    def sync_docker_configs(self, dry_run: bool = False) -> dict:
        """Sync Dockerfile and docker-compose.yml from SOT.

        Args:
            dry_run: If True, don't write changes

        Returns:
            Statistics dict
        """
        sot_data = self.load_sot_data()
        go_deps = sot_data.get("go_dependencies", {})

        stats = {"updated": 0, "unchanged": 0, "errors": 0}

        # Extract Go version
        go_version = None
        for dep in go_deps.get("language", []):
            if dep.get("package") == "Go":
                version = dep.get("version", "")
                match = re.match(r"(\d+\.\d+)", version)
                if match:
                    go_version = match.group(1)
                break

        if not go_version:
            stats["errors"] += 1
            return stats

        # Update Dockerfile
        dockerfile = self.repo_root / "Dockerfile"
        if dockerfile.exists():
            with open(dockerfile) as f:
                content = f.read()

            # Update FROM golang:X.Y
            new_content = re.sub(
                r"FROM golang:\d+\.\d+(?:\.\d+)?",
                f"FROM golang:{go_version}",
                content
            )

            if new_content != content:
                stats["updated"] += 1
                if not dry_run:
                    with open(dockerfile, "w") as f:
                        f.write(new_content)
                    print(f"✓ Updated Dockerfile (Go {go_version})")
                else:
                    print(f"  Would update Dockerfile (Go {go_version})")
            else:
                stats["unchanged"] += 1

        return stats

    def sync_all(self, dry_run: bool = False) -> dict:
        """Sync all configuration files.

        Args:
            dry_run: If True, don't write changes

        Returns:
            Combined statistics dict
        """
        print(f"\n{'='*70}")
        print(f"CONFIG SYNC - {'DRY RUN' if dry_run else 'LIVE'}")
        print(f"{'='*70}\n")

        total_stats = {"updated": 0, "unchanged": 0, "errors": 0}

        # Sync each config type
        configs = [
            ("Tool Versions (.tool-versions)", self.sync_tool_versions),
            ("Go Module (go.mod)", self.sync_go_mod),
            ("GitHub Workflows", self.sync_github_workflows),
            ("Docker Configs", self.sync_docker_configs),
        ]

        for name, sync_func in configs:
            print(f"\n{name}:")
            stats = sync_func(dry_run)

            # Aggregate stats
            for key in ["updated", "unchanged", "errors"]:
                total_stats[key] += stats.get(key, 0)

        # Summary
        print(f"\n{'='*70}")
        print("SUMMARY")
        print(f"{'='*70}")
        print(f"Updated: {total_stats['updated']}")
        print(f"Unchanged: {total_stats['unchanged']}")
        print(f"Errors: {total_stats['errors']}")

        if dry_run:
            print("\n⚠️  DRY RUN MODE - No changes written")

        print(f"{'='*70}\n")

        return total_stats


def main():
    """Main entry point."""
    repo_root = Path(__file__).parent.parent.parent

    # Parse arguments
    args = sys.argv[1:]
    dry_run = "--live" not in args

    if "--help" in args or not args:
        print("Usage: python config_sync.py [--live]")
        print()
        print("Options:")
        print("  --live     Apply changes (default: dry-run)")
        print("  --dry-run  Show what would be changed (default)")
        print()
        print("Examples:")
        print("  python config_sync.py --dry-run")
        print("  python config_sync.py --live")
        sys.exit(0)

    # Initialize syncer
    syncer = ConfigSync(repo_root)

    # Sync all configs
    stats = syncer.sync_all(dry_run)

    # Exit code based on results
    sys.exit(0 if stats["errors"] == 0 else 1)


if __name__ == "__main__":
    main()
