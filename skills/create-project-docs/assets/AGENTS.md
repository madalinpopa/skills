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
- Commit only after user approval; publish only on explicit request.
- Never add AI attribution.

## Specs

- Never invent requirements; fill specs progressively.
- Start specs empty; add sections only with known content. Omit placeholders.
- Follow the [spec writing guide](docs/specs/README.md): group related facts
  under headings and write short lists; avoid long paragraphs and prose-heavy tables.
- For requested spec work, always use the available question tool with concrete
  proposed options based on repository evidence.
- Ask focused rounds, reuse prior answers, and wait before writing dependent
  content. If no tool is available, explain and ask option-based questions in chat.
- ROADMAP.md contains only feature/sub-feature groups and task checklists;
  leave it empty when no work is identified.
- Requirements own behavior; plans own execution and checkpoints; validation
  owns evidence and user acceptance. Link shared facts.
- Mark completed tasks with evidence; completion does not imply user acceptance.

## Implementation

- Use STACK.md commands and working directories; report missing ones.
- Check current official docs and latest released library versions.
- Dependency additions, installations, replacements, and upgrades need approval.
- Follow [TDD](docs/conventions/TDD.md#workflow): failing test, implementation,
  independent adversarial review until clean, then user review.

## Read for the current task

- [MISSION.md](docs/specs/MISSION.md): goals, users, scope.
- [ROADMAP.md](docs/specs/ROADMAP.md): features, sub-features, tasks.
- [STACK.md](docs/specs/STACK.md): technologies, versions, commands.
- [CODE.md](docs/conventions/CODE.md): code, frontend, dependencies.
- [GOLANG.md](docs/conventions/GOLANG.md): Go work.
- [TDD.md](docs/conventions/TDD.md): behavior changes, tests, verification.
- [COMMITS.md](docs/conventions/COMMITS.md): commit proposals and delivery.
- [Decision guide](docs/decisions/README.md): architecture choices.
- [Feature guide](docs/features/README.md): planning, work, resume.
- Load historical records only for relevant dependencies.

## Replies

- Use short bullets for changes and verification; avoid repetition.
- End with **User must be aware**: material risks, limits, unresolved decisions,
  or "None".
- Record actual open questions where relevant; omit unknown topics.
