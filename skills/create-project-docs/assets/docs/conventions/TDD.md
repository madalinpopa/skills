# Test-driven changes

Test every visible behavior change, including fixes. Assert what users/callers
rely on, not field assignments, library behavior, private structure, or coverage.

## Coverage and readability

- Prefer the smallest integration boundary proving behavior, with real local
  dependencies and small fixtures where practical. Avoid full end-to-end setups
  when smaller boundaries suffice.
- Give each behavior one main test; inspect and extend existing coverage rather
  than duplicating scenarios across layers. Add focused unit tests only for logic
  integration tests miss; state the distinct gap. If integration is impractical
  or unreliable, use a unit test and record the integration gap/reason.
- Keep default tests deterministic, network-free, and independent of shared services.
  Gate external-service checks and document their commands.
- Make setup, action, and outcome obvious. Avoid generic harnesses and complicated
  stubs; use minimal stubs when real dependencies are impractical.
- Reuse small, descriptive helpers when clearer; keep tested behavior and expected
  results visible at the call site.
- Use small tables only for cases sharing setup/assertions, with descriptive names
  and explicit inputs/results. Separate success/failure groups; avoid flag-driven
  assertion branches or unrelated cases in one table.

## Workflow

Work one slice at a time; the user reviews once at its end.

1. Write/run the smallest useful failing test. Confirm missing behavior causes
   failure, not setup/environment; retain the command and output.
2. Implement the smallest clear change within approved scope; refactor only where
   useful while preserving tests and agreed behavior.
3. Run focused tests, then project format, lint, full-test, and build checks in
   [STACK.md](../specs/STACK.md). Run integration checks when the changed boundary needs them.
4. Delegate [adversarial review](#adversarial-review) to a different agent; wait
   before continuing or requesting user review.
5. Resolve findings below; after fixes, rerun affected checks and repeat step 4.
6. Once clean, stop for user review with changes, failing-test evidence, check
   results, findings/fixes, remaining concerns, and a proposed commit message.
   Commit only with approval for that commit.

Plan approval neither skips review nor authorizes commits; reuse explicit approval
within scope. Skip step 1 for documentation/configuration that gains nothing from
failing tests; behavior-changing configuration still needs tests. Report every
failed/skipped check and its cause; never claim an unrun check passed.

## Adversarial review

Default to one independent, read-only reviewer via the available subagent tool.
Provide the diff, approved scope, requirements, test evidence, and necessary code
context; omit implementer conclusions.

Reviewer instructions:

- Challenge correctness, regressions, boundary handling, security, and test adequacy.
  Inspect relevant callers, contracts, and existing tests beyond the diff; check
  whether tests fail for broken behavior.
- Verify findings with a reproducer, failing check, or clear code path. Report
  actionable issues introduced/exposed by this slice, each root cause once;
  distinguish questions from confirmed defects.
- Order findings by severity based on impact/reachability. Include location,
  trigger, impact, and evidence; distinguish blockers from optional improvements.
- Exclude style preferences; maintainability findings need concrete impact.
  Do not block on polish or equally valid designs.
- Return concise findings or “No actionable findings”, with material limitations,
  unrun checks, and relevant context not inspected.

Verify each finding; mark it fixed, rejected with evidence, or awaiting a user
decision. Do not repeat rejected findings without new evidence. Record findings,
dispositions, and checked state in the feature's validation.md when present;
avoid repeated review logs in replies.

Repeat review after fixes until no blocking findings remain; later code edits
require renewed review. Optional improvements do not block completion. If review
repeats a disputed finding without new evidence or stalls on blockers, stop and
request a user decision with evidence and remaining disagreement. Never claim
clean review with unresolved blockers; clean review does not prove defect-free code.

Review authorizes neither new scope nor dependency changes. Record unresolved
questions/decisions in requirements/plans and ask with the available question tool
before dependent work. If delegation is unavailable, report the blocker and request
a decision; do not substitute self-review or claim clean review.
