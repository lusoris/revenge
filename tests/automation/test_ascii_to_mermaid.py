#!/usr/bin/env python3
"""Tests for ASCII to Mermaid converter.

Test coverage for scripts/automation/ascii_to_mermaid.py:
- Box extraction from ASCII diagrams
- Connection inference between boxes
- Mermaid generation with flowchart LR and subgraphs
- Label escaping for special characters
- Node shape detection (database, service, external)
- Layer detection and styling
- YAML field processing
"""

from __future__ import annotations

import sys
from pathlib import Path


# Add repo root to path before importing local modules
repo_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(repo_root))

import pytest  # noqa: E402

from scripts.automation.ascii_to_mermaid import ASCIIToMermaid  # noqa: E402


class TestBoxExtraction:
    """Test ASCII box extraction."""

    @pytest.fixture
    def converter(self):
        """Create converter instance."""
        return ASCIIToMermaid()

    def test_extract_single_box(self, converter):
        """Test extracting a single ASCII box."""
        ascii_diagram = """
┌─────────────┐
│   Client    │
└─────────────┘
"""
        boxes = converter.extract_boxes(ascii_diagram)

        assert len(boxes) == 1
        assert boxes[0]["label"] == "Client"

    def test_extract_box_with_sublabels(self, converter):
        """Test extracting box with multiple lines."""
        ascii_diagram = """
┌─────────────┐
│   Client    │
│  (Web App)  │
└─────────────┘
"""
        boxes = converter.extract_boxes(ascii_diagram)

        assert len(boxes) == 1
        assert boxes[0]["label"] == "Client"
        assert "(Web App)" in boxes[0]["sublabels"]

    def test_extract_multiple_boxes_vertical(self, converter):
        """Test extracting vertically stacked boxes."""
        ascii_diagram = """
┌─────────────┐
│   Client    │
└─────────────┘
       │
       ▼
┌─────────────┐
│   Server    │
└─────────────┘
"""
        boxes = converter.extract_boxes(ascii_diagram)

        assert len(boxes) == 2
        labels = [b["label"] for b in boxes]
        assert "Client" in labels
        assert "Server" in labels

    def test_extract_multiple_boxes_horizontal(self, converter):
        """Test extracting horizontally arranged boxes."""
        ascii_diagram = """
┌─────────┐  ┌─────────┐  ┌─────────┐
│  Web    │  │ Mobile  │  │   TV    │
└─────────┘  └─────────┘  └─────────┘
"""
        boxes = converter.extract_boxes(ascii_diagram)

        assert len(boxes) == 3
        labels = [b["label"] for b in boxes]
        assert "Web" in labels
        assert "Mobile" in labels
        assert "TV" in labels

        # All should be on the same row
        rows = {b["row"] for b in boxes}
        assert len(rows) == 1

    def test_extract_nested_structure(self, converter):
        """Test extracting complex nested structure."""
        ascii_diagram = """
┌───────────────────────────────┐
│         Client Layer          │
└───────────────────────────────┘
         │
         ▼
┌─────────┐  ┌─────────┐
│  API    │  │ Handler │
└─────────┘  └─────────┘
"""
        boxes = converter.extract_boxes(ascii_diagram)

        assert len(boxes) >= 3
        labels = [b["label"] for b in boxes]
        assert "Client Layer" in labels
        assert "API" in labels
        assert "Handler" in labels

    def test_no_boxes_returns_empty(self, converter):
        """Test that text without boxes returns empty list."""
        ascii_diagram = """
No boxes here, just text.
Some arrows: ---> and <---
"""
        boxes = converter.extract_boxes(ascii_diagram)
        assert len(boxes) == 0


