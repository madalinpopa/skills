# Agent instructions

## Project

## Working rules

- Inspect relevant code, contracts, and diffs; preserve unrelated work and scope.
- Raise conflicting instructions or requirements before changing a contract.
- For production work, use small slices; verify and present each for review.
- Search narrowly; load applicable sections and expand through relevant dependencies.
- Preserve agreed requirements; record established facts and evidence, not placeholders.

## Permissions

- When `.jj` exists, invoke `use-jj`, even alongside `.git`. If unavailable,
  report it before mutations; do not substitute Git.
- Commit only on explicit request; plan/work approval is insufficient. Never add
  AI attribution. The [prototype workflow](docs/conventions/TDD.md#throwaway-prototypes)
  permits local prototype commits. Publish only on explicit request.
- Dependency additions, installations, replacements, and upgrades need approval.

## Read for the current task

- For scope decisions, read [MISSION.md](docs/specs/MISSION.md); for work selection
  or progress updates, read the relevant [ROADMAP.md](docs/specs/ROADMAP.md) group.
- For new or existing throwaway prototype work, read and use the
  [prototype workflow](docs/conventions/TDD.md#throwaway-prototypes) instead of
  production planning, TDD, and review gates.
- Before changing production behavior, find the affected feature contract and acceptance
  criteria in `requirements.md`; follow relevant shared specs and decisions.
- Before code, dependency, or trust-boundary changes, read [CODE.md](docs/conventions/CODE.md); before Go edits,
  read [GOLANG.md](docs/conventions/GOLANG.md).
- Before other behavior changes, test work, or verification, read [TDD workflow](docs/conventions/TDD.md#workflow),
  including its required review procedure. Before choosing/writing tests, also read
  [coverage guidance](docs/conventions/TDD.md#coverage-and-readability).
- Before checks, read commands and directories in [STACK.md](docs/specs/STACK.md).
  Before proposing commits or delivery, read [COMMITS.md](docs/conventions/COMMITS.md).
- Before writing specs, read [writing rules](docs/specs/README.md#working-rules)
  and the target document section; reading requirements alone does not trigger authoring.
- Before creating/planning records, read [ownership](docs/features/README.md#naming-and-ownership)
  and applicable writing sections; for resumption, read [Resume](docs/features/README.md#resume).
- Before recording architectural choices, read the [decision guide](docs/decisions/README.md).
- When input is needed, follow the [question guide](docs/conventions/QUESTIONS.md)
  for asking, recording blockers, and applying answers.

## Replies

- Report changes, evidence, failures, and limits in short bullets; completion is not acceptance.
- End with one concrete **Next Step**, linking the relevant task/question;
  proposals authorize no additional work.
