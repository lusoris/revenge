#!/usr/bin/env python3
"""Dependency update automation script.

Checks for and optionally updates dependencies across ecosystems:
- Go modules (go list -u -m all)
- npm packages (npm outdated)
- Python packages (pip list --outdated)

Supports:
- Check mode (list outdated)
- Update mode (automatic updates)
- Test after update (run tests to verify)
- PR creation (create PR for updates)

Requirements:
- go, npm, python installed
- gh CLI for PR creation
- Repository with proper CI setup

Usage:
    python scripts/automation/update_dependencies.py --check
    python scripts/automation/update_dependencies.py --update --test
    python scripts/automation/update_dependencies.py --update --pr
    python scripts/automation/update_dependencies.py --ecosystem go --update

Author: Automation System
Created: 2026-01-31
"""

import argparse
import json
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path


@dataclass
class Dependency:
    """Represents a dependency with current and available versions."""

    name: str
    current: str
    available: str
    ecosystem: str


class DependencyUpdater:
    """Manage dependency updates across ecosystems."""

    def __init__(self, dry_run: bool = False, verbose: bool = False):
        """Initialize dependency updater.

        Args:
            dry_run: If True, print actions without executing
            verbose: If True, print detailed output
        """
        self.dry_run = dry_run
        self.verbose = verbose
        self.root = Path.cwd()

    def run_command(self, cmd: list[str], cwd: Path | None = None) -> str:
        """Run command and return output.

        Args:
            cmd: Command to run
            cwd: Working directory

        Returns:
            Command output

        Raises:
            subprocess.CalledProcessError: If command fails
        """
        if self.dry_run:
            print(f"[DRY-RUN] Would run: {' '.join(cmd)}")
            return ""

        result = subprocess.run(
            cmd,
            cwd=cwd or self.root,
            capture_output=True,
            text=True,
            check=True,
        )
        return result.stdout.strip()

    def check_go_dependencies(self) -> list[Dependency]:
        """Check Go module dependencies for updates.

        Returns:
            List of outdated dependencies
        """
        if self.verbose:
            print("\n🔍 Checking Go dependencies...")

        try:
            # Get outdated modules
            output = self.run_command(["go", "list", "-u", "-m", "-json", "all"])

            if not output or self.dry_run:
                return []

            # Parse JSON output (one JSON object per line)
            outdated = []
            for line in output.split("\n"):
                if not line.strip():
                    continue

                try:
                    mod = json.loads(line)
                    if "Update" in mod:
                        outdated.append(
                            Dependency(
                                name=mod["Path"],
                                current=mod["Version"],
                                available=mod["Update"]["Version"],
                                ecosystem="go",
                            )
                        )
                except json.JSONDecodeError:
                    continue

            if self.verbose:
                print(f"   Found {len(outdated)} outdated Go modules")

            return outdated

        except subprocess.CalledProcessError as e:
            print(f"⚠️  Could not check Go dependencies: {e}")
            return []

    def check_npm_dependencies(self) -> list[Dependency]:
        """Check npm dependencies for updates.

        Returns:
            List of outdated dependencies
        """
        if self.verbose:
            print("\n🔍 Checking npm dependencies...")

        frontend_dir = self.root / "frontend"
        if not frontend_dir.exists():
            if self.verbose:
                print("   No frontend directory found")
            return []

        try:
            # Get outdated packages
            output = self.run_command(
                ["npm", "outdated", "--json"], cwd=frontend_dir
            )

            if not output or self.dry_run:
                return []

            # Parse JSON output
            outdated_data = json.loads(output)
            outdated = []

            for name, info in outdated_data.items():
                outdated.append(
                    Dependency(
                        name=name,
                        current=info["current"],
                        available=info["latest"],
                        ecosystem="npm",
                    )
                )

            if self.verbose:
                print(f"   Found {len(outdated)} outdated npm packages")

            return outdated

        except subprocess.CalledProcessError:
            # npm outdated exits with code 1 if there are outdated packages
            # Try again without check=True
            try:
                result = subprocess.run(
                    ["npm", "outdated", "--json"],
                    cwd=frontend_dir,
                    capture_output=True,
                    text=True,
                    check=False,
                )
                output = result.stdout.strip()
                if output:
                    outdated_data = json.loads(output)
                    outdated = []
                    for name, info in outdated_data.items():
                        outdated.append(
                            Dependency(
                                name=name,
                                current=info["current"],
                                available=info["latest"],
                                ecosystem="npm",
                            )
                        )
                    if self.verbose:
                        print(f"   Found {len(outdated)} outdated npm packages")
                    return outdated
            except (subprocess.CalledProcessError, json.JSONDecodeError):
                pass

            if self.verbose:
                print("   No outdated npm packages")
            return []

    def check_python_dependencies(self) -> list[Dependency]:
        """Check Python dependencies for updates.

        Returns:
            List of outdated dependencies
        """
        if self.verbose:
            print("\n🔍 Checking Python dependencies...")

        try:
            # Get outdated packages
            output = self.run_command(["pip", "list", "--outdated", "--format=json"])

            if not output or self.dry_run:
                return []

            # Parse JSON output
            outdated_data = json.loads(output)
            outdated = []

            for pkg in outdated_data:
                outdated.append(
                    Dependency(
                        name=pkg["name"],
                        current=pkg["version"],
                        available=pkg["latest_version"],
                        ecosystem="python",
                    )
                )

            if self.verbose:
                print(f"   Found {len(outdated)} outdated Python packages")

            return outdated

        except subprocess.CalledProcessError as e:
            print(f"⚠️  Could not check Python dependencies: {e}")
            return []

    def check_all_dependencies(
        self, ecosystem: str | None = None
    ) -> dict[str, list[Dependency]]:
        """Check all dependencies or specific ecosystem.

        Args:
            ecosystem: Specific ecosystem to check (go, npm, python) or None for all

        Returns:
            Dict mapping ecosystem to list of outdated dependencies
        """
        results = {}

        if ecosystem is None or ecosystem == "go":
            results["go"] = self.check_go_dependencies()

        if ecosystem is None or ecosystem == "npm":
            results["npm"] = self.check_npm_dependencies()

        if ecosystem is None or ecosystem == "python":
            results["python"] = self.check_python_dependencies()

        return results

    def update_go_dependencies(self, deps: list[Dependency]) -> bool:
        """Update Go dependencies.

        Args:
            deps: List of dependencies to update

        Returns:
            True if successful
        """
        if not deps:
            return True

        print(f"\n📦 Updating {len(deps)} Go dependencies...")

        if self.dry_run:
            for dep in deps:
                print(f"   [DRY-RUN] Would update {dep.name}: {dep.current} → {dep.available}")
            return True

        try:
            # Update all outdated modules
            for dep in deps:
                print(f"   Updating {dep.name}: {dep.current} → {dep.available}")
                self.run_command(["go", "get", f"{dep.name}@{dep.available}"])

            # Tidy up
            self.run_command(["go", "mod", "tidy"])
            print("   ✓ Go dependencies updated")
            return True

        except subprocess.CalledProcessError as e:
            print(f"   ✗ Failed to update Go dependencies: {e}")
            return False

    def update_npm_dependencies(self, deps: list[Dependency]) -> bool:
        """Update npm dependencies.

        Args:
            deps: List of dependencies to update

        Returns:
            True if successful
        """
        if not deps:
            return True

        print(f"\n📦 Updating {len(deps)} npm dependencies...")

        if self.dry_run:
            for dep in deps:
                print(f"   [DRY-RUN] Would update {dep.name}: {dep.current} → {dep.available}")
            return True

        frontend_dir = self.root / "frontend"
        try:
            # Update packages
            for dep in deps:
                print(f"   Updating {dep.name}: {dep.current} → {dep.available}")
                self.run_command(
                    ["npm", "update", dep.name], cwd=frontend_dir
                )

            print("   ✓ npm dependencies updated")
            return True

        except subprocess.CalledProcessError as e:
            print(f"   ✗ Failed to update npm dependencies: {e}")
            return False

    def update_python_dependencies(self, deps: list[Dependency]) -> bool:
        """Update Python dependencies.

        Args:
            deps: List of dependencies to update

        Returns:
            True if successful
        """
        if not deps:
            return True

        print(f"\n📦 Updating {len(deps)} Python dependencies...")

        if self.dry_run:
            for dep in deps:
                print(f"   [DRY-RUN] Would update {dep.name}: {dep.current} → {dep.available}")
            return True

        try:
            # Update packages
            for dep in deps:
                print(f"   Updating {dep.name}: {dep.current} → {dep.available}")
                self.run_command(
                    ["pip", "install", "--upgrade", dep.name]
                )

            # Update requirements.txt if it exists
            req_file = self.root / "scripts" / "requirements.txt"
            if req_file.exists():
                self.run_command(["pip", "freeze", ">", str(req_file)])

            print("   ✓ Python dependencies updated")
            return True

        except subprocess.CalledProcessError as e:
            print(f"   ✗ Failed to update Python dependencies: {e}")
            return False

    def run_tests(self) -> bool:
        """Run tests after dependency updates.

        Returns:
            True if tests pass
        """
        print("\n🧪 Running tests...")

        if self.dry_run:
            print("   [DRY-RUN] Would run tests")
            return True

        try:
            # Run Go tests
            print("   Running Go tests...")
            self.run_command(["go", "test", "./..."])

            # Run npm tests if frontend exists
            frontend_dir = self.root / "frontend"
            if frontend_dir.exists():
                print("   Running npm tests...")
                self.run_command(["npm", "test"], cwd=frontend_dir)

            # Run Python tests if they exist
            test_dir = self.root / "tests"
            if test_dir.exists():
                print("   Running Python tests...")
                self.run_command(["pytest"])

            print("   ✓ All tests passed")
            return True

        except subprocess.CalledProcessError as e:
            print(f"   ✗ Tests failed: {e}")
            return False

    def create_pr(self, deps: dict[str, list[Dependency]]) -> bool:
        """Create PR for dependency updates.

        Args:
            deps: Dict of dependencies by ecosystem

        Returns:
            True if PR created successfully
        """
        print("\n📝 Creating PR for dependency updates...")

        if self.dry_run:
            print("   [DRY-RUN] Would create PR")
            return True

        # Count total updates
        total = sum(len(d) for d in deps.values())
        if total == 0:
            print("   No dependencies to update")
            return True

        # Create branch
        branch = "chore/dependency-updates"
        try:
            self.run_command(["git", "checkout", "-b", branch])
        except subprocess.CalledProcessError:
            # Branch might exist, try to use it
            try:
                self.run_command(["git", "checkout", branch])
            except subprocess.CalledProcessError as e:
                print(f"   ✗ Failed to create branch: {e}")
                return False

        # Add changes
        try:
            self.run_command(["git", "add", "."])

            # Create commit message
            msg_lines = ["chore(deps): update dependencies", ""]
            for ecosystem, dep_list in deps.items():
                if dep_list:
                    msg_lines.append(f"{ecosystem.upper()}:")
                    for dep in dep_list:
                        msg_lines.append(f"- {dep.name}: {dep.current} → {dep.available}")
                    msg_lines.append("")

            msg_lines.append("Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>")

            self.run_command(["git", "commit", "-m", "\n".join(msg_lines)])

            # Push and create PR
            self.run_command(["git", "push", "-u", "origin", branch])

            pr_body = "## Dependency Updates\n\n"
            pr_body += f"Updated {total} dependencies across {len(deps)} ecosystems.\n\n"

            for ecosystem, dep_list in deps.items():
                if dep_list:
                    pr_body += f"### {ecosystem.upper()}\n\n"
                    for dep in dep_list:
                        pr_body += f"- **{dep.name}**: {dep.current} → {dep.available}\n"
                    pr_body += "\n"

            pr_body += "\n🤖 Generated with [Claude Code](https://claude.com/claude-code)"

            self.run_command([
                "gh", "pr", "create",
                "--title", "chore(deps): update dependencies",
                "--body", pr_body,
                "--label", "dependencies",
            ])

            print("   ✓ PR created successfully")
            return True

        except subprocess.CalledProcessError as e:
            print(f"   ✗ Failed to create PR: {e}")
            return False


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Manage dependency updates")
    parser.add_argument(
        "--check",
        action="store_true",
        help="Check for outdated dependencies",
    )
    parser.add_argument(
        "--update",
        action="store_true",
        help="Update outdated dependencies",
    )
    parser.add_argument(
        "--test",
        action="store_true",
        help="Run tests after updating (requires --update)",
    )
    parser.add_argument(
        "--pr",
        action="store_true",
        help="Create PR for updates (requires --update)",
    )
    parser.add_argument(
        "--ecosystem",
        choices=["go", "npm", "python"],
        help="Check/update specific ecosystem only",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Print actions without executing",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Verbose output",
    )

    args = parser.parse_args()

    # Validate arguments
    if not args.check and not args.update:
        print("❌ Error: Specify --check or --update")
        parser.print_help()
        sys.exit(1)

    if args.test and not args.update:
        print("❌ Error: --test requires --update")
        sys.exit(1)

    if args.pr and not args.update:
        print("❌ Error: --pr requires --update")
        sys.exit(1)

    # Create updater
    updater = DependencyUpdater(dry_run=args.dry_run, verbose=args.verbose)

    # Check dependencies
    print(f"\n{'='*70}")
    print(f"Dependency Update - {args.ecosystem or 'all ecosystems'}")
    print(f"{'='*70}\n")

    if args.dry_run:
        print("🔍 DRY-RUN MODE - No changes will be made\n")

    deps = updater.check_all_dependencies(args.ecosystem)

    # Print summary
    total = sum(len(d) for d in deps.values())
    print(f"\n📊 Found {total} outdated dependencies\n")

    for ecosystem, dep_list in deps.items():
        if dep_list:
            print(f"{ecosystem.upper()}:")
            for dep in dep_list:
                print(f"   {dep.name}: {dep.current} → {dep.available}")
            print()

    if args.check:
        # Just checking, we're done
        sys.exit(0)

    # Update dependencies
    if total == 0:
        print("✅ All dependencies are up to date!")
        sys.exit(0)

    success = True

    for ecosystem, dep_list in deps.items():
        if ecosystem == "go":
            success &= updater.update_go_dependencies(dep_list)
        elif ecosystem == "npm":
            success &= updater.update_npm_dependencies(dep_list)
        elif ecosystem == "python":
            success &= updater.update_python_dependencies(dep_list)

    if not success:
        print("\n❌ Some updates failed")
        sys.exit(1)

    # Run tests if requested
    if args.test:
        if not updater.run_tests():
            print("\n❌ Tests failed after updates")
            sys.exit(1)

    # Create PR if requested
    if args.pr:
        if not updater.create_pr(deps):
            print("\n❌ Failed to create PR")
            sys.exit(1)

    print("\n✅ Dependency updates complete!\n")


if __name__ == "__main__":
    main()
