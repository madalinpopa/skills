# Skills & Technical Knowledge

A personal library of reusable AI skills, prompts, and instructions, plus a
Go CLI that installs skills into your projects and keeps them up to date.
Supports Claude Code, Codex, and Gemini CLI.

Browse my [skill catalog](skills/). These skills are personal and shaped around
my daily workflow and tools. Take ideas from them, and use any that work for you.

## Install the CLI

```sh
go install github.com/madalinpopa/skills@latest
```

Or download a binary for macOS, Linux, or Windows from
[Releases](https://github.com/madalinpopa/skills/releases). You also need
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

The store is a local clone of your skills repository, shared by all your
projects. `sync` pulls new changes into it. `update` copies those changes into
your installed skills. You don't have to run `init` first: the CLI sets things
up the first time a command needs the store.

Your local edits are safe. If you changed an installed skill, `update` skips it.
Run `skills diff <skill>` to see what you changed. Add `--force` to `install`,
`update`, or `remove` to replace your edits; the CLI backs them up first.

Add `--dry-run` to see what would change without writing anything, and `-v` to
list each file. Run `skills <command> --help` for all options. Exit code `3`
means a skill needs your attention. `skills version` shows the CLI version and
the store commit.

## Examples

The output below is an example. Use a skill name from `skills ls` instead of
`demo`, and `claude` instead of `codex` if that is your agent. `$` marks the
command you type.

Install a skill for one agent:

```console
$ skills install demo --agent codex
  + demo   added

  1 skill, 1 changed
```

Pull new changes, then update a skill you have not edited. This assumes the
store has a newer version of `demo`:

```console
$ skills sync
  pulled store  a1b2c3d -> e4f5a6b

$ skills update demo --agent codex
  ~ demo   updated

  1 skill, 1 changed
```

If you edited the skill, the CLI keeps your version and tells you what to do:

```console
$ skills update demo --agent codex
  ! demo   skipped, you edited it

  1 skill, 0 changed, 1 needs attention
  Run 'skills diff demo --agent codex' to see your changes,
  or 'skills update demo --agent codex --force' to overwrite (backed up).
```

Run the `diff` command to see your edits. Use `--force` only when you want the
store version instead.

See what removing a skill would do, without changing any files:

```console
$ skills remove demo --agent codex --dry-run
  - demo   would remove
      would back up first

  1 skill, 1 would change
```

Drop `--dry-run` to remove it for real. Add `-v` to see each file under the skill.

## Use your own skills

You don't need to fork this repository to have your own skill library. Create
a Git repository with a `skills/` folder, one folder per skill:

```
skills/
  go-review/
    SKILL.md                the skill itself
    agents/openai.yaml      optional, copied only for Codex
    references/             optional, copied for every agent
```

The folder name is the skill name. The CLI reads only committed files under
`skills/`, so the rest of the repository can hold anything you like. Look at the
[demo skill](skills/demo/SKILL.md) for a full example.

Then point the CLI at your repository in `~/.config/skills/config.toml`:

```toml
[store]
repo = "https://github.com/your-name/your-skills"
branch = "main"
```

If you already ran the CLI with another repository, move the old store aside
first, since it may hold local work you want to keep:

```sh
mv ~/.config/skills/store ~/.config/skills/store.old
skills init
```

Skills you installed from the old repository stay as they are. `install` and
`update` skip them, even with `--force`. To switch one over, run `skills remove`
and then `skills install` again.

To publish a skill, commit and push it to your branch, then run `skills sync`.
A skill with `status: draft` in its frontmatter is not installed; a skill with
no status is published. Changing a skill never needs a new CLI release.

## Configuration and scope

The first run creates `~/.config/skills/config.toml` and a store at
`~/.config/skills/store/`. If `XDG_CONFIG_HOME` is set, the CLI uses it. All
projects on your machine share the config and the store.

By default, skills are installed for both agents:

| Agent | Project location | Global location |
| --- | --- | --- |
| Claude Code | `.claude/skills/` | `~/.claude/skills/` |
| Codex | `.agents/skills/` | `~/.agents/skills/` |

Gemini CLI also reads `.agents/skills/`, so `--agent codex` covers it too.
Use `--agent` to pick one agent, and `--global` to install for all your projects:

```sh
skills install demo --agent codex
skills install demo --global
skills ls --global
```

In a project, skills go into the repository root. Outside a repository, they go
into the current folder. Edit the config to change agent paths or defaults. The
[default configuration](docs/SPEC.md#configuration) shows the full format.

## Contributing

Fork this repository only if you want to work on the CLI itself.
[docs/SPEC.md](docs/SPEC.md) describes how the CLI behaves, and
[AGENTS.md](AGENTS.md) describes how changes are made.

Use [Task](https://taskfile.dev) for local checks:

```sh
task test:unit           # tests without network or Docker
task test:integration    # integration tests; requires Docker
task lint                # lint and dependency checks
task format              # format Go code
```

See [Taskfile.yml](Taskfile.yml) for all tasks.
