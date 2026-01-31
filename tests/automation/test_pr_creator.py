#!/usr/bin/env python3
"""Tests for PR creator.

Test coverage for pr_creator.py:
- Loop prevention (bot check, cooldown lock)
- Git operations (branch, commit, push)
- GitHub CLI (PR creation, auto-merge)
- File change detection
- Docs-only change detection
- Full workflow integration
"""

import subprocess
from datetime import datetime, timedelta
from unittest.mock import Mock, patch

import pytest

from scripts.automation.pr_creator import LoopPrevention, PRCreator


class TestLoopPrevention:
    """Test loop prevention system."""

    def test_is_bot_commit_true(self, tmp_path):
        """Test detecting bot commit."""
        loop = LoopPrevention(tmp_path)

        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock(stdout="revenge-bot\n")
            assert loop.is_bot_commit() is True

    def test_is_bot_commit_false(self, tmp_path):
        """Test detecting human commit."""
        loop = LoopPrevention(tmp_path)

        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock(stdout="human-user\n")
            assert loop.is_bot_commit() is False

    def test_is_bot_commit_error(self, tmp_path):
        """Test handling git error."""
        loop = LoopPrevention(tmp_path)

        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = subprocess.CalledProcessError(1, "git")
            assert loop.is_bot_commit() is False

    def test_is_locked_no_file(self, tmp_path):
        """Test no lock file."""
        loop = LoopPrevention(tmp_path)
        is_locked, reason = loop.is_locked()
        assert is_locked is False
        assert reason is None

    def test_is_locked_expired(self, tmp_path):
        """Test expired lock file."""
        loop = LoopPrevention(tmp_path)
        lock_file = tmp_path / ".automation-lock"

        # Create lock file with old timestamp (2 hours ago)
        old_time = datetime.now() - timedelta(hours=2)
        lock_file.write_text(old_time.isoformat())

        is_locked, reason = loop.is_locked()
        assert is_locked is False
        assert reason is None
        assert not lock_file.exists()  # Should be removed

    def test_is_locked_active(self, tmp_path):
        """Test active lock file."""
        loop = LoopPrevention(tmp_path)
        lock_file = tmp_path / ".automation-lock"

        # Create lock file with recent timestamp (30 minutes ago)
        recent_time = datetime.now() - timedelta(minutes=30)
        lock_file.write_text(recent_time.isoformat())

        is_locked, reason = loop.is_locked()
        assert is_locked is True
        assert "minutes remaining" in reason

    def test_is_locked_invalid_file(self, tmp_path):
        """Test invalid lock file."""
        loop = LoopPrevention(tmp_path)
        lock_file = tmp_path / ".automation-lock"

        # Create lock file with invalid content
        lock_file.write_text("invalid-timestamp")

        is_locked, reason = loop.is_locked()
        assert is_locked is False
        assert reason is None
        assert not lock_file.exists()  # Should be removed

    def test_create_lock(self, tmp_path):
        """Test creating lock file."""
        loop = LoopPrevention(tmp_path)
        lock_file = tmp_path / ".automation-lock"

        loop.create_lock()

        assert lock_file.exists()
        content = lock_file.read_text()
        # Should be valid ISO timestamp
        datetime.fromisoformat(content)

    def test_remove_lock(self, tmp_path):
        """Test removing lock file."""
        loop = LoopPrevention(tmp_path)
        lock_file = tmp_path / ".automation-lock"

        # Create lock file
        lock_file.write_text(datetime.now().isoformat())

        loop.remove_lock()

        assert not lock_file.exists()

    def test_can_proceed_bot_commit(self, tmp_path):
        """Test can_proceed when last commit is bot."""
        loop = LoopPrevention(tmp_path)

        with patch.object(loop, "is_bot_commit", return_value=True):
            can_proceed, reason = loop.can_proceed()
            assert can_proceed is False
            assert "revenge-bot" in reason

    def test_can_proceed_locked(self, tmp_path):
        """Test can_proceed when locked."""
        loop = LoopPrevention(tmp_path)

        with patch.object(loop, "is_bot_commit", return_value=False):
            with patch.object(
                loop, "is_locked", return_value=(True, "Cooldown active")
            ):
                can_proceed, reason = loop.can_proceed()
                assert can_proceed is False
                assert "Cooldown active" in reason

    def test_can_proceed_ok(self, tmp_path):
        """Test can_proceed when all checks pass."""
        loop = LoopPrevention(tmp_path)

        with patch.object(loop, "is_bot_commit", return_value=False):
            with patch.object(loop, "is_locked", return_value=(False, None)):
                can_proceed, reason = loop.can_proceed()
                assert can_proceed is True
                assert reason is None


