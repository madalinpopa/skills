# Test-driven changes

Use for observable behavior changes, including fixes. Test the contract at
the lowest useful layer; give each behavior one main test. Prefer realistic
inputs and outcomes over private implementation details or coverage targets.

1. **Red:** write the smallest useful test and run it. Confirm the failure
   demonstrates missing behavior, not an unrelated build or environment error.
2. **Test review:** show the test change, command, failure evidence, and proposed
   `test(<scope>): ...` message. Wait for explicit approval before implementation.
3. **Green:** write the smallest clear implementation within the approved scope.
4. **Refactor:** clean up only where useful and keep the tests passing.
5. **Implementation review:** run the project's required checks; show results,
   the diff, and a proposed `feat(<scope>): ...` or `fix(<scope>): ...` message.

Plan approval does not bypass test review. Reuse explicit approval within its
scope. Keep test and implementation changes separately reviewable for two
commits, following the project's commit conventions.

Docs or configuration changes that gain nothing from a failing test may use
one reviewed change. Classify by effect: configuration can change behavior.
Keep default tests deterministic and independent of network or shared services.
Add integration coverage only for a boundary lower-level tests cannot prove.
Report failed or skipped checks and their cause; do not call unrun checks passed.
