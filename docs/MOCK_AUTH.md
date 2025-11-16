# Mock JWT Authentication for OpenTDF MCP Server

## Overview

This document describes the mock JWT authentication system implemented for the `opentdf-mcp` server in the memo-buddy chatmode. This mock implementation demonstrates how agent-based authentication would work in a production environment using OAuth2-issued JWT tokens.

## Architecture

```
┌─────────────────┐
│  memo-buddy     │
│  Chatmode       │
└────────┬────────┘
         │
         │ Configured with Mock JWT
         ▼
┌─────────────────┐
│  MCP Config     │
│  (.vscode/      │
│   mcp.json)     │
└────────┬────────┘
         │
         │ OPENTDF_ACCESS_TOKEN env var
         ▼
┌─────────────────┐
│  opentdf-mcp    │
│  Server         │
└────────┬────────┘
         │
         │ Bearer Token Authentication
         ▼
┌─────────────────┐
│  OpenTDF        │
│  Platform       │
└─────────────────┘
```

## Mock JWT Token

### Token Structure

The mock JWT token is a standard three-part JWT (header.payload.signature):

```
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtZW1vLWJ1ZGR5LWFnZW50IiwiaXNzIjoiaHR0cHM6Ly9hdXRoLm1vY2stb3BlbnRkZi5jb20iLCJhdWQiOiJvcGVudGRmLXBsYXRmb3JtIiwiZXhwIjoyMDAwMDAwMDAwLCJpYXQiOjE3MDAwMDAwMDAsInNjb3BlIjoidGRmOmVuY3J5cHQgdGRmOmRlY3J5cHQgdGRmOmF0dHJpYnV0ZXM6cmVhZCIsImNsaWVudF9pZCI6Im1lbW8tYnVkZHktY2hhdG1vZGUiLCJqdGkiOiJtb2NrLWp3dC10b2tlbi1pZC0xMjM0NSJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

### Decoded Header

```json
{
  "alg": "RS256",
  "typ": "JWT"
}
```

### Decoded Payload

```json
{
  "sub": "memo-buddy-agent",
  "iss": "https://auth.mock-opentdf.com",
  "aud": "opentdf-platform",
  "exp": 2000000000,
  "iat": 1700000000,
  "scope": "tdf:encrypt tdf:decrypt tdf:attributes:read",
  "client_id": "memo-buddy-chatmode",
  "jti": "mock-jwt-token-id-12345"
}
```

#### Payload Claims Explained

- **sub** (Subject): `memo-buddy-agent` - Identifies the agent/principal making requests
- **iss** (Issuer): `https://auth.mock-opentdf.com` - Mock OAuth2 authorization server
- **aud** (Audience): `opentdf-platform` - The intended recipient of the token
- **exp** (Expiration): `2000000000` - Unix timestamp (May 18, 2033) - far future for demo purposes
- **iat** (Issued At): `1700000000` - Unix timestamp (November 14, 2023)
- **scope**: Space-delimited list of permissions granted to the agent:
  - `tdf:encrypt` - Permission to encrypt data
  - `tdf:decrypt` - Permission to decrypt data
  - `tdf:attributes:read` - Permission to list available attributes
- **client_id**: `memo-buddy-chatmode` - Identifies the chatmode/client application
- **jti** (JWT ID): `mock-jwt-token-id-12345` - Unique identifier for this token

### Signature

The signature is a mock value for demonstration purposes. In production, this would be:
- Generated using the RS256 algorithm (RSA + SHA-256)
- Signed with the authorization server's private key
- Verifiable using the server's public key (typically distributed via JWKS endpoint)

## Configuration

### MCP Server Configuration

The JWT token is passed to the `opentdf-mcp` server via environment variables in `.vscode/mcp.json`:

```json
{
  "servers": {
    "opentdf-mcp": {
      "type": "stdio",
      "command": "opentdf-mcp/opentdf-mcp-server",
      "args": [],
      "env": {
        "OPENTDF_PLATFORM_ENDPOINT": "http://localhost:8080",
        "OPENTDF_ACCESS_TOKEN": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
      }
    }
  }
}
```

### Server Implementation Requirements

The `opentdf-mcp` server should be modified to:

1. Check for the `OPENTDF_ACCESS_TOKEN` environment variable on startup
2. If present, use JWT bearer token authentication:
   ```
   Authorization: Bearer <OPENTDF_ACCESS_TOKEN>
   ```
3. If not present, fall back to OAuth2 client credentials flow (backward compatibility)

## Production Deployment

In a production environment, the authentication flow would work as follows:

