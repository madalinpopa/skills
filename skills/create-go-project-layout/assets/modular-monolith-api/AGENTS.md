# Agent instructions

`{{PROJECT_NAME}}` is a Go web API in `api/` with a separate frontend in
`web/`. [docs/SPEC.md](docs/SPEC.md) is the behavior source of truth.

## Change workflow

Before any change, inspect current code and re-read the relevant spec. Check
current official docs for every external library touched; do not rely on
remembered APIs or old examples.

For each behavioral change:

1. Write only necessary tests and run them. **Stop for review:** show the
   test changes, output proving failure for the right reason, and a proposed
   `test(<scope>): ...` commit message.
2. Implement only after explicit developer approval. Write the smallest clear
   code that passes the tests and fits the spec.
3. Run the verification below. **Stop for review:** show changes, results, and
   a proposed `feat(<scope>): ...` or `fix(<scope>): ...` commit message.

Each behavioral change has two commits: failing tests, then implementation.
Config, docs, and release changes that gain nothing from failing tests may use
one commit; developer review before committing still applies.

Propose a commit message for every requested change; never create Git or
Jujutsu commits unless asked. Keep PRs focused on one piece of work; open them
only after full developer review and an explicit request.
On “ready to create PR,” review the latest commits between `main` and `HEAD`
and provide a PR title. Never add AI attribution (including Co-Authored-By
trailers or generated-by signatures) to commits or PRs; preserve genuine human
attribution.

## Feature and phase planning

For one or more planned phases/commits or an explicit feature-planning request,
use [FEATURE.md](docs/templates/FEATURE.md) and
[PHASE.md](docs/templates/PHASE.md).

- Reuse an existing feature or create `docs/features/<feature-name>/FEATURE.md`
  with scope, non-goals, acceptance criteria, and roadmap.
- Detail only the active phase in `phase-NN-<outcome>.md` beside it; outline
  future phases. One reviewable outcome, including its test/implementation
  commit pair, belongs to one phase.
- Present the plan for review before changes; honor existing explicit approval
  within its scope. Plan approval never bypasses test review or authorizes
  commits, PRs, or additional scope.
- Follow template context-loading and review checkpoints. Update the checkpoint
  and roadmap at review boundaries and session end with evidence, developer
  decisions, and the exact next action.
- On resume, read the feature overview and active phase; verify recorded state
  against the working tree.

## Modern Go

- Before creating or editing any Go file, invoke
  `$modern-go-guidelines:use-modern-go`. Run `list` for that exact file, or Go
  {{GO_VERSION}} for a new file, and read the full, unfiltered output.
- Use `explain` only for IDs needing detail. Re-run guidance when moving to a
  different kind of file or package. Follow it unless it would change required
  behavior or break the build; explain exceptions.
- Write explicit, idiomatic, readable code with descriptive names and direct
  control flow. Keep packages small and focused; never create them to fill
  slots. Complex abstractions need a concrete reason.
- Trust internal invariants. Validate external input, config, filesystem
  state, and database results; avoid speculative defensive layers.
- Let code explain itself. Add short doc comments only for useful information
  not obvious from code; avoid long comment blocks.

## Package boundaries

- `api/cmd/server`: flags, logging setup, process context, the database pool,
  the Echo instance, and `server.New`. Nothing else.
- `api/internal/shared`: types every module shares, such as the UUID value
  type. Imports nothing from the project.
- `api/internal/server`: config loading and the module lifecycle.
- `api/internal/modules`: the `Module` interface, the router surface, and the
  cross-module contracts registry.
- `api/internal/modules/<name>`: one feature module. `domain` holds entities,
  value objects, and repository interfaces; `app` holds commands, queries, and
  the ports they need; `adapters` implement those ports for HTTP, Postgres, and
  external services; `module.go` wires them; `tests` drive the module through
  its HTTP API.
- `api/internal/platform`: errors, Echo setup, logging, and Postgres helpers
  shared by every module. Knows no module.
- `api/internal/testkit`: the integration environment for module tests.

Dependency rules:

- `domain` imports only `shared` and `platform/errs`.
- `app` imports `domain`; adapters implement the ports `app` defines.
- `adapters/*` import `app` and `domain`, never the reverse.
- Modules never import each other; they use `modules.Contracts`.
- `platform/*` never imports a module.
- Each module owns one Postgres schema. No foreign keys cross schemas; another
  module's identifier is stored as a plain UUID value.

Generated files (`*.gen.go`, `models/`) come from `openapi.yaml` and
`queries/*.sql` through `task gen`. Never edit them by hand.

## Tests

- Test observable behavior and realistic regressions, not plain field
  assignments, library behavior, or private helpers without real logic.
  Never add tests just for coverage.
- Give each behavior one main test at the lowest useful layer; avoid overlap.
  Domain and app unit tests own invariants and decisions. HTTP adapter tests
  own parsing, status codes, headers, and middleware. Postgres integration
  tests own SQL, migrations, constraints, and transactions. Module tests own
  wiring and a few critical journeys through real boundaries.
- Make Arrange, Act, Assert clear with comments and spacing.
- Use [Testify](https://github.com/stretchr/testify) consistently: `require`
  for setup that later checks depend on, `assert` for independent checks.
- Use small `map[string]struct` tables for shared setup and assertions.
  Split large or uneven cases into `_success` and `_failed` tests.
- Use `t.Parallel()` in tests and subtests except for process-state changes or
  shared external resources; briefly explain exceptions. Pass paths,
  environment values, clocks, and writers as dependencies.
- Use [Testcontainers for Go](https://golang.testcontainers.org/) for Postgres
  integration tests. At the very start of each container test, check
  `INTEGRATION=true`; otherwise skip with a clear message.
- Default tests must need neither network nor Docker.

## Verification

Run every command from the repository root. After each implementation, in order:

1. Run changed-package tests from `api/`.
2. Run `go fix ./...` in `api/` and review its diff.
3. Run `task format`; ensure no Go file remains unformatted.
4. Run `task test:unit`.
5. Run `go vet ./...` in `api/`.
6. Run `task lint`.
7. Run `task test:integration` only when touching persistence, module wiring,
   or another integration boundary.
8. Review the final diff; preserve unrelated user changes.

Report every failed command and whether the cause is the change, environment,
or an unavailable integration dependency.
