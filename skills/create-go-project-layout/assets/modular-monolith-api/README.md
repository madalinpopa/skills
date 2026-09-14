# {{PROJECT_NAME}}

TODO: one paragraph on what the application does.

## Project layout

- `api/` contains the Go API, its Dockerfile, and backend tooling.
- `web/` contains the frontend source files.
- `docs/` contains the specification and feature planning documents.
- `compose.yaml` runs the local development services.
- `Taskfile.yml` provides repository-level development commands.

## Development

Create `.envrc` from `envrc.template`, then run:

```sh
task init
task up
```

Run the API checks from the repository root:

```sh
task test
task lint
```

Regenerate handlers and queries after editing an `openapi.yaml` or a
`queries/*.sql` file:

```sh
task gen
```
