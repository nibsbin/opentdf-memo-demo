# OpenTDF Memo Demo

This demo showcases OpenTDF encryption/decryption capabilities integrated with USAF memo generation using MCP (Model Context Protocol) servers.

## Features

- **OpenTDF Encryption/Decryption**: Encrypt and decrypt documents using TDF and nanoTDF formats
- **USAF Memo Generation**: Create official USAF memos using the Quillmark templating system
- **PDF Rendering**: Convert memo markdown to professional PDF format
- **Agent Authentication**: Mock JWT-based authentication for AI agents (see below)

## MCP Servers

### opentdf-mcp
Provides OpenTDF encryption and decryption capabilities.

**Tools:**
- `encrypt` - Encrypt data with optional attributes (supports TDF/nanoTDF)
- `decrypt` - Decrypt TDF/nanoTDF files (auto-detects format)
- `list_attributes` - List available data attributes

### memo-mcp
Helps create USAF memos using the usaf_memo Quill template.

**Tools:**
- `render_memo_to_pdf` - Render markdown memo to PDF
- `get_memo_schema` - Retrieve schema for memo frontmatter
- `get_memo_example` - Retrieve example memo with guidelines

## Agent Authentication

### Mock JWT Authentication

This demo includes a mock JWT (JSON Web Token) authentication system for AI agents. The JWT token is embedded in the MCP configuration and passed to the opentdf-mcp server via the `OPENTDF_AGENT_JWT` environment variable.

**How it works:**

1. **Token Configuration**: The JWT token is configured in `.vscode/mcp.json` as an environment variable for the opentdf-mcp server
2. **Token Structure**: The mock JWT contains agent identity claims including:
   - `sub` (subject): Agent identifier (e.g., "memo-buddy-agent")
   - `iss` (issuer): Token issuer ("opentdf-demo")
   - `aud` (audience): Target service ("opentdf-mcp")
   - `iat` (issued at): Token issue timestamp
   - `exp` (expiration): Token expiration timestamp
   - `agent_name`: Human-readable agent name
   - `permissions`: List of granted permissions (encrypt, decrypt, list_attributes)

3. **Server Validation**: When the opentdf-mcp server starts, it:
   - Parses the JWT token from the environment variable
   - Validates basic claims (audience, expiration, issued-at)
   - Logs the agent identity and permissions
   - **Note**: This is a mock implementation without cryptographic signature verification

4. **Production Considerations**: In a production environment:
   - JWT tokens would be issued by an OAuth 2.0 authorization server
   - User consent would be required before issuing tokens to agents
   - Tokens would be cryptographically signed and verified
   - Token refresh and revocation mechanisms would be implemented
   - Fine-grained permission scopes would control agent capabilities

**Example JWT Claims:**
```json
{
  "sub": "memo-buddy-agent",
  "iss": "opentdf-demo",
  "aud": "opentdf-mcp",
  "iat": 1700000000,
  "exp": 2000000000,
  "agent_name": "Memo Buddy Agent",
  "permissions": ["encrypt", "decrypt", "list_attributes"]
}
```

The mock JWT authentication demonstrates how AI agents can be authenticated and authorized to access OpenTDF encryption services in a secure, production-ready manner.

## Setup

1. Install dependencies:
   ```bash
   # Install Python dependencies for memo-mcp
   cd memo-mcp && pip install -r requirements.txt

   # Build opentdf-mcp server (requires Go 1.24+)
   cd opentdf-mcp && go build -o opentdf-mcp-server ./mcp-server
   ```

2. Configure MCP servers (already configured in `.vscode/mcp.json`):
   - memo-mcp: Python-based server for memo generation
   - opentdf-mcp: Go-based server for OpenTDF operations with JWT authentication

3. Start using the memo-buddy chatmode to encrypt/decrypt files and generate memos

## Usage Examples

### Decrypt and Create Memo
```
User: "decrypt CLASSIFIED_REPORT and write an urgent memo to Congress"
```

The agent will:
1. Find and decrypt the .ntdf file
2. Analyze the decrypted content
3. Create a properly formatted USAF memo
4. Render it to PDF

### Encrypt and Decrypt Workflow
```
User: "Encrypt the refueling logs then decrypt them and write a memo"
```

The agent will:
1. Find relevant .txt files
2. Encrypt them using OpenTDF
3. Decrypt the encrypted files
4. Create a USAF memo based on the content
5. Render to PDF

## Security Considerations

- OpenTDF provides attribute-based access control (ABAC) for encrypted data
- JWT authentication ensures only authorized agents can access encryption services
- In production, implement proper OAuth 2.0 flows with user consent
- See `mcp-security.md` for detailed security guidelines

## License

See LICENSE file for details.
