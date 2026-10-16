---
name: create-project-docs
description: Initializes a project's AGENTS.md, CLAUDE.md, and documentation for conventions, product specs, decisions, and features. Use when starting a project or requesting this documentation foundation, especially for Go and TypeScript. Not for instructions-only setup (use create-agents-setup), application scaffolding, or completing an existing product specification.
compatibility: The copy helper requires Bash and standard Unix tools.
status: published
tags: [setup, docs, conventions, go, typescript]
---

# Create project docs

Copy a documentation foundation into a project, then fill in what the user and
the repository already make clear. Spec files start empty; add headings only
with real content. Omit placeholders, filler questions, and example entries.
Specs grow progressively as the user works with the agent.

## Inspect and initialize

Find the project root and this skill's installed directory, `<skill-dir>`. The
root must already exist and is never the skill's own directory. Read existing
instructions, README, docs, manifests, task scripts, and the diff when there is
one. In an empty project, work from the user's description.

Check for a `.jj` directory at the repository root. Its presence means the user
uses jj-vcs (Jujutsu), even when `.git` also exists. Invoke the `use-jj` skill
before version-control operations and follow it for history, changes, bookmarks,
and remotes. If the skill is unavailable, report the missing prerequisite before
version-control mutations; do not substitute Git mutations.

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
| [AGENTS.md](assets/AGENTS.md) | Project identity, engineering emphasis, context discipline, and task-based routes |
| [CLAUDE.md](assets/CLAUDE.md) | One `@AGENTS.md` import; no reverse import |
| [CODE.md](assets/docs/conventions/CODE.md) | Shared design and security rules, dependency approval, and applicable frontend rules |
| [COMMITS.md](assets/docs/conventions/COMMITS.md) | Commit messages and when a commit is allowed |
| [GOLANG.md](assets/docs/conventions/GOLANG.md) | Go conventions, toolchain, and `go fix` |
| [TDD.md](assets/docs/conventions/TDD.md) | Test-first workflow, review loop, and verification |
| [Question guide](assets/docs/conventions/QUESTIONS.md) | Simple questions grouped by purpose, answers, and closure in docs/QUESTIONS.md |
| [Spec writing guide](assets/docs/specs/README.md) | Section guidance and concise writing rules |
| [MISSION.md](assets/docs/specs/MISSION.md) | Problem, users, outcome, and non-goals |
| [ROADMAP.md](assets/docs/specs/ROADMAP.md) | Features and sub-features grouped with task checklists; no other content |
| [STACK.md](assets/docs/specs/STACK.md) | Technologies, supported versions, and project commands |
| [Decision guide](assets/docs/decisions/README.md) | Numbered decision records and section-writing guidance |
| [Feature guide](assets/docs/features/README.md) | Numbered feature folders and record-writing guidance |

Tailor engineering emphasis to the user's priorities and repository evidence.
For a Go web application, this may emphasize maintainable Go backends, secure
trust boundaries, and consistent UI/API/persistence behavior. Do not impose Go
or full-stack work on other projects. Keep only the short emphasis and task
routes in AGENTS.md; shared security and frontend rules belong in CODE.md,
Go mechanics in GOLANG.md, and tests/review gates in TDD.md. STACK.md records
selected technologies, not engineering instructions. Add a separate security
guide or skill only when substantial project-specific procedures justify it;
link it with an explicit read trigger rather than loading it for every task.

Setup is done when:

- AGENTS.md states the project's purpose when known; otherwise its Project
  section stays empty.
- When `.jj` exists at the repository root, newly created AGENTS.md records that
  the project uses jj-vcs and requires the `use-jj` skill for version-control
  work. Existing instruction files still require approval before editing.
- STACK.md records only technologies, versions, commands, and working directories
  supported by manifests, scripts, or user decisions. Omit unknown topics.
- Go or TypeScript/frontend guidance that does not apply is marked not
  applicable, and frontend rules are removed from Go-only projects.
- Spec files are empty unless established facts can be recorded. Do not seed
  titles, section outlines, or writing instructions in them. Add sections only
  with content, following the [spec writing guide](assets/docs/specs/README.md).
  Feature and decision guides also describe sections instead of empty templates.
- Assign each real feature the next unused four-digit number, `NNNN`, as its
  only ID across ROADMAP.md and feature records. Check both for used numbers;
  never reuse or renumber IDs. Detailed behavior belongs in requirements.md.
- ROADMAP.md contains only feature headings, nested sub-feature headings, and
  task/subtask checklists, with concise links to plans or evidence when useful.
  Do not add status tables, narrative, acceptance criteria, dependency notes,
  or workflow instructions. Those belong in feature records or guides.
- Read README, prototype files, existing docs, and earlier user answers for clear
  features and tasks. Record only work those sources support; a prototype is
  evidence of possible scope, not proof of agreed requirements or completion.
  If no feature or task is clear, leave ROADMAP.md empty. Do not create example
  decisions or features.

Use repository evidence and earlier answers first. Initialization does not need
an interview to fill every section. Never choose a language, framework, storage,
architecture, version, dependency, or deployment just to fill a field. Check
current official docs only for technologies being recorded; put verified
versions, official links, and verification dates in STACK.md. Compare selected
and latest released versions without upgrading dependencies automatically.
Record exact commands and working directories only when found; discovering a
command does not mean it was run. Document configuration names, never secrets.

## Fill specs progressively

Whenever the user asks to fill, refine, or work on specs or a specific section,
always use the available ask-user-question tool with concrete proposed options
to guide the work. This also applies when such work is requested during setup.
Read the relevant docs and repository evidence first, then ask a small, focused
round about the requested section. Base options on that evidence, explain
meaningful tradeoffs, and allow the user to supply another answer. Do not ask
again for decisions already made or turn one section into a whole-project
interview. Wait for answers before writing content that depends on them; keep
unanswered topics out of the document. If no question tool is available, report
that limit and ask equivalent option-based questions in chat.

