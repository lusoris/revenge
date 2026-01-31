#!/usr/bin/env python3
"""Tests for GitHub automation scripts.

Test coverage:
- github_projects.py: Project setup, field creation, workflow configuration
- github_discussions.py: Discussions setup, template creation, category listing

Tests use mocking to avoid actual GitHub API calls.
"""

import json
import os
import subprocess
from pathlib import Path
from unittest.mock import Mock, patch

import pytest

from scripts.automation.github_discussions import GitHubDiscussionsManager
from scripts.automation.github_labels import GitHubLabelsManager
from scripts.automation.github_milestones import GitHubMilestonesManager
from scripts.automation.github_projects import GitHubProjectsManager
from scripts.automation.github_security import GitHubSecurityManager


@pytest.fixture
def mock_subprocess():
    """Mock subprocess.run for gh CLI commands."""
    with patch("subprocess.run") as mock_run:
        # Default successful response
        mock_result = Mock()
        mock_result.stdout = ""
        mock_result.returncode = 0
        mock_run.return_value = mock_result
        yield mock_run


@pytest.fixture
def temp_git_repo(tmp_path):
    """Create temporary git repository structure."""
    # Create .git directory
    (tmp_path / ".git").mkdir()

    # Create .github directory
    (tmp_path / ".github").mkdir()

    return tmp_path


