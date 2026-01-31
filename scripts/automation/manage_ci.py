#!/usr/bin/env python3
"""CI/CD workflow management script.

Manages GitHub Actions workflows:
- List workflows
- Validate workflow syntax
- Trigger workflow runs
- Monitor workflow status
- Download workflow logs

Requires:
- gh CLI installed and authenticated
- actionlint (optional, for validation)

Usage:
    python scripts/automation/manage_ci.py --list
    python scripts/automation/manage_ci.py --validate
    python scripts/automation/manage_ci.py --trigger ci.yml
    python scripts/automation/manage_ci.py --status
    python scripts/automation/manage_ci.py --logs 123456

Author: Automation System
Created: 2026-01-31
"""

import argparse
import subprocess
import sys
from pathlib import Path


class CIManager:
    """Manage CI/CD workflows."""

    def __init__(self, dry_run: bool = False, verbose: bool = False):
        """Initialize CI manager.

        Args:
            dry_run: If True, print actions without executing
            verbose: If True, print detailed output
        """
        self.dry_run = dry_run
        self.verbose = verbose
        self.root = Path.cwd()
        self.workflows_dir = self.root / ".github" / "workflows"

    def run_command(self, cmd: list[str], cwd: Path | None = None) -> subprocess.CompletedProcess:
        """Run command and return result.

        Args:
            cmd: Command to run
            cwd: Working directory

        Returns:
            CompletedProcess result
        """
        if self.verbose:
            print(f"Running: {' '.join(cmd)}")

        if self.dry_run and cmd[0] == "gh":
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

    def check_gh_cli(self) -> bool:
        """Check if gh CLI is installed.

        Returns:
            True if gh CLI is available
        """
        result = self.run_command(["gh", "--version"])
        return result.returncode == 0

    def list_workflows(self) -> bool:
        """List all GitHub Actions workflows.

        Returns:
            True if successful
        """
        print("\n📋 Listing GitHub Actions workflows...")

        result = self.run_command(["gh", "workflow", "list"])

        if result.returncode != 0:
            print(f"❌ Failed to list workflows: {result.stderr}")
            return False

        print(result.stdout)
        return True

    def list_workflow_files(self) -> list[Path]:
        """List workflow files in .github/workflows/.

        Returns:
            List of workflow file paths
        """
        if not self.workflows_dir.exists():
            return []

        return sorted(self.workflows_dir.glob("*.yml")) + sorted(self.workflows_dir.glob("*.yaml"))

    def validate_workflows(self) -> bool:
        """Validate workflow syntax with actionlint.

        Returns:
            True if all workflows are valid
        """
        print("\n🔍 Validating GitHub Actions workflows...")

        workflow_files = self.list_workflow_files()
        if not workflow_files:
            print("⚠️  No workflow files found")
            return True

        print(f"Found {len(workflow_files)} workflow files")

        # Check if actionlint is available
        result = self.run_command(["actionlint", "-version"])
        if result.returncode != 0:
            print("⚠️  actionlint not found, using basic validation")
            # Basic YAML syntax check
            try:
                import yaml

                for workflow_file in workflow_files:
                    with open(workflow_file) as f:
                        yaml.safe_load(f)
                    print(f"   ✓ {workflow_file.name}")

                print("✅ Basic YAML validation passed")
                return True

            except ImportError:
                print("❌ PyYAML not installed")
                return False

        # Run actionlint
        result = self.run_command(["actionlint"])

        if result.returncode != 0:
            print(f"❌ Workflow validation failed:\n{result.stdout}")
            return False

        print("✅ All workflows are valid")
        return True

    def trigger_workflow(self, workflow: str, branch: str = "main") -> bool:
        """Trigger a workflow run.

        Args:
            workflow: Workflow filename or ID
            branch: Branch to run on

        Returns:
            True if successful
        """
        print(f"\n▶️  Triggering workflow: {workflow} on {branch}...")

        result = self.run_command([
            "gh", "workflow", "run",
            workflow,
            "--ref", branch,
        ])

        if result.returncode != 0:
            print(f"❌ Failed to trigger workflow: {result.stderr}")
            return False

        print(f"✅ Workflow '{workflow}' triggered successfully")
        return True

    def get_workflow_status(self, limit: int = 10) -> bool:
        """Get recent workflow run status.

        Args:
            limit: Number of recent runs to show

        Returns:
            True if successful
        """
        print(f"\n📊 Recent workflow runs (limit: {limit})...")

        result = self.run_command([
            "gh", "run", "list",
            "--limit", str(limit),
        ])

        if result.returncode != 0:
            print(f"❌ Failed to get workflow status: {result.stderr}")
            return False

        print(result.stdout)
        return True

    def watch_workflow(self, run_id: str) -> bool:
        """Watch a workflow run in real-time.

        Args:
            run_id: Workflow run ID

        Returns:
            True if successful
        """
        print(f"\n👁️  Watching workflow run: {run_id}...")

        result = self.run_command(["gh", "run", "watch", run_id])

        if result.returncode != 0:
            print(f"❌ Failed to watch workflow: {result.stderr}")
            return False

        print("✅ Workflow completed")
        return True

    def download_logs(self, run_id: str, output_dir: Path | None = None) -> bool:
        """Download workflow run logs.

        Args:
            run_id: Workflow run ID
            output_dir: Output directory for logs

        Returns:
            True if successful
        """
        print(f"\n📥 Downloading logs for run: {run_id}...")

        if output_dir is None:
            output_dir = self.root / "logs" / f"run-{run_id}"

        output_dir.mkdir(parents=True, exist_ok=True)

        result = self.run_command([
            "gh", "run", "download",
            run_id,
            "--dir", str(output_dir),
        ])

        if result.returncode != 0:
            print(f"❌ Failed to download logs: {result.stderr}")
            return False

        print(f"✅ Logs downloaded to: {output_dir}")
        return True

    def view_log(self, run_id: str, job: str | None = None) -> bool:
        """View workflow run log.

        Args:
            run_id: Workflow run ID
            job: Specific job name (optional)

        Returns:
            True if successful
        """
        print(f"\n📄 Viewing log for run: {run_id}...")

        cmd = ["gh", "run", "view", run_id, "--log"]
        if job:
            cmd.extend(["--job", job])

        result = self.run_command(cmd)

        if result.returncode != 0:
            print(f"❌ Failed to view log: {result.stderr}")
            return False

        print(result.stdout)
        return True

    def cancel_workflow(self, run_id: str) -> bool:
        """Cancel a workflow run.

        Args:
            run_id: Workflow run ID

        Returns:
            True if successful
        """
        print(f"\n⏹️  Canceling workflow run: {run_id}...")

        result = self.run_command(["gh", "run", "cancel", run_id])

        if result.returncode != 0:
            print(f"❌ Failed to cancel workflow: {result.stderr}")
            return False

        print(f"✅ Workflow run {run_id} canceled")
        return True

    def rerun_workflow(self, run_id: str, failed_only: bool = False) -> bool:
        """Rerun a workflow run.

        Args:
            run_id: Workflow run ID
            failed_only: If True, only rerun failed jobs

        Returns:
            True if successful
        """
        print(f"\n🔄 Rerunning workflow run: {run_id}...")

        cmd = ["gh", "run", "rerun", run_id]
        if failed_only:
            cmd.append("--failed")

        result = self.run_command(cmd)

        if result.returncode != 0:
            print(f"❌ Failed to rerun workflow: {result.stderr}")
            return False

        print(f"✅ Workflow run {run_id} restarted")
        return True


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Manage CI/CD workflows")
    parser.add_argument(
        "--list",
        action="store_true",
        help="List all workflows",
    )
    parser.add_argument(
        "--validate",
        action="store_true",
        help="Validate workflow syntax",
    )
    parser.add_argument(
        "--trigger",
        metavar="WORKFLOW",
        help="Trigger workflow run",
    )
    parser.add_argument(
        "--status",
        action="store_true",
        help="Show recent workflow runs",
    )
    parser.add_argument(
        "--watch",
        metavar="RUN_ID",
        help="Watch workflow run",
    )
    parser.add_argument(
        "--logs",
        metavar="RUN_ID",
        help="Download workflow logs",
    )
    parser.add_argument(
        "--view",
        metavar="RUN_ID",
        help="View workflow log",
    )
    parser.add_argument(
        "--cancel",
        metavar="RUN_ID",
        help="Cancel workflow run",
    )
    parser.add_argument(
        "--rerun",
        metavar="RUN_ID",
        help="Rerun workflow",
    )
    parser.add_argument(
        "--branch",
        default="main",
        help="Branch for workflow trigger (default: main)",
    )
    parser.add_argument(
        "--job",
        help="Specific job name (use with --view)",
    )
    parser.add_argument(
        "--output",
        type=Path,
        help="Output directory for logs (use with --logs)",
    )
    parser.add_argument(
        "--limit",
        type=int,
        default=10,
        help="Number of runs to show (use with --status, default: 10)",
    )
    parser.add_argument(
        "--failed-only",
        action="store_true",
        help="Only rerun failed jobs (use with --rerun)",
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
    manager = CIManager(dry_run=args.dry_run, verbose=args.verbose)

    # Check gh CLI
    if not manager.check_gh_cli():
        print("❌ Error: gh CLI not found")
        print("   Install: https://cli.github.com/")
        sys.exit(1)

    # Execute action
    success = True

    if args.list:
        success = manager.list_workflows()
    elif args.validate:
        success = manager.validate_workflows()
    elif args.trigger:
        success = manager.trigger_workflow(args.trigger, args.branch)
    elif args.status:
        success = manager.get_workflow_status(args.limit)
    elif args.watch:
        success = manager.watch_workflow(args.watch)
    elif args.logs:
        success = manager.download_logs(args.logs, args.output)
    elif args.view:
        success = manager.view_log(args.view, args.job)
    elif args.cancel:
        success = manager.cancel_workflow(args.cancel)
    elif args.rerun:
        success = manager.rerun_workflow(args.rerun, args.failed_only)
    else:
        print("❌ Error: Specify an action")
        parser.print_help()
        sys.exit(1)

    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
