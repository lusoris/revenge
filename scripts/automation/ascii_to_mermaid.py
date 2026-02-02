#!/usr/bin/env python3
"""ASCII to Mermaid Converter.

Converts ASCII box diagrams in YAML files to Mermaid syntax.
Handles architecture_diagram and data_flow_diagram fields.

Author: Automation System
Created: 2026-02-02
"""

import re
import sys
from pathlib import Path

import yaml


def str_representer(dumper, data):
    """Custom representer for multiline strings using literal block style."""
    if "\n" in data:
        return dumper.represent_scalar("tag:yaml.org,2002:str", data, style="|")
    return dumper.represent_scalar("tag:yaml.org,2002:str", data)


# Register the custom representer
yaml.add_representer(str, str_representer)


class ASCIIToMermaid:
    """Convert ASCII diagrams to Mermaid syntax."""

    def __init__(self):
        """Initialize converter."""
        # Patterns for ASCII elements
        self.box_pattern = re.compile(
            r"┌[─]+┐\s*\n"  # Top border
            r"│\s*(.+?)\s*│\s*\n"  # Content (capture group)
            r"(?:│\s*(.+?)\s*│\s*\n)*"  # Optional additional lines
            r"└[─]+┘",  # Bottom border
            re.MULTILINE,
        )
        # Simpler box pattern - just find box content
        self.simple_box = re.compile(r"│\s*([^│\n]+?)\s*│")
        # Arrow patterns
        self.arrow_down = re.compile(r"[│▼↓]")
        self.arrow_right = re.compile(r"[─▶→►]")
        self.arrow_left = re.compile(r"[◀←]")

        # Keywords for node type detection
        self.layer_keywords = {"LAYER", "TIER", "LEVEL"}
        self.database_keywords = {
            "PostgreSQL",
            "Database",
            "DB",
            "Cache",
            "Redis",
            "Storage",
        }
        self.service_keywords = {
            "Service",
            "API",
            "Handler",
            "Engine",
            "Manager",
            "Controller",
        }
        self.external_keywords = {
            "External",
            "Third-party",
            "Provider",
            "Client",
            "Web",
            "App",
        }

    def extract_boxes(self, ascii_text: str) -> list[dict]:
        """Extract box labels from ASCII diagram.

        Handles both vertical and horizontal box arrangements.
        Properly extracts multi-line text from within box borders.

        Args:
            ascii_text: ASCII diagram text

        Returns:
            List of dicts with box info
        """
        boxes = []
        lines = ascii_text.split("\n")
        box_id = 0

        # Box-drawing characters to filter from text
        box_chars = "┌┐└┘├┤┬┴┼─━═║╔╗╚╝╠╣╦╩╬"
        arrow_chars = "▶◀▼▲►◄→←↓↑"

        def clean_text(text: str) -> str:
            """Clean text by removing box/arrow characters."""
            for c in box_chars + arrow_chars + "│":
                text = text.replace(c, " ")
            return " ".join(text.split()).strip()

        def find_box_at(row: int, col: int) -> dict | None:
            """Find a complete box starting at (row, col).

            Returns box info or None if not a valid box start.
            """
            if row >= len(lines):
                return None
            line = lines[row]
            if col >= len(line) or line[col] != "┌":
                return None

            # Find the matching ┐ on the same line
            # Must handle the exact box, not grab text from adjacent boxes
            end_col = col + 1
            depth = 1
            while end_col < len(line) and depth > 0:
                if line[end_col] == "┌":
                    depth += 1
                elif line[end_col] == "┐":
                    depth -= 1
                end_col += 1

            if depth != 0:
                return None
            end_col -= 1  # Point to the ┐

            # Now find the bottom row - look for └ at col position
            bottom_row = row + 1
            while bottom_row < len(lines):
                bottom_line = lines[bottom_row]
                if len(bottom_line) > col and bottom_line[col] == "└":
                    break
                bottom_row += 1

            if bottom_row >= len(lines):
                return None

            # Verify the box corners match
            bottom_line = lines[bottom_row]
            if len(bottom_line) <= end_col or bottom_line[end_col] != "┘":
                # Try to find the actual ┘ near end_col
                for ec in range(
                    max(col, end_col - 2), min(len(bottom_line), end_col + 3)
                ):
                    if bottom_line[ec] == "┘":
                        end_col = ec
                        break

            # Extract text from content rows
            content_texts = []
            for r in range(row + 1, bottom_row):
                if r >= len(lines):
                    break
                content_line = lines[r]

                # Find text between │ characters within our column range
                if len(content_line) > col:
                    # Look for │ at start and end of box
                    segment = content_line[col : min(end_col + 1, len(content_line))]

                    # Check if this line has the vertical borders
                    if "│" in segment:
                        # Extract text between │ characters
                        first_pipe = segment.find("│")
                        last_pipe = segment.rfind("│")
                        if first_pipe != last_pipe:
                            inner_text = segment[first_pipe + 1 : last_pipe]
                            cleaned = clean_text(inner_text)
                            if cleaned:
                                content_texts.append(cleaned)
                        elif first_pipe >= 0:
                            # Single pipe or edge - try to extract any text
                            inner_text = segment.replace("│", " ")
                            cleaned = clean_text(inner_text)
                            if cleaned:
                                content_texts.append(cleaned)

            if content_texts:
                return {
                    "row": row,
                    "col": col,
                    "end_col": end_col,
                    "bottom_row": bottom_row,
                    "texts": content_texts,
                }
            return None

        # Scan for all box starts
        for row_idx, line in enumerate(lines):
            col = 0
            while col < len(line):
                if line[col] == "┌":
                    box_info = find_box_at(row_idx, col)
                    if box_info:
                        box_id += 1
                        boxes.append(
                            {
                                "id": f"node{box_id}",
                                "label": box_info["texts"][0],
                                "sublabels": (
                                    box_info["texts"][1:]
                                    if len(box_info["texts"]) > 1
                                    else []
                                ),
                                "row": box_info["row"],
                                "col": box_info["col"],
                                "bottom_row": box_info["bottom_row"],
                            }
                        )
                        # Skip to end of this box
                        col = box_info["end_col"] + 1
                    else:
                        col += 1
                else:
                    col += 1

        return boxes

    def infer_connections(
        self, boxes: list[dict], _ascii_text: str
    ) -> list[tuple[str, str, str]]:
        """Infer connections between boxes based on position and arrows.

        Args:
            boxes: List of box dicts
            ascii_text: Original ASCII text

        Returns:
            List of (from_id, to_id, label) tuples
        """
        connections = []

        if not boxes:
            return connections

        # Group boxes by row
        rows = {}
        for box in boxes:
            row = box.get("row", 0)
            if row not in rows:
                rows[row] = []
            rows[row].append(box)

        # Sort boxes in each row by column
        for row_boxes in rows.values():
            row_boxes.sort(key=lambda b: b.get("col", 0))

        # Connect boxes in same row horizontally (left to right)
        for row_boxes in rows.values():
            for i in range(len(row_boxes) - 1):
                from_box = row_boxes[i]
                to_box = row_boxes[i + 1]
                connections.append((from_box["id"], to_box["id"], ""))

        # Connect last box of each row to first box of next row
        sorted_rows = sorted(rows.keys())
        for i in range(len(sorted_rows) - 1):
            curr_row = sorted_rows[i]
            next_row = sorted_rows[i + 1]

            # Connect from rightmost of current to leftmost of next
            if rows[curr_row] and rows[next_row]:
                from_box = rows[curr_row][-1]
                to_box = rows[next_row][0]
                connections.append((from_box["id"], to_box["id"], ""))

        return connections

    def _is_layer_box(self, label: str) -> bool:
        """Check if a box represents a layer/tier heading."""
        upper = label.upper()
        return any(kw in upper for kw in self.layer_keywords)

    def _get_node_shape(self, label: str, sublabels: list[str]) -> tuple[str, str]:
        """Determine the Mermaid node shape based on content.

        Returns:
            Tuple of (open_bracket, close_bracket)
        """
        combined = f"{label} {' '.join(sublabels)}".upper()

        # Database/storage nodes - cylinder shape
        if any(kw.upper() in combined for kw in self.database_keywords):
            return "[(", ")]"

        # External/client nodes - stadium shape
        if any(kw.upper() in combined for kw in self.external_keywords):
            return "([", "])"

        # Service nodes - subroutine shape
        if any(kw.upper() in combined for kw in self.service_keywords):
            return "[[", "]]"

        # Default - rectangle
        return '["', '"]'

    def generate_mermaid(
        self,
        boxes: list[dict],
        connections: list[tuple[str, str, str]],
        diagram_type: str = "flowchart",
    ) -> str:
        """Generate Mermaid diagram from boxes and connections.

        Uses different node shapes based on content type for visual distinction.
        Layer boxes get a special banner shape for clear section headings.

        Args:
            boxes: List of box dicts
            connections: List of connection tuples
            diagram_type: Type of diagram (flowchart/graph)

        Returns:
            Mermaid diagram string
        """
        if not boxes:
            return ""

        lines = [f"{diagram_type} TD"]

        # Track layer boxes for styling
        layer_box_ids = []

        # Add all nodes with appropriate shapes
        for box in boxes:
            label = box["label"].replace('"', "'").replace("(", "[").replace(")", "]")

            # Add sublabels if present
            if box.get("sublabels"):
                sublabels = "<br/>".join(box["sublabels"][:2])
                label = f"{label}<br/>{sublabels}"

            # Check if this is a layer/tier box - use banner shape
            if self._is_layer_box(box["label"]):
                layer_box_ids.append(box["id"])
                # Use banner shape for layers: >label]
                lines.append(f'    {box["id"]}["{label}"]')
            else:
                # Use appropriate shape based on content
                open_br, close_br = self._get_node_shape(
                    box["label"], box.get("sublabels", [])
                )
                lines.append(f"    {box['id']}{open_br}{label}{close_br}")

        # Add connections
        for from_id, to_id, conn_label in connections:
            if conn_label:
                lines.append(f"    {from_id} -->|{conn_label}| {to_id}")
            else:
                lines.append(f"    {from_id} --> {to_id}")

        # Add styling for layer boxes to make them stand out
        if layer_box_ids:
            lines.append("")
            lines.append("    %% Layer styling")
            colors = ["#1976D2", "#388E3C", "#7B1FA2", "#F57C00", "#C2185B"]
            for i, layer_id in enumerate(layer_box_ids):
                color = colors[i % len(colors)]
                lines.append(
                    f"    style {layer_id} fill:{color},stroke:#fff,"
                    f"stroke-width:2px,color:#fff"
                )

        return "\n".join(lines)

    def convert_diagram(self, ascii_text: str) -> str:
        """Convert ASCII diagram to Mermaid.

        Args:
            ascii_text: ASCII diagram text

        Returns:
            Mermaid diagram string
        """
        # Extract boxes
        boxes = self.extract_boxes(ascii_text)

        if not boxes:
            # No boxes found, return original
            return ascii_text

        # Infer connections
        connections = self.infer_connections(boxes, ascii_text)

        # Generate Mermaid
        mermaid = self.generate_mermaid(boxes, connections)

        # Wrap in code block
        return f"```mermaid\n{mermaid}\n```"

    def has_ascii_diagram(self, text: str) -> bool:
        """Check if text contains ASCII box diagram.

        Args:
            text: Text to check

        Returns:
            True if contains ASCII diagram
        """
        # Check for box drawing characters
        box_chars = {"┌", "┐", "└", "┘", "│", "─", "├", "┤", "┬", "┴", "┼"}
        return any(c in text for c in box_chars)

    def process_yaml_field(self, value: str) -> tuple[str, bool]:
        """Process a YAML field that may contain ASCII diagram.

        Args:
            value: Field value

        Returns:
            Tuple of (processed_value, was_converted)
        """
        if not isinstance(value, str):
            return value, False

        if not self.has_ascii_diagram(value):
            return value, False

        # Check if already has mermaid
        if "```mermaid" in value:
            return value, False

        # Extract text before/after the code block
        # Pattern: text before ``` ... ``` text after
        code_block_pattern = re.compile(r"(.*?)```\n?(.*?)```(.*)", re.DOTALL)
        match = code_block_pattern.search(value)

        if match:
            before = match.group(1).strip()
            diagram = match.group(2)
            after = match.group(3).strip()

            # Convert the ASCII diagram
            mermaid = self.convert_diagram(diagram)

            # Reconstruct
            parts = []
            if before:
                parts.append(before)
            parts.append(mermaid)
            if after:
                parts.append(after)

            return "\n\n".join(parts), True

        # No code block, try converting the whole thing
        converted = self.convert_diagram(value)
        return converted, converted != value


