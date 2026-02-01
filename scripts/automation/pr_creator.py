#!/usr/bin/env python3
"""PR Creator - Create batched PRs for documentation updates.

This tool creates pull requests for automated documentation updates with:
- Batching by trigger type (SOT update, config sync, etc.)
- Loop prevention (bot user check, cooldown lock)
- Auto-merge for docs-only changes
- Proper commit messages and PR descriptions

Author: Automation System
Created: 2026-01-31
"""

import contextlib
import os
import subprocess
import sys
from datetime import datetime, timedelta
from pathlib import Path


class LoopPrevention:
    """Triple-safety loop prevention system."""

    def __init__(self, repo_root: Path):
        """Initialize loop prevention.

        Args:
            repo_root: Repository root path
        """
        self.repo_root = repo_root
        self.lock_file = repo_root / ".automation-lock"
        self.lock_timeout = timedelta(hours=1)
        self.bot_username = "revenge-bot"

    def is_bot_commit(self) -> bool:
        """Check if last commit was by the bot.

        Returns:
            True if last commit author is revenge-bot
        """
        try:
            result = subprocess.run(
                ["git", "log", "-1", "--pretty=%an"],
                cwd=self.repo_root,
                capture_output=True,
                text=True,
                check=True,
            )
            author = result.stdout.strip()
            return author == self.bot_username
        except subprocess.CalledProcessError:
            return False

    def is_locked(self) -> tuple[bool, str | None]:
        """Check if automation is locked (cooldown period).

        Returns:
            Tuple of (is_locked, reason)
        """
        if not self.lock_file.exists():
            return False, None

        # Read lock file
        try:
            with open(self.lock_file) as f:
                lock_time_str = f.read().strip()
                lock_time = datetime.fromisoformat(lock_time_str)
        except (ValueError, OSError):
            # Invalid lock file, remove it
            self.lock_file.unlink(missing_ok=True)
            return False, None

        # Check if lock has expired
        now = datetime.now()
        if now - lock_time > self.lock_timeout:
            # Lock expired, remove it
            self.lock_file.unlink(missing_ok=True)
            return False, None

        # Lock is still active
        time_remaining = self.lock_timeout - (now - lock_time)
        minutes_remaining = int(time_remaining.total_seconds() / 60)
        return True, f"Cooldown active ({minutes_remaining} minutes remaining)"

    def create_lock(self):
        """Create automation lock file."""
        with open(self.lock_file, "w") as f:
            f.write(datetime.now().isoformat())

    def remove_lock(self):
        """Remove automation lock file."""
        self.lock_file.unlink(missing_ok=True)

    def can_proceed(self) -> tuple[bool, str | None]:
        """Check if automation can proceed.

        Returns:
            Tuple of (can_proceed, reason_if_not)
        """
        # Check 1: Bot user check
        if self.is_bot_commit():
            return False, "Last commit was by revenge-bot (loop prevention)"

        # Check 2: Cooldown lock
        is_locked, reason = self.is_locked()
        if is_locked:
            return False, reason

        return True, None


