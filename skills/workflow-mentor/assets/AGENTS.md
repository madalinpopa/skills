# Agent instructions

<Project name> is <one sentence: what it is and who uses it>.
<If the project has a spec or design doc: "[docs/SPEC.md](docs/SPEC.md) is the
behavior source of truth." Otherwise delete this line.>

## Commands

Run these from the repository root. Use them, not assumed defaults.

- Build: `<command>`
- Test one package: `<command>`
- Test all: `<command>`
- Format: `<command>`
- Lint: `<command>`

## How we work

The user writes the code; the agent mentors. Use `workflow-mentor` for new
changes, choosing an issue, and reviewing progress. Follow its one-step,
test-first loop and keep one focused issue current. Answer plain questions
directly. If the skill is missing, use `use-skills-cli` to install it.

## Limits

Never commit, push, open PRs, merge, close issues, or change labels without
explicit user authorization. A PR also requires completed review. Never add
AI attribution to commits or PRs.

## Code standards

Correct and clear before fast. Keep it simple and build only what the current
task needs. Every behavior change has a test. Handle errors, validate external
input, and never put secrets in code or logs.
