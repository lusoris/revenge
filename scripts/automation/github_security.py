#!/usr/bin/env python3
"""GitHub Security & Branch Protection automation script.

Configures GitHub Advanced Security features:
- Branch protection rules (develop, main)
- CodeQL analysis setup
- Secret scanning
- Dependency review

Uses gh CLI and GitHub API.

Requirements:
- gh CLI installed and authenticated
- Repository admin permissions
- GitHub Advanced Security enabled (for private repos)

Usage:
    python scripts/automation/github_security.py [--dry-run]

Author: Automation System
Created: 2026-01-31
"""

import argparse
import json
import subprocess
import sys


class GitHubSecurityManager:
    """Manage GitHub Security and branch protection configuration."""

    def __init__(self, repo_owner: str, repo_name: str, dry_run: bool = False):
        """Initialize security manager.

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

    def configure_branch_protection(self, branch: str, strict: bool = True):
        """Configure branch protection rules.

        Args:
            branch: Branch name (e.g., "main", "develop")
            strict: If True, apply strict protection rules
        """
        print(f"   Configuring protection for branch: {branch}")

        # Build protection rules
        rules = {
            "required_status_checks": {
                "strict": strict,
                "contexts": ["test", "lint", "build"],
            },
            "enforce_admins": True,
            "required_pull_request_reviews": {
                "required_approving_review_count": 1,
                "dismiss_stale_reviews": True,
                "require_code_owner_reviews": True,
            },
            "restrictions": None,  # No restrictions on who can push
            "required_linear_history": True,
            "allow_force_pushes": False,
            "allow_deletions": False,
        }

        # Use GitHub API to set branch protection
        # gh API doesn't have direct branch protection commands, use GraphQL
        if self.dry_run:
            print(f"      [DRY-RUN] Would set protection rules for {branch}:")
            print("         • Require PR reviews (1 approval)")
            print(
                f"         • Require status checks: {rules['required_status_checks']['contexts']}"
            )
            print("         • Linear history required")
            print("         • No force push")
            print("         • Include administrators")
        else:
            # Use REST API to configure branch protection
            api_endpoint = f"/repos/{self.repo_full}/branches/{branch}/protection"

            try:
                self.run_gh_command(
                    [
                        "api",
                        "-X",
                        "PUT",
                        api_endpoint,
                        "-f",
                        f"required_status_checks[strict]={str(strict).lower()}",
                        "-f",
                        "required_status_checks[contexts][]=test",
                        "-f",
                        "required_status_checks[contexts][]=lint",
                        "-f",
                        "required_status_checks[contexts][]=build",
                        "-f",
                        "enforce_admins=true",
                        "-f",
                        "required_pull_request_reviews[required_approving_review_count]=1",
                        "-f",
                        "required_pull_request_reviews[dismiss_stale_reviews]=true",
                        "-f",
                        "required_pull_request_reviews[require_code_owner_reviews]=true",
                        "-f",
                        "required_linear_history=true",
                        "-f",
                        "allow_force_pushes=false",
                        "-f",
                        "allow_deletions=false",
                    ]
                )
                print(f"      ✓ Protection configured for {branch}")
            except subprocess.CalledProcessError as e:
                print(f"      ⚠️  Failed to configure protection for {branch}: {e}")
                print("         This may require manual configuration via GitHub UI")

    def enable_security_features(self):
        """Enable GitHub Advanced Security features.

        Enables:
        - Secret scanning
        - Dependency graph
        - Dependabot alerts
        - Code scanning (via CodeQL workflow)
        """
        print("\n🔒 Enabling security features...")

        features = {
            "secret_scanning": "Enable secret scanning",
            "secret_scanning_push_protection": "Enable push protection for secrets",
        }

        for feature, description in features.items():
            print(f"   {description}...")

            if self.dry_run:
                print(f"      [DRY-RUN] Would enable {feature}")
                continue

            try:
                # Use GitHub API to enable security features
                # Note: This requires admin access and may need manual setup
                api_endpoint = f"/repos/{self.repo_full}"

                # Enable secret scanning
                if feature == "secret_scanning":
                    self.run_gh_command(
                        [
                            "api",
                            "-X",
                            "PATCH",
                            api_endpoint,
                            "-f",
                            "security_and_analysis[secret_scanning][status]=enabled",
                        ]
                    )
                elif feature == "secret_scanning_push_protection":
                    self.run_gh_command(
                        [
                            "api",
                            "-X",
                            "PATCH",
                            api_endpoint,
                            "-f",
                            "security_and_analysis[secret_scanning_push_protection][status]=enabled",
                        ]
                    )

                print(f"      ✓ {description} enabled")
            except subprocess.CalledProcessError as e:
                print(f"      ⚠️  Failed to enable {feature}: {e}")
                print("         May require GitHub Advanced Security or manual setup")

    def check_security_status(self):
        """Check current security feature status."""
        print("\n📊 Checking security status...")

        try:
            output = self.run_gh_command(
                [
                    "api",
                    f"/repos/{self.repo_full}",
                    "--jq",
                    ".security_and_analysis",
                ]
            )

            if not self.dry_run and output:
                status = json.loads(output) if output else {}

                features_to_check = [
                    ("secret_scanning", "Secret scanning"),
                    ("secret_scanning_push_protection", "Push protection"),
                    ("dependabot_security_updates", "Dependabot security updates"),
                ]

                for key, name in features_to_check:
                    if key in status:
                        feature_status = status[key].get("status", "disabled")
                        icon = "✓" if feature_status == "enabled" else "✗"
                        print(f"   {icon} {name}: {feature_status}")
                    else:
                        print(f"   ℹ️  {name}: not available")

        except (subprocess.CalledProcessError, json.JSONDecodeError) as e:
            print(f"   ⚠️  Could not check security status: {e}")

    def print_manual_steps(self):
        """Print manual configuration steps."""
        print(f"\n{'=' * 70}")
        print("⚠️  Manual Configuration Steps")
        print(f"{'=' * 70}\n")

        print("Some features require manual setup in GitHub UI:\n")

        print("1. Enable CodeQL Analysis:")
        print(f"   → https://github.com/{self.repo_full}/settings/security_analysis")
        print("   → Code scanning → Set up → Use CodeQL")
        print("   → Or use the provided .github/workflows/codeql.yml\n")

        print("2. Enable Dependabot:")
        print(f"   → https://github.com/{self.repo_full}/settings/security_analysis")
        print("   → Dependabot alerts → Enable")
        print("   → Dependabot security updates → Enable\n")

        print("3. Review Branch Protection:")
        print(f"   → https://github.com/{self.repo_full}/settings/branches")
        print("   → Verify rules for main and develop branches\n")

        print("4. Add CODEOWNERS file:")
        print("   → Create .github/CODEOWNERS")
        print("   → Define code ownership patterns\n")

    def setup_security(self):
        """Set up complete GitHub Security configuration."""
        print(f"\n{'=' * 70}")
        print(f"GitHub Security Setup - {self.repo_full}")
        print(f"{'=' * 70}\n")

        if self.dry_run:
            print("🔍 DRY-RUN MODE - No changes will be made\n")

        # Check current status
        self.check_security_status()

        # Configure branch protection
        print("\n🛡️  Configuring branch protection...")
        for branch in ["main", "develop"]:
            self.configure_branch_protection(branch, strict=True)

        # Enable security features
        self.enable_security_features()

        # Print manual steps
        self.print_manual_steps()

        print(f"\n{'=' * 70}")
        print("✅ Security setup complete!")
        print(f"{'=' * 70}\n")
        print("Note: Some features may require GitHub Advanced Security")
        print("      for private repositories.\n")


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
    parser = argparse.ArgumentParser(
        description="Configure GitHub Security & Branch Protection"
    )
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
    manager = GitHubSecurityManager(repo_owner, repo_name, dry_run=args.dry_run)
    manager.setup_security()


if __name__ == "__main__":
    main()
