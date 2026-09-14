---
name: create-api-module
description: Adds a new feature module to a Go API project created with the create-go-project-layout skill. It scaffolds the module directory with domain, app, REST and Postgres adapters, an OpenAPI spec, a first migration and query, a generated client, and an integration test, then registers the module in the server. Use whenever the user wants to add, create, or scaffold a module, a bounded context, a feature slice, or a new resource with its own table and endpoints in such a project, even if they only say "add users" or "I need a products API". Not for adding a route to an existing module, for auth or other platform components, or for projects with a different layout.
status: published
tags: [go, echo, module, scaffold]
---

# Create API module

Scaffold one module inside a project made by `create-go-project-layout`.
The result compiles, migrates its own schema, serves two routes for a first
entity, and passes an integration test. Rename or extend the entity after.

## Inputs

| Input | Script flag | Example |
| --- | --- | --- |
| Module name, lowercase letters and digits | first argument | `brand` |
| Project root, the directory with `api/` | `--project`, default `.` | `./fia` |
| First entity, singular, lowercase | `--entity` | `profile` |
| Plural of the entity | `--plural`, default entity plus `s` | `profiles` |

Ask for the entity when the request does not name one. It becomes the first
table, the typed identifier, and the two routes, so a vague placeholder such
as `item` costs a rename later.

## Scaffold

1. Read [references/module.md](references/module.md). It explains every file
   the script writes and the rules a module follows.
2. Run the script from this skill's directory:

   ```sh
   scripts/scaffold-module.sh <name> --project <dir> --entity <singular>
   ```

   It checks that the project came from the layout skill, copies
   `assets/module/` into `api/internal/modules/<name>/`, replaces the tokens,
   registers the module in `server.go` and `config.go`, runs `go generate`
   for the module, then `go mod tidy`, `go build ./...`, and `go vet ./...`.
3. Read the script output. If a step fails, report it with its output instead
   of patching the generated files by hand.
4. Run the module's integration test when Docker is available:

   ```sh
   cd <project>/api && INTEGRATION=true go test ./internal/modules/<name>/... -count=1
   ```

   Report the result either way. Skipping it without Docker is fine; say so.
5. Do not commit unless the request asks for it.

## Report

List the files created, the two registration edits, the commands run with
their results, and the next step: replace the starter entity's fields, add
queries and routes, and regenerate with `task gen`.
