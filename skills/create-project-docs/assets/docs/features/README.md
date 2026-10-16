# Feature and codebase work records

Use this structure for product features, fixes, refactoring, migrations, tooling,
and documentation work. Create records only for real work that needs a plan;
never generate empty examples during setup.

## Naming and ownership

- Create `NNNN-<slug>/` with `requirements.md`, `plan.md`, and `validation.md`.
- Choose the next unused four-digit ID across [ROADMAP.md](../specs/ROADMAP.md)
  and feature folders. Use that ID everywhere; preserve existing IDs and names.
  Roadmap entries may precede their folders.
- `requirements.md` owns the agreed outcome, scope, constraints, approach,
  dependencies, decisions, and acceptance criteria.
- `plan.md` owns ordered tasks, progress, and the exact resume point.
- `validation.md` owns verification evidence, review results, actual commit
  references, and user acceptance.
- ROADMAP.md owns delivery order through feature/sub-feature headings and task
  checklists. Link to detailed records; omit narrative, status fields, and
  duplicated requirements or evidence.
- Keep these three records; do not add competing FEATURE.md or phase files.

## Writing rules

- Add only established facts and agreed work. Omit empty headings, speculative
  phases, filler, and generic questions. Follow the
  [spec writing guide](../specs/README.md).
- Use descriptive headings and short lists, one point per item. Reserve tables
  for compact comparisons. Link shared facts to their owner.
- For requested spec work, follow the
  [spec writing rules](../specs/README.md#working-rules), including guided
  questions and waiting before writing content that depends on answers.

## requirements.md

Add only sections with content:

- Outcome and scope: problem, affected users or maintainers, included work,
  non-goals, constraints, and the user's actual agreement.
- Contract: relevant inputs, preconditions, results, side effects, state,
  failures, and recovery. Cover repeated, conflicting, or partial execution only
  when relevant; do not invent retry or rollback requirements.
- Approach and dependencies: affected boundaries, agreed design choices,
  prerequisites, and links to shared specs or decisions.
- Acceptance criteria: concrete, verifiable scenarios with stable `R-001` IDs.
  For codebase work, state what must improve and what behavior must be preserved.
- Open questions: link relevant groups in [the question log](../QUESTIONS.md);
  follow the [question guide](../conventions/QUESTIONS.md), without local lists.

## plan.md

Use only a document title, ordered phase headings, and task/subtask checklists.
No prose sections, tables, separate status summaries, or checkpoint diary.
Put context in requirements.md and results in validation.md; link from tasks.

- Name each phase for one clear outcome. Order phases and tasks by dependency.
  Keep phases small: one to five planned commits; split larger phases without
  adding filler commits.
- Make each top-level task a reviewable outcome. Add concrete action subtasks
  only where they clarify execution; avoid repeating generic workflow steps.
- Include affected paths or boundaries and acceptance links where useful.
  Record prerequisite work as a task or link a blocking dependency inline.
- Detail the active phase enough to act without replanning. Keep identified
  later phases brief and expand them before starting; never invent future work.
- Include necessary tests, implementation or edits, verification, and review
  as concrete subtasks, following [TDD.md](../conventions/TDD.md). For docs or
  behavior-preserving work, choose checks that prove the outcome rather than
  manufacturing a failing test.
- Put each proposed commit message in a subtask and respect project policy on
  separate test/implementation commits. Distinguish a proposal from an actual
  commit with a reference in validation.md; checked work does not imply a commit.
- Mark completed actions `[x]` with a concise evidence link. Check a parent only
  when all required subtasks and reviews are complete. A phase is complete when
  all its required tasks are checked.
- Keep unfinished actions `[ ]`. Append a short note only when needed:
  `in progress`, `blocked: <reason>`, `awaiting review`, or `deferred: <reason>`.
  Partial work stays unchecked; split off completed subtasks when useful.
- Mark exactly one actionable unchecked task or subtask `next`. If progress
  requires a decision or review, mark that action next and name what is needed.
  This is the resume point; do not repeat it in a separate checkpoint.
- At review boundaries and session end, update checks, task notes, evidence
  links, and the next action. Completion does not imply feature acceptance or
  permission to commit or publish.

## validation.md

- Link evidence to acceptance criteria or plan tasks. Record the scenario,
  command or inspection, result, checked revision or working-tree state, and
  material limits. Include expected failure evidence when TDD applies.
- Distinguish change failures, environment problems, skipped checks, and checks
  not run. Passing commands alone do not prove every criterion.
- Record review decisions, actual commit references when commits exist, and
  the user's acceptance with its reference, accepted scope, and remaining work.
  Never infer acceptance from passing checks or completed tasks.

## Split large work

- Split work that cannot be implemented and reviewed in small slices into
  independently verifiable sub-features. Present the split for user review.
- Give each its own ID and sibling folder with the same three files. Link parent
  and children both ways; keep shared constraints in the parent.
- Show delivery order through ROADMAP.md grouping and task order; keep detailed
  dependencies in requirements.md and link them from blocked tasks.
- Implement one sub-feature at a time. Accept the parent only after child work
  and shared acceptance criteria have evidence and user acceptance.

## Resume

- Read the active phase and next action in plan.md, its linked requirements and
  acceptance criteria, and the latest applicable validation evidence. Compare
  recorded progress and checked revisions with the working tree before continuing.
- Read other sections only for relevant dependencies or unresolved decisions.
  Resuming existing work does not require the record-creation or splitting rules.
- Update plan.md and validation.md at review boundaries and session end; reflect proven progress
  in ROADMAP.md without duplicating the detailed plan.
