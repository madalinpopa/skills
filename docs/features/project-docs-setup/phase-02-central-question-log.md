# Phase 02: Central question log

- **Feature**: [Project documentation setup](FEATURE.md)
- **Depends on**: Existing initializer and generated documentation guides.
- **Outcome**: Actual unresolved questions have one checklist owner; answers
  update their relevant documents before questions are closed.

## Current checkpoint

- **Status**: awaiting test review.
- **Current step**: Initializer tests await review. Separately applied the
  explicitly requested documentation-only reply changes; those await review.
- **Approved scope**: User approved the proposed plan on 2026-10-15:
  "arata bine, continua cu implementarea". This does not bypass test review.
  Follow-up request replaces the mandatory awareness footer with Next Step
  and requires recording questions or blockers when discovered during work.
  A further explicit request requires the available ask-question tool whenever
  it supports a question, including questions outside spec work.
- **Working tree**: test_initializer.py, SKILL.md, assets/AGENTS.md, and these
  two planning records. Helper and new question assets remain unchanged.
- **Last verification**: The focused Python suite fails because the initializer
  does not create or preflight the two new question documents.
- **Pending decision**: Explicit approval of tests and implementation scope.
- **Next action**: User reviews the reply guidance and initializer tests;
  after test approval, implement the remaining question-log workflow.

## Scope and required context

- Add an initially empty docs/QUESTIONS.md and a concise persistent procedure
  in docs/conventions/QUESTIONS.md. Keep instructions out of the log.
- Update the initializer file list, SKILL.md, generated AGENTS.md, and affected
  spec/feature/decision guides and TDD reporting guidance. Replace local
  question lists with log links.
- Use stable question IDs, open-first task lists, required input, relevant
  target, and blocking impact when applicable. Keep resolved entries concise.
- Document complete or partial answers in their owning sections. Close only
  after successful updates; distinguish withdrawn or replaced questions.
- Preserve existing files and approvals. Consolidating existing questions
  requires authorization for affected edits; do not create a competing owner.
- Load open questions and relevant documents; consult resolved history only
  when needed. Report what input is needed next without repeating the log.
- Record questions and blockers when discovered, rather than waiting for the
  final reply. Explain required input and link the affected work. A blocker
  the agent can resolve within scope does not require an invented user question.
- Remove the mandatory User must be aware footer everywhere in this skill's
  output instructions. Report material failures or limits where relevant.
- End each work-slice reply with one concrete Next Step based on current state:
  the blocking question, required review, next task, or useful feature proposal.
  Reference an existing question ID when applicable. A suggestion does not
  authorize starting new scope; avoid inventing work when the request is done.
- No CLI behavior, Go, dependencies, commits, or publication changes.

Read [repository instructions](../../../AGENTS.md),
[TDD](../../conventions/TDD.md), and the source skill and affected assets.

## Acceptance criteria

- [ ] Initialization creates the empty log and its nonempty procedure guide.
- [ ] Reruns preserve existing questions and all generated document bytes.
- [ ] Symlinks and directory collisions at either new path fail before writes.
- [ ] Guidance has one question owner, stable IDs, and document-first closure.
- [ ] Partial answers, duplicates, contradictions, withdrawals, and resumption
  retain unresolved information without inventing decisions or authorization.
- [ ] Relative routes resolve and recurring guidance stays concise.
- [ ] Questions and blockers are recorded during work; replies propose one
  actionable Next Step without the mandatory awareness footer or duplicate logs.

## Change plan and review checkpoints

### Step 1: Initializer tests

- Extend existing public-entrypoint cases for the two new files, empty log,
  existing-log preservation, reruns, and unsafe destinations.
- Run the focused suite and stop for test review with meaningful failures.
- Test commit proposal:
  `test(create-project-docs): cover question log initialization and preservation`

### Step 2: Implementation after test approval

- Add assets and helper paths, then align the skill and generated routes.
- Write question lifecycle rules once in the procedure guide. Review realistic
  answer and resume scenarios; avoid wording-matching tests for instructions.
- Run focused helper tests, shell checks, skill validation, Markdown link and
  diff checks. Compare words/bytes without claiming measured token savings.
- Stop for implementation review; update this checkpoint and feature roadmap.
- Implementation commit proposal:
  `feat(create-project-docs): centralize documentation questions`

## Verification evidence

- Baseline: clean working copy at rmrmpomp / 19854b6f, parent f24c3f32.
- Command: `rtk proxy python3 -B -m unittest discover -s skills/create-project-docs/tests -v`.
- Result: 7 test methods, 9 expected assertion failures, no errors. New files
  are absent; symlinks and directory collisions at their paths are accepted by
  the old helper. Existing-log preservation checks are ready but cannot finish
  until the missing guide is created.
- Tests exercise observable filesystem behavior; helper and assets unchanged.
- Follow-up reporting changes do not alter helper behavior or test scope;
  review them with the skill instructions rather than wording-matching tests.
- Reply guidance applied in SKILL.md and assets/AGENTS.md on explicit user
  correction: remove the awareness footer, record questions as they arise,
  and end with one actionable Next Step. Scoped diff inspection and whitespace
  checks passed. Default validator lacked PyYAML; the existing temporary uv
  environment validated the skill offline successfully. No packages installed.
- Extended the same reply guidance to use an available ask-question tool for
  all questions it supports, with chat fallback and reuse of known answers.
- `git diff --check` passed. Initial `jj diff` hit sandbox snapshot restrictions;
  the permitted retry succeeded. No history changes were requested or made.

## Developer acceptance

- **Decision**: Pending test review.
- **Reference and scope**: Plan approved on 2026-10-15; implementation still
  requires explicit approval of the test review under AGENTS.md.
- **Remaining work**: Implementation, verification, and implementation review.