class TestPRCreator:
    """Test PR creator."""

    def test_init(self, tmp_path):
        """Test initialization."""
        creator = PRCreator(tmp_path)
        assert creator.repo_root == tmp_path
        assert creator.loop_prevention is not None

    def test_check_prerequisites_loop_prevention_fail(self, tmp_path):
        """Test prerequisites check when loop prevention fails."""
        creator = PRCreator(tmp_path)

        with patch.object(
            creator.loop_prevention,
            "can_proceed",
            return_value=(False, "Bot commit detected"),
        ):
            success, error = creator.check_prerequisites()
            assert success is False
            assert "Bot commit detected" in error

    def test_check_prerequisites_gh_not_installed(self, tmp_path):
        """Test prerequisites check when gh not installed."""
        creator = PRCreator(tmp_path)

        with patch.object(
            creator.loop_prevention, "can_proceed", return_value=(True, None)
        ), patch("subprocess.run", side_effect=FileNotFoundError):
            success, error = creator.check_prerequisites()
            assert success is False
            assert "gh CLI" in error

    def test_check_prerequisites_gh_not_authenticated(self, tmp_path):
        """Test prerequisites check when gh not authenticated."""
        creator = PRCreator(tmp_path)

        with patch.object(
            creator.loop_prevention, "can_proceed", return_value=(True, None)
        ), patch("subprocess.run") as mock_run:
            # gh --version succeeds, gh auth status fails
            mock_run.side_effect = [
                Mock(),  # gh --version
                subprocess.CalledProcessError(1, "gh"),  # gh auth status
            ]
            success, error = creator.check_prerequisites()
            assert success is False
            assert "authenticated" in error

    def test_check_prerequisites_ok(self, tmp_path):
        """Test prerequisites check when all pass."""
        creator = PRCreator(tmp_path)

        with patch.object(
            creator.loop_prevention, "can_proceed", return_value=(True, None)
        ), patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock()  # All commands succeed
            success, error = creator.check_prerequisites()
            assert success is True
            assert error is None

    def test_get_changed_files(self, tmp_path):
        """Test getting list of changed files."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run") as mock_run:
            # git status --porcelain format: "XY PATH" (2 status chars + space + filename)
            # Example: " M" = modified in worktree, " A" = added to index, " D" = deleted
            mock_run.return_value = Mock(
                stdout=" M file1.py\n A file2.md\n D file3.yaml\n"
            )
            files = creator.get_changed_files()
            assert files == ["file1.py", "file2.md", "file3.yaml"]

    def test_get_changed_files_error(self, tmp_path):
        """Test handling error getting changed files."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run", side_effect=subprocess.CalledProcessError(1, "git")):
            files = creator.get_changed_files()
            assert files == []

    def test_is_docs_only_change_true(self, tmp_path):
        """Test detecting docs-only changes."""
        creator = PRCreator(tmp_path)

        files = [
            "docs/design/FEATURE.md",
            "data/feature.yaml",
            "templates/feature.jinja2",
            "schemas/feature.json",
            ".github/workflows/doc-validation.yml",
        ]
        assert creator.is_docs_only_change(files) is True

    def test_is_docs_only_change_false(self, tmp_path):
        """Test detecting non-docs changes."""
        creator = PRCreator(tmp_path)

        files = [
            "docs/design/FEATURE.md",
            "src/main.py",  # Code file
        ]
        assert creator.is_docs_only_change(files) is False

    def test_create_branch_success(self, tmp_path):
        """Test creating branch."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock()
            assert creator.create_branch("test-branch") is True

    def test_create_branch_failure(self, tmp_path):
        """Test branch creation failure."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run", side_effect=subprocess.CalledProcessError(1, "git")):
            assert creator.create_branch("test-branch") is False

    def test_commit_changes_success(self, tmp_path):
        """Test committing changes."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock()
            assert creator.commit_changes("Test commit") is True
            # Should call git add and git commit
            assert mock_run.call_count == 2

    def test_commit_changes_failure(self, tmp_path):
        """Test commit failure."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run", side_effect=subprocess.CalledProcessError(1, "git")):
            assert creator.commit_changes("Test commit") is False

    def test_push_branch_success(self, tmp_path):
        """Test pushing branch."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock()
            assert creator.push_branch("test-branch") is True

    def test_push_branch_failure(self, tmp_path):
        """Test push failure."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run", side_effect=subprocess.CalledProcessError(1, "git")):
            assert creator.push_branch("test-branch") is False

    def test_create_pr_success(self, tmp_path):
        """Test creating PR."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock(stdout="https://github.com/repo/pull/1\n")
            success, url = creator.create_pr(
                title="Test PR",
                body="Test body",
                branch_name="test-branch",
                auto_merge=False,
            )
            assert success is True
            assert url == "https://github.com/repo/pull/1"

    def test_create_pr_with_automerge(self, tmp_path):
        """Test creating PR with auto-merge."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock(stdout="https://github.com/repo/pull/1\n")
            success, url = creator.create_pr(
                title="Test PR",
                body="Test body",
                branch_name="test-branch",
                auto_merge=True,
            )
            assert success is True
            # Should call gh pr create and gh pr merge
            assert mock_run.call_count == 2

    def test_create_pr_failure(self, tmp_path):
        """Test PR creation failure."""
        creator = PRCreator(tmp_path)

        with patch("subprocess.run", side_effect=subprocess.CalledProcessError(1, "gh")):
            success, error = creator.create_pr(
                title="Test PR",
                body="Test body",
                branch_name="test-branch",
            )
            assert success is False
            assert error is not None


