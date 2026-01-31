"""Tests for Dependabot configuration validation.

Ensures Dependabot config follows best practices and prevents issues like
duplicate scopes in commit messages.
"""

import re
from pathlib import Path

import pytest
import yaml


@pytest.fixture
def dependabot_config():
    """Load Dependabot configuration file."""
    config_path = Path(".github/dependabot.yml")
    assert config_path.exists(), "Dependabot config not found"
    with config_path.open() as f:
        return yaml.safe_load(f)


class TestDependabotConfigStructure:
    """Test basic structure of Dependabot configuration."""

    def test_version_is_2(self, dependabot_config):
        """Dependabot config must use version 2."""
        assert dependabot_config["version"] == 2

    def test_has_updates_section(self, dependabot_config):
        """Config must have updates section."""
        assert "updates" in dependabot_config
        assert isinstance(dependabot_config["updates"], list)
        assert len(dependabot_config["updates"]) > 0

    def test_all_updates_have_required_fields(self, dependabot_config):
        """Each update must have required fields."""
        required_fields = {"package-ecosystem", "directory", "schedule"}

        for update in dependabot_config["updates"]:
            for field in required_fields:
                assert field in update, (
                    f"Missing required field '{field}' in {update.get('package-ecosystem', 'unknown')} update"
                )


class TestDependabotCommitMessages:
    """Test commit message configuration to prevent duplicate scopes."""

    def test_no_duplicate_scope_in_prefix(self, dependabot_config):
        """Prevent duplicate scopes like 'chore(deps)(deps)'.

        When using 'include: scope', the prefix should NOT already contain
        a scope suffix like '(deps)'. It should be just the type.
        """
        for update in dependabot_config["updates"]:
            if "commit-message" not in update:
                continue

            commit_msg = update["commit-message"]
            ecosystem = update["package-ecosystem"]

            # If include: scope is used, prefix should not end with (...)
            if commit_msg.get("include") == "scope":
                prefix = commit_msg.get("prefix", "")

                # Prefix should not contain parentheses
                assert "(" not in prefix, (
                    f"{ecosystem}: prefix '{prefix}' contains '(' which will cause "
                    f"duplicate scope when used with 'include: scope'. "
                    f"Use just the type (e.g., 'chore' not 'chore(deps)')"
                )

                assert ")" not in prefix, (
                    f"{ecosystem}: prefix '{prefix}' contains ')' which will cause "
                    f"duplicate scope when used with 'include: scope'. "
                    f"Use just the type (e.g., 'chore' not 'chore(deps)')"
                )

    def test_conventional_commit_format(self, dependabot_config):
        """Commit message prefixes should follow conventional commits."""
        valid_types = {
            "feat",
            "fix",
            "docs",
            "style",
            "refactor",
            "perf",
            "test",
            "build",
            "ci",
            "chore",
            "revert",
        }

        for update in dependabot_config["updates"]:
            if "commit-message" not in update:
                continue

            prefix = update["commit-message"].get("prefix", "")
            ecosystem = update["package-ecosystem"]

            # Extract the type (part before any scope)
            commit_type = prefix.split("(")[0].strip()

            assert commit_type in valid_types, (
                f"{ecosystem}: commit type '{commit_type}' is not a valid "
                f"conventional commit type. Valid types: {valid_types}"
            )

    def test_consistent_commit_message_config(self, dependabot_config):
        """All ecosystems should have consistent commit message config."""
        ecosystems_with_commit_msg = [
            update
            for update in dependabot_config["updates"]
            if "commit-message" in update
        ]

        if not ecosystems_with_commit_msg:
            pytest.skip("No commit-message configurations found")

        # All should use include: scope
        for update in ecosystems_with_commit_msg:
            ecosystem = update["package-ecosystem"]
            commit_msg = update["commit-message"]

            assert "include" in commit_msg, (
                f"{ecosystem}: should specify 'include' field"
            )
            assert commit_msg["include"] == "scope", (
                f"{ecosystem}: should use 'include: scope' for consistency"
            )


class TestDependabotSchedules:
    """Test schedule configuration."""

    def test_all_have_schedule(self, dependabot_config):
        """Each update must have a schedule."""
        for update in dependabot_config["updates"]:
            assert "schedule" in update
            assert "interval" in update["schedule"]

    def test_schedule_intervals_are_valid(self, dependabot_config):
        """Schedule intervals must be valid values."""
        valid_intervals = {"daily", "weekly", "monthly"}

        for update in dependabot_config["updates"]:
            interval = update["schedule"]["interval"]
            ecosystem = update["package-ecosystem"]

            assert interval in valid_intervals, (
                f"{ecosystem}: interval '{interval}' is not valid. "
                f"Valid values: {valid_intervals}"
            )

    def test_weekly_schedules_have_day(self, dependabot_config):
        """Weekly schedules must specify a day."""
        for update in dependabot_config["updates"]:
            schedule = update["schedule"]
            ecosystem = update["package-ecosystem"]

            if schedule["interval"] == "weekly":
                assert "day" in schedule, (
                    f"{ecosystem}: weekly schedule must specify a day"
                )

                valid_days = {
                    "monday",
                    "tuesday",
                    "wednesday",
                    "thursday",
                    "friday",
                    "saturday",
                    "sunday",
                }
                assert schedule["day"] in valid_days, (
                    f"{ecosystem}: day '{schedule['day']}' is not valid"
                )


