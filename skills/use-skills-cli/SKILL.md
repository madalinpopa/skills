---
name: use-skills-cli
description: Installs a missing skill into the current project with the skills CLI so the agent can use it in the same session. Use whenever a task, a workflow skill, or the user names a skill that is not available to the current agent, when a skill invocation fails because the skill is unknown, or when the user asks to install, update, or refresh a skill from the skills store. Also use when a project's installed skill is behind the store and needs an update. Not for authoring or editing skills, and not for installing anything other than skills from the configured store.
compatibility: Requires the skills CLI and Git on PATH. Uses the store configured in ~/.config/skills/config.toml.
status: published
tags: [skills, cli, install, workflow]
---

# Use the skills CLI

Use this skill when a skill you need is not available. It installs the skill
into the current project with the `skills` CLI and makes it usable without
ending the session. Run every command from the project root. Never add
`--force`, `--global`, or edit files under the install targets by hand.

## 1. Check the CLI

Run `skills version`. If the command is missing, stop and tell the user how
to install it:

```sh
brew install --cask madalinpopa/tap/skills   # macOS
go install github.com/madalinpopa/skills@latest
```

Do not install the CLI yourself. The first CLI command creates the config and
clones the store when they do not exist yet.

## 2. Pick the agent and target

Select the current agent name: `claude` for Claude Code, `codex` for Codex and
Gemini CLI. The project target for each agent comes from the config; the
defaults are:

| Agent | Project target |
| --- | --- |
| `claude` | `.claude/skills/<skill>/` |
| `codex` | `.agents/skills/<skill>/` |

The skill is installed when `<target>/<skill>/.skill-lock.json` exists.
`skills ls --local` lists installed skills with their agents.

## 3. Sync the store and confirm the skill exists

Run `skills sync`. It fast-forwards the local store and never touches
installed skills. If it fails because the store is dirty or divergent, report
the message and stop; do not reset the store.

Run `skills ls` and check that the skill name is listed. A missing name means
the skill is not published in the store, for example because it is still
`status: draft`. Report this and stop. Do not copy the skill from a checkout
by hand.

## 4. Install or update

Replace `<skill>` and `<agent>`:

```sh
skills install <skill> --agent <agent>   # lock is absent
skills update <skill> --agent <agent>    # lock exists
```

Read the result line and the exit code:

| Exit | Meaning | What to do |
| --- | --- | --- |
| `0` | done, `+ added`, `~ updated`, or nothing to change | continue |
| `1` | runtime error or unknown name | report the message and stop |
| `2` | invalid command or agent name | fix the command and retry once |
| `3` | `!` needs attention: local edits, conflict, or unavailable | report and stop |

Exit `3` means someone edited the installed skill or the store version is
unavailable. The CLI keeps the local version. Show the user the printed hint,
such as `skills diff <skill> --agent <agent>`, and let them decide. Never
resolve it with `--force` on your own.

## 5. Make the skill usable

Agents load skill names at session start and may not see a new directory
right away. Try the skill first; if it is unknown, use the reload path for
the current agent:

- Claude Code: edits to an installed `SKILL.md` are picked up automatically.
  A new skill directory may need `/reload-plugins`, which the user types in
  the session. Ask for it only when the skill is still unknown after install.
- Codex and Gemini CLI: new skills are detected automatically. If the skill
  is still unknown, the user restarts the session.

When a reload is not possible right now, read the installed `SKILL.md` at
`<target>/<skill>/SKILL.md` and follow it directly, loading its references
on demand. Do not claim a skill ran when it was neither invoked nor read.

## 6. Report

Say which skill was installed or updated, for which agent, at which path,
and whether it was invoked or followed by reading its file. Whether the
installed files are committed or ignored is the user's choice; do not change
`.gitignore` or commit them.
