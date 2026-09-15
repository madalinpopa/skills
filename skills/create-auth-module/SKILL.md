---
name: create-auth-module
description: Scaffolds authentication, login, user accounts, and JWT sessions in a Go API built with create-go-project-layout. Use when adding an auth module. Not for extending existing auth, plain feature modules, roles and permissions, social login, or other project layouts.
status: published
tags: [go, echo, auth, jwt, module, scaffold]
---

# Create auth module

Add the fixed `auth` module, schema, and route prefix with register, login,
refresh, logout, and profile endpoints plus the `Auth` contract.

## Prerequisites

Resolve the project root containing `api/` from the request or context;
`--project` defaults to the current directory. Ask only if the target is
unclear. Require a compatible `create-go-project-layout` project, an absent
`api/internal/modules/auth/`, Go, and the module generator tools. The script
fetches dependencies and may need network access.

Inspect the project's instructions and existing registration/configuration
files before changing them. Routine scaffolding needs no full reference or
asset-content reads.

## Scaffold and verify

1. From the target project root, invoke the script by its absolute installed
   path; `<skill-dir>` contains this `SKILL.md`:

   ```sh
   "<skill-dir>/scripts/scaffold-auth.sh" --project .
   ```

   It copies templates, replaces the module path, registers auth in
   `api/internal/server/{server,config}.go` and `api/internal/modules/contracts.go`,
   and adds `AUTH_*` settings to `envrc.template` and `compose.yaml`. It fetches
   JWT/bcrypt dependencies, generates handlers/client/SQL models, then runs
   `go mod tidy`, `go build ./...`, and `go vet ./...`.
2. On failure, report the step, relevant output, and files already changed.
   For refusals of customized config/contracts, follow
   [manual registration](references/auth.md#registration); these checks happen
   before module files are copied. For other failures, stop rather than
   patching generated files or continuing to tests.
3. After successful scaffolding or manual integration, run unit tests from
   `api/`, then integration tests when Docker is available:

   ```sh
   go test ./internal/modules/auth/...
   INTEGRATION=true go test ./internal/modules/auth/... -count=1
   ```

   Report each result separately; state why an integration run was skipped.
   Do not commit unless requested.

## Configuration and follow-up

Before running the API, replace the template secret with a random signing key
of at least 32 bytes using `AUTH_SECRET_KEY` or `AUTH_SECRET_KEY_FILE`.
Unset `AUTH_SECRET_KEY` when using the file; a nonempty value takes precedence.
Refresh cookies default to Secure and SameSite=Lax; disable Secure only for
plain HTTP development. Other modules' routes require explicit integration
through the `Auth` contract; adding this module alone does not protect them.

Read only sections relevant to the task:

- [Configuration](references/auth.md#configuration): environment defaults,
  token lifetimes, and cross-site cookie configuration.
- [Consuming the contract](references/auth.md#consuming-the-contract): protecting
  another module's routes and wiring its tests.
- [Auth anatomy](references/auth.md): use its contents for route/token
  explanations, implementation details, or customization.

## Report

Summarize the module path, registration/configuration edits, check outcomes,
required secret setup, and next step. Include full file lists or command
output only when requested or needed to explain failure.
