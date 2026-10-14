---
name: create-project-docs
description: Initializes a project's AGENTS.md, CLAUDE.md, and documentation for conventions, product specs, decisions, and features. Use when starting a project or requesting this documentation foundation, especially for Go and TypeScript. Not for instructions-only setup (use create-agents-setup), application scaffolding, or completing an existing product specification.
compatibility: The copy helper requires Bash and standard Unix tools.
status: published
tags: [setup, docs, conventions, go, typescript]
---

# Create project docs

Copy a documentation foundation into a project, then fill in what the user and
the repository already make clear. Leave the rest as `Pending:`.

## Inspect and initialize

Find the project root and this skill's installed directory, `<skill-dir>`. The
root must already exist and is never the skill's own directory. Read existing
instructions, README, docs, manifests, task scripts, and the diff when there is
one. In an empty project, work from the user's description.

If existing docs already cover the same topic at another path, agree with the
user on which file owns it before adding a second one. Setup does not authorize
moving or rewriting existing docs.

Run the helper by its absolute path:

```sh
bash "<skill-dir>/scripts/init-project-docs.sh" "<project-root>"
```

It copies missing files from [assets](assets/), prints `created` or `kept` for
each path, and never replaces an existing file, even an empty one. Before
writing anything, it refuses symlinks, file/directory clashes, and another docs
layout (`docs/SPEC.md`, `docs/templates/FEATURE.md`, `docs/templates/PHASE.md`,
or `docs/features/*/FEATURE.md`). Add `--allow-existing-layout` before the root
only when the user chooses to keep both layouts. A failure after writing starts
can leave directories or a partial file: check the reported paths and fix the
cause before rerunning. Never delete or replace content to make a rerun pass.

## Tailor the new documents

Edit only files the helper reported as `created`. Change an existing document,
such as an instruction file that needs new routes or a CLAUDE import, only
after the user approves that exact edit. Keep Claude-only additions.

| File | Owns |
| --- | --- |
| [AGENTS.md](assets/AGENTS.md) | Project identity, working rules, and task-based routes |
| [CLAUDE.md](assets/CLAUDE.md) | One `@AGENTS.md` import; no reverse import |
| [CODE.md](assets/docs/conventions/CODE.md) | Shared code rules and dependency approval; TypeScript/frontend rules when used |
| [COMMITS.md](assets/docs/conventions/COMMITS.md) | Commit messages and when a commit is allowed |
| [GOLANG.md](assets/docs/conventions/GOLANG.md) | Go conventions, toolchain, and `go fix` |
| [TDD.md](assets/docs/conventions/TDD.md) | Test-first workflow, review loop, and verification |
| [MISSION.md](assets/docs/specs/MISSION.md) | Problem, users, outcome, and non-goals |
| [FEATURES.md](assets/docs/specs/FEATURES.md) | Feature index; each feature's `NNNN` number is its only ID |
| [ROADMAP.md](assets/docs/specs/ROADMAP.md) | Delivery order and the only record of feature status |
| [STACK.md](assets/docs/specs/STACK.md) | Technologies, supported versions, and project commands |
| [Decision guide](assets/docs/decisions/README.md) | Numbered decision records and their template |
| [Feature guide](assets/docs/features/README.md) | Numbered feature folders and their three templates |

Setup is done when:

- AGENTS.md states the project's purpose, or keeps it `Pending:`.
- STACK.md lists the languages and the real commands and working directories
  found in manifests and scripts; missing ones stay `Pending`.
- Go or TypeScript/frontend guidance that does not apply is marked not
  applicable, and frontend rules are removed from Go-only projects.
- Every other unknown stays `Pending:`. Fill MISSION, FEATURES, and ROADMAP
  only with facts the user already gave.

Use repository evidence and earlier answers first. Ask at most one short round
of questions, and only for what the repository cannot show, with the
ask-user-question tool or in chat. Go and TypeScript are likely choices, not
defaults: confirm them, and never pick a framework, storage, architecture,
version, dependency, or deployment just to fill a field. Check current official
docs only for technologies being recorded, and put links and versions in
STACK.md. Do not create example decisions or features.

When tailoring, keep these template rules, adapted only to existing project
policy:

- Ask only the minimum questions; record other unknowns as pending.
- Dependency changes need user approval; always check the latest release.
- The TDD workflow: failing test, implementation, then a review by a different
  agent repeated until clean, and one user review that approves the commit.
  This describes a reviewer role; it does not install agents or tools.
- Prefer integration coverage for behavior; add small unit tests only for distinct
  implementation gaps, without duplicating scenarios across layers.
- Roadmaps and plans use task and subtask checklists with short descriptions.
  Mark completed items `[x]` with evidence; keep pending or blocked items `[ ]`.
  Feature plans group detailed tasks into related phases of one to five planned
  commits each, with concrete subtasks and proposed commit messages. Roadmap
  feature acceptance and plan task completion have separate owners. Split roadmap
  features into small, independently verifiable slices with scope, dependencies,
  acceptance criteria, and links to plan tasks and evidence.
- Specs use clear sections and concise bullets, linking shared facts to their
  owner. Use tables only for compact comparisons or mappings, never long prose.
- The reply format with **User must be aware**.
- `go fix` guidance in GOLANG.md for Go projects.

Keep AGENTS.md near 300 words with plain Markdown links; do not import every
linked document.

Writing an existing `docs/SPEC.md` belongs to `create-spec` when available; do
not use it to rewrite this layout. Application scaffolding is a separate
request. This skill does not require either.

## Verify and report

Check relative links, task routes, the CLAUDE import, that commands are real,
that each topic has one owner, and that unknowns are marked pending. Check
that plans and roadmaps use checklists, feature phases contain detailed tasks
and one to five planned commits, roadmap features have small reviewable slices,
completion claims have evidence, spec text uses
sections and bullets, and test guidance avoids duplicate coverage.
Confirm existing files kept their bytes. Copied templates are not agreed requirements,
and copying docs does not mean any application tests passed.

Report in short `-` bullets, without full logs. Add **User must be aware** with
real limits, partial writes, and open questions, or "None". Record open items
in the relevant docs and ask the user to review them. Propose a commit message,
and leave changes uncommitted unless the user asks for a commit. Publishing
this skill does not authorize committing, pushing, or publishing the project.
