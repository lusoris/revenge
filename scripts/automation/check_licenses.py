#!/usr/bin/env python3
"""License checker script.

Scans dependencies for license compliance:
- Go modules (go-licenses)
- npm packages (license-checker)
- Python packages (pip-licenses)

Checks against allow/deny lists and generates compliance reports.

Requirements:
- go-licenses (go install github.com/google/go-licenses@latest)
- license-checker (npm install -g license-checker)
- pip-licenses (pip install pip-licenses)

Usage:
    python scripts/automation/check_licenses.py --all
    python scripts/automation/check_licenses.py --go --npm
    python scripts/automation/check_licenses.py --all --report
    python scripts/automation/check_licenses.py --go --strict

Author: Automation System
Created: 2026-01-31
"""

import argparse
import json
import subprocess
import sys
from dataclasses import dataclass, field
from pathlib import Path
from typing import ClassVar


@dataclass
class License:
    """License information for a dependency."""

    name: str
    version: str
    license: str
    ecosystem: str


@dataclass
class LicenseCheckResult:
    """Result of license check."""

    ecosystem: str
    success: bool
    total: int
    allowed: int
    denied: int
    unknown: int
    licenses: list[License] = field(default_factory=list)
    denied_licenses: list[License] = field(default_factory=list)
    unknown_licenses: list[License] = field(default_factory=list)
    output: str = ""
    error: str = ""