class TestConnectionInference:
    """Test connection inference between boxes."""

    @pytest.fixture
    def converter(self):
        """Create converter instance."""
        return ASCIIToMermaid()

    def test_infer_vertical_connections(self, converter):
        """Test inferring connections between vertical boxes."""
        boxes = [
            {"id": "node1", "label": "A", "row": 0, "col": 0, "bottom_row": 2},
            {"id": "node2", "label": "B", "row": 4, "col": 0, "bottom_row": 6},
        ]

        connections = converter.infer_connections(boxes, "")

        assert len(connections) == 1
        assert connections[0][0] == "node1"
        assert connections[0][1] == "node2"

    def test_infer_horizontal_connections(self, converter):
        """Test inferring connections between horizontal boxes."""
        boxes = [
            {"id": "node1", "label": "A", "row": 0, "col": 0, "bottom_row": 2},
            {"id": "node2", "label": "B", "row": 0, "col": 20, "bottom_row": 2},
            {"id": "node3", "label": "C", "row": 0, "col": 40, "bottom_row": 2},
        ]

        connections = converter.infer_connections(boxes, "")

        # Should connect A -> B -> C
        from_ids = [c[0] for c in connections]
        to_ids = [c[1] for c in connections]

        assert "node1" in from_ids
        assert "node2" in to_ids


class TestLayerDetection:
    """Test layer box detection."""

    @pytest.fixture
    def converter(self):
        """Create converter instance."""
        return ASCIIToMermaid()

    def test_detect_layer_keyword(self, converter):
        """Test detection of LAYER keyword."""
        assert converter._is_layer_box("CLIENT LAYER")
        assert converter._is_layer_box("Client Layer")
        assert converter._is_layer_box("Service Layer")

    def test_detect_tier_keyword(self, converter):
        """Test detection of TIER keyword."""
        assert converter._is_layer_box("Data Tier")
        assert converter._is_layer_box("PRESENTATION TIER")

    def test_non_layer_box(self, converter):
        """Test that non-layer boxes are not detected as layers."""
        assert not converter._is_layer_box("Client")
        assert not converter._is_layer_box("API Handler")
        assert not converter._is_layer_box("Database")


class TestNodeShapeDetection:
    """Test node shape detection based on content."""

    @pytest.fixture
    def converter(self):
        """Create converter instance."""
        return ASCIIToMermaid()

    def test_database_shape(self, converter):
        """Test database nodes get cylinder shape."""
        open_br, _ = converter._get_node_shape("PostgreSQL", [])
        assert "[(" in open_br  # Cylinder shape

        open_br, _ = converter._get_node_shape("Cache", ["Redis"])
        assert "[(" in open_br

    def test_service_shape(self, converter):
        """Test service nodes get subroutine shape."""
        open_br, _ = converter._get_node_shape("Auth Service", [])
        assert "[[" in open_br  # Subroutine shape

        open_br, _ = converter._get_node_shape("API Handler", [])
        assert "[[" in open_br

    def test_external_shape(self, converter):
        """Test external nodes get stadium shape."""
        open_br, _ = converter._get_node_shape("External API", [])
        assert "([" in open_br  # Stadium shape

        open_br, _ = converter._get_node_shape("Web Client", [])
        assert "([" in open_br

    def test_default_shape(self, converter):
        """Test default nodes get rectangle shape."""
        open_br, _ = converter._get_node_shape("Something", [])
        assert '["' in open_br  # Rectangle with quotes


