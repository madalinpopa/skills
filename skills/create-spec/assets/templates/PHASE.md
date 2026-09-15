# Phase [Number]: [Outcome]

- **Feature**: [Feature agreement](FEATURE.md)
- **Depends on**: [None, or phase link and the specific result needed]
- **Outcome**: [One observable result this phase delivers]

## Current checkpoint

Replace this snapshot at review boundaries and before ending a session; do not
append a session diary. On resume, compare it with the current working tree and
relevant code before proceeding. Preserve unrelated work. Reuse recorded
approval only for the scope it covers; never infer approval from a status.

- **Status**: [planned / awaiting test review / implementing / awaiting implementation review / accepted / blocked]
- **Current step**: [Step number and whether tests, implementation, or review is next]
- **Approved scope**: [Exact scope and developer statement/date; none until approved]
- **Working tree**: [Changed paths, uncommitted work, and unrelated edits to preserve]
- **Last verification**: [Command, result, and evidence reference; not run if absent]
- **Blocker or pending decision**: [Include expected failing tests, or none]
- **Next action**: [One concrete action; state whether developer approval is required]

## Scope and required context

- **Included**: [Behavior covered by this phase]
- **Non-goals**: [Related changes excluded from this phase]
- **Expected affected areas**: [Packages/files and their responsibilities]

Read repository instructions and the feature overview, then these references.
Use section anchors or symbols where useful. Inspect additional code when
needed to validate the actual behavior; avoid loading unrelated phase history.

| Reference | Why it is needed |
| --- | --- |
| [Relevant `docs/SPEC.md` section link] | [Contract this change must satisfy] |
| [Code/test path and symbol] | [Current behavior or reusable test setup] |
| [Dependency decision or external documentation link, if needed] | [Constraint or API to verify] |

## Acceptance criteria

Use observable behavior, including meaningful failure cases. Keep verification
commands in the evidence section. Do not weaken agreed criteria to make the
implementation pass; propose any scope change for developer review.

- [ ] [Given a specific input/state, the expected output or side effect occurs]
- [ ] [Relevant failure leaves the required state and reports the expected error]

## Change plan and review checkpoints

Follow `AGENTS.md` for engineering rules, test selection, and verification.
Confirm the phase scope with the developer before editing; an existing explicit
request covering it is sufficient. Make routine implementation decisions within
that scope. Stop for a decision when findings require changing agreed behavior,
acceptance criteria, or scope. Read-only investigation may continue while waiting.

### Step [Number]: [One reviewable change]

- **Change and reason**: [Responsibility/behavior to change and why]
- **Test cases**: [Smallest useful cases, test layer, and expected failure]
- **Test commit proposal**: `test(<scope>): <behavior covered>`
- **Implementation commit proposal**: `feat(<scope>): <outcome>` or `fix(<scope>): <outcome>`

1. Write the smallest useful tests and run them. Confirm they fail because the
   required behavior is missing, rather than a setup or environment error.
2. **Stop for test review.** Show the test diff, command, relevant failure output,
   proposed commit message, and exact implementation scope awaiting approval.
   Set status to `awaiting test review`; do not write implementation yet.
3. After explicit developer approval, record it in the checkpoint and set status
   to `implementing`. Write the smallest implementation within the approved scope.
4. Run the verification required by `AGENTS.md` after this implementation. Review
   the final diff and record results below, including failures and skipped checks.
5. **Stop for implementation review.** Show the diff, acceptance evidence,
   remaining concerns, proposed commit message, and proposed next action.
   Set status to `awaiting implementation review`.

Repeat this step structure only for necessary changes. Do not advance to another
step or phase without developer authorization covering that work. For config,
docs, or release work that gains nothing from a failing test, record why the
`AGENTS.md` exception applies, omit the test gate, and propose one commit message.
The implementation review still applies. Never commit or open a PR without an
explicit request; propose a PR title when requested using the actual reviewed work.

## Verification evidence

Keep concise results and the failure excerpt needed for review; link to durable
artifacts for long output. Tie results to the tested working-tree state, such as
a revision plus a description of uncommitted changes. Results predating relevant
edits do not verify those edits. Distinguish change failures from environment or
unavailable integration dependencies.

| Step / criterion | Command or inspection | Result and relevant evidence | Tested state |
| --- | --- | --- | --- |
| [Reference] | [Exact command or inspected behavior] | [Pass / expected fail / fail / not run, with reason] | [Revision and uncommitted changes] |

## Developer acceptance

- **Decision**: [Pending / accepted / changes requested]
- **Reference and scope**: [Developer statement/date and the result it accepts]
- **Remaining work**: [Unfinished step, requested changes, or none]

Set the phase to `accepted` only when all criteria have evidence, required
verification is complete, and the developer explicitly accepts the phase.
Update the roadmap in `FEATURE.md` at review boundaries. Checked tasks or passing
tests alone do not authorize acceptance, commits, or the next phase.
