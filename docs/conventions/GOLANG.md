# Go conventions

## Modern Go

- Before creating or editing any Go file, invoke
  `$modern-go-guidelines:use-modern-go`. Run `list` for that exact file, or Go
  1.27 for a new file, and read the full, unfiltered output.
- Use `explain` only for IDs needing detail. Re-run guidance when moving to a
  different kind of file or package. Follow it unless it would change required
  behavior or break the build; explain exceptions.
- Write explicit, idiomatic, readable code with descriptive names and direct
  control flow. Keep packages small and focused; never create them to fill
  slots. Complex abstractions need a concrete reason.
- Trust internal invariants. Validate external input, config, filesystem
  state, and Git results; avoid speculative defensive layers.
- Let code explain itself. Add short doc comments only for useful information
  not obvious from code; avoid long comment blocks.

## Package boundaries

- `main.go`: version injection, process context, command execution only.
- `cmd`: [Cobra](https://cobra.dev/) commands, argument checks, flags, help, output
  rendering, writers, and exit codes.
- `internal/config`: config types, defaults, [Viper](https://github.com/spf13/viper)
  TOML loading, and first-run save.
- `internal/project`: repository-root discovery with a fallback outside repos.
- `internal/store`: Git clone, fast-forward sync, commit lookup.
- `internal/skill`: catalog discovery, YAML metadata checks, per-agent transforms.
- `internal/install`: target trees, three-way planning, locks, filesystem
  changes, backups, update, removal.
- `internal/gittest`: helpers creating temporary Git repositories for tests.

Keep Cobra in `cmd` and Viper in `internal/config`. Domain packages must neither
depend on them nor print. Pass Git arguments directly to the process runner;
never construct shell command strings.
