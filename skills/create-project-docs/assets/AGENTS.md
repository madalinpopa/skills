# Agent instructions

## Project

## Working rules

- Inspect relevant code, contracts, and diffs; preserve unrelated work and scope.
- Raise conflicting instructions or requirements before changing a contract.
- Search narrowly; load applicable sections and relevant dependencies.
- Commit or publish only on explicit request, except local prototype commits
  permitted by [PROTOTYPE.md](docs/conventions/PROTOTYPE.md).

## Workflow

- Work one small slice at a time across code, docs, configuration, and refactoring.
  Stay within approved scope; preserve requirements and shared decisions. Record
  established facts and evidence, without placeholders.
- Write necessary failing tests before implementation; follow
  [TDD.md](docs/conventions/TDD.md#test-workflow) for test rules and exceptions.
- Implement the smallest clear change and run applicable project checks.
- Before user review, delegate independent review following
  [REVIEW.md](docs/conventions/REVIEW.md) and wait for results. Resolve findings, rerun affected checks,
  and repeat review after fixes until no blockers remain.
- Stop for user review with changes, evidence, findings/dispositions, remaining
  concerns, and a proposed commit message. Report failed/skipped checks and causes.
- Plan approval skips no reviews and authorizes no commits. Reuse approval within
  its scope; passing checks or completing tasks does not imply user acceptance.

## Read for the current task

- For throwaway prototypes, follow [PROTOTYPE.md](docs/conventions/PROTOTYPE.md)
  instead of the production workflow above.
- For scope, read [MISSION.md](docs/specs/MISSION.md); for work selection or
  progress, read the relevant [ROADMAP.md](docs/specs/ROADMAP.md) group.
- Before production behavior changes, read affected `spec.md` criteria and shared contracts.
  Before test design, read [coverage guidance](docs/conventions/TDD.md#coverage-and-readability).
- Before code, dependency, or trust-boundary changes, read [CODE.md](docs/conventions/CODE.md);
  before Go edits, read [GOLANG.md](docs/conventions/GOLANG.md).
- Before checks, read commands/directories in [STACK.md](docs/specs/STACK.md).
- Before version-control or delivery work, read [COMMITS.md](docs/conventions/COMMITS.md).
  When `.jj` exists, invoke `use-jj`, even alongside `.git`; if unavailable,
  report it before mutations and do not substitute Git.
- For feature spec authoring, use `use-task-spec`; for plans, `use-task-plan`; for task
  checklists, execution, or resumption, `use-tasks-checklists`.
  Follow [record ownership](docs/features/README.md#ownership-and-routes).
- Before writing shared specs or architectural decisions, read the relevant
  [writing guide](docs/specs/README.md) or [decision guide](docs/decisions/README.md).
- When input is needed, follow [QUESTIONS.md](docs/conventions/QUESTIONS.md).

## Replies

- Keep replies brief.
- End with one concrete **Next Step**, linking the relevant task/question.
