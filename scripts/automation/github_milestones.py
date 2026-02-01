#!/usr/bin/env python3
"""GitHub Milestones automation script.

Automates milestone management:
- Auto-create milestones based on versioning
- Auto-assign issues/PRs to milestones
- Auto-close completed milestones
- Move open issues to next milestone

Uses gh CLI for GitHub API access.

Requirements:
- gh CLI installed and authenticated
- Repository admin permissions

Usage:
    python scripts/automation/github_milestones.py --create-next
    python scripts/automation/github_milestones.py --auto-assign
    python scripts/automation/github_milestones.py --close-completed
    python scripts/automation/github_milestones.py --move-open

Author: Automation System
Created: 2026-01-31
"""

import argparse
import json
import re
import subprocess
import sys
from datetime import datetime, timedelta


class GitHubMilestonesManager:
    """Manage GitHub repository milestones."""

    def __init__(self, repo_owner: str, repo_name: str, dry_run: bool = False):
        """Initialize milestones manager.

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

    def get_milestones(self, state: str = "all") -> list[dict]:
        """Get milestones from repository.

        Args:
            state: Milestone state (open, closed, all)

        Returns:
            List of milestone dicts
        """
        try:
            output = self.run_gh_command(
                [
                    "api",
                    f"/repos/{self.repo_full}/milestones",
                    "-f",
                    f"state={state}",
                    "--paginate",
                    "--jq",
                    ".",
                ]
            )

            if not output or self.dry_run:
                return []

            return json.loads(output)

        except (subprocess.CalledProcessError, json.JSONDecodeError) as e:
            print(f"⚠️  Could not fetch milestones: {e}")
            return []

    def get_latest_milestone_version(self) -> str | None:
        """Get the latest milestone version.

        Returns:
            Version string (e.g., "v1.2.0") or None
        """
        milestones = self.get_milestones(state="all")

        # Extract version numbers from milestone titles
        versions = []
        for milestone in milestones:
            title = milestone["title"]
            # Match version patterns: v1.2.3, 1.2.3, v1.2, etc.
            match = re.match(r"v?(\d+)\.(\d+)(?:\.(\d+))?", title)
            if match:
                major = int(match.group(1))
                minor = int(match.group(2))
                patch = int(match.group(3)) if match.group(3) else 0
                versions.append((major, minor, patch))

        if not versions:
            return None

        # Get the latest version
        latest = max(versions)
        return f"v{latest[0]}.{latest[1]}.{latest[2]}"

    def increment_version(self, version: str, bump: str = "minor") -> str:
        """Increment version number.

        Args:
            version: Current version (e.g., "v1.2.3")
            bump: Version component to bump (major, minor, patch)

        Returns:
            New version string
        """
        match = re.match(r"v?(\d+)\.(\d+)\.(\d+)", version)
        if not match:
            return "v1.0.0"

        major = int(match.group(1))
        minor = int(match.group(2))
        patch = int(match.group(3))

        if bump == "major":
            major += 1
            minor = 0
            patch = 0
        elif bump == "minor":
            minor += 1
            patch = 0
        else:  # patch
            patch += 1

        return f"v{major}.{minor}.{patch}"

    def create_milestone(
        self, title: str, description: str = "", due_date: str | None = None
    ):
        """Create a new milestone.

        Args:
            title: Milestone title
            description: Milestone description
            due_date: Due date (ISO 8601 format: YYYY-MM-DD)
        """
        print(f"   Creating milestone: {title}")

        if self.dry_run:
            print(f"      [DRY-RUN] Description: {description}")
            if due_date:
                print(f"      [DRY-RUN] Due date: {due_date}")
            return

        try:
            cmd = [
                "api",
                "-X",
                "POST",
                f"/repos/{self.repo_full}/milestones",
                "-f",
                f"title={title}",
            ]

            if description:
                cmd.extend(["-f", f"description={description}"])

            if due_date:
                cmd.extend(["-f", f"due_on={due_date}T23:59:59Z"])

            self.run_gh_command(cmd)
            print(f"      ✓ Created: {title}")
        except subprocess.CalledProcessError as e:
            print(f"      ✗ Failed to create {title}: {e}")

    def close_milestone(self, milestone_number: int):
        """Close a milestone.

        Args:
            milestone_number: Milestone number
        """
        print(f"   Closing milestone: #{milestone_number}")

        if self.dry_run:
            print(f"      [DRY-RUN] Would close milestone #{milestone_number}")
            return

        try:
            self.run_gh_command(
                [
                    "api",
                    "-X",
                    "PATCH",
                    f"/repos/{self.repo_full}/milestones/{milestone_number}",
                    "-f",
                    "state=closed",
                ]
            )
            print(f"      ✓ Closed milestone #{milestone_number}")
        except subprocess.CalledProcessError as e:
            print(f"      ✗ Failed to close milestone: {e}")

    def assign_issue_to_milestone(self, issue_number: int, milestone_number: int):
        """Assign an issue or PR to a milestone.

        Args:
            issue_number: Issue/PR number
            milestone_number: Milestone number
        """
        if self.dry_run:
            print(
                f"      [DRY-RUN] Would assign #{issue_number} to milestone #{milestone_number}"
            )
            return

        try:
            self.run_gh_command(
                [
                    "api",
                    "-X",
                    "PATCH",
                    f"/repos/{self.repo_full}/issues/{issue_number}",
                    "-f",
                    f"milestone={milestone_number}",
                ]
            )
        except subprocess.CalledProcessError as e:
            print(f"      ✗ Failed to assign #{issue_number}: {e}")

    def create_next_milestone(self, bump: str = "minor"):
        """Create the next milestone based on versioning.

        Args:
            bump: Version component to bump (major, minor, patch)
        """
        print(f"\n{'=' * 70}")
        print(f"Create Next Milestone - {self.repo_full}")
        print(f"{'=' * 70}\n")

        if self.dry_run:
            print("🔍 DRY-RUN MODE - No changes will be made\n")

        # Get latest version
        latest_version = self.get_latest_milestone_version()
        if latest_version:
            print(f"Latest milestone: {latest_version}")
            next_version = self.increment_version(latest_version, bump)
        else:
            print("No existing milestones found")
            next_version = "v1.0.0"

        print(f"Next milestone: {next_version}\n")

        # Calculate due date (3 months from now)
        due_date = (datetime.now() + timedelta(days=90)).strftime("%Y-%m-%d")

        # Create milestone
        description = f"Release {next_version}"
        self.create_milestone(next_version, description, due_date)

        print("\n✅ Milestone creation complete!\n")

    def close_completed_milestones(self):
        """Close milestones that have all issues closed."""
        print(f"\n{'=' * 70}")
        print(f"Close Completed Milestones - {self.repo_full}")
        print(f"{'=' * 70}\n")

        if self.dry_run:
            print("🔍 DRY-RUN MODE - No changes will be made\n")

        milestones = self.get_milestones(state="open")
        print(f"Open milestones: {len(milestones)}\n")

        closed_count = 0
        for milestone in milestones:
            title = milestone["title"]
            number = milestone["number"]
            open_issues = milestone["open_issues"]
            closed_issues = milestone["closed_issues"]

            if open_issues == 0 and closed_issues > 0:
                print(f"Milestone: {title}")
                print(f"   Open: {open_issues}, Closed: {closed_issues}")
                self.close_milestone(number)
                closed_count += 1

        print(f"\n✅ Closed {closed_count} milestones\n")

    def move_open_issues(self, from_milestone: str, to_milestone: str):
        """Move open issues from one milestone to another.

        Args:
            from_milestone: Source milestone title
            to_milestone: Target milestone title
        """
        print(f"\n{'=' * 70}")
        print(f"Move Open Issues - {self.repo_full}")
        print(f"{'=' * 70}\n")

        if self.dry_run:
            print("🔍 DRY-RUN MODE - No changes will be made\n")

        # Get milestones
        milestones = self.get_milestones(state="all")
        from_ms = next((m for m in milestones if m["title"] == from_milestone), None)
        to_ms = next((m for m in milestones if m["title"] == to_milestone), None)

        if not from_ms:
            print(f"❌ Source milestone not found: {from_milestone}")
            return

        if not to_ms:
            print(f"❌ Target milestone not found: {to_milestone}")
            return

        print(f"From: {from_milestone} (#{from_ms['number']})")
        print(f"To: {to_milestone} (#{to_ms['number']})\n")

        # Get issues from source milestone
        try:
            output = self.run_gh_command(
                [
                    "api",
                    f"/repos/{self.repo_full}/issues",
                    "-f",
                    f"milestone={from_ms['number']}",
                    "-f",
                    "state=open",
                    "--paginate",
                    "--jq",
                    ".",
                ]
            )

            issues = [] if not output or self.dry_run else json.loads(output)

            print(f"Open issues to move: {len(issues)}\n")

            moved_count = 0
            for issue in issues:
                number = issue["number"]
                title = issue["title"]
                print(f"   Moving #{number}: {title}")
                self.assign_issue_to_milestone(number, to_ms["number"])
                moved_count += 1

            print(f"\n✅ Moved {moved_count} issues\n")

        except (subprocess.CalledProcessError, json.JSONDecodeError) as e:
            print(f"❌ Error moving issues: {e}")


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
            parts = (
                remote_url.split("github.com")[-1]
                .strip(":/")
                .replace(".git", "")
                .split("/")
            )
            if len(parts) >= 2:
                return parts[0], parts[1]

        raise ValueError("Could not parse GitHub repo from remote URL")

    except (subprocess.CalledProcessError, ValueError) as e:
        print(f"❌ Error: Could not determine repository info: {e}")
        print("   Make sure you're in a git repository with a GitHub remote")
        sys.exit(1)


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Manage GitHub milestones")
    parser.add_argument(
        "--create-next",
        action="store_true",
        help="Create next milestone based on versioning",
    )
    parser.add_argument(
        "--close-completed",
        action="store_true",
        help="Close milestones with all issues closed",
    )
    parser.add_argument(
        "--move-open",
        nargs=2,
        metavar=("FROM", "TO"),
        help="Move open issues from FROM milestone to TO milestone",
    )
    parser.add_argument(
        "--bump",
        choices=["major", "minor", "patch"],
        default="minor",
        help="Version component to bump (default: minor)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Print actions without executing",
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
    manager = GitHubMilestonesManager(repo_owner, repo_name, dry_run=args.dry_run)

    # Execute requested action
    if args.create_next:
        manager.create_next_milestone(bump=args.bump)
    elif args.close_completed:
        manager.close_completed_milestones()
    elif args.move_open:
        manager.move_open_issues(args.move_open[0], args.move_open[1])
    else:
        print(
            "❌ Error: Specify an action (--create-next, --close-completed, or --move-open)"
        )
        parser.print_help()
        sys.exit(1)


if __name__ == "__main__":
    main()
