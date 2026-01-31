"""Tests for dependency automation.

Tests:
- Dependabot configuration validation
- Release Please configuration validation
- Dependency update script functionality
"""

import json
from pathlib import Path
from unittest.mock import Mock, patch

import pytest
import yaml

from scripts.automation.update_dependencies import Dependency, DependencyUpdater


class TestDependabotConfig:
    """Test Dependabot configuration."""

    def test_dependabot_config_exists(self):
        """Test that dependabot.yml exists."""
        config_file = Path.cwd() / ".github" / "dependabot.yml"
        assert config_file.exists(), "dependabot.yml should exist"

    def test_dependabot_config_valid_yaml(self):
        """Test that dependabot.yml is valid YAML."""
        config_file = Path.cwd() / ".github" / "dependabot.yml"
        with open(config_file) as f:
            config = yaml.safe_load(f)
        assert config is not None
        assert "version" in config
        assert config["version"] == 2

    def test_dependabot_has_required_ecosystems(self):
        """Test that dependabot.yml has required ecosystems."""
        config_file = Path.cwd() / ".github" / "dependabot.yml"
        with open(config_file) as f:
            config = yaml.safe_load(f)

        ecosystems = {update["package-ecosystem"] for update in config["updates"]}
        required = {"gomod", "npm", "pip", "github-actions", "docker"}

        assert required.issubset(ecosystems), f"Missing ecosystems: {required - ecosystems}"

    def test_dependabot_weekly_schedule(self):
        """Test that all updates have weekly schedule."""
        config_file = Path.cwd() / ".github" / "dependabot.yml"
        with open(config_file) as f:
            config = yaml.safe_load(f)

        for update in config["updates"]:
            assert "schedule" in update
            assert update["schedule"]["interval"] == "weekly"

    def test_dependabot_has_labels(self):
        """Test that all updates have labels."""
        config_file = Path.cwd() / ".github" / "dependabot.yml"
        with open(config_file) as f:
            config = yaml.safe_load(f)

        for update in config["updates"]:
            assert "labels" in update
            assert "dependencies" in update["labels"]

    def test_dependabot_has_reviewers(self):
        """Test that all updates have reviewers."""
        config_file = Path.cwd() / ".github" / "dependabot.yml"
        with open(config_file) as f:
            config = yaml.safe_load(f)

        for update in config["updates"]:
            assert "reviewers" in update
            assert len(update["reviewers"]) > 0


class TestReleasePleaseConfig:
    """Test Release Please configuration."""

    def test_release_please_config_exists(self):
        """Test that release-please-config.json exists."""
        config_file = Path.cwd() / ".github" / "release-please-config.json"
        assert config_file.exists(), "release-please-config.json should exist"

    def test_release_please_config_valid_json(self):
        """Test that release-please-config.json is valid JSON."""
        config_file = Path.cwd() / ".github" / "release-please-config.json"
        with open(config_file) as f:
            config = json.load(f)
        assert config is not None
        assert "release-type" in config

    def test_release_please_has_packages(self):
        """Test that release-please-config.json has packages."""
        config_file = Path.cwd() / ".github" / "release-please-config.json"
        with open(config_file) as f:
            config = json.load(f)

        assert "packages" in config
        assert "." in config["packages"]

    def test_release_please_has_changelog_sections(self):
        """Test that release-please-config.json has changelog sections."""
        config_file = Path.cwd() / ".github" / "release-please-config.json"
        with open(config_file) as f:
            config = json.load(f)

        pkg_config = config["packages"]["."]
        assert "changelog-sections" in pkg_config

        # Check for required types
        types = {section["type"] for section in pkg_config["changelog-sections"]}
        required = {"feat", "fix", "docs"}
        assert required.issubset(types), f"Missing types: {required - types}"

    def test_release_please_manifest_exists(self):
        """Test that .release-please-manifest.json exists."""
        manifest_file = Path.cwd() / ".release-please-manifest.json"
        assert manifest_file.exists(), ".release-please-manifest.json should exist"

    def test_release_please_manifest_valid_json(self):
        """Test that .release-please-manifest.json is valid JSON."""
        manifest_file = Path.cwd() / ".release-please-manifest.json"
        with open(manifest_file) as f:
            manifest = json.load(f)
        assert manifest is not None
        assert "." in manifest

    def test_release_please_manifest_has_version(self):
        """Test that manifest has valid version."""
        manifest_file = Path.cwd() / ".release-please-manifest.json"
        with open(manifest_file) as f:
            manifest = json.load(f)

        version = manifest["."]
        # Check semver format
        assert version.count(".") == 2, "Version should be semver (x.y.z)"
        parts = version.split(".")
        assert all(p.isdigit() for p in parts), "Version parts should be numeric"


