# Tests and verification

## Behavioral changes

1. Write only necessary tests and run them. **Stop for review:** show the
   test changes, output proving failure for the right reason, and a proposed
   `test(<scope>): ...` commit message.
2. Implement only after explicit developer approval. Write the smallest clear
   code that passes the tests and fits the spec.
3. Run the verification below. **Stop for review:** show changes, results, and
   a proposed `feat(<scope>): ...` or `fix(<scope>): ...` commit message.

Config, docs, and release changes that gain nothing from failing tests may
omit the test gate. Implementation review still applies. Follow
[COMMITS.md](COMMITS.md) for commit and PR permissions.

## Test selection and style

- Test observable behavior and realistic regressions, not plain field
  assignments, library behavior, or private helpers without real logic.
  Never add tests just for coverage.
- Give each behavior one main test at the lowest useful layer; avoid overlap.
  Integration tests cover only what unit tests cannot prove.
- Make Arrange, Act, Assert clear with comments and spacing.
- Use [Testify](https://github.com/stretchr/testify) consistently: `require`
  for setup that later checks depend on, `assert` for independent checks.
- Use small `map[string]struct` tables for shared setup and assertions.
  Split large or uneven cases into `_success` and `_failed` tests.
- Use `t.Parallel()` in tests and subtests except for process-state changes or
  shared external resources; briefly explain exceptions. Do not force it with
  `t.Setenv`, working-directory changes, shared repositories, or shared
  containers. Pass paths, environment values, clocks, and writers as dependencies.
- Use temporary directories and local Git repositories for filesystem/Git
  tests: network-free and more informative than mocks.
- Use [Testcontainers for Go](https://golang.testcontainers.org/) only when a
  container adds coverage beyond temporary directories/local Git repos. Use
  an Ubuntu image with Git. At the very start of each container test, check
  `INTEGRATION=true`; otherwise skip with a clear message.
- Default tests must need neither network nor Docker.

## Go verification

After each Go implementation, run these commands from the repository root,
in order:

1. Run changed-package tests.
2. Run `go fix ./...` and review its diff.
3. Run `task format`; ensure no Go file remains unformatted.
4. Run `go test ./...`.
5. Run `go vet ./...`.
6. Run `task test:integration` only when touching a real integration boundary.
7. Review the final diff; preserve unrelated user changes.

## Other changes

- Documentation only: check Markdown links, routes, consistency, and the final
  diff. Unrelated Go checks are unnecessary.
- Repository skills: follow [SKILLS.md](SKILLS.md#verification).
- Configuration and release changes: check the affected configuration or
  workflow, and run Go verification if Go implementation is also affected.

[Taskfile.yml](../../Taskfile.yml) owns task definitions. Report every failed
command and whether the cause is the change, environment, or an unavailable
integration dependency. Report skipped checks and their reasons.
