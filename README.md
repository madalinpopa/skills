# Skills & Technical Knowledge

A personal library of reusable AI skills, prompts, and instructions, plus a
Go CLI that installs skills into your projects and keeps them up to date.
Supports Claude Code, Codex, and Gemini CLI.

Browse the [skill catalog](skills/), explore the [demo skill](skills/demo/SKILL.md),
or use the [documentation templates](docs/templates/) for your own projects.

## Install the CLI

```sh
go install github.com/madalinpopa/skills@latest
```

Or download a prebuilt binary for macOS, Linux, or Windows from
[Releases](https://github.com/madalinpopa/skills/releases). The CLI needs
[Git](https://git-scm.com/) 2.45 or newer on your `PATH`.

## Quick start

Run these commands from a project:

```sh
skills init                 # create the config and clone the skill store
skills ls                   # list published skills
skills install demo         # install a skill into this project
skills ls --local           # list this project's installed skills
skills sync                 # fetch the latest store content
skills update               # update installed skills from the store
skills remove demo          # uninstall a skill
```

The store is a local clone shared across projects. `sync` refreshes that clone;
`update` applies its content to installed skills. First-time setup also happens
automatically when a command needs the store.

Local edits are protected: changed skills are skipped during an update. Use
`skills diff <skill>` to inspect your edits. Add `--force` to an install, update,
or remove command to replace or remove edited content after a backup.

Use `--dry-run` to preview changes without writing, `-v` to show individual files,
and `skills <command> --help` for all options. Exit code `3` means a skill needs
attention; `skills version` shows the CLI version and current store commit.

## Configuration and scope

First run creates `~/.config/skills/config.toml` and a store at
`~/.config/skills/store/`. `XDG_CONFIG_HOME` is honored when set. Configuration
and the store are shared across projects on the machine.

By default, skills are installed for both configured agents:

| Agent | Project location | Global location |
| --- | --- | --- |
| Claude Code | `.claude/skills/` | `~/.claude/skills/` |
| Codex | `.agents/skills/` | `~/.agents/skills/` |

Gemini CLI also reads `.agents/skills/`, so `--agent codex` serves it too.
Use `--agent` to select an agent and `--global` to install for all your projects:

```sh
skills install demo --agent codex
skills install demo --global
skills ls --global
```

Project scope uses the repository root, or the current directory outside a
repository. Edit the config to change agent paths or defaults; see the
[default configuration](docs/SPEC.md#configuration) for the complete format.

## Make it your own

Fork and clone this repository to maintain your own skill library. The CLI
continues using its configured store; forking alone does not change the source.
To install from your fork, edit the `[store]` section of your CLI config:

```toml
[store]
repo = "https://github.com/your-name/skills"
branch = "main"
```

If you already initialized a store from another source, preserve any local work
and move the old store directory aside. Run `skills init` to clone the configured
fork, then install the skills you want. Existing installations retain their
original source; changing the config does not migrate them to the fork.

To add or update a skill, use [create-repo-skill](skills/create-repo-skill/SKILL.md).
It guides naming, metadata, packaging, and validation in this repository and
its forks. Skill sources live under `skills/`; [AGENTS.md](AGENTS.md#creating-and-updating-skills)
defines how an agent installs and uses the authoring skill.

Publish a skill by setting `status: published` and committing and pushing it to
your configured store branch. Then run `skills sync` and install or update it.
Skills marked `draft` cannot be installed. Skill content changes need no CLI release.

## Development

Use [Task](https://taskfile.dev) for local checks:

```sh
task test:unit           # tests without network or Docker
task test:integration    # integration tests; requires Docker
task lint               # lint and dependency checks
task format             # format Go code
```

See [Taskfile.yml](Taskfile.yml) for all tasks and
[docs/SPEC.md](docs/SPEC.md) for the CLI's full behavior contract.