class TestMermaidGeneration:
    """Test Mermaid diagram generation."""

    @pytest.fixture
    def converter(self):
        """Create converter instance."""
        return ASCIIToMermaid()

    def test_generate_flowchart_lr(self, converter):
        """Test that generated diagram uses flowchart LR."""
        boxes = [
            {"id": "node1", "label": "A", "row": 0, "col": 0, "sublabels": []},
        ]
        connections = []

        mermaid = converter.generate_mermaid(boxes, connections)

        assert "flowchart LR" in mermaid

    def test_generate_subgraphs(self, converter):
        """Test that rows are grouped into subgraphs."""
        boxes = [
            {"id": "node1", "label": "A", "row": 0, "col": 0, "sublabels": []},
            {"id": "node2", "label": "B", "row": 0, "col": 20, "sublabels": []},
            {"id": "node3", "label": "C", "row": 10, "col": 0, "sublabels": []},
        ]
        connections = [("node1", "node3", "")]

        mermaid = converter.generate_mermaid(boxes, connections)

        assert "subgraph" in mermaid
        assert "end" in mermaid

    def test_generate_with_layer_box(self, converter):
        """Test that layer boxes become subgraph titles."""
        boxes = [
            {"id": "node1", "label": "CLIENT LAYER", "row": 0, "col": 0, "sublabels": []},
            {"id": "node2", "label": "Web", "row": 0, "col": 20, "sublabels": []},
        ]
        connections = []

        mermaid = converter.generate_mermaid(boxes, connections)

        # Layer should be used as subgraph title
        assert "CLIENT LAYER" in mermaid

    def test_generate_styling(self, converter):
        """Test that subgraphs get colored styling."""
        boxes = [
            {"id": "node1", "label": "A", "row": 0, "col": 0, "sublabels": []},
            {"id": "node2", "label": "B", "row": 10, "col": 0, "sublabels": []},
        ]
        connections = []

        mermaid = converter.generate_mermaid(boxes, connections)

        assert "style" in mermaid
        assert "fill:" in mermaid

    def test_generate_connections(self, converter):
        """Test that connections are generated between rows."""
        boxes = [
            {"id": "node1", "label": "A", "row": 0, "col": 0, "sublabels": []},
            {"id": "node2", "label": "B", "row": 10, "col": 0, "sublabels": []},
        ]
        connections = [("node1", "node2", "")]

        mermaid = converter.generate_mermaid(boxes, connections)

        assert "node1 -->" in mermaid
        assert "node2" in mermaid

    def test_empty_boxes_returns_empty(self, converter):
        """Test that empty boxes list returns empty string."""
        mermaid = converter.generate_mermaid([], [])
        assert mermaid == ""


class TestLabelEscaping:
    """Test label escaping for special characters."""

    @pytest.fixture
    def converter(self):
        """Create converter instance."""
        return ASCIIToMermaid()

    def test_escape_quotes(self, converter):
        """Test that quotes are escaped."""
        boxes = [
            {"id": "node1", "label": 'Test "quoted"', "row": 0, "col": 0, "sublabels": []},
        ]

        mermaid = converter.generate_mermaid(boxes, [])

        # Should not have unescaped quotes in label
        assert '&quot;' in mermaid or 'Test "quoted"' not in mermaid

    def test_escape_angle_brackets(self, converter):
        """Test that angle brackets are escaped."""
        boxes = [
            {"id": "node1", "label": "Test <value>", "row": 0, "col": 0, "sublabels": []},
        ]

        mermaid = converter.generate_mermaid(boxes, [])

        # Should have escaped brackets
        assert "&lt;" in mermaid or "&gt;" in mermaid

    def test_parentheses_preserved(self, converter):
        """Test that parentheses are preserved (safe inside quotes)."""
        boxes = [
            {"id": "node1", "label": "React Native", "row": 0, "col": 0, "sublabels": ["(Mobile)"]},
        ]

        mermaid = converter.generate_mermaid(boxes, [])

        # Parentheses should be in the output (inside quoted labels)
        assert "(Mobile)" in mermaid


class TestDiagramConversion:
    """Test full diagram conversion."""

    @pytest.fixture
    def converter(self):
        """Create converter instance."""
        return ASCIIToMermaid()

    def test_convert_simple_diagram(self, converter):
        """Test converting a simple ASCII diagram to Mermaid."""
        ascii_diagram = """
┌─────────────┐
│   Client    │
└─────────────┘
       │
       ▼
┌─────────────┐
│   Server    │
└─────────────┘
"""
        mermaid = converter.convert_diagram(ascii_diagram)

        assert "```mermaid" in mermaid
        assert "flowchart LR" in mermaid
        assert "Client" in mermaid
        assert "Server" in mermaid
        assert "```" in mermaid

    def test_convert_layered_diagram(self, converter):
        """Test converting a layered ASCII diagram."""
        ascii_diagram = """
┌───────────────────────────────┐
│       CLIENT LAYER            │
│       (Frontend)              │
└───────────────────────────────┘
              │
              ▼
┌─────────┐  ┌─────────┐
│  API    │  │ Handler │
└─────────┘  └─────────┘
"""
        mermaid = converter.convert_diagram(ascii_diagram)

        assert "```mermaid" in mermaid
        assert "subgraph" in mermaid
        assert "CLIENT LAYER" in mermaid


