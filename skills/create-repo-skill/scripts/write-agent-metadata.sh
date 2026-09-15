#!/usr/bin/env bash
set -euo pipefail
script_dir="$(cd "$(dirname "$0")" && pwd)"
command -v python3 >/dev/null || { printf 'Python 3 with PyYAML is required\n' >&2; exit 1; }
exec python3 -B "$script_dir/authoring.py" metadata "$@"
