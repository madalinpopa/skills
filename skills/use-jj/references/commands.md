# jj command reference

Read only the section needed for the task. This file holds defaults, traps, and
choices that `jj <command> -h` does not make obvious. Use installed help for the
full flag list; it takes precedence over this file.

## Contents

- [Choose a command](#choose-a-command)
- [Defaults that change the result](#defaults-that-change-the-result)
- [Commands that block a session](#commands-that-block-a-session)
- [Bookmarks](#bookmarks)
- [Operations](#operations)
- [Workspaces](#workspaces)
- [Global options](#global-options)

## Choose a command

| Need | Command |
| --- | --- |
| Status, graph, one revision | `jj st`, `jj log`, `jj show <rev>` |
| Diff of `@`, a revision, or a range | `jj diff`, `jj diff -r <rev>`, `jj diff --from A --to B` |
| Compare two versions of a change | `jj interdiff --from <old> --to <new>` |
| Earlier versions of one change | `jj evolog -r <change>` |
| Blame, list, grep, or print files | `jj file annotate`, `jj file list`, `jj file search -p <pattern>`, `jj file show` |
| New change, or insert one mid-stack | `jj new [parents]`, `jj new -A <rev>` / `-B <rev>` |
| Set a message | `jj describe [rev] -m` |
| Finish `@` and start a new change | `jj commit -m` |
| Amend an existing revision directly | `jj edit <rev>` |
| Move content into another change | `jj squash`, `jj absorb` |
| Split a change | `jj split <paths>` |
| Discard edits | `jj restore [paths]`, `jj abandon [revs]` |
| Copy or undo a change | `jj duplicate <rev>`, `jj revert -r <rev>` |
| Move changes in the graph | `jj rebase`, `jj parallelize`, `jj simplify-parents` |
| Merge divergent versions | `jj converge -r <change> --no-interactive` |
| Change author or change ID only | `jj metaedit` |
| Run formatters over a stack | `jj fix -s <revs>` |
| Find the first bad revision | `jj bisect run -r <range> -- <command>` |
| Read settings | `jj config list`, `jj config get <name>` |

## Defaults that change the result

- `jj new` with no revision builds on `@`. Several revisions make a merge.
  `-A`/`-B` insert the change and rebase the rest of the stack on top.
- `jj squash` with no arguments moves all of `@` into its parent, like
  `git commit --amend`. `--from X --into Y` abandons X once it is empty.
- `jj split <paths>` gives the listed paths and the `-m` message to the first
  commit. The second keeps the rest and the original message. `-p` makes them
  siblings.
- `jj restore` with no paths drops every edit in `@`. `jj abandon` on `@`
  drops the whole change. Abandoned descendants rebase onto the parent.
- `jj rebase` with no selection flag uses `-b @`. `-r X` moves only X and
  leaves its descendants on X's old parent. `-s X` moves X and descendants.
  `-b X` moves everything since the fork point with the destination.
- Placement flags mean the same on `new`, `squash`, `split`, `rebase`,
  `duplicate`, and `revert`: `-o` makes children of the target, `-A` inserts
  after it, and `-B` inserts before it.
- `jj absorb` only targets mutable ancestors and leaves unplaced hunks in `@`.
- `jj file untrack` needs an ignore rule, or the next snapshot tracks the file
  again.

## Commands that block a session

These open an editor or interactive tool and stall a non-interactive run:

- `jj describe` or `jj commit` without `-m`.
- `jj squash` when both sides have descriptions, unless `-m` or `-u` is given.
- `jj split` without paths, `jj diffedit`, `jj resolve` without
  `--list`, `jj arrange`, `jj config edit`.
- Any `-i/--interactive`, `--tool`, or `--editor` flag.

Use `-m`, paths, and filesets instead. `jj describe --stdin` reads a message
from standard input.

## Bookmarks

- `jj bookmark create` fails if the name exists. `set` creates or moves.
  `move` changes only existing bookmarks. Moving backward or sideways needs
  `-B/--allow-backwards`.
- `jj bookmark delete` propagates the deletion on the next push.
  `jj bookmark forget` drops it locally only.
- `jj bookmark rename` leaves the old name on the remote.
- Bookmark arguments accept patterns such as `glob:"feature/*"`.

## Operations

- `jj undo` steps back one operation at a time; snapshots count as operations.
- `jj op restore <op>` resets the whole repository. `jj op revert <op>` undoes
  one operation and keeps later ones. `--what remote-tracking` limits restore
  to remote bookmark state.
- `jj op abandon` and `jj util gc` remove recovery data permanently. Ask
  before running them.

## Workspaces

The first workspace is `default`. `<name>@` addresses another workspace's
working copy in revsets. A workspace rewritten elsewhere becomes stale; run
`jj workspace update-stale`.

## Global options

| Option | Effect |
| --- | --- |
| `-R <path>` | Operate on another repository |
| `--ignore-immutable` | Allows rewriting immutable commits. Ask before using. |
| `--ignore-working-copy` | Skips the snapshot; can leave the working copy stale |
| `--at-op <op>` | Reads the repository as of an operation; not a read-only guard |
| `--no-integrate-operation` | Records an operation without making it current; `jj op integrate` applies it |
| `--config NAME=VALUE` | Overrides a setting for one command |
| `--color never`, `--no-pager` | Plain output for parsing |
