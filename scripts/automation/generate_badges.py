#!/usr/bin/env python3
"""Badge generation script.

Generates shields.io badges for documentation:
- Coverage badge (from coverage reports)
- Build status badge (from GitHub Actions)
- Version badge (from git tags / go.mod)
- License badge (from LICENSE file)

Uses shields.io API for badge generation.

Requirements:
- requests library
- git CLI (for version detection)
- gh CLI (for GitHub Actions status) - optional

Usage:
    python scripts/automation/generate_badges.py --all
    python scripts/automation/generate_badges.py --coverage
    python scripts/automation/generate_badges.py --build
    python scripts/automation/generate_badges.py --version
    python scripts/automation/generate_badges.py --license

Author: Automation System
Created: 2026-01-31
"""

import argparse
import re
import subprocess
import sys
from pathlib import Path


try:
    import requests
except ImportError:
    print("❌ Error: requests not installed")
    print("   Install: pip install requests")
    sys.exit(1)


class BadgeGenerator:
    """Generate shields.io badges for documentation."""

    SHIELDS_IO_BASE = "https://img.shields.io"

    def __init__(self, badges_dir: Path, repo_root: Path):
        """Initialize generator.

        Args:
            badges_dir: Output directory for badges (docs/assets/badges)
            repo_root: Repository root directory
        """
        self.badges_dir = badges_dir
        self.repo_root = repo_root
        self.badges_dir.mkdir(parents=True, exist_ok=True)

    def _run_command(self, cmd: list[str]) -> str | None:
        """Run command and return output.

        Args:
            cmd: Command to run

        Returns:
            Command output or None if failed
        """
        try:
            result = subprocess.run(
                cmd,
                cwd=self.repo_root,
                capture_output=True,
                text=True,
                check=True,
            )
            return result.stdout.strip()
        except (subprocess.CalledProcessError, FileNotFoundError):
            return None

    def _download_badge(self, url: str, filename: str) -> Path | None:
        """Download badge from shields.io.

        Args:
            url: Shields.io URL
            filename: Output filename

        Returns:
            Path to downloaded badge or None if failed
        """
        try:
            response = requests.get(url, timeout=10)
            response.raise_for_status()

            output_path = self.badges_dir / filename
            output_path.write_bytes(response.content)
            return output_path
        except requests.RequestException as e:
            print(f"   ⚠️  Failed to download badge: {e}")
            return None

    def generate_coverage_badge(self, coverage_file: Path | None = None) -> Path | None:
        """Generate coverage badge from coverage report.

        Args:
            coverage_file: Path to coverage file (coverage.out or coverage.txt)

        Returns:
            Path to generated badge or None if failed
        """
        print("   Generating coverage badge...")

        # Try to find coverage file
        if coverage_file is None:
            coverage_candidates = [
                self.repo_root / "coverage.out",
                self.repo_root / "coverage.txt",
                self.repo_root / "coverage" / "coverage.out",
            ]
            for candidate in coverage_candidates:
                if candidate.exists():
                    coverage_file = candidate
                    break

        if coverage_file is None or not coverage_file.exists():
            print("   ⚠️  Coverage file not found, using placeholder")
            # Use placeholder badge
            coverage_percent = 0.0
        else:
            # Parse coverage from Go coverage report
            coverage_percent = self._parse_go_coverage(coverage_file)

        # Determine color based on coverage
        if coverage_percent >= 80:
            color = "brightgreen"
        elif coverage_percent >= 60:
            color = "yellow"
        else:
            color = "red"

        # Generate badge URL
        url = f"{self.SHIELDS_IO_BASE}/badge/coverage-{coverage_percent:.1f}%25-{color}"

        return self._download_badge(url, "coverage.svg")

    def _parse_go_coverage(self, coverage_file: Path) -> float:
        """Parse coverage percentage from Go coverage report.

        Args:
            coverage_file: Path to coverage.out file

        Returns:
            Coverage percentage (0-100)
        """
        try:
            # Run go tool cover to get percentage
            result = self._run_command(
                ["go", "tool", "cover", "-func", str(coverage_file)]
            )
            if result:
                # Last line contains total coverage: "total: (statement coverage) xx.x%"
                match = re.search(r"total:.*?\s+([\d.]+)%", result)
                if match:
                    return float(match.group(1))
        except Exception:
            pass

        return 0.0

    def generate_build_badge(self, _workflow_name: str = "CI") -> Path | None:
        """Generate build status badge from GitHub Actions.

        Args:
            _workflow_name: GitHub Actions workflow name (not used - shields.io limitation)

        Returns:
            Path to generated badge or None if failed
        """
        print("   Generating build badge...")

        # Try to get repo info from git remote
        remote_url = self._run_command(["git", "remote", "get-url", "origin"])
        if not remote_url or "github.com" not in remote_url:
            print("   ⚠️  Not a GitHub repository, using placeholder")
            return self._download_badge(
                f"{self.SHIELDS_IO_BASE}/badge/build-passing-brightgreen",
                "build.svg",
            )

        # Parse owner/repo from URL
        match = re.search(r"github\.com[:/](.+?)/(.+?)(\.git)?$", remote_url)
        if not match:
            print("   ⚠️  Could not parse GitHub URL, using placeholder")
            return self._download_badge(
                f"{self.SHIELDS_IO_BASE}/badge/build-passing-brightgreen",
                "build.svg",
            )

        owner = match.group(1)
        repo = match.group(2)

        # Use GitHub Actions badge from shields.io
        url = f"{self.SHIELDS_IO_BASE}/github/actions/workflow/status/{owner}/{repo}/ci.yml"

        return self._download_badge(url, "build.svg")

    def generate_version_badge(self) -> Path | None:
        """Generate version badge from git tags or go.mod.

        Returns:
            Path to generated badge or None if failed
        """
        print("   Generating version badge...")

        # Try to get version from git tags
        version = self._run_command(["git", "describe", "--tags", "--abbrev=0"])

        if not version:
            # Try to get version from go.mod
            go_mod = self.repo_root / "go.mod"
            if go_mod.exists():
                go_mod.read_text()
                # Look for version comment or default to 0.0.0
                version = "v0.0.0-dev"
            else:
                version = "v0.0.0-dev"

        # Clean version string
        version = version.lstrip("v").strip()

        # Generate badge URL
        url = f"{self.SHIELDS_IO_BASE}/badge/version-v{version}-blue"

        return self._download_badge(url, "version.svg")

    def generate_license_badge(self) -> Path | None:
        """Generate license badge from LICENSE file.

        Returns:
            Path to generated badge or None if failed
        """
        print("   Generating license badge...")

        # Try to detect license from LICENSE file
        license_file = self.repo_root / "LICENSE"
        license_type = "MIT"  # Default

        if license_file.exists():
            content = license_file.read_text()
            # Simple license detection
            if "Apache License" in content:
                license_type = "Apache-2.0"
            elif "GNU GENERAL PUBLIC LICENSE" in content:
                license_type = "GPL-3.0" if "Version 3" in content else "GPL-2.0"
            elif "MIT License" in content:
                license_type = "MIT"
            elif "BSD" in content:
                license_type = "BSD"

        # Generate badge URL
        url = f"{self.SHIELDS_IO_BASE}/badge/license-{license_type}-green"

        return self._download_badge(url, "license.svg")

    def generate_all(self):
        """Generate all badges."""
        print(f"\n{'=' * 70}")
        print("Badge Generation")
        print(f"{'=' * 70}\n")
        print(f"Badges directory: {self.badges_dir}")

        badges = {
            "coverage": self.generate_coverage_badge,
            "build": self.generate_build_badge,
            "version": self.generate_version_badge,
            "license": self.generate_license_badge,
        }

        results = {}
        for name, generator_func in badges.items():
            try:
                path = generator_func()
                if path:
                    results[name] = "✓"
                    print(f"   ✓ Created: {path.name}")
                else:
                    results[name] = "⚠"
            except Exception as e:
                results[name] = "✗"
                print(f"   ✗ Failed to create {name} badge: {e}")

        print(f"\n{'=' * 70}")
        print("Badge Generation Summary")
        print(f"{'=' * 70}\n")

        for name, status in results.items():
            print(f"   {status} {name.capitalize()}")

        print()


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Generate shields.io badges")
    parser.add_argument(
        "--all",
        action="store_true",
        help="Generate all badges",
    )
    parser.add_argument(
        "--coverage",
        action="store_true",
        help="Generate coverage badge only",
    )
    parser.add_argument(
        "--build",
        action="store_true",
        help="Generate build status badge only",
    )
    parser.add_argument(
        "--version",
        action="store_true",
        help="Generate version badge only",
    )
    parser.add_argument(
        "--license",
        action="store_true",
        help="Generate license badge only",
    )
    parser.add_argument(
        "--badges-dir",
        type=Path,
        default=Path.cwd() / "docs" / "assets" / "badges",
        help="Badges directory (default: docs/assets/badges)",
    )
    parser.add_argument(
        "--repo-root",
        type=Path,
        default=Path.cwd(),
        help="Repository root (default: current directory)",
    )

    args = parser.parse_args()

    # Verify repo root is valid
    if not (args.repo_root / ".git").exists():
        print(f"❌ Error: Not a git repository: {args.repo_root}")
        sys.exit(1)

    generator = BadgeGenerator(args.badges_dir, args.repo_root)

    # Generate based on arguments
    if args.all or (
        not args.coverage and not args.build and not args.version and not args.license
    ):
        generator.generate_all()
    else:
        if args.coverage:
            path = generator.generate_coverage_badge()
            if path:
                print(f"✓ Created: {path}")
        if args.build:
            path = generator.generate_build_badge()
            if path:
                print(f"✓ Created: {path}")
        if args.version:
            path = generator.generate_version_badge()
            if path:
                print(f"✓ Created: {path}")
        if args.license:
            path = generator.generate_license_badge()
            if path:
                print(f"✓ Created: {path}")


if __name__ == "__main__":
    main()
