# Phase 2: Git store and portable skill metadata

One phase equals one pull request. Read [TODO.md](TODO.md),
[DESIGN.md](DESIGN.md) and the completed Phase 1 changes before starting. All
review pauses and the no-commit rule remain mandatory.

## Outcome

The CLI can initialise and safely synchronize its managed Git store, report the
full current commit, discover valid skill directories, and produce deterministic
Claude and OpenAI/Gemini variants from one YAML-frontmatter source.

## Dependencies and boundaries

- Keep Git process execution in `internal/store`.
- Pass Git arguments directly to the process runner; never construct a shell
  command string.
- Keep YAML parsing and transformation in `internal/skill`.
- Select a maintained YAML library only after checking its current upstream
  status, documentation and Go 1.27 compatibility.
- Use local temporary Git repositories. Do not use GitHub or another network
  service in the default tests.

## Planned commits

1. `test(store): define clone and fast-forward sync behavior`
2. `feat(store): manage the local Git skill store`
3. `test(skill): define metadata validation and agent transforms`
4. `feat(skill): parse and transform portable skills`

These are proposed logical commits and require review before creation.

## Slice 1: Managed Git store

- [x] Tests reviewed by the user.
- [x] Implementation and verification reviewed by the user.

### Test checkpoint

Define the store behavior with focused tests:

- initialisation clones the configured branch when the store is absent;
- initialisation is idempotent when the expected store already exists;
- a path that exists but is not the configured Git store produces a useful
  error instead of being overwritten;
- synchronization accepts a fast-forward and reports the old and new commits;
- synchronization refuses divergence and never creates a merge commit;
- dirty state is not reset or discarded;
- commit discovery returns the full hash used by lock files;
- command failures retain the useful Git error without dumping irrelevant
  process noise;
- `skills init` and `skills sync` route to the store service and propagate its
  result without embedding Git decisions in Cobra callbacks.

Use a small fake process runner only to verify exact boundary decisions that are
hard to provoke reliably, such as argument selection and process failures. Use
isolated local repositories for observable clone and fast-forward behavior.
Those repository tests may run in parallel because each owns every path it uses.

Testcontainers are not justified for this slice unless local Git cannot prove a
required behavior. If that changes, follow the container policy in `TODO.md`
and explain what unique boundary the Ubuntu container covers.

Run the focused tests, demonstrate the intended failure, stop for review, and
provide the proposed test commit message.

### Implementation checkpoint

- Add a narrow Git runner that accepts context, working directory and argument
  lists and captures the output required for helpful errors.
- Clone only the configured branch into the managed store path.
- Synchronize with an explicit fast-forward-only operation. Do not inherit a
  user's pull strategy in a way that can merge or rebase the store.
- Detect dirty or divergent state and explain how the user can inspect it; do
  not reset it automatically.
- Return old commit, new commit and whether anything changed as data. Leave
  human formatting to `cmd`.
- Add the `init` and `sync` Cobra commands. `init` makes first-run setup
  explicit; `sync` initialises a missing store and otherwise fast-forwards it.
- Preserve first-run behavior from Phase 1: commands that need the store may
  initialise it, while read-only dry-run paths may not.

Run all required verification, stop for review, and provide the proposed
implementation commit message.

## Slice 2: Skill metadata and agent transformation

- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.

### Test checkpoint

Use a small set of readable skill fixtures to cover:

- required `name` and `description` metadata;
- `published` and `draft` status behavior;
- tags and the requirement that directory name and skill name match;
- malformed frontmatter and unsupported status values;
- stripping `status` and `tags` from every installed copy;
- lifting `x-claude` into Claude frontmatter;
- dropping `x-claude` for the shared `.agents` target;
- copying `agents/openai.yaml` only to the `.agents` variant;
- preserving scripts, references, assets and ordinary files;
- deterministic transformed output from the same source bytes.

Do not create a fixture for every YAML syntax feature. The YAML library owns its
parser behavior; these tests own the store schema and transformations. Use a
small `map[string]struct` table where the cases share one clear shape, and use
separate success and failure tests otherwise.

Run the focused tests, show that they fail for the missing parser/transformer,
stop for review, and provide the proposed test commit message.

### Implementation checkpoint

- Discover only direct child directories of the store's `skills/` directory.
- Parse the frontmatter into an explicit metadata model and validate the
  portable store contract.
- Preserve the Markdown body separately from frontmatter transformation.
- Transform YAML structurally and serialize it deterministically. Do not use
  line-based YAML manipulation.
- Represent agent-specific output as data that later installation code can
  place without understanding YAML.
- Return draft skills to catalog code as drafts; do not silently discard them
  inside the parser.
- Return errors with skill and file context while avoiding redundant wrapping.

Run all required verification and inspect dependency changes for necessity.
Stop for review and provide the proposed implementation commit message.

## Phase acceptance criteria

- [ ] Store creation and synchronization use the configured branch only.
- [ ] Synchronization is fast-forward only and preserves dirty/divergent state.
- [ ] Full Git commit hashes are available to callers.
- [ ] Valid skills parse and invalid store metadata fails clearly.
- [ ] Claude and `.agents` outputs match `DESIGN.md` and are deterministic.
- [ ] Tests use local repositories and require neither network nor Docker.
- [ ] `go fix ./...`, formatting, `go test ./...` and `go vet ./...` pass.
- [ ] The user has reviewed all changes and received proposed commit messages.

## Out of scope

- User-facing catalog output.
- Three-way planning and filesystem installation.
- Lock files, updates, removal and backups.
- Final color and verbose output.

Proposed PR title: `Add the Git store and portable skill metadata`
