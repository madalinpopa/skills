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

The CLI needs [Git](https://git-scm.com/) 2.45 or newer on your `PATH`. The
store is a partial, sparse clone, so it downloads only what `skills/` needs.

## First run

The first command that needs the store creates `~/.config/skills/config.toml`
and clones this repository into `~/.config/skills/store/`. That local clone is
the store every project installs from, so one copy serves the whole machine.

```sh
skills init      # do the first run explicitly
```

`XDG_CONFIG_HOME` is honoured when set. `--dry-run` never creates the config or
clones the store. If the store is missing, it tells you to run `skills init`.

The config lists the agents and where each keeps its skills. A project path
must be relative and stay inside the project. A global path must be absolute
or start with `~/`. The default agent list must name at least one configured
agent. Anything else is a config error and no command runs.

### Default configuration

The default values written to `~/.config/skills/config.toml` on first run are:

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

## Use it

```sh
skills ls                   # published skills in the store
skills install go-review    # install into this project
skills ls --local           # what this project has installed
skills sync                 # fast-forward the store
skills update               # move installed skills to the synced content
```

Skills land in the directory each agent reads. 

```
.claude/skills/go-review/   Claude Code
.agents/skills/go-review/   Codex (Gemini CLI reads it too)
```

The default agents are `claude` and `codex`. Use `--agent` to narrow a command.
Gemini CLI reads the shared `.agents/skills` directory, so `--agent codex`
covers it. To use `--agent gemini`, add an `[agents.gemini]` entry with the
same paths to the config.
Use `--global` to install into your home directory instead, so a skill is
available in every project. Scope is explicit: a command never falls back from
the project to your home directory.

```sh
skills install go-review --agent claude
skills install go-review --global
```

## Examples

Install a skill. `-v` shows the files written for each agent:

```
$ skills install demo -v
  + demo   added
      .agents/skills/demo/SKILL.md
      .agents/skills/demo/references/guidelines.md
      .claude/skills/demo/SKILL.md
      .claude/skills/demo/references/guidelines.md

  1 skill, 1 changed
```

Update after you edited a skill. The edited skill is skipped and the command
exits 3:

```
$ skills update
  ~ go-review        updated
  ! demo             skipped, you edited it

  2 skills, 1 changed, 1 needs attention
  Run 'skills diff demo' to see your changes,
  or 'skills update demo --force' to overwrite (backed up).
```

See your changes, then overwrite them:

```
$ skills diff demo
--- a/.claude/skills/demo/SKILL.md
+++ b/.claude/skills/demo/SKILL.md
@@ -65,3 +65,4 @@
 permission checks.
+my local note

$ skills update demo --force
  ~ demo   updated
      backed up to ~/.config/skills/backups/20260913-004647/.../.claude/skills/demo

  1 skill, 1 changed
```

Preview a change with `--dry-run`. Nothing is written:

```
$ skills remove demo --dry-run
  - demo   would remove
      would back up first

  1 skill, 1 would change
```

Symbols carry the meaning, so nothing is lost when color is off:

```
+  added      ~  updated      -  removed      !  needs attention
```

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

```
--global      act on the home directory instead of the repository
--agent       narrow to certain agents (default: config defaults)
--force       overwrite or remove skills you have edited (backed up first)
--dry-run     show the plan, write nothing
-v            show individual files and their targets
--no-color    disable color (NO_COLOR is honoured too)
```

Exit codes: `0` success, `1` runtime error, `2` bad arguments, `3` a skill
needs attention.

## Writing a skill

Use [create-repo-skill](skills/create-repo-skill/SKILL.md) to create or update a
skill in this repository. It starts with `skill-creator` and contains the
repository's authoring workflow, CLI frontmatter rules, and agent options.
Read it directly from the source path; it does not need to be installed.

See its [naming and category legend](skills/create-repo-skill/references/naming.md)
for prefixes such as `create-`, `use-`, and `review-`, plus `workflow-` for
skills that coordinate other skills toward a larger outcome. All skill source
files belong under `skills/<skill-name>/`.

## Development

Common work runs through [Task](https://taskfile.dev):

```sh
task                     # unit tests, integration tests and lint
task test:unit           # go test ./... (needs neither network nor Docker)
task test:integration    # INTEGRATION=true go test ./... -count=1 (needs Docker)
task lint                # golangci-lint run and go mod tidy -diff
task format              # go fmt and golangci-lint fmt
task security            # govulncheck
```

Releases are tags. `task release:patch`, `task release:minor` or `task release:major` creates the next
version tag and pushes it. CI runs formatting, vet, lint and tests on every
pull request. Pushing a `v*` tag builds and publishes the release binaries
with GoReleaser.

## Specification

[docs/SPEC.md](docs/SPEC.md) is the full contract: config layout, scope rules,
the three-way comparison that protects local edits, backups, locks, exit
codes and the differences between agents.
