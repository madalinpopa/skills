---
name: use-tasks-checklists
description: Create phased tasks.md checklists from a plan, or execute and resume requested tasks with verifiable done checks, evidence, and a clear next action.
status: published
tags: [planning, tasks, workflow]
---

# Work with task checklists

Follow project instructions and approval gates. `spec.md` owns requirements,
`plan.md` owns approach and sequencing, and `tasks.md` owns execution, evidence,
and the resume point. Preserve existing record locations and agreed scope.

## Create a checklist

Read the relevant spec and plan. Write `tasks.md` beside them using
`# TASKS: <name>`, phase headings named for outcomes, and task checklists.

- Give each task a specific named change and **Done when:** an observable check
  understandable without reading the full plan. Link relevant requirements
  and approach rather than copying them.
- Order by dependencies; no task may depend on one listed later. Keep phases
  small and reviewable. Detail the active phase; keep identified later phases brief.
- Include necessary tests, implementation, verification, and required reviews.
  Follow project TDD policy; do not manufacture failing tests for documentation.
- Keep partial work unchecked. Use `in progress`, `blocked: <reason>`,
  `awaiting review`, or `deferred: <reason>` only when useful.
- While work remains, mark one actionable unchecked task or subtask `next`.
  A pending decision or required review can be the next action.

Creating a checklist does not authorize executing it. Present the result for
review with one concrete next action.

## Execute or resume

For requested execution, read [execution](references/execution.md). Start from
the active phase, next task, and nearby evidence; load only relevant spec and
plan sections and dependencies. Reading requirements does not require loading
their authoring skills.
