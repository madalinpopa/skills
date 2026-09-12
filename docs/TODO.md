# Implementation plan

This plan turns [DESIGN.md](DESIGN.md) into a small, correct and pleasant CLI.
It is split into phases so that each phase is one focused pull request with one
to five readable commits.

Read this file before starting any phase. Then read only `DESIGN.md` and the
file for the active phase. The design is authoritative; the phase files describe
delivery order, not a different architecture.

## Instructions for the implementing agent

### Progress tracking

- Use the checkboxes in the active phase as the source of truth.
- Mark a test checkpoint complete only after the user reviews the tests.
- Mark an implementation checkpoint complete only after the user reviews the
  implementation and its verification results.
- Mark the phase complete in this file only when every task in the phase is
  complete and the user accepts the phase.
- Keep unfinished and deferred work unchecked. Never mark work complete merely
  because code was written.

### Mandatory review and commit workflow

Every behavior-bearing slice has two logical commits:

1. `test(<scope>): ...` defines the behavior with the smallest useful failing
   tests.
2. `feat(<scope>): ...` or `fix(<scope>): ...` implements that behavior and
   makes the tests pass.

For every slice:

1. Inspect the current repository and re-read the relevant design section.
2. Check current official documentation for every affected external library
   before writing code. Do not rely on remembered APIs or old examples.
3. Write only the tests needed for this slice and run them to prove they fail
   for the intended reason.
4. Stop. Show the test changes, the failure evidence and the proposed test
   commit message to the user.
5. Do not implement until the user explicitly says to proceed.
6. After approval, implement the smallest clear solution that satisfies the
   reviewed tests and design.
7. Run the required verification, stop again, and show the implementation
   changes, results and proposed implementation commit message to the user.
8. Do not create a Git or Jujutsu commit unless the user explicitly asks after
   reviewing the relevant changes. Supplying a commit message is required;
   creating the commit is not authorized by this plan.

Configuration, documentation and release work that gains no value from a
failing test may use one logical commit. It still requires user review before
any commit is created.

One phase equals one pull request. Do not mix work from a later phase into the
current PR. Do not open a PR until the user has reviewed the complete phase and
explicitly asks for it.

### Modern Go requirements

- Invoke `$modern-go-guidelines:use-modern-go` before creating or editing any Go
  file.
- Run its `list` operation for the exact file being changed, or for Go 1.27
  before a new file exists. Read the complete, unfiltered output.
- Use `explain` only for specific returned guideline IDs that need clarification.
- Re-run the guidance when moving to a different kind of Go file or package.
- Treat applicable guidance as authoritative unless it would change required
  behavior or fail to compile; explain any deliberate exception.
- Run `go fix ./...` after every implementation checkpoint and review its diff.
  Also run formatting, focused tests and the phase verification commands.

The code should be explicit, idiomatic and easy to read. Use small packages with
clear responsibilities, descriptive names and direct control flow. Any complex
abstraction needs a concrete justification. Trust internal invariants; validate
external input, configuration, filesystem state and Git results without adding
speculative defensive layers.

### External libraries

Use these libraries where the design calls for them:

