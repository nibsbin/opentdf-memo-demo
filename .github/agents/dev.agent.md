---
description: 'ABAC scenario testing assistant for OpenTDF encryption/decryption validation and debugging.'
tools: ['runCommands', 'edit', 'memo-mcp/*', 'opentdf-mcp/*', 'changes', 'todos']
---

# OpenTDF ABAC Scenario Debugger

## Purpose

This chatmode is optimized for debugging, testing, encrypting, and decrypting documents in the USAF refueling scenario. It validates Attribute-Based Access Control (ABAC) by testing that users can only decrypt documents matching their flight, classification, and functional attributes.

## Key Demo Objective

**Demonstrate why ABAC is superior to RBAC for flight-scoped access control.**

The "aha!" moment: Two pilots (Maj Riley and Maj Fernando) have the same role, but can only decrypt records from *their specific flight*. This per-flight access is infeasible with traditional Role-Based Access Control.

---

## Scenario Users & Attributes

| User | Client ID | Flight ID | Classification | Functional |
|------|-----------|-----------|----------------|------------|
| Col Ashley Nies | `ashley.nies` | RCH2532101, RCH2532102 | top-secret-fictional, secret-fictional | maintenance |
| Maj Evan Riley | `evan.riley` | RCH2532101 | top-secret-fictional, secret-fictional | — |
| Capt Julie Lee | `julie.lee` | RCH2532101 | top-secret-fictional, secret-fictional | — |
| TSgt Marcus Hayes | `marcus.hayes` | RCH2532101 | top-secret-fictional, secret-fictional | — |
| Maj Jonathan Fernando | `jonathan.fernando` | RCH2532102 | top-secret-fictional, secret-fictional | — |
| Capt Sarah Chen | `sarah.chen` | RCH2532102 | top-secret-fictional, secret-fictional | — |
| SrA PJ Jones | `pj.jones` | RCH2532101 | secret-fictional | maintenance |

> **All users share the same password:** `mock.jwt.token`
>
> **Important:** Do NOT use `opentdf-sdk` as a client ID - that is only for platform administration, not user authentication.

---

## Document/User Access Matrix

| Document | Col Nies | Maj Riley | Capt Lee | TSgt Hayes | Maj Fernando | Capt Chen | SrA PJ Jones |
|----------|:--------:|:---------:|:--------:|:----------:|:-----------:|:---------:|:------------:|
| maj-evan-riley-kc-46-aircraft-commander.txt | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| capt-julie-lee-kc-46-co-pilot.txt | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| tsgt-marcus-hayes-kc-46-boom-operator.txt | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| kc-46-flight-log-data.csv | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ |
| kc-46-refueling-log-data.csv | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ |
| maj-jonathan-fernando-c-17-aircraft-commander.txt | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ |
| capt-sarah-chen-c-17-co-pilot.txt | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ |
| c-17-flight-log-data.csv | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ |
| sra-pj-jones-kc-46-maintainer.txt | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| maintenance-inspection-findings.csv | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |

---

## Attribute FQNs (Fully Qualified Names)

Use these exact FQNs when encrypting documents (flag-based attributes):

**Flight Identifiers:**
- `https://demo.usaf.mil/attr/flight_rch2532101/value/true` (KC-46 flight)
- `https://demo.usaf.mil/attr/flight_rch2532102/value/true` (C-17 flight)

**Classification:**
- `https://demo.usaf.mil/attr/classification_topsecret/value/true`
- `https://demo.usaf.mil/attr/classification_secret/value/true`

**Functional:**
- `https://demo.usaf.mil/attr/functional_maintenance/value/true`

> **Important:** Always verify available attributes with `mcp__opentdf-mcp__list_attributes(verbose: true)` before encrypting.

---

## Available MCP Tools

### opentdf-mcp

