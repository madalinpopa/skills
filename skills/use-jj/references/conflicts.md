# Conflicts and divergence

jj never stops a command for a conflict. A rebase or merge that conflicts
succeeds and marks the commit `(conflict)` in `jj log`. Resolve when
convenient. Descendants inherit the conflict until it is resolved, and the
resolution then flows down automatically.

## Contents

- [Find conflicts](#find-conflicts)
- [Marker format](#marker-format)
- [Resolve in place](#resolve-in-place)
- [Resolve with a tool](#resolve-with-a-tool)
- [Resolve by taking one side](#resolve-by-taking-one-side)
- [Divergent changes](#divergent-changes)
- [Conflicted bookmarks](#conflicted-bookmarks)

## Find conflicts

```sh
jj st                                   # lists conflicted files in @
jj log -r 'conflicts()'                 # every conflicted commit
jj resolve --list -r <x>                # conflicted paths in a revision
jj log -r '<x>::' -T 'change_id.short() ++ if(conflict, " CONFLICT", "") ++ "\n"'
```

Resolve the first conflicted commit in a chain. Descendants often resolve
themselves once the root of the problem is fixed.

## Marker format

Check the active style with `jj config get ui.conflict-marker-style`. The
upstream default, `diff`, shows one side as a snapshot and the other as a
diff to apply, so the two sides are easy to compare:

```
<<<<<<< conflict 1 of 1
%%%%%%% diff from: vpxusssl 38d49363 "merge base"
\\\\\\\        to: rtsqusxu 2768b0b9 "commit A"
 apple
-grape
+grapefruit
 orange
+++++++ ysrnknol 7a20f389 "commit B"
APPLE
GRAPE
ORANGE
>>>>>>> conflict 1 of 1 ends
```

Apply the diff to the other snapshot: the resolution here is `APPLE`,
`GRAPEFRUIT`, `ORANGE`. Set `ui.conflict-marker-style` to `snapshot` to see
every side in full, or `git` for two-sided diff3 markers. Markers grow extra
`<` characters when file content itself looks like a marker.

## Resolve in place

Working-copy conflict:

```sh
jj st                                   # see the files
# edit each file: keep the final content, delete every marker line
jj st                                   # "no conflicts" once all markers are gone
```

Conflict in an ancestor `<x>`:

```sh
jj new <x>                              # temporary change on top
# edit files until markers are gone
jj squash                               # resolution moves into x; descendants rebase
jj log -r 'conflicts()'                 # confirm the chain is clean
```

`jj edit <x>` also works and amends directly, but the `new` then `squash`
pattern lets you inspect the resolution before it lands.

## Resolve with a tool

`jj resolve` opens the configured merge tool per file. It blocks a
non-interactive session unless `ui.merge-editor` names a non-interactive
tool. Prefer editing markers by hand. `jj resolve --list` is safe anywhere.

## Resolve by taking one side

```sh
jj restore --from <x>- path/to/file     # take the parent's version
jj restore --from <other-parent> path/to/file
jj restore --from <x> --into @ path/to/file   # copy a known-good version in
```

Restore replaces the conflicted file in `@` with the chosen content. Then
squash into the conflicted commit as above.

## Divergent changes

A change is divergent when two visible commits share a change ID. `jj log`
marks them `(divergent)` and shows `change/0`, `change/1`. Causes: rewriting
a commit in two workspaces, fetching a rewritten commit, `jj new` or `jj
edit` on a hidden commit, or Git commands run in a colocated repository.

```sh
jj log -r 'divergent()'
jj evolog -r <change>/0                 # how each version got here
```

Pick one resolution:

```sh
jj abandon <commit-id>                  # drop the stale version
jj squash --from <commit-a> --into <commit-b>   # merge content into one
jj converge -r <change> --no-interactive # let jj merge the versions
jj metaedit --update-change-id <commit-id>      # keep both as separate changes
```

Use commit IDs rather than the change ID while it is divergent, because the
change ID is ambiguous.

## Conflicted bookmarks

A bookmark shows `name??` when jj cannot pick one target, usually after a
fetch where both the remote and the local bookmark moved.

```sh
jj bookmark list --conflicted
jj log -r 'bookmarks(exact:"main")'     # all candidate targets
jj bookmark move main --to <commit-id> -B
```

For a remote bookmark that is conflicted, `jj git fetch` reconciles it with
the remote. Push refuses conflicted bookmarks until resolved.
