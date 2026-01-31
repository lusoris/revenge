#!/usr/bin/env python3
"""
Test script to render DESIGN_TEMPLATE.md.jinja2 with test_data.yaml
Generates both Claude and Wiki versions for validation
"""

import sys
from pathlib import Path
import yaml
from jinja2 import Environment, FileSystemLoader, StrictUndefined

def render_template(template_path: Path, data_path: Path, target: str) -> str:
    """
    Render Jinja2 template with data

    Args:
        template_path: Path to .jinja2 template
        data_path: Path to YAML data file
        target: 'claude' or 'wiki'

    Returns:
        Rendered markdown content
    """
    # Load data
    with open(data_path, 'r') as f:
        data = yaml.safe_load(f)

    # Add target to data (for conditionals)
    data['claude'] = (target == 'claude')
    data['wiki'] = (target == 'wiki')

    # Set up Jinja2 environment
    env = Environment(
        loader=FileSystemLoader(template_path.parent),
        undefined=StrictUndefined,  # Fail on undefined variables
        trim_blocks=True,
        lstrip_blocks=True,
    )

    # Load template
    template = env.get_template(template_path.name)

    # Render
    return template.render(**data)

def main():
    """Test template rendering"""
    script_dir = Path(__file__).parent
    template_path = script_dir / "DESIGN_TEMPLATE.md.jinja2"
    data_path = script_dir / "test_data.yaml"

    if not template_path.exists():
        print(f"❌ Template not found: {template_path}", file=sys.stderr)
        return 1

    if not data_path.exists():
        print(f"❌ Data file not found: {data_path}", file=sys.stderr)
        return 1

    print("🔧 Testing template rendering...")
    print()

    # Test Claude version
    print("📝 Rendering Claude version...")
    try:
        claude_output = render_template(template_path, data_path, 'claude')
        claude_file = script_dir / "test_output_claude.md"
        with open(claude_file, 'w') as f:
            f.write(claude_output)
        print(f"✅ Claude version rendered successfully: {claude_file}")
        print(f"   Size: {len(claude_output)} bytes, {claude_output.count(chr(10))} lines")
    except Exception as e:
        print(f"❌ Claude rendering failed: {e}", file=sys.stderr)
        return 1

    print()

    # Test Wiki version
    print("📝 Rendering Wiki version...")
    try:
        wiki_output = render_template(template_path, data_path, 'wiki')
        wiki_file = script_dir / "test_output_wiki.md"
        with open(wiki_file, 'w') as f:
            f.write(wiki_output)
        print(f"✅ Wiki version rendered successfully: {wiki_file}")
        print(f"   Size: {len(wiki_output)} bytes, {wiki_output.count(chr(10))} lines")
    except Exception as e:
        print(f"❌ Wiki rendering failed: {e}", file=sys.stderr)
        return 1

    print()
    print("🎉 Template rendering test complete!")
    print()
    print("Next steps:")
    print("1. Review test_output_claude.md")
    print("2. Review test_output_wiki.md")
    print("3. Verify conditional content is correct")
    print("4. Check for rendering errors or missing sections")

    return 0

if __name__ == "__main__":
    sys.exit(main())