class TestGitHubProjectsManager:
    """Test GitHubProjectsManager class."""

    def test_init(self):
        """Test manager initialization."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=True)

        assert manager.repo_owner == "owner"
        assert manager.repo_name == "repo"
        assert manager.dry_run is True
        assert manager.repo_full == "owner/repo"

    def test_run_gh_command_dry_run(self, capsys):
        """Test gh command in dry-run mode."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=True)

        result = manager.run_gh_command(["project", "list"])

        captured = capsys.readouterr()
        assert "[DRY-RUN]" in captured.out
        assert "gh project list" in captured.out
        assert result == ""

    def test_run_gh_command_success(self, mock_subprocess):
        """Test successful gh command execution."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=False)

        mock_subprocess.return_value.stdout = "project output"

        result = manager.run_gh_command(["project", "list"])

        mock_subprocess.assert_called_once()
        call_args = mock_subprocess.call_args[0][0]
        assert call_args[0] == "gh"
        assert call_args[1] == "project"
        assert call_args[2] == "list"
        assert result == "project output"

    def test_run_gh_command_failure(self, mock_subprocess):
        """Test gh command failure."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=False)

        mock_subprocess.side_effect = subprocess.CalledProcessError(1, "gh")

        with pytest.raises(subprocess.CalledProcessError):
            manager.run_gh_command(["project", "create"])

    def test_create_project_dry_run(self, capsys):
        """Test project creation in dry-run mode."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=True)

        project_number = manager.create_project("Test Project", "Test description")

        assert project_number == "1"  # Dummy number in dry-run
        captured = capsys.readouterr()
        assert "Creating project" in captured.out

    def test_create_project_success(self, mock_subprocess):
        """Test successful project creation."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=False)

        # Mock project creation response
        project_data = {"number": 42, "title": "Test Project"}
        mock_subprocess.return_value.stdout = json.dumps(project_data)

        project_number = manager.create_project("Test Project", "Description")

        assert project_number == "42"
        mock_subprocess.assert_called_once()

    def test_add_field_text(self, mock_subprocess):
        """Test adding text field to project."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=False)

        manager.add_field("1", "Description", "TEXT")

        mock_subprocess.assert_called_once()
        call_args = mock_subprocess.call_args[0][0]
        assert "field-create" in call_args
        assert "Description" in call_args
        assert "TEXT" in call_args

    def test_add_field_single_select(self, mock_subprocess):
        """Test adding single select field with options."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=False)

        options = ["High", "Medium", "Low"]
        manager.add_field("1", "Priority", "SINGLE_SELECT", options)

        mock_subprocess.assert_called_once()
        call_args = mock_subprocess.call_args[0][0]
        assert "field-create" in call_args
        assert "Priority" in call_args
        assert "SINGLE_SELECT" in call_args
        # Check that options are included
        assert "--single-select-option" in call_args

    def test_configure_workflow(self, capsys):
        """Test workflow configuration."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=True)

        manager.configure_workflow("1")

        captured = capsys.readouterr()
        assert "Configuring workflow" in captured.out
        assert "Backlog" in captured.out
        assert "In Progress" in captured.out
        assert "Done" in captured.out

    def test_add_items_to_project_no_issues(self, mock_subprocess):
        """Test adding items when no issues exist."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=False)

        mock_subprocess.return_value.stdout = "[]"

        manager.add_items_to_project("1")

        # Should call gh issue list
        assert any(
            "issue" in str(call_args) for call_args in mock_subprocess.call_args_list
        )

    def test_add_items_to_project_with_issues(self, mock_subprocess):
        """Test adding existing issues to project."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=False)

        # Mock issue list response
        issues = [{"number": 1}, {"number": 2}]
        mock_subprocess.return_value.stdout = json.dumps(issues)

        # Mock successful item-add for both issues
        def run_side_effect(cmd, **kwargs):
            result = Mock()
            if "issue" in cmd and "list" in cmd:
                result.stdout = json.dumps(issues)
            else:
                result.stdout = ""
            result.returncode = 0
            return result

        mock_subprocess.side_effect = run_side_effect

        manager.add_items_to_project("1")

        # Should call issue list + 2 item-add calls
        assert mock_subprocess.call_count >= 1

    def test_setup_project_integration(self, mock_subprocess, capsys):
        """Test complete project setup workflow."""
        manager = GitHubProjectsManager("owner", "repo", dry_run=True)

        manager.setup_project()

        captured = capsys.readouterr()
        assert "GitHub Projects Setup" in captured.out
        assert "DRY-RUN MODE" in captured.out
        assert "Project setup complete" in captured.out


class TestGitHubDiscussionsManager:
    """Test GitHubDiscussionsManager class."""

    def test_init(self):
        """Test manager initialization."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=True)

        assert manager.repo_owner == "owner"
        assert manager.repo_name == "repo"
        assert manager.dry_run is True
        assert manager.repo_full == "owner/repo"

    def test_run_gh_command_dry_run(self, capsys):
        """Test gh command in dry-run mode."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=True)

        result = manager.run_gh_command(["api", "/repos/owner/repo"])

        captured = capsys.readouterr()
        assert "[DRY-RUN]" in captured.out
        assert "gh api /repos/owner/repo" in captured.out
        assert result == ""

    def test_enable_discussions_already_enabled(self, mock_subprocess):
        """Test checking discussions when already enabled."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=False)

        mock_subprocess.return_value.stdout = "true"

        result = manager.enable_discussions()

        assert result is True
        mock_subprocess.assert_called_once()

    def test_enable_discussions_not_enabled(self, mock_subprocess):
        """Test checking discussions when not enabled."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=False)

        mock_subprocess.return_value.stdout = "false"

        result = manager.enable_discussions()

        assert result is False

    def test_enable_discussions_api_error(self, mock_subprocess):
        """Test handling API error when checking discussions."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=False)

        mock_subprocess.side_effect = subprocess.CalledProcessError(1, "gh")

        result = manager.enable_discussions()

        assert result is False

    def test_list_categories_success(self, mock_subprocess):
        """Test listing discussion categories."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=False)

        # Mock GraphQL response
        graphql_response = {
            "data": {
                "repository": {
                    "discussionCategories": {
                        "nodes": [
                            {"id": "1", "name": "Ideas", "emoji": "💡", "description": "Feature ideas"},
                            {"id": "2", "name": "Q&A", "emoji": "❓", "description": "Questions"},
                        ]
                    }
                }
            }
        }
        mock_subprocess.return_value.stdout = json.dumps(graphql_response)

        categories = manager.list_categories()

        assert len(categories) == 2
        assert categories[0]["name"] == "Ideas"
        assert categories[1]["name"] == "Q&A"

    def test_list_categories_empty(self, mock_subprocess):
        """Test listing categories when none exist."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=False)

        graphql_response = {
            "data": {
                "repository": {
                    "discussionCategories": {
                        "nodes": []
                    }
                }
            }
        }
        mock_subprocess.return_value.stdout = json.dumps(graphql_response)

        categories = manager.list_categories()

        assert len(categories) == 0

    def test_list_categories_error(self, mock_subprocess):
        """Test handling error when listing categories."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=False)

        mock_subprocess.side_effect = subprocess.CalledProcessError(1, "gh")

        categories = manager.list_categories()

        assert categories == []

    def test_create_templates_dry_run(self, temp_git_repo, capsys):
        """Test template creation in dry-run mode."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=True)

        # Change to temp directory
        original_cwd = Path.cwd()
        os.chdir(temp_git_repo)

        try:
            manager.create_templates()

            captured = capsys.readouterr()
            assert "Creating discussion templates" in captured.out
            assert "[DRY-RUN]" in captured.out
            assert "ideas.yml" in captured.out
        finally:
            os.chdir(original_cwd)

    def test_create_templates_success(self, temp_git_repo):
        """Test successful template creation."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=False)

        # Change to temp directory
        original_cwd = Path.cwd()
        os.chdir(temp_git_repo)

        try:
            manager.create_templates()

            # Check that templates were created
            template_dir = temp_git_repo / ".github" / "DISCUSSION_TEMPLATE"
            assert template_dir.exists()
            assert (template_dir / "ideas.yml").exists()
            assert (template_dir / "question.yml").exists()
            assert (template_dir / "show-and-tell.yml").exists()
            assert (template_dir / "bug-report.yml").exists()

            # Check template content
            ideas_content = (template_dir / "ideas.yml").read_text()
            assert "Feature Idea" in ideas_content
            assert "Summary" in ideas_content
        finally:
            os.chdir(original_cwd)

    def test_create_templates_already_exists(self, temp_git_repo, capsys):
        """Test template creation when files already exist."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=False)

        original_cwd = Path.cwd()
        os.chdir(temp_git_repo)

        try:
            # Create template dir with existing file
            template_dir = temp_git_repo / ".github" / "DISCUSSION_TEMPLATE"
            template_dir.mkdir(parents=True, exist_ok=True)
            (template_dir / "ideas.yml").write_text("existing content")

            manager.create_templates()

            captured = capsys.readouterr()
            assert "already exists" in captured.out
        finally:
            os.chdir(original_cwd)

    def test_print_manual_steps(self, capsys):
        """Test printing manual configuration steps."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=True)

        manager.print_manual_steps()

        captured = capsys.readouterr()
        assert "Manual Configuration Required" in captured.out
        assert "Enable Discussions" in captured.out
        assert "Configure Discussion Categories" in captured.out

    def test_setup_discussions_integration(self, temp_git_repo, mock_subprocess, capsys):
        """Test complete discussions setup workflow."""
        manager = GitHubDiscussionsManager("owner", "repo", dry_run=True)

        original_cwd = Path.cwd()
        os.chdir(temp_git_repo)

        try:
            mock_subprocess.return_value.stdout = "false"

            manager.setup_discussions()

            captured = capsys.readouterr()
            assert "GitHub Discussions Setup" in captured.out
            assert "DRY-RUN MODE" in captured.out
            assert "Discussion templates created" in captured.out
        finally:
            os.chdir(original_cwd)


class TestHelperFunctions:
    """Test helper functions."""

    def test_get_repo_info_from_https_url(self):
        """Test parsing repo info from HTTPS URL."""
        from scripts.automation.github_projects import get_repo_info

        with patch("subprocess.run") as mock_run:
            mock_run.return_value.stdout = "https://github.com/owner/repo.git\n"
            mock_run.return_value.returncode = 0

            owner, repo = get_repo_info()

            assert owner == "owner"
            assert repo == "repo"

    def test_get_repo_info_from_ssh_url(self):
        """Test parsing repo info from SSH URL."""
        from scripts.automation.github_projects import get_repo_info

        with patch("subprocess.run") as mock_run:
            mock_run.return_value.stdout = "git@github.com:owner/repo.git\n"
            mock_run.return_value.returncode = 0

            owner, repo = get_repo_info()

            assert owner == "owner"
            assert repo == "repo"

    def test_get_repo_info_no_git_remote(self):
        """Test error when no git remote exists."""
        from scripts.automation.github_projects import get_repo_info

        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = subprocess.CalledProcessError(1, "git")

            with pytest.raises(SystemExit):
                get_repo_info()

    def test_get_repo_info_invalid_url(self):
        """Test error with invalid remote URL."""
        from scripts.automation.github_projects import get_repo_info

        with patch("subprocess.run") as mock_run:
            mock_run.return_value.stdout = "https://gitlab.com/owner/repo.git\n"
            mock_run.return_value.returncode = 0

            with pytest.raises(SystemExit):
                get_repo_info()


class TestGitHubSecurityManager:
    """Test GitHubSecurityManager class."""

    def test_init(self):
        """Test manager initialization."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=True)

        assert manager.repo_owner == "owner"
        assert manager.repo_name == "repo"
        assert manager.dry_run is True
        assert manager.repo_full == "owner/repo"

    def test_run_gh_command_dry_run(self, capsys):
        """Test gh command in dry-run mode."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=True)

        result = manager.run_gh_command(["api", "/repos/owner/repo"])

        captured = capsys.readouterr()
        assert "[DRY-RUN]" in captured.out
        assert "gh api /repos/owner/repo" in captured.out
        assert result == ""

    def test_run_gh_command_success(self, mock_subprocess):
        """Test successful gh command execution."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=False)

        mock_subprocess.return_value.stdout = "command output"

        result = manager.run_gh_command(["api", "/test"])

        mock_subprocess.assert_called_once()
        assert result == "command output"

    def test_configure_branch_protection_dry_run(self, capsys):
        """Test branch protection configuration in dry-run mode."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=True)

        manager.configure_branch_protection("main", strict=True)

        captured = capsys.readouterr()
        assert "Configuring protection for branch: main" in captured.out
        assert "[DRY-RUN]" in captured.out
        assert "Require PR reviews" in captured.out
        assert "Linear history" in captured.out

    def test_configure_branch_protection_success(self, mock_subprocess):
        """Test successful branch protection configuration."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=False)

        mock_subprocess.return_value.stdout = ""

        manager.configure_branch_protection("develop", strict=False)

        # Should call gh api to set branch protection
        mock_subprocess.assert_called_once()
        call_args = mock_subprocess.call_args[0][0]
        assert "api" in call_args
        assert "PUT" in call_args
        assert "branches/develop/protection" in " ".join(call_args)

    def test_configure_branch_protection_failure(self, mock_subprocess, capsys):
        """Test handling branch protection configuration failure."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=False)

        mock_subprocess.side_effect = subprocess.CalledProcessError(1, "gh")

        manager.configure_branch_protection("main")

        captured = capsys.readouterr()
        assert "Failed to configure protection" in captured.out

    def test_enable_security_features_dry_run(self, capsys):
        """Test enabling security features in dry-run mode."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=True)

        manager.enable_security_features()

        captured = capsys.readouterr()
        assert "Enabling security features" in captured.out
        assert "[DRY-RUN]" in captured.out
        assert "secret scanning" in captured.out.lower()

    def test_enable_security_features_success(self, mock_subprocess):
        """Test successful security feature enabling."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=False)

        mock_subprocess.return_value.stdout = ""

        manager.enable_security_features()

        # Should make API calls to enable features
        assert mock_subprocess.call_count >= 2  # At least 2 features

    def test_enable_security_features_partial_failure(self, mock_subprocess, capsys):
        """Test handling partial failures when enabling features."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=False)

        # First call succeeds, second fails
        mock_subprocess.side_effect = [
            Mock(stdout="", returncode=0),
            subprocess.CalledProcessError(1, "gh"),
        ]

        manager.enable_security_features()

        captured = capsys.readouterr()
        assert "Failed to enable" in captured.out

    def test_check_security_status_dry_run(self, capsys):
        """Test checking security status in dry-run mode."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=True)

        manager.check_security_status()

        captured = capsys.readouterr()
        assert "Checking security status" in captured.out

    def test_check_security_status_success(self, mock_subprocess):
        """Test successful security status check."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=False)

        security_status = {
            "secret_scanning": {"status": "enabled"},
            "secret_scanning_push_protection": {"status": "enabled"},
            "dependabot_security_updates": {"status": "disabled"},
        }
        mock_subprocess.return_value.stdout = json.dumps(security_status)

        manager.check_security_status()

        mock_subprocess.assert_called_once()

    def test_check_security_status_api_error(self, mock_subprocess, capsys):
        """Test handling API error when checking status."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=False)

        mock_subprocess.side_effect = subprocess.CalledProcessError(1, "gh")

        manager.check_security_status()

        captured = capsys.readouterr()
        assert "Could not check security status" in captured.out

    def test_check_security_status_invalid_json(self, mock_subprocess, capsys):
        """Test handling invalid JSON response."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=False)

        mock_subprocess.return_value.stdout = "invalid json"

        manager.check_security_status()

        captured = capsys.readouterr()
        assert "Could not check security status" in captured.out

    def test_print_manual_steps(self, capsys):
        """Test printing manual configuration steps."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=True)

        manager.print_manual_steps()

        captured = capsys.readouterr()
        assert "Manual Configuration Steps" in captured.out
        assert "CodeQL Analysis" in captured.out
        assert "Dependabot" in captured.out
        assert "Branch Protection" in captured.out
        assert "CODEOWNERS" in captured.out

    def test_setup_security_integration(self, mock_subprocess, capsys):
        """Test complete security setup workflow."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=True)

        manager.setup_security()

        captured = capsys.readouterr()
        assert "GitHub Security Setup" in captured.out
        assert "DRY-RUN MODE" in captured.out
        assert "Security setup complete" in captured.out

    def test_setup_security_live_mode(self, mock_subprocess):
        """Test security setup in live mode."""
        manager = GitHubSecurityManager("owner", "repo", dry_run=False)

        # Mock all API responses
        mock_subprocess.return_value.stdout = "{}"

        manager.setup_security()

        # Should make multiple API calls
        assert mock_subprocess.call_count >= 3  # status check + branches + features


