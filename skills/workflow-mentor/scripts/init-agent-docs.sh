#!/usr/bin/env bash
# Copy the bundled AGENTS.md and CLAUDE.md into a project root when missing.
# Existing files are never touched.
#
#   init-agent-docs.sh <project-root>
#
# Exit 0 when both files exist afterwards, 2 on usage or a missing root.
set -euo pipefail

usage() {
  echo "usage: init-agent-docs.sh <project-root>" >&2
  exit 2
}

[ $# -eq 1 ] || usage
root="$1"
[ -d "$root" ] || { echo "init-agent-docs: no such directory: $root" >&2; exit 2; }

assets="$(cd "$(dirname "$0")/.." && pwd)/assets"

created=0
for name in AGENTS.md CLAUDE.md; do
  if [ -e "$root/$name" ]; then
    echo "init-agent-docs: kept    $name"
  else
    cp "$assets/$name" "$root/$name"
    echo "init-agent-docs: created $name"
    created=$((created + 1))
  fi
done

if ! grep -q '@AGENTS.md' "$root/CLAUDE.md"; then
  echo "init-agent-docs: warning CLAUDE.md does not import @AGENTS.md"
fi
echo "init-agent-docs: $created created in $root"
