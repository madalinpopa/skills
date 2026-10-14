# Agent instructions

Project: Pending: purpose and users.

## Working rules

Work in small slices for every task, including docs, configuration, and
refactoring. Give each slice one outcome; verify and present it for user review
before starting the next. Keep changes within scope.

Inspect relevant code and diffs; preserve unrelated work. Surface conflicting
instructions, requirements, or code before changing the contract. Record unknowns
and unapproved choices as pending.

Use commands and working directories in STACK.md; report missing commands.
Verify current official documentation and latest released library versions.
External dependency adoption, installation, replacement, or upgrades need user
approval; checking the latest version does not authorize an upgrade.

Follow TDD's test and implementation review gates for behavior changes. Propose
a commit message after verified changes. Commit or publish only on explicit
request; never add AI attribution.

## Read only for the current task

| Task | Guidance |
| --- | --- |
| Goals, users, scope | [MISSION.md](docs/specs/MISSION.md) |
| Capability index | [FEATURES.md](docs/specs/FEATURES.md) |
| Delivery order | [ROADMAP.md](docs/specs/ROADMAP.md) |
| Technologies, versions, commands | [STACK.md](docs/specs/STACK.md) |
| Code, TypeScript/frontend, dependencies | [CODE.md](docs/conventions/CODE.md) |
| Go work | [GOLANG.md](docs/conventions/GOLANG.md) |
| Behavior changes, tests, verification | [TDD.md](docs/conventions/TDD.md) |
| Commit proposals, delivery | [COMMITS.md](docs/conventions/COMMITS.md) |
| Architecture choices | [Decision guide](docs/decisions/README.md) and relevant records |
| Feature planning, splitting, execution, resume | [Feature guide](docs/features/README.md), then active requirements.md, plan.md, and relevant validation.md |

Feature requirements own behavior, plans own steps/checkpoints, and validation
records own evidence. Follow the feature guide's checkpoint updates and resume
checks; load historical records only for a relevant dependency.
