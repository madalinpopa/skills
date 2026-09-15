#!/usr/bin/env bash
set -euo pipefail

SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TEMPLATE_DIR="$SKILL_DIR/assets/module"

usage() {
  cat <<'USAGE'
Usage: scaffold-auth.sh [--project <dir>]

Adds the auth module to a project created by create-go-project-layout.

Options:
  --project <dir>     project root, the directory holding api/; default .
  -h, --help          show this help
USAGE
}

fail() {
  echo "scaffold-auth: $*" >&2
  exit 1
}

project="."

while [ $# -gt 0 ]; do
  case "$1" in
    --project) project="${2:-}"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    *) fail "unknown argument: $1" ;;
  esac
done

command -v go >/dev/null || fail "go is required"

api_dir="$project/api"
server_go="$api_dir/internal/server/server.go"
config_go="$api_dir/internal/server/config.go"
contracts_go="$api_dir/internal/modules/contracts.go"
envrc="$project/envrc.template"
compose="$project/compose.yaml"
module_dir="$api_dir/internal/modules/auth"

[ -f "$api_dir/go.mod" ] || fail "no api/go.mod under '$project'; is this a project from create-go-project-layout?"
[ -f "$api_dir/internal/modules/module.go" ] || fail "missing internal/modules/module.go; scaffold the layout first"
[ -f "$api_dir/internal/platform/pgkit/migrations.go" ] || fail "missing pgkit/migrations.go; the layout predates the module skills"
[ -f "$api_dir/internal/testkit/env.go" ] || fail "missing testkit/env.go; the layout predates the module skills"
[ -f "$contracts_go" ] || fail "missing internal/modules/contracts.go"
grep -q 'allModules := \[\]modules.Module{' "$server_go" || fail "cannot find the allModules slice in $server_go"
grep -q '^type Config struct' "$config_go" || fail "cannot find the Config struct in $config_go"
grep -q '^	return Config{}, nil$' "$config_go" || fail "LoadConfigFromEnv in $config_go is no longer the stock one; add the auth loader by hand (see references/auth.md)"
grep -q '^type Contracts struct{}$' "$contracts_go" || fail "Contracts in $contracts_go already has fields; add the Auth contract by hand (see references/auth.md)"
grep -q '^import "context"$' "$contracts_go" || fail "$contracts_go is no longer the stock one; add the Auth contract by hand (see references/auth.md)"
[ ! -e "$module_dir" ] || fail "module directory already exists: $module_dir"

module_path="$(awk '/^module /{print $2; exit}' "$api_dir/go.mod")"
[ -n "$module_path" ] || fail "cannot read the module path from api/go.mod"

echo "==> Copying module template to $module_dir"
mkdir -p "$module_dir"
cp -R "$TEMPLATE_DIR/." "$module_dir/"

echo "==> Replacing placeholders"
while IFS= read -r file; do
  content="$(cat "$file")"
  printf '%s\n' "${content//\{\{MODULE_PATH\}\}/$module_path}" > "$file"
done < <(grep -rIl -e '{{' "$module_dir")

if leftover="$(grep -rIn -e '{{[A-Z_][A-Z_]*}}' "$module_dir")"; then
  echo "$leftover" >&2
  fail "unreplaced placeholders remain"
fi
gofmt -w "$module_dir"

echo "==> Registering the module in server.go and config.go"
import_line="\"$module_path/internal/modules/auth\""
tmp="$(mktemp)"

awk -v entry="		auth.NewModule(c.AuthConfig, pool)," -v imp="	$import_line" -v anchor="	\"$module_path/internal/modules\"" '
  $0 == "	allModules := []modules.Module{}" { print "	allModules := []modules.Module{"; print entry; print "	}"; next }
  $0 == "	allModules := []modules.Module{" { print; print entry; next }
  $0 == anchor { print; print imp; next }
  { print }
' "$server_go" > "$tmp" && mv "$tmp" "$server_go"

awk -v field="	AuthConfig auth.Config" -v imp="	$import_line" -v has_import="$(grep -c '^import (' "$config_go")" '
  $0 == "type Config struct{}" { print "type Config struct {"; print field; print "}"; next }
  $0 == "type Config struct {" { print; print field; next }
  $0 == "import (" && !done { print; print "	\"fmt\""; print ""; print imp; done = 1; next }
  $0 == "package server" && has_import == 0 { print; print ""; print "import ("; print "	\"fmt\""; print ""; print imp; print ")"; done = 1; next }
  $0 == "	return Config{}, nil" {
    print "	authConfig, err := auth.ConfigFromEnv()"
    print "	if err != nil {"
    print "		return Config{}, fmt.Errorf(\"loading auth config: %w\", err)"
    print "	}"
    print ""
    print "	return Config{AuthConfig: authConfig}, nil"
    next
  }
  { print }
