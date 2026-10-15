# Code conventions

Shared conventions for code changes. Language-specific guidance builds on these
rules; tests and verification follow [TDD.md](TDD.md).

## Responsibilities and boundaries

- Place behavior with the module that owns its rules and state. Keep transport,
  presentation, and persistence details out of business rules; follow existing
  project boundaries rather than introducing a new architecture for a small change.
- Keep functions focused on a coherent operation. Split them when doing so makes
  responsibilities or control flow clearer, not to meet an arbitrary size limit.
- Pass dependencies explicitly when they affect behavior; avoid hidden inputs
  through globals or ambient process state.

## Simplicity and abstractions

- Use names that express purpose and domain meaning. Prefer direct control flow
  over nested conditions or indirection that obscures the operation.
- Introduce an abstraction for a concrete shared responsibility or a needed
  substitution boundary. Avoid speculative extension points and generic frameworks
  for a single use case.
- Share code when callers share meaning and change for the same reason. Similar
  syntax alone is not enough; a little duplication can be clearer than coupling
  unrelated behavior.

## Validation and state

- Validate external data at the boundary where it enters the program, including
  configuration, persisted data, and responses from external systems. Check the
  constraints the application needs, beyond whether the data can be parsed.
- Once a boundary establishes an invariant, rely on it internally. Do not add
  repeated checks or silent defaults that hide broken program assumptions.
- Give mutable state a clear owner and source of truth. Derive values where
  practical instead of maintaining copies that must be synchronized.

## Errors and resource ownership

- Handle expected failures explicitly. Add context that identifies the failed
  operation while preserving the underlying cause when callers need it.
- Recover only when the required behavior defines a useful recovery path. Do not
  silently discard errors or turn failure into apparent success.
- Handle or report an error at the layer that can act on it. Avoid logging the
  same failure at every layer as it propagates.
- Make resource ownership and cleanup visible. Cover failure and cancellation
  paths as well as success; document transfers of ownership when they are unclear.

## Comments and formatting

- Explain non-obvious constraints, tradeoffs, and reasons in comments. Let names
  and code express the mechanics; update comments when the behavior changes.
- Follow project formatters and linters. Avoid unrelated style rewrites.

## External dependencies

Always verify the latest released version of an external library before using
it or proposing a dependency change. Use its official release source or package
registry and current official documentation; compare that version with any
existing pin and the project's supported toolchain. Verify current official
documentation for external APIs being changed.

Before adding, installing, replacing, or upgrading an external dependency,
present its name, proposed version, purpose, and compatibility implications,
then obtain explicit user approval. Existing approval counts only within its
stated scope. Checking the latest version does not authorize an upgrade.
If version verification is unavailable, report the unresolved check before
asking the user to approve the proposal. Record approved choices and verified
versions in [STACK.md](../specs/STACK.md).

## TypeScript and frontend work

Apply only when the project uses TypeScript or a frontend. Follow the selected
framework and commands in [STACK.md](../specs/STACK.md); no framework is assumed.

- Keep types explicit at public boundaries; prefer narrowing unknown external
  data over unchecked assertions or `any`.
- Add shared frontend state only for demonstrated cross-component needs.
- Represent loading, empty, success, and failure states where applicable.
- Use semantic controls, labels, and keyboard interactions.

## Project-specific boundaries and conventions
