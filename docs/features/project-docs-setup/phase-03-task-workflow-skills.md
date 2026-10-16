# Phase 03: Extract task workflow skills

- **Feature**: [Project documentation setup](FEATURE.md)
- **Outcome**: Three reusable task skills own feature requirements, technical
  planning, and execution checklists; project setup installs and routes to them.

## Current checkpoint

- **Status**: awaiting implementation review
- **Approved scope**: User approved the proposed extraction, CLI setup/update,
  narrow instruction refresh, and validation on 2026-10-16: “approve”.
- **Working tree**: Three new catalog skills; create-project-docs entrypoint,
  references, and routing assets changed. Existing prompt additions preserved.
- **Test gate**: Documentation-only instructions; no helper or CLI behavior changed.
  Failing wording tests would not provide useful evidence.
- **Next action**: Review the extraction and validation results; no commit authorized.

## Scope and acceptance

- [x] Extract the existing spec/plan/tasks ownership and writing rules into
  use-task-spec, use-task-plan, and use-tasks-checklists.
- [x] Incorporate prompt requirements: interview-before-draft when requested,
  repository-grounded planning, verifiable ordered tasks, bounded execution,
  and optional implementation notes.
- [x] Replace duplicated feature guidance with task-specific skill routes.
- [x] Document first-use CLI installation and subsequent version checks/updates;
  stop dependent setup on unavailable skills or local conflicts.
- [x] Refresh only reusable project instructions; preserve project facts,
  custom rules, records, and ambiguous passages.
- [x] Keep the copy helper and existing prompts unchanged.
- [x] Complete validation and record evidence below.
- [ ] Obtain explicit developer acceptance.

## Design decisions

- Existing numbered spec.md/plan.md/tasks.md ownership is retained. The source
  prompts are inputs, not files to rewrite.
- Entrypoints hold normal authoring paths. Execution and refresh/setup fallback
  details load only for those tasks. Generated instructions route to skills.
- Lost in the Middle (Liu et al., TACL 2024):
  https://aclanthology.org/2024.tacl-1.9/ reports positional sensitivity in long
  retrieval contexts. Selective loading is our design inference, not a measured
  model-performance improvement or a universal token limit.
- Copy-helper preservation behavior remains unchanged. Instruction refresh is
  an agent-guided merge, not automated replacement of existing documents.

## Verification evidence

- All four source skills pass the installed create-repo-skill validator using
  `uv run --offline --with pyyaml bash .../validate-skill.sh <skill-path>`.
  Initial direct runs failed because system Python lacks PyYAML; the cached
  offline uv environment resolved this without downloads.
- `python3 -m unittest discover -s skills/create-project-docs/tests -v`: all
  seven existing initializer checks pass, including byte preservation on rerun.
  The copy helper and tests are unchanged; generated Python cache removed.
- Checked 37 relative links/anchors across the 12 changed Markdown files;
  `git diff --check` passes.
- Reviewed routes for spec authoring versus consumption, planning without
  implementation, checklist creation versus execution, next-task versus
  all-tasks execution, blocked prerequisites, and project-policy preservation.
  This is an instruction review, not independent agent execution.
- Word counts: create-project-docs entrypoint 2,096 → 779; feature guide
  792 → 170; generated AGENTS.md 323 → 335; spec guide 515 → 521.
  New skill entrypoints are approximately 250–375 words each. Counts exclude
  conditional references and external skill/project instructions; no tokenizer
  or model-performance benchmark was run.
- Live installation/update of these new skills is not yet possible through the
  configured published store while they remain uncommitted. Conflict handling
  and instruction refresh were reviewed, not exercised against a live project.
- No Go changes or shared validation changes: unrelated Go/full-catalog tests
  omitted. Source diff inspected; existing prompt additions remain untouched.

Proposed commit: `feat(skills): extract task workflows from project docs`.

## Developer acceptance

Pending implementation review. Approval of the plan did not authorize a commit.
