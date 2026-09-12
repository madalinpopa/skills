# Fix plan: design drift

A review compared the code with [DESIGN.md](DESIGN.md) and found places where
they differ. This plan fixes them.

Each phase is one pull request with 1 to 5 commits. Follow the review and commit
workflow in [AGENTS.md](../AGENTS.md): failing tests first, stop for review,
then the implementation, then stop again. Never commit or open a PR without the
user asking.

Tick a box only after the user has reviewed that step.

## Phases

- [ ] [Phase 1: `sync --dry-run` writes nothing](#phase-1-sync---dry-run-writes-nothing)
- [ ] [Phase 2: `install --force`](#phase-2-install---force)
- [ ] [Phase 3: Honest dry-run and verbose output](#phase-3-honest-dry-run-and-verbose-output)
- [ ] [Phase 4: Scope hints and read-only config](#phase-4-scope-hints-and-read-only-config)
- [ ] [Phase 5: Agent format comes from config](#phase-5-agent-format-comes-from-config)
- [ ] [Phase 6: Document how installed skills are found](#phase-6-document-how-installed-skills-are-found)
- [ ] [Phase 7: Rename DESIGN.md to SPEC.md](#phase-7-rename-designmd-to-specmd)

---

## Phase 1: `sync --dry-run` writes nothing

### Problem

The design says `--dry-run` writes nothing. `skills sync --dry-run` still
fetches and fast-forwards the store, because `cmd/store.go:43` calls
`store.Sync` without looking at the flag. A test run moved the store from
`a784eec` to `92ddaff`.

### Fix

In dry-run, find the remote branch head without changing the store, for example
with `git ls-remote origin <branch>`. Then print what would happen:

```
$ skills sync --dry-run
  would pull store  a784eec -> 92ddaff
```

Print the "up to date" line when the heads match. Do not run `fetch` or `merge`.
The store's working tree, refs and objects must stay the same.

### Commits

1. `test(store): define a sync plan that does not change the store`
2. `fix(store): make sync --dry-run read-only`

### Checklist

- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.

### Done when

- A test with a local origin repo shows that `HEAD`, refs and the working tree
  are unchanged after `sync --dry-run`.
- The output shows the current and remote commits.

---

## Phase 2: `install --force`

### Problem

The design lists `--force` as a general flag, and after a conflict the output
says `or 'skills install --force' to overwrite (backed up)`. But `install` has
no `--force` flag (`cmd/install.go:43`), and `resolveInstall` always passes
`force=false` (`cmd/install.go:53`). Running the suggested command fails with
`unknown flag: --force` and exit code 2. The only way past an install conflict
is to delete the directory by hand.

### Fix

Add `--force` to `install` and pass it to the installer, the same way `update`
does. The installer already knows how to back up and overwrite.

### Commits

1. `test(cli): define install --force with backup`
2. `fix(cli): add --force to install`

### Checklist

- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.

### Done when

- `install --force` over an edited or untracked skill directory backs it up,
  prints the backup path, installs and exits 0.
- `install --force --dry-run` shows the plan and writes no backup.
- The conflict hint printed by `install` works when copied and run.

---

## Phase 3: Honest dry-run and verbose output

### Problems

1. **`remove -v` lists no files.** `-v` should show the file paths under each
   skill. `renderer.paths` reads `res.Targets` (`cmd/render.go:111`), but
   `Installer.remove` (`internal/install/backup.go:59`) never fills it.
2. **Dry-run output uses the past tense.** `remove --dry-run demo` prints
   `- demo   removed` and exits 0, though nothing was removed. `install` and
   `update` dry runs print "added" and "updated" the same way.

### Fix

1. Fill the removal plan's targets with the files that will be removed, so
   `-v` prints them like it does for install and update.
2. Pass the dry-run state to the renderer and use wording that says nothing has
   happened yet, for example `would add`, `would update`, `would remove`. Keep
   the symbols the same. Agree the exact words with the user before writing
   tests.
3. Add the dry-run wording to the Output section of `DESIGN.md`.

### Commits

1. `test(install): define verbose removal paths`
2. `fix(install): report the files a removal touches`
3. `test(cli): define dry-run result wording`
4. `fix(cli): say what a dry run would do`
5. `docs: describe dry-run output`

### Checklist

Slice 1: verbose removal paths

- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.

Slice 2: dry-run wording

- [ ] Wording agreed with the user.
- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.
- [ ] `DESIGN.md` update reviewed by the user.

### Done when

- `remove -v` lists every file under each removed skill, for every target.
- No dry-run line uses a past-tense result word.
- The summary line still counts skills the same way.

---

## Phase 4: Scope hints and read-only config

### Problems

1. **The "only installed globally" hint exists only in `remove`.** The design
   states it as a general scope rule: "If a requested skill exists only
   globally, the error says so and suggests re-running the command with
   `--global`." `update foo` and `diff foo` just say `not installed: foo`. The
   only place with the hint is `cmd/remove.go:29`.
2. **Commands that don't need the store still create the config.** `remove`,
   `diff`, `ls --local` and `ls --global` go through `config.Init`
   (`cmd/app.go:40`). The design says only commands that need the store run the
   first-run steps. A first `ls --local` writes to `~/.config/skills/`.

### Fix

1. Move the global lookup out of `remove` into a shared helper, and use it in
   `update` and `diff` when the named skills are not found in project scope.
2. Let each command say whether it needs the store. Commands that only scan
   installed skills load the config read-only with `config.Load`. `diff` needs
   the store to rebuild the base content, so it keeps the first-run steps.

### Commits

1. `test(cli): define the global-scope hint for update and diff`
2. `fix(cli): suggest --global in update and diff`
3. `test(cli): define which commands create the config`
4. `fix(cli): load the config read-only when the store is not needed`

### Checklist

Slice 1: global-scope hint

- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.

Slice 2: read-only config

- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.

### Done when

- `update foo` and `diff foo` suggest `--global` when `foo` is only installed
  globally, and give the plain error otherwise.
- `remove`, `ls --local` and `ls --global` on a machine with no config leave the
  config directory missing.
- `init`, `sync`, `ls`, `install`, `update` and `diff` still run the first-run
  steps.

---

## Phase 5: Agent format comes from config

### Problem

The format written for an agent is picked from its name. `variantOf` returns
the Claude format only for an agent named exactly `claude`
(`internal/install/targets.go:82`). The design says "Supporting a new tool is an
edit here, not a release". But an agent added to the config with a
Claude-style directory would get the other format, with no `x-claude` lifting.
Fixing that today needs a code change.

### Decision needed first

Agree with the user how the config says which format an agent uses. Options:

- An optional `format = "claude"` or `format = "agents"` key per agent, where a
  missing key means `agents`, except for the `claude` agent, which defaults to
  `claude` so existing configs keep working.
- A required `format` key, written into the default config, with an error for
  agents that lack it.

Do not write tests until the user picks one.

### Fix

Read the format from the agent's config entry, validate it in
`internal/config`, and use it in `install.Targets`. Keep the check that two
agents sharing a directory must use the same format. Update the Configuration
section of `DESIGN.md` and the README.

### Commits

1. `test(config): define the per-agent skill format`
2. `feat(config): read each agent's skill format from the config`
3. `docs: document the agent format setting`

### Checklist

- [ ] Config shape agreed with the user.
- [ ] Tests reviewed by the user.
- [ ] Implementation and verification reviewed by the user.
- [ ] Docs update reviewed by the user.

### Done when

- A new agent with the Claude format gets `x-claude` lifted and no
  `agents/openai.yaml`, whatever its name.
- An unknown format value is a config error.
- The default config behaves exactly as it does today.

---

## Phase 6: Document how installed skills are found

### Problem

The design says installed skills are found "by scanning the configured agent
directories". The code scans only the default agents' directories, or the ones
chosen with `--agent`. We keep the code as it is and change the design to match.

What the code does today:

- `ls --local`, `ls --global`, `update`, `remove` and `diff` scan the
  directories of `defaults.agents`.
- `update`, `remove` and `diff` accept `--agent` to scan other agents.
- `ls` does not accept `--agent`, so it lists only the default agents'
  installs.

### Fix

Update `DESIGN.md` so it says this plainly:

- Per-skill lock section: installed skills are found by scanning the directories
  of the default agents, or of the agents named with `--agent`.
- Commands section: the `--agent` flag applies to `install`, `update`,
  `remove` and `diff`, and `ls --local` and `ls --global` list only the default
  agents' installs.
- Say that a skill installed with `--agent` for a non-default agent needs the
  same `--agent` to be updated, removed or diffed.

Make the README match.

### Commits

1. `docs: describe installed-skill scanning as implemented`

### Checklist

- [ ] Docs update reviewed by the user.

### Done when

- No sentence in `DESIGN.md` or the README says every configured agent is
  scanned.
- The flag list shows which commands accept `--agent`.

---

## Phase 7: Rename DESIGN.md to SPEC.md

### Fix

1. Rename `docs/DESIGN.md` to `docs/SPEC.md` with `git mv`, so history follows
   the file.
2. Change its title from `# Design` to `# Spec`, and fix any wording that
   calls the document "the design" where that now reads wrong.
3. Update every reference:
   - `AGENTS.md` line 3
   - `README.md` line 204, including the "Design" section heading
   - `docs/TODO.md`, if this file still exists
4. Search the whole repository for `DESIGN` and `design doc` once more, and
   confirm nothing still points at the old path.

### Commits

1. `docs: rename DESIGN.md to SPEC.md`

### Checklist

- [ ] Rename and reference updates reviewed by the user.

### Done when

- `docs/DESIGN.md` no longer exists and `docs/SPEC.md` does.
- `grep -rn DESIGN .` finds nothing outside `.git`.
- Every link to the spec opens the right file.

---

## Not planned

- **Stale lock commit.** When an update finds no file changes, the lock keeps
  the old `commit` and `installed` values. `diff` still works because the
  content is the same. Revisit only if this causes a real problem.
