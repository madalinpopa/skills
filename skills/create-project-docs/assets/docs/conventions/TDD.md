# Test-driven changes

For production work, test every visible behavior change, including fixes.
Assert what users/callers rely on, not field assignments, library behavior,
private structure, or coverage.

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

Follow the overall slice and review workflow in
[AGENTS.md](../../AGENTS.md#workflow).

1. Before implementation, write/run the smallest useful failing test. Confirm
   missing behavior causes failure, not setup or environment; retain the command
   and relevant output.
2. Implement the smallest clear change that passes the tests; refactor only
   within approved scope while preserving agreed behavior.
3. Run focused tests, then applicable format, lint, full-test, and build checks
   from [STACK.md](../specs/STACK.md). Run integration checks when the changed
   boundary needs them. Record tested state and relevant failure evidence.

Skip failing tests for documentation/configuration that gains nothing from them;
behavior-changing configuration still needs tests. Report failed/skipped checks
and their causes; never claim an unrun check passed.
