#!/usr/bin/env python3
"""Coder workspace management script.

Manages Coder workspaces and templates:
- Update template from SOURCE_OF_TRUTH
- Create/start/stop/delete workspaces
- SSH connection helper
- Template validation

Requires:
- coder CLI installed and authenticated
- Terraform installed (for template validation)

Usage:
    python scripts/automation/manage_coder.py --update-template
    python scripts/automation/manage_coder.py --create my-workspace
    python scripts/automation/manage_coder.py --start my-workspace
    python scripts/automation/manage_coder.py --stop my-workspace
    python scripts/automation/manage_coder.py --ssh my-workspace
    python scripts/automation/manage_coder.py --list

Author: Automation System
Created: 2026-01-31
"""

import argparse
import subprocess
import sys
from pathlib import Path


class CoderManager:
    """Manage Coder workspaces and templates."""

    def __init__(self, dry_run: bool = False, verbose: bool = False):
        """Initialize Coder manager.

        Args:
            dry_run: If True, print actions without executing
            verbose: If True, print detailed output
        """
        self.dry_run = dry_run
        self.verbose = verbose
        self.root = Path.cwd()
        self.template_dir = self.root / ".coder"

    def run_command(
        self, cmd: list[str], cwd: Path | None = None
    ) -> subprocess.CompletedProcess:
        """Run command and return result.

        Args:
            cmd: Command to run
            cwd: Working directory

        Returns:
            CompletedProcess result
        """
        if self.verbose:
            print(f"Running: {' '.join(cmd)}")

        if self.dry_run:
            print(f"[DRY-RUN] Would run: {' '.join(cmd)}")
            return subprocess.CompletedProcess(
                args=cmd,
                returncode=0,
                stdout="",
                stderr="",
            )

        try:
            result = subprocess.run(
                cmd,
                cwd=cwd or self.root,
                capture_output=True,
                text=True,
                check=False,
            )
            return result

        except FileNotFoundError as e:
            return subprocess.CompletedProcess(
                args=cmd,
                returncode=127,
                stdout="",
                stderr=f"{cmd[0]} not found: {e}",
            )

    def check_coder_cli(self) -> bool:
        """Check if coder CLI is installed.

        Returns:
            True if coder CLI is available
        """
        result = self.run_command(["coder", "version"])
        return result.returncode == 0

    def list_workspaces(self) -> bool:
        """List all Coder workspaces.

        Returns:
            True if successful
        """
        print("\n📋 Listing Coder workspaces...")

        result = self.run_command(["coder", "ls"])

        if result.returncode != 0:
            print(f"❌ Failed to list workspaces: {result.stderr}")
            return False

        print(result.stdout)
        return True

    def create_workspace(self, name: str, template: str = "revenge") -> bool:
        """Create a new Coder workspace.

        Args:
            name: Workspace name
            template: Template name

        Returns:
            True if successful
        """
        print(f"\n🏗️  Creating workspace: {name}")

        if self.dry_run:
            print(
                f"[DRY-RUN] Would create workspace '{name}' from template '{template}'"
            )
            return True

        result = self.run_command(
            [
                "coder",
                "create",
                "--template",
                template,
                name,
            ]
        )

        if result.returncode != 0:
            print(f"❌ Failed to create workspace: {result.stderr}")
            return False

        print(f"✅ Workspace '{name}' created successfully")
        return True

    def start_workspace(self, name: str) -> bool:
        """Start a Coder workspace.

        Args:
            name: Workspace name

        Returns:
            True if successful
        """
        print(f"\n▶️  Starting workspace: {name}")

        result = self.run_command(["coder", "start", name])

        if result.returncode != 0:
            print(f"❌ Failed to start workspace: {result.stderr}")
            return False

        print(f"✅ Workspace '{name}' started successfully")
        return True

    def stop_workspace(self, name: str) -> bool:
        """Stop a Coder workspace.

        Args:
            name: Workspace name

        Returns:
            True if successful
        """
        print(f"\n⏹️  Stopping workspace: {name}")

        result = self.run_command(["coder", "stop", name])

        if result.returncode != 0:
            print(f"❌ Failed to stop workspace: {result.stderr}")
            return False

        print(f"✅ Workspace '{name}' stopped successfully")
        return True

    def delete_workspace(self, name: str, force: bool = False) -> bool:
        """Delete a Coder workspace.

        Args:
            name: Workspace name
            force: If True, skip confirmation

        Returns:
            True if successful
        """
        print(f"\n🗑️  Deleting workspace: {name}")

        cmd = ["coder", "delete"]
        if force or self.dry_run:
            cmd.append("--yes")
        cmd.append(name)

        result = self.run_command(cmd)

        if result.returncode != 0:
            print(f"❌ Failed to delete workspace: {result.stderr}")
            return False

        print(f"✅ Workspace '{name}' deleted successfully")
        return True

    def ssh_workspace(self, name: str, command: str | None = None) -> bool:
        """SSH into a Coder workspace.

        Args:
            name: Workspace name
            command: Optional command to run

        Returns:
            True if successful
        """
        print(f"\n🔑 SSH into workspace: {name}")

        cmd = ["coder", "ssh", name]
        if command:
            cmd.extend(["--", command])

        # For SSH, we want interactive mode
        if self.dry_run:
            print(f"[DRY-RUN] Would SSH: {' '.join(cmd)}")
            return True

        try:
            result = subprocess.run(cmd, check=False)
            return result.returncode == 0

        except KeyboardInterrupt:
            print("\n⚠️  SSH session interrupted")
            return True

    def validate_template(self) -> bool:
        """Validate Coder template with terraform validate.

        Returns:
            True if valid
        """
        print("\n🔍 Validating Coder template...")

        template_file = self.template_dir / "template.tf"
        if not template_file.exists():
            print(f"❌ Template file not found: {template_file}")
            return False

        # Run terraform validate
        result = self.run_command(
            ["terraform", "validate"],
            cwd=self.template_dir,
        )

        if result.returncode != 0:
            print(f"❌ Template validation failed:\n{result.stderr}")
            return False

        print("✅ Template is valid")
        return True

    def update_template_from_sot(self) -> bool:
        """Update Coder template from SOURCE_OF_TRUTH.

        Returns:
            True if successful
        """
        print("\n🔄 Updating Coder template from SOURCE_OF_TRUTH...")

        sot_file = self.root / "docs" / "dev" / "design" / "00_SOURCE_OF_TRUTH.md"
        if not sot_file.exists():
            print(f"❌ SOURCE_OF_TRUTH not found: {sot_file}")
            return False

        template_file = self.template_dir / "template.tf"
        if not template_file.exists():
            print(f"❌ Template file not found: {template_file}")
            return False

        print("⚠️  Template update from SOT not yet implemented")
        print("   Manual steps:")
        print("   1. Review SOURCE_OF_TRUTH for version updates")
        print("   2. Update .coder/template.tf accordingly")
        print("   3. Run --validate to check template")

        return True

    def push_template(self) -> bool:
        """Push template to Coder server.

        Returns:
            True if successful
        """
        print("\n📤 Pushing template to Coder...")

        result = self.run_command(
            [
                "coder",
                "templates",
                "push",
                "--directory",
                str(self.template_dir),
                "revenge",
            ]
        )

        if result.returncode != 0:
            print(f"❌ Failed to push template: {result.stderr}")
            return False

        print("✅ Template pushed successfully")
        return True


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Manage Coder workspaces")
    parser.add_argument(
        "--list",
        action="store_true",
        help="List all workspaces",
    )
    parser.add_argument(
        "--create",
        metavar="NAME",
        help="Create workspace with given name",
    )
    parser.add_argument(
        "--start",
        metavar="NAME",
        help="Start workspace",
    )
    parser.add_argument(
        "--stop",
        metavar="NAME",
        help="Stop workspace",
    )
    parser.add_argument(
        "--delete",
        metavar="NAME",
        help="Delete workspace",
    )
    parser.add_argument(
        "--ssh",
        metavar="NAME",
        help="SSH into workspace",
    )
    parser.add_argument(
        "--command",
        help="Command to run via SSH (use with --ssh)",
    )
    parser.add_argument(
        "--validate",
        action="store_true",
        help="Validate Coder template",
    )
    parser.add_argument(
        "--update-template",
        action="store_true",
        help="Update template from SOURCE_OF_TRUTH",
    )
    parser.add_argument(
        "--push-template",
        action="store_true",
        help="Push template to Coder server",
    )
    parser.add_argument(
        "--template",
        default="revenge",
        help="Template name (default: revenge)",
    )
    parser.add_argument(
        "--force",
        action="store_true",
        help="Force operation (skip confirmations)",
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

    # Create manager
    manager = CoderManager(dry_run=args.dry_run, verbose=args.verbose)

    # Check coder CLI
    if not manager.check_coder_cli():
        print("❌ Error: coder CLI not found")
        print("   Install: https://coder.com/docs/install")
        sys.exit(1)

    # Execute action
    success = True

    if args.list:
        success = manager.list_workspaces()
    elif args.create:
        success = manager.create_workspace(args.create, args.template)
    elif args.start:
        success = manager.start_workspace(args.start)
    elif args.stop:
        success = manager.stop_workspace(args.stop)
    elif args.delete:
        success = manager.delete_workspace(args.delete, args.force)
    elif args.ssh:
        success = manager.ssh_workspace(args.ssh, args.command)
    elif args.validate:
        success = manager.validate_template()
    elif args.update_template:
        success = manager.update_template_from_sot()
    elif args.push_template:
        success = manager.push_template()
    else:
        print("❌ Error: Specify an action")
        parser.print_help()
        sys.exit(1)

    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
