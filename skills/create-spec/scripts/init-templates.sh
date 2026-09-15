#!/usr/bin/env bash
# Copy the bundled SPEC.md, FEATURE.md, and PHASE.md templates into a
# project's docs/templates/ directory. Existing files are never touched.
#
#   init-templates.sh <project-root>
#
# Exit 0 when every template exists afterwards, 2 on usage or a missing root.
set -euo pipefail

usage() {
  echo "usage: init-templates.sh <project-root>" >&2
  exit 2
}

[ $# -eq 1 ] || usage
root="$1"
[ -d "$root" ] || { echo "init-templates: no such directory: $root" >&2; exit 2; }

assets="$(cd "$(dirname "$0")/.." && pwd)/assets/templates"
target="$root/docs/templates"
mkdir -p "$target"

created=0
for name in SPEC.md FEATURE.md PHASE.md; do
  if [ -e "$target/$name" ]; then
    echo "init-templates: kept    docs/templates/$name"
  else
    cp "$assets/$name" "$target/$name"
    echo "init-templates: created docs/templates/$name"
    created=$((created + 1))
  fi
done
echo "init-templates: $created created in $target"
