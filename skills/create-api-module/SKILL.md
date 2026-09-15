---
name: create-api-module
description: Scaffolds and registers a new feature module or resource with its own table and endpoints in a Go API built with create-go-project-layout. Not for extending existing modules, authentication, platform components, or projects with a different layout.
status: published
tags: [go, echo, module, scaffold]
---

# Create API module

Use the bundled script to create one module with a starter entity, its
Postgres schema, create/get routes, generated adapters, and integration tests.

For architecture explanations, customization, or troubleshooting, read only
relevant sections of [module anatomy](references/module.md). Routine
scaffolding needs no reference or asset-content reads.

## Inputs and prerequisites

Resolve inputs from the request and project context. Ask together for missing
required values; do not invent a placeholder entity. Use optional defaults.

| Input | Argument | Default or example |
| --- | --- | --- |
| Module name | first argument | required; `brand` |
| First entity, singular | `--entity` | required; `profile` |
| Entity plural | `--plural` | entity plus `s`; override irregular plurals |
| Project root containing `api/` | `--project` | current directory |

Names must start with a lowercase letter and contain only lowercase letters
and digits. Module and entity names must differ. The module directory must
not already exist.

Use a project compatible with `create-go-project-layout`, with Go and its
module generator tools available. Uncached dependencies may require network
access. Inspect the target project's instructions and the registration files
below before running the script; explain a layout mismatch instead of forcing
registration into another architecture.

## Scaffold and verify

1. From the target project root, invoke the script by its absolute installed
   path; `<skill-dir>` contains this `SKILL.md`:

   ```sh
   "<skill-dir>/scripts/scaffold-module.sh" <name> --project . --entity <singular>
   ```

   It checks layout prerequisites, copies and names module files, replaces
   tokens, and registers the module in `api/internal/server/server.go` and
   `config.go`. It runs module code generation, `go mod tidy`, `go build ./...`,
   and `go vet ./...`.
2. On failure, report the failed step and relevant output, including files
   already changed. Do not patch generated files by hand or continue to the
   integration test. For customization, edit OpenAPI or SQL sources and
   regenerate rather than editing generated adapters.
3. After success, run the module tests when Docker is available, from `api/`:

   ```sh
   INTEGRATION=true go test ./internal/modules/<name>/... -count=1
   ```

   Report passed, failed, or skipped; state why if Docker is unavailable.
   Do not commit unless requested.

## Report

Report the module path, entity and routes, registration edits, verification
outcomes, and next step: adapt starter fields, queries, and routes, then run
`task gen`. Include a full file list or command output only when requested or
needed to explain failure.
