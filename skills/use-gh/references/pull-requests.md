# Pull requests

`N` is a number, URL, or branch. Without it, `gh` uses the PR of the current
branch. Add `-R OWNER/REPO` outside the target clone.

## Create

The head branch must already be on the remote. Git: `git push -u origin
<branch>`. jj: push the bookmark with jj, then pass it as `-H`.

```sh
gh pr create -B main -H <branch> --title "..." --body-file - [--draft] [-r reviewer] [-a @me] [-l label] <<'MD'
...
MD
```

- `--fill` derives title and body from the commits; `--fill-first` uses only
  the first commit. Use them when the project accepts commit messages as PR
  text. `--dry-run` prints what would be created.
- `-T <file>` starts from a template file and cannot be combined with
  `--body` or `--body-file`.
- Success prints the PR URL. A failed run prints a `--recover <file>` hint.

## Inspect

```sh
gh pr status                                    # current branch, own PRs, review requests
gh pr list -L 20 --json number,title,headRefName,author,isDraft,reviewDecision
gh pr view N --json number,title,state,isDraft,baseRefName,headRefName,reviewDecision,mergeStateStatus,mergeable,url
gh pr diff N --name-only                        # or --patch for the full diff
gh pr view N --json reviews --jq '.reviews[] | {author: .author.login, state, body}'
gh api repos/{owner}/{repo}/pulls/N/comments --paginate --jq '.[] | {path, line, body, user: .user.login}'
```

- `list` filters: `-s open|closed|merged|all`, `-B base`, `-H head`,
  `-A author`, `-l label`, `-S "<search query>"`.
- `reviewDecision`: `APPROVED`, `CHANGES_REQUESTED`, `REVIEW_REQUIRED`, or
  empty when no review is required.
- `mergeStateStatus`: `CLEAN` mergeable; `BLOCKED` checks or reviews missing;
  `BEHIND` base moved; `DIRTY` conflicts; `UNSTABLE` non-required checks
  failing; `DRAFT`; `UNKNOWN` means GitHub is still computing, read again.
- Inline review comments live only in the API endpoint above, not in the
  `reviews` or `comments` fields.

## Checks

```sh
gh pr checks N [--required]                      # exit 0 pass, 1 failure, 8 pending
gh pr checks N --watch --fail-fast [-i 10]       # block until done
gh pr checks N --json name,state,bucket,link --jq '.[] | select(.bucket == "fail")'
```

`bucket` is `pass`, `fail`, `pending`, `skipping`, or `cancel`. The run id is
the number after `/runs/` in `link`; read its failure with the runs reference.

## Review, comment, edit

```sh
gh pr review N --approve | --request-changes --body-file - | --comment --body-file -
gh pr comment N --body-file - [--edit-last [--create-if-none]]
gh pr edit N --title "..." | --body-file - | -B base | --add-label l | --add-reviewer r | --add-assignee @me
gh pr ready N [--undo]                           # draft to ready, or back
gh pr close N [-c "comment"] ; gh pr reopen N
gh pr update-branch N [--rebase]                 # bring the head up to date with base
gh pr checkout N [-b name] [--worktree path]     # Git only; in jj fetch and start a change on the bookmark
```

You cannot approve your own PR. A single line comment needs the API:
`gh api -X POST repos/{owner}/{repo}/pulls/N/comments -f body=... -f commit_id=<head sha> -f path=<file> -F line=<n> -f side=RIGHT`.

## Merge

```sh
gh repo view --json squashMergeAllowed,mergeCommitAllowed,rebaseMergeAllowed,deleteBranchOnMerge
gh pr merge N --squash|--merge|--rebase [--delete-branch] [-t "subject" -b "body"] [--match-head-commit SHA]
gh pr merge N --squash --auto                    # merge once checks and reviews pass
gh pr merge N --disable-auto
gh pr view N --json state,mergedAt,mergeCommit   # state MERGED confirms
```

- One strategy flag is required non-interactively and must be allowed by the
  repository. For a merge queue no strategy is needed; `--auto` queues the PR
  once checks pass.
- `--match-head-commit <sha>` refuses to merge if new commits landed after
  the reviewed head.
- `--delete-branch` also deletes the local branch and checks out the base.
  Skip it in a jj repository and delete the bookmark with jj instead.
- `--admin` bypasses branch protection. Use it only on an explicit request.
