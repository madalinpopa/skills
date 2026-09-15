#!/usr/bin/env bash
# Check or uncheck one task in a GitHub issue body and leave the rest untouched.
#
#   issue-task.sh <issue> --check|--uncheck "<task text>" [-R OWNER/REPO]
#   issue-task.sh --file <body.md> --check|--uncheck "<task text>"   # offline: prints the new body
#
# <task text> is a plain-text fragment that appears in exactly one task line.
# Exit 0 done or already in that state, 2 usage, 3 no or several matching tasks,
# 1 gh failure or the edit did not stick.
set -euo pipefail

usage() {
  cat >&2 <<'USAGE'
usage: issue-task.sh <issue> (--check|--uncheck) "<task text>" [-R OWNER/REPO]
       issue-task.sh --file <body.md> (--check|--uncheck) "<task text>"
USAGE
  exit 2
}

issue="" file="" mode="" text="" repo=()
while [ $# -gt 0 ]; do
  case "$1" in
    --check|--uncheck)
      [ -z "$mode" ] && [ $# -ge 2 ] || usage
      mode="${1#--}"; text="$2"; shift 2 ;;
    --file) [ $# -ge 2 ] || usage; file="$2"; shift 2 ;;
    -R) [ $# -ge 2 ] || usage; repo=(-R "$2"); shift 2 ;;
    -*) usage ;;
    *) [ -z "$issue" ] || usage; issue="$1"; shift ;;
  esac
done
[ -n "$mode" ] && [ -n "$text" ] || usage
if [ -n "$file" ]; then
  [ -z "$issue" ] || usage
  [ -f "$file" ] || { echo "issue-task: no such file: $file" >&2; exit 2; }
else
  [ -n "$issue" ] || usage
fi

transform() {
  awk -v text="$text" -v mode="$mode" '
    {
      lines[NR] = $0
      if (match($0, /^[ \t]*[-*] \[[ xX]\] /) && index($0, text) > 0) { hits++; hit[hits] = NR }
    }
    END {
      if (hits == 0) { print "issue-task: no task line contains: " text > "/dev/stderr"; exit 3 }
      if (hits > 1) {
        print "issue-task: " hits " task lines contain: " text > "/dev/stderr"
        for (i = 1; i <= hits; i++) print "  " lines[hit[i]] > "/dev/stderr"
        exit 3
      }
      n = hit[1]
      if (mode == "check") { if (!sub(/\[ \]/, "[x]", lines[n])) print "issue-task: already checked" > "/dev/stderr" }
      else { if (!sub(/\[[xX]\]/, "[ ]", lines[n])) print "issue-task: already unchecked" > "/dev/stderr" }
      for (i = 1; i <= NR; i++) print lines[i]
      print "issue-task: " lines[n] > "/dev/stderr"
    }'
}

if [ -n "$file" ]; then
  transform < "$file"
  exit 0
fi

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
gh issue view "$issue" ${repo[@]+"${repo[@]}"} --json body --jq .body > "$tmp/body.md"
transform < "$tmp/body.md" > "$tmp/new.md"
if cmp -s "$tmp/body.md" "$tmp/new.md"; then
  echo "issue-task: #$issue unchanged"
  exit 0
fi
gh issue edit "$issue" ${repo[@]+"${repo[@]}"} --body-file "$tmp/new.md" > /dev/null
gh issue view "$issue" ${repo[@]+"${repo[@]}"} --json body --jq .body > "$tmp/after.md"
transform < "$tmp/after.md" > "$tmp/again.md" 2> /dev/null
if ! cmp -s "$tmp/after.md" "$tmp/again.md"; then
  echo "issue-task: edit did not stick on #$issue; re-run" >&2
  exit 1
fi
echo "issue-task: updated #$issue"
