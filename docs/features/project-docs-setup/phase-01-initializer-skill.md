# Phase 01: Initializer skill

- **Feature**: [Project documentation setup](FEATURE.md)
- **Depends on**: None.
- **Outcome**: One installed skill can initialize project guidance and docs.

## Current checkpoint

- **Status**: awaiting documentation review.
- **Current step**: Simplified the generated feature guide and aligned SKILL.md
  and assets/AGENTS.md with checklist-only plans for features and codebase work.
- **Approved scope**: User request, 2026-10-15: readable, complete feature
  instructions; logically phased tasks/subtasks; clear progress and resume
  point; codebase maintenance support; reduced context without losing intent.
- **Working tree**: Three source documentation files and these two planning
  records; no helper, test, Go, or external library changes.
- **Last verification**: Offline authoring validator, relative links/anchors,
  balanced Markdown fences, and diff whitespace checks passed. Manual review
  covered feature, refactor, docs-only, partial-work, blocked-review, and commit
  status cases. Feature guide reduced from 1,110 to 892 words (20%); these are
  word counts, not measured tokens. Executable checks were not needed.
- **Pending decision**: User review of this documentation slice.
- **Next action**: User reviews this documentation slice; commit only on a
  new explicit request.

## Scope and required context

Create `skills/create-project-docs/SKILL.md`, template assets, and a Bash copy
helper. Keep the skill self-contained and published in source metadata;
publishing to the configured store requires a separate explicit request.

