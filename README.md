# Skills & Technical Knowledge

My personal library of AI skills, prompts, and custom instructions, plus a small
Go CLI that installs them into a project and keeps them up to date.

The library is always a work in progress. The CLI exists so that evolving a
skill here costs nothing on the projects that already use it.

## Install the CLI

```sh
go install github.com/madalinpopa/skills/cmd/skills@latest
```

The skills are compiled into the binary, so the binary's version is the
library's version. There is no network access, no cache, and no second thing to
keep in sync.

## Use it

```sh
cd ~/some/project
skills install              # install everything into ./.claude
skills install go-review    # or just one skill
skills list                 # what is available, installed, outdated, modified
```

Later, when the library has moved on:

```sh
go install github.com/madalinpopa/skills/cmd/skills@latest
skills update
```

Add `--global` to any command to act on `~/.claude` instead of the current
project.

## How updating stays safe

`.claude/skills.lock.json` records a SHA-256 for every file the CLI wrote. An
update compares three things per file:

| | meaning |
| --- | --- |
| **want** | what this binary ships |
| **have** | what is on disk now |
| **base** | what the CLI last wrote, from the lock file |

If `have` still equals `base`, you never touched the file, so it is safe to
overwrite. Anything else is your own edit: the update reports a conflict and
leaves the file alone. With `--force` it overwrites, but copies the old content
to `.claude/.skills-backup/<timestamp>/` first.

`skills update --dry-run` prints the plan without writing anything.

## Layout

TODO
