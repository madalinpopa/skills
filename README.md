# Skills & Technical Knowledge

My personal library of AI skills, prompts, and custom instructions, plus a small
Go CLI that installs them into any project and keeps them up to date.

The skills work with Claude Code, Codex and Gemini CLI.

## Install the CLI

```sh
go install github.com/madalinpopa/skills@latest
```

The first command you run creates `~/.config/skills/` and clones this
repository into it. That local clone is the store every project installs from,
so one copy serves the whole machine.

## Use it

```sh
skills ls                   # skills available in the store
skills install go-review    # install into this project
skills ls --local           # what this project has installed
```

Skills land in the directory each agent reads:

```
.claude/skills/go-review/   Claude Code
.agents/skills/go-review/   Codex and Gemini CLI
```

Use `--global` to install into your home directory instead, so a skill is
available in every project:

```sh
skills install go-review --global
skills ls --global
```

To pick up newer skills:

```sh
skills sync      # git pull the store
skills update    # move installed skills to the pulled content
```

`sync` only pulls, so it never rewrites a file unexpectedly. `update` is the
step that touches your project.

## Commands

| command | what it does |
| --- | --- |
| `skills init` | create the config and clone the store |
| `skills sync` | git pull the store |
| `skills ls` | skills available in the store |
| `skills ls --local` | skills installed in this project |
| `skills ls --global` | skills installed in your home directory |
| `skills install <skill>` | install into this project |
| `skills remove <skill>` | uninstall, project first then global |
| `skills update` | move installed skills to the pulled content |
| `skills diff <skill>` | show your edits to an installed skill |
| `skills version` | CLI version and current store commit |

## Your edits are never lost

Each installed skill carries a `skill.lock.json` recording what the CLI wrote.
An update compares three things per file: what the store holds, what is on disk,
and what the CLI last wrote.

If the file on disk still matches what the CLI wrote, you never touched it, so
it is safe to overwrite. Anything else is your own edit. The update reports it
and leaves the file alone:

```
$ skills update
  + go-review        added
  ! sql-review       skipped, you edited it

  1 skill needs attention
```

`skills diff sql-review` shows what you changed. `skills update --force`
overwrites, after backing up the old content first.

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

`status` and `tags` are for the store and are stripped on install, so the agents
only ever see fields they understand. `status: draft` keeps a work in progress
out of `ls` and out of everyone's projects until it is ready — no release
needed, just push it.

Claude-only frontmatter goes under an `x-claude` key, which is lifted into place
for Claude Code and dropped for the others. Most skills need none of it.

The `description` is what decides whether a skill fires, so say both what it
does and when to use it.

## Design

[docs/DESIGN.md](docs/DESIGN.md) covers the whole design and the reasoning: the
store and config layout, how the three-way comparison protects local edits, what
each agent supports and where they differ, and how releases work.