def process_yaml_file(
    file_path: Path, converter: ASCIIToMermaid, dry_run: bool = True
) -> dict:
    """Process a YAML file and convert ASCII diagrams to Mermaid.

    Args:
        file_path: Path to YAML file
        converter: ASCIIToMermaid instance
        dry_run: If True, don't write changes

    Returns:
        Stats dict
    """
    stats = {"converted": 0, "fields": []}

    # Read YAML
    with open(file_path) as f:
        content = f.read()

    # Parse YAML
    try:
        data = yaml.safe_load(content)
    except yaml.YAMLError as e:
        print(f"  ❌ YAML parse error: {e}")
        return stats

    if not isinstance(data, dict):
        return stats

    # Fields that may contain diagrams
    diagram_fields = [
        "architecture_diagram",
        "data_flow_diagram",
        "component_description",
        "database_schema",
    ]

    modified = False

    for field in diagram_fields:
        if field in data and isinstance(data[field], str):
            original = data[field]
            converted, was_converted = converter.process_yaml_field(original)

            if was_converted:
                data[field] = converted
                stats["converted"] += 1
                stats["fields"].append(field)
                modified = True

    # Write back if modified
    if modified and not dry_run:
        with open(file_path, "w") as f:
            yaml.dump(
                data,
                f,
                default_flow_style=False,
                allow_unicode=True,
                sort_keys=False,
                width=120,
            )

    return stats


