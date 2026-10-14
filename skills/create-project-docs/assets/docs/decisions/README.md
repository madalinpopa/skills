# Decision records

Record significant architecture or technology choices and their reasons.
Create `0001-<decision-slug>.md`, then use the next unused number. Preserve
existing numbering; do not renumber history. Ordinary implementation details
need no decision record.

Use the template below when a real decision needs recording. Proposals remain
proposals until the user agrees. Link relevant [stack choices](../specs/STACK.md)
and [feature records](../features/README.md). Product requirements stay in each
feature's requirements.md, indexed by [FEATURES.md](../specs/FEATURES.md).

## Record template

Copy this template into the numbered file and replace pending fields:

```markdown
# Decision 0001: Pending title

- Status: proposed
- Date: Pending
- Decision reference: Pending: user approval or source of the agreed decision.
- Related work: Pending: relevant feature and spec links.

## Context

Pending: problem, constraints, and why a decision is necessary.

## Options considered

Pending: realistic alternatives and their tradeoffs.

## Decision

Pending: selected option and reason; do not imply approval before it exists.

## Consequences

Pending: benefits, costs, limitations, and follow-up obligations.
```

Statuses: proposed, accepted, rejected, superseded. To replace an accepted
decision, create a new record and link both directions; retain the old rationale.