class TestGitHubLabelsManager:
    """Test GitHubLabelsManager class."""

    def test_init(self):
        """Test manager initialization."""
        manager = GitHubLabelsManager("owner", "repo", dry_run=True)

        assert manager.repo_owner == "owner"
        assert manager.repo_name == "repo"
        assert manager.dry_run is True
        assert manager.repo_full == "owner/repo"

    def test_load_label_config(self, tmp_path):
        """Test loading label configuration from YAML."""
        manager = GitHubLabelsManager("owner", "repo")

        config_file = tmp_path / "labels.yml"
        config_file.write_text("""
- name: "bug"
  color: "d73a4a"
  description: "Something isn't working"
- name: "feature"
  color: "a2eeef"
  description: "New feature"
""")

        labels = manager.load_label_config(config_file)

        assert len(labels) == 2
        assert labels[0]["name"] == "bug"
        assert labels[0]["color"] == "d73a4a"
        assert labels[1]["name"] == "feature"

    def test_get_existing_labels(self, mock_subprocess):
        """Test getting existing labels from GitHub."""
        manager = GitHubLabelsManager("owner", "repo", dry_run=False)

        labels_data = [
            {"name": "bug", "color": "d73a4a", "description": "Bug"},
            {"name": "feature", "color": "a2eeef", "description": "Feature"},
        ]
        mock_subprocess.return_value.stdout = json.dumps(labels_data)

        labels = manager.get_existing_labels()

        assert len(labels) == 2
        assert "bug" in labels
        assert labels["bug"]["color"] == "d73a4a"

    def test_create_label(self, mock_subprocess):
        """Test creating a new label."""
        manager = GitHubLabelsManager("owner", "repo", dry_run=False)

        manager.create_label("test", "ffffff", "Test label")

        mock_subprocess.assert_called_once()
        call_args = mock_subprocess.call_args[0][0]
        assert "api" in call_args
        assert "POST" in call_args

    def test_update_label(self, mock_subprocess):
        """Test updating an existing label."""
        manager = GitHubLabelsManager("owner", "repo", dry_run=False)

        manager.update_label("test", "000000", "Updated description")

        mock_subprocess.assert_called_once()
        call_args = mock_subprocess.call_args[0][0]
        assert "api" in call_args
        assert "PATCH" in call_args

    def test_delete_label(self, mock_subprocess):
        """Test deleting a label."""
        manager = GitHubLabelsManager("owner", "repo", dry_run=False)

        manager.delete_label("test")

        mock_subprocess.assert_called_once()
        call_args = mock_subprocess.call_args[0][0]
        assert "api" in call_args
        assert "DELETE" in call_args

    def test_sync_labels_dry_run(self, tmp_path, mock_subprocess, capsys):
        """Test label sync in dry-run mode."""
        manager = GitHubLabelsManager("owner", "repo", dry_run=True)

        config_file = tmp_path / "labels.yml"
        config_file.write_text("""
- name: "new-label"
  color: "ffffff"
  description: "New label"
""")

        manager.sync_labels(config_file)

        captured = capsys.readouterr()
        assert "DRY-RUN MODE" in captured.out
        assert "new-label" in captured.out

    def test_check_labels(self, tmp_path, mock_subprocess, capsys):
        """Test checking labels without making changes."""
        manager = GitHubLabelsManager("owner", "repo", dry_run=False)

        config_file = tmp_path / "labels.yml"
        config_file.write_text("""
- name: "bug"
  color: "d73a4a"
  description: "Bug"
""")

        mock_subprocess.return_value.stdout = json.dumps([
            {"name": "feature", "color": "a2eeef", "description": "Feature"}
        ])

        manager.check_labels(config_file)

        captured = capsys.readouterr()
        assert "Label Status" in captured.out
        assert "Missing from GitHub" in captured.out


