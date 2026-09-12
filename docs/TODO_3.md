# Phase 3: Catalog listing and change planning

One phase equals one pull request. Read [TODO.md](TODO.md),
[DESIGN.md](DESIGN.md) and the completed earlier phases before starting.

## Outcome

Users can list published store skills, while the application has a pure,
well-tested planner that classifies desired, installed and last-written files
without touching the filesystem. This PR establishes the safety decisions that
all later mutations must obey.

## Dependencies and boundaries

- Catalog behavior belongs in `internal/skill`; Cobra only renders and routes.
- The three-map planner belongs in `internal/install` but performs no I/O.
- Use byte hashes and normalized relative paths as input data. Do not make the
  pure planner read directories, Git or configuration.
- Keep output basic in this phase. Final symbols, color and verbosity belong to
  Phase 6.

## Planned commits

1. `test(catalog): define published skill listing`
2. `feat(catalog): list skills from the store`
3. `test(install): define three-way change planning`
4. `feat(install): add the pure whole-skill planner`

These are proposed logical commits and require review before creation.

## Slice 1: Published catalog and `skills ls`

- [x] Tests reviewed by the user.
- [x] Implementation and verification reviewed by the user.

### Test checkpoint

Define catalog behavior without duplicating Phase 2 parser tests:

- published skills appear with name, description and tags;
- drafts remain absent from normal listing;
- results have a stable name-based order independent of filesystem order;
- an invalid published skill produces a contextual error rather than a partial,
  misleading catalog;
- an empty store produces a successful empty result;
- the `skills ls` command obtains the store through the application boundary and
  renders one line per skill;
- no-argument `ls` does not inspect project or global installations.

Assert stable fields and ordering rather than spacing that Phase 6 will change.
Use parallel subtests with isolated fixture stores.

Run the focused tests, prove the expected failures, stop for review, and provide
the proposed test commit message.

### Implementation checkpoint

- Build the catalog from Phase 2 metadata results and filter by publication
  status at the catalog boundary.
- Sort once before returning results.
- Add the Cobra `ls` command for the store view only; reserve local/global
  routing for Phase 5.
- Keep command behavior behind injected services and writers so tests do not
  mutate user configuration.
- Avoid search, filtering and JSON output until a real consumer needs them.

Run all required verification, stop for review, and provide the proposed
implementation commit message.

## Slice 2: Pure three-way and whole-skill planning

- [x] Tests reviewed by the user.
- [x] Implementation and verification reviewed by the user.

### Test checkpoint

Use a compact `map[string]struct` table for the core file states:

- desired file absent on disk becomes an add;
- `have` equal to `want` is unchanged;
- `have` equal to `base` becomes a safe update;
- `have` different from both `want` and `base` is a conflict;
- a file removed from the desired version is safely removed only when it still
  equals `base`;
- an upstream-removed file edited locally is a conflict;
- a local-only file absent from both `want` and `base` is treated as user data
  and prevents replacement rather than being deleted.

Add separate whole-skill tests proving:

- one conflicting file makes the entire skill conflict across all selected
  targets;
- a conflict does not erase the planned results for unrelated skills;
- an installed skill missing from the store is unavailable, not removed;
- a draft that was previously installed is unavailable, not removed;
- a lock from a different source requires attention;
- deterministic input produces a deterministic ordered plan.

Do not test hashing-library internals. Construct clear maps with meaningful file
names and content, then assert the user-visible decision. Keep Arrange, Act and
Assert obvious in every case.

Run the focused tests, confirm the planner is missing, stop for review, and
provide the proposed test commit message.

### Implementation checkpoint

- Define small plan/result types for add, update, remove, unchanged, conflict
  and unavailable states.
- Compare the union of paths from `want`, `have` and `base` so removed and
  local-only files cannot disappear from the decision.
- Aggregate file decisions into one skill decision only after every selected
  target has been evaluated.
- Make conflict and unavailable results descriptive enough for `cmd` to explain
  them without re-running the comparison.
- Keep ordering deterministic for readable output and stable tests.
- Perform no filesystem access, backups or writes in this package path.

Run all required verification and review the package API for clarity rather
than extensibility. Stop for user review and provide the proposed implementation
commit message.

## Phase acceptance criteria

- [x] `skills ls` lists only published skills in stable order.
- [x] Draft and invalid metadata behavior is clear and tested once.
- [x] Every three-map state in `DESIGN.md` has direct readable coverage.
- [x] File removal and local-only files cannot lose user work.
- [x] Any file conflict skips the whole skill across selected targets.
- [x] Missing, draft and foreign-source installs are attention states, not
  automatic mutations.
- [x] Planning is pure and deterministic.
- [x] `go fix ./...`, formatting, `go test ./...` and `go vet ./...` pass.
- [x] The user has reviewed all changes and received proposed commit messages.

## Out of scope

- Applying a plan to disk.
- Writing or scanning lock files.
- Local/global installed listings.
- Backups, force behavior and final output styling.

Proposed PR title: `Add skill catalog listing and safe change planning`