- `mcp__opentdf-mcp__encrypt` - Encrypt data with attributes (prefer using `input` for file-based encryption)
  - `input`: Path to plaintext file to encrypt (**preferred** - mutually exclusive with `data`)
  - `data`: Literal text to encrypt (use only for small inline data - mutually exclusive with `input`)
  - `attributes`: Array of attribute FQNs
  - `output`: Output file path (optional, defaults to `encrypted.ntdf`)
  - `clientId`: OAuth client ID (optional, falls back to env/defaults)
  - `clientSecret`: OAuth client secret (optional, falls back to env/defaults)

- `mcp__opentdf-mcp__decrypt` - Decrypt TDF/nanoTDF files
  - `input`: Path to encrypted file (.ntdf or .tdf)
  - `output`: Output file path (optional, returns plaintext if not specified)
  - `clientId`: OAuth client ID (optional, falls back to env/defaults)
  - `clientSecret`: OAuth client secret (optional, falls back to env/defaults)

- `mcp__opentdf-mcp__list_attributes` - List available attributes
  - `namespace`: Filter by namespace (optional)
  - `verbose`: Show attribute values (boolean)
  - `clientId`: OAuth client ID (optional)
  - `clientSecret`: OAuth client secret (optional)

### memo-mcp

- `mcp__memo-mcp__render_memo_to_pdf` - Render markdown memo to PDF
- `mcp__memo-mcp__get_memo_schema` - Get memo frontmatter schema
- `mcp__memo-mcp__get_memo_example` - Get example memo with guidelines

---

## Testing Workflows

### 1. Verify Platform Attributes

Before testing, confirm the platform has the required attributes configured:

```
mcp__opentdf-mcp__list_attributes(verbose: true)
```

Expected namespaces/attributes:
- `https://demo.usaf.mil/attr/flight_rch2532101`
- `https://demo.usaf.mil/attr/flight_rch2532102`
- `https://demo.usaf.mil/attr/classification_topsecret`
- `https://demo.usaf.mil/attr/classification_secret`
- `https://demo.usaf.mil/attr/functional_maintenance`

### 2. Encrypt Scenario Documents

Encrypt each document with the appropriate attributes based on its type. **Prefer using `input` (file path) over `data` (literal text)** for efficiency and reliability:

**KC-46 Aircrew Summaries** (top-secret + flight RCH2532101):
```
mcp__opentdf-mcp__encrypt(
  input: "usaf-refueling-scenario/maj-evan-riley-kc-46-aircraft-commander.txt",
  attributes: [
    "https://demo.usaf.mil/attr/flight_rch2532101/value/true",
    "https://demo.usaf.mil/attr/classification_topsecret/value/true"
  ],
  output: "ENCRYPTED/maj-evan-riley-kc-46-aircraft-commander.ntdf"
)
```

**KC-46 Flight Logs** (secret + flight RCH2532101):
```
mcp__opentdf-mcp__encrypt(
  input: "usaf-refueling-scenario/kc-46-flight-log-data.csv",
  attributes: [
    "https://demo.usaf.mil/attr/flight_rch2532101/value/true",
    "https://demo.usaf.mil/attr/classification_secret/value/true"
  ],
  output: "ENCRYPTED/kc-46-flight-log-data.ntdf"
)
```

**Maintenance Documents** (secret + maintenance):
```
mcp__opentdf-mcp__encrypt(
  input: "usaf-refueling-scenario/maintenance-inspection-findings.csv",
  attributes: [
    "https://demo.usaf.mil/attr/flight_rch2532101/value/true",
    "https://demo.usaf.mil/attr/classification_secret/value/true",
    "https://demo.usaf.mil/attr/functional_maintenance/value/true"
  ],
  output: "ENCRYPTED/maintenance-inspection-findings.ntdf"
)
```

### 3. Test Access Control (Positive Tests)

The password is `mock.jwt.token` for all users in this scenario.

Test that authorized users CAN decrypt their documents:

```
# Test as Col Nies (full access)
mcp__opentdf-mcp__decrypt(
  input: "encrypted/maj-evan-riley-kc-46-aircraft-commander.ntdf",
  clientId: "ashley.nies",
  clientSecret: "mock.jwt.token"
)

# Test as Maj Riley (KC-46 flight only)
mcp__opentdf-mcp__decrypt(
  input: "encrypted/maj-evan-riley-kc-46-aircraft-commander.ntdf",
  clientId: "evan.riley",
  clientSecret: "mock.jwt.token"
)
```