class PRCreator:
    """Create pull requests for automated documentation updates."""

    def __init__(self, repo_root: Path, bot_token: str | None = None):
        """Initialize PR creator.

        Args:
            repo_root: Repository root path
            bot_token: GitHub token for bot user (from env var if not provided)
        """
        self.repo_root = repo_root
        self.bot_token = bot_token or os.environ.get("GITHUB_TOKEN")
        self.loop_prevention = LoopPrevention(repo_root)

    def check_prerequisites(self) -> tuple[bool, str | None]:
        """Check if all prerequisites are met.

        Returns:
            Tuple of (success, error_message)
        """
        # Check loop prevention
        can_proceed, reason = self.loop_prevention.can_proceed()
        if not can_proceed:
            return False, f"Loop prevention: {reason}"

        # Check gh CLI
        try:
            subprocess.run(
                ["gh", "--version"],
                capture_output=True,
                check=True,
            )
        except (subprocess.CalledProcessError, FileNotFoundError):
            return False, "gh CLI not installed or not in PATH"

        # Check authentication
        try:
            subprocess.run(
                ["gh", "auth", "status"],
                capture_output=True,
                check=True,
            )
        except subprocess.CalledProcessError:
            return False, "gh CLI not authenticated (run: gh auth login)"

        return True, None

    def get_changed_files(self) -> list[str]:
        """Get list of changed files.

        Returns:
            List of changed file paths
        """
        try:
            result = subprocess.run(
                ["git", "status", "--porcelain"],
                cwd=self.repo_root,
                capture_output=True,
                text=True,
                check=True,
            )
            lines = result.stdout.split("\n")
            # Filter empty lines and extract filenames (skip first 3 chars: XY + space)
            files = [line[3:] for line in lines if line.strip()]
            return files
        except subprocess.CalledProcessError:
            return []

    def is_docs_only_change(self, files: list[str]) -> bool:
        """Check if changes are docs-only.

        Args:
            files: List of changed files

        Returns:
            True if all changes are in docs/, data/, templates/, schemas/
        """
        docs_prefixes = (
            "docs/",
            "data/",
            "templates/",
            "schemas/",
            ".github/workflows/doc-",
        )
        return all(any(f.startswith(prefix) for prefix in docs_prefixes) for f in files)

    def create_branch(self, branch_name: str) -> bool:
        """Create and checkout new branch.

        Args:
            branch_name: Name of branch to create

        Returns:
            True if successful
        """
        try:
            subprocess.run(
                ["git", "checkout", "-b", branch_name],
                cwd=self.repo_root,
                capture_output=True,
                check=True,
            )
            return True
        except subprocess.CalledProcessError:
            return False

    def commit_changes(self, message: str) -> bool:
        """Commit all changes.

        Args:
            message: Commit message

        Returns:
            True if successful
        """
        try:
            # Stage all changes
            subprocess.run(
                ["git", "add", "-A"],
                cwd=self.repo_root,
                check=True,
            )

            # Commit
            subprocess.run(
                ["git", "commit", "-m", message],
                cwd=self.repo_root,
                check=True,
            )
            return True
        except subprocess.CalledProcessError:
            return False

    def push_branch(self, branch_name: str) -> bool:
        """Push branch to remote.

        Args:
            branch_name: Name of branch to push

        Returns:
            True if successful
        """
        try:
            subprocess.run(
                ["git", "push", "-u", "origin", branch_name],
                cwd=self.repo_root,
                check=True,
            )
            return True
        except subprocess.CalledProcessError:
            return False

    def create_pr(
        self,
        title: str,
        body: str,
        branch_name: str,
        base_branch: str = "develop",
        auto_merge: bool = False,
    ) -> tuple[bool, str | None]:
        """Create pull request.

        Args:
            title: PR title
            body: PR body/description
            branch_name: Source branch
            base_branch: Target branch (default: develop)
            auto_merge: Whether to enable auto-merge

        Returns:
            Tuple of (success, pr_url)
        """
        try:
            # Create PR
            cmd = [
                "gh",
                "pr",
                "create",
                "--base",
                base_branch,
                "--head",
                branch_name,
                "--title",
                title,
                "--body",
                body,
            ]

            if auto_merge:
                cmd.append("--label")
                cmd.append("automerge")

            result = subprocess.run(
                cmd,
                cwd=self.repo_root,
                capture_output=True,
                text=True,
                check=True,
            )

            pr_url = result.stdout.strip()

            # Enable auto-merge if requested (ignore failure - PR still created)
            if auto_merge:
                with contextlib.suppress(subprocess.CalledProcessError):
                    subprocess.run(
                        ["gh", "pr", "merge", pr_url, "--auto", "--squash"],
                        cwd=self.repo_root,
                        check=True,
                    )

            return True, pr_url

        except subprocess.CalledProcessError as e:
            return False, str(e)

    def create_doc_update_pr(
        self,
        trigger_type: str,
        changed_files: list[str],
        dry_run: bool = False,
    ) -> tuple[bool, str | None]:
        """Create PR for documentation update.

        Args:
            trigger_type: Type of trigger (sot_update, config_sync, manual)
            changed_files: List of changed files
            dry_run: If True, don't actually create PR

        Returns:
            Tuple of (success, pr_url_or_error)
        """
        # Check prerequisites
        success, error = self.check_prerequisites()
        if not success:
            return False, error

        # Determine PR details based on trigger type
        timestamp = datetime.now().strftime("%Y%m%d-%H%M%S")

        if trigger_type == "sot_update":
            branch_name = f"docs/sot-update-{timestamp}"
            title = "docs: update from SOURCE_OF_TRUTH changes"
            body = f"""## Documentation Update

**Trigger**: SOURCE_OF_TRUTH.md update

**Changes**: {len(changed_files)} file(s)

**Files Updated**:
{"".join(f"- {f}\n" for f in changed_files[:20])}
{"..." if len(changed_files) > 20 else ""}

**Generated by**: Documentation automation system
**Can auto-merge**: {"Yes" if self.is_docs_only_change(changed_files) else "No (non-docs changes detected)"}

---

🤖 Generated with [Claude Code](https://claude.com/claude-code)
"""
        elif trigger_type == "config_sync":
            branch_name = f"config/sync-{timestamp}"
            title = "config: sync configs from SOURCE_OF_TRUTH"
            body = f"""## Configuration Sync

**Trigger**: Automatic config synchronization

**Changes**: {len(changed_files)} file(s)

**Files Updated**:
{"".join(f"- {f}\n" for f in changed_files[:20])}
{"..." if len(changed_files) > 20 else ""}

**Generated by**: Config sync automation

---

🤖 Generated with [Claude Code](https://claude.com/claude-code)
"""
        else:  # manual
            branch_name = f"docs/manual-update-{timestamp}"
            title = "docs: manual documentation update"
            body = f"""## Manual Documentation Update

**Changes**: {len(changed_files)} file(s)

**Files Updated**:
{"".join(f"- {f}\n" for f in changed_files[:20])}
{"..." if len(changed_files) > 20 else ""}

---

🤖 Generated with [Claude Code](https://claude.com/claude-code)
"""

        if dry_run:
            print(f"\n{'=' * 70}")
            print("DRY RUN - Would create PR:")
            print(f"{'=' * 70}")
            print(f"Branch: {branch_name}")
            print(f"Title: {title}")
            print(f"Body:\n{body}")
            print(f"Auto-merge: {self.is_docs_only_change(changed_files)}")
            print(f"{'=' * 70}\n")
            return True, "dry-run"

        # Create branch
        if not self.create_branch(branch_name):
            return False, "Failed to create branch"

        # Commit changes
        commit_msg = (
            title + "\n\nCo-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
        )
        if not self.commit_changes(commit_msg):
            return False, "Failed to commit changes"

        # Push branch
        if not self.push_branch(branch_name):
            return False, "Failed to push branch"

        # Create PR with auto-merge if docs-only
        auto_merge = self.is_docs_only_change(changed_files)
        success, result = self.create_pr(
            title=title,
            body=body,
            branch_name=branch_name,
            auto_merge=auto_merge,
        )

        if success:
            # Create lock to prevent loops
            self.loop_prevention.create_lock()
            return True, result
        return False, result


