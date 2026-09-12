# Plan: a correct skills CLI

[DESIGN.md](DESIGN.md) defines the intended contract. This plan closes the
remaining correctness gaps; it does not claim they are implemented.
Performance and optional features come later, when there is a concrete need.

Follow [AGENTS.md](../AGENTS.md): the smallest useful failing tests first,
stop for user review, then implement only after approval, run the required
verification and stop for review again. Proposed commit messages are not
permission to commit. Never commit or open a PR unless the user asks.

Each phase is one focused PR. When a phase lists two test/implementation pairs,
review each pair separately. Update relevant README guidance with the behavior
change. Tick review boxes only after the user's review.

## Specification decisions before implementation

The design was challenged against the code rather than turning every difference
into a feature. These are the proposed milestone decisions:

| Question | Decision and reason |
| --- | --- |
| Which agents are scanned? | Keep defaults, replaced by explicit `--agent`. Installed `ls` uses defaults. Document this instead of adding discovery flags. |
| Must formats be configurable? | No. Keep `claude` selecting Claude output and other names selecting shared output. Arbitrary transforms are outside this milestone. |
| Can dry-run use network? | Sync may query the remote, but writes nothing. A remote head does not necessarily prove fast-forward eligibility; say so. |
| Which commands initialise? | Only actual store consumers. Installed listing/removal load config read-only; empty update does not clone. Diff needs recorded store content. |
| What is whole-skill safety? | Known conflicts and invalid filesystem shapes prevent all writes for the skill. Stage before publishing. Unexpected I/O failures can leave explicitly reported partial progress. No rollback engine or crash journal. |
| What does force permit? | Resolve content conflicts after backup, preserving local-only regular files. Never bypass source, lock or filesystem validation. |
| What does the commit mean? | The actual committed content used for installation, including catalog selection. Do not label working-tree edits as `HEAD`. |
| Does discovery need valid frontmatter? | No. A valid lock identifies the installation. Missing/malformed local content stays inspectable and removable; a damaged lock is never an unmanaged overwrite candidate. |
| Are modes necessary? | Track executable state on POSIX so scripts run. Extend locks alongside existing byte hashes; unknown legacy modes need force/backup before mutation. No generalized permission management. |
| Are deleted files edits? | Update restores missing files still wanted upstream, as originally specified. Removal treats tracked deletion as an edit and requires force. |
| Must vendor advice be enforced? | No. Shared metadata and safe transforms are required. Vendor limits and authoring guidance remain reference material. |

The old unconditional "never partially updated" promise was too broad for
sequential filesystem writes. This clarification still requires the reproduced
directory/file collision to be caught before any write. It does not excuse a
known conflict as a runtime partial failure.

Reading a captured commit also avoids ignored files and checkout line-ending
conversion changing what a lock claims was installed. Reuse the existing
commit-reading capability; do not build another store/cache layer.

- [ ] Review the clarified design contract and this plan.

Proposed documentation commit:
`docs: define CLI correctness boundaries and repair plan`

## Phases

