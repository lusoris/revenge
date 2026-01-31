#!/usr/bin/env python3
"""Tests for config synchronization.

Test coverage for config_sync.py:
- SOT data loading
- Version extraction from SOT
- .tool-versions synchronization
- go.mod synchronization
- GitHub workflows synchronization
- Dockerfile synchronization
- Dry-run vs live mode
- Error handling
- Full sync workflow
"""


import pytest
import yaml

from scripts.automation.config_sync import ConfigSync


@pytest.fixture
def sync_setup(tmp_path):
    """Set up config sync test environment."""
    # Create directory structure
    (tmp_path / "data").mkdir()
    (tmp_path / "docs" / "dev" / "design").mkdir(parents=True)
    (tmp_path / ".github" / "workflows").mkdir(parents=True)

    # Create shared-sot.yaml with Go version
    shared_sot = {
        "go_dependencies": {
            "language": [
                {"package": "Go", "version": "1.25+", "purpose": "Backend"},
            ],
        },
    }
    (tmp_path / "data" / "shared-sot.yaml").write_text(
        yaml.dump(shared_sot, default_flow_style=False)
    )

    # Create .tool-versions
    (tmp_path / ".tool-versions").write_text("golang 1.23\nnode 20.0.0\n")

    # Create go.mod
    go_mod = """module github.com/example/revenge

go 1.23

require (
    github.com/example/pkg v1.0.0
)
"""
    (tmp_path / "go.mod").write_text(go_mod)

    # Create GitHub workflow
    workflow = """name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go test ./...
"""
    (tmp_path / ".github" / "workflows" / "ci.yml").write_text(workflow)

    # Create Dockerfile
    dockerfile = """FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN go build -o server .

FROM alpine:latest
COPY --from=builder /app/server /server
CMD ["/server"]
"""
    (tmp_path / "Dockerfile").write_text(dockerfile)

    return tmp_path


class TestInitialization:
    """Test ConfigSync initialization."""

    def test_init_sets_paths(self, sync_setup):
        """Test initialization sets correct paths."""
        syncer = ConfigSync(sync_setup)

        assert syncer.repo_root == sync_setup
        assert syncer.sot_file == sync_setup / "docs" / "dev" / "design" / "00_SOURCE_OF_TRUTH.md"
        assert syncer.shared_sot == sync_setup / "data" / "shared-sot.yaml"


class TestSOTDataLoading:
    """Test SOT data loading."""

    def test_load_existing_sot(self, sync_setup):
        """Test loading existing shared-sot.yaml."""
        syncer = ConfigSync(sync_setup)
        data = syncer.load_sot_data()

        assert data is not None
        assert "go_dependencies" in data
        assert "language" in data["go_dependencies"]

    def test_load_missing_sot(self, sync_setup):
        """Test handling missing shared-sot.yaml."""
        # Remove shared-sot.yaml
        (sync_setup / "data" / "shared-sot.yaml").unlink()

        syncer = ConfigSync(sync_setup)
        data = syncer.load_sot_data()

        assert data == {}

    def test_load_invalid_yaml(self, sync_setup):
        """Test handling invalid YAML."""
        # Write invalid YAML
        (sync_setup / "data" / "shared-sot.yaml").write_text("invalid: yaml: content:")

        syncer = ConfigSync(sync_setup)
        with pytest.raises(yaml.YAMLError):
            syncer.load_sot_data()


class TestToolVersionsSync:
    """Test .tool-versions synchronization."""

    def test_sync_tool_versions_dry_run(self, sync_setup):
        """Test dry-run mode doesn't write changes."""
        syncer = ConfigSync(sync_setup)

        # Read original content
        original = (sync_setup / ".tool-versions").read_text()

        # Sync in dry-run
        stats = syncer.sync_tool_versions(dry_run=True)

        # Content should be unchanged
        assert (sync_setup / ".tool-versions").read_text() == original
        assert stats["updated"] == 1  # Would update
        assert stats["errors"] == 0

    def test_sync_tool_versions_live(self, sync_setup):
        """Test live mode writes changes."""
        syncer = ConfigSync(sync_setup)

        # Sync in live mode
        stats = syncer.sync_tool_versions(dry_run=False)

        # Content should be updated
        content = (sync_setup / ".tool-versions").read_text()
        assert "golang 1.25" in content
        assert "golang 1.23" not in content
        assert stats["updated"] == 1
        assert stats["errors"] == 0

    def test_sync_tool_versions_no_change(self, sync_setup):
        """Test when version is already correct."""
        # Set correct version
        (sync_setup / ".tool-versions").write_text("golang 1.25\nnode 20.0.0\n")

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_tool_versions(dry_run=False)

        assert stats["unchanged"] == 1
        assert stats["updated"] == 0

    def test_sync_tool_versions_missing_file(self, sync_setup):
        """Test handling missing .tool-versions."""
        (sync_setup / ".tool-versions").unlink()

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_tool_versions(dry_run=False)

        assert stats["errors"] == 1

    def test_sync_tool_versions_missing_go_in_sot(self, sync_setup):
        """Test handling missing Go version in SOT."""
        # Create SOT without Go version
        (sync_setup / "data" / "shared-sot.yaml").write_text(
            yaml.dump({"go_dependencies": {"language": []}})
        )

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_tool_versions(dry_run=False)

        assert stats["errors"] == 1


