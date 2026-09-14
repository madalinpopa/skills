#!/usr/bin/env bash
set -euo pipefail

SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TEMPLATE_DIR="$SKILL_DIR/assets/module"

usage() {
  cat <<'USAGE'
Usage: scaffold-module.sh <name> --entity <singular> [options]

Adds a module to a project created by create-go-project-layout.

Arguments:
  <name>              module name: lowercase letters and digits, e.g. brand
  --entity <singular> first entity, lowercase letters and digits, e.g. profile

Options:
  --project <dir>     project root, the directory holding api/; default .
  --plural <plural>   plural of the entity; default entity plus "s"
  -h, --help          show this help
USAGE
}

fail() {
  echo "scaffold-module: $*" >&2
  exit 1
}

pascal() {
  local first rest
  first="$(printf '%s' "${1:0:1}" | tr '[:lower:]' '[:upper:]')"
  rest="${1:1}"
  printf '%s%s' "$first" "$rest"
}

name=""
entity=""
plural=""
project="."

while [ $# -gt 0 ]; do
  case "$1" in
    --entity) entity="${2:-}"; shift 2 ;;
    --plural) plural="${2:-}"; shift 2 ;;
    --project) project="${2:-}"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    -*) fail "unknown option: $1" ;;
    *)
      [ -z "$name" ] || fail "unexpected argument: $1"
      name="$1"; shift ;;
  esac
done

[ -n "$name" ] || { usage >&2; exit 1; }
[ -n "$entity" ] || fail "--entity is required"
[ -n "$plural" ] || plural="${entity}s"

for value in "$name" "$entity" "$plural"; do
  case "$value" in
    [a-z]*) ;;
    *) fail "'$value' must start with a lowercase letter" ;;
  esac
  case "$value" in
    *[!a-z0-9]*) fail "'$value' may contain only lowercase letters and digits" ;;
  esac
done
[ "$name" != "$entity" ] || fail "module name and entity must differ"

command -v go >/dev/null || fail "go is required"

api_dir="$project/api"
server_go="$api_dir/internal/server/server.go"
config_go="$api_dir/internal/server/config.go"
module_dir="$api_dir/internal/modules/$name"

[ -f "$api_dir/go.mod" ] || fail "no api/go.mod under '$project'; is this a project from create-go-project-layout?"
[ -f "$api_dir/internal/modules/module.go" ] || fail "missing internal/modules/module.go; scaffold the layout first"
[ -f "$api_dir/internal/platform/pgkit/migrations.go" ] || fail "missing pgkit/migrations.go; the layout predates the module skill"
[ -f "$api_dir/internal/testkit/env.go" ] || fail "missing testkit/env.go; the layout predates the module skill"
grep -q 'allModules := \[\]modules.Module{' "$server_go" || fail "cannot find the allModules slice in $server_go"
grep -q '^type Config struct' "$config_go" || fail "cannot find the Config struct in $config_go"
[ ! -e "$module_dir" ] || fail "module directory already exists: $module_dir"

module_path="$(awk '/^module /{print $2; exit}' "$api_dir/go.mod")"
[ -n "$module_path" ] || fail "cannot read the module path from api/go.mod"

module_pascal="$(pascal "$name")"
entity_pascal="$(pascal "$entity")"

echo "==> Copying module template to $module_dir"
mkdir -p "$module_dir"
cp -R "$TEMPLATE_DIR/." "$module_dir/"

echo "==> Naming files after the entity"
while IFS= read -r file; do
  base="$(basename "$file")"
  dir="$(dirname "$file")"
  case "$base" in
    *entities*) mv "$file" "$dir/${base//entities/$plural}" ;;
    *entity*) mv "$file" "$dir/${base//entity/$entity}" ;;
  esac
done < <(find "$module_dir" -type f \( -name '*entities*' -o -name '*entity*' \))

echo "==> Replacing placeholders"
replace_placeholders() {
  local file="$1" content
  content="$(cat "$file")"
  content="${content//\{\{MODULE_NAME\}\}/$name}"
  content="${content//\{\{MODULE_PASCAL\}\}/$module_pascal}"
  content="${content//\{\{MODULE_PATH\}\}/$module_path}"
  content="${content//\{\{ENTITY\}\}/$entity}"
  content="${content//\{\{ENTITY_PASCAL\}\}/$entity_pascal}"
  content="${content//\{\{ENTITIES\}\}/$plural}"
  printf '%s\n' "$content" > "$file"
}
while IFS= read -r file; do
  replace_placeholders "$file"
done < <(grep -rIl -e '{{' "$module_dir")

if leftover="$(grep -rIn -e '{{[A-Z_][A-Z_]*}}' "$module_dir")"; then
  echo "$leftover" >&2
  fail "unreplaced placeholders remain"
fi
gofmt -w "$module_dir"

echo "==> Registering the module in server.go and config.go"
import_line="\"$module_path/internal/modules/$name\""
tmp="$(mktemp)"

awk -v entry="		$name.NewModule(c.${module_pascal}Config, pool)," -v imp="	$import_line" -v anchor="	\"$module_path/internal/modules\"" '
  $0 == "	allModules := []modules.Module{}" { print "	allModules := []modules.Module{"; print entry; print "	}"; next }
  $0 == "	allModules := []modules.Module{" { print; print entry; next }
  $0 == anchor { print; print imp; next }
  { print }
' "$server_go" > "$tmp" && mv "$tmp" "$server_go"

awk -v field="	${module_pascal}Config $name.Config" -v imp="	$import_line" '
  $0 == "type Config struct{}" { print "type Config struct {"; print field; print "}"; next }
  $0 == "type Config struct {" { print; print field; next }
  $0 == "import (" && !done { print; print imp; done = 1; next }
  $0 == "package server" && !has_import { print; print ""; print "import ("; print imp; print ")"; done = 1; next }
  { print }
' has_import="$(grep -c '^import (' "$config_go")" "$config_go" > "$tmp" && mv "$tmp" "$config_go"

gofmt -w "$server_go" "$config_go"

echo "==> Generating handlers, client, and queries"
(
  cd "$api_dir"
  go generate "./internal/modules/$name/..."
  go mod tidy
  echo "==> Building and vetting"
  go build ./...
  go vet ./...
)

file_count="$(find "$module_dir" -type f | wc -l | tr -d ' ')"
cat <<SUMMARY

Module scaffold complete
  module:      $name ($module_dir)
  entity:      $entity, table $name.$plural
  routes:      POST /$name/$plural, GET /$name/$plural/{uuid}
  files:       $file_count
  registered:  $server_go, $config_go
  next:        cd $api_dir && INTEGRATION=true go test ./internal/modules/$name/... -count=1
SUMMARY
