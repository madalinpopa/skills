# Idiomatic Go

Use only for Go work. Read the applicable `go.mod` and toolchain configuration;
do not assume the agent's preferred Go version matches the module.

Use `modern-go-guidelines:use-modern-go` as selected by AGENTS.md. Follow its
installed instructions for version-specific guidance, using the existing file
or the target module version for a new file. Read the full applicable guidance;
use detailed explanations only where needed. Keep version-specific rules in
that skill instead of copying them here.

- Keep packages cohesive and APIs small. Prefer explicit error returns and
  follow the project's wrapping and error-handling conventions.
- Define interfaces where consumers need them; do not wrap concrete types
  without a real substitution or boundary requirement.
- Make resource ownership and cleanup explicit. Add concurrency only for a
  concrete need, with clear cancellation and goroutine lifetimes.
- Follow local context and test conventions. Prefer table tests when cases
  share setup and assertions; do not force unrelated cases into a table.
- Use the project's format, test, and vet commands. Do not introduce language
  features or library APIs beyond the supported toolchain.