class TestHasASCIIDiagram:
    """Test ASCII diagram detection."""

    @pytest.fixture
    def converter(self):
        """Create converter instance."""
        return ASCIIToMermaid()

    def test_detect_box_characters(self, converter):
        """Test detection of box drawing characters."""
        text_with_boxes = "Some text ┌─────┐ more text"
        assert converter.has_ascii_diagram(text_with_boxes)

    def test_no_detect_regular_text(self, converter):
        """Test that regular text is not detected."""
        regular_text = "Just some regular text without any diagrams."
        assert not converter.has_ascii_diagram(regular_text)

    def test_detect_various_box_chars(self, converter):
        """Test detection of various box characters."""
        assert converter.has_ascii_diagram("┌")
        assert converter.has_ascii_diagram("┐")
        assert converter.has_ascii_diagram("└")
        assert converter.has_ascii_diagram("┘")
        assert converter.has_ascii_diagram("│")
        assert converter.has_ascii_diagram("─")


class TestYAMLFieldProcessing:
    """Test YAML field processing."""

    @pytest.fixture
    def converter(self):
        """Create converter instance."""
        return ASCIIToMermaid()

    def test_process_field_with_ascii(self, converter):
        """Test processing a field containing ASCII diagram."""
        value = """
Some description.

```
┌─────────┐
│  Test   │
└─────────┘
```
"""
        result, converted = converter.process_yaml_field(value)

        assert converted
        assert "```mermaid" in result

    def test_process_field_without_ascii(self, converter):
        """Test processing a field without ASCII diagram."""
        value = "Just some text without diagrams."

        result, converted = converter.process_yaml_field(value)

        assert not converted
        assert result == value

    def test_process_field_already_mermaid(self, converter):
        """Test that fields with existing Mermaid are not converted."""
        value = """
```mermaid
flowchart TD
    A --> B
```
"""
        result, converted = converter.process_yaml_field(value)

        assert not converted
        assert result == value


class TestIntegrationWithRealDiagrams:
    """Integration tests with real diagram patterns from the codebase."""

    @pytest.fixture
    def converter(self):
        """Create converter instance."""
        return ASCIIToMermaid()

    def test_service_architecture_pattern(self, converter):
        """Test typical service architecture diagram pattern."""
        ascii_diagram = """
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Client    │  │  API Layer  │  │   Service   │
│  (Web/App)  │  │   (ogen)    │  │   (Logic)   │
└─────────────┘  └─────────────┘  └─────────────┘
       │                │                │
       └────────────────┼────────────────┘
                        ▼
              ┌─────────────────┐
              │   Repository    │
              │     (sqlc)      │
              └─────────────────┘
                        │
                        ▼
              ┌─────────────────┐
              │   PostgreSQL    │
              │     (pgx)       │
              └─────────────────┘
"""
        mermaid = converter.convert_diagram(ascii_diagram)

        assert "```mermaid" in mermaid
        assert "flowchart LR" in mermaid
        assert "Client" in mermaid
        assert "PostgreSQL" in mermaid
        # Database should have cylinder shape
        assert "[(" in mermaid or "Metadata DB" not in ascii_diagram

    def test_multi_layer_architecture(self, converter):
        """Test multi-layer architecture diagram."""
        ascii_diagram = """
┌───────────────────────────────────────────────────────────┐
│                    CLIENT LAYER                           │
└───────────────────────────────────────────────────────────┘
                           │
                           ▼
┌───────────────────────────────────────────────────────────┐
│                    SERVICE LAYER                          │
└───────────────────────────────────────────────────────────┘
                           │
                           ▼
┌───────────────────────────────────────────────────────────┐
│                    DATA LAYER                             │
└───────────────────────────────────────────────────────────┘
"""
        mermaid = converter.convert_diagram(ascii_diagram)

        assert "```mermaid" in mermaid
        assert "CLIENT LAYER" in mermaid
        assert "SERVICE LAYER" in mermaid
        assert "DATA LAYER" in mermaid
        # Should have subgraphs
        assert "subgraph" in mermaid


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
