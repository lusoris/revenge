#!/usr/bin/env python3
"""Test runner script.

Runs all tests across the codebase:
- Go tests (unit + integration)
- Python tests (pytest)
- Frontend tests (Vitest)

Supports:
- Coverage reporting (80% threshold)
- Watch mode for development
- Individual test suite selection
- Parallel execution

Requirements:
- go installed
- pytest installed (pip install pytest pytest-cov)
- npm installed (for frontend tests)

Usage:
    python scripts/automation/run_tests.py --all
    python scripts/automation/run_tests.py --go --python
    python scripts/automation/run_tests.py --all --coverage
    python scripts/automation/run_tests.py --python --watch

Author: Automation System
Created: 2026-01-31
"""

import argparse
import json
import re
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path


@dataclass
class TestResult:
    """Result of a test run."""

    suite: str
    success: bool
    passed: int
    failed: int
    skipped: int
    coverage: float | None
    output: str
    error: str
    exit_code: int


class TestRunner:
    """Run tests across the codebase."""

    def __init__(
        self,
        coverage: bool = False,
        watch: bool = False,
        verbose: bool = False,
        threshold: int = 80,
    ):
        """Initialize test runner.

        Args:
            coverage: If True, generate coverage reports
            watch: If True, run in watch mode
            verbose: If True, print detailed output
            threshold: Coverage threshold percentage
        """
        self.coverage = coverage
        self.watch = watch
        self.verbose = verbose
        self.threshold = threshold
        self.root = Path.cwd()

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
            # Return a fake result for missing commands
            fake_result = subprocess.CompletedProcess(
                args=cmd,
                returncode=127,
                stdout="",
                stderr=f"{cmd[0]} not found: {e}",
            )
            return fake_result

    def run_go_tests(self) -> TestResult:
        """Run Go tests with optional coverage.

        Returns:
            TestResult
        """
        print("\n🧪 Running Go tests...")

        cmd = ["go", "test"]

        if self.coverage:
            cmd.extend(["-coverprofile=coverage.out", "-covermode=atomic"])

        if self.verbose:
            cmd.append("-v")

        cmd.append("./...")

        result = self.run_command(cmd)

        # Parse output for test counts
        passed = 0
        failed = 0
        skipped = 0

        # Extract test results from output
        # Go test output format: "ok  \tpackage\t0.001s\tcoverage: 80.0% of statements"
        for line in result.stdout.split("\n"):
            if line.startswith("ok"):
                passed += 1
            elif line.startswith("FAIL"):
                failed += 1
            elif "SKIP" in line:
                skipped += 1

        # Extract coverage if available
        coverage_pct = None
        if self.coverage:
            coverage_match = re.search(r"coverage:\s+([\d.]+)%", result.stdout)
            if coverage_match:
                coverage_pct = float(coverage_match.group(1))
            else:
                # Try to get coverage from coverage.out
                coverage_file = self.root / "coverage.out"
                if coverage_file.exists():
                    cov_result = self.run_command(["go", "tool", "cover", "-func=coverage.out"])
                    total_match = re.search(r"total:.*?([\d.]+)%", cov_result.stdout)
                    if total_match:
                        coverage_pct = float(total_match.group(1))

        return TestResult(
            suite="go",
            success=result.returncode == 0,
            passed=passed,
            failed=failed,
            skipped=skipped,
            coverage=coverage_pct,
            output=result.stdout,
            error=result.stderr,
            exit_code=result.returncode,
        )

    def run_python_tests(self) -> TestResult:
        """Run Python tests with pytest.

        Returns:
            TestResult
        """
        print("\n🧪 Running Python tests...")

        cmd = ["pytest"]

        if self.coverage:
            cmd.extend([
                "--cov=scripts",
                "--cov-report=term-missing",
                "--cov-report=json",
            ])

        if self.verbose:
            cmd.append("-v")
        else:
            cmd.append("-q")

        if self.watch:
            # pytest-watch for watch mode
            cmd = ["ptw", "--", *cmd[1:]]

        result = self.run_command(cmd)

        # Parse pytest output for test counts
        passed = 0
        failed = 0
        skipped = 0

        # Extract test results from summary line
        # Format: "272 passed, 1 skipped in 1.38s"
        summary_match = re.search(
            r"(\d+)\s+passed(?:,\s+(\d+)\s+failed)?(?:,\s+(\d+)\s+skipped)?",
            result.stdout,
        )
        if summary_match:
            passed = int(summary_match.group(1))
            failed = int(summary_match.group(2) or 0)
            skipped = int(summary_match.group(3) or 0)

        # Extract coverage if available
        coverage_pct = None
        if self.coverage:
            # Try to read from coverage.json
            coverage_file = self.root / "coverage.json"
            if coverage_file.exists():
                try:
                    with open(coverage_file) as f:
                        cov_data = json.load(f)
                        coverage_pct = cov_data.get("totals", {}).get("percent_covered")
                except (json.JSONDecodeError, KeyError):
                    pass

            # Fallback to parsing output
            if coverage_pct is None:
                cov_match = re.search(r"TOTAL.*?(\d+)%", result.stdout)
                if cov_match:
                    coverage_pct = float(cov_match.group(1))

        return TestResult(
            suite="python",
            success=result.returncode == 0,
            passed=passed,
            failed=failed,
            skipped=skipped,
            coverage=coverage_pct,
            output=result.stdout,
            error=result.stderr,
            exit_code=result.returncode,
        )

    def run_frontend_tests(self) -> TestResult:
        """Run frontend tests with Vitest.

        Returns:
            TestResult
        """
        print("\n🧪 Running frontend tests...")

        frontend_dir = self.root / "frontend"
        if not frontend_dir.exists():
            print("   ⚠️  Frontend directory not found, skipping")
            return TestResult(
                suite="frontend",
                success=True,
                passed=0,
                failed=0,
                skipped=0,
                coverage=None,
                output="Frontend directory not found",
                error="",
                exit_code=0,
            )

        cmd = ["npm", "test"]

        if self.coverage:
            cmd = ["npm", "run", "test:coverage"]

        if self.watch:
            cmd = ["npm", "run", "test:watch"]

        result = self.run_command(cmd, cwd=frontend_dir)

        # Parse Vitest output for test counts
        passed = 0
        failed = 0
        skipped = 0

        # Extract test results
        test_match = re.search(
            r"Test Files\s+(\d+)\s+passed.*?Tests\s+(\d+)\s+passed(?:.*?(\d+)\s+failed)?",
            result.stdout,
            re.DOTALL,
        )
        if test_match:
            passed = int(test_match.group(2))
            failed = int(test_match.group(3) or 0)

        # Extract coverage if available
        coverage_pct = None
        if self.coverage:
            cov_match = re.search(r"All files.*?(\d+\.?\d*)%", result.stdout)
            if cov_match:
                coverage_pct = float(cov_match.group(1))

        return TestResult(
            suite="frontend",
            success=result.returncode == 0,
            passed=passed,
            failed=failed,
            skipped=skipped,
            coverage=coverage_pct,
            output=result.stdout,
            error=result.stderr,
            exit_code=result.returncode,
        )

    def run_all_tests(self, suites: list[str] | None = None) -> dict[str, TestResult]:
        """Run all or selected test suites.

        Args:
            suites: List of test suites to run (go, python, frontend)
                   If None, run all suites

        Returns:
            Dict mapping suite name to result
        """
        # Map suite names to methods
        suite_methods = {
            "go": self.run_go_tests,
            "python": self.run_python_tests,
            "frontend": self.run_frontend_tests,
        }

        # Filter suites if specified
        if suites:
            suite_methods = {
                name: method for name, method in suite_methods.items() if name in suites
            }

        results = {}

        # Run tests sequentially (not parallel due to resource usage)
        for suite_name, method in suite_methods.items():
            try:
                result = method()
                results[suite_name] = result

                # Print progress
                if result.success:
                    print(f"   ✓ {suite_name}: {result.passed} passed")
                else:
                    print(f"   ✗ {suite_name}: {result.failed} failed")

            except Exception as e:
                print(f"   ✗ {suite_name}: exception - {e}")
                results[suite_name] = TestResult(
                    suite=suite_name,
                    success=False,
                    passed=0,
                    failed=0,
                    skipped=0,
                    coverage=None,
                    output="",
                    error=str(e),
                    exit_code=1,
                )

        return results

    def print_summary(self, results: dict[str, TestResult]):
        """Print summary of test results.

        Args:
            results: Dict of test results
        """
        print(f"\n{'='*70}")
        print("Test Summary")
        print(f"{'='*70}\n")

        total_passed = sum(r.passed for r in results.values())
        total_failed = sum(r.failed for r in results.values())
        total_skipped = sum(r.skipped for r in results.values())

        for suite_name, result in sorted(results.items()):
            status = "✓ PASS" if result.success else "✗ FAIL"
            print(f"{status:10} {suite_name:10} P:{result.passed:4} F:{result.failed:4} S:{result.skipped:4}", end="")

            # Print coverage if available
            if result.coverage is not None:
                cov_status = "✓" if result.coverage >= self.threshold else "✗"
                print(f" Coverage:{result.coverage:6.2f}% {cov_status}", end="")

            print()

        print(f"\n{'='*70}")
        print(f"Total: P:{total_passed} F:{total_failed} S:{total_skipped}")

        if self.coverage:
            avg_coverage = sum(
                r.coverage for r in results.values() if r.coverage is not None
            ) / len([r for r in results.values() if r.coverage is not None])
            cov_status = "✓" if avg_coverage >= self.threshold else "✗"
            print(f"Average Coverage: {avg_coverage:.2f}% (threshold: {self.threshold}%) {cov_status}")

        if total_failed == 0:
            print("✅ All tests passed!")
        else:
            print(f"❌ {total_failed} test(s) failed")

        print(f"{'='*70}\n")

    def check_coverage_threshold(self, results: dict[str, TestResult]) -> bool:
        """Check if coverage meets threshold.

        Args:
            results: Dict of test results

        Returns:
            True if all coverage meets threshold
        """
        if not self.coverage:
            return True

        failed_suites = []
        for suite_name, result in results.items():
            if result.coverage is not None and result.coverage < self.threshold:
                failed_suites.append(f"{suite_name} ({result.coverage:.2f}%)")

        if failed_suites:
            print(f"❌ Coverage below threshold ({self.threshold}%):")
            for suite in failed_suites:
                print(f"   {suite}")
            return False

        return True


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Run tests on codebase")
    parser.add_argument(
        "--all",
        action="store_true",
        help="Run all test suites",
    )
    parser.add_argument(
        "--go",
        action="store_true",
        help="Run Go tests",
    )
    parser.add_argument(
        "--python",
        action="store_true",
        help="Run Python tests",
    )
    parser.add_argument(
        "--frontend",
        action="store_true",
        help="Run frontend tests",
    )
    parser.add_argument(
        "--coverage",
        action="store_true",
        help="Generate coverage reports",
    )
    parser.add_argument(
        "--watch",
        action="store_true",
        help="Run in watch mode (development)",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Verbose output",
    )
    parser.add_argument(
        "--threshold",
        type=int,
        default=80,
        help="Coverage threshold percentage (default: 80)",
    )

    args = parser.parse_args()

    # Determine which test suites to run
    suites = []
    if args.all:
        suites = None  # Run all
    else:
        if args.go:
            suites.append("go")
        if args.python:
            suites.append("python")
        if args.frontend:
            suites.append("frontend")

        if not suites:
            print("❌ Error: Specify at least one test suite or use --all")
            parser.print_help()
            sys.exit(1)

    # Create runner
    runner = TestRunner(
        coverage=args.coverage,
        watch=args.watch,
        verbose=args.verbose,
        threshold=args.threshold,
    )

    # Run tests
    print(f"\n{'='*70}")
    print(f"Running Tests - {len(suites) if suites else 'all suites'}")
    print(f"{'='*70}")

    if args.coverage:
        print(f"📊 Coverage reporting enabled (threshold: {args.threshold}%)")

    if args.watch:
        print("👁️  Watch mode enabled")

    results = runner.run_all_tests(suites)

    # Print summary
    runner.print_summary(results)

    # Check coverage threshold
    coverage_ok = runner.check_coverage_threshold(results)

    # Exit with error if any test failed or coverage below threshold
    if any(not r.success for r in results.values()) or not coverage_ok:
        sys.exit(1)


if __name__ == "__main__":
    main()
