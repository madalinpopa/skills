# Layout: modular monolith API

## Fit

Good fit: one Go HTTP API made of isolated feature modules, deployed as one
binary, with a frontend that lives in its own directory and talks to the API
over JSON. The frontend is a SPA or a mobile app served separately, so the API
needs CORS and a bearer token flow rather than sessions and templates. Each
module owns its OpenAPI spec, its Postgres schema, its migrations, and its
tests, and modules talk to each other only through small contract interfaces.

Not a fit: a server-rendered site with HTML templates, a project with several
deployable services, or a tiny API with a single resource where module
boundaries add more files than they save.

## Contents

- [Fit](#fit)
- [Tree](#tree)
- [How the server boots](#how-the-server-boots)
- [Module shape](#module-shape)
- [Generated code](#generated-code)
- [Platform packages](#platform-packages)
- [Dependencies](#dependencies)
- [Tooling and commands](#tooling-and-commands)
- [Environment](#environment)

## Tree

```
{{PROJECT_NAME}}/
├── AGENTS.md                 engineering rules and verification commands
├── CLAUDE.md                 points Claude Code at AGENTS.md
├── README.md                 layout and development commands
├── Taskfile.yml              every repository command, run from the root
├── compose.yaml              local API container with live reload plus Postgres
├── envrc.template            copy to .envrc; direnv loads it
├── .gitignore
├── .gitattributes            marks generated Go files for GitHub diffs
├── .github/workflows/api.yml test on feature branches, build image on main, release on tags
├── docs/
│   ├── SPEC.md               behavior source of truth, starts as a TODO
│   ├── features/             one folder per feature, from the templates
│   └── templates/            FEATURE.md and PHASE.md planning templates
├── api/                      the Go module
│   ├── Dockerfile            development, build, and production stages
│   ├── .dockerignore
│   ├── .golangci.yml         golangci-lint v2 with gofumpt
│   ├── reflex.conf           live reload command used by the development stage
│   ├── cmd/server/main.go    flags, logging, signals, pool, Echo, server.New, Run
│   └── internal/
│       ├── shared/            placeholder: types every module shares (UUID)
│       ├── server/
│       │   ├── config.go      Config and LoadConfigFromEnv; one field per module
│       │   └── server.go      module lifecycle and Run
│       ├── modules/
│       │   ├── module.go      the Module interface every module implements
│       │   ├── router.go      EchoRouter, the router surface modules register on
│       │   ├── contracts.go   cross-module interfaces and the Contracts registry
│       │   └── <name>/        one directory per module, created by the module skill
│       ├── platform/          shared code that knows no module
│       │   ├── errs/          placeholder: error type and Echo error handler
│       │   ├── httpserver/    placeholder: Echo constructor and middleware stack
│       │   ├── logging/       placeholder: slog setup and context helpers
│       │   └── pgkit/         MigrateDatabaseUp for one module schema
│       └── testkit/           Postgres container, Env, Wire, Serve for module tests
└── web/README.md             placeholder for the frontend
```

## How the server boots

`cmd/server/main.go` parses flags, sets the default `slog` logger, creates a
context that ends on SIGINT or SIGTERM, loads the config, opens one
`pgxpool.Pool` from `POSTGRES_URL`, creates the Echo instance with a
`GET /health` route, and hands everything to `server.New`.

`server.New` owns the module lifecycle. It builds the list of modules, then
runs five phases over all of them, in order:

1. `Init`: validate the module config, build repositories, use cases, and the
   HTTP handler, run the module's migrations.
2. `RegisterContracts`: publish the interfaces other modules may call.
3. `Contracts.Verify`: fail startup if any expected contract is missing.
4. `Connect`: take the contracts a module depends on, for example an identity
   verifier used to build an auth middleware.
5. `RegisterHTTP`: register routes on the Echo router.

The order matters. A module can only build things that depend on another
module in `Connect`, because the registry is complete only after every module
has published in phase 2. `Run` starts Echo with a graceful shutdown timeout
tied to the process context.

A module is registered by adding its constructor to the `allModules` slice in
`server.go` and its config as a field in `server.Config`. If it publishes a
contract, it also gets an interface and a field in `modules/contracts.go` and
a check in `Verify`. Nothing else outside the module changes.

## Module shape

The module skill creates these. The layout only reserves the place.

```
internal/modules/<name>/
├── module.go        Module implementation, embeds the migrations
├── config.go        module Config with Validate
├── contract.go      what this module publishes to others
├── domain/          entities, value objects, repository interfaces, typed IDs
├── app/             commands, queries, ports the adapters implement
├── adapters/
│   ├── rest/        openapi.yaml, generated strict server, handlers, middleware
│   └── pgstore/     sqlc.yaml, migrations/, queries/, generated models/, repos
└── tests/           module-level integration tests through the real HTTP API
```

Dependency rules, enforced by review:

- `domain` imports only `shared` and `platform/errs`.
- `app` imports `domain` and defines the ports it needs.
- `adapters/*` import `app` and `domain`, never the reverse.
- `module.go` is the only place that wires them together.
- Modules never import each other. They use `modules.Contracts`.
- `platform/*` never imports a module.
- Each module owns a Postgres schema named after it. No foreign keys cross
  schemas; another module's ID is stored as a plain UUID value.

## Generated code

Handlers are contract-first. Each module keeps an `openapi.yaml` and runs
`oapi-codegen` with `echo5-server`, `strict-server`, and `models` enabled, so
hand-written handlers implement the generated `StrictServerInterface`. A
second config generates a Go client used by the module's integration tests.

Queries are SQL-first. Each module keeps `queries/*.sql` and `migrations/*.up.sql`
and runs `sqlc` with `sql_package: pgx/v5`. Both generators are Go tools in
`go.mod` and run through `//go:generate` directives, so `task gen` regenerates
everything and formats the result. Generated files are never edited by hand.

## Platform packages

Two are real code because every module needs them from its first commit:

- `pgkit.MigrateDatabaseUp` creates the module's schema and applies its
  embedded `*.up.sql` files with golang-migrate, tracking them in a
  `schema_migrations` table inside that schema.
- `testkit` starts a Postgres container with Testcontainers, holds the pool,
  the Echo instance, and the contracts registry in an `Env`, wires one module
  through the same five phases as `server.New`, and serves it over `httptest`.
  `RequireIntegration` skips a test unless `INTEGRATION=true`.

The rest are placeholders. Their intended content:

- `errs`: an `Error` type with an HTTP status, a public slug and message, an
  internal error for logs, and optional details; constructors per status; an
  Echo error handler that renders it as JSON.
- `httpserver`: `NewEcho` with options for CORS origins and body limit, and the
  middleware stack: context timeout, recover, body limit, CORS, correlation
  ID, request log.
- `logging`: `slog` initialisation and context helpers for a request logger and
  correlation ID.
- `shared`: types every module needs that belong to no module. Identifiers
  use the standard library `uuid` package directly, so this may stay empty.

## Dependencies

The `.scaffold` manifest next to the assets lists them: Echo v5, pgx v5, the
oapi-codegen runtime, golang-migrate, Testify, and Testcontainers with its
Postgres module as libraries, plus oapi-codegen and sqlc as Go tools. The
script fetches the latest versions, so nothing is pinned here.

`go mod tidy` drops the libraries nothing imports yet; that is expected. The
tools stay in the `tool` block. A library comes back when a component or
module imports it.

## Tooling and commands

All commands run from the repository root through `Taskfile.yml`.

| Task | What it does |
| --- | --- |
| `task init` | download Go modules and install web dependencies |
| `task up` / `task down` | start or stop the API and Postgres with Compose |
| `task gen` | `go generate ./...` then format |
| `task test:unit` | `go test ./...`, Docker-free |
| `task test:integration` | `INTEGRATION=true go test ./... -count=1` |
| `task lint` | golangci-lint and a `go mod tidy -diff` check |
| `task format` | `go fmt` and `golangci-lint fmt` |
| `task security:vuln` | govulncheck |
| `task patch` / `minor` / `major` | tag and push the next release |

The development Compose stage mounts `api/` and runs `reflex`, so the server
restarts on every Go or SQL change. The production image is a static binary
on Alpine, running as a non-root user with a health check on `/health`.

## Environment

`envrc.template` lists every variable. `POSTGRES_URL` is the only one the
skeleton reads. Modules and components add their own, and the API validates
them at startup so a bad value fails boot instead of a request.
