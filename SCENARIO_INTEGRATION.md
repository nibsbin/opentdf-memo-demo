# Scenario Integration

Technical plan for ABAC [scenario](DEMO.md) integration with OpenTDF platform and Keycloak.

## Chatmode Prompts

Create base prompt template with placeholders for user data:
- Name
- JWT (hardcoded)

## Keycloak Config

Configure users from [scenario](DEMO.md). Apply attributes to users.

## MCP Tools

### opentdf-mcp
- `encrypt`: Encrypt data with optional attributes (supports TDF/nanoTDF)
- `decrypt`: Decrypt TDF/nanoTDF files (auto-detects format)
- `list_attributes`: List available data attributes
- New tool: Use `GetDecisionBulk` to test entitlements against files and send to LLM.


### memo-mcp
- `render_memo_to_pdf`: Render markdown memo to PDF
- `get_memo_schema`: Retrieve the schema for memo markdown frontmatter
- `get_memo_example`: Retrieve an example memo markdown file
- `get_memo_description`: Get a description of the usaf_memo Quill template