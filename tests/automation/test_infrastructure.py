"""Tests for infrastructure management.

Tests:
- Coder workspace management functionality
- Docker configuration management functionality
- CI/CD workflow management functionality
"""

from pathlib import Path
from unittest.mock import Mock, mock_open, patch

import pytest

from scripts.automation.manage_ci import CIManager
from scripts.automation.manage_coder import CoderManager
from scripts.automation.manage_docker import DockerManager


class TestCoderManager:
    """Test CoderManager class."""

    @pytest.fixture
    def manager(self):
        """Create manager instance."""
        return CoderManager(dry_run=True, verbose=False)

    @pytest.fixture
    def mock_subprocess(self):
        """Mock subprocess.run."""
        with patch("scripts.automation.manage_coder.subprocess.run") as mock:
            mock.return_value = Mock(stdout="", stderr="", returncode=0)
            yield mock

    def test_manager_initialization(self, manager):
        """Test manager initialization."""
        assert manager.dry_run is True
        assert manager.verbose is False
        assert manager.root == Path.cwd()

    def test_check_coder_cli(self, mock_subprocess):
        """Test checking coder CLI."""
        # Use non-dry-run manager for this test
        manager = CoderManager(dry_run=False, verbose=False)

        result = manager.check_coder_cli()
        assert result is True
        mock_subprocess.assert_called()

    def test_list_workspaces(self, manager, capsys):
        """Test listing workspaces."""
        result = manager.list_workspaces()
        assert result is True

        captured = capsys.readouterr()
        assert "Listing Coder workspaces" in captured.out

    def test_create_workspace(self, manager, capsys):
        """Test creating workspace."""
        result = manager.create_workspace("test-workspace")
        assert result is True

        captured = capsys.readouterr()
        assert "Creating workspace" in captured.out

    def test_start_workspace(self, manager, capsys):
        """Test starting workspace."""
        result = manager.start_workspace("test-workspace")
        assert result is True

        captured = capsys.readouterr()
        assert "Starting workspace" in captured.out

    def test_stop_workspace(self, manager, capsys):
        """Test stopping workspace."""
        result = manager.stop_workspace("test-workspace")
        assert result is True

        captured = capsys.readouterr()
        assert "Stopping workspace" in captured.out

    def test_delete_workspace(self, manager, capsys):
        """Test deleting workspace."""
        result = manager.delete_workspace("test-workspace", force=True)
        assert result is True

        captured = capsys.readouterr()
        assert "Deleting workspace" in captured.out

    def test_ssh_workspace(self, manager, capsys):
        """Test SSH into workspace."""
        result = manager.ssh_workspace("test-workspace")
        assert result is True

        captured = capsys.readouterr()
        assert "SSH into workspace" in captured.out

    def test_validate_template_no_file(self, manager, capsys):
        """Test validating template with no file."""
        with patch.object(Path, "exists", return_value=False):
            result = manager.validate_template()

        assert result is False
        captured = capsys.readouterr()
        assert "Template file not found" in captured.out

    def test_update_template_from_sot(self, manager, capsys):
        """Test updating template from SOT."""
        with (
            patch.object(Path, "exists", return_value=True),
        ):
            result = manager.update_template_from_sot()

        # Currently returns True but not implemented
        assert result is True
        captured = capsys.readouterr()
        assert "not yet implemented" in captured.out

    def test_push_template(self, manager, capsys):
        """Test pushing template."""
        result = manager.push_template()
        assert result is True

        captured = capsys.readouterr()
        assert "Pushing template" in captured.out


