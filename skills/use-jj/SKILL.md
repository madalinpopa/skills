---
name: use-jj
description: Uses Jujutsu (jj) for repository history, changes, bookmarks, remotes, and recovery. Apply when .jj exists or the user requests jj, including Git-worded tasks in jj repositories. Not for plain Git repositories without .jj.
compatibility: Requires the jj CLI on PATH. Colocated repositories also use Git.
status: published
tags: [jj, jujutsu, vcs, git]
---

# Use jj

Use `jj` for repository mutations and Git only for compatible read-only tools.
Answer explanation or review requests without creating or rewriting changes.

## Orient once

Use known workspace context; otherwise run `jj root`. A confirmed non-jj
repository uses Git normally. Missing CLI, access errors, or damaged `.jj`
state require diagnosis, not a fallback to Git mutations. Detached Git HEAD
beside `.jj` is expected; do not repair it with checkout.

Inspect `jj st` and `jj log -n 10` together before changing history. Read the
relevant diff before selecting files or revisions. Reuse this context until
edits, another actor, or a command changes it. Expand the graph only as needed.

## Essentials

- `@` is the working-copy commit; `@-` selects its parents. Most commands
  snapshot edits automatically; there is no staging step. Change IDs survive
  rewrites; commit hashes change. Use change IDs for current revisions and
  commit hashes when distinguishing historical or divergent versions.
- Rewrites rebase descendants and can leave stored conflicts. Bookmarks follow
  rewrites of their commit, but do not advance to newly created children.
- Preserve unrelated edits. Select explicit revisions/files; confirm the target
  before restore, abandon, or broad rewrites. Follow project commit rules and
  existing user authorization. Remote work must be within the requested scope;
  do not repeat permission questions already answered.
- Avoid editors: use `-m` for messages, paths/filesets for selection, and
  `squash -u` only when retaining the destination message is intended. Do not
  bypass immutable-history or push checks. Inspect failures before retrying.

## Choose and run

| Task | Starting command |
| --- | --- |
| Inspect changes | `jj diff --stat`, then the relevant paths |
| Describe / finish | `jj describe -m "message"` / `jj commit -m "message"` |
| Start from a known base | `jj new <base> -m "message"` |
| Split by path | `jj split -m "message" <paths>` |
| Move a stack | `jj rebase -s <change> -o <base>` |
| Point at the finished change | `jj bookmark set <name> -r @-` |

For unfamiliar flags, use `jj <command> -h`; expand to `--help` when needed.
Installed help takes precedence over examples. Read only the relevant section
of these references; familiar commands do not require another catalog read:

| Need | Reference |
| --- | --- |
| Command choice or flags | [Commands](references/commands.md) |
| Revision expression | [Revsets](references/revsets.md) |
| Multi-step recipe | [Workflows](references/workflows.md) |
| Fetch, push, PR, or Git translation | [Git interoperation](references/git-interop.md) |
| Conflict or divergence | [Conflicts](references/conflicts.md) |
| Recovery or workspace management | [Operations](references/operations.md) |
| Path patterns or custom output | [Filesets and templates](references/filesets-and-templates.md) |

Locate headings with `rg -n '^##' <reference>` when the section is unknown.
Batch independent inspections; run mutations sequentially and inspect their
results before dependent steps.

## Verify and report

After a logical change, check status and the affected graph/diff once. Resolve
or report conflicts. For recovery, inspect `jj op log` before selecting an
operation; preserve later work. Local undo does not undo a remote push.
Preserve `.jj`; do not delete or re-clone the workspace to recover.

Summarize the outcome, relevant change IDs/bookmarks, and remaining conflicts
or unpublished work. Include a recovery operation ID after history rewrites.
Keep full graphs, patches, and successful logs out of routine reports.