def main():
    """Main entry point."""
    repo_root = Path(__file__).parent.parent.parent

    # Parse arguments
    args = sys.argv[1:]

    trigger_type = "manual"
    dry_run = True

    for arg in args:
        if arg == "--sot-update":
            trigger_type = "sot_update"
        elif arg == "--config-sync":
            trigger_type = "config_sync"
        elif arg == "--manual":
            trigger_type = "manual"
        elif arg == "--live":
            dry_run = False

    if "--live" not in args and "--dry-run" not in args:
        dry_run = True  # Default to dry-run

    # Show usage
    if not args or "--help" in args:
        print("Usage: python pr_creator.py <trigger> [--live]")
        print()
        print("Triggers:")
        print("  --sot-update    PR for SOURCE_OF_TRUTH.md changes")
        print("  --config-sync   PR for config synchronization")
        print("  --manual        PR for manual updates")
        print()
        print("Options:")
        print("  --live          Create actual PR (default: dry-run)")
        print("  --dry-run       Show what would be done (default)")
        print()
        print("Examples:")
        print("  python pr_creator.py --sot-update --dry-run")
        print("  python pr_creator.py --manual --live")
        sys.exit(0)

    # Initialize PR creator
    creator = PRCreator(repo_root)

    # Get changed files
    changed_files = creator.get_changed_files()
    if not changed_files:
        print("No changes to commit")
        sys.exit(0)

    print(f"Found {len(changed_files)} changed file(s)")

    # Create PR
    success, result = creator.create_doc_update_pr(
        trigger_type=trigger_type,
        changed_files=changed_files,
        dry_run=dry_run,
    )

    if success:
        if not dry_run:
            print(f"\n✅ PR created: {result}")
        sys.exit(0)
    else:
        print(f"\n❌ Failed to create PR: {result}")
        sys.exit(1)


if __name__ == "__main__":
    main()
