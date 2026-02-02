"""Unit tests for sync-sot-status.py script."""

from __future__ import annotations

import subprocess
import sys
from pathlib import Path

import pytest

REPO_ROOT = Path(__file__).parent.parent.parent
SCRIPT_PATH = REPO_ROOT / "scripts" / "doc-pipeline" / "04-sync-sot-status.py"
SOT_FILE = REPO_ROOT / "docs" / "dev" / "design" / "00_SOURCE_OF_TRUTH.md"


def test_script_exists():
    """Test that the sync-sot-status.py script exists."""
    assert SCRIPT_PATH.exists(), f"Script not found at {SCRIPT_PATH}"


def test_script_runs_without_errors():
    """Test that script runs in dry-run mode without errors."""
    result = subprocess.run(
        [sys.executable, str(SCRIPT_PATH)],
        cwd=str(REPO_ROOT),
        capture_output=True,
        text=True,
    )

    assert result.returncode == 0, f"Script failed: {result.stderr}"
    assert "STATUS SYNC SUMMARY" in result.stdout


def test_script_has_apply_flag():
    """Test that script accepts --apply flag."""
    result = subprocess.run(
        [sys.executable, str(SCRIPT_PATH), "--help"],
        cwd=str(REPO_ROOT),
        capture_output=True,
        text=True,
    )

    assert "--apply" in result.stdout, "--apply flag not found in help"


def test_script_has_strict_flag():
    """Test that script accepts --strict flag."""
    result = subprocess.run(
        [sys.executable, str(SCRIPT_PATH), "--help"],
        cwd=str(REPO_ROOT),
        capture_output=True,
        text=True,
    )

    assert "--strict" in result.stdout, "--strict flag not found in help"


def test_script_has_verbose_flag():
    """Test that script accepts --verbose flag."""
    result = subprocess.run(
        [sys.executable, str(SCRIPT_PATH), "--help"],
        cwd=str(REPO_ROOT),
        capture_output=True,
        text=True,
    )

    assert "--verbose" in result.stdout or "-v" in result.stdout


def test_dry_run_does_not_modify_sot():
    """Test that dry-run mode doesn't modify SOURCE_OF_TRUTH.md."""
    # Get original modification time
    original_mtime = SOT_FILE.stat().st_mtime

    # Run in dry-run mode
    result = subprocess.run(
        [sys.executable, str(SCRIPT_PATH)],
        cwd=str(REPO_ROOT),
        capture_output=True,
        text=True,
    )

    assert result.returncode == 0

    # Check file wasn't modified
    new_mtime = SOT_FILE.stat().st_mtime
    assert original_mtime == new_mtime, "SOT file was modified in dry-run mode"


def test_script_detects_yaml_statuses():
    """Test that script finds and parses YAML files."""
    result = subprocess.run(
        [sys.executable, str(SCRIPT_PATH), "--verbose"],
        cwd=str(REPO_ROOT),
        capture_output=True,
        text=True,
    )

    assert result.returncode == 0
    assert "YAML files scanned:" in result.stdout

    # Extract number of YAML files found
    for line in result.stdout.split("\n"):
        if "YAML files scanned:" in line:
            num_files = int(line.split(":")[-1].strip())
            assert num_files > 0, "No YAML files found"
            break


def test_script_processes_all_tables():
    """Test that script processes all expected tables."""
    result = subprocess.run(
        [sys.executable, str(SCRIPT_PATH), "--verbose"],
        cwd=str(REPO_ROOT),
        capture_output=True,
        text=True,
    )

    assert result.returncode == 0

    expected_tables = [
        "Content Modules",
        "Backend Services",
        "Metadata Providers",
        "Arr Ecosystem",
    ]

    for table in expected_tables:
        assert f"Processing: {table}" in result.stdout, f"Table '{table}' not processed"


def test_verbose_mode_shows_changes():
    """Test that verbose mode displays detected changes."""
    result = subprocess.run(
        [sys.executable, str(SCRIPT_PATH), "--verbose"],
        cwd=str(REPO_ROOT),
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )

    assert result.returncode == 0

    # If there are changes, verbose should show them with arrow or hyphen
    if "Status changes detected: " in result.stdout:
        for line in result.stdout.split("\n"):
            if "Status changes detected:" in line:
                num_changes = int(line.split(":")[-1].strip())
                if num_changes > 0:
                    # Check for change indicator (arrow or hyphen as fallback)
                    has_indicator = ("→" in result.stdout or
                                    "->" in result.stdout or
                                    " - " in result.stdout)
                    assert has_indicator, "Verbose mode should show changes"
                break


def test_strict_mode_functionality():
    """Test that strict mode exits with code 1 when drift detected."""
    # First check if there's any drift
    check_result = subprocess.run(
        [sys.executable, str(SCRIPT_PATH)],
        cwd=str(REPO_ROOT),
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )

    # Extract number of changes
    num_changes = 0
    for line in check_result.stdout.split("\n"):
        if "Status changes detected:" in line:
            num_changes = int(line.split(":")[-1].strip())
            break

    # Run with --strict
    strict_result = subprocess.run(
        [sys.executable, str(SCRIPT_PATH), "--strict"],
        cwd=str(REPO_ROOT),
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )

    if num_changes > 0:
        assert strict_result.returncode == 1, "Strict mode should exit with code 1 when drift detected"
        assert strict_result.stdout and "strict" in strict_result.stdout.lower()
    else:
        assert strict_result.returncode == 0, "Strict mode should exit with code 0 when no drift"


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
