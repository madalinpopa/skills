# Agent instructions

Project purpose: Pending: describe the project and its users.

## Read for the current task

Load only the documents relevant to the requested work:

| Task | Guidance |
| --- | --- |
| Product goals, users, scope | [Mission](docs/specs/MISSION.md) |
| Behavior, requirements, acceptance | [Features](docs/specs/FEATURES.md) |
| Delivery order and dependencies | [Roadmap](docs/specs/ROADMAP.md) |
| Technologies, versions, commands | [Stack](docs/specs/STACK.md) |
| Code changes or review | [Code conventions](docs/conventions/CODE.md) |
| Go work | [Go conventions](docs/conventions/GOLANG.md) |
| Behavior changes, tests, verification | [TDD](docs/conventions/TDD.md) |
| Commit proposals or delivery | [Commit conventions](docs/conventions/COMMITS.md) |
| Architecture decisions and tradeoffs | [Decision records](docs/decisions/README.md) and relevant records |
| Feature work | [Feature records](docs/features/README.md): active requirements.md, plan.md, and validation.md |

Inspect relevant code and the current diff before editing; preserve unrelated
work. FEATURES.md indexes capabilities; each feature's requirements.md owns
its behavior, plan.md tracks work, and validation.md records evidence.
If existing code, docs, or instructions conflict, surface the conflict before
changing the contract.

## Working rules

Follow the TDD review gates for behavioral changes. Use commands and working
directories recorded in STACK.md; report missing commands rather than inventing
them. Verify current official documentation for external APIs being changed.
External dependency additions or changes require user approval and latest-version
verification under [code conventions](docs/conventions/CODE.md).
Read Go guidance only for Go work; apply TypeScript/frontend rules only when
the project uses them.

Record unknown requirements and unapproved choices as pending. Update the active
feature's checkpoint in plan.md and evidence in validation.md at review boundaries
and session end. Verify the checkpoint against the working tree on resume.
Load historical records only for a relevant dependency.

Propose a commit message after verified changes. Commit or publish only when
explicitly requested; never add AI attribution.
