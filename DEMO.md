# Demo Outline

## Flow

1. Maintainer, SrA PJ, writes memo using TtQ to his MX flight commander from his notes and maintenance findings
2. Wing Commander, Col Nies, writes a memo using TtQ to Air Mobility Command (AMC) Commander explaining what happened and why the tanker is going to be grounded for the forseable future.

## Attribute Based Access Control (ABAC)

- Flight related documents (individual summaries and logs) are encrypted with a unique flight identifier attribute
- KC-46 unique flight identifier: **RCH2532101** (Call Sign REACH, 2025, 321st day, 1st flight mission)
- C-17 unique flight identifier: **RCH2532102** (Call Sign REACH, 2025, 321st day, 2nd flight mission)

- Documents can be encypted with the **top-secret-fictional** attribute if they are top secret//fictional
- Documents can be encypted with the **secret-fictional** attribute if they are secret//fictional
- Documents can be encrypted with the **maintenance** attribute if they are relevant to maintainers

## Access

### Overview

- Wing commander can see everything
- Aircrew can see everyone on their flight's summary of the incident
- Aircrew can see their flight's logs
- Maintainers can see flight logs for planes they work on

### Col Ashley Nies: Wing Commander

Has attributes:

- RCH2532101
- RCH2532102
- maintenance
- top-secret-fictional
- secret-fictional

Therefore, can access:

- [capt-julie-lee-kc-46-co-pilot.txt](usaf-refueling-scenario/capt-julie-lee-kc-46-co-pilot.txt)
- [kc-46-flight-log-data.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv)
- [kc-46-refueling-log-data.csv](usaf-refueling-scenario/kc-46-refueling-log-data.csv)
- [maintenance-inspection-findings.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv)
- [maj-evan-riley-kc-46-aircraft-commander.txt](usaf-refueling-scenario/maj-evan-riley-kc-46-aircraft-commander.txt)
- [maj-jonathan-fernando-c-17-aircraft-commander](usaf-refueling-scenario/maj-jonathan-fernando-c-17-aircraft-commander.txt)
- [sra-pj-jones-kc-46-maintainer.txt](usaf-refueling-scenario/sra-pj-jones-kc-46-maintainer.txt)
- [tsgt-marcus-hayes-kc-46-boom-operator.txt](usaf-refueling-scenario/tsgt-marcus-hayes-kc-46-boom-operator.txt)

### Maj Jonathan Fernando: C-17 Commander

Has attributes:

- RCH2532102
- top-secret-fictional
- secret-fictional

Therefore, can access:

- [maj-jonathan-fernando-c-17-aircraft-commander](usaf-refueling-scenario/maj-jonathan-fernando-c-17-aircraft-commander.txt)

### Maj Evan Riley: KC-46 Commander

Has attributes:

- RCH2532101
- top-secret-fictional
- secret-fictional

Therefore, can access:

- [maj-evan-riley-kc-46-aircraft-commander.txt](usaf-refueling-scenario/maj-evan-riley-kc-46-aircraft-commander.txt)
- [capt-julie-lee-kc-46-co-pilot.txt](usaf-refueling-scenario/capt-julie-lee-kc-46-co-pilot.txt)
- [tsgt-marcus-hayes-kc-46-boom-operator.txt](usaf-refueling-scenario/tsgt-marcus-hayes-kc-46-boom-operator.txt)
- [kc-46-flight-log-data.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv)
- [kc-46-refueling-log-data.csv](usaf-refueling-scenario/kc-46-refueling-log-data.csv)

### Capt Julie Lee: KC-46 Co-Pilot

Has attributes:

- RCH2532101
- top-secret-fictional
- secret-fictional

Therefore, can access:

- [maj-evan-riley-kc-46-aircraft-commander.txt](usaf-refueling-scenario/maj-evan-riley-kc-46-aircraft-commander.txt)
- [capt-julie-lee-kc-46-co-pilot.txt](usaf-refueling-scenario/capt-julie-lee-kc-46-co-pilot.txt)
- [tsgt-marcus-hayes-kc-46-boom-operator.txt](usaf-refueling-scenario/tsgt-marcus-hayes-kc-46-boom-operator.txt)
- [kc-46-flight-log-data.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv)
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
- [kc-46-flight-log-data.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv)
- [kc-46-refueling-log-data.csv](usaf-refueling-scenario/kc-46-refueling-log-data.csv)

### SrA PJ Jones: KC-46 Maintainer

Has attributes:

- RCH2532101
- maintenance
- secret-fictional

Therefore, can access:

- [sra-pj-jones-kc-46-maintainer.txt](usaf-refueling-scenario/sra-pj-jones-kc-46-maintainer.txt)
- [kc-46-flight-log-data.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv)
- [kc-46-refueling-log-data.csv](usaf-refueling-scenario/kc-46-refueling-log-data.csv)
- [maintenance-inspection-findings.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv)

## Documents and their Attributes

- [capt-julie-lee-kc-46-co-pilot.txt](usaf-refueling-scenario/capt-julie-lee-kc-46-co-pilot.txt) has attributes RCH2532101 and top-secret-fictional
- [kc-46-flight-log-data.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv) has attributes RCH2532101 and secret-fictional
- [kc-46-refueling-log-data.csv](usaf-refueling-scenario/kc-46-refueling-log-data.csv) has attributes RCH2532101 and secret-fictional
- [maintenance-inspection-findings.csv](usaf-refueling-scenario/maintenance-inspection-findings.csv) has attributes maintenance and secret-fictional
- [maj-evan-riley-kc-46-aircraft-commander.txt](usaf-refueling-scenario/maj-evan-riley-kc-46-aircraft-commander.txt) has attributes RCH2532101 and top-secret-fictional
- [maj-jonathan-fernando-c-17-aircraft-commander](usaf-refueling-scenario/maj-jonathan-fernando-c-17-aircraft-commander.txt) has attributes RCH2532102 and top-secret-fictional
- [sra-pj-jones-kc-46-maintainer.txt](usaf-refueling-scenario/sra-pj-jones-kc-46-maintainer.txt) has attributes maintenance and secret-fictional
- [tsgt-marcus-hayes-kc-46-boom-operator.txt](usaf-refueling-scenario/tsgt-marcus-hayes-kc-46-boom-operator.txt) has attributes RCH2532101 and top-secret-fictional
