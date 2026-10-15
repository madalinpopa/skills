# Go conventions

Apply only to Go work. Read the relevant `go.mod` and toolchain configuration;
supported versions and module locations belong in [STACK.md](../specs/STACK.md).
Do not assume the machine's newest Go version is the project's target.
Apply [shared code conventions](CODE.md) and [TDD.md](TDD.md) alongside the
Go-specific rules below.

## Readable packages and APIs

- Name packages for the capability they provide; avoid catch-all `util`, `common`,
  or `helpers` packages. Keep the exported surface small and names clear at the
  call site, without repeating the package name.
- Use short names for small, obvious scopes and descriptive names for longer
  scopes or domain concepts. Keep initialisms consistent, such as `ID` and `URL`.
- Handle errors and exceptional cases early so the main operation reads straight
  through. Avoid shadowing an error or result when it obscures which value returns.
- Define small interfaces in the consuming package when a caller needs substitution.
  Prefer concrete return types unless an interface is the intended API; do not
  mirror every struct with an interface or add one solely to generate mocks.
- Prefer standard-library facilities when they fit. Use generics for genuinely
  shared algorithms or containers with useful type constraints; avoid `any`,
  reflection, or type assertions that discard information known at compile time.
- Document exported declarations and non-obvious contracts: mutation, ownership,
  zero-value behavior, and concurrent use. Keep implementation commentary selective.

