#!/usr/bin/env bash
set -euo pipefail

created=0
kept=0
fail() {
  printf 'init-project-docs: %s\n' "$*" >&2
  printf 'init-project-docs: stopped; %s created, %s kept; inspect reported paths before retrying\n' "$created" "$kept" >&2
  exit 1
}

if [[ $# -ne 1 ]]; then
  printf 'usage: init-project-docs.sh <existing-project-root>\n' >&2
  exit 2
fi
root="$1"
[[ ! -L "$root" ]] || fail "symlink project root: $root"
[[ -d "$root" ]] || fail "no such project directory: $root"
root="$(cd -- "$root" && pwd -P)"
assets="$(cd -- "$(dirname -- "$0")/../assets" && pwd -P)"

directories=(docs docs/conventions docs/specs docs/decisions docs/features)
files=(
  AGENTS.md CLAUDE.md
  docs/conventions/CODE.md docs/conventions/COMMITS.md
  docs/conventions/GOLANG.md docs/conventions/TDD.md
  docs/specs/MISSION.md docs/specs/FEATURES.md
  docs/specs/ROADMAP.md docs/specs/STACK.md
  docs/decisions/README.md docs/features/README.md
)

# Preflight every destination before creating any project documents.
for directory in "${directories[@]}"; do
  destination="$root/$directory"
  [[ ! -L "$destination" ]] || fail "symlink directory: $directory"
  if [[ -e "$destination" && ! -d "$destination" ]]; then
    fail "expected directory: $directory"
  fi
done
for file in "${files[@]}"; do
  destination="$root/$file"
  [[ ! -L "$destination" ]] || fail "symlink document: $file"
  if [[ -e "$destination" && ! -f "$destination" ]]; then
    fail "expected regular file: $file"
  fi
  [[ -f "$assets/$file" && ! -L "$assets/$file" ]] || fail "missing regular template: $file"
done

for directory in "${directories[@]}"; do
  mkdir -p -- "$root/$directory" || fail "could not create directory: $directory"
done
for file in "${files[@]}"; do
  destination="$root/$file"
  if [[ -f "$destination" ]]; then
    kept=$((kept + 1))
    printf 'init-project-docs: kept    %s\n' "$file"
  else
    # Refuse replacement even if a destination appeared after preflight.
    (set -o noclobber; cat -- "$assets/$file" > "$destination") || fail "could not copy $file; a partial file may remain"
    created=$((created + 1))
    printf 'init-project-docs: created %s\n' "$file"
  fi
done
printf 'init-project-docs: %s created, %s kept in %s\n' "$created" "$kept" "$root"
