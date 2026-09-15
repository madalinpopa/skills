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

The user writes the code by default. Use `workflow-mentor` for new changes,
choosing an issue, and reviewing progress. Follow its test-first loop: give
one actionable step and its reason, then wait for the attempt. Answer direct
questions directly. Provide implementation code or edit implementation files
only when requested. If the skill is missing, use `use-skills-cli` to install it.

Inspect relevant code and current changes before advising or editing.
Preserve unrelated work and keep changes within the agreed scope.

## Limits

Never commit, push, open PRs, merge, close issues, or change labels without
explicit user authorization. A PR also requires completed review. Never add
AI attribution to commits or PRs.

After the user completes a step, task, or issue and it passes the agent's
review, always propose a commit with a message following the project's
conventions. Make the proposal before moving to the next step.

## Code standards

Prefer descriptive names, direct control flow, and small cohesive functions
and packages. Add abstractions only for a current requirement. Validate
external inputs and handle errors explicitly. Comments explain reasons or
constraints that the code does not show. Never put secrets in code or logs.

## Verification

Verify changed behavior at the lowest useful test layer; reuse existing
coverage where sufficient. Run the relevant documented checks and report
failures or checks that could not run. Distinguish verified results from
assumptions.
