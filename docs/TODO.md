# Implementation TODO

Validated against the current code, tests, [SPEC.md](SPEC.md), and official
references on 2026-09-13. All phases are pending. This is the proposed plan;
the spec remains the current contract until each behavior change lands.

## Validation of the original items

| Original item | Finding | Planned work |
| --- | --- | --- |
| Download only `skills/` | Clone flags alone are insufficient: `Store.Tree` reads every repository blob. Sparse checkout limits working files; partial clone reduces transferred objects. | Phases 4–5: narrow committed reads, then optimize clone/sync. |
| Remove Gemini references | Its built-in alias duplicates Codex paths and the shared transform. Removing all mentions would hide continuing compatibility and upgrade implications. | Phase 1: remove the default alias and redundant examples, preserve custom agents, document compatibility. |
| Missing `status` and `tags` | `Parse` rejects absent status, blocking otherwise valid external skills. Missing tags already works and means no tags. | Phase 2: default absent status to published; preserve optional tags and validate explicit values. |
| Standard and Claude frontmatter | Unknown top-level fields already pass through. `x-claude` already lifts for Claude and drops for shared output. There is no vendor allowlist to remove. | Phase 3: add representative compatibility coverage and correct the documented limits. |

### Corrections and proposed decisions

- The format is **Agent Skills**, published at agentskills.io, rather than
  “agentio.” Gemini supports project/global `.agents/skills` aliases, with
  precedence over `.gemini/skills`.
  [Gemini discovery reference](https://geminicli.com/docs/cli/skills/#discovery-tiers)
- Standard fields are `name`, `description`, `license`, `compatibility`,
  `metadata`, and experimental `allowed-tools`. Author/version belong under
  `metadata`; `tools` and `hooks` are not standard fields.
  [Agent Skills specification](https://agentskills.io/specification)
- Claude adds fields such as `model`, `effort`, `context`, `hooks`, `paths`,
  and invocation controls. Preserve their values without interpreting or
  executing them. Claude accepts omitted identity fields; this CLI requires
  a directory-matching name and a description. Keep that portable identity
  requirement and document that some Claude files need adaptation.
  [Claude frontmatter reference](https://code.claude.com/docs/en/skills#frontmatter-reference)
- Keep existing transforms: top-level fields pass through; authors use
  `x-claude` for selective exclusion from shared output. Do not automatically
  relocate fields or introduce a vendor allowlist. Top-level `allowed-tools`
  is portable metadata, not automatically Claude-only.
- External stores still need `skills/<name>/SKILL.md`, valid YAML frontmatter,
  and supported regular files. Arbitrary Markdown, alternate repository roots,
  vendor runtime emulation, and a complete standards linter are separate work.
- Prefer full-history partial clone with cone-mode sparse checkout of
  `skills/`. Cone mode also includes repository-root files, and Git still
  downloads history metadata. Promise reduced unrelated content, not an absolute
  folder-only transfer. Strict root-file exclusion would need a separately
  justified non-cone design.
  [Sparse checkout](https://git-scm.com/docs/git-sparse-checkout)
- Avoid shallow history: `diff` needs the commit recorded in an installation
  lock, which may predate a fresh clone. Partial clone retains commit history
  but can fetch missing blobs later; phase 5 must define offline behavior and
  prevent downloads during read-only operations.
  [Clone options](https://git-scm.com/docs/git-clone),
  [partial clone behavior](https://git-scm.com/docs/partial-clone)

## Delivery rules

Each phase is one focused PR with 1–5 commits. Follow the order below;
phase 1 is independent, phase 3 follows phase 2, and phase 5 requires phase 4.
Commit messages and PR titles are proposals, not authorization to create them.

For each behavior change, follow [AGENTS.md](../AGENTS.md):

1. Re-read the relevant spec and current official dependency documentation.
   Invoke Modern Go guidance before editing Go files, using the exact file or
   Go 1.27 for a new file, and read the complete output.
2. Add the smallest useful tests, demonstrate the intended failure, and stop
   for user review with the diff, failure output, and test commit message.
3. Implement only after the user says to proceed. Update the spec in the
   implementation commit so behavior and contract land together.
4. Complete verification and stop again with the diff, results, and proposed
   implementation commit message. Create commits only when requested.
5. At phase completion, check its acceptance tasks and present the PR title
   listed at its end. Open the PR only when the user asks.

Phase 3 primarily verifies existing behavior. Report passing characterization
tests honestly; do not manufacture a failure or unnecessary implementation.

## Phase 1 — Simplify default agents

**Budget: 2 commits. Dependencies: none.**

Outcome: new configs offer `claude` and `codex`; Gemini continues to read the
shared target, and existing explicit agent configurations keep working.

- [ ] Update the smallest useful tests in `internal/config/config_test.go`
  and existing unknown-agent assertions in `internal/install/targets_test.go`
  and `cmd/install_test.go`. Verify both fresh defaults and written TOML omit
  Gemini, and available-agent diagnostics reflect the resulting map.
- [ ] Preserve de-duplication coverage by defining two custom shared agents in
  the target fixture. Replace the redundant Gemini transform case in
  `internal/skill/skill_test.go` with a custom shared-agent case.
- [ ] Verify an existing config explicitly defining `gemini` still loads and
  remains byte-for-byte unchanged during initialization. Keep such names as
  ordinary configured agents; do not reject or silently rewrite them.
- [ ] After test review, remove Gemini from both `Default` and `defaultTOML`
  in `internal/config/config.go`. Keep path selection and transform dispatch
  generic; no special-case implementation is needed in install or skill.
- [ ] Update SPEC Configuration / Agents and targets, README, and demo wording.
  Use shared-agent terminology for redundant examples; retain a compatibility
  note that `--agent codex` selects the shared path, while `--agent gemini`
  requires an explicit config entry after this change. Installs/locks need
  no migration.
- [ ] Acceptance: fresh defaults contain two agents, custom aliases still
  de-duplicate, existing configs survive, and diagnostics are accurate.
  Run common implementation verification below.

Proposed commits:

1. `test(config): specify defaults without the Gemini alias`
2. `feat(config): remove the redundant default Gemini alias`

**PR title after completion: Simplify default agents while preserving shared-path compatibility**

## Phase 2 — Accept optional publishing metadata

**Budget: 2 commits. Dependencies: none; recommended after phase 1.**

Outcome: a standard skill can be listed and installed without adding this
repository's publishing fields.

- [ ] Replace the missing-status failure expectation in
  `internal/skill/skill_test.go`. Cover absent status, explicit published/draft,
  and invalid explicit values, including empty/null. Only absence defaults to
  published; malformed or unsupported explicit values remain errors.
- [ ] Cover absent tags and an empty string list as no tags, a valid list, and
  invalid explicit shapes/types. Proposed rule: reject null and non-string
  elements rather than coercing bad metadata. Do not require an allocated empty
  slice when its representation has no user-visible effect.
- [ ] Extend catalog coverage to include a status-free skill while retaining
  draft exclusion. Add one command-level install case using a local external
  repository without status/tags; verify installed files and lock. Keep parser
  edge cases at the skill layer.
- [ ] After test review, update source validation in
  `internal/skill/skill.go` to distinguish absence from invalid explicit
  values. Keep installed-content handling in `Describe` separate from source
  publishing rules; damaged local metadata must not block diff or removal.
- [ ] Update SPEC Skill metadata / Drafts / Releases and demo authoring guidance.
  Omitting status publishes a valid skill, so work in progress must explicitly
  use `status: draft`. Retain explicit statuses in existing repository skills
  for clarity.
- [ ] Acceptance: status-free skills list/install, explicit drafts stay
  unavailable, invalid metadata reports its source path, and installed output
  strips store fields. No lock migration is needed. Run common verification,
  including existing damaged-installation tests.

Proposed commits:

1. `test(skill): specify optional publishing metadata`
2. `feat(skill): accept skills without publishing metadata`

**PR title after completion: Support external skills without store-only metadata**

## Phase 3 — Verify frontmatter compatibility

**Budget: 1 commit if current behavior passes; up to 3 if a defect is exposed.**
**Dependency: phase 2 for fixtures without status.**

Outcome: compatibility has explicit evidence and limits, without duplicating
vendor parsers or changing metadata handling that already works.

- [ ] Add representative fixtures in `internal/skill/skill_test.go`: standard
  optional fields with author/version metadata; native top-level Claude fields
  with scalar/list/nested mapping values; and scoped `x-claude` settings with
  invocation controls and hooks. Include an unknown future extension.
- [ ] Verify retained YAML values/types for Claude and shared destinations and
  unchanged body bytes. Include scalar/list forms of `allowed-tools` accepted
  by Claude. Avoid incidental formatting assertions or treating every
  Claude-specific field as standard-compatible.
- [ ] Reuse existing duplicate/reserved-key, collision, alias, CRLF,
  deterministic-output, and sidecar tests. Add only uncovered meaningful cases.
  Duplicate keys remain errors, not a compatibility feature to permit.
- [ ] If the fixtures pass, keep parser/transforms unchanged. If they reveal a
  concrete value-loss or validation defect, isolate the smallest failing case,
  stop for review, then make a local fix in its own implementation commit.
- [ ] Update SPEC Skill metadata / Agent differences and demo guidance to
  distinguish standard, Claude, and store fields. Native top-level Claude
  fields pass through to shared output too; `x-claude` enables selective
  exclusion. There is no configurable transform mechanism to document or add.
- [ ] Explain the identity boundary: metadata-less Claude skills need
  name/description before import. Do not claim full Claude import or complete
  standards validation. Preserving YAML does not prove destination runtime
  support for every retained setting.
- [ ] Acceptance: representative fields survive according to transform policy,
  malformed mappings remain rejected, and documentation matches these limits.
  Run `go test ./internal/skill` and `go test ./...`; use all common verification
  if implementation changes are needed.

Proposed commit when no defect is found:

1. `test(skill): document and verify frontmatter compatibility`

If a defect is found, use a reviewed `test(skill): ...` / `fix(skill): ...`
pair naming the defect, followed by the compatibility documentation commit.

**PR title after completion: Verify Agent Skills and Claude frontmatter compatibility**

## Phase 4 — Restrict committed reads to skills

**Budget: 2 commits. Dependencies: none; must precede phase 5.**

Outcome: catalog, install, update, and diff stop loading unrelated repository
content or causing its download in a partial clone.

- [ ] Extend `internal/store/store_test.go` with skills, unrelated directories,
  root files, and a similarly named sibling directory. Assert `Store.Tree`
  exposes only paths below the exact `skills/` directory, preserving the
  `skills/<name>/...` path shape expected by callers.
- [ ] Retain recorded-commit, executable-bit, symlink rejection,
  committed-content, checkout-conversion, and unknown-commit tests. Cover an
  empty/missing skills tree without hiding an invalid commit as an empty store.
- [ ] After test review, narrow tree enumeration and blob reads in
  `internal/store/store.go` before requesting object content. This belongs at
  the store boundary; loading everything then filtering in commands would
  leave the download problem intact.
- [ ] Trace `cmd/install.go:openStore` through catalog and install requests,
  its update caller, and `cmd/diff.go:diffSkill` against the narrowed filesystem.
  Preserve lock commit identity and committed bytes; keep mutable checkout
  content out of installations.
- [ ] Update SPEC Store layout and committed-content rules. This phase needs
  no config, lock, installed-file, or existing-clone migration.
- [ ] Acceptance: outside content is absent from returned trees; existing
  command behavior and historical diff remain correct. Run common verification
  and `task test:integration` for the Git boundary, accurately describing what
  coverage that task currently provides.

Proposed commits:

1. `test(store): restrict committed trees to the skills directory`
2. `fix(store): read only committed skill content`

**PR title after completion: Limit store reads to committed skill content**

## Phase 5 — Optimize clone and sync safely

**Budget: 4 commits, as two reviewed test/implementation pairs.**
**Dependency: phase 4.**

Outcome: new stores avoid unrelated subdirectories and unnecessary blobs while
preserving branch tracking, provenance, and explicit read/write behavior.

### Task A — Initialize a partial, sparse store

- [ ] Confirm the minimum supported Git version for sparse checkout and
  disabling lazy fetch before writing tests. Document it in README/SPEC and
  define a clear error for unsupported Git versions.
- [ ] Test actual filtering with temporary repositories and a `file://`
  transport whose server enables filtering. A plain local-path clone may
  bypass filtering and cannot prove transfer reduction.
- [ ] Cover a non-default branch, retained older commits, skills/root files
  present, unrelated directories absent, and unrelated blobs still missing.
  Object-presence checks must disable lazy fetching. Avoid timing, bandwidth,
  or pack-size thresholds.
- [ ] After test review, create new full-history single-branch clones without
  a full checkout, request blob filtering, and apply cone-mode sparse checkout
  for `skills/`. Complete setup before accepting the clone as usable; failed
  initialization must not leave a half-configured store accepted on retry.
- [ ] Leave existing full clones usable and unchanged on `Init`; do not
  automatically prune, convert, or reclone them. Test that initialization of
  an existing clone remains read-only: `app.ready` also calls it in dry runs.
- [ ] Test unsupported filtering. Proposed policy: accept a successful
  unfiltered transfer with sparse checkout, expose an optimization warning
  through `cmd` on stderr, and retain actual Git errors. Do not retry
  authentication/transport failures as a different clone strategy.

### Task B — Preserve sync and read-only history access

- [ ] Add focused regressions: sync adds/changes/deletes skills without fetching
  unrelated directory blobs; dirty/diverged stores preserve data; preview
  leaves refs, index, objects, and working files unchanged. Reuse current tests
  where they already prove these outcomes.
- [ ] After test review, preserve filtering and sparsity through fast-forward
  sync. Successful init/sync must leave current skill blobs locally available
  so listing, installation, and update need no implicit network access.
- [ ] Prevent implicit object downloads from `Store.Tree` reads. Git's
  [no-lazy-fetch control](https://git-scm.com/docs/git#Documentation/git.txt---no-lazy-fetch)
  separates reading from writing the object database. Missing content must
  identify the recorded commit and cause, never substitute HEAD or an empty tree.
- [ ] Define historical diff behavior in SPEC: downloaded bases work offline;
  an older base absent from a fresh partial clone produces a missing-content
  error. Provide a documented manual Git recovery path for fetching that
  recorded base before retrying. Automatic history recovery stays out of scope;
  reduced offline history availability is an explicit compatibility tradeoff.
- [ ] Test install followed by sync and offline diff against the original lock,
  plus a fresh clone lacking an older base. Dry-run install/update with missing
  blobs must fail without changing the object database. Snapshot store state
  as well as targets: a lazy fetch is a write even without an installation.
- [ ] Document rollout: only new clones gain transfer savings; full clones
  still benefit from phase 4. Keep config/lock formats unchanged. Document
  manual replacement only after preserving local work, retaining the old clone
  as the recovery option if the optimized clone is unsuitable.
- [ ] Acceptance: transport tests prove excluded blobs remain absent after
  init, sync, and reads; provenance/modes remain intact; fallback and
  missing-object errors are clear; dry runs preserve all inspected state.
  Run common verification and `task test:integration` for the Git boundary.

Proposed commits:

1. `test(store): specify partial and sparse initialization`
2. `feat(store): initialize skills-focused partial clones`
3. `test(store): protect sparse sync and read-only history access`
4. `fix(store): preserve sparse sync and prevent implicit object downloads`

**PR title after completion: Reduce store downloads with safe partial and sparse clones**

## Common implementation verification

- [ ] Run tests for each changed package.
- [ ] Run `go fix ./...` and review its diff.
- [ ] Run `task format`; confirm no Go file is left unformatted.
- [ ] Run `go test ./...` and `go vet ./...`.
- [ ] Run `task test:integration` when touching the Git integration boundary.
  Currently it reruns the suite with `INTEGRATION=true`; this repository has
  no Testcontainers dependency or container tests. Prefer local Git and
  transport fixtures. Add a container only for coverage those cannot provide,
  following AGENTS.md's Ubuntu/Git and environment-gate requirements.
- [ ] Review the full phase diff, confirm SPEC matches behavior, preserve
  unrelated edits, and report every failed or unavailable check.
- [ ] Present completed tasks, verification results, the proposed commit
  message, and the phase's PR title for review.

## Validation performed for this plan

- Read the original items and relevant spec sections; traced config loading,
  target selection, committed store reads, catalog validation, transforms,
  install/update, and lock-based diff.
- Ran `go test ./internal/skill ./internal/config ./internal/store ./internal/install ./cmd`:
  all five packages passed. No production code or tests were changed.
- Checked the official references linked above, plus the
  [YAML v3 API](https://pkg.go.dev/go.yaml.in/yaml/v3) and
  [Viper documentation](https://github.com/spf13/viper#reading-config-files).
- Partial-clone transport, vendor runtime behavior, and proposed new tests
  were not exercised in this planning pass. Baseline tests establish the
  current state; future behavior must be verified within its phase.