class TestFullWorkflow:
    """Test full PR creation workflow."""

    def test_create_doc_update_pr_prerequisites_fail(self, tmp_path):
        """Test workflow when prerequisites fail."""
        creator = PRCreator(tmp_path)

        with patch.object(
            creator, "check_prerequisites", return_value=(False, "gh not installed")
        ):
            success, error = creator.create_doc_update_pr(
                trigger_type="manual",
                changed_files=["docs/test.md"],
                dry_run=False,
            )
            assert success is False
            assert "gh not installed" in error

    def test_create_doc_update_pr_dry_run(self, tmp_path):
        """Test dry-run mode."""
        creator = PRCreator(tmp_path)

        with patch.object(creator, "check_prerequisites", return_value=(True, None)):
            success, result = creator.create_doc_update_pr(
                trigger_type="manual",
                changed_files=["docs/test.md"],
                dry_run=True,
            )
            assert success is True
            assert result == "dry-run"

    def test_create_doc_update_pr_sot_update(self, tmp_path):
        """Test SOT update PR creation."""
        creator = PRCreator(tmp_path)

        with patch.object(creator, "check_prerequisites", return_value=(True, None)):
            with patch.object(creator, "create_branch", return_value=True):
                with patch.object(creator, "commit_changes", return_value=True):
                    with patch.object(creator, "push_branch", return_value=True):
                        with patch.object(
                            creator,
                            "create_pr",
                            return_value=(
                                True,
                                "https://github.com/repo/pull/1",
                            ),
                        ):
                            with patch.object(creator, "is_docs_only_change", return_value=True):
                                success, url = creator.create_doc_update_pr(
                                    trigger_type="sot_update",
                                    changed_files=["docs/test.md", "data/test.yaml"],
                                    dry_run=False,
                                )
                                assert success is True
                                assert "github.com" in url

                                # Verify lock was created
                                lock_file = tmp_path / ".automation-lock"
                                assert lock_file.exists()

    def test_create_doc_update_pr_branch_creation_fails(self, tmp_path):
        """Test handling branch creation failure."""
        creator = PRCreator(tmp_path)

        with patch.object(creator, "check_prerequisites", return_value=(True, None)):
            with patch.object(creator, "create_branch", return_value=False):
                success, error = creator.create_doc_update_pr(
                    trigger_type="manual",
                    changed_files=["docs/test.md"],
                    dry_run=False,
                )
                assert success is False
                assert "branch" in error.lower()

    def test_create_doc_update_pr_commit_fails(self, tmp_path):
        """Test handling commit failure."""
        creator = PRCreator(tmp_path)

        with patch.object(creator, "check_prerequisites", return_value=(True, None)):
            with patch.object(creator, "create_branch", return_value=True):
                with patch.object(creator, "commit_changes", return_value=False):
                    success, error = creator.create_doc_update_pr(
                        trigger_type="manual",
                        changed_files=["docs/test.md"],
                        dry_run=False,
                    )
                    assert success is False
                    assert "commit" in error.lower()


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
