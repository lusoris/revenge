#!/usr/bin/env python3
"""Docker configuration management script.

Manages Docker configurations:
- Sync Dockerfile base images from SOURCE_OF_TRUTH
- Sync docker-compose.yml service versions
- Build images
- Push images to registry
- Validate configurations

Requires:
- docker CLI installed
- docker-compose or docker compose installed

Usage:
    python scripts/automation/manage_docker.py --sync
    python scripts/automation/manage_docker.py --build
    python scripts/automation/manage_docker.py --push
    python scripts/automation/manage_docker.py --validate

Author: Automation System
Created: 2026-01-31
"""

import argparse
import re
import subprocess
import sys
from pathlib import Path


class DockerManager:
    """Manage Docker configurations."""

    def __init__(self, dry_run: bool = False, verbose: bool = False):
        """Initialize Docker manager.

        Args:
            dry_run: If True, print actions without executing
            verbose: If True, print detailed output
        """
        self.dry_run = dry_run
        self.verbose = verbose
        self.root = Path.cwd()

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

        if self.dry_run and cmd[0] in ["docker"]:
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

    def check_docker(self) -> bool:
        """Check if docker is installed.

        Returns:
            True if docker is available
        """
        result = self.run_command(["docker", "version"])
        return result.returncode == 0

    def parse_sot_versions(self) -> dict[str, str]:
        """Parse versions from SOURCE_OF_TRUTH.

        Returns:
            Dict mapping component to version
        """
        sot_file = self.root / "docs" / "dev" / "design" / "00_SOURCE_OF_TRUTH.md"
        if not sot_file.exists():
            print(f"⚠️  SOURCE_OF_TRUTH not found: {sot_file}")
            return {}

        versions = {}

        with open(sot_file) as f:
            content = f.read()

            # Extract versions from infrastructure table
            # Format: | PostgreSQL | 18+ | Database | ✅ |
            infra_matches = re.findall(
                r"\|\s+(\w+)\s+\|\s+([\d.+]+)\s+\|",
                content,
            )

            for component, version in infra_matches:
                versions[component.lower()] = version.rstrip("+")

        return versions

    def sync_dockerfile(self) -> bool:
        """Sync Dockerfile base images from SOURCE_OF_TRUTH.

        Returns:
            True if successful
        """
        print("\n🔄 Syncing Dockerfile from SOURCE_OF_TRUTH...")

        dockerfile = self.root / "Dockerfile"
        if not dockerfile.exists():
            print(f"⚠️  Dockerfile not found: {dockerfile}")
            return False

        versions = self.parse_sot_versions()
        if not versions:
            print("⚠️  No versions found in SOURCE_OF_TRUTH")
            return False

        print(f"Found versions: {versions}")

        # Read Dockerfile
        with open(dockerfile) as f:
            content = f.read()

        # Update Go version if found
        if "go" in versions:
            go_version = versions["go"]
            content = re.sub(
                r"FROM golang:([\d.]+)",
                f"FROM golang:{go_version}",
                content,
            )
            print(f"   Updated Go base image to {go_version}")

        # Write back if not dry-run
        if not self.dry_run:
            with open(dockerfile, "w") as f:
                f.write(content)
            print("✅ Dockerfile updated")
        else:
            print("[DRY-RUN] Would update Dockerfile")

        return True

    def sync_compose(self) -> bool:
        """Sync docker-compose.yml from SOURCE_OF_TRUTH.

        Returns:
            True if successful
        """
        print("\n🔄 Syncing docker-compose.yml from SOURCE_OF_TRUTH...")

        compose_file = self.root / "docker-compose.yml"
        if not compose_file.exists():
            print(f"⚠️  docker-compose.yml not found: {compose_file}")
            return False

        versions = self.parse_sot_versions()
        if not versions:
            print("⚠️  No versions found in SOURCE_OF_TRUTH")
            return False

        # Read compose file
        with open(compose_file) as f:
            content = f.read()

        # Update PostgreSQL version if found
        if "postgresql" in versions:
            pg_version = versions["postgresql"]
            content = re.sub(
                r"image: postgres:([\d.]+)",
                f"image: postgres:{pg_version}",
                content,
            )
            print(f"   Updated PostgreSQL image to {pg_version}")

        # Write back if not dry-run
        if not self.dry_run:
            with open(compose_file, "w") as f:
                f.write(content)
            print("✅ docker-compose.yml updated")
        else:
            print("[DRY-RUN] Would update docker-compose.yml")

        return True

    def validate_configs(self) -> bool:
        """Validate Docker configurations.

        Returns:
            True if valid
        """
        print("\n🔍 Validating Docker configurations...")

        success = True

        # Validate Dockerfile
        dockerfile = self.root / "Dockerfile"
        if dockerfile.exists():
            result = self.run_command([
                "docker", "build",
                "--dry-run",
                "-f", str(dockerfile),
                ".",
            ])

            if result.returncode == 0:
                print("✅ Dockerfile is valid")
            else:
                print(f"❌ Dockerfile validation failed:\n{result.stderr}")
                success = False
        else:
            print("⚠️  Dockerfile not found")

        # Validate docker-compose.yml
        compose_file = self.root / "docker-compose.yml"
        if compose_file.exists():
            result = self.run_command(["docker", "compose", "config", "--quiet"])

            if result.returncode == 0:
                print("✅ docker-compose.yml is valid")
            else:
                print(f"❌ docker-compose.yml validation failed:\n{result.stderr}")
                success = False
        else:
            print("⚠️  docker-compose.yml not found")

        return success

    def build_images(self, tag: str = "latest") -> bool:
        """Build Docker images.

        Args:
            tag: Image tag

        Returns:
            True if successful
        """
        print(f"\n🏗️  Building Docker images (tag: {tag})...")

        result = self.run_command([
            "docker", "build",
            "-t", f"revenge:{tag}",
            ".",
        ])

        if result.returncode != 0:
            print(f"❌ Build failed:\n{result.stderr}")
            return False

        print(f"✅ Images built successfully: revenge:{tag}")
        return True

    def push_images(self, registry: str, tag: str = "latest") -> bool:
        """Push Docker images to registry.

        Args:
            registry: Registry URL
            tag: Image tag

        Returns:
            True if successful
        """
        print(f"\n📤 Pushing images to {registry}...")

        # Tag for registry
        local_tag = f"revenge:{tag}"
        remote_tag = f"{registry}/revenge:{tag}"

        # Tag image
        result = self.run_command([
            "docker", "tag",
            local_tag,
            remote_tag,
        ])

        if result.returncode != 0:
            print(f"❌ Tagging failed:\n{result.stderr}")
            return False

        # Push image
        result = self.run_command([
            "docker", "push",
            remote_tag,
        ])

        if result.returncode != 0:
            print(f"❌ Push failed:\n{result.stderr}")
            return False

        print(f"✅ Images pushed successfully: {remote_tag}")
        return True

    def compose_up(self, detach: bool = True) -> bool:
        """Start services with docker-compose.

        Args:
            detach: Run in detached mode

        Returns:
            True if successful
        """
        print("\n▶️  Starting services with docker-compose...")

        cmd = ["docker", "compose", "up"]
        if detach:
            cmd.append("-d")

        result = self.run_command(cmd)

        if result.returncode != 0:
            print(f"❌ Failed to start services:\n{result.stderr}")
            return False

        print("✅ Services started successfully")
        return True

    def compose_down(self) -> bool:
        """Stop services with docker-compose.

        Returns:
            True if successful
        """
        print("\n⏹️  Stopping services with docker-compose...")

        result = self.run_command(["docker", "compose", "down"])

        if result.returncode != 0:
            print(f"❌ Failed to stop services:\n{result.stderr}")
            return False

        print("✅ Services stopped successfully")
        return True


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Manage Docker configurations")
    parser.add_argument(
        "--sync",
        action="store_true",
        help="Sync configs from SOURCE_OF_TRUTH",
    )
    parser.add_argument(
        "--validate",
        action="store_true",
        help="Validate Docker configurations",
    )
    parser.add_argument(
        "--build",
        action="store_true",
        help="Build Docker images",
    )
    parser.add_argument(
        "--push",
        action="store_true",
        help="Push images to registry",
    )
    parser.add_argument(
        "--up",
        action="store_true",
        help="Start services with docker-compose",
    )
    parser.add_argument(
        "--down",
        action="store_true",
        help="Stop services with docker-compose",
    )
    parser.add_argument(
        "--tag",
        default="latest",
        help="Image tag (default: latest)",
    )
    parser.add_argument(
        "--registry",
        default="ghcr.io/lusoris",
        help="Registry URL (default: ghcr.io/lusoris)",
    )
    parser.add_argument(
        "--foreground",
        action="store_true",
        help="Run docker-compose in foreground (use with --up)",
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
    manager = DockerManager(dry_run=args.dry_run, verbose=args.verbose)

    # Check docker
    if not manager.check_docker():
        print("❌ Error: docker not found")
        print("   Install: https://docs.docker.com/get-docker/")
        sys.exit(1)

    # Execute action
    success = True

    if args.sync:
        success = manager.sync_dockerfile() and manager.sync_compose()
    elif args.validate:
        success = manager.validate_configs()
    elif args.build:
        success = manager.build_images(args.tag)
    elif args.push:
        success = manager.push_images(args.registry, args.tag)
    elif args.up:
        success = manager.compose_up(detach=not args.foreground)
    elif args.down:
        success = manager.compose_down()
    else:
        print("❌ Error: Specify an action")
        parser.print_help()
        sys.exit(1)

    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