class TestDockerManager:
    """Test DockerManager class."""

    @pytest.fixture
    def manager(self):
        """Create manager instance."""
        return DockerManager(dry_run=True, verbose=False)

    @pytest.fixture
    def mock_subprocess(self):
        """Mock subprocess.run."""
        with patch("scripts.automation.manage_docker.subprocess.run") as mock:
            mock.return_value = Mock(stdout="", stderr="", returncode=0)
            yield mock

    def test_manager_initialization(self, manager):
        """Test manager initialization."""
        assert manager.dry_run is True
        assert manager.verbose is False
        assert manager.root == Path.cwd()

    def test_check_docker(self, manager, mock_subprocess):
        """Test checking docker."""
        result = manager.check_docker()
        assert result is True

    def test_parse_sot_versions_no_file(self, manager):
        """Test parsing SOT versions with no file."""
        with patch.object(Path, "exists", return_value=False):
            versions = manager.parse_sot_versions()

        assert versions == {}

    def test_parse_sot_versions(self, manager):
        """Test parsing SOT versions."""
        sot_content = """
        ## Infrastructure Components

        | PostgreSQL | 18+ | Database | ✅ |
        | Dragonfly | latest | Cache | ✅ |
        """

        m = mock_open(read_data=sot_content)
        with (
            patch.object(Path, "exists", return_value=True),
            patch("builtins.open", m),
        ):
            versions = manager.parse_sot_versions()

        assert "postgresql" in versions
        assert versions["postgresql"] == "18"

    def test_sync_dockerfile_no_file(self, manager, capsys):
        """Test syncing Dockerfile with no file."""
        with patch.object(Path, "exists", return_value=False):
            result = manager.sync_dockerfile()

        assert result is False
        captured = capsys.readouterr()
        assert "Dockerfile not found" in captured.out

    def test_sync_compose_no_file(self, manager, capsys):
        """Test syncing docker-compose with no file."""
        with patch.object(Path, "exists", return_value=False):
            result = manager.sync_compose()

        assert result is False
        captured = capsys.readouterr()
        assert "docker-compose.yml not found" in captured.out

    def test_validate_configs(self, manager, capsys):
        """Test validating Docker configs."""
        with patch.object(Path, "exists", return_value=False):
            manager.validate_configs()

        # Returns True even if files don't exist (warnings only)
        captured = capsys.readouterr()
        assert "Validating Docker configurations" in captured.out

    def test_build_images(self, manager, capsys):
        """Test building images."""
        result = manager.build_images("latest")
        assert result is True

        captured = capsys.readouterr()
        assert "Building Docker images" in captured.out

    def test_push_images(self, manager, capsys):
        """Test pushing images."""
        result = manager.push_images("ghcr.io/test", "latest")
        assert result is True

        captured = capsys.readouterr()
        assert "Pushing images" in captured.out

    def test_compose_up(self, manager, capsys):
        """Test docker-compose up."""
        result = manager.compose_up(detach=True)
        assert result is True

        captured = capsys.readouterr()
        assert "Starting services" in captured.out

    def test_compose_down(self, manager, capsys):
        """Test docker-compose down."""
        result = manager.compose_down()
        assert result is True

        captured = capsys.readouterr()
        assert "Stopping services" in captured.out


