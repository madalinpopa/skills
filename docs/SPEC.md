# Specification

The behavior of the `skills` CLI and the reasons behind it. The binary
implements this contract; change the document when the behavior changes.

## Goal

The repository is a centralized store of agent skills. The CLI exists to get a
skill out of that store and into any project quickly, and to keep installed
skills current as the store changes.

Skills must work with Claude Code, Codex and Gemini CLI.

The first milestone is a correct CLI for listing, installing, inspecting,
updating and removing skills. Prefer small, explicit behavior over extension
frameworks. Performance work, a TUI, JSON output, configurable transforms,
automatic pruning and crash recovery wait for a demonstrated need.

## How it works

```
github repo          ~/.config/skills/store      a project
(the store)    --->  (local clone)         --->  .claude/skills/
                       skills sync               .agents/skills/
                                                 skills install
```

Three moving parts:

- **The store**: a git clone of the skills repository, kept in the user's home
  directory. One copy serves every project on the machine.
- **The config**: where the store came from and which agents exist.
- **The installs**: copies of individual skills inside a project, or in the
  user's home directory when installed globally.

The content is not compiled into the binary. Skills change far more often than
the CLI does, and building them in would mean a tag, a release and a binary
reinstall to publish a one-line edit. Keeping them in a synced clone means a
push to the store is enough.

## First run

Commands initialise the store only when they need to read it:

1. create the config directory and write a default config
2. clone the store
3. carry on with the command

`skills init` does the same thing explicitly. There is nothing to set up by
hand, but the first run needs network access.

| command | config and store access |
| --- | --- |
| `init`, `sync`, store `ls`, `install` | initialise config and store |
| `update`, `diff` | scan the selected installations first; initialise only if store content is needed |
| `remove`, `ls --local`, `ls --global` | load config or defaults without creating config or store |
| `version` | read existing config and store state; report an absent store |
| help or no arguments | no config or store access |

A removal still creates its backup directory when it actually removes a skill.
Read-only config loading does not prohibit those explicitly requested writes.
An empty `update` selection is a successful no-op and needs no clone.

`--dry-run` is the exception: it never creates the config or clones the store.
If the command cannot calculate a plan because the store is not initialised, it
exits with an explanation and tells the user to run `skills init` first.
It creates no backup, staging file or target directory either.

## Layout

```
~/.config/skills/config.toml            what the store is, which agents exist
~/.config/skills/store/                 git clone of the skills repository
```

Config and store live together under one directory, so there is a single place
to look, back up or delete. `XDG_CONFIG_HOME` is honoured when set.

The store is a clone rather than a cache: everything else depends on it. Version
one tracks a branch and updates it by fast-forward only. Tags and commit pinning
can be added later if there is a concrete need.

The CLI needs Git 2.45 or newer and checks the version before every command
that uses the store. An older Git is an error that names the version found;
there is no fallback.

A new store is a partial, sparse clone of the configured branch. It keeps the
full commit history but downloads file contents only as needed, and it checks
out only `skills/` plus the files at the repository root. If the server
ignores the filter, the clone downloads everything but stays sparse. The clone
is made in a temporary directory and moved into place only after setup
succeeds, so a failed setup never leaves a half-made store. An existing store,
including an older full clone, is used as it is and never converted.

### Store layout

One directory per skill, one `SKILL.md` inside it:

```
skills/
  go-review/
    SKILL.md                shared body and frontmatter
    agents/openai.yaml      optional, copied only to .agents/
    references/             optional, copied to both
```

The CLI reads only committed files under `skills/`, at the commit it installs
from or compares against. It never reads files at the repository root or in
other directories, so the store repository can hold anything else. A commit
without `skills/` is an empty store. A commit that does not exist is an error.

