# Test-driven changes

For production work, test every visible behavior change, including fixes.
Assert what users/callers rely on, not field assignments, library behavior,
private structure, or coverage.

## Throwaway prototypes

Use this workflow when the user explicitly requests a disposable prototype to
explore or learn about an idea, feature, design, or architecture. Also use it
when the user asks to update or continue files that repository evidence identifies
as an existing throwaway prototype, such as its README, dedicated directory, or
branch/bookmark context. A filename alone is insufficient; clarify uncertain
intent before bypassing the normal workflow.

- Implement the smallest useful experiment quickly. Skip TDD, automated test
  writing, the normal verification suite, adversarial review, and formal feature
  records. Do not require intermediate plan or review gates for this workflow.
- Preserve unrelated work. Use a separate `prototype/<name>` Git branch, or a
  separate jj change with a `prototype/<name>` bookmark; follow the project's
  version-control tooling. Reuse suitable existing prototype isolation when
  continuing work; isolate prototype files first if they are mixed with production.
- A request covered by this workflow authorizes local prototype commits under
  [COMMITS.md](COMMITS.md#scope-and-approval). Review the diff to include only the
  experiment; keep its bookmark pointing to the latest prototype commit.
- Run only the build or manual smoke check needed to demonstrate the idea.
  Report the check, what works, limitations, lessons learned, and branch/bookmark
  and commit references. Do not claim production readiness or passing tests.
- Keep dependency and external-action permissions. Push, publish, merge, and
  deletion still require explicit user requests. Discard prototype code when
  requested, preserving unrelated work.
- Treat production implementation as fresh work under the normal workflow,
  including requirements, tests, and review. Prototype results are evidence for
  decisions, not agreed requirements or completed production features.

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

For throwaway prototype work, use [its workflow](#throwaway-prototypes) instead.
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
questions in [the question log](../QUESTIONS.md), following the
[question guide](QUESTIONS.md), and link them from affected requirements/plans.
Ask with the available question tool
before dependent work. If delegation is unavailable, report the blocker and request
a decision; do not substitute self-review or claim clean review.
