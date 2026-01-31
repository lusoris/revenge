#!/usr/bin/env python3
"""GitHub Discussions automation script.

Configures GitHub Discussions with:
- Discussion categories (Ideas, Q&A, Announcements, Bugs)
- Discussion templates
- Auto-conversion rules

Uses gh CLI and GitHub API.

Requirements:
- gh CLI installed and authenticated
- Repository admin permissions
- Discussions feature enabled on repository

Usage:
    python scripts/automation/github_discussions.py [--dry-run]

Author: Automation System
Created: 2026-01-31
"""

import argparse
import json
import subprocess
import sys
from pathlib import Path


class GitHubDiscussionsManager:
    """Manage GitHub Discussions configuration."""

    def __init__(self, repo_owner: str, repo_name: str, dry_run: bool = False):
        """Initialize discussions manager.

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

    def enable_discussions(self) -> bool:
        """Enable GitHub Discussions on repository.

        Returns:
            True if discussions enabled, False otherwise

        Note:
            This requires GraphQL API access. gh CLI doesn't have
            a direct command for this. Returns False to indicate
            manual setup required.
        """
        print("📢 Checking Discussions status...")

        # Check if discussions are already enabled
        try:
            result = self.run_gh_command([
                "api",
                f"/repos/{self.repo_full}",
                "--jq", ".has_discussions",
            ])

            if not self.dry_run:
                has_discussions = result.strip().lower() == "true"
                if has_discussions:
                    print("   ✓ Discussions already enabled")
                    return True
                print("   ⚠️  Discussions not enabled")
                print("      Enable manually: Settings → Features → Discussions")
                return False

        except subprocess.CalledProcessError:
            print("   ⚠️  Could not check discussions status")

        return False

    def list_categories(self) -> list[dict[str, str]]:
        """List existing discussion categories.

        Returns:
            List of category dictionaries
        """
        print("   Listing discussion categories...")

        try:
            # Use GraphQL to list categories
            query = """
            query($owner: String!, $name: String!) {
                repository(owner: $owner, name: $name) {
                    discussionCategories(first: 10) {
                        nodes {
                            id
                            name
                            emoji
                            description
                        }
                    }
                }
            }
            """

            result = self.run_gh_command([
                "api", "graphql",
                "-f", f"owner={self.repo_owner}",
                "-f", f"name={self.repo_name}",
                "-f", f"query={query}",
            ])

            if self.dry_run or not result:
                return []

            data = json.loads(result)
            categories = data.get("data", {}).get("repository", {}).get("discussionCategories", {}).get("nodes", [])

            print(f"   ✓ Found {len(categories)} categories")
            return categories

        except (subprocess.CalledProcessError, json.JSONDecodeError, KeyError) as e:
            print(f"   ⚠️  Could not list categories: {e}")
            return []

    def create_templates(self):
        """Create discussion templates.

        Creates template files in .github/DISCUSSION_TEMPLATE/
        """
        print("\n📝 Creating discussion templates...")

        template_dir = Path(".github") / "DISCUSSION_TEMPLATE"
        template_dir.mkdir(parents=True, exist_ok=True)

        templates = {
            "ideas.yml": {
                "body": """---
name: 💡 Feature Idea
about: Suggest an idea for this project
labels: ["enhancement", "needs-triage"]
---

## Summary
<!-- Brief description of your idea -->

## Problem
<!-- What problem does this solve? -->

## Proposed Solution
<!-- How would this work? -->

## Alternatives Considered
<!-- What other approaches did you consider? -->

## Additional Context
<!-- Any other context, screenshots, or examples -->
""",
            },
            "question.yml": {
                "body": """---
name: ❓ Question
about: Ask a question about this project
labels: ["question"]
---

## Question
<!-- What would you like to know? -->

## Context
<!-- What are you trying to do? -->

## What I've Tried
<!-- What have you already tried? -->

## Environment
- OS:
- Version:
- Installation method:
""",
            },
            "show-and-tell.yml": {
                "body": """---
name: 🎨 Show and Tell
about: Share what you've built with Revenge
labels: ["showcase"]
---

## What I Built
<!-- Description of your project/setup -->

## Screenshots/Demo
<!-- Show us what you made! -->

## Tech Stack
<!-- What technologies did you use? -->

## Challenges & Solutions
<!-- What problems did you solve? -->
""",
            },
            "bug-report.yml": {
                "body": """---
name: 🐛 Bug Report (for Discussion)
about: Report a bug (will be converted to issue if confirmed)
labels: ["bug", "needs-triage"]
---

## Describe the Bug
<!-- Clear description of what's wrong -->

## Steps to Reproduce
1.
2.
3.

## Expected Behavior
<!-- What should happen -->

## Actual Behavior
<!-- What actually happens -->

## Environment
- OS:
- Version:
- Browser (if applicable):

## Logs
<!-- Paste relevant logs here -->

```
```
""",
            },
        }

        for filename, content in templates.items():
            template_file = template_dir / filename
            if template_file.exists() and not self.dry_run:
                print(f"   ⚠️  Template already exists: {filename}")
                continue

            if self.dry_run:
                print(f"   [DRY-RUN] Would create: {template_file}")
            else:
                template_file.write_text(content["body"])
                print(f"   ✓ Created template: {filename}")

    def print_manual_steps(self):
        """Print manual configuration steps."""
        print(f"\n{'='*70}")
        print("⚠️  Manual Configuration Required")
        print(f"{'='*70}\n")

        print("Since GitHub Discussions configuration via API is limited,")
        print("please complete the following steps manually:\n")

        print("1. Enable Discussions (if not already enabled):")
        print(f"   → https://github.com/{self.repo_full}/settings")
        print("   → Features → Discussions → Enable\n")

        print("2. Configure Discussion Categories:")
        print(f"   → https://github.com/{self.repo_full}/discussions/categories")
        print("   → Recommended categories:")
        print("      • 💡 Ideas (Ideas for features/improvements)")
        print("      • ❓ Q&A (Questions about using Revenge)")
        print("      • 📣 Announcements (Project updates)")
        print("      • 🐛 Bugs (Bug reports - convert to issues)")
        print("      • 🎨 Show and Tell (Share your setups)\n")

        print("3. Set up Auto-Convert Rules:")
        print("   → Category Settings → 🐛 Bugs")
        print("   → Enable 'Convert to issue' option\n")

        print("4. Pin Important Discussions:")
        print("   → Pin welcome/getting started discussions\n")

    def setup_discussions(self):
        """Set up complete GitHub Discussions."""
        print(f"\n{'='*70}")
        print(f"GitHub Discussions Setup - {self.repo_full}")
        print(f"{'='*70}\n")

        if self.dry_run:
            print("🔍 DRY-RUN MODE - No changes will be made\n")

        # Check if discussions are enabled
        discussions_enabled = self.enable_discussions()

        if discussions_enabled:
            # List existing categories
            categories = self.list_categories()
            if categories:
                print("\n   Existing categories:")
                for cat in categories:
                    emoji = cat.get("emoji", "")
                    name = cat.get("name", "")
                    print(f"      {emoji} {name}")

        # Create templates
        self.create_templates()

        # Print manual steps
        self.print_manual_steps()

        print(f"\n{'='*70}")
        print("✅ Discussion templates created!")
        print(f"{'='*70}\n")
        print("Templates location: .github/DISCUSSION_TEMPLATE/\n")


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
    parser = argparse.ArgumentParser(description="Configure GitHub Discussions")
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
    manager = GitHubDiscussionsManager(repo_owner, repo_name, dry_run=args.dry_run)
    manager.setup_discussions()


if __name__ == "__main__":
    main()
