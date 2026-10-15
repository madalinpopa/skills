# Agent instructions

## Project

## Working rules

- Work in small slices, including docs, configuration, and refactoring.
- Inspect relevant code and diffs; preserve unrelated work and stay in scope.
- Raise conflicts between instructions, requirements, and code before changing
  a contract.
- Verify each slice and present it for user review before the next.

## Version control

- When `.jj` exists, use jj-vcs and invoke `use-jj`, even alongside `.git`.
- If `use-jj` is unavailable, report it before mutations; do not substitute Git.
- Never commit unless explicitly asked; plan or implementation approval does
  not authorize a commit. Publish only on explicit request.
- Never add AI attribution.

## Specs

- Preserve agreed requirements; add only established content. Omit placeholders
  and empty sections.
- For requested spec work, always use the available question tool with concrete
  proposed options based on repository evidence.
- Ask focused rounds, reuse prior answers, and wait before writing dependent
  content. If no tool is available, explain and ask option-based questions in chat.
- Mark completed tasks with evidence; completion does not imply user acceptance.

## Implementation

- Check current official docs for external APIs being changed and latest
  releases for dependencies being proposed.
- Dependency additions, installations, replacements, and upgrades need approval.

## Read for the current task

- Before scope decisions, read [MISSION.md](docs/specs/MISSION.md).
- Before selecting work or updating progress, read [ROADMAP.md](docs/specs/ROADMAP.md).
- Before running project checks, read [STACK.md](docs/specs/STACK.md) for commands
  and working directories; report missing ones.
- Before code or dependency changes, read [CODE.md](docs/conventions/CODE.md).
- Before editing Go files, read [GOLANG.md](docs/conventions/GOLANG.md).
- Before behavior changes, test work, or verification, read
  [TDD.md](docs/conventions/TDD.md) and follow its applicable review gates.
- Before commit proposals or delivery, read [COMMITS.md](docs/conventions/COMMITS.md).
- Before spec edits, read the [spec writing guide](docs/specs/README.md).
- Before recording architecture decisions, read the [decision guide](docs/decisions/README.md).
- Before planning or resuming work, read the [feature guide](docs/features/README.md).
- Load historical records only for relevant dependencies.

## Replies

- Use short bullets for changes and verification; avoid repetition.
- End with **User must be aware**: material risks, limits, unresolved decisions,
  or "None".
- Record actual open questions where relevant; omit unknown topics.