- [x] [Phase 1: Read-only sync planning](#phase-1-read-only-sync-planning)
- [x] [Phase 2: Safe target trees](#phase-2-safe-target-trees)
- [x] [Phase 3: Runtime failure and backup reporting](#phase-3-runtime-failure-and-backup-reporting)
- [x] [Phase 4: Force and recovery commands](#phase-4-force-and-recovery-commands)
- [x] [Phase 5: Discovery independent of local metadata](#phase-5-discovery-independent-of-local-metadata)
- [ ] [Phase 6: A lock for every target](#phase-6-a-lock-for-every-target)
- [ ] [Phase 7: Content from the recorded commit](#phase-7-content-from-the-recorded-commit)
- [ ] [Phase 8: Valid portable frontmatter](#phase-8-valid-portable-frontmatter)
- [ ] [Phase 9: Executable state and legacy locks](#phase-9-executable-state-and-legacy-locks)
- [ ] [Phase 10: Dry-run and verbose output](#phase-10-dry-run-and-verbose-output)
- [ ] [Phase 11: Scope hints and first-run boundaries](#phase-11-scope-hints-and-first-run-boundaries)
- [ ] [Phase 12: Config and target validation](#phase-12-config-and-target-validation)

## Phase 1: Read-only sync planning

**Confirmed, high priority; original Phase 1.** `cmd/store.go` calls `Store.Sync`
even in dry-run. The CLI reproduction advanced the real local `HEAD` because
sync fetches and merges.

### Work and acceptance

- Query the configured remote branch head without changing the clone. Show
  local/remote commits or up to date. Do not promise ancestry that cannot be
  checked without fetching missing objects.
- Check dirty/wrong-source/wrong-branch stores consistently with real sync,
  without refreshing the Git index on disk.
- Use a local origin with a newer commit. Assert unchanged working tree,
  index, refs, object inventory and `HEAD`, not just unchanged installed files.
- Remote query failure exits 1. A missing store produces the `skills init`
  explanation with no config, clone or staging files.

### Review slices

1. `test(store): define read-only sync planning`
2. `fix(store): make sync dry-run leave the clone unchanged`

- [x] Failing tests reviewed.
- [x] Implementation and verification reviewed.

## Phase 2: Safe target trees

**Confirmed, high.** `internal/install/install.go:hashTree` ignores non-regular
entries. Staging follows parent symlinks. A reproduced update followed a
`references` symlink and overwrote external content, exiting 0 without backup.
A directory at a desired file path also escaped planning and failed after
other files and one target lock had changed.

### Work and acceptance

- Check all selected targets for unsupported entries and file/directory
  collisions before changing that skill. A skill-directory symlink must not
  disappear from discovery and become a fresh-install candidate.
- Reject target symlinks and special files as attention even with force.
  Leave the whole skill unchanged; healthy skills can proceed. Give the path
  and actual reason instead of always saying "you edited it".
- Keep writes, backup reads and removal inside validated target boundaries.
  Diff must explain unsupported paths instead of following them outside.
- Reproduce a symlink to external content and a directory at a desired file
  path. Assert no external writes, sibling target/lock changes or backup/staging
  side effects for the rejected skill. No symlink support framework is needed.

### Review slices

1. `test(install): define unsupported target and collision safety`
2. `fix(install): reject unsafe target trees before mutation`

- [x] Failing tests reviewed.
- [x] Implementation and verification reviewed.

## Phase 3: Runtime failure and backup reporting

**Confirmed, medium after Phase 2.** `staging.publish` can fail after publishing
files and currently writes locks before deferred removals finish. `Install`
and `Remove` discard earlier results on later errors; commands return before
rendering them. A reproduced two-skill removal deleted the first, failed the
second backup and printed none of the first skill's backup paths.

### Work and acceptance

- Keep stage-before-publish. Publish a target lock only after its content
  writes and removals succeed.
- Stop further mutations on runtime failure. Render completed work and every
  completed backup path, plus a partial-failure diagnostic identifying affected
  skill/targets. Exit 1; do not count a failed skill as completed.
- Complete all required backups for a skill before destructive writes/removal.
  If one backup fails, leave that skill installed and report any copies made.
- Test publication failure after progress and a later-skill backup failure.
  Assert lock ordering and visible recovery paths. Keep these separate from
  Phase 2's known conflicts. No rollback engine or crash simulation.

### Review slices

1. `test(install): define progress and lock state on runtime failure`
2. `fix(install): preserve partial results and publish locks last`

- [x] Failing tests reviewed.
- [x] Implementation and verification reviewed.

## Phase 4: Force and recovery commands

**Confirmed, medium; expands original Phase 2.** `install` has no force flag
and `resolveInstall` passes false. It exits 2 for the suggested flag. The
original claim that manual deletion was the only workaround was incorrect:
a managed skill can be removed with force and installed again.

`cmd/render.go` also drops required names, scope and explicit agents from
hints. A named update hint can broaden the request to all installs.
`Installer.force` substitutes all current files as base, making local-only
files eligible for deletion when another file conflicts; the clarified spec
requires preserving them.

### Work and acceptance

- Add and wire `install --force`. Cover edited managed content and conflicting
  unmanaged regular content, with backup before write and no writes in dry-run.
- Preserve local-only regular files during force; resolve the actual content
  conflicts. Keep foreign-source, malformed-lock and unsafe-tree guards.
- Hints include required names, `--global` and explicit `--agent` selection.
  Removal hints say remove. Do not suggest diff for an unmanaged conflict with
  no recorded base; explain that limitation.
- Keep existing backups safe on a same-second collision. A clear failure is
  sufficient; automatic backup-name retries are deferred.
- Test flag wiring/hints in `cmd` and file preservation in `internal/install`;
  do not duplicate all backup tests at the CLI layer.

### Review slices

1. `test(cli): define install force and scoped recovery commands`
2. `fix(cli): support install force and preserve recovery scope`
3. `test(install): preserve local-only files during forced updates`
4. `fix(install): force conflicts without deleting local-only files`

- [x] CLI failing tests reviewed.
- [x] CLI implementation and verification reviewed.
- [x] Installer failing tests reviewed.
- [x] Installer implementation and verification reviewed.

## Phase 5: Discovery independent of local metadata

**Confirmed, medium.** `internal/install/scan.go:Scan` parses `SKILL.md` before
commands select names. Malformed local frontmatter blocked its own diff and
forced removal, plus a named update of an unrelated healthy skill.

### Work and acceptance

- Use lock/directory identity for discovery; descriptions are optional display
  data. A missing or invalid local `SKILL.md` must not block diff, update
  planning or forced removal of a validly tracked installation.
- Installed listing includes damaged metadata and healthy rows, returning
  attention. Named operations isolate selected skill content.
- Keep absent locks distinct from damaged locks. Validate required provenance,
  identity and relative paths; force must not adopt a damaged or foreign lock.
  All-skills commands report damaged installations as attention and continue
  healthy ones. Runtime root/read failures still exit 1.
- Test a malformed local skill beside a healthy one through the command path,
  plus minimal lock validation at the installer boundary.
- Detect deleted tracked files as removal edits. Test a deleted ordinary file
  separately from metadata parsing; force backs up the remaining directory.

### Review slices

1. `test(install): define discovery with damaged local content`
2. `fix(install): discover managed skills independently of descriptions`
3. `test(install): treat tracked deletions as removal conflicts`
4. `fix(install): protect local deletions during removal`

- [x] Discovery/lock failing tests reviewed.
- [x] Discovery implementation and verification reviewed.
- [x] Removal failing tests reviewed.
- [x] Removal implementation and verification reviewed.

## Phase 6: A lock for every target

**Confirmed, medium.** `planSkill` treats a skill as installed if any target
has a base. With identical content in two targets and only one lock, install
reports zero changes and leaves the other unmanaged; listing omits it.

### Work and acceptance

- Treat missing target locks as work independently of content and sibling
  locks. Adopt identical unmanaged files without rewriting them.
- Extend the existing identical-adoption test with mixed managed/unmanaged
  targets. Verify both locks and discovery; dry-run writes no missing lock.
- Keep fully managed unchanged locks' commits/timestamps unchanged. Do not
  introduce a global lock refresh policy.

### Review slices

1. `test(install): require a lock for each selected installation`
2. `fix(install): adopt identical unmanaged targets independently`

- [ ] Failing tests reviewed.
- [ ] Implementation and verification reviewed.

## Phase 7: Content from the recorded commit

**Confirmed, medium.** `cmd/install.go:openStore` records `HEAD`, but catalog
and rendering read working-tree files. A dirty-store install succeeded and its
immediate diff reported the store edit as a local installation edit.

### Work and acceptance

- Capture one commit per store-reading command. Discover/render its tracked
  skills using the same representation as diff. Reuse `internal/store` commit
  reading rather than a second cache or a clean-check-only workaround.
- Exclude uncommitted edits, untracked skills and ignored files from published
  content. Do not reset or remove local store work. Sync still refuses dirt.
- Preserve tracked executable state and reject source symlinks rather than
  silently omitting them. Retain path/metadata validation at existing boundaries.
- Use local repositories to verify immediate empty diff, excluded unpublished
  bytes and a lock identifying the content used. Include checkout conversion
  so platform working-tree bytes cannot become a false base. No container.

### Review slices

1. `test(store): tie catalog and install bytes to the recorded revision`
2. `fix(store): read published skill content from one commit`

- [ ] Failing tests reviewed.
- [ ] Implementation and verification reviewed.

## Phase 8: Valid portable frontmatter

**Confirmed, medium; line endings clarified during spec review.** `Transform`
appends `x-claude` keys without collision checks. Install accepted duplicate
`name` keys that its own diff rejected, and nested `status` leaked into output.
`split` recognises only LF delimiters despite Windows distribution.

### Work and acceptance

- Validate source `x-claude` as a mapping with unique string keys even when
  selected output would drop it. Reject shared/reserved keys and top-level
  collisions; keep emitted YAML valid and store-only fields absent.
- Accept LF/CRLF delimiters, emit deterministic frontmatter and preserve body
  bytes. Cover this at the pure parser/transform layer.
- Keep ordinary extension pass-through; no vendor field allowlists, description
  budget checks or public validation command.

### Review slices

1. `test(skill): reject unsafe extension keys`
2. `fix(skill): validate extensions before transformation`
3. `test(skill): accept CRLF frontmatter delimiters`
4. `fix(skill): parse portable frontmatter line endings`

- [ ] Extension failing tests reviewed.
- [ ] Extension implementation and verification reviewed.
- [ ] Line-ending failing tests reviewed.
- [ ] Line-ending implementation and verification reviewed.

## Phase 9: Executable state and legacy locks

**Confirmed gap in both implementation and old spec.** Byte-only hashes ignore
mode changes. A synced executable-bit change left the installed file at `0644`
with zero changes reported. Bundled script execution is correctness work.

### Work and acceptance

- Compare content and executable state on POSIX in the pure three-way plan.
  Track one `executable` boolean per managed file alongside existing hashes.
  Apply upstream mode-only changes; protect local mode edits and display diffs.
- Read old locks without treating missing modes as false. Unknown POSIX base
  modes require attention before mutation unless force provides a backup.
  No-op legacy installs stay unchanged; new lock writes record complete modes.
- Preserve executable state in backups. Windows retains content behavior
  without POSIX enforcement. Ownership/ACLs/general permission management wait.
- Cover a real executable script update, local mode edit and old lock.
  Planner tests prove comparison; filesystem/CLI tests prove application/diff.
- Depends on force handling and committed modes in Phases 4 and 7. No separate
  migration command or generalized filesystem metadata model.

### Review slices

1. `test(install): define executable state and legacy lock safety`
2. `fix(install): track executable state in plans and locks`
3. `test(cli): show executable-only local changes in diff`
4. `fix(cli): report executable changes in skill diffs`

- [ ] Installer failing tests reviewed.
- [ ] Installer implementation and verification reviewed.
- [ ] Diff failing tests reviewed.
- [ ] Diff implementation and verification reviewed.

## Phase 10: Dry-run and verbose output

**Confirmed; original Phase 3.** Removal never fills result targets, so the
renderer has no paths. Dry runs say added/updated/removed and count changed
skills as if writes happened; existing tests explicitly assert that wording.

### Work and acceptance

- Removal plans list content paths and the lock under every selected target.
  Use the same plan for verbose output and dry-run.
- Use `would add`, `would update`, `would remove` and `N would change`.
  Count per skill, not per target/file. Skipped plans still exit 3.
- Say a backup would be required without claiming a backup path was created.
- Test paths at the installer boundary and wording in the renderer; keep only
  the command test needed to prove dry-run state reaches rendering.
- Align README output and safety language with the clarified runtime failure
  boundary; do not promise transaction or crash atomicity.

### Review slices

1. `test(install): define verbose removal paths`
2. `fix(install): report the complete removal plan`
3. `test(cli): define dry-run result wording`
4. `fix(cli): describe planned changes without claiming writes`

- [ ] Removal-path failing tests reviewed.
- [ ] Removal-path implementation and verification reviewed.
- [ ] Renderer failing tests reviewed.
- [ ] Renderer implementation, README and verification reviewed.

## Phase 11: Scope hints and first-run boundaries

**Confirmed; corrects original Phase 4 and absorbs Phase 6.** Only remove has
a global hint, and it looks up all requested names, so mixed requests can
mislabel local names or miss a useful hint. `newApp` creates config even for
installed-state commands. Diff was incorrectly included in the old problem
statement: it needs the store to reconstruct its base.

### Work and acceptance

- Share a helper checking only missing local names against selected global
  agents. Use it for update/diff/remove without changing scope.
- Cover a global-only name, truly missing name and mixed request. No lookup
  through unrelated agents and no hinted operation performed automatically.
- Load config read-only for installed listing/removal and initial selection.
  Initialise only when actual store content is needed. Installed listing leaves
  a fresh config directory absent; removal creates only necessary backups;
  empty update does not clone. Preserve store consumers and dry-run boundaries.
- Use local origins for initialization tests. Document exact flag applicability
  and default-agent scanning in README. No `ls --agent` or scan-all expansion.

### Review slices

1. `test(cli): define missing-name scope hints`
2. `fix(cli): suggest global scope only for missing local names`
3. `test(cli): define lazy first-run initialization`
4. `fix(cli): initialize config and store only when needed`

- [ ] Scope-hint failing tests reviewed.
- [ ] Scope-hint implementation and verification reviewed.
- [ ] Initialization failing tests reviewed.
- [ ] Initialization implementation, README and verification reviewed.

## Phase 12: Config and target validation

**Contract gaps found during spec review.** `Config.validate` checks nonempty
paths and defined default names, but not empty default selections or project
containment. `Targets` de-duplicates equal paths but accepts nested targets.
These checks prevent successful no-op installs and writes outside stated scope.

### Work and acceptance

- Reject empty defaults, escaping/absolute project paths and relative global
  paths. Preserve explicit agent selection and existing unknown-agent errors.
- Validate project containment after resolving configured roots too, so root
  links cannot redirect work outside the repository.
- Reject overlapping selected roots except compatible identical roots.
  Preserve shared Codex/Gemini targets and different-format rejection.
- Test the minimal invalid cases and compatible sharing at config/target
  boundaries. Keep Cobra out of those packages. Document constraints in README;
  do not add format configuration or restructure unrelated config loading.

### Review slices

1. `test(config): define valid agent selections and paths`
2. `fix(config): reject empty selections and invalid scope paths`
3. `test(install): reject escaping and overlapping targets`
4. `fix(install): validate resolved target boundaries`

- [ ] Config failing tests reviewed.
- [ ] Config implementation and verification reviewed.
- [ ] Target failing tests reviewed.
- [ ] Target implementation, README and verification reviewed.

## Completion check

- [ ] Complete AGENTS.md verification after each implementation.
- [ ] With a temporary local origin and project/global targets, exercise init,
  listing, install, sync, update, diff, force and remove against documented
  output and exit codes.
- [ ] Confirm dry-run leaves config, store, targets and backups unchanged.
- [ ] Review README and design against completed behavior.
- [ ] Confirm unrelated worktree changes remain untouched.

This is a release sanity exercise, not a second comprehensive end-to-end suite.
Temporary repositories suffice; use containers only for a boundary they uniquely
prove. Keep tests at the lowest useful layer without overlapping coverage.

## Deferred under YAGNI

- **Configurable formats (old Phase 5):** named transforms cover the three
  required tools. A per-agent format setting is an optional feature.
- **Rename DESIGN.md to SPEC.md (old Phase 7):** optional cleanup that does not
  improve management correctness. Keep existing links working.
- **Refresh unchanged lock commits/timestamps:** retain no-op behavior while
  repairing missing locks independently.
- **Rollback, crash recovery and concurrent mutation:** report runtime failure
  honestly now; no journal, automatic rollback or process coordination.
- **Performance, caches, parallelism and benchmarks:** measure a real bottleneck
  before changing the straightforward implementation.
- **More discovery flags, arbitrary transforms, JSON/TUI output, pruning, backup
  retention, backup-name retries, tag/commit pinning, public validation and
  automatic history recovery:** no demonstrated milestone need.

## Audit evidence and limits

The preceding review reproduced sync dry-run writes, missing install force,
output/scope gaps, symlink escape, partial publication, blocked discovery,
missing target locks, dirty-store provenance, duplicate transformed keys,
mode-only no-ops and lost successful-removal output after a later failure.
`go test ./...` and `go vet ./...` passed before these documentation edits;
existing tests do not cover those regressions.

Further requirements from spec review (CRLF handling, removal of locally
deleted files, force preservation, mode compatibility and config/target
validation) follow from the current code paths above. They still need reviewed
failing tests; this is not a claim each case was reproduced through the CLI.
Hosted branch protection and Windows runtime behavior were not verified.
No implementation phase is completed by this documentation update.
