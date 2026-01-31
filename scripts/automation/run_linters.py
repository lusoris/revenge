#!/usr/bin/env python3
"""Linter runner script.

Runs all linters across the codebase:
- golangci-lint (Go code)
- ruff (Python code)
- markdownlint (Markdown docs)
- prettier (TypeScript/JSON/YAML)

Supports:
- Parallel execution for speed
- Auto-fix mode
- Individual linter selection
- Summary reporting

Requirements:
- golangci-lint installed
- ruff installed (pip install ruff)
- markdownlint-cli installed (npm install -g markdownlint-cli)
- prettier installed (npm install -g prettier)

Usage:
    python scripts/automation/run_linters.py --all
    python scripts/automation/run_linters.py --go --python
    python scripts/automation/run_linters.py --all --fix
    python scripts/automation/run_linters.py --markdown --docs

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
class LintResult:
    """Result of a linter run."""

    linter: str
    success: bool
    output: str
    error: str
    exit_code: int


class LinterRunner:
    """Run linters across the codebase."""

    def __init__(self, fix: bool = False, verbose: bool = False):
        """Initialize linter runner.

        Args:
            fix: If True, auto-fix issues
            verbose: If True, print detailed output
        """
        self.fix = fix
        self.verbose = verbose
        self.root = Path.cwd()

    def run_command(self, cmd: list[str], cwd: Path | None = None) -> LintResult:
        """Run linter command and return result.

        Args:
            cmd: Command to run
            cwd: Working directory

        Returns:
            LintResult with command output
        """
        linter_name = cmd[0]

        if self.verbose:
            print(f"Running {linter_name}...")

        try:
            result = subprocess.run(
                cmd,
                cwd=cwd or self.root,
                capture_output=True,
                text=True,
                check=False,
            )

            return LintResult(
                linter=linter_name,
                success=result.returncode == 0,
                output=result.stdout,
                error=result.stderr,
                exit_code=result.returncode,
            )

        except FileNotFoundError:
            return LintResult(
                linter=linter_name,
                success=False,
                output="",
                error=f"{linter_name} not found. Please install it.",
                exit_code=127,
            )

    def run_golangci_lint(self) -> LintResult:
        """Run golangci-lint on Go code.

        Returns:
            LintResult
        """
        cmd = ["golangci-lint", "run"]

        if self.fix:
            cmd.append("--fix")

        if self.verbose:
            cmd.append("-v")

        cmd.append("./...")

        return self.run_command(cmd)

    def run_ruff(self) -> LintResult:
        """Run ruff on Python code.

        Returns:
            LintResult
        """
        cmd = ["ruff", "check", "."]

        if self.fix:
            cmd.extend(["--fix", "--unsafe-fixes"])

        return self.run_command(cmd)

    def run_markdownlint(self) -> LintResult:
        """Run markdownlint on Markdown docs.

        Returns:
            LintResult
        """
        cmd = ["markdownlint", "**/*.md"]

        if self.fix:
            cmd.append("--fix")

        # Ignore node_modules and .venv
        cmd.extend(["--ignore", "node_modules", "--ignore", ".venv"])

        return self.run_command(cmd)

    def run_prettier(self) -> LintResult:
        """Run prettier on TypeScript/JSON/YAML files.

        Returns:
            LintResult
        """
        cmd = [
            "prettier",
            "--check",
            "**/*.{ts,tsx,js,jsx,json,yaml,yml}",
        ]

        if self.fix:
            cmd[1] = "--write"

        # Ignore patterns
        cmd.extend([
            "--ignore-path",
            ".gitignore",
        ])

        return self.run_command(cmd)

    def run_all_linters(
        self, linters: list[str] | None = None
    ) -> dict[str, LintResult]:
        """Run all or selected linters in parallel.

        Args:
            linters: List of linters to run (go, python, markdown, prettier)
                    If None, run all linters

        Returns:
            Dict mapping linter name to result
        """
        # Map linter names to methods
        linter_methods = {
            "go": self.run_golangci_lint,
            "python": self.run_ruff,
            "markdown": self.run_markdownlint,
            "prettier": self.run_prettier,
        }

        # Filter linters if specified
        if linters:
            linter_methods = {
                name: method
                for name, method in linter_methods.items()
                if name in linters
            }

        results = {}

        # Run linters in parallel
        with ThreadPoolExecutor(max_workers=4) as executor:
            # Submit all linter tasks
            future_to_linter = {
                executor.submit(method): name
                for name, method in linter_methods.items()
            }

            # Collect results as they complete
            for future in as_completed(future_to_linter):
                linter_name = future_to_linter[future]
                try:
                    result = future.result()
                    results[linter_name] = result

                    # Print progress
                    if result.success:
                        print(f"✓ {linter_name}: passed")
                    else:
                        print(f"✗ {linter_name}: failed")

                except Exception as e:
                    print(f"✗ {linter_name}: exception - {e}")
                    results[linter_name] = LintResult(
                        linter=linter_name,
                        success=False,
                        output="",
                        error=str(e),
                        exit_code=1,
                    )

        return results

    def print_summary(self, results: dict[str, LintResult]):
        """Print summary of linter results.

        Args:
            results: Dict of linter results
        """
        print(f"\n{'='*70}")
        print("Linter Summary")
        print(f"{'='*70}\n")

        passed = sum(1 for r in results.values() if r.success)
        failed = len(results) - passed

        for linter_name, result in sorted(results.items()):
            status = "✓ PASS" if result.success else "✗ FAIL"
            print(f"{status:10} {linter_name}")

            # Print errors if failed
            if not result.success and result.error:
                print(f"           Error: {result.error[:100]}")

            # Print verbose output
            if self.verbose and result.output:
                print(f"           Output: {result.output[:200]}")

        print(f"\n{'='*70}")
        print(f"Total: {len(results)} | Passed: {passed} | Failed: {failed}")

        if failed == 0:
            print("✅ All linters passed!")
        else:
            print(f"❌ {failed} linter(s) failed")

        print(f"{'='*70}\n")

    def print_detailed_output(self, results: dict[str, LintResult]):
        """Print detailed output for failed linters.

        Args:
            results: Dict of linter results
        """
        failed_results = {
            name: result for name, result in results.items() if not result.success
        }

        if not failed_results:
            return

        print("\n" + "=" * 70)
        print("Detailed Output for Failed Linters")
        print("=" * 70 + "\n")

        for linter_name, result in sorted(failed_results.items()):
            print(f"\n{linter_name.upper()}")
            print("-" * 70)

            if result.error:
                print("STDERR:")
                print(result.error)

            if result.output:
                print("\nSTDOUT:")
                print(result.output)

            print("-" * 70)


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Run linters on codebase")
    parser.add_argument(
        "--all",
        action="store_true",
        help="Run all linters",
    )
    parser.add_argument(
        "--go",
        action="store_true",
        help="Run golangci-lint",
    )
    parser.add_argument(
        "--python",
        action="store_true",
        help="Run ruff",
    )
    parser.add_argument(
        "--markdown",
        action="store_true",
        help="Run markdownlint",
    )
    parser.add_argument(
        "--prettier",
        action="store_true",
        help="Run prettier",
    )
    parser.add_argument(
        "--fix",
        action="store_true",
        help="Auto-fix issues where possible",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Verbose output",
    )
    parser.add_argument(
        "--detailed",
        action="store_true",
        help="Show detailed output for failures",
    )

    args = parser.parse_args()

    # Determine which linters to run
    linters = []
    if args.all:
        linters = None  # Run all
    else:
        if args.go:
            linters.append("go")
        if args.python:
            linters.append("python")
        if args.markdown:
            linters.append("markdown")
        if args.prettier:
            linters.append("prettier")

        if not linters:
            print("❌ Error: Specify at least one linter or use --all")
            parser.print_help()
            sys.exit(1)

    # Create runner
    runner = LinterRunner(fix=args.fix, verbose=args.verbose)

    # Run linters
    print(f"\n{'='*70}")
    print(f"Running Linters - {len(linters) if linters else 'all'}")
    print(f"{'='*70}\n")

    if args.fix:
        print("🔧 Auto-fix mode enabled\n")

    results = runner.run_all_linters(linters)

    # Print summary
    runner.print_summary(results)

    # Print detailed output if requested
    if args.detailed:
        runner.print_detailed_output(results)

    # Exit with error if any linter failed
    if any(not r.success for r in results.values()):
        sys.exit(1)


if __name__ == "__main__":
    main()