- [Cobra](https://cobra.dev/) for commands, arguments, flags and help.
- [Viper](https://github.com/spf13/viper) for TOML configuration and defaults.
- [Testify](https://github.com/stretchr/testify) for test assertions.
- A maintained YAML library selected in Phase 2 after checking its current
  upstream documentation and Go 1.27 compatibility.
- [Testcontainers for Go](https://golang.testcontainers.org/) only when a real
  container boundary adds coverage that temporary directories and local Git
  repositories cannot provide.

Before writing code that uses Cobra or Viper, inspect their latest official
documentation and current recommended patterns. Keep Cobra inside the CLI
boundary and Viper inside the configuration boundary; domain packages must not
depend on either library.

### Test policy

- Follow test-driven development for behavior: reviewed failing test first,
  implementation second.
- Test observable functionality and realistic failure scenarios that could
  regress. Do not test trivial field assignment, library behavior or private
  helpers without meaningful logic.
- Avoid overlapping coverage. A behavior should have one primary test at the
  lowest useful layer; integration tests cover only boundaries that unit tests
  cannot prove.
- Follow Arrange, Act, Assert. Make the three sections visually easy to scan,
  using spacing or short comments where they improve readability.
- Use Testify consistently: `require` for prerequisites that make later checks
  unsafe, and `assert` for independent behavior checks.
- Use `map[string]struct` table tests when several independent cases share the
  same setup and assertion shape. Keep tables small and obvious. Split large or
  asymmetric cases into `_success` and `_failed` tests instead of building a
  complicated table.
- Call `t.Parallel()` for tests and subtests whenever they do not mutate process
  state or share external resources. Prefer passing paths, environment values,
  clocks and writers as dependencies so tests can remain parallel.
- Do not force parallelism around `t.Setenv`, working-directory changes, shared
  repositories, shared containers or other process-wide state. Explain the
  exception briefly.
- Use temporary directories and local temporary Git repositories for filesystem
  and Git boundary tests. They provide better feedback than mocks for these
  behaviors and need no network.
- If a container integration test becomes necessary, use Testcontainers for Go
  with an Ubuntu image and Git installed. At the very beginning of each such
  test, check for `INTEGRATION=true`; skip with a clear message when it is not
  set. Keep container tests focused and do not duplicate unit coverage.
- Do not add tests merely to increase coverage. Prefer a short suite that makes
  the contract clear and catches plausible bugs.

### Required verification

At each implementation checkpoint:

1. Run the focused package tests.
2. Run `go fix ./...` and inspect any changes it makes.
3. Run formatting and verify that no Go file remains unformatted.
4. Run `go test ./...`.
5. Run `go vet ./...`.
6. Run integration tests only when the active phase contains a justified
   integration boundary and `INTEGRATION=true` is available.
7. Inspect the final diff and confirm unrelated user changes were preserved.

Do not hide a failed command. Report whether the failure belongs to the change,
the environment or an intentionally unavailable integration dependency.

## Intended package boundaries

These boundaries keep dependencies pointing inward without creating a framework:

- `main.go`: version injection, process context and command execution only.
- `cmd`: Cobra construction, argument validation, writers and exit-code mapping.
- `internal/config`: configuration types, defaults, Viper loading and first-run
  persistence.
- `internal/project`: repository-root detection and non-repository fallback.
- `internal/store`: Git clone, fast-forward sync and commit discovery.
- `internal/skill`: catalog discovery, YAML metadata validation and per-agent
  transformation.
- `internal/install`: desired trees, three-way planning, locks, filesystem
  changes, backups, update and removal.

Do not create a package only to match this list. Keep formatting near `cmd`
until it becomes large enough to justify a separate package. Remove the empty
`skills` package when the first real structure replaces the scaffold.

## Phases

- [x] [Phase 1: CLI foundation, configuration and project scope](TODO_1.md) —
  four commits, one PR
- [ ] [Phase 2: Git store and portable skill metadata](TODO_2.md) — four
  commits, one PR
- [ ] [Phase 3: Catalog listing and change planning](TODO_3.md) — four commits,
  one PR
- [ ] [Phase 4: Safe installation and lock files](TODO_4.md) — four commits,
  one PR
- [ ] [Phase 5: Installed-skill lifecycle](TODO_5.md) — four commits, one PR
- [ ] [Phase 6: Diff, dry-run, polished output and release](TODO_6.md) — five
  commits, one PR

## Whole-project completion criteria

- Every command and safety guarantee in `DESIGN.md` is implemented.
- Project and global scopes remain explicit and never silently fall back.
- Local edits cannot be overwritten or removed without `--force` and a backup.
- Store synchronization is branch-based and fast-forward only.
- Human output, color behavior, dry-run behavior and exit codes match the
  documented contract.
- The focused test suite is readable, non-overlapping and green.
- The default test suite requires neither network nor Docker.
- CI and release workflows pass on the supported platforms.
- README usage matches the shipped CLI.
- No agent-created commit or PR exists without explicit user authorization.