class TestDependabotEcosystems:
    """Test package ecosystem configuration."""

    def test_critical_ecosystems_present(self, dependabot_config):
        """Ensure critical ecosystems are configured."""
        ecosystems = {
            update["package-ecosystem"]
            for update in dependabot_config["updates"]
        }

        critical = {"gomod", "npm", "github-actions"}
        missing = critical - ecosystems

        assert not missing, (
            f"Missing critical ecosystems: {missing}. "
            f"Found: {ecosystems}"
        )

    def test_ecosystem_directories_exist(self, dependabot_config):
        """Directories specified in config should exist (or be planned)."""
        for update in dependabot_config["updates"]:
            directory = update["directory"]
            ecosystem = update["package-ecosystem"]

            # Root directory always exists
            if directory == "/":
                continue

            # Remove leading slash for Path
            dir_path = Path(directory.lstrip("/"))

            # Skip check for npm frontend dir if not yet created
            if ecosystem == "npm" and directory == "/frontend" and not dir_path.exists():
                pytest.skip(
                    f"{ecosystem}: directory '{directory}' not yet created "
                    f"(frontend app to be implemented)"
                )
                continue

            assert dir_path.exists(), (
                f"{ecosystem}: directory '{directory}' does not exist"
            )


class TestDependabotLabels:
    """Test label configuration."""

    def test_all_have_labels(self, dependabot_config):
        """Each update should have labels configured."""
        for update in dependabot_config["updates"]:
            ecosystem = update["package-ecosystem"]
            assert "labels" in update, (
                f"{ecosystem}: should have labels configured"
            )
            assert isinstance(update["labels"], list)
            assert len(update["labels"]) > 0

    def test_dependencies_label_present(self, dependabot_config):
        """All updates should have 'dependencies' label."""
        for update in dependabot_config["updates"]:
            ecosystem = update["package-ecosystem"]
            labels = update.get("labels", [])

            assert "dependencies" in labels, (
                f"{ecosystem}: should have 'dependencies' label"
            )


class TestDependabotGroups:
    """Test dependency grouping configuration."""

    def test_gomod_has_grouping(self, dependabot_config):
        """Go modules should use dependency grouping."""
        gomod_updates = [
            update
            for update in dependabot_config["updates"]
            if update["package-ecosystem"] == "gomod"
        ]

        assert len(gomod_updates) > 0, "No gomod configuration found"

        for update in gomod_updates:
            assert "groups" in update, (
                "gomod should use dependency grouping for minor/patch updates"
            )

    def test_npm_has_grouping(self, dependabot_config):
        """npm should use dependency grouping."""
        npm_updates = [
            update
            for update in dependabot_config["updates"]
            if update["package-ecosystem"] == "npm"
        ]

        if npm_updates:
            for update in npm_updates:
                assert "groups" in update, (
                    "npm should use dependency grouping for minor/patch updates"
                )


class TestDependabotReviewers:
    """Test reviewer configuration."""

    def test_all_have_reviewers(self, dependabot_config):
        """Each update should have reviewers configured."""
        for update in dependabot_config["updates"]:
            ecosystem = update["package-ecosystem"]
            assert "reviewers" in update, (
                f"{ecosystem}: should have reviewers configured"
            )
            assert isinstance(update["reviewers"], list)
            assert len(update["reviewers"]) > 0


class TestDependabotPRLimits:
    """Test pull request limit configuration."""

    def test_all_have_pr_limits(self, dependabot_config):
        """Each update should have open-pull-requests-limit configured."""
        for update in dependabot_config["updates"]:
            ecosystem = update["package-ecosystem"]
            assert "open-pull-requests-limit" in update, (
                f"{ecosystem}: should have open-pull-requests-limit configured"
            )

    def test_pr_limits_are_reasonable(self, dependabot_config):
        """PR limits should be reasonable (not too high)."""
        for update in dependabot_config["updates"]:
            ecosystem = update["package-ecosystem"]
            limit = update["open-pull-requests-limit"]

            assert isinstance(limit, int), (
                f"{ecosystem}: PR limit must be an integer"
            )
            assert 1 <= limit <= 20, (
                f"{ecosystem}: PR limit {limit} is unreasonable. "
                f"Should be between 1 and 20"
            )


class TestDependabotYAMLValidity:
    """Test YAML file validity."""

    def test_yaml_is_valid(self):
        """Dependabot YAML must be valid."""
        config_path = Path(".github/dependabot.yml")

        with config_path.open() as f:
            try:
                yaml.safe_load(f)
            except yaml.YAMLError as e:
                pytest.fail(f"Invalid YAML: {e}")

    def test_no_tabs_in_yaml(self):
        """YAML should not contain tabs (use spaces)."""
        config_path = Path(".github/dependabot.yml")

        with config_path.open() as f:
            content = f.read()

        assert "\t" not in content, "YAML file contains tabs, should use spaces"

    def test_consistent_indentation(self):
        """YAML should use consistent 2-space indentation."""
        config_path = Path(".github/dependabot.yml")

        with config_path.open() as f:
            lines = f.readlines()

        for i, line in enumerate(lines, 1):
            if line.strip().startswith("#") or not line.strip():
                continue

            # Get leading spaces
            spaces = len(line) - len(line.lstrip(" "))

            # Should be multiple of 2
            if spaces > 0:
                assert spaces % 2 == 0, (
                    f"Line {i}: inconsistent indentation (not multiple of 2)"
                )