class TestGitHubMilestonesManager:
    """Test GitHubMilestonesManager class."""

    def test_init(self):
        """Test manager initialization."""
        manager = GitHubMilestonesManager("owner", "repo", dry_run=True)

        assert manager.repo_owner == "owner"
        assert manager.repo_name == "repo"
        assert manager.dry_run is True
        assert manager.repo_full == "owner/repo"

    def test_get_milestones(self, mock_subprocess):
        """Test getting milestones from repository."""
        manager = GitHubMilestonesManager("owner", "repo", dry_run=False)

        milestones_data = [
            {"title": "v1.0.0", "number": 1, "state": "open"},
            {"title": "v1.1.0", "number": 2, "state": "open"},
        ]
        mock_subprocess.return_value.stdout = json.dumps(milestones_data)

        milestones = manager.get_milestones()

        assert len(milestones) == 2
        assert milestones[0]["title"] == "v1.0.0"

    def test_get_latest_milestone_version(self, mock_subprocess):
        """Test getting latest milestone version."""
        manager = GitHubMilestonesManager("owner", "repo", dry_run=False)

        milestones_data = [
            {"title": "v1.0.0", "number": 1},
            {"title": "v1.2.0", "number": 2},
            {"title": "v1.1.0", "number": 3},
        ]
        mock_subprocess.return_value.stdout = json.dumps(milestones_data)

        latest = manager.get_latest_milestone_version()

        assert latest == "v1.2.0"

    def test_increment_version_minor(self):
        """Test incrementing version (minor)."""
        manager = GitHubMilestonesManager("owner", "repo")

        new_version = manager.increment_version("v1.2.3", "minor")

        assert new_version == "v1.3.0"

    def test_increment_version_major(self):
        """Test incrementing version (major)."""
        manager = GitHubMilestonesManager("owner", "repo")

        new_version = manager.increment_version("v1.2.3", "major")

        assert new_version == "v2.0.0"

    def test_increment_version_patch(self):
        """Test incrementing version (patch)."""
        manager = GitHubMilestonesManager("owner", "repo")

        new_version = manager.increment_version("v1.2.3", "patch")

        assert new_version == "v1.2.4"

    def test_create_milestone(self, mock_subprocess):
        """Test creating a milestone."""
        manager = GitHubMilestonesManager("owner", "repo", dry_run=False)

        manager.create_milestone("v1.0.0", "Release 1.0.0", "2026-12-31")

        mock_subprocess.assert_called_once()
        call_args = mock_subprocess.call_args[0][0]
        assert "api" in call_args
        assert "POST" in call_args

    def test_close_milestone(self, mock_subprocess):
        """Test closing a milestone."""
        manager = GitHubMilestonesManager("owner", "repo", dry_run=False)

        manager.close_milestone(1)

        mock_subprocess.assert_called_once()
        call_args = mock_subprocess.call_args[0][0]
        assert "api" in call_args
        assert "PATCH" in call_args

    def test_create_next_milestone_dry_run(self, mock_subprocess, capsys):
        """Test creating next milestone in dry-run mode."""
        # Create manager in non-dry-run mode first to get milestones
        temp_manager = GitHubMilestonesManager("owner", "repo", dry_run=False)

        mock_subprocess.return_value.stdout = json.dumps([
            {"title": "v1.0.0", "number": 1}
        ])

        # Now get latest version
        latest = temp_manager.get_latest_milestone_version()

        # Create dry-run manager
        manager = GitHubMilestonesManager("owner", "repo", dry_run=True)

        # Patch the get_latest_milestone_version to return the version we got
        with patch.object(manager, 'get_latest_milestone_version', return_value=latest):
            manager.create_next_milestone()

        captured = capsys.readouterr()
        assert "DRY-RUN MODE" in captured.out
        assert ("v1.1.0" in captured.out or "v1.0.0" in captured.out)

    def test_close_completed_milestones(self, mock_subprocess, capsys):
        """Test closing completed milestones."""
        manager = GitHubMilestonesManager("owner", "repo", dry_run=False)

        milestones_data = [
            {"title": "v1.0.0", "number": 1, "open_issues": 0, "closed_issues": 5},
            {"title": "v1.1.0", "number": 2, "open_issues": 3, "closed_issues": 2},
        ]
        mock_subprocess.return_value.stdout = json.dumps(milestones_data)

        manager.close_completed_milestones()

        captured = capsys.readouterr()
        assert "Closed 1 milestones" in captured.out


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