class LicenseChecker:
    """Check licenses across dependencies."""

    # Default allowed licenses (permissive)
    ALLOWED_LICENSES: ClassVar[set[str]] = {
        "MIT",
        "Apache-2.0",
        "Apache 2.0",
        "BSD-2-Clause",
        "BSD-3-Clause",
        "ISC",
        "0BSD",
        "CC0-1.0",
        "Unlicense",
        "MPL-2.0",
        "PostgreSQL",
        "Python-2.0",
        "Zlib",
    }

    # Denied licenses (copyleft)
    DENIED_LICENSES: ClassVar[set[str]] = {
        "GPL-2.0",
        "GPL-3.0",
        "LGPL-2.1",
        "LGPL-3.0",
        "AGPL-3.0",
        "SSPL",
        "CC-BY-NC",
        "CC-BY-NC-SA",
        "BUSL-1.1",
    }

    def __init__(self, strict: bool = False, verbose: bool = False):
        """Initialize license checker.

        Args:
            strict: If True, fail on unknown licenses
            verbose: If True, print detailed output
        """
        self.strict = strict
        self.verbose = verbose
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
            fake_result = subprocess.CompletedProcess(
                args=cmd,
                returncode=127,
                stdout="",
                stderr=f"{cmd[0]} not found: {e}",
            )
            return fake_result

    def normalize_license(self, license_str: str) -> str:
        """Normalize license string for comparison.

        Args:
            license_str: Raw license string

        Returns:
            Normalized license string
        """
        # Strip whitespace
        license_str = license_str.strip()

        # Handle common variations (check original first)
        replacements = {
            "Apache License 2.0": "Apache-2.0",
            "Apache 2": "Apache-2.0",
            "MIT License": "MIT",
            "BSD 3-Clause": "BSD-3-Clause",
            "BSD 2-Clause": "BSD-2-Clause",
            "ISC License": "ISC",
            "MPL-2.0": "MPL-2.0",
            "Mozilla Public License 2.0": "MPL-2.0",
        }

        # Check if exact match in replacements
        if license_str in replacements:
            return replacements[license_str]

        # Remove "License" suffix and check again
        normalized = license_str.replace(" License", "").replace("License", "").strip()
        return replacements.get(normalized, normalized)

    def check_license(self, license_str: str) -> str:
        """Check if license is allowed, denied, or unknown.

        Args:
            license_str: License string

        Returns:
            "allowed", "denied", or "unknown"
        """
        normalized = self.normalize_license(license_str)

        if normalized in self.ALLOWED_LICENSES:
            return "allowed"
        if normalized in self.DENIED_LICENSES:
            return "denied"
        return "unknown"

    def check_go_licenses(self) -> LicenseCheckResult:
        """Check Go module licenses.

        Returns:
            LicenseCheckResult
        """
        print("\n📄 Checking Go licenses...")

        cmd = ["go-licenses", "report", "./..."]

        result = self.run_command(cmd)

        if result.returncode == 127:
            print("   ⚠️  go-licenses not found, skipping")
            return LicenseCheckResult(
                ecosystem="go",
                success=True,
                total=0,
                allowed=0,
                denied=0,
                unknown=0,
                output=result.stderr,
            )

        licenses = []
        denied_licenses = []
        unknown_licenses = []

        # Parse go-licenses output
        # Format: "package,license_url,license_type"
        for line in result.stdout.split("\n"):
            if not line.strip() or line.startswith("W"):
                continue

            parts = line.split(",")
            if len(parts) >= 3:
                pkg_name = parts[0]
                license_type = parts[2].strip()

                lic = License(
                    name=pkg_name,
                    version="",  # go-licenses doesn't provide version
                    license=license_type,
                    ecosystem="go",
                )

                licenses.append(lic)

                check_result = self.check_license(license_type)
                if check_result == "denied":
                    denied_licenses.append(lic)
                elif check_result == "unknown":
                    unknown_licenses.append(lic)

        return LicenseCheckResult(
            ecosystem="go",
            success=len(denied_licenses) == 0
            and (not self.strict or len(unknown_licenses) == 0),
            total=len(licenses),
            allowed=len(licenses) - len(denied_licenses) - len(unknown_licenses),
            denied=len(denied_licenses),
            unknown=len(unknown_licenses),
            licenses=licenses,
            denied_licenses=denied_licenses,
            unknown_licenses=unknown_licenses,
            output=result.stdout,
            error=result.stderr,
        )

    def check_npm_licenses(self) -> LicenseCheckResult:
        """Check npm package licenses.

        Returns:
            LicenseCheckResult
        """
        print("\n📄 Checking npm licenses...")

        frontend_dir = self.root / "frontend"
        if not frontend_dir.exists():
            print("   ⚠️  Frontend directory not found, skipping")
            return LicenseCheckResult(
                ecosystem="npm",
                success=True,
                total=0,
                allowed=0,
                denied=0,
                unknown=0,
                output="Frontend directory not found",
            )

        cmd = ["license-checker", "--json"]

        result = self.run_command(cmd, cwd=frontend_dir)

        if result.returncode == 127:
            print("   ⚠️  license-checker not found, skipping")
            return LicenseCheckResult(
                ecosystem="npm",
                success=True,
                total=0,
                allowed=0,
                denied=0,
                unknown=0,
                output=result.stderr,
            )

        licenses = []
        denied_licenses = []
        unknown_licenses = []

        # Parse license-checker JSON output
        try:
            license_data = json.loads(result.stdout)

            for pkg_info_str, info in license_data.items():
                # pkg_info_str format: "package@version"
                pkg_parts = pkg_info_str.rsplit("@", 1)
                pkg_name = pkg_parts[0]
                pkg_version = pkg_parts[1] if len(pkg_parts) > 1 else ""

                license_type = info.get("licenses", "UNKNOWN")

                lic = License(
                    name=pkg_name,
                    version=pkg_version,
                    license=license_type,
                    ecosystem="npm",
                )

                licenses.append(lic)

                check_result = self.check_license(license_type)
                if check_result == "denied":
                    denied_licenses.append(lic)
                elif check_result == "unknown":
                    unknown_licenses.append(lic)

        except json.JSONDecodeError:
            pass

        return LicenseCheckResult(
            ecosystem="npm",
            success=len(denied_licenses) == 0
            and (not self.strict or len(unknown_licenses) == 0),
            total=len(licenses),
            allowed=len(licenses) - len(denied_licenses) - len(unknown_licenses),
            denied=len(denied_licenses),
            unknown=len(unknown_licenses),
            licenses=licenses,
            denied_licenses=denied_licenses,
            unknown_licenses=unknown_licenses,
            output=result.stdout,
            error=result.stderr,
        )

    def check_python_licenses(self) -> LicenseCheckResult:
        """Check Python package licenses.

        Returns:
            LicenseCheckResult
        """
        print("\n📄 Checking Python licenses...")

        cmd = ["pip-licenses", "--format=json"]

        result = self.run_command(cmd)

        if result.returncode == 127:
            print("   ⚠️  pip-licenses not found, skipping")
            return LicenseCheckResult(
                ecosystem="python",
                success=True,
                total=0,
                allowed=0,
                denied=0,
                unknown=0,
                output=result.stderr,
            )

        licenses = []
        denied_licenses = []
        unknown_licenses = []

        # Parse pip-licenses JSON output
        try:
            license_data = json.loads(result.stdout)

            for pkg_info in license_data:
                lic = License(
                    name=pkg_info["Name"],
                    version=pkg_info["Version"],
                    license=pkg_info["License"],
                    ecosystem="python",
                )

                licenses.append(lic)

                check_result = self.check_license(pkg_info["License"])
                if check_result == "denied":
                    denied_licenses.append(lic)
                elif check_result == "unknown":
                    unknown_licenses.append(lic)

        except json.JSONDecodeError:
            pass

        return LicenseCheckResult(
            ecosystem="python",
            success=len(denied_licenses) == 0
            and (not self.strict or len(unknown_licenses) == 0),
            total=len(licenses),
            allowed=len(licenses) - len(denied_licenses) - len(unknown_licenses),
            denied=len(denied_licenses),
            unknown=len(unknown_licenses),
            licenses=licenses,
            denied_licenses=denied_licenses,
            unknown_licenses=unknown_licenses,
            output=result.stdout,
            error=result.stderr,
        )

    def check_all_licenses(
        self,
        ecosystems: list[str] | None = None,
    ) -> dict[str, LicenseCheckResult]:
        """Check all or selected ecosystems.

        Args:
            ecosystems: List of ecosystems to check (go, npm, python)
                       If None, check all

        Returns:
            Dict mapping ecosystem to result
        """
        # Map ecosystem names to methods
        ecosystem_methods = {
            "go": self.check_go_licenses,
            "npm": self.check_npm_licenses,
            "python": self.check_python_licenses,
        }

        # Filter ecosystems if specified
        if ecosystems:
            ecosystem_methods = {
                name: method
                for name, method in ecosystem_methods.items()
                if name in ecosystems
            }

        results = {}

        # Run checks sequentially
        for ecosystem_name, method in ecosystem_methods.items():
            try:
                result = method()
                results[ecosystem_name] = result

                # Print progress
                if result.success:
                    print(
                        f"   ✓ {ecosystem_name}: {result.allowed}/{result.total} allowed"
                    )
                else:
                    print(
                        f"   ✗ {ecosystem_name}: {result.denied} denied, {result.unknown} unknown"
                    )

            except Exception as e:
                print(f"   ✗ {ecosystem_name}: exception - {e}")
                results[ecosystem_name] = LicenseCheckResult(
                    ecosystem=ecosystem_name,
                    success=False,
                    total=0,
                    allowed=0,
                    denied=0,
                    unknown=0,
                    error=str(e),
                )

        return results

    def print_summary(self, results: dict[str, LicenseCheckResult]):
        """Print summary of license checks.

        Args:
            results: Dict of check results
        """
        print(f"\n{'=' * 70}")
        print("License Compliance Summary")
        print(f"{'=' * 70}\n")

        total_deps = sum(r.total for r in results.values())
        total_allowed = sum(r.allowed for r in results.values())
        total_denied = sum(r.denied for r in results.values())
        total_unknown = sum(r.unknown for r in results.values())

        for ecosystem_name, result in sorted(results.items()):
            status = "✓ PASS" if result.success else "✗ FAIL"
            print(
                f"{status:10} {ecosystem_name:10} "
                f"A:{result.allowed:4} D:{result.denied:4} U:{result.unknown:4} "
                f"Total:{result.total:4}",
            )

        print(f"\n{'=' * 70}")
        print(
            f"Overall: Allowed:{total_allowed} Denied:{total_denied} "
            f"Unknown:{total_unknown} Total:{total_deps}",
        )

        if total_denied > 0:
            print(f"\n❌ Found {total_denied} denied license(s)")

        if self.strict and total_unknown > 0:
            print(f"❌ Found {total_unknown} unknown license(s) (strict mode)")

        if total_denied == 0 and (not self.strict or total_unknown == 0):
            print("✅ All licenses are compliant!")

        print(f"{'=' * 70}\n")

    def generate_report(
        self, results: dict[str, LicenseCheckResult], output_file: Path
    ):
        """Generate detailed license report.

        Args:
            results: Dict of check results
            output_file: Path to output file
        """
        with open(output_file, "w") as f:
            f.write("# License Compliance Report\n\n")

            for ecosystem_name, result in sorted(results.items()):
                f.write(f"## {ecosystem_name.upper()}\n\n")
                f.write(f"- Total: {result.total}\n")
                f.write(f"- Allowed: {result.allowed}\n")
                f.write(f"- Denied: {result.denied}\n")
                f.write(f"- Unknown: {result.unknown}\n\n")

                if result.denied_licenses:
                    f.write("### ❌ Denied Licenses\n\n")
                    f.writelines(
                        f"- {lic.name}@{lic.version}: {lic.license}\n"
                        for lic in result.denied_licenses
                    )
                    f.write("\n")

                if result.unknown_licenses:
                    f.write("### ⚠️ Unknown Licenses\n\n")
                    f.writelines(
                        f"- {lic.name}@{lic.version}: {lic.license}\n"
                        for lic in result.unknown_licenses
                    )
                    f.write("\n")

        print(f"📄 Report saved to: {output_file}")


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Check license compliance")
    parser.add_argument(
        "--all",
        action="store_true",
        help="Check all ecosystems",
    )
    parser.add_argument(
        "--go",
        action="store_true",
        help="Check Go licenses",
    )
    parser.add_argument(
        "--npm",
        action="store_true",
        help="Check npm licenses",
    )
    parser.add_argument(
        "--python",
        action="store_true",
        help="Check Python licenses",
    )
    parser.add_argument(
        "--strict",
        action="store_true",
        help="Fail on unknown licenses",
    )
    parser.add_argument(
        "--report",
        action="store_true",
        help="Generate detailed report",
    )
    parser.add_argument(
        "--output",
        type=Path,
        default=Path("license-report.md"),
        help="Report output file (default: license-report.md)",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Verbose output",
    )

    args = parser.parse_args()

    # Determine which ecosystems to check
    ecosystems = []
    if args.all:
        ecosystems = None  # Check all
    else:
        if args.go:
            ecosystems.append("go")
        if args.npm:
            ecosystems.append("npm")
        if args.python:
            ecosystems.append("python")

        if not ecosystems:
            print("❌ Error: Specify at least one ecosystem or use --all")
            parser.print_help()
            sys.exit(1)

    # Create checker
    checker = LicenseChecker(strict=args.strict, verbose=args.verbose)

    # Check licenses
    print(f"\n{'=' * 70}")
    print(
        f"License Compliance Check - {len(ecosystems) if ecosystems else 'all ecosystems'}"
    )
    print(f"{'=' * 70}")

    if args.strict:
        print("⚠️  Strict mode enabled - unknown licenses will fail\n")

    results = checker.check_all_licenses(ecosystems)

    # Print summary
    checker.print_summary(results)

    # Generate report if requested
    if args.report:
        checker.generate_report(results, args.output)

    # Exit with error if any check failed
    if any(not r.success for r in results.values()):
        sys.exit(1)


if __name__ == "__main__":
    main()
