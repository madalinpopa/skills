# Feature: Project documentation setup

Initialize a new project's agent instructions and documentation foundation.
The primary audience builds Go backends with TypeScript frontends; frontend
frameworks and product requirements remain project choices.

- **Active phase**: [03: Task workflow skills](phase-03-task-workflow-skills.md)
- **Feature acceptance**: Earlier work accepted for commit on 2026-10-14;
  phase 03 awaits implementation review and has no commit authorization.

## Using these documents

Read repository instructions, this overview, and the active phase on resume.
`docs/SPEC.md` owns this repository's CLI behavior and skill packaging contract.
This feature adds catalog content, without changing the CLI. Generated project
documents describe the target project, not this skill store.

## Scope

- **Included**: A new `skills/create-project-docs/` skill, bundled Markdown
  templates, and a Bash helper that copies missing project documents.
- **Non-goals**: Application scaffolding, dependency installation, choosing a
  frontend framework, writing an invented product spec, or publishing changes.
- **Constraints**: Keep generated `AGENTS.md` concise. `CLAUDE.md` imports
  `AGENTS.md`. Preserve existing files and project conventions on rerun.
- **Spec references**: [Store layout](../../SPEC.md#layout),
  [Skill metadata](../../SPEC.md#skill-metadata).

## Acceptance criteria

- [ ] Setup creates `AGENTS.md` and `CLAUDE.md`, plus the requested convention
  and specification documents, in the selected project root.
- [ ] Agent instructions route by task rather than loading every document.
- [ ] Templates separate product intent, observable behavior, stack choices,
  roadmap, engineering conventions, decisions, and feature execution.
- [ ] Go and TypeScript are defaults to discuss; React and other frameworks
  are selected only from project evidence or the user's choices.
- [ ] Reruns preserve existing bytes and report kept files. Unsafe destination
  paths fail before copying; later write failures report partial creation.
- [ ] Unknown project facts remain explicitly pending. No fictitious accepted
  decision or approved feature is created.
- [ ] Validation and focused helper checks pass; the developer reviews the result.

## Agreed decisions

| Decision | Reason | Approval or pending question |
| --- | --- | --- |
| Create a new skill for project initialization | User requested a new skill and a broader documentation structure | User request, 2026-10-14 |
| Name it `create-project-docs` | Distinguishes documentation setup from application scaffolding | Approved: "yes, proceed", 2026-10-14 |
| Use `TDD.md` and `STACK.md` | Consistent Markdown filenames; interprets `TDD.MD` and `STACKmd` | Approved: "yes, proceed", 2026-10-14 |
| Treat numbered paths as record naming patterns | Avoid invented project decisions and features | Approved: "yes, proceed", 2026-10-14 |
| Keep existing setup/spec skills unchanged | Their current narrower workflows remain independently useful | Approved: "yes, proceed", 2026-10-14 |
| Run Go fixes periodically; verify latest library versions and require approval for dependencies | User's project engineering policy | Explicit user request, 2026-10-14 |
| Prioritize readable tests, simple reusable helpers, and separate success/failure tables | User's test conventions | Explicit user request, 2026-10-14 |
| Numbered feature folders contain plan.md, requirements.md, and validation.md | Separate requirements, implementation planning, and evidence; replaces the generated FEATURE.md convention | Explicit user request, 2026-10-14; validation..md normalized to validation.md |
| All work proceeds in small slices with review between slices | Keep changes easy to review, including docs, configuration, and refactoring | Explicit user request, 2026-10-14; follow-up awaiting review |
| Split oversized features into sub-features with their own three-file records | Keep scope and outcomes independently reviewable | Explicit user request, 2026-10-14; follow-up awaiting review |
| Optimize context in small slices without losing instructions | Start with recurring AGENTS.md context, then conventions and feature/spec assets | Explicit user request, 2026-10-14; slice 1 awaiting review |
| Gather missing documentation facts with the available question tool, starting with the minimum | Reuse known facts, ask only for the next slice, and leave non-blocking unknowns pending | Explicit user request, 2026-10-14; follow-up awaiting review |
| Require independent adversarial review of code slices before user review | Resolve actionable findings within approved scope; re-review fixes; preserve test approval and dependency gates | Explicit user request, 2026-10-14; follow-up awaiting review |
| Keep replies minimal, with a User must be aware section and recorded pending decisions | Short change/check bullets, no repeated details; ask the user to review unresolved questions | Explicit user request, 2026-10-14; follow-up awaiting review |
| Use phased task-only plans for features and codebase work | Clear progress, one next action, linked evidence, and less duplicated context | Explicit user request, 2026-10-15; awaiting documentation review |
| Centralize actual open questions in docs/QUESTIONS.md | One visible checklist; document answers in their owning sections before closing questions | Plan approved: "arata bine, continua cu implementarea", 2026-10-15; test review remains required |
| Replace the mandatory awareness footer with one concrete Next Step | Record questions and decision blockers as they arise in QUESTIONS.md; propose the next useful action after each work slice | Explicit user request, 2026-10-15; included in phase 02 |
| Keep question checklists simple and grouped by purpose | Remove question IDs, HTML anchors, and mandatory metadata; document answers before checking tasks | Explicit user correction, 2026-10-15; documentation updated for review |

## Phase roadmap

| Phase | Outcome | Depends on | Status | Phase file |
| --- | --- | --- | --- | --- |
| 01 | A reusable initializer with templates and preservation checks | None | awaiting implementation review | [Phase 01](phase-01-initializer-skill.md) |
| 02 | One question log with document-first resolution and concise persistent guidance | Existing initializer | awaiting implementation review | [Phase 02](phase-02-central-question-log.md) |
| 03 | Extract task skills and refresh reusable project instructions safely | Existing documentation assets | awaiting implementation review | [Phase 03](phase-03-task-workflow-skills.md) |

Acceptance does not authorize commits or publication.
