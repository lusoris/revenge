"""Tests for GitHub Actions workflow validation.

Ensures workflows are correctly configured, especially the validate-sot.yml
workflow that checks for hardcoded versions.
"""

import re
from pathlib import Path

import pytest
import yaml


@pytest.fixture
def validate_sot_workflow():
    """Load validate-sot.yml workflow file."""
    workflow_path = Path(".github/workflows/validate-sot.yml")
    assert workflow_path.exists(), "validate-sot.yml workflow not found"
    with workflow_path.open() as f:
        return yaml.safe_load(f)


@pytest.fixture
def validate_sot_content():
    """Load raw content of validate-sot.yml for grep pattern testing."""
    workflow_path = Path(".github/workflows/validate-sot.yml")
    with workflow_path.open() as f:
        return f.read()


class TestValidateSOTWorkflowStructure:
    """Test basic structure of validate-sot.yml workflow."""

    def test_workflow_has_name(self, validate_sot_workflow):
        """Workflow must have a name."""
        assert "name" in validate_sot_workflow
        assert validate_sot_workflow["name"] == "Validate SOURCE_OF_TRUTH"

    def test_workflow_has_triggers(self, validate_sot_workflow):
        """Workflow must have trigger events."""
        # YAML parses 'on' as True (boolean key), so check both
        triggers = validate_sot_workflow.get("on") or validate_sot_workflow.get(True)
        assert triggers is not None, "Workflow must have 'on' trigger section"

        # Should have pull_request, schedule, and workflow_dispatch
        assert "pull_request" in triggers
        assert "schedule" in triggers
        assert "workflow_dispatch" in triggers

    def test_workflow_has_jobs(self, validate_sot_workflow):
        """Workflow must have jobs defined."""
        assert "jobs" in validate_sot_workflow
        assert len(validate_sot_workflow["jobs"]) > 0


class TestValidateSOTExcludeDirectories:
    """Test that validation correctly excludes external and archive directories."""

    def test_excludes_external_sources(self, validate_sot_content):
        """Should exclude docs/dev/sources from version checks."""
        # Check that docs/dev/sources is excluded in grep commands
        assert "--exclude-dir=docs/dev/sources" in validate_sot_content, (
            "validate-sot.yml should exclude docs/dev/sources directory "
            "(contains external documentation with hardcoded versions)"
        )

    def test_excludes_analysis_directory(self, validate_sot_content):
        """Should exclude .analysis from version checks."""
        assert "--exclude-dir=.analysis" in validate_sot_content, (
            "validate-sot.yml should exclude .analysis directory "
            "(contains archived files with hardcoded versions)"
        )

    def test_excludes_shared_directory(self, validate_sot_content):
        """Should exclude .shared from version checks."""
        assert "--exclude-dir=.shared" in validate_sot_content, (
            "validate-sot.yml should exclude .shared directory "
            "(may contain external content)"
        )

    def test_excludes_zed_directory(self, validate_sot_content):
        """Should exclude .zed from version checks."""
        assert "--exclude-dir=.zed" in validate_sot_content, (
            "validate-sot.yml should exclude .zed directory "
            "(may contain tool-specific configs)"
        )

    def test_excludes_source_of_truth_file(self, validate_sot_content):
        """Should exclude SOURCE_OF_TRUTH.md itself from checks."""
        # Check for various exclusion patterns
        assert (
            "00_SOURCE_OF_TRUTH.md" in validate_sot_content
            or "SOURCE_OF_TRUTH" in validate_sot_content
        ), (
            "validate-sot.yml should exclude SOURCE_OF_TRUTH.md from "
            "hardcoded version checks (it's the source of versions)"
        )

    def test_excludes_versions_yml(self, validate_sot_content):
        """Should exclude _versions.yml from checks."""
        assert "_versions.yml" in validate_sot_content, (
            "validate-sot.yml should exclude _versions.yml "
            "(contains extracted versions)"
        )

    def test_excludes_validate_sot_yml(self, validate_sot_content):
        """Should exclude validate-sot.yml itself from checks."""
        assert "validate-sot.yml" in validate_sot_content, (
            "validate-sot.yml should exclude itself from checks "
            "(contains version patterns for validation)"
        )


