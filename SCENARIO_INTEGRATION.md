# Scenario Integration

Technical plan for ABAC [scenario](DEMO.md) integration with OpenTDF platform and Keycloak.

## Key Demo Objective

**Demonstrate why ABAC is superior to RBAC for flight-scoped access control.**

The "aha!" moment: Two pilots (Maj Riley and Maj Fernando) have the same role, but can only decrypt records from *their specific flight*. This per-flight access is infeasible with traditional Role-Based Access Control, which would require creating a new role for every flight mission.

## Chatmode Prompts

Create base prompt template with placeholders for user data:
- Name
- JWT (hardcoded)

## Keycloak Config

Configure users from [scenario](DEMO.md). Apply flag-based attributes to users:

| User | flight_rch2532101 | flight_rch2532102 | classification_topsecret | classification_secret | functional_maintenance |
|------|:-----------------:|:-----------------:|:------------------------:|:---------------------:|:----------------------:|
| Col Ashley Nies | ✓ | ✓ | ✓ | | ✓ |
| Maj Evan Riley | ✓ | | ✓ | | |
| Capt Julie Lee | ✓ | | ✓ | | |
| TSgt Marcus Hayes | ✓ | | ✓ | | |
| Maj Jonathan Fernando | | ✓ | ✓ | | |
| Capt Sarah Chen | | ✓ | ✓ | | |
| SrA PJ Jones | ✓ | | | ✓ | ✓ |

> **Note:** Using flag-based attributes (value="true") instead of value-based attributes. This simplifies the subject mapping conditions and allows Col Nies to have access to BOTH flights.

### Keycloak Client Credentials

Each user is configured as an OAuth client in Keycloak. Use these credentials for MCP tool authentication:

| User | Client ID | Client Secret |
|------|-----------|---------------|
| Col Ashley Nies | `ashley.nies` | `password123` |
| Maj Evan Riley | `evan.riley` | `password123` |
| Capt Julie Lee | `julie.lee` | `password123` |
| TSgt Marcus Hayes | `marcus.hayes` | `password123` |
| Maj Jonathan Fernando | `jonathan.fernando` | `password123` |
| Capt Sarah Chen | `sarah.chen` | `password123` |
| SrA PJ Jones | `pj.jones` | `password123` |

> **Note:** Do NOT use `opentdf-sdk` as a client ID - that is only for platform administration, not user authentication.

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

## Demo Validation Checklist

Before the demo, verify:

1. **Positive access tests:**
   - Col Nies can decrypt ALL documents (10 total)
   - Maj Riley can decrypt KC-46 crew summaries and logs (RCH2532101)
   - Maj Fernando can decrypt C-17 crew summaries and logs (RCH2532102)
   - Capt Chen can decrypt C-17 documents (same flight as Maj Fernando)

2. **Negative access tests (the "aha!" moment):**
   - Maj Fernando CANNOT decrypt KC-46 documents (lacks RCH2532101)
   - Maj Riley CANNOT decrypt C-17 documents (lacks RCH2532102)
   - Capt Chen CANNOT decrypt KC-46 documents (lacks RCH2532101)
   - Capt Lee CANNOT decrypt C-17 documents (lacks RCH2532102)
   - SrA PJ Jones CANNOT decrypt aircrew summaries (lacks top-secret-fictional)

3. **Cross-flight isolation (strongest demo):**
   - Compare Capt Lee (KC-46 Co-Pilot) vs Capt Chen (C-17 Co-Pilot)
   - Same role, same rank, completely different access based on flight attribute
   - This is **impossible with RBAC** without creating per-flight roles

## Document/User Matrix

Below is a compact matrix showing which users can decrypt which documents in the scenario. ✅ = decrypt allowed, ❌ = decrypt denied.

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

> Note: Wing Commander (Col Nies) has attributes: RCH2532101, top-secret-fictional, and maintenance. Maintainers can view maintenance and engineering logs for flights they service. Aircrew see flight-scoped documents via RCH2532101 or RCH2532102 attributes. Each user has a single value per attribute (flattened model).

### How to Use the Matrix During the Demo

- Use `opentdf-mcp list_attributes` to verify configured attributes for each user.
- Use `opentdf-mcp GetDecisionBulk` (or the equivalent config in the MCP) to run the matrix table as a test across all documents/users and record the decision outputs (allow/deny).
- For a visual demo, show the result of `GetDecisionBulk` for two user sessions side-by-side: Capt Lee (KC-46) vs Capt Chen (C-17) — both are Co-Pilots but have different access.

> Tip: When running `GetDecisionBulk`, assert the result equals the matrix values. Any mismatch indicates attribute misconfiguration in Keycloak or an incorrect attribute on the encrypted TDF file.