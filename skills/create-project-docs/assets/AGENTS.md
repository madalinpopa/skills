# Agent instructions

Project: Pending: purpose and users.

## Working rules

Work in small slices with one outcome each, including docs, configuration, and
refactoring. Verify each slice and present it for user review before the next.
Stay in scope. Inspect relevant code and diffs; preserve unrelated work. Raise
conflicts between instructions, requirements, and code before changing a contract.

Never invent requirements. Reuse known facts and ask only the minimum questions
the next slice needs, with the ask-user-question tool or in chat. Record other
unknowns and unapproved choices as pending.

Use the commands and working directories in STACK.md; report missing ones.
Check current official docs and the latest released library versions. Adding,
installing, replacing, or upgrading a dependency needs user approval; a version
check is not approval.

For each change, follow [TDD's workflow](docs/conventions/TDD.md#workflow):
write a failing test and see it fail, implement, then have a different agent do
an adversarial review; fix and repeat until clean. Then the user reviews and
approves the commit. Commit only on that approval; publish only on explicit
request; never add AI attribution.

Reply in short `-` bullets for changes and verification, without repeats. End
with **User must be aware**: material risks, limits, and pending decisions, or
"None". Record pending items in the relevant docs and ask the user to review them.

## Read only for the current task

| Task | Guidance |
| --- | --- |
| Goals, users, scope | [MISSION.md](docs/specs/MISSION.md) |
| Feature index | [FEATURES.md](docs/specs/FEATURES.md) |
| Delivery order, feature status | [ROADMAP.md](docs/specs/ROADMAP.md) |
| Technologies, versions, commands | [STACK.md](docs/specs/STACK.md) |
| Code, TypeScript/frontend, dependencies | [CODE.md](docs/conventions/CODE.md) |
| Go work | [GOLANG.md](docs/conventions/GOLANG.md) |
| Behavior changes, tests, verification | [TDD.md](docs/conventions/TDD.md) |
| Commit proposals, delivery | [COMMITS.md](docs/conventions/COMMITS.md) |
| Architecture choices | [Decision guide](docs/decisions/README.md) |
| Feature planning, work, or resume | [Feature guide](docs/features/README.md) |

Use task/subtask checklists for roadmaps and plans; mark completed work with
evidence. Write specs as short bullets under clear sections; reserve tables for
compact mappings. ROADMAP.md owns feature status, requirements own behavior,
plans own task completion and the checkpoint, and validation owns evidence.
Load old feature records only for a relevant dependency.
