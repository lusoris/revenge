#!/usr/bin/env python3
"""Documentation Generator - Render Jinja2 templates with YAML data.

This generator:
1. Loads YAML data files (shared + doc-specific)
2. Renders Jinja2 templates (base + inheritance)
3. Generates dual output (Claude + Wiki)
4. Post-processes (TOC generation, formatting)
5. Atomic writes with validation

Author: Automation System
Created: 2026-01-31
"""

import sys
from pathlib import Path
from typing import Any

import yaml
from jinja2 import Environment, FileSystemLoader, StrictUndefined


# Add repo root to Python path for imports
repo_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(repo_root))

# Import TOCGenerator - handle both script and module contexts
try:
    from .toc_generator import TOCGenerator
except ImportError:
    from scripts.automation.toc_generator import TOCGenerator


class DocGenerator:
    """Generate documentation from YAML data and Jinja2 templates."""

    def __init__(self, repo_root: Path):
        """Initialize generator with repository root."""
        self.repo_root = repo_root
        self.templates_dir = repo_root / "templates"
        self.data_dir = repo_root / "data"
        self.output_dir_claude = repo_root / "docs" / "dev" / "design"
        self.output_dir_wiki = repo_root / "docs" / "wiki"

        # Initialize Jinja2 environment
        self.env = Environment(
            loader=FileSystemLoader(self.templates_dir),
            undefined=StrictUndefined,
            trim_blocks=False,
            lstrip_blocks=False,
        )

        # Initialize TOC generator
        self.toc_generator = TOCGenerator()

        # Load shared SOT data
        self.shared_data = self._load_shared_data()

        # Load sources mapping for local links
        self.sources_mapping = self._load_sources_mapping()

        # Add custom filter to Jinja2 environment
        self.env.filters["to_local_source"] = self._url_to_local_source

    def _load_shared_data(self) -> dict[str, Any]:
        """Load shared data from shared-sot.yaml."""
        shared_path = self.data_dir / "shared-sot.yaml"
        if not shared_path.exists():
            print(f"⚠️  Warning: {shared_path} not found, using empty shared data")
            return {}

        with open(shared_path, encoding="utf-8") as f:
            data = yaml.safe_load(f)
            print(f"✓ Loaded shared data from {shared_path}")
            return data

    def _load_sources_mapping(self) -> dict[str, str]:
        """Load SOURCES.yaml and create URL -> local path mapping."""
        sources_path = self.repo_root / "docs" / "dev" / "sources" / "SOURCES.yaml"
        if not sources_path.exists():
            print(f"⚠️  Warning: {sources_path} not found, using empty sources mapping")
            return {}

        with open(sources_path, encoding="utf-8") as f:
            sources_config = yaml.safe_load(f)

        # Build mapping: URL -> relative path in docs/dev/sources/
        url_to_path = {}
        for category, sources_list in sources_config.get("sources", {}).items():
            for source in sources_list:
                if "url" in source and "output" in source:
                    url_to_path[source["url"]] = source["output"]

        print(f"✓ Loaded {len(url_to_path)} source URL mappings")
        return url_to_path

    def _url_to_local_source(self, url: str, depth: int = 0) -> str:
        """Convert external URL to local source path if available.

        Args:
            url: External URL (e.g., 'https://pkg.go.dev/go.uber.org/fx')
            depth: Subdirectory depth (0 for root, 1 for one level deep, etc.)

        Returns:
            Local markdown link if URL is in mapping, otherwise original URL
        """
        if url in self.sources_mapping:
            local_path = self.sources_mapping[url]
            # Calculate relative path based on depth
            # depth=0: ../sources/ (from docs/dev/design/)
            # depth=1: ../../sources/ (from docs/dev/design/subdir/)
            prefix = "../" * (depth + 1)
            return f"{prefix}sources/{local_path}"
        return url

    def generate_doc(
        self,
        data_file: Path,
        template_name: str,
        output_subpath: str,
        render_both: bool = True,
    ) -> dict[str, Path]:
        """Generate documentation from data file and template.

        Args:
            data_file: Path to YAML data file
            template_name: Template name (e.g., 'feature.md.jinja2')
            output_subpath: Subdirectory path for output (e.g., 'features/video')
            render_both: If True, render both Claude and Wiki versions

        Returns:
            Dict with paths to generated files
        """
        print(f"\n📄 Generating: {data_file.stem}")

        # Load doc-specific data
        with open(data_file, encoding="utf-8") as f:
            doc_data = yaml.safe_load(f)

        # Merge shared data with doc-specific data
        merged_data = self._merge_data(self.shared_data, doc_data)

        # Calculate depth from output_subpath for relative linking
        # Empty or '.' = depth 0, 'features' = depth 1, 'features/video' = depth 2
        depth = (
            0
            if not output_subpath or output_subpath == "."
            else output_subpath.count("/") + 1
        )

        # Get template
        template = self.env.get_template(template_name)

        # Create depth-aware filter for this render
        def to_local_source_with_depth(url: str) -> str:
            return self._url_to_local_source(url, depth)

        # Temporarily override the filter for this render
        original_filter = self.env.filters.get("to_local_source")
        self.env.filters["to_local_source"] = to_local_source_with_depth

        generated_files = {}

        # Generate Claude version (always generated)
        if True:
            claude_data = {**merged_data, "claude": True, "wiki": False}
            claude_output = template.render(**claude_data)

            # Post-process: Add TOC
            claude_output = self.toc_generator.add_toc(claude_output)

            claude_path = self._save_output(
                claude_output,
                self.output_dir_claude / output_subpath / f"{data_file.stem}.md",
            )
            generated_files["claude"] = claude_path
            print(f"  ✓ Claude: {claude_path.relative_to(self.repo_root)}")

        # Generate Wiki version
        if render_both:
            wiki_data = {**merged_data, "claude": False, "wiki": True}
            wiki_output = template.render(**wiki_data)

            # Post-process: Add TOC
            wiki_output = self.toc_generator.add_toc(wiki_output)

            wiki_path = self._save_output(
                wiki_output,
                self.output_dir_wiki / output_subpath / f"{data_file.stem}.md",
            )
            generated_files["wiki"] = wiki_path
            print(f"  ✓ Wiki: {wiki_path.relative_to(self.repo_root)}")

        # Restore original filter
        if original_filter:
            self.env.filters["to_local_source"] = original_filter

        return generated_files

    def _merge_data(
        self,
        shared: dict[str, Any],
        doc: dict[str, Any],
    ) -> dict[str, Any]:
        """Merge shared data with doc-specific data.

        Doc-specific data takes precedence.
        """
        merged = {}

        # Add shared metadata (versions, etc.)
        if "metadata" in shared:
            merged.update(shared["metadata"])

        # Add shared infrastructure, dependencies, etc.
        merged["shared_infrastructure"] = shared.get("infrastructure", [])
        merged["shared_go_dependencies"] = shared.get("go_dependencies", {})
        merged["shared_design_principles"] = shared.get("design_principles", {})

        # Merge doc-specific data (takes precedence)
        merged.update(doc)

        return merged

    def _save_output(self, content: str, output_path: Path) -> Path:
        """Save generated content to file with atomic write.

        Creates parent directories if needed.
        """
        # Create parent directories
        output_path.parent.mkdir(parents=True, exist_ok=True)

        # Atomic write: write to temp, then rename
        temp_path = output_path.with_suffix(".tmp")
        try:
            with open(temp_path, "w", encoding="utf-8") as f:
                f.write(content)

            # Rename to final path (atomic on POSIX, requires unlink on Windows)
            if output_path.exists():
                output_path.unlink()
            temp_path.rename(output_path)

        except Exception as e:
            # Clean up temp file on error
            if temp_path.exists():
                temp_path.unlink()
            raise e

        return output_path


def main():
    """Main entry point - Test with MOVIE_MODULE."""
    repo_root = Path(__file__).parent.parent.parent

    # Initialize generator
    generator = DocGenerator(repo_root)

    # Test with MOVIE_MODULE.yaml
    data_file = repo_root / "data" / "features" / "video" / "MOVIE_MODULE.yaml"

    if not data_file.exists():
        print(f"❌ Error: {data_file} not found")
        sys.exit(1)

    # Generate documentation
    generated = generator.generate_doc(
        data_file=data_file,
        template_name="feature.md.jinja2",
        output_subpath="features/video",
        render_both=True,
    )

    print("\n✅ Documentation generated!")
    print(f"   Claude: {generated['claude']}")
    print(f"   Wiki: {generated['wiki']}")


if __name__ == "__main__":
    main()
