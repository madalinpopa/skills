---
name: create-go-project-layout
description: Scaffolds a new Go backend or web project from a supported layout, or helps choose a layout before scaffolding. Use when bootstrapping a project or preparing its skeleton for module skills. Not for adding components to existing projects.
status: published
tags: [go, echo, layout, scaffold]
---

# Create Go project layout

Create the skeleton with the bundled script. Leave modules, auth, frontend
implementation, and placeholder packages for later skills.

## Choose a layout

| Layout | Fit | Exclusions |
| --- | --- | --- |
| `modular-monolith-api` | One Go API binary with isolated modules, per-module Postgres schemas, and a separate frontend using JSON over HTTP. | Server-rendered sites, independently deployed services, or a tiny single-resource API. |

For layout advice, explain the fit without collecting scaffold inputs or
creating files. If the requirements do not fit a supported layout, explain
the mismatch; ask only when missing requirements affect the choice.

For architecture explanations, customization, or troubleshooting, read only
the relevant sections of the [layout reference](references/layouts/modular-monolith-api.md).
Routine scaffolding needs no reference or asset-content reads.

## Scaffold

1. Resolve the inputs from the request and available project context. Ask
   together for required values that remain unknown; use defaults for omitted
   optional values.

   | Input | Argument | Default |
   | --- | --- | --- |
   | Supported layout | first argument | chosen above |
   | Lowercase project slug | `--name` | required |
   | Go module path | `--module` | required |
   | Target directory | `--dir` | `./<name>` in the user's working directory |
   | Go major.minor version | `--go` | installed `go version` |

2. Require an absent or empty target directory. From the user's working
   directory, invoke the script by its absolute installed path; `<skill-dir>`
   is the directory containing this `SKILL.md`:

   ```sh
   "<skill-dir>/scripts/scaffold.sh" <layout> --name <slug> --module <path> --dir <target>
   ```

   The script copies assets, replaces and checks tokens, initializes the Go
   module, downloads manifest libraries and generator tools, then runs
   `go mod tidy`, `go build ./...`, and `go vet ./...`. Downloads require Go
   and network access. Tidy may remove unused libraries; generator tools stay.
   For a requested copy-only or offline scaffold, use `--skip-deps`; module
   initialization, dependency installation, build, and vet are skipped. Supply
   `--go <major.minor>` when Go is unavailable. `--list` lists layouts.
3. Report a failed step with its relevant output; do not patch the generated
   project by hand. Do not initialize version control, commit, or start
   containers unless requested.

Keep placeholder `doc.go` files and the TODO spec. Module migrations and the
integration test environment are real infrastructure, not placeholders.

## Report

Report the destination, layout, verification outcome (including skipped
checks), and next step, usually the first module. Include a full tree or
command output only when requested or needed to explain a failure.

## Maintain layouts

Read [adding a layout](references/adding-a-layout.md) only when extending this
skill's supported layouts.