' "$config_go" > "$tmp" && mv "$tmp" "$config_go"

echo "==> Publishing the Auth contract in contracts.go"
awk '
  $0 == "import \"context\"" { print "import ("; print "	\"context\""; print "	\"uuid\""; print ")"; next }
  $0 == "type Contracts struct{}" {
    print "// Auth verifies bearer access tokens for every other module. The caller"
    print "// crosses the module boundary as a plain UUID."
    print "type Auth interface {"
    print "	VerifyAccessToken(ctx context.Context, token string) (uuid.UUID, error)"
    print "}"
    print ""
    print "type Contracts struct {"
    print "	Auth Auth"
    print "}"
    next
  }
  { print }
' "$contracts_go" > "$tmp" && mv "$tmp" "$contracts_go"

gofmt -w "$server_go" "$config_go" "$contracts_go"

if [ -f "$envrc" ] && ! grep -q 'AUTH_SECRET_KEY' "$envrc"; then
  echo "==> Adding AUTH variables to envrc.template"
  awk '
    { print }
    $0 == "# startup, so a bad value fails boot instead of a request." {
      print ""
      print "# Auth module. The secret key signs access tokens and must have at least 32"
      print "# bytes; generate one with: openssl rand -base64 32"
      print "export AUTH_SECRET_KEY=\"change-me-to-a-random-value-of-32-bytes-or-more\""
      print "export AUTH_ISSUER=\"api\""
      print "export AUTH_AUDIENCE=\"api\""
      print "export AUTH_ACCESS_TOKEN_TTL=\"15m\""
      print "export AUTH_REFRESH_TOKEN_TTL=\"12h\""
      print "export AUTH_COOKIE_SECURE=\"true\""
      added = 1
    }
    END { if (!added) print "\n# Auth module\nexport AUTH_SECRET_KEY=\"change-me-to-a-random-value-of-32-bytes-or-more\"" }
  ' "$envrc" > "$tmp" && mv "$tmp" "$envrc"
fi

if [ -f "$compose" ] && ! grep -q 'AUTH_SECRET_KEY' "$compose"; then
  echo "==> Passing AUTH variables to the API container in compose.yaml"
  awk '
    { print }
    $0 == "      POSTGRES_URL: ${POSTGRES_URL}" {
      print "      AUTH_SECRET_KEY: ${AUTH_SECRET_KEY}"
      print "      AUTH_ISSUER: ${AUTH_ISSUER}"
      print "      AUTH_AUDIENCE: ${AUTH_AUDIENCE}"
      print "      AUTH_ACCESS_TOKEN_TTL: ${AUTH_ACCESS_TOKEN_TTL}"
      print "      AUTH_REFRESH_TOKEN_TTL: ${AUTH_REFRESH_TOKEN_TTL}"
      print "      AUTH_COOKIE_SECURE: ${AUTH_COOKIE_SECURE}"
    }
  ' "$compose" > "$tmp" && mv "$tmp" "$compose"
fi

echo "==> Fetching jwt and bcrypt, generating handlers, client, and queries"
(
  cd "$api_dir"
  go get github.com/golang-jwt/jwt/v5 golang.org/x/crypto
  go generate ./internal/modules/auth/...
  go mod tidy
  echo "==> Building and vetting"
  go build ./...
  go vet ./...
)

file_count="$(find "$module_dir" -type f | wc -l | tr -d ' ')"
cat <<SUMMARY

Auth module scaffold complete
  module:      auth ($module_dir)
  tables:      auth.users, auth.refresh_tokens
  routes:      POST /auth/register, POST /auth/login, POST /auth/refresh, POST /auth/logout, GET /auth/me
  contract:    modules.Contracts.Auth (VerifyAccessToken)
  files:       $file_count
  registered:  $server_go, $config_go, $contracts_go
  env:         AUTH_SECRET_KEY (required), AUTH_ISSUER, AUTH_AUDIENCE, AUTH_ACCESS_TOKEN_TTL, AUTH_REFRESH_TOKEN_TTL, AUTH_COOKIE_SECURE
  next:        cd $api_dir && go test ./internal/modules/auth/... && INTEGRATION=true go test ./internal/modules/auth/... -count=1
SUMMARY
