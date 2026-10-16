# Code conventions

Shared rules; language conventions extend them. Tests follow [TDD.md](TDD.md).

## Design and readability

- Put behavior with the module owning its rules/state. Separate business rules
  from transport, presentation, and persistence; respect existing boundaries.
- Keep functions focused; split for clarity, not size limits. Pass behavioral
  dependencies explicitly rather than hiding them in globals or process state.
- Use meaningful names and direct control flow. Add abstractions for concrete
  shared responsibilities or substitution needs, not speculative extension.
- Share code when meaning and reasons for change align; similar syntax alone
  does not justify coupling unrelated behavior.

## State, errors, and resources

- Validate external data at entry: configuration, persisted data, and external
  responses. Check business constraints, not just parsing. Trust established
  invariants internally; avoid repeated checks or defaults that hide broken assumptions.
- Give mutable state one owner/source of truth; derive values instead of
  maintaining synchronized copies where practical.
- Handle expected failures explicitly. Add operation context and preserve causes
  callers need. Recover only through a required, useful recovery path; never
  discard errors silently or report failure as success.
- Handle/report errors where action is possible; avoid duplicate logging as they
  propagate. Make resource ownership, transfers, and cleanup clear on success,
  failure, and cancellation.

## Security and application boundaries

- Identify affected trust boundaries and enforce authorization for the action
  and resource on the server; client-side checks only support the user experience.
- Use parameterized queries and context-appropriate output encoding. Use
  established authentication and cryptographic libraries rather than custom schemes.
- Keep secrets and sensitive data out of source, logs, and error responses.
  Return useful public errors without exposing internal details.
- Bound untrusted input, work, and resource consumption. Consider relevant abuse
  cases, including cross-user access and resource exhaustion, when changing a boundary.
- Trace affected contracts across UI, API, and persistence. Preserve data
  integrity and account for transaction boundaries and compatibility.

## Comments and formatting

- Explain non-obvious constraints, reasons, and tradeoffs; let code explain mechanics.
  Update comments with behavior. Follow project formatters/linters; avoid unrelated rewrites.

## External dependencies

- Before using a library or proposing dependency changes, verify its latest release
  through official releases/registries and current official docs. Compare existing
  pins and supported toolchains; check official docs for external APIs being changed.
- Before adding, installing, replacing, or upgrading a dependency, present its name,
  proposed version, purpose, and compatibility implications; obtain explicit user
  approval. Reuse approval only within its scope; checking versions authorizes no upgrade.
- Report unavailable version verification before requesting approval. Record approved
  choices and verified versions in [STACK.md](../specs/STACK.md).

## TypeScript and frontend

Apply only when used. Follow frameworks/commands in [STACK.md](../specs/STACK.md);
none is assumed.

- Keep public boundary types explicit; narrow unknown data rather than using
  unchecked assertions or `any`.
- Add shared state only for demonstrated cross-component needs.
- Represent loading, empty, success, and failure states where applicable.
- Use semantic controls, labels, and keyboard interactions.

## Project-specific boundaries and conventions
