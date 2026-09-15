---
name: create-auth-module
description: Adds an auth module to a Go API project created with the create-go-project-layout skill. It scaffolds users and refresh token tables, register, login, refresh, logout, and profile routes from an OpenAPI spec, JWT access tokens, rotated refresh tokens with replay detection, a bearer middleware, and the Auth contract other modules verify tokens through, then registers the module, its config, and its environment variables. Use whenever the user wants authentication, login, JWT, sessions, user accounts, identity, or protected routes in such a project, even if they only say "add auth" or "users need to log in". Not for a plain feature module (use create-api-module), for roles and permissions, for social login, or for projects with a different layout.
status: published
tags: [go, echo, auth, jwt, module, scaffold]
---

# Create auth module

Add the `auth` module to a project made by `create-go-project-layout`. The
result compiles, migrates its own schema, serves five routes, publishes the
`Auth` contract, and passes unit and integration tests. It owns user
accounts, login sessions, and the caller's profile. Roles and permissions are
not included; add them in the module after.

## Inputs

| Input | Script flag | Example |
| --- | --- | --- |
| Project root, the directory with `api/` | `--project`, default `.` | `./fia` |

The module name, schema, and route prefix are fixed to `auth`. Everything
else is configuration read at startup; see the reference.

## Scaffold

1. Read [references/auth.md](references/auth.md). It explains every file the
   script writes, the token flow, the configuration, and how another module
   consumes the contract.
2. Run the script from this skill's directory:

   ```sh
   scripts/scaffold-auth.sh --project <dir>
   ```

   It checks that the project came from the layout skill, copies
   `assets/module/` into `api/internal/modules/auth/`, replaces the module
   path, registers the module in `server.go` and `config.go`, publishes the
   contract in `modules/contracts.go`, adds the `AUTH_*` variables to
   `envrc.template` and `compose.yaml`, fetches `jwt` and `x/crypto`, runs
   `go generate` for the module, then `go mod tidy`, `go build ./...`, and
   `go vet ./...`.
3. Read the script output. If a step fails, report it with its output instead
   of patching the generated files by hand. When the script refuses because
   `config.go` or `contracts.go` were already edited, make the three edits
   listed in the reference and rerun the Go steps.
4. Run the unit tests, and the integration tests when Docker is available:

   ```sh
   cd <project>/api && go test ./internal/modules/auth/...
   cd <project>/api && INTEGRATION=true go test ./internal/modules/auth/... -count=1
   ```

   Report the result either way. Skipping the integration run without Docker
   is fine; say so.
5. Do not commit unless the request asks for it.

## Report

List the files created, the registration edits, the environment variables the
API now needs, the commands run with their results, and the next steps: set a
real `AUTH_SECRET_KEY`, and protect other modules' routes through the `Auth`
contract as shown in the reference.
