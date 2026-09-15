#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd)"
fail() { printf 'init-skill: %s\n' "$*" >&2; exit 1; }
usage() {
  printf '%s\n' 'Usage: init-skill.sh <name> --repo <checkout> --description <text> [--status published|draft] [--resources scripts,references,assets]'
}

if [[ "${1:-}" == --help || "${1:-}" == -h ]]; then usage; exit 0; fi
[[ $# -gt 0 ]] || { usage >&2; exit 1; }
skill_name="$1"
shift
repo=""
description=""
status=published
resources=""
while [[ $# -gt 0 ]]; do
  [[ $# -ge 2 ]] || fail "missing value for $1"
  case "$1" in
    --repo) repo="$2" ;;
    --description) description="$2" ;;
    --status) status="$2" ;;
    --resources) resources="$2" ;;
    *) fail "unknown option: $1" ;;
  esac
  shift 2
done

[[ "$skill_name" =~ ^[a-z0-9]+(-[a-z0-9]+)*$ && ${#skill_name} -le 64 ]] || fail 'invalid skill name'
[[ "$status" == published || "$status" == draft ]] || fail 'status must be published or draft'
[[ -n "$repo" && -f "$repo/AGENTS.md" && -f "$repo/docs/SPEC.md" && -d "$repo/skills" ]] || fail 'repo must contain AGENTS.md, docs/SPEC.md, and skills/'
[[ ! -L "$repo/skills" ]] || fail 'skills/ must not be a symlink'
repo="$(cd "$repo" && pwd -P)"
target="$repo/skills/$skill_name"
[[ ! -e "$target" && ! -L "$target" ]] || fail "skill already exists: $target"
resource_dirs=()
if [[ -n "$resources" ]]; then
  case "$resources" in ,*|*,|*,,*) fail 'invalid resources list' ;; esac
  IFS=, read -r -a resource_dirs <<< "$resources"
  for resource in "${resource_dirs[@]}"; do
    case "$resource" in scripts|references|assets) ;; *) fail "unknown resource: $resource" ;; esac
  done
fi
command -v python3 >/dev/null || fail 'Python 3 with PyYAML is required'
# Serialize data before creating anything; user strings never become shell code.
frontmatter="$(python3 -B "$script_dir/authoring.py" frontmatter "$skill_name" "$description" "$status")"
mkdir "$target"
printf '%s\n\n# %s\n\n[TODO: Add task-specific instructions.]\n' "$frontmatter" "$skill_name" > "$target/SKILL.md"
if [[ ${#resource_dirs[@]} -gt 0 ]]; then
  for resource in "${resource_dirs[@]}"; do mkdir -p "$target/$resource"; done
fi
printf 'Created %s (status: %s). Complete instructions, then validate.\n' "$target" "$status"
