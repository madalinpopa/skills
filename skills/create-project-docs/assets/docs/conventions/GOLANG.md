# Go conventions

Apply only to Go work. Read the relevant `go.mod` and toolchain configuration;
supported versions and module locations belong in [STACK.md](../specs/STACK.md).
Do not assume the machine's newest Go version is the project's target.
Apply [shared code conventions](CODE.md) and [TDD.md](TDD.md) alongside the
Go-specific rules below.

- Define interfaces at the consuming boundary when substitution is needed;
  do not wrap concrete types just to introduce an interface.
- Follow local error wrapping and handling rules.
- Pass request context through blocking operations.
- Add concurrency for a concrete need, with bounded work and clear goroutine
  lifetimes. Prefer simpler synchronous code when it meets the requirement.
- For table tests, prefer simple `map[string]struct` tables when appropriate;
  name success and failure groups `TestBehavior_success` and `TestBehavior_failed`.
- Mark Go test helpers with `t.Helper()` so failures point to the calling test.
- Use temporary directories and local fixtures for filesystem tests.
- Use parallel tests only when they do not mutate shared process state or resources.
- Include `go vet` in verification. Introduce language features and library APIs
  only within the supported versions.

When `modern-go-guidelines:use-modern-go` is available, follow its installed
version-specific guidance for the file/module being changed. If unavailable,
use current official Go documentation and report that limitation; this document
does not require installing tools automatically.

## Keep Go syntax current

Periodically run `go fix ./...` from each relevant module root, especially after
a supported Go version upgrade and before reviewing substantial Go changes.
For focused work, run `go fix ./path/to/package`; for a specific file, target
its enclosing package so fixes retain package and type context.

Use the project's supported toolchain and review the generated diff. Keep changes
within the agreed scope, preserve unrelated edits, then run formatting and the
affected tests. Do not raise the module's Go version just to apply newer syntax.
Available fixes depend on the toolchain; consult [official go fix documentation](https://pkg.go.dev/cmd/go#hdr-Apply_fixes_suggested_by_static_checkers)
and `go help fix` when selecting flags or a narrower scope.
