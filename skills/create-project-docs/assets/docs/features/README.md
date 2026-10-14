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
requirements; [ROADMAP.md](../specs/ROADMAP.md) owns delivery order and status
and links to plans. Keep feature status in ROADMAP.md; track task and subtask
completion in plan.md and acceptance evidence in validation.md. Link shared
decisions and constraints rather than duplicating them.

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

Use the templates below for a real feature. Replace pending fields with known
facts or explicit open questions. In plan.md, detail the tasks for every agreed
phase; keep unagreed future work as a brief checklist with pending decisions.
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
order in ROADMAP.md.
Keep the parent focused on the overall outcome and shared constraints; link to
the sub-feature owners instead of duplicating detailed requirements or evidence.

Present the proposed split for user review. Implement one sub-feature at a time
in small slices. Accept the parent only after its sub-features and any shared
end-to-end acceptance criteria have evidence and user acceptance.

## requirements.md template

```markdown
# Feature NNNN requirements: Pending title

- Agreement: Pending: user decision and approved scope.

## Problem and scope

- Actor and intended outcome: Pending
- Included: Pending
- Non-goals: Pending
- Constraints: Pending: relevant mission, stack, and decision links.

## Behavior contract

- Preconditions and input rules: Pending
- Observable results and side effects: Pending
- State and invariants: Pending
- Repeated, conflicting, or partial execution: Pending: only where relevant.
- Failures and recovery: Pending: errors, preserved state, and user action.

## Acceptance criteria

- R-001: Pending: concrete Given/When/Then outcome.
- R-002: Pending: meaningful failure or boundary case, if applicable.

Keep criterion IDs stable. Cite them in tests, validation, and other records
as `NNNN/R-001`. Recovery mechanisms such as retries or rollback are not implied
requirements.

## Open questions

- Pending: unresolved requirements; do not invent behavior to fill a section.

Execution belongs in [plan.md](plan.md); evidence belongs in
[validation.md](validation.md).
```

## plan.md template

```markdown
# Feature NNNN plan: Pending title

- Requirements: [requirements.md](requirements.md)
- Evidence and acceptance: [validation.md](validation.md)

## Approach and boundaries

- Approach: Pending: smallest implementation approach.
- Affected areas: Pending: relevant existing code and boundaries.
- Dependencies and decisions: Pending: links to agreed choices.
- Keep product requirements in requirements.md.

## Delivery phases

### Phase 1: Pending: related outcome

- Outcome and criteria: Pending: observable result and `NNNN/R-001` links.
- Dependencies: Pending: prerequisites or earlier phases.
- Planned commits: Pending: 1–5; follow project commit policy.
- [ ] Complete phase: all required tasks and reviews completed with evidence.

#### Planned commit 1: Pending: proposed commit message

- Affected areas: Pending: concrete files, packages, or boundaries.
- Commit state: planned; record user authorization and commit reference when made.
- [ ] Task: Pending: one reviewable outcome and relevant acceptance criteria.
  - [ ] Agree scope: Pending: exact user approval and boundaries.
  - [ ] Prove the gap: Pending: concrete scenario, necessary failing tests, and
    failure evidence; prefer integration coverage and identify distinct unit gaps.
  - [ ] Implement: Pending: specific changes needed to produce the outcome.
  - [ ] Verify: Pending: commands, expected results, and validation.md evidence.
  - [ ] Review: Pending: adversarial review, confirmed fixes, and review until clean.
  - [ ] User review: Pending: decision and any requested revisions.

Repeat the planned-commit section for the remaining commits, up to five per
phase. Add a phase section for each related group of agreed tasks, with detailed
subtasks. Keep unagreed later work separate:

### Future work awaiting agreement

- [ ] Pending: possible later outcome and the decision needed before planning.

### Execution rules

- Present the plan before changes. Follow ../../conventions/TDD.md: failing test,
  implementation, adversarial review until clean, then user review.
- Reuse explicit approval only within its stated scope. Plan approval does not
  remove per-slice review boundaries or authorize commits.
- Execute one small slice at a time, including docs, configuration, and
  refactoring. Verify and present it before starting the next slice.
- Omit the failing-test subtask for work that gains nothing from one, as allowed
  by TDD.md. Keep completion evidence beside tasks or linked in validation.md.

## Current checkpoint

- Active phase, planned commit, and task: Pending
- Approved scope: Pending: exact user decision and scope.
- Working tree: Pending: changed files and unrelated work to preserve.
- Last verification: not run; link evidence in validation.md when available.
- Pending question or blocker: Pending
- Next action: Pending: one concrete action, including required approval.
```

## validation.md template

```markdown
# Feature NNNN validation: Pending title

- Requirements: [requirements.md](requirements.md)
- Execution plan: [plan.md](plan.md)

## Verification evidence

### Criterion or task: Pending: `NNNN/R-001` or plan task

- Scenario or check: Pending: observable outcome.
- Command or inspection: Not run.
- Expected failure: Pending: command and failure for the intended reason.
- Final result: Not verified.
- Code checked: Pending: commit plus any uncommitted changes.
- Limitations: Pending: skipped checks, environment problems, or coverage gaps.

Link criteria to observable outcomes and tests. Retain the expected
failing-test evidence and subsequent passing results. Distinguish change failures
from environment problems and skipped checks; passing commands alone do not
prove every criterion. Record material limitations and unverified scenarios.

## User acceptance

- Decision: pending
- Reference and scope: Pending: explicit user decision and accepted criteria.
- Remaining work: Pending: unmet criteria or requested revisions.

Approval does not authorize commits or publication.
```

On resume, read requirements.md and the checkpoint in plan.md, then relevant
validation evidence and linked contracts. Compare the checkpoint with the
working tree. Update plan.md and validation.md at review boundaries and session
end, and update the feature's status in ROADMAP.md. Replace the checkpoint snapshot
rather than appending a diary. Load historical records only for a dependency.
