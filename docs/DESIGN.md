# Design

How the `skills` CLI works and why. This document is the plan; the code follows
it, not the other way round.

## Goal

The repository is a centralized store of agent skills. The CLI exists to get a
skill out of that store and into any project quickly, and to keep installed
skills current as the store changes.

Skills must work with Claude Code, Codex and Gemini CLI.

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

Any command that needs the store initialises it first:

1. create the config directory and write a default config
2. clone the store
3. carry on with the command

`skills init` does the same thing explicitly. There is nothing to set up by
hand, but the first run needs network access.

`--dry-run` is the exception: it never creates the config or clones the store.
If the command cannot calculate a plan because the store is not initialised, it
exits with an explanation and tells the user to run `skills init` first.

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

### Store layout

One directory per skill, one `SKILL.md` inside it:

```
skills/
  go-review/
    SKILL.md                shared body and frontmatter
    agents/openai.yaml      optional, copied only to .agents/
    references/             optional, copied to both
```

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

[agents.gemini]
project = ".agents/skills"
global  = "~/.agents/skills"

[defaults]
agents = ["claude", "codex"]
```

The agent map is the reason this is a file rather than a flag. Supporting a new
tool is an edit here, not a release.

Codex and Gemini share `.agents/skills`, so two agents resolve to one directory.
Installs are de-duplicated by resolved path: asking for both writes the files
once.

## Agents and targets

| agent | project | global |
| --- | --- | --- |
| Claude Code | `.claude/skills/` | `~/.claude/skills/` |
| Codex | `.agents/skills/` | `~/.agents/skills/` |
| Gemini CLI | `.agents/skills/` (alias, preferred over `.gemini/skills/`) | `~/.agents/skills/` (alias) |

All three use the same skill format: a directory holding a `SKILL.md` whose
frontmatter carries `name` and `description`, with optional `scripts/`,
`references/` and `assets/` subdirectories. All three load only the name and
description up front, read the body when the skill triggers, and read other
files on demand.

A plain skill is therefore portable as written. See
[Agent differences](#agent-differences) for where they diverge.

Skills are copied into each target as real files, never symlinked. Symlinks need
a git config flag to survive a clone, behave badly on Windows, and break when
one of the directories is gitignored.

## Scope

Project scope means the root of the current Git repository, regardless of which
subdirectory the command is run from. If the current directory is not inside a
Git repository, the CLI warns the user and uses the current directory as the
project root.

`--global` uses the user's home directories instead. Nothing else changes: the
same skill, the same lock, a different parent path.

Scope is always explicit. Commands default to project scope and never fall back
to global scope. If a requested skill exists only globally, the error says so
and suggests re-running the command with `--global`.

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
  "files": { "SKILL.md": "<sha256>" }
}
```

There is no central index. Installed skills are found by scanning the configured
agent directories for `.skill-lock.json`.

This is what makes the tool stateless about projects. A skill carries its own
provenance, deleting its directory deletes its state, copying it to another
project carries the record along, and per-agent installs each track themselves
without extra bookkeeping.

The `files` hashes are what make an update safe. Without a record of what the
CLI last wrote, it cannot tell its own output from an edit you made, and every
update becomes a guess.

The full Git commit is recorded rather than an abbreviated hash. During install
or update, if a lock's `source` differs from the configured store, the CLI leaves
the skill untouched, reports that it came from another source, and exits as
needing attention. It never silently adopts the skill into the new store.

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

Running `skills` with no arguments prints this help. It never writes to disk by
surprise.

`ls` with no flag lists the store; `--local` and `--global` list what is
installed, found by scanning the agent directories for `.skill-lock.json` and
reading the name and description from each `SKILL.md`.

```
$ skills ls
  go-review        Reviews Go code for correctness and idiom     go, review
  django-testing   Writes and structures Django tests            python, django

$ skills ls --local
  go-review        Reviews Go code for correctness and idiom     claude, agents
```

`install` defaults to the detected repository root and writes a copy for every
configured agent. `--global` writes the home directories instead.

`remove` acts only on the selected scope. It defaults to the repository and
requires `--global` to remove a global installation.

`sync` runs a fast-forward-only update of the configured branch. It never
creates a merge commit or resets a dirty or divergent store. It does not touch
installed skills, so syncing never surprises the user with a project rewrite.
`update` is the command that moves installed skills to the synced content.

## How install and update decide

For every file of every skill being acted on, compare three hashes:

| | meaning |
| --- | --- |
| **want** | what the store holds for that agent, with store-only fields stripped |
| **have** | what is on disk now |
| **base** | what the CLI last wrote, from `.skill-lock.json` |

