# Phase 03: Extract task workflow skills

- **Feature**: [Project documentation setup](FEATURE.md)
- **Outcome**: Three reusable task skills own feature requirements, technical
  planning, and execution checklists; project setup installs and routes to them.

## Current checkpoint

- **Status**: awaiting implementation review
- **Approved scope**: User approved the proposed extraction, CLI setup/update,
  narrow instruction refresh, and validation on 2026-10-16: “approve”.
  Follow-up split approved with “looks good”: core workflow in AGENTS.md,
  test specifics in TDD.md, detailed review procedure in a separate guide.
- **Working tree**: Workflow split is in parent b928b6d3. Current follow-up
  removes AGENTS.md backlinks from generated guides and updates the TDD anchor.
  Working tree was clean at the start of this follow-up.
- **Test gate**: Inventory test failed for missing REVIEW.md. User approved
  implementation with “approve” after reviewing that failure.
- **Next action**: Review the backlink removal and verification results.
  No commit authorized.

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

## Follow-up: Always-loaded workflow and separate review guide

- **Approved scope**: User requested an AGENTS.md trim, then approved keeping
  the core workflow always loaded, test specifics in TDD.md, and detailed
  review in a separate guide (“looks good”). Explicit test approval: “approve”.
- **Changes**: AGENTS.md owns small slices, approved scope, evidence discipline,
  independent review, user review, and acceptance/commit distinctions. TDD.md
  retains test selection, failing-test workflow, verification, and exceptions.
  REVIEW.md retains the detailed reviewer instructions and findings lifecycle.
- **Packaging and routes**: Added REVIEW.md to the initializer copy list and
  expected inventory. Updated feature, commit, Go, prototype, and setup guidance
  to point to each topic's owner. jj and prototype exceptions remain explicit.
- **Expected failure**: Before implementation,
  `python3 -B -m unittest discover -s skills/create-project-docs/tests -p test_initializer.py -k creates_documentation -v`
  failed at the generated file-set assertion: docs/conventions/REVIEW.md missing.
- **Verification**: `python3 -B -m unittest discover -s skills/create-project-docs/tests -v`
  passes all seven tests, including file creation and byte-preserving reruns.
  Source validator passes with cached offline uv/PyYAML. A fresh generated
  project resolves 34 relative links/anchors across seven affected documents.
  Diff whitespace checks pass. No Go changes; Go checks omitted.
- **Context**: AGENTS.md is approximately 345 words with the always-loaded core
  workflow restored. Detailed review remains conditional; no token or model
  performance claims. Earlier 213-word trim was superseded by this scope correction.
- **Review limits**: Inspected the source diff, generated routes, and preserved
  procedure content. No independent agent execution evaluation performed.
- **Status**: awaiting implementation review. Nothing committed.
- **Test commit proposal**: `test(create-project-docs): require generated review guide`.
- **Implementation commit proposal**: `docs(create-project-docs): separate workflow testing and review guidance`.
- **Next action**: User reviews the completed split and verification evidence.

## Follow-up: Remove AGENTS.md backlinks

- **Approved scope**: User requested removal because AGENTS.md is always loaded.
- **Changes**: Removed backlinks from TDD, review, commit, prototype, feature,
  and decision guides. Renamed the TDD section to Test workflow and updated its
  incoming anchor. Kept the required CLAUDE.md import and outward task routes.
- **Verification**: create-project-docs source validator passes with offline
  uv/PyYAML; 29 relative links/anchors resolve across seven affected assets;
  diff whitespace checks pass. Asset search finds only the CLAUDE.md import
  referencing AGENTS.md. No helper or file inventory changes; initializer and
  Go tests not rerun for this documentation-only edit.
- **Status**: awaiting implementation review; no commit authorized.
- **Commit proposal**: `docs(create-project-docs): remove redundant agent instruction backlinks`.
- **Next action**: Review the guide edits.
