#!/usr/bin/env python3
"""System health checker script.

Checks health of various system components:
- Automation system (dependencies, templates, schemas)
- Backend services (database, cache, search)
- Frontend build
- External integrations
- Resource usage

Supports:
- Individual component checks
- Full system health scan
- JSON output for monitoring tools
- Alert on failures (GitHub issues)

Usage:
    python scripts/automation/check_health.py --all
    python scripts/automation/check_health.py --automation
    python scripts/automation/check_health.py --services
    python scripts/automation/check_health.py --json

Author: Automation System
Created: 2026-01-31
"""

import argparse
import json
import subprocess
import sys
from dataclasses import dataclass, field
from pathlib import Path


@dataclass
class HealthCheck:
    """Result of a health check."""

    component: str
    status: str  # "healthy", "degraded", "unhealthy"
    message: str
    details: dict = field(default_factory=dict)


class HealthChecker:
    """Check system health."""

    def __init__(self, verbose: bool = False):
        """Initialize health checker.

        Args:
            verbose: If True, print detailed output
        """
        self.verbose = verbose
        self.root = Path.cwd()
        self.results: list[HealthCheck] = []

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
                timeout=10,
            )
            return result

        except (FileNotFoundError, subprocess.TimeoutExpired) as e:
            return subprocess.CompletedProcess(
                args=cmd,
                returncode=1,
                stdout="",
                stderr=str(e),
            )

    def check_python_dependencies(self) -> HealthCheck:
        """Check Python dependencies.

        Returns:
            HealthCheck result
        """
        if self.verbose:
            print("Checking Python dependencies...")

        req_file = self.root / "scripts" / "requirements.txt"
        if not req_file.exists():
            return HealthCheck(
                component="python-deps",
                status="degraded",
                message="requirements.txt not found",
            )

        # Try to import key dependencies
        missing = []
        try:
            import yaml  # noqa: F401
        except ImportError:
            missing.append("PyYAML")

        try:
            import pytest  # noqa: F401
        except ImportError:
            missing.append("pytest")

        if missing:
            return HealthCheck(
                component="python-deps",
                status="degraded",
                message=f"Missing dependencies: {', '.join(missing)}",
                details={"missing": missing},
            )

        return HealthCheck(
            component="python-deps",
            status="healthy",
            message="All Python dependencies available",
        )

    def check_templates(self) -> HealthCheck:
        """Check Jinja2 templates.

        Returns:
            HealthCheck result
        """
        if self.verbose:
            print("Checking templates...")

        templates_dir = self.root / "templates"
        if not templates_dir.exists():
            return HealthCheck(
                component="templates",
                status="unhealthy",
                message="Templates directory not found",
            )

        # Count template files
        template_files = list(templates_dir.rglob("*.jinja2"))
        if not template_files:
            return HealthCheck(
                component="templates",
                status="degraded",
                message="No template files found",
            )

        return HealthCheck(
            component="templates",
            status="healthy",
            message=f"Found {len(template_files)} template files",
            details={"count": len(template_files)},
        )

    def check_schemas(self) -> HealthCheck:
        """Check JSON schemas.

        Returns:
            HealthCheck result
        """
        if self.verbose:
            print("Checking schemas...")

        schemas_dir = self.root / "schemas"
        if not schemas_dir.exists():
            return HealthCheck(
                component="schemas",
                status="degraded",
                message="Schemas directory not found",
            )

        # Count schema files
        schema_files = list(schemas_dir.glob("*.schema.json"))
        if not schema_files:
            return HealthCheck(
                component="schemas",
                status="degraded",
                message="No schema files found",
            )

        return HealthCheck(
            component="schemas",
            status="healthy",
            message=f"Found {len(schema_files)} schema files",
            details={"count": len(schema_files)},
        )

    def check_database(self) -> HealthCheck:
        """Check database connectivity.

        Returns:
            HealthCheck result
        """
        if self.verbose:
            print("Checking database...")

        # Try to connect to PostgreSQL
        result = self.run_command(["docker", "compose", "ps", "postgres"])

        if result.returncode != 0:
            return HealthCheck(
                component="database",
                status="degraded",
                message="Cannot check database status (docker-compose not available)",
            )

        if "Up" in result.stdout:
            return HealthCheck(
                component="database",
                status="healthy",
                message="PostgreSQL is running",
            )

        return HealthCheck(
            component="database",
            status="unhealthy",
            message="PostgreSQL is not running",
        )

    def check_cache(self) -> HealthCheck:
        """Check cache service.

        Returns:
            HealthCheck result
        """
        if self.verbose:
            print("Checking cache...")

        # Try to check Dragonfly/Redis
        result = self.run_command(["docker", "compose", "ps", "dragonfly"])

        if result.returncode != 0:
            return HealthCheck(
                component="cache",
                status="degraded",
                message="Cannot check cache status (docker-compose not available)",
            )

        if "Up" in result.stdout:
            return HealthCheck(
                component="cache",
                status="healthy",
                message="Dragonfly cache is running",
            )

        return HealthCheck(
            component="cache",
            status="degraded",
            message="Dragonfly cache is not running",
        )

    def check_search(self) -> HealthCheck:
        """Check search service.

        Returns:
            HealthCheck result
        """
        if self.verbose:
            print("Checking search...")

        # Try to check Typesense
        result = self.run_command(["docker", "compose", "ps", "typesense"])

        if result.returncode != 0:
            return HealthCheck(
                component="search",
                status="degraded",
                message="Cannot check search status (docker-compose not available)",
            )

        if "Up" in result.stdout:
            return HealthCheck(
                component="search",
                status="healthy",
                message="Typesense search is running",
            )

        return HealthCheck(
            component="search",
            status="degraded",
            message="Typesense search is not running",
        )

    def check_frontend_build(self) -> HealthCheck:
        """Check frontend build.

        Returns:
            HealthCheck result
        """
        if self.verbose:
            print("Checking frontend build...")

        frontend_dir = self.root / "frontend"
        if not frontend_dir.exists():
            return HealthCheck(
                component="frontend",
                status="degraded",
                message="Frontend directory not found",
            )

        # Check if node_modules exists
        node_modules = frontend_dir / "node_modules"
        if not node_modules.exists():
            return HealthCheck(
                component="frontend",
                status="degraded",
                message="node_modules not found (run npm install)",
            )

        # Check if build directory exists
        build_dir = frontend_dir / "build"
        if not build_dir.exists():
            return HealthCheck(
                component="frontend",
                status="degraded",
                message="Build directory not found (run npm run build)",
            )

        return HealthCheck(
            component="frontend",
            status="healthy",
            message="Frontend build is ready",
        )

    def check_resource_usage(self) -> HealthCheck:
        """Check system resource usage.

        Returns:
            HealthCheck result
        """
        if self.verbose:
            print("Checking resource usage...")

        # Check disk usage
        result = self.run_command(["df", "-h", str(self.root)])

        if result.returncode != 0:
            return HealthCheck(
                component="resources",
                status="degraded",
                message="Cannot check resource usage",
            )

        # Parse disk usage percentage
        lines = result.stdout.split("\n")
        if len(lines) > 1:
            parts = lines[1].split()
            if len(parts) >= 5:
                usage = parts[4].rstrip("%")
                try:
                    usage_pct = int(usage)
                    if usage_pct > 90:
                        status = "unhealthy"
                    elif usage_pct > 75:
                        status = "degraded"
                    else:
                        status = "healthy"

                    return HealthCheck(
                        component="resources",
                        status=status,
                        message=f"Disk usage: {usage}%",
                        details={"disk_usage_pct": usage_pct},
                    )
                except ValueError:
                    pass

        return HealthCheck(
            component="resources",
            status="healthy",
            message="Resource usage within limits",
        )

    def check_automation_system(self) -> list[HealthCheck]:
        """Check automation system components.

        Returns:
            List of HealthCheck results
        """
        return [
            self.check_python_dependencies(),
            self.check_templates(),
            self.check_schemas(),
        ]

    def check_backend_services(self) -> list[HealthCheck]:
        """Check backend services.

        Returns:
            List of HealthCheck results
        """
        return [
            self.check_database(),
            self.check_cache(),
            self.check_search(),
        ]

    def check_all(self) -> list[HealthCheck]:
        """Run all health checks.

        Returns:
            List of all HealthCheck results
        """
        results = []
        results.extend(self.check_automation_system())
        results.extend(self.check_backend_services())
        results.append(self.check_frontend_build())
        results.append(self.check_resource_usage())
        return results

    def print_summary(self, results: list[HealthCheck]):
        """Print health check summary.

        Args:
            results: List of HealthCheck results
        """
        print(f"\n{'='*70}")
        print("System Health Check")
        print(f"{'='*70}\n")

        healthy = sum(1 for r in results if r.status == "healthy")
        degraded = sum(1 for r in results if r.status == "degraded")
        unhealthy = sum(1 for r in results if r.status == "unhealthy")

        for result in results:
            if result.status == "healthy":
                icon = "✅"
            elif result.status == "degraded":
                icon = "⚠️"
            else:
                icon = "❌"

            print(f"{icon} {result.component:20} {result.status.upper():12} {result.message}")

        print(f"\n{'='*70}")
        print(f"Summary: Healthy:{healthy} Degraded:{degraded} Unhealthy:{unhealthy}")

        if unhealthy > 0:
            print("❌ System health: UNHEALTHY")
        elif degraded > 0:
            print("⚠️  System health: DEGRADED")
        else:
            print("✅ System health: HEALTHY")

        print(f"{'='*70}\n")

    def create_github_issue(self, results: list[HealthCheck]) -> bool:
        """Create GitHub issue for unhealthy components.

        Args:
            results: List of HealthCheck results

        Returns:
            True if issue created
        """
        unhealthy = [r for r in results if r.status == "unhealthy"]
        if not unhealthy:
            return False

        print("\n📝 Creating GitHub issue for unhealthy components...")

        # Build issue body
        issue_body = "## System Health Alert\n\n"
        issue_body += "The following components are unhealthy:\n\n"

        for result in unhealthy:
            issue_body += f"- **{result.component}**: {result.message}\n"

        issue_body += "\n**Automated health check failed.**\n"
        issue_body += "\n🤖 Generated with [Claude Code](https://claude.com/claude-code)"

        # Create issue using gh CLI
        try:
            result = subprocess.run(
                [
                    "gh",
                    "issue",
                    "create",
                    "--title",
                    "🚨 System Health Alert - Unhealthy Components",
                    "--body",
                    issue_body,
                    "--label",
                    "bug,automation",
                ],
                capture_output=True,
                text=True,
                check=False,
            )

            if result.returncode == 0:
                print("✅ GitHub issue created")
                return True
            else:
                print(f"❌ Failed to create issue: {result.stderr}")
                return False

        except FileNotFoundError:
            print("❌ gh CLI not found, cannot create issue")
            return False


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Check system health")
    parser.add_argument(
        "--all",
        action="store_true",
        help="Check all components",
    )
    parser.add_argument(
        "--automation",
        action="store_true",
        help="Check automation system",
    )
    parser.add_argument(
        "--services",
        action="store_true",
        help="Check backend services",
    )
    parser.add_argument(
        "--frontend",
        action="store_true",
        help="Check frontend build",
    )
    parser.add_argument(
        "--resources",
        action="store_true",
        help="Check resource usage",
    )
    parser.add_argument(
        "--json",
        action="store_true",
        help="Output in JSON format",
    )
    parser.add_argument(
        "--alert",
        action="store_true",
        help="Create GitHub issue on failure",
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Verbose output",
    )

    args = parser.parse_args()

    # Create checker
    checker = HealthChecker(verbose=args.verbose)

    # Run checks
    results = []

    if args.all or not any([args.automation, args.services, args.frontend, args.resources]):
        results = checker.check_all()
    else:
        if args.automation:
            results.extend(checker.check_automation_system())
        if args.services:
            results.extend(checker.check_backend_services())
        if args.frontend:
            results.append(checker.check_frontend_build())
        if args.resources:
            results.append(checker.check_resource_usage())

    # Output results
    if args.json:
        output = {
            "status": "healthy" if all(r.status == "healthy" for r in results) else "degraded",
            "components": [
                {
                    "component": r.component,
                    "status": r.status,
                    "message": r.message,
                    "details": r.details,
                }
                for r in results
            ],
        }
        print(json.dumps(output, indent=2))
    else:
        checker.print_summary(results)

    # Create alert if requested
    if args.alert:
        checker.create_github_issue(results)

    # Exit with error if any component unhealthy
    if any(r.status == "unhealthy" for r in results):
        sys.exit(1)


if __name__ == "__main__":
    main()
