# Decision records

Record significant architecture or technology choices and their reasons.
Name each record `NNNN-<decision-slug>.md`, starting at `0001` and using the
next unused number. Never renumber existing records. Ordinary implementation
details need no decision record.

Use the template below when a real decision needs recording. Proposals remain
proposals until the user agrees. Link relevant [stack choices](../specs/STACK.md)
and [feature records](../features/README.md). Product requirements stay in each
feature's requirements.md, indexed by [FEATURES.md](../specs/FEATURES.md).

## Record template

Copy this section-only template into the numbered file, give it the actual
number and title, and add only established facts and agreed decisions. Leave
unknown sections empty. For requested spec work, use the question tool with
proposed options, following [AGENTS.md](../../AGENTS.md).

```markdown
# Decision

## Status

## Date

## Decision reference

## Related work

## Context

## Options considered

## Decision

## Consequences
```

Record the problem and constraints in Context, realistic alternatives and
tradeoffs in Options considered, the selected option and rationale in Decision,
and benefits, costs, limits, and follow-up work in Consequences. Include the
actual user approval or source in Decision reference; never imply agreement.

Statuses: proposed, accepted, rejected, superseded. To replace an accepted
decision, create a new record and link both directions; retain the old rationale.