See [Go code review comments](https://go.dev/wiki/CodeReviewComments) and
[Google's Go style decisions](https://google.github.io/styleguide/go/decisions).

## Values, mutation, and ownership

- Make zero values useful where practical. Use constructors when initialization
  must establish invariants or acquire dependencies; return errors for invalid
  external configuration rather than leaving a partially usable value.
- Choose pointer receivers for mutation or values that must not be copied. Keep
  receiver choices consistent for a type; do not copy a used mutex, `sync.Once`,
  `sync.WaitGroup`, or a struct containing them.
- Remember that copying a slice or map still shares its backing data. Copy at an
  ownership boundary when callers must mutate independently; `slices.Clone` and
  `maps.Clone` are shallow copies, so nested mutable values may still alias.
- Preserve observable distinctions between absent, zero, nil, and empty values,
  especially in serialization and updates. Initialize maps before writing to them;
  sort keys when output order is part of the contract.
- Keep acquisition and cleanup together. After successful acquisition, normally
  defer cleanup in the owning scope. In loops, use a smaller scope or explicit
  cleanup when waiting for the outer function to return would retain resources.
- Check write, flush, close, and commit errors when they determine whether an
  operation succeeded. Preserve the primary failure while reporting meaningful
  cleanup failures; deferred cleanup alone does not prove a write succeeded.

See the [`sync` contracts](https://pkg.go.dev/sync),
[`slices.Clone`](https://pkg.go.dev/slices#Clone), and
[Effective Go](https://go.dev/doc/effective_go).

## Errors as contracts

- Return `error` for expected failures; reserve panic for programming errors or
  documented unrecoverable conditions. Do not recover just to return success.
- Add operation-specific context when useful. Use `%w` when callers should inspect
  the cause, recognizing that wrapped errors become part of the API. Translate
  dependency errors at boundaries when exposing them would leak implementation
  details the package should remain free to change.
- Use `errors.Is` for sentinel matching and the supported toolchain's typed error
  inspection API, such as `errors.As`, for error types. Do not branch on error text.
  Introduce custom errors only when callers need a distinct recovery decision or
  structured information.
- Return an explicit `nil` error on success. A typed nil pointer stored in an
  `error` interface is non-nil; avoid returning one accidentally.
- Make partial results and partial side effects explicit in the contract. Do not
  hide failures behind empty results, retries, or fallback values unless required.

See [working with wrapped errors](https://go.dev/blog/go1.13-errors) and
[nil errors in the Go FAQ](https://go.dev/doc/faq#nil_error).

## Context and concurrency

- Accept `ctx context.Context` as the first parameter of operations that need
  cancellation or deadlines, and pass it to downstream work. Do not replace a
  request context with `context.Background()` to bypass cancellation.
- Pass contexts per call rather than storing them in long-lived structs, except
  when an existing API contract requires it. Use context values for request-scoped
  metadata, not optional arguments or dependency injection.
- Call the cancel function from `WithCancel`, `WithTimeout`, or `WithDeadline` on
  every path once the derived work is done. Cancellation signals a stop; it does
  not wait for goroutines to exit or interrupt code that ignores the context.
- Add concurrency for a concrete need. Bound in-flight work, and make each
  goroutine's owner, exit condition, error propagation, and completion observable.
  Join owned workers before releasing resources they still use.
- Use a mutex for shared state when it is clearest, and channels for communication
  or ownership transfer. Protect the entire invariant, not just individual fields;
  do not assume a concurrent map makes multi-step operations atomic.
- Give channel closure a clear owner that knows all sends are finished. Make
  potentially blocking sends and receives cancellable when work may be abandoned.
  Buffers and sleeps are not substitutes for coordination.

See the [`context` contract](https://pkg.go.dev/context) and
[goroutine lifetime guidance](https://go.dev/wiki/CodeReviewComments#goroutine-lifetimes).

## Go tests and verification

Follow [TDD.md](TDD.md) for coverage choices, review, and the test-first workflow.

- For table tests, prefer simple `map[string]struct` tables when appropriate;
  name success and failure groups `TestBehavior_success` and `TestBehavior_failed`.
  Use a slice when case order matters; never rely on map iteration order.
- Mark helpers with `t.Helper()`, register test-owned cleanup with `t.Cleanup()`,
  and use `t.TempDir()` and local fixtures for filesystem tests.
- Use `t.Parallel()` only for isolated tests. Do not combine it with `t.Setenv`,
  working-directory changes, or shared mutable resources. Coordinate concurrent
  tests with explicit signals instead of timing assumptions and arbitrary sleeps.
- Run the project's formatting, test, build, and `go vet` checks for relevant
  modules and build configurations. Exercise changed concurrent paths with
  `go test -race`; a passing run covers only paths executed under that workload.
- Consider fuzz tests for parsers, decoders, or transformations with useful
  invariants. Keep discovered regressions as seed cases; give exploratory fuzz
  runs a time budget rather than making ordinary test runs unbounded.
- Use the project's `govulncheck` workflow for vulnerability checks. Report missing
  tools or unsupported race-detector environments; follow [CODE.md](CODE.md) before
  installing dependencies. Benchmark and profile a demonstrated performance
  problem before adding pooling, caching, or synchronization complexity.

See [race detection](https://go.dev/doc/articles/race_detector),
[Go fuzzing](https://go.dev/doc/security/fuzz/), and
[Go security practices](https://go.dev/doc/security/best-practices).

## AI-assisted Go changes

- Before generating code, inspect the affected package, callers, existing tests,
  and module versions. Reuse established types and helpers; check unfamiliar APIs
  against the selected dependency version with `go doc` and official documentation.
- Review generated code for the same contracts as handwritten code: error causes,
  nil behavior, aliasing, cleanup, cancellation, and partial failure. Compilation
  establishes type correctness, not the required behavior.
- Derive expected results from requirements or independently checked examples.
  Tests that reproduce the generated algorithm can repeat its mistakes. Do not
  weaken assertions or change requirements merely to make generated code pass.
- Prefer toolchain diagnostics and deterministic refactoring tools over repeated
  speculative rewrites. Inspect each diff and run the relevant checks; generated
  explanations and claimed test results are not verification evidence.

These conventions apply
[Google's guidance on reviewing and verifying AI-assisted Go](https://developers.googleblog.com/why-go-is-an-ideal-language-for-ai-assisted-software-engineering/)
to this workflow; they do not require a particular AI tool.

## Keep Go syntax current

When `modern-go-guidelines:use-modern-go` is available, follow its installed
version-specific guidance for the file/module being changed. If unavailable,
use current official Go documentation and report that limitation; this document
does not require installing tools automatically.

Periodically run `go fix ./...` from each relevant module root, especially after
a supported Go version upgrade and before reviewing substantial Go changes.
For focused work, run `go fix ./path/to/package`; for a specific file, target
its enclosing package so fixes retain package and type context.

Use the project's supported toolchain and review the generated diff. Keep changes
within the agreed scope, preserve unrelated edits, then run formatting and the
affected tests. Do not raise the module's Go version just to apply newer syntax.
Available fixes depend on the toolchain; consult [official go fix documentation](https://pkg.go.dev/cmd/go#hdr-Apply_fixes_suggested_by_static_checkers)
and `go help fix` when selecting flags or a narrower scope.

Go 1.26 introduced a rewritten `go fix` with version-aware modernizers. Prefer
supported idioms when they improve clarity, but review semantic details such as
nil preservation, aliasing, and evaluation order. A modernization is not permission
to change behavior. See [using go fix to modernize Go code](https://go.dev/blog/gofix).