| condition | result |
| --- | --- |
| not on disk | add |
| `have` equals `want` | nothing to do |
| `have` equals `base` | update, you never touched it |
| anything else | conflict, this is your edit |

This comparison must stay a pure function over three maps. It is the part that
can lose work, so it has to be testable without touching disk.

Agents add no special cases. `.claude/skills/go-review/SKILL.md` and
`.agents/skills/go-review/SKILL.md` are two ordinary paths in the same three
maps.

A skill is the unit of change across all selected agent targets. The CLI
calculates the whole plan before writing anything. If any file conflicts, it
skips the entire skill; it never leaves one agent or one file partially updated.
Other non-conflicting skills in the same command can still proceed.

If an installed skill disappears from the store or becomes a draft, `update`
reports it as unavailable and leaves it installed. Only an explicit `remove`
command uninstalls a skill. A pruning command can be added later if it becomes
useful.

## Skill metadata

One `SKILL.md` serves every agent. Its frontmatter carries three things the
agents never see, and the CLI resolves them on install.

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
| `status` | `published` or `draft`. A draft is invisible to `ls` and cannot be installed. |
| `tags` | grouping for `ls` and, later, a TUI |
| `x-claude` | frontmatter only Claude should see |

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

The frontmatter is parsed and validated with a YAML library. The transformation
drops `status` and `tags`, then either lifts or drops `x-claude`. The output is
serialised deterministically. This remains a pure function from bytes and an
agent name to bytes, so it tests without touching disk.

This transformation is why there is no per-agent copy in the store. The CLI has
to rewrite frontmatter anyway to strip `status` and `tags`; resolving
`x-claude` there is still cheaper than maintaining a duplicated body in every
skill forever.

### Drafts

`status: draft` keeps a work in progress unpublished without a release. A draft
can be pushed, reviewed and iterated on in the open, but it is invisible to `ls`
and cannot be installed until the field flips to `published` and the next
`sync` picks it up.

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
  + go-review        added
  ! sql-review       skipped, you edited it

  1 skill needs attention
  Run 'skills diff sql-review' to see your changes,
  or 'skills update --force' to overwrite (backed up).
```

`--force` overwrites only after copying the old content into a timestamped
directory under the config root (`~/.config/skills/backups/` by default). The
CLI preserves enough of the target path to identify where the backup came from
and prints the exact backup path.

## Removals

No confirmation. Typing `skills remove go-review` is the confirmation, but it
acts only on the selected scope. An edited skill is a conflict unless `--force`
is present. Every successful removal is backed up first, and the exact backup
path is printed.

## Output

One line per skill, not per file. A skill installed for two agents should not
print two lines.

```
$ skills sync
  pulled store  a1b2c3d -> e5f6a7b

$ skills update
  + go-review        added
  ~ sql-review       updated
  ! django-testing   unavailable in store, left installed

  4 skills, 2 changed, 1 needs attention
```

`-v` adds the file paths under each skill, which is where the agents become
visible:

```
$ skills update -v
  + go-review        added
      .claude/skills/go-review/SKILL.md
      .agents/skills/go-review/SKILL.md
```

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
| `1` | runtime error |
| `2` | invalid command or arguments |
| `3` | attention required; one or more requested skills were skipped |

## Releases

Only the CLI is released. A skill is published by setting `status: published`
and pushing to the store; it reaches a project on the next `skills sync` plus
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

**CI**, on every pull request: `gofmt -l .`, `go vet ./...`, `go test ./...`.

**Release**, on `v*` tags: build with GoReleaser and publish binaries for macOS,
Linux and Windows on amd64 and arm64, plus checksums and generated notes.

### Version injection

`debug.ReadBuildInfo()` reports the right version for
`go install <module>@v0.3.0`, but `(devel)` for a binary built in CI from a
checkout. The release build sets it explicitly:

```
-ldflags "-X main.version={{.Version}}"
```

`main.go` is at the repository root, so the symbol is `main.version` with no
package path in front of it.

The CLI version and the store commit are separate facts, and `skills version`
prints both.

## Agent differences

Reference for the tools we install into. The shared format is described in
[Agents and targets](#agents-and-targets); this section records the divergences.

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

Claude Code keeps its extras in the frontmatter itself: `allowed-tools`,
`model`, `effort`, `context`, `hooks`, `paths`, `disable-model-invocation` and
others. Codex and Gemini document only `name` and `description`, and neither
says what it does with unknown keys.

The authoring rule:

- Claude-only frontmatter goes under `x-claude`, which is dropped for the other
  tools. See [Skill metadata](#skill-metadata).
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
4. **Exit codes.** Use `0` for complete success, `1` for runtime errors, `2` for
   usage errors and `3` when requested skills are skipped and need attention.
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
