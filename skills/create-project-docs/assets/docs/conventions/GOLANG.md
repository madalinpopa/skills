# Go conventions

Apply only to Go work, alongside [CODE.md](CODE.md) and [TDD.md](TDD.md).
Read the relevant `go.mod` and toolchain configuration; use supported versions,
not the machine's newest Go. Record versions and module roots in
[STACK.md](../specs/STACK.md).

- Before creating, editing, or reviewing Go files, invoke `use-go-conventions`.
  If it is unavailable, report it and follow official Go documentation.
- Coverage follows [TDD.md](TDD.md); review follows [REVIEW.md](REVIEW.md);
  dependency changes follow [CODE.md](CODE.md#external-dependencies).
- Run Go checks with the commands and module directories in
  [STACK.md](../specs/STACK.md).

## Project-specific Go conventions
