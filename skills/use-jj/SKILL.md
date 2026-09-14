---
name: use-jj
description: Works in repositories managed by Jujutsu (jj) version control. Covers inspecting history with revsets, creating and rewriting changes, splitting and squashing, rebasing stacks, bookmarks, conflicts, Git remotes in colocated repos, and undo through the operation log. Use whenever a repository has a .jj directory, Git reports a detached HEAD next to a .jj directory, or the user mentions jj, jujutsu, changes, bookmarks, revsets, or the op log, even when they phrase the task in Git words such as commit, branch, stash, amend, or rebase. Not for plain Git repositories without .jj.
compatibility: Requires the jj CLI on PATH. Colocated repositories also use Git.
status: published
tags: [jj, jujutsu, vcs, git]
---

# Use jj

Use this skill when a repository is managed by Jujutsu. Prefer `jj` over
`git` for every mutation in such a repository, and read the references below
before running a command family you have not used in this session.

## Detect a jj repository

Run `jj root` from the working directory. Success means the directory is
inside a jj workspace. Other signs: a `.jj` directory at the repository root,
or `git status` reporting `HEAD detached` in a repository that also has `.jj`.
That detached state is normal in a colocated repository: jj moves Git's HEAD
to the parent of the working-copy change. Never "fix" it with `git checkout`.

If `jj root` fails, this skill does not apply. Use Git as usual.

## Mental model

- The working copy is a commit, called `@`. Every jj command snapshots the
  working directory into `@` first. There is no staging area and no `add`.
- Each commit has a change ID (stable, letters `k`-`z`) and a commit ID
  (hash, changes on every rewrite). Address commits by change ID.
- Rewriting history is normal. Descendants rebase automatically. Conflicts are
  stored inside commits and never block a command.
- Bookmarks are named pointers, like Git branches, but they do not move on
  their own. Move them with `jj bookmark move` or `jj bookmark set`.
- Every command is an operation. `jj undo` reverts the last one, and
  `jj op restore` returns the whole repository to any earlier point.
- Commits reachable from `trunk()`, tags, and untracked remote bookmarks are
  immutable by default. jj refuses to rewrite them without `--ignore-immutable`.

## Rules for agent sessions

- Look before you rewrite: run `jj log` and `jj st`, then run the command,
  then run `jj log` again and report what changed with change IDs.
- Avoid commands that open an editor or a diff tool. Pass `-m` to `describe`,
  `commit`, `new`, `split`, and `squash`. Pass paths or filesets instead of
  `-i`. Resolve conflicts by editing the markers in the file.
- Never push, and never run any command that touches a remote, unless the
  user asks for it in this task. Follow the project's own rules about
  creating commits; `jj commit`, `jj describe`, and `jj new` write history.
- Do not rewrite immutable commits. If a command fails with an immutable
  error, stop and ask, rather than adding `--ignore-immutable`.
- In a colocated repository, use `git` only for read-only commands such as
  `git log` or `git show`. If a mutating `git` command was run by mistake,
  recover with `jj undo` or `jj op restore`.
- When something goes wrong, run `jj op log` first. Almost every state is
  recoverable; do not delete `.jj` or re-clone.

## Quick map

| Task | Command |
| --- | --- |
| See where you are | `jj st`, `jj log` |
| Start a change on top of main | `jj new main -m "message"` |
| Set or change a description | `jj describe -m "message"` |
| Finish this change and start the next | `jj commit -m "message"` or `jj new` |
| Fold working copy into its parent | `jj squash` |
| Move some files into the parent | `jj squash path/one path/two` |
| Split the working copy | `jj split path/for/first/commit` |
| Move a stack onto main | `jj rebase -s <change> -o main` |
| Point a bookmark at the working-copy parent | `jj bookmark set <name> -r @-` |
| Fetch and push | `jj git fetch`, `jj git push -b <name>` |
| Drop a change | `jj abandon <change>` |
| Undo the last command | `jj undo` |

## References

Read the reference that matches the task before running unfamiliar commands.
All paths are relative to this file.

- [references/commands.md](references/commands.md): every command and
  subcommand with the flags that matter. Read when you need an exact flag.
- [references/revsets.md](references/revsets.md): revision selection language,
  operators, functions, string and date patterns, aliases. Read before writing
  any `-r` expression beyond `@`, `@-`, or a bookmark name.
- [references/workflows.md](references/workflows.md): step-by-step use cases
  such as stacked changes, editing older commits, absorbing fixes, reviewing,
  bisecting, and cleaning history. Read when planning a multi-step task.
- [references/git-interop.md](references/git-interop.md): colocated repos,
  bookmarks and remotes, fetch and push rules, pull-request flows, and the
  Git-to-jj command table. Read for any remote or Git-flavored request.
- [references/conflicts.md](references/conflicts.md): conflict markers,
  resolving in place, divergent changes, conflicted bookmarks. Read when
  `jj log` shows `(conflict)`, `(divergent)`, or `??`.
- [references/operations.md](references/operations.md): operation log, undo,
  restore, revert, evolog, recovering hidden commits, and workspaces. Read
  before any recovery.
- [references/filesets-and-templates.md](references/filesets-and-templates.md):
  path selection language and output templates for scripted or precise
  output. Read when a command needs a path pattern or custom `-T` output.

## Reporting

When you finish, summarize the resulting graph in words: which changes exist,
their change IDs and descriptions, which bookmark points where, and whether
anything is conflicted or unpushed. Mention the operation ID from `jj op log`
when you performed a rewrite, so the user can undo it later.
