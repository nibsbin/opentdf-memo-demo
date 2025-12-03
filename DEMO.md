# Demo Outline

## Flow

1. Maintainer, SrA PJ, writes memo using TtQ to his MX flight commander from his notes and maintenance findings
2. Wing Commander, Col Nies, writes a memo using TtQ to Air Mobility Command (AMC) Commander explaining what happened and why the tanker is going to be grounded for the forseable future.

## Attribute Based Access Control (ABAC)

### The "Aha!" Moment: Why ABAC Beats RBAC

**Traditional Role-Based Access Control (RBAC) Problem:**
In RBAC, you would create roles like "Pilot", "Boom Operator", "Maintainer". But this fails because:
- All pilots would see all flight records
- All boom operators would see all refueling logs
- There's no way to scope access to *specific flights* without creating thousands of roles

**ABAC Solution:**
With Attribute-Based Access Control, each user is tagged with the specific flight identifiers they participated in. Access is granted based on matching attributes, not roles. This means:
- Maj Fernando (C-17 pilot) can only see C-17 flight RCH2532102 records
- Maj Riley (KC-46 pilot) can only see KC-46 flight RCH2532101 records
- Both are "pilots" but see completely different data based on *which flight they were on*

This is **infeasible with RBAC** because you'd need to create a new role for every single flight mission.

### Flight Identifier Attributes (flight_id)

- Flight-related documents (individual summaries and logs) are encrypted with a unique `flight_id` attribute
- KC-46 unique flight identifier: **RCH2532101** (Call Sign REACH, 2025, 321st day, 1st flight mission)
- C-17 unique flight identifier: **RCH2532102** (Call Sign REACH, 2025, 321st day, 2nd flight mission)
- Attribute namespace: `https://demo.usaf.mil`
- Full FQN example: `https://demo.usaf.mil/attr/flight_id/value/RCH2532101`

### Classification Attributes (classification)

- Documents can be encrypted with the **top-secret-fictional** attribute if they are top secret//fictional
- Documents can be encrypted with the **secret-fictional** attribute if they are secret//fictional
- Full FQN examples:
  - `https://demo.usaf.mil/attr/classification/value/top-secret-fictional`
  - `https://demo.usaf.mil/attr/classification/value/secret-fictional`

### Functional Attributes (functional)

- Documents can be encrypted with the **maintenance** attribute if they are relevant to maintainers
- Full FQN: `https://demo.usaf.mil/attr/functional/value/maintenance`

## Access

### Overview

- Wing commander can see everything (has all flight identifiers + all classifications)
- Aircrew can see **only their own flight's** personnel summaries and logs (flight-scoped access)
- Maintainers can see maintenance-specific documents and flight logs for planes they service

### Key ABAC Demonstration Points

| User | Role | flight_id | Can Decrypt KC-46 (RCH2532101) Files? | Can Decrypt C-17 (RCH2532102) Files? |
|------|------|-----------|---------------------------------------|--------------------------------------|
| Col Nies | Wing Commander | RCH2532101, RCH2532102 | ✅ Yes | ✅ Yes |
| Maj Riley | KC-46 Pilot | RCH2532101 | ✅ Yes | ❌ No |
| Capt Lee | KC-46 Co-Pilot | RCH2532101 | ✅ Yes | ❌ No |
| Maj Fernando | C-17 Pilot | RCH2532102 | ❌ No | ✅ Yes |
| Capt Chen | C-17 Co-Pilot | RCH2532102 | ❌ No | ✅ Yes |

**This is the "aha!" moment**: Maj Riley and Maj Fernando are both Aircraft Commanders. Capt Lee and Capt Chen are both Co-Pilots. But each crew can only access records from *their specific flight*—not the other aircraft's data.

### Col Ashley Nies: Wing Commander

Has attributes:

- RCH2532101
- RCH2532102
- maintenance
- top-secret-fictional
- secret-fictional

Therefore, can access **all documents** (demonstrates full oversight capability):

- All KC-46 crew summaries and logs (RCH2532101)
- All C-17 crew summaries and logs (RCH2532102)
- All maintenance documents

### Maj Jonathan Fernando: C-17 Commander

Has attributes:

- RCH2532102
- top-secret-fictional
- secret-fictional

Therefore, can access:

- [maj-jonathan-fernando-c-17-aircraft-commander.txt](usaf-refueling-scenario/maj-jonathan-fernando-c-17-aircraft-commander.txt)
- [capt-sarah-chen-c-17-co-pilot.txt](usaf-refueling-scenario/capt-sarah-chen-c-17-co-pilot.txt)
- [c-17-flight-log-data.csv](usaf-refueling-scenario/c-17-flight-log-data.csv)

**Cannot access** KC-46 flight records (RCH2532101) even though he is a pilot—demonstrating flight-scoped ABAC.

### Capt Sarah Chen: C-17 Co-Pilot

Has attributes:

- RCH2532102
- top-secret-fictional
- secret-fictional

Therefore, can access:

- [maj-jonathan-fernando-c-17-aircraft-commander.txt](usaf-refueling-scenario/maj-jonathan-fernando-c-17-aircraft-commander.txt)
- [capt-sarah-chen-c-17-co-pilot.txt](usaf-refueling-scenario/capt-sarah-chen-c-17-co-pilot.txt)
- [c-17-flight-log-data.csv](usaf-refueling-scenario/c-17-flight-log-data.csv)

**Cannot access** KC-46 flight records (RCH2532101)—demonstrating flight-scoped ABAC.

### Maj Evan Riley: KC-46 Commander

Has attributes:

- RCH2532101
- top-secret-fictional
- secret-fictional

Therefore, can access:

