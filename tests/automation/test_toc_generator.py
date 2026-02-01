#!/usr/bin/env python3
"""Tests for TOC generator."""

import sys
from pathlib import Path

# Add repo root to path
repo_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(repo_root))

from scripts.automation.toc_generator import TOCGenerator


def test_split_frontmatter_with_leading_blanks():
    """Test splitting frontmatter with leading blank lines."""
    generator = TOCGenerator()
    
    content = "\n\n---\nkey: value\n---\n\n# Title\n\nContent"
    frontmatter, body = generator.split_frontmatter(content)
    
    assert frontmatter == "\n\n---\nkey: value\n---\n"
    assert body == "\n# Title\n\nContent"


def test_split_frontmatter_no_blanks():
    """Test splitting frontmatter without leading blanks."""
    generator = TOCGenerator()
    
    content = "---\nkey: value\n---\n\n# Title"
    frontmatter, body = generator.split_frontmatter(content)
    
    assert frontmatter == "---\nkey: value\n---\n"
    assert body == "\n# Title"


def test_split_frontmatter_empty():
    """Test splitting empty frontmatter."""
    generator = TOCGenerator()
    
    content = "\n\n---\n---\n\n# Title"
    frontmatter, body = generator.split_frontmatter(content)
    
    assert frontmatter == "\n\n---\n---\n"
    assert body == "\n# Title"


def test_split_frontmatter_no_frontmatter():
    """Test content without frontmatter."""
    generator = TOCGenerator()
    
    content = "# Title\n\nContent"
    frontmatter, body = generator.split_frontmatter(content)
    
    assert frontmatter == ""
    assert body == "# Title\n\nContent"


def test_add_toc_preserves_frontmatter():
    """Test that TOC is added after frontmatter."""
    generator = TOCGenerator()
    
    content = "\n\n---\nsources:\n  - test\n---\n\n# Title\n\n## Section 1\n\n## Section 2"
    result = generator.add_toc(content)
    
    # Should start with frontmatter
    assert result.startswith("\n\n---\nsources:\n  - test\n---\n")
    # Should have TOC after frontmatter
    assert "## Table of Contents" in result
    # Should have title after TOC
    assert "# Title" in result


def test_extract_headers():
    """Test header extraction."""
    generator = TOCGenerator()
    
    content = "# Title\n\n## Section 1\n\n### Subsection\n\n## Section 2"
    headers = generator.extract_headers(content)

    assert len(headers) == 4
    assert headers[0] == (1, "Title")
    assert headers[1] == (2, "Section 1")
    assert headers[2] == (3, "Subsection")
    assert headers[3] == (2, "Section 2")


if __name__ == "__main__":
    print("Running TOC generator tests...")
    test_split_frontmatter_with_leading_blanks()
    print("✓ test_split_frontmatter_with_leading_blanks")
    test_split_frontmatter_no_blanks()
    print("✓ test_split_frontmatter_no_blanks")
    test_split_frontmatter_empty()
    print("✓ test_split_frontmatter_empty")
    test_split_frontmatter_no_frontmatter()
    print("✓ test_split_frontmatter_no_frontmatter")
    test_add_toc_preserves_frontmatter()
    print("✓ test_add_toc_preserves_frontmatter")
    test_extract_headers()
    print("✓ test_extract_headers")
    print("\n✅ All TOC generator tests passed!")
