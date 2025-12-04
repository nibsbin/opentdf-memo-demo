#!/usr/bin/env python3
"""
USAF Memo MCP Server

This server provides tools and resources for creating USAF memos using the
usaf_memo Quill template and Quillmark rendering engine.
"""

import os
import json
from pathlib import Path
from datetime import datetime

import quillmark
from mcp.server import FastMCP

# Initialize the MCP server
mcp = FastMCP(
    "usaf-memo",
    instructions="""
This server helps you create USAF memos using the usaf_memo Quill template.

Available resources:
- memo://schema - Get the field schema for USAF memos
- memo://example - Get an example with authoritative usage information for USAF memo markdown

Available tools:
- render_memo_to_pdf - Render a USAF memo markdown file to PDF
""",
)

# Load the usaf_memo quill
QUILL_PATH = Path(__file__).parent / "usaf_memo"
OUTPUT_DIR = Path(__file__).parent / "output"

# Create output directory if it doesn't exist
OUTPUT_DIR.mkdir(exist_ok=True)

# Load and register the quill
qm = quillmark.Quillmark()
quill = quillmark.Quill.from_path(str(QUILL_PATH))
qm.register_quill(quill)

@mcp.tool()
@mcp.resource("memo://schema")
def get_memo_schema() -> str:
    """
    Get the field schema for USAF memos.
    
    Returns a JSON object describing all available fields, their types,
    descriptions, and defaults.
    """
    return json.dumps(quill.schema, indent=2)

@mcp.tool()
@mcp.resource("memo://example")
def get_memo_example() -> str:
    """
    Get an example USAF memo in markdown format.
    
    This example demonstrates how to structure a USAF memo with frontmatter
    and content.
    """
    return quill.example or "No example available"

@mcp.tool()
@mcp.resource("memo://description")
def get_memo_description() -> str:
    """
    Get a description of the usaf_memo Quill template.
    
    Returns metadata about the Quill including its name, backend, and
    description.
    """
    info = {
        "name": quill.name,
        "backend": quill.backend,
        "description": quill.metadata.get("description", ""),
    }
    return json.dumps(info, indent=2)


@mcp.tool()
def render_memo_to_pdf(markdown_file_path: str) -> str:
    """
    Render a USAF memo markdown file to PDF.
    
    Takes a markdown file path as input, renders it using the usaf_memo Quill,
    and saves the resulting PDF to the output directory.
    
    Args:
        markdown_file_path: Path to the markdown file to render
        
    Returns:
        Path to the rendered PDF file
        
    Raises:
        FileNotFoundError: If the markdown file doesn't exist
        ValueError: If the markdown file doesn't have the correct QUILL tag
        Exception: If rendering fails
    """
    # Validate input file exists
    input_path = Path(markdown_file_path).resolve()
    if not input_path.exists():
        raise FileNotFoundError(f"Markdown file not found: {markdown_file_path}")
    
    # Read the markdown content
    with open(input_path, "r", encoding="utf-8") as f:
        markdown_content = f.read()
    
    # Parse the markdown
    try:
        parsed = quillmark.ParsedDocument.from_markdown(markdown_content)
    except Exception as e:
        raise ValueError(f"Failed to parse markdown: {e}")
    
    # Check if the document has the correct quill tag
    quill_tag = parsed.quill_tag()
    if quill_tag != "usaf_memo":
        raise ValueError(
            f"Markdown file must have 'QUILL: usaf_memo' in frontmatter, "
            f"but found 'QUILL: {quill_tag}'"
        )
    
    # Create workflow and render (quillmark v0.9.0 API)
    try:
        workflow = qm.workflow(quill_tag)
        result = workflow.render(parsed, quillmark.OutputFormat.PDF)
    except Exception as e:
        raise Exception(f"Failed to render PDF: {e}")
    
    # Check for warnings
    if result.warnings:
        warnings_msg = "\n".join(str(w) for w in result.warnings)
        print(f"Rendering warnings:\n{warnings_msg}")
    
    # Save the PDF
    if not result.artifacts:
        raise Exception("No artifacts generated during rendering")
    
    artifact = result.artifacts[0]
    
    # Generate output filename
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    output_filename = f"memo_{timestamp}.pdf"
    output_path = OUTPUT_DIR / output_filename
    
    # Save the artifact
    artifact.save(str(output_path))
    
    return str(output_path.resolve())


if __name__ == "__main__":
    # Run the server using stdio transport
    import asyncio
    asyncio.run(mcp.run_stdio_async())
