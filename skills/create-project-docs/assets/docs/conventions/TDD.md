# Test-driven changes

Apply to every change in visible behavior, including bug fixes. Test what users
and callers rely on, not field assignments, library behavior, private details,
or coverage numbers.

## Balance test coverage

- Prefer integration tests for observable behavior across the relevant components,
  using real local dependencies and small fixtures where practical.
- Give each behavior one main test. Inspect existing coverage before adding tests;
  extend or reuse it instead of repeating the same scenarios and assertions at
  both unit and integration levels.
- Add small unit tests only to prove implementation logic that integration tests
  do not adequately exercise, such as focused branches or boundary cases. State
  the distinct gap each unit test covers; avoid testing private structure.
- Keep the suite balanced: use the smallest integration boundary that proves the
  behavior, rather than a full end-to-end environment for every case. If that
  boundary is impractical or unreliable, use a focused unit test and record the
  integration gap and reason.
- Keep default runs deterministic and independent of network or shared services.
  Gate checks requiring external services and document how to run them.

## Readable tests

Readability is the priority. Make each test's setup, action, and expected outcome
easy to follow without tracing through a test framework or layers of helpers.

- Keep setup and assertions direct. Avoid unnecessary abstractions, generic test
  harnesses, and complicated stubs; use simple fixtures and minimal stubs when
  a real dependency cannot reasonably be exercised.
- Use helpers for repeated setup, fixture creation, or assertions when they make
  the test clearer. Reuse suitable existing helpers rather than duplicating them.
  Give helpers descriptive names and keep their responsibilities small; leave
  the behavior under test and its expected outcome visible at the call site.
- Use table tests when cases share setup and assertion structure. Give each case
  a descriptive name and explicit inputs and expected results. Keep tables small;
  do not force unrelated scenarios into one table.
- Group table tests into separate success and failure tests. Keep each group's
  assertions consistent rather than branching on flags inside a combined runner.

## Workflow

Work one slice at a time. The user reviews once, at the end of each slice.

1. Write the smallest useful failing test. Run it and confirm it fails because
   the required behavior is missing, not because of an environment or setup
   error. Keep the command and failure output as evidence.
2. Implement the smallest clear change that makes it pass, within the approved
   scope. Refactor only where useful, keeping tests and agreed behavior intact.
3. Run focused tests, then the format, lint, full-test, and build checks in
   [STACK.md](../specs/STACK.md). Run integration tests only when the changed
   boundary needs them.
4. Delegate the [adversarial review](#adversarial-review) to a different agent
   and wait for its result. Do not continue or ask for user review meanwhile.
5. If it reports issues, verify them, fix the confirmed ones, rerun affected
   checks, and go back to step 4 with the updated diff.
6. When the review is clean, stop for user review. Show the changes, the
   failing-test evidence, final check results, review findings and fixes,
   remaining concerns, and a proposed commit message. Commit only after the
   user approves that commit.

Plan approval does not skip the user review or authorize commits. Reuse explicit
approval within its stated scope. Documentation/configuration work that gains
nothing from a failing test skips step 1; configuration that changes behavior
still needs tests.

Keep tests deterministic and default runs independent of network or shared
services. Arrange setup, action, and assertions clearly. Report every failed
or skipped check and its cause; never describe an unrun check as passed.

## Adversarial review

Delegate to an independent, read-only reviewer using the environment's available
subagent tool. Supply the current diff, approved scope, relevant requirements,
test evidence, and necessary code context; omit the implementer's conclusions.

Reviewer brief: challenge correctness, regressions, boundary handling, security,
and test adequacy. Report only actionable issues introduced or exposed by this
slice, with location, triggering scenario, impact, and evidence. Skip stylistic
preferences. Return a concise findings list, or "No actionable findings", and
identify material review limitations.

Repeat the review until the final code state has no unresolved actionable
findings. Any later code edit requires renewed review. Record findings,
dispositions, and checked state in the feature's validation.md when present;
do not paste repeated review logs into user replies.

Review does not authorize new requirements or dependency changes. Record
unresolved questions or decisions in the relevant requirements/plan and ask the
user with the available question tool before dependent work. If delegation is
unavailable or the loop cannot resolve a finding, report the blocker and request
a decision; do not substitute self-review or claim a clean review. A clean
review is evidence, not proof of defect-free code.
