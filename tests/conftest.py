"""Shared pytest fixtures for pipeline tests."""

from __future__ import annotations

import shutil
import tempfile
from pathlib import Path
from typing import Generator

import pytest


@pytest.fixture
def temp_dir() -> Generator[Path, None, None]:
    """Create a temporary directory for tests."""
    temp_path = Path(tempfile.mkdtemp())
    yield temp_path
    shutil.rmtree(temp_path)


@pytest.fixture
def mock_design_dir(temp_dir: Path) -> Path:
    """Create a mock design directory structure."""
    design_dir = temp_dir / "docs" / "dev" / "design"
    design_dir.mkdir(parents=True)

    # Create SOT file
    sot_file = design_dir / "00_SOURCE_OF_TRUTH.md"
    sot_file.write_text(
        """# Source of Truth

> Master reference for the project

## Packages

| Package | Version |
|---------|---------|
| pgx | v5.7.2 |
| fx | v1.23.0 |
"""
    )

    # Create architecture directory
    arch_dir = design_dir / "architecture"
    arch_dir.mkdir()

    arch_file = arch_dir / "01_ARCHITECTURE.md"
    arch_file.write_text(
        """# Architecture

> Core system architecture

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete |
| Sources | 🟡 | Partial |

## Overview

This document describes the architecture.
"""
    )

    # Create services directory
    services_dir = design_dir / "services"
    services_dir.mkdir()

    auth_file = services_dir / "AUTH.md"
    auth_file.write_text(
        """# Authentication Service

> Handles user authentication

**Module**: `internal/service/auth`

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🟡 | In progress |
| Sources | 🔴 | Not started |

## Overview

Authentication service for the application.

See [Architecture](../architecture/01_ARCHITECTURE.md) for context.
"""
    )

    return design_dir


@pytest.fixture
def mock_sources_dir(temp_dir: Path) -> Path:
    """Create a mock sources directory structure."""
    sources_dir = temp_dir / "docs" / "dev" / "sources"
    sources_dir.mkdir(parents=True)

    # Create SOURCES.yaml
    sources_yaml = sources_dir / "SOURCES.yaml"
    sources_yaml.write_text(
        """sources:
  database:
    - id: pgx
      name: pgx
      url: https://pkg.go.dev/github.com/jackc/pgx/v5
      type: html
      output: database/pgx.md

  tooling:
    - id: fx
      name: Uber Fx
      url: https://pkg.go.dev/go.uber.org/fx
      type: html
      output: tooling/fx.md

fetch_config:
  delay_between_requests: 1
  timeout: 30
  retry_count: 3
"""
    )

    # Create INDEX.yaml
    index_yaml = sources_dir / "INDEX.yaml"
    index_yaml.write_text(
        """last_updated: "2026-01-31T00:00:00Z"
total_sources: 2
successful: 2
unchanged: 0
failed: 0
skipped: 0

sources:
  pgx:
    name: pgx
    url: https://pkg.go.dev/github.com/jackc/pgx/v5
    output: database/pgx.md
    status: success
    fetched_at: "2026-01-31T00:00:00Z"
    content_hash: abc123

  fx:
    name: Uber Fx
    url: https://pkg.go.dev/go.uber.org/fx
    output: tooling/fx.md
    status: success
    fetched_at: "2026-01-31T00:00:00Z"
    content_hash: def456
"""
    )

    return sources_dir


@pytest.fixture
def project_root() -> Path:
    """Get the actual project root."""
    return Path(__file__).parent.parent


@pytest.fixture
def scripts_dir(project_root: Path) -> Path:
    """Get the scripts directory."""
    return project_root / "scripts"


@pytest.fixture
def data_dir(project_root: Path) -> Path:
    """Get the YAML data directory."""
    return project_root / "data"


@pytest.fixture
def yaml_data_files(data_dir: Path) -> list[Path]:
    """Get all YAML data files (excluding templates and shared-sot.yaml).

    This fixture provides the canonical list of YAML files that should be
    validated and used for documentation generation.
    """
    yaml_files = list(data_dir.rglob("*.yaml"))
    return [
        f for f in yaml_files
        if ".templates" not in str(f) and "shared-sot.yaml" not in str(f)
    ]


@pytest.fixture
def feature_yaml_files(data_dir: Path) -> list[Path]:
    """Get feature YAML files (excluding INDEX files)."""
    files = list(data_dir.glob("features/**/*.yaml"))
    return [f for f in files if "INDEX.yaml" not in f.name]


@pytest.fixture
def service_yaml_files(data_dir: Path) -> list[Path]:
    """Get service YAML files."""
    return list(data_dir.glob("services/*.yaml"))


@pytest.fixture
def integration_yaml_files(data_dir: Path) -> list[Path]:
    """Get integration YAML files (excluding INDEX files)."""
    files = list(data_dir.glob("integrations/**/*.yaml"))
    return [f for f in files if "INDEX.yaml" not in f.name]


@pytest.fixture
def priority_yaml_patterns() -> list[str]:
    """Get glob patterns for priority documentation files.

    Priority 1-2 docs are critical for contributor onboarding and deployment.
    These must have no PLACEHOLDER content.
    """
    return [
        "operations/*.yaml",
        "technical/API.yaml",
        "technical/CONFIGURATION.yaml",
        "technical/FRONTEND.yaml",
    ]