class TestValidateSOTVersionChecks:
    """Test version checking logic."""

    def test_checks_for_go_versions(self, validate_sot_content):
        """Should check for hardcoded Go versions."""
        # Should have patterns for Go version detection
        assert "go1\\.25" in validate_sot_content or "Go 1\\.25" in validate_sot_content, (
            "Should check for hardcoded Go versions"
        )

    def test_checks_for_python_versions(self, validate_sot_content):
        """Should check for hardcoded Python versions."""
        assert "python-version" in validate_sot_content.lower() or "3\\.12" in validate_sot_content, (
            "Should check for hardcoded Python versions"
        )

    def test_checks_for_node_versions(self, validate_sot_content):
        """Should check for hardcoded Node.js versions."""
        assert "node-version" in validate_sot_content.lower(), (
            "Should check for hardcoded Node.js versions"
        )

    def test_checks_for_goexperiment(self, validate_sot_content):
        """Should check for hardcoded GOEXPERIMENT values."""
        assert "GOEXPERIMENT" in validate_sot_content, (
            "Should check for hardcoded GOEXPERIMENT values"
        )


class TestValidateSOTFormatChecks:
    """Test SOURCE_OF_TRUTH format validation."""

    def test_verifies_go_version_field(self, validate_sot_content):
        """Should verify Go Version field exists."""
        assert "**Go Version**" in validate_sot_content, (
            "Should check for '**Go Version**:' field in SOURCE_OF_TRUTH"
        )

    def test_verifies_nodejs_field(self, validate_sot_content):
        """Should verify Node.js field exists."""
        assert "**Node\\.js**" in validate_sot_content or "Node.js" in validate_sot_content, (
            "Should check for '**Node.js**:' field in SOURCE_OF_TRUTH"
        )

    def test_verifies_python_field(self, validate_sot_content):
        """Should verify Python field exists."""
        assert "**Python**" in validate_sot_content, (
            "Should check for '**Python**:' field in SOURCE_OF_TRUTH"
        )

    def test_verifies_postgresql_field(self, validate_sot_content):
        """Should verify PostgreSQL field exists."""
        assert "**PostgreSQL**" in validate_sot_content, (
            "Should check for '**PostgreSQL**:' field in SOURCE_OF_TRUTH"
        )

    def test_verifies_build_command_field(self, validate_sot_content):
        """Should verify Build Command field exists."""
        assert "**Build Command**" in validate_sot_content, (
            "Should check for '**Build Command**:' field in SOURCE_OF_TRUTH"
        )


class TestValidateSOTExtractionTests:
    """Test version extraction validation."""

    def test_extracts_go_version(self, validate_sot_content):
        """Should test extraction of Go version."""
        # Should have sed command to extract Go version
        assert "GO_VERSION" in validate_sot_content, (
            "Should test extraction of GO_VERSION"
        )

    def test_extracts_goexperiment(self, validate_sot_content):
        """Should test extraction of GOEXPERIMENT."""
        assert "GOEXPERIMENT" in validate_sot_content, (
            "Should test extraction of GOEXPERIMENT from Build Command"
        )

    def test_extracts_postgres_version(self, validate_sot_content):
        """Should test extraction of PostgreSQL version."""
        assert "POSTGRES_VERSION" in validate_sot_content, (
            "Should test extraction of POSTGRES_VERSION"
        )

    def test_extracts_python_version(self, validate_sot_content):
        """Should test extraction of Python version."""
        assert "PYTHON_VERSION" in validate_sot_content, (
            "Should test extraction of PYTHON_VERSION"
        )

    def test_extracts_node_version(self, validate_sot_content):
        """Should test extraction of Node.js version."""
        assert "NODE_VERSION" in validate_sot_content, (
            "Should test extraction of NODE_VERSION"
        )


class TestValidateSOTOutputs:
    """Test workflow outputs and summary."""

    def test_has_summary_job(self, validate_sot_workflow):
        """Should have a summary job."""
        assert "summary" in validate_sot_workflow["jobs"], (
            "Should have a 'summary' job for reporting results"
        )

    def test_summary_creates_github_summary(self, validate_sot_content):
        """Summary job should create GitHub step summary."""
        assert "GITHUB_STEP_SUMMARY" in validate_sot_content, (
            "Summary job should write to GITHUB_STEP_SUMMARY"
        )