def main():
    """Main entry point."""
    args = sys.argv[1:]

    if "--help" in args or "-h" in args:
        print("Usage: python ascii_to_mermaid.py [path] [--live]")
        print()
        print("Arguments:")
        print("  path       Path to YAML file or directory (default: data/)")
        print()
        print("Options:")
        print("  --live     Apply changes (default: dry-run)")
        print("  --dry-run  Show what would be changed (default)")
        print()
        print("Examples:")
        print("  python ascii_to_mermaid.py data/ --dry-run")
        print("  python ascii_to_mermaid.py data/services/AUTH.yaml --live")
        sys.exit(0)

    # Parse args
    path_arg = "data/"
    dry_run = True

    for arg in args:
        if arg == "--live":
            dry_run = False
        elif arg == "--dry-run":
            dry_run = True
        elif not arg.startswith("-"):
            path_arg = arg

    path = Path(path_arg)

    if not path.exists():
        print(f"❌ Path not found: {path}")
        sys.exit(1)

    # Initialize converter
    converter = ASCIIToMermaid()

    print(f"\n{'=' * 70}")
    print(f"ASCII TO MERMAID CONVERTER - {'DRY RUN' if dry_run else 'LIVE'}")
    print(f"{'=' * 70}\n")

    total_stats = {"files_processed": 0, "files_converted": 0, "total_fields": 0}

    if path.is_file():
        # Process single file
        print(f"Processing: {path}")
        stats = process_yaml_file(path, converter, dry_run)

        if stats["converted"] > 0:
            total_stats["files_converted"] += 1
            total_stats["total_fields"] += stats["converted"]
            if dry_run:
                print(f"  Would convert: {', '.join(stats['fields'])}")
            else:
                print(f"  ✓ Converted: {', '.join(stats['fields'])}")
        else:
            print("  No diagrams found or already converted")

        total_stats["files_processed"] = 1

    elif path.is_dir():
        # Process directory
        yaml_files = list(path.glob("**/*.yaml"))
        print(f"Processing: {len(yaml_files)} YAML files in {path}\n")

        for yaml_file in sorted(yaml_files):
            total_stats["files_processed"] += 1
            stats = process_yaml_file(yaml_file, converter, dry_run)

            if stats["converted"] > 0:
                total_stats["files_converted"] += 1
                total_stats["total_fields"] += stats["converted"]
                rel_path = yaml_file.relative_to(path)
                fields = ", ".join(stats["fields"])
                if dry_run:
                    print(f"  Would convert: {rel_path} ({fields})")
                else:
                    print(f"✓ Converted: {rel_path} ({fields})")

    print(f"\n{'=' * 70}")
    print("SUMMARY")
    print(f"{'=' * 70}")
    print(f"Files processed: {total_stats['files_processed']}")
    print(f"Files with conversions: {total_stats['files_converted']}")
    print(f"Total fields converted: {total_stats['total_fields']}")

    if dry_run:
        print("\n⚠️  DRY RUN MODE - No changes written")
        print("Run with --live to apply changes")

    print(f"{'=' * 70}\n")


if __name__ == "__main__":
    main()
