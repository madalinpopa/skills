# Phase 1: CLI foundation, configuration and project scope

One phase equals one pull request. Read [TODO.md](TODO.md) and
[DESIGN.md](DESIGN.md) before starting. The review pauses, Modern Go workflow,
test policy and no-commit rule in `TODO.md` are mandatory.

## Outcome

The repository has a thin executable, a testable Cobra command tree, TOML
configuration through Viper, and deterministic project-scope resolution. Help
works without writing anything. Commands that need configuration can initialise
it, while callers can explicitly request a read-only path for dry-run work.

## Dependencies and boundaries

- Add Cobra, Viper and Testify only when their first use is introduced.
- Check the latest Cobra and Viper official documentation immediately before
  their test or implementation checkpoint.
- Keep Cobra in `cmd` and Viper in `internal/config`.
- Keep `main.go` free of command behavior and configuration logic.
- Put repository detection in `internal/project`; do not use process-wide
  working-directory changes in tests.

## Planned commits

1. `test(cli): define root command behavior`
2. `feat(cli): add the Cobra application foundation`
3. `test(config): define configuration and project scope`
4. `feat(config): load defaults and resolve project scope`

These are proposed logical commits. The agent must stop for review and receive
explicit permission before creating each one.

## Slice 1: Root command

- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.

### Test checkpoint

Add focused command-construction tests that define:

- running `skills` with no arguments prints useful help and succeeds;
- help goes to the supplied output writer and does not touch configuration or
  the filesystem;
- invalid command syntax is distinguishable from a runtime failure;
- the root command exposes only the commands and shared flags needed by the
  current phase, without prematurely implementing later behavior.

Avoid a full help-text golden file. Assert the stable user contract—command
name, short purpose, usage and currently available command names—so harmless
Cobra formatting changes do not break the suite.

Run the focused test and prove that it fails because the command tree does not
exist. Stop for user review and provide the proposed test commit message.

### Implementation checkpoint

- Introduce a constructor for the root Cobra command with explicit input,
  output and error writers.
- Keep execution errors returned to one top-level place that maps them to exit
  codes; do not call process exit functions inside commands.
- Make `main.go` responsible only for process context, version injection and
  invoking the command tree.
- Disable noisy duplicate usage for runtime errors while preserving helpful
  usage for invalid arguments.
- Remove the empty `skills/skills.go` scaffold once it has no purpose.
- Implement no store or installation behavior in this slice.

Run the phase verification commands, stop for review, and provide the proposed
implementation commit message.

## Slice 2: Configuration and project scope

- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.

### Test checkpoint

Cover only the configuration and scope contract:

- defaults contain the configured repository, `main` branch, agent paths and
  default agents from `DESIGN.md`;
- `XDG_CONFIG_HOME` changes the config root, with the documented home-directory
  fallback when it is absent;
- valid TOML overrides defaults without exposing Viper outside the package;
- first-run persistence creates the default file and does not overwrite an
  existing file;
- a nested directory inside a Git repository resolves to that repository root;
- a directory outside Git resolves to itself and produces one clear warning;
- global scope resolves configured home targets and does not require repository
  detection;
- a read-only configuration request does not create directories or files.

Use per-test temporary directories and explicit environment/config inputs so
the tests can run in parallel. Do not use `os.Chdir`. Split success and failure
tests when that is easier to follow than one table.

Run the focused tests, confirm the intended failures, stop for review, and
provide the proposed test commit message.

### Implementation checkpoint

- Define small configuration structs representing the store, agents and
  defaults in `internal/config`.
- Use a local Viper instance. Do not use Viper globals or pass Viper through the
  application.
- Apply defaults before reading TOML and return validated domain values.
- Expand only the documented home prefix in configured global paths; avoid a
  general expression or environment-substitution language.
- Make first-run writing idempotent and ensure partial writes cannot replace a
  valid existing config.
- Detect the Git root from an explicit starting directory. If no Git repository
  exists, return the starting directory plus a warning rather than treating it
  as an error.
- Expose read-only versus initialise behavior explicitly so later dry-run work
  cannot write by accident.

Run all required verification and review the complete phase diff. Stop for user
review and provide the proposed implementation commit message.

## Phase acceptance criteria

- [ ] No-argument help succeeds and performs no writes.
- [ ] Configuration defaults and TOML overrides are covered by readable tests.
- [ ] Project scope is the Git root from any nested directory.
- [ ] Non-repository use warns once and uses the current directory.
- [ ] Global scope never performs repository detection.
- [ ] Cobra and Viper do not leak into domain packages.
- [ ] `go fix ./...`, formatting, `go test ./...` and `go vet ./...` pass.
- [ ] The user has reviewed all changes and received proposed commit messages.

## Out of scope

- Cloning or syncing the store.
- Reading skill metadata.
- Installation, lock files, backups or updates.
- Final output styling, color and release workflows.

Proposed PR title: `Add CLI configuration and project scope foundation`
