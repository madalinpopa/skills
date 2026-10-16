# Decision records

## When to create a record

- Record significant architecture or technology choices and their reasons.
- Skip ordinary implementation details; do not create examples during setup.
- Name records `NNNN-<decision-slug>.md` using the next unused number.
- Preserve numbering and names; never renumber existing records.
- Proposals remain proposals until the user agrees.
- Record unresolved questions in [the question log](../QUESTIONS.md), following
  the [question guide](../conventions/QUESTIONS.md); link the relevant group from the record.

## Writing a record

- Start empty; add the actual number, title, and sections only with known content.
- Follow the [spec writing guide](../specs/README.md): group related facts and
  write concise lists, without placeholders or empty outlines.
- Use the question tool with proposed options for requested spec work.
- Status and date: record the actual state and date of the decision.
- Decision reference: cite user approval or the source of the agreed choice.
- Related work: link relevant stack choices, feature records, or roadmap tasks.
- Context: list the problem, constraints, and reason a decision is needed.
- Options considered: list realistic alternatives and their tradeoffs.
- Decision: state the selected option and reason; never imply approval.
- Consequences: list benefits, costs, limits, and follow-up work.

## Ownership and revisions

- Keep technology choices in [STACK.md](../specs/STACK.md).
- Keep product behavior in [feature requirements](../features/README.md).
- Keep feature groups and tasks in [ROADMAP.md](../specs/ROADMAP.md).
- Use statuses: proposed, accepted, rejected, superseded.
- Replace accepted decisions with a new linked record; retain the old rationale.