Update only the requested scope with established facts and the user's answers.
Keep existing agreed content and document ownership. Record actual unresolved
questions and blockers needing user input in `docs/QUESTIONS.md` as they arise,
with the input needed and affected work. Reuse an existing entry instead of
duplicating it; do not seed generic questions.

Before recording questions or applying answers, read the
[question guide](assets/docs/conventions/QUESTIONS.md). It owns the log format
and lifecycle: group simple task questions by purpose, apply answers to their owning sections,
and check tasks only after verifying the updates. Partial or unapplied answers
stay open. Use no question IDs or HTML anchors; link other documents to the
relevant group in the log.

When tailoring, keep these template rules, adapted only to existing project
policy:

- Omit unknown topics and use guided question rounds for requested spec work.
- Dependency changes need user approval. Check current official docs for external
  APIs being changed and latest releases for dependencies being proposed.
- Route behavior changes, test work, and verification to TDD.md; keep its workflow,
  exceptions, and review gates there rather than repeating them in AGENTS.md.
  User review does not authorize a commit without an explicit request.
- Preserve TDD.md's throwaway prototype workflow for explicit disposable
  experiments and requests to continue existing prototypes identified by repository
  evidence. Route to it briefly from AGENTS.md; keep the local-commit exception in
  COMMITS.md and the formal-record exemption in the feature guide. Prototype work
  skips TDD, automated test writing, full verification, adversarial review, and
  intermediate planning/review gates; it uses isolated version control and a minimal
  demonstration check. Preserve normal permissions for dependencies and external
  actions. Production implementation returns to the normal workflow.
- Route behavior changes to the affected feature contract and acceptance criteria.
  Distinguish consuming requirements from authoring specs: only authoring loads
  writing rules and the target document's section. Keep guided spec questions
  mandatory as described above; reuse settled answers.
- Prefer integration coverage for behavior; add small unit tests only for distinct
  implementation gaps, without duplicating scenarios across layers.
- Roadmaps contain feature/sub-feature headings and task checklists. Feature
  records also cover fixes, refactoring, tooling, and documentation work.
  Follow the [feature guide](assets/docs/features/README.md) for ownership and
  planning: plan.md contains only phase headings and task/subtask checklists,
  with evidence links and one next action. Keep phases to one to five planned
  commits. Completion does not imply acceptance or commit permission.
- Group related spec content under clear sections and use concise lists.
  Link shared facts to their owner. Avoid long paragraphs and prose-heavy tables;
  reserve tables for compact comparisons or mappings.
- Short replies ending with one concrete **Next Step**, linked to the relevant
  question or task when available; proposals do not authorize new work.
- Whenever user input is needed, use the available ask-question tool when it
  supports the question; otherwise ask in chat. Apply this beyond spec work
  too, and reuse known answers instead of asking again.
- `go fix` guidance in GOLANG.md for Go projects.

Keep AGENTS.md near 300 words, using section headings and short bullet lists
only. Make task routes explicit: say when to read each guide, rather than only
listing its subject. Retain narrow searches and task-relevant reference loading
as context rules; do not duplicate detailed guides in the entrypoint.
Name the applicable section and trigger: TDD workflow before changes/checks,
coverage before test design, and feature Resume for existing work. Workflow
routes must preserve required reviews and the explicit prototype exception;
selective reading never waives an applicable approval or review gate.
A link or anchor locates guidance but does
not limit a tool's output: retrieve the applicable section and its dependencies.
Keep detailed writing, roadmap, record, and review rules in
their owning guides. Keep setup-only instructions, such as starting spec files
empty, in this skill; persistent instructions should require established content
without placeholders or empty sections. Avoid prose paragraphs and tables. Use
plain Markdown links; do not import every linked document.

Writing an existing `docs/SPEC.md` belongs to `create-spec` when available; do
not use it to rewrite this layout. Application scaffolding is a separate
request. This skill does not require either.

## Verify and report

Check `.jj` detection and the `use-jj` route when applicable. Check relative
links, task routes, the CLAUDE import, that commands are real,
that each topic has one owner, and that unknown specs are empty files. Check
that plans contain only phase headings and task/subtask checklists with a clear
next action, phases respect the feature guide, roadmaps contain only feature and
sub-feature groups with task checklists, completion claims have evidence, spec
text uses concise sectioned lists without empty headings, and test guidance
avoids duplicate coverage.
Check that behavior changes route to actual requirements, spec readers do not
activate authoring procedures, and resumption selects the active phase and
latest relevant evidence. Check that section routes retain all applicable
approval/review requirements. Check that new and existing throwaway prototypes
reach the exception, ambiguous prototype filenames do not bypass normal work,
and production implementation restores the normal workflow.
Distinguish word counts from token measurements
or measured model performance.
Confirm existing files kept their bytes. Copied templates are not agreed requirements,
and copying docs does not mean any application tests passed.

Report changes and verification in short `-` bullets without repetition or full
logs. Report material limits and partial writes alongside the relevant result.
Record actual questions and blockers needing user input in `docs/QUESTIONS.md`
during work rather than waiting for the final reply. Whenever user input is
needed, use the available ask-question tool when it supports the question;
otherwise ask in chat. End with **Next Step**:
one concrete action based on current progress, such as answering a blocking
question, reviewing changes, or working on a feature. Link the relevant question
or task when available; do not invent follow-up work or imply permission to start
new scope. Propose a commit message,
and leave changes uncommitted unless the user asks for a commit. Publishing
this skill does not authorize committing, pushing, or publishing the project.