class TestCIManager:
    """Test CIManager class."""

    @pytest.fixture
    def manager(self):
        """Create manager instance."""
        return CIManager(dry_run=True, verbose=False)

    @pytest.fixture
    def mock_subprocess(self):
        """Mock subprocess.run."""
        with patch("scripts.automation.manage_ci.subprocess.run") as mock:
            mock.return_value = Mock(stdout="", stderr="", returncode=0)
            yield mock

    def test_manager_initialization(self, manager):
        """Test manager initialization."""
        assert manager.dry_run is True
        assert manager.verbose is False
        assert manager.root == Path.cwd()
        assert manager.workflows_dir == Path.cwd() / ".github" / "workflows"

    def test_check_gh_cli(self, manager, mock_subprocess):
        """Test checking gh CLI."""
        result = manager.check_gh_cli()
        assert result is True

    def test_list_workflows(self, manager, capsys):
        """Test listing workflows."""
        result = manager.list_workflows()
        assert result is True

        captured = capsys.readouterr()
        assert "Listing GitHub Actions workflows" in captured.out

    def test_list_workflow_files_no_dir(self, manager):
        """Test listing workflow files with no directory."""
        with patch.object(Path, "exists", return_value=False):
            files = manager.list_workflow_files()

        assert files == []

    def test_list_workflow_files(self, manager):
        """Test listing workflow files."""
        test_files = [
            Path(".github/workflows/ci.yml"),
            Path(".github/workflows/release.yml"),
        ]

        with (
            patch.object(Path, "exists", return_value=True),
            patch.object(Path, "glob") as mock_glob,
        ):
            mock_glob.side_effect = [
                [test_files[0]],  # .yml files
                [test_files[1]],  # .yaml files
            ]

            files = manager.list_workflow_files()

        assert len(files) == 2

    def test_validate_workflows_no_files(self, manager, capsys):
        """Test validating workflows with no files."""
        with patch.object(manager, "list_workflow_files", return_value=[]):
            result = manager.validate_workflows()

        assert result is True
        captured = capsys.readouterr()
        assert "No workflow files found" in captured.out

    def test_trigger_workflow(self, manager, capsys):
        """Test triggering workflow."""
        result = manager.trigger_workflow("ci.yml", "main")
        assert result is True

        captured = capsys.readouterr()
        assert "Triggering workflow" in captured.out

    def test_get_workflow_status(self, manager, capsys):
        """Test getting workflow status."""
        result = manager.get_workflow_status(10)
        assert result is True

        captured = capsys.readouterr()
        assert "Recent workflow runs" in captured.out

    def test_watch_workflow(self, manager, capsys):
        """Test watching workflow."""
        result = manager.watch_workflow("123456")
        assert result is True

        captured = capsys.readouterr()
        assert "Watching workflow run" in captured.out

    def test_download_logs(self, manager, capsys):
        """Test downloading logs."""
        with patch.object(Path, "mkdir"):
            result = manager.download_logs("123456")

        assert result is True
        captured = capsys.readouterr()
        assert "Downloading logs" in captured.out

    def test_view_log(self, manager, capsys):
        """Test viewing log."""
        result = manager.view_log("123456")
        assert result is True

        captured = capsys.readouterr()
        assert "Viewing log" in captured.out

    def test_cancel_workflow(self, manager, capsys):
        """Test canceling workflow."""
        result = manager.cancel_workflow("123456")
        assert result is True

        captured = capsys.readouterr()
        assert "Canceling workflow run" in captured.out

    def test_rerun_workflow(self, manager, capsys):
        """Test rerunning workflow."""
        result = manager.rerun_workflow("123456", failed_only=False)
        assert result is True

        captured = capsys.readouterr()
        assert "Rerunning workflow run" in captured.out

    def test_rerun_workflow_failed_only(self, manager, capsys):
        """Test rerunning workflow (failed only)."""
        result = manager.rerun_workflow("123456", failed_only=True)
        assert result is True

        captured = capsys.readouterr()
        assert "Rerunning workflow run" in captured.out


class TestCommandExecution:
    """Test command execution with real subprocess."""

    def test_coder_manager_dry_run_command(self):
        """Test CoderManager dry-run command."""
        manager = CoderManager(dry_run=True, verbose=False)

        with patch("scripts.automation.manage_coder.subprocess.run") as mock:
            result = manager.run_command(["coder", "version"])

        # Should not call subprocess.run in dry-run mode
        mock.assert_not_called()
        assert result.returncode == 0

    def test_docker_manager_dry_run_command(self):
        """Test DockerManager dry-run command."""
        manager = DockerManager(dry_run=True, verbose=False)

        with patch("scripts.automation.manage_docker.subprocess.run") as mock:
            result = manager.run_command(["docker", "version"])

        # Should not call subprocess.run in dry-run mode for docker commands
        mock.assert_not_called()
        assert result.returncode == 0

    def test_ci_manager_dry_run_command(self):
        """Test CIManager dry-run command."""
        manager = CIManager(dry_run=True, verbose=False)

        with patch("scripts.automation.manage_ci.subprocess.run") as mock:
            result = manager.run_command(["gh", "version"])

        # Should not call subprocess.run in dry-run mode for gh commands
        mock.assert_not_called()
        assert result.returncode == 0

    def test_command_not_found(self):
        """Test handling of command not found."""
        manager = CoderManager(dry_run=False, verbose=False)

        with patch("scripts.automation.manage_coder.subprocess.run") as mock:
            mock.side_effect = FileNotFoundError("command not found")

            result = manager.run_command(["nonexistent"])

        assert result.returncode == 127
        assert "not found" in result.stderr