class TestGoModSync:
    """Test go.mod synchronization."""

    def test_sync_go_mod_dry_run(self, sync_setup):
        """Test dry-run mode for go.mod."""
        syncer = ConfigSync(sync_setup)
        original = (sync_setup / "go.mod").read_text()

        stats = syncer.sync_go_mod(dry_run=True)

        assert (sync_setup / "go.mod").read_text() == original
        assert stats["updated"] == 1

    def test_sync_go_mod_live(self, sync_setup):
        """Test live mode for go.mod."""
        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_go_mod(dry_run=False)

        content = (sync_setup / "go.mod").read_text()
        assert "go 1.25" in content
        assert "go 1.23" not in content
        assert stats["updated"] == 1

    def test_sync_go_mod_preserves_other_content(self, sync_setup):
        """Test that other go.mod content is preserved."""
        syncer = ConfigSync(sync_setup)
        syncer.sync_go_mod(dry_run=False)

        content = (sync_setup / "go.mod").read_text()
        assert "module github.com/example/revenge" in content
        assert "require" in content
        assert "github.com/example/pkg" in content

    def test_sync_go_mod_missing_file(self, sync_setup):
        """Test handling missing go.mod."""
        (sync_setup / "go.mod").unlink()

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_go_mod(dry_run=False)

        assert stats["errors"] == 1


class TestGitHubWorkflowsSync:
    """Test GitHub workflows synchronization."""

    def test_sync_workflows_dry_run(self, sync_setup):
        """Test dry-run mode for workflows."""
        syncer = ConfigSync(sync_setup)
        workflow_file = sync_setup / ".github" / "workflows" / "ci.yml"
        original = workflow_file.read_text()

        stats = syncer.sync_github_workflows(dry_run=True)

        assert workflow_file.read_text() == original
        assert stats["updated"] == 1

    def test_sync_workflows_live(self, sync_setup):
        """Test live mode for workflows."""
        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_github_workflows(dry_run=False)

        workflow_file = sync_setup / ".github" / "workflows" / "ci.yml"
        content = workflow_file.read_text()
        assert "go-version: '1.25'" in content
        assert "go-version: '1.23'" not in content
        assert stats["updated"] == 1

    def test_sync_workflows_multiple_files(self, sync_setup):
        """Test syncing multiple workflow files."""
        # Create second workflow
        workflow2 = """name: Release
on:
  push:
    tags: ['v*']
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"
"""
        (sync_setup / ".github" / "workflows" / "release.yml").write_text(workflow2)

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_github_workflows(dry_run=False)

        # Both files should be updated
        assert stats["updated"] == 2

    def test_sync_workflows_missing_directory(self, sync_setup):
        """Test handling missing workflows directory."""
        import shutil
        shutil.rmtree(sync_setup / ".github" / "workflows")

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_github_workflows(dry_run=False)

        assert stats["errors"] == 1

    def test_sync_workflows_preserves_other_content(self, sync_setup):
        """Test that other workflow content is preserved."""
        syncer = ConfigSync(sync_setup)
        syncer.sync_github_workflows(dry_run=False)

        content = (sync_setup / ".github" / "workflows" / "ci.yml").read_text()
        assert "name: CI" in content
        assert "on: [push, pull_request]" in content
        assert "runs-on: ubuntu-latest" in content


class TestDockerConfigSync:
    """Test Docker configuration synchronization."""

    def test_sync_dockerfile_dry_run(self, sync_setup):
        """Test dry-run mode for Dockerfile."""
        syncer = ConfigSync(sync_setup)
        original = (sync_setup / "Dockerfile").read_text()

        stats = syncer.sync_docker_configs(dry_run=True)

        assert (sync_setup / "Dockerfile").read_text() == original
        assert stats["updated"] == 1

    def test_sync_dockerfile_live(self, sync_setup):
        """Test live mode for Dockerfile."""
        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_docker_configs(dry_run=False)

        content = (sync_setup / "Dockerfile").read_text()
        assert "FROM golang:1.25 AS builder" in content
        assert "FROM golang:1.23" not in content
        assert stats["updated"] == 1

    def test_sync_dockerfile_preserves_other_content(self, sync_setup):
        """Test that other Dockerfile content is preserved."""
        syncer = ConfigSync(sync_setup)
        syncer.sync_docker_configs(dry_run=False)

        content = (sync_setup / "Dockerfile").read_text()
        assert "WORKDIR /app" in content
        assert "FROM alpine:latest" in content
        assert "CMD" in content

    def test_sync_dockerfile_missing_file(self, sync_setup):
        """Test handling missing Dockerfile."""
        (sync_setup / "Dockerfile").unlink()

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_docker_configs(dry_run=False)

        # Should not error, just return unchanged
        assert stats["errors"] == 0