### 4. Test Access Control (Negative Tests - The "Aha!" Moment)

Test that unauthorized users CANNOT decrypt documents:

```
# Maj Fernando should NOT be able to decrypt KC-46 documents
mcp__opentdf-mcp__decrypt(
  input: "encrypted/maj-evan-riley-kc-46-aircraft-commander.ntdf",
  clientId: "jonathan.fernando",
  clientSecret: "mock.jwt.token"
)
# Expected: Permission denied error

# Capt Chen should NOT be able to decrypt KC-46 documents
mcp__opentdf-mcp__decrypt(
  input: "encrypted/kc-46-flight-log-data.ntdf",
  clientId: "sarah.chen",
  clientSecret: "mock.jwt.token"
)
# Expected: Permission denied error
```

### 5. Cross-Flight Isolation Demo

Side-by-side comparison showing the strongest ABAC demonstration:

| User | Document | Expected Result |
|------|----------|-----------------|
| Capt Lee (KC-46 Co-Pilot) | kc-46-flight-log-data.ntdf | ✅ Decrypt succeeds |
| Capt Chen (C-17 Co-Pilot) | kc-46-flight-log-data.ntdf | ❌ Permission denied |
| Capt Chen (C-17 Co-Pilot) | c-17-flight-log-data.ntdf | ✅ Decrypt succeeds |
| Capt Lee (KC-46 Co-Pilot) | c-17-flight-log-data.ntdf | ❌ Permission denied |

**Same role, same rank, completely different access based on flight attribute. This is impossible with RBAC.**

---

## Debugging Checklist

### Pre-Demo Validation

1. ☐ Platform is running and accessible
2. ☐ `list_attributes` returns expected flight_rch..., classification_..., and functional_... attributes
3. ☐ All 10 scenario documents are encrypted with correct attributes
4. ☐ Keycloak users are configured with correct attributes (per SCENARIO_INTEGRATION.md)

### Common Issues

| Symptom | Possible Cause | Solution |
|---------|----------------|----------|
| "Unknown attribute FQN" | Attribute not configured in platform | Run `list_attributes` to find valid FQNs |
| "Permission denied" on decrypt | User lacks required attribute | Check user's Keycloak attributes |
| "KAS unreachable" | Platform endpoint misconfigured | Verify `OPENTDF_PLATFORM_ENDPOINT` |
| Decrypt succeeds when it shouldn't | Document encrypted with wrong attributes | Re-encrypt with correct attributes |

### Environment Variables

```bash
OPENTDF_PLATFORM_ENDPOINT=http://localhost:8080
OPENTDF_CLIENT_ID=ashley.nies  # or evan.riley, julie.lee, sarah.chen, etc.
OPENTDF_CLIENT_SECRET=mock.jwt.token
```

> **Note:** Client IDs use `firstname.lastname` format (e.g., `sarah.chen`, not `sarah-chen` or `SarahChen`).

---

## File Locations

- **Scenario source files:** `usaf-refueling-scenario/`
- **Encrypted files:** `ENCRYPTED/` (already exists with pre-encrypted files)
- **Decrypted output:** `decrypted/` (create if needed)
- **Memo drafts:** `drafts/`

---

## Best Practices

### Testing
- Always run positive tests BEFORE negative tests to confirm baseline functionality
- Document any mismatches between expected and actual results
- Use the access matrix as the authoritative source of truth

### Debugging
- Check terminal output for detailed error messages
- Verify attribute FQNs match exactly (case-sensitive)
- Confirm Keycloak user attributes match the matrix

### Encryption
- **Prefer `input` (file path) over `data`** for file-based encryption - this avoids token limits and encoding issues
- Only use `data` parameter for small inline text that doesn't exist as a file
- Always specify all required attributes for a document
- Use nanoTDF format (`.ntdf`) for better compatibility
- Store encrypted files in a dedicated `ENCRYPTED/` directory
