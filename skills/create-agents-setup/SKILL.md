---
name: create-agents-setup
description: Initializes concise AGENTS.md and a CLAUDE.md import, with project conventions and workflow-mentor routing. Use when either instruction file is missing or the user requests project instruction setup. Not for implementing features or writing the project specification.
compatibility: The copy helper requires Bash and standard Unix tools. Missing skills use use-skills-cli when available.
status: published
tags: [setup, conventions, docs, workflow]
---

# Create agents setup

Give agents a small project overview and explicit routes to information for
their current task. Shared rules live in `AGENTS.md`; `CLAUDE.md` imports it.
This skill creates instructions, not application code or feature plans.

## Inspect once

Resolve the user's project root and this installed skill's absolute directory
as `<skill-dir>`. Read existing instructions, relevant `docs/SPEC.md` sections,
README, task scripts, language manifests, and the current diff. Inspect only
enough code to establish layout and conventions. Preserve unrelated work.

Check both root instruction files, including empty files, and applicable
overrides or nested rules. Reuse existing conventions and feature records.
If both files already work and no revision was requested, report completion
and return. Before writing, check selected destinations and parent directories:
stop setup on symlinks (including dangling links), non-regular file targets,
or file/directory collisions. The helper does not enforce these guards.

## Create and tailor

Run from the user's project, using the helper's absolute path:

```sh
"<skill-dir>/scripts/init-agent-docs.sh" <project-root>
```

It copies missing [AGENTS.md](assets/AGENTS.md) and
[CLAUDE.md](assets/CLAUDE.md), keeping existing files. When only CLAUDE.md is
missing, its import is normally the whole change; do not scaffold unused docs.
When AGENTS.md is new, derive its rules from an existing CLAUDE.md if present.
Move shared rules out of CLAUDE.md only when that migration is authorized.

For a new AGENTS.md, copy applicable templates into `docs/conventions/`,
keeping their filenames. Reuse equivalent existing docs instead of duplicating
them; never overwrite a convention document. Copy assets with ordinary file
tools, then adapt only the new files to the project.

| Template | Include / read when |
| --- | --- |
| [TDD](assets/conventions/tdd.md) | Behavior changes or test review |
| [Commits](assets/conventions/commits.md) | Commit proposals or delivery |
| [Code](assets/conventions/code.md) | Code changes or review |
| [Go](assets/conventions/go.md) | A Go module exists; Go work only |
| [Workflow](assets/conventions/workflow.md) | Mentoring or feature work |

Fill project identity, command placeholders, and information paths from
evidence. Include working directories for commands in multi-module projects.
Remove unused rows and placeholders. Mark unknown facts pending rather than
inventing them; ask only when they block setup. Adapt the defaults to existing
review gates and commit policy without weakening or adding conflicting rules.

Keep `docs/SPEC.md` as the behavior reference when present. If absent, say it
is missing and name the current evidence; do not manufacture a spec. Route
requested spec authoring to `create-spec` when available. Reuse existing
feature paths, including `docs/features/<name>/FEATURE.md` and its active phase.
Otherwise use `docs/feature.md` for future feature work. Do not create an empty
plan during setup or claim an absent file exists.

Keep only relevant skill triggers: `workflow-mentor`, Go guidance for Go,
`use-jj` for jj, and `use-gh` for GitHub work. Inspect session skill metadata;
do not load every skill body or hard-code machine-specific installation paths.
For a required missing skill, use `use-skills-cli`; if unavailable or resolution
fails, report it and stop only dependent work. Do not install optional skills
preemptively. Setup never invokes workflow-mentor; return to its caller or
finish after reporting the files, avoiding a dependency cycle.

Existing files need targeted edits, not template replacement. Within authorized
setup scope, add a missing mentoring trigger or CLAUDE import after checking
for conflicting rules. Preserve Claude-specific additions. Empty files need
tailoring even though the helper reports them as kept. If a conflict needs a
user decision, show the concrete proposed edit and pause only that edit. A
rerun must not append duplicate sections or imports.

## Keep recurring context small

Aim for roughly 300 words in a new AGENTS.md; correctness takes priority.
Keep project facts, commands, essential limits, and task-based links there.
Detailed conventions belong in ordinary Markdown documents, not auto-loaded
instruction files. Do not import them with `@` or require every link at
startup. Never link AGENTS.md back to CLAUDE.md.

This follows [Codex instruction discovery](https://learn.chatgpt.com/docs/agent-configuration/agents-md)
and [Claude memory guidance](https://code.claude.com/docs/en/memory): project
instructions consume recurring context; Claude imports load their contents.
Use one import to share rules and ordinary links for conditional detail.

## Verify and return

Review text and diff. Check destinations, resolved links, command definitions,
placeholders, relevant skill availability, and conflicting instructions.
Distinguish command discovery from execution. Confirm the CLAUDE import has
no cycle and existing content was preserved. Review Go, non-Go, and existing-file
paths. Compare related workflows: keep each procedure in one owner and link
to it; retain only essential reminders in recurring context.

Report created, edited, and kept paths, missing dependencies, checks, and a
proposed commit message. On failure inspect outputs and report partial
creation; the helper is not transactional. Leave changes uncommitted.
