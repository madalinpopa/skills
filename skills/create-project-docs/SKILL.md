---
name: create-project-docs
description: Initializes or refreshes a project's AGENTS.md, CLAUDE.md, and documentation for conventions, product specs, decisions, and features. Use when starting a project or requesting this documentation foundation, especially for Go and TypeScript. Not for instructions-only setup (use create-agents-setup), application scaffolding, or completing an existing product specification.
compatibility: The copy helper requires Bash and standard Unix tools.
status: published
tags: [setup, docs, conventions, go, typescript]
---

# Create project docs

Initialize or refresh a project's documentation foundation. Preserve project
facts and custom rules; record only repository evidence and user decisions.
Keep recurring instructions short and route detailed workflows to skills.

## Prepare

Find the project root and this skill's installed directory, `<skill-dir>`.
The project root must exist and must not be the skill directory. Inspect existing
instructions, relevant docs, manifests, task scripts, and the diff. Reuse known
answers. If `.jj` exists, use `use-jj` for version-control work, including when
`.git` also exists; report a missing prerequisite before mutations.

On the first invocation in a project, use `use-skills-cli` to sync the store and
install `use-task-spec`, `use-task-plan`, and `use-tasks-checklists` for the
current agent in project scope. On subsequent invocations, sync and compare
installed versions with the store, updating outdated skills through the CLI.
Use configured targets and lock files; do not assume the default install path.
If `use-skills-cli` is not discoverable, follow
[skill setup](references/skill-setup.md). Never force an update, overwrite local
skill edits, or copy unpublished checkout content into install targets. Report
missing/unpublished skills or conflicts and stop dependent setup.

Installing the skills does not require loading all three. Read only the one
needed for the requested work. This setup request authorizes these workflow
skill installations; application dependency changes follow project policy.

## Initialize missing documents

If another path already owns the same topic, resolve ownership with the user
before creating a duplicate. Run:

```sh
bash "<skill-dir>/scripts/init-project-docs.sh" "<project-root>"
```

The helper creates missing assets and keeps every existing file, including empty
ones. It refuses symlinks, file/directory clashes, and competing layouts before
writes. Use `--allow-existing-layout` before the root only after the user chooses
to retain both layouts. Later write failures can leave partial output: inspect
reported paths and fix the cause before rerunning; never delete user content.

Tailor created documents using established facts. For existing documents, apply
[refresh rules](references/refresh.md); the copy helper never replaces them.

## Document ownership

- `AGENTS.md`: project identity, essential permissions, core slice workflow, and task routes.
  Keep it near 300 words. Preserve narrow searches and selective context loading.
- `CLAUDE.md`: `@AGENTS.md`; retain Claude-specific additions and no reverse import.
- `docs/conventions/`: shared code/security, commits, Go, test-specific TDD,
  detailed review procedures in REVIEW.md, prototypes, and question-log guidance.
  Keep the controlling workflow in AGENTS.md; retain review gates and prototype exceptions.
- `docs/specs/`: mission, stack, and roadmap. Use the
  [writing guide](assets/docs/specs/README.md) only for the document being authored.
- `docs/features/`: numbered records; the
  [feature index](assets/docs/features/README.md) routes authoring and execution.
- `docs/decisions/`: significant agreed choices under the
  [decision guide](assets/docs/decisions/README.md).

Route feature requirements to `use-task-spec`, approach to `use-task-plan`, and
checklists/execution/resumption to `use-tasks-checklists`. Keep actual project
requirements and records in project docs; do not copy skill instructions there.
Reading an existing spec does not invoke its authoring workflow.

Spec files and `docs/QUESTIONS.md` start empty unless real content is established;
add no empty outlines, examples, or invented requirements. Use the next unused
four-digit feature ID across ROADMAP.md and folders, preserving existing IDs.
ROADMAP.md contains only feature/sub-feature headings and task checklists with
concise links; progress needs evidence and does not imply acceptance.

Tailor engineering emphasis to the project. Keep applicable Go/TypeScript rules;
remove inapplicable frontend rules and mark unrelated language guides accordingly.
Record actual technologies, versions, official documentation links, verification
dates, and exact commands/directories in STACK.md. Compare selected and latest
versions without upgrading. Do not invent choices or record secrets; discovery
of a command does not prove it ran.

For requested mission/stack spec writing, ask focused option-based questions
using the available question tool, or chat if unavailable. Reuse settled answers;
wait before writing content dependent on unanswered decisions. Feature-spec
questions belong to `use-task-spec`. Record real unresolved questions in
`docs/QUESTIONS.md` using the
[question guide](assets/docs/conventions/QUESTIONS.md); apply answers before
checking them. Initialization does not require an interview to fill every field.

## Verify and report

Check relative links, skill routes and availability, CLAUDE import, document
ownership, real commands, empty unknown specs, and applicable review gates.
For refreshes, inspect the diff to confirm project content and custom rules are
preserved. Check copied asset paths and `.jj` routing. Copying docs proves no
application tests passed. Distinguish word counts from measured tokens or model
performance; load only relevant skills and reference sections.

Report changed paths, checks, failures, partial writes, and limits concisely.
Propose a commit message; commit or publish only on explicit request. End with
one concrete **Next Step**, linked to the relevant task or question when available.
