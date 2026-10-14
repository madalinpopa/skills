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
and links to plans. Do not record status in these three files. Link shared
decisions and constraints rather than duplicating them.

Use the templates below for a real feature. Replace pending fields with known
facts or explicit open questions. Detail only active work in plan.md; outline
future steps briefly. Each step has one reviewable outcome and follows
[TDD.md](../conventions/TDD.md). Keep the three files as the feature's records
rather than creating a competing FEATURE.md or separate phase files.

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

Pending: unresolved requirements; do not invent behavior to fill a section.

Execution belongs in [plan.md](plan.md); evidence belongs in
[validation.md](validation.md).
```

## plan.md template

```markdown
# Feature NNNN plan: Pending title

- Requirements: [requirements.md](requirements.md)
- Evidence and acceptance: [validation.md](validation.md)

## Approach and boundaries

Pending: implementation approach, affected areas, dependencies, and relevant
existing code or decisions. Keep product requirements in requirements.md.

## Delivery steps

Pending: one reviewable outcome per step, relevant `NNNN/R-001` criteria,
necessary failing tests, implementation scope, and a proposed commit message.
Detail only the active step; outline later work. Present the plan before changes
and follow the workflow in ../../conventions/TDD.md: failing test,
implementation, adversarial review until clean, then user review.
Existing explicit approval counts only within its stated scope.
Execute one small slice at a time. Keep its diff focused, verify it, and present
the result for user review before starting the next slice. Approval of the whole
plan does not remove these review boundaries. Apply the same approach to docs,
configuration, and refactoring; do not bundle later steps into the active slice.

## Current checkpoint

- Active step: Pending
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

| Criterion (`NNNN/R-001`) or step | Scenario or check | Command / inspection | Result | Code checked |
| --- | --- | --- | --- | --- |
| Pending | Pending | Not run | Not verified | Pending: commit plus any uncommitted changes |

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
