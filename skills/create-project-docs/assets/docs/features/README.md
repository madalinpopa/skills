# Feature records

Create one folder per real feature, using the next unused number and a
meaningful slug. Preserve existing numbering and names.

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

[FEATURES.md](../specs/FEATURES.md) indexes capabilities and links to their
requirements; [ROADMAP.md](../specs/ROADMAP.md) orders delivery and links to
plans. Link shared decisions and constraints rather than duplicating them.

Use the templates below for a real feature. Replace pending fields with known
facts or explicit open questions. Detail only active work in plan.md; outline
future steps briefly. Each step has one reviewable outcome and follows
[TDD.md](../conventions/TDD.md). Keep the three files as the feature's records
rather than creating a competing FEATURE.md or separate phase files.

## requirements.md template

```markdown
# Feature 0001 requirements: Pending title

- Capability ID: Pending: identifier from ../../specs/FEATURES.md.
- Status: proposed
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

Give criteria stable IDs for tests and validation evidence. Recovery mechanisms
such as retries or rollback are not implied requirements.

## Open questions

Pending: unresolved requirements; do not invent behavior to fill a section.

Execution belongs in [plan.md](plan.md); evidence belongs in
[validation.md](validation.md).
```

## plan.md template

```markdown
# Feature 0001 plan: Pending title

- Requirements: [requirements.md](requirements.md)
- Evidence and acceptance: [validation.md](validation.md)

## Approach and boundaries

Pending: implementation approach, affected areas, dependencies, and relevant
existing code or decisions. Keep product requirements in requirements.md.

## Delivery steps

Pending: one reviewable outcome per step, relevant requirement IDs, necessary
failing tests, implementation scope, and proposed test/implementation messages.
Detail only the active step; outline later work. Present the plan before changes
and follow ../../conventions/TDD.md for separate test and implementation review.
Existing explicit approval counts only within its stated scope.

## Current checkpoint

- Status: planned
- Active step: Pending
- Approved scope: Pending: exact user decision and scope.
- Working tree: Pending: changed files and unrelated work to preserve.
- Last verification: not run; link evidence in validation.md when available.
- Pending question or blocker: Pending
- Next action: Pending: one concrete action, including required approval.
```

## validation.md template

```markdown
# Feature 0001 validation: Pending title

- Requirements: [requirements.md](requirements.md)
- Execution plan: [plan.md](plan.md)

## Verification evidence

| Requirement / step | Scenario or check | Command / inspection | Result | Checked state |
| --- | --- | --- | --- | --- |
| Pending | Pending | Not run | Not verified | Pending: revision and relevant uncommitted changes |

Link requirement IDs to observable outcomes and tests. Retain the expected
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
end, and mirror delivery status in ROADMAP.md. Replace the checkpoint snapshot
rather than appending a diary. Load historical records only for a dependency.
