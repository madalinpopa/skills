# Code conventions

Keep changes within the agreed outcome and existing module responsibilities.
Verify current official documentation for external APIs being changed.

- Use descriptive names, direct control flow, and cohesive functions/modules.
- Add abstractions and dependencies for a demonstrated need, with clear ownership.
- Validate input at external boundaries; trust established internal invariants.
- Handle errors explicitly with useful context. Keep resource cleanup visible.
- Explain non-obvious constraints in comments; avoid narrating the code.
- Follow project formatters and linters. Avoid unrelated style rewrites.
- Protect observable behavior with focused tests, not coverage targets.

## External dependencies

Always verify the latest released version of an external library before using
it or proposing a dependency change. Use its official release source or package
registry and current official documentation; compare that version with any
existing pin and the project's supported toolchain.

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
- Keep server/domain logic separate from rendering and browser-specific code.
- Keep state ownership clear; derive values when possible instead of maintaining
  duplicate state. Add shared state only for demonstrated cross-component needs.
- Represent loading, empty, success, and failure states where applicable.
- Use semantic controls, labels, and keyboard interactions; test behavior through
  user-visible outcomes rather than component internals.

Pending: record project-specific boundaries and conventions when known.
