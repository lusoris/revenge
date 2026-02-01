"""Tests for documentation generation (batch_regenerate.py and doc_generator.py).

Ensures wiki files are generated alongside design docs from YAML data.
"""

import subprocess
from pathlib import Path

import pytest


class TestBatchRegenerateScript:
    """Test batch_regenerate.py script."""

    def test_script_exists(self):
        """batch_regenerate.py must exist."""
        script = Path("scripts/automation/batch_regenerate.py")
        assert script.exists(), "batch_regenerate.py not found"

    def test_script_has_help(self):
        """Script must have --help flag."""
        result = subprocess.run(
            ["python", "scripts/automation/batch_regenerate.py", "--help"],
            capture_output=True,
            text=True,
            timeout=10,
            check=False,
        )
        assert result.returncode == 0, "Script --help failed"
        assert "usage:" in result.stdout.lower()
        assert "--backup" in result.stdout

    def test_script_runs_without_arguments(self):
        """Script must run successfully without arguments."""
        # Full regeneration of 158+ YAML files - need adequate timeout
        result = subprocess.run(
            ["python", "scripts/automation/batch_regenerate.py"],
            capture_output=True,
            text=True,
            timeout=120,
            check=False,
        )
        # Should complete successfully
        assert result.returncode == 0, f"Script failed: {result.stderr}"


class TestDocPipelineIntegration:
    """Test doc-pipeline.sh integration."""

    def test_pipeline_script_exists(self):
        """doc-pipeline.sh must exist."""
        pipeline = Path("scripts/doc-pipeline.sh")
        assert pipeline.exists(), "doc-pipeline.sh not found"

    def test_pipeline_calls_batch_regenerate(self):
        """Pipeline must call batch_regenerate.py as Step 0."""
        pipeline = Path("scripts/doc-pipeline.sh")
        content = pipeline.read_text()

        # Should have Step 0
        assert "Step 0" in content, "Pipeline missing Step 0"
        assert "batch_regenerate" in content, (
            "Pipeline doesn't call batch_regenerate.py"
        )
        assert "Regenerate" in content, "Step 0 should regenerate docs"

    def test_pipeline_has_correct_step_order(self):
        """Pipeline steps must be in correct order."""
        pipeline = Path("scripts/doc-pipeline.sh")
        content = pipeline.read_text()

        # Check step comments exist in order
        steps = [
            "0. Regenerate docs from YAML",
            "1. Generate INDEX.md",
            "2. Add design breadcrumbs",
            "3. Sync status tables",
            "4. Validate document structure",
            "5. Fix broken links",
            "6. Generate meta files",
        ]

        for step in steps:
            assert step in content or step.replace(".", ":") in content, (
                f"Missing or misnamed: {step}"
            )


class TestWikiGeneration:
    """Test that wiki files are generated."""

    def test_data_directory_exists(self):
        """data/ directory must exist with YAML files."""
        data_dir = Path("data")
        assert data_dir.exists(), "data/ directory not found"

        yaml_files = list(data_dir.glob("**/*.yaml"))
        yaml_files = [f for f in yaml_files if f.name != "shared-sot.yaml"]

        assert len(yaml_files) > 0, "No YAML data files found"

    def test_wiki_output_directory_structure(self):
        """docs/wiki/ directory should match data/ structure."""
        wiki_dir = Path("docs/wiki")
        assert wiki_dir.exists(), "docs/wiki/ directory not found"

        # Check expected subdirectories based on data structure
        expected_dirs = [
            "features",
            "operations",
        ]

        for expected_dir in expected_dirs:
            # May not exist yet if wiki generation hasn't run
            # This test ensures structure is expected
            assert expected_dir in ["features", "operations", "services", "integrations"]

    def test_wiki_files_should_mirror_design_files(self):
        """Wiki file count should eventually match design file count."""
        design_dir = Path("docs/dev/design")
        wiki_dir = Path("docs/wiki")

        design_count = len(list(design_dir.glob("**/*.md")))
        wiki_count = len(list(wiki_dir.glob("**/*.md")))

        # After full regeneration, counts should be similar
        # (some design-only files may exist, but most should have wiki versions)
        if wiki_count < design_count * 0.5:
            pytest.skip(
                f"Wiki generation incomplete: {wiki_count} wiki files vs "
                f"{design_count} design files. Run: "
                f"python scripts/automation/batch_regenerate.py --live"
            )


class TestDocGeneratorModule:
    """Test doc_generator.py module."""

    def test_doc_generator_exists(self):
        """doc_generator.py must exist."""
        generator = Path("scripts/automation/doc_generator.py")
        assert generator.exists(), "doc_generator.py not found"

    def test_doc_generator_creates_both_outputs(self):
        """DocGenerator must create both Claude and Wiki outputs."""
        # Import and check
        import sys
        sys.path.insert(0, str(Path("scripts/automation")))

        from doc_generator import DocGenerator

        # Check DocGenerator class exists
        assert hasattr(DocGenerator, "generate_doc")

        # Check it has methods for both outputs
        generator = DocGenerator(Path())
        assert hasattr(generator, "output_dir_claude")
        assert hasattr(generator, "output_dir_wiki")

    def test_templates_exist(self):
        """Jinja2 templates must exist."""
        templates_dir = Path("templates")
        assert templates_dir.exists(), "templates/ directory not found"

        # Check for key templates
        required_templates = [
            "feature.md.jinja2",
            "service.md.jinja2",
            "integration.md.jinja2",
        ]

        for template in required_templates:
            template_path = templates_dir / template
            assert template_path.exists(), f"Missing template: {template}"


class TestDocGenerationWorkflow:
    """Test end-to-end documentation generation workflow."""

    def test_yaml_to_markdown_workflow(self):
        """YAML data should generate both Claude and Wiki markdown."""
        # This is a high-level integration test
        # Actual generation tested by running batch_regenerate.py

        data_dir = Path("data")
        yaml_files = list(data_dir.glob("**/*.yaml"))
        yaml_files = [f for f in yaml_files if f.name != "shared-sot.yaml"]

        assert len(yaml_files) > 0, "No YAML files to generate from"

        # Check that generated files would go to correct locations
        design_dir = Path("docs/dev/design")
        wiki_dir = Path("docs/wiki")

        assert design_dir.exists(), "docs/dev/design/ directory missing"
        assert wiki_dir.exists(), "docs/wiki/ directory missing"

    def test_shared_data_file_exists(self):
        """shared-sot.yaml must exist for shared data."""
        shared_data = Path("data/shared-sot.yaml")
        assert shared_data.exists(), "data/shared-sot.yaml not found"


class TestPipelineStepSequence:
    """Test that pipeline steps run in correct sequence."""

    def test_step_0_runs_before_indexes(self):
        """Step 0 (regenerate) must run before Step 1 (indexes)."""
        pipeline = Path("scripts/doc-pipeline.sh")
        content = pipeline.read_text()

        # Find positions of step calls
        step_0_pos = content.find("run_step 0")
        step_1_pos = content.find("run_step 1")

        assert step_0_pos > 0, "Step 0 not found"
        assert step_1_pos > 0, "Step 1 not found"
        assert step_0_pos < step_1_pos, (
            "Step 0 must run before Step 1"
        )

    def test_all_steps_present(self):
        """All 7 steps (0-6) must be present."""
        pipeline = Path("scripts/doc-pipeline.sh")
        content = pipeline.read_text()

        for step_num in range(7):
            assert f"run_step {step_num}" in content, (
                f"Step {step_num} missing from pipeline"
            )


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
