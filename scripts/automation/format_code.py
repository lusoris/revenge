#!/usr/bin/env python3
"""Code formatter script.

Formats code across all languages in the codebase:
- Go code (gofmt, goimports)
- Python code (ruff format)
- TypeScript/JavaScript code (prettier)
- JSON/YAML files (prettier)
- Markdown docs (prettier)

Supports:
- Check mode (verify formatting without changes)
- Fix mode (apply formatting)
- Individual language selection
- Parallel execution

Requirements:
- gofmt (included with Go)
- goimports (go install golang.org/x/tools/cmd/goimports@latest)
- ruff (pip install ruff)
- prettier (npm install -g prettier)

Usage:
    python scripts/automation/format_code.py --all
    python scripts/automation/format_code.py --go --python
    python scripts/automation/format_code.py --all --check
    python scripts/automation/format_code.py --frontend --fix

Author: Automation System
Created: 2026-01-31
"""

import argparse
import subprocess
import sys
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import dataclass
from pathlib import Path


@dataclass
class FormatResult:
    """Result of a formatter run."""

    formatter: str
    success: bool
    files_formatted: int
    output: str
    error: str
    exit_code: int


class CodeFormatter:
    """Format code across the codebase."""

    def __init__(self, check: bool = False, verbose: bool = False):
        """Initialize code formatter.

        Args:
            check: If True, check formatting without changes
            verbose: If True, print detailed output
        """
        self.check = check
        self.verbose = verbose
        self.root = Path.cwd()

    def run_command(self, cmd: list[str], cwd: Path | None = None) -> FormatResult:
        """Run formatter command and return result.

        Args:
            cmd: Command to run
            cwd: Working directory

        Returns:
            FormatResult with command output
        """
        formatter_name = cmd[0]

        if self.verbose:
            print(f"Running {formatter_name}...")

        try:
            result = subprocess.run(
                cmd,
                cwd=cwd or self.root,
                capture_output=True,
                text=True,
                check=False,
            )

            # Count formatted files from output
            files_formatted = result.stdout.count("\n") if result.stdout else 0

            return FormatResult(
                formatter=formatter_name,
                success=result.returncode == 0,
                files_formatted=files_formatted,
                output=result.stdout,
                error=result.stderr,
                exit_code=result.returncode,
            )

        except FileNotFoundError:
            return FormatResult(
                formatter=formatter_name,
                success=False,
                files_formatted=0,
                output="",
                error=f"{formatter_name} not found. Please install it.",
                exit_code=127,
            )

    def format_go_code(self) -> FormatResult:
        """Format Go code with gofmt and goimports.

        Returns:
            FormatResult
        """
        # First run gofmt
        gofmt_cmd = ["gofmt"]
        if not self.check:
            gofmt_cmd.append("-w")
        else:
            gofmt_cmd.append("-l")

        gofmt_cmd.append(".")

        gofmt_result = self.run_command(gofmt_cmd)

        # Then run goimports
        goimports_cmd = ["goimports"]
        if not self.check:
            goimports_cmd.append("-w")
        else:
            goimports_cmd.append("-l")

        goimports_cmd.append(".")

        goimports_result = self.run_command(goimports_cmd)

        # Combine results
        return FormatResult(
            formatter="go",
            success=gofmt_result.success and goimports_result.success,
            files_formatted=gofmt_result.files_formatted
            + goimports_result.files_formatted,
            output=f"{gofmt_result.output}\n{goimports_result.output}",
            error=f"{gofmt_result.error}\n{goimports_result.error}",
            exit_code=gofmt_result.exit_code or goimports_result.exit_code,
        )

    def format_python_code(self) -> FormatResult:
        """Format Python code with ruff.

        Returns:
            FormatResult
        """
        cmd = ["ruff", "format"]

        if self.check:
            cmd.append("--check")

        cmd.append(".")

        return self.run_command(cmd)

    def format_frontend_code(self) -> FormatResult:
        """Format frontend code with prettier.

        Returns:
            FormatResult
        """
        cmd = ["prettier"]

        if self.check:
            cmd.append("--check")
        else:
            cmd.append("--write")

        # Format TypeScript, JavaScript, JSON, YAML, Markdown
        cmd.extend(
            [
                "**/*.{ts,tsx,js,jsx,json,yaml,yml,md}",
                "--ignore-path",
                ".gitignore",
            ]
        )

        return self.run_command(cmd)

    def format_all_code(
        self,
        formatters: list[str] | None = None,
    ) -> dict[str, FormatResult]:
        """Format all or selected code.

        Args:
            formatters: List of formatters to run (go, python, frontend)
                       If None, run all formatters

        Returns:
            Dict mapping formatter name to result
        """
        # Map formatter names to methods
        formatter_methods = {
            "go": self.format_go_code,
            "python": self.format_python_code,
            "frontend": self.format_frontend_code,
        }

        # Filter formatters if specified
        if formatters:
            formatter_methods = {
                name: method
                for name, method in formatter_methods.items()
                if name in formatters
            }

        results = {}

        # Run formatters in parallel
        with ThreadPoolExecutor(max_workers=3) as executor:
            # Submit all formatter tasks
            future_to_formatter = {
                executor.submit(method): name
                for name, method in formatter_methods.items()
            }

            # Collect results as they complete
            for future in as_completed(future_to_formatter):
                formatter_name = future_to_formatter[future]
                try:
                    result = future.result()
                    results[formatter_name] = result

                    # Print progress
                    if result.success:
                        if self.check:
                            print(f"✓ {formatter_name}: formatting is correct")
                        else:
                            print(
                                f"✓ {formatter_name}: formatted {result.files_formatted} files"
                            )
                    else:
                        print(f"✗ {formatter_name}: formatting issues found")

                except Exception as e:
                    print(f"✗ {formatter_name}: exception - {e}")
                    results[formatter_name] = FormatResult(
                        formatter=formatter_name,
                        success=False,
                        files_formatted=0,
                        output="",
                        error=str(e),
                        exit_code=1,
                    )

        return results

    def print_summary(self, results: dict[str, FormatResult]):
        """Print summary of formatter results.

        Args:
            results: Dict of formatter results
        """
        print(f"\n{'=' * 70}")
        print("Formatter Summary")
        print(f"{'=' * 70}\n")

        total_files = sum(r.files_formatted for r in results.values())
        passed = sum(1 for r in results.values() if r.success)
        failed = len(results) - passed

        for formatter_name, result in sorted(results.items()):
            status = "✓ PASS" if result.success else "✗ FAIL"
            print(f"{status:10} {formatter_name:10} Files: {result.files_formatted}")

            # Print errors if failed
            if not result.success and result.error:
                print(f"           Error: {result.error[:100]}")

            # Print verbose output
            if self.verbose and result.output:
                print(f"           Output: {result.output[:200]}")

        print(f"\n{'=' * 70}")
        print(f"Total Formatters: {len(results)} | Passed: {passed} | Failed: {failed}")

        if not self.check:
            print(f"Total Files Formatted: {total_files}")

        if failed == 0:
            if self.check:
                print("✅ All code is properly formatted!")
            else:
                print("✅ All code formatted successfully!")
        elif self.check:
            print(f"❌ {failed} formatter(s) found formatting issues")
        else:
            print(f"❌ {failed} formatter(s) failed")

        print(f"{'=' * 70}\n")


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Format code across codebase")
    parser.add_argument(
        "--all",
        action="store_true",
        help="Format all code",
    )
    parser.add_argument(
        "--go",
        action="store_true",
        help="Format Go code",
    )
    parser.add_argument(
        "--python",
        action="store_true",
        help="Format Python code",
    )
    parser.add_argument(
        "--frontend",
        action="store_true",
        help="Format frontend code",
    )
    parser.add_argument(
        "--check",
        action="store_true",
        help="Check formatting without making changes",
    )
    parser.add_argument(
        "--fix",
        action="store_true",
        help="Apply formatting (default behavior)",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Verbose output",
    )

    args = parser.parse_args()

    # Determine which formatters to run
    formatters = []
    if args.all:
        formatters = None  # Run all
    else:
        if args.go:
            formatters.append("go")
        if args.python:
            formatters.append("python")
        if args.frontend:
            formatters.append("frontend")

        if not formatters:
            print("❌ Error: Specify at least one formatter or use --all")
            parser.print_help()
            sys.exit(1)

    # Create formatter
    formatter = CodeFormatter(check=args.check, verbose=args.verbose)

    # Run formatters
    print(f"\n{'=' * 70}")
    print(f"Code Formatting - {len(formatters) if formatters else 'all languages'}")
    print(f"{'=' * 70}\n")

    if args.check:
        print("🔍 Check mode - no changes will be made\n")
    else:
        print("🔧 Format mode - applying formatting\n")

    results = formatter.format_all_code(formatters)

    # Print summary
    formatter.print_summary(results)

    # Exit with error if any formatter failed
    if any(not r.success for r in results.values()):
        sys.exit(1)


if __name__ == "__main__":
    main()
