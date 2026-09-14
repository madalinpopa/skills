# Git interoperation

Most jj repositories back onto Git. Colocated repositories keep `.jj` and
`.git` side by side so Git tools still work. This file covers how the two
stay in sync, bookmarks versus branches, fetch and push, pull-request flows,
and the Git command table.

## Contents

- [Colocated repositories](#colocated-repositories)
- [Bookmarks versus branches](#bookmarks-versus-branches)
- [Fetch](#fetch)
- [Push](#push)
- [Pull-request flows](#pull-request-flows)
- [Remotes and cloning](#remotes-and-cloning)
- [Git command table](#git-command-table)

## Colocated repositories

`jj git init` and `jj git clone` create colocated repositories by default.
Check with `jj git colocation status`. In a colocated repository:

- Every jj command imports Git refs first and exports jj bookmarks to Git
  branches afterwards. `git branch` and `jj bookmark list` agree.
- Git's HEAD is kept detached at `@-`, the parent of the working copy. Git
  sees the working-copy change as uncommitted edits. This is expected.
- Bookmarks with no Git counterpart yet appear as `name@git` after export.
- Use `git` for read-only work: `git log`, `git show`, `git blame`, `git
  diff`. Mutating Git commands (`commit`, `checkout`, `rebase`, `reset`,
  `stash`) confuse the mapping and can create divergent changes. If one was
  run, `jj undo` or `jj op restore` recovers the jj view; `jj git import`
  pulls in refs Git changed.
- Git tools cannot read jj conflict markers as merged content, so resolve
  conflicts through jj before using Git tooling on those files.

Convert with `jj git colocation enable` or `disable`. Wrap an existing Git
checkout with `jj git init --git-repo=.` from its root, or `jj git init
--colocate` inside it.

## Bookmarks versus branches

| Git | jj |
| --- | --- |
| Branch | Bookmark. Does not advance when you commit. |
| Current branch | None. `@` is a commit, not a branch tip. |
| `origin/main` | `main@origin`, a remote bookmark |
| Tracking branch | Tracked remote bookmark: `jj bookmark track main@origin` |
| Deleted branch | `jj bookmark delete` propagates the deletion on push; `forget` does not |

Because bookmarks stay put, the common step before pushing is to move one:

```sh
jj bookmark set feature -r @-           # after jj commit, point at the finished change
jj bookmark move feature --to @         # when the working copy itself is the tip
jj bookmark move --from 'heads(::@- & bookmarks())' --to @-   # advance whichever bookmark sits below
```

Moving a bookmark backward or sideways needs `-B/--allow-backwards`.

Untracked remote bookmarks are immutable by default. Track the ones you
work on so that fetch updates the local bookmark:

```sh
jj bookmark list --all-remotes
jj bookmark track feature@origin
```

Set `remotes.origin.auto-track-bookmarks = "glob:*"` to track every fetched
bookmark from that remote.

## Fetch

```sh
jj git fetch                            # default remote, all bookmarks
jj git fetch --remote upstream
jj git fetch -b 'glob:"release/*"'
jj git fetch --all-remotes
```

Fetch updates `name@remote`. For tracked bookmarks it also moves the local
bookmark. If both sides moved, the local bookmark becomes conflicted and
shows as `name??`; see conflicts.md. Fetched commits that descend from
hidden local commits can make those visible again as divergent changes.

## Push

`jj git push` refuses to run outside safety checks. It only updates a remote
bookmark whose current position matches what was last fetched, like
`--force-with-lease`. It refuses conflicted bookmarks, commits with
conflicts, empty descriptions, and commits matching `git.private-commits`.

```sh
jj git push --dry-run                   # what would happen
jj git push -b feature                  # one bookmark
jj git push -b 'glob:"feat/*"'          # by pattern
jj git push -c @-                       # create push-<changeid> bookmark and push it
jj git push -r 'trunk()..@-'            # every bookmark on those revisions
jj git push --named pr-sync=@-          # create and push a bookmark in one step
jj git push --all                       # all bookmarks and tags
jj git push --tracked                   # only tracked bookmarks
jj git push --deleted                   # propagate local deletions
jj git push --remote upstream -b main
```

Defaults: without flags, push sends tracked bookmarks in
`remote_bookmarks(remote=<remote>)..@`. The remote comes from `git.push` or
`origin`. Pushing a new bookmark starts tracking it. A push never moves
bookmarks; move them first.

Always confirm the user wants a push. Show `--dry-run` output when the push
would create, move, or delete a remote bookmark that the user did not name.

## Pull-request flows

One change per PR:

```sh
jj new main -m "fix(config): default agents"
# edit
jj bookmark create fix-default-agents -r @
jj git push -b fix-default-agents
# after review feedback
# edit @ again, then:
jj git push -b fix-default-agents       # the bookmark follows the rewritten commit
```

Stacked PRs:

```sh
jj bookmark create part-1 -r <a>
jj bookmark create part-2 -r <b>
jj git push -b part-1 -b part-2
```

Open each PR with the previous bookmark as its base. After `main` moves,
rebase the stack, then push both bookmarks again.

Update a PR when the working copy is the tip: `jj bookmark move feature --to
@` before pushing. After `jj commit`, the finished change is `@-`, so use
`-r @-`.

Land and clean up after merge:

```sh
jj git fetch
jj log -r 'feature | main'              # feature now inside main?
jj bookmark delete feature              # local; deletion pushes with --deleted
jj new main
```

## Remotes and cloning

```sh
jj git clone https://example.com/repo.git
jj git clone --no-colocate <url> dir
jj git remote list
jj git remote add upstream <url>
jj git remote set-url origin <url>
jj git remote rename origin github
```

Clone tracks the default bookmark of the remote. `trunk()` resolves to that
bookmark on `upstream` if present, otherwise `origin`.

## Git command table

| Git | jj |
| --- | --- |
| `git init` | `jj git init` |
| `git clone <src> <dst>` | `jj git clone <src> <dst>` |
| `git fetch` | `jj git fetch` |
| `git push --all` | `jj git push --all` |
| `git push origin <branch>` | `jj git push -b <bookmark>` |
| `git remote add <n> <url>` | `jj git remote add <n> <url>` |
| `git status` | `jj st` |
| `git diff HEAD` | `jj diff` |
| `git diff <rev>^ <rev>` | `jj diff -r <rev>` |
| `git diff <rev>` | `jj diff --from <rev>` |
| `git diff A B` | `jj diff --from A --to B` |
| `git diff A...B` | `jj diff -r A..B` |
| `git show <rev>` | `jj show <rev>` |
| `git add <file>` | nothing; new files are tracked on snapshot |
| `git rm <file>` | `rm <file>` |
| `git rm --cached <file>` | `jj file untrack <file>` (needs an ignore rule) |
| `git commit -a` | `jj commit` |
| `git commit --amend -a` | `jj squash` |
| `git add -p; git commit --amend` | `jj squash -i` or `jj squash <paths>` |
| `git commit --fixup=X; git rebase --autosquash` | `jj squash --into X` |
| `git commit -p` | `jj split` |
| `git commit --amend --only` | `jj describe @-` |
| `git log --oneline --graph` | `jj log -r ::@` |
| `git log --graph --all` | `jj log -r 'all()'` |
| `git log --branches --not upstream/main` | `jj log` |
| `git log -G text` | `jj log -r 'diff_lines(regex:text)'` |
| `git ls-files` | `jj file list` |
| `git grep foo` | `jj file search --pattern=foo` |
| `git blame <file>` | `jj file annotate <file>` |
| `git reset --hard` (drop change) | `jj abandon` |
| `git reset --hard` (empty change) | `jj restore` |
| `git reset --soft HEAD~` | `jj squash --from @-` |
| `git restore <paths>` | `jj restore <paths>` |
| `git stash` | `jj new @-` |
| `git switch -c topic main` | `jj new main` |
| `git merge A` | `jj new @ A` |
| `git checkout v1.0` | `jj new v1.0` |
| `git rebase B A` | `jj rebase -b A -o B` |
| `git rebase --onto B A^ tip` | `jj rebase -s A -o B` |
| `git rebase -i` (reorder) | `jj rebase -r C --before B` or `jj arrange` |
| `git rebase -i` (edit) | `jj diffedit -r X`, `jj split -r X` |
| `git rebase --continue` | edit the file, then `jj squash` |
| `git cherry-pick X` | `jj duplicate X -o @` |
| `git revert X` | `jj revert -r X -B @` |
| `git rev-parse --show-toplevel` | `jj root` |
| `git branch` | `jj bookmark list` |
| `git branch <n> <rev>` | `jj bookmark create <n> -r <rev>` |
| `git branch -f <n> <rev>` | `jj bookmark move <n> --to <rev>` (`-B` to go backward) |
| `git branch -d <n>` | `jj bookmark delete <n>` |
| `git tag <n> <rev>` | `jj tag set <n> -r <rev>` |
| `git tag -d <n>` | `jj tag delete <n>` |
| `git tag --contains <rev>` | `jj tag list -r '<rev>::'` |
| `git reflog` | `jj op log`, `jj evolog` |
| `git reflog` restore | `jj undo`, `jj op restore <op>` |
