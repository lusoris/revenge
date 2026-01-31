#!/usr/bin/env python3
"""Log viewer script.

View and search automation logs:
- View workflow run logs
- Search logs for keywords
- Filter by success/failure
- Download logs
- Tail logs in real-time

Supports:
- GitHub Actions logs (via gh CLI)
- Local automation logs
- Docker container logs
- Application logs

Usage:
    python scripts/automation/view_logs.py --workflow
    python scripts/automation/view_logs.py --search "error"
    python scripts/automation/view_logs.py --failed
    python scripts/automation/view_logs.py --tail
    python scripts/automation/view_logs.py --docker postgres

Author: Automation System
Created: 2026-01-31
"""

import argparse
import re
import subprocess
import sys
from datetime import datetime
from pathlib import Path


class LogViewer:
    """View and search logs."""

    def __init__(self, verbose: bool = False):
        """Initialize log viewer.

        Args:
            verbose: If True, print detailed output
        """
        self.verbose = verbose
        self.root = Path.cwd()
        self.logs_dir = self.root / "logs"

    def run_command(self, cmd: list[str], cwd: Path | None = None) -> subprocess.CompletedProcess:
        """Run command and return result.

        Args:
            cmd: Command to run
            cwd: Working directory

        Returns:
            CompletedProcess result
        """
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

    def list_workflow_runs(self, limit: int = 20, status: str | None = None) -> bool:
        """List recent workflow runs.

        Args:
            limit: Number of runs to show
            status: Filter by status (success, failure, in_progress)

        Returns:
            True if successful
        """
        print(f"\n📋 Recent workflow runs (limit: {limit})...")

        cmd = ["gh", "run", "list", "--limit", str(limit)]

        if status:
            # Map status to gh CLI format
            status_map = {
                "success": "completed",
                "failure": "failure",
                "in_progress": "in_progress",
            }
            if status in status_map:
                cmd.extend(["--status", status_map[status]])

        result = self.run_command(cmd)

        if result.returncode != 0:
            print(f"❌ Failed to list workflow runs: {result.stderr}")
            return False

        print(result.stdout)
        return True

    def view_workflow_log(self, run_id: str, job: str | None = None) -> bool:
        """View workflow run log.

        Args:
            run_id: Workflow run ID
            job: Specific job name (optional)

        Returns:
            True if successful
        """
        print(f"\n📄 Viewing workflow run log: {run_id}...")

        cmd = ["gh", "run", "view", run_id, "--log"]
        if job:
            cmd.extend(["--job", job])

        result = self.run_command(cmd)

        if result.returncode != 0:
            print(f"❌ Failed to view log: {result.stderr}")
            return False

        print(result.stdout)
        return True

    def search_workflow_logs(
        self, run_id: str, pattern: str, case_sensitive: bool = False,
    ) -> bool:
        """Search workflow logs for pattern.

        Args:
            run_id: Workflow run ID
            pattern: Search pattern (regex)
            case_sensitive: If True, case-sensitive search

        Returns:
            True if successful
        """
        print(f"\n🔍 Searching workflow run {run_id} for: {pattern}...")

        # Get logs
        result = self.run_command(["gh", "run", "view", run_id, "--log"])

        if result.returncode != 0:
            print(f"❌ Failed to get logs: {result.stderr}")
            return False

        # Search logs
        flags = 0 if case_sensitive else re.IGNORECASE
        matches = []

        for line_num, line in enumerate(result.stdout.split("\n"), 1):
            if re.search(pattern, line, flags):
                matches.append((line_num, line))

        if matches:
            print(f"\nFound {len(matches)} matches:\n")
            for line_num, line in matches:
                print(f"{line_num:6}: {line}")
        else:
            print("No matches found")

        return True

    def download_workflow_logs(self, run_id: str, output_dir: Path | None = None) -> bool:
        """Download workflow logs.

        Args:
            run_id: Workflow run ID
            output_dir: Output directory (default: logs/run-{id})

        Returns:
            True if successful
        """
        print(f"\n📥 Downloading workflow logs: {run_id}...")

        if output_dir is None:
            output_dir = self.logs_dir / f"run-{run_id}"

        output_dir.mkdir(parents=True, exist_ok=True)

        result = self.run_command([
            "gh",
            "run",
            "download",
            run_id,
            "--dir",
            str(output_dir),
        ])

        if result.returncode != 0:
            print(f"❌ Failed to download logs: {result.stderr}")
            return False

        print(f"✅ Logs downloaded to: {output_dir}")
        return True

    def view_docker_logs(
        self, container: str, lines: int = 100, follow: bool = False,
    ) -> bool:
        """View Docker container logs.

        Args:
            container: Container name
            lines: Number of lines to show
            follow: If True, follow logs

        Returns:
            True if successful
        """
        print(f"\n📄 Viewing Docker logs for: {container}...")

        cmd = ["docker", "compose", "logs", "--tail", str(lines)]
        if follow:
            cmd.append("--follow")
        cmd.append(container)

        if follow:
            # For follow mode, use subprocess.run without capture
            try:
                subprocess.run(cmd, check=False)
                return True
            except KeyboardInterrupt:
                print("\n⚠️  Log streaming stopped")
                return True
        else:
            result = self.run_command(cmd)

            if result.returncode != 0:
                print(f"❌ Failed to view logs: {result.stderr}")
                return False

            print(result.stdout)
            return True

    def view_local_logs(
        self, log_file: str, lines: int | None = None, follow: bool = False,
    ) -> bool:
        """View local log files.

        Args:
            log_file: Log file path
            lines: Number of lines to show (None = all)
            follow: If True, follow logs

        Returns:
            True if successful
        """
        log_path = Path(log_file)
        if not log_path.exists():
            # Try in logs directory
            log_path = self.logs_dir / log_file
            if not log_path.exists():
                print(f"❌ Log file not found: {log_file}")
                return False

        print(f"\n📄 Viewing log file: {log_path}...")

        if follow:
            # Use tail -f for following
            cmd = ["tail", "-f", str(log_path)]
            try:
                subprocess.run(cmd, check=False)
                return True
            except KeyboardInterrupt:
                print("\n⚠️  Log streaming stopped")
                return True
        else:
            # Read file
            with open(log_path) as f:
                content = f.read()

            if lines:
                # Show last N lines
                content_lines = content.split("\n")
                content = "\n".join(content_lines[-lines:])

            print(content)
            return True

    def search_local_logs(
        self, pattern: str, directory: Path | None = None, case_sensitive: bool = False,
    ) -> bool:
        """Search local log files.

        Args:
            pattern: Search pattern
            directory: Directory to search (default: logs/)
            case_sensitive: If True, case-sensitive search

        Returns:
            True if successful
        """
        if directory is None:
            directory = self.logs_dir

        if not directory.exists():
            print(f"❌ Log directory not found: {directory}")
            return False

        print(f"\n🔍 Searching logs in {directory} for: {pattern}...")

        # Find all log files
        log_files = list(directory.rglob("*.log")) + list(directory.rglob("*.txt"))

        if not log_files:
            print("No log files found")
            return True

        flags = 0 if case_sensitive else re.IGNORECASE
        matches_found = False

        for log_file in log_files:
            try:
                with open(log_file) as f:
                    for line_num, line in enumerate(f, 1):
                        if re.search(pattern, line, flags):
                            if not matches_found:
                                matches_found = True
                            print(f"{log_file.name}:{line_num}: {line.rstrip()}")

            except (OSError, UnicodeDecodeError):
                continue

        if not matches_found:
            print("No matches found")

        return True

    def list_local_logs(self) -> bool:
        """List all local log files.

        Returns:
            True if successful
        """
        print(f"\n📋 Local log files in {self.logs_dir}...")

        if not self.logs_dir.exists():
            print("No logs directory found")
            return True

        log_files = list(self.logs_dir.rglob("*.log")) + list(
            self.logs_dir.rglob("*.txt"),
        )

        if not log_files:
            print("No log files found")
            return True

        for log_file in sorted(log_files):
            size = log_file.stat().st_size
            mtime = datetime.fromtimestamp(log_file.stat().st_mtime)

            # Format size
            if size < 1024:
                size_str = f"{size}B"
            elif size < 1024 * 1024:
                size_str = f"{size / 1024:.1f}KB"
            else:
                size_str = f"{size / (1024 * 1024):.1f}MB"

            print(
                f"{log_file.relative_to(self.logs_dir)}  "
                f"{size_str:>10}  {mtime.strftime('%Y-%m-%d %H:%M')}",
            )

        return True


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="View and search logs")
    parser.add_argument(
        "--workflow",
        action="store_true",
        help="List workflow runs",
    )
    parser.add_argument(
        "--view",
        metavar="RUN_ID",
        help="View workflow run log",
    )
    parser.add_argument(
        "--search",
        metavar="PATTERN",
        help="Search logs for pattern",
    )
    parser.add_argument(
        "--docker",
        metavar="CONTAINER",
        help="View Docker container logs",
    )
    parser.add_argument(
        "--local",
        metavar="FILE",
        help="View local log file",
    )
    parser.add_argument(
        "--list",
        action="store_true",
        help="List local log files",
    )
    parser.add_argument(
        "--download",
        metavar="RUN_ID",
        help="Download workflow logs",
    )
    parser.add_argument(
        "--run-id",
        help="Workflow run ID (use with --search)",
    )
    parser.add_argument(
        "--job",
        help="Specific job name (use with --view)",
    )
    parser.add_argument(
        "--status",
        choices=["success", "failure", "in_progress"],
        help="Filter by status (use with --workflow)",
    )
    parser.add_argument(
        "--limit",
        type=int,
        default=20,
        help="Number of runs to show (default: 20)",
    )
    parser.add_argument(
        "--lines",
        type=int,
        help="Number of lines to show",
    )
    parser.add_argument(
        "--follow",
        "-f",
        action="store_true",
        help="Follow logs in real-time",
    )
    parser.add_argument(
        "--case-sensitive",
        action="store_true",
        help="Case-sensitive search",
    )
    parser.add_argument(
        "--output",
        type=Path,
        help="Output directory (use with --download)",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Verbose output",
    )

    args = parser.parse_args()

    # Create viewer
    viewer = LogViewer(verbose=args.verbose)

    # Execute action
    success = True

    if args.workflow:
        success = viewer.list_workflow_runs(args.limit, args.status)
    elif args.view:
        success = viewer.view_workflow_log(args.view, args.job)
    elif args.search and args.run_id:
        success = viewer.search_workflow_logs(args.run_id, args.search, args.case_sensitive)
    elif args.search:
        success = viewer.search_local_logs(args.search, case_sensitive=args.case_sensitive)
    elif args.docker:
        success = viewer.view_docker_logs(
            args.docker, args.lines or 100, args.follow,
        )
    elif args.local:
        success = viewer.view_local_logs(args.local, args.lines, args.follow)
    elif args.download:
        success = viewer.download_workflow_logs(args.download, args.output)
    elif args.list:
        success = viewer.list_local_logs()
    else:
        print("❌ Error: Specify an action")
        parser.print_help()
        sys.exit(1)

    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
