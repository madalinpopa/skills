# Feature records

Create one folder per real feature, named `NNNN-<slug>`. `NNNN` is the
feature's number from [FEATURES.md](../specs/FEATURES.md) and its only ID.
Preserve existing numbering and names.

```text
0001-feature-x/
  plan.md
  requirements.md
  validation.md
```

Use `validation.md` with one dot before the extension. Create these records
only when a real feature is being planned, not as empty examples during setup.

## Document ownership

- `requirements.md`: feature scope, behavior, constraints, and acceptance criteria.
- `plan.md`: implementation approach, delivery steps, approvals, and current checkpoint.
- `validation.md`: verification evidence and the user's acceptance decision.

[FEATURES.md](../specs/FEATURES.md) indexes features and links to their
requirements, without tasks or progress. [ROADMAP.md](../specs/ROADMAP.md)
contains only feature and sub-feature groups with task/subtask checklists. Link
tasks to plan.md for execution details and validation.md for evidence when those
records exist. Keep acceptance decisions in validation.md, dependencies and the
checkpoint in plan.md, and behavior in requirements.md. Do not add status fields,
acceptance prose, or workflow notes to the roadmap. Link shared facts rather
than duplicating them. Use roadmap order for delivery order.

## Writing structure

- Use the template sections in a logical order: purpose and scope, behavior,
  acceptance criteria, then open questions; approach, delivery, then checkpoint.
- Write short bullets under sections. Split substantial behavior contracts into
  subsections for inputs, outcomes, state, and failures when needed.
- Use tables only for compact comparisons or mappings; keep requirements,
  rationale, tasks, and verification evidence in sectioned lists.
- Use `[ ]` task and subtask checklists for all plans and delivery sequences,
  including future work. Give every task a short outcome description and each
  subtask a concrete action. Group related tasks into named phases.
- Mark a subtask `[x]` when its action is complete and link relevant evidence.
  Mark its parent `[x]` only after all required subtasks and its user review are
  complete. Leave blocked, deferred, or unverified work unchecked with a reason.
- Update checklists at review boundaries and session end. Task completion does
  not imply feature acceptance or permission to commit.

Use the section-only templates below for a real feature. Add its actual number
and title, then fill sections progressively with established facts and user
answers. Leave unknown sections empty; do not copy writing guidance into them.
For requested spec work, always use the available question tool with proposed
options as described in [AGENTS.md](../../AGENTS.md); ask only about the requested
scope and wait for answers before filling dependent sections. In plan.md, detail
the tasks for agreed phases; add future work only when identified.
Each commit has one reviewable outcome and follows
[TDD.md](../conventions/TDD.md). Keep the three files as the feature's records
rather than creating a competing FEATURE.md or separate phase files.

## Plan phases and commits

- Group related implementation tasks into phases with a clear outcome and
  dependencies. Link tasks to the small delivery slices in ROADMAP.md; a phase
  may group related slices. Keep phases as sections in plan.md.
- Plan one to five commits per phase. Split a phase needing more than five into
  smaller related phases; do not add filler commits to meet a count.
- Give each planned commit a checklist task, a proposed message, affected areas,
  relevant acceptance criteria, and detailed action subtasks for tests,
  implementation, verification, and review as applicable.
- Describe concrete work, such as the behavior to add, the code boundary to
  change, and the scenario to verify; replace generic workflow labels with
  feature-specific steps when writing a real plan.
- Record the planned commit count per phase and respect project policy on whether
  tests and implementation share a commit or need separate commits.
- Mark a phase complete only after all required tasks and reviews are complete.
  Track whether each commit is merely planned, ready for approval, or actually
  committed; checked work alone does not claim a commit exists.

## Split oversized features

If a feature is too large to implement and review in small slices, split it into
sub-features before implementation. Give each a clear, independently verifiable
outcome, bounded requirements, and its own review checkpoints.

Give each sub-feature its own number in FEATURES.md and a sibling folder
containing `plan.md`, `requirements.md`, and `validation.md`. Link the parent
and sub-features in both directions and record their dependencies and delivery
order through its grouping and task order in ROADMAP.md; put dependency details
in plan.md.
Keep the parent focused on the overall outcome and shared constraints; link to
the sub-feature owners instead of duplicating detailed requirements or evidence.

Present the proposed split for user review. Implement one sub-feature at a time
in small slices. Accept the parent only after its sub-features and any shared
end-to-end acceptance criteria have evidence and user acceptance.

## Filling the records

- Requirements: record the actor, outcome, included scope, non-goals, constraints,
  and explicit approval when known. Describe observable behavior, meaningful
  failures, state invariants, and repeated/conflicting/partial cases only where
  relevant. Retries and rollback are not implied requirements.
- Give concrete acceptance scenarios stable IDs such as `R-001`; cite them as
  `NNNN/R-001` in tests and validation. Do not create criteria just to fill space.
- Plans: add sections for actual agreed phases and planned commits, including
  outcome, criteria, dependencies, commit count, proposed messages, affected
  areas, and concrete task/subtask checklists. Include necessary failing tests,
  implementation, verification, adversarial review, and user review per TDD.md.
  Omit failing tests for work that gains nothing from them, as TDD.md permits.
- Present plans before changes. Approval applies only within its stated scope;
  it does not authorize commits or remove per-slice review boundaries. Execute,
  verify, and present one small slice before starting the next, including docs,
  configuration, and refactoring.
- Checkpoints: record the active phase/commit/task, approved scope, working-tree
  changes and unrelated work to preserve, last verification, actual blockers,
  and one exact next action. Do not fill unknown fields with status guesses.
- Validation: link each criterion or task to its scenario, command or inspection,
  expected failure and final result, code checked, and actual limitations.
  Distinguish change failures, environment problems, skipped checks, and checks
  not run. Passing commands alone do not prove every criterion.
- User acceptance: record only an actual decision, its reference and scope,
  accepted criteria, and any remaining work. Approval does not authorize commits
  or publication. Keep acceptance separate from task completion.

## requirements.md template

```markdown
# Feature requirements

## Agreement

## Problem and scope

### Actor and intended outcome

### Included

### Non-goals

### Constraints

## Behavior contract

### Preconditions and input rules

### Observable results and side effects

### State and invariants

### Repeated, conflicting, or partial execution

### Failures and recovery

## Acceptance criteria

## Open questions

## Related records
```

## plan.md template

```markdown
# Feature plan

## Related records

## Approach and boundaries

### Approach

### Affected areas

### Dependencies and decisions

## Delivery phases

## Future work

## Current checkpoint

### Active phase, planned commit, and task

### Approved scope

### Working tree

### Last verification

### Questions and blockers

### Next action
```

## validation.md template

```markdown
# Feature validation

## Related records

## Verification evidence

## User acceptance

### Decision

### Reference and scope

### Remaining work
```

On resume, read requirements.md and the checkpoint in plan.md, then relevant
validation evidence and linked contracts. Compare the checkpoint with the
working tree. Update plan.md and validation.md at review boundaries and session
end, and update the corresponding roadmap checklists from that evidence. Replace
the checkpoint snapshot rather than appending a diary. Load historical records
only for a dependency.
