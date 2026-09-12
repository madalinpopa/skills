# Phase 4: Safe installation and lock files

One phase equals one pull request. Read [TODO.md](TODO.md),
[DESIGN.md](DESIGN.md) and the completed earlier phases before starting.

## Outcome

`skills install` resolves project or global targets, deduplicates shared agent
paths, transforms the selected skills, applies each whole-skill plan atomically,
and records full provenance and base hashes in `.skill-lock.json`.

## Dependencies and boundaries

- Target selection combines configuration and scope but returns ordinary paths
  and agent variants to `internal/install`.
- The installer consumes the pure Phase 3 plan. It must not reimplement state
  classification while writing.
- Lock encoding, hashing and filesystem application stay in `internal/install`.
- Use standard filesystem operations. Add no virtual filesystem dependency
  unless a concrete test cannot be expressed with temporary directories.

## Planned commits

1. `test(install): define agent target resolution`
2. `feat(install): resolve and deduplicate desired targets`
3. `test(install): define atomic installation and lock state`
4. `feat(install): apply installs and write skill locks`

These are proposed logical commits and require review before creation.

## Slice 1: Agent targets and desired trees

- [x] Tests reviewed by the user.
- [x] Implementation and verification reviewed by the user.

### Test checkpoint

Define target resolution with small, direct cases:

- no `--agent` selection uses the configured default agents;
- repeated agents do not create repeated work;
- Codex and Gemini resolve to one `.agents/skills` target and are written once;
- Claude receives the Claude frontmatter variant;
- the shared target receives the OpenAI/Gemini variant and optional
  `agents/openai.yaml`;
- project scope resolves beneath the detected repository root;
- global scope resolves beneath the configured home paths;
- an unknown agent is a usage error with the available names;
- selected skills are resolved by exact store name and missing names are
  reported before any write plan is applied.

Do not repeat YAML transformation tests. Assert which already-transformed
variant and destination each selected agent receives. Use a table only where
every case follows the same shape.

Run the focused tests, show the expected failures, stop for review, and provide
the proposed test commit message.

### Implementation checkpoint

- Add a target resolver that accepts explicit configuration, scope and agent
  names.
- Normalize target paths before deduplication.
- Preserve the relationship between a physical target and the agent transform
  it requires; reject a future ambiguous configuration rather than choosing one
  silently.
- Build complete desired file maps for each skill and selected physical target.
- Validate all requested skill names and agents before passing plans to the
  filesystem writer.
- Add the `install` Cobra command with exact-name arguments and project scope by
  default. Do not add interactive selection.

Run all required verification, stop for review, and provide the proposed
implementation commit message.

## Slice 2: Atomic installation and `.skill-lock.json`

- [x] Tests reviewed by the user.
- [x] Implementation and verification reviewed by the user.

### Test checkpoint

Use temporary project and home roots to cover observable installation behavior:

- a new skill installs every desired regular file and writes its lock;
- the lock contains the skill name, configured source, full commit, timestamp
  and hashes of exactly the files written by the CLI;
- lock JSON is deterministic apart from the injected timestamp;
- an existing identical unmanaged skill can be adopted by writing a lock without
  rewriting its content;
- a different unmanaged skill is a conflict and remains untouched;
- one file conflict prevents all writes for that skill across selected targets;
- an unrelated non-conflicting skill can still install in the same invocation;
- a failed filesystem application does not leave a partial skill or a lock that
  claims success;
- a successful no-op does not rewrite content or change its lock unnecessarily.

Inject the clock and the narrow failure point needed to test atomicity. Avoid a
general filesystem abstraction. Do not assert temporary staging names or other
implementation details.

Run the focused tests and prove they fail for missing installation behavior.
Stop for review and provide the proposed test commit message.

### Implementation checkpoint

- Calculate all selected skill plans before applying any of them.
- Apply one skill as a single unit across its selected targets. Stage replacement
  content beside its destination and publish it only after every file and lock
  is ready.
- Write the lock as CLI state, exclude it from the content hash map, and use
  stable relative slash-separated file names as keys.
- Preserve file permissions that matter to bundled executable scripts.
- Reject unsupported filesystem entries only at the copy boundary; do not add
  speculative validation elsewhere.
- Return structured per-skill results for later output and exit-code mapping.
- Implement no backup or `--force` behavior yet; conflicts remain untouched.

Run all required verification, inspect failure cleanup, and stop for user
review. Provide the proposed implementation commit message.

## Phase acceptance criteria

- [x] Project and global targets follow configuration and scope exactly.
- [x] Shared physical targets are written once.
- [x] Each skill installs atomically across all selected targets.
- [x] Conflicts leave the entire skill untouched.
- [x] Locks contain full provenance and accurate base hashes.
- [x] Existing identical content can be safely adopted.
- [x] Default tests use temporary directories and run in parallel where safe.
- [x] `go fix ./...`, formatting, `go test ./...` and `go vet ./...` pass.
- [x] The user has reviewed all changes and received proposed commit messages.

## Out of scope

- Scanning all installed skills.
- Update, removal, force and backups.
- `diff`, final color, verbose output and release configuration.

Proposed PR title: `Install skills atomically with provenance locks`
