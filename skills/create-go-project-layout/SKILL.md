---
name: create-go-project-layout
description: Scaffolds a new Go web project from a named layout. It creates the directory tree, placeholder Go packages, tooling (Task, Docker Compose, Dockerfile, golangci-lint, GitHub workflow), a docs folder with the spec and feature templates, and the agent instruction files. Use whenever the user wants to start, bootstrap, scaffold, or set up a new Go backend or web project, asks which layout fits an application (separate API and frontend, modular monolith, modules behind one frontend), or a workflow skill needs a project skeleton before modules are added. Not for adding a module, an auth component, or a frontend to an existing project, and not for projects in other languages.
status: published
tags: [go, echo, layout, scaffold]
---

# Create Go project layout

Scaffold a new project from one of the layouts in `references/layouts/`.
The result is a skeleton: the tree, the tooling, the docs folder, and
placeholder packages. Modules, auth, and other components are added later by
their own skills. Do not write them here.

## Inputs

Collect these before running anything. Ask for the ones the request does not
give.

| Input | Script flag | Example |
| --- | --- | --- |
| Layout name | first argument | `modular-monolith-api` |
| Project name, a short slug | `--name` | `fia` |
| Go module path | `--module` | `github.com/acme/fia` |
| Target directory | `--dir`, default `./<name>` | `./fia` |
| Go version, major and minor | `--go`, default from `go version` | `1.27` |

The target directory must not exist or must be empty. Never scaffold into an
existing project.

## Choose a layout

Every file in `references/layouts/` starts with a short fit summary. Read the
summaries, pick the layout that matches the application, and say why. When two
layouts could fit, ask.

| Layout | Good fit |
| --- | --- |
| [modular-monolith-api](references/layouts/modular-monolith-api.md) | One Go API built from isolated modules, with a separate frontend such as a SPA or a mobile app that talks JSON over HTTP. |

## Scaffold

1. Read the chosen layout reference in full. It explains what every directory
   is for and which files are placeholders.
2. Run the script from this skill's directory:

   ```sh
   scripts/scaffold.sh <layout> --name <slug> --module <path> --dir <target>
   ```

   It copies `assets/<layout>/`, replaces the placeholder tokens, checks that
   none is left, then runs `go mod init`, `go get` for the libraries and
   `go get -tool` for the generators listed in the layout manifest, `go mod
   tidy`, `go build ./...`, and `go vet ./...`. `--list` prints the layouts.
   `--skip-deps` copies the files only, for a machine without Go.
3. Read the script output. If a step fails, report it with its output instead
   of patching the generated project by hand.
4. Do not initialise version control, commit, or start containers unless the
   request asks for it.

## Placeholders

A placeholder package holds one `doc.go` that says what belongs there and which
kind of skill fills it. Leave these files in place; a component skill replaces
them. `docs/SPEC.md` starts with a TODO for the same reason. Migrations and
the integration test environment are real code, because the module skill
depends on them.

## Report

List the tree that was created, the commands run with their results, the
placeholders left for later skills, and the next step, which is usually to
create the first module.

## Adding a layout

Add three things and the script picks the layout up by name:

- `assets/<layout>/` with the project tree and a `.scaffold` manifest. The
  manifest holds `go_dir`, one `require` line per library, and one `tool` line
  per generator. The script deletes it from the target after copying.
- `references/layouts/<layout>.md` that starts with the fit summary and keeps
  the section order of the existing reference.
- One row in the table above.

Use only the three placeholder tokens the script knows. A layout that needs
another token also needs a script change.
