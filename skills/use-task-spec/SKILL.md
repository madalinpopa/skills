---
name: use-task-spec
description: Clarify and write a feature or change's spec.md with goals, constraints, and observable acceptance criteria. Use for task requirements, not technical plans or implementation.
status: published
tags: [docs, planning, specs]
---

# Write a task spec

Read the requested scope, relevant project instructions, existing requirements,
and repository evidence. Reuse settled answers. Record established facts;
omit empty headings, placeholders, and invented scope.

Use the project's existing record location. For a new record without a local
convention, use `docs/features/NNNN-<slug>/spec.md`; choose the next unused
four-digit ID across feature folders and `docs/specs/ROADMAP.md`. Never renumber
existing work. Propose independently verifiable sub-features when scope cannot
be reviewed in small slices; preserve parent/child links and shared constraints.
Give each sub-feature its own ID and records; parent acceptance requires child
evidence, shared criteria, and explicit user acceptance.

## Clarify, then draft

Use the available question tool for focused rounds with concrete options and
meaningful trade-offs; otherwise ask in chat. Ask only about unsettled decisions
affecting the requested spec. Cover edge cases, failures, recovery, repeated or
conflicting actions, partial results, and concurrency when relevant.
Wait for answers before writing dependent content. If the user requests an
interview before drafting, wait for their explicit instruction to write.

Follow the project's question-log convention. When using `docs/QUESTIONS.md`,
group simple unchecked questions by purpose; reuse existing entries. Apply
answers to their owning sections before checking questions and linking the answer.
Leave partial or conflicting answers open. Do not create generic questions.

Use `# SPEC: <name>` and include sections only when they have content:

- **Goal:** problem, affected users or maintainers, intended outcome.
- **Context:** existing behavior, boundaries, and shared decisions.
- **Constraints:** compatibility, invariants, permissions, agreed limits.
- **Acceptance Criteria:** checklists of observable outcomes, including relevant
  failure and recovery behavior. Add verification methods where needed for
  precision; do not assign a test layer to every criterion.
- **Non-Goals:** deliberate exclusions, never undecided scope.

Keep implementation choices in `plan.md`, progress and evidence in `tasks.md`.
Use concise sections and lists; link shared facts instead of copying them.
Check criteria only when task evidence proves them; completion is not acceptance.
Present the spec for review without implementing it. Report unresolved decisions
and one concrete next action; follow project commit and approval rules.
