# Go conventions

Apply only to Go work, alongside [CODE.md](CODE.md) and [TDD.md](TDD.md).
Read the relevant `go.mod` and toolchain configuration; use supported versions,
not the machine's newest Go. Record versions and module roots in
[STACK.md](../specs/STACK.md).

## Readability and APIs

- Name packages by capability; avoid catch-all utility packages and package-name
  repetition. Keep exports small, names proportional to scope, and initialisms
  consistent (`ID`, `URL`). Handle failures early; avoid confusing shadowing.
- Define small interfaces at the consuming boundary when substitution is needed.
  Prefer concrete returns unless an interface is the API; avoid mirrored
  interfaces and interfaces added solely for mocks.
- Prefer suitable standard-library APIs. Use generics for shared typed algorithms
  or containers; avoid `any`, reflection, and assertions that discard known types.
- Document exports and non-obvious contracts: ownership, mutation, zero values,
  partial results, side effects, and concurrent use. Keep other comments selective.

## Values and resources

- Prefer useful zero values. Use constructors to establish invariants or acquire
  dependencies; reject invalid external configuration before returning a usable value.
- Use pointer receivers for mutation or non-copyable state; keep receivers
  consistent. Never copy used synchronization values or structs containing them.
- Slice and map copies share data; clone when independent mutation is required.
  `slices.Clone` and `maps.Clone` are shallow: nested mutable values still alias.
- Preserve observable absent/zero/nil/empty distinctions, especially in serialized
  data and updates. Initialize maps before writes; sort keys for ordered output.
- Keep acquisition and cleanup together, normally with `defer` after success.
  Release per-iteration resources promptly rather than retaining them until return.
- Check write, flush, close, and commit errors that affect success. Preserve the
  primary failure and report meaningful cleanup failures; cleanup alone proves
  neither persistence nor success.

## Errors

- Return errors for expected failures; reserve panic for programming errors or
  documented unrecoverable conditions. Never recover merely to report success.
- Add useful operation context. Wrap with `%w` when callers should inspect the
  cause; wrapped errors become API commitments. Translate dependency errors when
  exposing them would violate the package boundary.
- Match sentinels with `errors.Is` and types with supported typed inspection APIs
  such as `errors.As`, never error strings. Add custom errors only for useful
  recovery decisions or structured information.
- Return explicit `nil` on success: a typed nil pointer in an `error` is non-nil.
  Do not hide failures with empty results, retries, or fallbacks unless required.

## Context and concurrency

- Pass `ctx context.Context` first where cancellation/deadlines apply and propagate
  it downstream. Never replace request context with `context.Background()` to
  bypass cancellation. Pass contexts per call unless an API requires storage;
  values carry request metadata, not options or dependencies.
- Always release derived contexts with their cancel function. Cancellation neither
  waits for goroutines nor interrupts work that ignores it.
- Add concurrency only for a concrete need. Bound work; define goroutine ownership,
  exit, error propagation, and completion. Join workers before releasing their resources.
- Use mutexes for shared state and channels for communication when clearest.
  Synchronize whole invariants; concurrent maps do not make compound actions atomic.
- Close channels only when all sends have finished, under a clear owner. Make
  blocking communication cancellable when work may be abandoned; buffers and
  sleeps do not replace coordination.

## Tests and verification

Coverage and review follow [TDD.md](TDD.md).

- Prefer small `map[string]struct` tables when suitable; use slices for ordered
  cases. Name groups `TestBehavior_success` and `TestBehavior_failed`.
- Use `t.Helper()`, `t.Cleanup()`, `t.TempDir()`, and local fixtures. Parallelize
  only isolated tests, never tests changing environment/working directory or
  shared mutable state. Coordinate with signals rather than sleeps.
- Run project format, test, build, and `go vet` checks for relevant modules/build
  configurations. Exercise changed concurrency with `go test -race`; it detects
  only races on executed paths. Report unavailable tools or unsupported environments.
- Consider fuzzing parsers, decoders, and transformations with useful invariants;
  retain regression seeds and time-bound exploratory runs.
- Use the project's `govulncheck` workflow; dependency installation follows
  [CODE.md](CODE.md). Benchmark/profile demonstrated performance problems before
  adding pools, caches, or synchronization complexity.

## AI-assisted changes and modernization

- Inspect affected packages, callers, tests, and versions before generating code.
  Reuse existing types/helpers; verify unfamiliar APIs for the selected version
  with `go doc` and official documentation.
- Review generated code against the contracts above. Compilation does not prove
  behavior. Derive test expectations independently from requirements or checked
  examples; never weaken assertions or requirements to fit generated code.
- Prefer tool diagnostics and deterministic refactoring over speculative rewrites.
  Inspect diffs and run checks; generated claims are not verification evidence.
- When available, follow `modern-go-guidelines:use-modern-go` for the file/module.
  Otherwise use official Go docs and report the limitation; do not auto-install tools.
- Periodically run `go fix ./...` per module, especially after supported toolchain
  upgrades and before substantial reviews. For focused work, target the affected
  package, including when changing one file. Consult `go help fix` for supported flags.
- Use supported idioms when clearer. Review fixes for nil preservation, aliasing,
  and evaluation order; preserve behavior, scope, and unrelated edits. Format and
  rerun affected tests. Do not raise the Go version merely to modernize syntax.

Sources: [Go review guidance](https://go.dev/wiki/CodeReviewComments),
[errors](https://go.dev/blog/go1.13-errors), [context](https://pkg.go.dev/context),
[modernization](https://go.dev/blog/gofix), and
[AI-assisted Go](https://developers.googleblog.com/why-go-is-an-ideal-language-for-ai-assisted-software-engineering/).
