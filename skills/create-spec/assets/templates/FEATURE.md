# Feature: [Name]

[Describe the problem, who it affects, and the intended outcome in 2–3 sentences.]

- **Active phase**: [Link to the active phase, or none]
- **Feature acceptance**: [Pending / accepted, with the developer's decision reference]

## Using these documents

Copy this template to `docs/features/<feature-name>/FEATURE.md`. Keep future
phases as brief roadmap entries; create a sibling phase file from
`docs/templates/PHASE.md` when that phase becomes active. Replace placeholders,
remove unused sections, and resolve repository paths from the repository root.

`AGENTS.md` owns engineering rules and verification commands. `docs/SPEC.md`
owns CLI behavior. This file records the feature agreement; phase files record
execution and evidence. Surface conflicts for a developer decision before
changing the affected behavior or contract.

At the start of a session, read repository instructions, this overview, the
active phase, and its linked spec sections. Inspect relevant code and tests as
needed. Load other phases only for a specific dependency. Link to longer
research and historical evidence instead of copying it into these documents.

Work on one phase and one reviewable change at a time. Use the phase's approval
checkpoints. Existing explicit authorization counts within its stated scope;
do not infer permission for later phases. Commits and PRs require separate
explicit requests under `AGENTS.md`.

## Scope

- **Included**: [Behavior this feature delivers]
- **Non-goals**: [Related behavior deliberately excluded]
- **Constraints**: [Compatibility requirements and invariants to preserve]
- **Spec references**: [Links to relevant sections of `docs/SPEC.md`]

## Acceptance criteria

Describe observable outcomes. Link each criterion to the phase evidence that
proves it; checked tasks and passing commands alone do not imply acceptance.

- [ ] [Expected behavior and link to verification evidence]

## Agreed decisions

Record only decisions that affect future work. Distinguish proposals from
developer-approved decisions; retain a short reason and decision reference.

| Decision | Reason | Approval or pending question |
| --- | --- | --- |
| [Decision] | [Why it matters] | [Developer statement/date or pending] |

## Phase roadmap

Each phase delivers one bounded, observable outcome with its own verification.
Link phase files when created. The phase checkpoint owns detailed state; this
table mirrors its status at review boundaries.

| Phase | Outcome | Depends on | Status | Phase file |
| --- | --- | --- | --- | --- |
| 01 | [Reviewable outcome] | None | planned | [Link when created] |
| 02 | [Next outcome] | 01 | planned | [Link when created] |

Statuses: `planned`, `awaiting test review`, `implementing`,
`awaiting implementation review`, `accepted`, `blocked`.

Accept the feature only after all phases are accepted, feature criteria have
evidence, relevant behavior documentation is current, and the developer
explicitly accepts the result. Acceptance does not mean committed or published.
