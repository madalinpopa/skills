# Operations, undo, and recovery

Every jj command that changes the repository records an operation: bookmark
targets, heads, working-copy commits, and Git refs at that moment. The log is
append-only, so any earlier state can be viewed or restored. Hidden commits
are never deleted until `jj util gc` runs.

## Contents

- [Read the operation log](#read-the-operation-log)
- [Undo and redo](#undo-and-redo)
- [Restore or revert a specific operation](#restore-or-revert-a-specific-operation)
- [View the past without changing it](#view-the-past-without-changing-it)
- [Recover a hidden or abandoned commit](#recover-a-hidden-or-abandoned-commit)
- [Evolution of one change](#evolution-of-one-change)
- [Stale working copies and workspaces](#stale-working-copies-and-workspaces)
- [Concurrent operations](#concurrent-operations)

## Read the operation log

```sh
jj op log                               # newest first, @ is current
jj op log -n 10 --no-graph
jj op show <op> -p                      # what one operation changed
jj op diff --from <op1> --to <op2>
jj op log --show-changes-in <change>    # operations that touched a change
```

Operation IDs are hex prefixes. `@` is the current operation, `@-` its
parent. Each entry names the command that ran, so the log explains how the
repository reached its state.

## Undo and redo

```sh
jj undo                                 # undo the latest operation
jj undo                                 # run again to go one further back
jj redo                                 # reapply the last undone operation
```

Undo is itself an operation. It undoes the previous operation's effect while
keeping the log intact. Snapshots count as operations too, so a plain `jj
undo` after editing files may only undo a snapshot; check `jj op log` first.

## Restore or revert a specific operation

```sh
jj op restore <op>                      # whole repository as of <op>
jj op restore @-                        # equivalent to one undo
jj op revert <op>                       # cancel only that operation, keep later ones
jj op restore <op> --what remote-tracking   # only remote bookmark state
```

Prefer `op restore` when several commands went wrong; prefer `op revert` when
one middle operation was the mistake. Both are safe to try because the
restore itself can be undone.

## View the past without changing it

```sh
jj log --at-op <op>
jj log --at-op <op> -r 'all()'
jj diff --at-op <op> -r <x>
jj log -r 'at_operation(<op>, @)'       # the working copy as of that op
```

`--at-op` disables snapshotting for that command, so it never alters the log.

## Recover a hidden or abandoned commit

Abandoned, squashed, or rewritten commits are hidden, not deleted.

```sh
jj evolog -r <change>                   # earlier versions of a change
jj log --at-op <op-before-abandon>      # the graph as it was, with IDs
jj log -r '<commit-id> & hidden()'      # hidden() needs an explicit ID
jj new <commit-id>                      # revive as a child, safest
jj bookmark set rescued -r <commit-id>  # or make it visible by naming it
jj op restore <op-before-abandon>       # or rewind everything
```

Reviving a hidden commit whose change ID still has a visible version creates
divergence; use the commit ID and see conflicts.md if that happens.

## Evolution of one change

```sh
jj evolog                               # versions of @
jj evolog -r <change> -p                # with diffs between versions
jj evolog -r <change> --reversed        # oldest first
```

Each row is a previous commit ID of the same change. Use `jj restore --from
<old-commit-id> path` to bring back a file from a previous version.

## Stale working copies and workspaces

A working copy becomes stale when its commit was rewritten from another
workspace or with `--ignore-working-copy`. jj reports this on the next
command.

```sh
jj workspace update-stale
jj workspace list
jj workspace add ../repo-review --name review -r main
jj workspace forget review
```

Each workspace has its own `@`, addressed as `review@` in revsets. All
workspaces share one operation log. Set `snapshot.auto-update-stale = true`
to refresh automatically.

## Concurrent operations

Two commands running at once create two operation heads. jj merges them on
the next command and reports divergent changes or conflicted bookmarks when
both sides touched the same thing. Resolve as in conflicts.md. Commands run
with `--no-integrate-operation` stay unintegrated until `jj op integrate`.

Old operations can be trimmed with `jj op abandon '..<op>'`, and unreachable
data with `jj util gc`. Both are irreversible; ask before running them.
