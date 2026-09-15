# Actions runs

`<id>` is the run's `databaseId`. Add `-R OWNER/REPO` outside the target
clone. Rerun, cancel, dispatch, and delete change shared state; task
authorization applies.

## Find runs

```sh
gh run list -L 5 --json databaseId,displayTitle,workflowName,headBranch,event,status,conclusion,url
gh workflow list                                  # names, states, ids
```

Filters: `-b branch`, `-w <workflow name or file>`, `-c SHA`,
`-e push|pull_request|workflow_dispatch`, `-u user`, `-s <status or
conclusion>` such as `in_progress`, `failure`, `success`. Runs behind a PR's
checks come from `gh pr checks N` in the pull requests reference.

## Status and logs

```sh
gh run view <id> --json status,conclusion,jobs --jq '.jobs[] | {name, status, conclusion, id: .databaseId}'
gh run view <id> --log-failed [-j <job-id>]      # only failed steps
gh run view <id> -v                               # job and step list
gh run watch <id> --exit-status --compact [-i 10] # block until done; non-zero on failure
```

- `status` is `queued`, `in_progress`, or `completed`; `conclusion` is set
  only when completed: `success`, `failure`, `cancelled`, `skipped`,
  `timed_out`, `action_required`, `startup_failure`.
- `--log` prints every step of every job. Pipe it through `rg` or
  `tail -n 200`; never print it whole.
- `run watch` does not work with fine-grained personal access tokens.

## Act

```sh
gh run rerun <id> --failed [-d]                   # failed jobs only; -d adds debug logging
gh run rerun <id> -j <job-id>
gh run cancel <id>
gh workflow run <name.yml> [-r ref] [-f key=value] [--json < inputs.json]
gh run list -w <name.yml> -L 1                    # find the dispatched run after a few seconds
gh run download <id> [-n artifact] [-p 'glob'] -D <dir>
```

`workflow run` needs a `workflow_dispatch` trigger in the workflow file.
`gh run delete <id>` is permanent; use it only on explicit request.

## Triage a failing PR

1. `gh pr checks N --json name,bucket,link --jq '.[] | select(.bucket == "fail")'`
   gives the failing checks and run links.
2. `gh run view <id> --log-failed | tail -n 100` shows the failing step.
3. Fix locally, push, then `gh pr checks N --watch --fail-fast`.

Rerun with `--failed` only for flaky infrastructure; a rerun never fixes a
real failure.
