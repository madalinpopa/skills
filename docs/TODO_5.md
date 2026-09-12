# Phase 5: Installed-skill lifecycle

One phase equals one pull request. Read [TODO.md](TODO.md),
[DESIGN.md](DESIGN.md) and the completed earlier phases before starting.

## Outcome

Users can inspect installed skills, update safe installations, understand
unavailable or foreign-source state, remove only from an explicit scope, and
use `--force` without losing the previous content because every destructive
operation creates a visible backup.

## Dependencies and boundaries

- Discover installations through `.skill-lock.json`; do not add a central
  project index.
- Reuse the Phase 3 planner and Phase 4 application path.
- Keep backup paths, clock use and filesystem mutation in `internal/install`.
- Keep scope and argument behavior in `cmd`; no domain package should inspect
  Cobra flags.

## Planned commits

1. `test(install): define installed listing and update behavior`
2. `feat(install): scan and update installed skills`
3. `test(install): define removal force and backup safety`
4. `feat(install): remove and force-update with backups`

These are proposed logical commits and require review before creation.

## Slice 1: Installed discovery, local/global listing and update

- [x] Tests reviewed by the user.
- [x] Implementation and verification reviewed by the user.

### Test checkpoint

Cover the lifecycle contract without repeating the planner matrix:

- project scanning starts at the detected repository root and finds only valid
  lock-bearing skill directories under configured targets;
- global scanning reads only global targets;
- shared Codex/Gemini targets are reported once;
- results have stable ordering and identify the physical target rather than
  pretending the shared target belongs to one agent;
- `ls --local` and `ls --global` never fall back to the other scope;
- combining mutually exclusive local and global listing flags is a usage error;
- updating all installed skills applies clean updates and no-ops unchanged
  skills;
- naming skills updates only those exact installations;
- locally edited content skips the entire skill and leaves it unchanged;
- a missing or newly draft store skill is reported unavailable and retained;
- a lock from another source is reported and retained;
- mixed results update safe skills while returning an attention result for
  skipped skills.

Use temporary roots with small installed fixtures. Reuse helpers for writing a
valid lock, but keep each test's Arrange section explicit enough to understand
the scenario. Parallelize tests that own separate roots.

Run the focused tests, show the intended failures, stop for review, and provide
the proposed test commit message.

### Implementation checkpoint

- Scan only configured target directories and treat the lock as the marker of a
  managed installation.
- Decode and validate lock provenance before constructing `base` state.
- Resolve base content from recorded state without weakening the hash checks.
- Add local and global modes to `ls` and add `update [skill...]` with project
  scope as the default.
- Reuse the target deduplication and whole-skill application paths.
- Return explicit unchanged, updated, conflict, unavailable and foreign-source
  results for output code.
- Do not delete unavailable skills or silently adopt a different source.

Run all required verification, stop for review, and provide the proposed
implementation commit message.

## Slice 2: Explicit removal, `--force` and backups

- [x] Tests reviewed by the user.
- [x] Implementation and verification reviewed by the user.

### Test checkpoint

Test the destructive paths carefully but without duplicating safe-update tests:

- `remove` defaults to project scope and never removes a global installation;
- when a skill exists only globally, project removal fails with a concise
  suggestion to use `--global`;
- every successful removal creates a timestamped backup before deleting the
  managed directory;
- an edited skill is retained and requires attention without `--force`;
- forced removal backs up the exact edited content and then removes it;
- forced update backs up the old content before replacement;
- backup paths use the configured XDG config root, preserve enough target
  identity for restoration, and are returned for display;
- a backup failure prevents removal or replacement;
- removing one skill cannot remove siblings or parent target directories;
- explicit removal can remove a foreign-source installation because it does not
  adopt or rewrite that installation.

Inject a clock so backup names are deterministic. Use real temporary files for
backup assertions and only the narrowest injectable failure needed for the
backup-before-mutation guarantee.

Run the focused tests, confirm the expected failures, stop for review, and
provide the proposed test commit message.

### Implementation checkpoint

- Add one backup operation shared by removal and forced replacement.
- Place backups under the resolved config root and include timestamp, skill and
  target identity without exposing unsafe path traversal.
- Create and verify the backup before any destructive change.
- Make `remove` inspect only the selected scope and selected agent targets.
- Preserve edited skills unless `--force` is explicitly present.
- Make forced update reuse the normal desired content and application path after
  backup; do not create a second update algorithm.
- Return exact backup paths and attention results as structured data.

Run all required verification and manually inspect one temporary backup layout
through a test artifact or focused test output. Stop for review and provide the
proposed implementation commit message.

## Phase acceptance criteria

- [x] Local and global listing are stable, explicit and deduplicated.
- [x] Clean installations update and edited installations remain untouched.
- [x] Missing, draft and foreign-source skills remain installed.
- [x] Removal never falls back from project to global scope.
- [x] Every removal and forced overwrite creates a backup first.
- [x] Backup failures cause no destructive change.
- [x] Mixed updates return enough information for exit code 3 in Phase 6.
- [x] `go fix ./...`, formatting, `go test ./...` and `go vet ./...` pass.
- [x] The user has reviewed all changes and received proposed commit messages.

## Out of scope

- Human diff rendering.
- Final symbols, color, verbose paths and dry-run presentation.
- Version output, CI and release automation.

Proposed PR title: `Manage installed skills without losing local changes`
