---
name: use-task-plan
description: Turn an existing spec.md into a repository-grounded technical plan.md. Use for approach, trade-offs, affected files, and sequencing, without implementation.
status: published
tags: [docs, planning]
---

# Plan a task

Read the relevant `spec.md`, project instructions, and existing decisions.
Explore affected code, tests, and dependency boundaries before proposing changes.
Preserve the project's record location and agreed requirements. If the spec is
missing or an unresolved decision changes the approach, surface that dependency
before writing the affected plan; use `use-task-spec` when spec authoring is needed.

Write `plan.md` beside the spec with `# PLAN: <name>`. Include only sections
with established content:

- **Approach:** how the change satisfies the spec, affected boundaries, dependencies.
- **Trade-Offs:** meaningful alternatives, chosen option, and reason.
- **Files Affected:** known paths or boundaries and intended changes.
- **Sequencing:** ordered implementation dependencies and delivery steps.

Link shared requirements rather than duplicating them. Keep progress and
executable checklists in `tasks.md`; use `use-tasks-checklists` to create them.
Detail the active scope and leave identified later work brief. Propose separate,
independently verifiable sub-features if work cannot be reviewed in small slices;
retain parent/child links and dependency ownership.

Sequencing must preserve the project's test-first workflow and review gates;
later integration verification does not defer necessary failing tests until
after implementation. Record actual unresolved questions using the project's
question convention; wait for decisions before dependent work.

Return the plan for review, with no implementation code or unrelated document
edits. State material uncertainties and one concrete next action. Plan approval
does not authorize commits or bypass required reviews.
