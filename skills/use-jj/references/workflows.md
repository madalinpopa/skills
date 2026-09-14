# Workflows

Step-by-step use cases. Every step uses `-m` or paths so no editor opens.
Replace `<x>` placeholders with change IDs from `jj log`. After any rewrite,
run `jj log` and confirm the graph matches the plan.

## Contents

- [Orient in an unfamiliar repository](#orient-in-an-unfamiliar-repository)
- [Start a new change](#start-a-new-change)
- [Squash workflow: refine as you go](#squash-workflow-refine-as-you-go)
- [Edit workflow: change an existing commit](#edit-workflow-change-an-existing-commit)
- [Stacked changes](#stacked-changes)
- [Fix an older commit in the stack](#fix-an-older-commit-in-the-stack)
- [Split one change into several](#split-one-change-into-several)
- [Combine or reorder changes](#combine-or-reorder-changes)
- [Move a stack onto updated main](#move-a-stack-onto-updated-main)
- [Set aside work and come back](#set-aside-work-and-come-back)
- [Discard work](#discard-work)
- [Cherry-pick and revert](#cherry-pick-and-revert)
- [Merge](#merge)
- [Review someone else's change](#review-someone-elses-change)
- [Find the commit that broke something](#find-the-commit-that-broke-something)
- [Clean history before sharing](#clean-history-before-sharing)
- [Scratch files and private commits](#scratch-files-and-private-commits)

## Orient in an unfamiliar repository

```sh
jj root
jj st
jj log                                  # default: local work plus trunk
jj log -r 'trunk()..@'                  # the stack under the working copy
jj bookmark list --all-remotes
jj op log -n 5
```

Read the output before changing anything. Note whether `@` is empty, whether
it has a description, and which bookmark is nearest.

## Start a new change

```sh
jj git fetch                            # only when the user wants remote state
jj new main -m "feat(cmd): add sync flag"
# edit files; jj snapshots them on the next command
jj st
jj diff
```

If `@` already holds unrelated edits, do not build on top of it. Run
`jj new main` to get a clean sibling; the old working copy remains as its own
change.

## Squash workflow: refine as you go

Work in a fresh empty `@` above the change you are building, then fold parts
down. This keeps the described commit clean while you experiment.

```sh
jj new -m "wip"                         # scratch change on top
# edit
jj squash                               # everything into the parent
jj squash internal/config               # only these paths into the parent
jj squash --into <x>                    # into a specific ancestor
```

`jj squash` keeps the parent's description unless you pass `-m`. Use
`-u/--use-destination-message` to silence the message prompt when both sides
have descriptions.

## Edit workflow: change an existing commit

```sh
jj edit <x>                             # @ now is that commit
# edit files; descendants rebase automatically after each snapshot
jj log                                  # check descendants have no (conflict)
jj new                                  # leave the commit, start a fresh one
```

Prefer the squash workflow when the commit is shared or when you want a
review step before content lands. `jj edit` amends immediately.

## Stacked changes

```sh
jj new main -m "refactor(store): extract clone helper"
# edit
jj new -m "feat(store): add shallow clone"
# edit
jj new -m "docs: describe shallow clone"
jj log -r 'trunk()..@'
```

Each `jj new` starts the next change on top of the previous one. Bookmarks
for pushing are set later, one per change or one at the top.

## Fix an older commit in the stack

Option A, targeted squash:

```sh
# make the fix in @
jj squash --into <x> path/to/file       # only that file moves down
```

Option B, absorb everything automatically:

```sh
# make several fixes in @
jj absorb                               # each hunk goes to the commit that last touched those lines
jj log -p -r 'trunk()..@'               # verify placement
```

`jj absorb` only targets mutable ancestors and leaves hunks it cannot place
in `@`. Report what stayed behind.

## Split one change into several

By path, no editor:

```sh
jj split -r <x> -m "test(store): add clone tests" 'glob:"**/*_test.go"'
```

The first commit gets the listed paths and the new message. The second keeps
the rest and the original description. Use `-p/--parallel` to make the two
commits siblings instead of a chain. Repeat `jj split` on the remainder for
three or more pieces.

## Combine or reorder changes

```sh
jj squash --from <b> --into <a>         # merge b into a, keep a's message
jj squash --from <b> --into <a> -m "combined message"
jj rebase -r <c> -B <b>                 # move c before b
jj rebase -r <c> -A <a>                 # move c right after a
jj parallelize <a>::<c>                 # make a chain into siblings
```

For a merge with both messages, run `jj describe <a> -m "..."` afterward.

## Move a stack onto updated main

```sh
jj git fetch
jj rebase -b @ -o main                  # whole branch containing @
# or
jj rebase -s 'roots(trunk()..@)' -o main
jj log -r 'conflicts()'                 # resolve any, see conflicts.md
```

`-b` finds the fork point with the destination and moves everything after it.
`-s` moves a chosen root and its descendants. Immutable commits are skipped
by definition because they are already on trunk.

## Set aside work and come back

There is no stash. The working copy is a change already.

```sh
jj describe -m "wip: half-done sync"    # name it so you can find it
jj new main                             # fresh working copy elsewhere
# later
jj log -r 'description(glob:"wip:*")'
jj edit <wip>                           # resume amending it
# or
jj new <wip>                            # build on it with a new change
```

## Discard work

```sh
jj restore path/to/file                 # drop edits to one file in @
jj restore                              # drop every edit in @
jj abandon                              # drop @ entirely, start empty
jj abandon <x>                          # hide x; children move to x's parent
jj abandon <x> --retain-bookmarks       # keep bookmarks that pointed at x
```

Abandoned commits remain in the operation log. See operations.md to recover.

## Cherry-pick and revert

```sh
jj duplicate <x> -o main                # copy x onto main, keep original
jj duplicate <x> -A @-                  # insert a copy into the stack
jj revert -r <x> -B @                   # new change that undoes x, before @
jj revert -r <x> -o @                   # the undo change on top of @
```

## Merge

```sh
jj new main feature -m "merge feature into main"
jj log -r 'conflicts()'
```

A clean merge shows as `(empty)` because its content equals the auto-merge of
its parents. To undo a merge, create a new change and restore from the first
parent: `jj new <merge>` then `jj restore --from '<merge>-'`.

## Review someone else's change

```sh
jj git fetch
jj log -r 'remote_bookmarks(glob:"pr/*")'
jj show <x>
jj diff -r '<base>..<tip>' --stat
jj interdiff --from <old> --to <new>    # what changed between two versions
jj file annotate -r <x> path/to/file
jj new <x>                              # try it locally without editing it
```

## Find the commit that broke something

```sh
jj bisect run -r 'main..@' -- go test ./internal/store
```

The command is run at each candidate; exit code 0 marks good. Use
`--find-good` to search for the first good commit instead. Manual bisection:
`jj new 'bisect(main..@)'`, test, then narrow the range.

## Clean history before sharing

```sh
jj log -r 'trunk()..@' -T builtin_log_compact
jj log -r '(trunk()..@) & empty()'      # empty changes to abandon
jj log -r '(trunk()..@) & description(exact:"")'   # missing messages
jj describe <x> -m "..."                # fix messages
jj squash --from <fixup> --into <target>
jj simplify-parents -r 'trunk()..@'     # drop redundant merge edges
jj fix -s 'trunk()..@'                  # run configured formatters, if any
```

Run tests at the top after cleaning. `jj git push --dry-run -b <name>` shows
what a push would do without doing it.

## Scratch files and private commits

Files match `.gitignore` as in Git. For untracked scratch files without ignore
rules, set `snapshot.auto-track = "none()"` and use `jj file track` for the
files that belong. For local-only changes that must never be pushed, keep them
in a change whose description starts with `private:` and set
`git.private-commits = "description(glob:'private:*')"` so push refuses them.
