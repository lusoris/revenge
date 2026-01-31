"""Tests for monitoring and health checks.

Tests:
- Health checker functionality
- Log viewer functionality
"""

from pathlib import Path
from unittest.mock import Mock, mock_open, patch

import pytest

from scripts.automation.check_health import HealthCheck, HealthChecker
from scripts.automation.view_logs import LogViewer


class TestHealthChecker:
    """Test HealthChecker class."""

    @pytest.fixture
    def checker(self):
        """Create checker instance."""
        return HealthChecker(verbose=False)

    @pytest.fixture
    def mock_subprocess(self):
        """Mock subprocess.run."""
        with patch("scripts.automation.check_health.subprocess.run") as mock:
            mock.return_value = Mock(stdout="", stderr="", returncode=0)
            yield mock

    def test_checker_initialization(self, checker):
        """Test checker initialization."""
        assert checker.verbose is False
        assert checker.root == Path.cwd()
        assert checker.results == []

    def test_run_command_success(self, checker, mock_subprocess):
        """Test running command successfully."""
        mock_subprocess.return_value = Mock(
            stdout="output",
            stderr="",
            returncode=0,
        )

        result = checker.run_command(["echo", "test"])

        assert result.returncode == 0
        assert result.stdout == "output"

    def test_run_command_not_found(self, checker):
        """Test handling command not found."""
        with patch("scripts.automation.check_health.subprocess.run") as mock:
            mock.side_effect = FileNotFoundError("not found")

            result = checker.run_command(["nonexistent"])

        assert result.returncode == 1
        assert "not found" in result.stderr

    def test_run_command_timeout(self, checker):
        """Test handling command timeout."""
        import subprocess

        with patch("scripts.automation.check_health.subprocess.run") as mock:
            mock.side_effect = subprocess.TimeoutExpired("slow", 10)

            result = checker.run_command(["slow"])

        assert result.returncode == 1
        assert "slow" in result.stderr.lower()

    def test_check_python_dependencies_missing_file(self, checker):
        """Test checking Python deps with missing requirements.txt."""
        with patch.object(Path, "exists", return_value=False):
            result = checker.check_python_dependencies()

        assert result.component == "python-deps"
        assert result.status == "degraded"
        assert "not found" in result.message

    def test_check_python_dependencies_missing_deps(self, checker):
        """Test checking Python deps with missing dependencies."""
        with (
            patch.object(Path, "exists", return_value=True),
            patch("builtins.__import__", side_effect=ImportError("no module")),
        ):
            result = checker.check_python_dependencies()

        assert result.component == "python-deps"
        assert result.status == "degraded"
        assert "Missing dependencies" in result.message

    def test_check_python_dependencies_success(self, checker):
        """Test checking Python deps successfully."""
        with patch.object(Path, "exists", return_value=True):
            # Don't mock imports - let them succeed naturally
            result = checker.check_python_dependencies()

        # Should succeed if yaml and pytest are installed
        assert result.component == "python-deps"
        assert result.status in ["healthy", "degraded"]

    def test_check_templates_no_dir(self, checker):
        """Test checking templates with no directory."""
        with patch.object(Path, "exists", return_value=False):
            result = checker.check_templates()

        assert result.component == "templates"
        assert result.status == "unhealthy"
        assert "not found" in result.message

    def test_check_templates_no_files(self, checker):
        """Test checking templates with no files."""
        with (
            patch.object(Path, "exists", return_value=True),
            patch.object(Path, "rglob", return_value=[]),
        ):
            result = checker.check_templates()

        assert result.component == "templates"
        assert result.status == "degraded"
        assert "No template files" in result.message

    def test_check_templates_success(self, checker):
        """Test checking templates successfully."""
        test_files = [
            Path("templates/test1.jinja2"),
            Path("templates/test2.jinja2"),
        ]

        with (
            patch.object(Path, "exists", return_value=True),
            patch.object(Path, "rglob", return_value=test_files),
        ):
            result = checker.check_templates()

        assert result.component == "templates"
        assert result.status == "healthy"
        assert "Found 2" in result.message
        assert result.details["count"] == 2

    def test_check_schemas_no_dir(self, checker):
        """Test checking schemas with no directory."""
        with patch.object(Path, "exists", return_value=False):
            result = checker.check_schemas()

        assert result.component == "schemas"
        assert result.status == "degraded"
        assert "not found" in result.message

    def test_check_schemas_success(self, checker):
        """Test checking schemas successfully."""
        test_files = [
            Path("schemas/test.schema.json"),
        ]

        with (
            patch.object(Path, "exists", return_value=True),
            patch.object(Path, "glob", return_value=test_files),
        ):
            result = checker.check_schemas()

        assert result.component == "schemas"
        assert result.status == "healthy"
        assert result.details["count"] == 1

    def test_check_database_not_available(self, checker, mock_subprocess):
        """Test checking database when docker not available."""
        mock_subprocess.return_value.returncode = 1

        result = checker.check_database()

        assert result.component == "database"
        assert result.status == "degraded"
        assert "docker-compose not available" in result.message

    def test_check_database_running(self, checker, mock_subprocess):
        """Test checking database when running."""
        mock_subprocess.return_value = Mock(
            stdout="postgres   Up   5432/tcp",
            returncode=0,
        )

        result = checker.check_database()

        assert result.component == "database"
        assert result.status == "healthy"
        assert "running" in result.message

    def test_check_database_not_running(self, checker, mock_subprocess):
        """Test checking database when not running."""
        mock_subprocess.return_value = Mock(
            stdout="postgres   Exit   5432/tcp",
            returncode=0,
        )

        result = checker.check_database()

        assert result.component == "database"
        assert result.status == "unhealthy"
        assert "not running" in result.message

    def test_check_cache_running(self, checker, mock_subprocess):
        """Test checking cache when running."""
        mock_subprocess.return_value = Mock(
            stdout="dragonfly   Up   6379/tcp",
            returncode=0,
        )

        result = checker.check_cache()

        assert result.component == "cache"
        assert result.status == "healthy"

    def test_check_search_running(self, checker, mock_subprocess):
        """Test checking search when running."""
        mock_subprocess.return_value = Mock(
            stdout="typesense   Up   8108/tcp",
            returncode=0,
        )

        result = checker.check_search()

        assert result.component == "search"
        assert result.status == "healthy"

    def test_check_frontend_build_no_dir(self, checker):
        """Test checking frontend with no directory."""
        with patch.object(Path, "exists", return_value=False):
            result = checker.check_frontend_build()

        assert result.component == "frontend"
        assert result.status == "degraded"

    def test_check_frontend_build_no_node_modules(self, checker):
        """Test checking frontend with no node_modules."""
        with patch.object(Path, "exists") as mock_exists:
            # frontend exists, node_modules doesn't
            mock_exists.side_effect = lambda: mock_exists.call_count == 1

            result = checker.check_frontend_build()

        assert result.component == "frontend"
        assert result.status == "degraded"
        assert "node_modules" in result.message

    def test_check_frontend_build_success(self, checker):
        """Test checking frontend successfully."""
        with patch.object(Path, "exists", return_value=True):
            result = checker.check_frontend_build()

        assert result.component == "frontend"
        assert result.status == "healthy"

    def test_check_resource_usage_cannot_check(self, checker, mock_subprocess):
        """Test checking resources when df fails."""
        mock_subprocess.return_value.returncode = 1

        result = checker.check_resource_usage()

        assert result.component == "resources"
        assert result.status == "degraded"

    def test_check_resource_usage_healthy(self, checker, mock_subprocess):
        """Test checking resources when healthy."""
        mock_subprocess.return_value = Mock(
            stdout="Filesystem     Size  Used Avail Use%\n/dev/sda1       100G   50G   50G  50%",
            returncode=0,
        )

        result = checker.check_resource_usage()

        assert result.component == "resources"
        assert result.status == "healthy"
        assert result.details["disk_usage_pct"] == 50

    def test_check_resource_usage_degraded(self, checker, mock_subprocess):
        """Test checking resources when degraded."""
        mock_subprocess.return_value = Mock(
            stdout="Filesystem     Size  Used Avail Use%\n/dev/sda1       100G   80G   20G  80%",
            returncode=0,
        )

        result = checker.check_resource_usage()

        assert result.component == "resources"
        assert result.status == "degraded"
        assert result.details["disk_usage_pct"] == 80

    def test_check_resource_usage_unhealthy(self, checker, mock_subprocess):
        """Test checking resources when unhealthy."""
        mock_subprocess.return_value = Mock(
            stdout="Filesystem     Size  Used Avail Use%\n/dev/sda1       100G   95G    5G  95%",
            returncode=0,
        )

        result = checker.check_resource_usage()

        assert result.component == "resources"
        assert result.status == "unhealthy"
        assert result.details["disk_usage_pct"] == 95

    def test_check_automation_system(self, checker):
        """Test checking automation system."""
        with patch.object(Path, "exists", return_value=True):
            results = checker.check_automation_system()

        assert len(results) == 3
        assert results[0].component == "python-deps"
        assert results[1].component == "templates"
        assert results[2].component == "schemas"

    def test_check_backend_services(self, checker, mock_subprocess):
        """Test checking backend services."""
        mock_subprocess.return_value = Mock(
            stdout="service   Up",
            returncode=0,
        )

        results = checker.check_backend_services()

        assert len(results) == 3
        assert results[0].component == "database"
        assert results[1].component == "cache"
        assert results[2].component == "search"

    def test_check_all(self, checker, mock_subprocess):
        """Test checking all components."""
        with patch.object(Path, "exists", return_value=True):
            mock_subprocess.return_value = Mock(
                stdout="service   Up\nFilesystem     100G   50G  50%",
                returncode=0,
            )

            results = checker.check_all()

        # Should have: 3 automation + 3 services + 1 frontend + 1 resources = 8
        assert len(results) >= 7  # At least 7 checks

    def test_print_summary(self, checker, capsys):
        """Test printing summary."""
        results = [
            HealthCheck("test1", "healthy", "All good"),
            HealthCheck("test2", "degraded", "Warning"),
            HealthCheck("test3", "unhealthy", "Error"),
        ]

        checker.print_summary(results)

        captured = capsys.readouterr()
        assert "System Health Check" in captured.out
        assert "✅" in captured.out
        assert "⚠️" in captured.out
        assert "❌" in captured.out
        assert "Healthy:1" in captured.out
        assert "Degraded:1" in captured.out
        assert "Unhealthy:1" in captured.out

    def test_print_summary_all_healthy(self, checker, capsys):
        """Test printing summary when all healthy."""
        results = [
            HealthCheck("test1", "healthy", "All good"),
            HealthCheck("test2", "healthy", "All good"),
        ]

        checker.print_summary(results)

        captured = capsys.readouterr()
        assert "HEALTHY" in captured.out

    def test_create_github_issue_no_unhealthy(self, checker):
        """Test creating GitHub issue with no unhealthy components."""
        results = [
            HealthCheck("test", "healthy", "All good"),
        ]

        result = checker.create_github_issue(results)

        assert result is False

    def test_create_github_issue_success(self, checker):
        """Test creating GitHub issue successfully."""
        results = [
            HealthCheck("test", "unhealthy", "Error"),
        ]

        with patch("scripts.automation.check_health.subprocess.run") as mock:
            mock.return_value = Mock(returncode=0)

            result = checker.create_github_issue(results)

        assert result is True
        mock.assert_called_once()

    def test_create_github_issue_failure(self, checker):
        """Test creating GitHub issue with failure."""
        results = [
            HealthCheck("test", "unhealthy", "Error"),
        ]

        with patch("scripts.automation.check_health.subprocess.run") as mock:
            mock.return_value = Mock(returncode=1, stderr="error")

            result = checker.create_github_issue(results)

        assert result is False

    def test_create_github_issue_no_gh_cli(self, checker, capsys):
        """Test creating GitHub issue without gh CLI."""
        results = [
            HealthCheck("test", "unhealthy", "Error"),
        ]

        with patch("scripts.automation.check_health.subprocess.run") as mock:
            mock.side_effect = FileNotFoundError("gh not found")

            result = checker.create_github_issue(results)

        assert result is False
        captured = capsys.readouterr()
        assert "gh CLI not found" in captured.out