class TestFullSyncWorkflow:
    """Test complete synchronization workflow."""

    def test_sync_all_dry_run(self, sync_setup):
        """Test sync_all in dry-run mode."""
        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_all(dry_run=True)

        # All files should be unchanged
        assert (sync_setup / ".tool-versions").read_text() == "golang 1.23\nnode 20.0.0\n"
        assert "go 1.23" in (sync_setup / "go.mod").read_text()
        assert "go-version: '1.23'" in (sync_setup / ".github" / "workflows" / "ci.yml").read_text()
        assert "FROM golang:1.23" in (sync_setup / "Dockerfile").read_text()

        # Stats should show updates would happen
        assert stats["updated"] == 4  # tool-versions, go.mod, workflow, dockerfile
        assert stats["errors"] == 0

    def test_sync_all_live(self, sync_setup):
        """Test sync_all in live mode."""
        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_all(dry_run=False)

        # All files should be updated
        assert "golang 1.25" in (sync_setup / ".tool-versions").read_text()
        assert "go 1.25" in (sync_setup / "go.mod").read_text()
        assert "go-version: '1.25'" in (sync_setup / ".github" / "workflows" / "ci.yml").read_text()
        assert "FROM golang:1.25" in (sync_setup / "Dockerfile").read_text()

        assert stats["updated"] == 4
        assert stats["errors"] == 0

    def test_sync_all_aggregates_stats(self, sync_setup):
        """Test that sync_all aggregates statistics correctly."""
        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_all(dry_run=False)

        assert "updated" in stats
        assert "unchanged" in stats
        assert "errors" in stats
        assert isinstance(stats["updated"], int)
        assert isinstance(stats["unchanged"], int)
        assert isinstance(stats["errors"], int)

    def test_sync_all_partial_success(self, sync_setup):
        """Test sync_all with some files missing."""
        # Remove some files
        (sync_setup / ".tool-versions").unlink()
        (sync_setup / "go.mod").unlink()

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_all(dry_run=False)

        # Should update remaining files
        assert stats["updated"] == 2  # workflows + dockerfile
        assert stats["errors"] == 2  # tool-versions + go.mod


class TestVersionExtraction:
    """Test Go version extraction logic."""

    def test_extract_version_with_plus(self, sync_setup):
        """Test extracting version with + suffix."""
        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_tool_versions(dry_run=False)

        # "1.25+" should be extracted as "1.25"
        content = (sync_setup / ".tool-versions").read_text()
        assert "golang 1.25" in content

    def test_extract_version_exact(self, sync_setup):
        """Test extracting exact version."""
        # Create SOT with exact version
        sot = {
            "go_dependencies": {
                "language": [{"package": "Go", "version": "1.24.3"}],
            },
        }
        (sync_setup / "data" / "shared-sot.yaml").write_text(yaml.dump(sot))

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_tool_versions(dry_run=False)

        content = (sync_setup / ".tool-versions").read_text()
        assert "golang 1.24" in content


class TestEdgeCases:
    """Test edge cases and error conditions."""

    def test_empty_sot_data(self, sync_setup):
        """Test handling empty SOT data."""
        (sync_setup / "data" / "shared-sot.yaml").write_text("")

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_all(dry_run=False)

        # Should handle gracefully with errors
        assert stats["errors"] > 0

    def test_malformed_config_files(self, sync_setup):
        """Test handling malformed config files."""
        # Create malformed go.mod (no version line)
        (sync_setup / "go.mod").write_text("module example\n")

        syncer = ConfigSync(sync_setup)
        stats = syncer.sync_go_mod(dry_run=False)

        # Should not error, just report unchanged
        assert stats["unchanged"] == 1

    def test_concurrent_sync_safety(self, sync_setup):
        """Test that sync operations don't interfere."""
        syncer = ConfigSync(sync_setup)

        # Run multiple syncs
        stats1 = syncer.sync_all(dry_run=False)
        stats2 = syncer.sync_all(dry_run=False)

        # First should update, second should report unchanged
        assert stats1["updated"] == 4
        assert stats2["unchanged"] == 4


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
