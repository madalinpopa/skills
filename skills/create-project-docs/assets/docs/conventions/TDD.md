# Test-driven changes

Apply to observable behavior changes, including bug fixes. Choose one main test
at the lowest useful layer; add integration coverage only for boundaries that
lower-level tests cannot prove. Test realistic contracts, not field assignments,
library behavior, private implementation details, or coverage numbers.

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

## Review workflow

Before either user review gate below, run the adversarial review for the changed
code slice, including tests and scripts. Expected failing tests at the red stage
are evidence, not implementation defects.

1. Write the smallest useful failing test. Run it and confirm it fails because
   required behavior is missing, rather than an environment or setup error.
2. Stop for test review: show the test changes, command, failure evidence, and a
   proposed `test(<scope>): ...` commit message. Wait for explicit approval.
3. Implement the smallest clear change within the approved scope.
4. Refactor only where useful, preserving passing tests and agreed behavior.
5. Run focused tests, then applicable format, lint, full-test, and build checks
   recorded in [STACK.md](../specs/STACK.md). Run integration tests only when the
   changed boundary requires them.
6. Stop for implementation review: show the changes, evidence, remaining concerns,
   and a proposed `feat(<scope>): ...` or `fix(<scope>): ...` message.

Plan approval does not bypass test review or authorize commits. Reuse explicit
approval within its stated scope. Documentation/configuration work that gains
nothing from a failing test may use one reviewed change; configuration that
changes behavior still needs appropriate tests.

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

Verify findings, fix confirmed issues within the approved scope, rerun affected
checks, and request another independent review of the updated diff. Repeat until
the final code state has no unresolved actionable findings. Any later code edit
requires renewed review. Record findings, dispositions, and checked state in
the feature's validation.md when present; do not paste repeated review logs into
user replies.

Review does not authorize implementation before test approval, new requirements,
or dependency changes. Record unresolved questions or decisions in the relevant
requirements/plan and ask the user with the available question tool before
dependent work. If delegation is unavailable or the loop cannot resolve a
finding, report the blocker and request a decision; do not substitute self-review
or claim a clean review. A clean review is evidence, not proof of defect-free code.
