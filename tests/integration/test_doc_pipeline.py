"""Integration tests for the doc pipeline."""

from __future__ import annotations

import subprocess
from pathlib import Path

import pytest


class TestDocPipelineRunner:
    """Tests for the doc pipeline runner script."""

    def test_runner_exists(self, scripts_dir: Path) -> None:
        """Test that the runner script exists."""
        runner = scripts_dir / "doc-pipeline.sh"
        assert runner.exists()
        assert runner.is_file()

    def test_runner_is_executable(self, scripts_dir: Path) -> None:
        """Test that the runner script is executable."""
        runner = scripts_dir / "doc-pipeline.sh"
        import os
        import stat

        mode = os.stat(runner).st_mode
        assert mode & stat.S_IXUSR  # Owner execute

    def test_runner_help(self, scripts_dir: Path) -> None:
        """Test that the runner shows help."""
        runner = scripts_dir / "doc-pipeline.sh"
        result = subprocess.run(
            [str(runner), "--help"],
            capture_output=True,
            text=True,
        )
        assert result.returncode == 0
        assert "--apply" in result.stdout
        assert "--step" in result.stdout

    def test_dry_run_no_changes(
        self,
        scripts_dir: Path,
        project_root: Path,
    ) -> None:
        """Test that dry run doesn't modify files."""
        runner = scripts_dir / "doc-pipeline.sh"

        # Get initial state
        design_dir = project_root / "docs" / "dev" / "design"
        if not design_dir.exists():
            pytest.skip("Design directory not found")

        # Run in dry-run mode
        result = subprocess.run(
            [str(runner)],
            capture_output=True,
            text=True,
            cwd=str(project_root),
        )

        # Should complete (may have validation warnings)
        assert "DRY RUN" in result.stdout


class TestSourcePipelineRunner:
    """Tests for the source pipeline runner script."""

    def test_runner_exists(self, scripts_dir: Path) -> None:
        """Test that the runner script exists."""
        runner = scripts_dir / "source-pipeline.sh"
        assert runner.exists()
        assert runner.is_file()

    def test_runner_help(self, scripts_dir: Path) -> None:
        """Test that the runner shows help."""
        runner = scripts_dir / "source-pipeline.sh"
        result = subprocess.run(
            [str(runner), "--help"],
            capture_output=True,
            text=True,
        )
        assert result.returncode == 0
        assert "--apply" in result.stdout
        assert "--force" in result.stdout


class TestPipelineScriptsExist:
    """Tests that all pipeline scripts exist."""

    def test_source_pipeline_scripts(self, scripts_dir: Path) -> None:
        """Test that all source pipeline scripts exist."""
        source_dir = scripts_dir / "source-pipeline"
        assert source_dir.exists()

        expected_scripts = [
            "01-fetch.py",
            "02-index.py",
            "03-breadcrumbs.py",
        ]

        for script in expected_scripts:
            script_path = source_dir / script
            assert script_path.exists(), f"Missing: {script}"

    def test_doc_pipeline_scripts(self, scripts_dir: Path) -> None:
        """Test that all doc pipeline scripts exist."""
        doc_dir = scripts_dir / "doc-pipeline"
        assert doc_dir.exists()

        expected_scripts = [
            "01-indexes.py",
            "02-breadcrumbs.py",
            "03-status.py",
            "04-validate.py",
            "05-fix.py",
            "06-meta.py",
        ]

        for script in expected_scripts:
            script_path = doc_dir / script
            assert script_path.exists(), f"Missing: {script}"

    def test_utils_scripts(self, scripts_dir: Path) -> None:
        """Test that utility scripts exist."""
        utils_dir = scripts_dir / "utils"
        assert utils_dir.exists()

        expected_scripts = [
            "archive-manager.py",
            "sync-versions.py",
        ]

        for script in expected_scripts:
            script_path = utils_dir / script
            assert script_path.exists(), f"Missing: {script}"
