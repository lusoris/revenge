"""Unit tests for sync-versions.py script."""

from __future__ import annotations

import subprocess
import sys
from pathlib import Path

import pytest

# Add scripts to path
REPO_ROOT = Path(__file__).parent.parent.parent
SYNC_VERSIONS_SCRIPT = REPO_ROOT / "scripts" / "sync-versions.py"


class TestSyncVersionsStrictMode:
    """Tests for --strict mode in sync-versions.py."""

    def test_strict_mode_exits_zero_when_no_drift(self) -> None:
        """Test that --strict mode exits with 0 when no version drift exists."""
        result = subprocess.run(
            [sys.executable, str(SYNC_VERSIONS_SCRIPT), "--strict"],
            capture_output=True,
            text=True,
            cwd=REPO_ROOT,
        )

        # Should exit with 0 when no drift
        assert result.returncode == 0, f"Expected exit 0, got {result.returncode}. Output: {result.stdout}\n{result.stderr}"
        assert "version drift" not in result.stdout or "0 version drift" in result.stdout

    def test_strict_mode_flag_exists(self) -> None:
        """Test that --strict flag is recognized."""
        result = subprocess.run(
            [sys.executable, str(SYNC_VERSIONS_SCRIPT), "--help"],
            capture_output=True,
            text=True,
        )

        assert "--strict" in result.stdout
        assert "Exit with error code if version drift found" in result.stdout or "CI mode" in result.stdout

    def test_script_runs_without_strict(self) -> None:
        """Test that script runs normally without --strict flag."""
        result = subprocess.run(
            [sys.executable, str(SYNC_VERSIONS_SCRIPT)],
            capture_output=True,
            text=True,
            cwd=REPO_ROOT,
        )

        # Should exit with 0 even if drift exists (non-strict mode)
        assert result.returncode == 0
        assert "Extracting versions from SOT" in result.stdout


class TestSyncVersionsBasic:
    """Basic functionality tests for sync-versions.py."""

    def test_script_imports_successfully(self) -> None:
        """Test that the script can be imported without errors."""
        # Just test the subprocess runs
        result = subprocess.run(
            [sys.executable, str(SYNC_VERSIONS_SCRIPT), "--help"],
            capture_output=True,
            text=True,
        )
        assert result.returncode == 0

    def test_help_flag_works(self) -> None:
        """Test that --help flag works."""
        result = subprocess.run(
            [sys.executable, str(SYNC_VERSIONS_SCRIPT), "--help"],
            capture_output=True,
            text=True,
        )

        assert result.returncode == 0
        assert "Sync versions from SOT" in result.stdout
        assert "--fix" in result.stdout
        assert "--report" in result.stdout
        assert "--verbose" in result.stdout

    def test_verbose_flag_recognized(self) -> None:
        """Test that --verbose flag is recognized."""
        result = subprocess.run(
            [sys.executable, str(SYNC_VERSIONS_SCRIPT), "--help"],
            capture_output=True,
            text=True,
        )

        assert "-v" in result.stdout or "--verbose" in result.stdout


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
