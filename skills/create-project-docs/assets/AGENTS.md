# Agent instructions

## Project

## Working rules

When the repository root contains a `.jj` directory, the project uses jj-vcs
(Jujutsu). Invoke `use-jj` for version-control work, even if `.git` also exists.
If the skill is unavailable, report it before version-control mutations; do not
substitute Git mutations.

Work in small slices with one outcome each, including docs, configuration, and
refactoring. Verify each slice and present it for user review before the next.
Stay in scope. Inspect relevant code and diffs; preserve unrelated work. Raise
conflicts between instructions, requirements, and code before changing a contract.

Never invent requirements. Fill specs progressively; leave unknown sections
empty, without placeholders or generic questions. When the user asks to fill or
work on specs or a specific section, always use the available ask-user-question
tool with concrete proposed options based on repository evidence. Ask focused
rounds, reuse prior answers, and wait before writing dependent content. If no
question tool is available, explain that and ask option-based questions in chat.

Use the commands and working directories in STACK.md; report missing ones.
Check current official docs and the latest released library versions. Adding,
installing, replacing, or upgrading a dependency needs user approval; a version
check is not approval.

For each change, follow [TDD's workflow](docs/conventions/TDD.md#workflow):
write a failing test and see it fail, implement, then have a different agent do
an adversarial review; fix and repeat until clean. Then the user reviews and
approves the commit. Commit only on that approval; publish only on explicit
request; never add AI attribution.

Reply in short `-` bullets for changes and verification, without repeats. End
with **User must be aware**: material risks, limits, and unresolved decisions, or
"None". Record actual open questions when relevant; leave other sections empty.

## Read only for the current task

| Task | Guidance |
| --- | --- |
| Goals, users, scope | [MISSION.md](docs/specs/MISSION.md) |
| Feature outcomes and requirement links | [FEATURES.md](docs/specs/FEATURES.md) |
| Features, sub-features, and task checklists | [ROADMAP.md](docs/specs/ROADMAP.md) |
| Technologies, versions, commands | [STACK.md](docs/specs/STACK.md) |
| Code, TypeScript/frontend, dependencies | [CODE.md](docs/conventions/CODE.md) |
| Go work | [GOLANG.md](docs/conventions/GOLANG.md) |
| Behavior changes, tests, verification | [TDD.md](docs/conventions/TDD.md) |
| Commit proposals, delivery | [COMMITS.md](docs/conventions/COMMITS.md) |
| Architecture choices | [Decision guide](docs/decisions/README.md) |
| Feature planning, work, or resume | [Feature guide](docs/features/README.md) |

Use task/subtask checklists for roadmaps and plans; mark completed work with
evidence. Write specs as short bullets under clear sections; reserve tables for
compact mappings. FEATURES.md catalogs outcomes and requirement links, without
progress. ROADMAP.md contains only feature/sub-feature groups and task checklists;
leave it title-only when no work is identified. Requirements own behavior, plans
own execution details and the checkpoint, and validation owns evidence and user
acceptance. Link shared facts rather than duplicating them.
Load old feature records only for a relevant dependency.
