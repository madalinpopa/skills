---
name: use-gh
description: Uses the GitHub CLI (gh) for issues, pull requests, reviews, merges, repositories, and Actions run status. Apply whenever a task touches GitHub from the terminal, including opening a PR after a push or triaging failing checks. Not for local Git or jj history operations.
compatibility: Requires the gh CLI on PATH, authenticated for the target host. Local history stays in Git or jj.
status: published
tags: [gh, github, cli, pull-requests, issues, actions]
---

# Use gh

Use `gh` for everything that lives on GitHub: issues, pull requests, reviews,
merges, repositories, and Actions runs. Keep local history in Git or jj.

## Essentials

- Run non-interactively. Pass every required value as a flag; `gh` prompts
  for a missing one and a prompt hangs the session. Never use `--web`,
  `--editor`, or `gh auth login`. Prefix a command with
  `GH_PROMPT_DISABLED=1` when it might still prompt, so it fails instead.
- Keep output small. Use `--json <fields> --jq <expr>` on `view`, `list`,
  and `status`; cap lists with `-L`; prefer `--log-failed` over `--log`.
  Running a command with `--json` and no field list prints the field names.
- Target the current clone by default. Outside a clone, or for another
  repository, pass `-R OWNER/REPO` on every command.
- Write Markdown bodies through stdin with `--body-file -` and a quoted
  heredoc. Avoid `-b` for multi-line or quoted text.
- Creating, editing, closing, commenting, merging, rerunning, and releasing
  change shared state. Do them only within the authorization the task already
  has; reuse approval already given rather than asking again. Follow project
  rules for titles, bodies, and attribution. Use `--admin`, `--force`,
  `delete`, or `--yes` only when the user asked for that exact action.
- Exit code `4` means authentication is needed: run `gh auth status`, report
  it, and stop; the user runs `gh auth login` themselves. Exit `1` is an
  error: read stderr before changing the command; do not retry it unchanged.

## Choose and run

| Task | Starting command |
| --- | --- |
| Repository and default branch | `gh repo view --json nameWithOwner,defaultBranchRef --jq '.nameWithOwner, .defaultBranchRef.name'` |
| List / read issues | `gh issue list -L 20 --json number,title,labels` / `gh issue view N --json title,body,state` |
| Open a PR from a pushed branch | `gh pr create -B <base> -H <head> --title "..." --body-file -` |
| PR state | `gh pr view [N] --json number,state,isDraft,reviewDecision,mergeStateStatus,url` |
| CI checks for a PR | `gh pr checks [N]`; add `--watch --fail-fast` to wait |
| Merge | `gh pr merge N --squash\|--merge\|--rebase [--delete-branch] [--auto]` |
| Recent runs / failure log | `gh run list -b <branch> -L 5` / `gh run view <id> --log-failed` |
| No subcommand exists | `gh api <endpoint> [-X METHOD] [-f k=v] [--paginate] --jq <expr>` |

Use `gh <command> <subcommand> --help` for unfamiliar flags. Read only the
reference for the area in hand; familiar commands need no reference:

| Area | Reference |
| --- | --- |
| Issues: search, create, triage, close | [Issues](references/issues.md) |
| Pull requests: create, inspect, review, checks, merge | [Pull requests](references/pull-requests.md) |
| Actions: run status, logs, rerun, cancel, dispatch, artifacts | [Runs](references/runs.md) |
| Repositories, labels, releases, and the API fallback | [Repos and API](references/repos-and-api.md) |

## Verify and report

After a write, confirm with one read: the URL printed by `create`, `pr view`
after a merge, `issue view --json state` after a close. Report the outcome,
the number or URL, and anything still pending such as unfinished checks or an
enabled auto-merge. Keep JSON dumps, full logs, and successful command output
out of the report.