The skill's name is the directory name. There is no per-agent copy: the
differences between agents are small enough to express as data in one file, and
the CLI applies them on install. See [Skill metadata](#skill-metadata).

Keeping one file is what stops the agents' copies drifting apart. It also means
a skill in the store is a real, readable `SKILL.md` that can be opened, reviewed
or used directly.

## Configuration

```toml
[store]
repo = "https://github.com/madalinpopa/skills"
branch = "main"

# Which agents exist, and where each one keeps its skills.
[agents.claude]
project = ".claude/skills"
global  = "~/.claude/skills"

[agents.codex]
project = ".agents/skills"
global  = "~/.agents/skills"

[defaults]
agents = ["claude", "codex"]
```

The map configures agent names and paths. The name `claude` selects the Claude
transform; all other names select the shared agents transform. A custom name
can reuse the shared format. Its directory name does not select a transform.
Supporting an additional transform is outside the first milestone; there is
no per-agent `format` key yet.

`defaults.agents` selects the agents used when `--agent` is absent; it must be
nonempty and contain defined names. An explicit `--agent` replaces that list
and may select any configured name, including a non-default agent.

Project paths must be relative paths contained beneath the project root.
Global paths must be absolute or begin with `~/`, expanded against the home
directory. Empty paths and paths that escape project scope are config errors.
Selected target roots must not overlap, except for identical paths that can
be de-duplicated. Reject different transforms sharing one target.

Several agents may share one directory, for example custom aliases that also
point at `.agents/skills`. Installs are de-duplicated by resolved path: asking
for all of them writes the files once.

Gemini CLI reads `.agents/skills` too, so `--agent codex` already serves it and
there is no built-in `gemini` alias. To select it by name, add an
`[agents.gemini]` entry with the same paths. A config that already defines one
keeps working unchanged.

## Agents and targets

| agent | project | global |
| --- | --- | --- |
| Claude Code | `.claude/skills/` | `~/.claude/skills/` |
| Codex | `.agents/skills/` | `~/.agents/skills/` |
| Gemini CLI | `.agents/skills/` (alias, preferred over `.gemini/skills/`) | `~/.agents/skills/` (alias) |

Only `claude` and `codex` are configured by default. Gemini CLI is served by
the shared `.agents/skills` directory; see [Configuration](#configuration).

All three use the same skill format: a directory holding a `SKILL.md` whose
frontmatter carries `name` and `description`, with optional `scripts/`,
`references/` and `assets/` subdirectories. All three load only the name and
description up front, read the body when the skill triggers, and read other
files on demand.

A plain skill is therefore portable as written. See
[Agent differences](#agent-differences) for where they diverge.

Skills contain regular files and directories; source symlinks and special
files are unsupported. The CLI writes real files. Before changing a skill,
check all selected target trees for symlinks, special files and file/directory
collisions. Report an unsupported target as needing attention and leave the
whole skill untouched, even with `--force`. Never follow a link inside an
installed skill to read, overwrite, back up or remove data outside that skill.
Configured target roots are resolved and checked before use; project targets
must remain within the project root after resolution. A skill directory itself
must not be a symlink.

Only content and the executable/non-executable distinction are managed for
regular files. On systems with POSIX executable bits, writes use `0755` for
executables and `0644` for other files. Other permissions, ownership, ACLs and
timestamps are outside the contract. Windows does not enforce POSIX executable
bits; it still copies and compares file content.

## Scope

Project scope means the root of the current Git repository, regardless of which
subdirectory the command is run from. If the current directory is not inside a
Git repository, the CLI warns the user and uses the current directory as the
project root.

`--global` uses the user's home directories instead. Nothing else changes: the
same skill, the same lock, a different parent path.

Scope is always explicit. Commands default to project scope and never fall back
to global scope. For `update`, `remove` and `diff`, a requested name missing
locally but found in the selected agents' global directories gets a `--global`
hint. Check only missing names; do not label locally installed names as global
only. A hint never changes the selected scope or performs the operation.

## Per-skill lock

Each installed skill records its own state, in the directory it was installed
to:

```
.claude/skills/go-review/
  SKILL.md
  .skill-lock.json
```

```json
{
  "name": "go-review",
  "source": "https://github.com/madalinpopa/skills",
  "commit": "a1b2c3d4e5f678901234567890abcdef12345678",
  "installed": "2026-09-12T18:40:00Z",
  "files": { "SKILL.md": "<sha256>" },
  "executable": { "SKILL.md": false }
}
```

There is no central index. Installed skills are found by scanning the default
agents' directories, or those selected by `--agent`, for `.skill-lock.json`.
Every installed target gets its own lock, including an identical unmanaged
copy adopted by `install` when another target is already managed.

Lock identity and provenance do not depend on parsing the installed `SKILL.md`.
Missing or malformed local frontmatter is still local content: it can be
diffed, checked for conflicts, or removed with a backup. Installed listings
show the directory name and mark an unavailable description. They report
metadata problems as attention without hiding healthy installations.

A missing lock means unmanaged content. A malformed lock is not equivalent
to a missing lock: report it and do not overwrite it automatically, including
with `--force`. Validate required provenance, directory/name agreement and
relative file paths. A named command need not inspect unrelated skill content.
An all-skills command reports a damaged installation as attention and proceeds
with healthy ones; it must not guess the damaged installation's base.

This is what makes the tool stateless about projects. A skill carries its own
provenance, deleting its directory deletes its state, copying it to another
project carries the record along, and per-agent installs each track themselves
without extra bookkeeping.

The `files` hashes are what make an update safe. Without a record of what the
CLI last wrote, it cannot tell its own output from an edit you made, and every
update becomes a guess.

`executable` records one boolean for each tracked file, alongside the existing
content hashes. Older locks without this field remain readable. On POSIX,
unknown base modes require attention before an update changes a file or a
removal proceeds; `--force` permits the change after backup. An unchanged
legacy installation can stay unchanged. An install/update that writes a lock
records complete mode metadata; no separate migration command is needed.

Only CLI-managed files belong in these maps; the lock itself and preserved
local-only files are excluded. A completely unchanged installation keeps its
commit and timestamp. Missing target locks still need to be written.

The root `.skill-lock.json` path is reserved for installation state; reject
source skills containing a file or directory there before writing any selected
skill. Files with that name in subdirectories are ordinary content and receive
the same copying and local-edit protection as other resources.

The full Git commit is recorded rather than an abbreviated hash. During install
or update, if a lock's `source` differs from the configured store, the CLI leaves
the skill untouched, reports that it came from another source, and exits as
needing attention. It never silently adopts the skill into the new store.
`--force` does not bypass source checks. Explicit removal may remove a foreign
installation with the normal conflict and backup checks.

## Commands

```
skills init                  create the config and clone the store
skills sync                  fast-forward the configured store branch
skills ls                    published skills in the store
skills ls --local            skills installed at the repository root
skills ls --global           skills installed in the home directories
skills install <skill>...    install at the repository root
skills remove <skill>...     uninstall from the selected scope
skills update [skill...]     re-install installed skills from the store
skills diff <skill>          show your edits to an installed skill
skills version               CLI version and current store commit

Flags:
  --global            act on the home directories instead of the repository
  --agent <name>...   narrow to certain agents (default: config defaults)
  --force             overwrite skills you have edited (backed up first)
  --dry-run           show the plan, write nothing
  -v                  show individual files and their targets
  --no-color          disable color (NO_COLOR is honoured too)
```

`--global` applies to `install`, `update`, `remove`, `diff` and installed `ls`.
`--agent` applies to `install`, `update`, `remove` and `diff`.
`--force` applies to `install`, `update` and `remove`. `ls --local` and
`ls --global` cannot be combined and list only the default agents' installs.

Running `skills` with no arguments prints this help. It never writes to disk by
surprise.

`ls` with no flag lists the store; `--local` and `--global` list what is
installed in the default agents' directories. Read descriptions from valid
local `SKILL.md` files without making valid frontmatter a discovery condition.

```
$ skills ls
  go-review        Reviews Go code for correctness and idiom     go, review
  django-testing   Writes and structures Django tests            python, django

$ skills ls --local
  go-review        Reviews Go code for correctness and idiom     agents, claude
```

`install` defaults to the detected repository root and writes a copy for every
selected agent. `--global` writes the configured global directories instead.
An identical unmanaged copy can be adopted by writing its lock. Different
unmanaged content needs `--force` and a backup.

`update` with no names refreshes all managed skills in the selected scope and
agents. It updates only targets where each skill is already installed; use
`install` to add a target. A skill installed for a non-default agent needs that
agent selected again to be updated, removed or diffed.

Repeated skill names are de-duplicated. An explicit name missing from the store
(`install`) or selected installations (`update`, `remove`, `diff`) is a
selection error: report all missing names and do no skill writes, exiting 1.
Invalid command syntax, flags and unknown agent names exit 2. A known installed
skill that became unavailable is attention, not a missing-name error.

`remove` acts only on the selected scope. It defaults to the repository and
requires `--global` to remove a global installation.

`sync` runs a fast-forward-only update of the configured branch. It never
creates a merge commit or resets a dirty or divergent store. It does not touch
installed skills, so syncing never surprises the user with a project rewrite.
`update` is the command that moves installed skills to the synced content.

`sync --dry-run` may query the remote branch head, but must not fetch, merge,
update refs, write objects or change the working tree. It reports the local
and remote heads. Equal heads mean up to date; otherwise report a proposed
sync. Without fetching missing history it may not be able to prove ancestry:
say that fast-forward eligibility will be checked during real sync. A remote
query failure exits 1 and must not be presented as up to date.

Store listings, installs and updates read tracked content from one captured
commit, which is also recorded in new locks. Uncommitted, untracked and ignored
working-tree files are not published content. They are never installed under
the identity of `HEAD`. Use the same committed-content representation for
install and diff; platform checkout conversion must not invent local edits.
`sync` still rejects a dirty store instead of resetting it.

`diff` compares each selected installation against its lock's recorded commit,
with the same agent transform, not the latest store head. It handles additions,
deletions and malformed local frontmatter as content; executable-bit changes
are shown on POSIX. No differences means no diff output and exit 0; differences
also exit 0. A missing recorded commit or foreign source produces an explanation
and exit 1, never a substitute base. No automatic source switching or history
recovery is required.

## How install and update decide

For every file of every skill being acted on, compare three file states
(content hash and, on POSIX, executable flag):

| | meaning |
| --- | --- |
| **want** | what the store holds for that agent, with store-only fields stripped |
| **have** | what is on disk now |
| **base** | what the CLI last wrote, from `.skill-lock.json` |

| condition | result |
| --- | --- |
| wanted by the store, not on disk | add |
| `have` equals `want` | nothing to do |
| `have` equals `base` | update, you never touched it |
| anything else | conflict, this is your edit |

The first row deliberately restores a deleted tracked file if the store still
contains it. For paths absent from `want`: remove an unchanged tracked file,
conflict on an edited tracked file, keep an untracked local-only file, and do
nothing if the path is also absent from disk. Empty directories carry no state.

This comparison must stay a pure function over three maps. It is the part that
can lose work, so it has to be testable without touching disk.

Agents add no special cases. `.claude/skills/go-review/SKILL.md` and
`.agents/skills/go-review/SKILL.md` are two ordinary paths in the same three
maps.

A skill is the unit of conflict checking across all selected agent targets.
Calculate its complete plan and validate filesystem shapes before writing.
Any known conflict skips the entire skill; other non-conflicting skills can
proceed. Stage replacement content for all its targets before publishing it.

This is not a transaction across filesystems. An unexpected I/O error during
publication or removal may leave partial progress. Stop further mutations,
exit 1 and report the affected skill, targets and every completed backup.
Never print a failed skill as successfully updated. Publish a target's new
lock only after its content writes and removals succeed. Keep completed earlier
results visible. Automatic rollback, crash-safe journals, concurrent writers
and protection against another process changing the tree during a command
are outside the first milestone; run one mutating command at a time.

If an installed skill disappears from the store or becomes a draft, `update`
reports it as unavailable and leaves it installed. Only an explicit `remove`
command uninstalls a skill. A pruning command can be added later if it becomes
useful.

## Skill metadata

One `SKILL.md` serves every agent. Its frontmatter may carry three things the
agents never see, and the CLI resolves them on install. All three are optional,
so a standard Agent Skills file installs without them.

```markdown
---
name: go-review
description: Reviews Go code for correctness and idiom. Use when...
status: published
tags: [go, review, testing]
x-claude:
  disable-model-invocation: true
  allowed-tools: Bash(go test *)
---

# Go review
...
```

| field | purpose |
| --- | --- |
| `status` | `published` or `draft`. Omitted means published. A draft is invisible to `ls` and cannot be installed. |
| `tags` | grouping for `ls` and, later, a TUI. Omitted or empty means no tags. |
| `x-claude` | frontmatter only Claude should see; dropped from the shared copy |

### Field kinds

Frontmatter holds three kinds of fields:

- **Standard fields** from the [Agent Skills specification](https://agentskills.io/specification):
  `name`, `description`, `license`, `compatibility`, `metadata` and the
  experimental `allowed-tools`. Author and version belong under `metadata`.
- **Claude fields** from the [Claude Code frontmatter reference](https://code.claude.com/docs/en/skills#frontmatter-reference),
  such as `when_to_use`, `model`, `effort`, `context`, `hooks`, `paths` and
  `disable-model-invocation`. At the top level they pass through to every
  agent. Put them under `x-claude` to keep them out of the shared copy.
- **Store fields**: `status`, `tags` and `x-claude`. Only this CLI reads them.

The CLI keeps the YAML value and type of every field it retains, but it does
not read, check or run those values. A kept field does not prove the target
agent supports it. Unknown fields pass through. There is no field allowlist and
no configurable transform. The CLI never moves a top-level field on its own.

Top-level `allowed-tools` is a standard field, so it stays in the shared copy.
The standard defines it as a space-separated string; Claude Code also accepts a
YAML list.

Claude Code accepts a skill without `name` or `description`. This CLI requires
both, so such a skill needs them added before it goes into a store. The CLI
checks what install needs; it is not a full Agent Skills validator and does not
promise to import every Claude skill unchanged.

### What install writes

| | `.claude/skills/` | `.agents/skills/` |
| --- | --- | --- |
| `status`, `tags` | stripped | stripped |
| `x-claude` | lifted to the top level | dropped |
| everything else | passed through | passed through |
| `agents/openai.yaml` | not copied | copied |
| other files | copied | copied |

Most skills have no `x-claude`, so their two installed copies differ only by the
absence of the store-only fields.

Validate source metadata before writing any selected skill: require nonempty
`name` and `description` and a directory-matching name. Only absence is a
default: an omitted `status` means `published`, but an explicit empty, null or
unsupported value is an error. If present, `tags` is a list of strings; null
and non-string elements are errors rather than coerced. If present, `x-claude`
is a mapping with unique string keys. Reject `name`, `description`, `status`,
`tags` and `x-claude` inside that mapping, and keys colliding with existing
top-level fields. Do not maintain a vendor-specific allowlist of every possible
extension key. Every mapping in the frontmatter, including nested values such
as `metadata` and `hooks`, must have unique keys. Emitted frontmatter must
remain valid with no duplicate keys.

Store frontmatter accepts LF and CRLF line endings; rewrites produce stable
frontmatter while preserving body bytes. Invalid source metadata is a runtime
error with its source path, not an unpublished skill silently omitted. The
catalog may fail on malformed store metadata; tolerance of edited installed
content does not relax source validation.

The frontmatter is parsed and validated with a YAML library. The transformation
drops `status` and `tags`, then either lifts or drops `x-claude`. The output is
serialised deterministically. This remains a pure function from bytes and an
agent name to bytes, so it tests without touching disk.

Expand YAML aliases in retained fields so removing or moving their anchors
cannot leave invalid frontmatter.

This transformation is why there is no per-agent copy in the store. The CLI has
to rewrite frontmatter anyway to strip `status` and `tags`; resolving
`x-claude` there is still cheaper than maintaining a duplicated body in every
skill forever.

### Drafts

`status: draft` keeps a work in progress unpublished without a release. A draft
can be pushed, reviewed and iterated on in the open, but it is invisible to `ls`
and cannot be installed until the field flips to `published` or is removed, and
the next `sync` picks it up. Omitting `status` publishes a valid skill, so work
in progress must say `status: draft` explicitly.

Registration is therefore a field, not a registry. There is no separate list of
known skills to keep in step with the directory, and no CLI release required to
publish one.

Content integrity needs nothing extra either: git already guarantees the store's
contents, and `.skill-lock.json` holds the per-file hashes that detect local
edits after install.

## Conflicts

A skill you edited is never overwritten silently. The tool skips the entire
skill, says so, and carries on with other skills. A command that leaves anything
needing attention exits with code 3, so scripts can distinguish partial
completion from success.

```
$ skills update
  ~ go-review        updated
  ! sql-review       skipped, you edited it

  2 skills, 1 changed, 1 needs attention
  Run 'skills diff sql-review' to see your changes,
  or 'skills update sql-review --force' to overwrite (backed up).
```

`--force` overwrites only after copying the old content into a timestamped
directory under the config root (`~/.config/skills/backups/` by default). The
CLI preserves enough of the target path to identify where the backup came from
and prints the exact backup path.

Force resolves detected content conflicts using the store version, after
backing up every existing selected target it will replace. Local-only regular
files remain user data and are preserved, including during a forced update.
An unchanged install/update needs no backup. Force is not a bypass for malformed locks,
foreign sources or unsupported filesystem shapes.

Backups never overwrite an earlier backup, including two operations in the
same second. Recovery hints include skill names and preserve `--global` and
explicit `--agent` selection; they never broaden a named request.

## Removals

No confirmation. Typing `skills remove go-review` is the confirmation, but it
acts only on the selected scope. An edited skill is a conflict unless `--force`
is present. Every successful removal is backed up first, and the exact backup
path is printed.

Added, modified or deleted tracked content, added local-only files and changed
executable bits on POSIX count as edits. A valid lock is sufficient to attempt
removal even if `SKILL.md` is missing or malformed. Back up all selected target
directories for the skill before removing any of them. A backup failure leaves
that skill installed; successful removals earlier in the command still appear
in the output. Unsupported entries need manual resolution even with force.

## Output

One line per skill, not per file. A skill installed for two agents should not
print two lines.

```
$ skills sync
  pulled store  a1b2c3d -> e5f6a7b

$ skills update
  ~ go-review        updated
  ~ sql-review       updated
  ! django-testing   unavailable in store, left installed

  4 skills, 2 changed, 1 needs attention
```

`-v` adds the file paths under each skill, which is where the agents become
visible:

```
$ skills update -v
  ~ go-review        updated
      .claude/skills/go-review/SKILL.md
      .agents/skills/go-review/SKILL.md
```

Unchanged skills are silent but included in the summary count.
Verbose output lists changed content paths, conflicts or paths being removed
under each target. Removal includes the lock file because that file is removed
too. Unchanged and preserved local-only paths are not install/update changes.

Dry-run action lines use `would add`, `would update` and `would remove`; its
summary uses `N would change`. Counts are per skill across targets. Conflict
and unavailable results still need attention and exit 3. Dry-run never claims
a backup was written; where required, it says a backup would precede the change.
The real backup path is printed after creation, including when a later step
fails. Result output goes to stdout; warnings and runtime diagnostics go to
stderr. A runtime error takes precedence over attention in the exit code.

Symbols carry the meaning, so nothing is lost when color is off:

```
+  added      ~  updated      -  removed      !  needs attention
```

Color is used only when stdout is a terminal, and `NO_COLOR` disables it.

There is no `--json` yet. Nothing consumes it, and the output shape is still
settling.

Exit codes are part of the command contract:

| code | meaning |
| --- | --- |
| `0` | every requested operation completed |
| `1` | runtime or missing-name selection error |
| `2` | invalid command or arguments |
| `3` | attention required; a skill was skipped or installed state needs repair |

## Releases

Only the CLI is released. A skill is published by pushing it to the store
without `status: draft`; it reaches a project on the next `skills sync` plus
`skills update`. Adding or editing a skill never needs a CLI release.

`main` is protected and work lands through a pull request. Merging a published
skill makes that skill available to the store immediately; a CLI release still
happens only when a tag is pushed.

```
skill: branch -> pull request -> merge to main -> available after sync
CLI:   branch -> pull request -> merge to main -> tag v0.3.0 -> release
```

### Branch protection

Repository configuration, not a file: require a pull request before merging,
require the CI checks to pass, block force pushes and branch deletion on `main`.

### Two workflows

**CI**, on every pull request: `gofmt -l .`, `go mod tidy -diff`,
`go vet ./...`, `golangci-lint run` and `go test -race ./...`.

**Release**, on `v*` tags: build with GoReleaser and publish binaries for macOS,
Linux and Windows on amd64 and arm64, plus checksums and generated notes.

### Version injection

`debug.ReadBuildInfo()` reports the right version for
`go install <module>@v0.3.0`, but `(devel)` or a pseudo-version for a binary
built from a checkout. The release build sets it explicitly:

```
-ldflags "-X main.version={{.Version}}"
```

`main.go` is at the repository root, so the symbol is `main.version` with no
package path in front of it.

The CLI version and the store commit are separate facts, and `skills version`
prints both.

## Agent differences

Authoring reference for the tools we install into, rather than additional CLI
validation requirements. Vendor limits and invocation syntax can change;
consult the linked vendor docs when authoring. The CLI's validation and
transformation contract is defined in [Skill metadata](#skill-metadata).

### Behaviour

| | Claude Code | Codex | Gemini CLI |
| --- | --- | --- | --- |
| explicit call | `/skill-name` | `$skill` in the CLI, `@skill` in ChatGPT | per its own docs |
| command identity | the **directory name**; `name` is display only | the `name` field | the `name` field |
| name rules | 64 chars max, lowercase letters, digits and hyphens; cannot contain "claude" or "anthropic"; cannot be `synced` | not documented | not documented |
| description limit | 1,024 characters in the spec; Claude Code truncates `description` plus `when_to_use` at 1,536 | the whole skill list is capped at 2% of the context window, or 8,000 characters; descriptions are shortened when there are many | not documented |
| body budget | keep under 500 lines | not specified | not specified |
| disable auto-trigger | `disable-model-invocation: true` in frontmatter | `policy.allow_implicit_invocation: false` in `agents/openai.yaml` | not documented |
| extra configuration | around twenty frontmatter fields | a separate `agents/openai.yaml` file | extension configuration |

### The extension problem

The tools extend the format in opposite directions, and only one of those is
safe to use here.

Codex keeps its extras in a sidecar file, `agents/openai.yaml`, holding display
name, icons, invocation policy and MCP tool dependencies. The other tools never
read it, so carrying it costs nothing.

Claude Code keeps its extras in the frontmatter itself: `model`, `effort`,
`context`, `hooks`, `paths`, `disable-model-invocation` and others. Codex and
Gemini do not say what they do with keys they do not know.

The authoring rule:

- Standard fields, including `allowed-tools`, stay at the top level.
- Claude fields at the top level are copied to every agent. Put them under
  `x-claude` to keep them out of the shared copy. See
  [Field kinds](#field-kinds).
- `status` and `tags` are store-only and are stripped on install.
- Codex-specific settings go in `agents/openai.yaml`, which the others ignore.

Because Claude Code takes the command name from the directory while the others
take it from the `name` field, the two must be identical or one skill answers to
two names.

### Authoring guidance

All three vendors agree on the main points:

- The description decides whether the skill fires. Say what it does and when to
  use it, lead with the trigger words, and write in the third person.
- Be concise. The body shares the context window with everything else.
- Keep the body under 500 lines; split into sibling files past that.
- Keep file references one level deep from `SKILL.md`. Nested references get
  partially read.
- Give reference files longer than 100 lines a table of contents.
- Use forward slashes in paths, and avoid dates and version notes that age.
- Use one term for one thing throughout a skill.
- Keep each skill focused on one job.

Anthropic and OpenAI differ slightly on scripts: Anthropic recommends bundling
utility scripts, since they are more reliable than generated code and cost no
context until run, while OpenAI prefers instructions unless the task needs
determinism or an external tool.

Sources: [Skill authoring best practices](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices),
[Claude Code skills reference](https://code.claude.com/docs/en/skills),
[Build skills](https://learn.chatgpt.com/docs/build-skills),
[Gemini CLI Agent Skills](https://geminicli.com/docs/cli/skills/).

## Decisions

1. **Config format.** Use TOML. This is a human-edited file, so readability and
   comments justify a small dependency.
2. **Store access.** Shell out to `git`, which reuses existing credentials and
   supports private stores without adding HTTP, TLS or authentication code.
   Version one tracks one branch and syncs it by fast-forward only.
3. **Lock file name.** Use `.skill-lock.json` so the CLI state stays out of the
   way inside the skill directory.
4. **Exit codes.** Use `0` for complete success, `1` for runtime or selection
   errors, `2` for usage errors and `3` when installed state needs attention.
5. **Backups.** Store timestamped backups under the config root
   (`~/.config/skills/backups/` by default) and always print the exact path.
6. **Binary name collision.** Keep the binary name `skills`. `go install`,
   `go run .`, `go vet ./...`, `go test ./...` and `go build ./...` work; only
   bare `go build .` fails, which does not justify a rename.
7. **Validation.** Listing, installation and updates validate skills internally.
   A public `skills check` command is deferred until an authoring or CI workflow
   needs it.
8. **Store and CLI repository.** Keep them together. The store location remains
   configurable, so they can be split later without changing the install model.

## Build order

Each step is testable before the next one depends on it.

1. Config and scope: load TOML, apply defaults, create it on first use, detect
   the repository root and warn when falling back to the current directory
2. Store: clone the configured branch, fast-forward it safely and report the
   full current commit
3. Frontmatter: parse and validate YAML; read `name`, `description`, `status`
   and `tags`; strip store-only fields and resolve `x-claude`, with tests
4. `ls` for the store
5. The three-way comparison and whole-skill plan, with tests
6. `install` and `remove` for one agent, writing `.skill-lock.json` and backing
   up destructive changes
7. Scanning installed skills: `ls --local`, `ls --global`, `update`
8. `diff`, dry-run, output formatting, color and exit codes
