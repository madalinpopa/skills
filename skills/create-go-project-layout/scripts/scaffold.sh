#!/usr/bin/env bash
set -euo pipefail

SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"
ASSETS_DIR="$SKILL_DIR/assets"
MANIFEST_NAME=".scaffold"

usage() {
  cat <<'USAGE'
Usage: scaffold.sh <layout> --name <slug> --module <path> [options]
       scaffold.sh --list

Creates a new project from one of the layouts under assets/.

Arguments:
  <layout>            layout name, a directory under assets/
  --name <slug>       project name, replaces {{PROJECT_NAME}}
  --module <path>     Go module path, replaces {{MODULE_PATH}}

Options:
  --dir <target>      target directory; default ./<slug>; must be empty or absent
  --go <major.minor>  Go version, replaces {{GO_VERSION}}; default from go version
  --skip-deps         copy and replace only; skip go mod init, go get, build, vet
  --list              print the available layouts and exit
  -h, --help          show this help
USAGE
}

fail() {
  echo "scaffold: $*" >&2
  exit 1
}

list_layouts() {
  local dir
  for dir in "$ASSETS_DIR"/*/; do
    [ -f "$dir$MANIFEST_NAME" ] || continue
    basename "$dir"
  done
}

layout=""
name=""
module=""
target=""
go_version=""
skip_deps=false

while [ $# -gt 0 ]; do
  case "$1" in
    --name) name="${2:-}"; shift 2 ;;
    --module) module="${2:-}"; shift 2 ;;
    --dir) target="${2:-}"; shift 2 ;;
    --go) go_version="${2:-}"; shift 2 ;;
    --skip-deps) skip_deps=true; shift ;;
    --list) list_layouts; exit 0 ;;
    -h|--help) usage; exit 0 ;;
    -*) fail "unknown option: $1" ;;
    *)
      [ -z "$layout" ] || fail "unexpected argument: $1"
      layout="$1"; shift ;;
  esac
done

[ -n "$layout" ] || { usage >&2; exit 1; }
[ -n "$name" ] || fail "--name is required"
[ -n "$module" ] || fail "--module is required"

case "$name" in
  *[!a-z0-9-]*|'') fail "--name must be a lowercase slug: letters, digits, hyphens" ;;
esac

layout_dir="$ASSETS_DIR/$layout"
manifest="$layout_dir/$MANIFEST_NAME"
[ -d "$layout_dir" ] || fail "unknown layout '$layout'; run with --list"
[ -f "$manifest" ] || fail "layout '$layout' has no $MANIFEST_NAME manifest"

if [ "$skip_deps" = false ]; then
  command -v go >/dev/null || fail "go is required; pass --skip-deps to copy files only"
fi

if [ -z "$go_version" ]; then
  if command -v go >/dev/null; then
    if [[ "$(go version)" =~ go([0-9]+\.[0-9]+) ]]; then
      go_version="${BASH_REMATCH[1]}"
    fi
  fi
  [ -n "$go_version" ] || fail "--go is required when go is not installed"
fi

[ -n "$target" ] || target="./$name"
if [ -e "$target" ]; then
  [ -d "$target" ] || fail "target '$target' exists and is not a directory"
  [ -z "$(ls -A "$target")" ] || fail "target '$target' is not empty"
fi

go_dir="."
requires=()
tools=()
while read -r key value; do
  case "$key" in
    ''|'#'*) ;;
    go_dir) go_dir="$value" ;;
    require) requires+=("$value") ;;
    tool) tools+=("$value") ;;
    *) fail "unknown manifest key '$key' in $manifest" ;;
  esac
done < "$manifest"

echo "==> Copying layout '$layout' to $target"
mkdir -p "$target"
cp -R "$layout_dir/." "$target/"
rm -f "$target/$MANIFEST_NAME"

echo "==> Replacing placeholders"
replace_placeholders() {
  local file="$1" content
  content="$(cat "$file")"
  content="${content//\{\{PROJECT_NAME\}\}/$name}"
  content="${content//\{\{MODULE_PATH\}\}/$module}"
  content="${content//\{\{GO_VERSION\}\}/$go_version}"
  printf '%s\n' "$content" > "$file"
}
while IFS= read -r file; do
  replace_placeholders "$file"
done < <(grep -rIl -e '{{' "$target")

if leftover="$(grep -rIn -e '{{[A-Z_][A-Z_]*}}' "$target")"; then
  echo "$leftover" >&2
  fail "unreplaced placeholders remain"
fi

if [ "$skip_deps" = false ]; then
  echo "==> Initialising Go module in $target/$go_dir"
  (
    cd "$target/$go_dir"
    go mod init "$module"
    if [ ${#requires[@]} -gt 0 ]; then
      go get "${requires[@]}"
    fi
    if [ ${#tools[@]} -gt 0 ]; then
      go get -tool "${tools[@]}"
    fi
    go mod tidy
    echo "==> Building and vetting"
    go build ./...
    go vet ./...
  )
fi

file_count="$(find "$target" -type f | wc -l | tr -d ' ')"
cat <<SUMMARY

Scaffold complete
  layout:   $layout
  target:   $target
  files:    $file_count
  module:   $module
  go dir:   $go_dir
  deps:     $([ "$skip_deps" = true ] && echo skipped || echo "installed, build and vet passed")
SUMMARY