class TestDependencyUpdater:
    """Test DependencyUpdater class."""

    @pytest.fixture
    def updater(self):
        """Create updater instance."""
        return DependencyUpdater(dry_run=True, verbose=False)

    @pytest.fixture
    def mock_subprocess(self):
        """Mock subprocess.run."""
        with patch("scripts.automation.update_dependencies.subprocess.run") as mock:
            mock.return_value = Mock(stdout="", returncode=0)
            yield mock

    def test_updater_initialization(self, updater):
        """Test updater initialization."""
        assert updater.dry_run is True
        assert updater.verbose is False
        assert updater.root == Path.cwd()

    def test_run_command_dry_run(self, updater, capsys):
        """Test run_command in dry-run mode."""
        result = updater.run_command(["echo", "test"])
        assert result == ""

        captured = capsys.readouterr()
        assert "[DRY-RUN]" in captured.out

    def test_check_go_dependencies_empty(self, updater, mock_subprocess):
        """Test checking Go dependencies with no updates."""
        mock_subprocess.return_value.stdout = ""
        deps = updater.check_go_dependencies()
        assert deps == []

    def test_check_go_dependencies_with_updates(self):
        """Test checking Go dependencies with updates available."""
        updater = DependencyUpdater(dry_run=False, verbose=False)

        go_output = json.dumps({
            "Path": "github.com/example/pkg",
            "Version": "v1.0.0",
            "Update": {"Version": "v1.1.0"},
        })

        with patch("scripts.automation.update_dependencies.subprocess.run") as mock:
            mock.return_value = Mock(stdout=go_output, returncode=0)
            deps = updater.check_go_dependencies()

        assert len(deps) == 1
        assert deps[0].name == "github.com/example/pkg"
        assert deps[0].current == "v1.0.0"
        assert deps[0].available == "v1.1.0"
        assert deps[0].ecosystem == "go"

    def test_check_npm_dependencies_no_frontend(self, updater):
        """Test checking npm dependencies with no frontend directory."""
        with patch.object(Path, "exists", return_value=False):
            deps = updater.check_npm_dependencies()
        assert deps == []

    def test_check_npm_dependencies_with_updates(self):
        """Test checking npm dependencies with updates available."""
        updater = DependencyUpdater(dry_run=False, verbose=False)

        npm_output = json.dumps({
            "react": {
                "current": "18.0.0",
                "wanted": "18.2.0",
                "latest": "18.2.0",
            }
        })

        with patch("scripts.automation.update_dependencies.subprocess.run") as mock:
            mock.return_value = Mock(stdout=npm_output, returncode=1)
            with patch.object(Path, "exists", return_value=True):
                deps = updater.check_npm_dependencies()

        assert len(deps) == 1
        assert deps[0].name == "react"
        assert deps[0].current == "18.0.0"
        assert deps[0].available == "18.2.0"
        assert deps[0].ecosystem == "npm"

    def test_check_python_dependencies_empty(self, updater, mock_subprocess):
        """Test checking Python dependencies with no updates."""
        mock_subprocess.return_value.stdout = "[]"
        deps = updater.check_python_dependencies()
        assert deps == []

    def test_check_python_dependencies_with_updates(self):
        """Test checking Python dependencies with updates available."""
        updater = DependencyUpdater(dry_run=False, verbose=False)

        pip_output = json.dumps([
            {
                "name": "pytest",
                "version": "7.0.0",
                "latest_version": "7.4.0",
            }
        ])

        with patch("scripts.automation.update_dependencies.subprocess.run") as mock:
            mock.return_value = Mock(stdout=pip_output, returncode=0)
            deps = updater.check_python_dependencies()

        assert len(deps) == 1
        assert deps[0].name == "pytest"
        assert deps[0].current == "7.0.0"
        assert deps[0].available == "7.4.0"
        assert deps[0].ecosystem == "python"

    def test_check_all_dependencies(self, updater):
        """Test checking all dependencies."""
        with (
            patch.object(updater, "check_go_dependencies", return_value=[]),
            patch.object(updater, "check_npm_dependencies", return_value=[]),
            patch.object(updater, "check_python_dependencies", return_value=[]),
        ):
            deps = updater.check_all_dependencies()

        assert "go" in deps
        assert "npm" in deps
        assert "python" in deps

    def test_check_specific_ecosystem(self, updater):
        """Test checking specific ecosystem only."""
        with (
            patch.object(updater, "check_go_dependencies", return_value=[]) as mock_go,
            patch.object(updater, "check_npm_dependencies") as mock_npm,
            patch.object(updater, "check_python_dependencies") as mock_py,
        ):
            deps = updater.check_all_dependencies(ecosystem="go")

        mock_go.assert_called_once()
        mock_npm.assert_not_called()
        mock_py.assert_not_called()
        assert "go" in deps
        assert "npm" not in deps

    def test_update_go_dependencies_empty(self, updater):
        """Test updating Go dependencies with empty list."""
        result = updater.update_go_dependencies([])
        assert result is True

    def test_update_go_dependencies_dry_run(self, updater, capsys):
        """Test updating Go dependencies in dry-run mode."""
        dep = Dependency(
            name="github.com/example/pkg",
            current="v1.0.0",
            available="v1.1.0",
            ecosystem="go",
        )

        result = updater.update_go_dependencies([dep])
        assert result is True

        captured = capsys.readouterr()
        assert "DRY-RUN" in captured.out
        assert "github.com/example/pkg" in captured.out

    def test_update_npm_dependencies_empty(self, updater):
        """Test updating npm dependencies with empty list."""
        result = updater.update_npm_dependencies([])
        assert result is True

    def test_update_python_dependencies_empty(self, updater):
        """Test updating Python dependencies with empty list."""
        result = updater.update_python_dependencies([])
        assert result is True

    def test_run_tests_dry_run(self, updater, capsys):
        """Test running tests in dry-run mode."""
        result = updater.run_tests()
        assert result is True

        captured = capsys.readouterr()
        assert "DRY-RUN" in captured.out

    def test_create_pr_dry_run(self, updater, capsys):
        """Test creating PR in dry-run mode."""
        deps = {
            "go": [
                Dependency("pkg1", "1.0.0", "1.1.0", "go"),
                Dependency("pkg2", "2.0.0", "2.1.0", "go"),
            ]
        }

        result = updater.create_pr(deps)
        assert result is True

        captured = capsys.readouterr()
        assert "DRY-RUN" in captured.out

    def test_create_pr_no_updates(self, updater):
        """Test creating PR with no updates."""
        result = updater.create_pr({"go": [], "npm": [], "python": []})
        assert result is True


class TestDependencyDataClass:
    """Test Dependency dataclass."""

    def test_dependency_creation(self):
        """Test creating Dependency instance."""
        dep = Dependency(
            name="example",
            current="1.0.0",
            available="2.0.0",
            ecosystem="go",
        )

        assert dep.name == "example"
        assert dep.current == "1.0.0"
        assert dep.available == "2.0.0"
        assert dep.ecosystem == "go"

    def test_dependency_equality(self):
        """Test Dependency equality."""
        dep1 = Dependency("example", "1.0.0", "2.0.0", "go")
        dep2 = Dependency("example", "1.0.0", "2.0.0", "go")

        assert dep1 == dep2

    def test_dependency_repr(self):
        """Test Dependency string representation."""
        dep = Dependency("example", "1.0.0", "2.0.0", "go")
        repr_str = repr(dep)

        assert "example" in repr_str
        assert "1.0.0" in repr_str
        assert "2.0.0" in repr_str
