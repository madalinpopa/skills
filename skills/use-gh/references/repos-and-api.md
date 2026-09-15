# Repositories, labels, releases, and the API

Add `-R OWNER/REPO` or a positional `OWNER/REPO` outside the target clone.
`delete`, `archive`, visibility changes, and `--force` are destructive or hard
to reverse; run them only on explicit request.

## Repositories

```sh
gh repo view [OWNER/REPO] --json nameWithOwner,defaultBranchRef,visibility,isFork,parent,url
gh repo list [owner] -L 30 --json name,visibility,isFork,updatedAt [--source] [--no-archived] [-l go] [--topic t]
gh repo clone OWNER/REPO [dir] [-- --depth 1]    # forks get an upstream remote
gh repo create NAME --private|--public|--internal [-d "desc"] [--add-readme] [-g Go] [-l mit] [-c]
gh repo create NAME --private --source . --push -r origin   # publish a local repo with commits
gh repo fork [OWNER/REPO] [--clone] ; gh repo sync [-b branch]
gh repo edit --default-branch main | --delete-branch-on-merge | --enable-auto-merge | --add-topic t
gh search repos "<query>" [--owner o] [--language go] -L 10 --json fullName,description,url
```

Merge policy fields on `repo view`: `squashMergeAllowed`,
`mergeCommitAllowed`, `rebaseMergeAllowed`, `deleteBranchOnMerge`,
`viewerPermission`. Changing `--visibility` also needs
`--accept-visibility-change-consequences`.

## Labels

```sh
gh label list --json name,color,description
gh label create NAME -c HEX -d "desc" [--force]   # --force updates an existing label
gh label clone SOURCE/REPO
```

## Releases

```sh
gh release list -L 10 --json tagName,isLatest,isDraft,publishedAt
gh release view <tag> --json tagName,name,body,assets
gh release create <tag> [files...] -t "title" --notes-file - | --generate-notes [--notes-start-tag prev] [--target sha] [--draft] [--prerelease] [--verify-tag]
gh release upload <tag> files... [--clobber]
```

`release create` creates the tag when it does not exist; `--verify-tag` aborts
unless the tag is already on the remote. `gh release delete <tag> --yes
[--cleanup-tag]` is permanent.

## gh api

Use for anything without a subcommand. `{owner}`, `{repo}`, and `{branch}`
in the endpoint fill in from the current clone or `GH_REPO`.

```sh
gh api repos/{owner}/{repo}/pulls/N/comments --paginate --jq '.[] | {path, line, body}'
gh api -X PATCH repos/{owner}/{repo}/issues/N -f state=closed -f state_reason=not_planned
gh api -X POST repos/{owner}/{repo}/issues --input body.json
gh api graphql -f query='query($owner:String!,$name:String!){repository(owner:$owner,name:$name){id}}' -F owner='{owner}' -F name='{repo}' --jq .data.repository.id
gh api rate_limit --jq .rate
```

- `-f` sends a string; `-F` sends typed values (numbers, booleans, `@file`,
  `@-` for stdin, placeholders). Any field switches the method to POST unless
  `-X GET` is given. `--input file.json` sends a nested body.
- `--paginate` follows every page; add `--slurp` for one combined array.
  `--cache 5m` avoids repeated identical reads. `-i` shows status and headers.
- Errors exit `1` with the JSON message on stderr. A `404` on a write usually
  means a missing token scope; `gh auth status` lists the scopes.
- Common needs: review threads and their resolution (GraphQL
  `reviewThreads`), check runs (`repos/{owner}/{repo}/commits/SHA/check-runs`),
  rulesets (`gh ruleset list`).
