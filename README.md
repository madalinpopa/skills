# Skills & Technical Knowledge

My personal library of AI skills, prompts, and custom instructions, plus a small
Go CLI that installs them into any project and keeps them up to date.

The skills work with Claude Code, Codex and Gemini CLI.

## Install the CLI

```sh
go install github.com/madalinpopa/skills@latest
```

Prebuilt binaries for macOS, Linux and Windows on amd64 and arm64 are attached
to every [release](https://github.com/madalinpopa/skills/releases), with a
checksum file next to them.

## First run

The first command that needs the store creates `~/.config/skills/config.toml`
and clones this repository into `~/.config/skills/store/`. That local clone is
the store every project installs from, so one copy serves the whole machine.

```sh
skills init      # do the first run explicitly
```

`XDG_CONFIG_HOME` is honoured when set. `--dry-run` never creates the config or
clones the store. If the store is missing, it tells you to run `skills init`.

## Use it

```sh
skills ls                   # published skills in the store
skills install go-review    # install into this project
skills ls --local           # what this project has installed
```

Skills land in the directory each agent reads:

```
.claude/skills/go-review/   Claude Code
.agents/skills/go-review/   Codex and Gemini CLI
```

Codex and Gemini share one directory, so their copy is written once. The
default agents are `claude` and `codex`. Use `--agent` to narrow a command to
some of them:

```sh
skills install go-review --agent claude
```

To pick up newer skills:

```sh
skills sync      # fast-forward the store
skills update    # move installed skills to the synced content
```

`sync` only moves the store forward, so it never rewrites a file in your
project. `update` is the step that touches your project.

`sync --dry-run` asks the remote for its branch head and prints the local and
remote commits. It does not fetch, so it cannot promise a fast-forward; the
real `sync` checks that.

Listing, install, update and diff read skills from the store's current commit,
which is the commit recorded in each lock. Uncommitted edits, untracked
directories and ignored files in the store clone are never installed.

## Scope

Commands act on the current Git repository, from any subdirectory. Outside a
repository they warn and use the current directory.

`--global` acts on your home directory instead, so a skill is available in
every project:

```sh
skills install go-review --global
skills ls --global
skills remove go-review --global
```

Scope is explicit. A command never falls back from the project to your home
directory. When `update`, `diff` or `remove` names a skill that is only
installed globally for the selected agents, the error tells you to add
`--global`. Only the missing names are checked, so a skill installed in the
project is never reported as global only.

## Commands

| command | what it does |
| --- | --- |
| `skills init` | create the config and clone the store |
| `skills sync` | fast-forward the configured store branch |
| `skills ls` | published skills in the store |
| `skills ls --local` | skills installed in this project |
| `skills ls --global` | skills installed in your home directory |
| `skills install <skill>...` | install into this project |
| `skills remove <skill>...` | uninstall from the selected scope |
| `skills update [skill...]` | re-install installed skills from the store |
| `skills diff <skill>` | show your edits to an installed skill |
| `skills version` | CLI version and current store commit |

Flags:

```
--global      act on the home directory instead of the repository
--agent       narrow to certain agents (default: config defaults)
--force       overwrite or remove skills you have edited (backed up first)
--dry-run     show the plan, write nothing
-v            show individual files and their targets
--no-color    disable color (NO_COLOR is honoured too)
```

## Your edits are never lost

Each installed skill carries a `.skill-lock.json` recording what the CLI wrote.
Every target directory gets its own lock. If a copy already sits in a target
with no lock and matches the store, `install` adopts it by writing the lock and
leaves the files as they are.
An update compares three things per file: what the store holds, what is on disk,
and what the CLI last wrote. On Linux and macOS the executable bit is part of
that comparison, so bundled scripts run after install and a local `chmod`
counts as an edit. A lock written before executable bits were recorded needs
`--force` the first time an update or removal would change that skill.

If the file on disk still matches what the CLI wrote, you never touched it, so
it is safe to overwrite. Anything else is your own edit. The update skips the
whole skill, says so, and carries on with the others:

```
$ skills update
  ~ go-review        updated
  ! sql-review       skipped, you edited it

  2 skills, 1 changed, 1 needs attention
  Run 'skills diff sql-review' to see your changes,
  or 'skills update sql-review --force' to overwrite (backed up).
```

`skills diff sql-review` shows your changes as a unified diff against the
content the CLI installed. A changed executable bit is shown as an old and
new mode line, the way Git shows it. `skills update --force` overwrites, after copying
the old content into a timestamped directory under `~/.config/skills/backups/`.
The exact backup path is printed. `skills install --force` does the same, and
also replaces a skill directory that `skills` did not install. Hints repeat the
skill names, `--global` and `--agent` you passed, so they never widen the
request.

`skills remove` also backs up before deleting, and refuses an edited skill
unless `--force` is present.

If a write or removal fails part way, the command stops, prints the skills it
completed and every backup path it created, and exits 1. The failed skill is
marked and never counted as changed. There is no automatic rollback.

A skill that disappears from the store, or becomes a draft, is reported as
unavailable and left installed. A skill installed from another store is
reported and left alone.

Symlinks and special files inside an installed skill are not supported. The
CLI reports the path and reason, leaves the whole skill untouched even with
`--force`, and never follows a link to read or write outside the skill.

A skill is found by its `.skill-lock.json`, not by its frontmatter, so a
broken `SKILL.md` can still be diffed, updated or removed. A damaged lock is
reported and never rewritten, even with `--force`. `ls --local` and
`ls --global` mark such rows and exit 3.

`--dry-run` prints the same plan without writing anything:

```sh
skills update --dry-run
```

Its lines say `would add`, `would update` or `would remove`, and the summary
counts skills that would change. Where a change would be backed up first it
says so, without creating a backup. Skipped skills still exit 3.

## Exit codes

| code | meaning |
| --- | --- |
| `0` | every requested operation completed |
| `1` | runtime error |
| `2` | invalid command or arguments |
| `3` | attention required; one or more requested skills were skipped |

## Writing a skill

One directory per skill, one `SKILL.md` inside it:

```
skills/
  go-review/
    SKILL.md                shared by every agent
    agents/openai.yaml      optional, Codex only
    references/             optional
```

```markdown
---
name: go-review
description: Reviews Go code for correctness and idiom. Use when reviewing a Go
  diff, a pull request touching .go files, or Go code quality.
status: published
tags: [go, review]
---

# Go review
...
```

`name` must match the directory name. `status` is required and is `published`
or `draft`. `status` and `tags` are for the store and are stripped on install,
so the agents only ever see fields they understand. `status: draft` keeps a
work in progress out of `ls` and out of everyone's projects until it is ready.
No release is needed, just push it.

Claude-only frontmatter goes under an `x-claude` key, which is lifted into place
for Claude Code and dropped for the others. Most skills need none of it. Its
keys must be unique strings, must not repeat a top-level field, and must not
be `name`, `description`, `status`, `tags` or `x-claude`. A skill that breaks
this is a store error for every agent, not just Claude. The frontmatter may
use LF or CRLF line endings; the installed copy always gets LF delimiters and
keeps the body bytes as written.

The `description` is what decides whether a skill fires, so say both what it
does and when to use it.

## Development

```sh
go test ./...        # the default suite needs neither network nor Docker
golangci-lint run    # the lint config lives in .golangci.yml
```

CI runs formatting, vet, lint and tests on every pull request. Pushing a `v*`
tag builds and publishes the release binaries with GoReleaser.

## Design

[docs/DESIGN.md](docs/DESIGN.md) covers the whole design and the reasoning: the
store and config layout, how the three-way comparison protects local edits, what
each agent supports and where they differ, and how releases work.