- [maj-evan-riley-kc-46-aircraft-commander.txt](usaf-refueling-scenario/maj-evan-riley-kc-46-aircraft-commander.txt)
- [capt-julie-lee-kc-46-co-pilot.txt](usaf-refueling-scenario/capt-julie-lee-kc-46-co-pilot.txt)
- [tsgt-marcus-hayes-kc-46-boom-operator.txt](usaf-refueling-scenario/tsgt-marcus-hayes-kc-46-boom-operator.txt)
- [kc-46-flight-log-data.csv](usaf-refueling-scenario/kc-46-flight-log-data.csv)
- [kc-46-refueling-log-data.csv](usaf-refueling-scenario/kc-46-refueling-log-data.csv)

**Cannot access** C-17 flight records (RCH2532102) or maintenance findings—demonstrating flight-scoped ABAC.

### Capt Julie Lee: KC-46 Co-Pilot

Has attributes:

- RCH2532101
- top-secret-fictional
- secret-fictional

Therefore, can access:

- [maj-evan-riley-kc-46-aircraft-commander.txt](usaf-refueling-scenario/maj-evan-riley-kc-46-aircraft-commander.txt)
- [capt-julie-lee-kc-46-co-pilot.txt](usaf-refueling-scenario/capt-julie-lee-kc-46-co-pilot.txt)
- [tsgt-marcus-hayes-kc-46-boom-operator.txt](usaf-refueling-scenario/tsgt-marcus-hayes-kc-46-boom-operator.txt)
- [kc-46-flight-log-data.csv](usaf-refueling-scenario/kc-46-flight-log-data.csv)
- [kc-46-refueling-log-data.csv](usaf-refueling-scenario/kc-46-refueling-log-data.csv)

### TSgt Marcus Hayes: KC-46 Boom Operator

Has attributes:

- RCH2532101
- top-secret-fictional
- secret-fictional

Therefore, can access:

- [maj-evan-riley-kc-46-aircraft-commander.txt](usaf-refueling-scenario/maj-evan-riley-kc-46-aircraft-commander.txt)
- [capt-julie-lee-kc-46-co-pilot.txt](usaf-refueling-scenario/capt-julie-lee-kc-46-co-pilot.txt)
- [tsgt-marcus-hayes-kc-46-boom-operator.txt](usaf-refueling-scenario/tsgt-marcus-hayes-kc-46-boom-operator.txt)
- [kc-46-flight-log-data.csv](usaf-refueling-scenario/kc-46-flight-log-data.csv)
- [kc-46-refueling-log-data.csv](usaf-refueling-scenario/kc-46-refueling-log-data.csv)

### SrA PJ Jones: KC-46 Maintainer

Has attributes:

- RCH2532101
- maintenance
- secret-fictional

Therefore, can access:

- [sra-pj-jones-kc-46-maintainer.txt](usaf-refueling-scenario/sra-pj-jones-kc-46-maintainer.txt)
- [kc-46-flight-log-data.csv](usaf-refueling-scenario/kc-46-flight-log-data.csv)
- [kc-46-refueling-log-data.csv](usaf-refueling-scenario/kc-46-refueling-log-data.csv)
- [maintenance-inspection-findings.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv)

**Cannot access** aircrew summaries (requires top-secret-fictional) or C-17 records (requires RCH2532102).

## Documents and their Attributes

### KC-46 Flight Documents (RCH2532101)

| Document | flight_id | classification | functional |
|----------|-----------|----------------|------------|
| [maj-evan-riley-kc-46-aircraft-commander.txt](usaf-refueling-scenario/maj-evan-riley-kc-46-aircraft-commander.txt) | RCH2532101 | top-secret-fictional | — |
| [capt-julie-lee-kc-46-co-pilot.txt](usaf-refueling-scenario/capt-julie-lee-kc-46-co-pilot.txt) | RCH2532101 | top-secret-fictional | — |
| [tsgt-marcus-hayes-kc-46-boom-operator.txt](usaf-refueling-scenario/tsgt-marcus-hayes-kc-46-boom-operator.txt) | RCH2532101 | top-secret-fictional | — |
| [kc-46-flight-log-data.csv](usaf-refueling-scenario/kc-46-flight-log-data.csv) | RCH2532101 | secret-fictional | — |
| [kc-46-refueling-log-data.csv](usaf-refueling-scenario/kc-46-refueling-log-data.csv) | RCH2532101 | secret-fictional | — |

### C-17 Flight Documents (RCH2532102)

| Document | flight_id | classification | functional |
|----------|-----------|----------------|------------|
| [maj-jonathan-fernando-c-17-aircraft-commander.txt](usaf-refueling-scenario/maj-jonathan-fernando-c-17-aircraft-commander.txt) | RCH2532102 | top-secret-fictional | — |
| [capt-sarah-chen-c-17-co-pilot.txt](usaf-refueling-scenario/capt-sarah-chen-c-17-co-pilot.txt) | RCH2532102 | top-secret-fictional | — |
| [c-17-flight-log-data.csv](usaf-refueling-scenario/c-17-flight-log-data.csv) | RCH2532102 | secret-fictional | — |

### Maintenance Documents

| Document | flight_id | classification | functional |
|----------|-----------|----------------|------------|
| [sra-pj-jones-kc-46-maintainer.txt](usaf-refueling-scenario/sra-pj-jones-kc-46-maintainer.txt) | RCH2532101 | secret-fictional | maintenance |
| [maintenance-inspection-findings.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv) | RCH2532101 | secret-fictional | maintenance |

### Attribute Combination Rules

To decrypt a document, a user must have **ALL** of the document's attributes:
- Flight ID attribute (RCH2532101 or RCH2532102)
- Classification attribute (top-secret-fictional or secret-fictional)
- Functional attribute (maintenance, if applicable)
