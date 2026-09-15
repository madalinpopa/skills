# jj command reference

Read only the section needed for the task. Use `jj <command> -h` for installed
flags and `--help` for detail; examples here are a lookup aid, not a required
read or an exhaustive CLI contract. Global options are listed at the end.

## Contents

- [Inspect](#inspect)
- [Create and describe changes](#create-and-describe-changes)
- [Move content between changes](#move-content-between-changes)
- [Move changes in the graph](#move-changes-in-the-graph)
- [Working copy and files](#working-copy-and-files)
- [Bookmarks and tags](#bookmarks-and-tags)
- [Git](#git)
- [Operations and recovery](#operations-and-recovery)
- [Workspaces and sparse checkouts](#workspaces-and-sparse-checkouts)
- [Config, signing, and utilities](#config-signing-and-utilities)
- [Global options](#global-options)
- [Flags shared across commands](#flags-shared-across-commands)

## Inspect

| Command | Purpose | Key flags |
| --- | --- | --- |
| `jj status` (`jj st`) | Working-copy changes, parent, conflicts, conflicted bookmarks | paths to limit |
| `jj log` | Graph of revisions; default revset is `revsets.log` | `-r <revset>`, `-n <limit>`, `--reversed`, `-G/--no-graph`, `-T <template>`, `-p/--patch`, `-s/--summary`, `--stat`, `--count` |
| `jj show [rev]` | Description and diff of one revision | `-T`, `-s`, `--stat`, `--git`, `--no-patch` |
| `jj diff` | Diff of `@` or between revisions | `-r <revsets>`, `-f/--from`, `-t/--to`, `-s`, `--stat`, `--git`, `--name-only`, `--types`, `--context <n>`, `-w`, `-b`, paths or filesets |
| `jj interdiff` | Difference between the diffs of two revisions | `-f/--from`, `-t/--to`, diff format flags |
| `jj evolog` | How a change was rewritten over time, including hidden versions | `-r <revsets>`, `-n`, `-p`, `--reversed`, `-G` |
| `jj file annotate <path>` | Line-by-line origin, like blame | `-r <rev>`, `-T` |
| `jj file list` | Tracked files in a revision | `-r <rev>`, `-T`, paths |
| `jj file search` | Grep tracked file contents | `-r <rev>`, `-p/--pattern`, `--name-only`, `-n/--line-number` |
| `jj file show <path>` | Print file contents at a revision | `-r <rev>` |
| `jj root` | Workspace root directory | |
| `jj version` | Installed version | |
| `jj help <command>` | Full help, also `-k` for docs topics such as `revsets` | |

Useful log calls:

```sh
jj log -r 'all()'                       # everything visible
jj log -r '::@'                         # ancestors of the working copy
jj log -r 'trunk()..@'                  # the current stack
jj log -r 'conflicts()'                 # conflicted commits
jj log --no-graph -T 'change_id.short() ++ " " ++ description.first_line() ++ "\n"'
```

## Create and describe changes

| Command | Purpose | Key flags |
| --- | --- | --- |
| `jj new [revs...]` | New empty change on top of the given parents (default `@`) and edit it. Several parents make a merge. | `-m <msg>`, `--no-edit`, `-A/--insert-after <revs>`, `-B/--insert-before <revs>` |
| `jj describe [rev]` (`jj desc`) | Set the description of `@` or another revision | `-m <msg>`, `--stdin`, `--editor` |
| `jj commit` (`jj ci`) | Describe `@` and create a new empty change on top | `-m <msg>`, `-i`, `--tool`, paths |
| `jj edit <rev>` | Make an existing revision the working copy. Later edits amend it directly. | |
| `jj next` | Move the working copy to the child revision | `-e/--edit`, `-n/--no-edit`, `--conflict` |
| `jj prev` | Move the working copy toward the parent | `-e/--edit`, `-n/--no-edit`, `--conflict` |
| `jj duplicate [revs]` | Copy revisions, like cherry-pick, keeping the originals | `-o/--onto <revs>`, `-A`, `-B` |
| `jj revert` | New change that reverses the given revisions | `-r <revs>`, `-o`, `-A`, `-B` |
| `jj metaedit [rev]` | Change author, timestamps, or change ID without touching content | `--update-change-id`, `--update-author`, `--author <name <email>>`, `--update-author-timestamp`, `--author-timestamp`, `-m` |
| `jj abandon [revs]` | Hide revisions; descendants rebase onto the abandoned parent | `--retain-bookmarks`, `--restore-descendants` |

`jj new` with `-A`/`-B` inserts a change in the middle of a stack and rebases
the rest automatically. `jj new a b` creates a merge of `a` and `b`.

## Move content between changes

| Command | Purpose | Key flags |
| --- | --- | --- |
| `jj squash` | Move changes from a revision (default `@`) into another (default parent). Paths limit which files move. | `-r <rev>`, `-f/--from <revs>`, `-t/--into <rev>`, `-o/--onto`, `-A`, `-B`, `-m <msg>`, `-u/--use-destination-message`, `-i`, `-k/--keep-emptied`, paths |
| `jj split [paths]` | Split a revision in two. With paths, the first commit gets those paths, the second gets the rest. Without paths it opens a diff editor. | `-r <rev>`, `-m <msg>` for the first part, `-p/--parallel`, `-o`, `-A`, `-B`, `-i`, `--tool` |
| `jj absorb` | Push each changed hunk from `@` into the nearest mutable ancestor that last touched those lines | `-f/--from <rev>`, `-t/--into <revs>`, `-i`, paths |
| `jj restore [paths]` | Copy paths from one revision into another, default from parent into `@` (discard working-copy edits) | `-f/--from <rev>`, `-t/--into <rev>`, `-c/--changes-in <rev>`, `-i`, `--restore-descendants` |
| `jj diffedit` | Edit the diff of a revision in a diff editor | `-r <rev>`, `-f`, `-t`, `--tool`, `--restore-descendants` |
| `jj fix` | Run configured formatters over changed files in a set of revisions | `-s/--source <revs>`, `--include-unchanged-files`, `-a/--all-lines`, paths |

`jj squash` without arguments is Git's `commit --amend`. `jj squash --from X
--into Y` moves all of X into Y and abandons X when it becomes empty.

## Move changes in the graph

| Command | Purpose | Key flags |
| --- | --- | --- |
| `jj rebase` | Move revisions to new parents | `-b/--branch <revs>` whole branch, `-s/--source <revs>` revision and descendants, `-r/--revision <revs>` only those revisions, `-o/--onto <revs>`, `-A/--insert-after`, `-B/--insert-before`, `--skip-emptied`, `--keep-divergent`, `--simplify-parents` |
| `jj parallelize <revs>` | Turn a linear chain into siblings sharing the same parent | |
| `jj simplify-parents` | Remove redundant parent edges | `-s <revs>`, `-r <revs>` |
| `jj arrange` | Interactive graph editor. Avoid in agent sessions. | |
| `jj converge` | Merge divergent versions of one change | `-r <revs>`, `--no-interactive` |
| `jj bisect run` | Run a command across a range to find the first bad revision | `-r/--range <revs>`, `--find-good`, `-- <command>` |

Rebase selection:

- `-r X` moves only X. Its descendants are rebased onto X's old parent.
- `-s X` moves X together with all descendants.
- `-b X` moves the whole branch containing X, that is everything since the
  fork point with the destination. This is the default when no flag is given.

Destinations: `-o` makes the revisions children of the target, `-A` inserts
after the target and rebases the target's children on top, `-B` inserts
before the target.

## Working copy and files

| Command | Purpose | Key flags |
| --- | --- | --- |
| `jj file track <paths>` | Start tracking paths when `snapshot.auto-track` excludes them | `--include-ignored` |
| `jj file untrack <paths>` | Stop tracking paths. They must match an ignore rule to stay untracked. | |
| `jj file chmod x|n <paths>` | Set or clear the executable bit | `-r <rev>` |
| `jj util snapshot` | Snapshot the working copy without another command | |
| `jj resolve [paths]` | Open a merge tool for conflicted files | `-r <rev>`, `-l/--list`, `--tool` |

`jj restore` with no arguments discards every working-copy edit. With paths it
discards only those paths. `jj abandon` on `@` drops the whole change and
starts a fresh empty one.

## Bookmarks and tags

`jj bookmark` is aliased to `jj b`. Subcommands accept unique prefixes.

| Command | Purpose | Key flags |
| --- | --- | --- |
| `jj bookmark list` | Local and tracked remote bookmarks | `-a/--all-remotes`, `--remote <name>`, `-t/--tracked`, `-c/--conflicted`, `-r <revs>`, `-T`, `--sort` |
| `jj bookmark create <name>` | New bookmark, fails if it exists | `-r <rev>` (default `@`) |
| `jj bookmark set <name>` | Create or move a bookmark | `-r <rev>`, `-B/--allow-backwards` |
| `jj bookmark move <names>` | Move existing bookmarks | `-f/--from <revs>` move all bookmarks on those revisions, `-t/--to <rev>`, `-B/--allow-backwards` |
| `jj bookmark advance <names>` | Move the closest bookmarks forward to a target | `-t/--to <rev>` |
| `jj bookmark rename <old> <new>` | Rename locally; the remote keeps the old name | `--overwrite-existing` |
| `jj bookmark delete <names>` | Delete locally and delete on remote at next push | patterns allowed |
| `jj bookmark forget <names>` | Drop the local bookmark without propagating a deletion | `--include-remotes` |
| `jj bookmark track <name>@<remote>` | Start tracking a remote bookmark | `--remote` |
| `jj bookmark untrack <name>@<remote>` | Stop tracking | `--remote` |
| `jj tag list` | Tags, optionally filtered | `-r <revs>`, `-T`, `--sort` |
| `jj tag set <name>` | Create or move a tag | `-r <rev>` |
| `jj tag delete <name>` | Delete a tag | |
| `jj tag track` / `jj tag untrack` | Track state for remote tags | `--remote` |

Bookmark names accept string patterns such as `glob:"feature/*"` where a
pattern is allowed.

## Git

| Command | Purpose | Key flags |
| --- | --- | --- |
| `jj git init [path]` | New repository, colocated by default | `--colocate`, `--no-colocate`, `--git-repo <path>` to wrap an existing Git repo |
| `jj git clone <url> [dest]` | Clone, colocated by default | `--remote <name>`, `--colocate`, `--no-colocate`, `--depth`, `-b/--branch`, `-t/--tag` |
| `jj git fetch` | Fetch bookmarks and tags | `-b/--branch <pattern>`, `-t/--tag`, `--tracked`, `--remote <name>`, `--all-remotes` |
| `jj git push` | Push bookmarks; default set is tracked bookmarks in `remote_bookmarks(remote=<remote>)..@` | `-b/--bookmark <pattern>`, `-t/--tag`, `--all`, `--tracked`, `--deleted`, `-r/--revision <revs>` bookmarks pointing at revs, `-c/--change <revs>` create `push-<changeid>` bookmarks, `--named <NAME=REV>`, `--remote <name>`, `--dry-run`, `--allow-empty-description`, `--allow-private`, `--allow-conflicts` |
| `jj git remote list|add|remove|rename|set-url` | Manage remotes | |
| `jj git import` | Import Git refs into jj (automatic in colocated repos) | |
| `jj git export` | Export jj bookmarks to Git refs (automatic in colocated repos) | |
| `jj git colocation status|enable|disable` | Inspect or change colocation | |
| `jj git root` | Path of the backing Git repository | |

## Operations and recovery

`jj operation` is aliased to `jj op`.

| Command | Purpose | Key flags |
| --- | --- | --- |
| `jj op log` | Operation history, newest first | `-n`, `--reversed`, `-G`, `-T`, `-d/--op-diff`, `-p`, `--show-changes-in <revs>` |
| `jj op show [op]` | Details and changes of one operation | `-p` |
| `jj op diff` | Compare repository state between two operations | `--from`, `--to`, `-p` |
| `jj undo` | Undo the most recent operation (sequential; repeat to go further back) | |
| `jj redo` | Redo the most recently undone operation | |
| `jj op restore <op>` | Reset the entire repository to the state after `<op>` | `--what repo|remote-tracking` |
| `jj op revert <op>` | Undo the effect of one specific operation, keeping later ones | `--what` |
| `jj op abandon <op>` | Drop old operations from the log | range syntax `..<op>` |
| `jj op integrate` | Integrate an operation created with `--no-integrate-operation` | |

Operation IDs come from `jj op log`. `@` is the current operation and `@-`
its parent, so `jj op restore @-` is one step back.

## Workspaces and sparse checkouts

| Command | Purpose | Key flags |
| --- | --- | --- |
| `jj workspace add <path>` | Second working directory sharing the same repository | `--name`, `-r <revs>`, `-m`, `--sparse-patterns` |
| `jj workspace list` | Workspaces and their working-copy commits | |
| `jj workspace forget [name]` | Stop tracking a workspace directory | |
| `jj workspace rename <new>` | Rename the current workspace | |
| `jj workspace root` | Same as `jj root` | `--workspace` |
| `jj workspace update-stale` | Refresh a workspace whose working copy was rewritten elsewhere | |
| `jj sparse list|set|reset|edit` | Limit which paths are materialized in the working copy | `--add`, `--remove`, `--clear` on `set` |

Workspaces are named `default`, then the name given at creation. `<name>@`
addresses another workspace's working copy in revsets.

## Config, signing, and utilities

| Command | Purpose | Key flags |
| --- | --- | --- |
| `jj config list` | Effective settings | `--include-defaults`, `--include-overridden`, `--user`, `--repo`, `--workspace` |
| `jj config get <name>` | One value | |
| `jj config set <name> <value>` | Write a setting | `--user`, `--repo`, `--workspace`, `--file` |
| `jj config unset <name>` | Remove a setting | scope flags |
| `jj config path` | Locations of config files | scope flags |
| `jj config edit` | Open config in an editor. Avoid in agent sessions. | |
| `jj sign` / `jj unsign` | Sign or unsign revisions | `-r <revs>`, `--key` |
| `jj util gc` | Garbage collect unreachable objects | `--aggressive` |
| `jj util exec -- <cmd>` | Run a command in the repository, useful for aliases | |
| `jj util completion <shell>` | Shell completions | |
| `jj util config-schema` | JSON schema of settings | |
| `jj run` | Run a command across revisions (experimental) | |
| `jj gerrit upload` | Gerrit review upload | see `jj help gerrit upload` |

## Global options

| Option | Effect |
| --- | --- |
| `-R, --repository <path>` | Operate on another repository |
| `--ignore-working-copy` | Do not snapshot the working copy before running |
| `--ignore-immutable` | Allow rewriting immutable commits. Ask before using. |
| `--at-operation <op>` (`--at-op`) | Load the repository as it was after an operation; implies no snapshot |
| `--no-integrate-operation` | Record the operation without making it the current one |
| `--config NAME=VALUE`, `--config-file <path>` | Override settings for one command |
| `--color always|never|auto|debug` | Color control; use `never` for parsing |
| `--no-pager`, `--quiet`, `--debug` | Output control |

## Flags shared across commands

- Diff output: `-s/--summary`, `--stat`, `--types`, `--name-only`, `--git`,
  `--color-words`, `--context <n>`, `-w/--ignore-all-space`,
  `-b/--ignore-space-change`, `--tool <name>`.
- Placement: `-o/--onto`, `-A/--insert-after`, `-B/--insert-before` on `new`,
  `squash`, `split`, `rebase`, `duplicate`, and `revert`.
- Interactivity: `-i/--interactive`, `--tool`, `--editor` open external
  programs and block a non-interactive session. Use `-m`, paths, and filesets.
- Message input: `-m <text>`. `describe --stdin` reads a message from stdin.