### 1. User Authorization

```
┌──────┐                                  ┌─────────────┐
│ User │                                  │   OAuth2    │
│      │─────(1) Authorize Agent────────▶│   Server    │
│      │                                  │             │
│      │◀────(2) Consent Screen──────────│             │
│      │                                  │             │
│      │─────(3) Grant Permission────────▶│             │
└──────┘                                  └──────┬──────┘
                                                 │
                                          (4) Issue JWT
                                                 │
                                                 ▼
                                         ┌───────────────┐
                                         │  JWT Token    │
                                         │  (short-lived)│
                                         └───────┬───────┘
                                                 │
                                          (5) Pass to Agent
                                                 │
                                                 ▼
                                         ┌───────────────┐
                                         │  memo-buddy   │
                                         │   Chatmode    │
                                         └───────────────┘
```

### 2. Dynamic Token Injection

Instead of embedding the token in `mcp.json`, production would use:

**Option A: Environment Variable from External Source**
```json
{
  "servers": {
    "opentdf-mcp": {
      "env": {
        "OPENTDF_ACCESS_TOKEN": "${env:OPENTDF_JWT_TOKEN}"
      }
    }
  }
}
```

**Option B: MCP Input (User Prompt)**
```json
{
  "inputs": [
    {
      "type": "promptString",
      "id": "opentdf-jwt-token",
      "description": "OpenTDF JWT Access Token",
      "password": true
    }
  ],
  "servers": {
    "opentdf-mcp": {
      "env": {
        "OPENTDF_ACCESS_TOKEN": "${input:opentdf-jwt-token}"
      }
    }
  }
}
```

**Option C: OAuth2 Flow Integration**
The MCP client would implement OAuth2 authorization code flow:
1. Redirect user to authorization endpoint
2. User consents to agent access
3. Receive authorization code
4. Exchange code for JWT access token
5. Inject token into server environment

### 3. Token Refresh

Production tokens should:
- Have short expiration times (e.g., 1 hour)
- Be automatically refreshed using refresh tokens
- Support token rotation for enhanced security

### 4. Scope-Based Access Control

The OpenTDF platform would validate token scopes:
- Reject requests with insufficient permissions
- Log access attempts for audit trails
- Support fine-grained permissions (e.g., specific data attributes)

## Security Considerations

### Current Mock Implementation

⚠️ **NOT FOR PRODUCTION USE** ⚠️

The current mock implementation:
- Uses a hardcoded, long-lived token
- Has no signature verification
- Grants broad permissions without user consent
- Is committed to version control (visible to all)

### Production Security Requirements

1. **Never commit tokens to version control**
   - Use environment variables or secure vaults
   - Rotate tokens regularly
   - Implement token revocation

2. **Token Validation**
   - Verify signature using public keys from JWKS endpoint
   - Check expiration (`exp` claim)
   - Validate issuer (`iss`) and audience (`aud`)
   - Ensure token hasn't been revoked

3. **Least Privilege**
   - Request minimal scopes needed for operation
   - Implement scope-based authorization
   - Support user-granular permissions

4. **Audit and Monitoring**
   - Log all token usage
   - Monitor for suspicious activity
   - Alert on token misuse

5. **Transport Security**
   - Always use HTTPS/TLS
   - Implement certificate pinning where appropriate
   - Secure token storage (encrypted at rest)

## Testing

To test the mock authentication:

1. Ensure the `opentdf-mcp` server is running with the mock token
2. Use the memo-buddy chatmode
3. Invoke encryption/decryption operations
4. Verify that requests include the JWT in the Authorization header

### Example: Inspecting Token Usage

If the `opentdf-mcp` server logs requests, you should see:

```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Future Enhancements

1. **Token Management UI**
   - Web interface for token issuance
   - User consent flow
   - Token lifecycle management

2. **Multi-Tenancy Support**
   - Per-user tokens with different permissions
   - Organization-level access control
   - Role-based access control (RBAC)

3. **Token Introspection**
   - Real-time token validation endpoint
   - Revocation checking
   - Claims enrichment

4. **Federated Identity**
   - Support multiple identity providers
   - SAML/OIDC integration
   - SSO capabilities

## References

- [RFC 7519: JSON Web Token (JWT)](https://tools.ietf.org/html/rfc7519)
- [RFC 6749: OAuth 2.0 Authorization Framework](https://tools.ietf.org/html/rfc6749)
- [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)
- [OpenTDF Documentation](https://github.com/opentdf/platform)

## Contact

For questions or issues regarding the authentication system, please file an issue in the project repository.