| Reference | Why it is needed |
| --- | --- |
| [Repository instructions](../../../AGENTS.md) | Authoring setup and review gates |
| [Skill metadata](../../SPEC.md#skill-metadata) | Portable source packaging |
| [Existing agent setup](../../../skills/create-agents-setup/SKILL.md) | Define routing boundary and reuse established guidance |
| [Existing spec workflow](../../../skills/create-spec/SKILL.md) | Avoid implicitly invoking its different `docs/SPEC.md` workflow |

## Proposed generated structure

```text
AGENTS.md
CLAUDE.md
docs/
  conventions/
    CODE.md
    COMMITS.md
    GOLANG.md
    TDD.md
  specs/
    MISSION.md
    FEATURES.md
    ROADMAP.md
    STACK.md
  decisions/
    README.md
  features/
    README.md
```

The directory guides hold reusable record templates and naming rules:
`docs/decisions/0001-<decision>.md` and
`docs/features/0001-<feature>/{plan.md,requirements.md,validation.md}`.
Create actual numbered records only
when the project supplies a real decision or feature. The developer approved
this interpretation of the example numbered paths during plan review.

`MISSION.md` owns the problem, users, intended outcome, and non-goals.
`FEATURES.md` indexes product capabilities and links to their requirements.
Each feature's requirements.md owns its behavior and acceptance contract,
plan.md owns implementation steps and checkpoint, and validation.md owns
verification evidence and user acceptance. `ROADMAP.md` owns delivery order
and links to these records.
`STACK.md` owns selected technologies, versions, verification commands, and
links to official documentation. Pending choices are explicit.

`CODE.md` covers common code guidance, including TypeScript boundaries and
frontend conventions when applicable. `GOLANG.md` covers Go-specific guidance.
`TDD.md` owns behavior-change test and implementation review gates.
`COMMITS.md` owns commit proposals and explicit publication authorization.

Generated `AGENTS.md` routes agents to these owners for their current task.
Generated `CLAUDE.md` contains `@AGENTS.md`, preserving any existing additions.
Do not auto-import convention files or mandate reading all docs at startup.
Use existing project commands when known; do not invent runnable commands for
an empty project. Do not impose mentoring or require unrelated skills.

## Acceptance criteria

- [ ] Empty projects receive all proposed files at the selected destination.
- [ ] Existing regular files are kept without replacement or duplicate content.
- [ ] Paths with spaces work; missing arguments and unsafe destinations fail
  clearly. Symlinks and file/directory collisions are rejected in preflight.
- [ ] Guidance handles Go-only, Go plus TypeScript, and an existing project
  without assuming React or manufacturing product requirements.
- [ ] Relative routes resolve; numbering guidance and document ownership agree.

## Change plan and review checkpoints

### Step 1: Tests for safe initialization

- **Change and reason**: Focused helper tests in isolated temporary directories
  establish the generated file set, preservation on rerun, and preflight refusal.
- **Test cases**: Empty destination including a path with spaces; customized
  existing instruction/convention files; symlink and collision; missing input.
  Confirm failures express missing initializer behavior rather than tool failures.
- **Test commit proposal**: `test(create-project-docs): cover safe documentation initialization`
- **Implementation commit proposal**: `feat(create-project-docs): initialize project guidance and documentation`

After plan approval, write and run the tests. Stop for test review with the diff
and meaningful failure evidence. Implement the helper only after explicit
approval. Skill prose and template assets gain nothing from artificial failing
tests; review their content and links alongside the verified helper.

### Step 2: Skill, templates, and helper

Use the installed authoring skill's validator and its initializer when useful.
The test-first skill directory already exists, so the initializer's refusal to
overwrite existing directories is respected; the small entrypoint was authored
directly. Add the templates
and helper within the approved scope. Inspect the chosen target before writes;
keep existing documents, adapt newly created templates, and report pending
facts and partial creation accurately. Clarify triggering boundaries with
`create-agents-setup`, `create-spec`, and application-layout skills in the new
skill's own instructions without changing those skills.

Run helper tests and catalog validation, check generated links and imports,
and review the final diff. No Go source or external library changes are planned;
Go formatting/build/test checks do not validate these Markdown/Bash artifacts.
If scope expands to Go changes, follow all repository Go gates and verification.
Stop for implementation review; leave changes uncommitted.

## Verification evidence

| Step / criterion | Command or inspection | Result and relevant evidence | Tested state |
| --- | --- | --- | --- |
| Authoring preparation | `skills sync` | Passed: store up to date at `ebfcb9e` | Before plan files |
| Authoring preparation | `skills update create-repo-skill --agent codex` | Passed: 1 skill, 0 changed | Before plan files |
| Baseline | `jj st` | Clean working copy | Before plan files |
| Plan checks | `git diff --check` | Passed | Planning files only; these new files also require direct inspection |
| Working-copy inspection | `jj diff --stat` | Environment failure: sandbox denied snapshot writes to `.git/objects`; no history mutation attempted | Planning files present |
| Step 1 | `python3 -B -m unittest discover -s skills/create-project-docs/tests -v` | Expected failure: 5 test methods, 10 failing cases; `Project documentation initializer not implemented: init-project-docs.sh` | New test file; helper and templates absent |
| Helper criteria | `python3 -B -m unittest discover -s skills/create-project-docs/tests -v` | Passed: all 5 methods, including preservation and refusal subcases | Helper, 12 assets, and approved tests |
| Shell checks | `shellcheck skills/create-project-docs/scripts/init-project-docs.sh`; `bash -n skills/create-project-docs/scripts/init-project-docs.sh` | Both passed | Completed helper |
| Validator environment | `bash .agents/skills/create-repo-skill/scripts/validate-skill.sh skills/create-project-docs` | Environment failure: default Python lacks PyYAML | Completed skill |
| Temporary environment | `UV_CACHE_DIR=/private/tmp/project-docs-uv-cache uv run --with pyyaml --no-project bash .agents/skills/create-repo-skill/scripts/validate-skill.sh skills/create-project-docs` | First attempt failed due sandbox DNS; permitted network retry passed, using temporary cache | Completed skill |
| Catalog format | Same uv validator command with `--offline` | Passed: `Valid skill: create-project-docs` | Final skill instructions and metadata |
| Generated links and imports | Isolated Python inspection of a generated project and source Markdown | Passed: local links, balanced code fences, whitespace, single Claude import, no dummy records | Completed assets and plan |
| Failure reporting | Isolated helper exercises under `/bin/bash` | Passed: root symlink and missing template refuse before writes; injected second-file copy failure reports 1 created and possible partial file | Completed helper |
| Recurring context | Python whitespace-separated word count of assets/AGENTS.md | Approximately 300 words after policy and feature-layout revisions; document routing uses ordinary Markdown links | Revised agent template |
| Trigger boundaries | Manual review of new-project Go/TypeScript, Go-only, existing-owner, instructions-only, and application-scaffold requests | Routes split-spec setup to this skill; preserves existing owners; no assumed React or tool installation; neighboring skills remain unchanged | Final skill instructions and templates |
| Scope-specific checks | Inspection of changed file inventory | No Go files or external libraries changed; Go fix/format/test/vet and integration runs omitted as in approved plan | Skill and planning docs only |
| Requested policy revision | Official Go command documentation and review of GOLANG.md, CODE.md, STACK.md, AGENTS.md, and SKILL.md | Documents periodic module/package fixes, latest-version verification, and explicit approval before dependency changes; focused file work uses its enclosing package | Revised documentation assets; helper and tests unchanged |
| Policy revision checks | Offline authoring validator, `git diff --check`, and direct Markdown link/whitespace/fence inspection | All passed; helper tests not repeated because executable behavior is unchanged | Revised skill, assets, and planning docs |
| Test readability revision | Review of TDD.md and GOLANG.md | Shared guidance prioritizes direct readable tests, minimal stubs, suitable helper reuse, descriptive small tables, and separate success/failure groups; Go guidance adds `t.Helper()` and test naming | Documentation assets only; helper and approved tests unchanged |
| Feature record revision | Review of feature guide, FEATURES.md, ROADMAP.md, decision guide, AGENTS.md, and SKILL.md | Numbered folders use plan.md, requirements.md, and validation.md; templates and routes agree on ownership; no dummy feature generated | Documentation assets only; helper and approved tests unchanged |
| Feature revision checks | Offline authoring validator, `git diff --check`, and Markdown link/anchor/whitespace/fence inspection | All passed; executable helper and test suite unchanged | Updated feature templates and routes |
| Commit review | Helper tests, offline skill validator, ShellCheck, and current diff review | All checks passed; 17 intended files; no unrelated changes | Final skill and planning records |
| Small-slice follow-up | Offline skill validator, `git diff --check`, and scoped diff review | Passed: one focused rule in AGENTS.md and matching execution guidance in the feature plan template; no executable changes | Two revised templates and two updated planning records after `91e4ecca` |
| Sub-feature follow-up | Offline skill validator, `git diff --check`, and scoped text review | Passed: one new section defines sibling numbered sub-feature records, independent outcomes, parent links, dependencies, and review before implementation | Feature guide and planning records only; helper and tests unchanged |
| Context optimization slice 1 | Word count, route inspection, instruction comparison, offline validator, and `git diff --check` | Passed: AGENTS.md reduced from 338 to 261 words (23%); all 10 routes resolve; small slices, review gates, dependency checks/approval, pending facts, preservation, contract conflicts, checkpoint routing, and commit limits retained | AGENTS.md only, plus planning updates; word counts are not measured tokens |
| Minimal-question follow-up | Offline skill validator, `git diff --check`, and scoped instruction review | Passed: setup and later documentation work use the available question tool, reuse facts, ask the minimum for the next slice, keep non-blocking unknowns pending, and fall back to chat if needed | SKILL.md and AGENTS.md plus planning records; no executable changes |
| Adversarial/reply follow-up | Offline validator, diff check, Markdown links/anchors, and independent `review_workflow` subagent inspection | Passed; reviewer reports no actionable findings, confirms approval gates and fix scope are preserved; review loop and minimal replies are documented. One reviewer file read targeted a nonexistent standalone validation template; the actual inline template was subsequently inspected. This was instruction review, not a generated-project code repair exercise. | SKILL.md, AGENTS.md, TDD.md, and planning records; no executable changes |

## Developer acceptance

- **Decision**: Prior work committed; adversarial-review and reply guidance awaits review.
- **Reference and scope**: "yes, proceed", 2026-10-14, accepts the proposed
  design, file ownership, naming patterns, and helper scope. "proceed",
  2026-10-14, approves reviewed tests and authorizes implementation.
  "do a commit for the current changes", 2026-10-14, authorizes one commit of
  the complete reviewed working tree, including tests and implementation together.
- **Remaining work**: Review adversarial-review/reply guidance, then optimize convention
  documents and feature/spec templates in separate reviewable slices.
  Earlier commit authorization does not cover these subsequent changes.
  No push or PR requested; publication remains separate.
