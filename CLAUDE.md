---
description: 'Assistant for writing and rendering USAF official memos using Quillmark.'
---

# OpenTDF Memo Demo

## Project Overview

This demo showcases OpenTDF encryption/decryption capabilities integrated with USAF memo generation. The project demonstrates how to:
- Encrypt documents with OpenTDF (TDF/nanoTDF formats)
- Decrypt OpenTDF-encrypted documents (TDF/nanoTDF formats)
- Generate official USAF memos using the usaf_memo Quill template
- Render memos to PDF format

## Available MCP Servers

### opentdf-mcp
Provides OpenTDF encryption and decryption capabilities.

**Tools:**
- `mcp__opentdf-mcp__encrypt` - Encrypt data with optional attributes (supports TDF/nanoTDF)
  - Use `input` parameter for file paths: encrypts the file contents
  - Use `data` parameter for literal text: encrypts the text directly
  - Parameters are mutually exclusive (specify one or the other)
- `mcp__opentdf-mcp__decrypt` - Decrypt TDF/nanoTDF files (auto-detects format)
- `mcp__opentdf-mcp__list_attributes` - List available data attributes

### memo-mcp
Helps create USAF memos using the usaf_memo Quill template.

**Tools:**
- `mcp__memo-mcp__render_memo_to_pdf` - Render markdown memo to PDF
- `mcp__memo-mcp__get_memo_schema` - Retrieve the schema for memo markdown frontmatter, ensuring proper structure and compliance.
- `mcp__memo-mcp__get_memo_example` - Retrieve an example memo markdown file with authoritative usage guidelines.

## Workflows

### Create USAF Memo from Encrypted File on the Contractor Scenario
1. Find encrypted file: e.g. `Glob` pattern `**/*.ntdf`
1. Decrypt: `mcp__opentdf-mcp__decrypt(input: "/path/to/file.ntdf")`
1. Get memo usage context: 
  - `mcp__memo-mcp__get_memo_schema()`. This ensures your memo frontmatter is structured correctly.
  - `mcp__memo-mcp__get_memo_example()`. This is an authoritative reference with writing guidelines.
1. Write memo markdown to a file saved in `drafts/`.
  - Do not carry over formatting from the original document.
  - FOLLOW THE INSTRUCTIONS IN THE EXAMPLE.
1. Render to PDF: `mcp__memo-mcp__render_memo_to_pdf(markdown_file_path: "...")`
1. Brief summary. Display links to markdown draft and PDF output. Offer the user to workshop/iterate.

**Example Task:**
```
User: "decrypt CLASSIFIED_REPORT and write an urgent memo to Congress"
→ Find .ntdf file
→ Decrypt with opentdf-mcp
→ Analyze decrypted content
- Read memo usage, schema, and example
→ Create USAF memo markdown
→ Render to PDF
```

### Create USAF Memo from Encrypted File on the Refueling Scenario
1. Find plain text files: e.g. `Glob` pattern `**/*.txt`
1. Encrypt each file to the associated encrypted folder: `mcp__opentdf-mcp__encrypt(input: "/path/to/file.txt")`
1. Find each encrypted file: e.g. `Glob` pattern `**/*.ntdf`
1. Decrypt `mcp__opentdf-mcp__decrypt(input: "/path/to/file.ntdf")` each file to the associated decrypted folder
1. Get memo usage context: 
  - `mcp__memo-mcp__get_memo_schema()`. This ensures your memo frontmatter is structured correctly.
  - `mcp__memo-mcp__get_memo_example()`. This is an authoritative reference with writing guidelines.
1. Write memo markdown to a file saved in `drafts/` by reading the decrypted content (not the original).
  - Do not carry over formatting from the original document.
  - FOLLOW THE INSTRUCTIONS IN THE EXAMPLE.
1. Render to PDF: `mcp__memo-mcp__render_memo_to_pdf(markdown_file_path: "...")`
1. Brief summary. Display links to markdown draft and PDF output. Offer the user to workshop/iterate.

**Example Task:**
```
User: "Encrypt the relevant logs then decrypt them and write a memo to the Commander of United States Indo-Pacific Command about the refueling operation"
→ Find .txt files
→ Encrypt each file to the associated encrypted folder: `mcp__opentdf-mcp__encrypt(input: "/path/to/file.txt")`
→ Find each encrypted file: e.g. `Glob` pattern `**/*.ntdf`
→ Decrypt each file to the associated decrypted folder: `mcp__opentdf-mcp__decrypt(input: "/path/to/file.ntdf")`
→ Analyze decrypted content
- Read memo usage, schema, and example
→ Create USAF memo markdown
→ Render to PDF
```

## Best Practices

### OpenTDF Encryption/Decryption
- **Encrypting files**: Use `input` parameter with file path
  - Example: `mcp__opentdf-mcp__encrypt(input: "/path/to/file.txt", format: "nano")`
  - This reads and encrypts the file contents (not the path string)
- **Encrypting literal text**: Use `data` parameter with text
  - Example: `mcp__opentdf-mcp__encrypt(data: "Secret message", format: "nano")`
- Prefer nanoTDF format (`format: "nano"`) for better compatibility
- Decrypt tool auto-detects TDF vs nanoTDF format
- Check available attributes with `list_attributes` before encrypting with attributes

### File Operations
- Check if files exist with Glob before attempting operations
- Use absolute paths for all file operations
- Read files before editing/writing to ensure correct handling
- Create markdown files in `drafts/` before rendering to pdf. Do not read other drafts.

### Writing Guidelines
- Specify `QUILL: usaf_memo` in frontmatter to use the USAF memo template.
- Follow the guidelines in the example memo retrieved via `mcp__memo-mcp__get_memo_example()`.
- If you are deriving from a classified source:
    - Add "//FICTIONAL" to classification banner. e.g. "SECRET//NOFORN//FICTIONAL"
    - Add portion markings at the beginning of each paragraph; e.g. "(U)", "(S)", etc.
