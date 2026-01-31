#!/usr/bin/env python3
"""GitHub Projects automation script.

Configures GitHub Projects V2 with:
- Project board with workflow columns
- Automation rules for issue/PR management
- Custom fields (Priority, Effort, Module)
- Multiple views (Board, Table, Roadmap)

Uses gh CLI for GitHub API access.

Requirements:
- gh CLI installed and authenticated
- Repository owner/admin permissions

Usage:
    python scripts/automation/github_projects.py [--dry-run]

Author: Automation System
Created: 2026-01-31
"""

import argparse
import json
import subprocess
import sys


class GitHubProjectsManager:
    """Manage GitHub Projects V2 configuration."""

    def __init__(self, repo_owner: str, repo_name: str, dry_run: bool = False):
        """Initialize project manager.

        Args:
            repo_owner: GitHub repository owner (user or org)
            repo_name: Repository name
            dry_run: If True, print commands without executing
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

    def create_project(self, title: str, _description: str) -> str:
        """Create a new GitHub Project.

        Args:
            title: Project title
            _description: Project description (not used - gh CLI limitation)

        Returns:
            Project number
        """
        print(f"📋 Creating project: {title}")

        output = self.run_gh_command([
            "project", "create",
            "--owner", self.repo_owner,
            "--title", title,
            "--format", "json",
        ])

        if self.dry_run:
            return "1"  # Dummy project number for dry-run

        project_data = json.loads(output)
        project_number = project_data["number"]

        print(f"   ✓ Created project #{project_number}")
        return str(project_number)

    def add_field(self, project_number: str, field_name: str, field_type: str, options: list[str] | None = None):
        """Add custom field to project.

        Args:
            project_number: Project number
            field_name: Field name
            field_type: Field type (TEXT, SINGLE_SELECT, etc.)
            options: Field options for SINGLE_SELECT
        """
        print(f"   Adding field: {field_name} ({field_type})")

        cmd = [
            "project", "field-create", project_number,
            "--owner", self.repo_owner,
            "--name", field_name,
            "--data-type", field_type,
        ]

        if options and field_type == "SINGLE_SELECT":
            for option in options:
                cmd.extend(["--single-select-option", option])

        self.run_gh_command(cmd)
        print(f"   ✓ Added field: {field_name}")

    def configure_workflow(self, _project_number: str):
        """Configure project workflow columns.

        Args:
            _project_number: Project number (not used - manual setup required)
        """
        print("   Configuring workflow columns...")

        # GitHub Projects V2 uses "Status" field by default
        # We'll configure the status field with our workflow states
        statuses = ["Backlog", "Todo", "In Progress", "Review", "Done"]

        for status in statuses:
            print(f"      • {status}")

        # Note: gh CLI doesn't directly support editing default Status field options
        # This would need to be done via GraphQL API
        # For now, we'll document this as a manual step

        print("   ⚠️  Status field configuration requires manual setup via GitHub UI")
        print("       Configure Status field with: Backlog, Todo, In Progress, Review, Done")

    def add_items_to_project(self, project_number: str):
        """Add existing issues/PRs to project.

        Args:
            project_number: Project number
        """
        print("   Adding existing issues to project...")

        # Get open issues
        issues_output = self.run_gh_command([
            "issue", "list",
            "--repo", self.repo_full,
            "--state", "open",
            "--limit", "100",
            "--json", "number",
        ])

        if self.dry_run or not issues_output:
            print("   ℹ️  No issues to add (or dry-run)")
            return

        issues = json.loads(issues_output)

        for issue in issues:
            issue_number = issue["number"]
            try:
                self.run_gh_command([
                    "project", "item-add", project_number,
                    "--owner", self.repo_owner,
                    "--url", f"https://github.com/{self.repo_full}/issues/{issue_number}",
                ])
                print(f"      ✓ Added issue #{issue_number}")
            except subprocess.CalledProcessError as e:
                print(f"      ⚠️  Failed to add issue #{issue_number}: {e}")

    def setup_project(self):
        """Set up complete GitHub Project."""
        print(f"\n{'='*70}")
        print(f"GitHub Projects Setup - {self.repo_full}")
        print(f"{'='*70}\n")

        if self.dry_run:
            print("🔍 DRY-RUN MODE - No changes will be made\n")

        # Create project
        project_number = self.create_project(
            title="Revenge Development",
            _description="Track development progress for Revenge media server",
        )

        # Add custom fields
        print("\n📝 Adding custom fields...")

        self.add_field(
            project_number,
            "Priority",
            "SINGLE_SELECT",
            ["🔴 Critical", "🟠 High", "🟡 Medium", "🟢 Low"],
        )

        self.add_field(
            project_number,
            "Effort",
            "SINGLE_SELECT",
            ["XS (< 1hr)", "S (1-4hrs)", "M (1-2 days)", "L (3-5 days)", "XL (1-2 weeks)"],
        )

        self.add_field(
            project_number,
            "Module",
            "SINGLE_SELECT",
            ["Core", "Frontend", "Backend", "API", "Infrastructure", "Docs", "Testing"],
        )

        # Configure workflow
        print("\n⚙️  Configuring workflow...")
        self.configure_workflow(project_number)

        # Add existing items
        print("\n📌 Adding items to project...")
        self.add_items_to_project(project_number)

        print(f"\n{'='*70}")
        print("✅ Project setup complete!")
        print(f"{'='*70}\n")
        print(f"Project URL: https://github.com/orgs/{self.repo_owner}/projects/{project_number}")
        print("\n⚠️  Manual steps required:")
        print("   1. Configure Status field options in GitHub UI")
        print("   2. Set up automation rules (auto-add, auto-move)")
        print("   3. Create additional views (Table, Roadmap)")


def get_repo_info() -> tuple[str, str]:
    """Get repository owner and name from git remote.

    Returns:
        Tuple of (owner, repo_name)
    """
    try:
        # Get remote URL
        result = subprocess.run(
            ["git", "remote", "get-url", "origin"],
            capture_output=True,
            text=True,
            check=True,
        )
        remote_url = result.stdout.strip()

        # Parse owner and repo from URL
        # Handles: git@github.com:owner/repo.git or https://github.com/owner/repo.git
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
    parser = argparse.ArgumentParser(description="Configure GitHub Projects")
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Print commands without executing",
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

    # Run setup
    manager = GitHubProjectsManager(repo_owner, repo_name, dry_run=args.dry_run)
    manager.setup_project()


if __name__ == "__main__":
    main()
