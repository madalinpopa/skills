# Issues

`N` is an issue number or URL. Add `-R OWNER/REPO` outside the target clone.

## Read and search

```sh
gh issue list -L 20 --json number,title,labels,assignees,updatedAt
gh issue view N --json number,title,body,state,stateReason,labels,comments,closedByPullRequestsReferences
gh search issues "<query>" --repo OWNER/REPO -L 20 --json number,title,state,url
```

- `list` filters: `-s open|closed|all`, `-l label`, `-a @me`, `-A author`,
  `-m milestone`, `--type name`. `-S "<query>"` takes GitHub search syntax
  such as `is:open label:bug no:assignee sort:updated-desc`.
- Comments: `--jq '.comments[] | {author: .author.login, body}'`. Label names:
  `--jq '.labels[].name'`. Add `--comments` only for plain-text reading.
- `search issues` spans repositories; `--include-prs` adds pull requests.
- Relationships are fields: `parent`, `subIssues`, `blockedBy`, `blocking`.

## Create

```sh
gh issue create --title "..." --body-file - [-l label] [-a @me] [-m milestone] [--type name] [--parent N] [--blocked-by N] <<'MD'
...
MD
```

`-T <template name>` starts from a repository issue template; it cannot be
combined with `--body`. Success prints the issue URL. A failed run prints a
`--recover <file>` hint; reuse it instead of retyping.

## Edit and triage

```sh
gh issue edit N --add-label bug --remove-label triage --add-assignee @me --milestone "v1.2"
gh issue edit 12 15 --add-label needs-info          # several at once
gh issue edit N --parent M | --remove-parent | --add-sub-issue M | --add-blocked-by M
gh issue comment N --body-file -                    # --edit-last updates your last comment
gh issue close N -r completed|"not planned"|duplicate [-c "comment"] [--duplicate-of M]
gh issue reopen N [-c "comment"]
```

`--edit-last --create-if-none` keeps one status comment per issue instead of
adding a new one on every run. `gh issue delete N --yes` removes the issue
permanently; use it only on explicit request.

## Link to work

A PR body line `Closes #N` or `Fixes #N` closes the issue on merge. Confirm the
link through the issue's `closedByPullRequestsReferences` field.
