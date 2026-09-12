# Phase 6: Diff, dry-run, polished output and release

One phase equals one pull request. Read [TODO.md](TODO.md),
[DESIGN.md](DESIGN.md) and the completed earlier phases before starting.

## Outcome

The CLI exposes the complete designed command set with predictable dry-run and
diff behavior, concise terminal-aware output, stable exit codes, version
reporting, CI and tagged cross-platform releases. README examples describe the
actual executable rather than the original scaffold.

## Dependencies and boundaries

- Keep rendering and exit-code mapping in `cmd` unless repetition clearly
  justifies a small internal output package.
- Inject output writers and terminal/color detection. Do not let domain packages
  print.
- Reconstruct diff base content from lock provenance and the recorded store
  commit; do not weaken the lock to store duplicate file bodies.
- Release configuration receives tests only where a test adds value. Workflow
  and documentation changes still require user review before commit.

## Planned commits

1. `test(cli): define diff dry-run and result output`
2. `feat(cli): add diff dry-run and polished result output`
3. `test(cli): define version and exit-code behavior`
4. `feat(cli): report versions and stable exit codes`
5. `build: add CI release automation and final documentation`

These are proposed logical commits and require review before creation.

## Slice 1: `diff`, dry-run and result presentation

- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.

### Test checkpoint

Define the remaining user-facing workflows:

- `skills diff <skill>` compares installed content with the exact transformed
  base recorded by its lock, showing only the user's edits;
- a missing base commit or mismatched source produces a helpful error and never
  substitutes current store content as the base;
- dry-run reports the same planned actions and attention states as a real
  command but writes no config, store, target, lock or backup;
- dry-run on an uninitialised store exits with the `skills init` instruction;
- default output uses one line per skill and aggregates shared targets;
- `-v` adds affected paths beneath the skill without changing the decision;
- symbols retain meaning when color is disabled;
- color is emitted only for a terminal and is disabled by `--no-color` or
  `NO_COLOR`;
- stdout contains normal results while stderr contains warnings and errors;
- summaries count total, changed and attention results correctly.

Avoid complete ANSI or help golden files. Assert stable messages, ordering,
symbols and absence of writes. Split success and failure tests to keep the flow
clear. Tests that manipulate `NO_COLOR` cannot be parallel; prefer injecting
the resolved setting into the renderer and test environment reading separately.

Run the focused tests, show the expected failures, stop for review, and provide
the proposed test commit message.

### Implementation checkpoint

- Add `diff` by loading the lock, reading the recorded commit from the matching
  store, applying the same agent transformation, and comparing it to `have`.
- Render a conventional unified diff with stable relative paths and no custom
  diff language.
- Thread a read-only/dry-run option through configuration, store and install
  orchestration. Reuse the normal planner and renderer; do not maintain a
  separate simulated algorithm.
- Add one renderer for per-skill results, verbose paths, summaries and warnings.
- Detect terminal capability at the CLI boundary and apply the documented color
  precedence.
- Keep wording concise and actionable, especially for conflict, unavailable,
  foreign-source and backup results.

Run all required verification, inspect representative output with and without
color, stop for review, and provide the proposed implementation commit message.

## Slice 2: Version reporting and exit codes

- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.

### Test checkpoint

Cover the final process contract:

- `skills version` reports the CLI version and full current store commit as
  separate facts;
- a development binary reports a clear development version;
- a missing store is reported without hiding the CLI version;
- complete success maps to exit code 0;
- runtime failure maps to exit code 1;
- invalid commands or arguments map to exit code 2;
- any skipped conflict, unavailable skill or foreign-source installation maps to
  exit code 3, even when other skills changed successfully;
- help and version paths do not accidentally initialise or mutate state.

Test the command executor's mapping directly and keep only one small process
smoke test if it catches wiring that command tests cannot. Do not duplicate each
command suite at the process level.

Run the focused tests, demonstrate the intended failures, stop for review, and
provide the proposed test commit message.

### Implementation checkpoint

- Resolve the CLI version from build information with the release-injected
  `main.version` override described in `DESIGN.md`.
- Let `version` report an absent or unavailable store separately rather than
  forcing first-run network access merely to print the binary version.
- Centralize exit-code mapping at the executable boundary.
- Ensure Cobra argument errors remain usage errors while service failures remain
  runtime errors and attention results remain code 3.
- Check every command and flag listed in `DESIGN.md` for consistent help,
  validation and scope wording.

Run all required verification, stop for review, and provide the proposed
implementation commit message.

## Slice 3: CI, release and documentation

- [ ] Configuration and documentation reviewed by the user.

This slice does not require a synthetic failing unit test. Before editing,
verify current official GitHub Actions, GoReleaser and Go documentation rather
than copying an old workflow.

- Add pull-request CI for formatting, `go vet ./...` and `go test ./...`.
- Add tagged `v*` release automation with GoReleaser for macOS, Linux and
  Windows on amd64 and arm64, plus checksums and generated release notes.
- Inject the release version into `main.version` exactly as documented.
- Keep integration tests out of default CI unless a justified integration suite
  now exists. If it exists, add a separate explicit job with
  `INTEGRATION=true` and the required Docker capability.
- Update README installation, first-run, commands, scope, conflict, backup,
  dry-run and update examples to match the implemented CLI.
- Re-check every statement in `DESIGN.md` against the finished behavior and
  correct documentation drift without redesigning the application.
- Run a local release configuration check and a snapshot build if the current
  GoReleaser documentation supports them.
- Provide the proposed build commit message and stop for user review. Do not
  commit, tag, publish or open the PR without explicit authorization.

## Phase acceptance criteria

- [ ] `diff` shows edits against the recorded transformed base.
- [ ] Dry-run and real execution share one planner, and dry-run writes nothing.
- [ ] Default, verbose, colored and colorless output match `DESIGN.md`.
- [ ] Exit codes 0, 1, 2 and 3 are stable and covered without redundant tests.
- [ ] Version output separates CLI and store versions.
- [ ] All documented commands and flags have consistent help and errors.
- [ ] CI covers formatting, vet and the default test suite.
- [ ] Tagged releases target every documented platform and architecture.
- [ ] README and design match the final behavior.
- [ ] `go fix ./...`, formatting, `go test ./...` and `go vet ./...` pass.
- [ ] The user has reviewed all changes and received all proposed commit
  messages.

## Out of scope

- JSON output, TUI, automatic pruning or pinned Git refs.
- Multiple stores, registries, daemons or background updates.
- Tests whose only purpose is increasing a coverage percentage.

Proposed PR title: `Complete the CLI experience and release pipeline`
