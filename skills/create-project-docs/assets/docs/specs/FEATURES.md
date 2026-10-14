# Product features

Index the project's actual capabilities and their agreed status. Requirements
must fit [MISSION.md](MISSION.md); delivery order belongs in
[ROADMAP.md](ROADMAP.md). Do not invent capabilities to populate the index.

## Capability index

Pending: list capabilities agreed with the user. Give each a stable ID such as
`F-001`, a name, a short intended outcome, and a status: proposed, agreed, or
implemented. Add links to its numbered feature folder when planning begins:

- `requirements.md`: detailed behavior, scope, and acceptance criteria.
- `plan.md`: delivery steps and current implementation checkpoint.
- `validation.md`: verification evidence and user acceptance.

Use the naming rules and three templates in the
[feature guide](../features/README.md). A capability can remain proposed in this
index without creating an empty feature folder.

## Contract ownership

Detailed contracts live in each feature's requirements.md. Include input rules,
observable results, state and invariants, meaningful failures, and relevant
repeated/conflicting/partial cases there. Acceptance criteria need stable IDs
and concrete scenarios; validation.md links those criteria to evidence.

Keep shared goals and constraints in MISSION.md and cross-feature choices in
[decision records](../decisions/README.md). Link these owners rather than
maintaining duplicate requirements here. Distinguish agreed behavior from
proposals and implementation progress. Changing an agreed contract requires
an explicit user decision; update related plans and validation expectations.
