# Agent instructions

`skills` is a Go CLI. [docs/DESIGN.md](docs/DESIGN.md) is the source of truth
for how it works. Read the relevant part of it before you change any behavior.

No AI LLM tool may append AI attribution (such as Co-Authored-By trailers or generated-by signatures)
to commits or pull requests. Genuine human attribution must be preserved.

Each change that adds or changes behavior has two commits:

1. `test(<scope>): ...` adds the smallest useful failing tests.
2. `feat(<scope>): ...` or `fix(<scope>): ...` makes those tests pass.

For each change:

1. Look at the current code and re-read the relevant design section.
2. Check the current official docs for every external library you touch. Do not
   rely on remembered APIs or old examples.
3. Write only the tests this change needs. Run them and show that they fail for
   the right reason.
4. Stop. Show the user the test changes, the failure output and a proposed
   commit message.
5. Do not write the implementation until the user says to go ahead.
6. Write the smallest clear code that passes the tests and fits the design.
7. Run the verification steps below. Stop again and show the user the changes,
   the results and a proposed commit message.
8. Do not create a Git or Jujutsu commit unless the user asks for it. Proposing
   a commit message is required. Creating the commit is not.

Config, docs and release work that gains nothing from a failing test can be a
single commit. It still needs user review before you commit.

Keep each pull request focused on one piece of work. Do not open a pull request
until the user has reviewed all of it and asks for one.

## Modern Go

- Invoke `$modern-go-guidelines:use-modern-go` before you create or edit any Go
  file.
- Run its `list` operation for the exact file you are changing, or for Go 1.27
  if the file does not exist yet. Read the full, unfiltered output.
- Use `explain` only for guideline IDs that need more detail.
- Run the guidance again when you move to a different kind of file or package.
- Follow the guidance unless it would change required behavior or break the
  build. Explain any exception.
- Run `go fix ./...` after each implementation and review its diff.

Write explicit, idiomatic code that is easy to read. Use small packages with
clear jobs, descriptive names and direct control flow. Any complex abstraction
needs a concrete reason. Trust internal invariants. Validate external input,
config, filesystem state and Git results, but do not add extra defensive
layers "just in case".

## Package boundaries

- `main.go`: version injection, process context and running the command only.
- `cmd`: Cobra commands, argument checks, output rendering, writers and exit
  codes.
- `internal/config`: config types, defaults, Viper loading and first-run save.
- `internal/project`: finding the repository root, with a fallback outside a
  repository.
- `internal/store`: Git clone, fast-forward sync and commit lookup.
- `internal/skill`: catalog discovery, YAML metadata checks and per-agent
  transforms.
- `internal/install`: target trees, three-way planning, locks, filesystem
  changes, backups, update and removal.
- `internal/gittest`: helpers that build temporary Git repositories for tests.

Keep Cobra inside `cmd` and Viper inside `internal/config`. Domain packages must
not depend on either one, and must not print. Pass Git arguments straight to the
process runner and never build a shell command string. Do not create a package
only to fill a slot.

## External libraries

- [Cobra](https://cobra.dev/) for commands, arguments, flags and help.
- [Viper](https://github.com/spf13/viper) for TOML config and defaults.
- [Testify](https://github.com/stretchr/testify) for test assertions.
- [Testcontainers for Go](https://golang.testcontainers.org/) only when a real
  container adds coverage that temporary directories and local Git repositories
  cannot.

## Tests

- Test first: reviewed failing test, then the implementation.
- Test behavior users can see and real failures that could come back. Do not
  test plain field assignment, library behavior or private helpers with no real
  logic.
- Avoid overlap. Each behavior gets one main test at the lowest useful layer.
  Integration tests cover only what unit tests cannot prove.
- Follow Arrange, Act, Assert, with spacing that makes the three parts easy to
  see.
- Use Testify consistently: `require` for setup that later checks depend on,
  `assert` for independent checks.
- Use `map[string]struct` table tests when cases share the same setup and
  assertions. Keep tables small. Split large or uneven cases into `_success` and
  `_failed` tests instead.
- Call `t.Parallel()` in tests and subtests unless they change process state or
  share external resources. Pass paths, environment values, clocks and writers
  in as dependencies so tests can stay parallel.
- Do not force parallel tests around `t.Setenv`, working-directory changes,
  shared repositories or shared containers. Briefly explain the exception.
- Use temporary directories and local Git repositories for filesystem and Git
  tests. They need no network and give better feedback than mocks.
- A container test must use Testcontainers for Go with an Ubuntu image that has
  Git. At the very start, check `INTEGRATION=true` and skip with a clear message
  if it is not set. Do not repeat unit coverage.
- The default test run must need neither network nor Docker.
- Do not add tests just to raise coverage.

## Verification

After each implementation:

1. Run the tests for the changed package.
2. Run `go fix ./...` and review what it changed.
3. Run `task format` and make sure no Go file is left unformatted.
4. Run `go test ./...`.
5. Run `go vet ./...`.
6. Run `task test:integration` only when the change touches a real integration
   boundary.
7. Review the final diff and make sure unrelated user changes are untouched.

Never hide a failed command. Say whether the failure comes from the change, the
environment or an integration dependency that is not available.