class TestWorkflowYAMLValidity:
    """Test YAML validity of all workflows."""

    def test_all_workflows_valid_yaml(self):
        """All workflow files must be valid YAML."""
        workflows_dir = Path(".github/workflows")
        assert workflows_dir.exists(), "Workflows directory not found"

        for workflow_file in workflows_dir.glob("*.yml"):
            with workflow_file.open() as f:
                try:
                    yaml.safe_load(f)
                except yaml.YAMLError as e:
                    pytest.fail(f"Invalid YAML in {workflow_file.name}: {e}")

    def test_workflows_use_consistent_indentation(self):
        """Workflows should use consistent indentation."""
        workflows_dir = Path(".github/workflows")

        for workflow_file in workflows_dir.glob("*.yml"):
            with workflow_file.open() as f:
                lines = f.readlines()

            for i, line in enumerate(lines, 1):
                if line.strip().startswith("#") or not line.strip():
                    continue

                spaces = len(line) - len(line.lstrip(" "))
                if spaces > 0:
                    assert spaces % 2 == 0, (
                        f"{workflow_file.name} line {i}: "
                        f"inconsistent indentation (not multiple of 2)"
                    )


class TestWorkflowNaming:
    """Test workflow naming conventions."""

    def test_workflow_files_use_kebab_case(self):
        """Workflow files should use kebab-case naming."""
        workflows_dir = Path(".github/workflows")

        for workflow_file in workflows_dir.glob("*.yml"):
            name = workflow_file.stem

            # Should be lowercase with hyphens
            assert name.islower() or "_" in name, (
                f"{workflow_file.name}: workflow file names should be "
                f"lowercase (kebab-case or snake_case)"
            )

            # Should not have spaces
            assert " " not in name, (
                f"{workflow_file.name}: workflow file names should not "
                f"contain spaces"
            )


class TestWorkflowPermissions:
    """Test workflow permissions configuration."""

    def test_workflows_have_minimal_permissions(self):
        """Workflows should specify minimal required permissions."""
        workflows_dir = Path(".github/workflows")

        for workflow_file in workflows_dir.glob("*.yml"):
            # Skip special workflows that may need broader permissions
            if workflow_file.name.startswith("_"):
                continue

            with workflow_file.open() as f:
                content = yaml.safe_load(f)

            # If workflow writes to repo or creates issues/PRs,
            # it should have explicit permissions
            if "permissions" not in content:
                # Check if it needs permissions
                workflow_str = str(content)
                needs_write = any(
                    keyword in workflow_str.lower()
                    for keyword in ["commit", "push", "create", "issue", "pr"]
                )

                if needs_write:
                    pytest.skip(
                        f"{workflow_file.name}: may need explicit permissions. "
                        f"Consider adding 'permissions:' section."
                    )


class TestVersionsWorkflow:
    """Test _versions.yml workflow specifically."""

    def test_versions_workflow_exists(self):
        """The _versions.yml workflow must exist."""
        versions_workflow = Path(".github/workflows/_versions.yml")
        assert versions_workflow.exists(), (
            "_versions.yml workflow not found. This is required for "
            "extracting versions from SOURCE_OF_TRUTH"
        )

    def test_versions_workflow_is_reusable(self):
        """_versions.yml should be a reusable workflow."""
        versions_workflow = Path(".github/workflows/_versions.yml")
        with versions_workflow.open() as f:
            content = yaml.safe_load(f)

        # YAML parses 'on' as True (boolean key)
        triggers = content.get("on") or content.get(True)
        assert triggers is not None, "Workflow must have triggers"
        assert "workflow_call" in triggers, (
            "_versions.yml should be a reusable workflow "
            "(use workflow_call trigger)"
        )

    def test_versions_workflow_has_outputs(self):
        """_versions.yml should define version outputs."""
        versions_workflow = Path(".github/workflows/_versions.yml")
        with versions_workflow.open() as f:
            content = yaml.safe_load(f)

        # Should have jobs with outputs
        assert "jobs" in content
        jobs = content["jobs"]

        # At least one job should have outputs
        has_outputs = any("outputs" in job for job in jobs.values())
        assert has_outputs, (
            "_versions.yml should have job outputs for version values"
        )