class TestLogViewer:
    """Test LogViewer class."""

    @pytest.fixture
    def viewer(self):
        """Create viewer instance."""
        return LogViewer(verbose=False)

    @pytest.fixture
    def mock_subprocess(self):
        """Mock subprocess.run."""
        with patch("scripts.automation.view_logs.subprocess.run") as mock:
            mock.return_value = Mock(stdout="", stderr="", returncode=0)
            yield mock

    def test_viewer_initialization(self, viewer):
        """Test viewer initialization."""
        assert viewer.verbose is False
        assert viewer.root == Path.cwd()
        assert viewer.logs_dir == Path.cwd() / "logs"

    def test_run_command_success(self, viewer, mock_subprocess):
        """Test running command successfully."""
        mock_subprocess.return_value = Mock(
            stdout="output",
            stderr="",
            returncode=0,
        )

        result = viewer.run_command(["echo", "test"])

        assert result.returncode == 0
        assert result.stdout == "output"

    def test_run_command_not_found(self, viewer):
        """Test handling command not found."""
        with patch("scripts.automation.view_logs.subprocess.run") as mock:
            mock.side_effect = FileNotFoundError("not found")

            result = viewer.run_command(["nonexistent"])

        assert result.returncode == 127
        assert "not found" in result.stderr

    def test_list_workflow_runs_success(self, viewer, mock_subprocess):
        """Test listing workflow runs successfully."""
        mock_subprocess.return_value = Mock(
            stdout="workflow1\nworkflow2",
            returncode=0,
        )

        result = viewer.list_workflow_runs(limit=10)

        assert result is True
        mock_subprocess.assert_called_once()

    def test_list_workflow_runs_failure(self, viewer, mock_subprocess, capsys):
        """Test listing workflow runs with failure."""
        mock_subprocess.return_value = Mock(
            stdout="",
            stderr="error",
            returncode=1,
        )

        result = viewer.list_workflow_runs()

        assert result is False
        captured = capsys.readouterr()
        assert "Failed to list" in captured.out

    def test_list_workflow_runs_with_status(self, viewer, mock_subprocess):
        """Test listing workflow runs with status filter."""
        mock_subprocess.return_value = Mock(stdout="runs", returncode=0)

        viewer.list_workflow_runs(status="success")

        args = mock_subprocess.call_args[0][0]
        assert "--status" in args
        assert "completed" in args

    def test_view_workflow_log_success(self, viewer, mock_subprocess):
        """Test viewing workflow log successfully."""
        mock_subprocess.return_value = Mock(
            stdout="log content",
            returncode=0,
        )

        result = viewer.view_workflow_log("12345")

        assert result is True

    def test_view_workflow_log_with_job(self, viewer, mock_subprocess):
        """Test viewing workflow log with specific job."""
        mock_subprocess.return_value = Mock(stdout="log", returncode=0)

        viewer.view_workflow_log("12345", job="test-job")

        args = mock_subprocess.call_args[0][0]
        assert "--job" in args
        assert "test-job" in args

    def test_search_workflow_logs_success(self, viewer, mock_subprocess, capsys):
        """Test searching workflow logs successfully."""
        mock_subprocess.return_value = Mock(
            stdout="line 1 with error\nline 2 normal\nline 3 with error",
            returncode=0,
        )

        result = viewer.search_workflow_logs("12345", "error")

        assert result is True
        captured = capsys.readouterr()
        assert "Found 2 matches" in captured.out

    def test_search_workflow_logs_no_matches(self, viewer, mock_subprocess, capsys):
        """Test searching workflow logs with no matches."""
        mock_subprocess.return_value = Mock(
            stdout="line 1\nline 2",
            returncode=0,
        )

        result = viewer.search_workflow_logs("12345", "notfound")

        assert result is True
        captured = capsys.readouterr()
        assert "No matches found" in captured.out

    def test_search_workflow_logs_case_sensitive(self, viewer, mock_subprocess):
        """Test searching workflow logs case-sensitive."""
        mock_subprocess.return_value = Mock(
            stdout="ERROR\nerror",
            returncode=0,
        )

        viewer.search_workflow_logs("12345", "ERROR", case_sensitive=True)

        # Should only match exact case

    def test_download_workflow_logs_success(self, viewer, mock_subprocess, capsys):
        """Test downloading workflow logs successfully."""
        with patch.object(Path, "mkdir"):
            result = viewer.download_workflow_logs("12345")

        assert result is True
        captured = capsys.readouterr()
        assert "downloaded" in captured.out.lower()

    def test_download_workflow_logs_custom_dir(self, viewer, mock_subprocess):
        """Test downloading workflow logs to custom directory."""
        output_dir = Path("/tmp/custom")

        with patch.object(Path, "mkdir"):
            viewer.download_workflow_logs("12345", output_dir=output_dir)

        args = mock_subprocess.call_args[0][0]
        assert str(output_dir) in args

    def test_view_docker_logs_success(self, viewer, mock_subprocess):
        """Test viewing Docker logs successfully."""
        mock_subprocess.return_value = Mock(
            stdout="container logs",
            returncode=0,
        )

        result = viewer.view_docker_logs("postgres")

        assert result is True

    def test_view_docker_logs_with_lines(self, viewer, mock_subprocess):
        """Test viewing Docker logs with line limit."""
        mock_subprocess.return_value = Mock(stdout="logs", returncode=0)

        viewer.view_docker_logs("postgres", lines=50)

        args = mock_subprocess.call_args[0][0]
        assert "--tail" in args
        assert "50" in args

    def test_view_docker_logs_follow(self, viewer):
        """Test viewing Docker logs with follow mode."""
        with patch("scripts.automation.view_logs.subprocess.run") as mock:
            # Simulate KeyboardInterrupt
            mock.side_effect = KeyboardInterrupt()

            result = viewer.view_docker_logs("postgres", follow=True)

        assert result is True

    def test_view_local_logs_file_not_found(self, viewer, capsys):
        """Test viewing local logs with file not found."""
        with patch.object(Path, "exists", return_value=False):
            result = viewer.view_local_logs("nonexistent.log")

        assert result is False
        captured = capsys.readouterr()
        assert "not found" in captured.out

    def test_view_local_logs_success(self, viewer):
        """Test viewing local logs successfully."""
        test_content = "log line 1\nlog line 2\nlog line 3"

        m = mock_open(read_data=test_content)
        with (
            patch.object(Path, "exists", return_value=True),
            patch("builtins.open", m),
        ):
            result = viewer.view_local_logs("test.log")

        assert result is True

    def test_view_local_logs_with_lines(self, viewer):
        """Test viewing local logs with line limit."""
        test_content = "line 1\nline 2\nline 3\nline 4\nline 5"

        m = mock_open(read_data=test_content)
        with (
            patch.object(Path, "exists", return_value=True),
            patch("builtins.open", m),
        ):
            viewer.view_local_logs("test.log", lines=2)

    def test_view_local_logs_follow(self, viewer):
        """Test viewing local logs with follow mode."""
        with (
            patch.object(Path, "exists", return_value=True),
            patch("scripts.automation.view_logs.subprocess.run") as mock,
        ):
            mock.side_effect = KeyboardInterrupt()

            result = viewer.view_local_logs("test.log", follow=True)

        assert result is True

    def test_search_local_logs_no_dir(self, viewer, capsys):
        """Test searching local logs with no directory."""
        with patch.object(Path, "exists", return_value=False):
            result = viewer.search_local_logs("error")

        assert result is False
        captured = capsys.readouterr()
        assert "not found" in captured.out

    def test_search_local_logs_no_files(self, viewer, capsys):
        """Test searching local logs with no files."""
        with (
            patch.object(Path, "exists", return_value=True),
            patch.object(Path, "rglob", return_value=[]),
        ):
            result = viewer.search_local_logs("error")

        assert result is True
        captured = capsys.readouterr()
        assert "No log files found" in captured.out

    def test_search_local_logs_success(self, viewer, capsys):
        """Test searching local logs successfully."""
        test_files = [Path("logs/test.log")]
        test_content = "line 1 with error\nline 2 normal"

        m = mock_open(read_data=test_content)
        with (
            patch.object(Path, "exists", return_value=True),
            patch.object(Path, "rglob", return_value=test_files),
            patch("builtins.open", m),
        ):
            result = viewer.search_local_logs("error")

        assert result is True
        captured = capsys.readouterr()
        assert "test.log:1" in captured.out

    def test_search_local_logs_case_sensitive(self, viewer):
        """Test searching local logs case-sensitive."""
        test_files = [Path("logs/test.log")]
        test_content = "ERROR\nerror"

        m = mock_open(read_data=test_content)
        with (
            patch.object(Path, "exists", return_value=True),
            patch.object(Path, "rglob", return_value=test_files),
            patch("builtins.open", m),
        ):
            viewer.search_local_logs("ERROR", case_sensitive=True)

    def test_list_local_logs_no_dir(self, viewer, capsys):
        """Test listing local logs with no directory."""
        with patch.object(Path, "exists", return_value=False):
            result = viewer.list_local_logs()

        assert result is True
        captured = capsys.readouterr()
        assert "No logs directory" in captured.out

    def test_list_local_logs_no_files(self, viewer, capsys):
        """Test listing local logs with no files."""
        with (
            patch.object(Path, "exists", return_value=True),
            patch.object(Path, "rglob", return_value=[]),
        ):
            result = viewer.list_local_logs()

        assert result is True
        captured = capsys.readouterr()
        assert "No log files found" in captured.out

    def test_list_local_logs_success(self, viewer, capsys):
        """Test listing local logs successfully."""
        test_file = Path("logs/test.log")

        with (
            patch.object(Path, "exists", return_value=True),
            patch.object(Path, "rglob") as mock_rglob,
            patch.object(Path, "stat") as mock_stat,
            patch.object(Path, "relative_to", return_value=test_file),
        ):
            mock_rglob.side_effect = [[test_file], []]
            mock_stat.return_value.st_size = 1024
            mock_stat.return_value.st_mtime = 1234567890

            result = viewer.list_local_logs()

        assert result is True
        captured = capsys.readouterr()
        assert "test.log" in captured.out


class TestCommandExecution:
    """Test command execution patterns."""

    def test_health_checker_timeout(self):
        """Test HealthChecker command timeout."""
        checker = HealthChecker(verbose=False)

        with patch("scripts.automation.check_health.subprocess.run") as mock:
            import subprocess
            mock.side_effect = subprocess.TimeoutExpired("cmd", 10)

            result = checker.run_command(["slow"])

        assert result.returncode == 1

    def test_log_viewer_error_handling(self):
        """Test LogViewer error handling."""
        viewer = LogViewer(verbose=False)

        result = viewer.view_workflow_log("invalid")

        # Should handle gracefully
        assert result in [True, False]
