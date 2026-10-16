# Feature and codebase work records

Use these records for features, fixes, refactoring, migrations, tooling, and
documentation work that needs a plan. Create only real, agreed work; never seed
example features. For throwaway prototypes, use
[PROTOTYPE.md](../conventions/PROTOTYPE.md) instead.

## Naming and ownership

- Create `NNNN-<slug>/` with `spec.md`, `plan.md`, and `tasks.md`.
- Choose the next unused four-digit ID across [ROADMAP.md](../specs/ROADMAP.md)
  and feature folders. Preserve existing IDs and names; roadmap entries may
  precede their folders.
- `spec.md` owns the outcome, context, constraints, acceptance criteria, and
  non-goals. `plan.md` owns the technical approach and delivery order.
  `tasks.md` owns execution, progress, evidence, and the resume point.
- ROADMAP.md owns delivery order across features through headings and task
  checklists. Link to these records rather than duplicating their details.
- Keep these three records; do not add competing feature, phase, or validation files.

## Writing rules

- Record established facts and agreed scope. Omit empty sections, speculative
  phases, and filler. Do not require criterion/task IDs.
- Use descriptive headings and concise lists; short explanatory paragraphs are
  useful for goals and approaches. Link shared facts to their owner.
- For requested spec work, follow the
  [spec writing rules](../specs/README.md#working-rules), including guided
  questions and waiting for answers that affect the content.
- Keep unresolved questions in [the question log](../QUESTIONS.md), following
  the [question guide](../conventions/QUESTIONS.md); link the affected section.

## spec.md

Use `# SPEC: <name>` and these sections when they have content:

- **Goal:** the problem, affected users or maintainers, and intended outcome.
- **Context:** relevant existing behavior, boundaries, and shared decisions.
- **Constraints:** compatibility, invariants, permissions, and other agreed limits.
- **Acceptance Criteria:** checklists of observable outcomes, including relevant
  failure and recovery behavior. Cover repeated, conflicting, or partial cases
  only when relevant. Include a verification method when needed to make a
  criterion precise; do not prescribe a test layer for every criterion.
- **Non-Goals:** deliberate exclusions, not undecided scope.

Keep implementation choices in plan.md. Check a criterion only when evidence
in tasks.md proves it; completion does not imply user acceptance.

## plan.md

Use `# PLAN: <name>` and these sections when they have content:

- **Approach:** how the implementation satisfies the spec; include affected
  boundaries and dependencies.
- **Trade-Offs:** meaningful alternatives, the chosen option, and its reason.
- **Files Affected:** known paths or boundaries and their intended changes.
- **Sequencing:** ordered implementation dependencies and delivery steps.

Keep progress and executable checklists in tasks.md. Sequencing does not override
[TDD.md](../conventions/TDD.md#workflow): write necessary failing tests before
implementation, even when a later step covers broader integration verification.

## tasks.md

Use `# TASKS: <name>`, phase headings named for outcomes, and task checklists.

- Give each task a concrete action and **Done when:** an observable completion
  condition. Link the relevant spec criterion or plan section rather than copying it.
- Order tasks by dependencies in plan.md. Keep phases small and reviewable;
  detail the active phase and leave identified later phases brief.
- Include necessary tests, implementation, verification, and review following
  TDD.md. Do not manufacture failing tests for documentation-only work.
- Keep concise evidence beside the task: command or inspection, result, checked
  revision or working-tree state, and material limits. Include expected failures
  when TDD applies; distinguish failures, skipped checks, and checks not run.
  Link long artifacts instead of embedding logs.
- Record review findings and dispositions, decisions, proposed commit messages,
  and actual commit references beside the affected task. Proposals are not commits.
- Mark `[x]` only when the completion condition and required reviews have evidence.
  Leave partial work `[ ]`; add `in progress`, `blocked: <reason>`,
  `awaiting review`, or `deferred: <reason>` only when useful.
- While work remains, mark one actionable unchecked task or subtask `next`.
  If blocked on a decision or review, make that the next action.
- Record explicit user acceptance with its reference, accepted scope, and remaining
  work. Neither passing checks nor completed tasks imply acceptance or permission
  to commit, publish, or continue beyond approved scope.

## Split large work

- Propose independently verifiable sub-features when work cannot be reviewed in
  small slices. Give each its own ID and sibling folder with the same three files.
- Link parent and children both ways. Keep shared constraints in the parent's
  spec.md and detailed dependencies in plan.md; link blockers from tasks.md.
- Accept the parent only after child work and shared criteria have evidence and
  user acceptance.

## Resume

- Read the active phase, next task, and nearby evidence in tasks.md first. Then
  read only the relevant spec criteria, plan sections, and dependencies.
- Compare recorded progress and checked state with the working tree before acting.
  Resuming does not require reloading record-creation or splitting guidance.
- At review boundaries and session end, update task checks, evidence, blockers,
  and the next action. Update spec.md or plan.md only when agreed scope or approach
  changes; reflect proven progress in ROADMAP.md without duplicating the task list.
